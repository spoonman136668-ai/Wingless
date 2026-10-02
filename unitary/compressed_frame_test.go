package unitary

import (
	"math"
	"testing"
)

func TestUP10CompressedFrameEncodesEntityValues(t *testing.T) {
	bank, err := makeCompressedFrameBank()
	if err != nil {
		t.Fatal(err)
	}
	for entity := 0; entity < 4; entity++ {
		previous := math.Inf(-1)
		for value := 0; value < 4; value++ {
			table := memoryTable{0, 0, 0, 0}
			table[entity] = value
			memory, err := encodeMemory(table)
			if err != nil {
				t.Fatal(err)
			}
			features, err := compressedFrameFeature(memory, bank, entity)
			if err != nil {
				t.Fatal(err)
			}
			if len(features) != 1 {
				t.Fatalf("entity=%d feature count=%d want=1", entity, len(features))
			}
			if features[0] <= previous {
				t.Fatalf("entity=%d value=%d feature=%g did not increase from %g", entity, value, features[0], previous)
			}
			previous = features[0]
		}
	}
}

func TestUP10CompressedFrameIsGlobalPhaseInvariant(t *testing.T) {
	bank, err := makeCompressedFrameBank()
	if err != nil {
		t.Fatal(err)
	}
	memory, err := encodeMemory(memoryTable{3, 2, 1, 0})
	if err != nil {
		t.Fatal(err)
	}
	rotated := rotateGlobalPhase(memory, 1.337)
	for entity := 0; entity < 4; entity++ {
		a, err := compressedFrameFeature(memory, bank, entity)
		if err != nil {
			t.Fatal(err)
		}
		b, err := compressedFrameFeature(rotated, bank, entity)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(a[0]-b[0]) > 1e-12 {
			t.Fatalf("entity=%d global phase changed feature %g vs %g", entity, a[0], b[0])
		}
	}
}

func TestUP10UnitaryCoevolutionPreservesCompressedFeature(t *testing.T) {
	bank, err := makeCompressedFrameBank()
	if err != nil {
		t.Fatal(err)
	}
	memory, err := encodeMemory(memoryTable{1, 3, 0, 2})
	if err != nil {
		t.Fatal(err)
	}
	memory = rotateGlobalPhase(memory, 0.771)
	block := stressProgram()
	evolvedMemory, err := applyStressUnitary(memory, block, 512)
	if err != nil {
		t.Fatal(err)
	}
	evolvedBank, err := transportCompressedFrame(bank, block, 512, applyStressUnitary, 1)
	if err != nil {
		t.Fatal(err)
	}
	for entity := 0; entity < 4; entity++ {
		before, err := compressedFrameFeature(memory, bank, entity)
		if err != nil {
			t.Fatal(err)
		}
		after, err := compressedFrameFeature(evolvedMemory, evolvedBank, entity)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(before[0]-after[0]) > 1e-12 {
			t.Fatalf("entity=%d coevolution changed feature %g vs %g", entity, before[0], after[0])
		}
	}
}
