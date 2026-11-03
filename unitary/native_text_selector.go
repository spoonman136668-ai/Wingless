package unitary

// Native whole-response selection is a bounded text input/output component.
// It is not novel-string generation, a semantic parser, or a language proof.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"sort"
	"unicode/utf8"
)

const nativeTextSchema = "wingless.native-whole-response-selector.v1"
const nativeTextEncoding = "utf8-thirteen-folded-percepts-twelve-edges.v1"

type NativeTextExample struct {
	Input    string `json:"input"`
	Response string `json:"response"`
}
type NativeTextPrediction struct {
	Response              string    `json:"response"`
	Responses             []string  `json:"responses"`
	NativeScores          []float64 `json:"native_scores"`
	OutputMode            string    `json:"output_mode"`
	LanguageProof         bool      `json:"language_proof"`
	NovelStringGeneration bool      `json:"novel_string_generation"`
}
type NativeTextSelector struct {
	responses   []string
	models      []*NativeSession
	trainingSHA string
}
type nativeTextPayload struct {
	Schema      string            `json:"schema"`
	Encoding    string            `json:"encoding"`
	SourceBlob  string            `json:"source_blob"`
	Responses   []string          `json:"responses"`
	TrainingSHA string            `json:"training_sha256"`
	Models      []json.RawMessage `json:"models"`
}
type nativeTextEnvelope struct {
	Payload nativeTextPayload `json:"payload"`
	SHA256  string            `json:"sha256"`
}

func nativeTextValid(s string, maximum int) bool {
	return utf8.ValidString(s) && len(s) > 0 && len(s) <= maximum
}
func nativeTextPercepts(text string) (NativeSessionCase, error) {
	if !nativeTextValid(text, 256) {
		return NativeSessionCase{}, ErrNativeSession
	}
	var nodes [13]int
	for i := range nodes {
		nodes[i] = i + 1
	}
	// All bytes are consumed. Folding is deliberately explicit and lossy; no
	// word boundaries, relation grammar, known answer or semantic labels enter.
	for i, b := range []byte(text) {
		slot := i % 13
		nodes[slot] = (nodes[slot]*257+int(b)+1)%999983 + 1
	}
	edges := make([]NativeSessionEdge, 12)
	for i := range edges {
		edges[i] = NativeSessionEdge{nodes[i], nodes[i+1]}
	}
	return NativeSessionCase{Edges: edges, Questions: [][]int{{nodes[0], nodes[4]}, {nodes[4], nodes[8]}}, Truth: make([]bool, 2)}, nil
}

func TrainNativeTextSelector(xs []NativeTextExample) (*NativeTextSelector, error) {
	if len(xs) < 2 || len(xs) > 16 {
		return nil, ErrNativeSession
	}
	total := 0
	seen := map[string]bool{}
	cases := make([]NativeSessionCase, len(xs))
	for i, x := range xs {
		var err error
		cases[i], err = nativeTextPercepts(x.Input)
		if err != nil || !nativeTextValid(x.Response, 128) {
			return nil, ErrNativeSession
		}
		total += len(x.Input)
		seen[x.Response] = true
	}
	if total > 1024 || len(seen) < 2 || len(seen) > 8 {
		return nil, ErrNativeSession
	}
	m := &NativeTextSelector{}
	for r := range seen {
		m.responses = append(m.responses, r)
	}
	sort.Strings(m.responses)
	raw, err := json.Marshal(xs)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(raw)
	m.trainingSHA = hex.EncodeToString(hash[:])
	for _, response := range m.responses {
		train := make([]NativeSessionCase, len(cases))
		copy(train, cases)
		for i := range train {
			label := xs[i].Response == response
			train[i].Truth = []bool{label, label}
		}
		native, err := TrainNativeSession(train, false)
		if err != nil {
			return nil, err
		}
		m.models = append(m.models, native)
	}
	return m, nil
}

func (m *NativeTextSelector) valid() bool {
	if m == nil || len(m.responses) < 2 || len(m.responses) > 8 || len(m.models) != len(m.responses) || len(m.trainingSHA) != 64 {
		return false
	}
	if _, err := hex.DecodeString(m.trainingSHA); err != nil {
		return false
	}
	for i, r := range m.responses {
		if !nativeTextValid(r, 128) || (i > 0 && m.responses[i-1] >= r) || !m.models[i].valid() || m.models[i].graphLoaded || m.models[i].memory != (nci11WMemory{}) {
			return false
		}
	}
	return true
}

func (m *NativeTextSelector) Predict(text string) (NativeTextPrediction, error) {
	input, err := nativeTextPercepts(text)
	if err != nil || !m.valid() {
		return NativeTextPrediction{}, ErrNativeSession
	}
	edges, err := nativeSessionEdges(input.Edges)
	if err != nil {
		return NativeTextPrediction{}, err
	}
	scores := make([]float64, len(m.models))
	best := 0
	for i, s := range m.models {
		// The original native learned projections build and query the associative
		// matrix. There is no alternate string/lookup classifier on this path.
		memory := nci11WBuild(edges, s.model)
		for _, q := range input.Questions {
			scores[i] += nci11WLogit(memory, s.model, q[0], q[1]) / 2
		}
		if math.IsNaN(scores[i]) || math.IsInf(scores[i], 0) {
			return NativeTextPrediction{}, ErrNativeSession
		}
		if scores[i] > scores[best] {
			best = i
		}
	}
	return NativeTextPrediction{Response: m.responses[best], Responses: append([]string(nil), m.responses...), NativeScores: scores, OutputMode: "closed-response-selection"}, nil
}

func (m *NativeTextSelector) Checkpoint() ([]byte, error) {
	if !m.valid() {
		return nil, ErrNativeSession
	}
	p := nativeTextPayload{Schema: nativeTextSchema, Encoding: nativeTextEncoding, SourceBlob: NativeSessionSource, Responses: m.responses, TrainingSHA: m.trainingSHA}
	for _, s := range m.models {
		raw, err := s.Checkpoint()
		if err != nil {
			return nil, err
		}
		p.Models = append(p.Models, json.RawMessage(raw))
	}
	raw, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(raw)
	out, err := json.Marshal(nativeTextEnvelope{p, hex.EncodeToString(hash[:])})
	if err != nil || len(out) > 131072 {
		return nil, ErrNativeSession
	}
	return out, nil
}

func RestoreNativeTextSelector(raw []byte) (*NativeTextSelector, error) {
	if len(raw) == 0 || len(raw) > 131072 || ValidateNativeSessionJSON(raw) != nil {
		return nil, ErrNativeSession
	}
	var envelope nativeTextEnvelope
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&envelope) != nil || d.Decode(new(any)) != io.EOF {
		return nil, ErrNativeSession
	}
	p := envelope.Payload
	if p.Schema != nativeTextSchema || p.Encoding != nativeTextEncoding || p.SourceBlob != NativeSessionSource || len(p.Models) < 2 || len(p.Models) > 8 {
		return nil, ErrNativeSession
	}
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, ErrNativeSession
	}
	hash := sha256.Sum256(payload)
	if envelope.SHA256 != hex.EncodeToString(hash[:]) {
		return nil, ErrNativeSession
	}
	m := &NativeTextSelector{responses: p.Responses, trainingSHA: p.TrainingSHA}
	for _, raw := range p.Models {
		s, err := RestoreNativeSession(raw)
		if err != nil {
			return nil, err
		}
		m.models = append(m.models, s)
	}
	if !m.valid() {
		return nil, ErrNativeSession
	}
	return m, nil
}
