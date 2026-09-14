package ckbplane

import (
	"strings"
	"testing"

	"github.com/spoonman136668-ai/Wingless/broker"
)

func validStage1ProjectionFixture() (WorkOrderContract, WorkspaceBinding, Stage1AuthorityEnvelope) {
	order := WorkOrderContract{
		ID:                   "WO_STAGE1_001",
		Title:                "bounded single-file repair",
		BaselineSHA:          strings.Repeat("a", 40),
		WorkerBranch:         "worker-stage1",
		AllowedPaths:         []string{"pkg/target.go"},
		ProtectAcceptedRefs:  true,
		RequireCleanBaseline: true,
	}
	workspace := WorkspaceBinding{
		WorkOrderID: order.ID,
		Attempt:     1,
		Path:        `C:\isolated\stage1`,
		BaselineSHA: order.BaselineSHA,
	}
	auth := Stage1AuthorityEnvelope{
		WorkOrderID:          order.ID,
		Attempt:              workspace.Attempt,
		BaselineSHA:          order.BaselineSHA,
		ProductionPath:       order.AllowedPaths[0],
		BeforeSHA256:         strings.Repeat("b", 64),
		TaskClass:            broker.Stage1TaskClassBoundedSingleFileRepair,
		CandidateFormat:      Stage1CandidateFormatWriteTextJSONSchema,
		ExactPreimage:        true,
		ApplicationAuthority: Stage1ApplicationAuthorityExternal,
		AcceptanceAuthority:  Stage1AcceptanceAuthorityExternalRequired,
	}
	return order, workspace, auth
}

func TestDecodeStage1AuthorityEnvelopeRejectsUnknownFieldsAndTrailingJSON(t *testing.T) {
	valid := `{
		"work_order_id":"WO_STAGE1_001",
		"attempt":1,
		"baseline_sha":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"production_path":"pkg/target.go",
		"before_sha256":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"task_class":"bounded_single_file_repair",
		"candidate_format":"write_text_json_schema",
		"exact_preimage":true,
		"tests_mutable":false,
		"application_authority":"external",
		"acceptance_authority":"external_required",
		"model_repair_retries":0,
		"changes_routing_architecture":false,
		"changes_authority":false,
		"live_action":false,
		"arbitrary_shell":false,
		"model_edit_authority":false,
		"model_test_authority":false
	}`

	if _, err := DecodeStage1AuthorityEnvelope([]byte(valid)); err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeStage1AuthorityEnvelope([]byte(strings.TrimSuffix(valid, "}") + `,"unknown":true}`)); err == nil {
		t.Fatal("unknown Stage-1 authority field accepted")
	}
	if _, err := DecodeStage1AuthorityEnvelope([]byte(valid + `{}`)); err == nil {
		t.Fatal("trailing JSON accepted")
	}
}

func TestProjectStage1LocalTaskAllowsBoundAuthorityEnvelope(t *testing.T) {
	order, workspace, auth := validStage1ProjectionFixture()

	task, decision, err := ProjectStage1LocalTask(order, workspace, auth)
	if err != nil {
		t.Fatal(err)
	}
	if !decision.Allowed {
		t.Fatalf("valid Stage-1 authority envelope denied: %+v", decision)
	}
	if decision.PreferredModel != broker.Stage1PreferredModel {
		t.Fatalf("preferred model = %q, want %q", decision.PreferredModel, broker.Stage1PreferredModel)
	}
	if decision.TrustCeiling != broker.Stage1TrustCeiling {
		t.Fatalf("trust ceiling = %q, want %q", decision.TrustCeiling, broker.Stage1TrustCeiling)
	}
	if task.ProductionFiles != 1 || !task.CandidateOnly || !task.StructuredCandidate || !task.ExactPreimage {
		t.Fatalf("unexpected projected task: %+v", task)
	}
	if task.TestsMutable || task.ModelRepairRetries != 0 ||
		task.ChangesRoutingOrArchitecture || task.ChangesAuthority ||
		task.LiveAction || task.ArbitraryShell ||
		task.ModelAppliesEdits || task.ModelControlsTests {
		t.Fatalf("projected task widened Stage-1 authority: %+v", task)
	}
}

func TestProjectStage1LocalTaskFailsClosedOutsideAuthorityEnvelope(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*WorkOrderContract, *WorkspaceBinding, *Stage1AuthorityEnvelope)
	}{
		{"unprotected_accepted_refs", func(o *WorkOrderContract, _ *WorkspaceBinding, _ *Stage1AuthorityEnvelope) {
			o.ProtectAcceptedRefs = false
		}},
		{"dirty_baseline_allowed", func(o *WorkOrderContract, _ *WorkspaceBinding, _ *Stage1AuthorityEnvelope) {
			o.RequireCleanBaseline = false
		}},
		{"work_order_identity_mismatch", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) { a.WorkOrderID = "OTHER" }},
		{"workspace_identity_mismatch", func(_ *WorkOrderContract, w *WorkspaceBinding, _ *Stage1AuthorityEnvelope) { w.WorkOrderID = "OTHER" }},
		{"attempt_mismatch", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) { a.Attempt = 2 }},
		{"workspace_baseline_mismatch", func(_ *WorkOrderContract, w *WorkspaceBinding, _ *Stage1AuthorityEnvelope) {
			w.BaselineSHA = strings.Repeat("c", 40)
		}},
		{"authority_baseline_mismatch", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.BaselineSHA = strings.Repeat("c", 40)
		}},
		{"multiple_allowed_paths", func(o *WorkOrderContract, _ *WorkspaceBinding, _ *Stage1AuthorityEnvelope) {
			o.AllowedPaths = append(o.AllowedPaths, "pkg/other.go")
		}},
		{"absolute_production_path", func(o *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			o.AllowedPaths = []string{`C:\target.go`}
			a.ProductionPath = `C:\target.go`
		}},
		{"production_path_not_exactly_bound", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.ProductionPath = "pkg/other.go"
		}},
		{"production_path_forbidden", func(o *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			o.ForbiddenPaths = []string{a.ProductionPath}
		}},
		{"invalid_preimage_sha256", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.BeforeSHA256 = "not-a-sha"
		}},
		{"wrong_task_class", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.TaskClass = "routing_repair"
		}},
		{"wrong_candidate_format", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.CandidateFormat = "free_text"
		}},
		{"no_exact_preimage", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) { a.ExactPreimage = false }},
		{"tests_mutable", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) { a.TestsMutable = true }},
		{"nonexternal_application", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.ApplicationAuthority = "model"
		}},
		{"nonexternal_acceptance", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.AcceptanceAuthority = "model"
		}},
		{"model_repair_retry", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) { a.ModelRepairRetries = 1 }},
		{"routing_architecture_change", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.ChangesRoutingArchitecture = true
		}},
		{"authority_change", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) { a.ChangesAuthority = true }},
		{"live_action", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) { a.LiveAction = true }},
		{"arbitrary_shell", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) { a.ArbitraryShell = true }},
		{"model_edit_authority", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.ModelEditAuthority = true
		}},
		{"model_test_authority", func(_ *WorkOrderContract, _ *WorkspaceBinding, a *Stage1AuthorityEnvelope) {
			a.ModelTestAuthority = true
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			order, workspace, auth := validStage1ProjectionFixture()
			tc.mutate(&order, &workspace, &auth)

			task, decision, err := ProjectStage1LocalTask(order, workspace, auth)
			if err == nil {
				t.Fatalf("out-of-envelope authority unexpectedly accepted: task=%+v decision=%+v", task, decision)
			}
			if decision.Allowed {
				t.Fatalf("rejected authority produced allowed decision: %+v", decision)
			}
		})
	}
}
