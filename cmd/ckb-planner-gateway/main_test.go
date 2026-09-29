package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/spoonman136668-ai/Wingless/reasoner"
)

type mockReasoner struct {
	last reasoner.Request
	text string
	err  error
}

func (m *mockReasoner) Invoke(_ context.Context, r reasoner.Request) (reasoner.Result, error) {
	m.last = r
	if m.err != nil {
		return reasoner.Result{}, m.err
	}
	return reasoner.Result{Text: m.text}, nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestModelsLoopbackOnly(t *testing.T) {
	g := &gateway{}
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/v1/models", nil)
	req.RemoteAddr = "203.0.113.7:1234"
	w := httptest.NewRecorder()
	g.models(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status=%d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "http://127.0.0.1/v1/models", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	w = httptest.NewRecorder()
	g.models(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), gatewayModel) {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestChatVerifiesPublicMaterialAndInvokesReasoner(t *testing.T) {
	content := "public evidence bytes\n"
	sha := sha256hex([]byte(content))
	commit := strings.Repeat("a", 40)

	planner := `{
  "schema":"ckb-plane.research-planner-request.v1",
  "instruction":"bounded planner instruction",
  "envelope":{
    "schema":"ckb-plane.research-planning-envelope.v1",
    "id":"UP_TEST",
    "project":"Wingless",
    "lane":"B",
    "authority_generation":6,
    "baseline_sha":"` + commit + `",
    "north_star":{"id":"north","path":"docs/north.md","sha256":"` + sha + `","commit":"` + commit + `"},
    "evidence":[{"id":"e1","path":"docs/evidence.md","sha256":"` + sha + `","commit":"` + commit + `"}],
    "capability_targets":["cap"],
    "allowed_harnesses":["harness"],
    "allowed_roots":["docs"],
    "max_candidates":2,
    "max_seeds":2,
    "max_compute_seconds":30,
    "min_evidence_refs":1,
    "execution_policy":"sealed-experiment",
    "cross_lane_allowed":false
  },
  "north_star":"` + strings.ReplaceAll(content, "\n", "\\n") + `",
  "evidence":[{"id":"e1","path":"docs/evidence.md","commit":"` + commit + `","sha256":"` + sha + `","content":"` + strings.ReplaceAll(content, "\n", "\\n") + `"}]
}`

	mock := &mockReasoner{text: `{"schema":"ckb-plane.research-experiment-proposal.v1"}`}
	g := &gateway{
		client: mock,
		publicClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(content)),
				Request:    r,
			}, nil
		})},
		repository: "spoonman136668-ai/Wingless",
		timeout:    2 * time.Second,
	}

	body := `{
  "model":"nemotron-3-ultra-planner",
  "messages":[
    {"role":"system","content":"system"},
    {"role":"user","content":` + quoteJSON(planner) + `}
  ],
  "temperature":0,
  "max_tokens":4096,
  "stream":false,
  "response_format":{"type":"json_object"}
}`

	req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1/v1/chat/completions", bytes.NewBufferString(body))
	req.RemoteAddr = "127.0.0.1:4321"
	w := httptest.NewRecorder()
	g.chat(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if mock.last.Schema != reasoner.RequestSchema || mock.last.Role != reasoner.RolePlan ||
		mock.last.DataClass != "public-repository" || mock.last.Project != "Wingless" {
		t.Fatalf("unexpected reasoner request: %+v", mock.last)
	}
	if !strings.Contains(w.Body.String(), "ckb-plane.research-experiment-proposal.v1") {
		t.Fatalf("response=%s", w.Body.String())
	}
}

func TestValidatePlannerRejectsUncommittedPublicMaterial(t *testing.T) {
	var p plannerRequest
	p.Schema = plannerSchema
	p.Envelope.ID = "UP_TEST"
	p.Envelope.Project = "Wingless"
	p.Envelope.Lane = "B"
	p.Envelope.AuthorityGeneration = 6
	p.Envelope.BaselineSHA = strings.Repeat("a", 40)
	p.Envelope.NorthStar = artifactRef{ID: "north", Path: "docs/north.md", SHA256: strings.Repeat("b", 64)}
	p.NorthStar = "x"
	if err := validatePlannerRequest(p); err == nil {
		t.Fatal("expected rejection")
	}
}

func quoteJSON(s string) string {
	var b bytes.Buffer
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString("\\\\")
		case '"':
			b.WriteString("\\\"")
		case '\n':
			b.WriteString("\\n")
		case '\r':
			b.WriteString("\\r")
		case '\t':
			b.WriteString("\\t")
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
