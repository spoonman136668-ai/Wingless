package unitary

const UP74ACompactFactorizedLowBudgetSchema = "wingless.up74a-compact-factorized-low-budget.v1"

type UP74ACompactFactorizedLowBudgetResult struct {
	Schema          string                     `json:"schema"`
	Experiment      string                     `json:"experiment"`
	SourceUP73ASeal string                     `json:"source_up73a_seal"`
	TrainingSteps   []int                      `json:"training_steps"`
	TrainingStates  int                        `json:"training_states"`
	HeldOutStates   int                        `json:"heldout_states"`
	Points          []UP73ATrainingBudgetPoint `json:"points"`
}

func RunUP74A() (UP74ACompactFactorizedLowBudgetResult, error) {
	steps := []int{25, 50, 100, 200}
	result := UP74ACompactFactorizedLowBudgetResult{
		Schema: UP74ACompactFactorizedLowBudgetSchema,
		Experiment: "UP-74A-compact-factorized-low-budget",
		SourceUP73ASeal: "6f9ea0d3f771e11831222e0742b7b799ca1e6892",
		TrainingSteps: append([]int(nil), steps...),
		TrainingStates: 54,
		HeldOutStates: 162,
	}
	for _, n := range steps {
		onehot, err := up73aRunPoint(n, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil {
			return UP74ACompactFactorizedLowBudgetResult{}, err
		}
		result.Points = append(result.Points, onehot)
		simplex, err := up73aRunPoint(n, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil {
			return UP74ACompactFactorizedLowBudgetResult{}, err
		}
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
