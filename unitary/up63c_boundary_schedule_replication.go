package unitary

const UP63CBoundaryScheduleReplicationSchema = "wingless.up63c-boundary-schedule-replication.v1"

type UP63CBoundaryScheduleReplicationResult struct {
	Schema                string             `json:"schema"`
	Experiment            string             `json:"experiment"`
	SourceUP61CSeal       string             `json:"source_up61c_seal"`
	Banks                 int                `json:"banks"`
	StateDimension        int                `json:"state_dimension"`
	ScheduleBases         []int              `json:"schedule_bases"`
	GoldenNoiseLevels     []float64          `json:"golden_noise_levels"`
	IrregularNoiseLevels  []float64          `json:"irregular_noise_levels"`
	SelectionPerformed    bool               `json:"selection_performed"`
	Metrics               []UP59CNoiseMetric `json:"metrics"`
}

func RunUP63C() (UP63CBoundaryScheduleReplicationResult, error) {
	schedules := []int{85000000, 86000000, 87000000, 88000000, 89000000, 90000000}
	goldenNoises := []float64{0.010, 0.011}
	irregularNoises := []float64{0.0085, 0.0090}

	var golden, irregular UP58CTagFamily
	for _, f := range up58cFamilies() {
		switch f.Name {
		case "golden_rotation":
			golden = f
		case "fixed_irregular":
			irregular = f
		}
	}

	result := UP63CBoundaryScheduleReplicationResult{
		Schema: UP63CBoundaryScheduleReplicationSchema,
		Experiment: "UP-63C-boundary-schedule-replication",
		SourceUP61CSeal: "854def68553ebc9cd6004acd83f78a691a8407b2",
		Banks: 6,
		StateDimension: 16,
		ScheduleBases: append([]int(nil), schedules...),
		GoldenNoiseLevels: append([]float64(nil), goldenNoises...),
		IrregularNoiseLevels: append([]float64(nil), irregularNoises...),
		SelectionPerformed: false,
	}

	for _, seedBase := range schedules {
		for _, noise := range goldenNoises {
			m, err := up59cEvaluate(golden.Tags, golden.Name, noise, seedBase)
			if err != nil {
				return UP63CBoundaryScheduleReplicationResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
		for _, noise := range irregularNoises {
			m, err := up59cEvaluate(irregular.Tags, irregular.Name, noise, seedBase)
			if err != nil {
				return UP63CBoundaryScheduleReplicationResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
	}
	return result, nil
}
