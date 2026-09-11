package modelhost

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if len(os.Args) > 2 && os.Args[1] == "--model" {
		args := map[string]string{}
		for i := 1; i+1 < len(os.Args); i += 2 {
			args[os.Args[i]] = os.Args[i+1]
		}
		http.HandleFunc("/v1/models", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, `{"data":[{"id":%q}]}`, args["--alias"]) })
		if http.ListenAndServe("127.0.0.1:"+args["--port"], nil) != nil {
			os.Exit(3)
		}
		return
	}
	os.Exit(m.Run())
}

type ample struct{}

func (ample) Snapshot() (resources.Metrics, error) {
	n := uint64(4 << 30)
	return resources.Metrics{RAMFree: &n}, nil
}
func TestOwnedProcessLifecycle(t *testing.T) {
	root := t.TempDir()
	name := "llama-server"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	exe := filepath.Join(root, name)
	self, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	src, e := os.Open(self)
	if e != nil {
		t.Fatal(e)
	}
	defer src.Close()
	dst, e := os.OpenFile(exe, os.O_CREATE|os.O_WRONLY, 0700)
	if e != nil {
		t.Fatal(e)
	}
	h := sha256.New()
	if _, e = io.Copy(io.MultiWriter(dst, h), src); e != nil {
		t.Fatal(e)
	}
	dst.Close()
	model := filepath.Join(root, "fixture.gguf")
	os.WriteFile(model, []byte("fixture"), 0600)
	mh := fmt.Sprintf("%x", sha256.Sum256([]byte("fixture")))
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	s, e := New(Config{Executable: exe, ExecutableSHA256: fmt.Sprintf("%x", h.Sum(nil)), Model: model, ModelSHA256: mh, ModelID: "fixture", Port: port, Threads: 1, ContextTokens: 128, StartupSeconds: 5, RuntimeSeconds: 10, MinRAM: 1 << 30}, ample{})
	if e != nil {
		t.Fatal(e)
	}
	ctx, c := context.WithTimeout(context.Background(), 15*time.Second)
	defer c()
	defer s.Stop(ctx)
	if e = s.Start(ctx); e != nil {
		t.Fatal(e, s.Status())
	}
	if s.Status().State != inference.Ready || s.Status().PID == 0 {
		t.Fatal(s.Status())
	}
	if e = s.Start(ctx); e == nil {
		t.Fatal("double start")
	}
	if e = s.Stop(ctx); e != nil {
		t.Fatal(e)
	}
	st := s.Status()
	if st.State != inference.Stopped || st.Reason != "operator_stop" || st.CompletedAt == nil || st.ExitCode == nil {
		t.Fatal(st)
	}
}
