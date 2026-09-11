package ckbplane

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

func validEnvelope() ReplayEnvelope {
	return ReplayEnvelope{
		Order: WorkOrderContract{
			ID:                   "V11110_STAGE_TEST_R1",
			Title:                "fixture",
			Instructions:         []string{"produce candidate only"},
			BaselineSHA:          "01260dda4bdebc454921cdf0c5d3d0eeafef7629",
			WorkerBranch:         "codex/test",
			AllowedPaths:         []string{"pkg/**"},
			ProtectAcceptedRefs:  true,
			RequireCleanBaseline: true,
		},
		Workspace: WorkspaceBinding{
			WorkOrderID: "V11110_STAGE_TEST_R1",
			Attempt:     2,
			Path:        `C:\plane\workspaces\V11110_STAGE_TEST_R1\2\worktree`,
			BaselineSHA: "01260dda4bdebc454921cdf0c5d3d0eeafef7629",
		},
		Context: "bounded source context",
		Policy: RequestPolicy{
			Role:            "implement",
			MaxContextBytes: 4096,
			MaxOutputTokens: 512,
			Deadline:        time.Now().Add(2 * time.Minute),
			Capabilities:    []string{"text"},
			AllowedBackends: []string{"mock-fast"},
			Resources:       resources.Policy{},
		},
	}
}

func TestTranslateValidPlaneBinding(t *testing.T) {
	e := validEnvelope()
	req, route, err := Translate(e)
	if err != nil {
		t.Fatal(err)
	}
	if req.ParentWorkID != e.Order.ID || req.Workspace != e.Workspace.Path {
		t.Fatalf("identity translation mismatch: %#v", req)
	}
	if req.ID != "V11110_STAGE_TEST_R1-wingless-2" {
		t.Fatalf("unexpected request id %q", req.ID)
	}
	if route.Repairs != 0 || route.MaxRepairs != 0 {
		t.Fatalf("plane-driven work must not consume Wingless repair budget: %#v", route)
	}
}

func TestTranslateTrustedRoutingAndResources(t *testing.T) {
	e := validEnvelope()
	e.Policy.AllowedBackends = []string{"fast", "deep"}
	e.Policy.AllowDeep = true
	e.Policy.AllowFallback = true
	e.Policy.Reason = "low_confidence"
	e.Policy.Resources = resources.Policy{MinRAM: 1024, MinVRAM: 2048, MinDisk: 4096, MaxCPU: 80}
	req, route, err := Translate(e)
	if err != nil {
		t.Fatal(err)
	}
	if req.Resources != e.Policy.Resources {
		t.Fatalf("resource policy widened or lost: got %#v want %#v", req.Resources, e.Policy.Resources)
	}
	if !route.AllowDeep || !route.AllowFallback || route.Reason != "low_confidence" || len(route.AllowedBackends) != 2 {
		t.Fatalf("trusted routing policy mismatch: %#v", route)
	}
	if route.MaxRepairs != 0 {
		t.Fatalf("plane owns work-order retry; got MaxRepairs=%d", route.MaxRepairs)
	}
}

func TestTranslateRejectsStaleWorkspaceRevision(t *testing.T) {
	e := validEnvelope()
	e.Workspace.BaselineSHA = strings.Repeat("a", 40)
	if _, _, err := Translate(e); err == nil || !strings.Contains(err.Error(), "stale workspace revision") {
		t.Fatalf("expected stale revision rejection, got %v", err)
	}
}

func TestTranslateRejectsWorkspaceIdentityMismatch(t *testing.T) {
	e := validEnvelope()
	e.Workspace.WorkOrderID = "OTHER"
	if _, _, err := Translate(e); err == nil || !strings.Contains(err.Error(), "workspace work-order mismatch") {
		t.Fatalf("expected work-order mismatch, got %v", err)
	}
}

func TestTranslateRejectsMissingPlaneProtectionGate(t *testing.T) {
	e := validEnvelope()
	e.Order.ProtectAcceptedRefs = false
	if _, _, err := Translate(e); err == nil || !strings.Contains(err.Error(), "protection gates") {
		t.Fatalf("expected protection gate rejection, got %v", err)
	}
}

func TestTranslateRejectsUntrustedBackendPolicy(t *testing.T) {
	e := validEnvelope()
	e.Policy.AllowedBackends = nil
	if _, _, err := Translate(e); err == nil || !strings.Contains(err.Error(), "backend policy") {
		t.Fatalf("expected backend policy rejection, got %v", err)
	}
}

func TestDecodeReplayEnvelopeRejectsUnknownField(t *testing.T) {
	e := validEnvelope()
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	raw := strings.TrimSuffix(string(b), "}") + `,"unknown":true}`
	if _, err := DecodeReplayEnvelope([]byte(raw)); err == nil {
		t.Fatal("unknown replay field must fail closed")
	}
}

type fakeCanceler struct {
	got string
	ok  bool
}

func (f *fakeCanceler) Cancel(id string) bool {
	f.got = id
	return f.ok
}

func TestCancelMapsPlaneAttemptToWinglessRequest(t *testing.T) {
	f := &fakeCanceler{ok: true}
	if err := Cancel(f, "V11110_STAGE_TEST_R1", 3); err != nil {
		t.Fatal(err)
	}
	if f.got != "V11110_STAGE_TEST_R1-wingless-3" {
		t.Fatalf("unexpected cancel id %q", f.got)
	}
}

func TestCandidateNeverClaimsPlaneAcceptance(t *testing.T) {
	out := broker.Outcome{
		Status:     "result_ready",
		Acceptance: "external_required",
		Events: []broker.Event{{Result: inference.Result{
			Status:    "completed",
			Text:      "candidate patch",
			BackendID: "mock-fast",
			ModelID:   "fixture",
		}}},
	}
	c, err := CandidateFromOutcome("V11110_STAGE_TEST_R1", "V11110_STAGE_TEST_R1-wingless-1", out)
	if err != nil {
		t.Fatal(err)
	}
	if c.Acceptance != "external_required" {
		t.Fatalf("candidate acceptance widened: %#v", c)
	}
}

func TestCandidatePreservesBlockedBackendOutcomeWithoutAcceptance(t *testing.T) {
	out := broker.Outcome{Status: "operator_blocked", Acceptance: "external_required"}
	c, err := CandidateFromOutcome("V11110_STAGE_TEST_R1", "V11110_STAGE_TEST_R1-wingless-1", out)
	if err != nil {
		t.Fatal(err)
	}
	if c.Status != "operator_blocked" || c.Acceptance != "external_required" || c.Text != "" {
		t.Fatalf("blocked outcome widened: %#v", c)
	}
}

func TestCandidateRejectsAcceptanceClaim(t *testing.T) {
	out := broker.Outcome{Status: "result_ready", Acceptance: "accepted"}
	if _, err := CandidateFromOutcome("V11110_STAGE_TEST_R1", "req", out); err == nil {
		t.Fatal("Wingless acceptance claim must be rejected")
	}
}

func TestCandidateRejectsUnknownTerminalState(t *testing.T) {
	out := broker.Outcome{Status: "completed", Acceptance: "external_required"}
	if _, err := CandidateFromOutcome("V11110_STAGE_TEST_R1", "req", out); err == nil {
		t.Fatal("unknown terminal state must fail closed")
	}
}
