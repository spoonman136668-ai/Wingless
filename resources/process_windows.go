//go:build windows

package resources

import (
	"fmt"
	"syscall"
	"unsafe"
)

var processMemory = syscall.NewLazyDLL("psapi.dll").NewProc("GetProcessMemoryInfo")
var processTimes = kernel.NewProc("GetProcessTimes")
var processIO = kernel.NewProc("GetProcessIoCounters")

func sampleProcess(pid int) (ProcessMetrics, string, error) {
	m := ProcessMetrics{PID: pid}
	h, e := syscall.OpenProcess(0x410, false, uint32(pid))
	if e != nil {
		return m, "", e
	}
	defer syscall.CloseHandle(h)
	var created, exit, kern, user syscall.Filetime
	ok, _, e := processTimes.Call(uintptr(h), uintptr(unsafe.Pointer(&created)), uintptr(unsafe.Pointer(&exit)), uintptr(unsafe.Pointer(&kern)), uintptr(unsafe.Pointer(&user)))
	if ok == 0 {
		return m, "", e
	}
	ticks := func(f syscall.Filetime) uint64 { return uint64(f.HighDateTime)<<32 | uint64(f.LowDateTime) }
	id := fmt.Sprint(ticks(created))
	cpu := (ticks(kern) + ticks(user)) * 100
	m.CPUTimeNS = &cpu
	mem := struct {
		Size, Faults                                                                    uint32
		Peak, Working, PeakPaged, Paged, PeakNonPaged, NonPaged, Pagefile, PeakPagefile uintptr
	}{}
	mem.Size = uint32(unsafe.Sizeof(mem))
	ok, _, _ = processMemory.Call(uintptr(h), uintptr(unsafe.Pointer(&mem)), uintptr(mem.Size))
	if ok != 0 {
		rss, peak := uint64(mem.Working), uint64(mem.Peak)
		m.RSSBytes = &rss
		m.PeakRSSBytes = &peak
	}
	var counters struct{ ReadOps, WriteOps, OtherOps, ReadBytes, WriteBytes, OtherBytes uint64 }
	ok, _, _ = processIO.Call(uintptr(h), uintptr(unsafe.Pointer(&counters)))
	if ok != 0 {
		m.DiskReadBytes = &counters.ReadBytes
		m.DiskWriteBytes = &counters.WriteBytes
	}
	return m, id, nil
}
