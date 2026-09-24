package unitary

import (
	"fmt"
)

const ProximalTaskBridgeSchema =
	"wingless.proximal-task-bridge.v1"

type ProximalTaskStep struct {
	Step                   int                       `json:"step"`
	Direction              []float64                 `json:"direction"`
	DirectionalDerivative  float64                   `json:"directional_derivative"`
	FullRankReconstruction bool                      `json:"full_rank_reconstruction"`
	GradientUpdateMaxAbs   float64                   `json:"gradient_update_max_abs"`
	TotalUpdateMaxAbs      float64                   `json:"total_update_max_abs"`
	OutputOffsets          []float64                 `json:"output_offsets"`
	ExactGroups            [][]int                   `json:"exact_groups"`
	ExactCapacity          int                       `json:"exact_capacity"`
	ExactFusion            bool                      `json:"exact_fusion"`
	FusionSeparatedLater   bool                      `json:"fusion_separated_later"`
	Proximal               ProximalFusionDiagnostics `json:"proximal"`
}

type ProximalTaskCheckpoint struct {
	Step          int                  `json:"step"`
	Offsets       []float64            `json:"offsets"`
	ExactGroups   [][]int              `json:"exact_groups"`
	ExactCapacity int                  `json:"exact_capacity"`
	ExactFusion   bool                 `json:"exact_fusion"`
	Evaluation    RichDirectEvaluation `json:"evaluation"`
}

type ProximalTaskDiagnosis struct {
	InnerSplitBalancedPass       bool    `json:"inner_split_balanced_pass"`
	SelectionUsesTrueHeldOut     bool    `json:"selection_uses_true_heldout"`
	FullRankReconstructionSteps  int     `json:"full_rank_reconstruction_steps"`
	SelectedCheckpointStep       int     `json:"selected_checkpoint_step"`
	InitialObjective             float64 `json:"initial_objective"`
	SelectedObjective            float64 `json:"selected_objective"`
	SelectedObjectiveGain        float64 `json:"selected_objective_gain"`
	ObjectiveImproved            bool    `json:"objective_improved"`
	FirstExactFusionStep         int     `json:"first_exact_fusion_step"`
	MaximumExactCapacity         int     `json:"maximum_exact_capacity"`
	SelectedExactCapacity        int     `json:"selected_exact_capacity"`
	SelectedExactFusion          bool    `json:"selected_exact_fusion"`
	AnyFusionSeparatedLater      bool    `json:"any_fusion_separated_later"`
	SelectedHeldOutAccuracy      float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy       float64 `json:"selected_commit_accuracy"`
	SelectedFinalAccuracy        float64 `json:"selected_final_accuracy"`
	SelectedRelationAccuracy     float64 `json:"selected_relation_accuracy"`
	HeldOutRetentionDelta        float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta         float64 `json:"commit_retention_delta"`
	CapabilityGates              bool    `json:"capability_gates"`
	ExactFusionAndCapabilityPass bool    `json:"exact_fusion_and_capability_pass"`
}

type ProximalTaskBridgeProbeResult struct {
	Schema                           string                  `json:"schema"`
	Experiment                       string                  `json:"experiment"`
	KnownTargetUsed                  bool                    `json:"known_target_used"`
	FinishedPartitionMenuProvided    bool                    `json:"finished_partition_menu_provided"`
	HardMergeCandidatesEvaluated     bool                    `json:"hard_merge_candidates_evaluated"`
	PairAffinityFieldUsed            bool                    `json:"pair_affinity_field_used"`
	RollingFullRankDirectionalMemory bool                    `json:"rolling_full_rank_directional_memory"`
	RollingBufferSize                int                     `json:"rolling_buffer_size"`
	StickyProjectionUsed             bool                    `json:"sticky_projection_used"`
	ConvexProximalFusionUsed         bool                    `json:"convex_proximal_fusion_used"`
	ProximalLambda                   float64                 `json:"proximal_lambda"`
	ProximalLambdaDerived            bool                    `json:"proximal_lambda_derived"`
	SqrtProbabilitySoftMemory        bool                    `json:"sqrt_probability_soft_memory"`
	GlobalSoftMemoryRenormalization  bool                    `json:"global_soft_memory_renormalization"`
	SelectionUsesHeldOutData         bool                    `json:"selection_uses_heldout_data"`
	ExactFusionUsedForSelection      bool                    `json:"exact_fusion_used_for_selection"`
	MaximumSteps                     int                     `json:"maximum_steps"`
	CheckpointSteps                  []int                   `json:"checkpoint_steps"`
	Perturbation                     float64                 `json:"perturbation"`
	LearningRate                     float64                 `json:"learning_rate"`
	MaximumCoordinateUpdate          float64                 `json:"maximum_coordinate_update"`
	ResourcePrice                    float64                 `json:"resource_price"`
	SoftCapacityTau                  float64                 `json:"soft_capacity_tau"`
	FitTables                        int                     `json:"fit_tables"`
	ValidationTables                 int                     `json:"validation_tables"`
	TrueHeldOutTables                int                     `json:"true_heldout_tables"`
	Initial                          RichDirectEvaluation    `json:"initial"`
	Steps                            []ProximalTaskStep       `json:"steps"`
	Checkpoints                      []ProximalTaskCheckpoint `json:"checkpoints"`
	SelectedCheckpoint               ProximalTaskCheckpoint  `json:"selected_checkpoint"`
	SelectedFinalArm                 MultiplicityDoseArm     `json:"selected_final_arm"`
	FullCapacityControl              MultiplicityDoseArm     `json:"full_capacity_control"`
	Diagnosis                        ProximalTaskDiagnosis    `json:"diagnosis"`
}

func exactPairSet(groups [][]int) map[[2]int]struct{} {
	out := make(map[[2]int]struct{})
	for _, group := range groups {
		for i := 0; i < len(group); i++ {
			for j := i + 1; j < len(group); j++ {
				a, b := group[i], group[j]
				if a > b {
					a, b = b, a
				}
				out[[2]int{a, b}] = struct{}{}
			}
		}
	}
	return out
}

func anyExactPairSeparated(
	previous, current map[[2]int]struct{},
) bool {
	for pair := range previous {
		if _, ok := current[pair]; !ok {
			return true
		}
	}
	return false
}

func maxOffsetDelta(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	delta := make([]float64, len(a))
	for i := range a {
		delta[i] = b[i] - a[i]
	}
	return maxAbsFloat64(delta)
}

func RunUP45() (ProximalTaskBridgeProbeResult, error) {
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

	currentOffsets := continuousFusionInitialOffsets()
	initial, err := evaluateSqrtFreeRunningDirectOffsets(
		"proximal_task_initial",
		currentOffsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return ProximalTaskBridgeProbeResult{}, err
	}
	if !richEvaluationStructurallyValid(initial) {
		return ProximalTaskBridgeProbeResult{}, fmt.Errorf(
			"proximal task initial evaluation structurally invalid",
		)
	}

	initialGroups, initialCapacity, err :=
		exactOffsetGroups(currentOffsets)
	if err != nil {
		return ProximalTaskBridgeProbeResult{}, err
	}
	selected := ProximalTaskCheckpoint{
		Step:          0,
		Offsets:       append([]float64(nil), currentOffsets...),
		ExactGroups:   initialGroups,
		ExactCapacity: initialCapacity,
		ExactFusion:   initialCapacity > 6,
		Evaluation:    initial,
	}

	observations := make(
		[]rollingDirectionalObservation, 0,
		rollingGradientBufferSize,
	)
	steps := make([]ProximalTaskStep, 0, rollingTaskMaximumSteps)
	checkpoints := make(
		[]ProximalTaskCheckpoint, 0,
		len(rollingTaskCheckpointSteps),
	)

	fullRankSteps := 0
	firstExactFusionStep := -1
	maximumExactCapacity := initialCapacity
	anySeparated := false
	previousExactPairs := exactPairSet(initialGroups)

	for step := 1; step <= rollingTaskMaximumSteps; step++ {
		union := newFusionUnion()
		direction := deterministicSPSADirection(
			step, union.groups(),
		)
		plusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			directOffsetPerturbation, &union,
		)
		if err != nil {
			return ProximalTaskBridgeProbeResult{}, err
		}
		minusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			-directOffsetPerturbation, &union,
		)
		if err != nil {
			return ProximalTaskBridgeProbeResult{}, err
		}

		plus, err := evaluateSqrtFreeRunningDirectOffsets(
			fmt.Sprintf("proximal_task_plus_%02d", step),
			plusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return ProximalTaskBridgeProbeResult{}, err
		}
		minus, err := evaluateSqrtFreeRunningDirectOffsets(
			fmt.Sprintf("proximal_task_minus_%02d", step),
			minusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return ProximalTaskBridgeProbeResult{}, err
		}
		if !richEvaluationStructurallyValid(plus) ||
			!richEvaluationStructurallyValid(minus) {
			return ProximalTaskBridgeProbeResult{}, fmt.Errorf(
				"proximal task perturbation structurally invalid step=%d",
				step,
			)
		}

		derivative :=
			(plus.Objective - minus.Objective) /
				(2 * directOffsetPerturbation)
		rankOne := directGradientEstimate(
			plus.Objective, minus.Objective, direction,
		)
		observations = append(
			observations,
			rollingDirectionalObservation{
				direction: append([]float64(nil), direction...),
				derivative: derivative,
			},
		)
		if len(observations) > rollingGradientBufferSize {
			observations = observations[
				len(observations)-rollingGradientBufferSize:
			]
		}

		gradient := rankOne
		reconstructed, fullRank, err :=
			rollingFullRankGradient(observations)
		if err != nil {
			return ProximalTaskBridgeProbeResult{}, err
		}
		if fullRank {
			gradient = reconstructed
			fullRankSteps++
		}

		nextOffsets, gradientUpdate, prox, err :=
			applyProximalFusionGradient(
				currentOffsets, gradient,
			)
		if err != nil {
			return ProximalTaskBridgeProbeResult{}, err
		}
		groups, capacity, err := exactOffsetGroups(nextOffsets)
		if err != nil {
			return ProximalTaskBridgeProbeResult{}, err
		}
		exactFusion := capacity > 6
		if exactFusion && firstExactFusionStep < 0 {
			firstExactFusionStep = step
		}
		if capacity > maximumExactCapacity {
			maximumExactCapacity = capacity
		}

		currentPairs := exactPairSet(groups)
		separated := anyExactPairSeparated(
			previousExactPairs, currentPairs,
		)
		if separated {
			anySeparated = true
		}
		previousExactPairs = currentPairs

		steps = append(steps, ProximalTaskStep{
			Step:                   step,
			Direction:              append([]float64(nil), direction...),
			DirectionalDerivative:  derivative,
			FullRankReconstruction: fullRank,
			GradientUpdateMaxAbs:   maxAbsFloat64(gradientUpdate),
			TotalUpdateMaxAbs:      maxOffsetDelta(currentOffsets, nextOffsets),
			OutputOffsets:          append([]float64(nil), nextOffsets...),
			ExactGroups:            groups,
			ExactCapacity:          capacity,
			ExactFusion:            exactFusion,
			FusionSeparatedLater:   separated,
			Proximal:               prox,
		})
		currentOffsets = nextOffsets

		if rollingTaskCheckpointStep(step) {
			evaluation, err :=
				evaluateSqrtFreeRunningDirectOffsets(
					fmt.Sprintf(
						"proximal_task_checkpoint_%02d",
						step,
					),
					currentOffsets,
					mixer,
					fitTables, validationTables,
					trainDepths, heldDepths, allDepths,
					memoryNoise,
				)
			if err != nil {
				return ProximalTaskBridgeProbeResult{}, err
			}
			if !richEvaluationStructurallyValid(evaluation) {
				return ProximalTaskBridgeProbeResult{}, fmt.Errorf(
					"proximal task checkpoint structurally invalid step=%d",
					step,
				)
			}
			checkpoint := ProximalTaskCheckpoint{
				Step:          step,
				Offsets:       append([]float64(nil), currentOffsets...),
				ExactGroups:   groups,
				ExactCapacity: capacity,
				ExactFusion:   exactFusion,
				Evaluation:    evaluation,
			}
			checkpoints = append(checkpoints, checkpoint)

			// Exact fusion is intentionally not part of selection.
			if evaluation.Objective > selected.Evaluation.Objective ||
				(evaluation.Objective == selected.Evaluation.Objective &&
					evaluation.SoftCapacity <
						selected.Evaluation.SoftCapacity) {
				selected = checkpoint
			}
		}
	}

	selectedFinal, err := evaluateMultiplicityDoseArm(
		fmt.Sprintf(
			"selected_proximal_task_checkpoint_%02d",
			selected.Step,
		),
		selected.Offsets,
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return ProximalTaskBridgeProbeResult{}, err
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
		return ProximalTaskBridgeProbeResult{}, err
	}

	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedFinal.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedFinal.Integration.CommitDecodeAccuracy
	capability :=
		selectedFinal.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedFinal.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedFinal.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedFinal.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	objectiveGain :=
		selected.Evaluation.Objective - initial.Objective

	return ProximalTaskBridgeProbeResult{
		Schema:                           ProximalTaskBridgeSchema,
		Experiment:                       "UP-45-reversible-convex-exact-fusion",
		KnownTargetUsed:                  false,
		FinishedPartitionMenuProvided:    false,
		HardMergeCandidatesEvaluated:     false,
		PairAffinityFieldUsed:            false,
		RollingFullRankDirectionalMemory: true,
		RollingBufferSize:                rollingGradientBufferSize,
		StickyProjectionUsed:             false,
		ConvexProximalFusionUsed:         true,
		ProximalLambda:                   proximalFusionLambda,
		ProximalLambdaDerived:            true,
		SqrtProbabilitySoftMemory:        true,
		GlobalSoftMemoryRenormalization:  false,
		SelectionUsesHeldOutData:         false,
		ExactFusionUsedForSelection:      false,
		MaximumSteps:                     rollingTaskMaximumSteps,
		CheckpointSteps:                  append([]int(nil), rollingTaskCheckpointSteps...),
		Perturbation:                     directOffsetPerturbation,
		LearningRate:                     directOffsetLearningRate,
		MaximumCoordinateUpdate:          directOffsetMaxUpdate,
		ResourcePrice:                    directOffsetResourcePrice,
		SoftCapacityTau:                  directOffsetSoftCapacityTau,
		FitTables:                        len(fitTables),
		ValidationTables:                 len(validationTables),
		TrueHeldOutTables:                len(trueHeld),
		Initial:                          initial,
		Steps:                            steps,
		Checkpoints:                      checkpoints,
		SelectedCheckpoint:               selected,
		SelectedFinalArm:                 selectedFinal,
		FullCapacityControl:              fullControl,
		Diagnosis: ProximalTaskDiagnosis{
			InnerSplitBalancedPass:       innerBalanced,
			SelectionUsesTrueHeldOut:     false,
			FullRankReconstructionSteps:  fullRankSteps,
			SelectedCheckpointStep:       selected.Step,
			InitialObjective:             initial.Objective,
			SelectedObjective:            selected.Evaluation.Objective,
			SelectedObjectiveGain:        objectiveGain,
			ObjectiveImproved:            objectiveGain > 0,
			FirstExactFusionStep:         firstExactFusionStep,
			MaximumExactCapacity:         maximumExactCapacity,
			SelectedExactCapacity:        selected.ExactCapacity,
			SelectedExactFusion:          selected.ExactFusion,
			AnyFusionSeparatedLater:      anySeparated,
			SelectedHeldOutAccuracy:      selectedFinal.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:       selectedFinal.Integration.CommitDecodeAccuracy,
			SelectedFinalAccuracy:        selectedFinal.Integration.ExactFinalTableAccuracy,
			SelectedRelationAccuracy:     selectedFinal.Integration.RelationalQueryAccuracy,
			HeldOutRetentionDelta:        heldDelta,
			CommitRetentionDelta:         commitDelta,
			CapabilityGates:              capability,
			ExactFusionAndCapabilityPass: selected.ExactFusion && capability,
		},
	}, nil
}
