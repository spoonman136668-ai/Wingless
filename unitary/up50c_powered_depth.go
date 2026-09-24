package unitary

import (
	"fmt"
	"math"
)

const UP50CPoweredDepthSchema = "wingless.up50c-powered-depth.v1"

type UP50CPoweredMetric struct {
	Dimension           int     `json:"dimension"`
	Depth               int     `json:"depth"`
	Samples             int     `json:"samples"`
	Accuracy            float64 `json:"accuracy"`
	MaxNormDrift        float64 `json:"max_norm_drift"`
	MaxGramError        float64 `json:"max_gram_error"`
	MaxPerturbationGain float64 `json:"max_perturbation_gain"`
	MaxRoundTripError   float64 `json:"max_round_trip_error"`
	Gate                bool    `json:"gate"`
}

type UP50CPoweredDepthResult struct {
	Schema                    string                 `json:"schema"`
	Experiment                string                 `json:"experiment"`
	SourceUP49CSeal           string                 `json:"source_up49c_seal"`
	OverlapDimension          int                    `json:"overlap_dimension"`
	OverlapDepth              int                    `json:"overlap_depth"`
	OverlapMaxStateError      float64                `json:"overlap_max_state_error"`
	OverlapEquivalenceGate    bool                   `json:"overlap_equivalence_gate"`
	DeepDepth                 int                    `json:"deep_depth"`
	DeepMetrics               []UP50CPoweredMetric   `json:"deep_metrics"`
	AllDeepGates              bool                   `json:"all_deep_gates"`
}

func up50cOneStepMatrix(dimension int) (latentMatrix, error) {
	block, err := up48cProgram(dimension)
	if err != nil {
		return nil, err
	}
	out := make(latentMatrix, dimension)
	for row := range out {
		out[row] = make([]complex128, dimension)
	}
	for column := 0; column < dimension; column++ {
		basis := make(State, dimension)
		basis[column] = 1
		evolved, err := applyStressUnitary(basis, block, 1)
		if err != nil {
			return nil, err
		}
		for row := 0; row < dimension; row++ {
			out[row][column] = evolved[row]
		}
	}
	return out, nil
}

func up50cPoweredMetric(dimension, depth int) (UP50CPoweredMetric, error) {
	prototypes := make([]State, 4)
	for class := range prototypes {
		state, err := up48cClassState(dimension, class)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		prototypes[class] = state
	}
	referenceGram, err := stressGram(prototypes)
	if err != nil {
		return UP50CPoweredMetric{}, err
	}
	samples, err := up49cSamples(prototypes, 2)
	if err != nil {
		return UP50CPoweredMetric{}, err
	}
	step, err := up50cOneStepMatrix(dimension)
	if err != nil {
		return UP50CPoweredMetric{}, err
	}
	powered, err := latentMatrixPower(step, depth)
	if err != nil {
		return UP50CPoweredMetric{}, err
	}
	adjoint, err := latentAdjoint(powered)
	if err != nil {
		return UP50CPoweredMetric{}, err
	}

	transformedPrototypes := make([]State, len(prototypes))
	var maxNorm, maxRound float64
	for i, prototype := range prototypes {
		state, err := latentMatrixVector(powered, prototype)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		transformedPrototypes[i] = state
		norm2, err := NormSquared(state)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNorm {
			maxNorm = drift
		}
		recovered, err := latentMatrixVector(adjoint, state)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		distance, err := L2Distance(prototype, recovered)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		if distance > maxRound {
			maxRound = distance
		}
	}

	gram, err := stressGram(transformedPrototypes)
	if err != nil {
		return UP50CPoweredMetric{}, err
	}
	gramError, err := stressGramError(referenceGram, gram)
	if err != nil {
		return UP50CPoweredMetric{}, err
	}

	hits := 0
	var maxPerturb float64
	for _, sample := range samples {
		transformed, err := latentMatrixVector(powered, sample.state)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		transformedClean, err := latentMatrixVector(powered, sample.clean)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		bestClass := 0
		bestScore := -1.0
		for class, prototype := range transformedPrototypes {
			score, err := stressFidelity(prototype, transformed)
			if err != nil {
				return UP50CPoweredMetric{}, err
			}
			if score > bestScore {
				bestScore = score
				bestClass = class
			}
		}
		if bestClass == sample.class {
			hits++
		}
		original, err := L2Distance(sample.state, sample.clean)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		after, err := L2Distance(transformed, transformedClean)
		if err != nil {
			return UP50CPoweredMetric{}, err
		}
		gain := after / original
		if gain > maxPerturb {
			maxPerturb = gain
		}
	}

	accuracy := float64(hits) / float64(len(samples))
	gate := accuracy == 1 &&
		maxNorm <= 1e-9 &&
		gramError <= 1e-9 &&
		math.Abs(maxPerturb-1) <= 1e-9 &&
		maxRound <= 1e-8

	return UP50CPoweredMetric{
		Dimension: dimension,
		Depth: depth,
		Samples: len(samples),
		Accuracy: accuracy,
		MaxNormDrift: maxNorm,
		MaxGramError: gramError,
		MaxPerturbationGain: maxPerturb,
		MaxRoundTripError: maxRound,
		Gate: gate,
	}, nil
}

func up50cOverlap(dimension, depth int) (float64, error) {
	prototypes := make([]State, 4)
	for class := range prototypes {
		state, err := up48cClassState(dimension, class)
		if err != nil {
			return 0, err
		}
		prototypes[class] = state
	}
	samples, err := up49cSamples(prototypes, 2)
	if err != nil {
		return 0, err
	}
	step, err := up50cOneStepMatrix(dimension)
	if err != nil {
		return 0, err
	}
	powered, err := latentMatrixPower(step, depth)
	if err != nil {
		return 0, err
	}
	block, err := up48cProgram(dimension)
	if err != nil {
		return 0, err
	}
	states := append([]State{}, prototypes...)
	for _, sample := range samples {
		states = append(states, sample.state, sample.clean)
	}
	var maximum float64
	for _, state := range states {
		fast, err := latentMatrixVector(powered, state)
		if err != nil {
			return 0, err
		}
		literal, err := applyStressUnitary(state, block, depth)
		if err != nil {
			return 0, err
		}
		distance, err := L2Distance(fast, literal)
		if err != nil {
			return 0, err
		}
		if distance > maximum {
			maximum = distance
		}
	}
	return maximum, nil
}

func RunUP50C() (UP50CPoweredDepthResult, error) {
	const (
		overlapDimension = 64
		overlapDepth = 2048
		deepDepth = 1048576
	)
	overlap, err := up50cOverlap(overlapDimension, overlapDepth)
	if err != nil {
		return UP50CPoweredDepthResult{}, err
	}
	equivalence := overlap <= 1e-8
	result := UP50CPoweredDepthResult{
		Schema: UP50CPoweredDepthSchema,
		Experiment: "UP-50C-powered-million-depth",
		SourceUP49CSeal: "c725f5b83a0f38c74acc3afe602b5468c9200c57",
		OverlapDimension: overlapDimension,
		OverlapDepth: overlapDepth,
		OverlapMaxStateError: overlap,
		OverlapEquivalenceGate: equivalence,
		DeepDepth: deepDepth,
		AllDeepGates: true,
	}
	if !equivalence {
		result.AllDeepGates = false
		return result, nil
	}
	for _, dimension := range []int{64, 128} {
		metric, err := up50cPoweredMetric(dimension, deepDepth)
		if err != nil {
			return UP50CPoweredDepthResult{}, err
		}
		if !metric.Gate {
			result.AllDeepGates = false
		}
		result.DeepMetrics = append(result.DeepMetrics, metric)
	}
	if len(result.DeepMetrics) != 2 {
		return UP50CPoweredDepthResult{}, fmt.Errorf("UP50C deep metric count mismatch")
	}
	return result, nil
}
