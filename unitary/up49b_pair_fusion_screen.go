package unitary

import (
	"fmt"
	"math"
)

const UP49BPairFusionSchema = "wingless.up49b-single-pair-fusion-screen.v1"

type UP49BPairCandidate struct {
	Pair          []int                `json:"pair"`
	Offsets       []float64            `json:"offsets"`
	ExactCapacity int                  `json:"exact_capacity"`
	Smooth        RichDirectEvaluation `json:"smooth"`
	Hard          MultiplicityDoseArm  `json:"hard"`
	HardMean      float64              `json:"hard_mean"`
	HardGate      bool                 `json:"hard_gate"`
}

type UP49BDiagnosis struct {
	SelectedPair                 []int   `json:"selected_pair"`
	SelectedObjective            float64 `json:"selected_objective"`
	BaseObjective                float64 `json:"base_objective"`
	SelectedObjectiveGain        float64 `json:"selected_objective_gain"`
	SelectedHardMean             float64 `json:"selected_hard_mean"`
	SelectedHardGate             bool    `json:"selected_hard_gate"`
	PosthocBestHardPair          []int   `json:"posthoc_best_hard_pair"`
	PosthocBestHardMean          float64 `json:"posthoc_best_hard_mean"`
	SelectedMatchesPosthocBest   bool    `json:"selected_matches_posthoc_best"`
	SmoothHardPearsonCorrelation float64 `json:"smooth_hard_pearson_correlation"`
}

type UP49BPairFusionResult struct {
	Schema                   string                   `json:"schema"`
	Experiment               string                   `json:"experiment"`
	SourceUP48BSeal          string                   `json:"source_up48b_seal"`
	SourceStep               int                      `json:"source_step"`
	OptimizerRun             bool                     `json:"optimizer_run"`
	SelectionUsesHeldOutData bool                     `json:"selection_uses_heldout_data"`
	CandidateCount           int                      `json:"candidate_count"`
	BaseSmooth               RichDirectEvaluation     `json:"base_smooth"`
	BaseHard                 MultiplicityDoseArm      `json:"base_hard"`
	Candidates               []UP49BPairCandidate     `json:"candidates"`
	SelectedIndex            int                      `json:"selected_index"`
	FullCapacityControl      MultiplicityDoseArm      `json:"full_capacity_control"`
	Diagnosis                UP49BDiagnosis           `json:"diagnosis"`
}

func up49bHardMean(arm MultiplicityDoseArm) float64 {
	return 0.25 * (
		arm.Static.HeldOutAccuracy +
			arm.Integration.CommitDecodeAccuracy +
			arm.Integration.ExactFinalTableAccuracy +
			arm.Integration.RelationalQueryAccuracy)
}

func up49bHardGate(arm, full MultiplicityDoseArm) bool {
	heldDelta := full.Static.HeldOutAccuracy - arm.Static.HeldOutAccuracy
	commitDelta := full.Integration.CommitDecodeAccuracy - arm.Integration.CommitDecodeAccuracy
	return arm.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
		arm.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
		arm.Integration.ExactFinalTableAccuracy >= 0.90 &&
		arm.Integration.RelationalQueryAccuracy >= 0.95 &&
		heldDelta <= 0.02 &&
		commitDelta <= 0.05
}

func up49bPearson(xs, ys []float64) float64 {
	if len(xs) != len(ys) || len(xs) < 2 {
		return math.NaN()
	}
	var mx, my float64
	for i := range xs {
		mx += xs[i]
		my += ys[i]
	}
	mx /= float64(len(xs))
	my /= float64(len(ys))
	var num, dx, dy float64
	for i := range xs {
		a := xs[i] - mx
		b := ys[i] - my
		num += a * b
		dx += a * a
		dy += b * b
	}
	if dx <= 0 || dy <= 0 {
		return 0
	}
	return num / math.Sqrt(dx*dy)
}

func RunUP49B() (UP49BPairFusionResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	step42, err := up48bFrozenOffsets(42)
	if err != nil {
		return UP49BPairFusionResult{}, err
	}
	mixer := fullLatentMixer()
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fitTables, validationTables := splitTaskAllocationTrainingPool(allTrain)

	baseSmooth, err := evaluateSqrtFreeRunningDirectOffsets(
		"up49b_step42_base_smooth",
		append([]float64(nil), step42...),
		mixer, fitTables, validationTables,
		trainDepths, heldDepths, allDepths, memoryNoise,
	)
	if err != nil {
		return UP49BPairFusionResult{}, err
	}
	if !richEvaluationStructurallyValid(baseSmooth) {
		return UP49BPairFusionResult{}, fmt.Errorf("UP49B base smooth invalid")
	}
	baseHard, err := evaluateMultiplicityDoseArm(
		"up49b_step42_base_hard",
		append([]float64(nil), step42...),
		mixer, allTrain, trueHeld,
		trainDepths, heldDepths, allDepths, memoryNoise,
	)
	if err != nil {
		return UP49BPairFusionResult{}, err
	}
	fullControl, err := evaluateMultiplicityDoseArm(
		"up49b_full_capacity_control",
		[]float64{0, 0, 0, 0, 0, 0},
		mixer, allTrain, trueHeld,
		trainDepths, heldDepths, allDepths, memoryNoise,
	)
	if err != nil {
		return UP49BPairFusionResult{}, err
	}

	type pending struct {
		pair    []int
		offsets []float64
		smooth  RichDirectEvaluation
		cap     int
	}
	var pendings []pending
	selected := -1
	for i := 0; i < len(step42); i++ {
		for j := i + 1; j < len(step42); j++ {
			offsets := append([]float64(nil), step42...)
			center := 0.5 * (offsets[i] + offsets[j])
			offsets[i], offsets[j] = center, center
			smooth, err := evaluateSqrtFreeRunningDirectOffsets(
				fmt.Sprintf("up49b_pair_%d_%d_smooth", i, j),
				offsets,
				mixer, fitTables, validationTables,
				trainDepths, heldDepths, allDepths, memoryNoise,
			)
			if err != nil {
				return UP49BPairFusionResult{}, err
			}
			if !richEvaluationStructurallyValid(smooth) {
				return UP49BPairFusionResult{}, fmt.Errorf("UP49B smooth invalid pair %d,%d", i, j)
			}
			_, cap, err := exactOffsetGroups(offsets)
			if err != nil {
				return UP49BPairFusionResult{}, err
			}
			pendings = append(pendings, pending{
				pair: []int{i, j}, offsets: offsets, smooth: smooth, cap: cap,
			})
			idx := len(pendings) - 1
			if selected < 0 ||
				smooth.Objective > pendings[selected].smooth.Objective ||
				(smooth.Objective == pendings[selected].smooth.Objective &&
					smooth.SoftCapacity < pendings[selected].smooth.SoftCapacity) {
				selected = idx
			}
		}
	}
	if len(pendings) != 15 || selected < 0 {
		return UP49BPairFusionResult{}, fmt.Errorf("UP49B candidate construction failed count=%d selected=%d", len(pendings), selected)
	}

	candidates := make([]UP49BPairCandidate, 0, len(pendings))
	smoothScores := make([]float64, 0, len(pendings))
	hardScores := make([]float64, 0, len(pendings))
	posthocBest := -1
	for idx, p := range pendings {
		hard, err := evaluateMultiplicityDoseArm(
			fmt.Sprintf("up49b_pair_%d_%d_hard", p.pair[0], p.pair[1]),
			append([]float64(nil), p.offsets...),
			mixer, allTrain, trueHeld,
			trainDepths, heldDepths, allDepths, memoryNoise,
		)
		if err != nil {
			return UP49BPairFusionResult{}, err
		}
		mean := up49bHardMean(hard)
		gate := up49bHardGate(hard, fullControl)
		candidates = append(candidates, UP49BPairCandidate{
			Pair:          append([]int(nil), p.pair...),
			Offsets:       append([]float64(nil), p.offsets...),
			ExactCapacity: p.cap,
			Smooth:        p.smooth,
			Hard:          hard,
			HardMean:      mean,
			HardGate:      gate,
		})
		smoothScores = append(smoothScores, p.smooth.Objective)
		hardScores = append(hardScores, mean)
		if posthocBest < 0 || mean > hardScores[posthocBest] {
			posthocBest = idx
		}
	}
	if posthocBest < 0 {
		return UP49BPairFusionResult{}, fmt.Errorf("UP49B no posthoc hard candidate")
	}

	selectedCandidate := candidates[selected]
	posthoc := candidates[posthocBest]
	return UP49BPairFusionResult{
		Schema:                   UP49BPairFusionSchema,
		Experiment:               "UP-49B-step42-single-pair-fusion-screen",
		SourceUP48BSeal:          "3f9789035ac9324364ae4718bf900b07189cafce",
		SourceStep:               42,
		OptimizerRun:             false,
		SelectionUsesHeldOutData: false,
		CandidateCount:           len(candidates),
		BaseSmooth:               baseSmooth,
		BaseHard:                 baseHard,
		Candidates:               candidates,
		SelectedIndex:            selected,
		FullCapacityControl:      fullControl,
		Diagnosis: UP49BDiagnosis{
			SelectedPair:                 append([]int(nil), selectedCandidate.Pair...),
			SelectedObjective:            selectedCandidate.Smooth.Objective,
			BaseObjective:                baseSmooth.Objective,
			SelectedObjectiveGain:        selectedCandidate.Smooth.Objective - baseSmooth.Objective,
			SelectedHardMean:             selectedCandidate.HardMean,
			SelectedHardGate:             selectedCandidate.HardGate,
			PosthocBestHardPair:          append([]int(nil), posthoc.Pair...),
			PosthocBestHardMean:          posthoc.HardMean,
			SelectedMatchesPosthocBest:   selected == posthocBest,
			SmoothHardPearsonCorrelation: up49bPearson(smoothScores, hardScores),
		},
	}, nil
}
