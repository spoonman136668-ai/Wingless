package unitary

import "strconv"

const UP228BParityGeometrySchema = "wingless.up228b-parity-geometry-determinism.v1"

type UP228BFamilySummary struct {
	Family string `json:"family"`
	EvenPhases int `json:"even_phases"`
	OddPhases int `json:"odd_phases"`
	DistinctEvenGeometries int `json:"distinct_even_geometries"`
	DistinctOddGeometries int `json:"distinct_odd_geometries"`
	SharedParityGeometries int `json:"shared_parity_geometries"`
}

type UP228BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP227BSeal string `json:"source_up227b_seal"`
	TrainingPhases []int `json:"training_phases"`
	Families []string `json:"families"`
	ParityUsedForModeling bool `json:"parity_used_for_modeling"`
	GateChanged bool `json:"gate_changed"`
	EvaluationStatesUsed bool `json:"evaluation_states_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP228BFamilySummary `json:"summaries"`
}

func up228bKey(x [4]float64) string {
	return strconv.FormatFloat(x[0], 'g', 17, 64) + "|" +
		strconv.FormatFloat(x[1], 'g', 17, 64) + "|" +
		strconv.FormatFloat(x[2], 'g', 17, 64) + "|" +
		strconv.FormatFloat(x[3], 'g', 17, 64)
}

func RunUP228B() (UP228BResult, error) {
	phases := []int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	known := []UP197BComposition{
		{Name:"mixed4", Canonical:[]int{0,5,1,6}},
		{Name:"observe4", Canonical:[]int{5,6,7,8}},
		{Name:"store4", Canonical:[]int{0,1,2,3}},
		{Name:"cross3", Canonical:[]int{0,5,13}},
	}
	res := UP228BResult{
		Schema:UP228BParityGeometrySchema,
		Experiment:"UP-228B-parity-geometry-determinism",
		SourceUP227BSeal:"2d1d85bf73e09f0497ac21128d5d0cf981485d51",
		TrainingPhases:phases,
		ParityUsedForModeling:false,
		GateChanged:false,
		EvaluationStatesUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFeatureUsed:false,
		LiveActivation:false,
	}
	for _,f := range known {
		res.Families = append(res.Families,f.Name)
		even := map[string]bool{}
		odd := map[string]bool{}
		sm := UP228BFamilySummary{Family:f.Name}
		for _,ph := range phases {
			k := up228bKey(up193bFeature(ph,f.Canonical))
			if ph%2==0 {
				sm.EvenPhases++
				even[k]=true
			} else {
				sm.OddPhases++
				odd[k]=true
			}
		}
		sm.DistinctEvenGeometries=len(even)
		sm.DistinctOddGeometries=len(odd)
		for k:=range even {if odd[k]{sm.SharedParityGeometries++}}
		res.Summaries=append(res.Summaries,sm)
	}
	return res,nil
}
