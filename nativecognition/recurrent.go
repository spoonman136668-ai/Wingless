// Package nativecognition contains an isolated, source-only native recurrent
// learning primitive. It is not wired to the worker or an LLM backend.
package nativecognition

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "math"
)

const Channels = 4
const checkpointVersion = 1

type cell struct {
    State float64 `json:"state"`
    Readout float64 `json:"readout"`
    Gate float64 `json:"gate"`
    PriorState float64 `json:"prior_state"`
    PriorInput float64 `json:"prior_input"`
    Primed bool `json:"primed"`
}

// Memory uses a fixed four-lane query-blind recurrent state. A readout learns
// current-input prediction; a bounded local eligibility approximation adjusts
// the recurrent gate. It is a development seam, not a claimed cognitive model.
type Memory struct {
    Rate float64 `json:"rate"`
    Cells [Channels]cell `json:"cells"`
}
type checkpoint struct {
    Version int `json:"version"`
    Memory Memory `json:"memory"`
}

func finite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }
func clamp(x, lo, hi float64) float64 {
    return math.Max(lo, math.Min(hi, x))
}

// New permits a zero-rate frozen-weight control but not negative or excessive learning.
func New(rate float64) (*Memory, error) {
    if !finite(rate) || rate < 0 || rate > .25 {
        return nil, errors.New("invalid learning rate")
    }
    m := &Memory{Rate: rate}
    for i := range m.Cells {
        m.Cells[i].Gate = .5
    }
    return m, nil
}

func (m *Memory) valid() error {
    if m == nil || !finite(m.Rate) || m.Rate < 0 || m.Rate > .25 {
        return errors.New("invalid memory learning rate")
    }
    for _, c := range m.Cells {
        vals := []float64{c.State, c.Readout, c.Gate, c.PriorState, c.PriorInput}
        for _, v := range vals {
            if !finite(v) { return errors.New("nonfinite memory state") }
        }
        if c.Gate < 0 || c.Gate > 1 || c.Readout < -4 || c.Readout > 4 {
            return errors.New("memory parameter outside hard bound")
        }
    }
    return nil
}

func (m *Memory) lane(channel int) (*cell,error) {
    if err := m.valid(); err != nil { return nil,err }
    if channel < 0 || channel >= Channels { return nil,errors.New("channel outside capacity") }
    return &m.Cells[channel],nil
}

// Predict reads the current state without learning or modifying memory.
func (m *Memory) Predict(channel int) (float64,error) {
    c,err := m.lane(channel)
    if err != nil { return 0,err }
    return c.Readout*c.State,nil
}

// Observe updates only from the present observation and previously retained
// state. No future query, answer, heldout label or global attention is consulted.
// Prediction is made BEFORE the observation updates the state.
func (m *Memory) Observe(channel int, value float64) (float64,error) {
    c,err := m.lane(channel)
    if err != nil { return 0,err }
    if !finite(value) || math.Abs(value)>1 { return 0,errors.New("invalid bounded observation") }
    previousState := c.State
    oldReadout := c.Readout
    prediction := oldReadout*previousState
    residual := value-prediction
    // A bounded one-step eligibility approximation (not full BPTT).
    if c.Primed {
        derivative := oldReadout*(c.PriorInput-c.PriorState)
        c.Gate = clamp(c.Gate+m.Rate*residual*derivative,0,1)
    }
    c.Readout = clamp(oldReadout+m.Rate*residual*previousState,-4,4)
    c.PriorState = previousState
    c.PriorInput = value
    c.Primed = true
    c.State = previousState+c.Gate*(value-previousState)
    return prediction,nil
}

// Checkpoint produces a value-copy, versioned, bounded JSON state.
func (m *Memory) Checkpoint() ([]byte,error) {
    if err := m.valid(); err != nil { return nil,err }
    return json.Marshal(checkpoint{Version:checkpointVersion,Memory:*m})
}

// FromCheckpoint rejects unknown fields, trailing content and identity drift.
func FromCheckpoint(data []byte) (*Memory,error) {
    if len(data)==0 || len(data)>8192 { return nil,errors.New("checkpoint size invalid") }
    d:=json.NewDecoder(bytes.NewReader(data))
    d.DisallowUnknownFields()
    var c checkpoint
    if err:=d.Decode(&c);err!=nil { return nil,fmt.Errorf("invalid checkpoint: %w",err) }
    var extra any
    if err:=d.Decode(&extra);err!=io.EOF { return nil,errors.New("checkpoint has trailing data") }
    if c.Version!=checkpointVersion { return nil,errors.New("checkpoint version drift") }
    if err:=c.Memory.valid();err!=nil { return nil,err }
    return &c.Memory,nil
}
