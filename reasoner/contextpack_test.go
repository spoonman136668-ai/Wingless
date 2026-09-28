package reasoner

import (
	"errors"
	"strings"
	"testing"
)

func digestFixture(id string) ExperimentDigest {
	return ExperimentDigest{
		Schema:            ExperimentDigestSchema,
		Project:           "Wingless",
		ExperimentID:      id,
		SourceRepository:  "spoonman136668-ai/Wingless",
		SourceRef:         strings.Repeat("a", 40),
		SourcePath:        "docs/experiments/" + id + ".md",
		SourceSHA256:      strings.Repeat("b", 64),
		Question:          "What mechanism should be tested next?",
		Observed:          []string{"A sealed result isolated a frame-dependent decoding failure."},
		Established:       []string{"The signal survives at a fixed frame."},
		NotEstablished:    []string{"General cross-depth decoding is not established."},
		FrozenConstraints: []string{"Do not retune thresholds or held-out splits."},
		NextSeams:         []string{"Test a compact co-evolving reference."},
		EvidenceRefs:      []string{"sealed-result:" + id},
	}
}

func TestBuildContextPackDeterministicAndBounded(t *testing.T) {
	a := digestFixture("UP8")
	b := digestFixture("UP23")
	selection := ContextSelection{
		Schema: ContextSelectionSchema,
		Provider: "mind-palace",
		Query: "frame ambiguity next experiment",
		ExperimentIDs: []string{"UP8", "UP23"},
	}
	p1, err := BuildContextPack(selection, []ExperimentDigest{a, b}, MaxContextPackBytes)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := BuildContextPack(selection, []ExperimentDigest{b, a}, MaxContextPackBytes)
	if err != nil {
		t.Fatal(err)
	}
	if p1.PackSHA256 != p2.PackSHA256 || p1.SourceIdentitySHA256 != p2.SourceIdentitySHA256 {
		t.Fatal("context pack identity changed with available digest ordering")
	}
	if strings.Join(p1.ExperimentIDs, ",") != "UP23,UP8" {
		t.Fatalf("selection not canonicalized: %#v", p1.ExperimentIDs)
	}
	ctx, err := p1.ContextForReasoner()
	if err != nil || !strings.Contains(ctx, "public-repository") {
		t.Fatalf("context generation failed: %v %s", err, ctx)
	}
}

func TestContextPackFailsRatherThanTruncates(t *testing.T) {
	d := digestFixture("UP8")
	selection := ContextSelection{
		Schema: ContextSelectionSchema,
		Provider: "mind-palace",
		Query: "frame ambiguity",
		ExperimentIDs: []string{"UP8"},
	}
	if _, err := BuildContextPack(selection, []ExperimentDigest{d}, 200); err == nil {
		t.Fatal("expected compact pack budget failure")
	}
}

func TestFreePackRejectsPrivateOrMissingSource(t *testing.T) {
	d := digestFixture("UP8")
	d.SourceRepository = "private/research"
	if err := d.Validate(); err == nil {
		t.Fatal("expected private source rejection")
	}
	selection := ContextSelection{
		Schema: ContextSelectionSchema,
		Provider: "mind-palace",
		Query: "frame ambiguity",
		ExperimentIDs: []string{"UP8"},
	}
	if _, err := BuildContextPack(selection, []ExperimentDigest{digestFixture("UP9")}, MaxContextPackBytes); err == nil {
		t.Fatal("expected missing selected digest rejection")
	}
}

type selectorStub struct {
	s   ContextSelection
	err error
}
func (s selectorStub) Select(string, int) (ContextSelection, error) { return s.s, s.err }

func TestLiveSelectionRequiresMindPalace(t *testing.T) {
	if _, err := RequireContextSelection(nil, "q", 4); !errors.Is(err, ErrContextSelectorUnavailable) {
		t.Fatalf("expected unavailable selector, got %v", err)
	}
	fallback := selectorStub{s: ContextSelection{
		Schema: ContextSelectionSchema,
		Provider: "explicit-public-fallback",
		Query: "q",
		ExperimentIDs: []string{"UP8"},
	}}
	if _, err := RequireContextSelection(fallback, "q", 4); err == nil {
		t.Fatal("live selection must not silently bypass Mind-Palace")
	}
	mp := selectorStub{s: ContextSelection{
		Schema: ContextSelectionSchema,
		Provider: "mind-palace",
		Query: "q",
		ExperimentIDs: []string{"UP8"},
	}}
	if _, err := RequireContextSelection(mp, "q", 4); err != nil {
		t.Fatal(err)
	}
}

func TestBuildReasonerRequestFromPackUsesPackIdentity(t *testing.T) {
	d := digestFixture("UP8")
	sel := ContextSelection{
		Schema: ContextSelectionSchema,
		Provider: "mind-palace",
		Query: "frame ambiguity",
		ExperimentIDs: []string{"UP8"},
	}
	pack, err := BuildContextPack(sel, []ExperimentDigest{d}, MaxContextPackBytes)
	if err != nil {
		t.Fatal(err)
	}
	req, err := BuildReasonerRequestFromPack(pack, "r1", "Wingless", "next", RolePlan, strings.Repeat("c", 64), 1024, 120)
	if err != nil {
		t.Fatal(err)
	}
	if req.FrontierSHA256 != pack.PackSHA256 || req.DataClass != "public-repository" {
		t.Fatalf("reasoner request not bound to reduced pack: %#v", req)
	}
}
