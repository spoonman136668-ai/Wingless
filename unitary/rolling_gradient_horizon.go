package unitary

import (
	"fmt"
	"math"
)

const RollingGradientHorizonSchema =
	"wingless.rolling-gradient-horizon.v1"

var rollingGradientHorizons = []int{12, 24, 48, 96}

const rollingGradientBufferSize = 7

type rollingDirectionalObservation struct {
	direction  []float64
	derivative float64
}

type RollingGradientHorizonRun struct {
	Estimator                   string  `json:"estimator"`
	Horizon                     int     `json:"horizon"`
	InitialTargetCosine         float64 `json:"initial_target_cosine"`
	BestStep                    int     `json:"best_step"`
	BestTargetCosine            float64 `json:"best_target_cosine"`
	FinalTargetCosine           float64 `json:"final_target_cosine"`
	FinalTargetGroupMaxSpread   float64 `json:"final_target_group_max_spread"`
	EnteredTargetFusionBasin    bool    `json:"entered_target_fusion_basin"`
	FirstTargetBasinStep        int     `json:"first_target_basin_step"`
	FullRankReconstructionSteps int     `json:"full_rank_reconstruction_steps"`
}

type RollingGradientHorizonDiagnosis struct {
	SequentialMinimumTestedHorizon int     `json:"sequential_minimum_tested_horizon"`
	RollingMinimumTestedHorizon    int     `json:"rolling_minimum_tested_horizon"`
	AnalyticMinimumTestedHorizon   int     `json:"analytic_minimum_tested_horizon"`
	SequentialFirstBasinStep       int     `json:"sequential_first_basin_step"`
	RollingFirstBasinStep          int     `json:"rolling_first_basin_step"`
	AnalyticFirstBasinStep         int     `json:"analytic_first_basin_step"`
	RollingMatchesAnalyticHorizon  bool    `json:"rolling_matches_analytic_horizon"`
	RollingBeatsSequentialAtMax    bool    `json:"rolling_beats_sequential_at_max"`
	SequentialMaxBestCosine        float64 `json:"sequential_max_best_cosine"`
	RollingMaxBestCosine           float64 `json:"rolling_max_best_cosine"`
	AnalyticMaxBestCosine          float64 `json:"analytic_max_best_cosine"`
	RollingEfficientCandidate      bool    `json:"rolling_efficient_candidate"`
}

type RollingGradientHorizonProbeResult struct {
	Schema                          string                            `json:"schema"`
	Experiment                      string                            `json:"experiment"`
	CalibrationOnly                 bool                              `json:"calibration_only"`
	TaskLabelsUsed                  bool                              `json:"task_labels_used"`
	HeldOutDataUsed                 bool                              `json:"heldout_data_used"`
	StickyProjectionUsed            bool                              `json:"sticky_projection_used"`
	ArchitectureSelectionAuthorized bool                              `json:"architecture_selection_authorized"`
	TargetSource                    string                            `json:"target_source"`
	TargetGroups                    [][]int                           `json:"target_groups"`
	TargetCapacity                  int                               `json:"target_capacity"`
	TargetOffsets                   []float64                         `json:"target_offsets"`
	Horizons                        []int                             `json:"horizons"`
	RollingBufferSize               int                               `json:"rolling_buffer_size"`
	Perturbation                    float64                           `json:"perturbation"`
	LearningRate                    float64                           `json:"learning_rate"`
	MaximumCoordinateUpdate         float64                           `json:"maximum_coordinate_update"`
	ExactFusionTolerance            float64                           `json:"exact_fusion_tolerance"`
	Runs                            []RollingGradientHorizonRun        `json:"runs"`
	Diagnosis                       RollingGradientHorizonDiagnosis    `json:"diagnosis"`
}

func calibrationDirectionalObservation(
	offsets, target []float64,
	step int,
) (rollingDirectionalObservation, []float64, error) {
	union := newFusionUnion()
	direction := deterministicSPSADirection(
		step, union.groups(),
	)
	plus, err := perturbDirectOffsets(
		offsets, direction,
		directOffsetPerturbation, &union,
	)
	if err != nil {
		return rollingDirectionalObservation{}, nil, err
	}
	minus, err := perturbDirectOffsets(
		offsets, direction,
		-directOffsetPerturbation, &union,
	)
	if err != nil {
		return rollingDirectionalObservation{}, nil, err
	}
	plusObjective, err :=
		offsetCosineSimilarity(plus, target)
	if err != nil {
		return rollingDirectionalObservation{}, nil, err
	}
	minusObjective, err :=
		offsetCosineSimilarity(minus, target)
	if err != nil {
		return rollingDirectionalObservation{}, nil, err
	}
	derivative :=
		(plusObjective - minusObjective) /
			(2 * directOffsetPerturbation)
	rankOne := directGradientEstimate(
		plusObjective, minusObjective, direction,
	)
	return rollingDirectionalObservation{
		direction:  append([]float64(nil), direction...),
		derivative: derivative,
	}, rankOne, nil
}

func solveFiveByFive(
	matrix [5][5]float64,
	right [5]float64,
) ([5]float64, bool) {
	var augmented [5][6]float64
	for row := 0; row < 5; row++ {
		for column := 0; column < 5; column++ {
			augmented[row][column] = matrix[row][column]
		}
		augmented[row][5] = right[row]
	}

	for column := 0; column < 5; column++ {
		pivot := column
		pivotMagnitude := math.Abs(augmented[pivot][column])
		for row := column + 1; row < 5; row++ {
			magnitude := math.Abs(augmented[row][column])
			if magnitude > pivotMagnitude {
				pivot = row
				pivotMagnitude = magnitude
			}
		}
		if pivotMagnitude <= 1e-12 {
			return [5]float64{}, false
		}
		if pivot != column {
			augmented[pivot], augmented[column] =
				augmented[column], augmented[pivot]
		}

		scale := augmented[column][column]
		for c := column; c < 6; c++ {
			augmented[column][c] /= scale
		}
		for row := 0; row < 5; row++ {
			if row == column {
				continue
			}
			factor := augmented[row][column]
			for c := column; c < 6; c++ {
				augmented[row][c] -=
					factor * augmented[column][c]
			}
		}
	}

	var solution [5]float64
	for row := 0; row < 5; row++ {
		solution[row] = augmented[row][5]
		if !finite(solution[row]) {
			return [5]float64{}, false
		}
	}
	return solution, true
}

func rollingFullRankGradient(
	observations []rollingDirectionalObservation,
) ([]float64, bool, error) {
	if len(observations) < 5 {
		return nil, false, nil
	}
	start := 0
	if len(observations) > rollingGradientBufferSize {
		start = len(observations) - rollingGradientBufferSize
	}
	window := observations[start:]

	var normal [5][5]float64
	var right [5]float64
	for _, observation := range window {
		if len(observation.direction) != compositeChannels ||
			!finite(observation.derivative) {
			return nil, false, fmt.Errorf(
				"rolling gradient observation invalid",
			)
		}
		var row [5]float64
		for coordinate := 0; coordinate < 5; coordinate++ {
			row[coordinate] =
				observation.direction[coordinate] -
					observation.direction[5]
		}
		for first := 0; first < 5; first++ {
			right[first] +=
				row[first] * observation.derivative
			for second := 0; second < 5; second++ {
				normal[first][second] +=
					row[first] * row[second]
			}
		}
	}

	coordinates, ok := solveFiveByFive(normal, right)
	if !ok {
		return nil, false, nil
	}
	gradient := make([]float64, compositeChannels)
	var sum float64
	for coordinate := 0; coordinate < 5; coordinate++ {
		gradient[coordinate] = coordinates[coordinate]
		sum += coordinates[coordinate]
	}
	gradient[5] = -sum

	for _, value := range gradient {
		if !finite(value) {
			return nil, false, fmt.Errorf(
				"rolling gradient reconstruction non-finite",
			)
		}
	}
	return gradient, true, nil
}

func runRollingHorizonCalibration(
	estimator string,
	horizon int,
	target []float64,
	targetGroups [][]int,
) (RollingGradientHorizonRun, error) {
	if horizon < 1 {
		return RollingGradientHorizonRun{}, fmt.Errorf(
			"rolling horizon invalid",
		)
	}
	current := continuousFusionInitialOffsets()
	initialCosine, err :=
		offsetCosineSimilarity(current, target)
	if err != nil {
		return RollingGradientHorizonRun{}, err
	}

	run := RollingGradientHorizonRun{
		Estimator:            estimator,
		Horizon:              horizon,
		InitialTargetCosine:  initialCosine,
		BestTargetCosine:     initialCosine,
		FirstTargetBasinStep: -1,
	}
	observations := make(
		[]rollingDirectionalObservation, 0,
		rollingGradientBufferSize,
	)

	for step := 1; step <= horizon; step++ {
		var gradient []float64
		switch estimator {
		case "sequential-projected-central-difference":
			_, rankOne, err :=
				calibrationDirectionalObservation(
					current, target, step,
				)
			if err != nil {
				return RollingGradientHorizonRun{}, err
			}
			gradient = rankOne

		case "rolling-full-rank-directional-memory":
			observation, rankOne, err :=
				calibrationDirectionalObservation(
					current, target, step,
				)
			if err != nil {
				return RollingGradientHorizonRun{}, err
			}
			observations = append(observations, observation)
			if len(observations) > rollingGradientBufferSize {
				observations = observations[
					len(observations)-rollingGradientBufferSize:
				]
			}
			reconstructed, fullRank, err :=
				rollingFullRankGradient(observations)
			if err != nil {
				return RollingGradientHorizonRun{}, err
			}
			if fullRank {
				gradient = reconstructed
				run.FullRankReconstructionSteps++
			} else {
				gradient = rankOne
			}

		case "analytic-cosine-reference":
			gradient, err =
				analyticOffsetCosineGradient(current, target)
			if err != nil {
				return RollingGradientHorizonRun{}, err
			}

		default:
			return RollingGradientHorizonRun{}, fmt.Errorf(
				"unknown rolling estimator=%q", estimator,
			)
		}

		next, _, err :=
			applyCalibrationGradientNoSticky(
				current, gradient,
			)
		if err != nil {
			return RollingGradientHorizonRun{}, err
		}
		cosine, err :=
			offsetCosineSimilarity(next, target)
		if err != nil {
			return RollingGradientHorizonRun{}, err
		}
		basin, err :=
			calibrationTargetFusionBasin(
				next, targetGroups,
			)
		if err != nil {
			return RollingGradientHorizonRun{}, err
		}
		if cosine > run.BestTargetCosine {
			run.BestTargetCosine = cosine
			run.BestStep = step
		}
		if basin && !run.EnteredTargetFusionBasin {
			run.EnteredTargetFusionBasin = true
			run.FirstTargetBasinStep = step
		}
		current = next
	}

	run.FinalTargetCosine, err =
		offsetCosineSimilarity(current, target)
	if err != nil {
		return RollingGradientHorizonRun{}, err
	}
	run.FinalTargetGroupMaxSpread, err =
		calibrationTargetGroupMaxSpread(
			current, targetGroups,
		)
	if err != nil {
		return RollingGradientHorizonRun{}, err
	}
	return run, nil
}

func minimumBasinHorizon(
	runs []RollingGradientHorizonRun,
	estimator string,
) (int, int) {
	minimumHorizon := 0
	firstStep := -1
	for _, run := range runs {
		if run.Estimator != estimator ||
			!run.EnteredTargetFusionBasin {
			continue
		}
		if minimumHorizon == 0 ||
			run.Horizon < minimumHorizon {
			minimumHorizon = run.Horizon
			firstStep = run.FirstTargetBasinStep
		}
	}
	return minimumHorizon, firstStep
}

func maxHorizonBestCosine(
	runs []RollingGradientHorizonRun,
	estimator string,
) float64 {
	maximumHorizon := -1
	best := 0.0
	for _, run := range runs {
		if run.Estimator != estimator {
			continue
		}
		if run.Horizon > maximumHorizon {
			maximumHorizon = run.Horizon
			best = run.BestTargetCosine
		}
	}
	return best
}

func RunUP42() (RollingGradientHorizonProbeResult, error) {
	targetSymmetryGroups, targetOffsets, err :=
		optimizerCalibrationTarget()
	if err != nil {
		return RollingGradientHorizonProbeResult{}, err
	}
	targetGroups := symmetryGroupsJSON(targetSymmetryGroups)
	estimators := []string{
		"sequential-projected-central-difference",
		"rolling-full-rank-directional-memory",
		"analytic-cosine-reference",
	}

	runs := make(
		[]RollingGradientHorizonRun, 0,
		len(estimators)*len(rollingGradientHorizons),
	)
	for _, estimator := range estimators {
		for _, horizon := range rollingGradientHorizons {
			run, err := runRollingHorizonCalibration(
				estimator, horizon,
				targetOffsets, targetGroups,
			)
			if err != nil {
				return RollingGradientHorizonProbeResult{}, err
			}
			runs = append(runs, run)
		}
	}

	sequentialHorizon, sequentialStep :=
		minimumBasinHorizon(
			runs,
			"sequential-projected-central-difference",
		)
	rollingHorizon, rollingStep :=
		minimumBasinHorizon(
			runs,
			"rolling-full-rank-directional-memory",
		)
	analyticHorizon, analyticStep :=
		minimumBasinHorizon(
			runs,
			"analytic-cosine-reference",
		)

	sequentialMax := maxHorizonBestCosine(
		runs,
		"sequential-projected-central-difference",
	)
	rollingMax := maxHorizonBestCosine(
		runs,
		"rolling-full-rank-directional-memory",
	)
	analyticMax := maxHorizonBestCosine(
		runs,
		"analytic-cosine-reference",
	)

	matchesReference :=
		rollingHorizon > 0 &&
			rollingHorizon == analyticHorizon
	rollingCandidate :=
		rollingHorizon > 0 &&
			(sequentialHorizon == 0 ||
				rollingHorizon < sequentialHorizon)

	return RollingGradientHorizonProbeResult{
		Schema:                          RollingGradientHorizonSchema,
		Experiment:                      "UP-42-rolling-full-rank-gradient-horizon-calibration",
		CalibrationOnly:                 true,
		TaskLabelsUsed:                  false,
		HeldOutDataUsed:                 false,
		StickyProjectionUsed:            false,
		ArchitectureSelectionAuthorized: false,
		TargetSource:                    optimizerCalibrationTargetSource,
		TargetGroups:                    targetGroups,
		TargetCapacity:                  symmetryGroupCapacity(targetSymmetryGroups),
		TargetOffsets:                   append([]float64(nil), targetOffsets...),
		Horizons:                        append([]int(nil), rollingGradientHorizons...),
		RollingBufferSize:               rollingGradientBufferSize,
		Perturbation:                    directOffsetPerturbation,
		LearningRate:                    directOffsetLearningRate,
		MaximumCoordinateUpdate:         directOffsetMaxUpdate,
		ExactFusionTolerance:            continuousFusionTolerance,
		Runs:                            runs,
		Diagnosis: RollingGradientHorizonDiagnosis{
			SequentialMinimumTestedHorizon: sequentialHorizon,
			RollingMinimumTestedHorizon:    rollingHorizon,
			AnalyticMinimumTestedHorizon:   analyticHorizon,
			SequentialFirstBasinStep:       sequentialStep,
			RollingFirstBasinStep:          rollingStep,
			AnalyticFirstBasinStep:         analyticStep,
			RollingMatchesAnalyticHorizon:  matchesReference,
			RollingBeatsSequentialAtMax:    rollingMax > sequentialMax,
			SequentialMaxBestCosine:        sequentialMax,
			RollingMaxBestCosine:           rollingMax,
			AnalyticMaxBestCosine:          analyticMax,
			RollingEfficientCandidate:      rollingCandidate,
		},
	}, nil
}
