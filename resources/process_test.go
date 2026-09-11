//go:build windows || linux

package resources

import (
	"errors"
	"os"
	"testing"
)

func TestOwnedProcessSampling(t *testing.T) {
	p, e := NewProcessProbe(os.Getpid())
	if errors.Is(e, ErrProcessNamespace) {
		t.Skip(e)
	}
	if e != nil {
		t.Fatal(e)
	}
	m, e := p.Snapshot()
	if e != nil || m.PID != os.Getpid() || m.RSSBytes == nil || m.PeakRSSBytes == nil || m.LifetimeMS == nil || m.GPUMemoryBytes != nil {
		t.Fatal(m, e)
	}
	p.identity = "wrong"
	if _, e = p.Snapshot(); e == nil {
		t.Fatal("identity drift accepted")
	}
}
