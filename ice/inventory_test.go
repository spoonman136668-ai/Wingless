package ice

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDerivedTamperingDenied(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "a.go"), []byte("package a\nfunc Safe(){}"), 0600)
	if e := Save(root); e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(root, ".ice", "test-map.json"), []byte("[]"), 0600)
	if Validate(root) == nil {
		t.Fatal("tampered companion accepted")
	}
}
func TestIntegrationInventory(t *testing.T) {
	b, e := os.ReadFile("../docs/integration-seams.json")
	if e != nil {
		t.Fatal(e)
	}
	var doc struct {
		Seams []map[string]string `json:"seams"`
	}
	if e = json.Unmarshal(b, &doc); e != nil {
		t.Fatal(e)
	}
	if len(doc.Seams) < 9 {
		t.Fatal("missing seams")
	}
	for _, s := range doc.Seams {
		for _, k := range []string{"module", "assumed_external_interface", "reason", "expected_plane_point", "data_contract", "authority_boundary", "failure_semantics", "risk_if_plane_differs", "reconciliation_action"} {
			if s[k] == "" {
				t.Fatal("missing", k)
			}
		}
	}
}
