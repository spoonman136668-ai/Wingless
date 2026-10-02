package unitary

import (
	"fmt"
	"math"
)

const RichDirectObjectiveSchema = "wingless.rich-direct-objective.v1"

const (
	richDirectSteps          = directOffsetSteps
	richDirectWrites         = 8
	richDirectScenarioCount  = 32
	richDirectEpsilon        = 1e-9
)

type RichSmoothSignals struct {
	NormalizedPhaseScore       float64 `json:"normalized_phase_score"`
	MeanCorrectValueProbability float64 `json:"mean_correct_value_probability"`
	MeanCorrectRelationProbability float64 `json:"mean_correct_relation_probability"`
	HarmonicTaskScore          float64 `json:"harmonic_task_score"`
}

type RichDirectEvaluation struct {
	Name                    string              `json:"name"`
	Offsets                 []float64           `json:"offsets"`
	SoftCapacity            float64             `json:"soft_capacity"`
	ResourcePenalty         float64             `json:"resource_penalty"`
	Signals                 RichSmoothSignals   `json:"signals"`
	Objective               float64             `json:"objective"`
	OrthogonalityError      float64             `json:"orthogonality_error"`
	DiscoveryCommutatorError float64            `json:"discovery_commutator_error"`
	FeatureDrift            float64             `json:"feature_drift"`
	Static                  DiscoveryBreadthArm `json:"static"`
	SelectedIndices         []int               `json:"selected_indices"`
}

type RichDirectStep struct {
	Step             int                  `json:"step"`
	InputOffsets     []float64            `json:"input_offsets"`
	InputGroups      [][]int              `json:"input_groups"`
	InputCapacity    int                  `json:"input_capacity"`
	Direction        []float64            `json:"direction"`
	Plus             RichDirectEvaluation `json:"plus"`
	Minus            RichDirectEvaluation `json:"minus"`
	GradientEstimate []float64            `json:"gradient_estimate"`
	Update           []float64            `json:"update"`
	OutputOffsets    []float64            `json:"output_offsets"`
	OutputGroups     [][]int              `json:"output_groups"`
	OutputCapacity   int                  `json:"output_capacity"`
	Output           RichDirectEvaluation `json:"output"`
}

type RichDirectDiagnosis struct {
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
	RichDirectPass              bool    `json:"rich_direct_pass"`
}

type RichDirectObjectiveProbeResult struct {
	Schema                        string               `json:"schema"`
	Experiment                    string               `json:"experiment"`
	LatentDimension               int                  `json:"latent_dimension"`
	RuntimeStateObjects           int                  `json:"runtime_state_objects"`
	FullCoordinateMixing          bool                 `json:"full_coordinate_mixing"`
	StartsFromIndependentOffsets  bool                 `json:"starts_from_independent_offsets"`
	FinishedPartitionMenuProvided bool                 `json:"finished_partition_menu_provided"`
	HardMergeCandidatesEvaluated  bool                 `json:"hard_merge_candidates_evaluated"`
	PairAffinityFieldUsed         bool                 `json:"pair_affinity_field_used"`
	DirectOffsetsOptimized        bool                 `json:"direct_offsets_optimized"`
	Optimizer                     string               `json:"optimizer"`
	TeacherForcedSmoothMutable    bool                 `json:"teacher_forced_smooth_mutable"`
	SmoothRelationProbability     bool                 `json:"smooth_relation_probability"`
	HarmonicBottleneckObjective   bool                 `json:"harmonic_bottleneck_objective"`
	StickyExactFusionProjection   bool                 `json:"sticky_exact_fusion_projection"`
	SelectionUsesHeldOutData      bool                 `json:"selection_uses_heldout_data"`
	ResourcePrice                 float64              `json:"resource_price"`
	SoftCapacityTau               float64              `json:"soft_capacity_tau"`
	Perturbation                  float64              `json:"perturbation"`
	LearningRate                  float64              `json:"learning_rate"`
	MaximumCoordinateUpdate       float64              `json:"maximum_coordinate_update"`
	MaximumSteps                  int                  `json:"maximum_steps"`
	FitTables                     int                  `json:"fit_tables"`
	ValidationTables              int                  `json:"validation_tables"`
	TrueHeldOutTables             int                  `json:"true_heldout_tables"`
	Initial                       RichDirectEvaluation `json:"initial"`
	Steps                         []RichDirectStep     `json:"steps"`
	SelectedStep                  RichDirectStep       `json:"selected_step"`
	SelectedFinalArm              MultiplicityDoseArm  `json:"selected_final_arm"`
	FullCapacityControl           MultiplicityDoseArm  `json:"full_capacity_control"`
	Diagnosis                     RichDirectDiagnosis  `json:"diagnosis"`
}

func harmonicMean3(a, b, c float64) float64 {
	if a <= 0 || b <= 0 || c <= 0 {
		return 0
	}
	return 3.0 / (
		1.0/math.Max(a, richDirectEpsilon) +
			1.0/math.Max(b, richDirectEpsilon) +
			1.0/math.Max(c, richDirectEpsilon))
}

func richTeacherForcedSignals(
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
		return 0, 0, fmt.Errorf("rich direct validation set empty")
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

		for writeIndex, write := range scenario.writes {
			canonical, err := encodeMemory(trueTable)
			if err != nil {
				return 0, 0, err
			}
			seed := 39000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
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
					"missing rich mutable depth=%d", write.gap,
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
			for entity := 0; entity < 4; entity++ {
				valueProbabilityTotal +=
					distributions[entity][trueTable[entity]]
				valueProbabilityCount++
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return 0, 0, err
			}
		}

		canonical, err := encodeMemory(trueTable)
		if err != nil {
			return 0, 0, err
		}
		seed := 39000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
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
				"missing rich final depth=%d", scenario.finalGap,
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
		for entity := 0; entity < 4; entity++ {
			valueProbabilityTotal +=
				distributions[entity][trueTable[entity]]
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
		return 0, 0, fmt.Errorf("rich direct objective produced no samples")
	}
	return valueProbabilityTotal / float64(valueProbabilityCount),
		relationProbabilityTotal / float64(relationProbabilityCount),
		nil
}

func evaluateRichDirectOffsets(
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
			"rich_"+name,
			fitTables, validationTables,
			trainDepths, heldDepths,
			ops, mixer, selectedObservables, memoryNoise,
		)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	valueProb, relationProb, err :=
		richTeacherForcedSignals(
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
			NormalizedPhaseScore:          phaseScore,
			MeanCorrectValueProbability:   valueProb,
			MeanCorrectRelationProbability: relationProb,
			HarmonicTaskScore:             taskScore,
		},
		Objective:                objective,
		OrthogonalityError:      orthError,
		DiscoveryCommutatorError: commutatorError,
		FeatureDrift:            featureDrift,
		Static:                  static,
		SelectedIndices:         append([]int(nil), selectedIndices...),
	}, nil
}

func richEvaluationStructurallyValid(
	evaluation RichDirectEvaluation,
) bool {
	return evaluation.OrthogonalityError <= 1e-10 &&
		evaluation.DiscoveryCommutatorError <= 1e-5 &&
		evaluation.FeatureDrift <= 5e-3 &&
		finite(evaluation.Objective)
}

func RunUP39() (RichDirectObjectiveProbeResult, error) {
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
	initial, err := evaluateRichDirectOffsets(
		"rich_direct_initial",
		currentOffsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return RichDirectObjectiveProbeResult{}, err
	}
	if !richEvaluationStructurallyValid(initial) {
		return RichDirectObjectiveProbeResult{}, fmt.Errorf(
			"rich direct initial evaluation structurally invalid",
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
			return RichDirectObjectiveProbeResult{}, err
		}
		minusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			-directOffsetPerturbation, &union,
		)
		if err != nil {
			return RichDirectObjectiveProbeResult{}, err
		}

		plus, err := evaluateRichDirectOffsets(
			fmt.Sprintf("rich_plus_%02d", step),
			plusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return RichDirectObjectiveProbeResult{}, err
		}
		minus, err := evaluateRichDirectOffsets(
			fmt.Sprintf("rich_minus_%02d", step),
			minusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return RichDirectObjectiveProbeResult{}, err
		}
		if !richEvaluationStructurallyValid(plus) ||
			!richEvaluationStructurallyValid(minus) {
			return RichDirectObjectiveProbeResult{}, fmt.Errorf(
				"rich direct perturbation structurally invalid step=%d",
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
			return RichDirectObjectiveProbeResult{}, err
		}
		outputGroups := union.groups()
		outputCapacity := fusionCapacity(outputGroups)

		output, err := evaluateRichDirectOffsets(
			fmt.Sprintf("rich_output_%02d", step),
			nextOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return RichDirectObjectiveProbeResult{}, err
		}
		if !richEvaluationStructurallyValid(output) {
			return RichDirectObjectiveProbeResult{}, fmt.Errorf(
				"rich direct output structurally invalid step=%d",
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

	selectedFinal, err := evaluateMultiplicityDoseArm(
		fmt.Sprintf("selected_rich_direct_step_%02d", selectedStepNumber),
		selectedEval.Offsets,
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return RichDirectObjectiveProbeResult{}, err
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
		return RichDirectObjectiveProbeResult{}, err
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

	return RichDirectObjectiveProbeResult{
		Schema:                        RichDirectObjectiveSchema,
		Experiment:                    "UP-39-rich-task-direct-spectral-offset-optimization",
		LatentDimension:               fullLatentDimension,
		RuntimeStateObjects:           1,
		FullCoordinateMixing:          true,
		StartsFromIndependentOffsets:  true,
		FinishedPartitionMenuProvided: false,
		HardMergeCandidatesEvaluated:  false,
		PairAffinityFieldUsed:         false,
		DirectOffsetsOptimized:        true,
		Optimizer:                     "deterministic-projected-central-difference",
		TeacherForcedSmoothMutable:    true,
		SmoothRelationProbability:     true,
		HarmonicBottleneckObjective:   true,
		StickyExactFusionProjection:   true,
		SelectionUsesHeldOutData:      false,
		ResourcePrice:                 directOffsetResourcePrice,
		SoftCapacityTau:               directOffsetSoftCapacityTau,
		Perturbation:                  directOffsetPerturbation,
		LearningRate:                  directOffsetLearningRate,
		MaximumCoordinateUpdate:       directOffsetMaxUpdate,
		MaximumSteps:                  richDirectSteps,
		FitTables:                     len(fitTables),
		ValidationTables:              len(validationTables),
		TrueHeldOutTables:             len(trueHeld),
		Initial:                       initial,
		Steps:                         steps,
		SelectedStep:                  selectedStep,
		SelectedFinalArm:              selectedFinal,
		FullCapacityControl:           fullControl,
		Diagnosis: RichDirectDiagnosis{
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
			RichDirectPass:              pass,
		},
	}, nil
}
