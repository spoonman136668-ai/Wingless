package inference

import (
	"context"
	"errors"
	"github.com/spoonman136668-ai/Wingless/resources"
	"time"
	"unicode/utf8"
)

type Request struct {
	ID              string           `json:"request_id"`
	ParentWorkID    string           `json:"parent_work_id"`
	Role            string           `json:"role"`
	Context         string           `json:"context"`
	MaxContextBytes int              `json:"max_context_bytes"`
	MaxOutputTokens int              `json:"max_output_tokens"`
	Deadline        time.Time        `json:"deadline"`
	Workspace       string           `json:"workspace"`
	Capabilities    []string         `json:"capabilities"`
	Resources       resources.Policy `json:"resources"`
}

func (r Request) Validate() error {
	if !utf8.ValidString(r.Context) || r.ID == "" || r.ParentWorkID == "" || r.Workspace == "" || r.Role == "" || len(r.Capabilities) == 0 || r.MaxContextBytes <= 0 || r.MaxContextBytes > 1048576 || len(r.Context) > r.MaxContextBytes || r.MaxOutputTokens <= 0 || r.MaxOutputTokens > 32768 || r.Deadline.IsZero() || time.Until(r.Deadline) <= 0 || time.Until(r.Deadline) > 15*time.Minute {
		return errors.New("invalid or unbounded request")
	}
	return nil
}

type Usage struct {
	PromptTokens *int `json:"prompt_tokens"`
	OutputTokens *int `json:"output_tokens"`
}
type Telemetry struct {
	TimeToFirstTokenMS *int64             `json:"time_to_first_token_ms"`
	GenerationMS       *int64             `json:"generation_ms"`
	ModelLoadMS        *int64             `json:"model_load_ms"`
	Before             *resources.Metrics `json:"before"`
	Peak               *resources.Metrics `json:"peak"`
	After              *resources.Metrics `json:"after"`
}
type Result struct {
	Telemetry   Telemetry         `json:"telemetry"`
	BackendID   string            `json:"backend_id"`
	ModelID     string            `json:"model_id"`
	Status      string            `json:"status"`
	Text        string            `json:"text"`
	Usage       Usage             `json:"usage"`
	LatencyMS   int64             `json:"latency_ms"`
	Termination string            `json:"termination_reason"`
	ErrorClass  string            `json:"error_class,omitempty"`
	Resources   resources.Metrics `json:"resources"`
}
type Cost struct {
	InputBytes      int
	MaxOutputTokens int
	Estimated       bool
}
type InferenceBackend interface {
	ID() string
	Capabilities() []string
	Health(context.Context) error
	EstimateCost(Request) Cost
	Invoke(context.Context, Request) (Result, error)
	Cancel(string) bool
}

func Classify(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "worker_timeout"
	}
	if errors.Is(err, context.Canceled) {
		return "canceled"
	}
	if err != nil {
		return "backend_failure"
	}
	return ""
}
