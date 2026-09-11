package resources

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrProcessNamespace = errors.New("procfs PID namespace differs from caller")

type ProcessMetrics struct {
	PID            int     `json:"pid"`
	RSSBytes       *uint64 `json:"rss_bytes"`
	PeakRSSBytes   *uint64 `json:"peak_rss_bytes"`
	CPUTimeNS      *uint64 `json:"cpu_time_ns"`
	LifetimeMS     *int64  `json:"observed_lifetime_ms"`
	DiskReadBytes  *uint64 `json:"disk_read_bytes"`
	DiskWriteBytes *uint64 `json:"disk_write_bytes"`
	GPUMemoryBytes *uint64 `json:"gpu_memory_bytes"`
}

// ProcessProbe detects PID identity changes; only the supervisor supplies an owned PID.
type ProcessProbe struct {
	mu       sync.Mutex
	pid      int
	identity string
	started  time.Time
}

func NewProcessProbe(pid int) (*ProcessProbe, error) {
	if pid <= 0 {
		return nil, fmt.Errorf("invalid process")
	}
	_, id, e := sampleProcess(pid)
	if e != nil {
		return nil, e
	}
	return &ProcessProbe{pid: pid, identity: id, started: time.Now()}, nil
}
func (p *ProcessProbe) Snapshot() (ProcessMetrics, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	m, id, e := sampleProcess(p.pid)
	if e != nil {
		return m, e
	}
	if id != p.identity {
		return ProcessMetrics{}, fmt.Errorf("process identity changed")
	}
	ms := time.Since(p.started).Milliseconds()
	m.LifetimeMS = &ms
	return m, nil
}
