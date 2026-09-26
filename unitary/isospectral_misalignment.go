package unitary

import (
	"fmt"
	"math/cmplx"
)

const IsospectralMisalignmentSchema = "wingless.isospectral-misalignment.v1"

type IsospectralMisalignmentDiagnosis struct {
	BaseControlPass                bool    `json:"base_control_pass"`
	MisalignedOrthogonalityPass    bool    `json:"misaligned_orthogonality_pass"`
	IsospectralEquivalencePass     bool    `json:"isospectral_equivalence_pass"`
	NaiveAlignmentBreakPass        bool    `json:"naive_alignment_break_pass"`
	CorrectIntertwinerPass         bool    `json:"correct_intertwiner_pass"`
	DiscoveryCommutatorPass        bool    `json:"discovery_commutator_pass"`
	FeatureInvariancePass          bool    `json:"feature_invariance_pass"`
	MisalignedTrainingPass         bool    `json:"misaligned_training_pass"`
	MisalignedUnseenDepthPass      bool    `json:"misaligned_unseen_depth_pass"`
	MisalignedMutablePass          bool    `json:"misaligned_mutable_pass"`
	AlignmentNotRequiredSupported  bool    `json:"alignment_not_required_supported"`
	MisalignedOrthogonalityError   float64 `json:"misaligned_orthogonality_error"`
	IsospectralEquivalenceError    float64 `json:"isospectral_equivalence_error"`
	NaiveShiftCommutator           float64 `json:"naive_shift_commutator"`
	CorrectShiftCommutator         float64 `json:"correct_shift_commutator"`
	DiscoveryCommutatorError       float64 `json:"discovery_commutator_error"`
	FeatureDrift                   float64 `json:"feature_drift"`
	MisalignedMixerParticipation   float64 `json:"misaligned_mixer_participation_ratio"`
}

type IsospectralMisalignmentProbeResult struct {
	Schema                           string                        `json:"schema"`
	Experiment                       string                        `json:"experiment"`
	LatentDimension                  int                           `json:"latent_dimension"`
	RuntimeStateObjects              int                           `json:"runtime_state_objects"`
	FullCoordinateMixing             bool                          `json:"full_coordinate_mixing"`
	SixfoldSpectrumPreserved         bool                          `json:"sixfold_spectrum_preserved"`
	HiddenPerCopyBasisChange         bool                          `json:"hidden_per_copy_basis_change"`
	EncoderTransportConjugatedTogether bool                        `json:"encoder_transport_conjugated_together"`
	ObserverUsesHiddenBasis          bool                          `json:"observer_uses_hidden_basis"`
	ObserverDiscoveryFromTransport   bool                          `json:"observer_discovery_from_transport"`
	SelectorUsesTrainingLabels       bool                          `json:"selector_uses_training_labels"`
	SelectorUsesHeldOutData          bool                          `json:"selector_uses_heldout_data"`
	SelectorUsesExplicitDepth        bool                          `json:"selector_uses_explicit_depth"`
	PhaseAlphabetSupervision         bool                          `json:"phase_alphabet_supervision"`
	CandidateObservableCount         int                           `json:"candidate_observable_count"`
	RuntimeObservableCount           int                           `json:"runtime_observable_count"`
	TrainingDepths                   []int                         `json:"training_depths"`
	HeldOutDepths                    []int                         `json:"held_out_depths"`
	BaseControl                      DiscoveryBreadthArm            `json:"base_control"`
	BaseIntegration                  DiscoveredIntegration          `json:"base_mutable_integration"`
	MisalignedArm                    DiscoveryBreadthArm            `json:"misaligned_arm"`
	MisalignedIntegration            DiscoveredIntegration          `json:"misaligned_mutable_integration"`
	MisalignedSelectedIndices        []int                          `json:"misaligned_selected_indices"`
	Diagnosis                        IsospectralMisalignmentDiagnosis `json:"diagnosis"`
}

func perCopyPermutationConjugator() (latentMatrix, error) {
	multipliers := []int{1, 3, 5, 7, 9, 11}
	offsets := []int{0, 1, 3, 5, 7, 9}
	if len(multipliers) != compositeChannels || len(offsets) != compositeChannels {
		return nil, fmt.Errorf("per-copy permutation configuration mismatch")
	}

	out := make(latentMatrix, fullLatentDimension)
	for row := range out {
		out[row] = make([]complex128, fullLatentDimension)
	}

	for channel := 0; channel < compositeChannels; channel++ {
		seen := make(map[int]bool, compositeChannelDim)
		for output := 0; output < compositeChannelDim; output++ {
			input := (multipliers[channel]*output + offsets[channel]) % compositeChannelDim
			if seen[input] {
				return nil, fmt.Errorf("per-copy permutation not bijective channel=%d", channel)
			}
			seen[input] = true
			row := channel*compositeChannelDim + output
			column := channel*compositeChannelDim + input
			out[row][column] = 1
		}
	}
	return out, nil
}

func isospectralMisalignedMixer() (latentMatrix, error) {
	base := fullLatentMixer()
	conjugator, err := perCopyPermutationConjugator()
	if err != nil {
		return nil, err
	}
	return latentMatrixMultiply(base, conjugator)
}

func maxLatentMatrixEntryDifference(a, b latentMatrix) (float64, error) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, fmt.Errorf("latent matrix difference dimension mismatch")
	}
	maximum := 0.0
	for row := range a {
		if len(a[row]) != len(b[row]) {
			return 0, fmt.Errorf("latent matrix difference row mismatch")
		}
		for column := range a[row] {
			delta := cmplx.Abs(a[row][column] - b[row][column])
			if delta > maximum {
				maximum = delta
			}
		}
	}
	return maximum, nil
}

func isospectralEquivalenceError(
	baseMixer, misalignedMixer latentMatrix,
	baseStep, misalignedStep latentMatrix,
) (float64, error) {
	baseAdjoint, err := latentAdjoint(baseMixer)
	if err != nil {
		return 0, err
	}
	change, err := latentMatrixMultiply(misalignedMixer, baseAdjoint)
	if err != nil {
		return 0, err
	}
	changeAdjoint, err := latentAdjoint(change)
	if err != nil {
		return 0, err
	}
	left, err := latentMatrixMultiply(change, baseStep)
	if err != nil {
		return 0, err
	}
	conjugated, err := latentMatrixMultiply(left, changeAdjoint)
	if err != nil {
		return 0, err
	}
	return maxLatentMatrixEntryDifference(conjugated, misalignedStep)
}

func RunUP31() (IsospectralMisalignmentProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	baseMixer := fullLatentMixer()
	misalignedMixer, err := isospectralMisalignedMixer()
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	baseStep, err := conjugatedLatentStep(baseMixer, applyStressUnitary)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	misalignedStep, err := conjugatedLatentStep(misalignedMixer, applyStressUnitary)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	misalignedReal, err := realifyMatrix(misalignedStep)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	misalignedOrthogonalityError, err := maxRealOrthogonalityError(misalignedReal)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	equivalenceError, err := isospectralEquivalenceError(
		baseMixer, misalignedMixer, baseStep, misalignedStep,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	naiveShift, err := fullLatentCrossChannelShift(baseMixer)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	correctShift, err := fullLatentCrossChannelShift(misalignedMixer)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	naiveShiftCommutator, err := maxWeylCommutatorEntry(
		[]latentMatrix{naiveShift}, misalignedStep,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	correctShiftCommutator, err := maxWeylCommutatorEntry(
		[]latentMatrix{correctShift}, misalignedStep,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	participation, err := fullLatentMixerParticipation(misalignedMixer)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	baseOps, err := latentDepthOperators(baseStep, allDepths)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	misalignedOps, err := latentDepthOperators(misalignedStep, allDepths)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	baseObservables, err := frozenUP28ObservableBank(baseStep)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	baseRegs, baseHeads, baseResult, err := trainAndEvaluateMultiplicityArm(
		"unitary_aligned_isospectral_control",
		"aligned_repeated_spectrum",
		trainTables, heldTables,
		trainDepths, heldDepths,
		baseOps, baseMixer, baseObservables, memoryNoise,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	baseIntegration, err := runBreadthIntegration(
		baseMixer, baseObservables, baseOps,
		baseRegs, baseHeads,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	candidates, _, err := discoverCommutingObservables(
		misalignedStep, interactionCandidateCount, interactionRounds,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	selectionStates, selectionTables, err := taskSelectedTrainingStates(
		trainTables, misalignedMixer, memoryNoise, 2,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	selectedIndices, selectedObservables, err :=
		selectInteractionRelevantObservables(
			selectionStates, selectionTables,
			candidates, interactionRuntimeCount,
		)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	discoveryCommutator, err := maxWeylCommutatorEntry(
		selectedObservables, misalignedStep,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	featureDrift, err := breadthFeatureDrift(
		misalignedMixer, selectedObservables, misalignedOps[1024],
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	misalignedRegs, misalignedHeads, misalignedResult, err :=
		trainAndEvaluateMultiplicityArm(
			"unitary_isospectral_hidden_basis",
			"isospectral_per_copy_misalignment",
			trainTables, heldTables,
			trainDepths, heldDepths,
			misalignedOps, misalignedMixer,
			selectedObservables, memoryNoise,
		)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}
	misalignedIntegration, err := runBreadthIntegration(
		misalignedMixer, selectedObservables, misalignedOps,
		misalignedRegs, misalignedHeads,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil {
		return IsospectralMisalignmentProbeResult{}, err
	}

	basePass :=
		baseResult.HeldOutAccuracy >= 0.99 &&
			baseIntegration.CommitDecodeAccuracy >= 0.99 &&
			baseIntegration.ExactFinalTableAccuracy >= 0.95 &&
			baseIntegration.RelationalQueryAccuracy >= 0.95
	orthogonalityPass := misalignedOrthogonalityError <= 1e-10
	isospectralPass := equivalenceError <= 1e-10
	naiveBreakPass := naiveShiftCommutator >= 1e-4
	correctIntertwinerPass := correctShiftCommutator <= 1e-10
	discoveryPass := discoveryCommutator <= 1e-5
	invariancePass := featureDrift <= 5e-3
	trainingPass := misalignedResult.TrainAccuracy >= 0.99
	unseenPass :=
		misalignedResult.HeldOutAccuracy >= 0.99 &&
			misalignedResult.MaxNormDrift <= 1e-10
	mutablePass :=
		misalignedIntegration.CommitDecodeAccuracy >= 0.99 &&
			misalignedIntegration.ExactFinalTableAccuracy >= 0.95 &&
			misalignedIntegration.RelationalQueryAccuracy >= 0.95 &&
			misalignedIntegration.MaxNormDrift <= 1e-10

	alignmentNotRequired :=
		basePass &&
			orthogonalityPass &&
			isospectralPass &&
			naiveBreakPass &&
			correctIntertwinerPass &&
			discoveryPass &&
			invariancePass &&
			trainingPass &&
			unseenPass &&
			mutablePass

	return IsospectralMisalignmentProbeResult{
		Schema: IsospectralMisalignmentSchema,
		Experiment: "UP-31-isospectral-hidden-per-copy-misalignment",
		LatentDimension: fullLatentDimension,
		RuntimeStateObjects: 1,
		FullCoordinateMixing: true,
		SixfoldSpectrumPreserved: true,
		HiddenPerCopyBasisChange: true,
		EncoderTransportConjugatedTogether: true,
		ObserverUsesHiddenBasis: false,
		ObserverDiscoveryFromTransport: true,
		SelectorUsesTrainingLabels: true,
		SelectorUsesHeldOutData: false,
		SelectorUsesExplicitDepth: false,
		PhaseAlphabetSupervision: true,
		CandidateObservableCount: interactionCandidateCount,
		RuntimeObservableCount: interactionRuntimeCount,
		TrainingDepths: append([]int(nil), trainDepths...),
		HeldOutDepths: append([]int(nil), heldDepths...),
		BaseControl: baseResult,
		BaseIntegration: baseIntegration,
		MisalignedArm: misalignedResult,
		MisalignedIntegration: misalignedIntegration,
		MisalignedSelectedIndices: append([]int(nil), selectedIndices...),
		Diagnosis: IsospectralMisalignmentDiagnosis{
			BaseControlPass: basePass,
			MisalignedOrthogonalityPass: orthogonalityPass,
			IsospectralEquivalencePass: isospectralPass,
			NaiveAlignmentBreakPass: naiveBreakPass,
			CorrectIntertwinerPass: correctIntertwinerPass,
			DiscoveryCommutatorPass: discoveryPass,
			FeatureInvariancePass: invariancePass,
			MisalignedTrainingPass: trainingPass,
			MisalignedUnseenDepthPass: unseenPass,
			MisalignedMutablePass: mutablePass,
			AlignmentNotRequiredSupported: alignmentNotRequired,
			MisalignedOrthogonalityError: misalignedOrthogonalityError,
			IsospectralEquivalenceError: equivalenceError,
			NaiveShiftCommutator: naiveShiftCommutator,
			CorrectShiftCommutator: correctShiftCommutator,
			DiscoveryCommutatorError: discoveryCommutator,
			FeatureDrift: featureDrift,
			MisalignedMixerParticipation: participation,
		},
	}, nil
}
