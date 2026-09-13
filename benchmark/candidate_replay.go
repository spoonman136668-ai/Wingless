package benchmark

import (
	"context"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"time"
)

type VerificationReplayRow struct {
	Task                 string `json:"task"`
	RawStrictCorrect     bool   `json:"raw_strict_correct"`
	RawSemanticCorrect   bool   `json:"raw_semantic_correct"`
	RawProtocolCompliant bool   `json:"raw_protocol_compliant"`
	CandidateVerified    bool   `json:"candidate_verified"`
	CandidateNormalized  bool   `json:"candidate_normalized"`
	RawModelText         string `json:"raw_model_text"`
	CandidateText        string `json:"candidate_text"`
	Verification         string `json:"verification"`
	Error                string `json:"error,omitempty"`
}

type VerificationReplay struct {
	Schema                string                  `json:"schema"`
	RawStrictCorrect      int                     `json:"raw_strict_correct"`
	RawSemanticCorrect    int                     `json:"raw_semantic_correct"`
	RawProtocolCompliant  int                     `json:"raw_protocol_compliant"`
	CandidateVerified     int                     `json:"candidate_verified"`
	Acceptance            string                  `json:"acceptance"`
	GeneratedCodeExecuted bool                    `json:"generated_code_executed"`
	Rows                  []VerificationReplayRow `json:"rows"`
}

type replayMetrics struct{}

func (replayMetrics) Snapshot() (resources.Metrics, error) { return resources.Metrics{}, nil }

type fixtureVerifier struct{ fixture localFixture }

func (v fixtureVerifier) Verify(_ context.Context, r inference.Result) error {
	if !v.fixture.protocol(r.Text) {
		return fmt.Errorf("fixture protocol mismatch")
	}
	if !v.fixture.check(r.Text) {
		return fmt.Errorf("fixture semantic mismatch")
	}
	return nil
}

// ReplayVerificationCandidates replays preserved raw model text through the
// production broker verification seam. It performs no model inference and
// executes no generated code or tools. Raw benchmark scores remain unchanged.
func ReplayVerificationCandidates(ctx context.Context, report LocalReport) (VerificationReplay, error) {
	out := VerificationReplay{Schema: "wingless.verified-candidate-replay.v1", Acceptance: "external_required", GeneratedCodeExecuted: false}
	if report.Version != 3 || report.Acceptance != "external_required" {
		return out, fmt.Errorf("unsupported source report")
	}
	fixtures := codingFixtures()
	if len(report.Rows) != len(fixtures) {
		return out, fmt.Errorf("source report must contain exactly %d rows", len(fixtures))
	}
	for i, fixture := range fixtures {
		row := report.Rows[i]
		if row.Task != fixture.id || row.Repetition != 0 || row.Error != "" {
			return out, fmt.Errorf("source row %d does not match fixed fixture contract", i)
		}
		if row.StrictCorrect {
			out.RawStrictCorrect++
		}
		if row.SemanticCorrect {
			out.RawSemanticCorrect++
		}
		if row.ProtocolCompliant {
			out.RawProtocolCompliant++
		}

		backendID := fmt.Sprintf("candidate-replay-%d", i)
		reg := &broker.Registry{}
		if err := reg.Register(broker.Entry{Backend: &inference.Mock{Name: backendID, Features: []string{"code"}, Reply: row.Result.Text}, Class: "mock", Tier: "fast"}); err != nil {
			return out, err
		}
		runner := broker.Runner{Registry: reg, Metrics: replayMetrics{}, Verifier: fixtureVerifier{fixture: fixture}}
		req := inference.Request{
			ID:               fmt.Sprintf("verified-candidate-replay-%d", i),
			ParentWorkID:     fixture.id,
			Role:             "code",
			Context:          fixture.prompt,
			MaxContextBytes:  4096,
			MaxOutputTokens:  256,
			Deadline:         time.Now().Add(5 * time.Second),
			Workspace:        "no-workspace-access",
			Capabilities:     []string{"code"},
		}
		policy := broker.Policy{AllowedBackends: []string{backendID}, MaxRepairs: 0, Reason: "routine"}
		outcome := runner.Run(ctx, req, policy)
		if len(outcome.Events) != 1 || outcome.Acceptance != "external_required" {
			return out, fmt.Errorf("unexpected broker replay shape for %s", fixture.id)
		}
		event := outcome.Events[0]
		rawText := row.Result.Text
		normalized := event.RawModelText != ""
		if normalized && event.RawModelText != rawText {
			return out, fmt.Errorf("raw model text not preserved for %s", fixture.id)
		}
		if !normalized && event.Result.Text != rawText {
			return out, fmt.Errorf("unverified normalization for %s", fixture.id)
		}
		verified := outcome.Status == "verified_candidate"
		if verified {
			out.CandidateVerified++
		}
		out.Rows = append(out.Rows, VerificationReplayRow{
			Task:                 fixture.id,
			RawStrictCorrect:     row.StrictCorrect,
			RawSemanticCorrect:   row.SemanticCorrect,
			RawProtocolCompliant: row.ProtocolCompliant,
			CandidateVerified:    verified,
			CandidateNormalized:  normalized,
			RawModelText:         rawText,
			CandidateText:        event.Result.Text,
			Verification:         event.Verification,
			Error:                event.Error,
		})
	}
	return out, nil
}
