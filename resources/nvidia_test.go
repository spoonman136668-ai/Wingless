package resources

import (
	"context"
	"fmt"
	"testing"
)

func TestNVIDIAParsing(t *testing.T) {
	g, e := parseGPU("GPU-abc, RTX fixture, 12288, 2048, 10240, 50, 20, 65, 100.5\n")
	if e != nil || *g.FreeBytes != 10240*1048576 || *g.PowerWatts != 100.5 {
		t.Fatal(g, e)
	}
	g, e = parseGPU("GPU-abc, RTX fixture, N/A, N/A, N/A, N/A, N/A, N/A, N/A\n")
	if e != nil || g.FreeBytes != nil || g.PowerWatts != nil {
		t.Fatal(g, e)
	}
	for _, s := range []string{"bad", "GPU-a, x, 1, 2, 3, NaN, 0, 0, 0", "GPU-a,x,1,0,1,200,0,0,0"} {
		if _, e = parseGPU(s); e == nil {
			t.Fatal(s)
		}
	}
}

type unavailableGPU struct{}

func (unavailableGPU) GPU(context.Context) (GPU, error) { return GPU{}, fmt.Errorf("unavailable") }

type emptyHost struct{}

func (emptyHost) Snapshot() (Metrics, error) { return Metrics{}, nil }
func TestOptionalGPUStillDeniesRequiredVRAM(t *testing.T) {
	m, e := (WithGPU{emptyHost{}, unavailableGPU{}}).Snapshot()
	if e != nil || m.GPUError == nil || m.VRAMFree != nil || Check(Policy{MinVRAM: 1}, m) == nil {
		t.Fatal(m, e)
	}
}
