package broker

import "testing"

func validStage1LocalTask() Stage1LocalTask {
	return Stage1LocalTask{
		TaskClass:           Stage1TaskClassBoundedSingleFileRepair,
		ProductionFiles:     1,
		CandidateOnly:       true,
		StructuredCandidate: true,
		ExactPreimage:       true,
		ExternalApplication: true,
		ExternalAcceptance:  true,
		ModelRepairRetries:  0,
	}
}

func TestStage1LocalSelectionAllowsOnlyQualifiedEnvelope(t *testing.T) {
	d, err := SelectStage1LocalModel(validStage1LocalTask())
	if err != nil {
		t.Fatal(err)
	}
	if !d.Allowed {
		t.Fatal("qualified Stage-1 envelope was denied")
	}
	if d.PreferredModel != Stage1PreferredModel {
		t.Fatalf("preferred model = %q, want %q", d.PreferredModel, Stage1PreferredModel)
	}
	if d.SecondaryModel != Stage1SecondaryModel {
		t.Fatalf("secondary model = %q, want %q", d.SecondaryModel, Stage1SecondaryModel)
	}
	if d.TrustCeiling != Stage1TrustCeiling {
		t.Fatalf("trust ceiling = %q, want %q", d.TrustCeiling, Stage1TrustCeiling)
	}
	if d.Reason != "qualified_stage1_candidate_generation" {
		t.Fatalf("reason = %q", d.Reason)
	}
}

func TestStage1LocalSelectionRejectsAnythingOutsideEnvelope(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Stage1LocalTask)
		reason string
	}{
		{
			name: "unknown_task_class",
			mutate: func(x *Stage1LocalTask) {
				x.TaskClass = "routing_repair"
			},
			reason: "task_class_not_allowed",
		},
		{
			name: "zero_production_files",
			mutate: func(x *Stage1LocalTask) {
				x.ProductionFiles = 0
			},
			reason: "production_scope_not_single_file",
		},
		{
			name: "multiple_production_files",
			mutate: func(x *Stage1LocalTask) {
				x.ProductionFiles = 2
			},
			reason: "production_scope_not_single_file",
		},
		{
			name: "not_candidate_only",
			mutate: func(x *Stage1LocalTask) {
				x.CandidateOnly = false
			},
			reason: "candidate_only_required",
		},
		{
			name: "unstructured_candidate",
			mutate: func(x *Stage1LocalTask) {
				x.StructuredCandidate = false
			},
			reason: "structured_candidate_required",
		},
		{
			name: "no_exact_preimage",
			mutate: func(x *Stage1LocalTask) {
				x.ExactPreimage = false
			},
			reason: "exact_preimage_required",
		},
		{
			name: "tests_mutable",
			mutate: func(x *Stage1LocalTask) {
				x.TestsMutable = true
			},
			reason: "tests_must_be_immutable",
		},
		{
			name: "model_applies_candidate",
			mutate: func(x *Stage1LocalTask) {
				x.ExternalApplication = false
			},
			reason: "external_application_required",
		},
		{
			name: "model_grants_acceptance",
			mutate: func(x *Stage1LocalTask) {
				x.ExternalAcceptance = false
			},
			reason: "external_acceptance_required",
		},
		{
			name: "model_repair_retry",
			mutate: func(x *Stage1LocalTask) {
				x.ModelRepairRetries = 1
			},
			reason: "model_repair_retries_not_allowed",
		},
		{
			name: "routing_or_architecture_change",
			mutate: func(x *Stage1LocalTask) {
				x.ChangesRoutingOrArchitecture = true
			},
			reason: "routing_or_architecture_change_not_allowed",
		},
		{
			name: "authority_change",
			mutate: func(x *Stage1LocalTask) {
				x.ChangesAuthority = true
			},
			reason: "authority_change_not_allowed",
		},
		{
			name: "live_action",
			mutate: func(x *Stage1LocalTask) {
				x.LiveAction = true
			},
			reason: "live_action_not_allowed",
		},
		{
			name: "arbitrary_shell",
			mutate: func(x *Stage1LocalTask) {
				x.ArbitraryShell = true
			},
			reason: "arbitrary_shell_not_allowed",
		},
		{
			name: "model_edit_authority",
			mutate: func(x *Stage1LocalTask) {
				x.ModelAppliesEdits = true
			},
			reason: "model_edit_authority_not_allowed",
		},
		{
			name: "model_test_authority",
			mutate: func(x *Stage1LocalTask) {
				x.ModelControlsTests = true
			},
			reason: "model_test_authority_not_allowed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := validStage1LocalTask()
			tc.mutate(&task)

			d, err := SelectStage1LocalModel(task)
			if err == nil {
				t.Fatalf("expected rejection: %+v", d)
			}
			if d.Allowed {
				t.Fatalf("rejected task marked allowed: %+v", d)
			}
			if d.Reason != tc.reason {
				t.Fatalf("reason = %q, want %q", d.Reason, tc.reason)
			}
			if d.PreferredModel != "" || d.SecondaryModel != "" {
				t.Fatalf("rejected decision leaked model selection: %+v", d)
			}
			if d.TrustCeiling != Stage1TrustCeiling {
				t.Fatalf("trust ceiling = %q, want %q", d.TrustCeiling, Stage1TrustCeiling)
			}
		})
	}
}

func TestStage1LocalSelectionZeroValueFailsClosed(t *testing.T) {
	d, err := SelectStage1LocalModel(Stage1LocalTask{})
	if err == nil {
		t.Fatalf("zero-value task unexpectedly allowed: %+v", d)
	}
	if d.Allowed {
		t.Fatalf("zero-value task marked allowed: %+v", d)
	}
	if d.Reason != "task_class_not_allowed" {
		t.Fatalf("reason = %q", d.Reason)
	}
}
