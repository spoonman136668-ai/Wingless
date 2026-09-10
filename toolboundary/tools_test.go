package toolboundary

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBoundary(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "ok.txt"), []byte("hello"), 0600)
	p := Policy{Workspace: root, AllowedTools: []string{"read_file"}, AllowedPaths: []string{"ok.txt", "../escape", "link"}, Timeout: time.Second, MaxOutputBytes: 100}
	o, e := Execute(context.Background(), p, Proposal{"read_file", "ok.txt"})
	if e != nil || o.Stdout != "hello" || o.ExitCode != 0 {
		t.Fatal(o, e)
	}
	for _, a := range []Proposal{{"shell", "ok.txt"}, {"read_file", "../escape"}, {"read_file", "unknown"}} {
		if _, e := Execute(context.Background(), p, a); e == nil {
			t.Fatal(a)
		}
	}
	outside := filepath.Join(t.TempDir(), "secret")
	os.WriteFile(outside, []byte("secret"), 0600)
	if e = os.Symlink(outside, filepath.Join(root, "link")); e == nil {
		if _, e = Execute(context.Background(), p, Proposal{"read_file", "link"}); e == nil {
			t.Fatal("symlink escaped")
		}
	}
	p.MaxOutputBytes = 2
	if _, e = Execute(context.Background(), p, Proposal{"read_file", "ok.txt"}); e == nil {
		t.Fatal("output bound")
	}
	ctx, c := context.WithCancel(context.Background())
	c()
	if _, e = Execute(ctx, p, Proposal{"read_file", "ok.txt"}); e == nil {
		t.Fatal("cancellation")
	}
}
