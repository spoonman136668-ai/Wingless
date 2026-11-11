package unitary

import (
	"math"
	"reflect"
	"testing"
)

func TestUP5RelationalReadWriteMemory(t *testing.T) {
	result, err := RunUP5()
	if err != nil {
		t.Fatal(err)
	}
	if result.Schema != MemoryProbeSchema {
		t.Fatalf("schema=%q", result.Schema)
	}
	if result.Entities != 4 || result.ValuesPerEntity != 4 || result.Dimension != 16 {
		t.Fatalf("unexpected memory shape: %#v", result)
	}
	if result.Scenarios != 32 || result.WritesPerScenario != 12 || result.CommitDecodesPerPath != 384 {
		t.Fatalf("unexpected workload: %#v", result)
	}
	if !reflect.DeepEqual(result.TransportDepths, []int{8, 32, 128, 512}) {
		t.Fatalf("transport depths=%v", result.TransportDepths)
	}
	if result.WriteBoundary != "observe_decode_overwrite_reencode" {
		t.Fatalf("write boundary=%q", result.WriteBoundary)
	}

	if result.Unitary.CommitDecodeAccuracy != 1 {
		t.Fatalf("unitary commit accuracy=%v want=1", result.Unitary.CommitDecodeAccuracy)
	}
	if result.Unitary.ExactFinalTableAccuracy != 1 {
		t.Fatalf("unitary final-table accuracy=%v want=1", result.Unitary.ExactFinalTableAccuracy)
	}
	if result.Unitary.RelationalQueryAccuracy != 1 {
		t.Fatalf("unitary relation accuracy=%v want=1", result.Unitary.RelationalQueryAccuracy)
	}
	if result.Unitary.MinDecodeMargin <= 0.1 {
		t.Fatalf("unitary decode margin=%g want>0.1", result.Unitary.MinDecodeMargin)
	}
	if result.Unitary.MaxNormDrift > 1e-12 {
		t.Fatalf("unitary norm drift=%g", result.Unitary.MaxNormDrift)
	}

	for label, value := range map[string]float64{
		"control commit accuracy": result.NonUnitary.CommitDecodeAccuracy,
		"control final accuracy": result.NonUnitary.ExactFinalTableAccuracy,
		"control relation accuracy": result.NonUnitary.RelationalQueryAccuracy,
		"control decode margin": result.NonUnitary.MinDecodeMargin,
		"control norm drift": result.NonUnitary.MaxNormDrift,
	} {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			t.Fatalf("%s is non-finite: %g", label, value)
		}
	}
}

func TestUP5OverwriteIsActuallyIrreversible(t *testing.T) {
	original := memoryTable{0, 1, 2, 3}
	a, err := applyMemoryWrite(original, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	b, err := applyMemoryWrite(memoryTable{0, 2, 2, 3}, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatalf("overwrite boundary unexpectedly retained displaced value: a=%v b=%v", a, b)
	}
}

func TestUP5IsDeterministic(t *testing.T) {
	a, err := RunUP5()
	if err != nil {
		t.Fatal(err)
	}
	b, err := RunUP5()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("UP-5 is nondeterministic:\nA=%#v\nB=%#v", a, b)
	}
}
