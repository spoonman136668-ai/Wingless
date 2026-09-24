package unitary

import "math"

const UP69AFactorizedCapacityAblationSchema = "wingless.up69a-factorized-capacity-ablation.v1"

type UP69AArm struct {
	Name                string            `json:"name"`
	FeatureDimension    int               `json:"feature_dimension"`
	TrainStates         int               `json:"train_states"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP69AFactorizedCapacityAblationResult struct {
	Schema          string     `json:"schema"`
	Experiment      string     `json:"experiment"`
	SourceUP67ASeal string     `json:"source_up67a_seal"`
	TrainingSplitRule string   `json:"training_split_rule"`
	Arms            []UP69AArm `json:"arms"`
}

func up69aRawOneHot(r up64aRoles) []float64 {
	out := make([]float64, 15)
	for role, value := range r {
		out[role*3+value] = 1
	}
	return out
}

func up69aSimplex(r up64aRoles) []float64 {
	out := make([]float64, 10)
	for role, value := range r {
		phase := 2 * math.Pi * float64(value) / 3
		out[2*role] = math.Cos(phase)
		out[2*role+1] = math.Sin(phase)
	}
	return out
}

func up69aRunArm(name string, dim int, featureFn func(up64aRoles) []float64) (UP69AArm, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		if up65aTrainState(r) {
			trainRoles = append(trainRoles, r)
		} else {
			heldRoles = append(heldRoles, r)
		}
	}
	arm := UP69AArm{
		Name: name, FeatureDimension: dim,
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
		head, metric, err := trainLinearSoftmax(train, 3, dim, 800, 1.0)
		if err != nil {
			return UP69AArm{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP69AArm{}, err
		}
		arm.PerRole = append(arm.PerRole, UP65ARoleMetric{Role: role, TrainAccuracy: metric.TrainAccuracy, HeldOutAccuracy: heldAcc})
		arm.MeanTrainAccuracy += metric.TrainAccuracy
		arm.MeanHeldOutAccuracy += heldAcc
		if heldAcc < 0.98 {
			arm.Gate = false
		}
	}
	arm.MeanTrainAccuracy /= 5
	arm.MeanHeldOutAccuracy /= 5
	return arm, nil
}

func RunUP69A() (UP69AFactorizedCapacityAblationResult, error) {
	specs := []struct {
		name string
		dim  int
		fn   func(up64aRoles) []float64
	}{
		{"raw_factorized_one_hot", 15, up69aRawOneHot},
		{"role_factorized_simplex", 10, up69aSimplex},
		{"dense_hadamard_factorized", 64, up65aHadamardFeatures},
		{"joint_fourier_entangled", 64, up65aEntangledFeatures},
	}
	result := UP69AFactorizedCapacityAblationResult{
		Schema: UP69AFactorizedCapacityAblationSchema,
		Experiment: "UP-69A-factorized-capacity-ablation",
		SourceUP67ASeal: "1c055b437255a40aa0158a97f981d0f7a094e9d4",
		TrainingSplitRule: "sum(role values) mod 3 == 0",
	}
	for _, spec := range specs {
		arm, err := up69aRunArm(spec.name, spec.dim, spec.fn)
		if err != nil {
			return UP69AFactorizedCapacityAblationResult{}, err
		}
		result.Arms = append(result.Arms, arm)
	}
	return result, nil
}
