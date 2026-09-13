package benchmark

import (
	"github.com/spoonman136668-ai/Wingless/resources"
	"testing"
)

func u64(v uint64) *uint64       { return &v }
func f64(v float64) *float64     { return &v }

func TestMergePeakPreservesWorstObservedPressure(t *testing.T) {
	a := resources.Metrics{
		RAMTotal: u64(32 << 30), RAMFree: u64(20 << 30), VRAMTotal: u64(12 << 30), VRAMFree: u64(8 << 30),
		CPUPercent: f64(20), GPUPercent: f64(40),
		GPU: &resources.GPU{Vendor: "NVIDIA", UUID: "GPU-fixture", UsedBytes: u64(4 << 30), FreeBytes: u64(8 << 30), Utilization: f64(40), TemperatureC: f64(55), PowerWatts: f64(60)},
	}
	b := resources.Metrics{
		RAMTotal: u64(32 << 30), RAMFree: u64(12 << 30), VRAMTotal: u64(12 << 30), VRAMFree: u64(2 << 30),
		CPUPercent: f64(70), GPUPercent: f64(95),
		GPU: &resources.GPU{Vendor: "NVIDIA", UUID: "GPU-fixture", UsedBytes: u64(10 << 30), FreeBytes: u64(2 << 30), Utilization: f64(95), TemperatureC: f64(78), PowerWatts: f64(115)},
	}
	peak := mergePeak(a, b)
	if peak.RAMFree == nil || *peak.RAMFree != 12<<30 || peak.VRAMFree == nil || *peak.VRAMFree != 2<<30 {
		t.Fatal("minimum free memory was not retained")
	}
	if peak.CPUPercent == nil || *peak.CPUPercent != 70 || peak.GPU == nil || peak.GPU.UsedBytes == nil || *peak.GPU.UsedBytes != 10<<30 {
		t.Fatal("maximum utilization/memory use was not retained")
	}
	if peak.GPU.TemperatureC == nil || *peak.GPU.TemperatureC != 78 || peak.GPU.PowerWatts == nil || *peak.GPU.PowerWatts != 115 {
		t.Fatal("thermal/power peak was not retained")
	}
}
