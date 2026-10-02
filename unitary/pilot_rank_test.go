package unitary

import (
	"math"
	"testing"
)

func TestUP11PilotRankFeatureBudget(t *testing.T) {
	memory, err := encodeMemory(memoryTable{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	for rank := 1; rank <= 4; rank++ {
		bank, err := makePilotRankBank(rank)
		if err != nil {
			t.Fatal(err)
		}
		if len(bank.codes) != rank {
			t.Fatalf("rank=%d code count=%d", rank, len(bank.codes))
		}
		features, err := pilotRankFeatures(memory, bank)
		if err != nil {
			t.Fatal(err)
		}
		if len(features) != 2*rank {
			t.Fatalf("rank=%d feature count=%d want=%d", rank, len(features), 2*rank)
		}
	}
}

func TestUP11PilotRankFeaturesAreGlobalPhaseInvariant(t *testing.T) {
	memory, err := encodeMemory(memoryTable{3, 1, 0, 2})
	if err != nil {
		t.Fatal(err)
	}
	bank, err := makePilotRankBank(4)
	if err != nil {
		t.Fatal(err)
	}
	a, err := pilotRankFeatures(memory, bank)
	if err != nil {
		t.Fatal(err)
	}
	b, err := pilotRankFeatures(rotateGlobalPhase(memory, 1.271), bank)
	if err != nil {
		t.Fatal(err)
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > 1e-12 {
			t.Fatalf("global phase changed rank feature %d: %g vs %g", i, a[i], b[i])
		}
	}
}

func TestUP11UnitaryCoevolutionPreservesRankFeatures(t *testing.T) {
	memory, err := encodeMemory(memoryTable{2, 0, 3, 1})
	if err != nil {
		t.Fatal(err)
	}
	memory = rotateGlobalPhase(memory, 0.441)
	bank, err := makePilotRankBank(3)
	if err != nil {
		t.Fatal(err)
	}
	before, err := pilotRankFeatures(memory, bank)
	if err != nil {
		t.Fatal(err)
	}
	block := stressProgram()
	evolvedMemory, err := applyStressUnitary(memory, block, 1024)
	if err != nil {
		t.Fatal(err)
	}
	evolvedBank, err := transportPilotRankBank(bank, block, 1024, applyStressUnitary)
	if err != nil {
		t.Fatal(err)
	}
	after, err := pilotRankFeatures(evolvedMemory, evolvedBank)
	if err != nil {
		t.Fatal(err)
	}
	for i := range before {
		if math.Abs(before[i]-after[i]) > 1e-12 {
			t.Fatalf("coevolution changed rank feature %d: %g vs %g", i, before[i], after[i])
		}
	}
}

func TestUP11SelectedSplitKeepsAllMarginals(t *testing.T) {
	train, err := selectObserverTables(true, 32)
	if err != nil {
		t.Fatal(err)
	}
	held, err := selectObserverTables(false, 32)
	if err != nil {
		t.Fatal(err)
	}
	if observerMarginalMinimum(train) <= 0 || observerMarginalMinimum(held) <= 0 {
		t.Fatal("rank-sweep table selection lost a marginal")
	}
}
