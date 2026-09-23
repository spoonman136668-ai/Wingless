package unitary

import (
	"math"
	"reflect"
	"testing"
)

func TestUP4SequentialComposition(t *testing.T) {
	result, err := RunUP4()
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != SequentialProbeSchema {
		t.Fatalf("schema=%q", result.Schema)
	}
	if result.Dimension != 4 || result.EventTypes != 3 || result.TrainSamples != 12 || result.HeldOutSamples != 96 {
		t.Fatalf("unexpected experiment shape: %#v", result)
	}
	if !reflect.DeepEqual(result.HeldOutLengths, []int{4, 8, 16, 32, 64, 128}) {
		t.Fatalf("held-out lengths=%v", result.HeldOutLengths)
	}
	if result.OrderSensitiveAB == result.OrderSensitiveBA {
		t.Fatalf("task is not order-sensitive: AB=%d BA=%d", result.OrderSensitiveAB, result.OrderSensitiveBA)
	}

	if result.Unitary.InitialTrainAccuracy != 0.5 {
		t.Fatalf("unitary initial train accuracy=%v want=0.5", result.Unitary.InitialTrainAccuracy)
	}
	if result.Unitary.TrainAccuracy != 1 {
		t.Fatalf("unitary train accuracy=%v want=1", result.Unitary.TrainAccuracy)
	}
	if result.Unitary.HeldOutAccuracy != 1 {
		t.Fatalf("unitary held-out accuracy=%v want=1", result.Unitary.HeldOutAccuracy)
	}
	if result.Unitary.FinalLoss > 1e-4 {
		t.Fatalf("unitary final loss=%g want<=1e-4", result.Unitary.FinalLoss)
	}
	for i, parameter := range result.Unitary.Parameters {
		if math.Abs(parameter-math.Pi/2) > 1e-2 {
			t.Fatalf("unitary parameter[%d]=%g target=%g", i, parameter, math.Pi/2)
		}
	}
	if result.Unitary.MaxNormDrift > 1e-12 {
		t.Fatalf("unitary norm drift=%g", result.Unitary.MaxNormDrift)
	}
	if result.Unitary.MaxRoundTripError > 1e-11 {
		t.Fatalf("unitary round-trip error=%g", result.Unitary.MaxRoundTripError)
	}

	if result.NonUnitary.InitialTrainAccuracy != 0.5 {
		t.Fatalf("control initial train accuracy=%v want=0.5", result.NonUnitary.InitialTrainAccuracy)
	}
	if result.NonUnitary.TrainAccuracy != 1 {
		t.Fatalf("control train accuracy=%v want=1", result.NonUnitary.TrainAccuracy)
	}
	for label, value := range map[string]float64{
		"held-out accuracy": result.NonUnitary.HeldOutAccuracy,
		"final loss": result.NonUnitary.FinalLoss,
		"norm drift": result.NonUnitary.MaxNormDrift,
	} {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			t.Fatalf("control %s is non-finite: %g", label, value)
		}
	}
}

func TestUP4TrainingLossFallsAtCheckpoints(t *testing.T) {
	result, err := RunUP4()
	if err != nil {
		t.Fatal(err)
	}
	previous := result.Unitary.InitialLoss
	for _, checkpoint := range result.Unitary.Checkpoints {
		if checkpoint.Loss > previous+1e-15 {
			t.Fatalf("unitary loss increased at step %d: previous=%g current=%g", checkpoint.Step, previous, checkpoint.Loss)
		}
		previous = checkpoint.Loss
	}
}

func TestUP4IsDeterministic(t *testing.T) {
	a, err := RunUP4()
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunUP4()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("UP-4 is nondeterministic:\nA=%#v\nB=%#v", a, b)
	}
}
