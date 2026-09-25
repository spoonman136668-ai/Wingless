package unitary

const UP78AOneStepCoverageBoundarySchema = "wingless.up78a-one-step-coverage-boundary.v1"

type UP78AOneStepCoverageBoundaryPoint struct {
	TrainingStates      int               `json:"training_states"`
	TrainingSteps       int               `json:"training_steps"`
	Representation      string            `json:"representation"`
	FeatureDimension    int               `json:"feature_dimension"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP78AOneStepCoverageBoundaryResult struct {
	Schema          string                              `json:"schema"`
	Experiment      string                              `json:"experiment"`
	SourceUP77ASeal string                              `json:"source_up77a_seal"`
	TrainingLevels  []int                               `json:"training_levels"`
	TrainingSteps   int                                 `json:"training_steps"`
	LearningRate    float64                             `json:"learning_rate"`
	HeldOutStates   int                                 `json:"heldout_states"`
	Points          []UP78AOneStepCoverageBoundaryPoint `json:"points"`
}

func up78aTrainAtLevel(r up64aRoles, level int) bool {
	sum := 0
	for _, v := range r {
		sum += v
	}
	if sum%3 != 0 {
		return false
	}
	s := up71aSecondary(r)
	t := up71aTertiary(r)
	switch level {
	case 27:
		return s == 0
	case 36:
		return s == 0 || (s == 1 && t == 0)
	case 45:
		return s == 0 || (s == 1 && t < 2)
	case 54:
		return s < 2
	default:
		return false
	}
}

func up78aRunPoint(level int, name string, dim int, featureFn func(up64aRoles) []float64) (UP78AOneStepCoverageBoundaryPoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r {
			sum += v
		}
		if up78aTrainAtLevel(r, level) {
			trainRoles = append(trainRoles, r)
		}
		if sum%3 != 0 {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP78AOneStepCoverageBoundaryPoint{
		TrainingStates: len(trainRoles), TrainingSteps: 1, Representation: name,
		FeatureDimension: dim, HeldOutStates: len(heldRoles), Gate: true,
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
		head, metric, err := trainLinearSoftmax(train, 3, dim, 1, 1.0)
		if err != nil {
			return UP78AOneStepCoverageBoundaryPoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP78AOneStepCoverageBoundaryPoint{}, err
		}
		point.PerRole = append(point.PerRole, UP65ARoleMetric{Role: role, TrainAccuracy: metric.TrainAccuracy, HeldOutAccuracy: heldAcc})
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

func RunUP78A() (UP78AOneStepCoverageBoundaryResult, error) {
	levels := []int{27, 36, 45, 54}
	result := UP78AOneStepCoverageBoundaryResult{
		Schema: UP78AOneStepCoverageBoundarySchema,
		Experiment: "UP-78A-one-step-coverage-boundary",
		SourceUP77ASeal: "e90596fcf15feef8ab13523e9a6bda5dc52feff9",
		TrainingLevels: append([]int(nil), levels...),
		TrainingSteps: 1,
		LearningRate: 1.0,
		HeldOutStates: 162,
	}
	for _, level := range levels {
		onehot, err := up78aRunPoint(level, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil {
			return UP78AOneStepCoverageBoundaryResult{}, err
		}
		result.Points = append(result.Points, onehot)
		simplex, err := up78aRunPoint(level, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil {
			return UP78AOneStepCoverageBoundaryResult{}, err
		}
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
