package unitary

const UP61CSixBankBoundaryRefinementSchema = "wingless.up61c-six-bank-boundary-refinement.v1"

type UP61CSixBankBoundaryRefinementResult struct {
	Schema             string             `json:"schema"`
	Experiment         string             `json:"experiment"`
	SourceUP60CSeal    string             `json:"source_up60c_seal"`
	Banks              int                `json:"banks"`
	StateDimension     int                `json:"state_dimension"`
	ScheduleBases      []int              `json:"schedule_bases"`
	GoldenNoiseLevels  []float64          `json:"golden_noise_levels"`
	IrregularNoiseLevels []float64        `json:"irregular_noise_levels"`
	SelectionPerformed bool               `json:"selection_performed"`
	Metrics            []UP59CNoiseMetric `json:"metrics"`
}

func RunUP61C() (UP61CSixBankBoundaryRefinementResult, error) {
	schedules := []int{77000000, 78000000}
	goldenNoises := []float64{0.010, 0.011, 0.012, 0.013, 0.014, 0.015}
	irregularNoises := []float64{0.0075, 0.0080, 0.0085, 0.0090, 0.0095, 0.0100}

	var golden, irregular UP58CTagFamily
	for _, f := range up58cFamilies() {
		switch f.Name {
		case "golden_rotation":
			golden = f
		case "fixed_irregular":
			irregular = f
		}
	}

	result := UP61CSixBankBoundaryRefinementResult{
		Schema: UP61CSixBankBoundaryRefinementSchema,
		Experiment: "UP-61C-six-bank-boundary-refinement",
		SourceUP60CSeal: "a157a3580438a28627ba885c13435f23bb53bd8a",
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
				return UP61CSixBankBoundaryRefinementResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
		for _, noise := range irregularNoises {
			m, err := up59cEvaluate(irregular.Tags, irregular.Name, noise, seedBase)
			if err != nil {
				return UP61CSixBankBoundaryRefinementResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
	}
	return result, nil
}
