package unitary

const UP73ACompactFactorizedTrainingBudgetSchema = "wingless.up73a-compact-factorized-training-budget.v1"

type UP73ATrainingBudgetPoint struct {
	TrainingSteps       int               `json:"training_steps"`
	TrainingStates      int               `json:"training_states"`
	Representation      string            `json:"representation"`
	FeatureDimension    int               `json:"feature_dimension"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP73ACompactFactorizedTrainingBudgetResult struct {
	Schema          string                     `json:"schema"`
	Experiment      string                     `json:"experiment"`
	SourceUP71ASeal string                     `json:"source_up71a_seal"`
	TrainingSteps   []int                      `json:"training_steps"`
	TrainingStates  int                        `json:"training_states"`
	HeldOutStates   int                        `json:"heldout_states"`
	Points          []UP73ATrainingBudgetPoint `json:"points"`
}

func up73aRunPoint(steps int, name string, dim int, featureFn func(up64aRoles) []float64) (UP73ATrainingBudgetPoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r {
			sum += v
		}
		if sum%3 == 0 && up71aSecondary(r) < 2 {
			trainRoles = append(trainRoles, r)
		}
		if sum%3 != 0 {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP73ATrainingBudgetPoint{
		TrainingSteps: steps,
		TrainingStates: len(trainRoles),
		Representation: name,
		FeatureDimension: dim,
		HeldOutStates: len(heldRoles),
		Gate: true,
	}
	for role := 0; role < 5; role++ {
		train := make([]headSample, 0, len(trainRoles))
		held := make([]headSample, 0, len(heldRoles))
		for _, r := range trainRoles {
			train = append(train, headSample{features: featureFn(r), target: r[role]})
		}
		for _, r := range heldRoles {
			held = append(held, headSample{features: featureFn(r), target: r[role]})
		}
		head, metric, err := trainLinearSoftmax(train, 3, dim, steps, 1.0)
		if err != nil {
			return UP73ATrainingBudgetPoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP73ATrainingBudgetPoint{}, err
		}
		point.PerRole = append(point.PerRole, UP65ARoleMetric{
			Role: role,
			TrainAccuracy: metric.TrainAccuracy,
			HeldOutAccuracy: heldAcc,
		})
		point.MeanTrainAccuracy += metric.TrainAccuracy
		point.MeanHeldOutAccuracy += heldAcc
		if heldAcc < 0.98 {
			point.Gate = false
		}
	}
	point.MeanTrainAccuracy /= 5
	point.MeanHeldOutAccuracy /= 5
	return point, nil
}

func RunUP73A() (UP73ACompactFactorizedTrainingBudgetResult, error) {
	steps := []int{200, 400, 800, 1600}
	result := UP73ACompactFactorizedTrainingBudgetResult{
		Schema: UP73ACompactFactorizedTrainingBudgetSchema,
		Experiment: "UP-73A-compact-factorized-training-budget",
		SourceUP71ASeal: "b9738d3d63e591d9d6194fdba43315344dc9dc88",
		TrainingSteps: append([]int(nil), steps...),
		TrainingStates: 54,
		HeldOutStates: 162,
	}
	for _, n := range steps {
		onehot, err := up73aRunPoint(n, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil {
			return UP73ACompactFactorizedTrainingBudgetResult{}, err
		}
		result.Points = append(result.Points, onehot)
		simplex, err := up73aRunPoint(n, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil {
			return UP73ACompactFactorizedTrainingBudgetResult{}, err
		}
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
