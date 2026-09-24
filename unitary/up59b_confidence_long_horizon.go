package unitary

const UP59BConfidenceLongHorizonSchema = "wingless.up59b-confidence-long-horizon.v1"

type UP59BConfidenceLongHorizonResult struct {
	Schema            string              `json:"schema"`
	Experiment        string              `json:"experiment"`
	SourceUP58BSeal   string              `json:"source_up58b_seal"`
	ScheduleBases     []int               `json:"schedule_bases"`
	NoiseLevels       []float64           `json:"noise_levels"`
	WriteLevels       []int               `json:"write_levels"`
	Thresholds        []float64           `json:"thresholds"`
	OracleUsed        bool                `json:"oracle_used"`
	TrainingChanged   bool                `json:"training_changed"`
	ThresholdSelected bool                `json:"threshold_selected"`
	Metrics           []UP57BPolicyMetric `json:"metrics"`
}

func RunUP59B() (UP59BConfidenceLongHorizonResult, error) {
	schedules := []int{80000000, 81000000, 82000000}
	noises := []float64{0.07, 0.075}
	writeLevels := []int{256, 512}
	thresholds := []float64{0.25, 0.5, 0.75}
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	offsets, err := up50bFusedOffsets()
	if err != nil {
		return UP59BConfidenceLongHorizonResult{}, err
	}

	result := UP59BConfidenceLongHorizonResult{
		Schema: UP59BConfidenceLongHorizonSchema,
		Experiment: "UP-59B-confidence-long-horizon",
		SourceUP58BSeal: "850e010027508a9c463656b9f4ce5e73542048dd",
		ScheduleBases: append([]int(nil), schedules...),
		NoiseLevels: append([]float64(nil), noises...),
		WriteLevels: append([]int(nil), writeLevels...),
		Thresholds: append([]float64(nil), thresholds...),
		OracleUsed: false,
		TrainingChanged: false,
		ThresholdSelected: false,
	}

	for _, noise := range noises {
		prepared, err := up52bPrepareArm("selected", offsets, mixer, trainTables, heldTables, trainDepths, heldDepths, allDepths, noise)
		if err != nil {
			return UP59BConfidenceLongHorizonResult{}, err
		}
		for _, writes := range writeLevels {
			for _, seedBase := range schedules {
				base, err := runUP57BPolicy(mixer, prepared, heldTables, heldDepths, noise, writes, seedBase, 0, false)
				if err != nil {
					return UP59BConfidenceLongHorizonResult{}, err
				}
				result.Metrics = append(result.Metrics, base)
				for _, threshold := range thresholds {
					m, err := runUP57BPolicy(mixer, prepared, heldTables, heldDepths, noise, writes, seedBase, threshold, true)
					if err != nil {
						return UP59BConfidenceLongHorizonResult{}, err
					}
					result.Metrics = append(result.Metrics, m)
				}
			}
		}
	}
	return result, nil
}
