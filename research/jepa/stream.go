// Package jepa contains an isolated development-only streaming paired-view
// predictive learner. It is not part of the runtime, executor or science ledger.
// The student's recurrent state learns to predict a fixed noncollapsed teacher
// representation of another observation of the same hidden state. Neither
// hidden state nor private rewards enter student training. No transformer,
// tokenizer, agent, model API, controller, or scientific first-use is involved.
package jepa

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
)

const MaxSteps = 256
const MaxEpochs = 300
const SnapshotSize = 4 + 16 + sha256.Size

var invalid = errors.New("WINGLESS_JEPA_INVALID_INPUT_OR_BUDGET")

// Pair is one source observation and one independent target view of a state.
type Pair struct {
	Observation float64
	OtherView   float64
}

// Recurrent is a two-parameter learned recurrent encoder plus predictor.
// State is reset at sequence boundaries; parameters survive cold restoration.
type Recurrent struct {
	Alpha  float64
	Weight float64
}

func New() Recurrent { return Recurrent{Alpha: 0.1, Weight: 1.0} }

func valid(pairs []Pair) bool {
	if len(pairs) < 4 || len(pairs) > MaxSteps {
		return false
	}
	mean := 0.0
	for _, p := range pairs {
		if math.IsNaN(p.Observation) || math.IsInf(p.Observation, 0) ||
			math.IsNaN(p.OtherView) || math.IsInf(p.OtherView, 0) ||
			math.Abs(p.Observation) > 10 || math.Abs(p.OtherView) > 10 {
			return false
		}
		mean += p.OtherView
	}
	mean /= float64(len(pairs))
	variance := 0.0
	for _, p := range pairs {
		d := p.OtherView - mean
		variance += d * d
	}
	return variance/float64(len(pairs)) > 0.00001
}

// DevelopmentSequence constructs two distinct public projections of one
// synthetic time-varying state. Only the fixture generator sees that state.
func DevelopmentSequence(seed uint64, length int) ([]Pair, error) {
	if length < 4 || length > MaxSteps {
		return nil, invalid
	}
	rows := make([]Pair, length)
	var state float64
	for i := range rows {
		seed = seed*6364136223846793005 + 1442695040888963407
		input := float64((seed>>32)%2001)/2000.0 - .5
		state = .6*state + input
		rows[i] = Pair{Observation: input, OtherView: state}
	}
	return rows, nil
}

// Loss returns the mean paired-view prediction MSE; teacher targets do not
// change when the student is trained.
func (m Recurrent) Loss(pairs []Pair) (float64, error) {
	if !valid(pairs) || math.IsNaN(m.Alpha) || math.IsNaN(m.Weight) ||
		m.Alpha < 0 || m.Alpha > .95 || m.Weight < 0 || m.Weight > 3 {
		return 0, invalid
	}
	state, sum := 0.0, 0.0
	for _, p := range pairs {
		state = m.Alpha*state + p.Observation
		d := m.Weight*state - p.OtherView
		sum += d * d
	}
	return sum / float64(len(pairs)), nil
}

// Fit updates recurrent memory retention and student projection by truncated
// sequence-local gradient descent. Each epoch receives the same observations
// and the same fixed teacher views. MemorylessAlpha locks the recurrence for a
// matched per-step conventional baseline; both still learn predictor weight.
func (m *Recurrent) Fit(sequences [][]Pair, epochs int, lr float64, memorylessAlpha bool) error {
	if m == nil || len(sequences) == 0 || len(sequences) > 256 || epochs < 1 ||
		epochs > MaxEpochs || lr <= 0 || lr > .2 || math.IsNaN(lr) || math.IsInf(lr, 0) {
		return invalid
	}
	for _, seq := range sequences {
		if !valid(seq) {
			return invalid
		}
	}
	if m.Alpha < 0 || m.Alpha > .95 || m.Weight < 0 || m.Weight > 3 ||
		math.IsNaN(m.Alpha) || math.IsNaN(m.Weight) {
		return invalid
	}
	if memorylessAlpha {
		m.Alpha = 0
	}
	for epoch := 0; epoch < epochs; epoch++ {
		gAlpha, gWeight, total := 0.0, 0.0, 0
		for _, seq := range sequences {
			state, derivative := 0.0, 0.0
			for _, p := range seq {
				prevState := state
				state = m.Alpha*state + p.Observation
				derivative = prevState + m.Alpha*derivative
				err := m.Weight*state - p.OtherView
				gWeight += err * state
				gAlpha += err * m.Weight * derivative
				total++
			}
		}
		change := lr / float64(total)
		m.Weight = math.Max(0, math.Min(3, m.Weight-change*gWeight))
		if !memorylessAlpha {
			m.Alpha = math.Max(0, math.Min(.95, m.Alpha-change*gAlpha))
		}
		if math.IsNaN(m.Alpha) || math.IsNaN(m.Weight) {
			return invalid
		}
	}
	return nil
}

// Snapshot is a versioned, authenticated exact two-weight cold checkpoint.
func (m Recurrent) Snapshot() ([]byte, error) {
	if m.Alpha < 0 || m.Alpha > .95 || m.Weight < 0 || m.Weight > 3 ||
		math.IsNaN(m.Alpha) || math.IsInf(m.Alpha, 0) ||
		math.IsNaN(m.Weight) || math.IsInf(m.Weight, 0) {
		return nil, invalid
	}
	raw := make([]byte, SnapshotSize)
	copy(raw[:4], []byte("WJP1"))
	binary.BigEndian.PutUint64(raw[4:12], math.Float64bits(m.Alpha))
	binary.BigEndian.PutUint64(raw[12:20], math.Float64bits(m.Weight))
	sum := sha256.Sum256(raw[:20])
	copy(raw[20:], sum[:])
	return raw, nil
}

func Restore(raw []byte) (Recurrent, error) {
	if len(raw) != SnapshotSize || string(raw[:4]) != "WJP1" {
		return Recurrent{}, invalid
	}
	sum := sha256.Sum256(raw[:20])
	if string(sum[:]) != string(raw[20:]) {
		return Recurrent{}, invalid
	}
	m := Recurrent{Alpha: math.Float64frombits(binary.BigEndian.Uint64(raw[4:12])),
		Weight: math.Float64frombits(binary.BigEndian.Uint64(raw[12:20]))}
	if _, err := m.Snapshot(); err != nil {
		return Recurrent{}, err
	}
	return m, nil
}
