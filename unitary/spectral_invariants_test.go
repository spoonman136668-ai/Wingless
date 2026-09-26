package unitary

import (
	"math"
	"testing"
)

func TestUP23SpectralMomentDimension(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	memory, err := encodeMemory(memoryTable{0,1,2,3})
	if err != nil { t.Fatal(err) }
	state, err := fullLatentEncode(memory, mixer)
	if err != nil { t.Fatal(err) }
	features, err := spectralMomentFeatures(state, step, spectralMomentCount)
	if err != nil { t.Fatal(err) }
	if len(features) != spectralFeatureDim {
		t.Fatalf("spectral feature dimension=%d want=%d", len(features), spectralFeatureDim)
	}
}

func TestUP23UnitarySpectralMomentsAreDepthInvariant(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	op, err := latentMatrixPower(step, 128)
	if err != nil { t.Fatal(err) }
	drift, err := spectralMomentDrift(mixer, step, op)
	if err != nil { t.Fatal(err) }
	if drift > 1e-10 {
		t.Fatalf("spectral moment drift=%g want<=1e-10", drift)
	}
}

func TestUP23SpectralMomentsIgnoreGlobalPhase(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	memory, err := encodeMemory(memoryTable{3,2,1,0})
	if err != nil { t.Fatal(err) }
	state, err := fullLatentEncode(memory, mixer)
	if err != nil { t.Fatal(err) }
	want, err := spectralMomentFeatures(state, step, spectralMomentCount)
	if err != nil { t.Fatal(err) }
	got, err := spectralMomentFeatures(rotateGlobalPhase(state, 1.234), step, spectralMomentCount)
	if err != nil { t.Fatal(err) }
	for i := range want {
		if math.Abs(want[i]-got[i]) > 1e-12 {
			t.Fatalf("global-phase spectral drift[%d]=%g", i, math.Abs(want[i]-got[i]))
		}
	}
}
