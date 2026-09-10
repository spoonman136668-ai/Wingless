package resources

import (
	"fmt"
	"sync"
)

// Host measures the OS host (not process/container reservations). GPU metrics remain unavailable.
type Host struct {
	Path     string
	mu       sync.Mutex
	previous *cpuTicks
}
type cpuTicks struct{ idle, total uint64 }

func (h *Host) Snapshot() (Metrics, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	m, t, e := hostSnapshot(h.Path)
	if e != nil {
		return m, e
	}
	if h.previous != nil && t.total > h.previous.total && t.idle >= h.previous.idle {
		total := t.total - h.previous.total
		idle := t.idle - h.previous.idle
		if idle > total {
			return m, fmt.Errorf("inconsistent CPU counters")
		}
		v := 100 * float64(total-idle) / float64(total)
		m.CPUPercent = &v
	}
	h.previous = &t
	return m, nil
}
