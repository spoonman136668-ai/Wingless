//go:build windows && amd64

package modelhost

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var jobs = syscall.NewLazyDLL("kernel32.dll")
var createJob = jobs.NewProc("CreateJobObjectW")
var setJob = jobs.NewProc("SetInformationJobObject")
var assignJob = jobs.NewProc("AssignProcessToJobObject")

type containment struct{ handle syscall.Handle }

func newContainment() (*containment, error) {
	h, _, e := createJob.Call(0, 0)
	if h == 0 {
		return nil, fmt.Errorf("create job: %w", e)
	}
	c := &containment{syscall.Handle(h)}
	info := struct {
		ProcessTime, JobTime                           int64
		Flags                                          uint32
		Padding                                        uint32
		MinWS, MaxWS                                   uintptr
		Active                                         uint32
		Pad                                            uint32
		Affinity                                       uintptr
		Priority, Scheduling                           uint32
		IO                                             [6]uint64
		ProcessMemory, JobMemory, PeakProcess, PeakJob uintptr
	}{}
	info.Flags = 0x2000 // JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	ok, _, e := setJob.Call(h, 9, uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	if ok == 0 {
		c.Close()
		return nil, fmt.Errorf("configure job: %w", e)
	}
	return c, nil
}
func (c *containment) Attach(p *os.Process) error {
	h, e := syscall.OpenProcess(0x101, false, uint32(p.Pid))
	if e != nil {
		return e
	}
	defer syscall.CloseHandle(h)
	ok, _, e := assignJob.Call(uintptr(c.handle), uintptr(h))
	if ok == 0 {
		return fmt.Errorf("assign job: %w", e)
	}
	return nil
}
func (c *containment) Close() {
	if c != nil && c.handle != 0 {
		syscall.CloseHandle(c.handle)
		c.handle = 0
	}
}
func containmentName() string { return "windows_job_kill_on_close_post_start" }
