package unitary

import "math"

const UP233BNativeParityDecodabilitySchema = "wingless.up233b-native-parity-decodability.v1"

type UP233BFamilySummary struct {
	Family                   string `json:"family"`
	EvaluationStates         int    `json:"evaluation_states"`
	PooledCorrect            int    `json:"pooled_correct"`
	FamilyConditionedCorrect int    `json:"family_conditioned_correct"`
}

type UP233BResult struct {
	Schema                          string                `json:"schema"`
	Experiment                      string                `json:"experiment"`
	SourceUP232BSeal                string                `json:"source_up232b_seal"`
	TrainingPhases                  []int                 `json:"training_phases"`
	EvaluationPhases                []int                 `json:"evaluation_phases"`
	Features                        []string              `json:"features"`
	TrainingPoints                  int                   `json:"training_points"`
	EvaluationPoints                int                   `json:"evaluation_points"`
	PooledParityAccuracy            float64               `json:"pooled_parity_accuracy"`
	FamilyConditionedParityAccuracy float64               `json:"family_conditioned_parity_accuracy"`
	TrainingParityLabelsUsed        bool                  `json:"training_parity_labels_used"`
	PhaseInputAtInferenceUsed       bool                  `json:"phase_input_at_inference_used"`
	EvaluationLabelFittingUsed      bool                  `json:"evaluation_label_fitting_used"`
	AdaptiveFeatureSelectionUsed    bool                  `json:"adaptive_feature_selection_used"`
	NonlinearClassifierUsed         bool                  `json:"nonlinear_classifier_used"`
	LiveActivation                  bool                  `json:"live_activation"`
	Summaries                       []UP233BFamilySummary `json:"summaries"`
}

func RunUP233B() (UP233BResult, error) {
	train := []int{58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72}
	eval := []int{121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152}
	specs := []UP197BComposition{
		{Name: "mixed4", Canonical: []int{0, 5, 1, 6}},
		{Name: "observe4", Canonical: []int{5, 6, 7, 8}},
		{Name: "store4", Canonical: []int{0, 1, 2, 3}},
		{Name: "cross3", Canonical: []int{0, 5, 13}},
	}
	type row struct {
		name  string
		phase int
		x     [4]float64
	}
	rows := []row{}
	for _, f := range specs {
		for _, ph := range train {
			rows = append(rows, row{name: f.Name, phase: ph, x: up193bFeature(ph, f.Canonical)})
		}
	}
	var mean, std [4]float64
	for _, r := range rows {
		for j := 0; j < 4; j++ {
			mean[j] += r.x[j]
		}
	}
	for j := 0; j < 4; j++ {
		mean[j] /= float64(len(rows))
	}
	for _, r := range rows {
		for j := 0; j < 4; j++ {
			d := r.x[j] - mean[j]
			std[j] += d * d
		}
	}
	for j := 0; j < 4; j++ {
		std[j] = math.Sqrt(std[j] / float64(len(rows)))
		if std[j] == 0 {
			std[j] = 1
		}
	}
	type accum struct {
		sum [4]float64
		n   int
	}
	pooled := map[string]accum{}
	family := map[string]accum{}
	parity := func(ph int) string {
		if ph%2 == 0 {
			return "even"
		}
		return "odd"
	}
	for _, r := range rows {
		z := [4]float64{}
		for j := 0; j < 4; j++ {
			z[j] = (r.x[j] - mean[j]) / std[j]
		}
		p := parity(r.phase)
		a := pooled[p]
		for j := 0; j < 4; j++ {
			a.sum[j] += z[j]
		}
		a.n++
		pooled[p] = a
		k := r.name + "|" + p
		b := family[k]
		for j := 0; j < 4; j++ {
			b.sum[j] += z[j]
		}
		b.n++
		family[k] = b
	}
	centroid := func(a accum) (c [4]float64) {
		for j := 0; j < 4; j++ {
			c[j] = a.sum[j] / float64(a.n)
		}
		return
	}
	pooledCent := map[string][4]float64{"even": centroid(pooled["even"]), "odd": centroid(pooled["odd"])}
	familyCent := map[string][4]float64{}
	for k, a := range family {
		familyCent[k] = centroid(a)
	}
	res := UP233BResult{
		Schema: UP233BNativeParityDecodabilitySchema, Experiment: "UP-233B-native-parity-decodability", SourceUP232BSeal: "57de3b873b33b2ad2165bc7ce1fab24aa14065c3",
		TrainingPhases: train, EvaluationPhases: eval, Features: []string{"native_correct_count", "mean_absolute_margin", "near_zero_margin_count", "min_absolute_margin"}, TrainingPoints: len(rows),
		TrainingParityLabelsUsed: true, PhaseInputAtInferenceUsed: false, EvaluationLabelFittingUsed: false, AdaptiveFeatureSelectionUsed: false, NonlinearClassifierUsed: false, LiveActivation: false,
	}
	pooledTotal, familyTotal := 0, 0
	for _, f := range specs {
		s := UP233BFamilySummary{Family: f.Name}
		for _, ph := range eval {
			x := up193bFeature(ph, f.Canonical)
			z := [4]float64{}
			for j := 0; j < 4; j++ {
				z[j] = (x[j] - mean[j]) / std[j]
			}
			bestP := ""
			bestPD := math.Inf(1)
			for _, p := range []string{"even", "odd"} {
				if d := up212bDist(z, pooledCent[p]); d < bestPD {
					bestPD = d
					bestP = p
				}
			}
			bestF := ""
			bestFD := math.Inf(1)
			for _, p := range []string{"even", "odd"} {
				if d := up212bDist(z, familyCent[f.Name+"|"+p]); d < bestFD {
					bestFD = d
					bestF = p
				}
			}
			want := parity(ph)
			s.EvaluationStates++
			res.EvaluationPoints++
			if bestP == want {
				s.PooledCorrect++
				pooledTotal++
			}
			if bestF == want {
				s.FamilyConditionedCorrect++
				familyTotal++
			}
		}
		res.Summaries = append(res.Summaries, s)
	}
	if res.EvaluationPoints > 0 {
		res.PooledParityAccuracy = float64(pooledTotal) / float64(res.EvaluationPoints)
		res.FamilyConditionedParityAccuracy = float64(familyTotal) / float64(res.EvaluationPoints)
	}
	return res, nil
}
