// Package ckbplane contains the inactive, pure translation boundary between the
// frozen ckb-plane contract and Wingless intelligence. It has no transport,
// queue, process, acceptance, or live-plane authority.
package ckbplane

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

const FrozenPlaneSourceCommit = "539d56fec273c8851aea131cf4d31425a62f7250"

var (
	safeID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,119}$`)
	sha40  = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

// WorkOrderContract is the subset of the frozen plane WorkOrder that Wingless
// is allowed to consume. It deliberately contains no queue or acceptance state.
type WorkOrderContract struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Instructions        []string `json:"instructions,omitempty"`
	BaselineSHA         string   `json:"baseline_sha"`
	WorkerBranch        string   `json:"worker_branch"`
	AllowedPaths        []string `json:"allowed_paths"`
	ForbiddenPaths      []string `json:"forbidden_paths,omitempty"`
	ProtectAcceptedRefs bool     `json:"protect_accepted_refs"`
	RequireCleanBaseline bool    `json:"require_clean_baseline"`
}

// WorkspaceBinding is plane-owned identity. Wingless may consume it but may not
// create, move, retry, clean up, or otherwise mutate the workspace lifecycle.
type WorkspaceBinding struct {
	WorkOrderID string `json:"work_order_id"`
	Attempt     int    `json:"attempt"`
	Path        string `json:"path"`
	BaselineSHA string `json:"baseline_sha"`
}

// RequestPolicy is trusted adapter configuration, not model output and not a
// claim that ckb-plane reserves resources or authorizes model routing.
type RequestPolicy struct {
	Role             string           `json:"role"`
	MaxContextBytes  int              `json:"max_context_bytes"`
	MaxOutputTokens  int              `json:"max_output_tokens"`
	Deadline         time.Time        `json:"deadline"`
	Capabilities     []string         `json:"capabilities"`
	AllowedBackends  []string         `json:"allowed_backends"`
	ExplicitBackend  string           `json:"explicit_backend,omitempty"`
	AllowDeep        bool             `json:"allow_deep"`
	AllowFallback    bool             `json:"allow_fallback"`
	Reason           string           `json:"reason,omitempty"`
	Resources        resources.Policy `json:"resources"`
}

// ReplayEnvelope is a test/replay normalization only. It is not a live plane
// wire format. DecodeReplayEnvelope rejects unknown fields to fail closed.
type ReplayEnvelope struct {
	Order     WorkOrderContract `json:"order"`
	Workspace WorkspaceBinding  `json:"workspace"`
	Context   string            `json:"context"`
	Policy    RequestPolicy     `json:"policy"`
}

func DecodeReplayEnvelope(data []byte) (ReplayEnvelope, error) {
	var e ReplayEnvelope
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&e); err != nil {
		return e, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		if err == nil {
			return e, errors.New("exactly one replay envelope required")
		}
		return e, err
	}
	return e, nil
}

func RequestID(workOrderID string, attempt int) (string, error) {
	if !safeID.MatchString(workOrderID) || attempt < 1 {
		return "", errors.New("invalid plane work identity")
	}
	return fmt.Sprintf("%s-wingless-%d", workOrderID, attempt), nil
}

// Translate validates the plane-owned work/workspace identity and produces one
// bounded inference request. Wingless repair attempts are always disabled for a
// plane-driven request because ckb-plane owns work-order retry/reconciliation.
func Translate(e ReplayEnvelope) (inference.Request, broker.Policy, error) {
	var req inference.Request
	var route broker.Policy
	if !safeID.MatchString(e.Order.ID) || strings.TrimSpace(e.Order.Title) == "" {
		return req, route, errors.New("invalid plane work order")
	}
	if !sha40.MatchString(e.Order.BaselineSHA) {
		return req, route, errors.New("invalid plane baseline sha")
	}
	if !e.Order.ProtectAcceptedRefs || !e.Order.RequireCleanBaseline {
		return req, route, errors.New("plane protection gates not satisfied")
	}
	if e.Workspace.WorkOrderID != e.Order.ID {
		return req, route, errors.New("workspace work-order mismatch")
	}
	if e.Workspace.BaselineSHA != e.Order.BaselineSHA {
		return req, route, errors.New("stale workspace revision")
	}
	if e.Workspace.Attempt < 1 || strings.TrimSpace(e.Workspace.Path) == "" {
		return req, route, errors.New("invalid plane workspace binding")
	}
	if len(e.Policy.AllowedBackends) == 0 {
		return req, route, errors.New("no trusted Wingless backend policy")
	}
	id, err := RequestID(e.Order.ID, e.Workspace.Attempt)
	if err != nil {
		return req, route, err
	}
	req = inference.Request{
		ID:              id,
		ParentWorkID:    e.Order.ID,
		Role:            e.Policy.Role,
		Context:         e.Context,
		MaxContextBytes: e.Policy.MaxContextBytes,
		MaxOutputTokens: e.Policy.MaxOutputTokens,
		Deadline:        e.Policy.Deadline,
		Workspace:       e.Workspace.Path,
		Capabilities:    append([]string(nil), e.Policy.Capabilities...),
		Resources:       e.Policy.Resources,
	}
	if err := req.Validate(); err != nil {
		return inference.Request{}, broker.Policy{}, err
	}
	route = broker.Policy{
		Explicit:        e.Policy.ExplicitBackend,
		AllowedBackends: append([]string(nil), e.Policy.AllowedBackends...),
		AllowDeep:       e.Policy.AllowDeep,
		AllowFallback:   e.Policy.AllowFallback,
		Reason:          e.Policy.Reason,
		Repairs:         0,
		MaxRepairs:      0,
	}
	return req, route, nil
}

type Canceler interface {
	Cancel(string) bool
}

// Cancel maps a plane work-order attempt to the deterministic Wingless request
// ID. The caller still owns the plane lifecycle and decides the terminal state.
func Cancel(c Canceler, workOrderID string, attempt int) error {
	if c == nil {
		return errors.New("nil cancel target")
	}
	id, err := RequestID(workOrderID, attempt)
	if err != nil {
		return err
	}
	if !c.Cancel(id) {
		return errors.New("Wingless request not cancelable")
	}
	return nil
}

// Candidate is non-authoritative output for the plane to inspect. There is no
// accepted=true field by design.
type Candidate struct {
	WorkOrderID string `json:"work_order_id"`
	RequestID   string `json:"request_id"`
	Status      string `json:"status"`
	Acceptance  string `json:"acceptance"`
	Text        string `json:"text,omitempty"`
	BackendID   string `json:"backend_id,omitempty"`
	ModelID     string `json:"model_id,omitempty"`
	ErrorClass  string `json:"error_class,omitempty"`
}

func CandidateFromOutcome(workOrderID, requestID string, out broker.Outcome) (Candidate, error) {
	c := Candidate{WorkOrderID: workOrderID, RequestID: requestID, Status: out.Status, Acceptance: "external_required"}
	if !safeID.MatchString(workOrderID) || requestID == "" {
		return Candidate{}, errors.New("invalid candidate identity")
	}
	if out.Acceptance != "external_required" {
		return Candidate{}, errors.New("Wingless outcome attempted to claim acceptance")
	}
	switch out.Status {
	case "result_ready", "verified_candidate", "operator_blocked":
	default:
		return Candidate{}, fmt.Errorf("unknown Wingless terminal status %q", out.Status)
	}
	if out.Status == "operator_blocked" {
		return c, nil
	}
	if len(out.Events) == 0 {
		return Candidate{}, errors.New("candidate result missing evidence event")
	}
	last := out.Events[len(out.Events)-1].Result
	c.Text = last.Text
	c.BackendID = last.BackendID
	c.ModelID = last.ModelID
	c.ErrorClass = last.ErrorClass
	return c, nil
}
