package unitary

const UP70ATrainingCoverageLadderSchema = "wingless.up70a-training-coverage-ladder.v1"

type UP70ACoveragePoint struct {
	TrainingStates      int               `json:"training_states"`
	Representation      string            `json:"representation"`
	FeatureDimension    int               `json:"feature_dimension"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP70ATrainingCoverageLadderResult struct {
	Schema             string               `json:"schema"`
	Experiment         string               `json:"experiment"`
	SourceUP68ASeal    string               `json:"source_up68a_seal"`
	TrainingLevels     []int                `json:"training_levels"`
	HeldOutStates      int                  `json:"heldout_states"`
	ObserverClassChanged bool               `json:"observer_class_changed"`
	Points             []UP70ACoveragePoint `json:"points"`
}

func up70aSecondary(r up64aRoles) int {
	return (r[0] + 2*r[1] + r[2] + 2*r[3] + r[4]) % 3
}

func up70aTrainAtLevel(r up64aRoles, bands int) bool {
	sum := 0
	for _, v := range r {
		sum += v
	}
	if sum%3 != 0 {
		return false
	}
	return up70aSecondary(r) < bands
}

func up70aRunPoint(bands int, name string, featureFn func(up64aRoles) []float64) (UP70ACoveragePoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r {
			sum += v
		}
		if up70aTrainAtLevel(r, bands) {
			trainRoles = append(trainRoles, r)
		}
		if sum%3 != 0 {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP70ACoveragePoint{
		TrainingStates: len(trainRoles), Representation: name, FeatureDimension: 64,
		HeldOutStates: len(heldRoles), Gate: true,
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
		head, metric, err := trainLinearSoftmax(train, 3, 64, 800, 1.0)
		if err != nil {
			return UP70ACoveragePoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP70ACoveragePoint{}, err
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

func RunUP70A() (UP70ATrainingCoverageLadderResult, error) {
	bands := []int{1, 2, 3}
	result := UP70ATrainingCoverageLadderResult{
		Schema: UP70ATrainingCoverageLadderSchema,
		Experiment: "UP-70A-training-coverage-ladder",
		SourceUP68ASeal: "49f8deb62b063ca5fd8f059867c0e82d1a5e9313",
		TrainingLevels: []int{27, 54, 81},
		HeldOutStates: 162,
		ObserverClassChanged: false,
	}
	for _, b := range bands {
		factorized, err := up70aRunPoint(b, "dense_hadamard_factorized", up65aHadamardFeatures)
		if err != nil {
			return UP70ATrainingCoverageLadderResult{}, err
		}
		result.Points = append(result.Points, factorized)
		entangled, err := up70aRunPoint(b, "joint_fourier_entangled", up65aEntangledFeatures)
		if err != nil {
			return UP70ATrainingCoverageLadderResult{}, err
		}
		result.Points = append(result.Points, entangled)
	}
	return result, nil
}
