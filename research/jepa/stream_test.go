package jepa

import (
	"math"
	"testing"
)

func dataset(from, count int) [][]Pair {
	rows := make([][]Pair, count)
	for i := range rows {
		v, err := DevelopmentSequence(uint64(from+i), 40)
		if err != nil {
			panic(err)
		}
		rows[i] = v
	}
	return rows
}

func mse(model Recurrent, sequences [][]Pair) float64 {
	sum := 0.0
	for _, seq := range sequences {
		v, err := model.Loss(seq)
		if err != nil {
			panic(err)
		}
		sum += v
	}
	return sum / float64(len(sequences))
}

func TestGoodViewsLearnTemporalStateBeyondMemorylessControl(t *testing.T) {
	train, holdout := dataset(10, 36), dataset(10000, 12)
	learned, fixed := New(), New()
	initial := mse(learned, holdout)
	if err := learned.Fit(train, 250, .16, false); err != nil {
		t.Fatal(err)
	}
	if err := fixed.Fit(train, 250, .16, true); err != nil {
		t.Fatal(err)
	}
	if after := mse(learned, holdout); after >= initial || after >= mse(fixed, holdout) {
		t.Fatalf("good-view temporal learning failed: initial=%f after=%f memoryless=%f", initial, after, mse(fixed, holdout))
	}
	if learned.Alpha <= .15 {
		t.Fatalf("recurrent memory did not learn: %+v", learned)
	}
}

func TestMispairedNegativeControl(t *testing.T) {
	train, holdout := dataset(10, 36), dataset(10000, 12)
	good, bad := New(), New()
	corrupt := make([][]Pair, len(train))
	for i, seq := range train {
		corrupt[i] = make([]Pair, len(seq))
		for j, p := range seq {
			corrupt[i][j] = Pair{Observation: p.Observation,
				OtherView: train[(i+len(train)/2)%len(train)][j].OtherView}
		}
	}
	if err := good.Fit(train, 250, .16, false); err != nil {
		t.Fatal(err)
	}
	if err := bad.Fit(corrupt, 250, .16, false); err != nil {
		t.Fatal(err)
	}
	if mse(good, holdout) >= mse(bad, holdout) {
		t.Fatal("mismatched good-view pairing went undetected")
	}
}

func TestColdRetentionCheckpointIsExactAndBounded(t *testing.T) {
	learned := New()
	if err := learned.Fit(dataset(10, 20), 100, .16, false); err != nil {
		t.Fatal(err)
	}
	raw, err := learned.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) != SnapshotSize {
		t.Fatalf("checkpoint bytes=%d", len(raw))
	}
	restored, err := Restore(raw)
	if err != nil || restored != learned {
		t.Fatalf("cold retention mismatch: %+v %v", restored, err)
	}
	score := mse(learned, dataset(5000, 5))
	if got := mse(restored, dataset(5000, 5)); got != score {
		t.Fatalf("cold outcome mismatch: %.12f != %.12f", got, score)
	}
	raw[13] ^= 0xff
	if _, err := Restore(raw); err == nil {
		t.Fatal("tampered checkpoint accepted")
	}
}

func TestInvalidViewsAndTrainingBudgetFailClosed(t *testing.T) {
	if _, err := DevelopmentSequence(1, MaxSteps+1); err == nil {
		t.Fatal("oversized sample count accepted")
	}
	rows := dataset(0, 2)
	m := New()
	if err := m.Fit(rows, MaxEpochs+1, .1, false); err == nil {
		t.Fatal("epoch budget bypass")
	}
	if err := m.Fit(rows, 10, math.NaN(), false); err == nil {
		t.Fatal("nan learning rate accepted")
	}
	rows[0][0].OtherView = math.Inf(1)
	if err := m.Fit(rows, 1, .1, false); err == nil {
		t.Fatal("nonfinite view accepted")
	}
	if _, err := Restore(nil); err == nil {
		t.Fatal("empty checkpoint accepted")
	}
}
