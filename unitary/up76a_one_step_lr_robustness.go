package unitary

const UP76AOneStepLRRobustnessSchema = "wingless.up76a-one-step-lr-robustness.v1"

type UP76AOneStepLRPoint struct {
	LearningRate        float64           `json:"learning_rate"`
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

type UP76AOneStepLRRobustnessResult struct {
	Schema          string                `json:"schema"`
	Experiment      string                `json:"experiment"`
	SourceUP75ASeal string                `json:"source_up75a_seal"`
	LearningRates   []float64             `json:"learning_rates"`
	TrainingSteps   int                   `json:"training_steps"`
	TrainingStates  int                   `json:"training_states"`
	HeldOutStates   int                   `json:"heldout_states"`
	Points          []UP76AOneStepLRPoint `json:"points"`
}

func up76aRunPoint(lr float64, name string, dim int, featureFn func(up64aRoles) []float64) (UP76AOneStepLRPoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r { sum += v }
		if sum%3 == 0 && up71aSecondary(r) < 2 { trainRoles = append(trainRoles, r) }
		if sum%3 != 0 { heldRoles = append(heldRoles, r) }
	}
	point := UP76AOneStepLRPoint{
		LearningRate: lr, TrainingSteps: 1, TrainingStates: len(trainRoles),
		Representation: name, FeatureDimension: dim, HeldOutStates: len(heldRoles), Gate: true,
	}
	for role := 0; role < 5; role++ {
		train := make([]headSample, 0, len(trainRoles))
		held := make([]headSample, 0, len(heldRoles))
		for _, r := range trainRoles { train = append(train, headSample{features: featureFn(r), target: r[role]}) }
		for _, r := range heldRoles { held = append(held, headSample{features: featureFn(r), target: r[role]}) }
		head, metric, err := trainLinearSoftmax(train, 3, dim, 1, lr)
		if err != nil { return UP76AOneStepLRPoint{}, err }
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil { return UP76AOneStepLRPoint{}, err }
		point.PerRole = append(point.PerRole, UP65ARoleMetric{Role: role, TrainAccuracy: metric.TrainAccuracy, HeldOutAccuracy: heldAcc})
		point.MeanTrainAccuracy += metric.TrainAccuracy
		point.MeanHeldOutAccuracy += heldAcc
		if heldAcc < 0.98 { point.Gate = false }
	}
	point.MeanTrainAccuracy /= 5
	point.MeanHeldOutAccuracy /= 5
	return point, nil
}

func RunUP76A() (UP76AOneStepLRRobustnessResult, error) {
	rates := []float64{0.125, 0.25, 0.5, 1.0, 2.0}
	result := UP76AOneStepLRRobustnessResult{
		Schema: UP76AOneStepLRRobustnessSchema, Experiment: "UP-76A-one-step-lr-robustness",
		SourceUP75ASeal: "8b0773a354f04bc73e6960a0c7ec197349ae46a1",
		LearningRates: append([]float64(nil), rates...), TrainingSteps: 1, TrainingStates: 54, HeldOutStates: 162,
	}
	for _, lr := range rates {
		onehot, err := up76aRunPoint(lr, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil { return UP76AOneStepLRRobustnessResult{}, err }
		result.Points = append(result.Points, onehot)
		simplex, err := up76aRunPoint(lr, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil { return UP76AOneStepLRRobustnessResult{}, err }
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
