package benchmark

import "testing"

func TestProfilesFailClosed(t *testing.T) {
	for _, name := range []string{"FAST_LOCAL", "DEEP_LOCAL", "HYBRID", "unknown"} {
		if Profile(name, false) == nil {
			t.Fatal(name)
		}
	}
	if Profile("GPU_RESIDENT_SMALL", false) == nil || Profile("CPU_ONLY", true) == nil {
		t.Fatal("silent substitution")
	}
	if Profile("CPU_ONLY", false) != nil || Profile("GPU_RESIDENT_SMALL", true) != nil {
		t.Fatal("valid profile")
	}
}
