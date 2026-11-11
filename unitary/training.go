package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const TrainingProbeSchema = "wingless.unitary-training-probe.v1"

type TrainingStep struct {
	Step     int     `json:"step"`
	Theta    float64 `json:"theta"`
	Loss     float64 `json:"loss"`
	Gradient float64 `json:"gradient"`
}

type TrainingEvaluation struct {
	ID            string `json:"id"`
	ExpectedRoute int    `json:"expected_route"`
	ObservedRoute int    `json:"observed_route"`
}

type TrainingProbeResult struct {
	Schema             string               `json:"schema"`
	Experiment         string               `json:"experiment"`
	Steps              int                  `json:"steps"`
	LearningRate       float64              `json:"learning_rate"`
	GradientEpsilon    float64              `json:"gradient_epsilon"`
	InitialTheta       float64              `json:"initial_theta"`
	LearnedTheta       float64              `json:"learned_theta"`
	TargetTheta        float64              `json:"target_theta"`
	InitialLoss        float64              `json:"initial_loss"`
	FinalLoss          float64              `json:"final_loss"`
	InitialAccuracy    float64              `json:"initial_accuracy"`
	TrainAccuracy      float64              `json:"train_accuracy"`
	HeldOutAccuracy    float64              `json:"held_out_accuracy"`
	MaxNormDrift       float64              `json:"max_norm_drift"`
	MaxRoundTripError  float64              `json:"max_round_trip_error"`
	Trace              []TrainingStep       `json:"trace"`
	HeldOutEvaluations []TrainingEvaluation `json:"held_out_evaluations"`
}

type routeSample struct {
	id     string
	state  State
	target int
}

func phaseRouteSample(id string, globalPhase float64, opposed bool) routeSample {
	r := 1 / math.Sqrt2
	g := cmplx.Rect(1, globalPhase)
	second := g * complex(r, 0)
	target := 1
	if opposed {
		second = -second
		target = 0
	}
	return routeSample{
		id:     id,
		state:  State{g * complex(r, 0), second},
		target: target,
	}
}

func up1TrainingSamples() []routeSample {
	return []routeSample{
		phaseRouteSample("train_in_phase_0", 0, false),
		phaseRouteSample("train_opposed_0", 0, true),
		phaseRouteSample("train_in_phase_pi_2", math.Pi/2, false),
		phaseRouteSample("train_opposed_pi_2", math.Pi/2, true),
	}
}

func up1HeldOutSamples() []routeSample {
	return []routeSample{
		phaseRouteSample("held_in_phase_037", 0.37, false),
		phaseRouteSample("held_opposed_037", 0.37, true),
		phaseRouteSample("held_in_phase_neg111", -1.11, false),
		phaseRouteSample("held_opposed_neg111", -1.11, true),
	}
}

func evaluateTheta(theta float64, samples []routeSample) (loss, accuracy, maxNormDrift, maxRoundTrip float64, evaluations []TrainingEvaluation, err error) {
	if !finite(theta) {
		err = fmt.Errorf("training theta must be finite")
		return
	}
	program := []Coupling{{A: 0, B: 1, Theta: theta}}
	var hits int
	for _, sample := range samples {
		before, e := NormSquared(sample.state)
		if e != nil {
			err = e
			return
		}
		evolved, e := Propagate(sample.state, program)
		if e != nil {
			err = e
			return
		}
		after, e := NormSquared(evolved)
		if e != nil {
			err = e
			return
		}
		drift := math.Abs(after - before)
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
		evaluations = append(evaluations, TrainingEvaluation{
			ID:            sample.id,
			ExpectedRoute: sample.target,
			ObservedRoute: observed,
		})

		recovered, e := Propagate(evolved, Invert(program))
		if e != nil {
			err = e
			return
		}
		roundTrip, e := L2Distance(sample.state, recovered)
		if e != nil {
			err = e
			return
		}
		if roundTrip > maxRoundTrip {
			maxRoundTrip = roundTrip
		}
	}
	loss /= float64(len(samples))
	accuracy = float64(hits) / float64(len(samples))
	return
}

func trainingLoss(theta float64, samples []routeSample) (float64, error) {
	loss, _, _, _, _, err := evaluateTheta(theta, samples)
	return loss, err
}

// RunUP1 tests whether the UP-0 coupling can be learned rather than supplied.
// The optimizer is deliberately tiny and deterministic: one trainable unitary
// angle, central-difference gradients, fixed learning rate, and fixed step
// budget. This establishes trainability of the primitive without introducing
// an autodiff framework or changing Wingless inference/worker authority.
func RunUP1() (TrainingProbeResult, error) {
	const (
		steps           = 20
		learningRate    = 0.15
		gradientEpsilon = 1e-6
		initialTheta    = 0.0
	)

	train := up1TrainingSamples()
	heldOut := up1HeldOutSamples()

	initialLoss, initialAccuracy, initialDrift, initialRoundTrip, _, err := evaluateTheta(initialTheta, heldOut)
	if err != nil {
		return TrainingProbeResult{}, err
	}

	theta := initialTheta
	trace := make([]TrainingStep, 0, steps)
	maxDrift := initialDrift
	maxRoundTrip := initialRoundTrip

	for step := 0; step < steps; step++ {
		plusLoss, err := trainingLoss(theta+gradientEpsilon, train)
		if err != nil {
			return TrainingProbeResult{}, err
		}
		minusLoss, err := trainingLoss(theta-gradientEpsilon, train)
		if err != nil {
			return TrainingProbeResult{}, err
		}
		gradient := (plusLoss - minusLoss) / (2 * gradientEpsilon)
		if !finite(gradient) {
			return TrainingProbeResult{}, fmt.Errorf("training gradient is not finite")
		}

		theta -= learningRate * gradient

		loss, _, drift, roundTrip, _, err := evaluateTheta(theta, train)
		if err != nil {
			return TrainingProbeResult{}, err
		}
		if drift > maxDrift {
			maxDrift = drift
		}
		if roundTrip > maxRoundTrip {
			maxRoundTrip = roundTrip
		}
		trace = append(trace, TrainingStep{
			Step:     step + 1,
			Theta:    theta,
			Loss:     loss,
			Gradient: gradient,
		})
	}

	finalLoss, trainAccuracy, trainDrift, trainRoundTrip, _, err := evaluateTheta(theta, train)
	if err != nil {
		return TrainingProbeResult{}, err
	}
	_, heldOutAccuracy, heldDrift, heldRoundTrip, heldEvaluations, err := evaluateTheta(theta, heldOut)
	if err != nil {
		return TrainingProbeResult{}, err
	}

	for _, v := range []float64{trainDrift, heldDrift} {
		if v > maxDrift {
			maxDrift = v
		}
	}
	for _, v := range []float64{trainRoundTrip, heldRoundTrip} {
		if v > maxRoundTrip {
			maxRoundTrip = v
		}
	}

	return TrainingProbeResult{
		Schema:             TrainingProbeSchema,
		Experiment:         "UP-1-trained-unitary-routing",
		Steps:              steps,
		LearningRate:       learningRate,
		GradientEpsilon:    gradientEpsilon,
		InitialTheta:       initialTheta,
		LearnedTheta:       theta,
		TargetTheta:        math.Pi / 4,
		InitialLoss:        initialLoss,
		FinalLoss:          finalLoss,
		InitialAccuracy:    initialAccuracy,
		TrainAccuracy:      trainAccuracy,
		HeldOutAccuracy:    heldOutAccuracy,
		MaxNormDrift:       maxDrift,
		MaxRoundTripError:  maxRoundTrip,
		Trace:              trace,
		HeldOutEvaluations: heldEvaluations,
	}, nil
}
