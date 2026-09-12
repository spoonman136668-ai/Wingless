package ckbplane

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

type listenerEngineFunc func(context.Context, TransportRequest) (CandidateEnvelope, error)

func (f listenerEngineFunc) Run(ctx context.Context, req TransportRequest) (CandidateEnvelope, error) {
	return f(ctx, req)
}

func listenerFixture(t *testing.T) (string, TransportRequest, CandidateEnvelope) {
	t.Helper()
	session := strings.Repeat("a", 64)
	workspace := t.TempDir()
	req := TransportRequest{
		Schema:               TransportRequestSchema,
		SourceCommit:         FrozenPlaneSourceCommit,
		RequestID:            "WO_TEST-wingless-1",
		WorkOrderID:          "WO_TEST",
		Attempt:              1,
		Workspace:            workspace,
		BaselineSHA:          strings.Repeat("b", 40),
		WorkOrderDigest:      strings.Repeat("c", 64),
		Title:                "fixture work",
		Instructions:         []string{"produce a bounded candidate"},
		ForbiddenActions:     []string{"accept", "deploy"},
		Preserve:             []string{"existing behavior"},
		ForbiddenPaths:       []string{"accepted/"},
		Deadline:             time.Now().Add(time.Minute).UTC(),
		Policy: TransportPolicy{
			ReadPaths:       []string{"src/a.go"},
			AllowedBackends: []string{"fast"},
			AllowFallback:   true,
			MaxContextBytes: 64 << 10,
			MaxOutputTokens: 512,
		},
		WorkOrderRetryBudget: 0,
	}
	edits := &EditEnvelope{
		SourceCommit:    req.SourceCommit,
		RequestID:       req.RequestID,
		WorkOrderID:     req.WorkOrderID,
		Attempt:         req.Attempt,
		Workspace:       req.Workspace,
		BaselineSHA:     req.BaselineSHA,
		WorkOrderDigest: req.WorkOrderDigest,
		Operations: []EditOperation{{Type: "write_text", Path: "src/new.go", Content: "package src\n"}},
	}
	candidate := CandidateEnvelope{
		Schema:          CandidateEnvelopeSchema,
		SourceCommit:    req.SourceCommit,
		RequestID:       req.RequestID,
		WorkOrderID:     req.WorkOrderID,
		Attempt:         req.Attempt,
		Workspace:       req.Workspace,
		BaselineSHA:     req.BaselineSHA,
		WorkOrderDigest: req.WorkOrderDigest,
		Status:          "result_ready",
		Acceptance:      "external_required",
		Text:            "fixture candidate",
		BackendID:       "fast",
		ModelID:         "fixture",
		Edits:           edits,
	}
	return session, req, candidate
}

func requestRecorder(t *testing.T, listener *Listener, req TransportRequest, session string) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	httpReq := httptest.NewRequest(http.MethodPost, "http://127.0.0.1"+TransportPath, bytes.NewReader(raw))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-CKB-Wingless-Session", session)
	w := httptest.NewRecorder()
	listener.ServeHTTP(w, httpReq)
	return w
}

func TestWinglessListenerAcceptsExactPlaneRequest(t *testing.T) {
	session, req, candidate := listenerFixture(t)
	var calls int32
	listener, err := NewListener(session, listenerEngineFunc(func(ctx context.Context, got TransportRequest) (CandidateEnvelope, error) {
		atomic.AddInt32(&calls, 1)
		if got.RequestID != req.RequestID || got.WorkOrderDigest != req.WorkOrderDigest || got.WorkOrderRetryBudget != 0 {
			t.Fatalf("request identity changed: %+v", got)
		}
		return candidate, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	w := requestRecorder(t, listener, req, session)
	if w.Code != http.StatusOK || w.Header().Get("X-CKB-Wingless-Session") != session || atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("response=%d session=%q calls=%d body=%s", w.Code, w.Header().Get("X-CKB-Wingless-Session"), calls, w.Body.String())
	}
	var got CandidateEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Acceptance != "external_required" || got.RequestID != req.RequestID || got.Edits == nil {
		t.Fatal(got)
	}
}

func TestWinglessListenerRejectsSessionMethodPathAndContentType(t *testing.T) {
	session, req, candidate := listenerFixture(t)
	listener, _ := NewListener(session, listenerEngineFunc(func(context.Context, TransportRequest) (CandidateEnvelope, error) { return candidate, nil }))
	raw, _ := json.Marshal(req)
	cases := []struct {
		method, path, session, contentType string
		want                              int
	}{
		{http.MethodGet, TransportPath, session, "application/json", http.StatusMethodNotAllowed},
		{http.MethodPost, "/wrong", session, "application/json", http.StatusNotFound},
		{http.MethodPost, TransportPath, strings.Repeat("b", 64), "application/json", http.StatusUnauthorized},
		{http.MethodPost, TransportPath, session, "text/plain", http.StatusUnsupportedMediaType},
	}
	for _, tc := range cases {
		r := httptest.NewRequest(tc.method, "http://127.0.0.1"+tc.path, bytes.NewReader(raw))
		r.Header.Set("Content-Type", tc.contentType)
		r.Header.Set("X-CKB-Wingless-Session", tc.session)
		w := httptest.NewRecorder()
		listener.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("%s %s: got %d want %d", tc.method, tc.path, w.Code, tc.want)
		}
	}
}

func TestWinglessListenerRejectsUnknownDuplicateAndTrailingJSON(t *testing.T) {
	session, req, candidate := listenerFixture(t)
	listener, _ := NewListener(session, listenerEngineFunc(func(context.Context, TransportRequest) (CandidateEnvelope, error) { return candidate, nil }))
	base, _ := json.Marshal(req)
	variants := [][]byte{
		[]byte(strings.Replace(string(base), `"schema":`, `"unknown":true,"schema":`, 1)),
		[]byte(strings.Replace(string(base), `"attempt":1`, `"attempt":1,"attempt":1`, 1)),
		append(append([]byte(nil), base...), []byte(` {}`)...),
	}
	for _, raw := range variants {
		r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1"+TransportPath, bytes.NewReader(raw))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("X-CKB-Wingless-Session", session)
		w := httptest.NewRecorder()
		listener.ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("invalid request accepted: %d %s", w.Code, w.Body.String())
		}
	}
}

type trackingBody struct {
	data   *bytes.Reader
	closed atomic.Bool
	eof    atomic.Bool
}

func newTrackingBody(data []byte) *trackingBody { return &trackingBody{data: bytes.NewReader(data)} }
func (b *trackingBody) Read(p []byte) (int, error) {
	n, err := b.data.Read(p)
	if err == io.EOF {
		b.eof.Store(true)
	}
	return n, err
}
func (b *trackingBody) Close() error { b.closed.Store(true); return nil }

func TestWinglessListenerConsumesAndClosesBodyBeforeInference(t *testing.T) {
	session, req, candidate := listenerFixture(t)
	raw, _ := json.Marshal(req)
	body := newTrackingBody(raw)
	listener, _ := NewListener(session, listenerEngineFunc(func(context.Context, TransportRequest) (CandidateEnvelope, error) {
		if !body.closed.Load() || !body.eof.Load() {
			t.Fatal("inference started before request body was fully consumed and closed")
		}
		return candidate, nil
	}))
	r := httptest.NewRequest(http.MethodPost, "http://127.0.0.1"+TransportPath, nil)
	r.Body = body
	r.ContentLength = int64(len(raw))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-CKB-Wingless-Session", session)
	w := httptest.NewRecorder()
	listener.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("got %d: %s", w.Code, w.Body.String())
	}
}

func TestWinglessListenerCancellationParentsInferenceToRequest(t *testing.T) {
	session, req, _ := listenerFixture(t)
	started := make(chan struct{})
	stopped := make(chan struct{})
	listener, _ := NewListener(session, listenerEngineFunc(func(ctx context.Context, req TransportRequest) (CandidateEnvelope, error) {
		close(started)
		<-ctx.Done()
		close(stopped)
		return CandidateEnvelope{}, ctx.Err()
	}))
	server := httptest.NewServer(listener.Handler())
	defer server.Close()
	raw, _ := json.Marshal(req)
	ctx, cancel := context.WithCancel(context.Background())
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+TransportPath, bytes.NewReader(raw))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-CKB-Wingless-Session", session)
	done := make(chan error, 1)
	go func() {
		resp, err := http.DefaultClient.Do(httpReq)
		if resp != nil {
			_ = resp.Body.Close()
		}
		done <- err
	}()
	<-started
	cancel()
	select {
	case <-stopped:
	case <-time.After(2 * time.Second):
		t.Fatal("inference context did not cancel")
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("client request did not return")
	}
}

func TestWinglessListenerRejectsAcceptanceClaimAndIdentityChange(t *testing.T) {
	session, req, candidate := listenerFixture(t)
	for _, mutate := range []func(*CandidateEnvelope){
		func(c *CandidateEnvelope) { c.Acceptance = "accepted" },
		func(c *CandidateEnvelope) { c.Attempt++ },
		func(c *CandidateEnvelope) { c.WorkOrderDigest = strings.Repeat("d", 64) },
		func(c *CandidateEnvelope) { c.Status = "unknown" },
	} {
		bad := candidate
		bad.EditsCopy(candidate.Edits)
		mutate(&bad)
		listener, _ := NewListener(session, listenerEngineFunc(func(context.Context, TransportRequest) (CandidateEnvelope, error) { return bad, nil }))
		w := requestRecorder(t, listener, req, session)
		if w.Code != http.StatusBadGateway {
			t.Fatalf("invalid candidate accepted: %d %s", w.Code, w.Body.String())
		}
	}
}

func (c *CandidateEnvelope) EditsCopy(src *EditEnvelope) {
	if src == nil {
		c.Edits = nil
		return
	}
	copyEdit := *src
	copyEdit.Operations = append([]EditOperation(nil), src.Operations...)
	c.Edits = &copyEdit
}

func TestWinglessListenerCallsIntelligenceExactlyOnce(t *testing.T) {
	session, req, candidate := listenerFixture(t)
	var calls int32
	listener, _ := NewListener(session, listenerEngineFunc(func(context.Context, TransportRequest) (CandidateEnvelope, error) {
		atomic.AddInt32(&calls, 1)
		return CandidateEnvelope{}, errors.New("fixture failure")
	}))
	w := requestRecorder(t, listener, req, session)
	if w.Code != http.StatusServiceUnavailable || atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("status=%d calls=%d", w.Code, calls)
	}
	_ = candidate
}

func TestWinglessListenerNoProcessListenQueueOrAcceptanceSurface(t *testing.T) {
	typ := reflect.TypeOf(&Listener{})
	for _, name := range []string{"Start", "RunServer", "Listen", "Register", "Accept", "Promote", "Retry"} {
		if _, ok := typ.MethodByName(name); ok {
			t.Fatal("unexpected authority method", name)
		}
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller unavailable")
	}
	raw, err := os.ReadFile(filepath.Join(filepath.Dir(file), "listener.go"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for _, forbidden := range []string{"ListenAndServe", "net.Listen", "os/exec", "exec.Command", "accepted_ref", "queue transition"} {
		if strings.Contains(text, forbidden) {
			t.Fatal("unexpected listener authority surface", forbidden)
		}
	}
}

type fixtureMetrics struct{}
func (fixtureMetrics) Snapshot() (resources.Metrics, error) { return resources.Metrics{}, nil }

type replayEngine struct{ runner broker.Runner }
func (e replayEngine) Run(ctx context.Context, req TransportRequest) (CandidateEnvelope, error) {
	prompt, _ := json.Marshal(struct {
		Title string `json:"title"`
		Instructions []string `json:"instructions"`
	}{req.Title, req.Instructions})
	infReq := inference.Request{
		ID: req.RequestID,
		ParentWorkID: req.WorkOrderID,
		Role: "code",
		Context: string(prompt),
		MaxContextBytes: req.Policy.MaxContextBytes,
		MaxOutputTokens: req.Policy.MaxOutputTokens,
		Deadline: req.Deadline,
		Workspace: req.Workspace,
		Capabilities: []string{"code"},
		Resources: resources.Policy{
			MinRAM: req.Policy.Resources.MinRAMBytes,
			MinVRAM: req.Policy.Resources.MinVRAMBytes,
			MinDisk: req.Policy.Resources.MinDiskBytes,
			MaxCPU: req.Policy.Resources.MaxCPUPercent,
		},
	}
	out := e.runner.Run(ctx, infReq, broker.Policy{
		Explicit: req.Policy.ExplicitBackend,
		AllowedBackends: append([]string(nil), req.Policy.AllowedBackends...),
		AllowDeep: req.Policy.AllowDeep,
		AllowFallback: req.Policy.AllowFallback,
		Repairs: 0,
		MaxRepairs: 0,
	})
	base, err := CandidateFromOutcome(req.WorkOrderID, req.RequestID, out)
	if err != nil {
		return CandidateEnvelope{}, err
	}
	if base.Status == "operator_blocked" {
		return CandidateEnvelope{
			Schema: CandidateEnvelopeSchema, SourceCommit: req.SourceCommit, RequestID: req.RequestID,
			WorkOrderID: req.WorkOrderID, Attempt: req.Attempt, Workspace: req.Workspace,
			BaselineSHA: req.BaselineSHA, WorkOrderDigest: req.WorkOrderDigest, Status: base.Status,
			Acceptance: "external_required", ErrorClass: "backend_unavailable",
		}, nil
	}
	return CandidateEnvelope{
		Schema: CandidateEnvelopeSchema, SourceCommit: req.SourceCommit, RequestID: req.RequestID,
		WorkOrderID: req.WorkOrderID, Attempt: req.Attempt, Workspace: req.Workspace,
		BaselineSHA: req.BaselineSHA, WorkOrderDigest: req.WorkOrderDigest, Status: base.Status,
		Acceptance: "external_required", Text: base.Text, BackendID: base.BackendID, ModelID: base.ModelID,
		Edits: &EditEnvelope{
			SourceCommit: req.SourceCommit, RequestID: req.RequestID, WorkOrderID: req.WorkOrderID,
			Attempt: req.Attempt, Workspace: req.Workspace, BaselineSHA: req.BaselineSHA,
			WorkOrderDigest: req.WorkOrderDigest,
			Operations: []EditOperation{{Type: "write_text", Path: "src/replay.go", Content: "package src\n"}},
		},
	}, nil
}

func TestWinglessListenerFullRequestResponseReplayThroughBroker(t *testing.T) {
	session, req, _ := listenerFixture(t)
	reg := &broker.Registry{}
	if err := reg.Register(broker.Entry{Backend: &inference.Mock{Name: "fast", Features: []string{"code"}, Reply: "fixture-ok"}, Class: "mock", Tier: "fast"}); err != nil {
		t.Fatal(err)
	}
	engine := replayEngine{runner: broker.Runner{Registry: reg, Metrics: fixtureMetrics{}}}
	listener, err := NewListener(session, engine)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(listener.Handler())
	defer server.Close()
	raw, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, server.URL+TransportPath, bytes.NewReader(raw))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-CKB-Wingless-Session", session)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || resp.Header.Get("X-CKB-Wingless-Session") != session {
		t.Fatalf("status=%d session=%q", resp.StatusCode, resp.Header.Get("X-CKB-Wingless-Session"))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var got CandidateEnvelope
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got.Status != "result_ready" || got.Acceptance != "external_required" || got.BackendID != "fast" || got.ModelID != "fixture" || got.Edits == nil || len(got.Edits.Operations) != 1 {
		t.Fatalf("unexpected replay candidate: %+v", got)
	}
}

func TestWinglessListenerIsolatedProof(t *testing.T) {
	t.Run("exact_request", TestWinglessListenerAcceptsExactPlaneRequest)
	t.Run("body_consumed", TestWinglessListenerConsumesAndClosesBodyBeforeInference)
	t.Run("cancellation", TestWinglessListenerCancellationParentsInferenceToRequest)
	t.Run("candidate_authority", TestWinglessListenerRejectsAcceptanceClaimAndIdentityChange)
	t.Run("single_invoke", TestWinglessListenerCallsIntelligenceExactlyOnce)
	t.Run("no_authority_surface", TestWinglessListenerNoProcessListenQueueOrAcceptanceSurface)
	t.Run("full_replay", TestWinglessListenerFullRequestResponseReplayThroughBroker)
}
