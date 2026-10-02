package unitary

import "testing"

func TestUP45FrozenControls(t *testing.T) {
	if proximalFusionLambda !=
		directOffsetLearningRate*directOffsetResourcePrice {
		t.Fatalf("proximal lambda is not derived")
	}
	if rollingTaskMaximumSteps != 48 {
		t.Fatalf("maximum steps=%d want=48", rollingTaskMaximumSteps)
	}
	want := []int{12, 24, 48}
	if len(rollingTaskCheckpointSteps) != len(want) {
		t.Fatalf("checkpoint count mismatch")
	}
	for i := range want {
		if rollingTaskCheckpointSteps[i] != want[i] {
			t.Fatalf(
				"checkpoint[%d]=%d want=%d",
				i, rollingTaskCheckpointSteps[i], want[i],
			)
		}
	}
}
