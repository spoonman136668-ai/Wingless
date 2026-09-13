package inference

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLocalHTTPJSONSchemaConstraint(t *testing.T) {
	const model = "model"
	var seen map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			http.NotFound(w, r)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(body, &seen); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"model":"model","choices":[{"message":{"content":"{\"value\":\"ok\"}"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1}}`)
	}))
	defer srv.Close()

	endpoint := strings.Replace(srv.URL, "localhost", "127.0.0.1", 1)
	backend, err := NewLocalHTTP("local", model, endpoint, []string{"code"})
	if err != nil {
		t.Fatal(err)
	}

	schema := json.RawMessage(`{"type":"object","additionalProperties":false,"properties":{"value":{"type":"string"}},"required":["value"]}`)
	req := Request{
		ID: "constraint",
		ParentWorkID: "parent",
		Role: "code",
		Context: "return json",
		MaxContextBytes: 4096,
		MaxOutputTokens: 32,
		Deadline: time.Now().Add(5 * time.Second),
		Workspace: "none",
		Capabilities: []string{"code"},
		OutputConstraint: &OutputConstraint{Kind: "json_schema", Schema: schema},
	}
	got, err := backend.Invoke(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "completed" || got.Text != `{"value":"ok"}` {
		t.Fatalf("unexpected result: %#v", got)
	}

	rf, ok := seen["response_format"].(map[string]any)
	if !ok || rf["type"] != "json_schema" {
		t.Fatalf("response_format missing or wrong: %#v", seen["response_format"])
	}
	js, ok := rf["json_schema"].(map[string]any)
	if !ok || js["name"] != "wingless_output" || js["strict"] != true {
		t.Fatalf("json_schema metadata wrong: %#v", rf["json_schema"])
	}
	gotSchema, ok := js["schema"].(map[string]any)
	if !ok || gotSchema["type"] != "object" {
		t.Fatalf("schema not forwarded: %#v", js["schema"])
	}
}

func TestLocalHTTPRejectsUnsupportedConstraintKind(t *testing.T) {
	backend, err := NewLocalHTTP("local", "model", "http://127.0.0.1:1", []string{"code"})
	if err != nil {
		t.Fatal(err)
	}
	req := Request{
		ID: "constraint-kind",
		ParentWorkID: "parent",
		Role: "code",
		Context: "x",
		MaxContextBytes: 4096,
		MaxOutputTokens: 8,
		Deadline: time.Now().Add(time.Second),
		Workspace: "none",
		Capabilities: []string{"code"},
		OutputConstraint: &OutputConstraint{Kind: "other", Schema: json.RawMessage(`{}`)},
	}
	_, err = backend.Invoke(context.Background(), req)
	if err == nil || !strings.Contains(err.Error(), "invalid bounded output constraint") {
		t.Fatalf("expected request validation failure, got %v", err)
	}
}
