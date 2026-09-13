package resources

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type GPU struct {
	Vendor            string   `json:"vendor"`
	UUID              string   `json:"uuid"`
	Model             *string  `json:"model"`
	TotalBytes        *uint64  `json:"total_bytes"`
	UsedBytes         *uint64  `json:"used_bytes"`
	FreeBytes         *uint64  `json:"free_bytes"`
	Utilization       *float64 `json:"utilization_percent"`
	MemoryUtilization *float64 `json:"memory_utilization_percent"`
	TemperatureC      *float64 `json:"temperature_c"`
	PowerWatts        *float64 `json:"power_watts"`
}
type GPUProvider interface {
	GPU(context.Context) (GPU, error)
}
type NVIDIA struct {
	Executable string
	Index      int
}
type limitedOutput struct{ data []byte }

func (b *limitedOutput) Write(p []byte) (int, error) {
	if len(b.data)+len(p) > 16384 {
		return 0, fmt.Errorf("GPU response too large")
	}
	b.data = append(b.data, p...)
	return len(p), nil
}

func nvidiaEnvironment() []string {
	env := []string{}
	for _, key := range []string{"SystemRoot", "WINDIR", "TEMP", "TMP", "ProgramW6432"} {
		if v := os.Getenv(key); v != "" {
			env = append(env, key+"="+v)
		}
	}
	return env
}

func (n NVIDIA) GPU(parent context.Context) (GPU, error) {
	base := strings.ToLower(filepath.Base(n.Executable))
	if !filepath.IsAbs(n.Executable) || (base != "nvidia-smi" && base != "nvidia-smi.exe") || n.Index < 0 || n.Index > 31 {
		return GPU{}, fmt.Errorf("explicit trusted nvidia-smi path and device index required")
	}
	ctx, c := context.WithTimeout(parent, 2*time.Second)
	defer c()
	cmd := exec.CommandContext(ctx, n.Executable, "--id="+strconv.Itoa(n.Index), "--query-gpu=uuid,name,memory.total,memory.used,memory.free,utilization.gpu,utilization.memory,temperature.gpu,power.draw", "--format=csv,noheader,nounits")
	cmd.Env = nvidiaEnvironment()
	var out, errout limitedOutput
	cmd.Stdout = &out
	cmd.Stderr = &errout
	cmd.WaitDelay = time.Second
	if e := cmd.Run(); e != nil {
		return GPU{}, fmt.Errorf("GPU query unavailable: %w", e)
	}
	return parseGPU(string(out.data))
}
func number(s string, max float64) (*float64, error) {
	s = strings.TrimSpace(s)
	if s == "N/A" || s == "[N/A]" || s == "[Not Supported]" || s == "Not Supported" {
		return nil, nil
	}
	v, e := strconv.ParseFloat(s, 64)
	if e != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 || v > max {
		return nil, fmt.Errorf("invalid GPU numeric field")
	}
	return &v, nil
}
func parseGPU(s string) (GPU, error) {
	var g GPU
	reader := csv.NewReader(strings.NewReader(s))
	row, e := reader.Read()
	if e != nil || len(row) != 9 {
		return g, fmt.Errorf("invalid GPU CSV shape")
	}
	if _, e = reader.Read(); e != io.EOF {
		return g, fmt.Errorf("ambiguous GPU rows")
	}
	g.Vendor = "NVIDIA"
	g.UUID = strings.TrimSpace(row[0])
	if !strings.HasPrefix(g.UUID, "GPU-") || len(g.UUID) > 80 || strings.ContainsAny(g.UUID, " ,\r\n") {
		return g, fmt.Errorf("GPU UUID required")
	}
	name := strings.TrimSpace(row[1])
	if name != "N/A" {
		g.Model = &name
	}
	values := make([]*float64, 7)
	for i := range values {
		max := float64(1 << 40)
		if i == 3 || i == 4 {
			max = 100
		}
		if i == 5 {
			max = 200
		}
		if i == 6 {
			max = 10000
		}
		values[i], e = number(row[i+2], max)
		if e != nil {
			return g, e
		}
	}
	bytes := func(v *float64) *uint64 {
		if v == nil {
			return nil
		}
		n := uint64(*v * 1048576)
		return &n
	}
	g.TotalBytes = bytes(values[0])
	g.UsedBytes = bytes(values[1])
	g.FreeBytes = bytes(values[2])
	g.Utilization = values[3]
	g.MemoryUtilization = values[4]
	g.TemperatureC = values[5]
	g.PowerWatts = values[6]
	if g.TotalBytes != nil && (g.UsedBytes != nil && *g.UsedBytes > *g.TotalBytes || g.FreeBytes != nil && *g.FreeBytes > *g.TotalBytes) {
		return GPU{}, fmt.Errorf("inconsistent GPU memory")
	}
	return g, nil
}

// WithGPU preserves optional collection errors separately; required VRAM policy still denies null.
type WithGPU struct {
	Host   Provider
	Device GPUProvider
}

func (p WithGPU) Snapshot() (Metrics, error) {
	if p.Host == nil || p.Device == nil {
		return Metrics{}, fmt.Errorf("telemetry provider not configured")
	}
	m, e := p.Host.Snapshot()
	if e != nil {
		return m, e
	}
	g, e := p.Device.GPU(context.Background())
	if e != nil {
		message := e.Error()
		m.GPUError = &message
		return m, nil
	}
	m.GPU = &g
	m.VRAMTotal = g.TotalBytes
	m.VRAMFree = g.FreeBytes
	m.GPUPercent = g.Utilization
	return m, nil
}
