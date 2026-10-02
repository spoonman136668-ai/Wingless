package unitary

import (
	"fmt"
	"math"
	"sort"
)

const CommutantCapacitySchema = "wingless.commutant-capacity-control.v1"

const (
	commutantCapacityHeldEquivalenceTolerance   = 0.03
	commutantCapacityCommitEquivalenceTolerance = 0.05
)

type CommutantCapacityVariant struct {
	Name               string              `json:"name"`
	Partition          []int               `json:"partition"`
	MaxMultiplicity    int                 `json:"max_multiplicity"`
	CommutantCapacity  int                 `json:"commutant_capacity"`
	Permutation        string              `json:"permutation"`
	Arm                MultiplicityDoseArm `json:"arm"`
}

type CommutantCapacityAggregate struct {
	Name                    string    `json:"name"`
	Partition               []int     `json:"partition"`
	MaxMultiplicity         int       `json:"max_multiplicity"`
	CommutantCapacity       int       `json:"commutant_capacity"`
	VariantCount            int       `json:"variant_count"`
	HeldOutAccuracies       []float64 `json:"held_out_accuracies"`
	CommitAccuracies        []float64 `json:"commit_accuracies"`
	MedianHeldOutAccuracy   float64   `json:"median_held_out_accuracy"`
	MeanHeldOutAccuracy     float64   `json:"mean_held_out_accuracy"`
	MedianCommitAccuracy    float64   `json:"median_commit_accuracy"`
	MeanCommitAccuracy      float64   `json:"mean_commit_accuracy"`
	MinimumHeldOutAccuracy  float64   `json:"minimum_held_out_accuracy"`
	MinimumCommitAccuracy   float64   `json:"minimum_commit_accuracy"`
}

type CommutantCapacityDiagnosis struct {
	StructuralValidityPass       bool    `json:"structural_validity_pass"`
	BaseControlPass              bool    `json:"base_control_pass"`
	Capacity18EquivalencePass    bool    `json:"capacity_18_equivalence_pass"`
	Capacity12EquivalencePass    bool    `json:"capacity_12_equivalence_pass"`
	HigherCapacityOrderingPass   bool    `json:"higher_capacity_ordering_pass"`
	CapacityHypothesisSupported  bool    `json:"capacity_hypothesis_supported"`
	Capacity18HeldDifference     float64 `json:"capacity_18_held_difference"`
	Capacity18CommitDifference   float64 `json:"capacity_18_commit_difference"`
	Capacity12HeldDifference     float64 `json:"capacity_12_held_difference"`
	Capacity12CommitDifference   float64 `json:"capacity_12_commit_difference"`
	MeanCapacity18Held           float64 `json:"mean_capacity_18_held"`
	MeanCapacity12Held           float64 `json:"mean_capacity_12_held"`
	MeanCapacity18Commit         float64 `json:"mean_capacity_18_commit"`
	MeanCapacity12Commit         float64 `json:"mean_capacity_12_commit"`
}

type CommutantCapacityProbeResult struct {
	Schema                           string                       `json:"schema"`
	Experiment                       string                       `json:"experiment"`
	LatentDimension                  int                          `json:"latent_dimension"`
	RuntimeStateObjects              int                          `json:"runtime_state_objects"`
	FullCoordinateMixing             bool                         `json:"full_coordinate_mixing"`
	SplitDirectionsAnonymous         bool                         `json:"split_directions_anonymous"`
	FixedNonzeroOffsetRMS            float64                      `json:"fixed_nonzero_offset_rms"`
	AllArmsNormPreserving            bool                         `json:"all_arms_norm_preserving"`
	AllArmsRealOrthogonalEquivalent  bool                         `json:"all_arms_real_orthogonal_equivalent"`
	AblationUsesKnownPartition       bool                         `json:"ablation_uses_known_partition"`
	ObserverDiscoveryFromTransport   bool                         `json:"observer_discovery_from_transport"`
	ObserverUsesKnownFactorization   bool                         `json:"observer_uses_known_factorization"`
	SelectorUsesTrainingLabels       bool                         `json:"selector_uses_training_labels"`
	SelectorUsesHeldOutData          bool                         `json:"selector_uses_heldout_data"`
	PhaseAlphabetSupervision         bool                         `json:"phase_alphabet_supervision"`
	CandidateObservableCount         int                          `json:"candidate_observable_count"`
	RuntimeObservableCount           int                          `json:"runtime_observable_count"`
	TrainingDepths                   []int                        `json:"training_depths"`
	HeldOutDepths                    []int                        `json:"held_out_depths"`
	PermutationNames                 []string                     `json:"permutation_names"`
	BaseControl                      MultiplicityDoseArm          `json:"base_control"`
	Variants                         []CommutantCapacityVariant   `json:"variants"`
	Aggregates                       []CommutantCapacityAggregate `json:"aggregates"`
	Diagnosis                        CommutantCapacityDiagnosis   `json:"diagnosis"`
}

type commutantCapacityPartitionSpec struct {
	name  string
	sizes []int
}

type commutantCapacityPermutation struct {
	name  string
	index [compositeChannels]int
}

func commutantCapacityPartition(sizes []int) int {
	total := 0
	for _, size := range sizes {
		total += size * size
	}
	return total
}

func commutantPartitionOffsets(
	sizes []int,
	sigma float64,
) ([]float64, error) {
	if len(sizes) < 2 {
		return nil, fmt.Errorf("capacity partition needs at least two groups")
	}
	total := 0
	var weightedIndex float64
	for index, size := range sizes {
		if size < 1 {
			return nil, fmt.Errorf("capacity partition contains nonpositive group")
		}
		total += size
		weightedIndex += float64(size * index)
	}
	if total != compositeChannels {
		return nil, fmt.Errorf(
			"capacity partition total=%d want=%d",
			total, compositeChannels,
		)
	}

	mean := weightedIndex / float64(total)
	raw := make([]float64, len(sizes))
	var variance float64
	for index, size := range sizes {
		raw[index] = float64(index) - mean
		variance += float64(size) * raw[index] * raw[index]
	}
	variance /= float64(total)
	if variance <= 0 {
		return nil, fmt.Errorf("capacity partition variance is zero")
	}
	scale := sigma / math.Sqrt(variance)

	offsets := make([]float64, 0, total)
	for index, size := range sizes {
		value := raw[index] * scale
		for count := 0; count < size; count++ {
			offsets = append(offsets, value)
		}
	}
	if len(offsets) != compositeChannels {
		return nil, fmt.Errorf("capacity offset length mismatch")
	}

	var sum float64
	for _, value := range offsets {
		sum += value
	}
	if math.Abs(sum) > 1e-12 {
		return nil, fmt.Errorf("capacity offsets are not zero mean: sum=%g", sum)
	}
	if math.Abs(doseOffsetRMS(offsets)-sigma) > 1e-12 {
		return nil, fmt.Errorf(
			"capacity offset RMS=%g want=%g",
			doseOffsetRMS(offsets), sigma,
		)
	}
	return offsets, nil
}

func commutantCapacityPermutations() []commutantCapacityPermutation {
	return []commutantCapacityPermutation{
		{
			name:  "identity",
			index: [compositeChannels]int{0, 1, 2, 3, 4, 5},
		},
		{
			name:  "rotate_two",
			index: [compositeChannels]int{2, 3, 4, 5, 0, 1},
		},
		{
			name:  "interleave",
			index: [compositeChannels]int{0, 3, 1, 4, 2, 5},
		},
	}
}

func permuteCapacityOffsets(
	offsets []float64,
	permutation commutantCapacityPermutation,
) ([]float64, error) {
	if len(offsets) != compositeChannels {
		return nil, fmt.Errorf("capacity permutation offset dimension mismatch")
	}
	seen := make(map[int]bool, compositeChannels)
	out := make([]float64, compositeChannels)
	for target, source := range permutation.index {
		if source < 0 || source >= compositeChannels {
			return nil, fmt.Errorf("capacity permutation source out of range")
		}
		if seen[source] {
			return nil, fmt.Errorf("capacity permutation duplicates source=%d", source)
		}
		seen[source] = true
		out[target] = offsets[source]
	}
	return out, nil
}

func medianFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	copyValues := append([]float64(nil), values...)
	sort.Float64s(copyValues)
	middle := len(copyValues) / 2
	if len(copyValues)%2 == 1 {
		return copyValues[middle]
	}
	return 0.5 * (copyValues[middle-1] + copyValues[middle])
}

func aggregateCapacityVariants(
	name string,
	partition []int,
	variants []CommutantCapacityVariant,
) (CommutantCapacityAggregate, error) {
	if len(variants) == 0 {
		return CommutantCapacityAggregate{}, fmt.Errorf(
			"capacity aggregate has no variants",
		)
	}
	held := make([]float64, 0, len(variants))
	commit := make([]float64, 0, len(variants))
	minHeld := math.Inf(1)
	minCommit := math.Inf(1)
	var meanHeld, meanCommit float64

	for _, variant := range variants {
		heldValue := variant.Arm.Static.HeldOutAccuracy
		commitValue := variant.Arm.Integration.CommitDecodeAccuracy
		held = append(held, heldValue)
		commit = append(commit, commitValue)
		meanHeld += heldValue
		meanCommit += commitValue
		if heldValue < minHeld {
			minHeld = heldValue
		}
		if commitValue < minCommit {
			minCommit = commitValue
		}
	}
	meanHeld /= float64(len(variants))
	meanCommit /= float64(len(variants))

	return CommutantCapacityAggregate{
		Name:                   name,
		Partition:              append([]int(nil), partition...),
		MaxMultiplicity:        doseMaxMultiplicityFromPartition(partition),
		CommutantCapacity:      commutantCapacityPartition(partition),
		VariantCount:           len(variants),
		HeldOutAccuracies:      held,
		CommitAccuracies:       commit,
		MedianHeldOutAccuracy:  medianFloat64(held),
		MeanHeldOutAccuracy:    meanHeld,
		MedianCommitAccuracy:   medianFloat64(commit),
		MeanCommitAccuracy:     meanCommit,
		MinimumHeldOutAccuracy: minHeld,
		MinimumCommitAccuracy:  minCommit,
	}, nil
}

func doseMaxMultiplicityFromPartition(partition []int) int {
	maximum := 0
	for _, size := range partition {
		if size > maximum {
			maximum = size
		}
	}
	return maximum
}

func RunUP33() (CommutantCapacityProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	baseControl, err := evaluateMultiplicityDoseArm(
		"multiplicity_6_control",
		[]float64{0, 0, 0, 0, 0, 0},
		mixer,
		trainTables, heldTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return CommutantCapacityProbeResult{}, err
	}

	specs := []commutantCapacityPartitionSpec{
		{name: "capacity18_m4_partition_4_1_1", sizes: []int{4, 1, 1}},
		{name: "capacity18_m3_partition_3_3", sizes: []int{3, 3}},
		{name: "capacity12_m3_partition_3_1_1_1", sizes: []int{3, 1, 1, 1}},
		{name: "capacity12_m2_partition_2_2_2", sizes: []int{2, 2, 2}},
	}
	permutations := commutantCapacityPermutations()

	variants := make([]CommutantCapacityVariant, 0, len(specs)*len(permutations))
	aggregates := make([]CommutantCapacityAggregate, 0, len(specs))
	structuralPass := true

	if baseControl.OrthogonalityError > 1e-10 ||
		baseControl.DiscoveryCommutatorError > 1e-5 ||
		baseControl.FeatureDrift > 5e-3 {
		structuralPass = false
	}

	for _, spec := range specs {
		baseOffsets, err := commutantPartitionOffsets(
			spec.sizes, multiplicityDoseRMS,
		)
		if err != nil {
			return CommutantCapacityProbeResult{}, err
		}
		specVariants := make([]CommutantCapacityVariant, 0, len(permutations))

		for _, permutation := range permutations {
			offsets, err := permuteCapacityOffsets(
				baseOffsets, permutation,
			)
			if err != nil {
				return CommutantCapacityProbeResult{}, err
			}
			armName := spec.name + "_" + permutation.name
			arm, err := evaluateMultiplicityDoseArm(
				armName, offsets,
				mixer,
				trainTables, heldTables,
				trainDepths, heldDepths, allDepths,
				memoryNoise,
			)
			if err != nil {
				return CommutantCapacityProbeResult{}, err
			}
			if arm.OrthogonalityError > 1e-10 ||
				arm.DiscoveryCommutatorError > 1e-5 ||
				arm.FeatureDrift > 5e-3 ||
				math.Abs(arm.OffsetRMS-multiplicityDoseRMS) > 1e-12 {
				structuralPass = false
			}
			variant := CommutantCapacityVariant{
				Name:              armName,
				Partition:         append([]int(nil), spec.sizes...),
				MaxMultiplicity:   doseMaxMultiplicityFromPartition(spec.sizes),
				CommutantCapacity: commutantCapacityPartition(spec.sizes),
				Permutation:       permutation.name,
				Arm:               arm,
			}
			variants = append(variants, variant)
			specVariants = append(specVariants, variant)
		}

		aggregate, err := aggregateCapacityVariants(
			spec.name, spec.sizes, specVariants,
		)
		if err != nil {
			return CommutantCapacityProbeResult{}, err
		}
		aggregates = append(aggregates, aggregate)
	}

	if len(aggregates) != 4 {
		return CommutantCapacityProbeResult{}, fmt.Errorf(
			"capacity aggregate count=%d want=4",
			len(aggregates),
		)
	}

	basePass :=
		baseControl.Static.HeldOutAccuracy >= 0.99 &&
			baseControl.Integration.CommitDecodeAccuracy >= 0.99 &&
			baseControl.Integration.ExactFinalTableAccuracy >= 0.95 &&
			baseControl.Integration.RelationalQueryAccuracy >= 0.95

	c18HeldDiff := math.Abs(
		aggregates[0].MedianHeldOutAccuracy -
			aggregates[1].MedianHeldOutAccuracy,
	)
	c18CommitDiff := math.Abs(
		aggregates[0].MedianCommitAccuracy -
			aggregates[1].MedianCommitAccuracy,
	)
	c12HeldDiff := math.Abs(
		aggregates[2].MedianHeldOutAccuracy -
			aggregates[3].MedianHeldOutAccuracy,
	)
	c12CommitDiff := math.Abs(
		aggregates[2].MedianCommitAccuracy -
			aggregates[3].MedianCommitAccuracy,
	)

	c18Equivalent :=
		c18HeldDiff <= commutantCapacityHeldEquivalenceTolerance &&
			c18CommitDiff <= commutantCapacityCommitEquivalenceTolerance
	c12Equivalent :=
		c12HeldDiff <= commutantCapacityHeldEquivalenceTolerance &&
			c12CommitDiff <= commutantCapacityCommitEquivalenceTolerance

	mean18Held := 0.5 * (
		aggregates[0].MeanHeldOutAccuracy +
			aggregates[1].MeanHeldOutAccuracy)
	mean12Held := 0.5 * (
		aggregates[2].MeanHeldOutAccuracy +
			aggregates[3].MeanHeldOutAccuracy)
	mean18Commit := 0.5 * (
		aggregates[0].MeanCommitAccuracy +
			aggregates[1].MeanCommitAccuracy)
	mean12Commit := 0.5 * (
		aggregates[2].MeanCommitAccuracy +
			aggregates[3].MeanCommitAccuracy)

	orderingPass :=
		mean18Held > mean12Held &&
			mean18Commit > mean12Commit

	capacitySupported :=
		structuralPass &&
			basePass &&
			c18Equivalent &&
			c12Equivalent &&
			orderingPass

	permutationNames := make([]string, 0, len(permutations))
	for _, permutation := range permutations {
		permutationNames = append(permutationNames, permutation.name)
	}

	return CommutantCapacityProbeResult{
		Schema:                          CommutantCapacitySchema,
		Experiment:                      "UP-33-equal-commutant-capacity-partition-control",
		LatentDimension:                 fullLatentDimension,
		RuntimeStateObjects:             1,
		FullCoordinateMixing:            true,
		SplitDirectionsAnonymous:        true,
		FixedNonzeroOffsetRMS:           multiplicityDoseRMS,
		AllArmsNormPreserving:           true,
		AllArmsRealOrthogonalEquivalent: true,
		AblationUsesKnownPartition:      true,
		ObserverDiscoveryFromTransport:  true,
		ObserverUsesKnownFactorization:  false,
		SelectorUsesTrainingLabels:      true,
		SelectorUsesHeldOutData:         false,
		PhaseAlphabetSupervision:        true,
		CandidateObservableCount:        interactionCandidateCount,
		RuntimeObservableCount:          interactionRuntimeCount,
		TrainingDepths:                  append([]int(nil), trainDepths...),
		HeldOutDepths:                   append([]int(nil), heldDepths...),
		PermutationNames:                permutationNames,
		BaseControl:                     baseControl,
		Variants:                        variants,
		Aggregates:                      aggregates,
		Diagnosis: CommutantCapacityDiagnosis{
			StructuralValidityPass:      structuralPass,
			BaseControlPass:             basePass,
			Capacity18EquivalencePass:   c18Equivalent,
			Capacity12EquivalencePass:   c12Equivalent,
			HigherCapacityOrderingPass:  orderingPass,
			CapacityHypothesisSupported: capacitySupported,
			Capacity18HeldDifference:    c18HeldDiff,
			Capacity18CommitDifference:  c18CommitDiff,
			Capacity12HeldDifference:    c12HeldDiff,
			Capacity12CommitDifference:  c12CommitDiff,
			MeanCapacity18Held:          mean18Held,
			MeanCapacity12Held:          mean12Held,
			MeanCapacity18Commit:        mean18Commit,
			MeanCapacity12Commit:        mean12Commit,
		},
	}, nil
}
