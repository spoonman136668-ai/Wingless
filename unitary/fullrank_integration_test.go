package unitary

import "testing"

func TestUP12FullBalancedPools(t *testing.T) {
	train := fullObserverTablePool(true)
	held := fullObserverTablePool(false)

	if len(train) != 128 || len(held) != 128 {
		t.Fatalf("pool sizes train=%d held=%d want=128/128", len(train), len(held))
	}
	if observerMarginalMinimum(train) != 32 {
		t.Fatalf("train minimum marginal=%d want=32", observerMarginalMinimum(train))
	}
	if observerMarginalMinimum(held) != 32 {
		t.Fatalf("held minimum marginal=%d want=32", observerMarginalMinimum(held))
	}

	seen := map[int]bool{}
	for _, table := range train {
		seen[memoryTableIndex(table)] = true
	}
	for _, table := range held {
		if seen[memoryTableIndex(table)] {
			t.Fatalf("table split overlap: %v", table)
		}
	}
}

func TestUP12FullRankFeatureBudget(t *testing.T) {
	bank, err := makePilotRankBank(4)
	if err != nil {
		t.Fatal(err)
	}
	state, err := encodeMemory(memoryTable{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	features, err := pilotRankFeatures(state, bank)
	if err != nil {
		t.Fatal(err)
	}
	if len(bank.codes) != 4 {
		t.Fatalf("code directions=%d want=4", len(bank.codes))
	}
	if len(features) != 8 {
		t.Fatalf("feature count=%d want=8", len(features))
	}
}
