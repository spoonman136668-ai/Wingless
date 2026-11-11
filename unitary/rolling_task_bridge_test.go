package unitary

import (
	"math"
	"testing"
)

func TestUP43CheckpointScheduleFrozen(t *testing.T) {
	want := []int{12, 24, 48}
	if len(rollingTaskCheckpointSteps) != len(want) {
		t.Fatalf("checkpoint count=%d", len(rollingTaskCheckpointSteps))
	}
	for index := range want {
		if rollingTaskCheckpointSteps[index] != want[index] {
			t.Fatalf(
				"checkpoint[%d]=%d want=%d",
				index,
				rollingTaskCheckpointSteps[index],
				want[index],
			)
		}
	}
	if rollingTaskMaximumSteps != 48 {
		t.Fatalf("maximum steps=%d want=48", rollingTaskMaximumSteps)
	}
}

func TestUP43ProvisionalGeometryDoesNotMutateOffsets(t *testing.T) {
	offsets := continuousFusionInitialOffsets()
	before := append([]float64(nil), offsets...)
	_, capacity, _, err :=
		reversibleProvisionalGeometry(offsets)
	if err != nil {
		t.Fatal(err)
	}
	if capacity != 6 {
		t.Fatalf("initial provisional capacity=%d want=6", capacity)
	}
	for index := range offsets {
		if offsets[index] != before[index] {
			t.Fatal("provisional geometry mutated offsets")
		}
	}
}

func TestUP43ProvisionalGeometryDetectsNearFusion(t *testing.T) {
	offsets := continuousFusionInitialOffsets()
	offsets[1] = offsets[0] + continuousFusionTolerance/2
	offsets, err := normalizeContinuousOffsets(offsets)
	if err != nil {
		t.Fatal(err)
	}
	groups, capacity, gap, err :=
		reversibleProvisionalGeometry(offsets)
	if err != nil {
		t.Fatal(err)
	}
	if capacity <= 6 {
		t.Fatalf("near fusion capacity=%d want>6 groups=%v", capacity, groups)
	}
	if gap > continuousFusionTolerance+1e-12 {
		t.Fatalf("nearest gap=%g exceeds tolerance", gap)
	}
	if math.Abs(doseOffsetRMS(offsets)-multiplicityDoseRMS) > 1e-12 {
		t.Fatalf("normalized rms=%g", doseOffsetRMS(offsets))
	}
}
