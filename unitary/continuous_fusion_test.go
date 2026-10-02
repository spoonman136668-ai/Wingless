package unitary

import (
	"math"
	"testing"
)

func TestUP36InitialOffsetsAreIndependent(t *testing.T) {
	offsets := continuousFusionInitialOffsets()
	if len(offsets) != compositeChannels {
		t.Fatalf("offset count=%d want=%d", len(offsets), compositeChannels)
	}
	if got := doseMaxMultiplicity(offsets); got != 1 {
		t.Fatalf("initial multiplicity=%d want=1", got)
	}
	if math.Abs(doseOffsetRMS(offsets)-multiplicityDoseRMS) > 1e-12 {
		t.Fatalf("initial RMS=%g want=%g", doseOffsetRMS(offsets), multiplicityDoseRMS)
	}
}

func TestUP36PairProbeIsContinuousNotMerged(t *testing.T) {
	offsets := continuousFusionInitialOffsets()
	probe, err := continuousPairProbeOffsets(
		offsets, 0, 1, continuousFusionProbeFraction,
	)
	if err != nil {
		t.Fatal(err)
	}
	if got := doseMaxMultiplicity(probe); got != 1 {
		t.Fatalf("probe created exact multiplicity=%d want=1", got)
	}
	if !(math.Abs(probe[0]-probe[1]) < math.Abs(offsets[0]-offsets[1])) {
		t.Fatal("pair probe did not reduce distance continuously")
	}
}

func TestUP36FusionFlowPreservesMeanAndRMS(t *testing.T) {
	offsets := continuousFusionInitialOffsets()
	var raw [compositeChannels][compositeChannels]float64
	raw[0][1], raw[1][0] = 1, 1
	weights := normalizeFusionWeights(raw)
	union := newFusionUnion()
	next, err := continuousFusionAdvance(offsets, weights, &union)
	if err != nil {
		t.Fatal(err)
	}
	var mean float64
	for _, value := range next {
		mean += value
	}
	mean /= float64(len(next))
	if math.Abs(mean) > 1e-12 {
		t.Fatalf("flow mean=%g want=0", mean)
	}
	if math.Abs(doseOffsetRMS(next)-multiplicityDoseRMS) > 1e-12 {
		t.Fatalf("flow RMS=%g want=%g", doseOffsetRMS(next), multiplicityDoseRMS)
	}
}

func TestUP36StickyFusionCreatesExactDegeneracy(t *testing.T) {
	offsets := []float64{-0.03, -0.029, -0.01, 0.01, 0.029, 0.03}
	offsets, err := normalizeContinuousOffsets(offsets)
	if err != nil {
		t.Fatal(err)
	}
	var raw [compositeChannels][compositeChannels]float64
	raw[0][1], raw[1][0] = 1, 1
	weights := normalizeFusionWeights(raw)
	union := newFusionUnion()
	for step := 0; step < 8; step++ {
		offsets, err = continuousFusionAdvance(offsets, weights, &union)
		if err != nil {
			t.Fatal(err)
		}
	}
	if union.find(0) != union.find(1) {
		t.Fatal("continuous flow did not trigger sticky fusion")
	}
	if math.Float64bits(offsets[0]) != math.Float64bits(offsets[1]) {
		t.Fatal("sticky fusion is not exact")
	}
}
