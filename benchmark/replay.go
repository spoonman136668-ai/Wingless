// Package benchmark replays exact fixtures. It does not measure Codex parity.
package benchmark

import (
	"context"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"time"
)

type Fixture struct {
	ID       string
	Prompt   string
	Expected string
}
type Row struct {
	Configuration string         `json:"configuration"`
	Fixture       string         `json:"fixture"`
	Correct       bool           `json:"correct"`
	WallMS        int64          `json:"wall_ms"`
	Outcome       broker.Outcome `json:"outcome"`
}
type exact string

func (e exact) Verify(_ context.Context, r inference.Result) error {
	if r.Text != string(e) {
		return fmt.Errorf("fixture output mismatch")
	}
	return nil
}
func Replay(ctx context.Context, r broker.Runner, p broker.Policy, config string, fixtures []Fixture) ([]Row, error) {
	switch config {
	case "B":
		p.AllowDeep = false
		p.MaxRepairs = 0
		p.Reason = "routine"
	case "C":
		p.AllowDeep = true
		p.Reason = "explicit_deep"
		p.AllowFallback = false
	case "D":
		p.AllowDeep = true
	default:
		return nil, fmt.Errorf("configuration %s unavailable: A requires reference adapter; E context fixture pipeline; F cache; G speculation", config)
	}
	if p.Explicit != "" {
		return nil, fmt.Errorf("benchmark configuration forbids explicit route override")
	}
	rows := []Row{}
	for i, f := range fixtures {
		r.Verifier = exact(f.Expected)
		req := inference.Request{ID: fmt.Sprintf("replay-%s-%d", config, i), ParentWorkID: f.ID, Role: "fixture", Context: f.Prompt, MaxContextBytes: 4096, MaxOutputTokens: 256, Deadline: time.Now().Add(10 * time.Second), Workspace: "mock-fixture", Capabilities: []string{"code"}}
		start := time.Now()
		out := r.Run(ctx, req, p)
		rows = append(rows, Row{config, f.ID, out.Status == "verified_candidate", time.Since(start).Milliseconds(), out})
	}
	return rows, nil
}
