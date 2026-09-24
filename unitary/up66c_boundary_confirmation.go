package unitary

const UP66CBoundaryConfirmationSchema = "wingless.up66c-boundary-confirmation.v1"

type UP66CBoundaryConfirmationResult struct {
	Schema               string             `json:"schema"`
	Experiment           string             `json:"experiment"`
	SourceUP65CSeal      string             `json:"source_up65c_seal"`
	Banks                int                `json:"banks"`
	StateDimension       int                `json:"state_dimension"`
	ScheduleBases        []int              `json:"schedule_bases"`
	GoldenNoiseLevels    []float64          `json:"golden_noise_levels"`
	IrregularNoiseLevels []float64          `json:"irregular_noise_levels"`
	SelectionPerformed   bool               `json:"selection_performed"`
	Metrics              []UP59CNoiseMetric `json:"metrics"`
}

func RunUP66C() (UP66CBoundaryConfirmationResult, error) {
	schedules := []int{99000000, 100000000, 101000000, 102000000, 103000000, 104000000}
	goldenNoises := []float64{0.0106, 0.0108, 0.0110}
	irregularNoises := []float64{0.0085, 0.0086, 0.0087}

	var golden, irregular UP58CTagFamily
	for _, f := range up58cFamilies() {
		switch f.Name {
		case "golden_rotation":
			golden = f
		case "fixed_irregular":
			irregular = f
		}
	}

	result := UP66CBoundaryConfirmationResult{
		Schema: UP66CBoundaryConfirmationSchema,
		Experiment: "UP-66C-boundary-confirmation",
		SourceUP65CSeal: "1895beda79f067ee1bf6646138445eed828fd3d8",
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
				return UP66CBoundaryConfirmationResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
		for _, noise := range irregularNoises {
			m, err := up59cEvaluate(irregular.Tags, irregular.Name, noise, seedBase)
			if err != nil {
				return UP66CBoundaryConfirmationResult{}, err
			}
			result.Metrics = append(result.Metrics, m)
		}
	}
	return result, nil
}
