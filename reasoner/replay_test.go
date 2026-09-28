package reasoner

import (
	"encoding/json"
	"strings"
	"testing"
)

func fixtureForTest() ReplayFixture {
	return ReplayFixture{
		Schema:             ReplayFixtureSchema,
		FixtureID:          "f1",
		Project:            "Wingless",
		HistoricalFrontier: "historical-ref",
		Question:           "What should run next?",
		Evidence:           []string{"a negative result isolated a representation ambiguity"},
		Constraints:        []string{"do not tune thresholds after results"},
		Candidates: []ReplayCandidate{
			{ID: "A", Proposal: "replay frozen states"},
			{ID: "B", Proposal: "change threshold"},
		},
		ExpectedCandidate: "A",
		ExpectedWhy:       "frozen replay preserves causal discipline",
	}
}

func TestReplayPromptWithholdsExpectedAnswer(t *testing.T) {
	f := fixtureForTest()
	prompt, err := f.Prompt()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(prompt, f.ExpectedWhy) {
		t.Fatal("prompt leaked expected rationale")
	}
	if strings.Contains(prompt, "ExpectedCandidate") || strings.Contains(prompt, "expected_candidate") {
		t.Fatal("prompt leaked expected candidate field")
	}
}

func TestReplayRequestPinsStructuredSchemaAndIdentity(t *testing.T) {
	f := fixtureForTest()
	r1, err := f.BuildRequest(1)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := f.BuildRequest(1)
	if err != nil {
		t.Fatal(err)
	}
	if r1.FrontierSHA256 != r2.FrontierSHA256 || r1.QualificationContractSHA256 != r2.QualificationContractSHA256 {
		t.Fatal("replay identity is not deterministic")
	}
	if r1.ResponseJSONSchema == nil || r1.ResponseJSONSchema.Name != "historical_replay_decision" {
		t.Fatal("structured output schema missing")
	}
	raw, err := json.Marshal(r1.ResponseJSONSchema.Schema)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "candidate_id") {
		t.Fatal("candidate schema missing")
	}
}

func TestReplayScore(t *testing.T) {
	f := fixtureForTest()
	ok, why := f.Score(ReplayDecision{CandidateID: "A", Rationale: "bounded", ConstraintViolation: false})
	if !ok || why != "pass" {
		t.Fatalf("expected pass, got %v %s", ok, why)
	}
	ok, why = f.Score(ReplayDecision{CandidateID: "B", Rationale: "bad", ConstraintViolation: false})
	if ok || why != "historical_choice_mismatch" {
		t.Fatalf("expected mismatch, got %v %s", ok, why)
	}
	ok, why = f.Score(ReplayDecision{CandidateID: "A", Rationale: "bad", ConstraintViolation: true})
	if ok || why != "constraint_violation" {
		t.Fatalf("expected violation, got %v %s", ok, why)
	}
}
