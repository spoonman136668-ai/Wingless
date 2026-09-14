package broker

import "fmt"

const (
	Stage1TaskClassBoundedSingleFileRepair = "bounded_single_file_repair"

	Stage1PreferredModel = "qwen3-coder-30ba3b-q3km"
	Stage1SecondaryModel = "devstral-small-2-24b-iq3m"
	Stage1TrustCeiling   = "stage_1_only"
)

type Stage1LocalTask struct {
	TaskClass                    string
	ProductionFiles              int
	CandidateOnly                bool
	StructuredCandidate          bool
	ExactPreimage                bool
	TestsMutable                 bool
	ExternalApplication          bool
	ExternalAcceptance           bool
	ModelRepairRetries           int
	ChangesRoutingOrArchitecture bool
	ChangesAuthority             bool
	LiveAction                   bool
	ArbitraryShell               bool
	ModelAppliesEdits            bool
	ModelControlsTests           bool
}

type Stage1LocalDecision struct {
	Allowed        bool   `json:"allowed"`
	PreferredModel string `json:"preferred_model,omitempty"`
	SecondaryModel string `json:"secondary_model,omitempty"`
	TrustCeiling   string `json:"trust_ceiling"`
	Reason         string `json:"reason"`
}

func denyStage1(reason string) (Stage1LocalDecision, error) {
	return Stage1LocalDecision{
		Allowed:      false,
		TrustCeiling: Stage1TrustCeiling,
		Reason:       reason,
	}, fmt.Errorf("stage1 local selection denied: %s", reason)
}

// SelectStage1LocalModel is an unwired, deterministic eligibility contract.
//
// It does not invoke a backend, apply edits, run tests, grant acceptance, mutate
// authority, or alter broker routing. A caller may use an allowed decision only
// to choose a candidate generator; external deterministic acceptance remains
// required.
func SelectStage1LocalModel(task Stage1LocalTask) (Stage1LocalDecision, error) {
	switch {
	case task.TaskClass != Stage1TaskClassBoundedSingleFileRepair:
		return denyStage1("task_class_not_allowed")
	case task.ProductionFiles != 1:
		return denyStage1("production_scope_not_single_file")
	case !task.CandidateOnly:
		return denyStage1("candidate_only_required")
	case !task.StructuredCandidate:
		return denyStage1("structured_candidate_required")
	case !task.ExactPreimage:
		return denyStage1("exact_preimage_required")
	case task.TestsMutable:
		return denyStage1("tests_must_be_immutable")
	case !task.ExternalApplication:
		return denyStage1("external_application_required")
	case !task.ExternalAcceptance:
		return denyStage1("external_acceptance_required")
	case task.ModelRepairRetries != 0:
		return denyStage1("model_repair_retries_not_allowed")
	case task.ChangesRoutingOrArchitecture:
		return denyStage1("routing_or_architecture_change_not_allowed")
	case task.ChangesAuthority:
		return denyStage1("authority_change_not_allowed")
	case task.LiveAction:
		return denyStage1("live_action_not_allowed")
	case task.ArbitraryShell:
		return denyStage1("arbitrary_shell_not_allowed")
	case task.ModelAppliesEdits:
		return denyStage1("model_edit_authority_not_allowed")
	case task.ModelControlsTests:
		return denyStage1("model_test_authority_not_allowed")
	}

	return Stage1LocalDecision{
		Allowed:        true,
		PreferredModel: Stage1PreferredModel,
		SecondaryModel: Stage1SecondaryModel,
		TrustCeiling:   Stage1TrustCeiling,
		Reason:         "qualified_stage1_candidate_generation",
	}, nil
}
