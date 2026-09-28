package reasoner

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	ExperimentDigestSchema = "wingless.experiment-digest.v1"
	ContextSelectionSchema = "wingless.context-selection.v1"
	ContextPackSchema       = "wingless.reasoner-context-pack.v1"

	MaxDigestBytes      = 32768
	MaxContextPackBytes = 65536
	MaxSelectedDigests  = 8
)

var publicResearchRepos = map[string]bool{
	"spoonman136668-ai/Wingless":  true,
	"spoonman136668-ai/Yggdrasil": true,
}

type ExperimentDigest struct {
	Schema            string   `json:"schema"`
	Project           string   `json:"project"`
	ExperimentID      string   `json:"experiment_id"`
	SourceRepository  string   `json:"source_repository"`
	SourceRef         string   `json:"source_ref"`
	SourcePath        string   `json:"source_path"`
	SourceSHA256      string   `json:"source_sha256"`
	Question          string   `json:"question"`
	Observed          []string `json:"observed"`
	Established       []string `json:"established"`
	NotEstablished    []string `json:"not_established"`
	FrozenConstraints []string `json:"frozen_constraints"`
	NextSeams         []string `json:"next_seams"`
	EvidenceRefs      []string `json:"evidence_refs"`
}

func (d ExperimentDigest) Validate() error {
	if d.Schema != ExperimentDigestSchema {
		return errors.New("experiment digest schema mismatch")
	}
	if d.Project != "Wingless" && d.Project != "Yggdrasil" {
		return errors.New("experiment digest project invalid")
	}
	if strings.TrimSpace(d.ExperimentID) == "" || strings.TrimSpace(d.Question) == "" {
		return errors.New("experiment digest identity or question missing")
	}
	if !publicResearchRepos[d.SourceRepository] {
		return errors.New("experiment digest source repository is not approved public research")
	}
	if !isGitSHA(d.SourceRef) || !isSHA256(d.SourceSHA256) {
		return errors.New("experiment digest source identity invalid")
	}
	if strings.TrimSpace(d.SourcePath) == "" || strings.Contains(d.SourcePath, "\\") ||
		strings.HasPrefix(d.SourcePath, "/") || strings.Contains(d.SourcePath, "..") {
		return errors.New("experiment digest source path invalid")
	}
	if len(d.Observed) == 0 || len(d.FrozenConstraints) == 0 || len(d.EvidenceRefs) == 0 {
		return errors.New("experiment digest missing evidence or constraints")
	}
	if err := validateDigestStrings(d); err != nil {
		return err
	}
	raw, err := json.Marshal(d)
	if err != nil {
		return err
	}
	if len(raw) > MaxDigestBytes {
		return errors.New("experiment digest exceeds compact record budget")
	}
	return nil
}

func validateDigestStrings(d ExperimentDigest) error {
	groups := [][]string{
		d.Observed, d.Established, d.NotEstablished,
		d.FrozenConstraints, d.NextSeams, d.EvidenceRefs,
	}
	count := 0
	for _, group := range groups {
		count += len(group)
		for _, s := range group {
			if strings.TrimSpace(s) == "" || !utf8.ValidString(s) || len(s) > 2048 {
				return errors.New("experiment digest contains invalid or unbounded string")
			}
		}
	}
	if count > 96 || !utf8.ValidString(d.Question) || len(d.Question) > 4096 {
		return errors.New("experiment digest content is unbounded")
	}
	return nil
}

func isGitSHA(s string) bool {
	if len(s) != 40 {
		return false
	}
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}

// ContextSelection is advisory retrieval metadata. Mind-Palace may select IDs,
// but no Mind-Palace model_context text is copied into free-route prompts.
type ContextSelection struct {
	Schema        string   `json:"schema"`
	Provider      string   `json:"provider"`
	Query         string   `json:"query"`
	ExperimentIDs []string `json:"experiment_ids"`
}

func (s ContextSelection) Validate() error {
	if s.Schema != ContextSelectionSchema {
		return errors.New("context selection schema mismatch")
	}
	if s.Provider != "mind-palace" && s.Provider != "explicit-public-fallback" {
		return errors.New("context selection provider invalid")
	}
	if strings.TrimSpace(s.Query) == "" || len(s.Query) > 256 || !utf8.ValidString(s.Query) {
		return errors.New("context selection query invalid")
	}
	if len(s.ExperimentIDs) == 0 || len(s.ExperimentIDs) > MaxSelectedDigests {
		return errors.New("context selection count invalid")
	}
	seen := map[string]bool{}
	for _, id := range s.ExperimentIDs {
		if strings.TrimSpace(id) == "" || len(id) > 160 || seen[id] {
			return errors.New("context selection IDs invalid")
		}
		seen[id] = true
	}
	return nil
}

type ContextPack struct {
	Schema               string             `json:"schema"`
	DataClass            string             `json:"data_class"`
	Query                string             `json:"query"`
	SelectionProvider    string             `json:"selection_provider"`
	ExperimentIDs        []string           `json:"experiment_ids"`
	Digests              []ExperimentDigest `json:"digests"`
	SourceIdentitySHA256 string             `json:"source_identity_sha256"`
	PackSHA256           string             `json:"pack_sha256"`
	EncodedBytes         int                `json:"encoded_bytes"`
}

func BuildContextPack(selection ContextSelection, available []ExperimentDigest, maxBytes int) (ContextPack, error) {
	var out ContextPack
	if err := selection.Validate(); err != nil {
		return out, err
	}
	if maxBytes <= 0 || maxBytes > MaxContextPackBytes {
		return out, errors.New("context pack budget invalid")
	}

	byID := make(map[string]ExperimentDigest, len(available))
	for _, d := range available {
		if err := d.Validate(); err != nil {
			return out, fmt.Errorf("%s: %w", d.ExperimentID, err)
		}
		if _, exists := byID[d.ExperimentID]; exists {
			return out, fmt.Errorf("duplicate experiment digest %s", d.ExperimentID)
		}
		byID[d.ExperimentID] = d
	}

	ids := append([]string(nil), selection.ExperimentIDs...)
	sort.Strings(ids)
	digests := make([]ExperimentDigest, 0, len(ids))
	for _, id := range ids {
		d, ok := byID[id]
		if !ok {
			return out, fmt.Errorf("selected experiment digest missing: %s", id)
		}
		digests = append(digests, d)
	}

	identity := make([]map[string]string, 0, len(digests))
	for _, d := range digests {
		identity = append(identity, map[string]string{
			"experiment_id": d.ExperimentID,
			"repository":    d.SourceRepository,
			"source_ref":    d.SourceRef,
			"source_path":   d.SourcePath,
			"source_sha256": d.SourceSHA256,
		})
	}
	idRaw, err := json.Marshal(identity)
	if err != nil {
		return out, err
	}

	out = ContextPack{
		Schema:               ContextPackSchema,
		DataClass:            "public-repository",
		Query:                selection.Query,
		SelectionProvider:    selection.Provider,
		ExperimentIDs:        ids,
		Digests:              digests,
		SourceIdentitySHA256: sha256hex(idRaw),
	}
	preimage := out
	preimage.PackSHA256 = ""
	preimage.EncodedBytes = 0
	raw, err := json.Marshal(preimage)
	if err != nil {
		return ContextPack{}, err
	}
	if len(raw) > maxBytes {
		return ContextPack{}, fmt.Errorf("context pack exceeds budget: bytes=%d max=%d; narrow selection rather than truncate", len(raw), maxBytes)
	}
	out.PackSHA256 = sha256hex(raw)

	finalRaw, err := json.Marshal(out)
	if err != nil {
		return ContextPack{}, err
	}
	if len(finalRaw) > maxBytes {
		return ContextPack{}, fmt.Errorf("context pack metadata exceeds budget: bytes=%d max=%d", len(finalRaw), maxBytes)
	}
	out.EncodedBytes = len(finalRaw)
	return out, nil
}

func (p ContextPack) ContextForReasoner() (string, error) {
	if p.Schema != ContextPackSchema || p.DataClass != "public-repository" ||
		!isSHA256(p.SourceIdentitySHA256) || !isSHA256(p.PackSHA256) {
		return "", errors.New("context pack identity invalid")
	}
	raw, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	if len(raw) > MaxContextPackBytes {
		return "", errors.New("context pack exceeds reasoner bound")
	}
	return string(raw), nil
}

type ContextSelector interface {
	Select(query string, maxExperiments int) (ContextSelection, error)
}

var ErrContextSelectorUnavailable = errors.New("context selector unavailable")

func RequireContextSelection(selector ContextSelector, query string, maxExperiments int) (ContextSelection, error) {
	if selector == nil {
		return ContextSelection{}, ErrContextSelectorUnavailable
	}
	s, err := selector.Select(query, maxExperiments)
	if err != nil {
		return ContextSelection{}, err
	}
	if err := s.Validate(); err != nil {
		return ContextSelection{}, err
	}
	if s.Provider != "mind-palace" {
		return ContextSelection{}, errors.New("live context selection requires mind-palace provider")
	}
	return s, nil
}

func BuildReasonerRequestFromPack(pack ContextPack, requestID, project, experimentID string, role Role, qualificationContractSHA256 string, maxOutputTokens, timeoutSeconds int) (Request, error) {
	ctx, err := pack.ContextForReasoner()
	if err != nil {
		return Request{}, err
	}
	return Request{
		Schema:                      RequestSchema,
		RequestID:                   requestID,
		Project:                     project,
		ExperimentID:                experimentID,
		Role:                        role,
		FrontierSHA256:              pack.PackSHA256,
		QualificationContractSHA256: qualificationContractSHA256,
		Context:                     ctx,
		MaxOutputTokens:             maxOutputTokens,
		TimeoutSeconds:              timeoutSeconds,
		DataClass:                   "public-repository",
	}, nil
}
