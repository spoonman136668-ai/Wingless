package unitary

// This structured-data runtime reuses NCI11. It is not a language model,
// source-authentication service, scientific evaluator or continual learner.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
)

const NativeSessionSource = "f248fd6dec93d04571f7cb4c3ccd721017812bfd"
const nativeSessionSchema = "wingless.native-learned-session.v1"

var ErrNativeSession = errors.New("WINGLESS_NATIVE_SESSION_INVALID")

type NativeSessionEdge struct {
	From int `json:"from"`
	To   int `json:"to"`
}
type NativeSessionCase struct {
	Edges     []NativeSessionEdge `json:"edges"`
	Questions [][]int             `json:"questions"`
	Truth     []bool              `json:"truth"`
}
type NativeSession struct {
	model       nci11WModel
	memory      nci11WMemory
	graphLoaded bool
}
type nativeSessionPayload struct {
	Schema      string       `json:"schema"`
	SourceBlob  string       `json:"source_blob"`
	Model       nci11WModel  `json:"model"`
	Memory      nci11WMemory `json:"memory"`
	GraphLoaded bool         `json:"graph_loaded"`
}
type nativeSessionEnvelope struct {
	Payload nativeSessionPayload `json:"payload"`
	SHA256  string               `json:"sha256"`
}

func nativeSessionID(x int) bool { return x > 0 && x <= 1000000 }
func nativeSessionEdges(xs []NativeSessionEdge) ([12]nci02WEdge, error) {
	var out [12]nci02WEdge
	if len(xs) != 12 {
		return out, ErrNativeSession
	}
	for i, e := range xs {
		if !nativeSessionID(e.From) || !nativeSessionID(e.To) {
			return out, ErrNativeSession
		}
		out[i] = nci02WEdge{From: e.From, To: e.To}
	}
	return out, nil
}

// TrainNativeSession uses the unchanged NCI11 from-scratch training routine.
// Caller-supplied labels must have legitimate provenance; this API grants none.
func TrainNativeSession(xs []NativeSessionCase, leaky bool) (*NativeSession, error) {
	if len(xs) < 1 || len(xs) > 16 {
		return nil, ErrNativeSession
	}
	cases := make([]nci11WCase, len(xs))
	for i, x := range xs {
		edges, err := nativeSessionEdges(x.Edges)
		if err != nil || len(x.Questions) != 2 || len(x.Truth) != 2 {
			return nil, ErrNativeSession
		}
		for j, q := range x.Questions {
			if len(q) != 2 || !nativeSessionID(q[0]) || !nativeSessionID(q[1]) {
				return nil, ErrNativeSession
			}
			cases[i].Questions[j] = [2]int{q[0], q[1]}
			cases[i].Truth[j] = x.Truth[j]
		}
		cases[i].Edges = edges
	}
	s := &NativeSession{model: nci11WLearn(cases, leaky)}
	if !s.valid() {
		return nil, ErrNativeSession
	}
	return s, nil
}
func (s *NativeSession) valid() bool {
	if s == nil || math.IsNaN(s.model.Gain) || math.IsInf(s.model.Gain, 0) ||
		s.model.Gain < .1 || s.model.Gain > 100 || math.IsNaN(s.model.Bias) ||
		math.IsInf(s.model.Bias, 0) || math.Abs(s.model.Bias) > 10000 {
		return false
	}
	for _, row := range s.model.Scale {
		for _, v := range row {
			if math.IsNaN(v) || math.IsInf(v, 0) || v < .1 || v > 3 {
				return false
			}
		}
	}
	for _, row := range s.memory {
		for _, v := range row {
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) || v < -1 || v > 1 {
				return false
			}
		}
	}
	return true
}
func (s *NativeSession) ObserveGraph(xs []NativeSessionEdge) error {
	edges, err := nativeSessionEdges(xs)
	if !s.valid() || err != nil {
		return ErrNativeSession
	}
	s.memory = nci11WBuild(edges, s.model)
	s.graphLoaded = true
	return nil
}
func (s *NativeSession) Predict(from, to int) (bool, float64, error) {
	if !s.valid() || !s.graphLoaded || !nativeSessionID(from) || !nativeSessionID(to) {
		return false, 0, ErrNativeSession
	}
	logit := nci11WLogit(s.memory, s.model, from, to)
	if math.IsNaN(logit) || math.IsInf(logit, 0) {
		return false, 0, ErrNativeSession
	}
	return logit >= 0, logit, nil
}
func (s *NativeSession) Checkpoint() ([]byte, error) {
	if !s.valid() {
		return nil, ErrNativeSession
	}
	p := nativeSessionPayload{nativeSessionSchema, NativeSessionSource, s.model, s.memory, s.graphLoaded}
	raw, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(raw)
	return json.Marshal(nativeSessionEnvelope{p, hex.EncodeToString(hash[:])})
}

// ValidateNativeSessionJSON rejects ambiguous duplicate keys, nulls, trailing
// values and excessive nesting. Checkpoint digests provide integrity, not trust.
func ValidateNativeSessionJSON(raw []byte) error {
	if len(raw) == 0 || len(raw) > 65536 {
		return ErrNativeSession
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var visit func(int) error
	visit = func(depth int) error {
		if depth > 16 {
			return ErrNativeSession
		}
		tok, err := d.Token()
		if err != nil || tok == nil {
			return ErrNativeSession
		}
		if delim, ok := tok.(json.Delim); ok {
			switch delim {
			case '{':
				seen := map[string]bool{}
				for d.More() {
					key, err := d.Token()
					k, ok := key.(string)
					if err != nil || !ok || seen[k] {
						return ErrNativeSession
					}
					seen[k] = true
					if visit(depth+1) != nil {
						return ErrNativeSession
					}
				}
				end, err := d.Token()
				if err != nil || end != json.Delim('}') {
					return ErrNativeSession
				}
			case '[':
				for d.More() {
					if visit(depth+1) != nil {
						return ErrNativeSession
					}
				}
				end, err := d.Token()
				if err != nil || end != json.Delim(']') {
					return ErrNativeSession
				}
			default:
				return ErrNativeSession
			}
		}
		return nil
	}
	if visit(0) != nil {
		return ErrNativeSession
	}
	if _, err := d.Token(); err != io.EOF {
		return ErrNativeSession
	}
	return nil
}
func nativeSessionObject(v any, keys ...string) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	if !ok || len(m) != len(keys) {
		return nil, false
	}
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return nil, false
		}
	}
	return m, true
}
func nativeSessionMatrix(v any, rows, cols int) bool {
	a, ok := v.([]any)
	if !ok || len(a) != rows {
		return false
	}
	for _, r := range a {
		b, ok := r.([]any)
		if !ok || len(b) != cols {
			return false
		}
		for _, x := range b {
			if _, ok := x.(float64); !ok {
				return false
			}
		}
	}
	return true
}
func RestoreNativeSession(raw []byte) (*NativeSession, error) {
	if ValidateNativeSessionJSON(raw) != nil {
		return nil, ErrNativeSession
	}
	var shape any
	if json.Unmarshal(raw, &shape) != nil {
		return nil, ErrNativeSession
	}
	top, ok := nativeSessionObject(shape, "payload", "sha256")
	if !ok {
		return nil, ErrNativeSession
	}
	p, ok := nativeSessionObject(top["payload"], "schema", "source_blob", "model", "memory", "graph_loaded")
	if !ok {
		return nil, ErrNativeSession
	}
	m, ok := nativeSessionObject(p["model"], "Scale", "Gain", "Bias", "Leaky")
	if !ok || !nativeSessionMatrix(m["Scale"], 4, 12) || !nativeSessionMatrix(p["memory"], 12, 12) {
		return nil, ErrNativeSession
	}
	var e nativeSessionEnvelope
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&e) != nil || d.Decode(new(any)) != io.EOF {
		return nil, ErrNativeSession
	}
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return nil, ErrNativeSession
	}
	hash := sha256.Sum256(payload)
	if e.Payload.Schema != nativeSessionSchema || e.Payload.SourceBlob != NativeSessionSource ||
		e.SHA256 != hex.EncodeToString(hash[:]) {
		return nil, ErrNativeSession
	}
	s := &NativeSession{e.Payload.Model, e.Payload.Memory, e.Payload.GraphLoaded}
	if !s.valid() {
		return nil, ErrNativeSession
	}
	return s, nil
}
