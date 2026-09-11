//go:build linux

package resources

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func processIdentity(base string) (string, error) {
	b, e := os.ReadFile(base + "/stat")
	if e != nil {
		return "", e
	}
	i := strings.LastIndex(string(b), ")")
	if i < 0 {
		return "", fmt.Errorf("invalid proc stat")
	}
	f := strings.Fields(string(b[i+1:]))
	if len(f) < 20 {
		return "", fmt.Errorf("short proc stat")
	}
	return f[19], nil
}
func sampleProcess(pid int) (ProcessMetrics, string, error) {
	m := ProcessMetrics{PID: pid}
	self, e := os.Readlink("/proc/self")
	if e != nil {
		return m, "", e
	}
	if self != strconv.Itoa(os.Getpid()) {
		return m, "", ErrProcessNamespace
	}
	base := "/proc/" + strconv.Itoa(pid)
	id, e := processIdentity(base)
	if e != nil {
		return m, "", e
	}
	b, e := os.ReadFile(base + "/status")
	if e != nil {
		return m, id, e
	}
	for _, line := range strings.Split(string(b), "\n") {
		f := strings.Fields(line)
		if len(f) == 3 && (f[0] == "VmRSS:" || f[0] == "VmHWM:") {
			n, e := strconv.ParseUint(f[1], 10, 64)
			if e == nil {
				n *= 1024
				if f[0] == "VmRSS:" {
					m.RSSBytes = &n
				} else {
					m.PeakRSSBytes = &n
				}
			}
		}
	}
	// schedstat is per-thread; do not mislabel it as aggregate process CPU time.

	if b, e = os.ReadFile(base + "/io"); e == nil {
		for _, line := range strings.Split(string(b), "\n") {
			f := strings.Fields(line)
			if len(f) == 2 {
				n, e := strconv.ParseUint(f[1], 10, 64)
				if e == nil {
					if f[0] == "read_bytes:" {
						m.DiskReadBytes = &n
					}
					if f[0] == "write_bytes:" {
						m.DiskWriteBytes = &n
					}
				}
			}
		}
	}
	end, e := processIdentity(base)
	if e != nil || end != id {
		return ProcessMetrics{}, "", fmt.Errorf("process exited or changed during sample")
	}
	return m, id, nil
}
