package reasoner

import (
	"strings"
	"testing"
)

func TestParseMindPalaceContextEnforcesAdvisoryBoundary(t *testing.T) {
	raw := `{
	  "schema":"ckb-plane.mind-palace-context.v1",
	  "status":"succeeded",
	  "mind_palace_head":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	  "records_sha256":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	  "model_context":"private advisory summary",
	  "primary_ids":["lesson/a"],
	  "support_ids":["lesson/b"],
	  "advisory_only":true,
	  "authorizes_execution":false,
	  "authorizes_acceptance":false,
	  "durable_context_written":false,
	  "git_written":false
	}`
	mp, err := ParseMindPalaceContext(raw)
	if err != nil {
		t.Fatal(err)
	}
	if mp.ModelContext == "" || len(mp.PrimaryIDs) != 1 {
		t.Fatal("parsed context incomplete")
	}
}

func TestResolveMindPalaceContextUsesIDsNotPrivateText(t *testing.T) {
	mp := MindPalaceContext{
		Schema: MindPalaceContextSchema,
		Status: "succeeded",
		MindPalaceHead: strings.Repeat("a",40),
		RecordsSHA256: strings.Repeat("b",64),
		ModelContext: "THIS PRIVATE TEXT MUST NEVER ENTER THE NEMOTRON PACK",
		PrimaryIDs: []string{"lesson/frame"},
		SupportIDs: []string{"lesson/history"},
		AdvisoryOnly: true,
	}
	index := PublicEvidenceIndex{
		Schema: PublicEvidenceIndexSchema,
		Entries: []PublicEvidenceIndexEntry{
			{RecordID:"lesson/frame", ExperimentID:"UP8", Project:"Wingless"},
			{RecordID:"lesson/history", ExperimentID:"YGG-STAB17", Project:"Yggdrasil"},
		},
	}
	sel, err := ResolveMindPalaceContext(mp,index,"next experiment",4)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.Join(sel.ExperimentIDs,","), "PRIVATE") {
		t.Fatal("private model context leaked into selection")
	}
	if strings.Join(sel.ExperimentIDs,",") != "UP8,YGG-STAB17" {
		t.Fatalf("unexpected resolution: %#v",sel.ExperimentIDs)
	}
}

func TestResolveMindPalaceContextFailsWhenNoPublicMappingExists(t *testing.T) {
	mp := MindPalaceContext{
		Schema: MindPalaceContextSchema,
		Status: "succeeded",
		PrimaryIDs: []string{"private/only"},
		AdvisoryOnly: true,
	}
	index := PublicEvidenceIndex{
		Schema: PublicEvidenceIndexSchema,
		Entries: []PublicEvidenceIndexEntry{
			{RecordID:"lesson/public", ExperimentID:"UP8", Project:"Wingless"},
		},
	}
	if _,err:=ResolveMindPalaceContext(mp,index,"q",4); err==nil {
		t.Fatal("expected unresolved public evidence failure")
	}
}
