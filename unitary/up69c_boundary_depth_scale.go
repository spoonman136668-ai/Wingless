package unitary

const UP69CBoundaryDepthScaleSchema = "wingless.up69c-boundary-depth-scale.v1"

type UP69CBoundaryDepthScaleResult struct {
	Schema               string             `json:"schema"`
	Experiment           string             `json:"experiment"`
	SourceUP68CSeal      string             `json:"source_up68c_seal"`
	Banks                int                `json:"banks"`
	StateDimension       int                `json:"state_dimension"`
	ScheduleBases        []int              `json:"schedule_bases"`
	GoldenNoiseLevels    []float64          `json:"golden_noise_levels"`
	IrregularNoiseLevels []float64          `json:"irregular_noise_levels"`
	DepthLevels          []int              `json:"depth_levels"`
	SelectionPerformed   bool               `json:"selection_performed"`
	Metrics              []UP67CDepthMetric `json:"metrics"`
}

func RunUP69C() (UP69CBoundaryDepthScaleResult, error) {
	schedules := []int{109000000, 110000000}
	goldenNoises := []float64{0.0108, 0.0110}
	irregularNoises := []float64{0.0087}
	depths := []int{1024, 2048, 4096}
	var golden, irregular UP58CTagFamily
	for _, f := range up58cFamilies() {
		switch f.Name {
		case "golden_rotation": golden = f
		case "fixed_irregular": irregular = f
		}
	}
	result := UP69CBoundaryDepthScaleResult{
		Schema: UP69CBoundaryDepthScaleSchema,
		Experiment: "UP-69C-boundary-depth-scale",
		SourceUP68CSeal: "f336e130724d0a61dc279b51fdaad55c6ae3eb8a",
		Banks: 6, StateDimension: 16,
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
				if err != nil { return UP69CBoundaryDepthScaleResult{}, err }
				result.Metrics = append(result.Metrics, m)
			}
			for _, noise := range irregularNoises {
				m, err := up67cEvaluate(irregular.Tags, irregular.Name, noise, seedBase, depth)
				if err != nil { return UP69CBoundaryDepthScaleResult{}, err }
				result.Metrics = append(result.Metrics, m)
			}
		}
	}
	return result, nil
}
