package r53prefreeze

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

var frozen = []struct {
	id ArmID
	groups [6]int
	parameters int
}{
	{A0, [6]int{0, 1, 2, 3, 4, 5}, 7},
	{A1, [6]int{0, 0, 1, 2, 3, 4}, 6},
	{A2, [6]int{0, 1, 2, 2, 2, 2}, 4},
	{A3, [6]int{0, 0, 1, 1, 1, 1}, 3},
}

func TestFourFrozenSlopeBlocks(t *testing.T) {
	if got := Arms(); got != [4]ArmID{A0, A1, A2, A3} {
		t.Fatalf("unexpected arm set: %v", got)
	}
	for _, want := range frozen {
		got, count, err := Groups(want.id)
		if err != nil || got != want.groups || count != want.parameters {
			t.Fatalf("arm %s: groups %v count %d err %v", want.id, got, count, err)
		}
	}
}
func TestSingleInterceptAndUnchangedResourceCeiling(t *testing.T) {
	if StateDimension != 6 || MaxReadoutParameters != 7 ||
		AdaptationPerArm != 436 || TotalAdaptationUnits != 4*AdaptationPerArm ||
		ExternalModelCalls != 0 || TokenizerCalls != 0 {
		t.Fatal("frozen scientific resource ceilings changed")
	}
	x := [6]float64{0, 0, 0, 0, 0, 0}
	for _, arm := range Arms() {
		row, err := DesignRow(arm, x)
		if err != nil || row != [7]float64{1, 0, 0, 0, 0, 0, 0} {
			t.Fatalf("%s multiple intercepts: %v %v", arm, row, err)
		}
	}
}
func TestImmutableSelectedStateFormerAlpha(t *testing.T) {
	fixed := [6]float64{0.75, 0.75, 0.125, 0.125, 1, 1}
	if SelectedAlpha() != fixed { t.Fatal("R48-selected alpha changed") }
	copy := SelectedAlpha()
	copy[0] = 99
	if SelectedAlpha() != fixed { t.Fatal("alpha configuration mutable") }
}
func TestR50HistoricalDesignRowAlgebra(t *testing.T) {
	// Frozen R50 source: z[0]=1, then z[groups[i]+1]+=e.x[i].
	input := [6]float64{3, -2, 7, 11, 13, -5}
	want := map[ArmID][7]float64{
		A0: {1, 3, -2, 7, 11, 13, -5},
		A1: {1, 1, 7, 11, 13, -5, 0},
		A2: {1, 3, -2, 26, 0, 0, 0},
		A3: {1, 1, 26, 0, 0, 0, 0},
	}
	for _, arm := range Arms() {
		got, err := DesignRow(arm, input)
		if err != nil || got != want[arm] {
			t.Fatalf("%s R50 feature parity mismatch %v want %v err %v",
				arm, got, want[arm], err)
		}
	}
}
func TestEarlyPairTieOnly(t *testing.T) {
	x := [6]float64{1, 2, 3, 4, 5, 6}
	a, _ := DesignRow(A1, x)
	if a != [7]float64{1, 3, 3, 4, 5, 6, 0} {
		t.Fatalf("early-only tie changed: %v", a)
	}
	g, _, _ := Groups(A1)
	if g[2] == g[3] || g[3] == g[4] || g[4] == g[5] {
		t.Fatal("late slopes unexpectedly tied")
	}
}
func TestLateFourTieOnly(t *testing.T) {
	x := [6]float64{1, 2, 3, 4, 5, 6}
	a, _ := DesignRow(A2, x)
	if a != [7]float64{1, 1, 2, 18, 0, 0, 0} {
		t.Fatalf("late-only tie changed: %v", a)
	}
	g, _, _ := Groups(A2)
	if g[0] == g[1] || g[1] == g[2] {
		t.Fatal("early slopes unexpectedly tied")
	}
}
func TestParentR52TwoGroupsExact(t *testing.T) {
	x := [6]float64{1, 2, 3, 4, 5, 6}
	a, _ := DesignRow(A3, x)
	if a != [7]float64{1, 3, 18, 0, 0, 0, 0} {
		t.Fatalf("R52 selected-two-group2 comparator changed: %v", a)
	}
}
func TestCoefficientExpansionMatchesHistoricalR50(t *testing.T) {
	theta := map[ArmID][7]float64{
		A0: {2, 3, 5, 7, 11, 13, 17},
		A1: {2, 3, 5, 7, 11, 13, 0},
		A2: {2, 3, 5, 7, 0, 0, 0},
		A3: {2, 3, 5, 0, 0, 0, 0},
	}
	want := map[ArmID][7]float64{
		A0: {2, 3, 5, 7, 11, 13, 17},
		A1: {2, 3, 3, 5, 7, 11, 13},
		A2: {2, 3, 5, 7, 7, 7, 7},
		A3: {2, 3, 3, 5, 5, 5, 5},
	}
	for _, arm := range Arms() {
		got, err := ExpandReadout(arm, theta[arm])
		if err != nil || got != want[arm] {
			t.Fatalf("%s beta parity: %v want %v err %v", arm, got, want[arm], err)
		}
	}
}
func TestCompressedExpandedPredictParity(t *testing.T) {
	inputs := [][6]float64{
		{1, 2, 3, 4, 5, 6},
		{-3.5, 0.125, 6, 17, -8, 0},
		{0, 0, 0, 0, 0, 0},
		{31, -2, -12, 1, 0, 0.5},
	}
	for _, arm := range Arms() {
		_, n, err := Groups(arm)
		if err != nil { t.Fatal(err) }
		theta := [7]float64{}
		for i := 0; i < n; i++ { theta[i] = float64(i+1)/7.0 }
		for _, x := range inputs {
			expanded, err := ExpandReadout(arm, theta)
			if err != nil { t.Fatal(err) }
			row, err := DesignRow(arm, x)
			if err != nil { t.Fatal(err) }
			expected := 0.0
			for i := 0; i < n; i++ { expected += row[i]*theta[i] }
			uncompressed := expanded[0]
			for i, v := range x { uncompressed += v*expanded[i+1] }
			result, err := Predict(arm, x, theta)
			if err != nil || math.Abs(result-expected) > 1e-11 ||
				math.Abs(result-uncompressed) > 1e-11 {
				t.Fatalf("%s prediction diverges from original R50 algebra: result=%v compressed=%v expanded=%v err=%v", arm, result, expected, uncompressed, err)
			}
		}
	}
}
func TestFrozenGroupsAreCopiedNotMutable(t *testing.T) {
	original, _, _ := Groups(A2)
	modified, _, _ := Groups(A2)
	modified[2] = 99
	next, _, _ := Groups(A2)
	if original != next { t.Fatal("mutable global group mapping") }
}
func TestRejectUnpreregisteredArm(t *testing.T) {
	if _, _, err := Groups("A4"); !errors.Is(err, ErrUnfrozenArm) { t.Fatal(err) }
	if _, err := DesignRow("A4", [6]float64{}); !errors.Is(err, ErrUnfrozenArm) { t.Fatal(err) }
	if _, err := ExpandReadout("A4", [7]float64{}); !errors.Is(err, ErrUnfrozenArm) { t.Fatal(err) }
	if _, err := Predict("A4", [6]float64{}, [7]float64{}); !errors.Is(err, ErrUnfrozenArm) { t.Fatal(err) }
}
func TestFailClosedNonfiniteInputAndCoefficients(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		x := [6]float64{1, 2, 3, 4, 5, bad}
		if _, err := DesignRow(A1, x); !errors.Is(err, ErrNonfiniteState) {
			t.Fatalf("nonfinite state accepted: %v", err)
		}
		beta := [7]float64{0, 1, 0, 2, 3, bad}
		if _, err := ExpandReadout(A1, beta); !errors.Is(err, ErrNonfiniteCoefficient) {
			t.Fatalf("nonfinite theta accepted: %v", err)
		}
	}
}
func TestRejectUnbudgetedUnusedCoefficients(t *testing.T) {
	for _, arm := range []ArmID{A1, A2, A3} {
		_, n, _ := Groups(arm)
		theta := [7]float64{}
		theta[n] = 0.125
		if _, err := ExpandReadout(arm, theta); !errors.Is(err, ErrHiddenParameter) {
			t.Fatalf("%s hidden coefficient accepted: %v", arm, err)
		}
	}
}
func TestStateReadoutOverflowFailsClosed(t *testing.T) {
	x := [6]float64{1e308, 1e308, 0, 0, 0, 0}
	if _, err := DesignRow(A1, x); !errors.Is(err, ErrNonfiniteState) { t.Fatal(err) }
	theta := [7]float64{1e308, 1e308, 0, 0, 0, 0, 0}
	if _, err := Predict(A3, [6]float64{10, 10, 0, 0, 0, 0}, theta); !errors.Is(err, ErrNonfiniteCoefficient) {
		t.Fatalf("infinite score accepted: %v", err)
	}
}
func TestNoTrainingOrSourceStateMutated(t *testing.T) {
	x := [6]float64{1, 2, 3, 4, 5, 6}
	theta := [7]float64{2, 1, 2, 3, 4, 5, 0}
	srcX, srcTheta := x, theta
	_, err := Predict(A1, x, theta)
	if err != nil { t.Fatal(err) }
	if !reflect.DeepEqual(x, srcX) || !reflect.DeepEqual(theta, srcTheta) {
		t.Fatal("mutated caller data")
	}
}
