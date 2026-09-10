//go:build windows

package resources

import (
	"fmt"
	"syscall"
	"unsafe"
)

var kernel = syscall.NewLazyDLL("kernel32.dll")
var memoryStatus = kernel.NewProc("GlobalMemoryStatusEx")
var diskSpace = kernel.NewProc("GetDiskFreeSpaceExW")
var systemTimes = kernel.NewProc("GetSystemTimes")

func hostSnapshot(path string) (Metrics, cpuTicks, error) {
	var m Metrics
	var t cpuTicks
	mem := struct {
		Length, Load                                                               uint32
		Total, Available, PageTotal, PageFree, VirtualTotal, VirtualFree, Extended uint64
	}{}
	mem.Length = uint32(unsafe.Sizeof(mem))
	ok, _, e := memoryStatus.Call(uintptr(unsafe.Pointer(&mem)))
	if ok == 0 {
		return m, t, fmt.Errorf("GlobalMemoryStatusEx: %w", e)
	}
	m.RAMTotal = &mem.Total
	m.RAMFree = &mem.Available
	if path == "" {
		path = "."
	}
	p, e := syscall.UTF16PtrFromString(path)
	if e != nil {
		return m, t, e
	}
	var free, total, allFree uint64
	ok, _, e = diskSpace.Call(uintptr(unsafe.Pointer(p)), uintptr(unsafe.Pointer(&free)), uintptr(unsafe.Pointer(&total)), uintptr(unsafe.Pointer(&allFree)))
	if ok == 0 {
		return m, t, fmt.Errorf("GetDiskFreeSpaceExW: %w", e)
	}
	m.DiskFree = &free
	var idle, kern, user syscall.Filetime
	ok, _, e = systemTimes.Call(uintptr(unsafe.Pointer(&idle)), uintptr(unsafe.Pointer(&kern)), uintptr(unsafe.Pointer(&user)))
	if ok == 0 {
		return m, t, fmt.Errorf("GetSystemTimes: %w", e)
	}
	ticks := func(f syscall.Filetime) uint64 { return uint64(f.HighDateTime)<<32 | uint64(f.LowDateTime) }
	t = cpuTicks{ticks(idle), ticks(kern) + ticks(user)}
	return m, t, nil
}
