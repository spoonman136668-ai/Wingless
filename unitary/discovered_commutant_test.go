package unitary

import "testing"

func TestUP25DiscoverySeedDoesNotNeedHiddenStructure(t *testing.T) {
	seed, err := discoverySeedMatrix(0)
	if err != nil { t.Fatal(err) }
	if len(seed) != fullLatentDimension {
		t.Fatalf("seed dimension=%d want=%d", len(seed), fullLatentDimension)
	}
	norm, err := matrixFrobeniusNorm(seed)
	if err != nil { t.Fatal(err) }
	if norm < 0.999999 || norm > 1.000001 {
		t.Fatalf("seed frobenius norm=%g want~1", norm)
	}
}

func TestUP25BinaryConjugationAverageReducesCommutator(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	seed, err := discoverySeedMatrix(0)
	if err != nil { t.Fatal(err) }
	initial, err := maxWeylCommutatorEntry([]latentMatrix{seed}, step)
	if err != nil { t.Fatal(err) }
	discovered, _, err := discoverCommutingObservables(step, 1, 8)
	if err != nil { t.Fatal(err) }
	final, err := maxWeylCommutatorEntry(discovered, step)
	if err != nil { t.Fatal(err) }
	if !(final < initial) {
		t.Fatalf("commutator did not shrink initial=%g final=%g", initial, final)
	}
}

func TestUP25DiscoveredFeatureDimensions(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	observables, _, err := discoverCommutingObservables(step, discoveredObservableCount, 2)
	if err != nil { t.Fatal(err) }
	memory, err := encodeMemory(memoryTable{0,1,2,3})
	if err != nil { t.Fatal(err) }
	state, err := fullLatentEncode(memory, mixer)
	if err != nil { t.Fatal(err) }
	raw, err := discoveredRawFeatures(state, observables)
	if err != nil { t.Fatal(err) }
	if len(raw) != discoveredRawFeatureDim {
		t.Fatalf("raw dimension=%d want=%d", len(raw), discoveredRawFeatureDim)
	}
	quadratic, err := discoveredQuadraticFeatures(state, observables)
	if err != nil { t.Fatal(err) }
	if len(quadratic) != discoveredQuadraticFeatureDim {
		t.Fatalf("quadratic dimension=%d want=%d", len(quadratic), discoveredQuadraticFeatureDim)
	}
}
