package unitary

import "math"

const UP224BManifoldDegeneracySchema = "wingless.up224b-training-manifold-degeneracy.v1"

type UP224BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP223BSeal string `json:"source_up223b_seal"`
	TrainingStates int `json:"training_states"`
	TrainingCorrectStates int `json:"training_correct_states"`
	UniqueCorrectStates int `json:"unique_correct_states"`
	DuplicateCorrectStates int `json:"duplicate_correct_states"`
	MaxDuplicateMultiplicity int `json:"max_duplicate_multiplicity"`
	OriginalManifoldThreshold float64 `json:"original_manifold_threshold"`
	UniqueNearestMin float64 `json:"unique_nearest_min"`
	UniqueNearestMax float64 `json:"unique_nearest_max"`
	ZeroDistancePairs int `json:"zero_distance_pairs"`
	EvaluationStatesUsed bool `json:"evaluation_states_used"`
	GateChanged bool `json:"gate_changed"`
	EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	LiveActivation bool `json:"live_activation"`
}

func up224bSame(a, b [4]float64) bool {
	return up223bVDist(a, b) <= 1e-12
}

func RunUP224B() (UP224BResult, error) {
	train := []int{58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72}
	known := []UP197BComposition{
		{Name: "mixed4", Canonical: []int{0, 5, 1, 6}},
		{Name: "observe4", Canonical: []int{5, 6, 7, 8}},
		{Name: "store4", Canonical: []int{0, 1, 2, 3}},
		{Name: "cross3", Canonical: []int{0, 5, 13}},
	}
	type raw struct {
		name string
		x [4]float64
	}
	raws := []raw{}
	for _, f := range known {
		for _, ph := range train {
			raws = append(raws, raw{name: f.Name, x: up193bFeature(ph, f.Canonical)})
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
			if rr.name != f.Name {
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
	correct := make([][4]float64, 0, len(raws))
	for fi, f := range known {
		for _, ph := range train {
			z := zfun(up193bFeature(ph, f.Canonical))
			bestIdx := -1
			best := math.Inf(1)
			for gi, c := range centroids {
				if d := up212bDist(z, c); d < best {
					best = d
					bestIdx = gi
				}
			}
			if bestIdx == fi {
				correct = append(correct, z)
			}
		}
	}
	originalT := 0.0
	zeroPairs := 0
	for i, z := range correct {
		near := math.Inf(1)
		for j, q := range correct {
			if i == j {
				continue
			}
			d := up223bVDist(z, q)
			if d <= 1e-12 {
				zeroPairs++
			}
			if d < near {
				near = d
			}
		}
		if !math.IsInf(near, 1) && near > originalT {
			originalT = near
		}
	}
	zeroPairs /= 2
	unique := make([][4]float64, 0, len(correct))
	multiplicity := []int{}
	for _, z := range correct {
		found := -1
		for i, q := range unique {
			if up224bSame(z, q) {
				found = i
				break
			}
		}
		if found >= 0 {
			multiplicity[found]++
		} else {
			unique = append(unique, z)
			multiplicity = append(multiplicity, 1)
		}
	}
	maxMult := 0
	for _, n := range multiplicity {
		if n > maxMult {
			maxMult = n
		}
	}
	uniqueMin := 0.0
	uniqueMax := 0.0
	if len(unique) > 1 {
		uniqueMin = math.Inf(1)
		for i, z := range unique {
			near := math.Inf(1)
			for j, q := range unique {
				if i == j {
					continue
				}
				if d := up223bVDist(z, q); d < near {
					near = d
				}
			}
			if near < uniqueMin {
				uniqueMin = near
			}
			if near > uniqueMax {
				uniqueMax = near
			}
		}
	}
	return UP224BResult{
		Schema: UP224BManifoldDegeneracySchema,
		Experiment: "UP-224B-training-manifold-degeneracy",
		SourceUP223BSeal: "f8aa129045967e4b40629df2db02541525d92134",
		TrainingStates: len(raws),
		TrainingCorrectStates: len(correct),
		UniqueCorrectStates: len(unique),
		DuplicateCorrectStates: len(correct) - len(unique),
		MaxDuplicateMultiplicity: maxMult,
		OriginalManifoldThreshold: originalT,
		UniqueNearestMin: uniqueMin,
		UniqueNearestMax: uniqueMax,
		ZeroDistancePairs: zeroPairs,
		EvaluationStatesUsed: false,
		GateChanged: false,
		EvaluationDerivedThresholdUsed: false,
		NewNativeFeatureUsed: false,
		LiveActivation: false,
	}, nil
}
