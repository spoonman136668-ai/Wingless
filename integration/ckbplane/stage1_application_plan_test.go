package ckbplane

import (
	"encoding/json"
	"strings"
	"testing"
)

func validStage1ApplicationPlanFixture(t *testing.T) (Stage1CandidateRequest, Stage1ValidatedCandidate, string) {
	t.Helper()

	req, sourceText, raw := validStage1ValidatedCandidateFixture(t)
	candidate, err := ValidateStage1Candidate(req, sourceText, raw)
	if err != nil {
		t.Fatal(err)
	}
	return req, candidate, sourceText
}

func TestPlanStage1ExternalApplicationBindsPreimageAndPostimage(t *testing.T) {
	req, candidate, sourceText := validStage1ApplicationPlanFixture(t)

	plan, err := PlanStage1ExternalApplication(req, candidate, sourceText)
	if err != nil {
		t.Fatal(err)
	}

	if plan.WorkOrderID != req.Request.ParentWorkID ||
		plan.RequestID != req.Request.ID ||
		plan.Workspace != req.Request.Workspace {
		t.Fatalf("application identity mismatch: %+v", plan)
	}
	if plan.Path != req.ProductionPath {
		t.Fatalf("path = %q, want %q", plan.Path, req.ProductionPath)
	}
	if plan.ExpectedBeforeSHA256 != req.BeforeSHA256 {
		t.Fatalf("before sha = %q, want %q", plan.ExpectedBeforeSHA256, req.BeforeSHA256)
	}
	if plan.ExpectedAfterSHA256 != candidate.ContentSHA256 {
		t.Fatalf("after sha = %q, want %q", plan.ExpectedAfterSHA256, candidate.ContentSHA256)
	}
	if plan.Content != candidate.Operation.Content {
		t.Fatal("application content changed")
	}
	if plan.ApplicationAuthority != Stage1ApplicationAuthorityExternal {
		t.Fatalf("application authority = %q", plan.ApplicationAuthority)
	}
	if plan.AcceptanceAuthority != Stage1AcceptanceAuthorityExternalRequired {
		t.Fatalf("acceptance authority = %q", plan.AcceptanceAuthority)
	}
	if plan.TrustCeiling != "stage_1_only" {
		t.Fatalf("trust ceiling = %q", plan.TrustCeiling)
	}
}

func TestPlanStage1ExternalApplicationRejectsStaleSource(t *testing.T) {
	req, candidate, _ := validStage1ApplicationPlanFixture(t)
	stale := "package target\n\nfunc Value() int { return 99 }\n"

	if _, err := PlanStage1ExternalApplication(req, candidate, stale); err == nil {
		t.Fatal("stale current source accepted")
	}
}

func TestPlanStage1ExternalApplicationRejectsForgedValidatedCandidate(t *testing.T) {
	req, candidate, sourceText := validStage1ApplicationPlanFixture(t)

	tests := []struct {
		name   string
		mutate func(*Stage1ValidatedCandidate)
	}{
		{
			name: "wrong_operation",
			mutate: func(x *Stage1ValidatedCandidate) {
				x.Operation.Type = "shell"
			},
		},
		{
			name: "wrong_path",
			mutate: func(x *Stage1ValidatedCandidate) {
				x.Operation.Path = "pkg/other.go"
			},
		},
		{
			name: "wrong_preimage",
			mutate: func(x *Stage1ValidatedCandidate) {
				x.Operation.ExpectedBeforeSHA256 = strings.Repeat("f", 64)
			},
		},
		{
			name: "empty_content",
			mutate: func(x *Stage1ValidatedCandidate) {
				x.Operation.Content = ""
			},
		},
		{
			name: "nul_content",
			mutate: func(x *Stage1ValidatedCandidate) {
				x.Operation.Content = "a\x00b"
				x.ContentSHA256 = exactSHA256(x.Operation.Content)
			},
		},
		{
			name: "no_op",
			mutate: func(x *Stage1ValidatedCandidate) {
				x.Operation.Content = sourceText
				x.ContentSHA256 = exactSHA256(sourceText)
			},
		},
		{
			name: "wrong_postimage_hash",
			mutate: func(x *Stage1ValidatedCandidate) {
				x.ContentSHA256 = strings.Repeat("e", 64)
			},
		},
		{
			name: "oversized_content",
			mutate: func(x *Stage1ValidatedCandidate) {
				x.Operation.Content = strings.Repeat("x", maxStage1ReplacementBytes+1)
				x.ContentSHA256 = exactSHA256(x.Operation.Content)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			changed := candidate
			tc.mutate(&changed)
			if _, err := PlanStage1ExternalApplication(req, changed, sourceText); err == nil {
				t.Fatalf("forged validated candidate accepted: %s", tc.name)
			}
		})
	}
}

func TestPlanStage1ExternalApplicationRejectsTamperedRequestBinding(t *testing.T) {
	req, candidate, sourceText := validStage1ApplicationPlanFixture(t)

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
			name: "path_changed",
			mutate: func(x *Stage1CandidateRequest) {
				x.ProductionPath = "pkg/other.go"
			},
		},
		{
			name: "preimage_changed",
			mutate: func(x *Stage1CandidateRequest) {
				x.BeforeSHA256 = strings.Repeat("f", 64)
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

			if _, err := PlanStage1ExternalApplication(changed, candidate, sourceText); err == nil {
				t.Fatalf("tampered request accepted: %s", tc.name)
			}
		})
	}
}

func TestPlanStage1ExternalApplicationDoesNotMutateInputs(t *testing.T) {
	req, candidate, sourceText := validStage1ApplicationPlanFixture(t)

	beforeSource := sourceText
	beforeContent := candidate.Operation.Content
	beforeHash := candidate.ContentSHA256

	if _, err := PlanStage1ExternalApplication(req, candidate, sourceText); err != nil {
		t.Fatal(err)
	}
	if sourceText != beforeSource {
		t.Fatal("planner mutated current source")
	}
	if candidate.Operation.Content != beforeContent || candidate.ContentSHA256 != beforeHash {
		t.Fatal("planner mutated validated candidate")
	}
}
