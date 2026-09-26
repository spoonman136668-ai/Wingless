package unitary

import "testing"

func TestUP30SplitTransportRemainsOrthogonalEquivalent(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := splitFullLatentStep(mixer)
	if err != nil { t.Fatal(err) }
	realStep, err := realifyMatrix(step)
	if err != nil { t.Fatal(err) }
	maximum, err := maxRealOrthogonalityError(realStep)
	if err != nil { t.Fatal(err) }
	if maximum > 1e-10 {
		t.Fatalf("split orthogonality error=%g", maximum)
	}
}

func TestUP30SplitBreaksCrossChannelShiftSymmetry(t *testing.T) {
	mixer := fullLatentMixer()
	base, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { t.Fatal(err) }
	split, err := splitFullLatentStep(mixer)
	if err != nil { t.Fatal(err) }
	shift, err := fullLatentCrossChannelShift(mixer)
	if err != nil { t.Fatal(err) }
	baseError, err := maxWeylCommutatorEntry([]latentMatrix{shift}, base)
	if err != nil { t.Fatal(err) }
	splitError, err := maxWeylCommutatorEntry([]latentMatrix{shift}, split)
	if err != nil { t.Fatal(err) }
	if baseError > 1e-10 {
		t.Fatalf("base shift no longer commutes: %g", baseError)
	}
	if splitError < 1e-4 {
		t.Fatalf("split shift commutator too small: %g", splitError)
	}
}
