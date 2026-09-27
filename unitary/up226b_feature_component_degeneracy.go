package unitary

import (
	"math"
	"strconv"
)

const UP226BFeatureDegeneracySchema = "wingless.up226b-feature-component-degeneracy.v1"

type UP226BFamilySummary struct {
	Family string `json:"family"`
	CorrectStates int `json:"correct_states"`
	DistinctJointGeometries int `json:"distinct_joint_geometries"`
	DistinctNativeCorrectCount int `json:"distinct_native_correct_count"`
	DistinctMeanAbsoluteMargin int `json:"distinct_mean_absolute_margin"`
	DistinctNearZeroMarginCount int `json:"distinct_near_zero_margin_count"`
	DistinctMinAbsoluteMargin int `json:"distinct_min_absolute_margin"`
}

type UP226BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP225BSeal string `json:"source_up225b_seal"`
	TrainingStates int `json:"training_states"`
	TrainingCorrectStates int `json:"training_correct_states"`
	Features []string `json:"features"`
	EvaluationStatesUsed bool `json:"evaluation_states_used"`
	GateChanged bool `json:"gate_changed"`
	EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	LiveActivation bool `json:"live_activation"`
	FamilySummaries []UP226BFamilySummary `json:"family_summaries"`
}

func up226bFloatKey(v float64) string {
	return strconv.FormatFloat(v, 'g', 17, 64)
}

func up226bJointKey(z [4]float64) string {
	return up226bFloatKey(z[0]) + "|" + up226bFloatKey(z[1]) + "|" + up226bFloatKey(z[2]) + "|" + up226bFloatKey(z[3])
}

func RunUP226B() (UP226BResult, error) {
	train := []int{58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72}
	known := []UP197BComposition{
		{Name: "mixed4", Canonical: []int{0, 5, 1, 6}},
		{Name: "observe4", Canonical: []int{5, 6, 7, 8}},
		{Name: "store4", Canonical: []int{0, 1, 2, 3}},
		{Name: "cross3", Canonical: []int{0, 5, 13}},
	}
	type raw struct {
		family string
		x [4]float64
	}
	raws := make([]raw, 0, len(train)*len(known))
	for _, f := range known {
		for _, ph := range train {
			raws = append(raws, raw{family: f.Name, x: up193bFeature(ph, f.Canonical)})
		}
	}

	var mean, std [4]float64
	for _, rr := range raws {
		for j := 0; j < 4; j++ {
			mean[j] += rr.x[j]
		}
	}
	for j := 0; j < 4; j++ {
		mean[j] /= float64(len(raws))
	}
	for _, rr := range raws {
		for j := 0; j < 4; j++ {
			d := rr.x[j] - mean[j]
			std[j] += d * d
		}
	}
	for j := 0; j < 4; j++ {
		std[j] = math.Sqrt(std[j] / float64(len(raws)))
		if std[j] == 0 {
			std[j] = 1
		}
	}
	zfun := func(x [4]float64) (z [4]float64) {
		for j := 0; j < 4; j++ {
			z[j] = (x[j] - mean[j]) / std[j]
		}
		return
	}

	centroids := make([][4]float64, len(known))
	counts := make([]int, len(known))
	for fi, f := range known {
		for _, rr := range raws {
			if rr.family != f.Name {
				continue
			}
			z := zfun(rr.x)
			for j := 0; j < 4; j++ {
				centroids[fi][j] += z[j]
			}
			counts[fi]++
		}
		for j := 0; j < 4; j++ {
			centroids[fi][j] /= float64(counts[fi])
		}
	}

	type point struct {
		family string
		x [4]float64
		z [4]float64
	}
	correct := make([]point, 0, len(raws))
	for fi, f := range known {
		for _, ph := range train {
			x := up193bFeature(ph, f.Canonical)
			z := zfun(x)
			bestIdx := -1
			best := math.Inf(1)
			for gi, c := range centroids {
				if d := up212bDist(z, c); d < best {
					best = d
					bestIdx = gi
				}
			}
			if bestIdx == fi {
				correct = append(correct, point{family: f.Name, x: x, z: z})
			}
		}
	}

	res := UP226BResult{
		Schema: UP226BFeatureDegeneracySchema,
		Experiment: "UP-226B-feature-component-degeneracy",
		SourceUP225BSeal: "cf0e5ec6fdcbbb718c9383952234fa2beab794c4",
		TrainingStates: len(raws),
		TrainingCorrectStates: len(correct),
		Features: []string{"native_correct_count", "mean_absolute_margin", "near_zero_margin_count", "min_absolute_margin"},
		EvaluationStatesUsed: false,
		GateChanged: false,
		EvaluationDerivedThresholdUsed: false,
		AdaptiveFeatureSelectionUsed: false,
		NewNativeFeatureUsed: false,
		PhaseInputUsed: false,
		LiveActivation: false,
	}

	for _, f := range known {
		joint := map[string]bool{}
		components := [4]map[string]bool{{}, {}, {}, {}}
		n := 0
		for _, p := range correct {
			if p.family != f.Name {
				continue
			}
			n++
			joint[up226bJointKey(p.z)] = true
			for j := 0; j < 4; j++ {
				components[j][up226bFloatKey(p.x[j])] = true
			}
		}
		res.FamilySummaries = append(res.FamilySummaries, UP226BFamilySummary{
			Family: f.Name,
			CorrectStates: n,
			DistinctJointGeometries: len(joint),
			DistinctNativeCorrectCount: len(components[0]),
			DistinctMeanAbsoluteMargin: len(components[1]),
			DistinctNearZeroMarginCount: len(components[2]),
			DistinctMinAbsoluteMargin: len(components[3]),
		})
	}
	return res, nil
}
