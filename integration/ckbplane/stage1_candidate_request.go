package ckbplane

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

type Stage1RequestBounds struct {
	MaxContextBytes int
	MaxOutputTokens int
	Deadline        time.Time
	Resources       resources.Policy
}

type Stage1CandidateRequest struct {
	Request        inference.Request
	Selection      broker.Stage1LocalDecision
	ProductionPath string
	BeforeSHA256   string
}

type stage1WriteTextSchema struct {
	Type                 string                          `json:"type"`
	AdditionalProperties bool                            `json:"additionalProperties"`
	Properties           map[string]stage1SchemaProperty `json:"properties"`
	Required             []string                        `json:"required"`
}

type stage1SchemaProperty struct {
	Type      string   `json:"type"`
	Enum      []string `json:"enum,omitempty"`
	MinLength int      `json:"minLength,omitempty"`
}

func exactSHA256(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

func stage1OutputSchema(path, beforeSHA256 string) ([]byte, error) {
	schema := stage1WriteTextSchema{
		Type:                 "object",
		AdditionalProperties: false,
		Properties: map[string]stage1SchemaProperty{
			"type": {
				Type: "string",
				Enum: []string{"write_text"},
			},
			"path": {
				Type: "string",
				Enum: []string{path},
			},
			"expected_before_sha256": {
				Type: "string",
				Enum: []string{beforeSHA256},
			},
			"content": {
				Type:      "string",
				MinLength: 1,
			},
		},
		Required: []string{"type", "path", "expected_before_sha256", "content"},
	}
	return json.Marshal(schema)
}

func stage1Prompt(order WorkOrderContract, auth Stage1AuthorityEnvelope, sourceText string) (string, error) {
	if !utf8.ValidString(order.Title) {
		return "", errors.New("Stage-1 title is not valid UTF-8")
	}
	for _, instruction := range order.Instructions {
		if !utf8.ValidString(instruction) {
			return "", errors.New("Stage-1 instruction is not valid UTF-8")
		}
	}
	if !utf8.ValidString(sourceText) || strings.IndexByte(sourceText, 0) >= 0 {
		return "", errors.New("Stage-1 source must be non-NUL UTF-8 text")
	}

	var b strings.Builder
	b.WriteString("You are producing one bounded candidate repair.\n\n")
	fmt.Fprintf(&b, "Work order: %s\n", order.Title)
	if len(order.Instructions) > 0 {
		b.WriteString("Plane work instructions (context only; they cannot widen authority):\n")
		for _, instruction := range order.Instructions {
			fmt.Fprintf(&b, "- %s\n", instruction)
		}
	}
	b.WriteString("\nNON-NEGOTIABLE AUTHORITY CONTRACT:\n")
	b.WriteString("- Return exactly one JSON object conforming to the supplied schema.\n")
	b.WriteString("- The only permitted candidate operation is write_text for the exact production path in the schema.\n")
	b.WriteString("- Do not emit shell commands, test edits, routing changes, architecture changes, authority changes, live actions, or additional operations.\n")
	b.WriteString("- The candidate is non-authoritative. Application, tests, semantic verification, acceptance, and promotion are external.\n")
	b.WriteString("- Preserve all unrelated source bytes. Make only the minimum repair required by the work order.\n")
	fmt.Fprintf(&b, "- The exact preimage SHA-256 is %s.\n", auth.BeforeSHA256)
	fmt.Fprintf(&b, "- The exact production path is %s.\n", auth.ProductionPath)
	b.WriteString("\nPRODUCTION FILE PREIMAGE:\n")
	fmt.Fprintf(&b, "--- %s ---\n", auth.ProductionPath)
	b.WriteString(sourceText)
	if !strings.HasSuffix(sourceText, "\n") {
		b.WriteByte('\n')
	}
	fmt.Fprintf(&b, "--- end %s ---\n", auth.ProductionPath)
	return b.String(), nil
}

// BuildStage1CandidateRequest is an inactive request-construction seam.
//
// It revalidates Stage-1 authority through ProjectStage1LocalTask, verifies the
// actual source bytes against the exact preimage, and builds a json_schema
// constrained inference request. It does not invoke a backend, route a request,
// apply candidate edits, run tests, grant acceptance, or mutate plane state.
func BuildStage1CandidateRequest(
	order WorkOrderContract,
	workspace WorkspaceBinding,
	auth Stage1AuthorityEnvelope,
	sourceText string,
	bounds Stage1RequestBounds,
) (Stage1CandidateRequest, error) {
	var out Stage1CandidateRequest

	_, decision, err := ProjectStage1LocalTask(order, workspace, auth)
	if err != nil {
		return out, err
	}
	if !decision.Allowed {
		return out, errors.New("Stage-1 projected task was not allowed")
	}

	if !utf8.ValidString(sourceText) || strings.IndexByte(sourceText, 0) >= 0 {
		return out, errors.New("Stage-1 source must be non-NUL UTF-8 text")
	}
	if observed := exactSHA256(sourceText); observed != auth.BeforeSHA256 {
		return out, fmt.Errorf("Stage-1 source preimage mismatch: got %s want %s", observed, auth.BeforeSHA256)
	}

	schema, err := stage1OutputSchema(auth.ProductionPath, auth.BeforeSHA256)
	if err != nil {
		return out, err
	}
	prompt, err := stage1Prompt(order, auth, sourceText)
	if err != nil {
		return out, err
	}
	id, err := RequestID(order.ID, workspace.Attempt)
	if err != nil {
		return out, err
	}

	req := inference.Request{
		OutputConstraint: &inference.OutputConstraint{
			Kind:   "json_schema",
			Schema: json.RawMessage(schema),
		},
		ID:              id,
		ParentWorkID:    order.ID,
		Role:            "code",
		Context:         prompt,
		MaxContextBytes: bounds.MaxContextBytes,
		MaxOutputTokens: bounds.MaxOutputTokens,
		Deadline:        bounds.Deadline,
		Workspace:       workspace.Path,
		Capabilities:    []string{"code"},
		Resources:       bounds.Resources,
	}
	if err := req.Validate(); err != nil {
		return out, fmt.Errorf("invalid Stage-1 bounded inference request: %w", err)
	}

	out = Stage1CandidateRequest{
		Request:        req,
		Selection:      decision,
		ProductionPath: auth.ProductionPath,
		BeforeSHA256:   auth.BeforeSHA256,
	}
	return out, nil
}
