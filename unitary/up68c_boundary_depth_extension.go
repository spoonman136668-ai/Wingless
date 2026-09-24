package unitary

const UP68CBoundaryDepthExtensionSchema = "wingless.up68c-boundary-depth-extension.v1"

type UP68CBoundaryDepthExtensionResult struct {
	Schema               string             `json:"schema"`
	Experiment           string             `json:"experiment"`
	SourceUP67CSeal      string             `json:"source_up67c_seal"`
	Banks                int                `json:"banks"`
	StateDimension       int                `json:"state_dimension"`
	ScheduleBases        []int              `json:"schedule_bases"`
	GoldenNoiseLevels    []float64          `json:"golden_noise_levels"`
	IrregularNoiseLevels []float64          `json:"irregular_noise_levels"`
	DepthLevels          []int              `json:"depth_levels"`
	SelectionPerformed   bool               `json:"selection_performed"`
	Metrics              []UP67CDepthMetric `json:"metrics"`
}

func RunUP68C() (UP68CBoundaryDepthExtensionResult, error) {
	schedules := []int{107000000, 108000000}
	goldenNoises := []float64{0.0106, 0.0108, 0.0110}
	irregularNoises := []float64{0.0085, 0.0086, 0.0087}
	depths := []int{256, 512, 1024}

	var golden, irregular UP58CTagFamily
	for _, f := range up58cFamilies() {
		switch f.Name {
		case "golden_rotation":
			golden = f
		case "fixed_irregular":
			irregular = f
		}
	}

	result := UP68CBoundaryDepthExtensionResult{
		Schema: UP68CBoundaryDepthExtensionSchema,
		Experiment: "UP-68C-boundary-depth-extension",
		SourceUP67CSeal: "7aa5f2cf56a5d61aef33807ee86ae25f0644c7f6",
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
				m, err := up67cEvaluate(golden.Tags, golden.Name, noise, seedBase, depth)
				if err != nil {
					return UP68CBoundaryDepthExtensionResult{}, err
				}
				result.Metrics = append(result.Metrics, m)
			}
			for _, noise := range irregularNoises {
				m, err := up67cEvaluate(irregular.Tags, irregular.Name, noise, seedBase, depth)
				if err != nil {
					return UP68CBoundaryDepthExtensionResult{}, err
				}
				result.Metrics = append(result.Metrics, m)
			}
		}
	}
	return result, nil
}
