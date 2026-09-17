package nr1

import (
	"strings"
	"testing"
)

func syntheticTrace() []TraceEvent {
	return []TraceEvent{
		{Schema: TraceSchema, SessionID: "s1", WorkloadID: "w1", TaskFamily: "repair", Phase: "prefill", TokenIndex: 0, Layer: 0, Experts: []int{0, 1}},
		{Schema: TraceSchema, SessionID: "s1", WorkloadID: "w1", TaskFamily: "repair", Phase: "prefill", TokenIndex: 0, Layer: 1, Experts: []int{2, 3}},
		{Schema: TraceSchema, SessionID: "s1", WorkloadID: "w1", TaskFamily: "repair", Phase: "decode", TokenIndex: 1, Layer: 0, Experts: []int{0, 1}},
		{Schema: TraceSchema, SessionID: "s1", WorkloadID: "w1", TaskFamily: "repair", Phase: "decode", TokenIndex: 1, Layer: 1, Experts: []int{2, 4}},
		{Schema: TraceSchema, SessionID: "s1", WorkloadID: "w1", TaskFamily: "repair", Phase: "decode", TokenIndex: 2, Layer: 0, Experts: []int{0, 5}},
		{Schema: TraceSchema, SessionID: "s1", WorkloadID: "w1", TaskFamily: "repair", Phase: "decode", TokenIndex: 2, Layer: 1, Experts: []int{2, 3}},
	}
}

func TestReadJSONLStrict(t *testing.T) {
	input := strings.Join([]string{
		`{"schema":"wingless.nr1.router-trace.v2","session_id":"s1","workload_id":"w1","task_family":"code","phase":"prefill","token_index":0,"layer":0,"experts":[1,2],"scores":[0.7,0.3]}`,
		`{"schema":"wingless.nr1.router-trace.v2","session_id":"s1","workload_id":"w1","task_family":"code","phase":"decode","token_index":1,"layer":0,"experts":[1,3]}`,
	}, "\n")
	got, err := ReadJSONL(strings.NewReader(input), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[1].Experts[1] != 3 || got[1].Phase != "decode" {
		t.Fatalf("unexpected trace: %+v", got)
	}

	bad := `{"schema":"wingless.nr1.router-trace.v2","session_id":"s1","workload_id":"w1","task_family":"code","phase":"prefill","token_index":0,"layer":0,"experts":[1,1]}`
	if _, err = ReadJSONL(strings.NewReader(bad), 10); err == nil {
		t.Fatal("duplicate expert accepted")
	}

	unknown := `{"schema":"wingless.nr1.router-trace.v2","session_id":"s1","workload_id":"w1","task_family":"code","phase":"prefill","token_index":0,"layer":0,"experts":[1],"extra":true}`
	if _, err = ReadJSONL(strings.NewReader(unknown), 10); err == nil {
		t.Fatal("unknown field accepted")
	}

	missingPhase := `{"schema":"wingless.nr1.router-trace.v2","session_id":"s1","workload_id":"w1","task_family":"code","token_index":0,"layer":0,"experts":[1]}`
	if _, err = ReadJSONL(strings.NewReader(missingPhase), 10); err == nil {
		t.Fatal("missing phase accepted")
	}
}

func TestAnalyzeLocality(t *testing.T) {
	report, err := Analyze(syntheticTrace())
	if err != nil {
		t.Fatal(err)
	}
	if report.Schema != LocalityReportSchema || report.Events != 6 || report.Tokens != 3 || report.Sessions != 1 || len(report.Layers) != 2 {
		t.Fatalf("unexpected report header: %+v", report)
	}
	if report.AdjacentPairs != 4 {
		t.Fatalf("adjacent pairs=%d want=4", report.AdjacentPairs)
	}
	// Layer 0 overlaps 2/2 then 1/2; layer 1 overlaps 1/2 then 1/2.
	// Mean over all four adjacent layer-token pairs = 0.625.
	if report.AdjacentReuseRate != 0.625 {
		t.Fatalf("reuse=%v want=0.625", report.AdjacentReuseRate)
	}
	if report.Layers[0].Coverage[0].Experts < 2 {
		t.Fatalf("coverage unexpectedly narrow: %+v", report.Layers[0].Coverage)
	}
	if len(report.TopTransitions) == 0 {
		t.Fatal("missing transitions")
	}
}

func TestResidencySimulation(t *testing.T) {
	events := syntheticTrace()
	const expertBytes = uint64(1024)
	report, err := Simulate(events, ResidencyConfig{
		ExpertBytes:   expertBytes,
		ReservedBytes: 2048,
		BudgetsBytes:  []uint64{2048, 4096, 16384},
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.Schema != ResidencyReportSchema || len(report.Results) != 15 {
		t.Fatalf("unexpected report: %+v", report)
	}
	for _, result := range report.Results {
		if result.Hits+result.Misses != result.Accesses {
			t.Fatalf("accounting mismatch: %+v", result)
		}
		if result.BudgetBytes == 2048 && result.Hits != 0 {
			t.Fatalf("zero-slot cache hit: %+v", result)
		}
		if result.BudgetBytes == 16384 && result.Slots >= 8 && result.Policy != "static-global" && result.Policy != "session-static" && result.HitRate <= 0 {
			t.Fatalf("large cache made no reuse: %+v", result)
		}
		switch result.Policy {
		case "static-global":
			if !result.UsesFutureTrace || result.CacheLifecycle != "persistent_across_sessions" {
				t.Fatalf("static-global oracle metadata invalid: %+v", result)
			}
		case "session-static":
			if !result.UsesFutureTrace || result.CacheLifecycle != "reset_per_session" {
				t.Fatalf("session-static oracle metadata invalid: %+v", result)
			}
		default:
			if result.UsesFutureTrace || result.CacheLifecycle != "persistent_across_sessions" {
				t.Fatalf("online policy metadata invalid: %+v", result)
			}
		}
	}
}

func TestSessionStaticIsOracleBound(t *testing.T) {
	events := syntheticTrace()
	report, err := Simulate(events, ResidencyConfig{ExpertBytes: 1, BudgetsBytes: []uint64{4}})
	if err != nil {
		t.Fatal(err)
	}
	var sessionStatic, lru *ResidencyResult
	for i := range report.Results {
		r := &report.Results[i]
		switch r.Policy {
		case "session-static":
			sessionStatic = r
		case "lru":
			lru = r
		}
	}
	if sessionStatic == nil || lru == nil {
		t.Fatal("missing policies")
	}
	if sessionStatic.HitRate < 0 || sessionStatic.HitRate > 1 || lru.HitRate < 0 || lru.HitRate > 1 {
		t.Fatal("invalid hit rate")
	}
	if !sessionStatic.UsesFutureTrace || lru.UsesFutureTrace {
		t.Fatal("oracle classification mismatch")
	}
}
