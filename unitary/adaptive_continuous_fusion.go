package unitary

import (
	"fmt"
	"math"
)

const AdaptiveContinuousFusionSchema = "wingless.adaptive-continuous-fusion.v1"

type AdaptiveFusionPairProbe struct {
	Step             int                 `json:"step"`
	FirstGroup       []int               `json:"first_group"`
	SecondGroup      []int               `json:"second_group"`
	ProbeOffsets     []float64           `json:"probe_offsets"`
	ProbeScore       float64             `json:"probe_score"`
	BaselineScore    float64             `json:"baseline_score"`
	RawGain          float64             `json:"raw_gain"`
	AttractionWeight float64             `json:"attraction_weight"`
	ValidationArm    MultiplicityDoseArm `json:"validation_arm"`
}

type AdaptiveFusionStep struct {
	Step                  int                 `json:"step"`
	InputOffsets          []float64           `json:"input_offsets"`
	InputGroups           [][]int             `json:"input_groups"`
	InputCapacity         int                 `json:"input_capacity"`
	BaselineScore         float64             `json:"baseline_score"`
	PairProbeCount        int                 `json:"pair_probe_count"`
	PositiveAffinityCount int                 `json:"positive_affinity_count"`
	OutputOffsets         []float64           `json:"output_offsets"`
	OutputGroups          [][]int             `json:"output_groups"`
	OutputCapacity        int                 `json:"output_capacity"`
	OutputScore           float64             `json:"output_score"`
	OutputStaticAccuracy  float64             `json:"output_static_accuracy"`
	OutputCommitAccuracy  float64             `json:"output_commit_accuracy"`
	OutputValidationArm   MultiplicityDoseArm `json:"output_validation_arm"`
}

type AdaptiveContinuousFusionDiagnosis struct {
	InnerSplitBalancedPass      bool    `json:"inner_split_balanced_pass"`
	SelectionUsesTrueHeldOut    bool    `json:"selection_uses_true_heldout"`
	AffinityReestimationCount   int     `json:"affinity_reestimation_count"`
	TotalPositiveAffinityCount  int     `json:"total_positive_affinity_count"`
	SelectedStep                int     `json:"selected_step"`
	StartCapacity               int     `json:"start_capacity"`
	SelectedCapacity            int     `json:"selected_capacity"`
	MaxCapacity                 int     `json:"max_capacity"`
	AdditionalFusionBeyondUP36  bool    `json:"additional_fusion_beyond_up36"`
	SelectedHeldOutAccuracy     float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy      float64 `json:"selected_commit_accuracy"`
	FullCapacityHeldOutAccuracy float64 `json:"full_capacity_heldout_accuracy"`
	FullCapacityCommitAccuracy  float64 `json:"full_capacity_commit_accuracy"`
	HeldOutRetentionDelta       float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta        float64 `json:"commit_retention_delta"`
	AdaptiveFusionPass          bool    `json:"adaptive_fusion_pass"`
}

type AdaptiveContinuousFusionProbeResult struct {
	Schema                          string                           `json:"schema"`
	Experiment                      string                           `json:"experiment"`
	LatentDimension                 int                              `json:"latent_dimension"`
	RuntimeStateObjects             int                              `json:"runtime_state_objects"`
	FullCoordinateMixing            bool                             `json:"full_coordinate_mixing"`
	StartsFromIndependentOffsets    bool                             `json:"starts_from_independent_offsets"`
	FinishedPartitionMenuProvided   bool                             `json:"finished_partition_menu_provided"`
	HardMergeCandidatesEvaluated    bool                             `json:"hard_merge_candidates_evaluated"`
	PairProbesContinuousOnly        bool                             `json:"pair_probes_continuous_only"`
	AffinityRecomputedEveryStep     bool                             `json:"affinity_recomputed_every_step"`
	FusionFlowSimultaneous          bool                             `json:"fusion_flow_simultaneous"`
	FusionStickyAfterTolerance      bool                             `json:"fusion_sticky_after_tolerance"`
	SelectionUsesHeldOutData        bool                             `json:"selection_uses_heldout_data"`
	CapacityPrice                   float64                          `json:"capacity_price"`
	ProbeFraction                   float64                          `json:"probe_fraction"`
	FlowStepRate                    float64                          `json:"flow_step_rate"`
	FusionTolerance                 float64                          `json:"fusion_tolerance"`
	MaximumFlowSteps                int                              `json:"maximum_flow_steps"`
	FitTables                       int                              `json:"fit_tables"`
	ValidationTables                int                              `json:"validation_tables"`
	TrueHeldOutTables               int                              `json:"true_heldout_tables"`
	InitialOffsets                  []float64                        `json:"initial_offsets"`
	PairProbes                      []AdaptiveFusionPairProbe        `json:"pair_probes"`
	FlowSteps                       []AdaptiveFusionStep             `json:"flow_steps"`
	SelectedStep                    AdaptiveFusionStep               `json:"selected_step"`
	SelectedFinalArm                MultiplicityDoseArm              `json:"selected_final_arm"`
	FullCapacityControl             MultiplicityDoseArm              `json:"full_capacity_control"`
	Diagnosis                       AdaptiveContinuousFusionDiagnosis `json:"diagnosis"`
}

func fusionGroupValue(
	offsets []float64,
	group []int,
) (float64, error) {
	if len(group) == 0 {
		return 0, fmt.Errorf("fusion group empty")
	}
	var total float64
	for _, member := range group {
		if member < 0 || member >= len(offsets) {
			return 0, fmt.Errorf("fusion group member out of range=%d", member)
		}
		total += offsets[member]
	}
	return total / float64(len(group)), nil
}

func adaptiveGroupProbeOffsets(
	offsets []float64,
	groups [][]int,
	first, second int,
	fraction float64,
) ([]float64, error) {
	if len(offsets) != compositeChannels ||
		first < 0 || second < 0 ||
		first >= len(groups) || second >= len(groups) ||
		first == second ||
		fraction <= 0 || fraction >= 1 {
		return nil, fmt.Errorf("invalid adaptive fusion probe")
	}

	firstValue, err := fusionGroupValue(offsets, groups[first])
	if err != nil {
		return nil, err
	}
	secondValue, err := fusionGroupValue(offsets, groups[second])
	if err != nil {
		return nil, err
	}
	midpoint := 0.5 * (firstValue + secondValue)
	nextFirst := firstValue + fraction*(midpoint-firstValue)
	nextSecond := secondValue + fraction*(midpoint-secondValue)

	out := append([]float64(nil), offsets...)
	for _, member := range groups[first] {
		out[member] = nextFirst
	}
	for _, member := range groups[second] {
		out[member] = nextSecond
	}
	return normalizeContinuousOffsets(out)
}

func normalizeAdaptiveGroupWeights(
	raw [][]float64,
) [][]float64 {
	out := make([][]float64, len(raw))
	var total float64
	for i := 0; i < len(raw); i++ {
		out[i] = make([]float64, len(raw))
		for j := i + 1; j < len(raw); j++ {
			if raw[i][j] > 0 {
				total += raw[i][j]
			}
		}
	}
	if total <= 0 {
		return out
	}
	for i := 0; i < len(raw); i++ {
		for j := i + 1; j < len(raw); j++ {
			value := raw[i][j] / total
			out[i][j] = value
			out[j][i] = value
		}
	}
	return out
}

func adaptiveFusionAdvance(
	offsets []float64,
	union *fusionUnion,
	groupWeights [][]float64,
) ([]float64, error) {
	groups := union.groups()
	if len(groupWeights) != len(groups) {
		return nil, fmt.Errorf("adaptive group weight dimension mismatch")
	}
	values := make([]float64, len(groups))
	for i, group := range groups {
		value, err := fusionGroupValue(offsets, group)
		if err != nil {
			return nil, err
		}
		values[i] = value
		if len(groupWeights[i]) != len(groups) {
			return nil, fmt.Errorf("adaptive group weight row mismatch")
		}
	}

	nextValues := append([]float64(nil), values...)
	for i := range groups {
		var delta float64
		for j := range groups {
			if i == j {
				continue
			}
			delta += groupWeights[i][j] * (values[j] - values[i])
		}
		nextValues[i] += continuousFusionStepRate * delta
	}

	proposal := append([]float64(nil), offsets...)
	for i, group := range groups {
		for _, member := range group {
			proposal[member] = nextValues[i]
		}
	}
	normalized, err := normalizeContinuousOffsets(proposal)
	if err != nil {
		return nil, err
	}
	normalized, err = applyFusionGroups(normalized, union)
	if err != nil {
		return nil, err
	}

	currentGroups := union.groups()
	currentValues := make([]float64, len(currentGroups))
	for i, group := range currentGroups {
		value, err := fusionGroupValue(normalized, group)
		if err != nil {
			return nil, err
		}
		currentValues[i] = value
	}
	for i := 0; i < len(currentGroups); i++ {
		for j := i + 1; j < len(currentGroups); j++ {
			if math.Abs(currentValues[i]-currentValues[j]) <= continuousFusionTolerance {
				union.union(currentGroups[i][0], currentGroups[j][0])
			}
		}
	}
	return applyFusionGroups(normalized, union)
}

func evaluateAdaptivePairField(
	step int,
	offsets []float64,
	union *fusionUnion,
	baselineScore float64,
	mixer latentMatrix,
	fitTables, validationTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (
	[]AdaptiveFusionPairProbe,
	[][]float64,
	int,
	error,
) {
	groups := union.groups()
	raw := make([][]float64, len(groups))
	for i := range raw {
		raw[i] = make([]float64, len(groups))
	}
	probes := make([]AdaptiveFusionPairProbe, 0, len(groups)*(len(groups)-1)/2)
	positive := 0

	for first := 0; first < len(groups); first++ {
		for second := first + 1; second < len(groups); second++ {
			probeOffsets, err := adaptiveGroupProbeOffsets(
				offsets, groups,
				first, second,
				continuousFusionProbeFraction,
			)
			if err != nil {
				return nil, nil, 0, err
			}

			arm, err := evaluateMultiplicityDoseArm(
				fmt.Sprintf("adaptive_probe_s%02d_g%d_g%d", step, first, second),
				probeOffsets,
				mixer,
				fitTables, validationTables,
				trainDepths, heldDepths, allDepths,
				memoryNoise,
			)
			if err != nil {
				return nil, nil, 0, err
			}
			if !structurallyValidContinuousArm(arm) {
				return nil, nil, 0, fmt.Errorf(
					"adaptive pair probe step=%d groups=%d,%d structurally invalid",
					step, first, second,
				)
			}

			capacity := fusionCapacity(groups)
			score := continuousArmScore(arm, capacity)
			gain := score - baselineScore
			weight := 0.0
			if gain > 0 {
				weight = gain
				positive++
			}
			raw[first][second] = weight
			raw[second][first] = weight
			probes = append(probes, AdaptiveFusionPairProbe{
				Step:             step,
				FirstGroup:       append([]int(nil), groups[first]...),
				SecondGroup:      append([]int(nil), groups[second]...),
				ProbeOffsets:     probeOffsets,
				ProbeScore:       score,
				BaselineScore:    baselineScore,
				RawGain:          gain,
				AttractionWeight: weight,
				ValidationArm:    arm,
			})
		}
	}
	return probes, normalizeAdaptiveGroupWeights(raw), positive, nil
}

func RunUP37() (AdaptiveContinuousFusionProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fitTables, validationTables := splitTaskAllocationTrainingPool(allTrain)

	minFit := minimumTableMarginalCount(fitTables)
	minValidation := minimumTableMarginalCount(validationTables)
	innerBalanced := minFit >= 12 && minValidation >= 4

	initialOffsets := continuousFusionInitialOffsets()
	union := newFusionUnion()
	currentOffsets := append([]float64(nil), initialOffsets...)

	currentArm, err := evaluateMultiplicityDoseArm(
		"adaptive_continuous_baseline",
		currentOffsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return AdaptiveContinuousFusionProbeResult{}, err
	}
	if !structurallyValidContinuousArm(currentArm) {
		return AdaptiveContinuousFusionProbeResult{}, fmt.Errorf(
			"adaptive continuous baseline structurally invalid",
		)
	}

	initialGroups := union.groups()
	currentCapacity := fusionCapacity(initialGroups)
	currentScore := continuousArmScore(currentArm, currentCapacity)

	selected := AdaptiveFusionStep{
		Step:                 0,
		InputOffsets:         append([]float64(nil), currentOffsets...),
		InputGroups:          initialGroups,
		InputCapacity:        currentCapacity,
		BaselineScore:        currentScore,
		PairProbeCount:       0,
		PositiveAffinityCount: 0,
		OutputOffsets:        append([]float64(nil), currentOffsets...),
		OutputGroups:         initialGroups,
		OutputCapacity:       currentCapacity,
		OutputScore:          currentScore,
		OutputStaticAccuracy: currentArm.Static.HeldOutAccuracy,
		OutputCommitAccuracy: currentArm.Integration.CommitDecodeAccuracy,
		OutputValidationArm:  currentArm,
	}

	allProbes := make([]AdaptiveFusionPairProbe, 0, 60)
	steps := make([]AdaptiveFusionStep, 0, continuousFusionSteps)
	totalPositive := 0
	reestimations := 0

	for step := 1; step <= continuousFusionSteps; step++ {
		inputGroups := union.groups()
		if len(inputGroups) <= 1 {
			break
		}
		inputCapacity := fusionCapacity(inputGroups)

		probes, weights, positive, err :=
			evaluateAdaptivePairField(
				step,
				currentOffsets,
				&union,
				currentScore,
				mixer,
				fitTables, validationTables,
				trainDepths, heldDepths, allDepths,
				memoryNoise,
			)
		if err != nil {
			return AdaptiveContinuousFusionProbeResult{}, err
		}
		reestimations++
		totalPositive += positive
		allProbes = append(allProbes, probes...)

		nextOffsets, err := adaptiveFusionAdvance(
			currentOffsets, &union, weights,
		)
		if err != nil {
			return AdaptiveContinuousFusionProbeResult{}, err
		}
		outputGroups := union.groups()
		outputCapacity := fusionCapacity(outputGroups)

		nextArm, err := evaluateMultiplicityDoseArm(
			fmt.Sprintf("adaptive_flow_%02d", step),
			nextOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return AdaptiveContinuousFusionProbeResult{}, err
		}
		if !structurallyValidContinuousArm(nextArm) {
			return AdaptiveContinuousFusionProbeResult{}, fmt.Errorf(
				"adaptive flow step=%d structurally invalid",
				step,
			)
		}
		nextScore := continuousArmScore(
			nextArm, outputCapacity,
		)

		entry := AdaptiveFusionStep{
			Step:                  step,
			InputOffsets:          append([]float64(nil), currentOffsets...),
			InputGroups:           inputGroups,
			InputCapacity:         inputCapacity,
			BaselineScore:         currentScore,
			PairProbeCount:        len(probes),
			PositiveAffinityCount: positive,
			OutputOffsets:         append([]float64(nil), nextOffsets...),
			OutputGroups:          outputGroups,
			OutputCapacity:        outputCapacity,
			OutputScore:           nextScore,
			OutputStaticAccuracy:  nextArm.Static.HeldOutAccuracy,
			OutputCommitAccuracy:  nextArm.Integration.CommitDecodeAccuracy,
			OutputValidationArm:   nextArm,
		}
		steps = append(steps, entry)

		if nextScore > selected.OutputScore ||
			(nextScore == selected.OutputScore &&
				outputCapacity < selected.OutputCapacity) {
			selected = entry
		}

		currentOffsets = nextOffsets
		currentArm = nextArm
		currentScore = nextScore
	}

	selectedFinal, err := evaluateMultiplicityDoseArm(
		fmt.Sprintf("selected_adaptive_step_%02d", selected.Step),
		selected.OutputOffsets,
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return AdaptiveContinuousFusionProbeResult{}, err
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
		return AdaptiveContinuousFusionProbeResult{}, err
	}

	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedFinal.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedFinal.Integration.CommitDecodeAccuracy
	additionalFusion := selected.OutputCapacity > 8

	pass :=
		innerBalanced &&
			reestimations >= 2 &&
			totalPositive >= 1 &&
			selected.Step >= 2 &&
			additionalFusion &&
			selected.OutputCapacity < 36 &&
			selectedFinal.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedFinal.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedFinal.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedFinal.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	return AdaptiveContinuousFusionProbeResult{
		Schema:                         AdaptiveContinuousFusionSchema,
		Experiment:                     "UP-37-adaptive-task-weighted-continuous-symmetry-fusion",
		LatentDimension:                fullLatentDimension,
		RuntimeStateObjects:            1,
		FullCoordinateMixing:           true,
		StartsFromIndependentOffsets:   true,
		FinishedPartitionMenuProvided:  false,
		HardMergeCandidatesEvaluated:   false,
		PairProbesContinuousOnly:       true,
		AffinityRecomputedEveryStep:    true,
		FusionFlowSimultaneous:         true,
		FusionStickyAfterTolerance:     true,
		SelectionUsesHeldOutData:       false,
		CapacityPrice:                  taskAllocationCapacityPrice,
		ProbeFraction:                  continuousFusionProbeFraction,
		FlowStepRate:                   continuousFusionStepRate,
		FusionTolerance:                continuousFusionTolerance,
		MaximumFlowSteps:               continuousFusionSteps,
		FitTables:                      len(fitTables),
		ValidationTables:               len(validationTables),
		TrueHeldOutTables:              len(trueHeld),
		InitialOffsets:                 append([]float64(nil), initialOffsets...),
		PairProbes:                     allProbes,
		FlowSteps:                      steps,
		SelectedStep:                   selected,
		SelectedFinalArm:               selectedFinal,
		FullCapacityControl:            fullControl,
		Diagnosis: AdaptiveContinuousFusionDiagnosis{
			InnerSplitBalancedPass:      innerBalanced,
			SelectionUsesTrueHeldOut:    false,
			AffinityReestimationCount:   reestimations,
			TotalPositiveAffinityCount:  totalPositive,
			SelectedStep:                selected.Step,
			StartCapacity:               6,
			SelectedCapacity:            selected.OutputCapacity,
			MaxCapacity:                 36,
			AdditionalFusionBeyondUP36:  additionalFusion,
			SelectedHeldOutAccuracy:     selectedFinal.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:      selectedFinal.Integration.CommitDecodeAccuracy,
			FullCapacityHeldOutAccuracy: fullControl.Static.HeldOutAccuracy,
			FullCapacityCommitAccuracy:  fullControl.Integration.CommitDecodeAccuracy,
			HeldOutRetentionDelta:       heldDelta,
			CommitRetentionDelta:        commitDelta,
			AdaptiveFusionPass:          pass,
		},
	}, nil
}
