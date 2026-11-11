package unitary

import (
	"math"
	"reflect"
	"testing"
)

func TestUP1LearnsUnitaryRouting(t *testing.T) {
	result, err := RunUP1()
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != TrainingProbeSchema {
		t.Fatalf("schema=%q", result.Schema)
	}
	if result.Steps != 20 {
		t.Fatalf("steps=%d want=20", result.Steps)
	}
	if result.InitialAccuracy != 0.5 {
		t.Fatalf("initial accuracy=%v want=0.5", result.InitialAccuracy)
	}
	if result.TrainAccuracy != 1 {
		t.Fatalf("train accuracy=%v want=1", result.TrainAccuracy)
	}
	if result.HeldOutAccuracy != 1 {
		t.Fatalf("held-out accuracy=%v want=1", result.HeldOutAccuracy)
	}
	if result.FinalLoss >= result.InitialLoss {
		t.Fatalf("loss did not improve: initial=%g final=%g", result.InitialLoss, result.FinalLoss)
	}
	if result.FinalLoss > 1e-5 {
		t.Fatalf("final loss=%g want<=1e-5", result.FinalLoss)
	}
	if math.Abs(result.LearnedTheta-math.Pi/4) > 1e-3 {
		t.Fatalf("learned theta=%g target=%g", result.LearnedTheta, math.Pi/4)
	}
	if result.MaxNormDrift > 1e-12 {
		t.Fatalf("norm drift=%g", result.MaxNormDrift)
	}
	if result.MaxRoundTripError > 1e-12 {
		t.Fatalf("round-trip error=%g", result.MaxRoundTripError)
	}
}

func TestUP1IsDeterministic(t *testing.T) {
	a, err := RunUP1()
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunUP1()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("training probe is nondeterministic:\nA=%#v\nB=%#v", a, b)
	}
}

func TestUP1LossFallsMonotonically(t *testing.T) {
	result, err := RunUP1()
	if err != nil {
		t.Fatal(err)
	}
	previous := result.InitialLoss
	for _, step := range result.Trace {
		if step.Loss > previous+1e-15 {
			t.Fatalf("loss increased at step %d: previous=%g current=%g", step.Step, previous, step.Loss)
		}
		previous = step.Loss
	}
}
