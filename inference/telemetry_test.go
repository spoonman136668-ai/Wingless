package inference

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTelemetrySerialization(t *testing.T) {
	r := Result{BackendID: "mock", Status: "completed", Text: "accepted=true"}
	a, e := json.Marshal(r)
	if e != nil {
		t.Fatal(e)
	}
	b, _ := json.Marshal(r)
	if string(a) != string(b) || !strings.Contains(string(a), `"ram_free":null`) || !strings.Contains(string(a), `"time_to_first_token_ms":null`) {
		t.Fatal(string(a))
	}
}
