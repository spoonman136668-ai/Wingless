package ice

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeclarationsAndOperationalFiles(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "a.go"), []byte("package a\nconst Budget=2\nfunc External()\n"), 0600)
	os.MkdirAll(filepath.Join(root, ".github", "workflows"), 0700)
	path := filepath.Join(root, ".github", "workflows", "test.yml")
	os.WriteFile(path, []byte("name: original"), 0600)
	if e := Save(root); e != nil {
		t.Fatal(e)
	}
	q, e := Query(root, "Budget")
	if e != nil || len(q) != 1 {
		t.Fatal(q, e)
	}
	os.WriteFile(path, []byte("name: changed"), 0600)
	if Validate(root) == nil {
		t.Fatal("workflow drift missed")
	}
}
