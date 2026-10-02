package unitary

import (
	"fmt"
	"math"
)

const SequentialProbeSchema = "wingless.unitary-sequential-probe.v1"

var sequentialPairs = [][2]int{{0, 1}, {1, 2}, {2, 3}}

type SequentialCheckpoint struct {
	Step   int       `json:"step"`
	Params []float64 `json:"params"`
	Loss   float64   `json:"loss"`
}

type SequentialPathResult struct {
	Name                  string                 `json:"name"`
	InitialTrainAccuracy  float64                `json:"initial_train_accuracy"`
	TrainAccuracy         float64                `json:"train_accuracy"`
	HeldOutAccuracy       float64                `json:"held_out_accuracy"`
	InitialLoss           float64                `json:"initial_loss"`
	FinalLoss             float64                `json:"final_loss"`
	Parameters            []float64              `json:"parameters"`
	MaxNormDrift          float64                `json:"max_norm_drift"`
	MaxRoundTripError     float64                `json:"max_round_trip_error"`
	Checkpoints           []SequentialCheckpoint `json:"checkpoints"`
}

type SequentialProbeResult struct {
	Schema             string               `json:"schema"`
	Experiment         string               `json:"experiment"`
	Dimension          int                  `json:"dimension"`
	EventTypes         int                  `json:"event_types"`
	TrainSamples       int                  `json:"train_samples"`
	HeldOutSamples     int                  `json:"held_out_samples"`
	HeldOutLengths     []int                `json:"held_out_lengths"`
	Steps              int                  `json:"steps"`
	LearningRate       float64              `json:"learning_rate"`
	GradientEpsilon    float64              `json:"gradient_epsilon"`
	Unitary            SequentialPathResult `json:"unitary"`
	NonUnitary         SequentialPathResult `json:"non_unitary_matched"`
	OrderSensitiveAB   int                  `json:"order_sensitive_ab_target"`
	OrderSensitiveBA   int                  `json:"order_sensitive_ba_target"`
}

type sequentialSample struct {
	initial int
	events  []int
	target  int
}

func sequentialTarget(initial int, events []int) (int, error) {
	if initial < 0 || initial >= 4 {
		return 0, fmt.Errorf("sequential initial state out of range")
	}
	state := initial
	for _, event := range events {
		if event < 0 || event >= len(sequentialPairs) {
			return 0, fmt.Errorf("sequential event out of range")
		}
		pair := sequentialPairs[event]
		switch state {
		case pair[0]:
			state = pair[1]
		case pair[1]:
			state = pair[0]
		}
	}
	return state, nil
}

func sequentialBasis(index int) (State, error) {
	if index < 0 || index >= 4 {
		return nil, fmt.Errorf("sequential basis index out of range")
	}
	state := make(State, 4)
	state[index] = 1
	return state, nil
}

func sequentialTrainSamples() ([]sequentialSample, error) {
	out := make([]sequentialSample, 0, 12)
	for initial := 0; initial < 4; initial++ {
		for event := 0; event < len(sequentialPairs); event++ {
			target, err := sequentialTarget(initial, []int{event})
			if err != nil {
				return nil, err
			}
			out = append(out, sequentialSample{
				initial: initial,
				events:  []int{event},
				target:  target,
			})
		}
	}
	return out, nil
}

func sequentialHeldOutSamples(lengths []int) ([]sequentialSample, error) {
	var out []sequentialSample
	for _, length := range lengths {
		if length <= 1 {
			return nil, fmt.Errorf("held-out sequence length must exceed one")
		}
		for initial := 0; initial < 4; initial++ {
			for seed := 0; seed < 4; seed++ {
				events := make([]int, length)
				for i := range events {
					events[i] = (seed + i*i + 2*i + initial) % len(sequentialPairs)
				}
				target, err := sequentialTarget(initial, events)
				if err != nil {
					return nil, err
				}
				out = append(out, sequentialSample{
					initial: initial,
					events:  events,
					target:  target,
				})
			}
		}
	}
	return out, nil
}

type sequentialApply func(State, []float64, []int) (State, error)

func applySequentialUnitary(initial State, params []float64, events []int) (State, error) {
	if len(params) != len(sequentialPairs) {
		return nil, fmt.Errorf("sequential unitary parameter count=%d want=%d", len(params), len(sequentialPairs))
	}
	out := append(State(nil), initial...)
	for _, event := range events {
		if event < 0 || event >= len(sequentialPairs) {
			return nil, fmt.Errorf("sequential unitary event out of range")
		}
		pair := sequentialPairs[event]
		var err error
		out, err = Propagate(out, []Coupling{{
			A: pair[0], B: pair[1], Theta: params[event],
		}})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func applySequentialNonUnitary(initial State, params []float64, events []int) (State, error) {
	if err := validateState(initial); err != nil {
		return nil, err
	}
	if len(params) != len(sequentialPairs) {
		return nil, fmt.Errorf("sequential control parameter count=%d want=%d", len(params), len(sequentialPairs))
	}
	out := append(State(nil), initial...)
	for _, event := range events {
		if event < 0 || event >= len(sequentialPairs) {
			return nil, fmt.Errorf("sequential control event out of range")
		}
		g := params[event]
		if !finite(g) {
			return nil, fmt.Errorf("sequential control parameter %d is not finite", event)
		}
		pair := sequentialPairs[event]
		a := out[pair[0]]
		b := out[pair[1]]
		out[pair[0]] = a - complex(g, 0)*b
		out[pair[1]] = complex(g, 0)*a + b
	}
	return out, nil
}

func inverseSequentialUnitary(state State, params []float64, events []int) (State, error) {
	out := append(State(nil), state...)
	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		pair := sequentialPairs[event]
		var err error
		out, err = Propagate(out, []Coupling{{
			A: pair[0], B: pair[1], Theta: -params[event],
		}})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func evaluateSequential(params []float64, samples []sequentialSample, apply sequentialApply, roundTrip bool) (loss, accuracy, maxNormDrift, maxRoundTrip float64, err error) {
	var hits int
	for _, sample := range samples {
		initial, e := sequentialBasis(sample.initial)
		if e != nil {
			err = e
			return
		}
		evolved, e := apply(initial, params, sample.events)
		if e != nil {
			err = e
			return
		}

		probabilities, e := Probabilities(evolved)
		if e != nil {
			err = e
			return
		}
		p := probabilities[sample.target]
		if p < 1e-15 {
			p = 1e-15
		}
		loss += -math.Log(p)

		observed, e := ArgMax(evolved)
		if e != nil {
			err = e
			return
		}
		if observed == sample.target {
			hits++
		}

		norm2, e := NormSquared(evolved)
		if e != nil {
			err = e
			return
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}

		if roundTrip {
			recovered, e := inverseSequentialUnitary(evolved, params, sample.events)
			if e != nil {
				err = e
				return
			}
			distance, e := L2Distance(initial, recovered)
			if e != nil {
				err = e
				return
			}
			if distance > maxRoundTrip {
				maxRoundTrip = distance
			}
		}
	}
	loss /= float64(len(samples))
	accuracy = float64(hits) / float64(len(samples))
	return
}

func sequentialLoss(params []float64, samples []sequentialSample, apply sequentialApply) (float64, error) {
	loss, _, _, _, err := evaluateSequential(params, samples, apply, false)
	return loss, err
}

func trainSequential(name string, train, held []sequentialSample, apply sequentialApply, roundTrip bool, steps int, learningRate, epsilon float64) (SequentialPathResult, error) {
	params := []float64{0.1, 0.1, 0.1}
	initialLoss, initialAccuracy, _, _, err := evaluateSequential(params, train, apply, false)
	if err != nil {
		return SequentialPathResult{}, err
	}

	checkpoints := make([]SequentialCheckpoint, 0, 6)
	for step := 1; step <= steps; step++ {
		gradient := make([]float64, len(params))
		for j := range params {
			plus := append([]float64(nil), params...)
			minus := append([]float64(nil), params...)
			plus[j] += epsilon
			minus[j] -= epsilon
			plusLoss, e := sequentialLoss(plus, train, apply)
			if e != nil {
				return SequentialPathResult{}, e
			}
			minusLoss, e := sequentialLoss(minus, train, apply)
			if e != nil {
				return SequentialPathResult{}, e
			}
			gradient[j] = (plusLoss - minusLoss) / (2 * epsilon)
			if !finite(gradient[j]) {
				return SequentialPathResult{}, fmt.Errorf("%s gradient %d is not finite", name, j)
			}
		}
		for j := range params {
			params[j] -= learningRate * gradient[j]
		}

		if step == 1 || step == 5 || step == 10 || step == 20 || step == 50 || step == steps {
			currentLoss, e := sequentialLoss(params, train, apply)
			if e != nil {
				return SequentialPathResult{}, e
			}
			checkpoints = append(checkpoints, SequentialCheckpoint{
				Step:   step,
				Params: append([]float64(nil), params...),
				Loss:   currentLoss,
			})
		}
	}

	finalLoss, trainAccuracy, trainNorm, trainRoundTrip, err := evaluateSequential(params, train, apply, roundTrip)
	if err != nil {
		return SequentialPathResult{}, err
	}
	_, heldAccuracy, heldNorm, heldRoundTrip, err := evaluateSequential(params, held, apply, roundTrip)
	if err != nil {
		return SequentialPathResult{}, err
	}

	return SequentialPathResult{
		Name:                 name,
		InitialTrainAccuracy: initialAccuracy,
		TrainAccuracy:        trainAccuracy,
		HeldOutAccuracy:      heldAccuracy,
		InitialLoss:          initialLoss,
		FinalLoss:            finalLoss,
		Parameters:           append([]float64(nil), params...),
		MaxNormDrift:         math.Max(trainNorm, heldNorm),
		MaxRoundTripError:    math.Max(trainRoundTrip, heldRoundTrip),
		Checkpoints:          checkpoints,
	}, nil
}

// RunUP4 learns three order-sensitive local state-transition operations from
// single-step supervision only, then composes them over unseen sequences up to
// 128 events.
//
// The unitary and matched non-unitary paths have the same three trainable
// scalars, event topology, training examples, finite-difference optimizer and
// optimization budget. Long-sequence composition is never present in training.
func RunUP4() (SequentialProbeResult, error) {
	const (
		steps        = 100
		learningRate = 0.15
		epsilon      = 1e-6
	)
	lengths := []int{4, 8, 16, 32, 64, 128}

	train, err := sequentialTrainSamples()
	if err != nil {
		return SequentialProbeResult{}, err
	}
	held, err := sequentialHeldOutSamples(lengths)
	if err != nil {
		return SequentialProbeResult{}, err
	}

	unitaryResult, err := trainSequential(
		"unitary",
		train,
		held,
		applySequentialUnitary,
		true,
		steps,
		learningRate,
		epsilon,
	)
	if err != nil {
		return SequentialProbeResult{}, err
	}
	controlResult, err := trainSequential(
		"non_unitary_matched",
		train,
		held,
		applySequentialNonUnitary,
		false,
		steps,
		learningRate,
		epsilon,
	)
	if err != nil {
		return SequentialProbeResult{}, err
	}

	ab, err := sequentialTarget(0, []int{0, 1})
	if err != nil {
		return SequentialProbeResult{}, err
	}
	ba, err := sequentialTarget(0, []int{1, 0})
	if err != nil {
		return SequentialProbeResult{}, err
	}

	return SequentialProbeResult{
		Schema:           SequentialProbeSchema,
		Experiment:       "UP-4-sequential-composition",
		Dimension:        4,
		EventTypes:       len(sequentialPairs),
		TrainSamples:     len(train),
		HeldOutSamples:   len(held),
		HeldOutLengths:   append([]int(nil), lengths...),
		Steps:            steps,
		LearningRate:     learningRate,
		GradientEpsilon:  epsilon,
		Unitary:          unitaryResult,
		NonUnitary:       controlResult,
		OrderSensitiveAB: ab,
		OrderSensitiveBA: ba,
	}, nil
}
