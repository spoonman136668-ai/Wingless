package contextbroker

import (
	"github.com/spoonman136668-ai/Wingless/ice"
	"os"
	"path/filepath"
	"testing"
)

type oversized struct{}

func (oversized) Retrieve(string, int) ([]Snippet, error) { return []Snippet{{Text: "12345"}}, nil }
func TestContextBounds(t *testing.T) {
	if _, e := Collect(oversized{}, "q", 4); e == nil {
		t.Fatal("provider escaped bound")
	}
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "a.go"), []byte("package a\nfunc Small() {}\n"), 0600)
	if e := ice.Save(root); e != nil {
		t.Fatal(e)
	}
	s, e := Collect(ICE{root}, "Small", 100)
	if e != nil || len(s) != 1 || s[0].Text != "func Small() {}" {
		t.Fatal(s, e)
	}
	s, e = Collect(ICE{root}, "Small", 2)
	if e != nil || len(s) != 0 {
		t.Fatal(s, e)
	}
}
