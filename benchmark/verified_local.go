package benchmark

import (
	"context"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"runtime"
	"time"
)

type VerifiedLocalRow struct {
	Task                 string           `json:"task"`
	Repetition           int              `json:"repetition"`
	RawText              string           `json:"raw_text"`
	CandidateText        string           `json:"candidate_text"`
	CandidateNormalized  bool             `json:"candidate_normalized"`
	RawSemanticCorrect   bool             `json:"raw_semantic_correct"`
	RawProtocolCompliant bool             `json:"raw_protocol_compliant"`
	RawStrictCorrect     bool             `json:"raw_strict_correct"`
	CandidateVerified    bool             `json:"candidate_verified"`
	Verification         string           `json:"verification"`
	Error                string           `json:"error,omitempty"`
	Result               inference.Result `json:"result"`
}

type VerifiedLocalReport struct {
	Provenance Provenance        `json:"provenance"`
	Profile    string            `json:"profile"`
	Rows       []VerifiedLocalRow `json:"rows"`

	Version                int       `json:"version"`
	OS                     string    `json:"os"`
	Arch                   string    `json:"arch"`
	CPUs                   int       `json:"logical_cpus"`
	StartedAt              time.Time `json:"started_at"`
	RawStrictCorrect       int       `json:"raw_strict_correct"`
	RawSemanticCorrect     int       `json:"raw_semantic_correct"`
	RawProtocolCompliant   int       `json:"raw_protocol_compliant"`
	CandidateVerified      int       `json:"candidate_verified"`
	CandidateNormalized    int       `json:"candidate_normalized"`
	GeneratedCodeExecuted  bool      `json:"generated_code_executed"`
	GeneratedToolsExecuted bool      `json:"generated_tools_executed"`
	Acceptance             string    `json:"acceptance"`
}

type localFixtureVerifier struct{ fixture localFixture }

func (v localFixtureVerifier) Verify(_ context.Context, r inference.Result) error {
	if !v.fixture.protocol(r.Text) || !v.fixture.check(r.Text) {
		return fmt.Errorf("fixture output mismatch")
	}
	return nil
}

// RunLocalVerified executes the fixed public microfixtures through the real
// broker/verifier path. Raw model output remains separately scored and retained;
// only the trusted deterministic verifier sees a bounded verification candidate.
// It executes no generated code or tools and permits no repair attempts.
func RunLocalVerified(ctx context.Context, b inference.InferenceBackend, h resources.Provider, repeats int) (VerifiedLocalReport, error) {
	out := VerifiedLocalReport{
		Version:    1,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		CPUs:       runtime.NumCPU(),
		StartedAt:  time.Now().UTC(),
		Acceptance: "external_required",
	}
	if b == nil || h == nil {
		return out, fmt.Errorf("benchmark backend/provider required")
	}
	if repeats < 1 || repeats > 3 {
		return out, fmt.Errorf("repeats must be 1..3")
	}

	reg := &broker.Registry{}
	if err := reg.Register(broker.Entry{Backend: b, Class: "local-openai-compatible", Tier: "fast"}); err != nil {
		return out, err
	}
	runner := broker.Runner{Registry: reg, Metrics: h}
	policy := broker.Policy{AllowedBackends: []string{b.ID()}, AllowDeep: false, AllowFallback: false, MaxRepairs: 0}
	fixtures := codingFixtures()

	for n := 0; n < repeats; n++ {
		plan := ""
		for _, base := range fixtures {
			f := base
			if f.id == "plan-implementation" {
				f.prompt += "\nProposed plan (untrusted):\n" + plan
			}
			if err := ctx.Err(); err != nil {
				return out, err
			}

			runner.Verifier = localFixtureVerifier{fixture: f}
			req := inference.Request{
				ID:              fmt.Sprintf("local-verified-%d-%s", n, f.id),
				ParentWorkID:    "local-public-microfixtures",
				Role:            "code",
				Context:         f.prompt,
				MaxContextBytes: 4096,
				MaxOutputTokens: 128,
				Deadline:        time.Now().Add(60 * time.Second),
				Workspace:       "no-workspace-access",
				Capabilities:    []string{"code"},
				Resources:       resources.Policy{MinRAM: 256 << 20},
			}

			outcome := runner.Run(ctx, req, policy)
			row := VerifiedLocalRow{Task: f.id, Repetition: n}
			if len(outcome.Events) == 0 {
				row.Error = "broker produced no event"
				out.Rows = append(out.Rows, row)
				continue
			}

			event := outcome.Events[len(outcome.Events)-1]
			row.Result = event.Result
			row.CandidateText = event.Result.Text
			row.RawText = event.Result.Text
			if event.RawModelText != "" {
				row.RawText = event.RawModelText
				row.CandidateNormalized = true
				out.CandidateNormalized++
			}
			row.Verification = event.Verification
			row.Error = event.Error
			row.RawSemanticCorrect = f.check(semanticCandidate(row.RawText))
			row.RawProtocolCompliant = f.protocol(row.RawText)
			row.RawStrictCorrect = row.RawSemanticCorrect && row.RawProtocolCompliant
			row.CandidateVerified = outcome.Status == "verified_candidate" && event.Verification == "fixture_pass"

			if row.RawSemanticCorrect {
				out.RawSemanticCorrect++
			}
			if row.RawProtocolCompliant {
				out.RawProtocolCompliant++
			}
			if row.RawStrictCorrect {
				out.RawStrictCorrect++
			}
			if row.CandidateVerified {
				out.CandidateVerified++
			}
			if f.id == "plan" {
				plan = row.RawText
			}
			out.Rows = append(out.Rows, row)
		}
	}
	return out, nil
}
