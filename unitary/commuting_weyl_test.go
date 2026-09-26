package unitary

import "testing"

func TestUP24WeylObservableCountAndFeatureDimensions(t *testing.T) {
	mixer := fullLatentMixer()
	observables, err := fullLatentWeylObservables(mixer)
	if err != nil { t.Fatal(err) }
	if len(observables) != weylObservableCount {
		t.Fatalf("observable count=%d want=%d", len(observables), weylObservableCount)
	}
	memory, err := encodeMemory(memoryTable{0,1,2,3})
	if err != nil { t.Fatal(err) }
	state, err := fullLatentEncode(memory, mixer)
	if err != nil { t.Fatal(err) }
	raw, err := weylRawFeatures(state, observables)
	if err != nil { t.Fatal(err) }
	if len(raw) != weylRawFeatureDim {
		t.Fatalf("raw feature dimension=%d want=%d", len(raw), weylRawFeatureDim)
	}
	quadratic, err := weylQuadraticFeatures(state, observables)
	if err != nil { t.Fatal(err) }
	if len(quadratic) != weylQuadraticFeatureDim {
		t.Fatalf("quadratic feature dimension=%d want=%d", len(quadratic), weylQuadraticFeatureDim)
	}
}

func TestUP24WeylObservablesCommuteWithFullLatentTransport(t *testing.T) {
	mixer := fullLatentMixer()
	observables, err := fullLatentWeylObservables(mixer)
	if err != nil { t.Fatal(err) }
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	maximum, err := maxWeylCommutatorEntry(observables, step)
	if err != nil { t.Fatal(err) }
	if maximum > 1e-10 {
		t.Fatalf("weyl commutator error=%g want<=1e-10", maximum)
	}
}

func TestUP24WeylFeaturesAreDepthInvariant(t *testing.T) {
	mixer := fullLatentMixer()
	observables, err := fullLatentWeylObservables(mixer)
	if err != nil { t.Fatal(err) }
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	op, err := latentMatrixPower(step, 128)
	if err != nil { t.Fatal(err) }
	drift, err := weylFeatureDrift(mixer, observables, op)
	if err != nil { t.Fatal(err) }
	if drift > 1e-9 {
		t.Fatalf("weyl feature drift=%g want<=1e-9", drift)
	}
}
