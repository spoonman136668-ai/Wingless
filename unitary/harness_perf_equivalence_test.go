package unitary

import "testing"

func TestHarnessPerfMinimumSpeedupFrozen(t *testing.T) {
	if harnessPerfMinimumSpeedup != 1.05 {
		t.Fatalf(
			"minimum speedup=%g want=1.05",
			harnessPerfMinimumSpeedup,
		)
	}
}
