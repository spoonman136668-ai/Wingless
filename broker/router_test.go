package broker

import (
	"context"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"testing"
	"time"
)

func fixture() (*Registry, inference.Request, Policy) {
	r := &Registry{}
	for _, tier := range []string{"fast", "deep"} {
		r.Register(Entry{&inference.Mock{Name: tier, Features: []string{"code"}, Reply: tier}, "mock", tier})
	}
	return r, inference.Request{ID: "r", ParentWorkID: "p", Role: "code", Context: "x", MaxContextBytes: 4096, MaxOutputTokens: 32, Deadline: time.Now().Add(time.Second), Workspace: "fixture", Capabilities: []string{"code"}}, Policy{AllowedBackends: []string{"fast", "deep"}, AllowDeep: true, MaxRepairs: 1}
}
func TestRoutes(t *testing.T) {
	for _, kind := range []string{"fast", "deep", "explicit", "unauthorized", "capability", "health", "fallback", "budget", "resource"} {
		t.Run(kind, func(t *testing.T) {
			r, q, p := fixture()
			want := "fast"
			denied := false
			switch kind {
			case "deep":
				p.Reason = "hard_failure"
				want = "deep"
			case "explicit":
				p.Explicit = "fast"
				p.Reason = "hard_failure"
			case "unauthorized":
				p.AllowedBackends = nil
				denied = true
			case "capability":
				q.Capabilities = []string{"shell"}
				denied = true
			case "health", "fallback":
				p.Reason = "hard_failure"
				r.entries["deep"].Backend.(*inference.Mock).Unhealthy = true
				if kind == "fallback" {
					p.AllowFallback = true
				} else {
					denied = true
				}
			case "budget":
				p.Repairs = 2
				denied = true
			case "resource":
				q.Resources.MinRAM = 1
				denied = true
			}
			b, d, e := Select(context.Background(), r, q, p, resources.Metrics{})
			if denied {
				if e == nil {
					t.Fatal("expected denial", d)
				}
				return
			}
			if e != nil || b.ID() != want || d.Backend != want {
				t.Fatal(d, e)
			}
		})
	}
}
func TestRegistration(t *testing.T) {
	r, _, _ := fixture()
	if e := r.Register(r.entries["fast"]); e == nil {
		t.Fatal("duplicate")
	}
	if len(r.Entries()[0].Backend.Capabilities()) != 1 {
		t.Fatal("capabilities")
	}
}

type metrics struct{}

func (metrics) Snapshot() (resources.Metrics, error) { return resources.Metrics{}, nil }

type verifier struct{}

func (verifier) Verify(_ context.Context, r inference.Result) error {
	if r.Text == "deep" {
		return nil
	}
	return fmt.Errorf("fixture failed")
}
func TestBoundedRepairAndAuthority(t *testing.T) {
	r, q, p := fixture()
	out := (Runner{r, metrics{}, verifier{}}).Run(context.Background(), q, p)
	if out.Status != "verified_candidate" || len(out.Events) != 2 || out.Events[1].Route.Backend != "deep" || out.Acceptance != "external_required" {
		t.Fatal(out)
	}
	p.MaxRepairs = 0
	out = (Runner{r, metrics{}, verifier{}}).Run(context.Background(), q, p)
	if out.Status != "operator_blocked" || len(out.Events) != 1 {
		t.Fatal(out)
	}
	out = (Runner{r, metrics{}, nil}).Run(context.Background(), q, p)
	if out.Status != "result_ready" || out.Acceptance != "external_required" {
		t.Fatal(out)
	}
}
func TestPlanThenFast(t *testing.T) {
	r, q, p := fixture()
	a, b := (Runner{r, metrics{}, nil}).PlanThenImplement(context.Background(), q, p)
	if a.Status != "result_ready" || b.Status != "result_ready" || a.Events[0].Route.Tier != "deep" || b.Events[0].Route.Tier != "fast" {
		t.Fatal(a, b)
	}
}
