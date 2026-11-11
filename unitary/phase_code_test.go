package unitary

import (
	"math"
	"testing"
)

func TestUP21LearnedPhaseTargetsAreUnitCircle(t *testing.T) {
	for value := 0; value < 4; value++ {
		target, err := learnedPhaseTarget(value)
		if err != nil {
			t.Fatal(err)
		}
		norm := math.Hypot(target[0], target[1])
		if math.Abs(norm-1) > 1e-12 {
			t.Fatalf("value=%d phase target norm=%g want=1", value, norm)
		}
	}
}

func TestUP21PhaseRegressorRuntimeIsPrimalOnly(t *testing.T) {
	rows, err := buildRidgeAnonymousRows(
		selectFirstMemoryTables(allMemoryTables(), 8),
		0.05,
		1,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	regressors, err := trainPhaseCodeRegressors(
		rows, ridgeAnonymousLambda,
	)
	if err != nil {
		t.Fatal(err)
	}
	for entity := 0; entity < 4; entity++ {
		for output := 0; output < 2; output++ {
			if len(regressors[entity].weights[output]) !=
				anonymousGramQuadraticDim {
				t.Fatalf(
					"entity=%d output=%d dim=%d want=%d",
					entity,
					output,
					len(regressors[entity].weights[output]),
					anonymousGramQuadraticDim,
				)
			}
		}
	}
}
