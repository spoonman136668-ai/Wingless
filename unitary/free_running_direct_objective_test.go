package unitary

import "testing"

func TestUP40OneHotSoftMemoryMatchesCanonicalWrite(t *testing.T) {
	table := memoryTable{0, 1, 2, 3}
	var distributions [4][]float64
	for entity, value := range table {
		distributions[entity] = []float64{0, 0, 0, 0}
		distributions[entity][value] = 1
	}
	got, err := softMemoryAfterCommandedWrite(distributions, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	wantTable, err := applyMemoryWrite(table, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	want, err := encodeMemory(wantTable)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("dimension=%d want=%d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("state[%d]=%v want=%v", i, got[i], want[i])
		}
	}
}

func TestUP40SoftMemoryPreservesUncertaintyAndExactWrite(t *testing.T) {
	distributions := [4][]float64{
		{0.5, 0.5, 0, 0},
		{0.1, 0.2, 0.3, 0.4},
		{0.25, 0.25, 0.25, 0.25},
		{0, 0, 0.75, 0.25},
	}
	state, err := softMemoryAfterCommandedWrite(distributions, 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if state[0] == 0 || state[1] == 0 {
		t.Fatal("unwritten entity uncertainty was collapsed")
	}
	for value := 0; value < 4; value++ {
		got := state[4+value]
		if value == 3 {
			if got == 0 {
				t.Fatal("commanded value was not written")
			}
			continue
		}
		if got != 0 {
			t.Fatalf("written entity retained probability at value=%d: %v", value, got)
		}
	}
	norm2, err := NormSquared(state)
	if err != nil {
		t.Fatal(err)
	}
	if norm2 < 1-1e-12 || norm2 > 1+1e-12 {
		t.Fatalf("soft memory norm=%g want=1", norm2)
	}
}

func TestUP40SoftMemoryRejectsMalformedDistribution(t *testing.T) {
	distributions := [4][]float64{
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{1, 0, 0, 0},
		{0.4, 0.4, 0, 0},
	}
	if _, err := softMemoryAfterCommandedWrite(
		distributions, 0, 2,
	); err == nil {
		t.Fatal("malformed distribution unexpectedly accepted")
	}
}
