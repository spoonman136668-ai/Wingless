package unitary

const UP72ACompactFactorized54SplitReplicationSchema = "wingless.up72a-compact-factorized-54-split-replication.v1"

type UP72ASplitPoint struct {
	TrainingResidue      int               `json:"training_residue"`
	TrainingStates       int               `json:"training_states"`
	Representation       string            `json:"representation"`
	FeatureDimension     int               `json:"feature_dimension"`
	HeldOutStates        int               `json:"heldout_states"`
	PerRole              []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy    float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy  float64           `json:"mean_heldout_accuracy"`
	Gate                 bool              `json:"gate"`
}

type UP72ACompactFactorized54SplitReplicationResult struct {
	Schema          string            `json:"schema"`
	Experiment      string            `json:"experiment"`
	SourceUP71ASeal string            `json:"source_up71a_seal"`
	TrainingResidues []int            `json:"training_residues"`
	TrainingStates  int               `json:"training_states"`
	HeldOutStates   int               `json:"heldout_states"`
	Points          []UP72ASplitPoint `json:"points"`
}

func up72aRunPoint(residue int, name string, dim int, featureFn func(up64aRoles) []float64) (UP72ASplitPoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		sum := 0
		for _, v := range r {
			sum += v
		}
		if sum%3 == residue && up71aSecondary(r) < 2 {
			trainRoles = append(trainRoles, r)
		}
		if sum%3 != residue {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP72ASplitPoint{
		TrainingResidue: residue,
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
		head, metric, err := trainLinearSoftmax(train, 3, dim, 800, 1.0)
		if err != nil {
			return UP72ASplitPoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP72ASplitPoint{}, err
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

func RunUP72A() (UP72ACompactFactorized54SplitReplicationResult, error) {
	residues := []int{0, 1, 2}
	result := UP72ACompactFactorized54SplitReplicationResult{
		Schema: UP72ACompactFactorized54SplitReplicationSchema,
		Experiment: "UP-72A-compact-factorized-54-split-replication",
		SourceUP71ASeal: "b9738d3d63e591d9d6194fdba43315344dc9dc88",
		TrainingResidues: append([]int(nil), residues...),
		TrainingStates: 54,
		HeldOutStates: 162,
	}
	for _, residue := range residues {
		onehot, err := up72aRunPoint(residue, "raw_factorized_one_hot", 15, up69aRawOneHot)
		if err != nil {
			return UP72ACompactFactorized54SplitReplicationResult{}, err
		}
		result.Points = append(result.Points, onehot)
		simplex, err := up72aRunPoint(residue, "role_factorized_simplex", 10, up69aSimplex)
		if err != nil {
			return UP72ACompactFactorized54SplitReplicationResult{}, err
		}
		result.Points = append(result.Points, simplex)
	}
	return result, nil
}
