package unitary

const UP79ACoverageSlicesSchema = "wingless.up79a-coverage-slices.v1"

type UP79ACoverageSlicePoint struct {
	TertiarySlice       int               `json:"tertiary_slice"`
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

type UP79ACoverageSlicesResult struct {
	Schema          string                    `json:"schema"`
	Experiment      string                    `json:"experiment"`
	SourceUP78ASeal string                    `json:"source_up78a_seal"`
	TertiarySlices  []int                     `json:"tertiary_slices"`
	TrainingSteps   int                       `json:"training_steps"`
	LearningRate    float64                   `json:"learning_rate"`
	HeldOutStates   int                       `json:"heldout_states"`
	Points          []UP79ACoverageSlicePoint `json:"points"`
}

func up79aTrainInSlice(r up64aRoles, tertiarySlice int) bool {
	sum := 0
	for _, v := range r {
		sum += v
	}
	if sum%3 != 0 {
		return false
	}
	s := up71aSecondary(r)
	t := up71aTertiary(r)
	return s == 0 || (s == 1 && t == tertiarySlice)
}

func up79aRunPoint(tertiarySlice int, name string, dim int, featureFn func(up64aRoles) []float64) (UP79ACoverageSlicePoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r {
			sum += v
		}
		if up79aTrainInSlice(r, tertiarySlice) {
			trainRoles = append(trainRoles, r)
		}
		if sum%3 != 0 {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP79ACoverageSlicePoint{
		TertiarySlice: tertiarySlice,
		TrainingStates: len(trainRoles),
		TrainingSteps: 1,
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
		head, metric, err := trainLinearSoftmax(train, 3, dim, 1, 1.0)
		if err != nil {
			return UP79ACoverageSlicePoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP79ACoverageSlicePoint{}, err
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

func RunUP79A() (UP79ACoverageSlicesResult, error) {
	slices := []int{0, 1, 2}
	result := UP79ACoverageSlicesResult{
		Schema: UP79ACoverageSlicesSchema,
		Experiment: "UP-79A-coverage-slices",
		SourceUP78ASeal: "2a35f521daa7898d714959a602994f5b04b15580",
		TertiarySlices: append([]int(nil), slices...),
		TrainingSteps: 1,
		LearningRate: 1.0,
		HeldOutStates: 162,
	}
	for _, slice := range slices {
		onehot, err := up79aRunPoint(slice, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil {
			return UP79ACoverageSlicesResult{}, err
		}
		result.Points = append(result.Points, onehot)
		simplex, err := up79aRunPoint(slice, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil {
			return UP79ACoverageSlicesResult{}, err
		}
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
