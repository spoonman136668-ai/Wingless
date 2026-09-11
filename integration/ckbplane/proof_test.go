package ckbplane

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"github.com/spoonman136668-ai/Wingless/worker"
)

type proofMetrics struct{ m resources.Metrics }

func (p proofMetrics) Snapshot() (resources.Metrics, error) { return p.m, nil }

func proofU64(v uint64) *uint64    { return &v }
func proofF64(v float64) *float64  { return &v }

// recordedPlaneEnvelope normalizes the immutable identity/scope from the frozen
// ckb-r4a-source-proof profile into the inactive adapter's replay contract. It
// is recorded input, not a live plane request.
func recordedPlaneEnvelope() ReplayEnvelope {
	return ReplayEnvelope{
		Order: WorkOrderContract{
			ID:           "CKB_R4A_SOURCE_ONLY_WORKER_PROOF_R1",
			Title:        "R4a source-only worker proof: create docs/CKB_PLANE_R4A_SOURCE_ONLY_PROOF.md and pkg/r4a_source_only_worker_proof_test.go only",
			Instructions: []string{
				"Create docs/CKB_PLANE_R4A_SOURCE_ONLY_PROOF.md containing the marker CKB_R4A_SOURCE_ONLY_WORKER_PROOF_R1.",
				"Create pkg/r4a_source_only_worker_proof_test.go with deterministic source-only proof coverage.",
			},
			BaselineSHA: "01260dda4bdebc454921cdf0c5d3d0eeafef7629",
			WorkerBranch: "codex/r4a-source-only-worker-proof-r1",
			AllowedPaths: []string{
				"docs/CKB_PLANE_R4A_SOURCE_ONLY_PROOF.md",
				"pkg/r4a_source_only_worker_proof_test.go",
			},
			ForbiddenPaths: []string{
				".git/**", ".github/**", "runtime/**", "**/*coinbase*", "**/*broker*",
				"**/*credential*", "**/*secret*", "**/*tradeguard*", "**/*accounting*",
				"**/*reconciliation*", "**/*deployment*", "**/*.exe", "**/*.dll",
			},
			ProtectAcceptedRefs:  true,
			RequireCleanBaseline: true,
		},
		Workspace: WorkspaceBinding{
			WorkOrderID: "CKB_R4A_SOURCE_ONLY_WORKER_PROOF_R1",
			Attempt:     1,
			Path:        `C:\ckb-plane-isolated-proof\CKB_R4A_SOURCE_ONLY_WORKER_PROOF_R1\1\worktree`,
			BaselineSHA: "01260dda4bdebc454921cdf0c5d3d0eeafef7629",
		},
		Context: "Recorded frozen-plane source-only proof context. Produce candidate reasoning only; plane acceptance remains external.",
		Policy: RequestPolicy{
			Role:             "implement",
			MaxContextBytes:  8192,
			MaxOutputTokens:  1024,
			Deadline:         time.Now().Add(2 * time.Minute),
			Capabilities:     []string{"text"},
			AllowedBackends:  []string{"mock-fast"},
			AllowDeep:        false,
			AllowFallback:    false,
			Reason:           "routine",
			Resources: resources.Policy{
				MinRAM:  1024,
				MinVRAM: 1024,
				MinDisk: 1024,
				MaxCPU:  90,
			},
		},
	}
}

func proofProvider() proofMetrics {
	return proofMetrics{m: resources.Metrics{
		RAMFree:    proofU64(16 << 30),
		VRAMFree:   proofU64(8 << 30),
		DiskFree:   proofU64(100 << 30),
		CPUPercent: proofF64(10),
	}}
}

func proofRegistry(t *testing.T, entries ...broker.Entry) *broker.Registry {
	t.Helper()
	r := &broker.Registry{}
	for _, e := range entries {
		if err := r.Register(e); err != nil {
			t.Fatalf("register backend: %v", err)
		}
	}
	return r
}

func TestIsolatedProofCandidateFlow(t *testing.T) {
	fast := &inference.Mock{Name: "mock-fast", Features: []string{"text"}, Reply: "candidate-only patch plan"}
	reg := proofRegistry(t, broker.Entry{Backend: fast, Class: "mock", Tier: "fast"})
	runner := broker.Runner{Registry: reg, Metrics: proofProvider()}

	env := recordedPlaneEnvelope()
	req, policy, err := Translate(env)
	if err != nil {
		t.Fatal(err)
	}
	if policy.MaxRepairs != 0 || policy.Repairs != 0 {
		t.Fatalf("plane-driven proof acquired Wingless retry budget: %#v", policy)
	}

	out := runner.Run(context.Background(), req, policy)
	if out.Status != "result_ready" || out.Acceptance != "external_required" {
		t.Fatalf("unexpected candidate outcome: %#v", out)
	}
	if len(out.Events) != 1 || out.Events[0].Route.Backend != "mock-fast" {
		t.Fatalf("expected exactly one fast candidate attempt: %#v", out.Events)
	}
	candidate, err := CandidateFromOutcome(env.Order.ID, req.ID, out)
	if err != nil {
		t.Fatal(err)
	}
	if candidate.Text != "candidate-only patch plan" || candidate.Acceptance != "external_required" {
		t.Fatalf("candidate projection widened or lost output: %#v", candidate)
	}
}

func TestIsolatedProofAuthorizedDeepFallback(t *testing.T) {
	deep := &inference.Mock{Name: "mock-deep", Features: []string{"text"}, Unhealthy: true}
	fast := &inference.Mock{Name: "mock-fast", Features: []string{"text"}, Reply: "fallback candidate"}
	reg := proofRegistry(t,
		broker.Entry{Backend: deep, Class: "mock", Tier: "deep"},
		broker.Entry{Backend: fast, Class: "mock", Tier: "fast"},
	)
	runner := broker.Runner{Registry: reg, Metrics: proofProvider()}

	env := recordedPlaneEnvelope()
	env.Policy.AllowedBackends = []string{"mock-deep", "mock-fast"}
	env.Policy.AllowDeep = true
	env.Policy.AllowFallback = true
	env.Policy.Reason = "low_confidence"
	req, policy, err := Translate(env)
	if err != nil {
		t.Fatal(err)
	}
	out := runner.Run(context.Background(), req, policy)
	if out.Status != "result_ready" || len(out.Events) != 1 {
		t.Fatalf("fallback did not produce one candidate attempt: %#v", out)
	}
	route := out.Events[0].Route
	if route.Backend != "mock-fast" || route.Tier != "fast" || route.Reason != "authorized_fast_fallback" {
		t.Fatalf("unexpected fallback route: %#v", route)
	}
	if len(route.Rejected) == 0 || !strings.Contains(strings.Join(route.Rejected, " "), "mock-deep: health") {
		t.Fatalf("deep failure evidence missing: %#v", route)
	}
}

func TestIsolatedProofBackendFailureNoDuplicateRetry(t *testing.T) {
	fast := &inference.Mock{Name: "mock-fast", Features: []string{"text"}, Failure: errors.New("fixture backend failure")}
	reg := proofRegistry(t, broker.Entry{Backend: fast, Class: "mock", Tier: "fast"})
	runner := broker.Runner{Registry: reg, Metrics: proofProvider()}

	env := recordedPlaneEnvelope()
	req, policy, err := Translate(env)
	if err != nil {
		t.Fatal(err)
	}
	out := runner.Run(context.Background(), req, policy)
	if out.Status != "operator_blocked" || out.Acceptance != "external_required" {
		t.Fatalf("backend failure escaped candidate-only boundary: %#v", out)
	}
	if len(out.Events) != 1 {
		t.Fatalf("plane-owned retry was duplicated by Wingless: attempts=%d %#v", len(out.Events), out.Events)
	}
	if !strings.Contains(out.Events[0].Error, "fixture backend failure") {
		t.Fatalf("backend failure evidence missing: %#v", out.Events[0])
	}
	candidate, err := CandidateFromOutcome(env.Order.ID, req.ID, out)
	if err != nil {
		t.Fatal(err)
	}
	if candidate.Status != "operator_blocked" || candidate.Acceptance != "external_required" {
		t.Fatalf("blocked candidate widened authority: %#v", candidate)
	}
}

func TestIsolatedProofCancellationPropagation(t *testing.T) {
	fast := &inference.Mock{Name: "mock-fast", Features: []string{"text"}, Reply: "must not complete", Delay: 5 * time.Second}
	reg := proofRegistry(t, broker.Entry{Backend: fast, Class: "mock", Tier: "fast"})
	runner := broker.Runner{Registry: reg, Metrics: proofProvider()}

	env := recordedPlaneEnvelope()
	env.Workspace.Attempt = 3
	req, policy, err := Translate(env)
	if err != nil {
		t.Fatal(err)
	}
	svc := worker.New(runner, policy)
	id, err := svc.SubmitWork(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if id != req.ID {
		t.Fatalf("service changed deterministic adapter request ID: got %q want %q", id, req.ID)
	}
	if err := Cancel(svc, env.Order.ID, env.Workspace.Attempt); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		status, err := svc.Status(req.ID)
		if err != nil {
			t.Fatal(err)
		}
		if status != "running" {
			if status != "operator_blocked" {
				t.Fatalf("unexpected post-cancel Wingless status %q", status)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("cancellation did not propagate to Wingless request")
		}
		time.Sleep(10 * time.Millisecond)
	}

	out, err := svc.Result(req.ID)
	if err != nil {
		t.Fatal(err)
	}
	if out.Acceptance != "external_required" || len(out.Events) != 1 {
		t.Fatalf("cancel outcome widened authority or retried: %#v", out)
	}
}
