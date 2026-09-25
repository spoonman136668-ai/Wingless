package unitary

const UP62BMarginLongHorizonSchema = "wingless.up62b-margin-long-horizon.v1"

type UP62BMarginLongHorizonResult struct {
	Schema            string                   `json:"schema"`
	Experiment        string                   `json:"experiment"`
	SourceUP61BSeal   string                   `json:"source_up61b_seal"`
	ScheduleBases     []int                    `json:"schedule_bases"`
	NoiseLevels       []float64                `json:"noise_levels"`
	WriteHorizons     []int                    `json:"write_horizons"`
	BinEdges          []float64                `json:"bin_edges"`
	OracleUsed        bool                     `json:"oracle_used"`
	PolicyChanged     bool                     `json:"policy_changed"`
	TrainingChanged   bool                     `json:"training_changed"`
	ThresholdSelected bool                     `json:"threshold_selected"`
	Metrics           []UP61BCalibrationMetric `json:"metrics"`
}

func RunUP62B() (UP62BMarginLongHorizonResult, error) {
	schedules := []int{92000000, 93000000}
	noises := []float64{0.08, 0.085, 0.09}
	writes := []int{128, 256}
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}
	mixer := fullLatentMixer()
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	offsets, err := up50bFusedOffsets()
	if err != nil { return UP62BMarginLongHorizonResult{}, err }
	result := UP62BMarginLongHorizonResult{
		Schema: UP62BMarginLongHorizonSchema, Experiment: "UP-62B-margin-long-horizon",
		SourceUP61BSeal: "3ec4c5975388b109892a572e850c8d8729b2d1c8",
		ScheduleBases: append([]int(nil), schedules...), NoiseLevels: append([]float64(nil), noises...),
		WriteHorizons: append([]int(nil), writes...), BinEdges: []float64{0,0.25,0.5,0.75,1},
		OracleUsed: false, PolicyChanged: false, TrainingChanged: false, ThresholdSelected: false,
	}
	for _, noise := range noises {
		prepared, err := up52bPrepareArm("selected", offsets, mixer, trainTables, heldTables, trainDepths, heldDepths, allDepths, noise)
		if err != nil { return UP62BMarginLongHorizonResult{}, err }
		for _, writeCount := range writes {
			for _, seedBase := range schedules {
				m, err := up61bRunCalibration(mixer, prepared, heldTables, heldDepths, noise, writeCount, seedBase)
				if err != nil { return UP62BMarginLongHorizonResult{}, err }
				result.Metrics = append(result.Metrics, m)
			}
		}
	}
	return result, nil
}
