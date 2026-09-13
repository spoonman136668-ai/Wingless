package resources

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestNVIDIAParsing(t *testing.T) {
	g, e := parseGPU("GPU-abc, RTX fixture, 12288, 2048, 10240, 50, 20, 65, 100.5\n")
	if e != nil || *g.FreeBytes != 10240*1048576 || *g.PowerWatts != 100.5 {
		t.Fatal(g, e)
	}
	g, e = parseGPU("GPU-abc, RTX fixture, N/A, N/A, N/A, N/A, N/A, N/A, N/A\n")
	if e != nil || g.FreeBytes != nil || g.PowerWatts != nil {
		t.Fatal(g, e)
	}
	for _, s := range []string{"bad", "GPU-a, x, 1, 2, 3, NaN, 0, 0, 0", "GPU-a,x,1,0,1,200,0,0,0"} {
		if _, e = parseGPU(s); e == nil {
			t.Fatal(s)
		}
	}
}

func TestNVIDIAEnvironmentIsBounded(t *testing.T) {
	t.Setenv("SystemRoot", `C:\Windows`)
	t.Setenv("WINDIR", `C:\Windows`)
	t.Setenv("TEMP", `C:\Temp`)
	t.Setenv("TMP", `C:\Tmp`)
	t.Setenv("ProgramW6432", `C:\Program Files`)
	t.Setenv("PATH", `C:\sensitive-path`)
	t.Setenv("OPENAI_API_KEY", "do-not-inherit")
	t.Setenv("HTTPS_PROXY", "http://do-not-inherit")

	want := []string{
		`SystemRoot=C:\Windows`,
		`WINDIR=C:\Windows`,
		`TEMP=C:\Temp`,
		`TMP=C:\Tmp`,
		`ProgramW6432=C:\Program Files`,
	}
	if got := nvidiaEnvironment(); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected NVIDIA environment: got=%q want=%q", got, want)
	}
}

type unavailableGPU struct{}

func (unavailableGPU) GPU(context.Context) (GPU, error) { return GPU{}, fmt.Errorf("unavailable") }

type emptyHost struct{}

func (emptyHost) Snapshot() (Metrics, error) { return Metrics{}, nil }
func TestOptionalGPUStillDeniesRequiredVRAM(t *testing.T) {
	m, e := (WithGPU{emptyHost{}, unavailableGPU{}}).Snapshot()
	if e != nil || m.GPUError == nil || m.VRAMFree != nil || Check(Policy{MinVRAM: 1}, m) == nil {
		t.Fatal(m, e)
	}
}
