package unitary

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestUP22FullLatentMixerIsUnitaryAndDistributed(t *testing.T) {
	mixer := fullLatentMixer()
	adjoint, err := latentAdjoint(mixer)
	if err != nil {
		t.Fatal(err)
	}
	product, err := latentMatrixMultiply(mixer, adjoint)
	if err != nil {
		t.Fatal(err)
	}
	maximum := 0.0
	for row := 0; row < fullLatentDimension; row++ {
		for column := 0; column < fullLatentDimension; column++ {
			want := complex(0, 0)
			if row == column {
				want = 1
			}
			delta := cmplx.Abs(product[row][column] - want)
			if delta > maximum {
				maximum = delta
			}
		}
	}
	if maximum > 1e-11 {
		t.Fatalf("full latent mixer orthogonality error=%g", maximum)
	}

	participation, err := fullLatentMixerParticipation(mixer)
	if err != nil {
		t.Fatal(err)
	}
	if participation < 90 {
		t.Fatalf(
			"minimum mixer participation=%g want>=90",
			participation,
		)
	}
}

func TestUP22HermitianFeaturesIgnoreGlobalPhase(t *testing.T) {
	memory, err := encodeMemory(memoryTable{3, 0, 2, 1})
	if err != nil {
		t.Fatal(err)
	}
	mixer := fullLatentMixer()
	state, err := fullLatentEncode(memory, mixer)
	if err != nil {
		t.Fatal(err)
	}
	want, err := fullLatentHermitianFeatures(state)
	if err != nil {
		t.Fatal(err)
	}
	got, err := fullLatentHermitianFeatures(
		rotateGlobalPhase(state, 1.231),
	)
	if err != nil {
		t.Fatal(err)
	}
	maximum := 0.0
	for i := range want {
		delta := math.Abs(want[i] - got[i])
		if delta > maximum {
			maximum = delta
		}
	}
	if maximum > 1e-11 {
		t.Fatalf("global-phase feature drift=%g", maximum)
	}
}

func TestUP22ConjugatedTransportMatchesPackedTransport(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil {
		t.Fatal(err)
	}
	operators, err := latentDepthOperators(
		step, []int{0, 1, 32, 128},
	)
	if err != nil {
		t.Fatal(err)
	}
	errorValue, err := fullLatentRecoveryError(
		mixer,
		operators,
		[]int{0, 1, 32, 128},
	)
	if err != nil {
		t.Fatal(err)
	}
	if errorValue > 1e-10 {
		t.Fatalf(
			"conjugated transport recovery error=%g want<=1e-10",
			errorValue,
		)
	}
}
