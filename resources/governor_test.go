package resources

import (
	"math"
	"testing"
)

func TestFloors(t *testing.T) {
	low := uint64(3)
	high := uint64(8)
	for _, m := range []Metrics{{}, {RAMFree: &low}} {
		if Check(Policy{MinRAM: 4}, m) == nil {
			t.Fatal("floor allowed")
		}
	}
	if e := Check(Policy{MinRAM: 4}, Metrics{RAMFree: &high}); e != nil {
		t.Fatal(e)
	}
	if Check(Policy{MaxCPU: math.NaN()}, Metrics{}) == nil {
		t.Fatal("NaN policy")
	}
}
