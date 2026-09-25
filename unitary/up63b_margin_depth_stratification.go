package unitary

const UP63BMarginDepthStratificationSchema = "wingless.up63b-margin-depth-stratification.v1"

type UP63BDepthCalibrationMetric struct {
	ScheduleBase      int                 `json:"schedule_base"`
	MemoryNoise       float64             `json:"memory_noise"`
	WritesPerScenario int                 `json:"writes_per_scenario"`
	Depth             int                 `json:"depth"`
	Buckets           []UP61BMarginBucket `json:"buckets"`
	MaxNormDrift      float64             `json:"max_norm_drift"`
}

type UP63BMarginDepthStratificationResult struct {
	Schema            string                        `json:"schema"`
	Experiment        string                        `json:"experiment"`
	SourceUP62BSeal   string                        `json:"source_up62b_seal"`
	ScheduleBases     []int                         `json:"schedule_bases"`
	NoiseLevels       []float64                     `json:"noise_levels"`
	DepthLevels       []int                         `json:"depth_levels"`
	WritesPerScenario int                           `json:"writes_per_scenario"`
	BinEdges          []float64                     `json:"bin_edges"`
	OracleUsed        bool                          `json:"oracle_used"`
	PolicyChanged     bool                          `json:"policy_changed"`
	TrainingChanged   bool                          `json:"training_changed"`
	ThresholdSelected bool                          `json:"threshold_selected"`
	Metrics           []UP63BDepthCalibrationMetric `json:"metrics"`
}

func RunUP63B() (UP63BMarginDepthStratificationResult, error) {
	schedules := []int{94000000, 95000000}
	noises := []float64{0.085, 0.09}
	depths := []int{32, 128, 512, 1024}
	const writes = 256
	trainDepths := []int{0}
	allDepths := []int{0, 32, 128, 512, 1024}
	mixer := fullLatentMixer()
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	offsets, err := up50bFusedOffsets()
	if err != nil { return UP63BMarginDepthStratificationResult{}, err }
	result := UP63BMarginDepthStratificationResult{
		Schema: UP63BMarginDepthStratificationSchema, Experiment: "UP-63B-margin-depth-stratification",
		SourceUP62BSeal: "9fe4b71e267e34af346fe771da0025c2542d5555",
		ScheduleBases: append([]int(nil), schedules...), NoiseLevels: append([]float64(nil), noises...),
		DepthLevels: append([]int(nil), depths...), WritesPerScenario: writes,
		BinEdges: []float64{0,0.25,0.5,0.75,1}, OracleUsed: false, PolicyChanged: false,
		TrainingChanged: false, ThresholdSelected: false,
	}
	for _, noise := range noises {
		prepared, err := up52bPrepareArm("selected", offsets, mixer, trainTables, heldTables, trainDepths, depths, allDepths, noise)
		if err != nil { return UP63BMarginDepthStratificationResult{}, err }
		for _, depth := range depths {
			fixedDepth := []int{depth}
			for _, seedBase := range schedules {
				m, err := up61bRunCalibration(mixer, prepared, heldTables, fixedDepth, noise, writes, seedBase)
				if err != nil { return UP63BMarginDepthStratificationResult{}, err }
				result.Metrics = append(result.Metrics, UP63BDepthCalibrationMetric{
					ScheduleBase: m.ScheduleBase, MemoryNoise: m.MemoryNoise, WritesPerScenario: m.WritesPerScenario,
					Depth: depth, Buckets: m.Buckets, MaxNormDrift: m.MaxNormDrift,
				})
			}
		}
	}
	return result, nil
}
