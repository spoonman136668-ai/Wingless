package unitary

import (
	"math"
	"testing"
)

func TestUP17DenseMixerIsOrthogonalAndDense(t *testing.T) {
	mixer := denseChannelMixer()
	if err := mixerOrthogonalityError(mixer); err > 1e-12 {
		t.Fatalf("mixer orthogonality error=%g", err)
	}
	for row := 0; row < 6; row++ {
		for column := 0; column < 6; column++ {
			if math.Abs(mixer[row][column]) < 1e-6 {
				t.Fatalf(
					"mixer row=%d col=%d is not dense: %g",
					row, column, mixer[row][column],
				)
			}
		}
	}
}

func TestUP17AnonymousMixingPreservesCompositeNorm(t *testing.T) {
	memory, err := encodeMemory(memoryTable{0, 1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	mixed, err := blindMixedComposite(memory)
	if err != nil {
		t.Fatal(err)
	}
	norm2, err := NormSquared(mixed)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(norm2-1) > 1e-12 {
		t.Fatalf("mixed composite norm=%g want=1", norm2)
	}
}

func TestUP17IdentityReadoutDoesNotUndoDenseMixer(t *testing.T) {
	composed := channelMatrixProduct(
		identityChannelMatrix(),
		denseChannelMixer(),
	)
	alignments := roleAlignments(
		identityChannelMatrix(),
		denseChannelMixer(),
	)
	mean, _ := blindMeanMin(alignments)
	if mean >= 0.50 {
		t.Fatalf(
			"identity readout mean role alignment=%g unexpectedly high; matrix=%v",
			mean, composed,
		)
	}
}
