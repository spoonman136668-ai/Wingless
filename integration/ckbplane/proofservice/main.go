package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ckbplane "github.com/spoonman136668-ai/Wingless/integration/ckbplane"
)

type proofEngine struct {
	marker string
}

func (e proofEngine) Run(ctx context.Context, req ckbplane.TransportRequest) (ckbplane.CandidateEnvelope, error) {
	if strings.Contains(req.Title, "CANCEL_PROOF") {
		if err := os.WriteFile(e.marker, []byte("active\n"), 0600); err != nil {
			return ckbplane.CandidateEnvelope{}, err
		}
		defer os.Remove(e.marker)
		<-ctx.Done()
		return ckbplane.CandidateEnvelope{}, ctx.Err()
	}
	backend := "fast"
	if len(req.Policy.AllowedBackends) > 0 {
		backend = req.Policy.AllowedBackends[0]
	}
	edits := &ckbplane.EditEnvelope{
		SourceCommit:    req.SourceCommit,
		RequestID:       req.RequestID,
		WorkOrderID:     req.WorkOrderID,
		Attempt:         req.Attempt,
		Workspace:       req.Workspace,
		BaselineSHA:     req.BaselineSHA,
		WorkOrderDigest: req.WorkOrderDigest,
		Operations: []ckbplane.EditOperation{{
			Type:    "write_text",
			Path:    "src/serviceproof.go",
			Content: "package src\n\nconst WinglessServiceProof = true\n",
		}},
	}
	return ckbplane.CandidateEnvelope{
		Schema:          ckbplane.CandidateEnvelopeSchema,
		SourceCommit:    req.SourceCommit,
		RequestID:       req.RequestID,
		WorkOrderID:     req.WorkOrderID,
		Attempt:         req.Attempt,
		Workspace:       req.Workspace,
		BaselineSHA:     req.BaselineSHA,
		WorkOrderDigest: req.WorkOrderDigest,
		Status:          "result_ready",
		Acceptance:      "external_required",
		Text:            "isolated process-boundary candidate",
		BackendID:       backend,
		ModelID:         "wingless-proof-service",
		Edits:           edits,
	}, nil
}

func main() {
	session := os.Getenv("CKB_WINGLESS_SESSION")
	if session == "" {
		fmt.Fprintln(os.Stderr, "WINGLESS_PROOF_SERVICE_SESSION_REQUIRED")
		os.Exit(2)
	}
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	marker := filepath.Join(filepath.Dir(exe), "wingless-proof-active.marker")
	_ = os.Remove(marker)
	if err := ckbplane.ServeLoopbackService(context.Background(), session, proofEngine{marker: marker}, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
