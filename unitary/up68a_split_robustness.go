package unitary

const UP68ASplitRobustnessSchema = "wingless.up68a-split-robustness.v1"

type UP68ASplitPoint struct {
	Residue             int               `json:"residue"`
	Representation      string            `json:"representation"`
	FeatureDimension    int               `json:"feature_dimension"`
	TrainStates         int               `json:"train_states"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP68ASplitRobustnessResult struct {
	Schema               string            `json:"schema"`
	Experiment           string            `json:"experiment"`
	SourceUP66ASeal      string            `json:"source_up66a_seal"`
	Residues             []int             `json:"residues"`
	ObserverClassChanged bool              `json:"observer_class_changed"`
	Points               []UP68ASplitPoint `json:"points"`
}

func up68aTrainState(r up64aRoles, residue int) bool {
	sum := 0
	for _, v := range r {
		sum += v
	}
	return sum%3 == residue
}

func up68aRunPoint(residue int, name string, featureFn func(up64aRoles) []float64) (UP68ASplitPoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		if up68aTrainState(r, residue) {
			trainRoles = append(trainRoles, r)
		} else {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP68ASplitPoint{
		Residue: residue, Representation: name, FeatureDimension: 64,
		TrainStates: len(trainRoles), HeldOutStates: len(heldRoles), Gate: true,
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
			return UP68ASplitPoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP68ASplitPoint{}, err
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

func RunUP68A() (UP68ASplitRobustnessResult, error) {
	residues := []int{0, 1, 2}
	result := UP68ASplitRobustnessResult{
		Schema: UP68ASplitRobustnessSchema,
		Experiment: "UP-68A-split-robustness",
		SourceUP66ASeal: "a4a00334e101f95ef1b25218058834183526ce69",
		Residues: append([]int(nil), residues...),
		ObserverClassChanged: false,
	}
	for _, residue := range residues {
		factorized, err := up68aRunPoint(residue, "dense_hadamard_factorized", up65aHadamardFeatures)
		if err != nil {
			return UP68ASplitRobustnessResult{}, err
		}
		result.Points = append(result.Points, factorized)
		entangled, err := up68aRunPoint(residue, "joint_fourier_entangled", up65aEntangledFeatures)
		if err != nil {
			return UP68ASplitRobustnessResult{}, err
		}
		result.Points = append(result.Points, entangled)
	}
	return result, nil
}
