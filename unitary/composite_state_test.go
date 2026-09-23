package unitary

import (
	"math"
	"testing"
)

func TestUP16CompositePackingHasOneUnitNormState(t *testing.T) {
	memory, err := encodeMemory(memoryTable{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		t.Fatal(err)
	}
	composite, err := packCompositeState(memory, bank)
	if err != nil {
		t.Fatal(err)
	}
	if len(composite) != 96 {
		t.Fatalf("composite dimension=%d want=96", len(composite))
	}
	norm2, err := NormSquared(composite)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(norm2-1) > 1e-12 {
		t.Fatalf("composite norm=%g want=1", norm2)
	}
}

func TestUP16CompositeObservableMatchesSeparateFrame(t *testing.T) {
	for _, depth := range []int{1, 8, 128, 1024} {
		errValue, err := compositeObservableEquivalence(depth)
		if err != nil {
			t.Fatal(err)
		}
		if errValue > 1e-12 {
			t.Fatalf(
				"depth=%d composite observable error=%g want<=1e-12",
				depth, errValue,
			)
		}
	}
}

func TestUP16UnitaryCompositeRoundTrip(t *testing.T) {
	memory, err := encodeMemory(memoryTable{3, 2, 1, 0})
	if err != nil {
		t.Fatal(err)
	}
	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		t.Fatal(err)
	}
	composite, err := packCompositeState(memory, bank)
	if err != nil {
		t.Fatal(err)
	}
	block := stressProgram()
	forward, err := evolveCompositeState(
		composite, block, 128, applyStressUnitary,
	)
	if err != nil {
		t.Fatal(err)
	}
	reversed := reverseCouplings(block)
	roundTrip, err := evolveCompositeState(
		forward, reversed, 128, applyStressUnitary,
	)
	if err != nil {
		t.Fatal(err)
	}
	distance, err := L2Distance(composite, roundTrip)
	if err != nil {
		t.Fatal(err)
	}
	if distance > 1e-11 {
		t.Fatalf("composite round-trip error=%g want<=1e-11", distance)
	}
}
