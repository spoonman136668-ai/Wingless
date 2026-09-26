package unitary

import (
	"fmt"
	"math"
)

const UP48ASizeTransferSchema = "wingless.up48a-size-transfer-composition.v1"

type UP48AInstruction struct {
	Kind  int `json:"kind"`
	Start int `json:"start"`
}

type up48aSample struct {
	dimension int
	initial   int
	program   []UP48AInstruction
	target    int
}

type UP48ACheckpoint struct {
	Step   int       `json:"step"`
	Params []float64 `json:"params"`
	Loss   float64   `json:"loss"`
}

type UP48ADimensionMetric struct {
	Dimension int     `json:"dimension"`
	Samples   int     `json:"samples"`
	Accuracy  float64 `json:"accuracy"`
}

type UP48APathResult struct {
	Name                 string                 `json:"name"`
	InitialTrainAccuracy float64                `json:"initial_train_accuracy"`
	TrainAccuracy        float64                `json:"train_accuracy"`
	HeldOutAccuracy      float64                `json:"held_out_accuracy"`
	HeldOutByDimension   []UP48ADimensionMetric `json:"held_out_by_dimension"`
	InitialLoss          float64                `json:"initial_loss"`
	FinalLoss            float64                `json:"final_loss"`
	Parameters           []float64              `json:"parameters"`
	MaxNormDrift         float64                `json:"max_norm_drift"`
	MaxRoundTripError    float64                `json:"max_round_trip_error"`
	Checkpoints          []UP48ACheckpoint      `json:"checkpoints"`
}

type UP48ADiagnosis struct {
	UnitaryTrainAccuracy              float64 `json:"unitary_train_accuracy"`
	UnitaryHeldOutAccuracy            float64 `json:"unitary_heldout_accuracy"`
	MinimumUnitaryDimensionAccuracy   float64 `json:"minimum_unitary_dimension_accuracy"`
	UnitaryAlgorithmicTransferGate    bool    `json:"unitary_algorithmic_transfer_gate"`
	NonUnitaryHeldOutAccuracy         float64 `json:"nonunitary_heldout_accuracy"`
	UnitaryHeldOutAdvantage           float64 `json:"unitary_heldout_advantage"`
}

type UP48ASizeTransferResult struct {
	Schema             string            `json:"schema"`
	Experiment         string            `json:"experiment"`
	SourceUP47Head     string            `json:"source_up47_head"`
	TrainDimensions    []int             `json:"train_dimensions"`
	HeldOutDimensions  []int             `json:"heldout_dimensions"`
	OperatorSpans      []int             `json:"operator_spans"`
	TrainProgramsLong  bool              `json:"train_programs_long"`
	HeldOutLengths     []int             `json:"heldout_lengths"`
	Steps              int               `json:"steps"`
	LearningRate       float64           `json:"learning_rate"`
	GradientEpsilon    float64           `json:"gradient_epsilon"`
	Unitary            UP48APathResult   `json:"unitary"`
	NonUnitary         UP48APathResult   `json:"nonunitary_matched"`
	Diagnosis          UP48ADiagnosis    `json:"diagnosis"`
}

func up48aPair(dimension int, instruction UP48AInstruction) (int, int, error) {
	if dimension < 4 {
		return 0, 0, fmt.Errorf("UP48A dimension too small")
	}
	if instruction.Kind < 0 || instruction.Kind > 1 {
		return 0, 0, fmt.Errorf("UP48A instruction kind out of range")
	}
	if instruction.Start < 0 || instruction.Start >= dimension {
		return 0, 0, fmt.Errorf("UP48A instruction start out of range")
	}
	span := instruction.Kind + 1
	a := instruction.Start
	b := (instruction.Start + span) % dimension
	if a == b {
		return 0, 0, fmt.Errorf("UP48A degenerate pair")
	}
	return a, b, nil
}

func up48aTarget(dimension, initial int, program []UP48AInstruction) (int, error) {
	if initial < 0 || initial >= dimension {
		return 0, fmt.Errorf("UP48A initial out of range")
	}
	state := initial
	for _, instruction := range program {
		a, b, err := up48aPair(dimension, instruction)
		if err != nil {
			return 0, err
		}
		switch state {
		case a:
			state = b
		case b:
			state = a
		}
	}
	return state, nil
}

func up48aBasis(dimension, index int) (State, error) {
	if index < 0 || index >= dimension {
		return nil, fmt.Errorf("UP48A basis index out of range")
	}
	state := make(State, dimension)
	state[index] = 1
	return state, nil
}

func up48aTrainSamples(dimensions []int) ([]up48aSample, error) {
	var out []up48aSample
	for _, dimension := range dimensions {
		for initial := 0; initial < dimension; initial++ {
			for kind := 0; kind < 2; kind++ {
				for start := 0; start < dimension; start++ {
					program := []UP48AInstruction{{Kind: kind, Start: start}}
					target, err := up48aTarget(dimension, initial, program)
					if err != nil {
						return nil, err
					}
					out = append(out, up48aSample{
						dimension: dimension,
						initial:   initial,
						program:   program,
						target:    target,
					})
				}
			}
		}
	}
	return out, nil
}

func up48aHeldSamples(dimensions, lengths []int) ([]up48aSample, error) {
	var out []up48aSample
	for _, dimension := range dimensions {
		for _, length := range lengths {
			for initial := 0; initial < dimension; initial++ {
				for seed := 0; seed < 4; seed++ {
					program := make([]UP48AInstruction, length)
					for i := range program {
						program[i] = UP48AInstruction{
							Kind:  (seed + i + initial) % 2,
							Start: (3*seed + i*i + 2*i + 5*initial) % dimension,
						}
					}
					target, err := up48aTarget(dimension, initial, program)
					if err != nil {
						return nil, err
					}
					out = append(out, up48aSample{
						dimension: dimension,
						initial:   initial,
						program:   program,
						target:    target,
					})
				}
			}
		}
	}
	return out, nil
}

type up48aApply func(State, []float64, []UP48AInstruction) (State, error)

func up48aApplyUnitary(initial State, params []float64, program []UP48AInstruction) (State, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("UP48A unitary parameter count=%d want=2", len(params))
	}
	out := append(State(nil), initial...)
	for _, instruction := range program {
		a, b, err := up48aPair(len(out), instruction)
		if err != nil {
			return nil, err
		}
		out, err = Propagate(out, []Coupling{{A: a, B: b, Theta: params[instruction.Kind]}})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func up48aApplyNonUnitary(initial State, params []float64, program []UP48AInstruction) (State, error) {
	if err := validateState(initial); err != nil {
		return nil, err
	}
	if len(params) != 2 {
		return nil, fmt.Errorf("UP48A nonunitary parameter count=%d want=2", len(params))
	}
	out := append(State(nil), initial...)
	for _, instruction := range program {
		a, b, err := up48aPair(len(out), instruction)
		if err != nil {
			return nil, err
		}
		g := params[instruction.Kind]
		if !finite(g) {
			return nil, fmt.Errorf("UP48A nonunitary parameter not finite")
		}
		va := out[a]
		vb := out[b]
		out[a] = va - complex(g, 0)*vb
		out[b] = complex(g, 0)*va + vb
	}
	return out, nil
}

func up48aInverseUnitary(state State, params []float64, program []UP48AInstruction) (State, error) {
	out := append(State(nil), state...)
	for i := len(program) - 1; i >= 0; i-- {
		instruction := program[i]
		a, b, err := up48aPair(len(out), instruction)
		if err != nil {
			return nil, err
		}
		out, err = Propagate(out, []Coupling{{A: a, B: b, Theta: -params[instruction.Kind]}})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func up48aEvaluate(params []float64, samples []up48aSample, apply up48aApply, roundTrip bool) (
	loss, accuracy, maxNormDrift, maxRoundTrip float64,
	byDimension []UP48ADimensionMetric,
	err error,
) {
	hits := 0
	type counts struct{ samples, hits int }
	perDim := map[int]*counts{}
	order := []int{}
	for _, sample := range samples {
		initial, e := up48aBasis(sample.dimension, sample.initial)
		if e != nil {
			err = e
			return
		}
		evolved, e := apply(initial, params, sample.program)
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
		if _, ok := perDim[sample.dimension]; !ok {
			perDim[sample.dimension] = &counts{}
			order = append(order, sample.dimension)
		}
		perDim[sample.dimension].samples++
		if observed == sample.target {
			hits++
			perDim[sample.dimension].hits++
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
			recovered, e := up48aInverseUnitary(evolved, params, sample.program)
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
	if len(samples) == 0 {
		err = fmt.Errorf("UP48A empty sample set")
		return
	}
	loss /= float64(len(samples))
	accuracy = float64(hits) / float64(len(samples))
	for _, dimension := range order {
		c := perDim[dimension]
		byDimension = append(byDimension, UP48ADimensionMetric{
			Dimension: dimension,
			Samples:   c.samples,
			Accuracy:  float64(c.hits) / float64(c.samples),
		})
	}
	return
}

func up48aLoss(params []float64, samples []up48aSample, apply up48aApply) (float64, error) {
	loss, _, _, _, _, err := up48aEvaluate(params, samples, apply, false)
	return loss, err
}

func up48aTrain(
	name string,
	train, held []up48aSample,
	apply up48aApply,
	roundTrip bool,
	steps int,
	learningRate, epsilon float64,
) (UP48APathResult, error) {
	params := []float64{0.1, 0.1}
	initialLoss, initialAccuracy, _, _, _, err := up48aEvaluate(params, train, apply, false)
	if err != nil {
		return UP48APathResult{}, err
	}
	checkpoints := make([]UP48ACheckpoint, 0, 7)
	for step := 1; step <= steps; step++ {
		gradient := make([]float64, len(params))
		for j := range params {
			plus := append([]float64(nil), params...)
			minus := append([]float64(nil), params...)
			plus[j] += epsilon
			minus[j] -= epsilon
			plusLoss, e := up48aLoss(plus, train, apply)
			if e != nil {
				return UP48APathResult{}, e
			}
			minusLoss, e := up48aLoss(minus, train, apply)
			if e != nil {
				return UP48APathResult{}, e
			}
			gradient[j] = (plusLoss - minusLoss) / (2 * epsilon)
			if !finite(gradient[j]) {
				return UP48APathResult{}, fmt.Errorf("UP48A %s gradient %d not finite", name, j)
			}
		}
		for j := range params {
			params[j] -= learningRate * gradient[j]
		}
		if step == 1 || step == 5 || step == 10 || step == 20 || step == 40 || step == 80 || step == steps {
			currentLoss, e := up48aLoss(params, train, apply)
			if e != nil {
				return UP48APathResult{}, e
			}
			checkpoints = append(checkpoints, UP48ACheckpoint{
				Step:   step,
				Params: append([]float64(nil), params...),
				Loss:   currentLoss,
			})
		}
	}

	finalLoss, trainAccuracy, trainNorm, trainRound, _, err := up48aEvaluate(params, train, apply, roundTrip)
	if err != nil {
		return UP48APathResult{}, err
	}
	_, heldAccuracy, heldNorm, heldRound, byDimension, err := up48aEvaluate(params, held, apply, roundTrip)
	if err != nil {
		return UP48APathResult{}, err
	}

	return UP48APathResult{
		Name:                 name,
		InitialTrainAccuracy: initialAccuracy,
		TrainAccuracy:        trainAccuracy,
		HeldOutAccuracy:      heldAccuracy,
		HeldOutByDimension:   byDimension,
		InitialLoss:          initialLoss,
		FinalLoss:            finalLoss,
		Parameters:           append([]float64(nil), params...),
		MaxNormDrift:         math.Max(trainNorm, heldNorm),
		MaxRoundTripError:    math.Max(trainRound, heldRound),
		Checkpoints:          checkpoints,
	}, nil
}

func RunUP48A() (UP48ASizeTransferResult, error) {
	const (
		steps        = 120
		learningRate = 0.15
		epsilon      = 1e-6
	)
	trainDimensions := []int{4, 6}
	heldDimensions := []int{8, 12}
	heldLengths := []int{4, 16, 64, 128}

	train, err := up48aTrainSamples(trainDimensions)
	if err != nil {
		return UP48ASizeTransferResult{}, err
	}
	held, err := up48aHeldSamples(heldDimensions, heldLengths)
	if err != nil {
		return UP48ASizeTransferResult{}, err
	}

	unitaryResult, err := up48aTrain(
		"unitary",
		train, held,
		up48aApplyUnitary,
		true,
		steps, learningRate, epsilon,
	)
	if err != nil {
		return UP48ASizeTransferResult{}, err
	}
	nonUnitaryResult, err := up48aTrain(
		"nonunitary_matched",
		train, held,
		up48aApplyNonUnitary,
		false,
		steps, learningRate, epsilon,
	)
	if err != nil {
		return UP48ASizeTransferResult{}, err
	}

	minDimensionAccuracy := 1.0
	for _, metric := range unitaryResult.HeldOutByDimension {
		if metric.Accuracy < minDimensionAccuracy {
			minDimensionAccuracy = metric.Accuracy
		}
	}
	gate :=
		unitaryResult.TrainAccuracy >= 0.99 &&
			unitaryResult.HeldOutAccuracy >= 0.98 &&
			minDimensionAccuracy >= 0.95 &&
			unitaryResult.MaxNormDrift <= 1e-9 &&
			unitaryResult.MaxRoundTripError <= 1e-8

	return UP48ASizeTransferResult{
		Schema:            UP48ASizeTransferSchema,
		Experiment:        "UP-48A-size-transfer-compositional-cognition",
		SourceUP47Head:    "474ebb42082d1a9dd022517edb3c48c4f16dbd1d",
		TrainDimensions:   append([]int(nil), trainDimensions...),
		HeldOutDimensions: append([]int(nil), heldDimensions...),
		OperatorSpans:     []int{1, 2},
		TrainProgramsLong: false,
		HeldOutLengths:    append([]int(nil), heldLengths...),
		Steps:             steps,
		LearningRate:      learningRate,
		GradientEpsilon:   epsilon,
		Unitary:           unitaryResult,
		NonUnitary:        nonUnitaryResult,
		Diagnosis: UP48ADiagnosis{
			UnitaryTrainAccuracy:            unitaryResult.TrainAccuracy,
			UnitaryHeldOutAccuracy:          unitaryResult.HeldOutAccuracy,
			MinimumUnitaryDimensionAccuracy: minDimensionAccuracy,
			UnitaryAlgorithmicTransferGate:  gate,
			NonUnitaryHeldOutAccuracy:       nonUnitaryResult.HeldOutAccuracy,
			UnitaryHeldOutAdvantage:         unitaryResult.HeldOutAccuracy - nonUnitaryResult.HeldOutAccuracy,
		},
	}, nil
}
