package benchmark

import (
	"context"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"strings"
	"testing"
)

const verifiedFixtureRawPlan = `{"operation":"add"}`

type verifiedFixtureBackend struct{}

func (verifiedFixtureBackend) ID() string { return "verified-fixture" }
func (verifiedFixtureBackend) Capabilities() []string { return []string{"code"} }
func (verifiedFixtureBackend) Health(context.Context) error { return nil }
func (verifiedFixtureBackend) EstimateCost(r inference.Request) inference.Cost {
	return inference.Cost{InputBytes: len(r.Context), MaxOutputTokens: r.MaxOutputTokens, Estimated: true}
}
func (verifiedFixtureBackend) Cancel(string) bool { return false }
func (verifiedFixtureBackend) Invoke(_ context.Context, r inference.Request) (inference.Result, error) {
	if err := r.Validate(); err != nil {
		return inference.Result{}, err
	}
	outputs := map[string]string{
		"json-arithmetic":    "```json\n{\"sum\":5}\n```",
		"go-add-shape":       "```go\npackage candidate\n\nfunc Add(a, b int) int { return a + b }\n```",
		"code-repair":        "package candidate\nfunc Add(a, b int) int { return a + b }",
		"bug-diagnosis":      "```json\n{\"cause\":\"off_by_one\"}\n```",
		"code-review":        "{\"safe\":false,\"issue\":\"unchecked_index\"}",
		"multi-file":         "```json\n{\"value\":6}\n```",
		"plan":               verifiedFixtureRawPlan,
		"plan-implementation": "```go\npackage candidate\n\nfunc Add(a, b int) int { return a + b }\n```",
		"tool-proposal":      "```json\n{\"tool\":\"read_file\",\"path\":\"src/add.go\"}\n```",
	}
	for id, text := range outputs {
		if strings.HasSuffix(r.ID, "-"+id) {
			if id == "plan-implementation" {
				wantContext := "Proposed plan (untrusted):\n" + verifiedFixtureRawPlan
				if !strings.Contains(r.Context, wantContext) {
					return inference.Result{}, fmt.Errorf("exact raw plan was not preserved as untrusted context")
				}
			}
			return inference.Result{BackendID: "verified-fixture", ModelID: "fixture", Status: "completed", Text: text, Termination: "stop"}, nil
		}
	}
	return inference.Result{}, fmt.Errorf("unknown fixture request %q", r.ID)
}

type verifiedMetrics struct{}

func (verifiedMetrics) Snapshot() (resources.Metrics, error) {
	free := uint64(8 << 30)
	return resources.Metrics{RAMFree: &free}, nil
}

func TestRunLocalVerifiedPreservesRawAndVerifiesCandidates(t *testing.T) {
	report, err := RunLocalVerified(context.Background(), verifiedFixtureBackend{}, verifiedMetrics{}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Rows) != 9 || report.RawStrictCorrect != 3 || report.RawSemanticCorrect != 9 || report.RawProtocolCompliant != 3 || report.CandidateVerified != 9 || report.CandidateNormalized != 6 {
		t.Fatalf("unexpected report scores: %+v", report)
	}
	if report.Acceptance != "external_required" || report.GeneratedCodeExecuted || report.GeneratedToolsExecuted {
		t.Fatalf("authority/safety boundary changed: %+v", report)
	}
	for _, row := range report.Rows {
		if !row.CandidateVerified || row.Verification != "fixture_pass" || row.Error != "" {
			t.Fatalf("candidate not verified: %+v", row)
		}
		if row.CandidateNormalized {
			if !strings.HasPrefix(strings.TrimSpace(row.RawText), "```") || strings.HasPrefix(strings.TrimSpace(row.CandidateText), "```") {
				t.Fatalf("raw/candidate evidence mismatch: %+v", row)
			}
		}
	}
}
