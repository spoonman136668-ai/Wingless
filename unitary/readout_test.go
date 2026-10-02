package unitary

import (
	"math"
	"reflect"
	"testing"
)

func TestUP6LearnedReadoutQuery(t *testing.T) {
	result, err := RunUP6()
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != LearnedReadoutProbeSchema {
		t.Fatalf("schema=%q", result.Schema)
	}
	if result.RuntimePrototypeLookup {
		t.Fatal("runtime prototype lookup must be disabled")
	}
	if result.Dimension != 16 || result.Entities != 4 || result.ValuesPerEntity != 4 {
		t.Fatalf("unexpected shape: %#v", result)
	}
	if result.Scenarios != 48 || result.WritesPerScenario != 16 {
		t.Fatalf("unexpected workload: %#v", result)
	}
	if !reflect.DeepEqual(result.TransportDepths, []int{32, 128, 512, 1024}) {
		t.Fatalf("transport depths=%v", result.TransportDepths)
	}
	if result.ValueDecoder.TrainAccuracy != 1 {
		t.Fatalf("value decoder train accuracy=%v", result.ValueDecoder.TrainAccuracy)
	}
	if result.RelationHead.TrainAccuracy != 1 {
		t.Fatalf("relation head train accuracy=%v", result.RelationHead.TrainAccuracy)
	}

	if result.Unitary.CommitDecodeAccuracy != 1 {
		t.Fatalf("unitary commit accuracy=%v", result.Unitary.CommitDecodeAccuracy)
	}
	if result.Unitary.ExactFinalTableAccuracy != 1 {
		t.Fatalf("unitary final-table accuracy=%v", result.Unitary.ExactFinalTableAccuracy)
	}
	if result.Unitary.RelationalQueryAccuracy != 1 {
		t.Fatalf("unitary relation accuracy=%v", result.Unitary.RelationalQueryAccuracy)
	}
	if result.Unitary.MaxRoundTripError > 1e-10 {
		t.Fatalf("unitary round-trip error=%g", result.Unitary.MaxRoundTripError)
	}
	if result.Unitary.MaxForwardNormDrift > 1e-12 {
		t.Fatalf("unitary forward norm drift=%g", result.Unitary.MaxForwardNormDrift)
	}
	if result.Unitary.MinValueMargin <= 0.5 {
		t.Fatalf("unitary value margin=%g want>0.5", result.Unitary.MinValueMargin)
	}
	if result.Unitary.MinRelationMargin <= 0.5 {
		t.Fatalf("unitary relation margin=%g want>0.5", result.Unitary.MinRelationMargin)
	}

	for label, value := range map[string]float64{
		"control commit accuracy": result.NonUnitary.CommitDecodeAccuracy,
		"control final accuracy": result.NonUnitary.ExactFinalTableAccuracy,
		"control relation accuracy": result.NonUnitary.RelationalQueryAccuracy,
		"control round-trip": result.NonUnitary.MaxRoundTripError,
		"control forward norm": result.NonUnitary.MaxForwardNormDrift,
		"control value margin": result.NonUnitary.MinValueMargin,
		"control relation margin": result.NonUnitary.MinRelationMargin,
	} {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			t.Fatalf("%s is non-finite: %g", label, value)
		}
	}
}

func TestUP6HeldOutEntityPairsAreReserved(t *testing.T) {
	result, err := RunUP6()
	if err != nil {
		t.Fatal(err)
	}
	want := [][2]int{{0, 3}, {3, 0}, {2, 0}, {3, 1}}
	if !reflect.DeepEqual(result.HeldOutEntityPairs, want) {
		t.Fatalf("held-out entity pairs=%v", result.HeldOutEntityPairs)
	}
}

func TestUP6IsDeterministic(t *testing.T) {
	a, err := RunUP6()
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunUP6()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("UP-6 is nondeterministic:\nA=%#v\nB=%#v", a, b)
	}
}
