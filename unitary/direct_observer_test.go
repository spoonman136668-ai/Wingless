package unitary

import (
	"math"
	"reflect"
	"testing"
)

func TestUP7DirectObserverHarness(t *testing.T) {
	result, err := RunUP7()
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != DirectObserverProbeSchema {
		t.Fatalf("schema=%q", result.Schema)
	}
	if result.RuntimePrototypeLookup {
		t.Fatal("runtime prototype lookup must be disabled")
	}
	if result.ExplicitInverseReadout {
		t.Fatal("explicit inverse readout must be disabled")
	}
	if result.DepthFeatureProvided {
		t.Fatal("transport depth must not be supplied to observer")
	}
	if result.Dimension != 16 || result.Entities != 4 || result.ValuesPerEntity != 4 {
		t.Fatalf("unexpected shape: %#v", result)
	}
	if result.TrainTables != 128 || result.HeldOutTables != 128 {
		t.Fatalf("unexpected table split: train=%d held=%d", result.TrainTables, result.HeldOutTables)
	}
	if !reflect.DeepEqual(result.TrainDepths, []int{8, 24, 72, 216, 432, 648}) {
		t.Fatalf("train depths=%v", result.TrainDepths)
	}
	if !reflect.DeepEqual(result.HeldOutDepths, []int{32, 128, 512, 1024}) {
		t.Fatalf("held depths=%v", result.HeldOutDepths)
	}
	if result.RelationHead.TrainAccuracy != 1 {
		t.Fatalf("relation head train accuracy=%v", result.RelationHead.TrainAccuracy)
	}

	for label, value := range map[string]float64{
		"unitary observer train accuracy": result.UnitaryObserver.TrainAccuracy,
		"unitary observer train loss": result.UnitaryObserver.TrainLoss,
		"unitary static held accuracy": result.UnitaryObserver.StaticHeldOutAcc,
		"unitary commit accuracy": result.Unitary.CommitDecodeAccuracy,
		"unitary final accuracy": result.Unitary.ExactFinalTableAccuracy,
		"unitary relation accuracy": result.Unitary.RelationalQueryAccuracy,
		"unitary value margin": result.Unitary.MinValueMargin,
		"unitary relation margin": result.Unitary.MinRelationMargin,
		"unitary norm drift": result.Unitary.MaxForwardNormDrift,
		"control observer train accuracy": result.ControlObserver.TrainAccuracy,
		"control observer train loss": result.ControlObserver.TrainLoss,
		"control static held accuracy": result.ControlObserver.StaticHeldOutAcc,
		"control commit accuracy": result.NonUnitary.CommitDecodeAccuracy,
		"control final accuracy": result.NonUnitary.ExactFinalTableAccuracy,
		"control relation accuracy": result.NonUnitary.RelationalQueryAccuracy,
		"control value margin": result.NonUnitary.MinValueMargin,
		"control relation margin": result.NonUnitary.MinRelationMargin,
		"control norm drift": result.NonUnitary.MaxForwardNormDrift,
	} {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			t.Fatalf("%s is non-finite: %g", label, value)
		}
	}
}

func TestUP7TableSplitIsDisjoint(t *testing.T) {
	for _, table := range allMemoryTables() {
		index := memoryTableIndex(table)
		if directTrainTable(table) != (index%2 == 0) {
			t.Fatalf("unexpected table split at index=%d", index)
		}
	}
}

func TestUP7IsDeterministic(t *testing.T) {
	a, err := RunUP7()
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunUP7()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("UP-7 is nondeterministic:\nA=%#v\nB=%#v", a, b)
	}
}
