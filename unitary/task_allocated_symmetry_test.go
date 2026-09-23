package unitary

import (
	"math"
	"testing"
)

func TestUP34InnerSplitDoesNotUseTrueHeldOutPool(t *testing.T) {
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fit, validation := splitTaskAllocationTrainingPool(allTrain)

	if len(fit)+len(validation) != len(allTrain) {
		t.Fatalf(
			"inner split size=%d want=%d",
			len(fit)+len(validation), len(allTrain),
		)
	}
	seen := make(map[int]bool)
	for _, table := range fit {
		seen[memoryTableIndex(table)] = true
	}
	for _, table := range validation {
		index := memoryTableIndex(table)
		if seen[index] {
			t.Fatalf("table index=%d appears in fit and validation", index)
		}
		seen[index] = true
	}
	for _, table := range trueHeld {
		index := memoryTableIndex(table)
		if seen[index] {
			t.Fatalf("true held-out table index=%d leaked into allocator", index)
		}
	}
}

func TestUP34CandidateCapacities(t *testing.T) {
	want := map[string]int{
		"capacity6_singletons": 6,
		"capacity8_2_1_1_1_1": 8,
		"capacity10_2_2_1_1": 10,
		"capacity12_3_1_1_1": 12,
		"capacity12_2_2_2": 12,
		"capacity18_4_1_1": 18,
		"capacity18_3_3": 18,
		"capacity36_6": 36,
	}
	for _, spec := range taskAllocationSpecs() {
		got := commutantCapacityPartition(spec.partition)
		if got != want[spec.name] {
			t.Fatalf("%s capacity=%d want=%d", spec.name, got, want[spec.name])
		}
	}
}

func TestUP34AllocationScorePaysForCapacity(t *testing.T) {
	high := taskAllocationScore(0.99, 0.99, 36)
	low := taskAllocationScore(0.99, 0.99, 12)
	if !(low > high) {
		t.Fatalf("same-performance low-capacity score=%g want>%g", low, high)
	}
	if math.Abs((low-high)-taskAllocationCapacityPrice*(24.0/36.0)) > 1e-12 {
		t.Fatalf("capacity price delta=%g", low-high)
	}
}
