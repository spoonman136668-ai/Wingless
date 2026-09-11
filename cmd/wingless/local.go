package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/benchmark"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/modelhost"
	"github.com/spoonman136668-ai/Wingless/resources"
	"io"
	"os"
	"os/signal"
	"time"
)

// runLocal is an operator-only local inference experiment, never a plane worker/tool.
func runLocal(args []string) error {
	flags := flag.NewFlagSet("local-benchmark", flag.ContinueOnError)
	config := flags.String("config", "", "pinned llama-server/model JSON configuration")
	endpoint := flags.String("endpoint", "", "existing numeric loopback endpoint (no process launch)")
	model := flags.String("model", "", "existing endpoint model ID")
	repeats := flags.Int("repeats", 1, "fixed fixture repetitions, 1..3")
	if e := flags.Parse(args); e != nil {
		return e
	}
	if flags.NArg() != 0 || (*config == "") == (*endpoint == "") {
		return fmt.Errorf("choose exactly one of --config or --endpoint")
	}
	if *repeats < 1 || *repeats > 3 {
		return fmt.Errorf("repeats must be 1..3")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	host := &resources.Host{Path: "."}
	var supervisor *modelhost.Supervisor
	url, id := *endpoint, *model
	if *config != "" {
		f, e := os.Open(*config)
		if e != nil {
			return e
		}
		raw, e := io.ReadAll(io.LimitReader(f, 16385))
		f.Close()
		if e != nil {
			return e
		}
		if len(raw) > 16384 {
			return fmt.Errorf("configuration too large")
		}
		var cfg modelhost.Config
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		if e = dec.Decode(&cfg); e != nil {
			return e
		}
		if e = dec.Decode(new(any)); e != io.EOF {
			return fmt.Errorf("trailing config data")
		}
		supervisor, e = modelhost.New(cfg, host)
		if e != nil {
			return e
		}
		defer func() {
			stopCtx, c := context.WithTimeout(context.Background(), 5*time.Second)
			defer c()
			supervisor.Stop(stopCtx)
		}()
		if e = supervisor.Start(ctx); e != nil {
			json.NewEncoder(os.Stderr).Encode(supervisor.Status())
			return e
		}
		url, id = supervisor.Endpoint(), cfg.ModelID
	}
	b, e := inference.NewStreamingLocalHTTP("local-benchmark", id, url, []string{"code"})
	if e != nil {
		return e
	}
	defer b.Close()
	report, e := benchmark.RunLocal(ctx, b, host, *repeats)
	if supervisor != nil {
		stopCtx, c := context.WithTimeout(context.Background(), 5*time.Second)
		stopErr := supervisor.Stop(stopCtx)
		c()
		if stopErr != nil {
			return stopErr
		}
	}
	envelope := struct {
		Report  benchmark.LocalReport `json:"report"`
		Process *modelhost.Evidence   `json:"process,omitempty"`
	}{Report: report}
	if supervisor != nil {
		state := supervisor.Status()
		envelope.Process = &state
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if encodeErr := enc.Encode(envelope); encodeErr != nil {
		return encodeErr
	}
	return e
}
