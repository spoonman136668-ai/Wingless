package unitary

const UP75ACompactFactorizedStepFloorSchema = "wingless.up75a-compact-factorized-step-floor.v1"

type UP75ACompactFactorizedStepFloorResult struct {
	Schema          string                     `json:"schema"`
	Experiment      string                     `json:"experiment"`
	SourceUP74ASeal string                     `json:"source_up74a_seal"`
	TrainingSteps   []int                      `json:"training_steps"`
	TrainingStates  int                        `json:"training_states"`
	HeldOutStates   int                        `json:"heldout_states"`
	Points          []UP73ATrainingBudgetPoint `json:"points"`
}

func RunUP75A() (UP75ACompactFactorizedStepFloorResult, error) {
	steps := []int{1, 2, 4, 8, 16}
	result := UP75ACompactFactorizedStepFloorResult{
		Schema: UP75ACompactFactorizedStepFloorSchema,
		Experiment: "UP-75A-compact-factorized-step-floor",
		SourceUP74ASeal: "04cd3f9b2e77f706e624790902fcd8bc022d104b",
		TrainingSteps: append([]int(nil), steps...),
		TrainingStates: 54,
		HeldOutStates: 162,
	}
	for _, n := range steps {
		onehot, err := up73aRunPoint(n, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil { return UP75ACompactFactorizedStepFloorResult{}, err }
		result.Points = append(result.Points, onehot)
		simplex, err := up73aRunPoint(n, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil { return UP75ACompactFactorizedStepFloorResult{}, err }
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
