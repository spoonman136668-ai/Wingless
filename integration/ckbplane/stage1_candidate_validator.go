package ckbplane

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
)

const maxStage1CandidateJSONBytes = 262144

type Stage1WriteTextCandidate struct {
	Type                 string `json:"type"`
	Path                 string `json:"path"`
	ExpectedBeforeSHA256 string `json:"expected_before_sha256"`
	Content              string `json:"content"`
}

type Stage1ValidatedCandidate struct {
	Operation     Stage1WriteTextCandidate
	ContentSHA256 string
}

func decodeStage1WriteTextCandidate(raw string) (Stage1WriteTextCandidate, error) {
	var op Stage1WriteTextCandidate

	if raw == "" {
		return op, errors.New("empty Stage-1 candidate")
	}
	if len(raw) > maxStage1CandidateJSONBytes {
		return op, errors.New("Stage-1 candidate exceeds bounded JSON size")
	}
	if !utf8.ValidString(raw) || strings.IndexByte(raw, 0) >= 0 {
		return op, errors.New("Stage-1 candidate must be non-NUL UTF-8 JSON")
	}

	d := json.NewDecoder(strings.NewReader(raw))
	d.DisallowUnknownFields()
	if err := d.Decode(&op); err != nil {
		return op, fmt.Errorf("invalid Stage-1 candidate JSON: %w", err)
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		if err == nil {
			return op, errors.New("exactly one Stage-1 candidate JSON object required")
		}
		return op, fmt.Errorf("invalid trailing Stage-1 candidate data: %w", err)
	}
	return op, nil
}

func validateStage1BuiltRequestBinding(req Stage1CandidateRequest) error {
	if !req.Selection.Allowed ||
		req.Selection.TrustCeiling != broker.Stage1TrustCeiling ||
		req.ProductionPath == "" ||
		!validStage1ProductionPath(req.ProductionPath) ||
		!sha256Lower.MatchString(req.BeforeSHA256) {
		return errors.New("invalid Stage-1 built request binding")
	}
	if req.Request.OutputConstraint == nil || req.Request.OutputConstraint.Kind != "json_schema" {
		return errors.New("Stage-1 request missing json_schema output constraint")
	}
	if err := req.Request.Validate(); err != nil {
		return fmt.Errorf("invalid Stage-1 inference request: %w", err)
	}

	expectedSchema, err := stage1OutputSchema(req.ProductionPath, req.BeforeSHA256)
	if err != nil {
		return err
	}
	if !bytes.Equal(req.Request.OutputConstraint.Schema, expectedSchema) {
		return errors.New("Stage-1 request schema no longer matches path/preimage binding")
	}
	return nil
}

// ValidateStage1Candidate strictly validates raw model output against the
// already-authorized Stage-1 candidate request.
//
// It returns candidate data only. It does not write files, apply edits, run
// tests, grant acceptance, invoke a backend, or mutate plane state.
func ValidateStage1Candidate(
	req Stage1CandidateRequest,
	currentSource string,
	raw string,
) (Stage1ValidatedCandidate, error) {
	var out Stage1ValidatedCandidate

	if err := validateStage1BuiltRequestBinding(req); err != nil {
		return out, err
	}
	if !utf8.ValidString(currentSource) || strings.IndexByte(currentSource, 0) >= 0 {
		return out, errors.New("Stage-1 current source must be non-NUL UTF-8 text")
	}
	if observed := exactSHA256(currentSource); observed != req.BeforeSHA256 {
		return out, fmt.Errorf("Stage-1 current source preimage mismatch: got %s want %s", observed, req.BeforeSHA256)
	}

	op, err := decodeStage1WriteTextCandidate(raw)
	if err != nil {
		return out, err
	}
	if op.Type != "write_text" {
		return out, errors.New("Stage-1 candidate operation must be write_text")
	}
	if op.Path != req.ProductionPath {
		return out, errors.New("Stage-1 candidate path mismatch")
	}
	if op.ExpectedBeforeSHA256 != req.BeforeSHA256 {
		return out, errors.New("Stage-1 candidate preimage mismatch")
	}
	if op.Content == "" || !utf8.ValidString(op.Content) || strings.IndexByte(op.Content, 0) >= 0 {
		return out, errors.New("Stage-1 candidate content must be non-empty non-NUL UTF-8 text")
	}

	contentSHA := exactSHA256(op.Content)
	if contentSHA == req.BeforeSHA256 {
		return out, errors.New("Stage-1 candidate is a no-op replacement")
	}

	out = Stage1ValidatedCandidate{
		Operation:     op,
		ContentSHA256: contentSHA,
	}
	return out, nil
}

// Compile-time assertion that the request constraint remains the expected type.
var _ = (*inference.OutputConstraint)(nil)
