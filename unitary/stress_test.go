package unitary

import (
	"math"
	"reflect"
	"testing"
)

func TestUP3DepthRetention(t *testing.T) {
	result, err := RunUP3()
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != StressProbeSchema {
		t.Fatalf("schema=%q", result.Schema)
	}
	if result.Dimension != 16 || result.Classes != 4 || result.DistractorDimensions != 12 {
		t.Fatalf("unexpected shape: %#v", result)
	}
	if result.Samples != 32 || result.CouplingsPerBlock != 32 {
		t.Fatalf("unexpected sample/coupling count: samples=%d couplings=%d", result.Samples, result.CouplingsPerBlock)
	}
	if !reflect.DeepEqual(result.Depths, []int{0, 1, 8, 32, 128}) {
		t.Fatalf("depths=%v", result.Depths)
	}

	for _, metric := range result.Unitary.Metrics {
		if metric.Accuracy != 1 {
			t.Fatalf("unitary accuracy at depth %d = %v", metric.Depth, metric.Accuracy)
		}
		if metric.MaxNormDrift > 1e-12 {
			t.Fatalf("unitary norm drift at depth %d = %g", metric.Depth, metric.MaxNormDrift)
		}
		if metric.MaxGramError > 1e-12 {
			t.Fatalf("unitary gram error at depth %d = %g", metric.Depth, metric.MaxGramError)
		}
		if math.Abs(metric.MaxPerturbationGain-1) > 1e-12 {
			t.Fatalf("unitary perturbation gain at depth %d = %g", metric.Depth, metric.MaxPerturbationGain)
		}
	}
	if result.UnitaryMaxRoundTripError > 1e-11 {
		t.Fatalf("unitary round-trip error=%g", result.UnitaryMaxRoundTripError)
	}

	if len(result.NonUnitary.Metrics) != len(result.Unitary.Metrics) {
		t.Fatalf("matched control depth count mismatch")
	}
	for _, metric := range result.NonUnitary.Metrics {
		for label, value := range map[string]float64{
			"accuracy": metric.Accuracy,
			"norm": metric.MaxNormDrift,
			"gram": metric.MaxGramError,
			"perturbation": metric.MaxPerturbationGain,
		} {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				t.Fatalf("non-unitary %s at depth %d is non-finite: %g", label, metric.Depth, value)
			}
		}
	}
}

func TestUP3IsDeterministic(t *testing.T) {
	a, err := RunUP3()
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunUP3()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("UP-3 is nondeterministic:\nA=%#v\nB=%#v", a, b)
	}
}

func TestUP3ControlIsMatchedAtDepthZero(t *testing.T) {
	result, err := RunUP3()
	if err != nil {
		t.Fatal(err)
	}
	u := result.Unitary.Metrics[0]
	n := result.NonUnitary.Metrics[0]
	if u.Depth != 0 || n.Depth != 0 {
		t.Fatalf("first metrics are not depth zero")
	}
	if !reflect.DeepEqual(u, n) {
		t.Fatalf("depth-zero paths differ: unitary=%#v control=%#v", u, n)
	}
}
