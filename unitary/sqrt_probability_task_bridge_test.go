package unitary

import (
	"math"
	"testing"
)

func TestUP44SqrtSoftMemoryOneHotMatchesCanonicalWrite(t *testing.T) {
	distributions := [4][]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	got, err := sqrtSoftMemoryAfterCommandedWrite(
		distributions, 1, 3,
	)
	if err != nil {
		t.Fatal(err)
	}
	want, err := encodeMemory(memoryTable{0, 3, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("dimension=%d want=%d", len(got), len(want))
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf(
				"state[%d]=%v want=%v",
				index, got[index], want[index],
			)
		}
	}
}

func TestUP44SqrtSoftMemoryIsUnitNormWithoutGlobalNormalize(t *testing.T) {
	distributions := [4][]float64{
		{0.7, 0.1, 0.1, 0.1},
		{0.4, 0.3, 0.2, 0.1},
		{0.25, 0.25, 0.25, 0.25},
		{0.05, 0.15, 0.30, 0.50},
	}
	state, err := sqrtSoftMemoryAfterCommandedWrite(
		distributions, 2, 1,
	)
	if err != nil {
		t.Fatal(err)
	}
	norm2, err := NormSquared(state)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(norm2-1) > 1e-12 {
		t.Fatalf("norm2=%g want=1", norm2)
	}
}

func TestUP44SqrtSoftMemoryPreservesUnwrittenProbabilities(t *testing.T) {
	distributions := [4][]float64{
		{0.7, 0.1, 0.1, 0.1},
		{0.4, 0.3, 0.2, 0.1},
		{0.25, 0.25, 0.25, 0.25},
		{0.05, 0.15, 0.30, 0.50},
	}
	state, err := sqrtSoftMemoryAfterCommandedWrite(
		distributions, 2, 1,
	)
	if err != nil {
		t.Fatal(err)
	}
	for entity := 0; entity < 4; entity++ {
		if entity == 2 {
			continue
		}
		for value := 0; value < 4; value++ {
			amplitude := real(state[entity*4+value])
			got := 4 * amplitude * amplitude
			want := distributions[entity][value]
			if math.Abs(got-want) > 1e-12 {
				t.Fatalf(
					"entity=%d value=%d recovered=%g want=%g",
					entity, value, got, want,
				)
			}
		}
	}
}

func TestUP44KeepsUP43OptimizerAndHorizon(t *testing.T) {
	if rollingTaskMaximumSteps != 48 {
		t.Fatalf(
			"maximum steps=%d want=48",
			rollingTaskMaximumSteps,
		)
	}
	if rollingGradientBufferSize != 7 {
		t.Fatalf(
			"rolling buffer=%d want=7",
			rollingGradientBufferSize,
		)
	}
	want := []int{12, 24, 48}
	if len(rollingTaskCheckpointSteps) != len(want) {
		t.Fatalf(
			"checkpoint count=%d want=%d",
			len(rollingTaskCheckpointSteps), len(want),
		)
	}
	for index := range want {
		if rollingTaskCheckpointSteps[index] != want[index] {
			t.Fatalf(
				"checkpoint[%d]=%d want=%d",
				index,
				rollingTaskCheckpointSteps[index],
				want[index],
			)
		}
	}
}
