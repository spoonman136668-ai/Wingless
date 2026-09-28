package reasoner

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	RequestSchema        = "wingless.experiment-reasoning-request.v1"
	ResultSchema         = "wingless.experiment-reasoning-result.v1"
	DefaultModel         = "nvidia/nemotron-3-ultra-550b-a55b"
	DefaultEndpoint      = "https://openrouter.ai/api/v1/chat/completions"
	MaxContextBytes      = 1 << 20
	MaxCompletionTokens = 16384
)

const systemPrompt = `You are an advisory scientific reasoner. You may propose hypotheses, controls, discriminating experiments, predicted observations, confounds, and interpretations. You do not authorize execution, acceptance, promotion, deployment, threshold changes, retries, or scope expansion. Treat the supplied preregistration and frozen scientific constraints as immutable. Never reinterpret a completed experiment by changing its criteria. Distinguish observations from hypotheses and uncertainty. Return only the requested scientific analysis.`

type Role string

const (
	RolePlan      Role = "plan-next-experiment"
	RoleInterpret Role = "interpret-sealed-result"
	RoleCritique  Role = "critique-proposed-experiment"
)

type JSONSchemaConstraint struct {
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema"`
}

type Request struct {
	Schema                      string                `json:"schema"`
	RequestID                   string                `json:"request_id"`
	Project                     string                `json:"project"`
	ExperimentID                string                `json:"experiment_id"`
	Role                        Role                  `json:"role"`
	FrontierSHA256              string                `json:"frontier_sha256"`
	QualificationContractSHA256 string                `json:"qualification_contract_sha256"`
	Context                     string                `json:"context"`
	MaxOutputTokens             int                   `json:"max_output_tokens"`
	TimeoutSeconds              int                   `json:"timeout_seconds"`
	ResponseJSONSchema          *JSONSchemaConstraint `json:"response_json_schema,omitempty"`
}

func (r Request) Validate() error {
	if r.ResponseJSONSchema != nil {
		if r.ResponseJSONSchema.Name == "" || len(r.ResponseJSONSchema.Schema) == 0 {
			return errors.New("response JSON schema invalid")
		}
	}
	if r.Schema != RequestSchema {
		return errors.New("reasoner request schema mismatch")
	}
	if r.RequestID == "" || r.Project == "" || r.ExperimentID == "" {
		return errors.New("reasoner request identity missing")
	}
	if r.ResponseJSONSchema != nil && (r.ResponseJSONSchema.Name == "" || len(r.ResponseJSONSchema.Schema) == 0) {
		return errors.New("response JSON schema invalid")
	}
	switch r.Role {
	case RolePlan, RoleInterpret, RoleCritique:
	default:
		return errors.New("reasoner role invalid")
	}
	if !isSHA256(r.FrontierSHA256) || !isSHA256(r.QualificationContractSHA256) {
		return errors.New("reasoner qualification identity invalid")
	}
	if !utf8.ValidString(r.Context) || len(r.Context) == 0 || len(r.Context) > MaxContextBytes {
		return errors.New("reasoner context invalid or unbounded")
	}
	if r.MaxOutputTokens <= 0 || r.MaxOutputTokens > MaxCompletionTokens {
		return errors.New("reasoner output budget invalid")
	}
	if r.TimeoutSeconds <= 0 || r.TimeoutSeconds > 900 {
		return errors.New("reasoner timeout invalid")
	}
	return nil
}

type Config struct {
	Model                  string  `json:"model"`
	Endpoint               string  `json:"endpoint"`
	Seed                   int64   `json:"seed"`
	ReasoningEffort        string  `json:"reasoning_effort"`
	Temperature            float64 `json:"temperature"`
	TopP                   float64 `json:"top_p"`
	DataCollection         string  `json:"data_collection"`
	ZeroDataRetention      bool    `json:"zero_data_retention"`
	AllowProviderFallbacks bool    `json:"allow_provider_fallbacks"`
}

func DefaultConfig() Config {
	return Config{
		Model:                  DefaultModel,
		Endpoint:               DefaultEndpoint,
		Seed:                   0,
		ReasoningEffort:        "high",
		Temperature:            0,
		TopP:                   0.95,
		DataCollection:         "deny",
		ZeroDataRetention:      true,
		AllowProviderFallbacks: false,
	}
}

func (c Config) Validate() error {
	if c.Model == "" || strings.ContainsAny(c.Model, " 	
") {
		return errors.New("reasoner model invalid")
	}
	if c.Endpoint != DefaultEndpoint {
		return errors.New("reasoner endpoint must be the pinned OpenRouter chat endpoint")
	}
	switch c.ReasoningEffort {
	case "xhigh", "high", "medium", "low", "minimal":
	default:
		return errors.New("reasoning effort invalid")
	}
	if c.Temperature < 0 || c.Temperature > 2 || c.TopP <= 0 || c.TopP > 1 {
		return errors.New("sampling configuration invalid")
	}
	if c.DataCollection != "deny" {
		return errors.New("data collection must be denied")
	}
	if !c.ZeroDataRetention {
		return errors.New("zero data retention must remain enabled")
	}
	return nil
}

type Usage struct {
	PromptTokens     *int `json:"prompt_tokens,omitempty"`
	CompletionTokens *int `json:"completion_tokens,omitempty"`
	ReasoningTokens  *int `json:"reasoning_tokens,omitempty"`
}

type Result struct {
	Schema                      string `json:"schema"`
	RequestID                   string `json:"request_id"`
	Project                     string `json:"project"`
	ExperimentID                string `json:"experiment_id"`
	Role                        Role   `json:"role"`
	RequestedModel              string `json:"requested_model"`
	ReturnedModel               string `json:"returned_model"`
	Provider                    string `json:"provider,omitempty"`
	Seed                        int64  `json:"seed"`
	ReasoningEffort             string `json:"reasoning_effort"`
	FrontierSHA256              string `json:"frontier_sha256"`
	QualificationContractSHA256 string `json:"qualification_contract_sha256"`
	RequestSHA256               string `json:"request_sha256"`
	ConfigSHA256                string `json:"config_sha256"`
	ResponseSHA256              string `json:"response_sha256"`
	QualificationIdentitySHA256 string `json:"qualification_identity_sha256"`
	Text                        string `json:"text"`
	Usage                       Usage  `json:"usage"`
	LatencyMS                   int64  `json:"latency_ms"`
	CompletedAtUTC              string `json:"completed_at_utc"`
}

type Client struct {
	token    string
	config   Config
	endpoint string
	http     *http.Client
}

func NewClient(token string, config Config) (*Client, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("OpenRouter API key missing")
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy:                 nil,
		ResponseHeaderTimeout: 60 * time.Second,
		MaxResponseHeaderBytes: 16384,
	}
	return &Client{
		token:    token,
		config:   config,
		endpoint: config.Endpoint,
		http: &http.Client{
			Transport: transport,
			Timeout:   15 * time.Minute,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errors.New("reasoner redirect denied")
			},
		},
	}, nil
}

func (c *Client) Invoke(parent context.Context, r Request) (Result, error) {
	var out Result
	if err := r.Validate(); err != nil {
		return out, err
	}
	if err := c.config.Validate(); err != nil {
		return out, err
	}

	reqBytes, err := json.Marshal(r)
	if err != nil {
		return out, err
	}
	configBytes, err := json.Marshal(c.config)
	if err != nil {
		return out, err
	}
	reqHash := sha256hex(reqBytes)
	configHash := sha256hex(configBytes)

	payload := struct {
		Model       string              `json:"model"`
		Messages    []map[string]string `json:"messages"`
		MaxTokens   int                 `json:"max_tokens"`
		Temperature float64             `json:"temperature"`
		TopP        float64             `json:"top_p"`
		Seed        int64               `json:"seed"`
		Reasoning   map[string]any      `json:"reasoning"`
		Provider       map[string]any      `json:"provider"`
		ResponseFormat map[string]any      `json:"response_format,omitempty"`
		Stream         bool                `json:"stream"`
	}{
		Model: c.config.Model,
		Messages: []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": r.Context},
		},
		MaxTokens:   r.MaxOutputTokens,
		Temperature: c.config.Temperature,
		TopP:        c.config.TopP,
		Seed:        c.config.Seed,
		Reasoning: map[string]any{
			"effort":  c.config.ReasoningEffort,
			"exclude": true,
		},
		Provider: map[string]any{
			"allow_fallbacks":  c.config.AllowProviderFallbacks,
			"data_collection":  c.config.DataCollection,
			"zdr":              c.config.ZeroDataRetention,
			"require_parameters": true,
		},
		Stream: false,
	}
	if r.ResponseJSONSchema != nil {
		payload.ResponseFormat = map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name":   r.ResponseJSONSchema.Name,
				"strict": true,
				"schema": r.ResponseJSONSchema.Schema,
			},
		}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return out, err
	}

	timeout := time.Duration(r.TimeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return out, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.token)
	httpReq.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.http.Do(httpReq)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("OpenRouter inference HTTP %d", resp.StatusCode)
	}

	limit := int64(r.MaxOutputTokens*24 + 131072)
	raw, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return out, err
	}
	if int64(len(raw)) > limit || !utf8.Valid(raw) {
		return out, errors.New("OpenRouter response invalid or oversized")
	}

	var wire struct {
		Model    string `json:"model"`
		Provider string `json:"provider"`
		Choices  []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens            *int `json:"prompt_tokens"`
			CompletionTokens        *int `json:"completion_tokens"`
			CompletionTokensDetails *struct {
				ReasoningTokens *int `json:"reasoning_tokens"`
			} `json:"completion_tokens_details"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		return out, err
	}
	if wire.Model != c.config.Model {
		return out, fmt.Errorf("model identity drift: requested=%s returned=%s", c.config.Model, wire.Model)
	}
	if len(wire.Choices) != 1 || wire.Choices[0].FinishReason != "stop" {
		return out, errors.New("incomplete or ambiguous reasoner generation")
	}
	text := wire.Choices[0].Message.Content
	if text == "" || !utf8.ValidString(text) || len(text) > r.MaxOutputTokens*24 {
		return out, errors.New("reasoner output invalid or oversized")
	}

	var reasoningTokens *int
	if wire.Usage.CompletionTokensDetails != nil {
		reasoningTokens = wire.Usage.CompletionTokensDetails.ReasoningTokens
	}
	responseHash := sha256hex([]byte(text))
	identityBytes, err := json.Marshal(struct {
		RequestSHA256               string `json:"request_sha256"`
		ConfigSHA256                string `json:"config_sha256"`
		ReturnedModel               string `json:"returned_model"`
		Provider                    string `json:"provider"`
		FrontierSHA256              string `json:"frontier_sha256"`
		QualificationContractSHA256 string `json:"qualification_contract_sha256"`
	}{
		RequestSHA256:               reqHash,
		ConfigSHA256:                configHash,
		ReturnedModel:               wire.Model,
		Provider:                    wire.Provider,
		FrontierSHA256:              r.FrontierSHA256,
		QualificationContractSHA256: r.QualificationContractSHA256,
	})
	if err != nil {
		return out, err
	}

	out = Result{
		Schema:                      ResultSchema,
		RequestID:                   r.RequestID,
		Project:                     r.Project,
		ExperimentID:                r.ExperimentID,
		Role:                        r.Role,
		RequestedModel:              c.config.Model,
		ReturnedModel:               wire.Model,
		Provider:                    wire.Provider,
		Seed:                        c.config.Seed,
		ReasoningEffort:             c.config.ReasoningEffort,
		FrontierSHA256:              r.FrontierSHA256,
		QualificationContractSHA256: r.QualificationContractSHA256,
		RequestSHA256:               reqHash,
		ConfigSHA256:                configHash,
		ResponseSHA256:              responseHash,
		QualificationIdentitySHA256: sha256hex(identityBytes),
		Text:                        text,
		Usage: Usage{
			PromptTokens:     wire.Usage.PromptTokens,
			CompletionTokens: wire.Usage.CompletionTokens,
			ReasoningTokens:  reasoningTokens,
		},
		LatencyMS:      latency,
		CompletedAtUTC: time.Now().UTC().Format(time.RFC3339Nano),
	}
	return out, nil
}

func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func isSHA256(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}
