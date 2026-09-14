package ckbplane

import (
	"encoding/json"
	"strings"
	"testing"
)

func validStage1ValidatedCandidateFixture(t *testing.T) (Stage1CandidateRequest, string, string) {
	t.Helper()

	order, workspace, auth, sourceText, bounds := validStage1RequestFixture(t)
	req, err := BuildStage1CandidateRequest(order, workspace, auth, sourceText, bounds)
	if err != nil {
		t.Fatal(err)
	}

	rawBytes, err := json.Marshal(Stage1WriteTextCandidate{
		Type:                 "write_text",
		Path:                 auth.ProductionPath,
		ExpectedBeforeSHA256: auth.BeforeSHA256,
		Content:              "package target\n\nfunc Value() int { return 2 }\n",
	})
	if err != nil {
		t.Fatal(err)
	}
	return req, sourceText, string(rawBytes)
}

func TestValidateStage1CandidateAcceptsExactBoundOperation(t *testing.T) {
	req, sourceText, raw := validStage1ValidatedCandidateFixture(t)

	got, err := ValidateStage1Candidate(req, sourceText, raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.Operation.Type != "write_text" ||
		got.Operation.Path != req.ProductionPath ||
		got.Operation.ExpectedBeforeSHA256 != req.BeforeSHA256 {
		t.Fatalf("unexpected validated operation: %+v", got.Operation)
	}
	if got.ContentSHA256 == "" || got.ContentSHA256 == req.BeforeSHA256 {
		t.Fatalf("unexpected content hash: %q", got.ContentSHA256)
	}
}

func TestValidateStage1CandidateRejectsAmbiguousOrWrappedOutput(t *testing.T) {
	req, sourceText, raw := validStage1ValidatedCandidateFixture(t)

	tests := []struct {
		name string
		raw  string
	}{
		{"markdown_fence", "```json\n" + raw + "\n```"},
		{"prefix_text", "candidate:\n" + raw},
		{"suffix_text", raw + "\nthanks"},
		{"second_object", raw + "\n{}"},
		{"array_wrapper", "[" + raw + "]"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ValidateStage1Candidate(req, sourceText, tc.raw); err == nil {
				t.Fatalf("ambiguous/wrapped candidate accepted: %s", tc.name)
			}
		})
	}
}

func TestValidateStage1CandidateRejectsUnknownFields(t *testing.T) {
	req, sourceText, raw := validStage1ValidatedCandidateFixture(t)

	var obj map[string]any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		t.Fatal(err)
	}
	obj["shell"] = "powershell"
	withUnknown, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ValidateStage1Candidate(req, sourceText, string(withUnknown)); err == nil {
		t.Fatal("candidate with unknown authority-like field accepted")
	}
}

func TestValidateStage1CandidateRejectsBindingViolations(t *testing.T) {
	req, sourceText, raw := validStage1ValidatedCandidateFixture(t)

	var base Stage1WriteTextCandidate
	if err := json.Unmarshal([]byte(raw), &base); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		mutate func(*Stage1WriteTextCandidate)
	}{
		{"wrong_operation", func(x *Stage1WriteTextCandidate) { x.Type = "shell" }},
		{"wrong_path", func(x *Stage1WriteTextCandidate) { x.Path = "pkg/other.go" }},
		{"wrong_preimage", func(x *Stage1WriteTextCandidate) { x.ExpectedBeforeSHA256 = strings.Repeat("f", 64) }},
		{"empty_content", func(x *Stage1WriteTextCandidate) { x.Content = "" }},
		{"nul_content", func(x *Stage1WriteTextCandidate) { x.Content = "a\x00b" }},
		{"no_op_content", func(x *Stage1WriteTextCandidate) { x.Content = sourceText }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			op := base
			tc.mutate(&op)
			b, err := json.Marshal(op)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := ValidateStage1Candidate(req, sourceText, string(b)); err == nil {
				t.Fatalf("binding violation accepted: %s", tc.name)
			}
		})
	}
}

func TestValidateStage1CandidateRejectsStaleCurrentSource(t *testing.T) {
	req, _, raw := validStage1ValidatedCandidateFixture(t)
	stale := "package target\n\nfunc Value() int { return 99 }\n"

	if _, err := ValidateStage1Candidate(req, stale, raw); err == nil {
		t.Fatal("stale current source accepted")
	}
}

func TestValidateStage1CandidateRejectsTamperedBuiltRequest(t *testing.T) {
	req, sourceText, raw := validStage1ValidatedCandidateFixture(t)

	tests := []struct {
		name   string
		mutate func(*Stage1CandidateRequest)
	}{
		{
			name: "selection_denied",
			mutate: func(x *Stage1CandidateRequest) {
				x.Selection.Allowed = false
			},
		},
		{
			name: "trust_ceiling_widened",
			mutate: func(x *Stage1CandidateRequest) {
				x.Selection.TrustCeiling = "stage_2"
			},
		},
		{
			name: "path_binding_changed",
			mutate: func(x *Stage1CandidateRequest) {
				x.ProductionPath = "pkg/other.go"
			},
		},
		{
			name: "preimage_binding_changed",
			mutate: func(x *Stage1CandidateRequest) {
				x.BeforeSHA256 = strings.Repeat("f", 64)
			},
		},
		{
			name: "schema_removed",
			mutate: func(x *Stage1CandidateRequest) {
				x.Request.OutputConstraint = nil
			},
		},
		{
			name: "schema_tampered",
			mutate: func(x *Stage1CandidateRequest) {
				x.Request.OutputConstraint.Schema = json.RawMessage(`{"type":"object"}`)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			changed := req
			if req.Request.OutputConstraint != nil {
				copiedConstraint := *req.Request.OutputConstraint
				copiedSchema := append([]byte(nil), req.Request.OutputConstraint.Schema...)
				copiedConstraint.Schema = copiedSchema
				changed.Request.OutputConstraint = &copiedConstraint
			}
			tc.mutate(&changed)

			if _, err := ValidateStage1Candidate(changed, sourceText, raw); err == nil {
				t.Fatalf("tampered built request accepted: %s", tc.name)
			}
		})
	}
}

func TestValidateStage1CandidateDoesNotMutateInputs(t *testing.T) {
	req, sourceText, raw := validStage1ValidatedCandidateFixture(t)
	beforeSource := sourceText
	beforeRaw := raw

	if _, err := ValidateStage1Candidate(req, sourceText, raw); err != nil {
		t.Fatal(err)
	}
	if sourceText != beforeSource {
		t.Fatal("validator mutated source input")
	}
	if raw != beforeRaw {
		t.Fatal("validator mutated raw candidate input")
	}
}
