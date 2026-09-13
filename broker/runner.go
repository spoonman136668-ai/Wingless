package broker

import (
	"context"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

type Event struct {
	Attempt      int              `json:"attempt"`
	Route        Decision         `json:"route"`
	Result       inference.Result `json:"result"`
	Verification string           `json:"verification"`
	RawModelText string           `json:"raw_model_text,omitempty"`
	Error        string           `json:"error,omitempty"`
}
type Outcome struct {
	Status     string  `json:"status"`
	Events     []Event `json:"events"`
	Acceptance string  `json:"acceptance"`
}

// Verifier is trusted deterministic caller code; it never consumes a model's pass claim as a verdict.
type Verifier interface {
	Verify(context.Context, inference.Result) error
}
type Runner struct {
	Registry *Registry
	Metrics  resources.Provider
	Verifier Verifier
}

func (r Runner) Run(parent context.Context, req inference.Request, p Policy) (out Outcome) {
	out = Outcome{Status: "operator_blocked", Acceptance: "external_required"}
	if err := req.Validate(); err != nil {
		out.Events = append(out.Events, Event{Error: err.Error()})
		return
	}
	if r.Registry == nil || r.Metrics == nil || p.MaxRepairs < 0 || p.MaxRepairs > 8 || p.Repairs < 0 || p.Repairs > p.MaxRepairs {
		out.Events = append(out.Events, Event{Error: "invalid runner policy/dependencies"})
		return
	}
	ctx, c := context.WithDeadline(parent, req.Deadline)
	defer c()
	for n := p.Repairs; n <= p.MaxRepairs; n++ {
		if err := ctx.Err(); err != nil {
			out.Events = append(out.Events, Event{Attempt: n, Error: err.Error()})
			return
		}
		m, e := r.Metrics.Snapshot()
		if e != nil {
			out.Events = append(out.Events, Event{Attempt: n, Error: "resource snapshot unavailable"})
			return
		}
		p.Repairs = n
		b, d, e := Select(ctx, r.Registry, req, p, m)
		event := Event{Attempt: n, Route: d}
		if e != nil {
			event.Error = e.Error()
			out.Events = append(out.Events, event)
			return
		}
		event.Result, e = b.Invoke(ctx, req)
		event.Result.Resources = m
		if e == nil {
			e = ctx.Err()
		}
		if e == nil && event.Result.Status != "completed" {
			e = fmt.Errorf("backend returned noncompleted result")
		}
		if e == nil {
			if r.Verifier == nil {
				event.Verification = "not_run"
				out.Events = append(out.Events, event)
				out.Status = "result_ready"
				return
			}
			verifyResult := event.Result
			if candidate, normalized := inference.VerificationCandidate(event.Result.Text); normalized {
				event.RawModelText = event.Result.Text
				verifyResult.Text = candidate
				event.Result.Text = candidate
			}
			e = r.Verifier.Verify(ctx, verifyResult)
			if e == nil {
				e = ctx.Err()
			}
			if e == nil {
				event.Verification = "fixture_pass"
				out.Events = append(out.Events, event)
				out.Status = "verified_candidate"
				return
			}
			event.Verification = "fixture_fail"
		} else {
			event.Verification = "not_run"
		}
		event.Error = e.Error()
		out.Events = append(out.Events, event)
		// Only deterministic fixture failure can authorize a bounded repair.
		if ctx.Err() != nil || event.Verification != "fixture_fail" {
			return
		}
		if n == p.MaxRepairs {
			return
		}
		feedback := "\nDeterministic fixture failure (untrusted diagnostic): " + event.Error
		if len(req.Context)+len(feedback) > req.MaxContextBytes {
			return
		}
		req.Context += feedback
		p.Reason = "hard_failure"
		if !p.AllowDeep {
			return
		}
	}
	return
}

// PlanThenImplement preserves the heavy-plan / fast-generation seam with a shared deadline.
// Both phases remain candidates, and the plan is untrusted context, not tool authorization.
func (r Runner) PlanThenImplement(ctx context.Context, req inference.Request, p Policy) (Outcome, Outcome) {
	planning := r
	planning.Verifier = nil
	pp := p
	pp.Explicit = ""
	pp.Reason = "explicit_deep"
	pp.AllowFallback = false
	pp.MaxRepairs = 0
	pp.Repairs = 0
	req.Role = "plan"
	plan := planning.Run(ctx, req, pp)
	if plan.Status != "result_ready" {
		return plan, Outcome{Status: "operator_blocked", Acceptance: "external_required"}
	}
	text := plan.Events[len(plan.Events)-1].Result.Text
	req.Context += "\nUntrusted proposed plan:\n" + text
	req.Role = "implement"
	p.Explicit = ""
	p.Reason = "routine"
	return plan, r.Run(ctx, req, p)
}
