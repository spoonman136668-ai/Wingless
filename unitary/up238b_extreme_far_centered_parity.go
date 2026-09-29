package unitary

import "math"

const UP238BExtremeParitySchema = "wingless.up238b-extreme-far-centered-parity.v1"

type UP238BFamilySummary struct {
	Family  string `json:"family"`
	States  int    `json:"states"`
	Correct int    `json:"correct"`
}

type UP238BResult struct {
	Schema                       string                `json:"schema"`
	Experiment                   string                `json:"experiment"`
	SourceUP237BSeal             string                `json:"source_up237b_seal"`
	ParentClassification         string                `json:"parent_classification"`
	ParentFarAccuracy            float64               `json:"parent_far_accuracy"`
	TrainingPhases               []int                 `json:"training_phases"`
	ExtremeEvaluationPhases      []int                 `json:"extreme_evaluation_phases"`
	Feature                      string                `json:"feature"`
	ExtremeAccuracy              float64               `json:"extreme_accuracy"`
	Classification               string                `json:"classification"`
	EvaluationLabelFittingUsed   bool                  `json:"evaluation_label_fitting_used"`
	PhaseInputAtInferenceUsed    bool                  `json:"phase_input_at_inference_used"`
	AdaptiveFeatureSelectionUsed bool                  `json:"adaptive_feature_selection_used"`
	NonlinearClassifierUsed      bool                  `json:"nonlinear_classifier_used"`
	LiveActivation               bool                  `json:"live_activation"`
	Summaries                    []UP238BFamilySummary `json:"summaries"`
}

func RunUP238B() (UP238BResult, error) {
	p, err := RunUP237B()
	if err != nil {
		return UP238BResult{}, err
	}
	train := []int{58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72}
	extreme := []int{249, 250, 251, 252, 253, 254, 255, 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280}
	specs := []UP197BComposition{
		{Name: "mixed4", Canonical: []int{0, 5, 1, 6}},
		{Name: "observe4", Canonical: []int{5, 6, 7, 8}},
		{Name: "store4", Canonical: []int{0, 1, 2, 3}},
		{Name: "cross3", Canonical: []int{0, 5, 13}},
	}
	feat := func(ph int, f UP197BComposition) float64 { return up193bFeature(ph, f.Canonical)[2] }
	par := func(ph int) string {
		if ph%2 == 0 {
			return "even"
		}
		return "odd"
	}
	famSum := map[string]float64{}
	famN := map[string]int{}
	type row struct {
		name string
		ph   int
		x    float64
	}
	rows := []row{}
	for _, f := range specs {
		for _, ph := range train {
			x := feat(ph, f)
			rows = append(rows, row{f.Name, ph, x})
			famSum[f.Name] += x
			famN[f.Name]++
		}
	}
	famMean := map[string]float64{}
	for _, f := range specs {
		famMean[f.Name] = famSum[f.Name] / float64(famN[f.Name])
	}
	std := 0.0
	for _, x := range rows {
		d := x.x - famMean[x.name]
		std += d * d
	}
	std = math.Sqrt(std / float64(len(rows)))
	if std == 0 {
		std = 1
	}
	sum := map[string]float64{"even": 0, "odd": 0}
	n := map[string]int{"even": 0, "odd": 0}
	for _, x := range rows {
		q := par(x.ph)
		sum[q] += (x.x - famMean[x.name]) / std
		n[q]++
	}
	cent := map[string]float64{
		"even": sum["even"] / float64(n["even"]),
		"odd":  sum["odd"] / float64(n["odd"]),
	}
	r := UP238BResult{
		Schema:                   UP238BExtremeParitySchema,
		Experiment:               "UP-238B-extreme-far-centered-parity",
		SourceUP237BSeal:         "f8fbe98595b3dd4427fc0b282dfa7a490bd65d7e",
		ParentClassification:     p.Classification,
		ParentFarAccuracy:        p.FarAccuracy,
		TrainingPhases:           train,
		ExtremeEvaluationPhases:  extreme,
		Feature:                  "near_zero_margin_count",
	}
	correct := 0
	for _, f := range specs {
		sm := UP238BFamilySummary{Family: f.Name}
		for _, ph := range extreme {
			z := (feat(ph, f) - famMean[f.Name]) / std
			pred := "even"
			if math.Abs(z-cent["odd"]) < math.Abs(z-cent["even"]) {
				pred = "odd"
			}
			sm.States++
			if pred == par(ph) {
				sm.Correct++
				correct++
			}
		}
		r.Summaries = append(r.Summaries, sm)
	}
	r.ExtremeAccuracy = float64(correct) / 128.0
	if p.Classification != "FAR_TRANSFER_PERFECT" || p.FarAccuracy != 1.0 {
		r.Classification = "ANCHOR_NOT_REPRODUCED"
	} else if r.ExtremeAccuracy == 1.0 {
		r.Classification = "EXTREME_TRANSFER_PERFECT"
	} else {
		r.Classification = "EXTREME_TRANSFER_DECAY"
	}
	return r, nil
}
