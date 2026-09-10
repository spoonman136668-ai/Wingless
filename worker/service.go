// Package worker is a provisional in-memory plane boundary, not a durable queue.
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/broker"
	"github.com/spoonman136668-ai/Wingless/inference"
	"sync"
)

type Boundary interface {
	SubmitWork(context.Context, inference.Request) (string, error)
	Status(string) (string, error)
	Cancel(string) bool
	Result(string) (broker.Outcome, error)
	Health() error
	Capabilities() []string
}
type job struct {
	status string
	cancel context.CancelFunc
	result broker.Outcome
}
type Service struct {
	mu     sync.Mutex
	runner broker.Runner
	policy broker.Policy
	jobs   map[string]*job
	active bool
}

var _ Boundary = (*Service)(nil)

func New(r broker.Runner, p broker.Policy) *Service {
	p.AllowedBackends = append([]string(nil), p.AllowedBackends...)
	return &Service{runner: r, policy: p, jobs: map[string]*job{}}
}
func (s *Service) Capabilities() []string {
	return []string{"bounded_inference", "in_memory_status", "cancel", "candidate_result"}
}
func (s *Service) Health() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runner.Registry == nil || s.runner.Metrics == nil {
		return fmt.Errorf("unconfigured service")
	}
	if s.active {
		return fmt.Errorf("busy")
	}
	if len(s.jobs) >= 256 {
		return fmt.Errorf("result retention full")
	}
	return nil
}
func (s *Service) SubmitWork(parent context.Context, r inference.Request) (string, error) {
	if e := r.Validate(); e != nil {
		return "", e
	}
	if e := parent.Err(); e != nil {
		return "", e
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		return "", fmt.Errorf("single worker busy")
	}
	if _, ok := s.jobs[r.ID]; ok {
		return "", fmt.Errorf("duplicate request id")
	}
	if len(s.jobs) >= 256 {
		return "", fmt.Errorf("result retention full")
	}
	ctx, c := context.WithDeadline(parent, r.Deadline)
	r.Capabilities = append([]string(nil), r.Capabilities...)
	j := &job{status: "running", cancel: c}
	s.jobs[r.ID] = j
	s.active = true
	go func() {
		out := s.runner.Run(ctx, r, s.policy)
		c()
		s.mu.Lock()
		defer s.mu.Unlock()
		j.result = out
		j.status = out.Status
		s.active = false
	}()
	return r.ID, nil
}
func (s *Service) Status(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok {
		return "", fmt.Errorf("unknown request")
	}
	return j.status, nil
}
func (s *Service) Cancel(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok || j.status != "running" {
		return false
	}
	j.cancel()
	return true
}
func (s *Service) Result(id string) (broker.Outcome, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok || j.status == "running" {
		return broker.Outcome{}, fmt.Errorf("result unavailable")
	}
	data, err := json.Marshal(j.result)
	if err != nil {
		return broker.Outcome{}, err
	}
	var out broker.Outcome
	err = json.Unmarshal(data, &out)
	return out, err
}
