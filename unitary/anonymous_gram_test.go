package unitary

import (
	"math"
	"testing"
)

func TestUP19AnonymousFeatureDimensions(t *testing.T) {
	memory, err := encodeMemory(memoryTable{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	mixed, err := blindMixedComposite(memory)
	if err != nil {
		t.Fatal(err)
	}
	linear, err := anonymousGramFeatures(mixed)
	if err != nil {
		t.Fatal(err)
	}
	if len(linear) != anonymousGramLinearDim {
		t.Fatalf("linear dimension=%d want=%d", len(linear), anonymousGramLinearDim)
	}
	quadratic, err := anonymousQuadraticFeatures(mixed)
	if err != nil {
		t.Fatal(err)
	}
	if len(quadratic) != anonymousGramQuadraticDim {
		t.Fatalf("quadratic dimension=%d want=%d", len(quadratic), anonymousGramQuadraticDim)
	}
	for _, value := range quadratic {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			t.Fatal("quadratic feature is non-finite")
		}
	}
}

func TestUP19AnonymousGramIgnoresCommonGlobalPhase(t *testing.T) {
	memory, err := encodeMemory(memoryTable{3, 1, 0, 2})
	if err != nil {
		t.Fatal(err)
	}
	mixed, err := blindMixedComposite(memory)
	if err != nil {
		t.Fatal(err)
	}
	want, err := anonymousQuadraticFeatures(mixed)
	if err != nil {
		t.Fatal(err)
	}
	got, err := anonymousQuadraticFeatures(
		rotateGlobalPhase(mixed, 1.234),
	)
	if err != nil {
		t.Fatal(err)
	}
	var maximum float64
	for i := range want {
		delta := math.Abs(want[i] - got[i])
		if delta > maximum {
			maximum = delta
		}
	}
	if maximum > 1e-12 {
		t.Fatalf("common-phase feature drift=%g want<=1e-12", maximum)
	}
}

func TestUP19AnonymousGramUnitaryDepthInvariant(t *testing.T) {
	drift, err := anonymousFeatureDrift(applyStressUnitary, 1024)
	if err != nil {
		t.Fatal(err)
	}
	if drift > 1e-12 {
		t.Fatalf("unitary anonymous feature drift=%g want<=1e-12", drift)
	}
}
