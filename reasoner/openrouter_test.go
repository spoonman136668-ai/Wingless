package reasoner

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvokePinsScientificControlsAndCapturesIdentity(t *testing.T) {
	var gotAuth string
	var gotPayload map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		gotAuth = req.Header.Get("Authorization")
		if err := json.NewDecoder(req.Body).Decode(&gotPayload); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
		  "model":"nvidia/nemotron-3-ultra-550b-a55b:free",
		  "provider":"NVIDIA",
		  "choices":[{"message":{"content":"bounded scientific analysis"},"finish_reason":"stop"}],
		  "usage":{"prompt_tokens":100,"completion_tokens":20,"completion_tokens_details":{"reasoning_tokens":10}}
		}`))
	}))
	defer srv.Close()

	cfg := DefaultConfig()
	client, err := NewClient("test-key", cfg)
	if err != nil {
		t.Fatal(err)
	}
	client.endpoint = srv.URL

	r := Request{
		Schema:                      RequestSchema,
		RequestID:                   "req-1",
		Project:                     "Wingless",
		ExperimentID:                "fixture-1",
		Role:                        RolePlan,
		FrontierSHA256:              strings.Repeat("a", 64),
		QualificationContractSHA256: strings.Repeat("b", 64),
		Context:                     "Use only evidence available before the held-out experiment.",
		MaxOutputTokens:             1024,
		TimeoutSeconds:              10,
		DataClass:                   "public-repository",
	}
	out, err := client.Invoke(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer test-key" {
		t.Fatalf("authorization header mismatch")
	}
	if gotPayload["model"] != DefaultModel {
		t.Fatalf("model not pinned: %#v", gotPayload["model"])
	}
	if gotPayload["seed"] != float64(0) || gotPayload["temperature"] != float64(0) {
		t.Fatalf("sampling controls not pinned: %#v", gotPayload)
	}
	reasoning := gotPayload["reasoning"].(map[string]any)
	if reasoning["effort"] != "high" || reasoning["exclude"] != true {
		t.Fatalf("reasoning controls mismatch: %#v", reasoning)
	}
	provider := gotPayload["provider"].(map[string]any)
	if provider["allow_fallbacks"] != false || provider["require_parameters"] != true {
		t.Fatalf("provider controls mismatch: %#v", provider)
	}
	if _, ok := provider["data_collection"]; ok {
		t.Fatalf("free route must not claim paid-route privacy controls: %#v", provider)
	}
	if _, ok := provider["zdr"]; ok {
		t.Fatalf("free route must not claim ZDR: %#v", provider)
	}
	if out.Provider != "NVIDIA" || out.ReturnedModel != DefaultModel {
		t.Fatalf("provenance not captured: %#v", out)
	}
	if out.RequestSHA256 == "" || out.ConfigSHA256 == "" || out.ResponseSHA256 == "" || out.QualificationIdentitySHA256 == "" {
		t.Fatalf("qualification hashes missing: %#v", out)
	}
}

func TestRequestRejectsMutableOrUnboundedInputs(t *testing.T) {
	r := Request{
		Schema:                      RequestSchema,
		RequestID:                   "req-2",
		Project:                     "Wingless",
		ExperimentID:                "fixture-2",
		Role:                        RolePlan,
		FrontierSHA256:              strings.Repeat("a", 64),
		QualificationContractSHA256: strings.Repeat("b", 64),
		Context:                     "x",
		MaxOutputTokens:             MaxCompletionTokens + 1,
		TimeoutSeconds:              10,
		DataClass:                   "public-repository",
	}
	if err := r.Validate(); err == nil {
		t.Fatal("expected output budget rejection")
	}
}

func TestFreeRouteRejectsPrivateDataClass(t *testing.T) {
	cfg := DefaultConfig()
	client, err := NewClient("test-key", cfg)
	if err != nil {
		t.Fatal(err)
	}
	r := Request{
		Schema:                      RequestSchema,
		RequestID:                   "req-private",
		Project:                     "Wingless",
		ExperimentID:                "private-fixture",
		Role:                        RolePlan,
		FrontierSHA256:              strings.Repeat("a", 64),
		QualificationContractSHA256: strings.Repeat("b", 64),
		Context:                     "private context must never reach free endpoint",
		MaxOutputTokens:             128,
		TimeoutSeconds:              10,
		DataClass:                   "private",
	}
	if _, err := client.Invoke(context.Background(), r); err == nil ||
		!strings.Contains(err.Error(), "public-repository") {
		t.Fatalf("expected free-route public-repository guard, got %v", err)
	}
}

func TestFreeRouteDoesNotSendResponseFormat(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
		  "model":"nvidia/nemotron-3-ultra-550b-a55b:free",
		  "provider":"NVIDIA",
		  "choices":[{"message":{"content":"{\"candidate_id\":\"A\",\"constraint_violation\":false,\"rationale\":\"bounded\"}"},"finish_reason":"stop"}],
		  "usage":{"prompt_tokens":10,"completion_tokens":10}
		}`))
	}))
	defer srv.Close()

	cfg := DefaultConfig()
	client, err := NewClient("test-key", cfg)
	if err != nil {
		t.Fatal(err)
	}
	client.endpoint = srv.URL
	r := Request{
		Schema:                      RequestSchema,
		RequestID:                   "req-json",
		Project:                     "Wingless",
		ExperimentID:                "fixture-json",
		Role:                        RolePlan,
		FrontierSHA256:              strings.Repeat("a", 64),
		QualificationContractSHA256: strings.Repeat("b", 64),
		Context:                     "Return the requested JSON only.",
		MaxOutputTokens:             128,
		TimeoutSeconds:              10,
		DataClass:                   "public-repository",
		ResponseJSONSchema: &JSONSchemaConstraint{
			Name: "fixture",
			Schema: map[string]any{
				"type": "object",
			},
		},
	}
	if _, err := client.Invoke(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	if _, ok := got["response_format"]; ok {
		t.Fatalf("free endpoint does not support response_format: %#v", got["response_format"])
	}
}
