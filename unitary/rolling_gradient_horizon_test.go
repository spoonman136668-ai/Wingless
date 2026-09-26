package unitary

import (
	"math"
	"testing"
)

func TestUP42RollingReconstructionRecoversKnownTangentGradient(t *testing.T) {
	want := []float64{1, -2, 3, -4, 5, -3}
	var mean float64
	for _, value := range want {
		mean += value
	}
	if mean != 0 {
		t.Fatal("test gradient must be zero mean")
	}

	union := newFusionUnion()
	observations := make([]rollingDirectionalObservation, 0, 7)
	for step := 1; step <= 7; step++ {
		direction := deterministicSPSADirection(
			step, union.groups(),
		)
		var derivative float64
		for index, value := range direction {
			derivative += value * want[index]
		}
		observations = append(
			observations,
			rollingDirectionalObservation{
				direction: direction,
				derivative: derivative,
			},
		)
	}

	got, fullRank, err :=
		rollingFullRankGradient(observations)
	if err != nil {
		t.Fatal(err)
	}
	if !fullRank {
		t.Fatal("seven-direction buffer did not reach full rank")
	}
	for index := range want {
		if math.Abs(got[index]-want[index]) > 1e-10 {
			t.Fatalf(
				"gradient[%d]=%g want=%g",
				index, got[index], want[index],
			)
		}
	}
}

func TestUP42HorizonScheduleIsFrozenDoublingSequence(t *testing.T) {
	want := []int{12, 24, 48, 96}
	if len(rollingGradientHorizons) != len(want) {
		t.Fatalf("horizon count=%d", len(rollingGradientHorizons))
	}
	for index := range want {
		if rollingGradientHorizons[index] != want[index] {
			t.Fatalf(
				"horizon[%d]=%d want=%d",
				index,
				rollingGradientHorizons[index],
				want[index],
			)
		}
	}
}

func TestUP42RollingBufferBound(t *testing.T) {
	if rollingGradientBufferSize != 7 {
		t.Fatalf(
			"rolling buffer=%d want=7",
			rollingGradientBufferSize,
		)
	}
}
