package worker

import (
	"context"
	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"testing"
	"time"
)

type metrics struct{}

func (metrics) Snapshot() (resources.Metrics, error) { return resources.Metrics{}, nil }
func TestSingleWorkerCancel(t *testing.T) {
	reg := &broker.Registry{}
	reg.Register(broker.Entry{Backend: &inference.Mock{Name: "fast", Features: []string{"code"}, Delay: time.Second}, Class: "mock", Tier: "fast"})
	s := New(broker.Runner{Registry: reg, Metrics: metrics{}}, broker.Policy{AllowedBackends: []string{"fast"}})
	r := inference.Request{ID: "1", ParentWorkID: "p", Role: "code", MaxContextBytes: 1, MaxOutputTokens: 1, Deadline: time.Now().Add(2 * time.Second), Workspace: "fixture", Capabilities: []string{"code"}}
	id, e := s.SubmitWork(context.Background(), r)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.SubmitWork(context.Background(), r); e == nil {
		t.Fatal("duplicate")
	}
	if !s.Cancel(id) {
		t.Fatal("cancel")
	}
	until := time.Now().Add(time.Second)
	for {
		state, _ := s.Status(id)
		if state != "running" {
			break
		}
		if time.Now().After(until) {
			t.Fatal("no terminal state")
		}
		time.Sleep(time.Millisecond)
	}
	out, e := s.Result(id)
	if e != nil || out.Acceptance != "external_required" || out.Status != "operator_blocked" {
		t.Fatal(out, e)
	}
}
