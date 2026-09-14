package ckbplane

import (
	"strings"
	"testing"
)

func validStage1PostApplicationFixture(t *testing.T) (Stage1CandidateRequest, Stage1ExternalApplicationPlan) {
	t.Helper()

	req, candidate, sourceText := validStage1ApplicationPlanFixture(t)
	plan, err := PlanStage1ExternalApplication(req, candidate, sourceText)
	if err != nil {
		t.Fatal(err)
	}
	return req, plan
}

func TestVerifyStage1ExternalApplicationAcceptsExactObservedPostimage(t *testing.T) {
	req, plan := validStage1PostApplicationFixture(t)

	got, err := VerifyStage1ExternalApplication(req, plan, plan.Content)
	if err != nil {
		t.Fatal(err)
	}

	if got.WorkOrderID != plan.WorkOrderID ||
		got.RequestID != plan.RequestID ||
		got.Workspace != plan.Workspace ||
		got.Path != plan.Path {
		t.Fatalf("post-application identity mismatch: %+v", got)
	}
	if got.ExpectedBeforeSHA256 != plan.ExpectedBeforeSHA256 ||
		got.ExpectedAfterSHA256 != plan.ExpectedAfterSHA256 ||
		got.ObservedAfterSHA256 != plan.ExpectedAfterSHA256 {
		t.Fatalf("post-application hashes mismatch: %+v", got)
	}
	if got.Status != Stage1PostApplicationStatusVerified {
		t.Fatalf("status = %q", got.Status)
	}
	if got.Acceptance != Stage1AcceptanceAuthorityExternalRequired {
		t.Fatalf("acceptance = %q", got.Acceptance)
	}
	if got.TrustCeiling != "stage_1_only" {
		t.Fatalf("trust ceiling = %q", got.TrustCeiling)
	}
}

func TestVerifyStage1ExternalApplicationRejectsWrongObservedBytes(t *testing.T) {
	req, plan := validStage1PostApplicationFixture(t)

	wrong := plan.Content + "\n"
	if _, err := VerifyStage1ExternalApplication(req, plan, wrong); err == nil {
		t.Fatal("wrong observed postimage accepted")
	}
}

func TestVerifyStage1ExternalApplicationRejectsNULObservedSource(t *testing.T) {
	req, plan := validStage1PostApplicationFixture(t)

	if _, err := VerifyStage1ExternalApplication(req, plan, "a\x00b"); err == nil {
		t.Fatal("NUL-containing observed source accepted")
	}
}

func TestVerifyStage1ExternalApplicationRejectsTamperedPlan(t *testing.T) {
	req, plan := validStage1PostApplicationFixture(t)

	tests := []struct {
		name   string
		mutate func(*Stage1ExternalApplicationPlan)
	}{
		{
			name: "wrong_work_order",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.WorkOrderID = "OTHER"
			},
		},
		{
			name: "wrong_request",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.RequestID = "OTHER"
			},
		},
		{
			name: "wrong_workspace",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.Workspace = `C:\other`
			},
		},
		{
			name: "wrong_path",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.Path = "pkg/other.go"
			},
		},
		{
			name: "wrong_before_hash",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.ExpectedBeforeSHA256 = strings.Repeat("f", 64)
			},
		},
		{
			name: "wrong_after_hash",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.ExpectedAfterSHA256 = strings.Repeat("e", 64)
			},
		},
		{
			name: "content_changed",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.Content += "\n"
			},
		},
		{
			name: "application_authority_widened",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.ApplicationAuthority = "model"
			},
		},
		{
			name: "acceptance_authority_widened",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.AcceptanceAuthority = "model"
			},
		},
		{
			name: "trust_ceiling_widened",
			mutate: func(x *Stage1ExternalApplicationPlan) {
				x.TrustCeiling = "stage_2"
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			changed := plan
			tc.mutate(&changed)

			if _, err := VerifyStage1ExternalApplication(req, changed, plan.Content); err == nil {
				t.Fatalf("tampered plan accepted: %s", tc.name)
			}
		})
	}
}

func TestVerifyStage1ExternalApplicationRejectsTamperedRequest(t *testing.T) {
	req, plan := validStage1PostApplicationFixture(t)
	req.Selection.Allowed = false

	if _, err := VerifyStage1ExternalApplication(req, plan, plan.Content); err == nil {
		t.Fatal("tampered Stage-1 request accepted")
	}
}

func TestVerifyStage1ExternalApplicationDoesNotGrantAcceptance(t *testing.T) {
	req, plan := validStage1PostApplicationFixture(t)

	got, err := VerifyStage1ExternalApplication(req, plan, plan.Content)
	if err != nil {
		t.Fatal(err)
	}
	if got.Acceptance != "external_required" {
		t.Fatalf("post-application verifier granted acceptance: %q", got.Acceptance)
	}
	if got.Status == "accepted" || got.Status == "promoted" {
		t.Fatalf("post-application verifier emitted authority-bearing status: %q", got.Status)
	}
}
