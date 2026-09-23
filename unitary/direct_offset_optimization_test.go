package unitary

import (
	"math"
	"testing"
)

func TestUP38SoftCapacityMatchesExactLimits(t *testing.T) {
	full, err := directSoftCapacity(
		[]float64{0, 0, 0, 0, 0, 0},
	)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(full-36) > 1e-12 {
		t.Fatalf("full soft capacity=%g want=36", full)
	}

	singletons := continuousFusionInitialOffsets()
	soft, err := directSoftCapacity(singletons)
	if err != nil {
		t.Fatal(err)
	}
	if !(soft >= 6 && soft < 8) {
		t.Fatalf("singleton soft capacity=%g want in [6,8)", soft)
	}
}

func TestUP38DirectionPreservesFusedGroups(t *testing.T) {
	groups := [][]int{{0, 1, 2}, {3}, {4}, {5}}
	direction := deterministicSPSADirection(3, groups)
	if direction[0] != direction[1] ||
		direction[1] != direction[2] {
		t.Fatal("SPSA direction splits fused group")
	}
	var mean float64
	for _, value := range direction {
		mean += value
	}
	mean /= float64(len(direction))
	if math.Abs(mean) > 1e-12 {
		t.Fatalf("direction mean=%g want=0", mean)
	}
}

func TestUP38DirectObjectiveChargesSymmetry(t *testing.T) {
	var arm MultiplicityDoseArm
	arm.Static.MeanPhaseCosine = 0.9

	independent := continuousFusionInitialOffsets()
	full := []float64{0, 0, 0, 0, 0, 0}

	independentObjective, _, _, err :=
		directObjective(arm, independent)
	if err != nil {
		t.Fatal(err)
	}
	fullObjective, _, _, err :=
		directObjective(arm, full)
	if err != nil {
		t.Fatal(err)
	}
	if !(independentObjective > fullObjective) {
		t.Fatalf(
			"resource objective independent=%g want>full=%g",
			independentObjective, fullObjective,
		)
	}
}

func TestUP38GradientUpdateIsBounded(t *testing.T) {
	offsets := continuousFusionInitialOffsets()
	gradient := []float64{100, -100, 50, -50, 25, -25}
	union := newFusionUnion()
	next, update, err := applyDirectGradient(
		offsets, gradient, &union,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(next) != compositeChannels {
		t.Fatalf("next dimension=%d", len(next))
	}
	for _, value := range update {
		if math.Abs(value) > directOffsetMaxUpdate+1e-15 {
			t.Fatalf("update=%g exceeds max=%g", value, directOffsetMaxUpdate)
		}
	}
}
