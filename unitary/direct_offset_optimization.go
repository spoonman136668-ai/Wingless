package unitary

import (
	"fmt"
	"math"
)

const DirectOffsetOptimizationSchema = "wingless.direct-offset-optimization.v1"

const (
	directOffsetSteps            = 12
	directOffsetPerturbation      = 0.002
	directOffsetLearningRate      = 0.002
	directOffsetMaxUpdate         = 0.004
	directOffsetSoftCapacityTau   = continuousFusionTolerance
	directOffsetResourcePrice     = taskAllocationCapacityPrice
)

type DirectOffsetEvaluation struct {
	Name             string              `json:"name"`
	Offsets          []float64           `json:"offsets"`
	SoftCapacity     float64             `json:"soft_capacity"`
	ResourcePenalty float64             `json:"resource_penalty"`
	TaskScore        float64             `json:"task_score"`
	Objective        float64             `json:"objective"`
	Arm              MultiplicityDoseArm `json:"arm"`
}

type DirectOffsetStep struct {
	Step             int                    `json:"step"`
	InputOffsets     []float64              `json:"input_offsets"`
	InputGroups      [][]int                `json:"input_groups"`
	InputCapacity    int                    `json:"input_capacity"`
	Direction        []float64              `json:"direction"`
	Plus             DirectOffsetEvaluation `json:"plus"`
	Minus            DirectOffsetEvaluation `json:"minus"`
	GradientEstimate []float64              `json:"gradient_estimate"`
	Update           []float64              `json:"update"`
	OutputOffsets    []float64              `json:"output_offsets"`
	OutputGroups     [][]int                `json:"output_groups"`
	OutputCapacity   int                    `json:"output_capacity"`
	Output           DirectOffsetEvaluation `json:"output"`
}

type DirectOffsetDiagnosis struct {
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
	FullCapacityHeldOutAccuracy float64 `json:"full_capacity_heldout_accuracy"`
	FullCapacityCommitAccuracy  float64 `json:"full_capacity_commit_accuracy"`
	HeldOutRetentionDelta       float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta        float64 `json:"commit_retention_delta"`
	DirectOptimizationPass      bool    `json:"direct_optimization_pass"`
}

type DirectOffsetOptimizationProbeResult struct {
	Schema                         string                 `json:"schema"`
	Experiment                     string                 `json:"experiment"`
	LatentDimension                int                    `json:"latent_dimension"`
	RuntimeStateObjects            int                    `json:"runtime_state_objects"`
	FullCoordinateMixing           bool                   `json:"full_coordinate_mixing"`
	StartsFromIndependentOffsets   bool                   `json:"starts_from_independent_offsets"`
	FinishedPartitionMenuProvided  bool                   `json:"finished_partition_menu_provided"`
	HardMergeCandidatesEvaluated   bool                   `json:"hard_merge_candidates_evaluated"`
	PairAffinityFieldUsed          bool                   `json:"pair_affinity_field_used"`
	DirectOffsetsOptimized         bool                   `json:"direct_offsets_optimized"`
	Optimizer                      string                 `json:"optimizer"`
	SmoothTaskObjective            bool                   `json:"smooth_task_objective"`
	SmoothResourceObjective        bool                   `json:"smooth_resource_objective"`
	StickyExactFusionProjection    bool                   `json:"sticky_exact_fusion_projection"`
	SelectionUsesHeldOutData       bool                   `json:"selection_uses_heldout_data"`
	ResourcePrice                  float64                `json:"resource_price"`
	SoftCapacityTau                float64                `json:"soft_capacity_tau"`
	Perturbation                   float64                `json:"perturbation"`
	LearningRate                   float64                `json:"learning_rate"`
	MaximumCoordinateUpdate        float64                `json:"maximum_coordinate_update"`
	MaximumSteps                   int                    `json:"maximum_steps"`
	FitTables                      int                    `json:"fit_tables"`
	ValidationTables               int                    `json:"validation_tables"`
	TrueHeldOutTables              int                    `json:"true_heldout_tables"`
	Initial                        DirectOffsetEvaluation `json:"initial"`
	Steps                          []DirectOffsetStep     `json:"steps"`
	SelectedStep                   DirectOffsetStep       `json:"selected_step"`
	SelectedFinalArm               MultiplicityDoseArm    `json:"selected_final_arm"`
	FullCapacityControl            MultiplicityDoseArm    `json:"full_capacity_control"`
	Diagnosis                      DirectOffsetDiagnosis  `json:"diagnosis"`
}

func directSoftCapacity(offsets []float64) (float64, error) {
	if len(offsets) != compositeChannels {
		return 0, fmt.Errorf("direct soft capacity dimension mismatch")
	}
	capacity := float64(compositeChannels)
	tau2 := directOffsetSoftCapacityTau * directOffsetSoftCapacityTau
	if tau2 <= 0 {
		return 0, fmt.Errorf("direct soft capacity tau invalid")
	}
	for i := 0; i < compositeChannels; i++ {
		for j := i + 1; j < compositeChannels; j++ {
			diff := offsets[i] - offsets[j]
			capacity += 2 * math.Exp(
				-(diff*diff)/(2*tau2),
			)
		}
	}
	return capacity, nil
}

func directObjective(
	arm MultiplicityDoseArm,
	offsets []float64,
) (float64, float64, float64, error) {
	softCapacity, err := directSoftCapacity(offsets)
	if err != nil {
		return 0, 0, 0, err
	}
	resourcePenalty :=
		directOffsetResourcePrice * softCapacity / 36.0
	taskScore := arm.Static.MeanPhaseCosine
	objective := taskScore - resourcePenalty
	if !finite(objective) {
		return 0, 0, 0, fmt.Errorf(
			"direct objective non-finite",
		)
	}
	return objective, taskScore, softCapacity, nil
}

func evaluateDirectOffsets(
	name string,
	offsets []float64,
	mixer latentMatrix,
	fitTables, validationTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (DirectOffsetEvaluation, error) {
	arm, err := evaluateMultiplicityDoseArm(
		name,
		offsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return DirectOffsetEvaluation{}, err
	}
	if !structurallyValidContinuousArm(arm) {
		return DirectOffsetEvaluation{}, fmt.Errorf(
			"direct offset evaluation %s structurally invalid",
			name,
		)
	}
	objective, taskScore, softCapacity, err :=
		directObjective(arm, offsets)
	if err != nil {
		return DirectOffsetEvaluation{}, err
	}
	return DirectOffsetEvaluation{
		Name:             name,
		Offsets:          append([]float64(nil), offsets...),
		SoftCapacity:     softCapacity,
		ResourcePenalty: directOffsetResourcePrice * softCapacity / 36.0,
		TaskScore:        taskScore,
		Objective:        objective,
		Arm:              arm,
	}, nil
}

func deterministicSPSADirection(
	step int,
	groups [][]int,
) []float64 {
	// Non-constant rows of an 8x8 Sylvester-Hadamard matrix, restricted
	// to the six available group coordinates. After zero-mean projection
	// these rows span the full five-dimensional tangent space when all six
	// groups are independent, and the full (k-1)-dimensional tangent space
	// for every k=2..6 prefix of active groups.
	patterns := [7][compositeChannels]float64{
		{1, -1, 1, -1, 1, -1},
		{1, 1, -1, -1, 1, 1},
		{1, -1, -1, 1, 1, -1},
		{1, 1, 1, 1, -1, -1},
		{1, -1, 1, -1, -1, 1},
		{1, 1, -1, -1, -1, -1},
		{1, -1, -1, 1, -1, 1},
	}
	row := patterns[(step-1+len(patterns))%len(patterns)]

	direction := make([]float64, compositeChannels)
	for groupIndex, group := range groups {
		sign := row[groupIndex]
		for _, member := range group {
			direction[member] = sign
		}
	}

	// Remove the member-weighted mean so the perturbation remains tangent
	// to the zero-mean manifold even after groups have fused.
	var mean float64
	for _, value := range direction {
		mean += value
	}
	mean /= float64(len(direction))
	for i := range direction {
		direction[i] -= mean
	}
	return direction
}

func perturbDirectOffsets(
	offsets, direction []float64,
	scale float64,
	union *fusionUnion,
) ([]float64, error) {
	if len(offsets) != compositeChannels ||
		len(direction) != compositeChannels {
		return nil, fmt.Errorf("direct perturbation dimension mismatch")
	}
	out := append([]float64(nil), offsets...)
	for i := range out {
		out[i] += scale * direction[i]
	}
	normalized, err := normalizeContinuousOffsets(out)
	if err != nil {
		return nil, err
	}
	return applyFusionGroups(normalized, union)
}

func directGradientEstimate(
	objectivePlus, objectiveMinus float64,
	direction []float64,
) []float64 {
	gradient := make([]float64, len(direction))
	denominator := 2 * directOffsetPerturbation
	delta := (objectivePlus - objectiveMinus) / denominator
	for i, sign := range direction {
		gradient[i] = delta * sign
	}
	return gradient
}

func clampDirectUpdate(value float64) float64 {
	if value > directOffsetMaxUpdate {
		return directOffsetMaxUpdate
	}
	if value < -directOffsetMaxUpdate {
		return -directOffsetMaxUpdate
	}
	return value
}

func applyDirectGradient(
	offsets, gradient []float64,
	union *fusionUnion,
) ([]float64, []float64, error) {
	if len(offsets) != compositeChannels ||
		len(gradient) != compositeChannels {
		return nil, nil, fmt.Errorf("direct update dimension mismatch")
	}
	update := make([]float64, compositeChannels)
	proposal := append([]float64(nil), offsets...)
	for i := range proposal {
		update[i] = clampDirectUpdate(
			directOffsetLearningRate * gradient[i],
		)
		proposal[i] += update[i]
	}
	normalized, err := normalizeContinuousOffsets(proposal)
	if err != nil {
		return nil, nil, err
	}
	normalized, err = applyFusionGroups(normalized, union)
	if err != nil {
		return nil, nil, err
	}

	groups := union.groups()
	values := make([]float64, len(groups))
	for i, group := range groups {
		value, err := fusionGroupValue(normalized, group)
		if err != nil {
			return nil, nil, err
		}
		values[i] = value
	}
	for i := 0; i < len(groups); i++ {
		for j := i + 1; j < len(groups); j++ {
			if math.Abs(values[i]-values[j]) <= continuousFusionTolerance {
				union.union(groups[i][0], groups[j][0])
			}
		}
	}
	normalized, err = applyFusionGroups(normalized, union)
	if err != nil {
		return nil, nil, err
	}
	return normalized, update, nil
}

func RunUP38() (DirectOffsetOptimizationProbeResult, error) {
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
	initial, err := evaluateDirectOffsets(
		"direct_initial",
		currentOffsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return DirectOffsetOptimizationProbeResult{}, err
	}

	selectedEval := initial
	selectedStepNumber := 0
	var selectedStep DirectOffsetStep
	steps := make([]DirectOffsetStep, 0, directOffsetSteps)

	for step := 1; step <= directOffsetSteps; step++ {
		inputGroups := union.groups()
		inputCapacity := fusionCapacity(inputGroups)
		direction := deterministicSPSADirection(step, inputGroups)

		plusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			directOffsetPerturbation, &union,
		)
		if err != nil {
			return DirectOffsetOptimizationProbeResult{}, err
		}
		minusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			-directOffsetPerturbation, &union,
		)
		if err != nil {
			return DirectOffsetOptimizationProbeResult{}, err
		}

		plus, err := evaluateDirectOffsets(
			fmt.Sprintf("direct_plus_%02d", step),
			plusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return DirectOffsetOptimizationProbeResult{}, err
		}
		minus, err := evaluateDirectOffsets(
			fmt.Sprintf("direct_minus_%02d", step),
			minusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return DirectOffsetOptimizationProbeResult{}, err
		}

		gradient := directGradientEstimate(
			plus.Objective, minus.Objective, direction,
		)
		nextOffsets, update, err := applyDirectGradient(
			currentOffsets, gradient, &union,
		)
		if err != nil {
			return DirectOffsetOptimizationProbeResult{}, err
		}
		outputGroups := union.groups()
		outputCapacity := fusionCapacity(outputGroups)
		output, err := evaluateDirectOffsets(
			fmt.Sprintf("direct_output_%02d", step),
			nextOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return DirectOffsetOptimizationProbeResult{}, err
		}

		entry := DirectOffsetStep{
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

	// If no optimized checkpoint beat initialization, keep a complete
	// selected-step record for schema determinism.
	if selectedStepNumber == 0 {
		selectedStep = DirectOffsetStep{
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
		fmt.Sprintf("selected_direct_step_%02d", selectedStepNumber),
		selectedEval.Offsets,
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return DirectOffsetOptimizationProbeResult{}, err
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
		return DirectOffsetOptimizationProbeResult{}, err
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

	return DirectOffsetOptimizationProbeResult{
		Schema:                        DirectOffsetOptimizationSchema,
		Experiment:                    "UP-38-direct-continuous-spectral-offset-optimization",
		LatentDimension:               fullLatentDimension,
		RuntimeStateObjects:           1,
		FullCoordinateMixing:          true,
		StartsFromIndependentOffsets:  true,
		FinishedPartitionMenuProvided: false,
		HardMergeCandidatesEvaluated:  false,
		PairAffinityFieldUsed:         false,
		DirectOffsetsOptimized:        true,
		Optimizer:                     "deterministic-spsa",
		SmoothTaskObjective:           true,
		SmoothResourceObjective:       true,
		StickyExactFusionProjection:   true,
		SelectionUsesHeldOutData:      false,
		ResourcePrice:                 directOffsetResourcePrice,
		SoftCapacityTau:               directOffsetSoftCapacityTau,
		Perturbation:                  directOffsetPerturbation,
		LearningRate:                  directOffsetLearningRate,
		MaximumCoordinateUpdate:       directOffsetMaxUpdate,
		MaximumSteps:                  directOffsetSteps,
		FitTables:                     len(fitTables),
		ValidationTables:              len(validationTables),
		TrueHeldOutTables:             len(trueHeld),
		Initial:                       initial,
		Steps:                         steps,
		SelectedStep:                  selectedStep,
		SelectedFinalArm:              selectedFinal,
		FullCapacityControl:           fullControl,
		Diagnosis: DirectOffsetDiagnosis{
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
			FullCapacityHeldOutAccuracy: fullControl.Static.HeldOutAccuracy,
			FullCapacityCommitAccuracy:  fullControl.Integration.CommitDecodeAccuracy,
			HeldOutRetentionDelta:       heldDelta,
			CommitRetentionDelta:        commitDelta,
			DirectOptimizationPass:      pass,
		},
	}, nil
}

func initialSymmetryGroupsJSON() [][]int {
	out := make([][]int, compositeChannels)
	for i := 0; i < compositeChannels; i++ {
		out[i] = []int{i}
	}
	return out
}
