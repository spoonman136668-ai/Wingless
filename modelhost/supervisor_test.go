package modelhost

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/resources"
	"os"
	"path/filepath"
	"testing"
)

type unknown struct{}

func (unknown) Snapshot() (resources.Metrics, error) { return resources.Metrics{}, nil }
func TestPinnedAndResourceDenial(t *testing.T) {
	root := t.TempDir()
	exe := filepath.Join(root, "llama-server")
	model := filepath.Join(root, "small.gguf")
	data := []byte("fixture")
	os.WriteFile(exe, data, 0700)
	os.WriteFile(model, data, 0600)
	sum := fmt.Sprintf("%x", sha256.Sum256(data))
	if pinned(model, sum, 100) != nil {
		t.Fatal("valid pin")
	}
	if pinned(model, sum, 2) == nil {
		t.Fatal("size accepted")
	}
	c := Config{Executable: exe, Model: model, ExecutableSHA256: sum, ModelSHA256: sum, ModelID: "fixture", Port: 18080, Threads: 1, ContextTokens: 512, StartupSeconds: 1, RuntimeSeconds: 2, MinRAM: 1 << 30}
	s, e := New(c, unknown{})
	if e != nil {
		t.Fatal(e)
	}
	if e = s.Start(context.Background()); e == nil || s.Status().PID != 0 {
		t.Fatal("spawned without measured RAM")
	}
	c.Executable = filepath.Join(root, "shell")
	if _, e = New(c, unknown{}); e == nil {
		t.Fatal("arbitrary command accepted")
	}
}
func TestTailBound(t *testing.T) {
	var b tail
	data := make([]byte, 100000)
	n, _ := b.Write(data)
	if n != len(data) || len(b.text()) != 8192 {
		t.Fatal("unbounded logs")
	}
}
