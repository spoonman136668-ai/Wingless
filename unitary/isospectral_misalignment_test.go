package unitary

import "testing"

func TestUP31PerCopyConjugatorIsUnitaryPermutation(t *testing.T) {
	conjugator, err := perCopyPermutationConjugator()
	if err != nil { t.Fatal(err) }
	adjoint, err := latentAdjoint(conjugator)
	if err != nil { t.Fatal(err) }
	product, err := latentMatrixMultiply(conjugator, adjoint)
	if err != nil { t.Fatal(err) }
	identity := identityLatentMatrix(fullLatentDimension)
	errValue, err := maxLatentMatrixEntryDifference(product, identity)
	if err != nil { t.Fatal(err) }
	if errValue > 1e-12 {
		t.Fatalf("per-copy conjugator unitarity error=%g", errValue)
	}
}

func TestUP31MisalignedTransportIsExactlyConjugate(t *testing.T) {
	baseMixer := fullLatentMixer()
	misalignedMixer, err := isospectralMisalignedMixer()
	if err != nil { t.Fatal(err) }
	baseStep, err := conjugatedLatentStep(baseMixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	misalignedStep, err := conjugatedLatentStep(misalignedMixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	errValue, err := isospectralEquivalenceError(
		baseMixer, misalignedMixer, baseStep, misalignedStep,
	)
	if err != nil { t.Fatal(err) }
	if errValue > 1e-10 {
		t.Fatalf("isospectral conjugacy error=%g", errValue)
	}
}

func TestUP31BreaksNaiveButPreservesCorrectShiftIntertwiner(t *testing.T) {
	baseMixer := fullLatentMixer()
	misalignedMixer, err := isospectralMisalignedMixer()
	if err != nil { t.Fatal(err) }
	step, err := conjugatedLatentStep(misalignedMixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	naive, err := fullLatentCrossChannelShift(baseMixer)
	if err != nil { t.Fatal(err) }
	correct, err := fullLatentCrossChannelShift(misalignedMixer)
	if err != nil { t.Fatal(err) }

	naiveError, err := maxWeylCommutatorEntry([]latentMatrix{naive}, step)
	if err != nil { t.Fatal(err) }
	correctError, err := maxWeylCommutatorEntry([]latentMatrix{correct}, step)
	if err != nil { t.Fatal(err) }

	if naiveError < 1e-4 {
		t.Fatalf("naive alignment was not broken: %g", naiveError)
	}
	if correctError > 1e-10 {
		t.Fatalf("correct hidden intertwiner no longer commutes: %g", correctError)
	}
}
