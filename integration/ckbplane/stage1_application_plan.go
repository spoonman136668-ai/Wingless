package ckbplane

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const maxStage1ReplacementBytes = 262144

type Stage1ExternalApplicationPlan struct {
	WorkOrderID          string `json:"work_order_id"`
	RequestID            string `json:"request_id"`
	Workspace            string `json:"workspace"`
	Path                 string `json:"path"`
	ExpectedBeforeSHA256 string `json:"expected_before_sha256"`
	ExpectedAfterSHA256  string `json:"expected_after_sha256"`
	Content              string `json:"content"`
	ApplicationAuthority string `json:"application_authority"`
	AcceptanceAuthority  string `json:"acceptance_authority"`
	TrustCeiling         string `json:"trust_ceiling"`
}

// PlanStage1ExternalApplication converts an already validated Stage-1
// candidate into immutable external-application data.
//
// It deliberately contains no filesystem write, rename, chmod, shell,
// subprocess, test, acceptance, promotion, queue, or live-runtime behavior.
func PlanStage1ExternalApplication(
	req Stage1CandidateRequest,
	candidate Stage1ValidatedCandidate,
	currentSource string,
) (Stage1ExternalApplicationPlan, error) {
	var out Stage1ExternalApplicationPlan

	if err := validateStage1BuiltRequestBinding(req); err != nil {
		return out, err
	}
	if !utf8.ValidString(currentSource) || strings.IndexByte(currentSource, 0) >= 0 {
		return out, errors.New("Stage-1 current source must be non-NUL UTF-8 text")
	}
	if observed := exactSHA256(currentSource); observed != req.BeforeSHA256 {
		return out, fmt.Errorf("Stage-1 current source preimage mismatch: got %s want %s", observed, req.BeforeSHA256)
	}

	op := candidate.Operation
	if op.Type != "write_text" {
		return out, errors.New("Stage-1 application plan requires write_text")
	}
	if op.Path != req.ProductionPath {
		return out, errors.New("Stage-1 application plan path mismatch")
	}
	if op.ExpectedBeforeSHA256 != req.BeforeSHA256 {
		return out, errors.New("Stage-1 application plan preimage mismatch")
	}
	if op.Content == "" || !utf8.ValidString(op.Content) || strings.IndexByte(op.Content, 0) >= 0 {
		return out, errors.New("Stage-1 application content must be non-empty non-NUL UTF-8 text")
	}
	if len(op.Content) > maxStage1ReplacementBytes {
		return out, errors.New("Stage-1 application content exceeds bounded replacement size")
	}

	afterSHA := exactSHA256(op.Content)
	if afterSHA == req.BeforeSHA256 {
		return out, errors.New("Stage-1 application plan is a no-op")
	}
	if candidate.ContentSHA256 == "" || candidate.ContentSHA256 != afterSHA {
		return out, errors.New("Stage-1 validated candidate postimage hash mismatch")
	}

	if req.Request.ParentWorkID == "" || req.Request.ID == "" || strings.TrimSpace(req.Request.Workspace) == "" {
		return out, errors.New("Stage-1 request identity incomplete")
	}

	out = Stage1ExternalApplicationPlan{
		WorkOrderID:          req.Request.ParentWorkID,
		RequestID:            req.Request.ID,
		Workspace:            req.Request.Workspace,
		Path:                 req.ProductionPath,
		ExpectedBeforeSHA256: req.BeforeSHA256,
		ExpectedAfterSHA256:  afterSHA,
		Content:              op.Content,
		ApplicationAuthority: Stage1ApplicationAuthorityExternal,
		AcceptanceAuthority:  Stage1AcceptanceAuthorityExternalRequired,
		TrustCeiling:         "stage_1_only",
	}
	return out, nil
}
