package unitary

import (
	"fmt"
	"math"
	"sort"
)

const TaskSelectedCommutantSchema = "wingless.unitary-task-selected-commutant.v1"

const (
	taskSelectedCandidateCount = 128
	taskSelectedRuntimeCount   = 64
	taskSelectedRounds         = discoveredProjectionRounds
	taskSelectedLambda         = discoveredRidgeLambda
)

type TaskSelectedDiagnosis struct {
	Prefix64ReproductionPass bool    `json:"prefix64_reproduction_pass"`
	SelectionTrainingOnlyPass bool   `json:"selection_training_only_pass"`
	DiscoveryCommutatorPass  bool    `json:"discovery_commutator_pass"`
	FeatureInvariancePass    bool    `json:"feature_invariance_pass"`
	PhaseCodeLearningPass    bool    `json:"phase_code_learning_pass"`
	UnseenDepthPass          bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass   bool    `json:"mutable_integration_pass"`
	Prefix64HeldOutDelta     float64 `json:"prefix64_heldout_delta"`
	SelectedOverlapWithPrefix64 int  `json:"selected_overlap_with_prefix64"`
	MaxCommutatorEntryError  float64 `json:"max_commutator_entry_error"`
	MaxFeatureDrift          float64 `json:"max_feature_drift"`
	SelectedTrainAccuracy    float64 `json:"selected_train_accuracy"`
	SelectedHeldOutAccuracy  float64 `json:"selected_heldout_accuracy"`
	MatchedControlAccuracy   float64 `json:"matched_control_accuracy"`
	MeanHeldPhaseCosine      float64 `json:"mean_held_phase_cosine"`
}

type TaskSelectedCommutantProbeResult struct {
	Schema                           string                   `json:"schema"`
	Experiment                       string                   `json:"experiment"`
	LatentDimension                  int                      `json:"latent_dimension"`
	RuntimeStateObjects              int                      `json:"runtime_state_objects"`
	VisibleChannelBlocks             bool                     `json:"visible_channel_blocks"`
	FullCoordinateMixing             bool                     `json:"full_coordinate_mixing"`
	ObservableDiscoveryFromTransport bool                     `json:"observable_discovery_from_transport"`
	DiscoveryUsesHiddenMultiplicity  bool                     `json:"discovery_uses_hidden_multiplicity"`
	DiscoveryUsesHiddenMixer         bool                     `json:"discovery_uses_hidden_mixer"`
	SelectionUsesTrainingLabels      bool                     `json:"selection_uses_training_labels"`
	SelectionUsesHeldOutData         bool                     `json:"selection_uses_heldout_data"`
	SelectionUsesExplicitDepth       bool                     `json:"selection_uses_explicit_depth"`
	CandidateObservableCount         int                      `json:"candidate_observable_count"`
	RuntimeObservableCount           int                      `json:"runtime_observable_count"`
	ProjectionRounds                 int                      `json:"projection_rounds"`
	EffectiveOrbitSize               int                      `json:"effective_orbit_size"`
	SelectedIndices                  []int                    `json:"selected_indices"`
	RidgeLambda                      float64                  `json:"ridge_lambda"`
	TrainingDepths                   []int                    `json:"training_depths"`
	HeldOutDepths                    []int                    `json:"held_out_depths"`
	Prefix64                         DiscoveryBreadthArm       `json:"prefix_64"`
	Selected64                       DiscoveryBreadthArm       `json:"selected_64"`
	MatchedControl64                 DiscoveryBreadthArm       `json:"matched_control_64"`
	Integration64                    DiscoveredIntegration    `json:"unitary_mutable_integration_64"`
	Diagnosis                        TaskSelectedDiagnosis     `json:"diagnosis"`
}

type observableRelevance struct {
	index int
	score float64
}

func taskSelectedTrainingStates(
	tables []memoryTable,
	mixer latentMatrix,
	memoryNoise float64,
	trials int,
) ([]State, []memoryTable, error) {
	states := make([]State, 0, len(tables)*trials)
	targetTables := make([]memoryTable, 0, len(tables)*trials)
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, nil, err
		}
		for trial := 0; trial < trials; trial++ {
			seed := memoryTableIndex(table)*1000 + trial*17
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return nil, nil, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.271*float64(seed+1), 2*math.Pi),
			)
			state, err := fullLatentEncode(memory, mixer)
			if err != nil {
				return nil, nil, err
			}
			states = append(states, state)
			targetTables = append(targetTables, table)
		}
	}
	return states, targetTables, nil
}

func phaseTargetMatrix(tables []memoryTable) ([][]float64, error) {
	targets := make([][]float64, len(tables))
	for row, table := range tables {
		targets[row] = make([]float64, 8)
		for entity := 0; entity < 4; entity++ {
			target, err := learnedPhaseTarget(table[entity])
			if err != nil {
				return nil, err
			}
			targets[row][entity*2] = target[0]
			targets[row][entity*2+1] = target[1]
		}
	}
	return targets, nil
}

func centeredVariance(values []float64) (float64, float64) {
	var mean float64
	for _, value := range values {
		mean += value
	}
	mean /= float64(len(values))
	var sum float64
	for _, value := range values {
		delta := value - mean
		sum += delta * delta
	}
	return mean, sum
}

func normalizedCovarianceScore(
	values []float64,
	targetColumns [][]float64,
) float64 {
	mean, variance := centeredVariance(values)
	if variance <= 1e-18 {
		return 0
	}
	var score float64
	for _, target := range targetColumns {
		targetMean, targetVariance := centeredVariance(target)
		if targetVariance <= 1e-18 {
			continue
		}
		var covariance float64
		for row, value := range values {
			covariance += (value - mean) * (target[row] - targetMean)
		}
		score += covariance * covariance / (variance * targetVariance)
	}
	return score
}

func selectTaskRelevantObservables(
	states []State,
	tables []memoryTable,
	candidates []latentMatrix,
	count int,
) ([]int, []latentMatrix, error) {
	if len(states) == 0 || len(states) != len(tables) {
		return nil, nil, fmt.Errorf("task selection training rows mismatch")
	}
	if count < 1 || count > len(candidates) {
		return nil, nil, fmt.Errorf("invalid selected observable count=%d", count)
	}

	targets, err := phaseTargetMatrix(tables)
	if err != nil {
		return nil, nil, err
	}
	targetColumns := make([][]float64, 8)
	for column := 0; column < 8; column++ {
		targetColumns[column] = make([]float64, len(targets))
		for row := range targets {
			targetColumns[column][row] = targets[row][column]
		}
	}

	ranked := make([]observableRelevance, 0, len(candidates))
	for candidateIndex, candidate := range candidates {
		realValues := make([]float64, len(states))
		imagValues := make([]float64, len(states))
		for row, state := range states {
			applied, err := latentMatrixVector(candidate, state)
			if err != nil {
				return nil, nil, err
			}
			value, err := stateInner(state, applied)
			if err != nil {
				return nil, nil, err
			}
			realValues[row] = real(value)
			imagValues[row] = imag(value)
		}
		score :=
			normalizedCovarianceScore(realValues, targetColumns) +
				normalizedCovarianceScore(imagValues, targetColumns)
		if !finite(score) {
			return nil, nil, fmt.Errorf(
				"task relevance score non-finite candidate=%d",
				candidateIndex,
			)
		}
		ranked = append(ranked, observableRelevance{
			index: candidateIndex,
			score: score,
		})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score == ranked[j].score {
			return ranked[i].index < ranked[j].index
		}
		return ranked[i].score > ranked[j].score
	})

	indices := make([]int, count)
	selected := make([]latentMatrix, count)
	for i := 0; i < count; i++ {
		indices[i] = ranked[i].index
		selected[i] = candidates[ranked[i].index]
	}
	return indices, selected, nil
}

func selectedPrefixOverlap(indices []int, prefix int) int {
	var overlap int
	for _, index := range indices {
		if index < prefix {
			overlap++
		}
	}
	return overlap
}

func RunUP27() (TaskSelectedCommutantProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	unitaryStep, err := conjugatedLatentStep(
		mixer, applyStressUnitary,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	controlStep, err := conjugatedLatentStep(
		mixer, applyStressNonUnitary,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}

	candidates, _, err := discoverCommutingObservables(
		unitaryStep,
		taskSelectedCandidateCount,
		taskSelectedRounds,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	prefix64 := candidates[:taskSelectedRuntimeCount]

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	selectionStates, selectionTables, err := taskSelectedTrainingStates(
		trainTables, mixer, memoryNoise, 2,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	selectedIndices, selected64, err := selectTaskRelevantObservables(
		selectionStates,
		selectionTables,
		candidates,
		taskSelectedRuntimeCount,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}

	unitaryOps, err := latentDepthOperators(unitaryStep, allDepths)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	controlOps, err := latentDepthOperators(controlStep, heldDepths)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}

	prefixTrainRows, prefixTrainDrift, err := buildBreadthRows(
		trainTables, trainDepths, unitaryOps,
		mixer, prefix64, memoryNoise, 2, 0,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	prefixRegs, prefixHeads, prefixBase, err := trainBreadthModel(
		prefixTrainRows,
		taskSelectedRuntimeCount,
		"unitary_prefix_64_control",
		"unitary_prefix_64_reproduction",
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	prefixBase.MaxNormDrift = prefixTrainDrift
	prefixHeldRows, prefixHeldDrift, err := buildBreadthRows(
		heldTables, heldDepths, unitaryOps,
		mixer, prefix64, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	prefixResult, err := evaluateBreadthModel(
		prefixBase, prefixRegs, prefixHeads,
		prefixHeldRows, prefixHeldDrift,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}

	selectedTrainRows, selectedTrainDrift, err := buildBreadthRows(
		trainTables, trainDepths, unitaryOps,
		mixer, selected64, memoryNoise, 2, 0,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	selectedRegs, selectedHeads, selectedBase, err := trainBreadthModel(
		selectedTrainRows,
		taskSelectedRuntimeCount,
		"unitary_task_selected_64",
		"unitary_training_selected_64",
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	selectedBase.MaxNormDrift = selectedTrainDrift

	selectedHeldRows, selectedHeldDrift, err := buildBreadthRows(
		heldTables, heldDepths, unitaryOps,
		mixer, selected64, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	selectedResult, err := evaluateBreadthModel(
		selectedBase, selectedRegs, selectedHeads,
		selectedHeldRows, selectedHeldDrift,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}

	controlRows, controlDrift, err := buildBreadthRows(
		heldTables, heldDepths, controlOps,
		mixer, selected64, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	controlResult, err := evaluateBreadthModel(
		DiscoveryBreadthArm{
			Name:                       "non_unitary_task_selected_64",
			Path:                       "non_unitary_same_selected_model",
			ObservableCount:            taskSelectedRuntimeCount,
			RawFeatureDimension:        breadthRawFeatureDim(taskSelectedRuntimeCount),
			QuadraticFeatureDimension:  breadthQuadraticFeatureDim(taskSelectedRuntimeCount),
			TrainAccuracy:              selectedBase.TrainAccuracy,
			PerEntityTrain:             append([]float64(nil), selectedBase.PerEntityTrain...),
		},
		selectedRegs, selectedHeads,
		controlRows, controlDrift,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}

	integration, err := runBreadthIntegration(
		mixer, selected64, unitaryOps,
		selectedRegs, selectedHeads,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}

	commutatorError, err := maxWeylCommutatorEntry(
		selected64, unitaryStep,
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}
	featureDrift, err := breadthFeatureDrift(
		mixer, selected64, unitaryOps[1024],
	)
	if err != nil {
		return TaskSelectedCommutantProbeResult{}, err
	}

	prefixDelta := math.Abs(
		prefixResult.HeldOutAccuracy - 0.998779296875,
	)
	prefixPass := prefixDelta <= 1e-12
	selectionTrainingOnlyPass := true
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

	return TaskSelectedCommutantProbeResult{
		Schema:                           TaskSelectedCommutantSchema,
		Experiment:                       "UP-27-training-selected-discovered-commutant",
		LatentDimension:                  fullLatentDimension,
		RuntimeStateObjects:              1,
		VisibleChannelBlocks:             false,
		FullCoordinateMixing:             true,
		ObservableDiscoveryFromTransport: true,
		DiscoveryUsesHiddenMultiplicity:  false,
		DiscoveryUsesHiddenMixer:         false,
		SelectionUsesTrainingLabels:      true,
		SelectionUsesHeldOutData:         false,
		SelectionUsesExplicitDepth:       false,
		CandidateObservableCount:         taskSelectedCandidateCount,
		RuntimeObservableCount:           taskSelectedRuntimeCount,
		ProjectionRounds:                 taskSelectedRounds,
		EffectiveOrbitSize:               1 << taskSelectedRounds,
		SelectedIndices:                  append([]int(nil), selectedIndices...),
		RidgeLambda:                      taskSelectedLambda,
		TrainingDepths:                   append([]int(nil), trainDepths...),
		HeldOutDepths:                    append([]int(nil), heldDepths...),
		Prefix64:                         prefixResult,
		Selected64:                       selectedResult,
		MatchedControl64:                 controlResult,
		Integration64:                    integration,
		Diagnosis: TaskSelectedDiagnosis{
			Prefix64ReproductionPass:  prefixPass,
			SelectionTrainingOnlyPass: selectionTrainingOnlyPass,
			DiscoveryCommutatorPass:   commutatorPass,
			FeatureInvariancePass:     invariancePass,
			PhaseCodeLearningPass:     learningPass,
			UnseenDepthPass:           unseenPass,
			MutableIntegrationPass:    mutablePass,
			Prefix64HeldOutDelta:      prefixDelta,
			SelectedOverlapWithPrefix64: selectedPrefixOverlap(
				selectedIndices, taskSelectedRuntimeCount,
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
