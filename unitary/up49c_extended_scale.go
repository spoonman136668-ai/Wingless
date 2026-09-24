package unitary

import "fmt"

const UP49CExtendedScaleSchema = "wingless.up49c-extended-scale-boundary.v1"

type UP49CSpec struct {
	Dimension int `json:"dimension"`
	Depth     int `json:"depth"`
}

type UP49CCase struct {
	Spec                     UP49CSpec        `json:"spec"`
	Samples                  int              `json:"samples"`
	CouplingsPerBlock        int              `json:"couplings_per_block"`
	Unitary                  StressPathResult `json:"unitary"`
	UnitaryMaxRoundTripError float64          `json:"unitary_max_round_trip_error"`
	Gate                     bool             `json:"gate"`
}

type UP49CExtendedScaleResult struct {
	Schema             string      `json:"schema"`
	Experiment         string      `json:"experiment"`
	SourceUP48CSeal    string      `json:"source_up48c_seal"`
	ScreeningOnly      bool        `json:"screening_only"`
	MatchedControlRun  bool        `json:"matched_control_run"`
	TrialsPerClass     int         `json:"trials_per_class"`
	Cases              []UP49CCase `json:"cases"`
	AllGates           bool        `json:"all_gates"`
	FirstFailingCase   int         `json:"first_failing_case"`
}

func up49cSamples(prototypes []State, trialsPerClass int) ([]stressSample, error) {
	if len(prototypes) != 4 || trialsPerClass <= 0 {
		return nil, fmt.Errorf("UP49C invalid sample configuration")
	}
	base, err := up48cSamples(prototypes)
	if err != nil {
		return nil, err
	}
	out := make([]stressSample, 0, 4*trialsPerClass)
	seen := make([]int, 4)
	for _, sample := range base {
		if seen[sample.class] >= trialsPerClass {
			continue
		}
		out = append(out, sample)
		seen[sample.class]++
	}
	if len(out) != 4*trialsPerClass {
		return nil, fmt.Errorf("UP49C sample count=%d want=%d", len(out), 4*trialsPerClass)
	}
	return out, nil
}

func RunUP49C() (UP49CExtendedScaleResult, error) {
	const trialsPerClass = 2
	specs := []UP49CSpec{
		{Dimension: 128, Depth: 8192},
		{Dimension: 256, Depth: 4096},
		{Dimension: 512, Depth: 2048},
		{Dimension: 1024, Depth: 1024},
	}
	result := UP49CExtendedScaleResult{
		Schema:            UP49CExtendedScaleSchema,
		Experiment:        "UP-49C-extended-scale-boundary-screen",
		SourceUP48CSeal:   "d28cf67488c33509c7f42c6f523755e00031c467",
		ScreeningOnly:     true,
		MatchedControlRun: false,
		TrialsPerClass:    trialsPerClass,
		AllGates:          true,
		FirstFailingCase:  -1,
	}

	for index, spec := range specs {
		prototypes := make([]State, 4)
		for class := range prototypes {
			state, err := up48cClassState(spec.Dimension, class)
			if err != nil {
				return UP49CExtendedScaleResult{}, err
			}
			prototypes[class] = state
		}
		referenceGram, err := stressGram(prototypes)
		if err != nil {
			return UP49CExtendedScaleResult{}, err
		}
		samples, err := up49cSamples(prototypes, trialsPerClass)
		if err != nil {
			return UP49CExtendedScaleResult{}, err
		}
		block, err := up48cProgram(spec.Dimension)
		if err != nil {
			return UP49CExtendedScaleResult{}, err
		}
		depths := []int{0, spec.Depth}
		unitaryPath, err := evaluateStressPath(
			fmt.Sprintf("up49c_unitary_dim_%d_depth_%d", spec.Dimension, spec.Depth),
			prototypes, samples, referenceGram, block, depths,
			applyStressUnitary,
		)
		if err != nil {
			return UP49CExtendedScaleResult{}, err
		}
		roundTrip, err := up48cRoundTrip(prototypes, block, spec.Depth)
		if err != nil {
			return UP49CExtendedScaleResult{}, err
		}
		gate := up48cGate(unitaryPath, roundTrip)
		if !gate && result.FirstFailingCase < 0 {
			result.FirstFailingCase = index
			result.AllGates = false
		}
		result.Cases = append(result.Cases, UP49CCase{
			Spec:                     spec,
			Samples:                  len(samples),
			CouplingsPerBlock:        len(block),
			Unitary:                  unitaryPath,
			UnitaryMaxRoundTripError: roundTrip,
			Gate:                     gate,
		})
	}
	return result, nil
}
