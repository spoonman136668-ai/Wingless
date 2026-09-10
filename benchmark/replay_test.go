package benchmark

import (
	"context"
	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"testing"
)

type metrics struct{}

func (metrics) Snapshot() (resources.Metrics, error) { return resources.Metrics{}, nil }
func TestReplay(t *testing.T) {
	reg := &broker.Registry{}
	reg.Register(broker.Entry{Backend: &inference.Mock{Name: "fast", Features: []string{"code"}, Reply: "ok"}, Class: "mock", Tier: "fast"})
	r := broker.Runner{Registry: reg, Metrics: metrics{}}
	p := broker.Policy{AllowedBackends: []string{"fast"}}
	for _, config := range []string{"A", "E", "F", "G"} {
		if _, e := Replay(context.Background(), r, p, config, nil); e == nil {
			t.Fatal(config)
		}
	}
	rows, e := Replay(context.Background(), r, p, "B", []Fixture{{"same-task", "reply", "ok"}, {"same-task", "reply", "wrong"}})
	if e != nil || len(rows) != 2 || !rows[0].Correct || rows[1].Correct {
		t.Fatal(rows, e)
	}
}
