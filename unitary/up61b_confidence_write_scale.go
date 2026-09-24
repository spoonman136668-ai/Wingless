package unitary

const UP61BConfidenceWriteScaleSchema = "wingless.up61b-confidence-write-scale.v1"

type UP61BConfidenceWriteScaleResult struct {
	Schema            string              `json:"schema"`
	Experiment        string              `json:"experiment"`
	SourceUP60BSeal   string              `json:"source_up60b_seal"`
	ScheduleBases     []int               `json:"schedule_bases"`
	NoiseLevels       []float64           `json:"noise_levels"`
	WriteLevels       []int               `json:"write_levels"`
	Thresholds        []float64           `json:"thresholds"`
	OracleUsed        bool                `json:"oracle_used"`
	TrainingChanged   bool                `json:"training_changed"`
	ThresholdSelected bool                `json:"threshold_selected"`
	Metrics           []UP57BPolicyMetric `json:"metrics"`
}

func RunUP61BConfidenceWriteScale() (UP61BConfidenceWriteScaleResult, error) {
	schedules := []int{89000000, 90000000, 91000000}
	noises := []float64{0.08, 0.085, 0.09}
	writeLevels := []int{128, 256, 512}
	thresholds := []float64{0.25, 0.5, 0.75}
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	offsets, err := up50bFusedOffsets()
	if err != nil {
		return UP61BConfidenceWriteScaleResult{}, err
	}

	result := UP61BConfidenceWriteScaleResult{
		Schema: UP61BConfidenceWriteScaleSchema,
		Experiment: "UP-61B-confidence-write-scale",
		SourceUP60BSeal: "f639951748145749a75403a6cb0e59bd375aea15",
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
			return UP61BConfidenceWriteScaleResult{}, err
		}
		for _, writes := range writeLevels {
			for _, seedBase := range schedules {
				base, err := runUP57BPolicy(mixer, prepared, heldTables, heldDepths, noise, writes, seedBase, 0, false)
				if err != nil {
					return UP61BConfidenceWriteScaleResult{}, err
				}
				result.Metrics = append(result.Metrics, base)
				for _, threshold := range thresholds {
					m, err := runUP57BPolicy(mixer, prepared, heldTables, heldDepths, noise, writes, seedBase, threshold, true)
					if err != nil {
						return UP61BConfidenceWriteScaleResult{}, err
					}
					result.Metrics = append(result.Metrics, m)
				}
			}
		}
	}
	return result, nil
}
