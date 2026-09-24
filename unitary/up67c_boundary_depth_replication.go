package unitary

const UP67CBoundaryDepthReplicationSchema = "wingless.up67c-boundary-depth-replication.v1"

type UP67CBoundaryDepthReplicationResult struct {
	Schema               string             `json:"schema"`
	Experiment           string             `json:"experiment"`
	SourceUP65CSeal      string             `json:"source_up65c_seal"`
	Banks                int                `json:"banks"`
	StateDimension       int                `json:"state_dimension"`
	ScheduleBases        []int              `json:"schedule_bases"`
	GoldenNoiseLevels    []float64          `json:"golden_noise_levels"`
	IrregularNoiseLevels []float64          `json:"irregular_noise_levels"`
	DepthLevels          []int              `json:"depth_levels"`
	SelectionPerformed   bool               `json:"selection_performed"`
	Metrics              []UP64CDepthMetric `json:"metrics"`
}

func RunUP67C() (UP67CBoundaryDepthReplicationResult, error) {
	schedules := []int{105000000, 106000000}
	goldenNoises := []float64{0.0106, 0.0108, 0.0110}
	irregularNoises := []float64{0.0085, 0.0086, 0.0087}
	depths := []int{32, 64, 128}

	var golden, irregular UP58CTagFamily
	for _, f := range up58cFamilies() {
		switch f.Name {
		case "golden_rotation":
			golden = f
		case "fixed_irregular":
			irregular = f
		}
	}

	result := UP67CBoundaryDepthReplicationResult{
		Schema: UP67CBoundaryDepthReplicationSchema,
		Experiment: "UP-67C-boundary-depth-replication",
		SourceUP65CSeal: "1895beda79f067ee1bf6646138445eed828fd3d8",
		Banks: 6,
		StateDimension: 16,
		ScheduleBases: append([]int(nil), schedules...),
		GoldenNoiseLevels: append([]float64(nil), goldenNoises...),
		IrregularNoiseLevels: append([]float64(nil), irregularNoises...),
		DepthLevels: append([]int(nil), depths...),
		SelectionPerformed: false,
	}

	for _, seedBase := range schedules {
		for _, depth := range depths {
			for _, noise := range goldenNoises {
				m, err := up64cEvaluate(golden.Tags, golden.Name, noise, seedBase, depth)
				if err != nil {
					return UP67CBoundaryDepthReplicationResult{}, err
				}
				result.Metrics = append(result.Metrics, m)
			}
			for _, noise := range irregularNoises {
				m, err := up64cEvaluate(irregular.Tags, irregular.Name, noise, seedBase, depth)
				if err != nil {
					return UP67CBoundaryDepthReplicationResult{}, err
				}
				result.Metrics = append(result.Metrics, m)
			}
		}
	}
	return result, nil
}
