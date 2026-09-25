package unitary

const UP70CGoldenBankCapacitySchema = "wingless.up70c-golden-bank-capacity.v1"

type UP70CGoldenBankCapacityResult struct {
	Schema             string                  `json:"schema"`
	Experiment         string                  `json:"experiment"`
	SourceUP69CSeal    string                  `json:"source_up69c_seal"`
	StateDimension     int                     `json:"state_dimension"`
	ScheduleBases      []int                   `json:"schedule_bases"`
	BankCounts         []int                   `json:"bank_counts"`
	NoiseLevels        []float64               `json:"noise_levels"`
	SelectionPerformed bool                    `json:"selection_performed"`
	Points             []UP57CGoldenScalePoint `json:"points"`
}

func RunUP70C() (UP70CGoldenBankCapacityResult, error) {
	schedules := []int{111000000, 112000000}
	banks := []int{6, 7, 8, 10}
	noises := []float64{0, 0.004, 0.008}
	result := UP70CGoldenBankCapacityResult{
		Schema: UP70CGoldenBankCapacitySchema, Experiment: "UP-70C-golden-bank-capacity",
		SourceUP69CSeal: "7c59d784b7a9cca434cd0e9787177ffa9a1ebb44",
		StateDimension: 16, ScheduleBases: append([]int(nil), schedules...),
		BankCounts: append([]int(nil), banks...), NoiseLevels: append([]float64(nil), noises...), SelectionPerformed: false,
	}
	for _, seedBase := range schedules {
		for _, bankCount := range banks {
			for _, noise := range noises {
				p, err := up57cEvaluate(bankCount, noise, seedBase)
				if err != nil { return UP70CGoldenBankCapacityResult{}, err }
				result.Points = append(result.Points, p)
			}
		}
	}
	return result, nil
}
