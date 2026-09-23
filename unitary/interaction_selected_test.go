package unitary

import (
	"math"
	"testing"
)

func TestUP28PolynomialFeatureDotMatchesExplicitLift(t *testing.T) {
	a := []float64{0.2, -0.3, 0.7, 0.1}
	b := []float64{-0.4, 0.5, 0.2, -0.6}
	got, err := polynomialFeatureDot(a, b)
	if err != nil { t.Fatal(err) }

	var want float64
	for i := range a {
		want += a[i] * b[i]
	}
	for first := 0; first < len(a); first++ {
		for second := first; second < len(a); second++ {
			want += a[first] * a[second] * b[first] * b[second]
		}
	}
	if math.Abs(got-want) > 1e-12 {
		t.Fatalf("polynomial dot=%g want=%g", got, want)
	}
}

func TestUP28InteractionSelectionDeterministic(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	candidates, _, err := discoverCommutingObservables(step, 16, 2)
	if err != nil { t.Fatal(err) }
	tables := selectFirstMemoryTables(fullObserverTablePool(true), 8)
	states, targetTables, err := taskSelectedTrainingStates(tables, mixer, 0.05, 1)
	if err != nil { t.Fatal(err) }
	first, _, err := selectInteractionRelevantObservables(states, targetTables, candidates, 8)
	if err != nil { t.Fatal(err) }
	second, _, err := selectInteractionRelevantObservables(states, targetTables, candidates, 8)
	if err != nil { t.Fatal(err) }
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("interaction selection nondeterministic at %d", i)
		}
	}
}
