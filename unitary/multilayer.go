package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const MultilayerProbeSchema = "wingless.unitary-multilayer-probe.v1"

type OptimizerCheckpoint struct {
	Step   int       `json:"step"`
	Params []float64 `json:"params"`
	Loss   float64   `json:"loss"`
}

type ComparatorResult struct {
	Name            string                `json:"name"`
	InitialAccuracy float64               `json:"initial_accuracy"`
	TrainAccuracy   float64               `json:"train_accuracy"`
	HeldOutAccuracy float64               `json:"held_out_accuracy"`
	InitialLoss     float64               `json:"initial_loss"`
	FinalLoss       float64               `json:"final_loss"`
	Parameters      []float64             `json:"parameters"`
	MaxNormDrift    float64               `json:"max_norm_drift"`
	Checkpoints     []OptimizerCheckpoint `json:"checkpoints"`
}

type MultilayerProbeResult struct {
	Schema                    string           `json:"schema"`
	Experiment                string           `json:"experiment"`
	Dimension                 int              `json:"dimension"`
	Classes                   int              `json:"classes"`
	Couplings                 int              `json:"couplings"`
	TrainSamples              int              `json:"train_samples"`
	HeldOutSamples            int              `json:"held_out_samples"`
	Steps                     int              `json:"steps"`
	LearningRate              float64          `json:"learning_rate"`
	GradientEpsilon           float64          `json:"gradient_epsilon"`
	Unitary                    ComparatorResult `json:"unitary"`
	NonUnitary                ComparatorResult `json:"non_unitary"`
	UnitaryMaxRoundTripError  float64          `json:"unitary_max_round_trip_error"`
}

type multiRouteSample struct {
	id     string
	state  State
	target int
}

var up2Pairs = [][2]int{{0, 1}, {2, 3}, {0, 2}, {1, 3}}

var up2Signs = [4][4]float64{
	{1, 1, 1, 1},
	{1, -1, 1, -1},
	{1, 1, -1, -1},
	{1, -1, -1, 1},
}

func up2Sample(id string, class int, globalPhase float64) multiRouteSample {
	g := cmplx.Rect(0.5, globalPhase)
	state := make(State, 4)
	for i, sign := range up2Signs[class] {
		state[i] = complex(sign, 0) * g
	}
	return multiRouteSample{id: id, state: state, target: class}
}

func up2Samples(phases []float64, prefix string) []multiRouteSample {
	out := make([]multiRouteSample, 0, len(phases)*4)
	for phaseIndex, phase := range phases {
		for class := 0; class < 4; class++ {
			out = append(out, up2Sample(
				fmt.Sprintf("%s_p%d_c%d", prefix, phaseIndex, class),
				class,
				phase,
			))
		}
	}
	return out
}

func applyUP2Unitary(initial State, params []float64) (State, error) {
	if len(params) != len(up2Pairs) {
		return nil, fmt.Errorf("unitary parameter count=%d want=%d", len(params), len(up2Pairs))
	}
	program := make([]Coupling, len(up2Pairs))
	for i, pair := range up2Pairs {
		program[i] = Coupling{A: pair[0], B: pair[1], Theta: params[i]}
	}
	return Propagate(initial, program)
}

// applyUP2NonUnitary uses the exact same pair topology and one scalar per pair
// as the unitary path, but removes the trigonometric normalization constraint.
//
//     [a']   [1  -g] [a]
//     [b'] = [g   1] [b]
//
// The normalized direction can still solve the routing task, while the raw
// state norm is free to grow or shrink. This is a matched functional control,
// not an intentionally incapable baseline.
func applyUP2NonUnitary(initial State, params []float64) (State, error) {
	if err := validateState(initial); err != nil {
		return nil, err
	}
	if len(params) != len(up2Pairs) {
		return nil, fmt.Errorf("non-unitary parameter count=%d want=%d", len(params), len(up2Pairs))
	}
	out := append(State(nil), initial...)
	for i, pair := range up2Pairs {
		g := params[i]
		if !finite(g) {
			return nil, fmt.Errorf("non-unitary parameter %d is not finite", i)
		}
		a := out[pair[0]]
		b := out[pair[1]]
		out[pair[0]] = a - complex(g, 0)*b
		out[pair[1]] = complex(g, 0)*a + b
	}
	return out, nil
}

type up2Apply func(State, []float64) (State, error)

func evaluateUP2(params []float64, samples []multiRouteSample, apply up2Apply, roundTrip bool) (loss, accuracy, maxNormDrift, maxRoundTrip float64, err error) {
	var hits int
	for _, sample := range samples {
		beforeNorm, e := NormSquared(sample.state)
		if e != nil {
			err = e
			return
		}
		evolved, e := apply(sample.state, params)
		if e != nil {
			err = e
			return
		}
		afterNorm, e := NormSquared(evolved)
		if e != nil {
			err = e
			return
		}
		drift := math.Abs(afterNorm - beforeNorm)
		if drift > maxNormDrift {
			maxNormDrift = drift
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

		if roundTrip {
			program := make([]Coupling, len(up2Pairs))
			for i, pair := range up2Pairs {
				program[i] = Coupling{A: pair[0], B: pair[1], Theta: params[i]}
			}
			recovered, e := Propagate(evolved, Invert(program))
			if e != nil {
				err = e
				return
			}
			distance, e := L2Distance(sample.state, recovered)
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

func up2Loss(params []float64, samples []multiRouteSample, apply up2Apply) (float64, error) {
	loss, _, _, _, err := evaluateUP2(params, samples, apply, false)
	return loss, err
}

func trainUP2(name string, train, held []multiRouteSample, apply up2Apply, roundTrip bool, steps int, learningRate, epsilon float64) (ComparatorResult, float64, error) {
	params := make([]float64, len(up2Pairs))
	initialLoss, initialAccuracy, _, _, err := evaluateUP2(params, held, apply, false)
	if err != nil {
		return ComparatorResult{}, 0, err
	}

	checkpoints := make([]OptimizerCheckpoint, 0, 7)
	for step := 1; step <= steps; step++ {
		gradient := make([]float64, len(params))
		for j := range params {
			plus := append([]float64(nil), params...)
			minus := append([]float64(nil), params...)
			plus[j] += epsilon
			minus[j] -= epsilon

			plusLoss, e := up2Loss(plus, train, apply)
			if e != nil {
				return ComparatorResult{}, 0, e
			}
			minusLoss, e := up2Loss(minus, train, apply)
			if e != nil {
				return ComparatorResult{}, 0, e
			}
			gradient[j] = (plusLoss - minusLoss) / (2 * epsilon)
			if !finite(gradient[j]) {
				return ComparatorResult{}, 0, fmt.Errorf("%s gradient %d is not finite", name, j)
			}
		}

		for j := range params {
			params[j] -= learningRate * gradient[j]
		}

		if step == 1 || step%10 == 0 {
			loss, e := up2Loss(params, train, apply)
			if e != nil {
				return ComparatorResult{}, 0, e
			}
			checkpoints = append(checkpoints, OptimizerCheckpoint{
				Step:   step,
				Params: append([]float64(nil), params...),
				Loss:   loss,
			})
		}
	}

	finalLoss, trainAccuracy, trainDrift, trainRoundTrip, err := evaluateUP2(params, train, apply, roundTrip)
	if err != nil {
		return ComparatorResult{}, 0, err
	}
	_, heldAccuracy, heldDrift, heldRoundTrip, err := evaluateUP2(params, held, apply, roundTrip)
	if err != nil {
		return ComparatorResult{}, 0, err
	}

	maxDrift := math.Max(trainDrift, heldDrift)
	maxRoundTrip := math.Max(trainRoundTrip, heldRoundTrip)

	return ComparatorResult{
		Name:            name,
		InitialAccuracy: initialAccuracy,
		TrainAccuracy:   trainAccuracy,
		HeldOutAccuracy: heldAccuracy,
		InitialLoss:     initialLoss,
		FinalLoss:       finalLoss,
		Parameters:      append([]float64(nil), params...),
		MaxNormDrift:    maxDrift,
		Checkpoints:     checkpoints,
	}, maxRoundTrip, nil
}

// RunUP2 raises the primitive from binary/two-dimensional routing to a
// four-class/four-dimensional Walsh phase code. Four learned pair couplings
// must compose into a two-stage interference network.
//
// A parameter-matched non-unitary mixer is trained beside it under the same
// data, optimizer and step budget. The experiment tests capability and
// numerical behavior; it does not assume the unitary path must outperform.
func RunUP2() (MultilayerProbeResult, error) {
	const (
		steps        = 60
		learningRate = 0.15
		epsilon      = 1e-6
	)

	train := up2Samples([]float64{0, math.Pi / 2}, "train")
	held := up2Samples([]float64{0.37, -1.11, 2.22}, "held")

	unitaryResult, maxRoundTrip, err := trainUP2(
		"unitary",
		train,
		held,
		applyUP2Unitary,
		true,
		steps,
		learningRate,
		epsilon,
	)
	if err != nil {
		return MultilayerProbeResult{}, err
	}

	nonUnitaryResult, _, err := trainUP2(
		"non_unitary_matched",
		train,
		held,
		applyUP2NonUnitary,
		false,
		steps,
		learningRate,
		epsilon,
	)
	if err != nil {
		return MultilayerProbeResult{}, err
	}

	return MultilayerProbeResult{
		Schema:                   MultilayerProbeSchema,
		Experiment:               "UP-2-multilayer-phase-routing",
		Dimension:                4,
		Classes:                  4,
		Couplings:                len(up2Pairs),
		TrainSamples:             len(train),
		HeldOutSamples:           len(held),
		Steps:                    steps,
		LearningRate:             learningRate,
		GradientEpsilon:          epsilon,
		Unitary:                   unitaryResult,
		NonUnitary:                nonUnitaryResult,
		UnitaryMaxRoundTripError: maxRoundTrip,
	}, nil
}
