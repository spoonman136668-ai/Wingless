package unitary

const UP77AOneStepCoverageLadderSchema = "wingless.up77a-one-step-coverage-ladder.v1"

type UP77AOneStepCoveragePoint struct {
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

type UP77AOneStepCoverageLadderResult struct {
	Schema          string                      `json:"schema"`
	Experiment      string                      `json:"experiment"`
	SourceUP76ASeal string                      `json:"source_up76a_seal"`
	TrainingLevels  []int                       `json:"training_levels"`
	TrainingSteps   int                         `json:"training_steps"`
	LearningRate    float64                     `json:"learning_rate"`
	HeldOutStates   int                         `json:"heldout_states"`
	Points          []UP77AOneStepCoveragePoint `json:"points"`
}

func up77aRunPoint(level int, name string, dim int, featureFn func(up64aRoles) []float64) (UP77AOneStepCoveragePoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r { sum += v }
		if up71aTrainAtLevel(r, level) { trainRoles = append(trainRoles, r) }
		if sum%3 != 0 { heldRoles = append(heldRoles, r) }
	}
	point := UP77AOneStepCoveragePoint{
		TrainingStates: len(trainRoles), TrainingSteps: 1, Representation: name,
		FeatureDimension: dim, HeldOutStates: len(heldRoles), Gate: true,
	}
	for role := 0; role < 5; role++ {
		train := make([]headSample, 0, len(trainRoles))
		held := make([]headSample, 0, len(heldRoles))
		for _, r := range trainRoles { train = append(train, headSample{features: featureFn(r), target: r[role]}) }
		for _, r := range heldRoles { held = append(held, headSample{features: featureFn(r), target: r[role]}) }
		head, metric, err := trainLinearSoftmax(train, 3, dim, 1, 1.0)
		if err != nil { return UP77AOneStepCoveragePoint{}, err }
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil { return UP77AOneStepCoveragePoint{}, err }
		point.PerRole = append(point.PerRole, UP65ARoleMetric{Role: role, TrainAccuracy: metric.TrainAccuracy, HeldOutAccuracy: heldAcc})
		point.MeanTrainAccuracy += metric.TrainAccuracy
		point.MeanHeldOutAccuracy += heldAcc
		if heldAcc < 0.98 { point.Gate = false }
	}
	point.MeanTrainAccuracy /= 5
	point.MeanHeldOutAccuracy /= 5
	return point, nil
}

func RunUP77A() (UP77AOneStepCoverageLadderResult, error) {
	levels := []int{9, 18, 27, 54, 81}
	result := UP77AOneStepCoverageLadderResult{
		Schema: UP77AOneStepCoverageLadderSchema, Experiment: "UP-77A-one-step-coverage-ladder",
		SourceUP76ASeal: "5616cd22352031edd2da1c37a3cac81fd8cb4574",
		TrainingLevels: append([]int(nil), levels...), TrainingSteps: 1, LearningRate: 1.0, HeldOutStates: 162,
	}
	for _, level := range levels {
		onehot, err := up77aRunPoint(level, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil { return UP77AOneStepCoverageLadderResult{}, err }
		result.Points = append(result.Points, onehot)
		simplex, err := up77aRunPoint(level, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil { return UP77AOneStepCoverageLadderResult{}, err }
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
