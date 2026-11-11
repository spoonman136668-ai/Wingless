package unitary

import (
	"math"
	"testing"
)

func TestUP39HarmonicMeanPenalizesWeakestSignal(t *testing.T) {
	balanced := harmonicMean3(0.8, 0.8, 0.8)
	weakRelation := harmonicMean3(0.95, 0.95, 0.4)
	if !(balanced > weakRelation) {
		t.Fatalf("balanced=%g want>weakRelation=%g", balanced, weakRelation)
	}
}

func TestUP39HarmonicMeanEqualInputs(t *testing.T) {
	got := harmonicMean3(0.75, 0.75, 0.75)
	if math.Abs(got-0.75) > 1e-12 {
		t.Fatalf("harmonic equal=%g want=0.75", got)
	}
}

func TestUP39UsesSameDirectOptimizerGeometry(t *testing.T) {
	groups := initialSymmetryGroupsJSON()
	rows := make([][]float64, 0, 7)
	for step := 1; step <= 7; step++ {
		rows = append(rows, deterministicSPSADirection(step, groups))
	}
	if got := directDirectionRank(rows); got != compositeChannels-1 {
		t.Fatalf("rich direct tangent rank=%d want=%d", got, compositeChannels-1)
	}
}

func TestUP39NoTrueHeldOutInInnerSplit(t *testing.T) {
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fit, validation := splitTaskAllocationTrainingPool(allTrain)
	seen := make(map[int]bool)
	for _, table := range fit {
		seen[memoryTableIndex(table)] = true
	}
	for _, table := range validation {
		seen[memoryTableIndex(table)] = true
	}
	for _, table := range trueHeld {
		if seen[memoryTableIndex(table)] {
			t.Fatal("true held-out table leaked into rich objective")
		}
	}
}
