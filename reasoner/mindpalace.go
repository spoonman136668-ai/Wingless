package reasoner

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	MindPalaceContextSchema = "ckb-plane.mind-palace-context.v1"
	PublicEvidenceIndexSchema = "wingless.public-evidence-index.v1"
)

type MindPalaceContext struct {
	Schema         string   `json:"schema"`
	Status         string   `json:"status"`
	MindPalaceHead string   `json:"mind_palace_head"`
	RecordsSHA256  string   `json:"records_sha256"`
	ModelContext   string   `json:"model_context"`
	GovernanceIDs  []string `json:"governance_ids"`
	PrimaryIDs     []string `json:"primary_ids"`
	SupportIDs     []string `json:"support_ids"`
	AdvisoryOnly   bool     `json:"advisory_only"`
	AuthorizesExecution bool `json:"authorizes_execution"`
	AuthorizesAcceptance bool `json:"authorizes_acceptance"`
	DurableContextWritten bool `json:"durable_context_written"`
	GitWritten bool `json:"git_written"`
}

func ParseMindPalaceContext(raw string) (MindPalaceContext, error) {
	var out MindPalaceContext
	if strings.TrimSpace(raw) == "" {
		return out, errors.New("mind-palace context empty")
	}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return out, err
	}
	if out.Schema != MindPalaceContextSchema || out.Status != "succeeded" {
		return out, errors.New("mind-palace context schema or status invalid")
	}
	if len(out.MindPalaceHead) != 40 || !isSHA256(out.RecordsSHA256) {
		return out, errors.New("mind-palace provenance invalid")
	}
	if !out.AdvisoryOnly || out.AuthorizesExecution || out.AuthorizesAcceptance ||
		out.DurableContextWritten || out.GitWritten {
		return out, errors.New("mind-palace authority boundary invalid")
	}
	return out, nil
}

type PublicEvidenceIndexEntry struct {
	RecordID     string `json:"record_id"`
	ExperimentID string `json:"experiment_id"`
	Project      string `json:"project"`
}

type PublicEvidenceIndex struct {
	Schema  string                     `json:"schema"`
	Entries []PublicEvidenceIndexEntry `json:"entries"`
}

func (i PublicEvidenceIndex) Validate() error {
	if i.Schema != PublicEvidenceIndexSchema {
		return errors.New("public evidence index schema mismatch")
	}
	if len(i.Entries) == 0 || len(i.Entries) > 4096 {
		return errors.New("public evidence index size invalid")
	}
	seen := map[string]bool{}
	for _, e := range i.Entries {
		if strings.TrimSpace(e.RecordID) == "" || strings.TrimSpace(e.ExperimentID) == "" {
			return errors.New("public evidence index identity missing")
		}
		if e.Project != "Wingless" && e.Project != "Yggdrasil" {
			return errors.New("public evidence index project invalid")
		}
		key := e.RecordID + "\x00" + e.ExperimentID
		if seen[key] {
			return errors.New("public evidence index duplicate")
		}
		seen[key] = true
	}
	return nil
}

// ResolveMindPalaceContext converts only durable record identifiers into public
// experiment IDs. ModelContext is intentionally ignored so private prose never
// becomes free-route prompt material.
func ResolveMindPalaceContext(mp MindPalaceContext, index PublicEvidenceIndex, query string, maxExperiments int) (ContextSelection, error) {
	if mp.Schema != MindPalaceContextSchema || mp.Status != "succeeded" {
		return ContextSelection{}, errors.New("mind-palace context invalid")
	}
	if err := index.Validate(); err != nil {
		return ContextSelection{}, err
	}
	if maxExperiments <= 0 || maxExperiments > MaxSelectedDigests {
		return ContextSelection{}, errors.New("mind-palace selection bound invalid")
	}

	byRecord := map[string][]PublicEvidenceIndexEntry{}
	for _, e := range index.Entries {
		byRecord[e.RecordID] = append(byRecord[e.RecordID], e)
	}
	records := append([]string(nil), mp.PrimaryIDs...)
	records = append(records, mp.SupportIDs...)

	seenExperiment := map[string]bool{}
	var experiments []string
	for _, recordID := range records {
		entries := append([]PublicEvidenceIndexEntry(nil), byRecord[recordID]...)
		sort.Slice(entries, func(a, b int) bool {
			if entries[a].Project == entries[b].Project {
				return entries[a].ExperimentID < entries[b].ExperimentID
			}
			return entries[a].Project < entries[b].Project
		})
		for _, e := range entries {
			if seenExperiment[e.ExperimentID] {
				continue
			}
			seenExperiment[e.ExperimentID] = true
			experiments = append(experiments, e.ExperimentID)
			if len(experiments) == maxExperiments {
				break
			}
		}
		if len(experiments) == maxExperiments {
			break
		}
	}
	if len(experiments) == 0 {
		return ContextSelection{}, fmt.Errorf("mind-palace returned no resolvable public experiment IDs")
	}
	return ContextSelection{
		Schema:        ContextSelectionSchema,
		Provider:      "mind-palace",
		Query:         query,
		ExperimentIDs: experiments,
	}, nil
}
