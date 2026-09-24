package unitary

import "math"

const UP61BMarginCalibrationSchema = "wingless.up61b-margin-calibration.v1"

type UP61BMarginBucket struct {
	LowerBound float64 `json:"lower_bound"`
	UpperBound float64 `json:"upper_bound"`
	Events     int     `json:"events"`
	Correct    int     `json:"correct"`
	Accuracy   float64 `json:"accuracy"`
}

type UP61BCalibrationMetric struct {
	ScheduleBase      int                 `json:"schedule_base"`
	MemoryNoise       float64             `json:"memory_noise"`
	WritesPerScenario int                 `json:"writes_per_scenario"`
	Buckets           []UP61BMarginBucket `json:"buckets"`
	MaxNormDrift      float64             `json:"max_norm_drift"`
}

type UP61BMarginCalibrationResult struct {
	Schema            string                   `json:"schema"`
	Experiment        string                   `json:"experiment"`
	SourceUP58BSeal   string                   `json:"source_up58b_seal"`
	ScheduleBases     []int                    `json:"schedule_bases"`
	NoiseLevels       []float64                `json:"noise_levels"`
	WritesPerScenario int                      `json:"writes_per_scenario"`
	BinEdges          []float64                `json:"bin_edges"`
	OracleUsed        bool                     `json:"oracle_used"`
	PolicyChanged     bool                     `json:"policy_changed"`
	TrainingChanged   bool                     `json:"training_changed"`
	ThresholdSelected bool                     `json:"threshold_selected"`
	Metrics           []UP61BCalibrationMetric `json:"metrics"`
}

func up61bBucketIndex(margin float64) int {
	switch {
	case margin < 0.25:
		return 0
	case margin < 0.50:
		return 1
	case margin < 0.75:
		return 2
	default:
		return 3
	}
}

func up61bRunCalibration(
	mixer latentMatrix,
	prepared up52bPreparedArm,
	heldTables []memoryTable,
	depths []int,
	noise float64,
	writes int,
	seedBase int,
) (UP61BCalibrationMetric, error) {
	const scenarios = 48
	buckets := []UP61BMarginBucket{
		{LowerBound: 0.00, UpperBound: 0.25},
		{LowerBound: 0.25, UpperBound: 0.50},
		{LowerBound: 0.50, UpperBound: 0.75},
		{LowerBound: 0.75, UpperBound: 1.00},
	}
	var maxNormDrift float64
	for scenarioIndex := 0; scenarioIndex < scenarios; scenarioIndex++ {
		scenario := makeFullRankScenario(scenarioIndex, heldTables[(scenarioIndex*7)%len(heldTables)], writes, depths)
		pathTable := scenario.initial
		trueTable := scenario.initial
		for writeIndex, write := range scenario.writes {
			canonical, err := encodeMemory(pathTable)
			if err != nil {
				return UP61BCalibrationMetric{}, err
			}
			seed := seedBase + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, noise)
			if err != nil {
				return UP61BCalibrationMetric{}, err
			}
			memory = rotateGlobalPhase(memory, math.Mod(0.271*float64(seed+1), 2*math.Pi))
			state, err := fullLatentEncode(memory, mixer)
			if err != nil {
				return UP61BCalibrationMetric{}, err
			}
			operator, ok := prepared.Ops[write.gap]
			if !ok {
				return UP61BCalibrationMetric{}, &up61bMissingDepthError{Depth: write.gap}
			}
			state, err = latentMatrixVector(operator, state)
			if err != nil {
				return UP61BCalibrationMetric{}, err
			}
			norm2, err := NormSquared(state)
			if err != nil {
				return UP61BCalibrationMetric{}, err
			}
			if d := math.Abs(norm2 - 1); d > maxNormDrift {
				maxNormDrift = d
			}
			decoded, _, margin, err := decodeBreadthTable(state, prepared.Observables, prepared.Regressors, prepared.Classifiers)
			if err != nil {
				return UP61BCalibrationMetric{}, err
			}
			bi := up61bBucketIndex(margin)
			buckets[bi].Events++
			if decoded == trueTable {
				buckets[bi].Correct++
			}
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil {
				return UP61BCalibrationMetric{}, err
			}
			pathTable, err = applyMemoryWrite(decoded, write.entity, write.value)
			if err != nil {
				return UP61BCalibrationMetric{}, err
			}
		}
	}
	for i := range buckets {
		if buckets[i].Events > 0 {
			buckets[i].Accuracy = float64(buckets[i].Correct) / float64(buckets[i].Events)
		}
	}
	return UP61BCalibrationMetric{
		ScheduleBase: seedBase,
		MemoryNoise: noise,
		WritesPerScenario: writes,
		Buckets: buckets,
		MaxNormDrift: maxNormDrift,
	}, nil
}

type up61bMissingDepthError struct {
	Depth int
}

func (e *up61bMissingDepthError) Error() string {
	return "UP61B missing prepared depth"
}

func RunUP61B() (UP61BMarginCalibrationResult, error) {
	schedules := []int{86000000, 87000000, 88000000}
	noises := []float64{0.06, 0.07, 0.08, 0.09}
	const writes = 64
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}
	mixer := fullLatentMixer()
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	offsets, err := up50bFusedOffsets()
	if err != nil {
		return UP61BMarginCalibrationResult{}, err
	}
	result := UP61BMarginCalibrationResult{
		Schema: UP61BMarginCalibrationSchema,
		Experiment: "UP-61B-margin-calibration",
		SourceUP58BSeal: "850e010027508a9c463656b9f4ce5e73542048dd",
		ScheduleBases: append([]int(nil), schedules...),
		NoiseLevels: append([]float64(nil), noises...),
		WritesPerScenario: writes,
		BinEdges: []float64{0, 0.25, 0.5, 0.75, 1},
		OracleUsed: false,
		PolicyChanged: false,
		TrainingChanged: false,
		ThresholdSelected: false,
	}
	for _, noise := range noises {
		prepared, err := up52bPrepareArm("selected", offsets, mixer, trainTables, heldTables, trainDepths, heldDepths, allDepths, noise)
		if err != nil {
			return UP61BMarginCalibrationResult{}, err
		}
		for _, seedBase := range schedules {
			m, err := up61bRunCalibration(mixer, prepared, heldTables, heldDepths, noise, writes, seedBase)
			if err != nil {
				return UP61BMarginCalibrationResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
	}
	return result, nil
}
