package unitary

import "math"

const UP67AEntangledIndexRobustnessSchema = "wingless.up67a-entangled-index-robustness.v1"

type UP67AIndexMap struct {
	Name       string `json:"name"`
	Multiplier int    `json:"multiplier"`
	Offset     int    `json:"offset"`
}

type UP67AIndexPoint struct {
	Map                 UP67AIndexMap     `json:"map"`
	FrequencyPairs      int               `json:"frequency_pairs"`
	FeatureDimension    int               `json:"feature_dimension"`
	TrainStates         int               `json:"train_states"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP67AEntangledIndexRobustnessResult struct {
	Schema               string            `json:"schema"`
	Experiment           string            `json:"experiment"`
	SourceUP66ASeal      string            `json:"source_up66a_seal"`
	TrainingSplitRule    string            `json:"training_split_rule"`
	FrequencyPairs       int               `json:"frequency_pairs"`
	ObserverClassChanged bool              `json:"observer_class_changed"`
	IndexMaps            []UP67AIndexMap   `json:"index_maps"`
	Points               []UP67AIndexPoint `json:"points"`
}

func up67aMappedIndex(r up64aRoles, m UP67AIndexMap) int {
	index, _ := up64aIndex(r)
	return (m.Multiplier*index + m.Offset) % 243
}

func up67aEntangledFeatures(r up64aRoles, m UP67AIndexMap) []float64 {
	const pairs = 32
	index := up67aMappedIndex(r, m)
	out := make([]float64, 2*pairs)
	scale := 1 / math.Sqrt(float64(pairs))
	for f := 1; f <= pairs; f++ {
		phase := 2 * math.Pi * float64(f*index) / 243
		out[2*(f-1)] = scale * math.Cos(phase)
		out[2*(f-1)+1] = scale * math.Sin(phase)
	}
	return out
}

func up67aRunPoint(m UP67AIndexMap) (UP67AIndexPoint, error) {
	all := up64aAllRoles()
	var trainRoles, heldRoles []up64aRoles
	for _, r := range all {
		if up65aTrainState(r) {
			trainRoles = append(trainRoles, r)
		} else {
			heldRoles = append(heldRoles, r)
		}
	}
	point := UP67AIndexPoint{
		Map: m, FrequencyPairs: 32, FeatureDimension: 64,
		TrainStates: len(trainRoles), HeldOutStates: len(heldRoles), Gate: true,
	}
	for role := 0; role < 5; role++ {
		train := make([]headSample, 0, len(trainRoles))
		held := make([]headSample, 0, len(heldRoles))
		for _, r := range trainRoles {
			train = append(train, headSample{features: up67aEntangledFeatures(r, m), target: r[role]})
		}
		for _, r := range heldRoles {
			held = append(held, headSample{features: up67aEntangledFeatures(r, m), target: r[role]})
		}
		head, metric, err := trainLinearSoftmax(train, 3, 64, 800, 1.0)
		if err != nil {
			return UP67AIndexPoint{}, err
		}
		_, heldAcc, err := evaluateHead(head, held)
		if err != nil {
			return UP67AIndexPoint{}, err
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

func RunUP67A() (UP67AEntangledIndexRobustnessResult, error) {
	maps := []UP67AIndexMap{
		{Name: "lexicographic", Multiplier: 1, Offset: 0},
		{Name: "affine_2_17", Multiplier: 2, Offset: 17},
		{Name: "affine_80_31", Multiplier: 80, Offset: 31},
		{Name: "affine_241_7", Multiplier: 241, Offset: 7},
	}
	result := UP67AEntangledIndexRobustnessResult{
		Schema: UP67AEntangledIndexRobustnessSchema,
		Experiment: "UP-67A-entangled-index-robustness",
		SourceUP66ASeal: "a4a00334e101f95ef1b25218058834183526ce69",
		TrainingSplitRule: "sum(role values) mod 3 == 0",
		FrequencyPairs: 32,
		ObserverClassChanged: false,
		IndexMaps: append([]UP67AIndexMap(nil), maps...),
	}
	for _, m := range maps {
		p, err := up67aRunPoint(m)
		if err != nil {
			return UP67AEntangledIndexRobustnessResult{}, err
		}
		result.Points = append(result.Points, p)
	}
	return result, nil
}
