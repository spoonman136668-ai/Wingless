package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/spoonman136668-ai/Wingless/research/nr1"
)

const reportSchema = "wingless.nr1a.report.v1"

type report struct {
	Schema     string              `json:"schema"`
	TracePath  string              `json:"trace_path"`
	Locality   nr1.LocalityReport  `json:"locality"`
	Residency  nr1.ResidencyReport `json:"residency"`
}

func main() {
	var (
		tracePath     string
		outPath       string
		budgetsGiB    string
		expertBytes   uint64
		reservedBytes uint64
		maxEvents     int
	)
	flag.StringVar(&tracePath, "trace", "", "NR-1 router trace JSONL")
	flag.StringVar(&outPath, "out", "", "output report JSON (stdout when empty)")
	flag.StringVar(&budgetsGiB, "budgets-gib", "4,6,8,12,16", "comma-separated total residency budgets in GiB")
	flag.Uint64Var(&expertBytes, "expert-bytes", 0, "encoded bytes per layer/expert object; required")
	flag.Uint64Var(&reservedBytes, "reserved-bytes", 0, "bytes reserved for core, KV cache, and runtime buffers")
	flag.IntVar(&maxEvents, "max-events", 10_000_000, "maximum trace events")
	flag.Parse()

	if tracePath == "" || expertBytes == 0 || maxEvents < 1 {
		fatal(errors.New("-trace, -expert-bytes, and positive -max-events are required"))
	}
	budgets, err := parseBudgets(budgetsGiB)
	if err != nil {
		fatal(err)
	}

	f, err := os.Open(tracePath)
	if err != nil {
		fatal(err)
	}
	events, err := nr1.ReadJSONL(f, maxEvents)
	_ = f.Close()
	if err != nil {
		fatal(err)
	}

	locality, err := nr1.Analyze(events)
	if err != nil {
		fatal(err)
	}
	residency, err := nr1.Simulate(events, nr1.ResidencyConfig{
		ExpertBytes: expertBytes, ReservedBytes: reservedBytes, BudgetsBytes: budgets,
	})
	if err != nil {
		fatal(err)
	}

	payload, err := json.MarshalIndent(report{
		Schema: reportSchema, TracePath: tracePath, Locality: locality, Residency: residency,
	}, "", "  ")
	if err != nil {
		fatal(err)
	}
	payload = append(payload, '\n')
	if outPath == "" {
		_, err = os.Stdout.Write(payload)
	} else {
		err = os.WriteFile(outPath, payload, 0o600)
	}
	if err != nil {
		fatal(err)
	}
}

func parseBudgets(raw string) ([]uint64, error) {
	parts := strings.Split(raw, ",")
	if len(parts) == 0 {
		return nil, errors.New("no budgets")
	}
	out := make([]uint64, 0, len(parts))
	for _, part := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 || v > 1024 {
			return nil, fmt.Errorf("invalid GiB budget %q", part)
		}
		bytes := v * 1024 * 1024 * 1024
		if bytes > math.MaxUint64 {
			return nil, fmt.Errorf("GiB budget overflow %q", part)
		}
		out = append(out, uint64(math.Round(bytes)))
	}
	return out, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "nr1a:", err)
	os.Exit(1)
}
