package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spoonman136668-ai/Wingless/reasoner"
)

const (
	gatewayModel = "nemotron-3-ultra-planner"
	plannerSchema = "ckb-plane.research-planner-request.v1"
	contractText = "ckb-plane.nemotron-planner-gateway.v1|public-repository-only|no-authority|no-retry"
	maxBodyBytes = 8 << 20
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []chatMessage `json:"messages"`
	Temperature    float64       `json:"temperature"`
	MaxTokens      int           `json:"max_tokens"`
	Stream         bool          `json:"stream"`
	ResponseFormat struct {
		Type string `json:"type"`
	} `json:"response_format"`
}

type artifactRef struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Commit string `json:"commit,omitempty"`
}

type plannerEvidence struct {
	ID      string `json:"id"`
	Path    string `json:"path"`
	Commit  string `json:"commit,omitempty"`
	SHA256  string `json:"sha256"`
	Content string `json:"content"`
}

type plannerRequest struct {
	Schema    string `json:"schema"`
	Envelope  struct {
		ID                  string        `json:"id"`
		Project             string        `json:"project"`
		Lane                string        `json:"lane"`
		AuthorityGeneration int64         `json:"authority_generation"`
		BaselineSHA         string        `json:"baseline_sha"`
		NorthStar           artifactRef   `json:"north_star"`
		Evidence            []artifactRef `json:"evidence"`
	} `json:"envelope"`
	NorthStar string            `json:"north_star"`
	Evidence  []plannerEvidence `json:"evidence"`
}

type reasonerClient interface {
	Invoke(context.Context, reasoner.Request) (reasoner.Result, error)
}

type gateway struct {
	client       reasonerClient
	publicClient *http.Client
	repository   string
	timeout      time.Duration
}

func main() {
	listen := flag.String("listen", "127.0.0.1:18791", "numeric IPv4 loopback listen address")
	timeoutSeconds := flag.Int("timeout-seconds", 600, "single planning request timeout")
	flag.Parse()

	if *listen == "" || !strings.HasPrefix(*listen, "127.0.0.1:") {
		die(errors.New("listen must be numeric IPv4 loopback"))
	}
	host, port, err := net.SplitHostPort(*listen)
	if err != nil || host != "127.0.0.1" || port == "" {
		die(errors.New("listen must be numeric IPv4 loopback with port"))
	}
	if *timeoutSeconds < 1 || *timeoutSeconds > 900 {
		die(errors.New("timeout-seconds out of range"))
	}
	key := os.Getenv("WINGLESS_REASONER_OPENROUTER_API_KEY")
	if strings.TrimSpace(key) == "" {
		die(errors.New("WINGLESS_REASONER_OPENROUTER_API_KEY is not set"))
	}
	cfg := reasoner.DefaultConfig()
	client, err := reasoner.NewClient(key, cfg)
	if err != nil {
		die(err)
	}
	publicTransport := &http.Transport{
		Proxy:                 nil,
		ResponseHeaderTimeout: 20 * time.Second,
		MaxResponseHeaderBytes: 16 << 10,
	}
	defer publicTransport.CloseIdleConnections()
	g := &gateway{
		client: client,
		publicClient: &http.Client{
			Transport: publicTransport,
			Timeout:   30 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errors.New("public evidence redirect denied")
			},
		},
		repository: "spoonman136668-ai/Wingless",
		timeout:    time.Duration(*timeoutSeconds) * time.Second,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/models", g.models)
	mux.HandleFunc("/v1/chat/completions", g.chat)
	server := &http.Server{
		Addr:              *listen,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      time.Duration(*timeoutSeconds+30) * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	fmt.Printf("CKB_NEMOTRON_PLANNER_GATEWAY=READY listen=%s model=%s\n", *listen, gatewayModel)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		die(err)
	}
}

func (g *gateway) models(w http.ResponseWriter, r *http.Request) {
	if !loopbackRequest(r) || r.Method != http.MethodGet {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"object": "list",
		"data": []map[string]any{{"id": gatewayModel, "object": "model"}},
	})
}

func (g *gateway) chat(w http.ResponseWriter, r *http.Request) {
	if !loopbackRequest(r) || r.Method != http.MethodPost {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil || len(raw) == 0 || len(raw) > maxBodyBytes {
		http.Error(w, "request invalid", http.StatusBadRequest)
		return
	}
	var req chatRequest
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil || req.Model != gatewayModel || req.Stream ||
		req.Temperature != 0 || req.MaxTokens < 256 || req.MaxTokens > reasoner.MaxCompletionTokens ||
		req.ResponseFormat.Type != "json_object" || len(req.Messages) != 2 ||
		req.Messages[0].Role != "system" || req.Messages[1].Role != "user" {
		http.Error(w, "request contract invalid", http.StatusBadRequest)
		return
	}
	var planner plannerRequest
	if err := json.Unmarshal([]byte(req.Messages[1].Content), &planner); err != nil {
		http.Error(w, "planner request invalid", http.StatusBadRequest)
		return
	}
	if err := validatePlannerRequest(planner); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), g.timeout)
	defer cancel()
	if err := g.verifyPublicMaterial(ctx, planner); err != nil {
		http.Error(w, "public evidence verification failed", http.StatusBadRequest)
		return
	}

	frontier := sha256hex([]byte(req.Messages[1].Content))
	contract := sha256hex([]byte(contractText))
	reasonerReq := reasoner.Request{
		Schema:                      reasoner.RequestSchema,
		RequestID:                   planner.Envelope.ID,
		Project:                     planner.Envelope.Project,
		ExperimentID:                planner.Envelope.ID,
		Role:                        reasoner.RolePlan,
		FrontierSHA256:              frontier,
		QualificationContractSHA256: contract,
		Context: "Return exactly one JSON object matching ckb-plane.research-experiment-proposal.v1. " +
			"Use only the following verified CKB-plane planner request. Do not add markdown or commentary.\n" +
			req.Messages[1].Content,
		MaxOutputTokens: req.MaxTokens,
		TimeoutSeconds:  int(g.timeout / time.Second),
		DataClass:       "public-repository",
	}
	result, err := g.client.Invoke(ctx, reasonerReq)
	if err != nil {
		http.Error(w, "reasoner unavailable", http.StatusBadGateway)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id":      "ckb-nemotron-" + planner.Envelope.ID,
		"object":  "chat.completion",
		"model":   gatewayModel,
		"choices": []map[string]any{{"index": 0, "message": map[string]any{"role": "assistant", "content": result.Text}, "finish_reason": "stop"}},
	})
}

func validatePlannerRequest(p plannerRequest) error {
	if p.Schema != plannerSchema || p.Envelope.Project != "Wingless" || p.Envelope.ID == "" ||
		p.Envelope.Lane == "" || p.Envelope.AuthorityGeneration < 1 || !isSHA1(p.Envelope.BaselineSHA) {
		return errors.New("planner identity invalid")
	}
	if !isSHA256(p.Envelope.NorthStar.SHA256) || !isSHA1(p.Envelope.NorthStar.Commit) ||
		p.Envelope.NorthStar.Path == "" || len(p.NorthStar) == 0 ||
		sha256hex([]byte(p.NorthStar)) != p.Envelope.NorthStar.SHA256 {
		return errors.New("north star binding invalid")
	}
	if len(p.Evidence) != len(p.Envelope.Evidence) {
		return errors.New("evidence set mismatch")
	}
	refs := map[string]artifactRef{}
	for _, ref := range p.Envelope.Evidence {
		if ref.ID == "" || ref.Path == "" || !isSHA256(ref.SHA256) || !isSHA1(ref.Commit) {
			return errors.New("evidence reference invalid")
		}
		if _, ok := refs[ref.ID]; ok {
			return errors.New("evidence reference duplicate")
		}
		refs[ref.ID] = ref
	}
	for _, e := range p.Evidence {
		ref, ok := refs[e.ID]
		if !ok || e.Path != ref.Path || e.Commit != ref.Commit || e.SHA256 != ref.SHA256 ||
			sha256hex([]byte(e.Content)) != e.SHA256 {
			return errors.New("evidence binding invalid")
		}
		delete(refs, e.ID)
	}
	if len(refs) != 0 {
		return errors.New("evidence set mismatch")
	}
	return nil
}

func (g *gateway) verifyPublicMaterial(ctx context.Context, p plannerRequest) error {
	if err := g.verifyPublic(ctx, p.Envelope.NorthStar.Commit, p.Envelope.NorthStar.Path, p.Envelope.NorthStar.SHA256); err != nil {
		return err
	}
	for _, e := range p.Evidence {
		if err := g.verifyPublic(ctx, e.Commit, e.Path, e.SHA256); err != nil {
			return err
		}
	}
	return nil
}

func (g *gateway) verifyPublic(ctx context.Context, commit, path, want string) error {
	parts := strings.Split(path, "/")
	for i := range parts {
		if parts[i] == "" || parts[i] == "." || parts[i] == ".." {
			return errors.New("public evidence path invalid")
		}
		parts[i] = url.PathEscape(parts[i])
	}
	u := "https://raw.githubusercontent.com/" + g.repository + "/" + commit + "/" + strings.Join(parts, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	res, err := g.publicClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("public evidence HTTP %d", res.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(res.Body, 4<<20+1))
	if err != nil || len(raw) > 4<<20 || sha256hex(raw) != want {
		return errors.New("public evidence hash mismatch")
	}
	return nil
}

func loopbackRequest(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	raw, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "encode failure", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(raw)
}

func sha256hex(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func isSHA256(v string) bool {
	return len(v) == 64 && strings.Trim(v, "0123456789abcdef") == ""
}

func isSHA1(v string) bool {
	return len(v) == 40 && strings.Trim(v, "0123456789abcdef") == ""
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "ckb-planner-gateway:", err)
	os.Exit(1)
}
