package unitary

import (
	"math"
	"reflect"
	"testing"
)

func TestProximalFusionLambdaDerivedFromFrozenConstants(t *testing.T) {
	want := directOffsetLearningRate * directOffsetResourcePrice
	if proximalFusionLambda != want {
		t.Fatalf("lambda=%g want=%g", proximalFusionLambda, want)
	}
	if proximalFusionLambda != 0.00004 {
		t.Fatalf("lambda=%g want=0.00004", proximalFusionLambda)
	}
}

func TestCompleteGraphFusionProxCreatesExactPool(t *testing.T) {
	input := []float64{-0.075, -0.045, -0.015, 0.015, 0.04499, 0.04501}
	output, diagnosis, err :=
		completeGraphFusionProx(input, proximalFusionLambda)
	if err != nil {
		t.Fatal(err)
	}
	if !diagnosis.ExactFusionOccurred {
		t.Fatalf("expected exact fusion: %+v", diagnosis)
	}
	if diagnosis.ExactCapacity <= 6 {
		t.Fatalf("capacity=%d want>6", diagnosis.ExactCapacity)
	}
	if output[4] != output[5] {
		t.Fatalf("pair not exactly fused: %.17g %.17g", output[4], output[5])
	}
	if diagnosis.ProxObjectiveAfter > diagnosis.ProxObjectiveBefore+1e-15 {
		t.Fatalf(
			"prox objective increased before=%g after=%g",
			diagnosis.ProxObjectiveBefore,
			diagnosis.ProxObjectiveAfter,
		)
	}
	if diagnosis.MaximumMonotonicViolation > 1e-15 {
		t.Fatalf(
			"monotonic violation=%g",
			diagnosis.MaximumMonotonicViolation,
		)
	}
	if diagnosis.MaximumBlockMeanResidual > 1e-15 {
		t.Fatalf(
			"block residual=%g",
			diagnosis.MaximumBlockMeanResidual,
		)
	}
}

func TestCompleteGraphFusionProxDoesNotNeedStickyState(t *testing.T) {
	closeInput := []float64{-0.075, -0.045, -0.015, 0.015, 0.04499, 0.04501}
	fused, _, err := completeGraphFusionProx(
		closeInput, proximalFusionLambda,
	)
	if err != nil {
		t.Fatal(err)
	}
	if fused[4] != fused[5] {
		t.Fatal("setup did not fuse")
	}

	separatedInput := append([]float64(nil), fused...)
	separatedInput[4] -= 0.01
	separatedInput[5] += 0.01
	separated, _, err := completeGraphFusionProx(
		separatedInput, proximalFusionLambda,
	)
	if err != nil {
		t.Fatal(err)
	}
	if separated[4] == separated[5] {
		t.Fatal("prox retained fusion despite separated new evidence")
	}
}

func TestCompleteGraphFusionProxDeterministic(t *testing.T) {
	input := continuousFusionInitialOffsets()
	first, firstDiagnosis, err :=
		completeGraphFusionProx(input, proximalFusionLambda)
	if err != nil {
		t.Fatal(err)
	}
	second, secondDiagnosis, err :=
		completeGraphFusionProx(input, proximalFusionLambda)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("output nondeterministic")
	}
	if !reflect.DeepEqual(firstDiagnosis, secondDiagnosis) {
		t.Fatalf("diagnosis nondeterministic")
	}
	rms := doseOffsetRMS(first)
	if math.Abs(rms-multiplicityDoseRMS) > 1e-12 && rms != 0 {
		t.Fatalf("rms=%g want=%g or 0", rms, multiplicityDoseRMS)
	}
}
