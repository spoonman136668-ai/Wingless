package unitary

import (
	"math"
	"testing"
)

func TestUP18OracleTransposeExactlyRecoversDenseMixer(t *testing.T) {
	mixer := denseChannelMixer()
	inverse := transposeChannelMatrix(mixer)
	product := channelMatrixProduct(inverse, mixer)

	maxError := 0.0
	for row := 0; row < 6; row++ {
		for column := 0; column < 6; column++ {
			want := 0.0
			if row == column {
				want = 1
			}
			delta := math.Abs(product[row][column] - want)
			if delta > maxError {
				maxError = delta
			}
		}
	}
	if maxError > 1e-12 {
		t.Fatalf("oracle inverse error=%g want<=1e-12", maxError)
	}
}

func TestUP18StructuredDemixerRemainsOrthogonal(t *testing.T) {
	var angles GivensAngles
	for i := range angles {
		angles[i] = -0.31 + 0.047*float64(i)
	}
	matrix := orthogonalDemixer(angles)
	if err := mixerOrthogonalityError(matrix); err > 1e-12 {
		t.Fatalf("structured demixer orthogonality error=%g", err)
	}
}

func TestUP18OracleAngleFamilyContainsExactInverse(t *testing.T) {
	got := orthogonalDemixer(oracleInverseAngles())
	want := transposeChannelMatrix(denseChannelMixer())
	if distance := channelMatrixL2(got, want); distance > 1e-12 {
		t.Fatalf("oracle-angle matrix distance=%g want<=1e-12", distance)
	}
}

func TestUP18OracleRecoversObservableFeatures(t *testing.T) {
	errValue, err := oracleFeatureRecoveryError()
	if err != nil {
		t.Fatal(err)
	}
	if errValue > 1e-12 {
		t.Fatalf("oracle feature recovery error=%g want<=1e-12", errValue)
	}
}
