package inference

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// LocalHTTP only contacts numeric loopback endpoints. No proxy, redirects or credentials.
type LocalHTTP struct {
	streaming             bool
	name, model, endpoint string
	features              []string
	client                *http.Client
	mu                    sync.Mutex
	active                map[string]context.CancelFunc
}

var _ InferenceBackend = (*LocalHTTP)(nil)

func NewLocalHTTP(id, model, endpoint string, features []string) (*LocalHTTP, error) {
	u, e := url.Parse(endpoint)
	if e != nil {
		return nil, e
	}
	ip := net.ParseIP(u.Hostname())
	if id == "" || model == "" || u.Scheme != "http" || ip == nil || !ip.IsLoopback() || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" && u.Path != "/" {
		return nil, fmt.Errorf("explicit numeric loopback HTTP origin required")
	}
	tr := &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: 5 * time.Second}).DialContext, ResponseHeaderTimeout: 30 * time.Second, MaxResponseHeaderBytes: 16384}
	return &LocalHTTP{name: id, model: model, endpoint: strings.TrimRight(endpoint, "/"), features: append([]string(nil), features...), client: &http.Client{Transport: tr, Timeout: 15 * time.Minute, CheckRedirect: func(*http.Request, []*http.Request) error { return fmt.Errorf("redirect denied") }}, active: map[string]context.CancelFunc{}}, nil
}
func (b *LocalHTTP) ID() string             { return b.name }
func (b *LocalHTTP) Capabilities() []string { return append([]string(nil), b.features...) }
func (b *LocalHTTP) EstimateCost(r Request) Cost {
	return Cost{len(r.Context), r.MaxOutputTokens, true}
}
func (b *LocalHTTP) Health(ctx context.Context) error {
	ctx, c := context.WithTimeout(ctx, 3*time.Second)
	defer c()
	req, e := http.NewRequestWithContext(ctx, "GET", b.endpoint+"/v1/models", nil)
	if e != nil {
		return e
	}
	res, e := b.client.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("health HTTP %d", res.StatusCode)
	}
	var body struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	raw, e := io.ReadAll(io.LimitReader(res.Body, 65537))
	if e != nil {
		return e
	}
	if len(raw) > 65536 || !utf8.Valid(raw) {
		return fmt.Errorf("invalid/oversize health response")
	}
	if e = json.Unmarshal(raw, &body); e != nil {
		return e
	}
	for _, m := range body.Data {
		if m.ID == b.model {
			return nil
		}
	}
	return fmt.Errorf("configured model unavailable")
}
func (b *LocalHTTP) Cancel(id string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	c, ok := b.active[id]
	if ok {
		c()
	}
	return ok
}
func (b *LocalHTTP) Invoke(parent context.Context, r Request) (out Result, err error) {
	out = Result{BackendID: b.name, ModelID: b.model, Status: "failed"}
	start := time.Now()
	defer func() {
		out.LatencyMS = time.Since(start).Milliseconds()
		if err != nil {
			out.Status = "failed"
			if out.ErrorClass == "" {
				out.ErrorClass = Classify(err)
			}
			out.Termination = out.ErrorClass
		}
	}()
	if err = r.Validate(); err != nil {
		return
	}
	if r.OutputConstraint != nil {
		err = fmt.Errorf("output constraint capability unavailable")
		return
	}

	ctx, c := context.WithDeadline(parent, r.Deadline)
	defer c()
	b.mu.Lock()
	if _, ok := b.active[r.ID]; ok {
		b.mu.Unlock()
		err = fmt.Errorf("duplicate active request")
		return
	}
	b.active[r.ID] = c
	b.mu.Unlock()
	defer func() { b.mu.Lock(); delete(b.active, r.ID); b.mu.Unlock() }()
	payload := struct {
		Model           string              `json:"model"`
		Messages        []map[string]string `json:"messages"`
		MaxTokens       int                 `json:"max_tokens"`
		Stream          bool                `json:"stream"`
		StreamOptions   map[string]bool     `json:"stream_options,omitempty"`
		TimingsPerToken bool                `json:"timings_per_token,omitempty"`
	}{b.model, []map[string]string{{"role": "user", "content": r.Context}}, r.MaxOutputTokens, b.streaming, nil, b.streaming}
	if b.streaming {
		payload.StreamOptions = map[string]bool{"include_usage": true}
	}
	data, _ := json.Marshal(payload)
	req, e := http.NewRequestWithContext(ctx, "POST", b.endpoint+"/v1/chat/completions", bytes.NewReader(data))
	if e != nil {
		err = e
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res, e := b.client.Do(req)
	if e != nil {
		err = e
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		switch res.StatusCode {
		case 401, 403:
			out.ErrorClass = "auth_failed"
		case 429:
			out.ErrorClass = "rate_limited"
		case 404, 503:
			out.ErrorClass = "model_unavailable"
		}
		err = fmt.Errorf("inference HTTP %d", res.StatusCode)
		return
	}
	if b.streaming {
		if !strings.HasPrefix(strings.ToLower(res.Header.Get("Content-Type")), "text/event-stream") {
			err = fmt.Errorf("SSE content type required")
			return
		}
		out, err = b.readStream(res.Body, r, start)
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return
	}
	limit := int64(r.MaxOutputTokens*16 + 65536)
	data, err = io.ReadAll(io.LimitReader(res.Body, limit+1))
	if err != nil {
		return
	}
	if int64(len(data)) > limit {
		err = fmt.Errorf("response byte limit exceeded")
		return
	}
	var wire struct {
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Finish string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			Prompt     *int `json:"prompt_tokens"`
			Completion *int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if !utf8.Valid(data) {
		err = fmt.Errorf("invalid UTF-8 response")
		return
	}
	if err = json.Unmarshal(data, &wire); err != nil {
		return
	}
	if wire.Model != b.model || len(wire.Choices) != 1 {
		err = fmt.Errorf("unexpected model or choices")
		return
	}
	if wire.Usage.Completion != nil && (*wire.Usage.Completion < 0 || *wire.Usage.Completion > r.MaxOutputTokens) || wire.Usage.Prompt != nil && *wire.Usage.Prompt < 0 {
		err = fmt.Errorf("invalid token usage")
		return
	}
	if wire.Choices[0].Finish != "stop" {
		err = fmt.Errorf("incomplete generation: %s", wire.Choices[0].Finish)
		return
	}
	if len(wire.Choices[0].Message.Content) > r.MaxOutputTokens*16 {
		err = fmt.Errorf("output byte limit")
		return
	}
	out.Text = wire.Choices[0].Message.Content
	out.Usage = Usage{wire.Usage.Prompt, wire.Usage.Completion}
	out.Status = "completed"
	out.Termination = "stop"
	return
}
