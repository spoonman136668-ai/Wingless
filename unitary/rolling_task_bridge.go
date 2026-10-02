package unitary

import (
	"fmt"
	"math"
)

const RollingTaskBridgeSchema = "wingless.rolling-task-bridge.v1"

var rollingTaskCheckpointSteps = []int{12, 24, 48}

const rollingTaskMaximumSteps = 48

type RollingTaskStep struct {
	Step                  int       `json:"step"`
	Direction             []float64 `json:"direction"`
	DirectionalDerivative float64   `json:"directional_derivative"`
	FullRankReconstruction bool     `json:"full_rank_reconstruction"`
	UpdateMaxAbs          float64   `json:"update_max_abs"`
	OutputOffsets         []float64 `json:"output_offsets"`
	ProvisionalGroups     [][]int   `json:"provisional_groups"`
	ProvisionalCapacity   int       `json:"provisional_capacity"`
	NearestOffsetGap      float64   `json:"nearest_offset_gap"`
}

type RollingTaskCheckpoint struct {
	Step                int                  `json:"step"`
	Offsets             []float64            `json:"offsets"`
	ProvisionalGroups   [][]int              `json:"provisional_groups"`
	ProvisionalCapacity int                  `json:"provisional_capacity"`
	NearestOffsetGap    float64              `json:"nearest_offset_gap"`
	Evaluation          RichDirectEvaluation `json:"evaluation"`
}

type RollingTaskBridgeDiagnosis struct {
	InnerSplitBalancedPass             bool    `json:"inner_split_balanced_pass"`
	SelectionUsesTrueHeldOut           bool    `json:"selection_uses_true_heldout"`
	FullRankReconstructionSteps        int     `json:"full_rank_reconstruction_steps"`
	SelectedCheckpointStep             int     `json:"selected_checkpoint_step"`
	InitialObjective                   float64 `json:"initial_objective"`
	SelectedObjective                  float64 `json:"selected_objective"`
	SelectedObjectiveGain              float64 `json:"selected_objective_gain"`
	ObjectiveImproved                  bool    `json:"objective_improved"`
	FirstProvisionalFusionStep         int     `json:"first_provisional_fusion_step"`
	MaximumProvisionalCapacity         int     `json:"maximum_provisional_capacity"`
	SelectedProvisionalCapacity        int     `json:"selected_provisional_capacity"`
	SelectedNearestOffsetGap           float64 `json:"selected_nearest_offset_gap"`
	TaskGradientEnteredFusionBasin     bool    `json:"task_gradient_entered_fusion_basin"`
	SelectedHeldOutAccuracy            float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy             float64 `json:"selected_commit_accuracy"`
	SelectedFinalAccuracy              float64 `json:"selected_final_accuracy"`
	SelectedRelationAccuracy           float64 `json:"selected_relation_accuracy"`
	HeldOutRetentionDelta              float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta               float64 `json:"commit_retention_delta"`
	CapabilityGatesWithoutExactFusion  bool    `json:"capability_gates_without_exact_fusion"`
}

type RollingTaskBridgeProbeResult struct {
	Schema                         string                      `json:"schema"`
	Experiment                     string                      `json:"experiment"`
	KnownTargetUsed                bool                        `json:"known_target_used"`
	RollingFullRankDirectionalMemory bool                      `json:"rolling_full_rank_directional_memory"`
	RollingBufferSize              int                         `json:"rolling_buffer_size"`
	StickyProjectionUsed           bool                        `json:"sticky_projection_used"`
	TeacherForcedSmoothMutable     bool                        `json:"teacher_forced_smooth_mutable"`
	FreeRunningProbabilisticMutable bool                       `json:"free_running_probabilistic_mutable"`
	SelectionUsesHeldOutData       bool                        `json:"selection_uses_heldout_data"`
	MaximumSteps                   int                         `json:"maximum_steps"`
	CheckpointSteps                []int                       `json:"checkpoint_steps"`
	Perturbation                   float64                     `json:"perturbation"`
	LearningRate                   float64                     `json:"learning_rate"`
	MaximumCoordinateUpdate        float64                     `json:"maximum_coordinate_update"`
	ExactFusionTolerance           float64                     `json:"exact_fusion_tolerance"`
	ResourcePrice                 float64                     `json:"resource_price"`
	SoftCapacityTau               float64                     `json:"soft_capacity_tau"`
	FitTables                     int                         `json:"fit_tables"`
	ValidationTables              int                         `json:"validation_tables"`
	TrueHeldOutTables             int                         `json:"true_heldout_tables"`
	Initial                       RichDirectEvaluation        `json:"initial"`
	Steps                         []RollingTaskStep            `json:"steps"`
	Checkpoints                   []RollingTaskCheckpoint      `json:"checkpoints"`
	SelectedCheckpoint            RollingTaskCheckpoint       `json:"selected_checkpoint"`
	SelectedFinalArm              MultiplicityDoseArm         `json:"selected_final_arm"`
	FullCapacityControl           MultiplicityDoseArm         `json:"full_capacity_control"`
	Diagnosis                     RollingTaskBridgeDiagnosis   `json:"diagnosis"`
}

func reversibleProvisionalGeometry(
	offsets []float64,
) ([][]int, int, float64, error) {
	if len(offsets) != compositeChannels {
		return nil, 0, 0, fmt.Errorf(
			"provisional geometry dimension mismatch",
		)
	}
	union := newFusionUnion()
	nearest := math.Inf(1)
	for first := 0; first < compositeChannels; first++ {
		for second := first + 1; second < compositeChannels; second++ {
			gap := math.Abs(offsets[first] - offsets[second])
			if gap < nearest {
				nearest = gap
			}
			if gap <= continuousFusionTolerance {
				union.union(first, second)
			}
		}
	}
	groups := union.groups()
	return groups, fusionCapacity(groups), nearest, nil
}

func rollingTaskCheckpointStep(step int) bool {
	for _, candidate := range rollingTaskCheckpointSteps {
		if step == candidate {
			return true
		}
	}
	return false
}

func maxAbsFloat64(values []float64) float64 {
	var maximum float64
	for _, value := range values {
		if magnitude := math.Abs(value); magnitude > maximum {
			maximum = magnitude
		}
	}
	return maximum
}

func RunUP43() (RollingTaskBridgeProbeResult, error) {
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
	initial, err := evaluateFreeRunningDirectOffsets(
		"rolling_task_initial",
		currentOffsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return RollingTaskBridgeProbeResult{}, err
	}
	if !richEvaluationStructurallyValid(initial) {
		return RollingTaskBridgeProbeResult{}, fmt.Errorf(
			"rolling task initial evaluation structurally invalid",
		)
	}

	initialGroups, initialCapacity, initialGap, err :=
		reversibleProvisionalGeometry(currentOffsets)
	if err != nil {
		return RollingTaskBridgeProbeResult{}, err
	}
	selected := RollingTaskCheckpoint{
		Step:                0,
		Offsets:             append([]float64(nil), currentOffsets...),
		ProvisionalGroups:   initialGroups,
		ProvisionalCapacity: initialCapacity,
		NearestOffsetGap:    initialGap,
		Evaluation:          initial,
	}

	observations := make(
		[]rollingDirectionalObservation, 0,
		rollingGradientBufferSize,
	)
	steps := make([]RollingTaskStep, 0, rollingTaskMaximumSteps)
	checkpoints := make(
		[]RollingTaskCheckpoint, 0,
		len(rollingTaskCheckpointSteps),
	)

	fullRankSteps := 0
	firstProvisionalFusionStep := -1
	maximumProvisionalCapacity := initialCapacity

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
			return RollingTaskBridgeProbeResult{}, err
		}
		minusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			-directOffsetPerturbation, &union,
		)
		if err != nil {
			return RollingTaskBridgeProbeResult{}, err
		}

		plus, err := evaluateFreeRunningDirectOffsets(
			fmt.Sprintf("rolling_task_plus_%02d", step),
			plusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return RollingTaskBridgeProbeResult{}, err
		}
		minus, err := evaluateFreeRunningDirectOffsets(
			fmt.Sprintf("rolling_task_minus_%02d", step),
			minusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return RollingTaskBridgeProbeResult{}, err
		}
		if !richEvaluationStructurallyValid(plus) ||
			!richEvaluationStructurallyValid(minus) {
			return RollingTaskBridgeProbeResult{}, fmt.Errorf(
				"rolling task perturbation structurally invalid step=%d",
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
			return RollingTaskBridgeProbeResult{}, err
		}
		if fullRank {
			gradient = reconstructed
			fullRankSteps++
		}

		nextOffsets, update, err :=
			applyCalibrationGradientNoSticky(
				currentOffsets, gradient,
			)
		if err != nil {
			return RollingTaskBridgeProbeResult{}, err
		}
		groups, capacity, nearest, err :=
			reversibleProvisionalGeometry(nextOffsets)
		if err != nil {
			return RollingTaskBridgeProbeResult{}, err
		}
		if capacity > 6 &&
			firstProvisionalFusionStep < 0 {
			firstProvisionalFusionStep = step
		}
		if capacity > maximumProvisionalCapacity {
			maximumProvisionalCapacity = capacity
		}

		steps = append(steps, RollingTaskStep{
			Step:                   step,
			Direction:              append([]float64(nil), direction...),
			DirectionalDerivative:  derivative,
			FullRankReconstruction: fullRank,
			UpdateMaxAbs:           maxAbsFloat64(update),
			OutputOffsets:          append([]float64(nil), nextOffsets...),
			ProvisionalGroups:      groups,
			ProvisionalCapacity:    capacity,
			NearestOffsetGap:       nearest,
		})
		currentOffsets = nextOffsets

		if rollingTaskCheckpointStep(step) {
			evaluation, err := evaluateFreeRunningDirectOffsets(
				fmt.Sprintf("rolling_task_checkpoint_%02d", step),
				currentOffsets,
				mixer,
				fitTables, validationTables,
				trainDepths, heldDepths, allDepths,
				memoryNoise,
			)
			if err != nil {
				return RollingTaskBridgeProbeResult{}, err
			}
			if !richEvaluationStructurallyValid(evaluation) {
				return RollingTaskBridgeProbeResult{}, fmt.Errorf(
					"rolling task checkpoint structurally invalid step=%d",
					step,
				)
			}
			checkpoint := RollingTaskCheckpoint{
				Step:                step,
				Offsets:             append([]float64(nil), currentOffsets...),
				ProvisionalGroups:   groups,
				ProvisionalCapacity: capacity,
				NearestOffsetGap:    nearest,
				Evaluation:          evaluation,
			}
			checkpoints = append(checkpoints, checkpoint)
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
			"selected_rolling_task_checkpoint_%02d",
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
		return RollingTaskBridgeProbeResult{}, err
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
		return RollingTaskBridgeProbeResult{}, err
	}

	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedFinal.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedFinal.Integration.CommitDecodeAccuracy
	capabilityWithoutExactFusion :=
		selectedFinal.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedFinal.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedFinal.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedFinal.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	objectiveGain :=
		selected.Evaluation.Objective - initial.Objective

	return RollingTaskBridgeProbeResult{
		Schema:                           RollingTaskBridgeSchema,
		Experiment:                       "UP-43-rolling-gradient-free-running-task-bridge",
		KnownTargetUsed:                  false,
		RollingFullRankDirectionalMemory: true,
		RollingBufferSize:                rollingGradientBufferSize,
		StickyProjectionUsed:             false,
		TeacherForcedSmoothMutable:       false,
		FreeRunningProbabilisticMutable:  true,
		SelectionUsesHeldOutData:         false,
		MaximumSteps:                     rollingTaskMaximumSteps,
		CheckpointSteps:                  append([]int(nil), rollingTaskCheckpointSteps...),
		Perturbation:                     directOffsetPerturbation,
		LearningRate:                     directOffsetLearningRate,
		MaximumCoordinateUpdate:          directOffsetMaxUpdate,
		ExactFusionTolerance:             continuousFusionTolerance,
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
		Diagnosis: RollingTaskBridgeDiagnosis{
			InnerSplitBalancedPass:            innerBalanced,
			SelectionUsesTrueHeldOut:          false,
			FullRankReconstructionSteps:       fullRankSteps,
			SelectedCheckpointStep:            selected.Step,
			InitialObjective:                  initial.Objective,
			SelectedObjective:                 selected.Evaluation.Objective,
			SelectedObjectiveGain:             objectiveGain,
			ObjectiveImproved:                 objectiveGain > 0,
			FirstProvisionalFusionStep:        firstProvisionalFusionStep,
			MaximumProvisionalCapacity:        maximumProvisionalCapacity,
			SelectedProvisionalCapacity:       selected.ProvisionalCapacity,
			SelectedNearestOffsetGap:          selected.NearestOffsetGap,
			TaskGradientEnteredFusionBasin:    maximumProvisionalCapacity > 6,
			SelectedHeldOutAccuracy:           selectedFinal.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:            selectedFinal.Integration.CommitDecodeAccuracy,
			SelectedFinalAccuracy:             selectedFinal.Integration.ExactFinalTableAccuracy,
			SelectedRelationAccuracy:          selectedFinal.Integration.RelationalQueryAccuracy,
			HeldOutRetentionDelta:             heldDelta,
			CommitRetentionDelta:              commitDelta,
			CapabilityGatesWithoutExactFusion: capabilityWithoutExactFusion,
		},
	}, nil
}
