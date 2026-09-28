package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spoonman136668-ai/Wingless/reasoner"
)

const dryRunContract = "wingless.nemotron-reduced-context-dryrun.v1|mind-palace-id-selection|public-reconstruction-only|programs-independent|no-authority"

type selectionFixture struct {
	Schema                       string   `json:"schema"`
	SourceRequestID              string   `json:"source_request_id"`
	MindPalaceHead               string   `json:"mind_palace_head"`
	RecordsSHA256                string   `json:"records_sha256"`
	Query                        string   `json:"query"`
	SelectedRecordIDs            []string `json:"selected_record_ids"`
	ExperimentIDs                []string `json:"experiment_ids"`
	PrivateModelContextIncluded  bool     `json:"private_model_context_included"`
	PublicReconstructionRequired bool     `json:"public_reconstruction_required"`
}

type answer struct {
	WinglessNextQuestion      string `json:"wingless_next_question"`
	YggdrasilNextQuestion     string `json:"yggdrasil_next_question"`
	ConstraintViolation      bool   `json:"constraint_violation"`
	CrossProjectLineageClaim bool   `json:"cross_project_lineage_claim"`
}

type evidence struct {
	Schema                       string            `json:"schema"`
	Model                        string            `json:"model"`
	SourceRequestID              string            `json:"source_request_id"`
	MindPalaceHead               string            `json:"mind_palace_head"`
	RecordsSHA256                string            `json:"records_sha256"`
	SelectedRecordIDs            []string          `json:"selected_record_ids"`
	ExperimentIDs                []string          `json:"experiment_ids"`
	PackSHA256                   string            `json:"pack_sha256"`
	SourceIdentitySHA256         string            `json:"source_identity_sha256"`
	PackBytes                    int               `json:"pack_bytes"`
	QualificationContractSHA256  string            `json:"qualification_contract_sha256"`
	RequestSHA256                string            `json:"request_sha256"`
	ConfigSHA256                 string            `json:"config_sha256"`
	ResponseSHA256               string            `json:"response_sha256"`
	QualificationIdentitySHA256  string            `json:"qualification_identity_sha256"`
	ReturnedModel                string            `json:"returned_model"`
	Provider                     string            `json:"provider,omitempty"`
	ReasoningTokens              *int              `json:"reasoning_tokens,omitempty"`
	LatencyMS                    int64             `json:"latency_ms"`
	Answer                       answer            `json:"answer"`
	PublicRepositoryOnly         bool              `json:"public_repository_only"`
	PrivateMindPalaceTextSent     bool              `json:"private_mind_palace_text_sent"`
	AdvisoryOnly                 bool              `json:"advisory_only"`
}

func hashText(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func loadFixture(path string) (selectionFixture, error) {
	var f selectionFixture
	raw, err := os.ReadFile(path)
	if err != nil {
		return f, err
	}
	if err := json.Unmarshal(raw, &f); err != nil {
		return f, err
	}
	if f.Schema != "wingless.mind-palace-public-selection-fixture.v1" ||
		f.PrivateModelContextIncluded ||
		!f.PublicReconstructionRequired ||
		len(f.MindPalaceHead) != 40 ||
		len(f.RecordsSHA256) != 64 ||
		len(f.ExperimentIDs) == 0 ||
		len(f.ExperimentIDs) > reasoner.MaxSelectedDigests {
		return f, errors.New("selection fixture invalid")
	}
	return f, nil
}

func loadDigests(dir string, ids []string) ([]reasoner.ExperimentDigest, error) {
	out := make([]reasoner.ExperimentDigest, 0, len(ids))
	for _, id := range ids {
		if strings.ContainsAny(id, "/\\") || id == "" {
			return nil, errors.New("experiment ID invalid")
		}
		raw, err := os.ReadFile(filepath.Join(dir, id+".json"))
		if err != nil {
			return nil, err
		}
		var d reasoner.ExperimentDigest
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, err
		}
		if d.ExperimentID != id {
			return nil, fmt.Errorf("digest identity mismatch for %s", id)
		}
		if err := d.Validate(); err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		out = append(out, d)
	}
	return out, nil
}

func main() {
	fixturePath := flag.String("selection", "reasoner/fixtures/context/mind-palace-live-selection-20260928.json", "frozen Mind-Palace selection fixture")
	digestDir := flag.String("digests", "reasoner/digests/history", "public digest directory")
	outPath := flag.String("out", "evidence/nemotron-reduced-context-dryrun.json", "evidence output")
	flag.Parse()

	if os.Getenv("WINGLESS_REMOTE_REASONER_ENABLE") != "1" {
		die(errors.New("remote reasoner disabled"))
	}
	key := os.Getenv("WINGLESS_REASONER_OPENROUTER_API_KEY")
	if key == "" {
		die(errors.New("dedicated research OpenRouter key missing"))
	}

	fixture, err := loadFixture(*fixturePath)
	if err != nil {
		die(err)
	}
	digests, err := loadDigests(*digestDir, fixture.ExperimentIDs)
	if err != nil {
		die(err)
	}
	pack, err := reasoner.BuildContextPack(reasoner.ContextSelection{
		Schema:        reasoner.ContextSelectionSchema,
		Provider:      "mind-palace",
		Query:         fixture.Query,
		ExperimentIDs: fixture.ExperimentIDs,
	}, digests, reasoner.MaxContextPackBytes)
	if err != nil {
		die(err)
	}

	contractHash := hashText(dryRunContract)
	req, err := reasoner.BuildReasonerRequestFromPack(
		pack,
		"nemotron-context-dryrun-20260928",
		"Wingless-Yggdrasil",
		"MP-NEMO-CONTEXT-DRYRUN-20260928",
		reasoner.RolePlan,
		contractHash,
		2048,
		600,
	)
	if err != nil {
		die(err)
	}
	client, err := reasoner.NewClient(key, reasoner.DefaultConfig())
	if err != nil {
		die(err)
	}
	result, err := client.Invoke(context.Background(), req)
	if err != nil {
		die(err)
	}

	var a answer
	if err := json.Unmarshal([]byte(result.Text), &a); err != nil {
		die(fmt.Errorf("Nemotron reduced-context answer must be exact JSON: %w", err))
	}
	if strings.TrimSpace(a.WinglessNextQuestion) == "" ||
		strings.TrimSpace(a.YggdrasilNextQuestion) == "" {
		die(errors.New("Nemotron omitted a project-specific next question"))
	}
	if a.ConstraintViolation {
		die(errors.New("Nemotron reported a scientific constraint violation"))
	}
	if a.CrossProjectLineageClaim {
		die(errors.New("Nemotron merged independent research lineages"))
	}

	ev := evidence{
		Schema:                      "wingless.nemotron-reduced-context-dryrun-result.v1",
		Model:                       reasoner.DefaultModel,
		SourceRequestID:             fixture.SourceRequestID,
		MindPalaceHead:              fixture.MindPalaceHead,
		RecordsSHA256:               fixture.RecordsSHA256,
		SelectedRecordIDs:           fixture.SelectedRecordIDs,
		ExperimentIDs:               fixture.ExperimentIDs,
		PackSHA256:                  pack.PackSHA256,
		SourceIdentitySHA256:        pack.SourceIdentitySHA256,
		PackBytes:                   pack.EncodedBytes,
		QualificationContractSHA256: contractHash,
		RequestSHA256:               result.RequestSHA256,
		ConfigSHA256:                result.ConfigSHA256,
		ResponseSHA256:              result.ResponseSHA256,
		QualificationIdentitySHA256: result.QualificationIdentitySHA256,
		ReturnedModel:               result.ReturnedModel,
		Provider:                    result.Provider,
		ReasoningTokens:             result.Usage.ReasoningTokens,
		LatencyMS:                   result.LatencyMS,
		Answer:                      a,
		PublicRepositoryOnly:        true,
		PrivateMindPalaceTextSent:   false,
		AdvisoryOnly:                true,
	}
	raw, err := json.MarshalIndent(ev, "", "  ")
	if err != nil {
		die(err)
	}
	raw = append(raw, '\n')
	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		die(err)
	}
	if err := os.WriteFile(*outPath, raw, 0o644); err != nil {
		die(err)
	}

	fmt.Printf("NEMOTRON_REDUCED_CONTEXT_DRYRUN_PASS experiments=%d pack_bytes=%d pack_sha=%s\n",
		len(fixture.ExperimentIDs), pack.EncodedBytes, pack.PackSHA256)
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "nemotron-context-dryrun:", err)
	os.Exit(1)
}
