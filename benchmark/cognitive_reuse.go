package benchmark

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spoonman136668-ai/Wingless/cognitive"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

type CognitiveCost struct {
	ModelCalls       int                       `json:"model_calls"`
	InputTokens      int                       `json:"input_tokens"`
	OutputTokens     int                       `json:"output_tokens"`
	LatencyMS        int64                     `json:"latency_ms"`
	MemoryRetrievals int                       `json:"memory_retrievals"`
	SkillReuses      int                       `json:"skill_reuses"`
	CognitivePasses  int                       `json:"cognitive_passes"`
	RAMFreeMinBytes  *uint64                   `json:"ram_free_min_bytes"`
	VRAMFreeMinBytes *uint64                   `json:"vram_free_min_bytes"`
	ModelTelemetry   *resources.ModelTelemetry `json:"model_telemetry"`
}

type CognitiveReuseReport struct {
	Schema                  string                   `json:"schema"`
	Repetitions             int                      `json:"repetitions"`
	Control                 CognitiveCost            `json:"control"`
	Experiment              CognitiveCost            `json:"experiment"`
	FirstSolutionModelCalls int                      `json:"first_solution_model_calls"`
	ReuseSolutionModelCalls int                      `json:"reuse_solution_model_calls"`
	SemanticCorrect         bool                     `json:"semantic_correct"`
	StrictCorrect           bool                     `json:"strict_correct"`
	ProtocolCorrect         bool                     `json:"protocol_correct"`
	SkillReuseCount         int                      `json:"skill_reuse_count"`
	ExperimentCheckpoints   map[string]CognitiveCost `json:"experiment_checkpoints"`
}

type cognitiveFixtureBackend struct {
	mu    sync.Mutex
	calls int
}

func (b *cognitiveFixtureBackend) ID() string                       { return "cr1a-fixture" }
func (b *cognitiveFixtureBackend) Capabilities() []string           { return []string{"text"} }
func (b *cognitiveFixtureBackend) Health(ctx context.Context) error { return ctx.Err() }
func (b *cognitiveFixtureBackend) EstimateCost(r inference.Request) inference.Cost {
	return inference.Cost{InputBytes: len(r.Context), MaxOutputTokens: r.MaxOutputTokens, Estimated: true}
}
func (b *cognitiveFixtureBackend) Cancel(string) bool { return false }
func (b *cognitiveFixtureBackend) Invoke(ctx context.Context, r inference.Request) (inference.Result, error) {
	if err := ctx.Err(); err != nil {
		return inference.Result{}, err
	}
	b.mu.Lock()
	b.calls++
	b.mu.Unlock()
	inTok, outTok := 12, 5
	return inference.Result{
		BackendID:   "cr1a-fixture",
		ModelID:     "deterministic-fixture",
		Status:      "completed",
		Text:        `{"sum":5}`,
		Usage:       inference.Usage{PromptTokens: &inTok, OutputTokens: &outTok},
		LatencyMS:   1,
		Termination: "stop",
	}, nil
}

func RunCognitiveReuseFixture(ctx context.Context, repetitions int) (CognitiveReuseReport, error) {
	if repetitions < 2 || repetitions > 1000 {
		return CognitiveReuseReport{}, errors.New("repetitions must be between 2 and 1000")
	}
	report := CognitiveReuseReport{
		Schema:                "wingless.cognitive-runtime-reuse-benchmark.v1",
		Repetitions:           repetitions,
		SemanticCorrect:       true,
		StrictCorrect:         true,
		ProtocolCorrect:       true,
		ExperimentCheckpoints: map[string]CognitiveCost{},
	}
	// CONTROL uses the existing inference backend seam directly with no cognitive reuse.
	controlBackend := &cognitiveFixtureBackend{}
	for i := 0; i < repetitions; i++ {
		req := fixtureRequest("control", i)
		result, err := controlBackend.Invoke(ctx, req.Base)
		if err != nil {
			return report, err
		}
		semantic, strict, protocol := validateFixtureOutput(result.Text)
		report.SemanticCorrect = report.SemanticCorrect && semantic
		report.StrictCorrect = report.StrictCorrect && strict
		report.ProtocolCorrect = report.ProtocolCorrect && protocol
		accumulateInferenceCost(&report.Control, result)
	}

	dir, err := os.MkdirTemp("", "wingless-cr1a-*")
	if err != nil {
		return report, err
	}
	defer os.RemoveAll(dir)
	store, err := cognitive.NewFileStore(dir)
	if err != nil {
		return report, err
	}
	experimentBackend := &cognitiveFixtureBackend{}
	experimentRuntime := cognitive.Runtime{
		Backend: experimentBackend,
		Skills:  store,
		Policy: cognitive.Policy{
			Version: "cr1a-experiment-v1", MaxPasses: 1, EnableSkills: true,
		},
	}

	first, err := experimentRuntime.Run(ctx, fixtureRequest("experiment", 0))
	if err != nil {
		return report, err
	}
	semantic, strict, protocol := validateFixtureOutput(first.Candidate.Text)
	report.SemanticCorrect = report.SemanticCorrect && semantic
	report.StrictCorrect = report.StrictCorrect && strict
	report.ProtocolCorrect = report.ProtocolCorrect && protocol
	if !strict || !protocol {
		return report, errors.New("fixture first solution failed deterministic validation")
	}
	report.FirstSolutionModelCalls = first.Evidence.ModelCalls
	firstCost := CognitiveCost{}
	accumulateCost(&firstCost, first.Evidence)
	report.ExperimentCheckpoints["1"] = firstCost
	accumulateCost(&report.Experiment, first.Evidence)

	skill, err := store.PutSkillCandidate(ctx, cognitive.Skill{
		Version:         1,
		Source:          "cr1a-deterministic-benchmark-validator",
		TaskFamily:      "json-arithmetic",
		Inputs:          []string{"sum request"},
		ExpectedOutputs: []string{`{"sum":5}`},
		Kind:            cognitive.SkillStaticText,
		Body:            first.Candidate.Text,
	})
	if err != nil {
		return report, err
	}
	sum := sha256.Sum256([]byte("cr1a-strict-json:" + first.Candidate.Text))
	_, err = store.PromoteSkill(ctx, skill.ID, cognitive.ValidationEvidence{
		ValidatorID:    "cr1a-strict-json-fixture",
		EvidenceSHA256: hex.EncodeToString(sum[:]),
	})
	if err != nil {
		return report, err
	}

	for i := 1; i < repetitions; i++ {
		res, err := experimentRuntime.Run(ctx, fixtureRequest("experiment", i))
		if err != nil {
			return report, err
		}
		semantic, strict, protocol := validateFixtureOutput(res.Candidate.Text)
		report.SemanticCorrect = report.SemanticCorrect && semantic
		report.StrictCorrect = report.StrictCorrect && strict
		report.ProtocolCorrect = report.ProtocolCorrect && protocol
		if res.Evidence.Reuse {
			report.SkillReuseCount++
		}
		report.ReuseSolutionModelCalls += res.Evidence.ModelCalls
		runCost := CognitiveCost{}
		accumulateCost(&runCost, res.Evidence)
		ordinal := i + 1
		if ordinal == 10 || ordinal == 50 || ordinal == 100 {
			report.ExperimentCheckpoints[itoa(ordinal)] = runCost
		}
		accumulateCost(&report.Experiment, res.Evidence)
	}
	return report, nil
}

func fixtureRequest(prefix string, i int) cognitive.RunRequest {
	id := prefix + "-cr1a-" + itoa(i)
	return cognitive.RunRequest{
		RequestID:  id,
		TaskFamily: "json-arithmetic",
		Base: inference.Request{
			ID:              id,
			ParentWorkID:    "cr1a-benchmark",
			Role:            "research_fixture",
			Context:         `Return only {"sum":5}.`,
			MaxContextBytes: 4096,
			MaxOutputTokens: 64,
			Deadline:        time.Now().Add(30 * time.Second),
			Workspace:       "cr1a-fixture",
			Capabilities:    []string{"text"},
		},
	}
}

func validateFixtureOutput(text string) (semantic bool, strict bool, protocol bool) {
	strict = text == `{"sum":5}`
	dec := json.NewDecoder(strings.NewReader(text))
	var obj map[string]json.RawMessage
	if err := dec.Decode(&obj); err != nil {
		return false, strict, false
	}
	raw, ok := obj["sum"]
	if !ok {
		return false, strict, false
	}
	var sum int
	if err := json.Unmarshal(raw, &sum); err != nil || sum != 5 {
		return false, strict, false
	}
	semantic = true
	if len(obj) != 1 {
		return semantic, strict, false
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return semantic, strict, false
	}
	return semantic, strict, true
}

func accumulateInferenceCost(dst *CognitiveCost, result inference.Result) {
	previousCalls := dst.ModelCalls
	dst.ModelCalls++
	if result.Usage.PromptTokens != nil {
		dst.InputTokens += *result.Usage.PromptTokens
	}
	if result.Usage.OutputTokens != nil {
		dst.OutputTokens += *result.Usage.OutputTokens
	}
	dst.LatencyMS += result.LatencyMS
	observeInferenceTelemetry(dst, result.Telemetry)
	if previousCalls == 0 {
		dst.ModelTelemetry = result.Telemetry.Model
	} else {
		// The benchmark does not assume per-call residency/counter values can be
		// combined across calls without backend-declared aggregation semantics.
		dst.ModelTelemetry = nil
	}
}

func accumulateCost(dst *CognitiveCost, ev cognitive.RunEvidence) {
	previousCalls := dst.ModelCalls
	dst.ModelCalls += ev.ModelCalls
	if ev.TotalInputTokens != nil {
		dst.InputTokens += *ev.TotalInputTokens
	}
	if ev.TotalOutputTokens != nil {
		dst.OutputTokens += *ev.TotalOutputTokens
	}
	dst.LatencyMS += ev.TotalLatencyMS
	dst.MemoryRetrievals += ev.Session.MemoryRetrievals
	dst.SkillReuses += ev.Session.SkillReuses
	dst.CognitivePasses += ev.Session.CognitivePasses
	for _, pass := range ev.Passes {
		observeInferenceTelemetry(dst, pass.Telemetry)
	}
	if previousCalls == 0 && ev.ModelCalls == 1 {
		dst.ModelTelemetry = ev.Session.ResourceTelemetry
	} else if ev.ModelCalls > 0 && dst.ModelCalls > 1 {
		dst.ModelTelemetry = nil
	}
}

func observeInferenceTelemetry(dst *CognitiveCost, telemetry inference.Telemetry) {
	for _, snapshot := range []*resources.Metrics{telemetry.Before, telemetry.Peak, telemetry.After} {
		if snapshot == nil {
			continue
		}
		dst.RAMFreeMinBytes = minKnownUint(dst.RAMFreeMinBytes, snapshot.RAMFree)
		dst.VRAMFreeMinBytes = minKnownUint(dst.VRAMFreeMinBytes, snapshot.VRAMFree)
	}
}

func minKnownUint(current, candidate *uint64) *uint64 {
	if candidate == nil {
		return current
	}
	if current == nil || *candidate < *current {
		v := *candidate
		return &v
	}
	return current
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
