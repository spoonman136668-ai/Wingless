package unitary

import (
	"math"
	"reflect"
	"testing"
)

func TestUP2MultilayerPhaseRouting(t *testing.T) {
	result, err := RunUP2()
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != MultilayerProbeSchema {
		t.Fatalf("schema=%q", result.Schema)
	}
	if result.Dimension != 4 || result.Classes != 4 || result.Couplings != 4 {
		t.Fatalf("unexpected shape: dimension=%d classes=%d couplings=%d", result.Dimension, result.Classes, result.Couplings)
	}
	if result.Unitary.InitialAccuracy != 0.25 {
		t.Fatalf("unitary initial accuracy=%v want=0.25", result.Unitary.InitialAccuracy)
	}
	if result.Unitary.TrainAccuracy != 1 || result.Unitary.HeldOutAccuracy != 1 {
		t.Fatalf("unitary accuracy train=%v held=%v", result.Unitary.TrainAccuracy, result.Unitary.HeldOutAccuracy)
	}
	if result.NonUnitary.InitialAccuracy != 0.25 {
		t.Fatalf("non-unitary initial accuracy=%v want=0.25", result.NonUnitary.InitialAccuracy)
	}
	if result.NonUnitary.TrainAccuracy != 1 || result.NonUnitary.HeldOutAccuracy != 1 {
		t.Fatalf("non-unitary accuracy train=%v held=%v", result.NonUnitary.TrainAccuracy, result.NonUnitary.HeldOutAccuracy)
	}
	if result.Unitary.FinalLoss > 1e-6 {
		t.Fatalf("unitary final loss=%g want<=1e-6", result.Unitary.FinalLoss)
	}
	for i, theta := range result.Unitary.Parameters {
		if math.Abs(theta+math.Pi/4) > 1e-3 {
			t.Fatalf("unitary theta[%d]=%g target=%g", i, theta, -math.Pi/4)
		}
	}
	if result.Unitary.MaxNormDrift > 1e-12 {
		t.Fatalf("unitary norm drift=%g", result.Unitary.MaxNormDrift)
	}
	if result.UnitaryMaxRoundTripError > 1e-12 {
		t.Fatalf("unitary round-trip error=%g", result.UnitaryMaxRoundTripError)
	}
	if !finite(result.NonUnitary.FinalLoss) || !finite(result.NonUnitary.MaxNormDrift) {
		t.Fatalf("non-unitary control produced non-finite metrics: %#v", result.NonUnitary)
	}
}

func TestUP2MatchedParameterCountAndOptimizerBudget(t *testing.T) {
	result, err := RunUP2()
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Unitary.Parameters) != len(result.NonUnitary.Parameters) {
		t.Fatalf("parameter count mismatch: unitary=%d non-unitary=%d", len(result.Unitary.Parameters), len(result.NonUnitary.Parameters))
	}
	if result.Steps != 60 || result.LearningRate != 0.15 || result.GradientEpsilon != 1e-6 {
		t.Fatalf("unexpected optimizer budget: %#v", result)
	}
}

func TestUP2IsDeterministic(t *testing.T) {
	a, err := RunUP2()
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunUP2()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("UP-2 is nondeterministic:\nA=%#v\nB=%#v", a, b)
	}
}

func TestUP2UnitaryLossFallsAtEveryCheckpoint(t *testing.T) {
	result, err := RunUP2()
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
