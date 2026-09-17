package nr1

import (
	"strings"
	"testing"
)

func qualifiedTrace(tokens int) []TraceEvent {
	out := make([]TraceEvent, 0, tokens*QualifiedQwen3CoderLayers)
	for token := 0; token < tokens; token++ {
		phase := "prefill"
		if token > 0 {
			phase = "decode"
		}
		for layer := 0; layer < QualifiedQwen3CoderLayers; layer++ {
			experts := make([]int, QualifiedQwen3CoderExpertsPerToken)
			for i := range experts {
				experts[i] = (layer*QualifiedQwen3CoderExpertsPerToken + i) % QualifiedQwen3CoderExpertsPerLayer
			}
			out = append(out, TraceEvent{
				Schema: TraceSchema, SessionID: "s1", WorkloadID: "w1", TaskFamily: "coding", Phase: phase,
				TokenIndex: int64(token), Layer: layer, Experts: experts,
			})
		}
	}
	return out
}

func TestValidateQualifiedQwen3CoderTrace(t *testing.T) {
	if err := ValidateQualifiedQwen3CoderTrace(qualifiedTrace(3)); err != nil {
		t.Fatal(err)
	}
}

func TestValidateQualifiedQwen3CoderTraceRejectsIncompleteLayer(t *testing.T) {
	events := qualifiedTrace(2)
	events = append(events[:17], events[18:]...)
	if err := ValidateQualifiedQwen3CoderTrace(events); err == nil || !strings.Contains(err.Error(), "INCOMPLETE_LAYERS") {
		t.Fatalf("expected incomplete-layer rejection, got %v", err)
	}
}

func TestValidateQualifiedQwen3CoderTraceRejectsDuplicateLayer(t *testing.T) {
	events := qualifiedTrace(1)
	events = append(events, events[0])
	if err := ValidateQualifiedQwen3CoderTrace(events); err == nil || !strings.Contains(err.Error(), "DUPLICATE_LAYER") {
		t.Fatalf("expected duplicate-layer rejection, got %v", err)
	}
}

func TestValidateQualifiedQwen3CoderTraceRejectsWrongTopK(t *testing.T) {
	events := qualifiedTrace(1)
	events[0].Experts = events[0].Experts[:7]
	if err := ValidateQualifiedQwen3CoderTrace(events); err == nil || !strings.Contains(err.Error(), "TOPK_INVALID") {
		t.Fatalf("expected top-k rejection, got %v", err)
	}
}

func TestValidateQualifiedQwen3CoderTraceRejectsExpertOutOfRange(t *testing.T) {
	events := qualifiedTrace(1)
	events[0].Experts[0] = QualifiedQwen3CoderExpertsPerLayer
	if err := ValidateQualifiedQwen3CoderTrace(events); err == nil || !strings.Contains(err.Error(), "EXPERT_ID_INVALID") {
		t.Fatalf("expected expert-id rejection, got %v", err)
	}
}

func TestValidateQualifiedQwen3CoderTraceRejectsTokenGap(t *testing.T) {
	events := qualifiedTrace(3)
	filtered := events[:0]
	for _, event := range events {
		if event.TokenIndex != 1 {
			filtered = append(filtered, event)
		}
	}
	if err := ValidateQualifiedQwen3CoderTrace(filtered); err == nil || !strings.Contains(err.Error(), "TOKEN_SEQUENCE_INVALID") {
		t.Fatalf("expected token-sequence rejection, got %v", err)
	}
}
