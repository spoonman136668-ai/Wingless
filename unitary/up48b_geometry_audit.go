package unitary

import (
	"fmt"
	"math"
)

const UP48BGeometryAuditSchema = "wingless.up48b-frozen-geometry-audit.v1"

type UP48BArm struct {
	Name          string               `json:"name"`
	Source        string               `json:"source"`
	Offsets       []float64            `json:"offsets"`
	ExactCapacity int                  `json:"exact_capacity"`
	ExactFusion   bool                 `json:"exact_fusion"`
	Smooth        RichDirectEvaluation `json:"smooth"`
	Hard          MultiplicityDoseArm  `json:"hard"`
}

type UP48BMetricDelta struct {
	Objective        float64 `json:"objective"`
	HeldOutAccuracy  float64 `json:"heldout_accuracy"`
	CommitAccuracy   float64 `json:"commit_accuracy"`
	FinalAccuracy    float64 `json:"final_accuracy"`
	RelationAccuracy float64 `json:"relation_accuracy"`
}

type UP48BDiagnosis struct {
	Step29MinusMatchedSplit UP48BMetricDelta `json:"step29_minus_matched_split"`
	Step42MinusStep29       UP48BMetricDelta `json:"step42_minus_step29"`
	Step42FusedMinusStep42  UP48BMetricDelta `json:"step42_fused_minus_step42"`
}

type UP48BGeometryAuditResult struct {
	Schema                   string           `json:"schema"`
	Experiment               string           `json:"experiment"`
	SourceUP47Head           string           `json:"source_up47_head"`
	OptimizerRun             bool             `json:"optimizer_run"`
	SelectionUsesHeldOutData bool             `json:"selection_uses_heldout_data"`
	ControlGapSourceStep     int              `json:"control_gap_source_step"`
	Step42FusedPair          []int            `json:"step42_fused_pair"`
	Arms                     []UP48BArm       `json:"arms"`
	FullCapacityControl      MultiplicityDoseArm `json:"full_capacity_control"`
	Diagnosis                UP48BDiagnosis   `json:"diagnosis"`
}

func up48bFrozenOffsets(step int) ([]float64, error) {
	if step < 0 || step >= len(up45FrozenDenseTrajectory) {
		return nil, fmt.Errorf("UP48B frozen step out of range: %d", step)
	}
	state := up45FrozenDenseTrajectory[step]
	if state.Step != step {
		return nil, fmt.Errorf("UP48B frozen step mismatch index=%d state=%d", step, state.Step)
	}
	return append([]float64(nil), state.Offsets...), nil
}

func up48bArm(
	name, source string,
	offsets []float64,
	mixer latentMatrix,
	allTrain, trueHeld []memoryTable,
	fitTables, validationTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (UP48BArm, error) {
	smooth, err := evaluateSqrtFreeRunningDirectOffsets(
		name+"_smooth",
		append([]float64(nil), offsets...),
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return UP48BArm{}, err
	}
	if !richEvaluationStructurallyValid(smooth) {
		return UP48BArm{}, fmt.Errorf("UP48B smooth evaluation invalid for %s", name)
	}

	hard, err := evaluateMultiplicityDoseArm(
		name+"_hard",
		append([]float64(nil), offsets...),
		mixer,
		allTrain, trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return UP48BArm{}, err
	}

	_, capacity, err := exactOffsetGroups(offsets)
	if err != nil {
		return UP48BArm{}, err
	}

	return UP48BArm{
		Name:          name,
		Source:        source,
		Offsets:       append([]float64(nil), offsets...),
		ExactCapacity: capacity,
		ExactFusion:   capacity > compositeChannels,
		Smooth:        smooth,
		Hard:          hard,
	}, nil
}

func up48bDelta(a, b UP48BArm) UP48BMetricDelta {
	return UP48BMetricDelta{
		Objective:        a.Smooth.Objective - b.Smooth.Objective,
		HeldOutAccuracy:  a.Hard.Static.HeldOutAccuracy - b.Hard.Static.HeldOutAccuracy,
		CommitAccuracy:   a.Hard.Integration.CommitDecodeAccuracy - b.Hard.Integration.CommitDecodeAccuracy,
		FinalAccuracy:    a.Hard.Integration.ExactFinalTableAccuracy - b.Hard.Integration.ExactFinalTableAccuracy,
		RelationAccuracy: a.Hard.Integration.RelationalQueryAccuracy - b.Hard.Integration.RelationalQueryAccuracy,
	}
}

func RunUP48B() (UP48BGeometryAuditResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	step24, err := up48bFrozenOffsets(24)
	if err != nil {
		return UP48BGeometryAuditResult{}, err
	}
	step28, err := up48bFrozenOffsets(28)
	if err != nil {
		return UP48BGeometryAuditResult{}, err
	}
	step29, err := up48bFrozenOffsets(29)
	if err != nil {
		return UP48BGeometryAuditResult{}, err
	}
	step42, err := up48bFrozenOffsets(42)
	if err != nil {
		return UP48BGeometryAuditResult{}, err
	}

	// Matched exact-equality ablation: preserve the step-29 pair mean and all
	// other coordinates, but restore the pair gap observed one step earlier at
	// frozen step 28. This magnitude is trajectory-derived, not tuned.
	step29Split := append([]float64(nil), step29...)
	gap28 := step28[0] - step28[3]
	center29 := 0.5 * (step29[0] + step29[3])
	step29Split[0] = center29 + 0.5*gap28
	step29Split[3] = center29 - 0.5*gap28
	if math.Abs(step29Split[0]-step29Split[3]) < 1e-12 {
		return UP48BGeometryAuditResult{}, fmt.Errorf("UP48B matched split unexpectedly fused")
	}

	// Step 42's nearest frozen pair is [4,5]. Fuse only that pair at its mean
	// while preserving every other coordinate.
	nearestI, nearestJ := 0, 1
	nearestGap := math.Inf(1)
	for i := 0; i < len(step42); i++ {
		for j := i + 1; j < len(step42); j++ {
			g := math.Abs(step42[i] - step42[j])
			if g < nearestGap {
				nearestGap = g
				nearestI, nearestJ = i, j
			}
		}
	}
	if nearestI != 4 || nearestJ != 5 {
		return UP48BGeometryAuditResult{}, fmt.Errorf(
			"UP48B step42 nearest pair=(%d,%d) want=(4,5)",
			nearestI, nearestJ,
		)
	}
	step42Fused := append([]float64(nil), step42...)
	center42 := 0.5 * (step42[4] + step42[5])
	step42Fused[4] = center42
	step42Fused[5] = center42

	mixer := fullLatentMixer()
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fitTables, validationTables := splitTaskAllocationTrainingPool(allTrain)

	specs := []struct {
		name, source string
		offsets      []float64
	}{
		{"step24_sparse_selected", "UP45 step 24", step24},
		{"step29_exact_fused", "UP45 step 29", step29},
		{"step29_matched_split", "UP45 step 29 with step28 pair gap", step29Split},
		{"step42_dense_selected", "UP47 step 42", step42},
		{"step42_nearest_pair_fused", "UP47 step 42 with pair [4,5] fused at mean", step42Fused},
	}

	arms := make([]UP48BArm, 0, len(specs))
	for _, spec := range specs {
		arm, err := up48bArm(
			spec.name, spec.source, spec.offsets,
			mixer, allTrain, trueHeld, fitTables, validationTables,
			trainDepths, heldDepths, allDepths, memoryNoise,
		)
		if err != nil {
			return UP48BGeometryAuditResult{}, err
		}
		arms = append(arms, arm)
	}

	fullControl, err := evaluateMultiplicityDoseArm(
		"up48b_full_capacity_control",
		[]float64{0, 0, 0, 0, 0, 0},
		mixer, allTrain, trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return UP48BGeometryAuditResult{}, err
	}

	return UP48BGeometryAuditResult{
		Schema:                   UP48BGeometryAuditSchema,
		Experiment:               "UP-48B-frozen-geometry-causal-audit",
		SourceUP47Head:           "474ebb42082d1a9dd022517edb3c48c4f16dbd1d",
		OptimizerRun:             false,
		SelectionUsesHeldOutData: false,
		ControlGapSourceStep:     28,
		Step42FusedPair:          []int{4, 5},
		Arms:                     arms,
		FullCapacityControl:      fullControl,
		Diagnosis: UP48BDiagnosis{
			Step29MinusMatchedSplit: up48bDelta(arms[1], arms[2]),
			Step42MinusStep29:       up48bDelta(arms[3], arms[1]),
			Step42FusedMinusStep42:  up48bDelta(arms[4], arms[3]),
		},
	}, nil
}
