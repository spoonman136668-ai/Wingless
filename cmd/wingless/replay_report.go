package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/benchmark"
	"io"
	"os"
)

// runReplayReport deterministically replays preserved local-benchmark model
// outputs through the trusted broker verification seam. It launches no model,
// performs no inference, and executes no generated code or tools.
func runReplayReport(args []string) error {
	flags := flag.NewFlagSet("replay-report", flag.ContinueOnError)
	reportPath := flags.String("report", "", "preserved local-benchmark JSON envelope")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || *reportPath == "" {
		return fmt.Errorf("--report FILE required")
	}
	f, err := os.Open(*reportPath)
	if err != nil {
		return err
	}
	raw, err := io.ReadAll(io.LimitReader(f, 8<<20+1))
	f.Close()
	if err != nil {
		return err
	}
	if len(raw) > 8<<20 {
		return fmt.Errorf("report too large")
	}
	var envelope struct {
		Report  benchmark.LocalReport `json:"report"`
		Process json.RawMessage       `json:"process,omitempty"`
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&envelope); err != nil {
		return err
	}
	if err = dec.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("trailing report data")
	}
	replay, err := benchmark.ReplayVerificationCandidates(context.Background(), envelope.Report)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(replay)
}
