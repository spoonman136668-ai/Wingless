package unitary

import (
	"math"
	"testing"
)

func TestUP0InterferenceRouting(t *testing.T) {
	result, err := RunUP0()
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != ProbeSchema {
		t.Fatalf("schema=%q", result.Schema)
	}
	if result.BaselineAccuracy != 0.5 {
		t.Fatalf("baseline accuracy=%v want=0.5", result.BaselineAccuracy)
	}
	if result.UnitaryAccuracy != 1 {
		t.Fatalf("unitary accuracy=%v want=1", result.UnitaryAccuracy)
	}
	if result.MaxNormDrift > 1e-12 {
		t.Fatalf("norm drift=%g", result.MaxNormDrift)
	}
	if result.MaxRoundTripError > 1e-12 {
		t.Fatalf("round-trip error=%g", result.MaxRoundTripError)
	}
}

func TestPropagatePreservesNormAndInverts(t *testing.T) {
	initial, err := Normalize(State{
		complex(1.25, -0.5),
		complex(-0.75, 0.25),
		complex(0.5, 1.5),
		complex(-1, -0.25),
	})
	if err != nil {
		t.Fatal(err)
	}
	program := []Coupling{
		{A: 0, B: 1, Theta: 0.37, Phi: 0.19},
		{A: 2, B: 3, Theta: -0.81, Phi: -0.43},
		{A: 1, B: 2, Theta: 1.17, Phi: 0.72},
	}

	before, err := NormSquared(initial)
	if err != nil {
		t.Fatal(err)
	}
	evolved, err := Propagate(initial, program)
	if err != nil {
		t.Fatal(err)
	}
	after, err := NormSquared(evolved)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(before-after) > 1e-12 {
		t.Fatalf("norm changed: before=%0.16g after=%0.16g", before, after)
	}

	recovered, err := Propagate(evolved, Invert(program))
	if err != nil {
		t.Fatal(err)
	}
	distance, err := L2Distance(initial, recovered)
	if err != nil {
		t.Fatal(err)
	}
	if distance > 1e-12 {
		t.Fatalf("inverse failed: distance=%g", distance)
	}
}

func TestPropagateFailsClosedOnInvalidCoupling(t *testing.T) {
	state := State{1, 0}
	for _, program := range [][]Coupling{
		{{A: 0, B: 0, Theta: 0.1}},
		{{A: -1, B: 1, Theta: 0.1}},
		{{A: 0, B: 2, Theta: 0.1}},
		{{A: 0, B: 1, Theta: math.NaN()}},
		{{A: 0, B: 1, Theta: 0.1, Phi: math.Inf(1)}},
	} {
		if _, err := Propagate(state, program); err == nil {
			t.Fatalf("expected invalid coupling to fail: %#v", program)
		}
	}
}

func TestRelativePhaseAffectsMeasurementOnlyAfterPropagation(t *testing.T) {
	r := 1 / math.Sqrt2
	inPhase := State{complex(r, 0), complex(r, 0)}
	opposed := State{complex(r, 0), complex(-r, 0)}

	beforeA, err := Probabilities(inPhase)
	if err != nil {
		t.Fatal(err)
	}
	beforeB, err := Probabilities(opposed)
	if err != nil {
		t.Fatal(err)
	}
	for i := range beforeA {
		if math.Abs(beforeA[i]-beforeB[i]) > 1e-15 {
			t.Fatalf("pre-propagation probabilities differ: %v vs %v", beforeA, beforeB)
		}
	}

	program := []Coupling{{A: 0, B: 1, Theta: math.Pi / 4}}
	afterA, err := Propagate(inPhase, program)
	if err != nil {
		t.Fatal(err)
	}
	afterB, err := Propagate(opposed, program)
	if err != nil {
		t.Fatal(err)
	}
	routeA, err := ArgMax(afterA)
	if err != nil {
		t.Fatal(err)
	}
	routeB, err := ArgMax(afterB)
	if err != nil {
		t.Fatal(err)
	}
	if routeA == routeB {
		t.Fatalf("relative phase did not separate routes: %d == %d", routeA, routeB)
	}
}
