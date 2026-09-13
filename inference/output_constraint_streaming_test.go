package inference

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestLocalHTTPRejectsConstraintWithStreaming(t *testing.T) {
	backend, err := NewLocalHTTP("local", "model", "http://127.0.0.1:1", []string{"code"})
	if err != nil {
		t.Fatal(err)
	}
	backend.streaming = true
	req := Request{
		ID: "constraint-stream",
		ParentWorkID: "parent",
		Role: "code",
		Context: "x",
		MaxContextBytes: 4096,
		MaxOutputTokens: 8,
		Deadline: time.Now().Add(time.Second),
		Workspace: "none",
		Capabilities: []string{"code"},
		OutputConstraint: &OutputConstraint{Kind: "json_schema", Schema: json.RawMessage(`{"type":"object"}`)},
	}
	_, err = backend.Invoke(context.Background(), req)
	if err == nil || !strings.Contains(err.Error(), "output constraint unavailable for streaming") {
		t.Fatalf("expected streaming constraint rejection, got %v", err)
	}
}
