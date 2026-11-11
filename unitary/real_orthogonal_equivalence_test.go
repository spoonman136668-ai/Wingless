package unitary

import (
	"math"
	"testing"
)

func TestUP29RealificationPreservesMatrixVectorAction(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	realStep, err := realifyMatrix(step)
	if err != nil { t.Fatal(err) }

	memory, err := encodeMemory(memoryTable{0,1,2,3})
	if err != nil { t.Fatal(err) }
	state, err := fullLatentEncode(memory, mixer)
	if err != nil { t.Fatal(err) }

	complexNext, err := latentMatrixVector(step, state)
	if err != nil { t.Fatal(err) }
	realState, err := realifyState(state)
	if err != nil { t.Fatal(err) }
	realNext, err := realMatrixVector(realStep, realState)
	if err != nil { t.Fatal(err) }
	converted, err := complexifyRealState(realNext)
	if err != nil { t.Fatal(err) }
	distance, err := L2Distance(converted, complexNext)
	if err != nil { t.Fatal(err) }
	if distance > 1e-12 {
		t.Fatalf("realification action error=%g", distance)
	}
}

func TestUP29RealifiedUnitaryIsOrthogonal(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	realStep, err := realifyMatrix(step)
	if err != nil { t.Fatal(err) }
	maximum, err := maxRealOrthogonalityError(realStep)
	if err != nil { t.Fatal(err) }
	if maximum > 1e-10 {
		t.Fatalf("orthogonality error=%g", maximum)
	}
}

func TestUP29RealObservableExpectationMatchesComplex(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	bank, err := frozenUP28ObservableBank(step)
	if err != nil { t.Fatal(err) }
	realObs, err := realifyObservable(bank[0])
	if err != nil { t.Fatal(err) }

	memory, err := encodeMemory(memoryTable{3,1,0,2})
	if err != nil { t.Fatal(err) }
	state, err := fullLatentEncode(memory, mixer)
	if err != nil { t.Fatal(err) }
	applied, err := latentMatrixVector(bank[0], state)
	if err != nil { t.Fatal(err) }
	want, err := stateInner(state, applied)
	if err != nil { t.Fatal(err) }

	realState, err := realifyState(state)
	if err != nil { t.Fatal(err) }
	gotRe, gotIm, err := realObservableExpectation(realState, realObs)
	if err != nil { t.Fatal(err) }
	if math.Abs(gotRe-real(want)) > 1e-12 ||
		math.Abs(gotIm-imag(want)) > 1e-12 {
		t.Fatalf(
			"observable mismatch got=(%g,%g) want=(%g,%g)",
			gotRe, gotIm, real(want), imag(want),
		)
	}
}
