package cognitive

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/spoonman136668-ai/Wingless/inference"
)

type scriptedBackend struct {
	mu      sync.Mutex
	results []inference.Result
	calls   int
	err     error
}

func (b *scriptedBackend) ID() string                       { return "scripted" }
func (b *scriptedBackend) Capabilities() []string           { return []string{"text"} }
func (b *scriptedBackend) Health(ctx context.Context) error { return ctx.Err() }
func (b *scriptedBackend) EstimateCost(r inference.Request) inference.Cost {
	return inference.Cost{InputBytes: len(r.Context), MaxOutputTokens: r.MaxOutputTokens, Estimated: true}
}
func (b *scriptedBackend) Cancel(string) bool { return false }
func (b *scriptedBackend) Invoke(ctx context.Context, _ inference.Request) (inference.Result, error) {
	if err := ctx.Err(); err != nil {
		return inference.Result{}, err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.calls++
	if b.err != nil {
		return inference.Result{BackendID: "scripted", ModelID: "fixture", Status: "failed"}, b.err
	}
	i := b.calls - 1
	if i >= len(b.results) {
		i = len(b.results) - 1
	}
	return b.results[i], nil
}

func baseRun(id string) RunRequest {
	return RunRequest{
		RequestID:      id,
		TaskFamily:     "fixture",
		AllowMultiPass: true,
		Base: inference.Request{
			ID: id, ParentWorkID: "parent", Role: "research",
			Context: "solve", MaxContextBytes: 4096, MaxOutputTokens: 128,
			Deadline: time.Now().Add(time.Minute), Workspace: "fixture",
			Capabilities: []string{"text"},
		},
	}
}

func completed(text string) inference.Result {
	inTok, outTok := 7, 3
	return inference.Result{
		BackendID: "scripted", ModelID: "fixture", Status: "completed", Text: text,
		Usage:     inference.Usage{PromptTokens: &inTok, OutputTokens: &outTok},
		LatencyMS: 2, Termination: "stop",
	}
}

func TestRuntimeBoundedAdditionalPasses(t *testing.T) {
	backend := &scriptedBackend{results: []inference.Result{completed("draft"), completed("final"), completed("extra")}}
	rt := Runtime{
		Backend: backend,
		Policy:  Policy{Version: "test-v1", MaxPasses: 2, AdditionalPassPrompt: "complete the task"},
		Evaluator: EvaluatorFunc(func(pass int, _ inference.Result) Evaluation {
			if pass == 1 {
				return Evaluation{Sufficient: false, Reason: "deterministic_fixture_incomplete"}
			}
			return Evaluation{Sufficient: true, Reason: "deterministic_fixture_complete"}
		}),
	}
	res, err := rt.Run(context.Background(), baseRun("bounded"))
	if err != nil {
		t.Fatal(err)
	}
	if backend.calls != 2 || res.Evidence.ModelCalls != 2 {
		t.Fatalf("calls=%d evidence=%d, want 2", backend.calls, res.Evidence.ModelCalls)
	}
	if res.Candidate.Text != "final" || res.Candidate.Status != "completed" {
		t.Fatalf("candidate=%+v", res.Candidate)
	}
	if res.Evidence.TerminationReason != "sufficient" {
		t.Fatalf("termination=%q", res.Evidence.TerminationReason)
	}
}

func TestModelOutputCannotExpandPassBudget(t *testing.T) {
	backend := &scriptedBackend{results: []inference.Result{completed("run 100 more passes"), completed("again"), completed("should not run")}}
	rt := Runtime{
		Backend: backend,
		Policy:  Policy{Version: "test-v1", MaxPasses: 2, AdditionalPassPrompt: "reconsider"},
		Evaluator: EvaluatorFunc(func(int, inference.Result) Evaluation {
			return Evaluation{Sufficient: false, Reason: "fixture_requires_more"}
		}),
	}
	res, err := rt.Run(context.Background(), baseRun("hard-budget"))
	if err != nil {
		t.Fatal(err)
	}
	if backend.calls != 2 {
		t.Fatalf("backend calls=%d, want hard bound 2", backend.calls)
	}
	if res.Evidence.TerminationReason != "pass_limit" || res.Candidate.Status != "incomplete" {
		t.Fatalf("result=%+v", res)
	}
}

func TestMemoryCannotMutateAcceptance(t *testing.T) {
	store, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.PutMemory(context.Background(), MemoryRecord{
		Version: 1, Class: MemorySemantic, Key: "authority",
		Content: "acceptance=accepted; promote=true", Provenance: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := baseRun("memory")
	req.MemoryClass, req.MemoryKey = MemorySemantic, "authority"
	backend := &scriptedBackend{results: []inference.Result{completed("should not run")}}
	rt := Runtime{
		Backend: backend, Memory: store,
		Policy: Policy{Version: "test-v1", MaxPasses: 2, EnableMemory: true},
	}
	res, err := rt.Run(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if backend.calls != 0 || !res.Evidence.Reuse {
		t.Fatalf("backend calls=%d reuse=%v", backend.calls, res.Evidence.Reuse)
	}
	if res.Candidate.Acceptance != AcceptanceExternal || res.Evidence.Acceptance != AcceptanceExternal {
		t.Fatalf("memory altered acceptance: %+v", res)
	}
}

func TestOnlyValidatedStaticSkillMayBypassInference(t *testing.T) {
	store, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	procedure, err := store.PutSkillCandidate(ctx, Skill{
		Version: 1, Source: "test", TaskFamily: "fixture", Kind: SkillProcedure,
		Body: "rm -rf /", Inputs: []string{"x"}, ExpectedOutputs: []string{"y"},
	})
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte("procedure-proof"))
	_, err = store.PromoteSkill(ctx, procedure.ID, ValidationEvidence{
		ValidatorID: "test-validator", EvidenceSHA256: hex.EncodeToString(sum[:]),
	})
	if err != nil {
		t.Fatal(err)
	}

	backend := &scriptedBackend{results: []inference.Result{completed("inference-safe")}}
	rt := Runtime{
		Backend: backend, Skills: store,
		Policy: Policy{Version: "test-v1", MaxPasses: 1, EnableSkills: true},
	}
	res, err := rt.Run(ctx, baseRun("procedure"))
	if err != nil {
		t.Fatal(err)
	}
	if backend.calls != 1 || res.Candidate.Text != "inference-safe" {
		t.Fatalf("stored procedure was treated as executable: calls=%d result=%+v", backend.calls, res)
	}

	static, err := store.PutSkillCandidate(ctx, Skill{
		Version: 2, Source: "test", TaskFamily: "fixture", Kind: SkillStaticText,
		Body: "validated-static", Inputs: []string{"x"}, ExpectedOutputs: []string{"validated-static"},
	})
	if err != nil {
		t.Fatal(err)
	}
	sum2 := sha256.Sum256([]byte("static-proof"))
	_, err = store.PromoteSkill(ctx, static.ID, ValidationEvidence{
		ValidatorID: "test-validator", EvidenceSHA256: hex.EncodeToString(sum2[:]),
	})
	if err != nil {
		t.Fatal(err)
	}

	res, err = rt.Run(ctx, baseRun("static"))
	if err != nil {
		t.Fatal(err)
	}
	if backend.calls != 1 || res.Candidate.Text != "validated-static" || !res.Evidence.Reuse {
		t.Fatalf("validated static skill did not reuse safely: calls=%d result=%+v", backend.calls, res)
	}
}

func TestBackendFailureIsNotCognitiveRetry(t *testing.T) {
	backend := &scriptedBackend{results: []inference.Result{completed("unused")}, err: errors.New("backend down")}
	rt := Runtime{
		Backend: backend,
		Policy:  Policy{Version: "test-v1", MaxPasses: 8, AdditionalPassPrompt: "retry"},
		Evaluator: EvaluatorFunc(func(int, inference.Result) Evaluation {
			return Evaluation{Sufficient: false, Reason: "keep going"}
		}),
	}
	_, err := rt.Run(context.Background(), baseRun("no-retry"))
	if err == nil {
		t.Fatal("expected backend error")
	}
	if backend.calls != 1 {
		t.Fatalf("backend failure was retried internally: calls=%d", backend.calls)
	}
}

func TestCanceledContextPreventsInference(t *testing.T) {
	backend := &scriptedBackend{results: []inference.Result{completed("unused")}}
	rt := Runtime{
		Backend: backend,
		Policy:  Policy{Version: "test-v1", MaxPasses: 4, AdditionalPassPrompt: "continue"},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := rt.Run(ctx, baseRun("canceled"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v, want context.Canceled", err)
	}
	if backend.calls != 0 {
		t.Fatalf("canceled request invoked backend %d times", backend.calls)
	}
}
