package unitary

const UP227BPhasePeriodicitySchema = "wingless.up227b-phase-periodicity-diagnostic.v1"

type UP227BPeriodSummary struct {
	Family string `json:"family"`
	Period int `json:"period"`
	Comparisons int `json:"comparisons"`
	ExactGeometryMatches int `json:"exact_geometry_matches"`
}

type UP227BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP226BSeal string `json:"source_up226b_seal"`
	TrainingPhases []int `json:"training_phases"`
	Families []string `json:"families"`
	CandidatePeriods []int `json:"candidate_periods"`
	EvaluationStatesUsed bool `json:"evaluation_states_used"`
	GateChanged bool `json:"gate_changed"`
	EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`
	AdaptivePeriodSelectionUsed bool `json:"adaptive_period_selection_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP227BPeriodSummary `json:"summaries"`
}

func RunUP227B() (UP227BResult, error) {
	phases := []int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	known := []UP197BComposition{
		{Name:"mixed4", Canonical:[]int{0,5,1,6}},
		{Name:"observe4", Canonical:[]int{5,6,7,8}},
		{Name:"store4", Canonical:[]int{0,1,2,3}},
		{Name:"cross3", Canonical:[]int{0,5,13}},
	}
	periods := []int{1,2,3,4,5,6,7,8}
	res := UP227BResult{
		Schema:UP227BPhasePeriodicitySchema,
		Experiment:"UP-227B-phase-periodicity-diagnostic",
		SourceUP226BSeal:"1b459506b88827a70a3aa3a1341dc127635d53eb",
		TrainingPhases:phases,
		CandidatePeriods:periods,
		EvaluationStatesUsed:false,
		GateChanged:false,
		EvaluationDerivedThresholdUsed:false,
		AdaptivePeriodSelectionUsed:false,
		NewNativeFeatureUsed:false,
		LiveActivation:false,
	}
	for _,f := range known {
		res.Families = append(res.Families, f.Name)
		vals := map[int][4]float64{}
		for _,ph := range phases {
			vals[ph] = up193bFeature(ph, f.Canonical)
		}
		for _,period := range periods {
			s := UP227BPeriodSummary{Family:f.Name,Period:period}
			for _,ph := range phases {
				v2,ok := vals[ph+period]
				if !ok { continue }
				s.Comparisons++
				if vals[ph] == v2 { s.ExactGeometryMatches++ }
			}
			res.Summaries = append(res.Summaries,s)
		}
	}
	return res,nil
}
