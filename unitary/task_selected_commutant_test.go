package unitary

import "testing"

func TestUP27SelectionIsDeterministicAndBounded(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	candidates, _, err := discoverCommutingObservables(step, 16, 2)
	if err != nil { t.Fatal(err) }
	tables := selectFirstMemoryTables(fullObserverTablePool(true), 8)
	states, targetTables, err := taskSelectedTrainingStates(tables, mixer, 0.05, 1)
	if err != nil { t.Fatal(err) }
	first, _, err := selectTaskRelevantObservables(states, targetTables, candidates, 8)
	if err != nil { t.Fatal(err) }
	second, _, err := selectTaskRelevantObservables(states, targetTables, candidates, 8)
	if err != nil { t.Fatal(err) }
	if len(first) != 8 || len(second) != 8 {
		t.Fatalf("selection size first=%d second=%d", len(first), len(second))
	}
	seen := map[int]bool{}
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("selection nondeterministic at %d: %d vs %d", i, first[i], second[i])
		}
		if first[i] < 0 || first[i] >= 16 {
			t.Fatalf("selected index out of bounds: %d", first[i])
		}
		if seen[first[i]] {
			t.Fatalf("duplicate selected index: %d", first[i])
		}
		seen[first[i]] = true
	}
}

func TestUP27PhaseTargetMatrixHasEightComponents(t *testing.T) {
	tables := []memoryTable{{0,1,2,3},{3,2,1,0}}
	targets, err := phaseTargetMatrix(tables)
	if err != nil { t.Fatal(err) }
	if len(targets) != 2 || len(targets[0]) != 8 || len(targets[1]) != 8 {
		t.Fatalf("phase target matrix dimensions invalid")
	}
}
