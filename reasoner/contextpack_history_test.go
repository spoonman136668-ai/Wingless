package reasoner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type liveSelectionFixture struct {
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

func readHistoryDigest(t *testing.T, id string) ExperimentDigest {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("digests", "history", id+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var d ExperimentDigest
	if err := json.Unmarshal(raw, &d); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatalf("%s digest invalid: %v", id, err)
	}
	return d
}

func TestLiveMindPalaceSelectionBuildsPublicOnlyReducedPack(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("fixtures", "context", "mind-palace-live-selection-20260928.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixture liveSelectionFixture
	if err := json.Unmarshal(raw, &fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.Schema != "wingless.mind-palace-public-selection-fixture.v1" {
		t.Fatalf("fixture schema=%q", fixture.Schema)
	}
	if len(fixture.MindPalaceHead) != 40 || !isSHA256(fixture.RecordsSHA256) {
		t.Fatal("fixture Mind-Palace provenance invalid")
	}
	if fixture.PrivateModelContextIncluded || !fixture.PublicReconstructionRequired {
		t.Fatal("fixture violates public reconstruction boundary")
	}
	if len(fixture.ExperimentIDs) != 4 {
		t.Fatalf("experiment count=%d want=4", len(fixture.ExperimentIDs))
	}

	digests := make([]ExperimentDigest, 0, len(fixture.ExperimentIDs))
	for _, id := range fixture.ExperimentIDs {
		digests = append(digests, readHistoryDigest(t, id))
	}
	pack, err := BuildContextPack(ContextSelection{
		Schema:        ContextSelectionSchema,
		Provider:      "mind-palace",
		Query:         fixture.Query,
		ExperimentIDs: fixture.ExperimentIDs,
	}, digests, MaxContextPackBytes)
	if err != nil {
		t.Fatal(err)
	}
	if pack.EncodedBytes <= 0 || pack.EncodedBytes > MaxContextPackBytes {
		t.Fatalf("pack bytes=%d", pack.EncodedBytes)
	}
	ctx, err := pack.ContextForReasoner()
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range append([]string{
		fixture.MindPalaceHead,
		fixture.RecordsSHA256,
		fixture.SourceRequestID,
		"model_context",
	}, fixture.SelectedRecordIDs...) {
		if strings.Contains(ctx, forbidden) {
			t.Fatalf("private/orchestration provenance leaked into free-route pack: %q", forbidden)
		}
	}
	if !strings.Contains(ctx, "UP-155B") ||
		!strings.Contains(ctx, "UP-134C") ||
		!strings.Contains(ctx, "UP-LM2A") ||
		!strings.Contains(ctx, "YGG-B36") {
		t.Fatal("public selected experiment evidence missing from pack")
	}
}
