package unitary

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestUP9FrameFeatureIsGlobalPhaseInvariant(t *testing.T) {
	bank, err := makeFrameBank()
	if err != nil {
		t.Fatal(err)
	}
	memory, err := encodeMemory(memoryTable{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	a, err := frameFeature(memory, bank.anchor, bank.phase[2])
	if err != nil {
		t.Fatal(err)
	}
	rotated := rotateGlobalPhase(memory, 1.234)
	b, err := frameFeature(rotated, bank.anchor, bank.phase[2])
	if err != nil {
		t.Fatal(err)
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > 1e-12 {
			t.Fatalf("global phase changed frame feature %d: %g vs %g", i, a[i], b[i])
		}
	}
}

func TestUP9UnitaryCoevolutionPreservesFrameFeature(t *testing.T) {
	bank, err := makeFrameBank()
	if err != nil {
		t.Fatal(err)
	}
	memory, err := encodeMemory(memoryTable{3, 2, 1, 0})
	if err != nil {
		t.Fatal(err)
	}
	memory = rotateGlobalPhase(memory, 0.731)
	before, err := frameFeature(memory, bank.anchor, bank.phase[1])
	if err != nil {
		t.Fatal(err)
	}
	block := stressProgram()
	evolvedMemory, err := applyStressUnitary(memory, block, 128)
	if err != nil {
		t.Fatal(err)
	}
	evolvedBank, err := transportFrameBank(bank, block, 128, applyStressUnitary, true)
	if err != nil {
		t.Fatal(err)
	}
	after, err := frameFeature(evolvedMemory, evolvedBank.anchor, evolvedBank.phase[1])
	if err != nil {
		t.Fatal(err)
	}
	for i := range before {
		if math.Abs(before[i]-after[i]) > 1e-12 {
			t.Fatalf("unitary coevolution changed frame feature %d: %g vs %g", i, before[i], after[i])
		}
	}
}

func TestUP9FrameBankHasExpectedPhaseCode(t *testing.T) {
	bank, err := makeFrameBank()
	if err != nil {
		t.Fatal(err)
	}
	for value := 0; value < 4; value++ {
		table := memoryTable{0, 0, value, 0}
		memory, err := encodeMemory(table)
		if err != nil {
			t.Fatal(err)
		}
		features, err := frameFeature(memory, bank.anchor, bank.phase[2])
		if err != nil {
			t.Fatal(err)
		}
		got := cmplx.Phase(complex(features[0], features[1]))
		want := -math.Pi * float64(value) / 2
		delta := math.Atan2(math.Sin(got-want), math.Cos(got-want))
		if math.Abs(delta) > 1e-12 {
			t.Fatalf("value=%d phase=%g want=%g", value, got, want)
		}
	}
}
