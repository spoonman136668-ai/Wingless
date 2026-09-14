package ckbplane

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/spoonman136668-ai/Wingless/broker"
)

const Stage1PostApplicationStatusVerified = "postimage_verified"

type Stage1PostApplicationEvidence struct {
	WorkOrderID          string `json:"work_order_id"`
	RequestID            string `json:"request_id"`
	Workspace            string `json:"workspace"`
	Path                 string `json:"path"`
	ExpectedBeforeSHA256 string `json:"expected_before_sha256"`
	ExpectedAfterSHA256  string `json:"expected_after_sha256"`
	ObservedAfterSHA256  string `json:"observed_after_sha256"`
	Status               string `json:"status"`
	Acceptance           string `json:"acceptance"`
	TrustCeiling         string `json:"trust_ceiling"`
}

func validateStage1ExternalApplicationPlanBinding(
	req Stage1CandidateRequest,
	plan Stage1ExternalApplicationPlan,
) error {
	if err := validateStage1BuiltRequestBinding(req); err != nil {
		return err
	}
	if plan.WorkOrderID == "" || plan.WorkOrderID != req.Request.ParentWorkID {
		return errors.New("Stage-1 application plan work-order identity mismatch")
	}
	if plan.RequestID == "" || plan.RequestID != req.Request.ID {
		return errors.New("Stage-1 application plan request identity mismatch")
	}
	if strings.TrimSpace(plan.Workspace) == "" || plan.Workspace != req.Request.Workspace {
		return errors.New("Stage-1 application plan workspace mismatch")
	}
	if !validStage1ProductionPath(plan.Path) || plan.Path != req.ProductionPath {
		return errors.New("Stage-1 application plan path mismatch")
	}
	if !sha256Lower.MatchString(plan.ExpectedBeforeSHA256) ||
		plan.ExpectedBeforeSHA256 != req.BeforeSHA256 {
		return errors.New("Stage-1 application plan preimage mismatch")
	}
	if !sha256Lower.MatchString(plan.ExpectedAfterSHA256) ||
		plan.ExpectedAfterSHA256 == plan.ExpectedBeforeSHA256 {
		return errors.New("invalid Stage-1 application plan postimage")
	}
	if plan.Content == "" ||
		!utf8.ValidString(plan.Content) ||
		strings.IndexByte(plan.Content, 0) >= 0 ||
		len(plan.Content) > maxStage1ReplacementBytes {
		return errors.New("invalid Stage-1 application plan content")
	}
	if observed := exactSHA256(plan.Content); observed != plan.ExpectedAfterSHA256 {
		return errors.New("Stage-1 application plan content/postimage mismatch")
	}
	if plan.ApplicationAuthority != Stage1ApplicationAuthorityExternal {
		return errors.New("Stage-1 application authority must remain external")
	}
	if plan.AcceptanceAuthority != Stage1AcceptanceAuthorityExternalRequired {
		return errors.New("Stage-1 acceptance authority must remain external_required")
	}
	if plan.TrustCeiling != broker.Stage1TrustCeiling {
		return errors.New("Stage-1 application plan trust ceiling mismatch")
	}
	return nil
}

// VerifyStage1ExternalApplication verifies bytes observed after an external
// application step. It proves only that the observed postimage matches the
// already-bound plan.
//
// It does not write files, invoke a backend, run tests, grant acceptance,
// promote refs, mutate queue state, or activate live Wingless.
func VerifyStage1ExternalApplication(
	req Stage1CandidateRequest,
	plan Stage1ExternalApplicationPlan,
	observedSource string,
) (Stage1PostApplicationEvidence, error) {
	var out Stage1PostApplicationEvidence

	if err := validateStage1ExternalApplicationPlanBinding(req, plan); err != nil {
		return out, err
	}
	if !utf8.ValidString(observedSource) || strings.IndexByte(observedSource, 0) >= 0 {
		return out, errors.New("Stage-1 observed source must be non-NUL UTF-8 text")
	}

	observedSHA := exactSHA256(observedSource)
	if observedSHA != plan.ExpectedAfterSHA256 {
		return out, fmt.Errorf(
			"Stage-1 observed postimage mismatch: got %s want %s",
			observedSHA,
			plan.ExpectedAfterSHA256,
		)
	}
	if observedSource != plan.Content {
		return out, errors.New("Stage-1 observed source bytes differ from planned content")
	}

	out = Stage1PostApplicationEvidence{
		WorkOrderID:          plan.WorkOrderID,
		RequestID:            plan.RequestID,
		Workspace:            plan.Workspace,
		Path:                 plan.Path,
		ExpectedBeforeSHA256: plan.ExpectedBeforeSHA256,
		ExpectedAfterSHA256:  plan.ExpectedAfterSHA256,
		ObservedAfterSHA256:  observedSHA,
		Status:               Stage1PostApplicationStatusVerified,
		Acceptance:           Stage1AcceptanceAuthorityExternalRequired,
		TrustCeiling:         broker.Stage1TrustCeiling,
	}
	return out, nil
}
