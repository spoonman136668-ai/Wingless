package unitary

import "math"

const UP234BFamilyCenteredParitySchema = "wingless.up234b-family-centered-parity.v1"

type UP234BFamilySummary struct {
	Family   string `json:"family"`
	States   int    `json:"states"`
	Correct  int    `json:"correct"`
}

type UP234BResult struct {
	Schema                          string                `json:"schema"`
	Experiment                      string                `json:"experiment"`
	SourceUP233BSeal                string                `json:"source_up233b_seal"`
	TrainingPhases                  []int                 `json:"training_phases"`
	EvaluationPhases                []int                 `json:"evaluation_phases"`
	Features                        []string              `json:"features"`
	TrainingPoints                  int                   `json:"training_points"`
	EvaluationPoints                int                   `json:"evaluation_points"`
	ParentPooledParityAccuracy      float64               `json:"parent_pooled_parity_accuracy"`
	ParentFamilyParityAccuracy      float64               `json:"parent_family_conditioned_parity_accuracy"`
	FamilyCenteredPooledAccuracy    float64               `json:"family_centered_pooled_parity_accuracy"`
	TrainingParityLabelsUsed        bool                  `json:"training_parity_labels_used"`
	FamilyCenteringUsed             bool                  `json:"family_centering_used"`
	FamilySpecificClassifierUsed    bool                  `json:"family_specific_classifier_used"`
	PhaseInputAtInferenceUsed       bool                  `json:"phase_input_at_inference_used"`
	EvaluationLabelFittingUsed      bool                  `json:"evaluation_label_fitting_used"`
	AdaptiveFeatureSelectionUsed    bool                  `json:"adaptive_feature_selection_used"`
	NonlinearClassifierUsed         bool                  `json:"nonlinear_classifier_used"`
	LiveActivation                  bool                  `json:"live_activation"`
	Summaries                       []UP234BFamilySummary `json:"summaries"`
}

func RunUP234B() (UP234BResult, error) {
	parent, err := RunUP233B()
	if err != nil {
		return UP234BResult{}, err
	}
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
	familySum := map[string][4]float64{}
	familyN := map[string]int{}
	for _, f := range specs {
		for _, ph := range train {
			x := up193bFeature(ph, f.Canonical)
			rows = append(rows, row{name: f.Name, phase: ph, x: x})
			a := familySum[f.Name]
			for j := 0; j < 4; j++ {
				a[j] += x[j]
			}
			familySum[f.Name] = a
			familyN[f.Name]++
		}
	}
	familyMean := map[string][4]float64{}
	for _, f := range specs {
		a := familySum[f.Name]
		for j := 0; j < 4; j++ {
			a[j] /= float64(familyN[f.Name])
		}
		familyMean[f.Name] = a
	}
	var std [4]float64
	for _, r := range rows {
		m := familyMean[r.name]
		for j := 0; j < 4; j++ {
			d := r.x[j] - m[j]
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
	parity := func(ph int) string {
		if ph%2 == 0 {
			return "even"
		}
		return "odd"
	}
	for _, r := range rows {
		m := familyMean[r.name]
		z := [4]float64{}
		for j := 0; j < 4; j++ {
			z[j] = (r.x[j] - m[j]) / std[j]
		}
		p := parity(r.phase)
		a := pooled[p]
		for j := 0; j < 4; j++ {
			a.sum[j] += z[j]
		}
		a.n++
		pooled[p] = a
	}
	centroid := func(a accum) (c [4]float64) {
		for j := 0; j < 4; j++ {
			c[j] = a.sum[j] / float64(a.n)
		}
		return
	}
	cent := map[string][4]float64{
		"even": centroid(pooled["even"]),
		"odd":  centroid(pooled["odd"]),
	}
	res := UP234BResult{
		Schema: UP234BFamilyCenteredParitySchema,
		Experiment: "UP-234B-family-centered-parity",
		SourceUP233BSeal: "1ae40d9bb91c4788af66bce6e484885cc7348810",
		TrainingPhases: train,
		EvaluationPhases: eval,
		Features: []string{"native_correct_count", "mean_absolute_margin", "near_zero_margin_count", "min_absolute_margin"},
		TrainingPoints: len(rows),
		ParentPooledParityAccuracy: parent.PooledParityAccuracy,
		ParentFamilyParityAccuracy: parent.FamilyConditionedParityAccuracy,
		TrainingParityLabelsUsed: true,
		FamilyCenteringUsed: true,
		FamilySpecificClassifierUsed: false,
		PhaseInputAtInferenceUsed: false,
		EvaluationLabelFittingUsed: false,
		AdaptiveFeatureSelectionUsed: false,
		NonlinearClassifierUsed: false,
		LiveActivation: false,
	}
	correct := 0
	for _, f := range specs {
		sm := UP234BFamilySummary{Family: f.Name}
		m := familyMean[f.Name]
		for _, ph := range eval {
			x := up193bFeature(ph, f.Canonical)
			z := [4]float64{}
			for j := 0; j < 4; j++ {
				z[j] = (x[j] - m[j]) / std[j]
			}
			best := ""
			bestD := math.Inf(1)
			for _, p := range []string{"even", "odd"} {
				if d := up212bDist(z, cent[p]); d < bestD {
					bestD = d
					best = p
				}
			}
			sm.States++
			res.EvaluationPoints++
			if best == parity(ph) {
				sm.Correct++
				correct++
			}
		}
		res.Summaries = append(res.Summaries, sm)
	}
	if res.EvaluationPoints > 0 {
		res.FamilyCenteredPooledAccuracy = float64(correct) / float64(res.EvaluationPoints)
	}
	return res, nil
}
