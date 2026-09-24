package unitary

const UP65CBoundaryFineGridSchema = "wingless.up65c-boundary-fine-grid.v1"

type UP65CBoundaryFineGridResult struct {
	Schema                string             `json:"schema"`
	Experiment            string             `json:"experiment"`
	SourceUP63CSeal       string             `json:"source_up63c_seal"`
	Banks                 int                `json:"banks"`
	StateDimension        int                `json:"state_dimension"`
	ScheduleBases         []int              `json:"schedule_bases"`
	GoldenNoiseLevels     []float64          `json:"golden_noise_levels"`
	IrregularNoiseLevels  []float64          `json:"irregular_noise_levels"`
	SelectionPerformed    bool               `json:"selection_performed"`
	Metrics               []UP59CNoiseMetric `json:"metrics"`
}

func RunUP65C() (UP65CBoundaryFineGridResult, error) {
	schedules := []int{93000000, 94000000, 95000000, 96000000, 97000000, 98000000}
	goldenNoises := []float64{0.0100, 0.0102, 0.0104, 0.0106, 0.0108, 0.0110}
	irregularNoises := []float64{0.0085, 0.0086, 0.0087, 0.0088, 0.0089, 0.0090}

	var golden, irregular UP58CTagFamily
	for _, f := range up58cFamilies() {
		switch f.Name {
		case "golden_rotation":
			golden = f
		case "fixed_irregular":
			irregular = f
		}
	}

	result := UP65CBoundaryFineGridResult{
		Schema: UP65CBoundaryFineGridSchema,
		Experiment: "UP-65C-boundary-fine-grid",
		SourceUP63CSeal: "66e705ac348294947625093728d4c6b792874dbb",
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
				return UP65CBoundaryFineGridResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
		for _, noise := range irregularNoises {
			m, err := up59cEvaluate(irregular.Tags, irregular.Name, noise, seedBase)
			if err != nil {
				return UP65CBoundaryFineGridResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
	}
	return result, nil
}
