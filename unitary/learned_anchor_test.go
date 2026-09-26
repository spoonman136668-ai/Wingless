package unitary

import (
	"math"
	"testing"
)

func TestUP15InitialAnchorIsNormalizedAndIrregular(t *testing.T) {
	params := initialLearnedAnchorParams()
	state, err := anchorStateFromParams(params)
	if err != nil {
		t.Fatal(err)
	}
	norm2, err := NormSquared(state)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(norm2-1) > 1e-12 {
		t.Fatalf("anchor norm=%g want=1", norm2)
	}
	concentration, err := anchorPhaseConcentration(params)
	if err != nil {
		t.Fatal(err)
	}
	if concentration >= 0.95 {
		t.Fatalf("initial anchor phase concentration=%g unexpectedly high", concentration)
	}
}

func TestUP15AnchorGaugeIsCanonical(t *testing.T) {
	params := initialLearnedAnchorParams()
	state, err := anchorStateFromParams(params)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(imag(state[0])) > 1e-12 {
		t.Fatalf("anchor gauge imag(state[0])=%g", imag(state[0]))
	}
	if real(state[0]) <= 0 {
		t.Fatalf("anchor gauge real(state[0])=%g want positive", real(state[0]))
	}
}

func TestUP15LearnedAnchorBankRemainsGlobalPhaseInvariant(t *testing.T) {
	params := initialLearnedAnchorParams()
	bank, err := makeTaskLearnedAnchorBank(params)
	if err != nil {
		t.Fatal(err)
	}
	memory, err := encodeMemory(memoryTable{3, 2, 1, 0})
	if err != nil {
		t.Fatal(err)
	}
	for entity := 0; entity < 4; entity++ {
		a, err := frameFeature(
			memory, bank.anchor, bank.phase[entity],
		)
		if err != nil {
			t.Fatal(err)
		}
		b, err := frameFeature(
			rotateGlobalPhase(memory, 1.337),
			bank.anchor,
			bank.phase[entity],
		)
		if err != nil {
			t.Fatal(err)
		}
		for i := range a {
			if math.Abs(a[i]-b[i]) > 1e-12 {
				t.Fatalf(
					"entity=%d feature=%d global phase changed %g vs %g",
					entity, i, a[i], b[i],
				)
			}
		}
	}
}

func TestUP15UnitaryCoevolutionPreservesLearnedAnchorFeature(t *testing.T) {
	params := initialLearnedAnchorParams()
	bank, err := makeTaskLearnedAnchorBank(params)
	if err != nil {
		t.Fatal(err)
	}
	memory, err := encodeMemory(memoryTable{1, 3, 0, 2})
	if err != nil {
		t.Fatal(err)
	}
	memory = rotateGlobalPhase(memory, 0.771)
	block := stressProgram()
	evolvedMemory, err := applyStressUnitary(
		memory, block, 1024,
	)
	if err != nil {
		t.Fatal(err)
	}
	evolvedBank, err := transportFrameBank(
		bank, block, 1024,
		applyStressUnitary, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	for entity := 0; entity < 4; entity++ {
		before, err := frameFeature(
			memory, bank.anchor, bank.phase[entity],
		)
		if err != nil {
			t.Fatal(err)
		}
		after, err := frameFeature(
			evolvedMemory,
			evolvedBank.anchor,
			evolvedBank.phase[entity],
		)
		if err != nil {
			t.Fatal(err)
		}
		for i := range before {
			if math.Abs(before[i]-after[i]) > 1e-12 {
				t.Fatalf(
					"entity=%d feature=%d coevolution changed %g vs %g",
					entity, i, before[i], after[i],
				)
			}
		}
	}
}
