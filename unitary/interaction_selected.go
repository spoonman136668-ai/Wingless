package unitary

import (
	"fmt"
	"math"
	"sort"
)

const InteractionSelectedSchema = "wingless.unitary-interaction-selected.v1"

const (
	interactionCandidateCount = 128
	interactionRuntimeCount   = 64
	interactionRounds         = discoveredProjectionRounds
	interactionLambda         = discoveredRidgeLambda
)

type interactionScore struct {
	index int
	score float64
}

type InteractionSelectedDiagnosis struct {
	Prefix64ReproductionPass bool    `json:"prefix64_reproduction_pass"`
	SelectorTrainingOnlyPass bool    `json:"selector_training_only_pass"`
	DiscoveryCommutatorPass  bool    `json:"discovery_commutator_pass"`
	FeatureInvariancePass    bool    `json:"feature_invariance_pass"`
	PhaseCodeLearningPass    bool    `json:"phase_code_learning_pass"`
	UnseenDepthPass          bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass   bool    `json:"mutable_integration_pass"`
	Prefix64HeldOutDelta     float64 `json:"prefix64_heldout_delta"`
	SelectedOverlapWithPrefix64 int  `json:"selected_overlap_with_prefix64"`
	SelectedOverlapWithUP27  int     `json:"selected_overlap_with_up27"`
	MaxCommutatorEntryError  float64 `json:"max_commutator_entry_error"`
	MaxFeatureDrift          float64 `json:"max_feature_drift"`
	SelectedTrainAccuracy    float64 `json:"selected_train_accuracy"`
	SelectedHeldOutAccuracy  float64 `json:"selected_heldout_accuracy"`
	MatchedControlAccuracy   float64 `json:"matched_control_accuracy"`
	MeanHeldPhaseCosine      float64 `json:"mean_held_phase_cosine"`
}

type InteractionSelectedProbeResult struct {
	Schema                           string                   `json:"schema"`
	Experiment                       string                   `json:"experiment"`
	LatentDimension                  int                      `json:"latent_dimension"`
	RuntimeStateObjects              int                      `json:"runtime_state_objects"`
	VisibleChannelBlocks             bool                     `json:"visible_channel_blocks"`
	FullCoordinateMixing             bool                     `json:"full_coordinate_mixing"`
	ObservableDiscoveryFromTransport bool                     `json:"observable_discovery_from_transport"`
	DiscoveryUsesHiddenMultiplicity  bool                     `json:"discovery_uses_hidden_multiplicity"`
	DiscoveryUsesHiddenMixer         bool                     `json:"discovery_uses_hidden_mixer"`
	SelectorUsesTrainingLabels       bool                     `json:"selector_uses_training_labels"`
	SelectorUsesHeldOutData          bool                     `json:"selector_uses_heldout_data"`
	SelectorUsesExplicitDepth        bool                     `json:"selector_uses_explicit_depth"`
	SelectorUsesQuadraticInteractions bool                    `json:"selector_uses_quadratic_interactions"`
	CandidateObservableCount         int                      `json:"candidate_observable_count"`
	RuntimeObservableCount           int                      `json:"runtime_observable_count"`
	ProjectionRounds                 int                      `json:"projection_rounds"`
	EffectiveOrbitSize               int                      `json:"effective_orbit_size"`
	SelectedIndices                  []int                    `json:"selected_indices"`
	RidgeLambda                      float64                  `json:"ridge_lambda"`
	TrainingDepths                   []int                    `json:"training_depths"`
	HeldOutDepths                    []int                    `json:"held_out_depths"`
	Prefix64                         DiscoveryBreadthArm       `json:"prefix_64"`
	Selected64                       DiscoveryBreadthArm       `json:"interaction_selected_64"`
	MatchedControl64                 DiscoveryBreadthArm       `json:"matched_control_64"`
	Integration64                    DiscoveredIntegration    `json:"unitary_mutable_integration_64"`
	Diagnosis                        InteractionSelectedDiagnosis `json:"diagnosis"`
}

func interactionRawRows(
	states []State,
	candidates []latentMatrix,
) ([][]float64, error) {
	if len(states) == 0 || len(candidates) == 0 {
		return nil, fmt.Errorf("interaction selector requires rows and candidates")
	}
	rawDim := breadthRawFeatureDim(len(candidates))
	rows := make([][]float64, len(states))
	for row, state := range states {
		features, err := breadthRawFeatures(state, candidates)
		if err != nil {
			return nil, err
		}
		if len(features) != rawDim {
			return nil, fmt.Errorf(
				"interaction raw dimension=%d want=%d",
				len(features), rawDim,
			)
		}
		rows[row] = features
	}
	return rows, nil
}

func polynomialFeatureDot(a, b []float64) (float64, error) {
	if len(a) == 0 || len(a) != len(b) {
		return 0, fmt.Errorf("polynomial feature dot dimension mismatch")
	}
	var dot float64
	var diagonalProducts float64
	for i := range a {
		product := a[i] * b[i]
		dot += product
		diagonalProducts += product * product
	}
	value := dot + 0.5*(dot*dot+diagonalProducts)
	if !finite(value) {
		return 0, fmt.Errorf("polynomial feature dot non-finite")
	}
	return value, nil
}

func interactionDualAlpha(
	rawRows [][]float64,
	tables []memoryTable,
	lambda float64,
) ([][]float64, error) {
	if len(rawRows) == 0 || len(rawRows) != len(tables) || lambda <= 0 {
		return nil, fmt.Errorf("invalid interaction selector training data")
	}
	n := len(rawRows)
	const outputs = 8
	gram := make([][]float64, n)
	right := make([][]float64, n)

	for row := 0; row < n; row++ {
		gram[row] = make([]float64, n)
		right[row] = make([]float64, outputs)
		for entity := 0; entity < 4; entity++ {
			target, err := learnedPhaseTarget(tables[row][entity])
			if err != nil {
				return nil, err
			}
			right[row][entity*2] = target[0]
			right[row][entity*2+1] = target[1]
		}
	}

	for first := 0; first < n; first++ {
		self, err := polynomialFeatureDot(rawRows[first], rawRows[first])
		if err != nil {
			return nil, err
		}
		gram[first][first] = 1 + lambda + self
		for second := first + 1; second < n; second++ {
			value, err := polynomialFeatureDot(
				rawRows[first], rawRows[second],
			)
			if err != nil {
				return nil, err
			}
			value += 1
			gram[first][second] = value
			gram[second][first] = value
		}
	}
	return solveDenseMultiple(gram, right)
}

func interactionObservableScores(
	rawRows [][]float64,
	alpha [][]float64,
	observableCount int,
) ([]float64, error) {
	if len(rawRows) == 0 || len(rawRows) != len(alpha) {
		return nil, fmt.Errorf("interaction score row mismatch")
	}
	rawDim := breadthRawFeatureDim(observableCount)
	if len(rawRows[0]) != rawDim {
		return nil, fmt.Errorf(
			"interaction raw dimension=%d want=%d",
			len(rawRows[0]), rawDim,
		)
	}
	const outputs = 8
	scores := make([]float64, observableCount)

	for dimension := 0; dimension < rawDim; dimension++ {
		var weights [outputs]float64
		for row := range rawRows {
			value := rawRows[row][dimension]
			for output := 0; output < outputs; output++ {
				weights[output] += alpha[row][output] * value
			}
		}
		observable := dimension / 2
		for output := 0; output < outputs; output++ {
			scores[observable] += weights[output] * weights[output]
		}
	}

	for first := 0; first < rawDim; first++ {
		for second := first; second < rawDim; second++ {
			var weights [outputs]float64
			for row := range rawRows {
				value := rawRows[row][first] * rawRows[row][second]
				for output := 0; output < outputs; output++ {
					weights[output] += alpha[row][output] * value
				}
			}
			firstObservable := first / 2
			secondObservable := second / 2
			var energy float64
			for output := 0; output < outputs; output++ {
				energy += weights[output] * weights[output]
			}
			if firstObservable == secondObservable {
				scores[firstObservable] += energy
			} else {
				scores[firstObservable] += 0.5 * energy
				scores[secondObservable] += 0.5 * energy
			}
		}
	}

	for _, score := range scores {
		if !finite(score) {
			return nil, fmt.Errorf("interaction importance is non-finite")
		}
	}
	return scores, nil
}

func selectInteractionRelevantObservables(
	states []State,
	tables []memoryTable,
	candidates []latentMatrix,
	count int,
) ([]int, []latentMatrix, error) {
	rawRows, err := interactionRawRows(states, candidates)
	if err != nil {
		return nil, nil, err
	}
	alpha, err := interactionDualAlpha(
		rawRows, tables, interactionLambda,
	)
	if err != nil {
		return nil, nil, err
	}
	scores, err := interactionObservableScores(
		rawRows, alpha, len(candidates),
	)
	if err != nil {
		return nil, nil, err
	}

	ranked := make([]interactionScore, len(candidates))
	for index := range candidates {
		ranked[index] = interactionScore{
			index: index,
			score: scores[index],
		}
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score == ranked[j].score {
			return ranked[i].index < ranked[j].index
		}
		return ranked[i].score > ranked[j].score
	})

	if count < 1 || count > len(ranked) {
		return nil, nil, fmt.Errorf("invalid interaction selected count=%d", count)
	}
	indices := make([]int, count)
	selected := make([]latentMatrix, count)
	for i := 0; i < count; i++ {
		indices[i] = ranked[i].index
		selected[i] = candidates[ranked[i].index]
	}
	return indices, selected, nil
}

func selectedSetOverlap(a, b []int) int {
	set := make(map[int]struct{}, len(a))
	for _, value := range a {
		set[value] = struct{}{}
	}
	var overlap int
	for _, value := range b {
		if _, ok := set[value]; ok {
			overlap++
		}
	}
	return overlap
}

func up27SelectedIndices() []int {
	return []int{
		2,21,63,13,126,19,64,69,35,118,30,51,45,22,109,68,
		99,97,0,15,102,41,20,117,88,83,42,25,101,93,1,107,
		91,24,80,57,39,47,114,5,116,3,71,74,103,53,40,7,
		98,120,110,58,78,65,82,84,33,59,75,61,115,52,127,104,
	}
}

func RunUP28() (InteractionSelectedProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	unitaryStep, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	controlStep, err := conjugatedLatentStep(mixer, applyStressNonUnitary)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}

	candidates, _, err := discoverCommutingObservables(
		unitaryStep,
		interactionCandidateCount,
		interactionRounds,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	prefix64 := candidates[:interactionRuntimeCount]

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	selectionStates, selectionTables, err := taskSelectedTrainingStates(
		trainTables, mixer, memoryNoise, 2,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	selectedIndices, selected64, err := selectInteractionRelevantObservables(
		selectionStates,
		selectionTables,
		candidates,
		interactionRuntimeCount,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}

	unitaryOps, err := latentDepthOperators(unitaryStep, allDepths)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	controlOps, err := latentDepthOperators(controlStep, heldDepths)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}

	prefixTrainRows, prefixTrainDrift, err := buildBreadthRows(
		trainTables, trainDepths, unitaryOps,
		mixer, prefix64, memoryNoise, 2, 0,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	prefixRegs, prefixHeads, prefixBase, err := trainBreadthModel(
		prefixTrainRows,
		interactionRuntimeCount,
		"unitary_prefix_64_control",
		"unitary_prefix_64_reproduction",
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	prefixBase.MaxNormDrift = prefixTrainDrift
	prefixHeldRows, prefixHeldDrift, err := buildBreadthRows(
		heldTables, heldDepths, unitaryOps,
		mixer, prefix64, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	prefixResult, err := evaluateBreadthModel(
		prefixBase, prefixRegs, prefixHeads,
		prefixHeldRows, prefixHeldDrift,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}

	selectedTrainRows, selectedTrainDrift, err := buildBreadthRows(
		trainTables, trainDepths, unitaryOps,
		mixer, selected64, memoryNoise, 2, 0,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	selectedRegs, selectedHeads, selectedBase, err := trainBreadthModel(
		selectedTrainRows,
		interactionRuntimeCount,
		"unitary_interaction_selected_64",
		"unitary_quadratic_importance_selected_64",
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	selectedBase.MaxNormDrift = selectedTrainDrift

	selectedHeldRows, selectedHeldDrift, err := buildBreadthRows(
		heldTables, heldDepths, unitaryOps,
		mixer, selected64, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	selectedResult, err := evaluateBreadthModel(
		selectedBase, selectedRegs, selectedHeads,
		selectedHeldRows, selectedHeldDrift,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}

	controlRows, controlDrift, err := buildBreadthRows(
		heldTables, heldDepths, controlOps,
		mixer, selected64, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	controlResult, err := evaluateBreadthModel(
		DiscoveryBreadthArm{
			Name:                      "non_unitary_interaction_selected_64",
			Path:                      "non_unitary_same_selected_model",
			ObservableCount:           interactionRuntimeCount,
			RawFeatureDimension:       breadthRawFeatureDim(interactionRuntimeCount),
			QuadraticFeatureDimension: breadthQuadraticFeatureDim(interactionRuntimeCount),
			TrainAccuracy:             selectedBase.TrainAccuracy,
			PerEntityTrain:            append([]float64(nil), selectedBase.PerEntityTrain...),
		},
		selectedRegs, selectedHeads,
		controlRows, controlDrift,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}

	integration, err := runBreadthIntegration(
		mixer, selected64, unitaryOps,
		selectedRegs, selectedHeads,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}

	commutatorError, err := maxWeylCommutatorEntry(
		selected64, unitaryStep,
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}
	featureDrift, err := breadthFeatureDrift(
		mixer, selected64, unitaryOps[1024],
	)
	if err != nil {
		return InteractionSelectedProbeResult{}, err
	}

	prefixDelta := math.Abs(
		prefixResult.HeldOutAccuracy - 0.998779296875,
	)
	prefixPass := prefixDelta <= 1e-12
	selectorTrainingOnlyPass := true
	commutatorPass := commutatorError <= 1e-5
	invariancePass := featureDrift <= 5e-3
	learningPass := selectedResult.TrainAccuracy >= 0.99
	unseenPass :=
		selectedResult.HeldOutAccuracy >= 0.99 &&
			selectedResult.MaxNormDrift <= 1e-10
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-10

	return InteractionSelectedProbeResult{
		Schema:                           InteractionSelectedSchema,
		Experiment:                       "UP-28-interaction-aware-discovered-commutant-selection",
		LatentDimension:                  fullLatentDimension,
		RuntimeStateObjects:              1,
		VisibleChannelBlocks:             false,
		FullCoordinateMixing:             true,
		ObservableDiscoveryFromTransport: true,
		DiscoveryUsesHiddenMultiplicity:  false,
		DiscoveryUsesHiddenMixer:         false,
		SelectorUsesTrainingLabels:       true,
		SelectorUsesHeldOutData:          false,
		SelectorUsesExplicitDepth:        false,
		SelectorUsesQuadraticInteractions: true,
		CandidateObservableCount:         interactionCandidateCount,
		RuntimeObservableCount:           interactionRuntimeCount,
		ProjectionRounds:                 interactionRounds,
		EffectiveOrbitSize:               1 << interactionRounds,
		SelectedIndices:                  append([]int(nil), selectedIndices...),
		RidgeLambda:                      interactionLambda,
		TrainingDepths:                   append([]int(nil), trainDepths...),
		HeldOutDepths:                    append([]int(nil), heldDepths...),
		Prefix64:                         prefixResult,
		Selected64:                       selectedResult,
		MatchedControl64:                 controlResult,
		Integration64:                    integration,
		Diagnosis: InteractionSelectedDiagnosis{
			Prefix64ReproductionPass: prefixPass,
			SelectorTrainingOnlyPass: selectorTrainingOnlyPass,
			DiscoveryCommutatorPass:  commutatorPass,
			FeatureInvariancePass:    invariancePass,
			PhaseCodeLearningPass:    learningPass,
			UnseenDepthPass:          unseenPass,
			MutableIntegrationPass:   mutablePass,
			Prefix64HeldOutDelta:     prefixDelta,
			SelectedOverlapWithPrefix64: selectedPrefixOverlap(
				selectedIndices, interactionRuntimeCount,
			),
			SelectedOverlapWithUP27: selectedSetOverlap(
				selectedIndices, up27SelectedIndices(),
			),
			MaxCommutatorEntryError: commutatorError,
			MaxFeatureDrift:         featureDrift,
			SelectedTrainAccuracy:   selectedResult.TrainAccuracy,
			SelectedHeldOutAccuracy: selectedResult.HeldOutAccuracy,
			MatchedControlAccuracy:  controlResult.HeldOutAccuracy,
			MeanHeldPhaseCosine:     selectedResult.MeanPhaseCosine,
		},
	}, nil
}
