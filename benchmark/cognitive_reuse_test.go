package benchmark

import (
	"context"
	"testing"
)

func TestCognitiveReuseFixture(t *testing.T) {
	report, err := RunCognitiveReuseFixture(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if !report.StrictCorrect || !report.ProtocolCorrect {
		t.Fatalf("correctness regression: %+v", report)
	}
	if report.Control.ModelCalls != 10 {
		t.Fatalf("control model calls=%d, want 10", report.Control.ModelCalls)
	}
	if report.Experiment.ModelCalls != 1 || report.FirstSolutionModelCalls != 1 || report.ReuseSolutionModelCalls != 0 {
		t.Fatalf("experiment did not reduce repeated inference: %+v", report)
	}
	if report.SkillReuseCount != 9 {
		t.Fatalf("skill reuse=%d, want 9", report.SkillReuseCount)
	}
	if report.Experiment.InputTokens >= report.Control.InputTokens || report.Experiment.OutputTokens >= report.Control.OutputTokens {
		t.Fatalf("token cost did not decrease: %+v", report)
	}
}
