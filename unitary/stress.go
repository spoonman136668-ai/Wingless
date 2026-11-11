package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const StressProbeSchema = "wingless.unitary-stress-probe.v1"

type StressDepthMetric struct {
	Depth               int     `json:"depth"`
	Accuracy            float64 `json:"accuracy"`
	MaxNormDrift        float64 `json:"max_norm_drift"`
	MaxGramError        float64 `json:"max_gram_error"`
	MaxPerturbationGain float64 `json:"max_perturbation_gain"`
}

type StressPathResult struct {
	Name                string              `json:"name"`
	Metrics             []StressDepthMetric `json:"metrics"`
	MaxNormDrift        float64             `json:"max_norm_drift"`
	MaxGramError        float64             `json:"max_gram_error"`
	MaxPerturbationGain float64             `json:"max_perturbation_gain"`
}

type StressProbeResult struct {
	Schema                    string           `json:"schema"`
	Experiment                string           `json:"experiment"`
	Dimension                 int              `json:"dimension"`
	Classes                   int              `json:"classes"`
	DistractorDimensions      int              `json:"distractor_dimensions"`
	Samples                   int              `json:"samples"`
	CouplingsPerBlock         int              `json:"couplings_per_block"`
	Depths                    []int            `json:"depths"`
	Unitary                   StressPathResult `json:"unitary"`
	NonUnitary                StressPathResult `json:"non_unitary_matched"`
	UnitaryMaxRoundTripError  float64          `json:"unitary_max_round_trip_error"`
}

type stressSample struct {
	class int
	state State
	clean State
}

func stressClassState(class int) (State, error) {
	if class < 0 || class >= 4 {
		return nil, fmt.Errorf("stress class out of range")
	}
	signs := up2Signs[class]
	state := make(State, 16)
	for i := 0; i < 4; i++ {
		state[i] = complex(signs[i], 0)
	}
	for i := 4; i < len(state); i++ {
		phase := 0.31*float64(i) + 0.07*float64(i*i)
		state[i] = cmplx.Rect(0.18, phase)
	}
	return Normalize(state)
}

func stressProgram() []Coupling {
	const dimension = 16
	var program []Coupling
	for _, stride := range []int{1, 2, 4, 8} {
		for base := 0; base < dimension; base += 2 * stride {
			for offset := 0; offset < stride; offset++ {
				a := base + offset
				b := a + stride
				index := len(program)
				magnitude := 0.055 + 0.006*float64((index*7+3)%9)
				theta := magnitude
				if index%2 == 1 {
					theta = -theta
				}
				program = append(program, Coupling{A: a, B: b, Theta: theta})
			}
		}
	}
	return program
}

func applyStressUnitary(initial State, block []Coupling, depth int) (State, error) {
	out := append(State(nil), initial...)
	var err error
	for i := 0; i < depth; i++ {
		out, err = Propagate(out, block)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func applyStressNonUnitary(initial State, block []Coupling, depth int) (State, error) {
	if err := validateState(initial); err != nil {
		return nil, err
	}
	out := append(State(nil), initial...)
	for step := 0; step < depth; step++ {
		for _, coupling := range block {
			if err := validateCoupling(coupling, len(out)); err != nil {
				return nil, err
			}
			g := coupling.Theta
			a := out[coupling.A]
			b := out[coupling.B]
			out[coupling.A] = a - complex(g, 0)*b
			out[coupling.B] = complex(g, 0)*a + b
		}
	}
	return out, nil
}

func stressFidelity(a, b State) (float64, error) {
	if err := validateState(a); err != nil {
		return 0, err
	}
	if err := validateState(b); err != nil {
		return 0, err
	}
	if len(a) != len(b) {
		return 0, fmt.Errorf("stress fidelity dimension mismatch")
	}
	var inner complex128
	for i := range a {
		inner += cmplx.Conj(a[i]) * b[i]
	}
	an, err := NormSquared(a)
	if err != nil {
		return 0, err
	}
	bn, err := NormSquared(b)
	if err != nil {
		return 0, err
	}
	if an <= 0 || bn <= 0 {
		return 0, fmt.Errorf("stress fidelity requires positive norms")
	}
	value := cmplx.Abs(inner) / math.Sqrt(an*bn)
	if !finite(value) {
		return 0, fmt.Errorf("stress fidelity is not finite")
	}
	return value, nil
}

func stressGram(states []State) ([][]float64, error) {
	out := make([][]float64, len(states))
	for i := range states {
		out[i] = make([]float64, len(states))
		for j := range states {
			value, err := stressFidelity(states[i], states[j])
			if err != nil {
				return nil, err
			}
			out[i][j] = value
		}
	}
	return out, nil
}

func stressGramError(reference, observed [][]float64) (float64, error) {
	if len(reference) != len(observed) {
		return 0, fmt.Errorf("stress gram size mismatch")
	}
	var maximum float64
	for i := range reference {
		if len(reference[i]) != len(observed[i]) {
			return 0, fmt.Errorf("stress gram row mismatch")
		}
		for j := range reference[i] {
			delta := math.Abs(reference[i][j] - observed[i][j])
			if delta > maximum {
				maximum = delta
			}
		}
	}
	return maximum, nil
}

func stressNoise(class, trial, dimension int) State {
	out := make(State, dimension)
	for i := range out {
		amplitude := 0.02 * (1 + 0.25*float64((i+class+trial)%3))
		phase := 0.41*float64(i) + 0.23*float64(class) + 0.17*float64(trial)
		out[i] = cmplx.Rect(amplitude, phase)
	}
	return out
}

func stressSamples(prototypes []State) ([]stressSample, error) {
	phases := []float64{0, 0.37, -1.11, 2.22}
	out := make([]stressSample, 0, len(prototypes)*8)
	for class, prototype := range prototypes {
		for trial := 0; trial < 8; trial++ {
			phase := cmplx.Rect(1, phases[trial%len(phases)])
			clean := make(State, len(prototype))
			for i := range prototype {
				clean[i] = prototype[i] * phase
			}
			noise := stressNoise(class, trial, len(prototype))
			noisy := make(State, len(prototype))
			for i := range prototype {
				noisy[i] = clean[i] + noise[i]
			}
			normalized, err := Normalize(noisy)
			if err != nil {
				return nil, err
			}
			out = append(out, stressSample{class: class, state: normalized, clean: clean})
		}
	}
	return out, nil
}

type stressApply func(State, []Coupling, int) (State, error)

func evaluateStressPath(name string, prototypes []State, samples []stressSample, referenceGram [][]float64, block []Coupling, depths []int, apply stressApply) (StressPathResult, error) {
	result := StressPathResult{Name: name}
	for _, depth := range depths {
		transformedPrototypes := make([]State, len(prototypes))
		var maxNormDrift float64
		for i, prototype := range prototypes {
			state, err := apply(prototype, block, depth)
			if err != nil {
				return StressPathResult{}, err
			}
			transformedPrototypes[i] = state
			norm2, err := NormSquared(state)
			if err != nil {
				return StressPathResult{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}
		}

		gram, err := stressGram(transformedPrototypes)
		if err != nil {
			return StressPathResult{}, err
		}
		gramError, err := stressGramError(referenceGram, gram)
		if err != nil {
			return StressPathResult{}, err
		}

		var hits int
		var maxPerturbationGain float64
		for _, sample := range samples {
			transformed, err := apply(sample.state, block, depth)
			if err != nil {
				return StressPathResult{}, err
			}
			transformedClean, err := apply(sample.clean, block, depth)
			if err != nil {
				return StressPathResult{}, err
			}

			bestClass := 0
			bestScore := -1.0
			for class, prototype := range transformedPrototypes {
				score, err := stressFidelity(prototype, transformed)
				if err != nil {
					return StressPathResult{}, err
				}
				if score > bestScore {
					bestScore = score
					bestClass = class
				}
			}
			if bestClass == sample.class {
				hits++
			}

			originalPerturbation, err := L2Distance(sample.state, sample.clean)
			if err != nil {
				return StressPathResult{}, err
			}
			transformedPerturbation, err := L2Distance(transformed, transformedClean)
			if err != nil {
				return StressPathResult{}, err
			}
			if originalPerturbation <= 0 {
				return StressPathResult{}, fmt.Errorf("stress perturbation unexpectedly zero")
			}
			gain := transformedPerturbation / originalPerturbation
			if !finite(gain) {
				return StressPathResult{}, fmt.Errorf("stress perturbation gain is not finite")
			}
			if gain > maxPerturbationGain {
				maxPerturbationGain = gain
			}
		}

		metric := StressDepthMetric{
			Depth:               depth,
			Accuracy:            float64(hits) / float64(len(samples)),
			MaxNormDrift:        maxNormDrift,
			MaxGramError:        gramError,
			MaxPerturbationGain: maxPerturbationGain,
		}
		result.Metrics = append(result.Metrics, metric)
		if metric.MaxNormDrift > result.MaxNormDrift {
			result.MaxNormDrift = metric.MaxNormDrift
		}
		if metric.MaxGramError > result.MaxGramError {
			result.MaxGramError = metric.MaxGramError
		}
		if metric.MaxPerturbationGain > result.MaxPerturbationGain {
			result.MaxPerturbationGain = metric.MaxPerturbationGain
		}
	}
	return result, nil
}

// RunUP3 tests depth retention rather than learning. The same deterministic
// 16-dimensional states and the same 32 pair couplings are propagated through
// both paths at increasing depth. Twelve dimensions are distractors.
//
// The unitary path is expected to preserve norm, pairwise fidelity geometry,
// perturbation magnitude, and nearest-prototype classification. The matched
// non-unitary path is measured rather than required to fail.
func RunUP3() (StressProbeResult, error) {
	depths := []int{0, 1, 8, 32, 128}
	block := stressProgram()

	prototypes := make([]State, 4)
	for class := range prototypes {
		state, err := stressClassState(class)
		if err != nil {
			return StressProbeResult{}, err
		}
		prototypes[class] = state
	}
	referenceGram, err := stressGram(prototypes)
	if err != nil {
		return StressProbeResult{}, err
	}
	samples, err := stressSamples(prototypes)
	if err != nil {
		return StressProbeResult{}, err
	}

	unitaryResult, err := evaluateStressPath(
		"unitary",
		prototypes,
		samples,
		referenceGram,
		block,
		depths,
		applyStressUnitary,
	)
	if err != nil {
		return StressProbeResult{}, err
	}

	nonUnitaryResult, err := evaluateStressPath(
		"non_unitary_matched",
		prototypes,
		samples,
		referenceGram,
		block,
		depths,
		applyStressNonUnitary,
	)
	if err != nil {
		return StressProbeResult{}, err
	}

	maxDepth := depths[len(depths)-1]
	var maxRoundTrip float64
	inverse := Invert(block)
	for _, prototype := range prototypes {
		evolved, err := applyStressUnitary(prototype, block, maxDepth)
		if err != nil {
			return StressProbeResult{}, err
		}
		recovered := evolved
		for i := 0; i < maxDepth; i++ {
			recovered, err = Propagate(recovered, inverse)
			if err != nil {
				return StressProbeResult{}, err
			}
		}
		distance, err := L2Distance(prototype, recovered)
		if err != nil {
			return StressProbeResult{}, err
		}
		if distance > maxRoundTrip {
			maxRoundTrip = distance
		}
	}

	return StressProbeResult{
		Schema:                   StressProbeSchema,
		Experiment:               "UP-3-depth-retention",
		Dimension:                16,
		Classes:                  4,
		DistractorDimensions:     12,
		Samples:                  len(samples),
		CouplingsPerBlock:        len(block),
		Depths:                   append([]int(nil), depths...),
		Unitary:                  unitaryResult,
		NonUnitary:               nonUnitaryResult,
		UnitaryMaxRoundTripError: maxRoundTrip,
	}, nil
}
