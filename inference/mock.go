package inference

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Mock has immutable configuration after registration. Replies are test fixtures, not acceptance.
type Mock struct {
	Name      string
	Features  []string
	Reply     string
	Delay     time.Duration
	Unhealthy bool
	Failure   error
	mu        sync.Mutex
	active    map[string]context.CancelFunc
}

var _ InferenceBackend = (*Mock)(nil)

func (m *Mock) ID() string             { return m.Name }
func (m *Mock) Capabilities() []string { return append([]string(nil), m.Features...) }
func (m *Mock) Health(ctx context.Context) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	if m.Unhealthy {
		return fmt.Errorf("backend unavailable")
	}
	return nil
}
func (m *Mock) EstimateCost(r Request) Cost { return Cost{len(r.Context), r.MaxOutputTokens, true} }
func (m *Mock) Cancel(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.active[id]
	if ok {
		c()
	}
	return ok
}
func (m *Mock) Invoke(parent context.Context, r Request) (out Result, err error) {
	out = Result{BackendID: m.Name, ModelID: "fixture", Status: "failed"}
	start := time.Now()
	defer func() {
		out.LatencyMS = time.Since(start).Milliseconds()
		out.ErrorClass = Classify(err)
		if err != nil {
			out.Termination = out.ErrorClass
		}
	}()
	if err = r.Validate(); err != nil {
		return
	}
	if err = m.Health(parent); err != nil {
		return
	}
	ctx, cancel := context.WithDeadline(parent, r.Deadline)
	defer cancel()
	m.mu.Lock()
	if m.active == nil {
		m.active = map[string]context.CancelFunc{}
	}
	if _, ok := m.active[r.ID]; ok {
		m.mu.Unlock()
		err = fmt.Errorf("duplicate active request")
		return
	}
	m.active[r.ID] = cancel
	m.mu.Unlock()
	defer func() { m.mu.Lock(); delete(m.active, r.ID); m.mu.Unlock() }()
	timer := time.NewTimer(m.Delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case <-timer.C:
	}
	if err = ctx.Err(); err != nil {
		return
	}
	if m.Failure != nil {
		err = m.Failure
		return
	}
	if len(m.Reply) > r.MaxOutputTokens*4 {
		err = fmt.Errorf("fixture exceeds output bound")
		return
	}
	out.Text = m.Reply
	out.Status = "completed"
	out.Termination = "stop"
	return
}
