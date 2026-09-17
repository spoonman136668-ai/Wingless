package cognitive

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

const (
	EvidenceSchema         = "wingless.cognitive-runtime-run.v1"
	AcceptanceExternal     = "external_required"
	RoutingDecisionVersion = 1
)

type MemoryClass string

const (
	MemorySemantic   MemoryClass = "semantic"
	MemoryEpisodic   MemoryClass = "episodic"
	MemoryProcedural MemoryClass = "procedural"
	MemoryFailure    MemoryClass = "failure"
)

func (c MemoryClass) Valid() bool {
	switch c {
	case MemorySemantic, MemoryEpisodic, MemoryProcedural, MemoryFailure:
		return true
	default:
		return false
	}
}

type MemoryRecord struct {
	ID         string      `json:"id"`
	Version    int         `json:"version"`
	Class      MemoryClass `json:"class"`
	Key        string      `json:"key"`
	Content    string      `json:"content"`
	Provenance string      `json:"provenance"`
	CreatedAt  time.Time   `json:"created_at"`
}

type SkillKind string

const (
	SkillStaticText SkillKind = "static_text"
	SkillProcedure  SkillKind = "procedure"
	SkillTemplate   SkillKind = "prompt_template"
)

type ValidationState string

const (
	ValidationCandidate ValidationState = "candidate"
	ValidationValidated ValidationState = "validated"
)

type ValidationEvidence struct {
	ValidatorID    string    `json:"validator_id"`
	EvidenceSHA256 string    `json:"evidence_sha256"`
	ValidatedAt    time.Time `json:"validated_at"`
}

type Skill struct {
	ID                 string              `json:"id"`
	Version            int                 `json:"version"`
	Source             string              `json:"source"`
	TaskFamily         string              `json:"task_family"`
	Inputs             []string            `json:"inputs"`
	ExpectedOutputs    []string            `json:"expected_outputs"`
	Kind               SkillKind           `json:"kind"`
	Body               string              `json:"body"`
	Validation         ValidationState     `json:"validation"`
	ValidationEvidence *ValidationEvidence `json:"validation_evidence,omitempty"`
	UsageCount         uint64              `json:"usage_count"`
	SuccessCount       uint64              `json:"success_count"`
	FailureCount       uint64              `json:"failure_count"`
	CreatedAt          time.Time           `json:"created_at"`
	LastUseAt          *time.Time          `json:"last_use_at,omitempty"`
}

type MemoryStore interface {
	PutMemory(context.Context, MemoryRecord) (MemoryRecord, error)
	LookupMemory(context.Context, MemoryClass, string) (MemoryRecord, bool, error)
	DeleteMemory(context.Context, string) error
}

type SkillStore interface {
	PutSkillCandidate(context.Context, Skill) (Skill, error)
	PromoteSkill(context.Context, string, ValidationEvidence) (Skill, error)
	ListValidatedSkills(context.Context, string) ([]Skill, error)
	RecordSkillUse(context.Context, string, bool) error
}

type Route string

const (
	RouteMemory                  Route = "memory"
	RouteSkill                   Route = "skill"
	RouteInferenceSingle         Route = "single_pass_inference"
	RouteInferenceMulti          Route = "multi_pass_inference"
	RouteAuthorizedDeeperBackend Route = "authorized_deeper_backend"
	RouteToolProposal            Route = "tool_proposal"
)

type RoutingDecision struct {
	Version     int      `json:"version"`
	ID          string   `json:"id"`
	Policy      string   `json:"policy"`
	RequestID   string   `json:"request_id"`
	TaskFamily  string   `json:"task_family"`
	Route       Route    `json:"route"`
	Reason      string   `json:"reason"`
	MemoryIDs   []string `json:"memory_ids,omitempty"`
	SkillIDs    []string `json:"skill_ids,omitempty"`
	InputDigest string   `json:"input_digest"`
}

type Policy struct {
	Version              string `json:"version"`
	MaxPasses            int    `json:"max_passes"`
	EnableMemory         bool   `json:"enable_memory"`
	EnableSkills         bool   `json:"enable_skills"`
	AdditionalPassPrompt string `json:"additional_pass_prompt"`
}

func (p Policy) Validate() error {
	if strings.TrimSpace(p.Version) == "" {
		return errors.New("cognitive policy version is required")
	}
	if p.MaxPasses < 1 || p.MaxPasses > 8 {
		return errors.New("cognitive max passes must be between 1 and 8")
	}
	if len(p.AdditionalPassPrompt) > 1024 {
		return errors.New("additional-pass prompt exceeds hard bound")
	}
	return nil
}

type Evaluation struct {
	Sufficient bool   `json:"sufficient"`
	Reason     string `json:"reason"`
}

type Evaluator interface {
	Evaluate(pass int, result inference.Result) Evaluation
}

type EvaluatorFunc func(pass int, result inference.Result) Evaluation

func (f EvaluatorFunc) Evaluate(pass int, result inference.Result) Evaluation {
	return f(pass, result)
}

type CompletedNonEmptyEvaluator struct{}

func (CompletedNonEmptyEvaluator) Evaluate(_ int, result inference.Result) Evaluation {
	if result.Status == "completed" && strings.TrimSpace(result.Text) != "" {
		return Evaluation{Sufficient: true, Reason: "completed_non_empty"}
	}
	return Evaluation{Sufficient: false, Reason: "incomplete_or_empty"}
}

type RunRequest struct {
	RequestID      string            `json:"request_id"`
	TaskFamily     string            `json:"task_family"`
	MemoryClass    MemoryClass       `json:"memory_class,omitempty"`
	MemoryKey      string            `json:"memory_key,omitempty"`
	AllowMultiPass bool              `json:"allow_multi_pass"`
	Base           inference.Request `json:"base"`
}

type Candidate struct {
	Status     string `json:"status"`
	Text       string `json:"text"`
	Source     string `json:"source"`
	Acceptance string `json:"acceptance"`
}

type PassEvidence struct {
	Pass             int                 `json:"pass"`
	Reason           string              `json:"reason"`
	BackendID        string              `json:"backend_id"`
	ModelID          string              `json:"model_id"`
	InputTokens      *int                `json:"input_tokens"`
	OutputTokens     *int                `json:"output_tokens"`
	LatencyMS        int64               `json:"latency_ms"`
	Termination      string              `json:"termination_reason"`
	EvaluationReason string              `json:"evaluation_reason"`
	Telemetry        inference.Telemetry `json:"telemetry"`
}

type InferenceSessionEvidence struct {
	SessionID           string                    `json:"session_id"`
	BackendID           *string                   `json:"backend_id"`
	ModelID             *string                   `json:"model_id"`
	StartedAt           time.Time                 `json:"started_at"`
	FinishedAt          time.Time                 `json:"finished_at"`
	ModelCalls          int                       `json:"model_calls"`
	CognitivePasses     int                       `json:"cognitive_passes"`
	MemoryRetrievals    int                       `json:"memory_retrievals"`
	SkillReuses         int                       `json:"skill_reuses"`
	TerminationReason   string                    `json:"termination_reason"`
	ResourceTelemetry   *resources.ModelTelemetry `json:"resource_telemetry"`
}

type RunEvidence struct {
	Schema            string                   `json:"schema"`
	RequestID         string                   `json:"request_id"`
	Decision          RoutingDecision          `json:"decision"`
	Session           InferenceSessionEvidence `json:"session"`
	Passes            []PassEvidence           `json:"passes"`
	ModelCalls        int                      `json:"model_calls"`
	TotalInputTokens  *int                     `json:"total_input_tokens"`
	TotalOutputTokens *int                     `json:"total_output_tokens"`
	TotalLatencyMS    int64                    `json:"total_latency_ms"`
	TerminationReason string                   `json:"termination_reason"`
	Reuse             bool                     `json:"reuse"`
	Acceptance        string                   `json:"acceptance"`
	StartedAt         time.Time                `json:"started_at"`
	FinishedAt        time.Time                `json:"finished_at"`
}

type RunResult struct {
	Candidate Candidate   `json:"candidate"`
	Evidence  RunEvidence `json:"evidence"`
}
