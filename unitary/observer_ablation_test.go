package unitary

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestUP8BalancedSplitMarginals(t *testing.T) {
	train, err := selectObserverTables(true, 64)
	if err != nil {
		t.Fatal(err)
	}
	held, err := selectObserverTables(false, 64)
	if err != nil {
		t.Fatal(err)
	}
	if len(train) != 64 || len(held) != 64 {
		t.Fatalf("unexpected selected table counts train=%d held=%d", len(train), len(held))
	}
	if observerMarginalMinimum(train) <= 0 || observerMarginalMinimum(held) <= 0 {
		t.Fatal("a selected partition lost an entity/value marginal")
	}
	seen := map[int]bool{}
	for _, table := range train {
		if !balancedObserverTrainTable(table) {
			t.Fatalf("train table in held partition: %v", table)
		}
		seen[memoryTableIndex(table)] = true
	}
	for _, table := range held {
		if balancedObserverTrainTable(table) {
			t.Fatalf("held table in train partition: %v", table)
		}
		if seen[memoryTableIndex(table)] {
			t.Fatalf("table split overlap: %v", table)
		}
	}
}

func TestUP8CoherenceFeaturesPreserveGlobalPhase(t *testing.T) {
	state, err := Normalize(State{
		complex(1, 2), complex(-2, 0.5), complex(0.25, -1), complex(3, 0),
	})
	if err != nil {
		t.Fatal(err)
	}
	rotated := make(State, len(state))
	global := cmplx.Rect(1, 0.731)
	for i := range state {
		rotated[i] = state[i] * global
	}
	a, err := observerCoherenceFeatures(state)
	if err != nil {
		t.Fatal(err)
	}
	b, err := observerCoherenceFeatures(rotated)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 16 && len(state) == 4 {
		// The production UP-8 state is 16-D. This branch intentionally avoids
		// asserting 256 features for the smaller fixture below.
	}
	if len(a) != len(b) {
		t.Fatalf("feature length mismatch %d vs %d", len(a), len(b))
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > 1e-12 {
			t.Fatalf("global phase changed coherence feature %d: %g vs %g", i, a[i], b[i])
		}
	}
}

func TestUP8ProductionCoherenceFeatureLength(t *testing.T) {
	state, err := encodeMemory(memoryTable{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	features, err := observerCoherenceFeatures(state)
	if err != nil {
		t.Fatal(err)
	}
	if len(features) != 256 {
		t.Fatalf("feature length=%d want=256", len(features))
	}
}
