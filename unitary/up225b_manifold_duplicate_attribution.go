package unitary

import "math"

const UP225BDuplicateAttributionSchema = "wingless.up225b-manifold-duplicate-attribution.v1"

type UP225BFamilySummary struct {
	Family string `json:"family"`
	CorrectStates int `json:"correct_states"`
	DistinctGeometries int `json:"distinct_geometries"`
	DuplicateStates int `json:"duplicate_states"`
}

type UP225BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP224BSeal string `json:"source_up224b_seal"`
	TrainingStates int `json:"training_states"`
	TrainingCorrectStates int `json:"training_correct_states"`
	UniqueCorrectGeometries int `json:"unique_correct_geometries"`
	DuplicateCorrectStates int `json:"duplicate_correct_states"`
	DuplicateGeometryGroups int `json:"duplicate_geometry_groups"`
	CrossFamilyDuplicateGroups int `json:"cross_family_duplicate_groups"`
	WithinFamilyDuplicateGroups int `json:"within_family_duplicate_groups"`
	MaxDuplicateMultiplicity int `json:"max_duplicate_multiplicity"`
	FamilySummaries []UP225BFamilySummary `json:"family_summaries"`
	EvaluationStatesUsed bool `json:"evaluation_states_used"`
	GateChanged bool `json:"gate_changed"`
	EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	LiveActivation bool `json:"live_activation"`
}

type up225bPoint struct {
	family string
	z [4]float64
}

type up225bGroup struct {
	z [4]float64
	count int
	families map[string]bool
}

func RunUP225B() (UP225BResult, error) {
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

	correct := make([]up225bPoint, 0, len(raws))
	familyCorrect := map[string]int{}
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
				correct = append(correct, up225bPoint{family: f.Name, z: z})
				familyCorrect[f.Name]++
			}
		}
	}

	groups := make([]up225bGroup, 0, len(correct))
	pointGroup := make([]int, 0, len(correct))
	for _, p := range correct {
		found := -1
		for i := range groups {
			if up224bSame(p.z, groups[i].z) {
				found = i
				break
			}
		}
		if found < 0 {
			groups = append(groups, up225bGroup{z: p.z, families: map[string]bool{}})
			found = len(groups) - 1
		}
		groups[found].count++
		groups[found].families[p.family] = true
		pointGroup = append(pointGroup, found)
	}

	duplicateGroups := 0
	crossFamily := 0
	withinFamily := 0
	maxMultiplicity := 0
	for _, g := range groups {
		if g.count > maxMultiplicity {
			maxMultiplicity = g.count
		}
		if g.count <= 1 {
			continue
		}
		duplicateGroups++
		if len(g.families) > 1 {
			crossFamily++
		} else {
			withinFamily++
		}
	}

	res := UP225BResult{
		Schema: UP225BDuplicateAttributionSchema,
		Experiment: "UP-225B-manifold-duplicate-attribution",
		SourceUP224BSeal: "7b25c8e574fcbb957d3045e147598576a33954f4",
		TrainingStates: len(raws),
		TrainingCorrectStates: len(correct),
		UniqueCorrectGeometries: len(groups),
		DuplicateCorrectStates: len(correct) - len(groups),
		DuplicateGeometryGroups: duplicateGroups,
		CrossFamilyDuplicateGroups: crossFamily,
		WithinFamilyDuplicateGroups: withinFamily,
		MaxDuplicateMultiplicity: maxMultiplicity,
		EvaluationStatesUsed: false,
		GateChanged: false,
		EvaluationDerivedThresholdUsed: false,
		NewNativeFeatureUsed: false,
		LiveActivation: false,
	}
	for _, f := range known {
		seen := map[int]bool{}
		for i, p := range correct {
			if p.family == f.Name {
				seen[pointGroup[i]] = true
			}
		}
		res.FamilySummaries = append(res.FamilySummaries, UP225BFamilySummary{
			Family: f.Name,
			CorrectStates: familyCorrect[f.Name],
			DistinctGeometries: len(seen),
			DuplicateStates: familyCorrect[f.Name] - len(seen),
		})
	}
	return res, nil
}
