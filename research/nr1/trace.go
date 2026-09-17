package nr1

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
)

const TraceSchema = "wingless.nr1.router-trace.v2"

var traceID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$`)

// TraceEvent is one passive observation of MoE routing for one token at one
// routed layer. It is research telemetry only and carries no execution or
// acceptance authority.
type TraceEvent struct {
	Schema     string    `json:"schema"`
	SessionID  string    `json:"session_id"`
	WorkloadID string    `json:"workload_id"`
	TaskFamily string    `json:"task_family"`
	Phase      string    `json:"phase"`
	TokenIndex int64     `json:"token_index"`
	Layer      int       `json:"layer"`
	Experts    []int     `json:"experts"`
	Scores     []float64 `json:"scores,omitempty"`
}

func (e TraceEvent) Validate() error {
	if e.Schema != TraceSchema {
		return errors.New("NR1_TRACE_SCHEMA_INVALID")
	}
	if !traceID.MatchString(e.SessionID) || !traceID.MatchString(e.WorkloadID) || !traceID.MatchString(e.TaskFamily) {
		return errors.New("NR1_TRACE_ID_INVALID")
	}
	if e.Phase != "prefill" && e.Phase != "decode" {
		return errors.New("NR1_TRACE_PHASE_INVALID")
	}
	if e.TokenIndex < 0 || e.Layer < 0 || e.Layer > 255 {
		return errors.New("NR1_TRACE_POSITION_INVALID")
	}
	if len(e.Experts) < 1 || len(e.Experts) > 128 {
		return errors.New("NR1_TRACE_EXPERT_COUNT_INVALID")
	}
	if len(e.Scores) != 0 && len(e.Scores) != len(e.Experts) {
		return errors.New("NR1_TRACE_SCORE_COUNT_INVALID")
	}
	seen := make(map[int]struct{}, len(e.Experts))
	for i, expert := range e.Experts {
		if expert < 0 || expert > 1023 {
			return errors.New("NR1_TRACE_EXPERT_ID_INVALID")
		}
		if _, ok := seen[expert]; ok {
			return errors.New("NR1_TRACE_DUPLICATE_EXPERT")
		}
		seen[expert] = struct{}{}
		if len(e.Scores) != 0 && (math.IsNaN(e.Scores[i]) || math.IsInf(e.Scores[i], 0)) {
			return errors.New("NR1_TRACE_SCORE_INVALID")
		}
	}
	return nil
}

// ReadJSONL strictly parses a bounded JSONL trace. Unknown fields are rejected
// so experimental producer drift cannot silently contaminate residency results.
func ReadJSONL(r io.Reader, maxEvents int) ([]TraceEvent, error) {
	if r == nil || maxEvents < 1 || maxEvents > 10_000_000 {
		return nil, errors.New("NR1_TRACE_READ_BOUNDS_INVALID")
	}
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 64*1024), 1024*1024)
	out := make([]TraceEvent, 0, min(maxEvents, 4096))
	line := 0
	for s.Scan() {
		line++
		raw := bytes.TrimSpace(s.Bytes())
		if len(raw) == 0 {
			continue
		}
		if len(out) >= maxEvents {
			return nil, errors.New("NR1_TRACE_EVENT_LIMIT_EXCEEDED")
		}
		var e TraceEvent
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("NR1_TRACE_JSON_INVALID line=%d: %w", line, err)
		}
		if err := dec.Decode(new(any)); err != io.EOF {
			return nil, fmt.Errorf("NR1_TRACE_JSON_TRAILING line=%d", line)
		}
		if err := e.Validate(); err != nil {
			return nil, fmt.Errorf("%w line=%d", err, line)
		}
		out = append(out, e)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, errors.New("NR1_TRACE_EMPTY")
	}
	return out, nil
}
