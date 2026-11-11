package unitary

import (
	"math"
	"testing"
)

func TestUP32DoseOffsetsHaveExpectedMultiplicityAndRMS(t *testing.T) {
	specs := multiplicityDoseOffsets()
	wantMultiplicity := []int{6, 4, 3, 2, 1}
	if len(specs) != len(wantMultiplicity) {
		t.Fatalf("dose arm count=%d want=%d", len(specs), len(wantMultiplicity))
	}
	for index, spec := range specs {
		gotMultiplicity := doseMaxMultiplicity(spec.offsets)
		if gotMultiplicity != wantMultiplicity[index] {
			t.Fatalf(
				"arm=%s multiplicity=%d want=%d",
				spec.name, gotMultiplicity, wantMultiplicity[index],
			)
		}
		if index == 0 {
			if doseOffsetRMS(spec.offsets) != 0 {
				t.Fatalf("control RMS=%g want=0", doseOffsetRMS(spec.offsets))
			}
			continue
		}
		if math.Abs(doseOffsetRMS(spec.offsets)-multiplicityDoseRMS) > 1e-12 {
			t.Fatalf(
				"arm=%s RMS=%g want=%g",
				spec.name, doseOffsetRMS(spec.offsets), multiplicityDoseRMS,
			)
		}
	}
}

func TestUP32AnonymousSplitBasisIsOrthogonal(t *testing.T) {
	matrix := denseChannelMixer()
	if errValue := mixerOrthogonalityError(matrix); errValue > 1e-12 {
		t.Fatalf("anonymous channel mixer orthogonality error=%g", errValue)
	}
}

func TestUP32AllDoseTransportsRemainRealOrthogonallyEquivalent(t *testing.T) {
	mixer := fullLatentMixer()
	for _, spec := range multiplicityDoseOffsets() {
		step, err := fullLatentDoseStep(mixer, spec.offsets)
		if err != nil { t.Fatal(err) }
		realStep, err := realifyMatrix(step)
		if err != nil { t.Fatal(err) }
		errValue, err := maxRealOrthogonalityError(realStep)
		if err != nil { t.Fatal(err) }
		if errValue > 1e-10 {
			t.Fatalf("arm=%s orthogonality error=%g", spec.name, errValue)
		}
	}
}
