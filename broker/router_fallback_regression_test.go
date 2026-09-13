package broker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

type transientFastBackend struct {
	healthCalls int
}

func (b *transientFastBackend) ID() string { return "fast" }

func (b *transientFastBackend) Capabilities() []string {
	return []string{"code"}
}

func (b *transientFastBackend) Health(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b.healthCalls++
	if b.healthCalls == 1 {
		return fmt.Errorf("transient fast health failure")
	}
	return nil
}

func (b *transientFastBackend) EstimateCost(r inference.Request) inference.Cost {
	return inference.Cost{
		InputBytes:      len(r.Context),
		MaxOutputTokens: r.MaxOutputTokens,
		Estimated:       true,
	}
}

func (b *transientFastBackend) Invoke(context.Context, inference.Request) (inference.Result, error) {
	return inference.Result{}, fmt.Errorf("not used")
}

func (b *transientFastBackend) Cancel(string) bool { return false }

func TestRoutineFastDoesNotRetryAsFallback(t *testing.T) {
	reg := &Registry{}
	fast := &transientFastBackend{}
	if err := reg.Register(Entry{
		Backend: fast,
		Class:   "mock",
		Tier:    "fast",
	}); err != nil {
		t.Fatal(err)
	}

	req := inference.Request{
		ID:              "routine-fast-no-retry",
		ParentWorkID:    "qualification-regression",
		Role:            "code",
		Context:         "x",
		MaxContextBytes: 4096,
		MaxOutputTokens: 32,
		Deadline:        time.Now().Add(time.Second),
		Workspace:       "fixture",
		Capabilities:    []string{"code"},
	}

	policy := Policy{
		AllowedBackends: []string{"fast"},
		AllowFallback:   true,
		MaxRepairs:      1,
	}

	_, decision, err := Select(context.Background(), reg, req, policy, resources.Metrics{})
	if err == nil {
		t.Fatalf(
			"routine fast selection must fail after one authorized fast attempt; decision=%+v health_calls=%d",
			decision,
			fast.healthCalls,
		)
	}
	if fast.healthCalls != 1 {
		t.Fatalf("routine fast selection must perform exactly one fast health check; got %d", fast.healthCalls)
	}
}
