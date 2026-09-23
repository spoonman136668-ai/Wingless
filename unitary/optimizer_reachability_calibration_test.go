package unitary

import (
	"math"
	"testing"
)

func TestUP41TargetReproducesAcceptedCapacity18Geometry(t *testing.T) {
	groups, offsets, err := optimizerCalibrationTarget()
	if err != nil {
		t.Fatal(err)
	}
	if got := symmetryGroupCapacity(groups); got != 18 {
		t.Fatalf("target capacity=%d want=18", got)
	}
	if math.Abs(doseOffsetRMS(offsets)-multiplicityDoseRMS) > 1e-12 {
		t.Fatalf(
			"target rms=%g want=%g",
			doseOffsetRMS(offsets), multiplicityDoseRMS,
		)
	}
	var sum float64
	for _, value := range offsets {
		sum += value
	}
	if math.Abs(sum) > 1e-12 {
		t.Fatalf("target mean sum=%g want=0", sum)
	}
	for _, pair := range [][2]int{{0, 1}, {1, 2}, {2, 5}} {
		if offsets[pair[0]] != offsets[pair[1]] {
			t.Fatalf("accepted target group is not exactly fused")
		}
	}
}

func TestUP41AnalyticCosineGradientIsTangent(t *testing.T) {
	offsets := continuousFusionInitialOffsets()
	_, target, err := optimizerCalibrationTarget()
	if err != nil {
		t.Fatal(err)
	}
	gradient, err := analyticOffsetCosineGradient(offsets, target)
	if err != nil {
		t.Fatal(err)
	}
	var mean, radial float64
	for index, value := range gradient {
		mean += value
		radial += value * offsets[index]
	}
	mean /= float64(len(gradient))
	if math.Abs(mean) > 1e-12 {
		t.Fatalf("gradient mean=%g want=0", mean)
	}
	if math.Abs(radial) > 1e-12 {
		t.Fatalf("gradient radial component=%g want=0", radial)
	}
}

func TestUP41WrongFusionDetection(t *testing.T) {
	groups, _, err := optimizerCalibrationTarget()
	if err != nil {
		t.Fatal(err)
	}
	target := symmetryGroupsJSON(groups)

	wrong, err := calibrationWrongCrossTargetFusion(
		[][]int{{0}, {1}, {2}, {3, 5}, {4}},
		target,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !wrong {
		t.Fatal("cross-target fusion was not detected")
	}

	wrong, err = calibrationWrongCrossTargetFusion(
		[][]int{{0, 1}, {2, 5}, {3}, {4}},
		target,
	)
	if err != nil {
		t.Fatal(err)
	}
	if wrong {
		t.Fatal("target-compatible fusion was incorrectly rejected")
	}
}

func TestUP41NoStickyUpdatePreservesSingletonGeometry(t *testing.T) {
	offsets := continuousFusionInitialOffsets()
	gradient := []float64{1, -1, 2, -2, 3, -3}
	next, update, err :=
		applyCalibrationGradientNoSticky(offsets, gradient)
	if err != nil {
		t.Fatal(err)
	}
	if len(next) != compositeChannels ||
		len(update) != compositeChannels {
		t.Fatal("no-sticky update dimension mismatch")
	}
	for _, value := range update {
		if math.Abs(value) > directOffsetMaxUpdate+1e-15 {
			t.Fatalf("update=%g exceeds bound", value)
		}
	}
	if math.Abs(doseOffsetRMS(next)-multiplicityDoseRMS) > 1e-12 {
		t.Fatalf("updated rms=%g", doseOffsetRMS(next))
	}
}
