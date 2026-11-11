package unitary

import (
	"fmt"
	"math/cmplx"
)

const MultiplicityAblationSchema = "wingless.multiplicity-ablation.v1"

var multiplicitySplitOffsets = []float64{
	0.0,
	0.017,
	-0.029,
	0.043,
	-0.061,
	0.079,
}

type MultiplicityAblationDiagnosis struct {
	BaseControlPass              bool    `json:"base_control_pass"`
	BaseOrthogonalityPass        bool    `json:"base_orthogonality_pass"`
	SplitOrthogonalityPass       bool    `json:"split_orthogonality_pass"`
	DegeneracyBreakPass          bool    `json:"degeneracy_break_pass"`
	SplitDiscoveryCommutatorPass bool    `json:"split_discovery_commutator_pass"`
	SplitFeatureInvariancePass   bool    `json:"split_feature_invariance_pass"`
	SplitTrainingPass            bool    `json:"split_training_pass"`
	SplitUnseenDepthPass         bool    `json:"split_unseen_depth_pass"`
	SplitMutablePass             bool    `json:"split_mutable_pass"`
	MaterialAblationEffect       bool    `json:"material_ablation_effect"`
	MultiplicityDependenceSupported bool `json:"multiplicity_dependence_supported"`
	BaseCrossShiftCommutator     float64 `json:"base_cross_shift_commutator"`
	SplitCrossShiftCommutator    float64 `json:"split_cross_shift_commutator"`
	BaseOrthogonalityError       float64 `json:"base_orthogonality_error"`
	SplitOrthogonalityError      float64 `json:"split_orthogonality_error"`
	SplitDiscoveryCommutatorError float64 `json:"split_discovery_commutator_error"`
	SplitFeatureDrift            float64 `json:"split_feature_drift"`
	HeldOutAccuracyDrop          float64 `json:"heldout_accuracy_drop"`
	MutableCommitDrop            float64 `json:"mutable_commit_drop"`
}

type MultiplicityAblationProbeResult struct {
	Schema                           string                  `json:"schema"`
	Experiment                       string                  `json:"experiment"`
	LatentDimension                  int                     `json:"latent_dimension"`
	RuntimeStateObjects              int                     `json:"runtime_state_objects"`
	FullCoordinateMixing             bool                    `json:"full_coordinate_mixing"`
	BothTransportsNormPreserving     bool                    `json:"both_transports_norm_preserving"`
	BothHaveRealOrthogonalEquivalent bool                    `json:"both_have_real_orthogonal_equivalent"`
	AblationUsesKnownFactorization   bool                    `json:"ablation_uses_known_factorization"`
	ObserverUsesKnownFactorization   bool                    `json:"observer_uses_known_factorization"`
	ObserverDiscoveryFromTransport   bool                    `json:"observer_discovery_from_transport"`
	PhaseAlphabetSupervision         bool                    `json:"phase_alphabet_supervision"`
	SplitOffsets                     []float64               `json:"split_offsets_radians_per_step"`
	CandidateObservableCount         int                     `json:"candidate_observable_count"`
	RuntimeObservableCount           int                     `json:"runtime_observable_count"`
	TrainingDepths                   []int                   `json:"training_depths"`
	HeldOutDepths                    []int                   `json:"held_out_depths"`
	BaseControl                      DiscoveryBreadthArm     `json:"base_control"`
	BaseIntegration                  DiscoveredIntegration   `json:"base_mutable_integration"`
	SplitArm                         DiscoveryBreadthArm     `json:"split_arm"`
	SplitIntegration                 DiscoveredIntegration   `json:"split_mutable_integration"`
	SplitSelectedIndices             []int                   `json:"split_selected_indices"`
	Diagnosis                        MultiplicityAblationDiagnosis `json:"diagnosis"`
}

func splitCompositeOneStepMatrix() (latentMatrix, error) {
	if len(multiplicitySplitOffsets) != compositeChannels {
		return nil, fmt.Errorf(
			"split offset count=%d want=%d",
			len(multiplicitySplitOffsets), compositeChannels,
		)
	}
	base, err := compositeOneStepMatrix(applyStressUnitary)
	if err != nil {
		return nil, err
	}
	out := cloneLatentMatrix(base)
	for channel := 0; channel < compositeChannels; channel++ {
		phase := cmplx.Rect(1, multiplicitySplitOffsets[channel])
		start := channel * compositeChannelDim
		end := start + compositeChannelDim
		for row := start; row < end; row++ {
			for column := 0; column < fullLatentDimension; column++ {
				out[row][column] *= phase
			}
		}
	}
	return out, nil
}

func conjugateIntoFullLatent(
	mixer latentMatrix,
	matrix latentMatrix,
) (latentMatrix, error) {
	adjoint, err := latentAdjoint(mixer)
	if err != nil {
		return nil, err
	}
	left, err := latentMatrixMultiply(mixer, matrix)
	if err != nil {
		return nil, err
	}
	return latentMatrixMultiply(left, adjoint)
}

func splitFullLatentStep(mixer latentMatrix) (latentMatrix, error) {
	split, err := splitCompositeOneStepMatrix()
	if err != nil {
		return nil, err
	}
	return conjugateIntoFullLatent(mixer, split)
}

func fullLatentCrossChannelShift(mixer latentMatrix) (latentMatrix, error) {
	raw, err := compositeWeylOperator(1, 0)
	if err != nil {
		return nil, err
	}
	return conjugateIntoFullLatent(mixer, raw)
}

func trainAndEvaluateMultiplicityArm(
	name, path string,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	ops map[int]latentMatrix,
	mixer latentMatrix,
	observables []latentMatrix,
	memoryNoise float64,
) (
	[4]breadthRegressor,
	[4]linearSoftmaxHead,
	DiscoveryBreadthArm,
	error,
) {
	trainRows, trainDrift, err := buildBreadthRows(
		trainTables, trainDepths, ops,
		mixer, observables, memoryNoise, 2, 0,
	)
	if err != nil {
		return [4]breadthRegressor{}, [4]linearSoftmaxHead{}, DiscoveryBreadthArm{}, err
	}
	regressors, classifiers, base, err := trainBreadthModel(
		trainRows, len(observables), name, path,
	)
	if err != nil {
		return regressors, classifiers, DiscoveryBreadthArm{}, err
	}
	base.MaxNormDrift = trainDrift
	heldRows, heldDrift, err := buildBreadthRows(
		heldTables, heldDepths, ops,
		mixer, observables, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return regressors, classifiers, DiscoveryBreadthArm{}, err
	}
	result, err := evaluateBreadthModel(
		base, regressors, classifiers, heldRows, heldDrift,
	)
	if err != nil {
		return regressors, classifiers, DiscoveryBreadthArm{}, err
	}
	return regressors, classifiers, result, nil
}

func RunUP30() (MultiplicityAblationProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	baseStep, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	splitStep, err := splitFullLatentStep(mixer)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}

	baseReal, err := realifyMatrix(baseStep)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	splitReal, err := realifyMatrix(splitStep)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	baseOrthError, err := maxRealOrthogonalityError(baseReal)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	splitOrthError, err := maxRealOrthogonalityError(splitReal)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}

	shift, err := fullLatentCrossChannelShift(mixer)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	baseShiftCommutator, err := maxWeylCommutatorEntry(
		[]latentMatrix{shift}, baseStep,
	)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	splitShiftCommutator, err := maxWeylCommutatorEntry(
		[]latentMatrix{shift}, splitStep,
	)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}

	baseOps, err := latentDepthOperators(baseStep, allDepths)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	splitOps, err := latentDepthOperators(splitStep, allDepths)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	baseObservables, err := frozenUP28ObservableBank(baseStep)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	baseRegs, baseHeads, baseResult, err :=
		trainAndEvaluateMultiplicityArm(
			"unitary_repeated_spectrum_control",
			"base_repeated_spectrum",
			trainTables, heldTables,
			trainDepths, heldDepths,
			baseOps, mixer, baseObservables, memoryNoise,
		)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	baseIntegration, err := runBreadthIntegration(
		mixer, baseObservables, baseOps,
		baseRegs, baseHeads,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}

	splitCandidates, _, err := discoverCommutingObservables(
		splitStep, interactionCandidateCount, interactionRounds,
	)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	selectionStates, selectionTables, err := taskSelectedTrainingStates(
		trainTables, mixer, memoryNoise, 2,
	)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	splitIndices, splitObservables, err :=
		selectInteractionRelevantObservables(
			selectionStates, selectionTables,
			splitCandidates, interactionRuntimeCount,
		)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}

	splitDiscoveryCommutator, err := maxWeylCommutatorEntry(
		splitObservables, splitStep,
	)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	splitFeatureDrift, err := breadthFeatureDrift(
		mixer, splitObservables, splitOps[1024],
	)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}

	splitRegs, splitHeads, splitResult, err :=
		trainAndEvaluateMultiplicityArm(
			"unitary_split_spectrum_ablation",
			"split_repeated_spectrum_ablation",
			trainTables, heldTables,
			trainDepths, heldDepths,
			splitOps, mixer, splitObservables, memoryNoise,
		)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}
	splitIntegration, err := runBreadthIntegration(
		mixer, splitObservables, splitOps,
		splitRegs, splitHeads,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil {
		return MultiplicityAblationProbeResult{}, err
	}

	basePass :=
		baseResult.HeldOutAccuracy >= 0.99 &&
			baseIntegration.CommitDecodeAccuracy >= 0.99 &&
			baseIntegration.ExactFinalTableAccuracy >= 0.95 &&
			baseIntegration.RelationalQueryAccuracy >= 0.95
	baseOrthPass := baseOrthError <= 1e-10
	splitOrthPass := splitOrthError <= 1e-10
	degeneracyBreakPass :=
		baseShiftCommutator <= 1e-10 &&
			splitShiftCommutator >= 1e-4
	splitDiscoveryPass := splitDiscoveryCommutator <= 1e-5
	splitInvariancePass := splitFeatureDrift <= 5e-3
	splitTrainingPass := splitResult.TrainAccuracy >= 0.99
	splitUnseenPass :=
		splitResult.HeldOutAccuracy >= 0.99 &&
			splitResult.MaxNormDrift <= 1e-10
	splitMutablePass :=
		splitIntegration.CommitDecodeAccuracy >= 0.99 &&
			splitIntegration.ExactFinalTableAccuracy >= 0.95 &&
			splitIntegration.RelationalQueryAccuracy >= 0.95 &&
			splitIntegration.MaxNormDrift <= 1e-10

	heldDrop := baseResult.HeldOutAccuracy - splitResult.HeldOutAccuracy
	commitDrop :=
		baseIntegration.CommitDecodeAccuracy -
			splitIntegration.CommitDecodeAccuracy
	materialEffect := heldDrop >= 0.10 || commitDrop >= 0.10
	dependenceSupported :=
		basePass &&
			baseOrthPass &&
			splitOrthPass &&
			degeneracyBreakPass &&
			splitDiscoveryPass &&
			splitInvariancePass &&
			materialEffect

	return MultiplicityAblationProbeResult{
		Schema: MultiplicityAblationSchema,
		Experiment: "UP-30-orthogonal-transport-multiplicity-ablation",
		LatentDimension: fullLatentDimension,
		RuntimeStateObjects: 1,
		FullCoordinateMixing: true,
		BothTransportsNormPreserving: true,
		BothHaveRealOrthogonalEquivalent: true,
		AblationUsesKnownFactorization: true,
		ObserverUsesKnownFactorization: false,
		ObserverDiscoveryFromTransport: true,
		PhaseAlphabetSupervision: true,
		SplitOffsets: append([]float64(nil), multiplicitySplitOffsets...),
		CandidateObservableCount: interactionCandidateCount,
		RuntimeObservableCount: interactionRuntimeCount,
		TrainingDepths: append([]int(nil), trainDepths...),
		HeldOutDepths: append([]int(nil), heldDepths...),
		BaseControl: baseResult,
		BaseIntegration: baseIntegration,
		SplitArm: splitResult,
		SplitIntegration: splitIntegration,
		SplitSelectedIndices: append([]int(nil), splitIndices...),
		Diagnosis: MultiplicityAblationDiagnosis{
			BaseControlPass: basePass,
			BaseOrthogonalityPass: baseOrthPass,
			SplitOrthogonalityPass: splitOrthPass,
			DegeneracyBreakPass: degeneracyBreakPass,
			SplitDiscoveryCommutatorPass: splitDiscoveryPass,
			SplitFeatureInvariancePass: splitInvariancePass,
			SplitTrainingPass: splitTrainingPass,
			SplitUnseenDepthPass: splitUnseenPass,
			SplitMutablePass: splitMutablePass,
			MaterialAblationEffect: materialEffect,
			MultiplicityDependenceSupported: dependenceSupported,
			BaseCrossShiftCommutator: baseShiftCommutator,
			SplitCrossShiftCommutator: splitShiftCommutator,
			BaseOrthogonalityError: baseOrthError,
			SplitOrthogonalityError: splitOrthError,
			SplitDiscoveryCommutatorError: splitDiscoveryCommutator,
			SplitFeatureDrift: splitFeatureDrift,
			HeldOutAccuracyDrop: heldDrop,
			MutableCommitDrop: commitDrop,
		},
	}, nil
}
