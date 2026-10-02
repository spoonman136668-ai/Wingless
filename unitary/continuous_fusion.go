package unitary

import (
	"fmt"
	"math"
	"sort"
)

const ContinuousFusionSchema = "wingless.continuous-fusion.v1"

const (
	continuousFusionProbeFraction = 0.50
	continuousFusionStepRate      = 0.45
	continuousFusionTolerance     = 0.0025
	continuousFusionSteps         = 10
)

type FusionPairProbe struct {
	First              int                 `json:"first"`
	Second             int                 `json:"second"`
	ProbeOffsets       []float64           `json:"probe_offsets"`
	ProbeScore         float64             `json:"probe_score"`
	BaselineScore      float64             `json:"baseline_score"`
	RawGain            float64             `json:"raw_gain"`
	AttractionWeight   float64             `json:"attraction_weight"`
	ValidationArm      MultiplicityDoseArm `json:"validation_arm"`
}

type ContinuousFusionStep struct {
	Step                 int                 `json:"step"`
	Offsets              []float64           `json:"offsets"`
	Groups               [][]int             `json:"groups"`
	CommutantCapacity    int                 `json:"commutant_capacity"`
	ValidationScore      float64             `json:"validation_score"`
	StaticAccuracy       float64             `json:"static_accuracy"`
	CommitAccuracy       float64             `json:"commit_accuracy"`
	ValidationArm        MultiplicityDoseArm `json:"validation_arm"`
}

type ContinuousFusionDiagnosis struct {
	InnerSplitBalancedPass       bool    `json:"inner_split_balanced_pass"`
	SelectionUsesTrueHeldOut     bool    `json:"selection_uses_true_heldout"`
	PositiveAffinityCount        int     `json:"positive_affinity_count"`
	SelectedStep                 int     `json:"selected_step"`
	StartCapacity                int     `json:"start_capacity"`
	SelectedCapacity             int     `json:"selected_capacity"`
	MaxCapacity                  int     `json:"max_capacity"`
	ContinuousFusionOccurred     bool    `json:"continuous_fusion_occurred"`
	SelectedHeldOutAccuracy      float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy       float64 `json:"selected_commit_accuracy"`
	FullCapacityHeldOutAccuracy  float64 `json:"full_capacity_heldout_accuracy"`
	FullCapacityCommitAccuracy   float64 `json:"full_capacity_commit_accuracy"`
	HeldOutRetentionDelta        float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta         float64 `json:"commit_retention_delta"`
	ContinuousFusionPass         bool    `json:"continuous_fusion_pass"`
}

type ContinuousFusionProbeResult struct {
	Schema                         string                     `json:"schema"`
	Experiment                     string                     `json:"experiment"`
	LatentDimension                int                        `json:"latent_dimension"`
	RuntimeStateObjects            int                        `json:"runtime_state_objects"`
	FullCoordinateMixing           bool                       `json:"full_coordinate_mixing"`
	StartsFromIndependentOffsets   bool                       `json:"starts_from_independent_offsets"`
	FinishedPartitionMenuProvided  bool                       `json:"finished_partition_menu_provided"`
	HardMergeCandidatesEvaluated   bool                       `json:"hard_merge_candidates_evaluated"`
	PairProbesContinuousOnly       bool                       `json:"pair_probes_continuous_only"`
	FusionFlowSimultaneous         bool                       `json:"fusion_flow_simultaneous"`
	FusionStickyAfterTolerance     bool                       `json:"fusion_sticky_after_tolerance"`
	SelectionUsesHeldOutData       bool                       `json:"selection_uses_heldout_data"`
	CapacityPrice                  float64                    `json:"capacity_price"`
	ProbeFraction                  float64                    `json:"probe_fraction"`
	FlowStepRate                   float64                    `json:"flow_step_rate"`
	FusionTolerance                float64                    `json:"fusion_tolerance"`
	MaximumFlowSteps               int                        `json:"maximum_flow_steps"`
	FitTables                      int                        `json:"fit_tables"`
	ValidationTables               int                        `json:"validation_tables"`
	TrueHeldOutTables              int                        `json:"true_heldout_tables"`
	InitialOffsets                 []float64                  `json:"initial_offsets"`
	BaselineValidationArm          MultiplicityDoseArm        `json:"baseline_validation_arm"`
	PairProbes                     []FusionPairProbe          `json:"pair_probes"`
	FlowSteps                      []ContinuousFusionStep     `json:"flow_steps"`
	SelectedStep                   ContinuousFusionStep       `json:"selected_step"`
	SelectedFinalArm               MultiplicityDoseArm        `json:"selected_final_arm"`
	FullCapacityControl            MultiplicityDoseArm        `json:"full_capacity_control"`
	Diagnosis                      ContinuousFusionDiagnosis  `json:"diagnosis"`
}

type fusionUnion struct {
	parent [compositeChannels]int
}

func newFusionUnion() fusionUnion {
	var out fusionUnion
	for i := 0; i < compositeChannels; i++ {
		out.parent[i] = i
	}
	return out
}

func (u *fusionUnion) find(value int) int {
	root := value
	for u.parent[root] != root {
		root = u.parent[root]
	}
	for u.parent[value] != value {
		next := u.parent[value]
		u.parent[value] = root
		value = next
	}
	return root
}

func (u *fusionUnion) union(first, second int) {
	a := u.find(first)
	b := u.find(second)
	if a == b {
		return
	}
	if a < b {
		u.parent[b] = a
	} else {
		u.parent[a] = b
	}
}

func (u *fusionUnion) groups() [][]int {
	groupMap := make(map[int][]int)
	for i := 0; i < compositeChannels; i++ {
		root := u.find(i)
		groupMap[root] = append(groupMap[root], i)
	}
	roots := make([]int, 0, len(groupMap))
	for root := range groupMap {
		roots = append(roots, root)
	}
	sort.Ints(roots)
	out := make([][]int, 0, len(roots))
	for _, root := range roots {
		members := append([]int(nil), groupMap[root]...)
		sort.Ints(members)
		out = append(out, members)
	}
	return out
}

func fusionCapacity(groups [][]int) int {
	total := 0
	for _, group := range groups {
		total += len(group) * len(group)
	}
	return total
}

func normalizeContinuousOffsets(offsets []float64) ([]float64, error) {
	if len(offsets) != compositeChannels {
		return nil, fmt.Errorf(
			"continuous offset dimension=%d want=%d",
			len(offsets), compositeChannels,
		)
	}
	out := append([]float64(nil), offsets...)
	var mean float64
	for _, value := range out {
		if !finite(value) {
			return nil, fmt.Errorf("continuous offset is non-finite")
		}
		mean += value
	}
	mean /= float64(len(out))
	for i := range out {
		out[i] -= mean
	}
	rms := doseOffsetRMS(out)
	if rms <= 1e-15 {
		return make([]float64, compositeChannels), nil
	}
	scale := multiplicityDoseRMS / rms
	for i := range out {
		out[i] *= scale
	}
	return out, nil
}

func continuousFusionInitialOffsets() []float64 {
	specs := multiplicityDoseOffsets()
	return append([]float64(nil), specs[len(specs)-1].offsets...)
}

func continuousPairProbeOffsets(
	offsets []float64,
	first, second int,
	fraction float64,
) ([]float64, error) {
	if len(offsets) != compositeChannels ||
		first < 0 || second < 0 ||
		first >= compositeChannels || second >= compositeChannels ||
		first == second ||
		fraction <= 0 || fraction >= 1 {
		return nil, fmt.Errorf("invalid continuous pair probe")
	}
	out := append([]float64(nil), offsets...)
	midpoint := 0.5 * (out[first] + out[second])
	out[first] += fraction * (midpoint - out[first])
	out[second] += fraction * (midpoint - out[second])
	return normalizeContinuousOffsets(out)
}

func continuousArmScore(
	arm MultiplicityDoseArm,
	capacity int,
) float64 {
	return taskAllocationScore(
		arm.Static.HeldOutAccuracy,
		arm.Integration.CommitDecodeAccuracy,
		capacity,
	)
}

func normalizeFusionWeights(
	raw [compositeChannels][compositeChannels]float64,
) [compositeChannels][compositeChannels]float64 {
	var total float64
	for i := 0; i < compositeChannels; i++ {
		for j := i + 1; j < compositeChannels; j++ {
			if raw[i][j] > 0 {
				total += raw[i][j]
			}
		}
	}
	if total <= 0 {
		return [compositeChannels][compositeChannels]float64{}
	}
	var out [compositeChannels][compositeChannels]float64
	for i := 0; i < compositeChannels; i++ {
		for j := i + 1; j < compositeChannels; j++ {
			value := raw[i][j] / total
			out[i][j] = value
			out[j][i] = value
		}
	}
	return out
}

func applyFusionGroups(
	offsets []float64,
	union *fusionUnion,
) ([]float64, error) {
	if len(offsets) != compositeChannels {
		return nil, fmt.Errorf("fusion group offset dimension mismatch")
	}
	out := append([]float64(nil), offsets...)
	for _, group := range union.groups() {
		var mean float64
		for _, member := range group {
			mean += out[member]
		}
		mean /= float64(len(group))
		for _, member := range group {
			out[member] = mean
		}
	}
	return normalizeContinuousOffsets(out)
}

func continuousFusionAdvance(
	offsets []float64,
	weights [compositeChannels][compositeChannels]float64,
	union *fusionUnion,
) ([]float64, error) {
	if len(offsets) != compositeChannels {
		return nil, fmt.Errorf("fusion flow offset dimension mismatch")
	}
	proposal := append([]float64(nil), offsets...)
	for i := 0; i < compositeChannels; i++ {
		var delta float64
		for j := 0; j < compositeChannels; j++ {
			if i == j {
				continue
			}
			delta += weights[i][j] * (offsets[j] - offsets[i])
		}
		proposal[i] += continuousFusionStepRate * delta
	}
	normalized, err := normalizeContinuousOffsets(proposal)
	if err != nil {
		return nil, err
	}
	normalized, err = applyFusionGroups(normalized, union)
	if err != nil {
		return nil, err
	}

	for i := 0; i < compositeChannels; i++ {
		for j := i + 1; j < compositeChannels; j++ {
			if math.Abs(normalized[i]-normalized[j]) <= continuousFusionTolerance {
				union.union(i, j)
			}
		}
	}
	return applyFusionGroups(normalized, union)
}

func structurallyValidContinuousArm(arm MultiplicityDoseArm) bool {
	rmsValid := math.Abs(arm.OffsetRMS-multiplicityDoseRMS) <= 1e-12
	if arm.MaxMultiplicity == compositeChannels {
		rmsValid = math.Abs(arm.OffsetRMS) <= 1e-12
	}
	return arm.OrthogonalityError <= 1e-10 &&
		arm.DiscoveryCommutatorError <= 1e-5 &&
		arm.FeatureDrift <= 5e-3 &&
		rmsValid
}

func RunUP36() (ContinuousFusionProbeResult, error) {
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
	baselineArm, err := evaluateMultiplicityDoseArm(
		"continuous_singleton_baseline",
		initialOffsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return ContinuousFusionProbeResult{}, err
	}
	if !structurallyValidContinuousArm(baselineArm) {
		return ContinuousFusionProbeResult{}, fmt.Errorf(
			"continuous singleton baseline structurally invalid",
		)
	}
	baselineScore := continuousArmScore(baselineArm, 6)

	var rawWeights [compositeChannels][compositeChannels]float64
	pairProbes := make([]FusionPairProbe, 0, 15)
	positiveAffinityCount := 0

	for first := 0; first < compositeChannels; first++ {
		for second := first + 1; second < compositeChannels; second++ {
			probeOffsets, err := continuousPairProbeOffsets(
				initialOffsets,
				first, second,
				continuousFusionProbeFraction,
			)
			if err != nil {
				return ContinuousFusionProbeResult{}, err
			}
			arm, err := evaluateMultiplicityDoseArm(
				fmt.Sprintf("continuous_probe_%d_%d", first, second),
				probeOffsets,
				mixer,
				fitTables, validationTables,
				trainDepths, heldDepths, allDepths,
				memoryNoise,
			)
			if err != nil {
				return ContinuousFusionProbeResult{}, err
			}
			if !structurallyValidContinuousArm(arm) {
				return ContinuousFusionProbeResult{}, fmt.Errorf(
					"continuous pair probe %d,%d structurally invalid",
					first, second,
				)
			}
			score := continuousArmScore(arm, 6)
			gain := score - baselineScore
			weight := 0.0
			if gain > 0 {
				weight = gain
				positiveAffinityCount++
			}
			rawWeights[first][second] = weight
			rawWeights[second][first] = weight
			pairProbes = append(pairProbes, FusionPairProbe{
				First:            first,
				Second:           second,
				ProbeOffsets:     probeOffsets,
				ProbeScore:       score,
				BaselineScore:    baselineScore,
				RawGain:          gain,
				AttractionWeight: weight,
				ValidationArm:    arm,
			})
		}
	}

	weights := normalizeFusionWeights(rawWeights)
	union := newFusionUnion()
	currentOffsets := append([]float64(nil), initialOffsets...)
	steps := make([]ContinuousFusionStep, 0, continuousFusionSteps)

	selected := ContinuousFusionStep{
		Step:              -1,
		Offsets:           append([]float64(nil), initialOffsets...),
		Groups:            union.groups(),
		CommutantCapacity: 6,
		ValidationScore:   baselineScore,
		StaticAccuracy:    baselineArm.Static.HeldOutAccuracy,
		CommitAccuracy:    baselineArm.Integration.CommitDecodeAccuracy,
		ValidationArm:     baselineArm,
	}

	for step := 0; step < continuousFusionSteps; step++ {
		nextOffsets, err := continuousFusionAdvance(
			currentOffsets, weights, &union,
		)
		if err != nil {
			return ContinuousFusionProbeResult{}, err
		}
		currentOffsets = nextOffsets
		groups := union.groups()
		capacity := fusionCapacity(groups)

		arm, err := evaluateMultiplicityDoseArm(
			fmt.Sprintf("continuous_flow_%02d", step+1),
			currentOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return ContinuousFusionProbeResult{}, err
		}
		if !structurallyValidContinuousArm(arm) {
			return ContinuousFusionProbeResult{}, fmt.Errorf(
				"continuous flow step=%d structurally invalid",
				step+1,
			)
		}
		score := continuousArmScore(arm, capacity)
		entry := ContinuousFusionStep{
			Step:              step + 1,
			Offsets:           append([]float64(nil), currentOffsets...),
			Groups:            groups,
			CommutantCapacity: capacity,
			ValidationScore:   score,
			StaticAccuracy:    arm.Static.HeldOutAccuracy,
			CommitAccuracy:    arm.Integration.CommitDecodeAccuracy,
			ValidationArm:     arm,
		}
		steps = append(steps, entry)

		if score > selected.ValidationScore ||
			(score == selected.ValidationScore &&
				capacity < selected.CommutantCapacity) {
			selected = entry
		}
	}

	selectedFinal, err := evaluateMultiplicityDoseArm(
		fmt.Sprintf("selected_continuous_step_%02d", selected.Step),
		selected.Offsets,
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return ContinuousFusionProbeResult{}, err
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
		return ContinuousFusionProbeResult{}, err
	}

	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedFinal.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedFinal.Integration.CommitDecodeAccuracy

	fusionOccurred := selected.CommutantCapacity > 6
	pass :=
		innerBalanced &&
			positiveAffinityCount >= 1 &&
			selected.Step >= 1 &&
			fusionOccurred &&
			selected.CommutantCapacity < 36 &&
			selectedFinal.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedFinal.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedFinal.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedFinal.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	return ContinuousFusionProbeResult{
		Schema:                        ContinuousFusionSchema,
		Experiment:                    "UP-36-task-weighted-continuous-symmetry-fusion",
		LatentDimension:               fullLatentDimension,
		RuntimeStateObjects:           1,
		FullCoordinateMixing:          true,
		StartsFromIndependentOffsets:  true,
		FinishedPartitionMenuProvided: false,
		HardMergeCandidatesEvaluated:  false,
		PairProbesContinuousOnly:      true,
		FusionFlowSimultaneous:        true,
		FusionStickyAfterTolerance:    true,
		SelectionUsesHeldOutData:      false,
		CapacityPrice:                 taskAllocationCapacityPrice,
		ProbeFraction:                 continuousFusionProbeFraction,
		FlowStepRate:                  continuousFusionStepRate,
		FusionTolerance:               continuousFusionTolerance,
		MaximumFlowSteps:              continuousFusionSteps,
		FitTables:                     len(fitTables),
		ValidationTables:              len(validationTables),
		TrueHeldOutTables:             len(trueHeld),
		InitialOffsets:                append([]float64(nil), initialOffsets...),
		BaselineValidationArm:         baselineArm,
		PairProbes:                    pairProbes,
		FlowSteps:                     steps,
		SelectedStep:                  selected,
		SelectedFinalArm:              selectedFinal,
		FullCapacityControl:           fullControl,
		Diagnosis: ContinuousFusionDiagnosis{
			InnerSplitBalancedPass:      innerBalanced,
			SelectionUsesTrueHeldOut:    false,
			PositiveAffinityCount:       positiveAffinityCount,
			SelectedStep:                selected.Step,
			StartCapacity:               6,
			SelectedCapacity:            selected.CommutantCapacity,
			MaxCapacity:                 36,
			ContinuousFusionOccurred:    fusionOccurred,
			SelectedHeldOutAccuracy:     selectedFinal.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:      selectedFinal.Integration.CommitDecodeAccuracy,
			FullCapacityHeldOutAccuracy: fullControl.Static.HeldOutAccuracy,
			FullCapacityCommitAccuracy:  fullControl.Integration.CommitDecodeAccuracy,
			HeldOutRetentionDelta:       heldDelta,
			CommitRetentionDelta:        commitDelta,
			ContinuousFusionPass:        pass,
		},
	}, nil
}
