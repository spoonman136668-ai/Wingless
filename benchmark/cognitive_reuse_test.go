package benchmark

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

func TestCognitiveReuseFixture(t *testing.T) {
	report, err := RunCognitiveReuseFixture(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if !report.SemanticCorrect || !report.StrictCorrect || !report.ProtocolCorrect {
		t.Fatalf("correctness regression: %+v", report)
	}
	if report.Control.ModelCalls != 10 {
		t.Fatalf("control model calls=%d, want 10", report.Control.ModelCalls)
	}
	if report.Experiment.ModelCalls != 1 || report.FirstSolutionModelCalls != 1 || report.ReuseSolutionModelCalls != 0 {
		t.Fatalf("experiment did not reduce repeated inference: %+v", report)
	}
	if report.SkillReuseCount != 9 || report.Experiment.SkillReuses != 9 {
		t.Fatalf("skill reuse mismatch: report=%d cost=%d", report.SkillReuseCount, report.Experiment.SkillReuses)
	}
	if report.Control.CognitivePasses != 0 || report.Experiment.CognitivePasses != 1 {
		t.Fatalf("cognitive pass accounting regressed: control=%d experiment=%d", report.Control.CognitivePasses, report.Experiment.CognitivePasses)
	}
	if report.Control.MemoryRetrievals != 0 || report.Experiment.MemoryRetrievals != 0 {
		t.Fatalf("fixture fabricated memory retrievals: %+v", report)
	}
	if report.Experiment.InputTokens >= report.Control.InputTokens || report.Experiment.OutputTokens >= report.Control.OutputTokens {
		t.Fatalf("token cost did not decrease: %+v", report)
	}
	if report.Control.RAMFreeMinBytes != nil || report.Control.VRAMFreeMinBytes != nil || report.Control.ModelTelemetry != nil {
		t.Fatalf("fixture fabricated unavailable resource telemetry: %+v", report.Control)
	}
	b, err := json.Marshal(report.Control)
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	if !strings.Contains(text, `"ram_free_min_bytes":null`) || !strings.Contains(text, `"vram_free_min_bytes":null`) || !strings.Contains(text, `"model_telemetry":null`) {
		t.Fatalf("unavailable benchmark telemetry was not explicit null: %s", text)
	}
}

func TestCognitiveCostCarriesMeasuredTelemetry(t *testing.T) {
	ramBefore := uint64(900)
	ramPeak := uint64(700)
	vramBefore := uint64(400)
	vramPeak := uint64(250)
	logical := uint64(1_000)
	resident := uint64(600)
	model := resources.NewModelTelemetry()
	model.Capacity = &resources.ModelCapacityTelemetry{
		LogicalModelBytes:  &logical,
		ResidentModelBytes: &resident,
		Provenance: &resources.TelemetryProvenance{
			Origin: resources.TelemetryOriginBackend, ProviderID: "measured-fixture",
		},
	}
	result := inference.Result{
		Telemetry: inference.Telemetry{
			Before: &resources.Metrics{RAMFree: &ramBefore, VRAMFree: &vramBefore},
			Peak:   &resources.Metrics{RAMFree: &ramPeak, VRAMFree: &vramPeak},
			Model:  &model,
		},
	}
	cost := CognitiveCost{}
	accumulateInferenceCost(&cost, result)
	if cost.RAMFreeMinBytes == nil || *cost.RAMFreeMinBytes != ramPeak {
		t.Fatalf("RAM telemetry not propagated: %+v", cost)
	}
	if cost.VRAMFreeMinBytes == nil || *cost.VRAMFreeMinBytes != vramPeak {
		t.Fatalf("VRAM telemetry not propagated: %+v", cost)
	}
	if cost.ModelTelemetry == nil || cost.ModelTelemetry.Capacity == nil || cost.ModelTelemetry.Capacity.ResidentModelBytes == nil || *cost.ModelTelemetry.Capacity.ResidentModelBytes != resident {
		t.Fatalf("model telemetry not propagated: %+v", cost.ModelTelemetry)
	}
}
