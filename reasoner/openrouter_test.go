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
		  "model":"nvidia/nemotron-3-ultra-550b-a55b",
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
	if provider["allow_fallbacks"] != false || provider["data_collection"] != "deny" || provider["zdr"] != true || provider["require_parameters"] != true {
		t.Fatalf("provider controls mismatch: %#v", provider)
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
	}
	if err := r.Validate(); err == nil {
		t.Fatal("expected output budget rejection")
	}
}

func TestConfigFailsClosedWithoutPrivacyControls(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ZeroDataRetention = false
	if _, err := NewClient("test-key", cfg); err == nil {
		t.Fatal("expected ZDR rejection")
	}

	cfg = DefaultConfig()
	cfg.DataCollection = "allow"
	if _, err := NewClient("test-key", cfg); err == nil {
		t.Fatal("expected data collection rejection")
	}
}
