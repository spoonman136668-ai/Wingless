//go:build linux || windows

package resources

import (
	"testing"
	"time"
)

func TestNativeHost(t *testing.T) {
	h := &Host{Path: t.TempDir()}
	first, e := h.Snapshot()
	if e != nil {
		t.Fatal(e)
	}
	if first.RAMTotal == nil || first.RAMFree == nil || first.DiskFree == nil || *first.RAMTotal == 0 || *first.RAMFree > *first.RAMTotal || first.CPUPercent != nil {
		t.Fatalf("bad first snapshot: %+v", first)
	}
	time.Sleep(30 * time.Millisecond)
	second, e := h.Snapshot()
	if e != nil || second.CPUPercent == nil || *second.CPUPercent < 0 || *second.CPUPercent > 100 {
		t.Fatal(second, e)
	}
	if second.VRAMFree != nil || second.GPUPercent != nil {
		t.Fatal("fabricated GPU metrics")
	}
}
