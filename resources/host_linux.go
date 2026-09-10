//go:build linux

package resources

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func hostSnapshot(path string) (Metrics, cpuTicks, error) {
	var m Metrics
	var t cpuTicks
	b, e := os.ReadFile("/proc/meminfo")
	if e != nil {
		return m, t, e
	}
	values := map[string]uint64{}
	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 3 && fields[2] == "kB" {
			n, e := strconv.ParseUint(fields[1], 10, 64)
			if e != nil {
				return m, t, e
			}
			values[strings.TrimSuffix(fields[0], ":")] = n * 1024
		}
	}
	total, ok := values["MemTotal"]
	if !ok {
		return m, t, fmt.Errorf("MemTotal unavailable")
	}
	free, ok := values["MemAvailable"]
	if !ok {
		return m, t, fmt.Errorf("MemAvailable unavailable")
	}
	m.RAMTotal = &total
	m.RAMFree = &free
	if path == "" {
		path = "."
	}
	var st syscall.Statfs_t
	if e = syscall.Statfs(path, &st); e != nil {
		return m, t, e
	}
	disk := st.Bavail * uint64(st.Bsize)
	m.DiskFree = &disk
	b, e = os.ReadFile("/proc/stat")
	if e != nil {
		return m, t, e
	}
	fields := strings.Fields(strings.SplitN(string(b), "\n", 2)[0])
	if len(fields) < 9 || fields[0] != "cpu" {
		return m, t, fmt.Errorf("CPU counters unavailable")
	}
	for i := 1; i <= 8; i++ {
		v, e := strconv.ParseUint(fields[i], 10, 64)
		if e != nil {
			return m, t, e
		}
		t.total += v
		if i == 4 || i == 5 {
			t.idle += v
		}
	}
	return m, t, nil
}
