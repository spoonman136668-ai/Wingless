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


func directDirectionRank(rows [][]float64) int {
	if len(rows) == 0 {
		return 0
	}
	matrix := make([][]float64, len(rows))
	for i := range rows {
		matrix[i] = append([]float64(nil), rows[i]...)
	}
	rank := 0
	columnCount := len(matrix[0])
	for column := 0; column < columnCount && rank < len(matrix); column++ {
		pivot := -1
		for row := rank; row < len(matrix); row++ {
			if math.Abs(matrix[row][column]) > 1e-12 {
				pivot = row
				break
			}
		}
		if pivot < 0 {
			continue
		}
		matrix[rank], matrix[pivot] = matrix[pivot], matrix[rank]
		pivotValue := matrix[rank][column]
		for c := column; c < columnCount; c++ {
			matrix[rank][c] /= pivotValue
		}
		for row := 0; row < len(matrix); row++ {
			if row == rank {
				continue
			}
			factor := matrix[row][column]
			if math.Abs(factor) <= 1e-12 {
				continue
			}
			for c := column; c < columnCount; c++ {
				matrix[row][c] -= factor * matrix[rank][c]
			}
		}
		rank++
	}
	return rank
}

func TestUP38DeterministicDirectionsSpanZeroMeanTangent(t *testing.T) {
	groups := initialSymmetryGroupsJSON()
	rows := make([][]float64, 0, 7)
	for step := 1; step <= 7; step++ {
		direction := deterministicSPSADirection(step, groups)
		var mean float64
		for _, value := range direction {
			mean += value
		}
		mean /= float64(len(direction))
		if math.Abs(mean) > 1e-12 {
			t.Fatalf("step=%d direction mean=%g want=0", step, mean)
		}
		rows = append(rows, direction)
	}
	if got := directDirectionRank(rows); got != compositeChannels-1 {
		t.Fatalf("direction tangent rank=%d want=%d", got, compositeChannels-1)
	}
}
