package unitary

import (
	"fmt"
	"math"
	"sort"
)

const TaskAllocatedSymmetrySchema = "wingless.task-allocated-symmetry.v1"

const (
	taskAllocationCapacityPrice = 0.02
	taskAllocationMinHeld       = 0.985
	taskAllocationMinCommit     = 0.95
)

type TaskAllocationCandidate struct {
	Name               string              `json:"name"`
	Partition          []int               `json:"partition"`
	Permutation        string              `json:"permutation"`
	MaxMultiplicity    int                 `json:"max_multiplicity"`
	CommutantCapacity  int                 `json:"commutant_capacity"`
	ValidationHeld     float64              `json:"validation_held_accuracy"`
	ValidationCommit   float64              `json:"validation_commit_accuracy"`
	ResourcePenalty    float64              `json:"resource_penalty"`
	AllocationScore    float64              `json:"allocation_score"`
	StructuralPass     bool                 `json:"structural_pass"`
	ValidationArm      MultiplicityDoseArm  `json:"validation_arm"`
}

type TaskAllocationDiagnosis struct {
	InnerSplitBalancedPass         bool    `json:"inner_split_balanced_pass"`
	SelectionUsesTrueHeldOut       bool    `json:"selection_uses_true_heldout"`
	SelectedBelowMaxCapacity       bool    `json:"selected_below_max_capacity"`
	SelectedCapacity               int     `json:"selected_capacity"`
	MaxAvailableCapacity           int     `json:"max_available_capacity"`
	CapacityReductionFraction      float64 `json:"capacity_reduction_fraction"`
	SelectedValidationScore        float64 `json:"selected_validation_score"`
	SelectedHeldOutAccuracy        float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy         float64 `json:"selected_commit_accuracy"`
	FullCapacityHeldOutAccuracy    float64 `json:"full_capacity_heldout_accuracy"`
	FullCapacityCommitAccuracy     float64 `json:"full_capacity_commit_accuracy"`
	HeldOutRetentionDelta          float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta           float64 `json:"commit_retention_delta"`
	TaskAllocationPass             bool    `json:"task_allocation_pass"`
}

type TaskAllocatedSymmetryProbeResult struct {
	Schema                          string                    `json:"schema"`
	Experiment                      string                    `json:"experiment"`
	LatentDimension                 int                       `json:"latent_dimension"`
	RuntimeStateObjects             int                       `json:"runtime_state_objects"`
	FullCoordinateMixing            bool                      `json:"full_coordinate_mixing"`
	CandidateMenuHandSpecified      bool                      `json:"candidate_menu_hand_specified"`
	AllocationUsesTrainingOnly      bool                      `json:"allocation_uses_training_only"`
	SelectionUsesHeldOutData        bool                      `json:"selection_uses_heldout_data"`
	CapacityPrice                   float64                   `json:"capacity_price"`
	FitTables                       int                       `json:"fit_tables"`
	ValidationTables                int                       `json:"validation_tables"`
	TrueHeldOutTables               int                       `json:"true_heldout_tables"`
	MinFitMarginalCount             int                       `json:"min_fit_marginal_count"`
	MinValidationMarginalCount      int                       `json:"min_validation_marginal_count"`
	Candidates                      []TaskAllocationCandidate `json:"candidates"`
	SelectedCandidate               TaskAllocationCandidate   `json:"selected_candidate"`
	SelectedFinalArm                MultiplicityDoseArm       `json:"selected_final_arm"`
	FullCapacityControl             MultiplicityDoseArm       `json:"full_capacity_control"`
	Diagnosis                       TaskAllocationDiagnosis   `json:"diagnosis"`
}

type taskAllocationSpec struct {
	name      string
	partition []int
}

func splitTaskAllocationTrainingPool(
	tables []memoryTable,
) ([]memoryTable, []memoryTable) {
	fit := make([]memoryTable, 0, len(tables)*3/4)
	validation := make([]memoryTable, 0, len(tables)/4)
	for index, table := range tables {
		if index%4 == 0 {
			validation = append(validation, table)
		} else {
			fit = append(fit, table)
		}
	}
	return fit, validation
}

func minimumTableMarginalCount(tables []memoryTable) int {
	if len(tables) == 0 {
		return 0
	}
	counts := [4][4]int{}
	for _, table := range tables {
		for entity := 0; entity < 4; entity++ {
			counts[entity][table[entity]]++
		}
	}
	minimum := len(tables)
	for entity := 0; entity < 4; entity++ {
		for value := 0; value < 4; value++ {
			if counts[entity][value] < minimum {
				minimum = counts[entity][value]
			}
		}
	}
	return minimum
}

func taskAllocationSpecs() []taskAllocationSpec {
	return []taskAllocationSpec{
		{name: "capacity6_singletons", partition: []int{1, 1, 1, 1, 1, 1}},
		{name: "capacity8_2_1_1_1_1", partition: []int{2, 1, 1, 1, 1}},
		{name: "capacity10_2_2_1_1", partition: []int{2, 2, 1, 1}},
		{name: "capacity12_3_1_1_1", partition: []int{3, 1, 1, 1}},
		{name: "capacity12_2_2_2", partition: []int{2, 2, 2}},
		{name: "capacity18_4_1_1", partition: []int{4, 1, 1}},
		{name: "capacity18_3_3", partition: []int{3, 3}},
		{name: "capacity36_6", partition: []int{6}},
	}
}

func offsetsForTaskAllocationSpec(
	partition []int,
	permutation commutantCapacityPermutation,
) ([]float64, error) {
	if len(partition) == 1 && partition[0] == compositeChannels {
		return []float64{0, 0, 0, 0, 0, 0}, nil
	}
	base, err := commutantPartitionOffsets(partition, multiplicityDoseRMS)
	if err != nil {
		return nil, err
	}
	return permuteCapacityOffsets(base, permutation)
}

func taskAllocationScore(
	held, commit float64,
	capacity int,
) float64 {
	performance := math.Min(held, commit)
	penalty := taskAllocationCapacityPrice *
		float64(capacity) / 36.0
	return performance - penalty
}

func evaluateTaskAllocationCandidate(
	spec taskAllocationSpec,
	permutation commutantCapacityPermutation,
	mixer latentMatrix,
	fitTables, validationTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (TaskAllocationCandidate, error) {
	offsets, err := offsetsForTaskAllocationSpec(
		spec.partition, permutation,
	)
	if err != nil {
		return TaskAllocationCandidate{}, err
	}
	name := spec.name + "_" + permutation.name
	if len(spec.partition) == 1 {
		name = spec.name
	}
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
		return TaskAllocationCandidate{}, err
	}

	capacity := commutantCapacityPartition(spec.partition)
	structural :=
		arm.OrthogonalityError <= 1e-10 &&
			arm.DiscoveryCommutatorError <= 1e-5 &&
			arm.FeatureDrift <= 5e-3
	if len(spec.partition) > 1 {
		structural = structural &&
			math.Abs(arm.OffsetRMS-multiplicityDoseRMS) <= 1e-12
	}

	penalty := taskAllocationCapacityPrice *
		float64(capacity) / 36.0
	return TaskAllocationCandidate{
		Name:              name,
		Partition:         append([]int(nil), spec.partition...),
		Permutation:       permutation.name,
		MaxMultiplicity:   doseMaxMultiplicityFromPartition(spec.partition),
		CommutantCapacity: capacity,
		ValidationHeld:    arm.Static.HeldOutAccuracy,
		ValidationCommit:  arm.Integration.CommitDecodeAccuracy,
		ResourcePenalty:   penalty,
		AllocationScore: taskAllocationScore(
			arm.Static.HeldOutAccuracy,
			arm.Integration.CommitDecodeAccuracy,
			capacity,
		),
		StructuralPass: structural,
		ValidationArm:  arm,
	}, nil
}

func chooseTaskAllocationCandidate(
	candidates []TaskAllocationCandidate,
) (TaskAllocationCandidate, error) {
	if len(candidates) == 0 {
		return TaskAllocationCandidate{}, fmt.Errorf(
			"task allocation has no candidates",
		)
	}
	sorted := append([]TaskAllocationCandidate(nil), candidates...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].StructuralPass != sorted[j].StructuralPass {
			return sorted[i].StructuralPass
		}
		if sorted[i].AllocationScore != sorted[j].AllocationScore {
			return sorted[i].AllocationScore > sorted[j].AllocationScore
		}
		if sorted[i].CommutantCapacity != sorted[j].CommutantCapacity {
			return sorted[i].CommutantCapacity < sorted[j].CommutantCapacity
		}
		return sorted[i].Name < sorted[j].Name
	})
	if !sorted[0].StructuralPass {
		return TaskAllocationCandidate{}, fmt.Errorf(
			"no structurally valid allocation candidate",
		)
	}
	return sorted[0], nil
}

func RunUP34() (TaskAllocatedSymmetryProbeResult, error) {
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

	permutations := commutantCapacityPermutations()
	specs := taskAllocationSpecs()
	candidates := make([]TaskAllocationCandidate, 0, 22)

	for _, spec := range specs {
		if len(spec.partition) == 1 {
			candidate, err := evaluateTaskAllocationCandidate(
				spec,
				permutations[0],
				mixer,
				fitTables, validationTables,
				trainDepths, heldDepths, allDepths,
				memoryNoise,
			)
			if err != nil {
				return TaskAllocatedSymmetryProbeResult{}, err
			}
			candidates = append(candidates, candidate)
			continue
		}
		for _, permutation := range permutations {
			candidate, err := evaluateTaskAllocationCandidate(
				spec,
				permutation,
				mixer,
				fitTables, validationTables,
				trainDepths, heldDepths, allDepths,
				memoryNoise,
			)
			if err != nil {
				return TaskAllocatedSymmetryProbeResult{}, err
			}
			candidates = append(candidates, candidate)
		}
	}

	selected, err := chooseTaskAllocationCandidate(candidates)
	if err != nil {
		return TaskAllocatedSymmetryProbeResult{}, err
	}

	var selectedPermutation commutantCapacityPermutation
	foundPermutation := false
	for _, permutation := range permutations {
		if permutation.name == selected.Permutation {
			selectedPermutation = permutation
			foundPermutation = true
			break
		}
	}
	if !foundPermutation {
		return TaskAllocatedSymmetryProbeResult{}, fmt.Errorf(
			"selected permutation not found: %s",
			selected.Permutation,
		)
	}
	selectedOffsets, err := offsetsForTaskAllocationSpec(
		selected.Partition, selectedPermutation,
	)
	if err != nil {
		return TaskAllocatedSymmetryProbeResult{}, err
	}

	selectedFinal, err := evaluateMultiplicityDoseArm(
		"selected_"+selected.Name,
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
		return TaskAllocatedSymmetryProbeResult{}, err
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
		return TaskAllocatedSymmetryProbeResult{}, err
	}

	maxCapacity := 36
	reduction := 1 -
		float64(selected.CommutantCapacity)/float64(maxCapacity)
	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedFinal.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedFinal.Integration.CommitDecodeAccuracy

	allocationPass :=
		innerBalanced &&
			selected.CommutantCapacity < maxCapacity &&
			selectedFinal.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedFinal.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedFinal.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedFinal.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	return TaskAllocatedSymmetryProbeResult{
		Schema:                         TaskAllocatedSymmetrySchema,
		Experiment:                     "UP-34-task-allocated-commutant-capacity",
		LatentDimension:                fullLatentDimension,
		RuntimeStateObjects:            1,
		FullCoordinateMixing:           true,
		CandidateMenuHandSpecified:     true,
		AllocationUsesTrainingOnly:     true,
		SelectionUsesHeldOutData:       false,
		CapacityPrice:                  taskAllocationCapacityPrice,
		FitTables:                      len(fitTables),
		ValidationTables:               len(validationTables),
		TrueHeldOutTables:              len(trueHeld),
		MinFitMarginalCount:            minFit,
		MinValidationMarginalCount:     minValidation,
		Candidates:                     candidates,
		SelectedCandidate:              selected,
		SelectedFinalArm:               selectedFinal,
		FullCapacityControl:            fullControl,
		Diagnosis: TaskAllocationDiagnosis{
			InnerSplitBalancedPass:      innerBalanced,
			SelectionUsesTrueHeldOut:    false,
			SelectedBelowMaxCapacity:    selected.CommutantCapacity < maxCapacity,
			SelectedCapacity:            selected.CommutantCapacity,
			MaxAvailableCapacity:        maxCapacity,
			CapacityReductionFraction:   reduction,
			SelectedValidationScore:     selected.AllocationScore,
			SelectedHeldOutAccuracy:     selectedFinal.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:      selectedFinal.Integration.CommitDecodeAccuracy,
			FullCapacityHeldOutAccuracy: fullControl.Static.HeldOutAccuracy,
			FullCapacityCommitAccuracy:  fullControl.Integration.CommitDecodeAccuracy,
			HeldOutRetentionDelta:       heldDelta,
			CommitRetentionDelta:        commitDelta,
			TaskAllocationPass:          allocationPass,
		},
	}, nil
}
