package unitary

const UP71ACompactFactorizedCoverageSchema = "wingless.up71a-compact-factorized-coverage.v1"

type UP71ACoveragePoint struct {
	TrainingStates      int               `json:"training_states"`
	Representation      string            `json:"representation"`
	FeatureDimension    int               `json:"feature_dimension"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP71ACompactFactorizedCoverageResult struct {
	Schema          string               `json:"schema"`
	Experiment      string               `json:"experiment"`
	SourceUP69ASeal string               `json:"source_up69a_seal"`
	TrainingLevels  []int                `json:"training_levels"`
	HeldOutStates   int                  `json:"heldout_states"`
	Points          []UP71ACoveragePoint `json:"points"`
}

func up71aSecondary(r up64aRoles) int {
	return (r[0] + 2*r[1] + r[2] + 2*r[3] + r[4]) % 3
}

func up71aTertiary(r up64aRoles) int {
	return (r[0] + r[1] + 2*r[2] + 2*r[3]) % 3
}

func up71aTrainAtLevel(r up64aRoles, level int) bool {
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
	case 9:
		return s == 0 && t == 0
	case 18:
		return s == 0 && t < 2
	case 27:
		return s == 0
	case 54:
		return s < 2
	case 81:
		return true
	default:
		return false
	}
}

func up71aRunPoint(level int, name string, dim int, featureFn func(up64aRoles) []float64) (UP71ACoveragePoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r {
			sum += v
		}
		if up71aTrainAtLevel(r, level) {
			trainRoles = append(trainRoles, r)
		}
		if sum%3 != 0 {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP71ACoveragePoint{
		TrainingStates: len(trainRoles), Representation: name, FeatureDimension: dim,
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
		head, metric, err := trainLinearSoftmax(train, 3, dim, 800, 1.0)
		if err != nil {
			return UP71ACoveragePoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP71ACoveragePoint{}, err
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

func RunUP71A() (UP71ACompactFactorizedCoverageResult, error) {
	levels := []int{9, 18, 27, 54, 81}
	result := UP71ACompactFactorizedCoverageResult{
		Schema: UP71ACompactFactorizedCoverageSchema,
		Experiment: "UP-71A-compact-factorized-coverage",
		SourceUP69ASeal: "d0b5ccd3566141b090b88c4a033f4c924a9a11b2",
		TrainingLevels: append([]int(nil), levels...),
		HeldOutStates: 162,
	}
	for _, level := range levels {
		onehot, err := up71aRunPoint(level, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil {
			return UP71ACompactFactorizedCoverageResult{}, err
		}
		result.Points = append(result.Points, onehot)
		simplex, err := up71aRunPoint(level, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil {
			return UP71ACompactFactorizedCoverageResult{}, err
		}
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
