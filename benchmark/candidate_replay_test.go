package benchmark

import (
	"context"
	"github.com/spoonman136668-ai/Wingless/inference"
	"testing"
)

func qwenReplayFixtureReport() LocalReport {
	outputs := map[string]string{
		"json-arithmetic":     "```json\n{\"sum\":5}\n```",
		"go-add-shape":        "```go\npackage candidate\n\nfunc Add(a, b int) int {\n    return a + b\n}\n```",
		"code-repair":         "package candidate\nfunc Add(a,b int) int {return a+b}",
		"bug-diagnosis":       "```json\n{\n  \"cause\": \"off_by_one\"\n}\n```",
		"code-review":         "{\"safe\":false,\"issue\":\"unchecked_index\"}",
		"multi-file":          "```json\n{\n  \"value\": 6\n}\n```",
		"plan":                "{\"operation\":\"add\"}",
		"plan-implementation": "```go\npackage candidate\n\nfunc Add(a, b int) int {\n    return a + b\n}\n```",
		"tool-proposal":       "```json\n{\n  \"tool\": \"read_file\",\n  \"path\": \"src/add.go\"\n}\n```",
	}
	strict := map[string]bool{"code-repair": true, "code-review": true, "plan": true}
	report := LocalReport{Version: 3, Acceptance: "external_required"}
	for _, fixture := range codingFixtures() {
		report.Rows = append(report.Rows, LocalRow{
			Task:              fixture.id,
			Repetition:        0,
			SemanticCorrect:   true,
			ProtocolCompliant: strict[fixture.id],
			StrictCorrect:     strict[fixture.id],
			Correct:           strict[fixture.id],
			Result:            inference.Result{Text: outputs[fixture.id], Status: "completed"},
		})
	}
	return report
}

func TestReplayVerificationCandidatesQwenShape(t *testing.T) {
	replay, err := ReplayVerificationCandidates(context.Background(), qwenReplayFixtureReport())
	if err != nil {
		t.Fatal(err)
	}
	if replay.RawStrictCorrect != 3 || replay.RawSemanticCorrect != 9 || replay.RawProtocolCompliant != 3 || replay.CandidateVerified != 9 || replay.Acceptance != "external_required" || replay.GeneratedCodeExecuted {
		t.Fatal(replay)
	}
	normalized := 0
	for _, row := range replay.Rows {
		if !row.CandidateVerified || row.Verification != "fixture_pass" {
			t.Fatal(row)
		}
		if row.CandidateNormalized {
			normalized++
			if row.RawModelText == "" || row.RawModelText == row.CandidateText {
				t.Fatal("normalized row did not preserve distinct raw model text", row)
			}
		}
	}
	if normalized != 6 {
		t.Fatalf("normalized=%d want=6", normalized)
	}
}

func TestReplayVerificationCandidatesRejectsChangedFixtureOrder(t *testing.T) {
	report := qwenReplayFixtureReport()
	report.Rows[0], report.Rows[1] = report.Rows[1], report.Rows[0]
	if _, err := ReplayVerificationCandidates(context.Background(), report); err == nil {
		t.Fatal("expected fixed fixture order rejection")
	}
}
