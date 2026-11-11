package unitary

import (
	"fmt"
	"math"
)

const OptimizerReachabilityCalibrationSchema =
	"wingless.optimizer-reachability-calibration.v1"

const optimizerCalibrationTargetSource =
	"UP-35-seal-09284aa70203adf8881be42f67084c49b338dabe"

type OptimizerCalibrationStep struct {
	Step                    int         `json:"step"`
	InputOffsets            []float64   `json:"input_offsets"`
	InputGroups             [][]int     `json:"input_groups"`
	Direction               []float64   `json:"direction,omitempty"`
	GradientEstimate        []float64   `json:"gradient_estimate"`
	Update                  []float64   `json:"update"`
	OutputOffsets           []float64   `json:"output_offsets"`
	OutputGroups            [][]int     `json:"output_groups"`
	TargetCosine            float64     `json:"target_cosine"`
	TargetGroupMaxSpread    float64     `json:"target_group_max_spread"`
	TargetFusionBasin       bool        `json:"target_fusion_basin"`
	WrongCrossTargetFusion  bool        `json:"wrong_cross_target_fusion"`
}

type OptimizerCalibrationArm struct {
	Name                       string                     `json:"name"`
	Estimator                  string                     `json:"estimator"`
	StickyProjection           bool                       `json:"sticky_projection"`
	InitialTargetCosine        float64                    `json:"initial_target_cosine"`
	BestStep                   int                        `json:"best_step"`
	BestTargetCosine           float64                    `json:"best_target_cosine"`
	FinalTargetCosine          float64                    `json:"final_target_cosine"`
	FinalTargetGroupMaxSpread  float64                    `json:"final_target_group_max_spread"`
	EnteredTargetFusionBasin   bool                       `json:"entered_target_fusion_basin"`
	FirstTargetBasinStep       int                        `json:"first_target_basin_step"`
	ExactTargetGroupingEver    bool                       `json:"exact_target_grouping_ever"`
	WrongCrossTargetFusionEver bool                       `json:"wrong_cross_target_fusion_ever"`
	FinalGroups                [][]int                    `json:"final_groups"`
	Steps                      []OptimizerCalibrationStep `json:"steps"`
}

type OptimizerCalibrationDiagnosis struct {
	CurrentOptimizerEntersTargetBasin          bool    `json:"current_optimizer_enters_target_basin"`
	CurrentOptimizerExactTargetGrouping        bool    `json:"current_optimizer_exact_target_grouping"`
	AnalyticReferenceEntersTargetBasin         bool    `json:"analytic_reference_enters_target_basin"`
	SequentialEstimatorTrailsAnalyticReference bool    `json:"sequential_estimator_trails_analytic_reference"`
	StickyProjectionCrossTargetFusion          bool    `json:"sticky_projection_cross_target_fusion"`
	TwelveStepReachabilityLimited              bool    `json:"twelve_step_reachability_limited"`
	SequentialNoStickyBestCosine               float64 `json:"sequential_no_sticky_best_cosine"`
	AnalyticNoStickyBestCosine                 float64 `json:"analytic_no_sticky_best_cosine"`
	EstimatorBestCosineDelta                   float64 `json:"estimator_best_cosine_delta"`
	OptimizerMechanicsImplicated               bool    `json:"optimizer_mechanics_implicated"`
}

type OptimizerReachabilityCalibrationProbeResult struct {
	Schema                         string                         `json:"schema"`
	Experiment                     string                         `json:"experiment"`
	CalibrationOnly                bool                           `json:"calibration_only"`
	TaskLabelsUsed                 bool                           `json:"task_labels_used"`
	HeldOutDataUsed                bool                           `json:"heldout_data_used"`
	ArchitectureSelectionAuthorized bool                         `json:"architecture_selection_authorized"`
	TargetSource                   string                         `json:"target_source"`
	TargetGroups                   [][]int                        `json:"target_groups"`
	TargetCapacity                 int                            `json:"target_capacity"`
	TargetOffsets                  []float64                      `json:"target_offsets"`
	StartOffsets                   []float64                      `json:"start_offsets"`
	ZeroMeanConstraint             bool                           `json:"zero_mean_constraint"`
	FixedRMS                       float64                        `json:"fixed_rms"`
	Perturbation                   float64                        `json:"perturbation"`
	LearningRate                   float64                        `json:"learning_rate"`
	MaximumCoordinateUpdate        float64                        `json:"maximum_coordinate_update"`
	OptimizationSteps              int                            `json:"optimization_steps"`
	ExactFusionTolerance           float64                        `json:"exact_fusion_tolerance"`
	Arms                           []OptimizerCalibrationArm       `json:"arms"`
	Diagnosis                      OptimizerCalibrationDiagnosis   `json:"diagnosis"`
}

func optimizerCalibrationTarget() ([]symmetryGroup, []float64, error) {
	groups := []symmetryGroup{
		{members: []int{0, 1, 2, 5}},
		{members: []int{3}},
		{members: []int{4}},
	}
	normalized, err := normalizeSymmetryGroups(groups)
	if err != nil {
		return nil, nil, err
	}
	offsets, err := offsetsForSymmetryGroups(
		normalized, multiplicityDoseRMS,
	)
	if err != nil {
		return nil, nil, err
	}
	return normalized, offsets, nil
}

func offsetCosineSimilarity(first, second []float64) (float64, error) {
	if len(first) != compositeChannels ||
		len(second) != compositeChannels {
		return 0, fmt.Errorf("offset cosine dimension mismatch")
	}
	var dot, first2, second2 float64
	for index := 0; index < compositeChannels; index++ {
		if !finite(first[index]) || !finite(second[index]) {
			return 0, fmt.Errorf("offset cosine non-finite input")
		}
		dot += first[index] * second[index]
		first2 += first[index] * first[index]
		second2 += second[index] * second[index]
	}
	if first2 <= 1e-15 || second2 <= 1e-15 {
		return 0, fmt.Errorf("offset cosine zero norm")
	}
	value := dot / math.Sqrt(first2*second2)
	if !finite(value) {
		return 0, fmt.Errorf("offset cosine non-finite")
	}
	return value, nil
}

func analyticOffsetCosineGradient(
	offsets, target []float64,
) ([]float64, error) {
	if len(offsets) != compositeChannels ||
		len(target) != compositeChannels {
		return nil, fmt.Errorf("analytic cosine gradient dimension mismatch")
	}
	var dot, offset2, target2 float64
	for index := 0; index < compositeChannels; index++ {
		dot += offsets[index] * target[index]
		offset2 += offsets[index] * offsets[index]
		target2 += target[index] * target[index]
	}
	if offset2 <= 1e-15 || target2 <= 1e-15 {
		return nil, fmt.Errorf("analytic cosine gradient zero norm")
	}
	offsetNorm := math.Sqrt(offset2)
	targetNorm := math.Sqrt(target2)
	cosine := dot / (offsetNorm * targetNorm)
	gradient := make([]float64, compositeChannels)
	for index := range gradient {
		gradient[index] =
			target[index]/(offsetNorm*targetNorm) -
				cosine*offsets[index]/offset2
	}
	return gradient, nil
}

func projectGradientToCalibrationGroups(
	gradient []float64,
	groups [][]int,
) ([]float64, error) {
	if len(gradient) != compositeChannels {
		return nil, fmt.Errorf("calibration gradient dimension mismatch")
	}
	out := append([]float64(nil), gradient...)
	for _, group := range groups {
		if len(group) == 0 {
			return nil, fmt.Errorf("calibration projection empty group")
		}
		var mean float64
		for _, member := range group {
			if member < 0 || member >= compositeChannels {
				return nil, fmt.Errorf("calibration projection member out of range")
			}
			mean += out[member]
		}
		mean /= float64(len(group))
		for _, member := range group {
			out[member] = mean
		}
	}
	var mean float64
	for _, value := range out {
		mean += value
	}
	mean /= float64(len(out))
	for index := range out {
		out[index] -= mean
	}
	return out, nil
}

func applyCalibrationGradientNoSticky(
	offsets, gradient []float64,
) ([]float64, []float64, error) {
	if len(offsets) != compositeChannels ||
		len(gradient) != compositeChannels {
		return nil, nil, fmt.Errorf("calibration update dimension mismatch")
	}
	update := make([]float64, compositeChannels)
	proposal := append([]float64(nil), offsets...)
	for index := range proposal {
		update[index] = clampDirectUpdate(
			directOffsetLearningRate * gradient[index],
		)
		proposal[index] += update[index]
	}
	normalized, err := normalizeContinuousOffsets(proposal)
	if err != nil {
		return nil, nil, err
	}
	return normalized, update, nil
}

func calibrationTargetGroupMaxSpread(
	offsets []float64,
	targetGroups [][]int,
) (float64, error) {
	if len(offsets) != compositeChannels {
		return 0, fmt.Errorf("target spread offset dimension mismatch")
	}
	maximum := 0.0
	for _, group := range targetGroups {
		if len(group) < 2 {
			continue
		}
		minimum := math.Inf(1)
		maximumValue := math.Inf(-1)
		for _, member := range group {
			if member < 0 || member >= compositeChannels {
				return 0, fmt.Errorf("target spread member out of range")
			}
			value := offsets[member]
			if value < minimum {
				minimum = value
			}
			if value > maximumValue {
				maximumValue = value
			}
		}
		spread := maximumValue - minimum
		if spread > maximum {
			maximum = spread
		}
	}
	return maximum, nil
}

func calibrationTargetFusionBasin(
	offsets []float64,
	targetGroups [][]int,
) (bool, error) {
	spread, err :=
		calibrationTargetGroupMaxSpread(offsets, targetGroups)
	if err != nil {
		return false, err
	}
	if spread > continuousFusionTolerance {
		return false, nil
	}

	groupOf := make([]int, compositeChannels)
	for index := range groupOf {
		groupOf[index] = -1
	}
	for groupIndex, group := range targetGroups {
		for _, member := range group {
			groupOf[member] = groupIndex
		}
	}
	for first := 0; first < compositeChannels; first++ {
		for second := first + 1; second < compositeChannels; second++ {
			if groupOf[first] < 0 || groupOf[second] < 0 {
				return false, fmt.Errorf("target groups incomplete")
			}
			if groupOf[first] == groupOf[second] {
				continue
			}
			if math.Abs(offsets[first]-offsets[second]) <=
				continuousFusionTolerance {
				return false, nil
			}
		}
	}
	return true, nil
}

func calibrationWrongCrossTargetFusion(
	currentGroups, targetGroups [][]int,
) (bool, error) {
	targetOf := make([]int, compositeChannels)
	for index := range targetOf {
		targetOf[index] = -1
	}
	for targetIndex, target := range targetGroups {
		for _, member := range target {
			if member < 0 || member >= compositeChannels {
				return false, fmt.Errorf("target group member out of range")
			}
			if targetOf[member] >= 0 {
				return false, fmt.Errorf("target group member duplicated")
			}
			targetOf[member] = targetIndex
		}
	}
	for _, current := range currentGroups {
		if len(current) < 2 {
			continue
		}
		targetIndex := targetOf[current[0]]
		for _, member := range current[1:] {
			if targetOf[member] != targetIndex {
				return true, nil
			}
		}
	}
	return false, nil
}

func calibrationGroupsEqual(first, second [][]int) bool {
	if len(first) != len(second) {
		return false
	}
	for groupIndex := range first {
		if len(first[groupIndex]) != len(second[groupIndex]) {
			return false
		}
		for memberIndex := range first[groupIndex] {
			if first[groupIndex][memberIndex] !=
				second[groupIndex][memberIndex] {
				return false
			}
		}
	}
	return true
}

func runOptimizerCalibrationArm(
	name, estimator string,
	sticky bool,
	target []float64,
	targetGroups [][]int,
) (OptimizerCalibrationArm, error) {
	current := continuousFusionInitialOffsets()
	union := newFusionUnion()

	initialCosine, err :=
		offsetCosineSimilarity(current, target)
	if err != nil {
		return OptimizerCalibrationArm{}, err
	}

	arm := OptimizerCalibrationArm{
		Name:                name,
		Estimator:           estimator,
		StickyProjection:    sticky,
		InitialTargetCosine: initialCosine,
		BestStep:            0,
		BestTargetCosine:    initialCosine,
		FirstTargetBasinStep: -1,
		Steps:               make([]OptimizerCalibrationStep, 0, directOffsetSteps),
	}

	for step := 1; step <= directOffsetSteps; step++ {
		inputGroups := union.groups()
		var direction []float64
		var gradient []float64

		switch estimator {
		case "sequential-projected-central-difference":
			direction = deterministicSPSADirection(
				step, inputGroups,
			)
			plusOffsets, err := perturbDirectOffsets(
				current, direction,
				directOffsetPerturbation, &union,
			)
			if err != nil {
				return OptimizerCalibrationArm{}, err
			}
			minusOffsets, err := perturbDirectOffsets(
				current, direction,
				-directOffsetPerturbation, &union,
			)
			if err != nil {
				return OptimizerCalibrationArm{}, err
			}
			plusObjective, err :=
				offsetCosineSimilarity(plusOffsets, target)
			if err != nil {
				return OptimizerCalibrationArm{}, err
			}
			minusObjective, err :=
				offsetCosineSimilarity(minusOffsets, target)
			if err != nil {
				return OptimizerCalibrationArm{}, err
			}
			gradient = directGradientEstimate(
				plusObjective, minusObjective, direction,
			)

		case "analytic-cosine-reference":
			gradient, err =
				analyticOffsetCosineGradient(current, target)
			if err != nil {
				return OptimizerCalibrationArm{}, err
			}
			if sticky {
				gradient, err =
					projectGradientToCalibrationGroups(
						gradient, inputGroups,
					)
				if err != nil {
					return OptimizerCalibrationArm{}, err
				}
			}

		default:
			return OptimizerCalibrationArm{}, fmt.Errorf(
				"unknown calibration estimator=%q", estimator,
			)
		}

		var next, update []float64
		if sticky {
			next, update, err =
				applyDirectGradient(current, gradient, &union)
		} else {
			next, update, err =
				applyCalibrationGradientNoSticky(
					current, gradient,
				)
		}
		if err != nil {
			return OptimizerCalibrationArm{}, err
		}

		outputGroups := union.groups()
		cosine, err := offsetCosineSimilarity(next, target)
		if err != nil {
			return OptimizerCalibrationArm{}, err
		}
		spread, err :=
			calibrationTargetGroupMaxSpread(
				next, targetGroups,
			)
		if err != nil {
			return OptimizerCalibrationArm{}, err
		}
		basin, err :=
			calibrationTargetFusionBasin(
				next, targetGroups,
			)
		if err != nil {
			return OptimizerCalibrationArm{}, err
		}
		wrongFusion, err :=
			calibrationWrongCrossTargetFusion(
				outputGroups, targetGroups,
			)
		if err != nil {
			return OptimizerCalibrationArm{}, err
		}

		arm.Steps = append(
			arm.Steps,
			OptimizerCalibrationStep{
				Step:                   step,
				InputOffsets:           append([]float64(nil), current...),
				InputGroups:            inputGroups,
				Direction:              append([]float64(nil), direction...),
				GradientEstimate:       append([]float64(nil), gradient...),
				Update:                 append([]float64(nil), update...),
				OutputOffsets:          append([]float64(nil), next...),
				OutputGroups:           outputGroups,
				TargetCosine:           cosine,
				TargetGroupMaxSpread:   spread,
				TargetFusionBasin:      basin,
				WrongCrossTargetFusion: wrongFusion,
			},
		)

		if cosine > arm.BestTargetCosine {
			arm.BestTargetCosine = cosine
			arm.BestStep = step
		}
		if basin && !arm.EnteredTargetFusionBasin {
			arm.EnteredTargetFusionBasin = true
			arm.FirstTargetBasinStep = step
		}
		if calibrationGroupsEqual(outputGroups, targetGroups) {
			arm.ExactTargetGroupingEver = true
		}
		if wrongFusion {
			arm.WrongCrossTargetFusionEver = true
		}
		current = next
	}

	finalCosine, err :=
		offsetCosineSimilarity(current, target)
	if err != nil {
		return OptimizerCalibrationArm{}, err
	}
	finalSpread, err :=
		calibrationTargetGroupMaxSpread(
			current, targetGroups,
		)
	if err != nil {
		return OptimizerCalibrationArm{}, err
	}
	arm.FinalTargetCosine = finalCosine
	arm.FinalTargetGroupMaxSpread = finalSpread
	arm.FinalGroups = union.groups()
	return arm, nil
}

func RunUP41() (
	OptimizerReachabilityCalibrationProbeResult,
	error,
) {
	targetSymmetryGroups, targetOffsets, err :=
		optimizerCalibrationTarget()
	if err != nil {
		return OptimizerReachabilityCalibrationProbeResult{}, err
	}
	targetGroups := symmetryGroupsJSON(targetSymmetryGroups)
	targetCapacity := symmetryGroupCapacity(targetSymmetryGroups)
	startOffsets := continuousFusionInitialOffsets()

	specs := []struct {
		name      string
		estimator string
		sticky    bool
	}{
		{
			name: "sequential_sticky_current",
			estimator: "sequential-projected-central-difference",
			sticky: true,
		},
		{
			name: "sequential_no_sticky",
			estimator: "sequential-projected-central-difference",
			sticky: false,
		},
		{
			name: "analytic_reference_sticky",
			estimator: "analytic-cosine-reference",
			sticky: true,
		},
		{
			name: "analytic_reference_no_sticky",
			estimator: "analytic-cosine-reference",
			sticky: false,
		},
	}

	arms := make([]OptimizerCalibrationArm, 0, len(specs))
	for _, spec := range specs {
		arm, err := runOptimizerCalibrationArm(
			spec.name, spec.estimator, spec.sticky,
			targetOffsets, targetGroups,
		)
		if err != nil {
			return OptimizerReachabilityCalibrationProbeResult{}, err
		}
		arms = append(arms, arm)
	}

	current := arms[0]
	sequentialNoSticky := arms[1]
	analyticSticky := arms[2]
	analyticNoSticky := arms[3]

	estimatorDelta :=
		analyticNoSticky.BestTargetCosine -
			sequentialNoSticky.BestTargetCosine
	estimatorTrails := estimatorDelta > 0
	stickyCrossTarget :=
		current.WrongCrossTargetFusionEver ||
			analyticSticky.WrongCrossTargetFusionEver
	reachabilityLimited :=
		!analyticNoSticky.EnteredTargetFusionBasin
	mechanicsImplicated :=
		reachabilityLimited ||
			estimatorTrails ||
			stickyCrossTarget

	return OptimizerReachabilityCalibrationProbeResult{
		Schema:                           OptimizerReachabilityCalibrationSchema,
		Experiment:                       "UP-41-direct-optimizer-reachability-calibration",
		CalibrationOnly:                  true,
		TaskLabelsUsed:                   false,
		HeldOutDataUsed:                  false,
		ArchitectureSelectionAuthorized:  false,
		TargetSource:                     optimizerCalibrationTargetSource,
		TargetGroups:                     targetGroups,
		TargetCapacity:                   targetCapacity,
		TargetOffsets:                    append([]float64(nil), targetOffsets...),
		StartOffsets:                     append([]float64(nil), startOffsets...),
		ZeroMeanConstraint:               true,
		FixedRMS:                         multiplicityDoseRMS,
		Perturbation:                     directOffsetPerturbation,
		LearningRate:                     directOffsetLearningRate,
		MaximumCoordinateUpdate:          directOffsetMaxUpdate,
		OptimizationSteps:                directOffsetSteps,
		ExactFusionTolerance:             continuousFusionTolerance,
		Arms:                             arms,
		Diagnosis: OptimizerCalibrationDiagnosis{
			CurrentOptimizerEntersTargetBasin:
				current.EnteredTargetFusionBasin,
			CurrentOptimizerExactTargetGrouping:
				current.ExactTargetGroupingEver,
			AnalyticReferenceEntersTargetBasin:
				analyticNoSticky.EnteredTargetFusionBasin,
			SequentialEstimatorTrailsAnalyticReference:
				estimatorTrails,
			StickyProjectionCrossTargetFusion:
				stickyCrossTarget,
			TwelveStepReachabilityLimited:
				reachabilityLimited,
			SequentialNoStickyBestCosine:
				sequentialNoSticky.BestTargetCosine,
			AnalyticNoStickyBestCosine:
				analyticNoSticky.BestTargetCosine,
			EstimatorBestCosineDelta:
				estimatorDelta,
			OptimizerMechanicsImplicated:
				mechanicsImplicated,
		},
	}, nil
}
