package unitary

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const GreedySymmetryInductionSchema = "wingless.greedy-symmetry-induction.v1"

const (
	greedySymmetryCapacityPrice = taskAllocationCapacityPrice
	greedySymmetryMinGain       = 0.005
)

type symmetryGroup struct {
	members []int
}

type GreedySymmetryStep struct {
	Step                     int                `json:"step"`
	CurrentGroups            [][]int            `json:"current_groups"`
	CurrentCapacity          int                `json:"current_capacity"`
	CurrentScore             float64            `json:"current_score"`
	CandidateMergeCount      int                `json:"candidate_merge_count"`
	BestMergeA               []int              `json:"best_merge_a"`
	BestMergeB               []int              `json:"best_merge_b"`
	BestCandidateGroups      [][]int            `json:"best_candidate_groups"`
	BestCandidateCapacity    int                `json:"best_candidate_capacity"`
	BestCandidateScore       float64            `json:"best_candidate_score"`
	ScoreGain                float64            `json:"score_gain"`
	Accepted                 bool               `json:"accepted"`
	BestValidationArm        MultiplicityDoseArm `json:"best_validation_arm"`
}

type GreedySymmetryDiagnosis struct {
	InnerSplitBalancedPass       bool    `json:"inner_split_balanced_pass"`
	SelectionUsesTrueHeldOut     bool    `json:"selection_uses_true_heldout"`
	AcceptedMergeCount           int     `json:"accepted_merge_count"`
	SelectedCapacity             int     `json:"selected_capacity"`
	StartCapacity                int     `json:"start_capacity"`
	MaxCapacity                  int     `json:"max_capacity"`
	CapacityCreatedFraction      float64 `json:"capacity_created_fraction"`
	StoppedByFrozenGainRule      bool    `json:"stopped_by_frozen_gain_rule"`
	SelectedHeldOutAccuracy      float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy       float64 `json:"selected_commit_accuracy"`
	FullCapacityHeldOutAccuracy  float64 `json:"full_capacity_heldout_accuracy"`
	FullCapacityCommitAccuracy   float64 `json:"full_capacity_commit_accuracy"`
	HeldOutRetentionDelta        float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta         float64 `json:"commit_retention_delta"`
	GreedyInductionPass          bool    `json:"greedy_induction_pass"`
}

type GreedySymmetryInductionProbeResult struct {
	Schema                         string                 `json:"schema"`
	Experiment                     string                 `json:"experiment"`
	LatentDimension                int                    `json:"latent_dimension"`
	RuntimeStateObjects            int                    `json:"runtime_state_objects"`
	FullCoordinateMixing           bool                   `json:"full_coordinate_mixing"`
	StartsFromMinimumSymmetry      bool                   `json:"starts_from_minimum_symmetry"`
	FinishedPartitionMenuProvided  bool                   `json:"finished_partition_menu_provided"`
	PairwiseMergeOnly              bool                   `json:"pairwise_merge_only"`
	MergeSelectionTrainingOnly     bool                   `json:"merge_selection_training_only"`
	SelectionUsesHeldOutData       bool                   `json:"selection_uses_heldout_data"`
	CapacityPrice                  float64                `json:"capacity_price"`
	MinimumAcceptedGain            float64                `json:"minimum_accepted_gain"`
	FitTables                      int                    `json:"fit_tables"`
	ValidationTables               int                    `json:"validation_tables"`
	TrueHeldOutTables              int                    `json:"true_heldout_tables"`
	InitialGroups                  [][]int                `json:"initial_groups"`
	Steps                          []GreedySymmetryStep   `json:"steps"`
	SelectedGroups                 [][]int                `json:"selected_groups"`
	SelectedValidationArm          MultiplicityDoseArm    `json:"selected_validation_arm"`
	SelectedFinalArm               MultiplicityDoseArm    `json:"selected_final_arm"`
	FullCapacityControl            MultiplicityDoseArm    `json:"full_capacity_control"`
	Diagnosis                      GreedySymmetryDiagnosis `json:"diagnosis"`
}

func initialSymmetryGroups() []symmetryGroup {
	out := make([]symmetryGroup, compositeChannels)
	for index := 0; index < compositeChannels; index++ {
		out[index] = symmetryGroup{members: []int{index}}
	}
	return out
}

func cloneSymmetryGroups(groups []symmetryGroup) []symmetryGroup {
	out := make([]symmetryGroup, len(groups))
	for i, group := range groups {
		out[i] = symmetryGroup{members: append([]int(nil), group.members...)}
	}
	return out
}

func normalizeSymmetryGroups(groups []symmetryGroup) ([]symmetryGroup, error) {
	if len(groups) < 1 {
		return nil, fmt.Errorf("symmetry groups empty")
	}
	seen := make(map[int]bool, compositeChannels)
	out := cloneSymmetryGroups(groups)
	for i := range out {
		if len(out[i].members) < 1 {
			return nil, fmt.Errorf("symmetry group empty")
		}
		sort.Ints(out[i].members)
		for _, member := range out[i].members {
			if member < 0 || member >= compositeChannels {
				return nil, fmt.Errorf("symmetry member out of range=%d", member)
			}
			if seen[member] {
				return nil, fmt.Errorf("symmetry member duplicated=%d", member)
			}
			seen[member] = true
		}
	}
	if len(seen) != compositeChannels {
		return nil, fmt.Errorf("symmetry groups cover=%d want=%d", len(seen), compositeChannels)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].members[0] < out[j].members[0]
	})
	return out, nil
}

func symmetryGroupsJSON(groups []symmetryGroup) [][]int {
	out := make([][]int, len(groups))
	for i, group := range groups {
		out[i] = append([]int(nil), group.members...)
	}
	return out
}

func symmetryGroupCapacity(groups []symmetryGroup) int {
	total := 0
	for _, group := range groups {
		total += len(group.members) * len(group.members)
	}
	return total
}

func symmetryGroupName(groups []symmetryGroup) string {
	parts := make([]string, len(groups))
	for i, group := range groups {
		members := make([]string, len(group.members))
		for j, member := range group.members {
			members[j] = fmt.Sprintf("%d", member)
		}
		parts[i] = strings.Join(members, "")
	}
	return strings.Join(parts, "_")
}

func offsetsForSymmetryGroups(
	groups []symmetryGroup,
	sigma float64,
) ([]float64, error) {
	normalized, err := normalizeSymmetryGroups(groups)
	if err != nil {
		return nil, err
	}
	if len(normalized) == 1 {
		return make([]float64, compositeChannels), nil
	}

	var weightedIndex float64
	for index, group := range normalized {
		weightedIndex += float64(index * len(group.members))
	}
	mean := weightedIndex / float64(compositeChannels)

	raw := make([]float64, len(normalized))
	var variance float64
	for index, group := range normalized {
		raw[index] = float64(index) - mean
		variance += float64(len(group.members)) * raw[index] * raw[index]
	}
	variance /= float64(compositeChannels)
	if variance <= 0 {
		return nil, fmt.Errorf("symmetry group variance is zero")
	}
	scale := sigma / math.Sqrt(variance)

	offsets := make([]float64, compositeChannels)
	for index, group := range normalized {
		value := raw[index] * scale
		for _, member := range group.members {
			offsets[member] = value
		}
	}
	if delta := math.Abs(doseOffsetRMS(offsets)-sigma); delta > 1e-12 {
		return nil, fmt.Errorf("symmetry group RMS delta=%g", delta)
	}
	var sum float64
	for _, value := range offsets {
		sum += value
	}
	if math.Abs(sum) > 1e-12 {
		return nil, fmt.Errorf("symmetry group offsets not zero mean sum=%g", sum)
	}
	return offsets, nil
}

func mergeSymmetryGroups(
	groups []symmetryGroup,
	first, second int,
) ([]symmetryGroup, error) {
	if first < 0 || second < 0 || first >= len(groups) || second >= len(groups) || first == second {
		return nil, fmt.Errorf("invalid symmetry merge")
	}
	if first > second {
		first, second = second, first
	}
	out := make([]symmetryGroup, 0, len(groups)-1)
	for index, group := range groups {
		switch index {
		case first:
			merged := append([]int(nil), group.members...)
			merged = append(merged, groups[second].members...)
			sort.Ints(merged)
			out = append(out, symmetryGroup{members: merged})
		case second:
			continue
		default:
			out = append(out, symmetryGroup{members: append([]int(nil), group.members...)})
		}
	}
	return normalizeSymmetryGroups(out)
}

func evaluateGreedySymmetryGrouping(
	groups []symmetryGroup,
	mixer latentMatrix,
	fitTables, validationTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (MultiplicityDoseArm, float64, error) {
	offsets, err := offsetsForSymmetryGroups(groups, multiplicityDoseRMS)
	if err != nil {
		return MultiplicityDoseArm{}, 0, err
	}
	name := "groups_" + symmetryGroupName(groups)
	arm, err := evaluateMultiplicityDoseArm(
		name,
		offsets,
		mixer,
		fitTables,
		validationTables,
		trainDepths,
		heldDepths,
		allDepths,
		memoryNoise,
	)
	if err != nil {
		return MultiplicityDoseArm{}, 0, err
	}
	capacity := symmetryGroupCapacity(groups)
	score := taskAllocationScore(
		arm.Static.HeldOutAccuracy,
		arm.Integration.CommitDecodeAccuracy,
		capacity,
	)
	return arm, score, nil
}

func structurallyValidGreedyArm(arm MultiplicityDoseArm) bool {
	rmsValid := math.Abs(arm.OffsetRMS-multiplicityDoseRMS) <= 1e-12
	if arm.MaxMultiplicity == compositeChannels {
		rmsValid = math.Abs(arm.OffsetRMS) <= 1e-12
	}
	return arm.OrthogonalityError <= 1e-10 &&
		arm.DiscoveryCommutatorError <= 1e-5 &&
		arm.FeatureDrift <= 5e-3 &&
		rmsValid
}

func RunUP35() (GreedySymmetryInductionProbeResult, error) {
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

	initialGroups := initialSymmetryGroups()
	currentGroups, err := normalizeSymmetryGroups(initialGroups)
	if err != nil {
		return GreedySymmetryInductionProbeResult{}, err
	}
	currentArm, currentScore, err := evaluateGreedySymmetryGrouping(
		currentGroups,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return GreedySymmetryInductionProbeResult{}, err
	}
	if !structurallyValidGreedyArm(currentArm) {
		return GreedySymmetryInductionProbeResult{}, fmt.Errorf("initial singleton arm structurally invalid")
	}

	steps := make([]GreedySymmetryStep, 0, compositeChannels-1)
	acceptedMerges := 0
	stoppedByGain := false

	for step := 0; len(currentGroups) > 1; step++ {
		bestScore := math.Inf(-1)
		var bestGroups []symmetryGroup
		var bestArm MultiplicityDoseArm
		var bestA, bestB []int
		candidateCount := 0

		for first := 0; first < len(currentGroups); first++ {
			for second := first + 1; second < len(currentGroups); second++ {
				candidateCount++
				candidateGroups, err := mergeSymmetryGroups(
					currentGroups, first, second,
				)
				if err != nil {
					return GreedySymmetryInductionProbeResult{}, err
				}
				arm, score, err := evaluateGreedySymmetryGrouping(
					candidateGroups,
					mixer,
					fitTables, validationTables,
					trainDepths, heldDepths, allDepths,
					memoryNoise,
				)
				if err != nil {
					return GreedySymmetryInductionProbeResult{}, err
				}
				if !structurallyValidGreedyArm(arm) {
					continue
				}
				if score > bestScore ||
					(score == bestScore &&
						symmetryGroupName(candidateGroups) < symmetryGroupName(bestGroups)) {
					bestScore = score
					bestGroups = candidateGroups
					bestArm = arm
					bestA = append([]int(nil), currentGroups[first].members...)
					bestB = append([]int(nil), currentGroups[second].members...)
				}
			}
		}

		if len(bestGroups) == 0 {
			return GreedySymmetryInductionProbeResult{}, fmt.Errorf(
				"no structurally valid merge candidate at step=%d",
				step,
			)
		}
		gain := bestScore - currentScore
		accepted := gain >= greedySymmetryMinGain
		steps = append(steps, GreedySymmetryStep{
			Step:                  step,
			CurrentGroups:         symmetryGroupsJSON(currentGroups),
			CurrentCapacity:       symmetryGroupCapacity(currentGroups),
			CurrentScore:          currentScore,
			CandidateMergeCount:   candidateCount,
			BestMergeA:            bestA,
			BestMergeB:            bestB,
			BestCandidateGroups:   symmetryGroupsJSON(bestGroups),
			BestCandidateCapacity: symmetryGroupCapacity(bestGroups),
			BestCandidateScore:    bestScore,
			ScoreGain:             gain,
			Accepted:              accepted,
			BestValidationArm:     bestArm,
		})
		if !accepted {
			stoppedByGain = true
			break
		}
		currentGroups = bestGroups
		currentArm = bestArm
		currentScore = bestScore
		acceptedMerges++
	}

	selectedOffsets, err := offsetsForSymmetryGroups(
		currentGroups, multiplicityDoseRMS,
	)
	if err != nil {
		return GreedySymmetryInductionProbeResult{}, err
	}
	selectedFinal, err := evaluateMultiplicityDoseArm(
		"selected_greedy_"+symmetryGroupName(currentGroups),
		selectedOffsets,
		mixer,
		allTrain,
		trueHeld,
		trainDepths,
		heldDepths,
		allDepths,
		memoryNoise,
	)
	if err != nil {
		return GreedySymmetryInductionProbeResult{}, err
	}
	fullControl, err := evaluateMultiplicityDoseArm(
		"full_capacity_control",
		[]float64{0, 0, 0, 0, 0, 0},
		mixer,
		allTrain,
		trueHeld,
		trainDepths,
		heldDepths,
		allDepths,
		memoryNoise,
	)
	if err != nil {
		return GreedySymmetryInductionProbeResult{}, err
	}

	selectedCapacity := symmetryGroupCapacity(currentGroups)
	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedFinal.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedFinal.Integration.CommitDecodeAccuracy
	createdFraction := float64(selectedCapacity-6) / float64(36-6)

	pass :=
		innerBalanced &&
			acceptedMerges >= 1 &&
			selectedCapacity > 6 &&
			selectedCapacity < 36 &&
			stoppedByGain &&
			selectedFinal.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedFinal.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedFinal.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedFinal.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	return GreedySymmetryInductionProbeResult{
		Schema:                        GreedySymmetryInductionSchema,
		Experiment:                    "UP-35-greedy-task-driven-symmetry-induction",
		LatentDimension:               fullLatentDimension,
		RuntimeStateObjects:           1,
		FullCoordinateMixing:          true,
		StartsFromMinimumSymmetry:     true,
		FinishedPartitionMenuProvided: false,
		PairwiseMergeOnly:             true,
		MergeSelectionTrainingOnly:    true,
		SelectionUsesHeldOutData:      false,
		CapacityPrice:                 greedySymmetryCapacityPrice,
		MinimumAcceptedGain:           greedySymmetryMinGain,
		FitTables:                     len(fitTables),
		ValidationTables:              len(validationTables),
		TrueHeldOutTables:             len(trueHeld),
		InitialGroups:                 symmetryGroupsJSON(initialGroups),
		Steps:                         steps,
		SelectedGroups:                symmetryGroupsJSON(currentGroups),
		SelectedValidationArm:         currentArm,
		SelectedFinalArm:              selectedFinal,
		FullCapacityControl:           fullControl,
		Diagnosis: GreedySymmetryDiagnosis{
			InnerSplitBalancedPass:      innerBalanced,
			SelectionUsesTrueHeldOut:    false,
			AcceptedMergeCount:          acceptedMerges,
			SelectedCapacity:            selectedCapacity,
			StartCapacity:               6,
			MaxCapacity:                 36,
			CapacityCreatedFraction:     createdFraction,
			StoppedByFrozenGainRule:     stoppedByGain,
			SelectedHeldOutAccuracy:     selectedFinal.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:      selectedFinal.Integration.CommitDecodeAccuracy,
			FullCapacityHeldOutAccuracy: fullControl.Static.HeldOutAccuracy,
			FullCapacityCommitAccuracy:  fullControl.Integration.CommitDecodeAccuracy,
			HeldOutRetentionDelta:       heldDelta,
			CommitRetentionDelta:        commitDelta,
			GreedyInductionPass:         pass,
		},
	}, nil
}
