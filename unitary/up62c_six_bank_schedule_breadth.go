package unitary

const UP62CSixBankScheduleBreadthSchema = "wingless.up62c-six-bank-schedule-breadth.v1"

type UP62CSixBankScheduleBreadthResult struct {
	Schema             string             `json:"schema"`
	Experiment         string             `json:"experiment"`
	SourceUP60CSeal    string             `json:"source_up60c_seal"`
	Banks              int                `json:"banks"`
	StateDimension     int                `json:"state_dimension"`
	ScheduleBases      []int              `json:"schedule_bases"`
	NoiseLevels        []float64          `json:"noise_levels"`
	Families           []string           `json:"families"`
	SelectionPerformed bool               `json:"selection_performed"`
	Metrics            []UP59CNoiseMetric `json:"metrics"`
}

func RunUP62C() (UP62CSixBankScheduleBreadthResult, error) {
	schedules := []int{79000000, 80000000, 81000000, 82000000, 83000000, 84000000}
	noises := []float64{0.0075, 0.0100}
	all := up58cFamilies()
	var families []UP58CTagFamily
	for _, f := range all {
		if f.Name == "golden_rotation" || f.Name == "fixed_irregular" {
			families = append(families, f)
		}
	}

	result := UP62CSixBankScheduleBreadthResult{
		Schema: UP62CSixBankScheduleBreadthSchema,
		Experiment: "UP-62C-six-bank-schedule-breadth",
		SourceUP60CSeal: "a157a3580438a28627ba885c13435f23bb53bd8a",
		Banks: 6,
		StateDimension: 16,
		ScheduleBases: append([]int(nil), schedules...),
		NoiseLevels: append([]float64(nil), noises...),
		SelectionPerformed: false,
	}
	for _, f := range families {
		result.Families = append(result.Families, f.Name)
	}
	for _, seedBase := range schedules {
		for _, f := range families {
			for _, noise := range noises {
				m, err := up59cEvaluate(f.Tags, f.Name, noise, seedBase)
				if err != nil {
					return UP62CSixBankScheduleBreadthResult{}, err
				}
				result.Metrics = append(result.Metrics, m)
			}
		}
	}
	return result, nil
}
