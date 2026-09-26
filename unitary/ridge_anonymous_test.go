package unitary

import (
	"math"
	"testing"
)

func TestUP20DenseMultipleSolve(t *testing.T) {
	matrix := [][]float64{
		{3, 1},
		{1, 2},
	}
	right := [][]float64{
		{9, 1},
		{8, 0},
	}
	got, err := solveDenseMultiple(matrix, right)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]float64{
		{2, 0.4},
		{3, -0.2},
	}
	for row := range want {
		for column := range want[row] {
			if math.Abs(got[row][column]-want[row][column]) > 1e-12 {
				t.Fatalf(
					"solve[%d][%d]=%g want=%g",
					row, column, got[row][column], want[row][column],
				)
			}
		}
	}
}

func TestUP20OracleNumeratorLivesInQuadraticAnonymousSpan(t *testing.T) {
	errValue, err := oracleNumeratorWitnessError()
	if err != nil {
		t.Fatal(err)
	}
	if errValue > 1e-12 {
		t.Fatalf("oracle quadratic witness error=%g want<=1e-12", errValue)
	}
}

func TestUP20RidgeRuntimeDoesNotNeedSupportVectors(t *testing.T) {
	rows, err := buildRidgeAnonymousRows(
		selectFirstMemoryTables(allMemoryTables(), 8),
		0.05,
		1,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	heads, _, err := trainRidgeAnonymousHeads(rows, ridgeAnonymousLambda)
	if err != nil {
		t.Fatal(err)
	}
	for entity := 0; entity < 4; entity++ {
		if len(heads[entity].weights) != 4 {
			t.Fatalf("entity=%d class count=%d", entity, len(heads[entity].weights))
		}
		if len(heads[entity].weights[0]) != anonymousGramQuadraticDim {
			t.Fatalf(
				"entity=%d runtime feature dim=%d want=%d",
				entity, len(heads[entity].weights[0]), anonymousGramQuadraticDim,
			)
		}
	}
}

func selectFirstMemoryTables(tables []memoryTable, count int) []memoryTable {
	if count > len(tables) {
		count = len(tables)
	}
	return append([]memoryTable(nil), tables[:count]...)
}
