package unitary

import (
	"math"
	"testing"
)

func TestUP35StartsFromSingletonCapacity(t *testing.T) {
	groups := initialSymmetryGroups()
	if got := symmetryGroupCapacity(groups); got != 6 {
		t.Fatalf("singleton capacity=%d want=6", got)
	}
	offsets, err := offsetsForSymmetryGroups(groups, multiplicityDoseRMS)
	if err != nil {
		t.Fatal(err)
	}
	if got := doseMaxMultiplicity(offsets); got != 1 {
		t.Fatalf("singleton max multiplicity=%d want=1", got)
	}
	if delta := math.Abs(doseOffsetRMS(offsets)-multiplicityDoseRMS); delta > 1e-12 {
		t.Fatalf("singleton RMS delta=%g", delta)
	}
}

func TestUP35PairMergeRaisesCapacityByTwiceProduct(t *testing.T) {
	groups := initialSymmetryGroups()
	before := symmetryGroupCapacity(groups)
	merged, err := mergeSymmetryGroups(groups, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	after := symmetryGroupCapacity(merged)
	if after-before != 2 {
		t.Fatalf("first merge capacity gain=%d want=2", after-before)
	}
	merged2, err := mergeSymmetryGroups(merged, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got := symmetryGroupCapacity(merged2); got != 12 {
		t.Fatalf("second merge capacity=%d want=12", got)
	}
}

func TestUP35InnerSplitRemainsBalancedAndLeakFree(t *testing.T) {
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fit, validation := splitTaskAllocationTrainingPool(allTrain)
	if minimumTableMarginalCount(fit) < 12 {
		t.Fatalf("fit marginal minimum=%d want>=12", minimumTableMarginalCount(fit))
	}
	if minimumTableMarginalCount(validation) < 4 {
		t.Fatalf("validation marginal minimum=%d want>=4", minimumTableMarginalCount(validation))
	}
	seen := make(map[int]bool)
	for _, table := range fit {
		seen[memoryTableIndex(table)] = true
	}
	for _, table := range validation {
		seen[memoryTableIndex(table)] = true
	}
	for _, table := range trueHeld {
		if seen[memoryTableIndex(table)] {
			t.Fatal("true held-out table leaked into induction pool")
		}
	}
}

func TestUP35FrozenGainRule(t *testing.T) {
	current := 0.75
	if !(0.756-current >= greedySymmetryMinGain) {
		t.Fatal("gain above frozen threshold should be accepted")
	}
	if 0.754-current >= greedySymmetryMinGain {
		t.Fatal("gain below frozen threshold should not be accepted")
	}
}
