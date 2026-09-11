package resources

import (
	"fmt"
	"math"
)

type Metrics struct {
	GPU           *GPU     `json:"gpu"`
	GPUError      *string  `json:"gpu_error"`
	RAMTotal      *uint64  `json:"ram_total"`
	RAMFree       *uint64  `json:"ram_free"`
	VRAMTotal     *uint64  `json:"vram_total"`
	VRAMFree      *uint64  `json:"vram_free"`
	DiskFree      *uint64  `json:"disk_free"`
	CPUPercent    *float64 `json:"cpu_percent"`
	GPUPercent    *float64 `json:"gpu_percent"`
	Resident      *bool    `json:"resident"`
	DiskBytesRead *uint64  `json:"disk_bytes_read"`
	EnergyJoules  *float64 `json:"energy_joules"`
}
type Policy struct {
	MinRAM  uint64  `json:"min_ram"`
	MinVRAM uint64  `json:"min_vram"`
	MinDisk uint64  `json:"min_disk"`
	MaxCPU  float64 `json:"max_cpu"`
}
type Provider interface{ Snapshot() (Metrics, error) }

func Check(p Policy, m Metrics) error {
	for _, x := range []struct {
		name   string
		min    uint64
		actual *uint64
	}{{"RAM", p.MinRAM, m.RAMFree}, {"VRAM", p.MinVRAM, m.VRAMFree}, {"disk", p.MinDisk, m.DiskFree}} {
		if x.min > 0 && (x.actual == nil || *x.actual < x.min) {
			return fmt.Errorf("resource denied: %s unavailable or below floor", x.name)
		}
	}
	if math.IsNaN(p.MaxCPU) || p.MaxCPU < 0 || p.MaxCPU > 100 {
		return fmt.Errorf("invalid CPU ceiling")
	}
	if p.MaxCPU > 0 && (m.CPUPercent == nil || math.IsNaN(*m.CPUPercent) || *m.CPUPercent > p.MaxCPU || *m.CPUPercent < 0) {
		return fmt.Errorf("CPU unavailable or above ceiling")
	}
	return nil
}
