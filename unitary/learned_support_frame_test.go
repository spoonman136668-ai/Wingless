package unitary

import (
	"math"
	"testing"
)

func TestUP14InitialSupportIsDenseAndLowPurity(t *testing.T) {
	matrix, err := normalizedSupportMatrix(initialSupportMatrix())
	if err != nil {
		t.Fatal(err)
	}
	purities := supportPurities(matrix)
	for row, purity := range purities {
		if purity >= 0.30 {
			t.Fatalf("row=%d initial target purity=%g want<0.30", row, purity)
		}
		for column, value := range matrix[row] {
			if math.Abs(value) < 1e-12 {
				t.Fatalf(
					"row=%d column=%d initial support unexpectedly zero",
					row, column,
				)
			}
		}
	}
}

func TestUP14SupportBankIsGlobalPhaseInvariant(t *testing.T) {
	matrix, err := normalizedSupportMatrix(initialSupportMatrix())
	if err != nil {
		t.Fatal(err)
	}
	bank, err := makeLearnedSupportBank(
		up13LearnedPhases(), matrix,
	)
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

func TestUP14UnitaryCoevolutionPreservesMixedSupportFeature(t *testing.T) {
	matrix, err := normalizedSupportMatrix(initialSupportMatrix())
	if err != nil {
		t.Fatal(err)
	}
	bank, err := makeLearnedSupportBank(
		up13LearnedPhases(), matrix,
	)
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
