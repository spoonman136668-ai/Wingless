package cognitive

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func TestFileStoreMemoryPersistenceAndDeletion(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileStore(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	rec, err := store.PutMemory(ctx, MemoryRecord{
		Version: 1, Class: MemorySemantic, Key: "answer",
		Content: "42", Provenance: "fixture",
	})
	if err != nil {
		t.Fatal(err)
	}

	reopened, err := NewFileStore(root)
	if err != nil {
		t.Fatal(err)
	}
	got, ok, err := reopened.LookupMemory(ctx, MemorySemantic, "answer")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || got.ID != rec.ID || got.Content != "42" {
		t.Fatalf("persisted memory mismatch: ok=%v got=%+v", ok, got)
	}
	if err := reopened.DeleteMemory(ctx, rec.ID); err != nil {
		t.Fatal(err)
	}
	_, ok, err = reopened.LookupMemory(ctx, MemorySemantic, "answer")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("deleted memory remained visible")
	}
}

func TestSkillCannotSelfDeclareValidation(t *testing.T) {
	store, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	_, err = store.PutSkillCandidate(ctx, Skill{
		Version: 1, Source: "model", TaskFamily: "fixture",
		Kind: SkillStaticText, Body: "answer", Validation: ValidationValidated,
	})
	if err == nil {
		t.Fatal("model-shaped candidate self-declared validation")
	}

	candidate, err := store.PutSkillCandidate(ctx, Skill{
		Version: 1, Source: "model", TaskFamily: "fixture",
		Kind: SkillStaticText, Body: "answer",
	})
	if err != nil {
		t.Fatal(err)
	}
	skills, err := store.ListValidatedSkills(ctx, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 0 {
		t.Fatalf("candidate visible as validated: %+v", skills)
	}

	sum := sha256.Sum256([]byte("deterministic-evidence"))
	_, err = store.PromoteSkill(ctx, candidate.ID, ValidationEvidence{
		ValidatorID:    "fixture-validator",
		EvidenceSHA256: hex.EncodeToString(sum[:]),
	})
	if err != nil {
		t.Fatal(err)
	}
	skills, err = store.ListValidatedSkills(ctx, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].ID != candidate.ID {
		t.Fatalf("validated skill missing: %+v", skills)
	}
}

func TestRoutingDecisionReplayStable(t *testing.T) {
	store, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	_, err = store.PutMemory(ctx, MemoryRecord{
		Version: 1, Class: MemorySemantic, Key: "known",
		Content: "value", Provenance: "fixture",
	})
	if err != nil {
		t.Fatal(err)
	}
	req := baseRun("replay")
	req.MemoryClass, req.MemoryKey = MemorySemantic, "known"
	policy := Policy{Version: "replay-v1", MaxPasses: 2, EnableMemory: true}
	router := Router{Policy: policy, Memory: store}
	a, _, _, err := router.Decide(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	b, _, _, err := router.Decide(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != b.ID || a.InputDigest != b.InputDigest || a.Route != b.Route {
		t.Fatalf("routing decision not replay-stable: a=%+v b=%+v", a, b)
	}
}

func TestStoreRejectsUnboundedRecords(t *testing.T) {
	store, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	_, err = store.PutMemory(ctx, MemoryRecord{
		Version: 1, Class: MemorySemantic, Key: "oversize",
		Content: strings.Repeat("x", maxMemoryContentBytes+1), Provenance: "fixture",
	})
	if err == nil {
		t.Fatal("oversized memory record was accepted")
	}
	_, err = store.PutSkillCandidate(ctx, Skill{
		Version: 1, Source: "fixture", TaskFamily: "fixture",
		Kind: SkillStaticText, Body: strings.Repeat("x", maxSkillBodyBytes+1),
	})
	if err == nil {
		t.Fatal("oversized skill record was accepted")
	}
}

func TestContentAddressedPutIsIdempotent(t *testing.T) {
	store, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	first, err := store.PutMemory(ctx, MemoryRecord{
		Version: 1, Class: MemorySemantic, Key: "stable",
		Content: "value", Provenance: "fixture",
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.PutMemory(ctx, MemoryRecord{
		Version: 1, Class: MemorySemantic, Key: "stable",
		Content: "value", Provenance: "fixture",
	})
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID || !first.CreatedAt.Equal(second.CreatedAt) {
		t.Fatalf("content-addressed duplicate was not idempotent: first=%+v second=%+v", first, second)
	}
}
