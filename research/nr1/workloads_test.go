package nr1

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNR1AWorkloadCorpusComplete(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller path unavailable")
	}
	root := filepath.Join(filepath.Dir(file), "workloads")
	required := []string{
		"code_generation.txt",
		"debugging.txt",
		"repair.txt",
		"repo_reasoning.txt",
		"planning.txt",
		"review.txt",
		"instruction_following.txt",
		"long_coding.txt",
	}
	for _, name := range required {
		info, err := os.Stat(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("required workload %s unavailable: %v", name, err)
		}
		if !info.Mode().IsRegular() || info.Size() < 128 || info.Size() > 16*1024 {
			t.Fatalf("required workload %s has invalid size/type: mode=%v bytes=%d", name, info.Mode(), info.Size())
		}
	}

	system, err := os.Stat(filepath.Join(root, "system.txt"))
	if err != nil {
		t.Fatalf("system workload prompt unavailable: %v", err)
	}
	if !system.Mode().IsRegular() || system.Size() < 64 || system.Size() > 4096 {
		t.Fatalf("system workload prompt has invalid size/type: mode=%v bytes=%d", system.Mode(), system.Size())
	}
}
