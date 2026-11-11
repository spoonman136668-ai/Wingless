package unitary

import (
	"math"
	"testing"
)

func TestUP13InitialPhaseAlphabetIsClustered(t *testing.T) {
	phases := []float64{0, 0.12, 0.24, 0.36}
	if got := minPhaseSeparation(phases); math.Abs(got-0.12) > 1e-12 {
		t.Fatalf("initial minimum separation=%g want=0.12", got)
	}
}

func TestUP13PhaseBankIsGlobalPhaseInvariant(t *testing.T) {
	phases := []float64{0, -1.6, 1.05, 2.35}
	bank, err := makeLearnedPhaseBank(phases)
	if err != nil {
		t.Fatal(err)
	}
	memory, err := encodeMemory(memoryTable{3, 2, 1, 0})
	if err != nil {
		t.Fatal(err)
	}
	for entity := 0; entity < 4; entity++ {
		a, err := frameFeature(memory, bank.anchor, bank.phase[entity])
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

func TestUP13UnitaryCoevolutionPreservesLearnedFrameFeature(t *testing.T) {
	phases := []float64{0, -1.6, 1.05, 2.35}
	bank, err := makeLearnedPhaseBank(phases)
	if err != nil {
		t.Fatal(err)
	}
	memory, err := encodeMemory(memoryTable{1, 3, 0, 2})
	if err != nil {
		t.Fatal(err)
	}
	memory = rotateGlobalPhase(memory, 0.771)
	block := stressProgram()
	evolvedMemory, err := applyStressUnitary(memory, block, 1024)
	if err != nil {
		t.Fatal(err)
	}
	evolvedBank, err := transportFrameBank(
		bank, block, 1024, applyStressUnitary, true,
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
