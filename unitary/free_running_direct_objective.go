package unitary

import (
	"fmt"
	"math"
)

const FreeRunningDirectObjectiveSchema = "wingless.free-running-direct-objective.v1"

const freeRunningSoftWriteBoundary =
	"decode_probability_state_commanded_write_renormalize"

type FreeRunningDirectDiagnosis struct {
	InnerSplitBalancedPass      bool    `json:"inner_split_balanced_pass"`
	SelectionUsesTrueHeldOut    bool    `json:"selection_uses_true_heldout"`
	OptimizationSteps           int     `json:"optimization_steps"`
	SelectedStep                int     `json:"selected_step"`
	StartCapacity               int     `json:"start_capacity"`
	SelectedCapacity            int     `json:"selected_capacity"`
	MaxCapacity                 int     `json:"max_capacity"`
	ExactFusionOccurred         bool    `json:"exact_fusion_occurred"`
	SelectedHeldOutAccuracy     float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy      float64 `json:"selected_commit_accuracy"`
	SelectedFinalAccuracy       float64 `json:"selected_final_accuracy"`
	SelectedRelationAccuracy    float64 `json:"selected_relation_accuracy"`
	FullCapacityHeldOutAccuracy float64 `json:"full_capacity_heldout_accuracy"`
	FullCapacityCommitAccuracy  float64 `json:"full_capacity_commit_accuracy"`
	HeldOutRetentionDelta       float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta        float64 `json:"commit_retention_delta"`
	FreeRunningDirectPass       bool    `json:"free_running_direct_pass"`
}

type FreeRunningDirectObjectiveProbeResult struct {
	Schema                           string                      `json:"schema"`
	Experiment                       string                      `json:"experiment"`
	LatentDimension                  int                         `json:"latent_dimension"`
	RuntimeStateObjects              int                         `json:"runtime_state_objects"`
	FullCoordinateMixing             bool                        `json:"full_coordinate_mixing"`
	StartsFromIndependentOffsets     bool                        `json:"starts_from_independent_offsets"`
	FinishedPartitionMenuProvided    bool                        `json:"finished_partition_menu_provided"`
	HardMergeCandidatesEvaluated     bool                        `json:"hard_merge_candidates_evaluated"`
	PairAffinityFieldUsed            bool                        `json:"pair_affinity_field_used"`
	DirectOffsetsOptimized           bool                        `json:"direct_offsets_optimized"`
	Optimizer                        string                      `json:"optimizer"`
	TeacherForcedSmoothMutable       bool                        `json:"teacher_forced_smooth_mutable"`
	FreeRunningProbabilisticMutable  bool                        `json:"free_running_probabilistic_mutable"`
	ArgmaxFeedbackUsed               bool                        `json:"argmax_feedback_used"`
	SoftMemoryDimension              int                         `json:"soft_memory_dimension"`
	SoftWriteBoundary                string                      `json:"soft_write_boundary"`
	SmoothRelationProbability        bool                        `json:"smooth_relation_probability"`
	HarmonicBottleneckObjective      bool                        `json:"harmonic_bottleneck_objective"`
	StickyExactFusionProjection      bool                        `json:"sticky_exact_fusion_projection"`
	SelectionUsesHeldOutData         bool                        `json:"selection_uses_heldout_data"`
	ResourcePrice                    float64                     `json:"resource_price"`
	SoftCapacityTau                  float64                     `json:"soft_capacity_tau"`
	Perturbation                     float64                     `json:"perturbation"`
	LearningRate                     float64                     `json:"learning_rate"`
	MaximumCoordinateUpdate          float64                     `json:"maximum_coordinate_update"`
	MaximumSteps                     int                         `json:"maximum_steps"`
	FitTables                        int                         `json:"fit_tables"`
	ValidationTables                 int                         `json:"validation_tables"`
	TrueHeldOutTables                int                         `json:"true_heldout_tables"`
	Initial                          RichDirectEvaluation        `json:"initial"`
	Steps                            []RichDirectStep             `json:"steps"`
	SelectedStep                     RichDirectStep               `json:"selected_step"`
	SelectedFinalArm                 MultiplicityDoseArm          `json:"selected_final_arm"`
	FullCapacityControl              MultiplicityDoseArm          `json:"full_capacity_control"`
	Diagnosis                        FreeRunningDirectDiagnosis   `json:"diagnosis"`
}

func softMemoryAfterCommandedWrite(
	distributions [4][]float64,
	entity, value int,
) (State, error) {
	if entity < 0 || entity >= 4 || value < 0 || value >= 4 {
		return nil, fmt.Errorf("soft memory write out of range")
	}
	state := make(State, 16)
	for currentEntity := 0; currentEntity < 4; currentEntity++ {
		if len(distributions[currentEntity]) != 4 {
			return nil, fmt.Errorf(
				"soft memory distribution dimension entity=%d got=%d",
				currentEntity, len(distributions[currentEntity]),
			)
		}
		var sum float64
		for currentValue, probability :=
			range distributions[currentEntity] {
			if !finite(probability) || probability < 0 {
				return nil, fmt.Errorf(
					"soft memory probability invalid entity=%d value=%d",
					currentEntity, currentValue,
				)
			}
			sum += probability
			state[currentEntity*4+currentValue] =
				complex(0.5*probability, 0)
		}
		if math.Abs(sum-1.0) > 1e-6 {
			return nil, fmt.Errorf(
				"soft memory probabilities not normalized entity=%d sum=%g",
				currentEntity, sum,
			)
		}
	}

	start := entity * 4
	for currentValue := 0; currentValue < 4; currentValue++ {
		state[start+currentValue] = 0
	}
	state[start+value] = complex(0.5, 0)

	return Normalize(state)
}

func freeRunningProbabilisticSignals(
	mixer latentMatrix,
	observables []latentMatrix,
	unitaryOps map[int]latentMatrix,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
	tables []memoryTable,
	depths []int,
	memoryNoise float64,
) (float64, float64, error) {
	if len(tables) == 0 || len(depths) == 0 {
		return 0, 0, fmt.Errorf("free-running direct validation set empty")
	}

	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return 0, 0, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return 0, 0, err
	}

	scenarios := richDirectScenarioCount
	if len(tables) < scenarios {
		scenarios = len(tables)
	}

	var valueProbabilityTotal float64
	var valueProbabilityCount int
	var relationProbabilityTotal float64
	var relationProbabilityCount int

	for scenarioIndex := 0; scenarioIndex < scenarios; scenarioIndex++ {
		scenario := makeFullRankScenario(
			scenarioIndex,
			tables[(scenarioIndex*7)%len(tables)],
			richDirectWrites,
			depths,
		)
		trueTable := scenario.initial
		softMemory, err := encodeMemory(scenario.initial)
		if err != nil {
			return 0, 0, err
		}

		for writeIndex, write := range scenario.writes {
			seed := 39000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(softMemory, seed, memoryNoise)
			if err != nil {
				return 0, 0, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.271*float64(seed+1), 2*math.Pi),
			)
			state, err := fullLatentEncode(memory, mixer)
			if err != nil {
				return 0, 0, err
			}
			operator, ok := unitaryOps[write.gap]
			if !ok {
				return 0, 0, fmt.Errorf(
					"missing free-running mutable depth=%d", write.gap,
				)
			}
			state, err = latentMatrixVector(operator, state)
			if err != nil {
				return 0, 0, err
			}
			_, distributions, _, err := decodeBreadthTable(
				state, observables, regressors, classifiers,
			)
			if err != nil {
				return 0, 0, err
			}
			for currentEntity := 0; currentEntity < 4; currentEntity++ {
				valueProbabilityTotal +=
					distributions[currentEntity][trueTable[currentEntity]]
				valueProbabilityCount++
			}

			softMemory, err = softMemoryAfterCommandedWrite(
				distributions, write.entity, write.value,
			)
			if err != nil {
				return 0, 0, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return 0, 0, err
			}
		}

		seed := 39000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(softMemory, seed, memoryNoise)
		if err != nil {
			return 0, 0, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.271*float64(seed+1), 2*math.Pi),
		)
		state, err := fullLatentEncode(memory, mixer)
		if err != nil {
			return 0, 0, err
		}
		operator, ok := unitaryOps[scenario.finalGap]
		if !ok {
			return 0, 0, fmt.Errorf(
				"missing free-running final depth=%d", scenario.finalGap,
			)
		}
		state, err = latentMatrixVector(operator, state)
		if err != nil {
			return 0, 0, err
		}
		_, distributions, _, err := decodeBreadthTable(
			state, observables, regressors, classifiers,
		)
		if err != nil {
			return 0, 0, err
		}
		for currentEntity := 0; currentEntity < 4; currentEntity++ {
			valueProbabilityTotal +=
				distributions[currentEntity][trueTable[currentEntity]]
			valueProbabilityCount++
		}

		relationInput, err := relationFeatures(
			distributions[scenario.queryA],
			distributions[scenario.queryB],
		)
		if err != nil {
			return 0, 0, err
		}
		relationProbabilities, err :=
			relationHead.probabilities(relationInput)
		if err != nil {
			return 0, 0, err
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return 0, 0, err
		}
		relationProbabilityTotal +=
			relationProbabilities[wantRelation]
		relationProbabilityCount++
	}

	if valueProbabilityCount == 0 || relationProbabilityCount == 0 {
		return 0, 0, fmt.Errorf(
			"free-running direct objective produced no samples",
		)
	}
	return valueProbabilityTotal / float64(valueProbabilityCount),
		relationProbabilityTotal / float64(relationProbabilityCount),
		nil
}

func evaluateFreeRunningDirectOffsets(
	name string,
	offsets []float64,
	mixer latentMatrix,
	fitTables, validationTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (RichDirectEvaluation, error) {
	step, err := fullLatentDoseStep(mixer, offsets)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	realStep, err := realifyMatrix(step)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	orthError, err := maxRealOrthogonalityError(realStep)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	ops, err := latentDepthOperators(step, allDepths)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	candidates, _, err := discoverCommutingObservables(
		step, interactionCandidateCount, interactionRounds,
	)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	selectionStates, selectionTables, err :=
		taskSelectedTrainingStates(
			fitTables, mixer, memoryNoise, 2,
		)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	selectedIndices, selectedObservables, err :=
		selectInteractionRelevantObservables(
			selectionStates, selectionTables,
			candidates, interactionRuntimeCount,
		)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	commutatorError, err := maxWeylCommutatorEntry(
		selectedObservables, step,
	)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	featureDrift, err := breadthFeatureDrift(
		mixer, selectedObservables, ops[1024],
	)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	regressors, classifiers, static, err :=
		trainAndEvaluateMultiplicityArm(
			name,
			"free_running_"+name,
			fitTables, validationTables,
			trainDepths, heldDepths,
			ops, mixer, selectedObservables, memoryNoise,
		)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	valueProb, relationProb, err :=
		freeRunningProbabilisticSignals(
			mixer, selectedObservables, ops,
			regressors, classifiers,
			validationTables, heldDepths, memoryNoise,
		)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	phaseScore := 0.5 * (static.MeanPhaseCosine + 1.0)
	if phaseScore < 0 {
		phaseScore = 0
	}
	if phaseScore > 1 {
		phaseScore = 1
	}
	taskScore := harmonicMean3(
		phaseScore, valueProb, relationProb,
	)
	softCapacity, err := directSoftCapacity(offsets)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	resourcePenalty :=
		directOffsetResourcePrice * softCapacity / 36.0
	objective := taskScore - resourcePenalty

	return RichDirectEvaluation{
		Name:                    name,
		Offsets:                 append([]float64(nil), offsets...),
		SoftCapacity:            softCapacity,
		ResourcePenalty:         resourcePenalty,
		Signals: RichSmoothSignals{
			NormalizedPhaseScore:            phaseScore,
			MeanCorrectValueProbability:     valueProb,
			MeanCorrectRelationProbability:  relationProb,
			HarmonicTaskScore:               taskScore,
		},
		Objective:                 objective,
		OrthogonalityError:       orthError,
		DiscoveryCommutatorError: commutatorError,
		FeatureDrift:             featureDrift,
		Static:                   static,
		SelectedIndices:          append([]int(nil), selectedIndices...),
	}, nil
}

func RunUP40() (FreeRunningDirectObjectiveProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fitTables, validationTables :=
		splitTaskAllocationTrainingPool(allTrain)

	minFit := minimumTableMarginalCount(fitTables)
	minValidation := minimumTableMarginalCount(validationTables)
	innerBalanced := minFit >= 12 && minValidation >= 4

	union := newFusionUnion()
	currentOffsets := continuousFusionInitialOffsets()
	initial, err := evaluateFreeRunningDirectOffsets(
		"free_running_direct_initial",
		currentOffsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return FreeRunningDirectObjectiveProbeResult{}, err
	}
	if !richEvaluationStructurallyValid(initial) {
		return FreeRunningDirectObjectiveProbeResult{}, fmt.Errorf(
			"free-running direct initial evaluation structurally invalid",
		)
	}

	selectedEval := initial
	selectedStepNumber := 0
	var selectedStep RichDirectStep
	steps := make([]RichDirectStep, 0, richDirectSteps)

	for step := 1; step <= richDirectSteps; step++ {
		inputGroups := union.groups()
		inputCapacity := fusionCapacity(inputGroups)
		direction := deterministicSPSADirection(step, inputGroups)

		plusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			directOffsetPerturbation, &union,
		)
		if err != nil {
			return FreeRunningDirectObjectiveProbeResult{}, err
		}
		minusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			-directOffsetPerturbation, &union,
		)
		if err != nil {
			return FreeRunningDirectObjectiveProbeResult{}, err
		}

		plus, err := evaluateFreeRunningDirectOffsets(
			fmt.Sprintf("free_running_plus_%02d", step),
			plusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return FreeRunningDirectObjectiveProbeResult{}, err
		}
		minus, err := evaluateFreeRunningDirectOffsets(
			fmt.Sprintf("free_running_minus_%02d", step),
			minusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return FreeRunningDirectObjectiveProbeResult{}, err
		}
		if !richEvaluationStructurallyValid(plus) ||
			!richEvaluationStructurallyValid(minus) {
			return FreeRunningDirectObjectiveProbeResult{}, fmt.Errorf(
				"free-running direct perturbation structurally invalid step=%d",
				step,
			)
		}

		gradient := directGradientEstimate(
			plus.Objective, minus.Objective, direction,
		)
		nextOffsets, update, err := applyDirectGradient(
			currentOffsets, gradient, &union,
		)
		if err != nil {
			return FreeRunningDirectObjectiveProbeResult{}, err
		}
		outputGroups := union.groups()
		outputCapacity := fusionCapacity(outputGroups)

		output, err := evaluateFreeRunningDirectOffsets(
			fmt.Sprintf("free_running_output_%02d", step),
			nextOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return FreeRunningDirectObjectiveProbeResult{}, err
		}
		if !richEvaluationStructurallyValid(output) {
			return FreeRunningDirectObjectiveProbeResult{}, fmt.Errorf(
				"free-running direct output structurally invalid step=%d",
				step,
			)
		}

		entry := RichDirectStep{
			Step:             step,
			InputOffsets:     append([]float64(nil), currentOffsets...),
			InputGroups:      inputGroups,
			InputCapacity:    inputCapacity,
			Direction:        append([]float64(nil), direction...),
			Plus:             plus,
			Minus:            minus,
			GradientEstimate: gradient,
			Update:           update,
			OutputOffsets:    append([]float64(nil), nextOffsets...),
			OutputGroups:     outputGroups,
			OutputCapacity:   outputCapacity,
			Output:           output,
		}
		steps = append(steps, entry)

		if output.Objective > selectedEval.Objective ||
			(output.Objective == selectedEval.Objective &&
				output.SoftCapacity < selectedEval.SoftCapacity) {
			selectedEval = output
			selectedStep = entry
			selectedStepNumber = step
		}
		currentOffsets = nextOffsets
	}

	if selectedStepNumber == 0 {
		selectedStep = RichDirectStep{
			Step:           0,
			InputOffsets:   append([]float64(nil), initial.Offsets...),
			InputGroups:    initialSymmetryGroupsJSON(),
			InputCapacity:  6,
			OutputOffsets:  append([]float64(nil), initial.Offsets...),
			OutputGroups:   initialSymmetryGroupsJSON(),
			OutputCapacity: 6,
			Output:         initial,
		}
	}

	// Final scientific evaluation is intentionally unchanged: it uses the
	// actual hard recurrent mutable loop with decoded states fed back.
	selectedFinal, err := evaluateMultiplicityDoseArm(
		fmt.Sprintf(
			"selected_free_running_direct_step_%02d",
			selectedStepNumber,
		),
		selectedEval.Offsets,
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return FreeRunningDirectObjectiveProbeResult{}, err
	}
	fullControl, err := evaluateMultiplicityDoseArm(
		"full_capacity_control",
		[]float64{0, 0, 0, 0, 0, 0},
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return FreeRunningDirectObjectiveProbeResult{}, err
	}

	selectedCapacity := fusionCapacity(selectedStep.OutputGroups)
	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedFinal.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedFinal.Integration.CommitDecodeAccuracy
	exactFusion := selectedCapacity > 6

	pass :=
		innerBalanced &&
			selectedStepNumber >= 1 &&
			exactFusion &&
			selectedCapacity < 36 &&
			selectedFinal.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedFinal.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedFinal.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedFinal.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	return FreeRunningDirectObjectiveProbeResult{
		Schema:                          FreeRunningDirectObjectiveSchema,
		Experiment:                      "UP-40-free-running-probabilistic-direct-objective",
		LatentDimension:                 fullLatentDimension,
		RuntimeStateObjects:             1,
		FullCoordinateMixing:            true,
		StartsFromIndependentOffsets:    true,
		FinishedPartitionMenuProvided:   false,
		HardMergeCandidatesEvaluated:    false,
		PairAffinityFieldUsed:           false,
		DirectOffsetsOptimized:          true,
		Optimizer:                       "deterministic-projected-central-difference",
		TeacherForcedSmoothMutable:      false,
		FreeRunningProbabilisticMutable: true,
		ArgmaxFeedbackUsed:              false,
		SoftMemoryDimension:             16,
		SoftWriteBoundary:               freeRunningSoftWriteBoundary,
		SmoothRelationProbability:       true,
		HarmonicBottleneckObjective:     true,
		StickyExactFusionProjection:     true,
		SelectionUsesHeldOutData:        false,
		ResourcePrice:                   directOffsetResourcePrice,
		SoftCapacityTau:                 directOffsetSoftCapacityTau,
		Perturbation:                    directOffsetPerturbation,
		LearningRate:                    directOffsetLearningRate,
		MaximumCoordinateUpdate:         directOffsetMaxUpdate,
		MaximumSteps:                    richDirectSteps,
		FitTables:                       len(fitTables),
		ValidationTables:                len(validationTables),
		TrueHeldOutTables:               len(trueHeld),
		Initial:                         initial,
		Steps:                           steps,
		SelectedStep:                    selectedStep,
		SelectedFinalArm:                selectedFinal,
		FullCapacityControl:             fullControl,
		Diagnosis: FreeRunningDirectDiagnosis{
			InnerSplitBalancedPass:      innerBalanced,
			SelectionUsesTrueHeldOut:    false,
			OptimizationSteps:           len(steps),
			SelectedStep:                selectedStepNumber,
			StartCapacity:               6,
			SelectedCapacity:            selectedCapacity,
			MaxCapacity:                 36,
			ExactFusionOccurred:         exactFusion,
			SelectedHeldOutAccuracy:     selectedFinal.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:      selectedFinal.Integration.CommitDecodeAccuracy,
			SelectedFinalAccuracy:       selectedFinal.Integration.ExactFinalTableAccuracy,
			SelectedRelationAccuracy:    selectedFinal.Integration.RelationalQueryAccuracy,
			FullCapacityHeldOutAccuracy: fullControl.Static.HeldOutAccuracy,
			FullCapacityCommitAccuracy:  fullControl.Integration.CommitDecodeAccuracy,
			HeldOutRetentionDelta:       heldDelta,
			CommitRetentionDelta:        commitDelta,
			FreeRunningDirectPass:       pass,
		},
	}, nil
}
