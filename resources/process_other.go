//go:build !windows && !linux

package resources

import "fmt"

func sampleProcess(pid int) (ProcessMetrics, string, error) {
	return ProcessMetrics{PID: pid}, "", fmt.Errorf("process metrics unavailable")
}
