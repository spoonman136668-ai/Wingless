package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const MultiplicityDoseSchema = "wingless.multiplicity-dose.v1"

const multiplicityDoseRMS = 0.05

type MultiplicityDoseArm struct {
	Name                     string                `json:"name"`
	MaxMultiplicity          int                   `json:"max_multiplicity"`
	Offsets                  []float64             `json:"offsets_radians_per_step"`
	OffsetRMS                float64               `json:"offset_rms"`
	OrthogonalityError       float64               `json:"orthogonality_error"`
	DiscoveryCommutatorError float64               `json:"discovery_commutator_error"`
	FeatureDrift             float64               `json:"feature_drift"`
	Static                   DiscoveryBreadthArm   `json:"static"`
	Integration              DiscoveredIntegration `json:"mutable_integration"`
	SelectedIndices          []int                 `json:"selected_indices"`
}

type MultiplicityDoseDiagnosis struct {
	StructuralValidityPass       bool    `json:"structural_validity_pass"`
	BaseControlPass              bool    `json:"base_control_pass"`
	SingletonMaterialEffectPass  bool    `json:"singleton_material_effect_pass"`
	HeldOrderedSteps             int     `json:"held_ordered_steps"`
	CommitOrderedSteps           int     `json:"commit_ordered_steps"`
	MultiplicityDoseResponsePass bool    `json:"multiplicity_dose_response_pass"`
	BaseHeldAccuracy             float64 `json:"base_held_accuracy"`
	SingletonHeldAccuracy        float64 `json:"singleton_held_accuracy"`
	BaseCommitAccuracy           float64 `json:"base_commit_accuracy"`
	SingletonCommitAccuracy      float64 `json:"singleton_commit_accuracy"`
	HeldEndpointDrop             float64 `json:"held_endpoint_drop"`
	CommitEndpointDrop           float64 `json:"commit_endpoint_drop"`
}

type MultiplicityDoseProbeResult struct {
	Schema                           string                     `json:"schema"`
	Experiment                       string                     `json:"experiment"`
	LatentDimension                  int                        `json:"latent_dimension"`
	RuntimeStateObjects              int                        `json:"runtime_state_objects"`
	FullCoordinateMixing             bool                       `json:"full_coordinate_mixing"`
	SplitDirectionsAnonymous         bool                       `json:"split_directions_anonymous"`
	FixedNonzeroOffsetRMS            float64                    `json:"fixed_nonzero_offset_rms"`
	AllArmsNormPreserving            bool                       `json:"all_arms_norm_preserving"`
	AllArmsRealOrthogonalEquivalent  bool                       `json:"all_arms_real_orthogonal_equivalent"`
	ObserverDiscoveryFromTransport   bool                       `json:"observer_discovery_from_transport"`
	ObserverUsesKnownFactorization   bool                       `json:"observer_uses_known_factorization"`
	SelectorUsesTrainingLabels       bool                       `json:"selector_uses_training_labels"`
	SelectorUsesHeldOutData          bool                       `json:"selector_uses_heldout_data"`
	PhaseAlphabetSupervision         bool                       `json:"phase_alphabet_supervision"`
	CandidateObservableCount         int                        `json:"candidate_observable_count"`
	RuntimeObservableCount           int                        `json:"runtime_observable_count"`
	TrainingDepths                   []int                      `json:"training_depths"`
	HeldOutDepths                    []int                      `json:"held_out_depths"`
	Arms                             []MultiplicityDoseArm      `json:"arms"`
	Diagnosis                        MultiplicityDoseDiagnosis  `json:"diagnosis"`
}

func channelAxisLatentMatrix(matrix ChannelMatrix) latentMatrix {
	out := make(latentMatrix, fullLatentDimension)
	for row := range out {
		out[row] = make([]complex128, fullLatentDimension)
	}
	for rowChannel := 0; rowChannel < compositeChannels; rowChannel++ {
		for columnChannel := 0; columnChannel < compositeChannels; columnChannel++ {
			coefficient := complex(matrix[rowChannel][columnChannel], 0)
			for coordinate := 0; coordinate < compositeChannelDim; coordinate++ {
				row := rowChannel*compositeChannelDim + coordinate
				column := columnChannel*compositeChannelDim + coordinate
				out[row][column] = coefficient
			}
		}
	}
	return out
}

func compositeStepWithAnonymousOffsets(offsets []float64) (latentMatrix, error) {
	if len(offsets) != compositeChannels {
		return nil, fmt.Errorf("dose offset count=%d want=%d", len(offsets), compositeChannels)
	}
	base, err := compositeOneStepMatrix(applyStressUnitary)
	if err != nil {
		return nil, err
	}

	// The split basis is a dense orthogonal mixture of the historical channel
	// axis. Offsets therefore attach to anonymous linear combinations rather
	// than directly to memory/anchor/pilot roles.
	channelMix := channelAxisLatentMatrix(denseChannelMixer())
	channelAdjoint, err := latentAdjoint(channelMix)
	if err != nil {
		return nil, err
	}

	phased := cloneLatentMatrix(base)
	for channel, offset := range offsets {
		phase := cmplx.Rect(1, offset)
		start := channel * compositeChannelDim
		end := start + compositeChannelDim
		for row := start; row < end; row++ {
			for column := 0; column < fullLatentDimension; column++ {
				phased[row][column] *= phase
			}
		}
	}

	left, err := latentMatrixMultiply(channelAdjoint, phased)
	if err != nil {
		return nil, err
	}
	return latentMatrixMultiply(left, channelMix)
}

func fullLatentDoseStep(mixer latentMatrix, offsets []float64) (latentMatrix, error) {
	composite, err := compositeStepWithAnonymousOffsets(offsets)
	if err != nil {
		return nil, err
	}
	return conjugateIntoFullLatent(mixer, composite)
}

func doseOffsetRMS(offsets []float64) float64 {
	if len(offsets) == 0 {
		return 0
	}
	var sum float64
	for _, offset := range offsets {
		sum += offset * offset
	}
	return math.Sqrt(sum / float64(len(offsets)))
}

func doseMaxMultiplicity(offsets []float64) int {
	counts := make(map[uint64]int)
	maximum := 0
	for _, offset := range offsets {
		key := math.Float64bits(offset)
		counts[key]++
		if counts[key] > maximum {
			maximum = counts[key]
		}
	}
	return maximum
}

func multiplicityDoseOffsets() []struct {
	name    string
	offsets []float64
} {
	sigma := multiplicityDoseRMS

	// Every non-control arm has zero mean and the same RMS magnitude.
	// The only intended structural change is the maximum repeated group size.
	a42 := sigma / math.Sqrt(2)
	a222 := sigma * math.Sqrt(3.0/2.0)
	c111111 := sigma / math.Sqrt(35.0/3.0)

	return []struct {
		name    string
		offsets []float64
	}{
		{
			name:    "multiplicity_6_control",
			offsets: []float64{0, 0, 0, 0, 0, 0},
		},
		{
			name: "multiplicity_4_plus_2",
			offsets: []float64{
				a42, a42, a42, a42,
				-2 * a42, -2 * a42,
			},
		},
		{
			name: "multiplicity_3_plus_3",
			offsets: []float64{
				sigma, sigma, sigma,
				-sigma, -sigma, -sigma,
			},
		},
		{
			name: "multiplicity_2_plus_2_plus_2",
			offsets: []float64{
				-a222, -a222,
				0, 0,
				a222, a222,
			},
		},
		{
			name: "multiplicity_1_each",
			offsets: []float64{
				-5 * c111111,
				-3 * c111111,
				-c111111,
				c111111,
				3 * c111111,
				5 * c111111,
			},
		},
	}
}

func evaluateMultiplicityDoseArm(
	name string,
	offsets []float64,
	mixer latentMatrix,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (MultiplicityDoseArm, error) {
	step, err := fullLatentDoseStep(mixer, offsets)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}
	realStep, err := realifyMatrix(step)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}
	orthError, err := maxRealOrthogonalityError(realStep)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}

	ops, err := latentDepthOperators(step, allDepths)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}

	candidates, _, err := discoverCommutingObservables(
		step, interactionCandidateCount, interactionRounds,
	)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}
	selectionStates, selectionTables, err := taskSelectedTrainingStates(
		trainTables, mixer, memoryNoise, 2,
	)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}
	selectedIndices, selectedObservables, err :=
		selectInteractionRelevantObservables(
			selectionStates, selectionTables,
			candidates, interactionRuntimeCount,
		)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}

	commutatorError, err := maxWeylCommutatorEntry(selectedObservables, step)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}
	featureDrift, err := breadthFeatureDrift(
		mixer, selectedObservables, ops[1024],
	)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}

	regressors, classifiers, static, err :=
		trainAndEvaluateMultiplicityArm(
			name,
			"dose_"+name,
			trainTables, heldTables,
			trainDepths, heldDepths,
			ops, mixer, selectedObservables, memoryNoise,
		)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}
	integration, err := runBreadthIntegration(
		mixer, selectedObservables, ops,
		regressors, classifiers,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil {
		return MultiplicityDoseArm{}, err
	}

	return MultiplicityDoseArm{
		Name:                     name,
		MaxMultiplicity:          doseMaxMultiplicity(offsets),
		Offsets:                  append([]float64(nil), offsets...),
		OffsetRMS:                doseOffsetRMS(offsets),
		OrthogonalityError:       orthError,
		DiscoveryCommutatorError: commutatorError,
		FeatureDrift:             featureDrift,
		Static:                   static,
		Integration:              integration,
		SelectedIndices:          append([]int(nil), selectedIndices...),
	}, nil
}

func RunUP32() (MultiplicityDoseProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	specs := multiplicityDoseOffsets()
	arms := make([]MultiplicityDoseArm, 0, len(specs))
	for _, spec := range specs {
		arm, err := evaluateMultiplicityDoseArm(
			spec.name, spec.offsets,
			mixer,
			trainTables, heldTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return MultiplicityDoseProbeResult{}, err
		}
		arms = append(arms, arm)
	}

	structuralPass := true
	for index, arm := range arms {
		if arm.OrthogonalityError > 1e-10 ||
			arm.DiscoveryCommutatorError > 1e-5 ||
			arm.FeatureDrift > 5e-3 {
			structuralPass = false
		}
		if index > 0 && math.Abs(arm.OffsetRMS-multiplicityDoseRMS) > 1e-12 {
			structuralPass = false
		}
	}

	base := arms[0]
	singleton := arms[len(arms)-1]
	basePass :=
		base.Static.HeldOutAccuracy >= 0.99 &&
			base.Integration.CommitDecodeAccuracy >= 0.99 &&
			base.Integration.ExactFinalTableAccuracy >= 0.95 &&
			base.Integration.RelationalQueryAccuracy >= 0.95

	heldDrop := base.Static.HeldOutAccuracy - singleton.Static.HeldOutAccuracy
	commitDrop := base.Integration.CommitDecodeAccuracy -
		singleton.Integration.CommitDecodeAccuracy
	singletonEffect := heldDrop >= 0.25 && commitDrop >= 0.25

	heldOrdered := 0
	commitOrdered := 0
	for i := 0; i < len(arms)-1; i++ {
		// A 0.03 tolerance allows finite-bank learner noise without allowing
		// a large reversal of the multiplicity trend.
		if arms[i+1].Static.HeldOutAccuracy <=
			arms[i].Static.HeldOutAccuracy+0.03 {
			heldOrdered++
		}
		if arms[i+1].Integration.CommitDecodeAccuracy <=
			arms[i].Integration.CommitDecodeAccuracy+0.03 {
			commitOrdered++
		}
	}

	dosePass :=
		structuralPass &&
			basePass &&
			singletonEffect &&
			heldOrdered >= 3 &&
			commitOrdered >= 3

	return MultiplicityDoseProbeResult{
		Schema:                          MultiplicityDoseSchema,
		Experiment:                      "UP-32-repeated-spectrum-multiplicity-dose-response",
		LatentDimension:                 fullLatentDimension,
		RuntimeStateObjects:             1,
		FullCoordinateMixing:            true,
		SplitDirectionsAnonymous:        true,
		FixedNonzeroOffsetRMS:           multiplicityDoseRMS,
		AllArmsNormPreserving:           true,
		AllArmsRealOrthogonalEquivalent: true,
		ObserverDiscoveryFromTransport:  true,
		ObserverUsesKnownFactorization:  false,
		SelectorUsesTrainingLabels:      true,
		SelectorUsesHeldOutData:         false,
		PhaseAlphabetSupervision:        true,
		CandidateObservableCount:        interactionCandidateCount,
		RuntimeObservableCount:          interactionRuntimeCount,
		TrainingDepths:                  append([]int(nil), trainDepths...),
		HeldOutDepths:                   append([]int(nil), heldDepths...),
		Arms:                            arms,
		Diagnosis: MultiplicityDoseDiagnosis{
			StructuralValidityPass:       structuralPass,
			BaseControlPass:              basePass,
			SingletonMaterialEffectPass:  singletonEffect,
			HeldOrderedSteps:             heldOrdered,
			CommitOrderedSteps:           commitOrdered,
			MultiplicityDoseResponsePass: dosePass,
			BaseHeldAccuracy:             base.Static.HeldOutAccuracy,
			SingletonHeldAccuracy:        singleton.Static.HeldOutAccuracy,
			BaseCommitAccuracy:           base.Integration.CommitDecodeAccuracy,
			SingletonCommitAccuracy:      singleton.Integration.CommitDecodeAccuracy,
			HeldEndpointDrop:             heldDrop,
			CommitEndpointDrop:           commitDrop,
		},
	}, nil
}
