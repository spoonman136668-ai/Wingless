package ckbplane

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/spoonman136668-ai/Wingless/broker"
)

const (
	Stage1CandidateFormatWriteTextJSONSchema  = "write_text_json_schema"
	Stage1ApplicationAuthorityExternal        = "external"
	Stage1AcceptanceAuthorityExternalRequired = "external_required"
)

var sha256Lower = regexp.MustCompile(`^[0-9a-f]{64}$`)

// Stage1AuthorityEnvelope is inactive, plane-owned eligibility metadata for a
// future Stage-1 local-model integration. It is not model output and it does not
// itself authorize inference, application, acceptance, promotion, or live use.
//
// The envelope is deliberately bound to the existing frozen WorkOrderContract
// and WorkspaceBinding. ProjectStage1LocalTask fails closed unless the work
// order exposes exactly one production path and all Stage-1 trust boundaries
// are explicit.
type Stage1AuthorityEnvelope struct {
	WorkOrderID                string `json:"work_order_id"`
	Attempt                    int    `json:"attempt"`
	BaselineSHA                string `json:"baseline_sha"`
	ProductionPath             string `json:"production_path"`
	BeforeSHA256               string `json:"before_sha256"`
	TaskClass                  string `json:"task_class"`
	CandidateFormat            string `json:"candidate_format"`
	ExactPreimage              bool   `json:"exact_preimage"`
	TestsMutable               bool   `json:"tests_mutable"`
	ApplicationAuthority       string `json:"application_authority"`
	AcceptanceAuthority        string `json:"acceptance_authority"`
	ModelRepairRetries         int    `json:"model_repair_retries"`
	ChangesRoutingArchitecture bool   `json:"changes_routing_architecture"`
	ChangesAuthority           bool   `json:"changes_authority"`
	LiveAction                 bool   `json:"live_action"`
	ArbitraryShell             bool   `json:"arbitrary_shell"`
	ModelEditAuthority         bool   `json:"model_edit_authority"`
	ModelTestAuthority         bool   `json:"model_test_authority"`
}

func DecodeStage1AuthorityEnvelope(data []byte) (Stage1AuthorityEnvelope, error) {
	var e Stage1AuthorityEnvelope
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&e); err != nil {
		return e, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		if err == nil {
			return e, errors.New("exactly one Stage-1 authority envelope required")
		}
		return e, err
	}
	return e, nil
}

// validStage1ProductionPath validates repository paths, not host filesystem
// paths. Git/repository paths use slash semantics on every platform.
func validStage1ProductionPath(p string) bool {
	if p == "" || strings.HasPrefix(p, "/") || strings.ContainsAny(p, `\:`) {
		return false
	}
	clean := path.Clean(p)
	return clean == p && clean != "." && clean != ".." && !strings.HasPrefix(clean, "../")
}

func containsExact(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}

// ProjectStage1LocalTask validates inactive plane-owned metadata and projects it
// into the broker's pure Stage-1 local-model eligibility contract.
//
// This function does not invoke any backend and does not alter existing
// ckb-plane adapter translation or broker routing.
func ProjectStage1LocalTask(
	order WorkOrderContract,
	workspace WorkspaceBinding,
	auth Stage1AuthorityEnvelope,
) (broker.Stage1LocalTask, broker.Stage1LocalDecision, error) {
	var task broker.Stage1LocalTask
	var decision broker.Stage1LocalDecision

	if !safeID.MatchString(order.ID) || strings.TrimSpace(order.Title) == "" {
		return task, decision, errors.New("invalid plane work order")
	}
	if !sha40.MatchString(order.BaselineSHA) {
		return task, decision, errors.New("invalid plane baseline sha")
	}
	if !order.ProtectAcceptedRefs || !order.RequireCleanBaseline {
		return task, decision, errors.New("plane protection gates not satisfied")
	}
	if workspace.WorkOrderID != order.ID || auth.WorkOrderID != order.ID {
		return task, decision, errors.New("Stage-1 work-order identity mismatch")
	}
	if workspace.Attempt < 1 || auth.Attempt != workspace.Attempt {
		return task, decision, errors.New("Stage-1 attempt mismatch")
	}
	if workspace.BaselineSHA != order.BaselineSHA || auth.BaselineSHA != order.BaselineSHA {
		return task, decision, errors.New("Stage-1 baseline mismatch")
	}
	if strings.TrimSpace(workspace.Path) == "" {
		return task, decision, errors.New("invalid plane workspace binding")
	}
	if len(order.AllowedPaths) != 1 {
		return task, decision, errors.New("Stage-1 requires exactly one allowed production path")
	}
	if !validStage1ProductionPath(auth.ProductionPath) {
		return task, decision, errors.New("invalid Stage-1 production path")
	}
	if order.AllowedPaths[0] != auth.ProductionPath {
		return task, decision, errors.New("Stage-1 production path not exactly bound to work order")
	}
	if containsExact(order.ForbiddenPaths, auth.ProductionPath) {
		return task, decision, errors.New("Stage-1 production path is forbidden")
	}
	if !sha256Lower.MatchString(auth.BeforeSHA256) {
		return task, decision, errors.New("invalid Stage-1 exact preimage sha256")
	}
	if auth.TaskClass != broker.Stage1TaskClassBoundedSingleFileRepair {
		return task, decision, errors.New("Stage-1 task class not allowed")
	}
	if auth.CandidateFormat != Stage1CandidateFormatWriteTextJSONSchema {
		return task, decision, errors.New("Stage-1 candidate format not allowed")
	}
	if !auth.ExactPreimage {
		return task, decision, errors.New("Stage-1 exact preimage required")
	}
	if auth.TestsMutable {
		return task, decision, errors.New("Stage-1 tests must be immutable")
	}
	if auth.ApplicationAuthority != Stage1ApplicationAuthorityExternal {
		return task, decision, errors.New("Stage-1 external application required")
	}
	if auth.AcceptanceAuthority != Stage1AcceptanceAuthorityExternalRequired {
		return task, decision, errors.New("Stage-1 external acceptance required")
	}
	if auth.ModelRepairRetries != 0 {
		return task, decision, errors.New("Stage-1 model repair retries not allowed")
	}
	if auth.ChangesRoutingArchitecture {
		return task, decision, errors.New("Stage-1 routing or architecture change not allowed")
	}
	if auth.ChangesAuthority {
		return task, decision, errors.New("Stage-1 authority change not allowed")
	}
	if auth.LiveAction {
		return task, decision, errors.New("Stage-1 live action not allowed")
	}
	if auth.ArbitraryShell {
		return task, decision, errors.New("Stage-1 arbitrary shell not allowed")
	}
	if auth.ModelEditAuthority {
		return task, decision, errors.New("Stage-1 model edit authority not allowed")
	}
	if auth.ModelTestAuthority {
		return task, decision, errors.New("Stage-1 model test authority not allowed")
	}

	task = broker.Stage1LocalTask{
		TaskClass:                    auth.TaskClass,
		ProductionFiles:              1,
		CandidateOnly:                true,
		StructuredCandidate:          true,
		ExactPreimage:                true,
		TestsMutable:                 false,
		ExternalApplication:          true,
		ExternalAcceptance:           true,
		ModelRepairRetries:           0,
		ChangesRoutingOrArchitecture: false,
		ChangesAuthority:             false,
		LiveAction:                   false,
		ArbitraryShell:               false,
		ModelAppliesEdits:            false,
		ModelControlsTests:           false,
	}

	decision, err := broker.SelectStage1LocalModel(task)
	if err != nil {
		return broker.Stage1LocalTask{}, decision, fmt.Errorf("Stage-1 broker selection rejected projected task: %w", err)
	}
	return task, decision, nil
}
