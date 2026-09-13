// wingless runs mock demonstrations and explicitly configured local inference experiments.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/benchmark"
	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"os"
	"time"
)

type unknown struct{}

func (unknown) Snapshot() (resources.Metrics, error) { return resources.Metrics{}, nil }
func main() {
	if len(os.Args) > 1 && os.Args[1] == "gpu" {
		if e := runGPU(os.Args[2:]); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "local-benchmark" {
		if e := runLocal(os.Args[2:]); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "replay-report" {
		if e := runReplayReport(os.Args[2:]); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) == 2 && os.Args[1] == "host" {
		h := &resources.Host{Path: "."}
		if _, e := h.Snapshot(); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		time.Sleep(100 * time.Millisecond)
		m, e := h.Snapshot()
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		if e = json.NewEncoder(os.Stdout).Encode(m); e != nil {
			os.Exit(1)
		}
		return
	}

	if len(os.Args) != 2 || (os.Args[1] != "demo" && os.Args[1] != "benchmark") {
		fmt.Fprintln(os.Stderr, "usage: wingless host|demo|benchmark|local-benchmark (--config FILE or --endpoint URL --model ID)|replay-report --report FILE")
		os.Exit(2)
	}
	reg := &broker.Registry{}
	for _, tier := range []string{"fast", "deep"} {
		if e := reg.Register(broker.Entry{Backend: &inference.Mock{Name: tier, Features: []string{"code"}, Reply: "fixture-ok"}, Class: "mock", Tier: tier}); e != nil {
			panic(e)
		}
	}
	runner := broker.Runner{Registry: reg, Metrics: unknown{}}
	p := broker.Policy{AllowedBackends: []string{"fast", "deep"}, AllowDeep: true, MaxRepairs: 1}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if os.Args[1] == "benchmark" {
		rows, e := benchmark.Replay(context.Background(), runner, p, "D", []benchmark.Fixture{{ID: "fixture-1", Prompt: "Return fixture-ok", Expected: "fixture-ok"}})
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
		if e = enc.Encode(rows); e != nil {
			panic(e)
		}
		return
	}
	out := runner.Run(context.Background(), inference.Request{ID: "demo-1", ParentWorkID: "offline-demo", Role: "code", Context: "Return fixture-ok", MaxContextBytes: 4096, MaxOutputTokens: 128, Deadline: time.Now().Add(time.Second), Workspace: "mock-workspace", Capabilities: []string{"code"}}, p)
	if e := enc.Encode(out); e != nil {
		panic(e)
	}
}
