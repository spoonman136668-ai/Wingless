package ckbplane

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	TransportRequestSchema = "ckb-plane.wingless-transport-request.v1"
	CandidateEnvelopeSchema = "ckb-plane.wingless-candidate-envelope.v1"
	TransportPath = "/v1/wingless/candidate"
	maxTransportRequestBytes = 2 << 20
	maxCandidateEnvelopeBytes = 7 << 20
	maxCandidateTextBytes = 1 << 20
	maxEditOperations = 64
	maxEditFileBytes = 1 << 20
	maxEditAggregateBytes = 4 << 20
)

var sha64 = regexp.MustCompile(`^[0-9a-f]{64}$`)

type TransportResources struct {
	MinRAMBytes   uint64  `json:"min_ram_bytes"`
	MinVRAMBytes  uint64  `json:"min_vram_bytes"`
	MinDiskBytes  uint64  `json:"min_disk_bytes"`
	MaxCPUPercent float64 `json:"max_cpu_percent"`
}

type TransportPolicy struct {
	ReadPaths       []string           `json:"read_paths"`
	AllowedBackends []string           `json:"allowed_backends"`
	ExplicitBackend string             `json:"explicit_backend,omitempty"`
	AllowDeep       bool               `json:"allow_deep"`
	AllowFallback   bool               `json:"allow_fallback"`
	MaxContextBytes int                `json:"max_context_bytes"`
	MaxOutputTokens int                `json:"max_output_tokens"`
	Resources       TransportResources `json:"resources"`
}

type TransportRequest struct {
	Schema               string          `json:"schema"`
	SourceCommit         string          `json:"source_commit"`
	RequestID            string          `json:"request_id"`
	WorkOrderID          string          `json:"work_order_id"`
	Attempt              int             `json:"attempt"`
	Workspace            string          `json:"workspace"`
	BaselineSHA          string          `json:"baseline_sha"`
	WorkOrderDigest      string          `json:"work_order_digest"`
	Title                string          `json:"title"`
	Instructions         []string        `json:"instructions"`
	ForbiddenActions     []string        `json:"forbidden_actions"`
	Preserve             []string        `json:"preserve"`
	ForbiddenPaths       []string        `json:"forbidden_paths"`
	Deadline             time.Time       `json:"deadline"`
	Policy               TransportPolicy `json:"policy"`
	WorkOrderRetryBudget int             `json:"work_order_retry_budget"`
}

type EditOperation struct {
	Type                 string `json:"type"`
	Path                 string `json:"path"`
	ExpectedBeforeSHA256 string `json:"expected_before_sha256,omitempty"`
	Content              string `json:"content,omitempty"`
}

type EditEnvelope struct {
	SourceCommit    string          `json:"source_commit"`
	RequestID       string          `json:"request_id"`
	WorkOrderID     string          `json:"work_order_id"`
	Attempt         int             `json:"attempt"`
	Workspace       string          `json:"workspace"`
	BaselineSHA     string          `json:"baseline_sha"`
	WorkOrderDigest string          `json:"work_order_digest"`
	Operations      []EditOperation `json:"operations"`
}

type CandidateEnvelope struct {
	Schema          string        `json:"schema"`
	SourceCommit    string        `json:"source_commit"`
	RequestID       string        `json:"request_id"`
	WorkOrderID     string        `json:"work_order_id"`
	Attempt         int           `json:"attempt"`
	Workspace       string        `json:"workspace"`
	BaselineSHA     string        `json:"baseline_sha"`
	WorkOrderDigest string        `json:"work_order_digest"`
	Status          string        `json:"status"`
	Acceptance      string        `json:"acceptance"`
	Text            string        `json:"text,omitempty"`
	BackendID       string        `json:"backend_id"`
	ModelID         string        `json:"model_id"`
	ErrorClass      string        `json:"error_class,omitempty"`
	Edits           *EditEnvelope `json:"edits,omitempty"`
}

type Intelligence interface {
	Run(context.Context, TransportRequest) (CandidateEnvelope, error)
}

type Listener struct {
	session string
	engine  Intelligence
}

func NewListener(session string, engine Intelligence) (*Listener, error) {
	if !validSession(session) {
		return nil, errors.New("WINGLESS_LISTENER_SESSION_INVALID")
	}
	if engine == nil {
		return nil, errors.New("WINGLESS_LISTENER_ENGINE_REQUIRED")
	}
	return &Listener{session: session, engine: engine}, nil
}

func validSession(session string) bool {
	if len(session) != 64 {
		return false
	}
	for _, r := range session {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return false
		}
	}
	return true
}

func (l *Listener) Handler() http.Handler {
	return http.HandlerFunc(l.ServeHTTP)
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if l == nil || l.engine == nil {
		writeListenerError(w, http.StatusServiceUnavailable, "WINGLESS_LISTENER_UNAVAILABLE")
		return
	}
	if r.URL.Path != TransportPath {
		writeListenerError(w, http.StatusNotFound, "WINGLESS_LISTENER_PATH_INVALID")
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeListenerError(w, http.StatusMethodNotAllowed, "WINGLESS_LISTENER_METHOD_INVALID")
		return
	}
	gotSession := r.Header.Get("X-CKB-Wingless-Session")
	if len(gotSession) != len(l.session) || subtle.ConstantTimeCompare([]byte(gotSession), []byte(l.session)) != 1 {
		writeListenerError(w, http.StatusUnauthorized, "WINGLESS_LISTENER_SESSION_MISMATCH")
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeListenerError(w, http.StatusUnsupportedMediaType, "WINGLESS_LISTENER_CONTENT_TYPE")
		return
	}
	if r.ContentLength > maxTransportRequestBytes {
		writeListenerError(w, http.StatusRequestEntityTooLarge, "WINGLESS_LISTENER_REQUEST_LIMIT")
		return
	}

	// The request body is fully consumed and closed before intelligence starts.
	// This is required so HTTP request-context cancellation can propagate while
	// inference is running instead of being masked by an unread HTTP/1 body.
	raw, readErr := io.ReadAll(io.LimitReader(r.Body, maxTransportRequestBytes+1))
	closeErr := r.Body.Close()
	if readErr != nil || closeErr != nil {
		writeListenerError(w, http.StatusBadRequest, "WINGLESS_LISTENER_REQUEST_INVALID")
		return
	}
	if len(raw) == 0 || len(raw) > maxTransportRequestBytes || !utf8.Valid(raw) {
		writeListenerError(w, http.StatusRequestEntityTooLarge, "WINGLESS_LISTENER_REQUEST_LIMIT")
		return
	}

	req, err := decodeTransportRequest(raw)
	if err != nil || validateTransportRequest(req) != nil {
		writeListenerError(w, http.StatusBadRequest, "WINGLESS_LISTENER_REQUEST_INVALID")
		return
	}

	ctx, cancel := context.WithDeadline(r.Context(), req.Deadline)
	defer cancel()
	candidate, err := l.engine.Run(ctx, req) // exactly one intelligence invocation
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		writeListenerError(w, http.StatusServiceUnavailable, "WINGLESS_LISTENER_ENGINE_FAILURE")
		return
	}
	if ctx.Err() != nil {
		return
	}
	if err := validateCandidateEnvelope(candidate, req); err != nil {
		writeListenerError(w, http.StatusBadGateway, "WINGLESS_LISTENER_CANDIDATE_INVALID")
		return
	}
	payload, err := json.Marshal(candidate)
	if err != nil || len(payload) > maxCandidateEnvelopeBytes {
		writeListenerError(w, http.StatusBadGateway, "WINGLESS_LISTENER_CANDIDATE_INVALID")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CKB-Wingless-Session", l.session)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(payload, '\n'))
}

func writeListenerError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": code})
}

func decodeTransportRequest(data []byte) (TransportRequest, error) {
	var req TransportRequest
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return req, err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&req); err != nil {
		return req, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return TransportRequest{}, errors.New("trailing JSON")
	}
	return req, nil
}

func validateTransportRequest(req TransportRequest) error {
	if req.Schema != TransportRequestSchema || req.SourceCommit != FrozenPlaneSourceCommit {
		return errors.New("transport schema or source mismatch")
	}
	if !safeID.MatchString(req.WorkOrderID) || req.RequestID != fmt.Sprintf("%s-wingless-%d", req.WorkOrderID, req.Attempt) || req.Attempt < 1 {
		return errors.New("transport identity mismatch")
	}
	if !sha40.MatchString(req.BaselineSHA) || !sha64.MatchString(req.WorkOrderDigest) || strings.TrimSpace(req.Workspace) == "" || strings.TrimSpace(req.Title) == "" {
		return errors.New("transport identity invalid")
	}
	if req.WorkOrderRetryBudget != 0 || req.Deadline.IsZero() || !req.Deadline.After(time.Now()) || time.Until(req.Deadline) > 15*time.Minute {
		return errors.New("transport authority or deadline invalid")
	}
	p := req.Policy
	if p.MaxContextBytes < 1 || p.MaxContextBytes > 1<<20 || p.MaxOutputTokens < 1 || p.MaxOutputTokens > 32768 || len(p.AllowedBackends) == 0 || len(p.AllowedBackends) > 32 || len(p.ReadPaths) > 256 {
		return errors.New("transport policy invalid")
	}
	seen := map[string]bool{}
	for _, id := range p.AllowedBackends {
		if !safeID.MatchString(id) || seen[id] {
			return errors.New("transport backend policy invalid")
		}
		seen[id] = true
	}
	if p.ExplicitBackend != "" && !seen[p.ExplicitBackend] {
		return errors.New("transport explicit backend invalid")
	}
	for _, path := range p.ReadPaths {
		if strings.TrimSpace(path) == "" || strings.ContainsAny(path, "*?[]\x00\r\n") {
			return errors.New("transport read path invalid")
		}
	}
	return nil
}

func validateCandidateEnvelope(c CandidateEnvelope, req TransportRequest) error {
	if c.Schema != CandidateEnvelopeSchema || c.SourceCommit != req.SourceCommit || c.RequestID != req.RequestID || c.WorkOrderID != req.WorkOrderID || c.Attempt != req.Attempt || c.Workspace != req.Workspace || c.BaselineSHA != req.BaselineSHA || c.WorkOrderDigest != req.WorkOrderDigest {
		return errors.New("candidate identity mismatch")
	}
	if c.Acceptance != "external_required" || !utf8.ValidString(c.Text) || len([]byte(c.Text)) > maxCandidateTextBytes || len(c.BackendID) > 120 || len(c.ModelID) > 120 || len(c.ErrorClass) > 120 {
		return errors.New("candidate authority or size invalid")
	}
	switch c.Status {
	case "result_ready", "verified_candidate":
		if c.Edits == nil || c.ModelID == "" || c.ErrorClass != "" {
			return errors.New("candidate result incomplete")
		}
		allowed := false
		for _, id := range req.Policy.AllowedBackends {
			if id == c.BackendID {
				allowed = true
				break
			}
		}
		if !allowed || (req.Policy.ExplicitBackend != "" && c.BackendID != req.Policy.ExplicitBackend && !req.Policy.AllowFallback) {
			return errors.New("candidate backend policy mismatch")
		}
	case "operator_blocked":
		if c.Edits != nil || c.ErrorClass == "" {
			return errors.New("blocked candidate invalid")
		}
	default:
		return errors.New("unknown candidate terminal state")
	}
	if c.Edits != nil {
		if err := validateEditEnvelope(*c.Edits, req); err != nil {
			return err
		}
	}
	return nil
}

func validateEditEnvelope(e EditEnvelope, req TransportRequest) error {
	if e.SourceCommit != req.SourceCommit || e.RequestID != req.RequestID || e.WorkOrderID != req.WorkOrderID || e.Attempt != req.Attempt || e.Workspace != req.Workspace || e.BaselineSHA != req.BaselineSHA || e.WorkOrderDigest != req.WorkOrderDigest {
		return errors.New("edit identity mismatch")
	}
	if len(e.Operations) == 0 || len(e.Operations) > maxEditOperations {
		return errors.New("edit operation count invalid")
	}
	seen := map[string]bool{}
	aggregate := 0
	for _, op := range e.Operations {
		if strings.TrimSpace(op.Path) == "" || strings.ContainsAny(op.Path, "*?[]\x00\r\n") || seen[op.Path] {
			return errors.New("edit path invalid")
		}
		seen[op.Path] = true
		switch op.Type {
		case "write_text":
			if !utf8.ValidString(op.Content) || len([]byte(op.Content)) > maxEditFileBytes {
				return errors.New("edit content invalid")
			}
			aggregate += len([]byte(op.Content))
		case "delete_file":
			if op.Content != "" || !sha64.MatchString(op.ExpectedBeforeSHA256) {
				return errors.New("delete edit invalid")
			}
		default:
			return errors.New("unsupported edit operation")
		}
		if op.ExpectedBeforeSHA256 != "" && !sha64.MatchString(op.ExpectedBeforeSHA256) {
			return errors.New("edit preimage invalid")
		}
	}
	if aggregate > maxEditAggregateBytes {
		return errors.New("edit aggregate limit")
	}
	return nil
}

func rejectDuplicateJSONKeys(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	var walk func() error
	walk = func() error {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		delim, ok := tok.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				key, err := d.Token()
				if err != nil {
					return err
				}
				name, ok := key.(string)
				if !ok || seen[name] {
					return errors.New("duplicate JSON field")
				}
				seen[name] = true
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = d.Token()
			return err
		case '[':
			for d.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = d.Token()
			return err
		default:
			return errors.New("invalid JSON container")
		}
	}
	if err := walk(); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return errors.New("trailing JSON")
	}
	return nil
}
