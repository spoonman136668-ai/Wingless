package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const UP48CScaleStressSchema = "wingless.up48c-joint-scale-stress.v1"

type UP48CScaleCase struct {
	Dimension                int              `json:"dimension"`
	DistractorDimensions     int              `json:"distractor_dimensions"`
	Samples                  int              `json:"samples"`
	CouplingsPerBlock        int              `json:"couplings_per_block"`
	Depths                   []int            `json:"depths"`
	Unitary                  StressPathResult `json:"unitary"`
	NonUnitary               StressPathResult `json:"non_unitary_matched"`
	UnitaryMaxRoundTripError float64          `json:"unitary_max_round_trip_error"`
	UnitaryScaleGate         bool             `json:"unitary_scale_gate"`
}

type UP48CScaleStressResult struct {
	Schema         string           `json:"schema"`
	Experiment     string           `json:"experiment"`
	SourceUP47Head string           `json:"source_up47_head"`
	Dimensions     []int            `json:"dimensions"`
	Depths         []int            `json:"depths"`
	Cases          []UP48CScaleCase `json:"cases"`
	AllScaleGates  bool             `json:"all_scale_gates"`
}

func up48cClassState(dimension, class int) (State, error) {
	if dimension < 8 || dimension&(dimension-1) != 0 {
		return nil, fmt.Errorf("UP48C dimension must be power of two >= 8")
	}
	if class < 0 || class >= 4 {
		return nil, fmt.Errorf("UP48C class out of range")
	}
	signs := up2Signs[class]
	state := make(State, dimension)
	for i := 0; i < 4; i++ {
		state[i] = complex(signs[i], 0)
	}
	for i := 4; i < dimension; i++ {
		phase := 0.17*float64(i) + 0.013*float64(i*i) + 0.19*float64(class)
		amplitude := 0.18 * (1 + 0.05*float64(i%5))
		state[i] = cmplx.Rect(amplitude, phase)
	}
	return Normalize(state)
}

func up48cProgram(dimension int) ([]Coupling, error) {
	if dimension < 8 || dimension&(dimension-1) != 0 {
		return nil, fmt.Errorf("UP48C program dimension invalid")
	}
	var program []Coupling
	for stride := 1; stride < dimension; stride *= 2 {
		for base := 0; base < dimension; base += 2 * stride {
			for offset := 0; offset < stride; offset++ {
				a := base + offset
				b := a + stride
				index := len(program)
				magnitude := 0.045 + 0.004*float64((index*11+dimension)%11)
				theta := magnitude
				if (index+dimension)%2 == 1 {
					theta = -theta
				}
				program = append(program, Coupling{A: a, B: b, Theta: theta})
			}
		}
	}
	return program, nil
}

func up48cNoise(class, trial, dimension int) State {
	out := make(State, dimension)
	for i := range out {
		amplitude := 0.012 * (1 + 0.20*float64((i+2*class+trial)%5))
		phase := 0.31*float64(i) + 0.23*float64(class) + 0.17*float64(trial)
		out[i] = cmplx.Rect(amplitude, phase)
	}
	return out
}

func up48cSamples(prototypes []State) ([]stressSample, error) {
	if len(prototypes) != 4 {
		return nil, fmt.Errorf("UP48C requires four prototypes")
	}
	phases := []float64{0, 0.37, -1.11, 2.22}
	out := make([]stressSample, 0, 32)
	for class, prototype := range prototypes {
		for trial := 0; trial < 8; trial++ {
			phase := cmplx.Rect(1, phases[trial%len(phases)])
			clean := make(State, len(prototype))
			for i := range prototype {
				clean[i] = prototype[i] * phase
			}
			noise := up48cNoise(class, trial, len(prototype))
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

func up48cRoundTrip(prototypes []State, block []Coupling, depth int) (float64, error) {
	inverse := Invert(block)
	var maximum float64
	for _, prototype := range prototypes {
		evolved, err := applyStressUnitary(prototype, block, depth)
		if err != nil {
			return 0, err
		}
		recovered := evolved
		for i := 0; i < depth; i++ {
			recovered, err = Propagate(recovered, inverse)
			if err != nil {
				return 0, err
			}
		}
		distance, err := L2Distance(prototype, recovered)
		if err != nil {
			return 0, err
		}
		if distance > maximum {
			maximum = distance
		}
	}
	return maximum, nil
}

func up48cGate(path StressPathResult, roundTrip float64) bool {
	if roundTrip > 1e-8 {
		return false
	}
	for _, metric := range path.Metrics {
		if metric.Accuracy != 1 {
			return false
		}
		if metric.MaxNormDrift > 1e-9 || metric.MaxGramError > 1e-9 {
			return false
		}
		if math.Abs(metric.MaxPerturbationGain-1) > 1e-9 {
			return false
		}
	}
	return true
}

func RunUP48C() (UP48CScaleStressResult, error) {
	dimensions := []int{16, 32, 64, 128}
	depths := []int{0, 128, 512, 2048}
	result := UP48CScaleStressResult{
		Schema:         UP48CScaleStressSchema,
		Experiment:     "UP-48C-joint-dimension-depth-stress",
		SourceUP47Head: "474ebb42082d1a9dd022517edb3c48c4f16dbd1d",
		Dimensions:     append([]int(nil), dimensions...),
		Depths:         append([]int(nil), depths...),
		AllScaleGates:  true,
	}

	for _, dimension := range dimensions {
		prototypes := make([]State, 4)
		for class := range prototypes {
			state, err := up48cClassState(dimension, class)
			if err != nil {
				return UP48CScaleStressResult{}, err
			}
			prototypes[class] = state
		}
		referenceGram, err := stressGram(prototypes)
		if err != nil {
			return UP48CScaleStressResult{}, err
		}
		samples, err := up48cSamples(prototypes)
		if err != nil {
			return UP48CScaleStressResult{}, err
		}
		block, err := up48cProgram(dimension)
		if err != nil {
			return UP48CScaleStressResult{}, err
		}

		unitaryResult, err := evaluateStressPath(
			fmt.Sprintf("unitary_dim_%d", dimension),
			prototypes, samples, referenceGram, block, depths,
			applyStressUnitary,
		)
		if err != nil {
			return UP48CScaleStressResult{}, err
		}
		nonUnitaryResult, err := evaluateStressPath(
			fmt.Sprintf("nonunitary_dim_%d", dimension),
			prototypes, samples, referenceGram, block, depths,
			applyStressNonUnitary,
		)
		if err != nil {
			return UP48CScaleStressResult{}, err
		}
		roundTrip, err := up48cRoundTrip(
			prototypes, block, depths[len(depths)-1],
		)
		if err != nil {
			return UP48CScaleStressResult{}, err
		}
		gate := up48cGate(unitaryResult, roundTrip)
		if !gate {
			result.AllScaleGates = false
		}
		result.Cases = append(result.Cases, UP48CScaleCase{
			Dimension:                dimension,
			DistractorDimensions:     dimension - 4,
			Samples:                  len(samples),
			CouplingsPerBlock:        len(block),
			Depths:                   append([]int(nil), depths...),
			Unitary:                  unitaryResult,
			NonUnitary:               nonUnitaryResult,
			UnitaryMaxRoundTripError: roundTrip,
			UnitaryScaleGate:         gate,
		})
	}
	return result, nil
}
