package unitary

import (
	"math"
	"testing"
)

func TestUP26BreadthFeatureDimensions(t *testing.T) {
	if got := breadthRawFeatureDim(32); got != 64 {
		t.Fatalf("raw32=%d want=64", got)
	}
	if got := breadthQuadraticFeatureDim(32); got != 2144 {
		t.Fatalf("quad32=%d want=2144", got)
	}
	if got := breadthRawFeatureDim(64); got != 128 {
		t.Fatalf("raw64=%d want=128", got)
	}
	if got := breadthQuadraticFeatureDim(64); got != 8384 {
		t.Fatalf("quad64=%d want=8384", got)
	}
}

func TestUP26First32DiscoveryIsPrefixInvariant(t *testing.T) {
	mixer := fullLatentMixer()
	step, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil {
		t.Fatal(err)
	}
	first32, _, err := discoverCommutingObservables(step, 32, 3)
	if err != nil {
		t.Fatal(err)
	}
	first64, _, err := discoverCommutingObservables(step, 64, 3)
	if err != nil {
		t.Fatal(err)
	}
	for observable := 0; observable < 32; observable++ {
		for row := 0; row < fullLatentDimension; row++ {
			for column := 0; column < fullLatentDimension; column++ {
				delta := math.Abs(real(
					first32[observable][row][column] -
						first64[observable][row][column],
				))
				deltaImag := math.Abs(imag(
					first32[observable][row][column] -
						first64[observable][row][column],
				))
				if delta > 1e-14 || deltaImag > 1e-14 {
					t.Fatalf(
						"prefix changed observable=%d row=%d col=%d",
						observable, row, column,
					)
				}
			}
		}
	}
}
