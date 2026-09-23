package unitary

import (
	"math"
	"testing"
)

func TestUP33MatchedCapacityPairs(t *testing.T) {
	if got := commutantCapacityPartition([]int{4, 1, 1}); got != 18 {
		t.Fatalf("capacity 4+1+1=%d want=18", got)
	}
	if got := commutantCapacityPartition([]int{3, 3}); got != 18 {
		t.Fatalf("capacity 3+3=%d want=18", got)
	}
	if got := commutantCapacityPartition([]int{3, 1, 1, 1}); got != 12 {
		t.Fatalf("capacity 3+1+1+1=%d want=12", got)
	}
	if got := commutantCapacityPartition([]int{2, 2, 2}); got != 12 {
		t.Fatalf("capacity 2+2+2=%d want=12", got)
	}
}

func TestUP33PartitionOffsetsMatchRMSAndMultiplicity(t *testing.T) {
	cases := []struct {
		sizes []int
		max   int
	}{
		{[]int{4, 1, 1}, 4},
		{[]int{3, 3}, 3},
		{[]int{3, 1, 1, 1}, 3},
		{[]int{2, 2, 2}, 2},
	}
	for _, tc := range cases {
		offsets, err := commutantPartitionOffsets(
			tc.sizes, multiplicityDoseRMS,
		)
		if err != nil {
			t.Fatal(err)
		}
		if got := doseMaxMultiplicity(offsets); got != tc.max {
			t.Fatalf(
				"partition=%v max multiplicity=%d want=%d",
				tc.sizes, got, tc.max,
			)
		}
		if delta := math.Abs(
			doseOffsetRMS(offsets)-multiplicityDoseRMS,
		); delta > 1e-12 {
			t.Fatalf("partition=%v RMS delta=%g", tc.sizes, delta)
		}
		var sum float64
		for _, value := range offsets {
			sum += value
		}
		if math.Abs(sum) > 1e-12 {
			t.Fatalf("partition=%v offset sum=%g", tc.sizes, sum)
		}
	}
}

func TestUP33PermutationsAreBijective(t *testing.T) {
	for _, permutation := range commutantCapacityPermutations() {
		seen := make(map[int]bool)
		for _, index := range permutation.index {
			if index < 0 || index >= compositeChannels {
				t.Fatalf("%s index out of range=%d", permutation.name, index)
			}
			if seen[index] {
				t.Fatalf("%s duplicate index=%d", permutation.name, index)
			}
			seen[index] = true
		}
	}
}
