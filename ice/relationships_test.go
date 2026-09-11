package ice

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTypeResolvedRelationships(t *testing.T) {
	root := t.TempDir()
	source := "package fixture\ntype I interface { F() }\ntype T struct{}\nfunc(T) F(){}\nfunc Use(t T){t.F()}\n"
	os.WriteFile(filepath.Join(root, "fixture.go"), []byte(source), 0600)
	if e := Save(root); e != nil {
		t.Fatal(e)
	}
	r, e := Relationships(root)
	if e != nil || len(r.Diagnostics) != 0 {
		t.Fatal(r, e)
	}
	impl, call := false, false
	for _, x := range r.Relationships {
		impl = impl || x.Kind == "same_package_implementation"
		call = call || x.Kind == "resolved_call_target"
	}
	if !impl || !call {
		t.Fatal(r)
	}
	os.WriteFile(filepath.Join(root, "fixture.go"), []byte(source+"var X=missing\n"), 0600)
	if _, e = Relationships(root); e == nil {
		t.Fatal("stale source accepted")
	}
}
