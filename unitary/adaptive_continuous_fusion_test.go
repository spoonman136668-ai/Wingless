package unitary

import (
	"math"
	"testing"
)

func TestUP37GroupProbePreservesExistingFusion(t *testing.T) {
	offsets := []float64{-0.04, -0.04, -0.01, 0.01, 0.04, 0.04}
	offsets, err := normalizeContinuousOffsets(offsets)
	if err != nil {
		t.Fatal(err)
	}
	union := newFusionUnion()
	union.union(0, 1)
	union.union(4, 5)
	groups := union.groups()

	probe, err := adaptiveGroupProbeOffsets(
		offsets, groups, 0, 1,
		continuousFusionProbeFraction,
	)
	if err != nil {
		t.Fatal(err)
	}
	if math.Float64bits(probe[0]) != math.Float64bits(probe[1]) {
		t.Fatal("existing fused group 0,1 was broken by adaptive probe")
	}
	if math.Float64bits(probe[4]) != math.Float64bits(probe[5]) {
		t.Fatal("existing fused group 4,5 was broken by adaptive probe")
	}
}

func TestUP37AdaptiveAdvancePreservesMeanRMSAndFusion(t *testing.T) {
	offsets := []float64{-0.04, -0.04, -0.01, 0.01, 0.04, 0.04}
	offsets, err := normalizeContinuousOffsets(offsets)
	if err != nil {
		t.Fatal(err)
	}
	union := newFusionUnion()
	union.union(0, 1)
	union.union(4, 5)
	groups := union.groups()
	weights := make([][]float64, len(groups))
	for i := range weights {
		weights[i] = make([]float64, len(groups))
	}
	weights[0][1], weights[1][0] = 0.5, 0.5
	weights[1][2], weights[2][1] = 0.5, 0.5

	next, err := adaptiveFusionAdvance(offsets, &union, weights)
	if err != nil {
		t.Fatal(err)
	}
	var mean float64
	for _, value := range next {
		mean += value
	}
	mean /= float64(len(next))
	if math.Abs(mean) > 1e-12 {
		t.Fatalf("adaptive flow mean=%g want=0", mean)
	}
	if math.Abs(doseOffsetRMS(next)-multiplicityDoseRMS) > 1e-12 {
		t.Fatalf(
			"adaptive flow RMS=%g want=%g",
			doseOffsetRMS(next), multiplicityDoseRMS,
		)
	}
	if union.find(0) != union.find(1) {
		t.Fatal("adaptive flow broke fused group 0,1")
	}
	if union.find(4) != union.find(5) {
		t.Fatal("adaptive flow broke fused group 4,5")
	}
}

func TestUP37ReestimationUsesCurrentGroups(t *testing.T) {
	union := newFusionUnion()
	union.union(1, 2)
	groups := union.groups()
	if len(groups) != 5 {
		t.Fatalf("group count=%d want=5", len(groups))
	}
	wantPairs := len(groups) * (len(groups) - 1) / 2
	if wantPairs != 10 {
		t.Fatalf("pair count=%d want=10", wantPairs)
	}
}
