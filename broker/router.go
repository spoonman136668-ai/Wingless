package broker

import (
	"context"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"sort"
	"sync"
)

type Entry struct {
	Backend inference.InferenceBackend
	Class   string
	Tier    string
}
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

func (r *Registry) Register(e Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e.Backend == nil || e.Backend.ID() == "" || (e.Tier != "fast" && e.Tier != "deep" && e.Tier != "external") {
		return fmt.Errorf("invalid registration")
	}
	switch e.Class {
	case "mock", "codex-adapter", "local-openai-compatible", "local-fast-coder", "local-deep-reasoner":
	default:
		return fmt.Errorf("unknown backend class")
	}
	if r.entries == nil {
		r.entries = map[string]Entry{}
	}
	if _, ok := r.entries[e.Backend.ID()]; ok {
		return fmt.Errorf("duplicate backend")
	}
	r.entries[e.Backend.ID()] = e
	return nil
}
func (r *Registry) Entries() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Entry{}
	for _, e := range r.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Backend.ID() < out[j].Backend.ID() })
	return out
}

type Policy struct {
	Explicit        string
	AllowedBackends []string
	AllowDeep       bool
	AllowFallback   bool
	Reason          string
	Repairs         int
	MaxRepairs      int
}
type Decision struct {
	Backend  string   `json:"backend"`
	Tier     string   `json:"tier"`
	Reason   string   `json:"reason"`
	Rejected []string `json:"rejected"`
}

func contains(a []string, x string) bool {
	for _, s := range a {
		if s == x {
			return true
		}
	}
	return false
}
func Select(ctx context.Context, reg *Registry, req inference.Request, p Policy, m resources.Metrics) (inference.InferenceBackend, Decision, error) {
	d := Decision{Reason: "automatic_fast"}
	if e := req.Validate(); e != nil {
		return nil, d, e
	}
	if p.Repairs < 0 || p.MaxRepairs < 0 || p.MaxRepairs > 8 || p.Repairs > p.MaxRepairs {
		return nil, d, fmt.Errorf("repair budget exhausted")
	}
	if e := resources.Check(req.Resources, m); e != nil {
		return nil, d, e
	}
	tier := "fast"
	switch p.Reason {
	case "", "routine":
	case "low_confidence", "hard_failure", "protocol_failure", "ambiguous_architecture", "explicit_deep":
		if !p.AllowDeep && p.Explicit == "" {
			return nil, d, fmt.Errorf("deep escalation unauthorized")
		}
		tier = "deep"
		d.Reason = p.Reason
	default:
		return nil, d, fmt.Errorf("unknown escalation reason")
	}
	if p.Explicit != "" {
		d.Reason = "explicit_selection"
	}
	entries := reg.Entries()
	for pass := 0; pass < 2; pass++ {
		for _, e := range entries {
			b := e.Backend
			if !contains(p.AllowedBackends, b.ID()) {
				continue
			}
			if p.Explicit != "" {
				if b.ID() != p.Explicit {
					continue
				}
			} else if e.Tier != tier {
				continue
			}
			if e.Tier == "deep" && !p.AllowDeep {
				d.Rejected = append(d.Rejected, b.ID()+": unauthorized deep")
				continue
			}
			ok := true
			for _, c := range req.Capabilities {
				if !contains(b.Capabilities(), c) {
					ok = false
				}
			}
			if !ok {
				d.Rejected = append(d.Rejected, b.ID()+": capability")
				continue
			}
			if err := b.Health(ctx); err != nil {
				d.Rejected = append(d.Rejected, b.ID()+": health")
				continue
			}
			d.Backend = b.ID()
			d.Tier = e.Tier
			return b, d, nil
		}
		if p.Explicit != "" || !p.AllowFallback || tier != "deep" {
			break
		}
		tier = "fast"
		d.Reason = "authorized_fast_fallback"
	}
	return nil, d, fmt.Errorf("no authorized healthy capable backend")
}
