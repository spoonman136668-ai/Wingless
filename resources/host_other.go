//go:build !linux && !windows

package resources

import "fmt"

func hostSnapshot(string) (Metrics, cpuTicks, error) {
	return Metrics{}, cpuTicks{}, fmt.Errorf("host collector unsupported on this OS")
}
