package benchmark

import (
	"context"
	"github.com/spoonman136668-ai/Wingless/resources"
	"sync"
	"time"
)

const resourceSampleInterval = 500 * time.Millisecond

type resourceSampler struct {
	provider resources.Provider
	stop     chan struct{}
	done     chan struct{}
	once     sync.Once
	mu       sync.Mutex
	peak     *resources.Metrics
	err      error
}

func startResourceSampler(ctx context.Context, provider resources.Provider, initial resources.Metrics) *resourceSampler {
	s := &resourceSampler{provider: provider, stop: make(chan struct{}), done: make(chan struct{})}
	copy := initial
	s.peak = &copy
	go s.run(ctx)
	return s
}

func (s *resourceSampler) run(ctx context.Context) {
	defer close(s.done)
	ticker := time.NewTicker(resourceSampleInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stop:
			return
		case <-ticker.C:
			m, err := s.provider.Snapshot()
			s.mu.Lock()
			if err != nil {
				if s.err == nil {
					s.err = err
				}
			} else {
				merged := mergePeak(*s.peak, m)
				s.peak = &merged
			}
			s.mu.Unlock()
		}
	}
}

func (s *resourceSampler) finish(final *resources.Metrics) (*resources.Metrics, error) {
	s.once.Do(func() { close(s.stop) })
	<-s.done
	s.mu.Lock()
	defer s.mu.Unlock()
	if final != nil {
		merged := mergePeak(*s.peak, *final)
		s.peak = &merged
	}
	copy := *s.peak
	return &copy, s.err
}

func maxFloat(a, b *float64) *float64 {
	if a == nil {
		return b
	}
	if b == nil || *a >= *b {
		return a
	}
	return b
}

func minUint(a, b *uint64) *uint64 {
	if a == nil {
		return b
	}
	if b == nil || *a <= *b {
		return a
	}
	return b
}

func maxUint(a, b *uint64) *uint64 {
	if a == nil {
		return b
	}
	if b == nil || *a >= *b {
		return a
	}
	return b
}

func mergePeak(a, b resources.Metrics) resources.Metrics {
	out := a
	out.RAMTotal = maxUint(a.RAMTotal, b.RAMTotal)
	out.RAMFree = minUint(a.RAMFree, b.RAMFree)
	out.VRAMTotal = maxUint(a.VRAMTotal, b.VRAMTotal)
	out.VRAMFree = minUint(a.VRAMFree, b.VRAMFree)
	out.DiskFree = minUint(a.DiskFree, b.DiskFree)
	out.CPUPercent = maxFloat(a.CPUPercent, b.CPUPercent)
	out.GPUPercent = maxFloat(a.GPUPercent, b.GPUPercent)
	out.DiskBytesRead = maxUint(a.DiskBytesRead, b.DiskBytesRead)
	out.EnergyJoules = maxFloat(a.EnergyJoules, b.EnergyJoules)
	if b.GPUError != nil {
		out.GPUError = b.GPUError
	}
	if b.Resident != nil {
		out.Resident = b.Resident
	}
	if a.GPU == nil {
		out.GPU = b.GPU
	} else if b.GPU != nil {
		g := *a.GPU
		if g.Model == nil {
			g.Model = b.GPU.Model
		}
		if g.UUID == "" {
			g.UUID = b.GPU.UUID
		}
		if g.Vendor == "" {
			g.Vendor = b.GPU.Vendor
		}
		g.TotalBytes = maxUint(a.GPU.TotalBytes, b.GPU.TotalBytes)
		g.UsedBytes = maxUint(a.GPU.UsedBytes, b.GPU.UsedBytes)
		g.FreeBytes = minUint(a.GPU.FreeBytes, b.GPU.FreeBytes)
		g.Utilization = maxFloat(a.GPU.Utilization, b.GPU.Utilization)
		g.MemoryUtilization = maxFloat(a.GPU.MemoryUtilization, b.GPU.MemoryUtilization)
		g.TemperatureC = maxFloat(a.GPU.TemperatureC, b.GPU.TemperatureC)
		g.PowerWatts = maxFloat(a.GPU.PowerWatts, b.GPU.PowerWatts)
		out.GPU = &g
	}
	return out
}
