package main

import (
	"encoding/json"
	"testing"
)

func TestStrictBoolAcceptsOnlyExactBooleanSemantics(t *testing.T) {
	for _, raw := range []string{"true", "false", "\"true\"", "\"false\""} {
		var v strictBool
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			t.Fatalf("%s rejected: %v", raw, err)
		}
	}
	for _, raw := range []string{"1", "0", "\"yes\"", "\"False\"", "null", "{}"} {
		var v strictBool
		if err := json.Unmarshal([]byte(raw), &v); err == nil {
			t.Fatalf("%s unexpectedly accepted", raw)
		}
	}
}
