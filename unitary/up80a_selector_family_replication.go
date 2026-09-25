package unitary

const UP80ASelectorFamilyReplicationSchema = "wingless.up80a-selector-family-replication.v1"

type UP80ASelectorFamily struct {
	Name       string `json:"name"`
	Secondary  [5]int `json:"secondary"`
	Tertiary   [5]int `json:"tertiary"`
}

type UP80ASelectorFamilyPoint struct {
	SelectorFamily       string            `json:"selector_family"`
	TrainingStates       int               `json:"training_states"`
	TrainingSteps        int               `json:"training_steps"`
	Representation       string            `json:"representation"`
	FeatureDimension     int               `json:"feature_dimension"`
	HeldOutStates        int               `json:"heldout_states"`
	PerRole              []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy    float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy  float64           `json:"mean_heldout_accuracy"`
	Gate                 bool              `json:"gate"`
}

type UP80ASelectorFamilyReplicationResult struct {
	Schema          string                         `json:"schema"`
	Experiment      string                         `json:"experiment"`
	SourceUP79ASeal string                         `json:"source_up79a_seal"`
	Families        []UP80ASelectorFamily          `json:"families"`
	TrainingSteps   int                            `json:"training_steps"`
	LearningRate    float64                        `json:"learning_rate"`
	HeldOutStates   int                            `json:"heldout_states"`
	Points          []UP80ASelectorFamilyPoint     `json:"points"`
}

func up80aLinear(r up64aRoles, coeff [5]int) int {
	v := 0
	for i := 0; i < 5; i++ {
		v += coeff[i] * r[i]
	}
	v %= 3
	if v < 0 {
		v += 3
	}
	return v
}

func up80aTrainInFamily(r up64aRoles, family UP80ASelectorFamily) bool {
	sum := 0
	for _, v := range r {
		sum += v
	}
	if sum%3 != 0 {
		return false
	}
	s := up80aLinear(r, family.Secondary)
	t := up80aLinear(r, family.Tertiary)
	return s == 0 || (s == 1 && t == 0)
}

func up80aRunPoint(family UP80ASelectorFamily, name string, dim int, featureFn func(up64aRoles) []float64) (UP80ASelectorFamilyPoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r {
			sum += v
		}
		if up80aTrainInFamily(r, family) {
			trainRoles = append(trainRoles, r)
		}
		if sum%3 != 0 {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP80ASelectorFamilyPoint{
		SelectorFamily: family.Name,
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
			return UP80ASelectorFamilyPoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP80ASelectorFamilyPoint{}, err
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

func RunUP80A() (UP80ASelectorFamilyReplicationResult, error) {
	families := []UP80ASelectorFamily{
		{Name: "family_a", Secondary: [5]int{1,2,1,2,1}, Tertiary: [5]int{1,1,2,2,0}},
		{Name: "family_b", Secondary: [5]int{2,1,2,1,2}, Tertiary: [5]int{1,2,2,1,0}},
		{Name: "family_c", Secondary: [5]int{1,1,2,1,2}, Tertiary: [5]int{2,1,1,2,0}},
	}
	result := UP80ASelectorFamilyReplicationResult{
		Schema: UP80ASelectorFamilyReplicationSchema,
		Experiment: "UP-80A-selector-family-replication",
		SourceUP79ASeal: "145f5730c4e86bced6f31e9a6bdb29eb25b03c56",
		Families: append([]UP80ASelectorFamily(nil), families...),
		TrainingSteps: 1,
		LearningRate: 1.0,
		HeldOutStates: 162,
	}
	for _, family := range families {
		onehot, err := up80aRunPoint(family, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil {
			return UP80ASelectorFamilyReplicationResult{}, err
		}
		result.Points = append(result.Points, onehot)
		simplex, err := up80aRunPoint(family, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil {
			return UP80ASelectorFamilyReplicationResult{}, err
		}
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
