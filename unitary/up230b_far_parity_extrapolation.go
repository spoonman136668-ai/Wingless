package unitary

const UP230BFarParitySchema = "wingless.up230b-far-parity-extrapolation.v1"

type UP230BFamilySummary struct {
	Family string `json:"family"`
	HoldoutStates int `json:"holdout_states"`
	SameParityExactMatches int `json:"same_parity_exact_matches"`
	OppositeParityExactMatches int `json:"opposite_parity_exact_matches"`
	UnmatchedHoldoutStates int `json:"unmatched_holdout_states"`
}

type UP230BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP229BSeal string `json:"source_up229b_seal"`
	ReferencePhases []int `json:"reference_phases"`
	HoldoutPhases []int `json:"holdout_phases"`
	Families []string `json:"families"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	ParityUsedForModeling bool `json:"parity_used_for_modeling"`
	GateChanged bool `json:"gate_changed"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP230BFamilySummary `json:"summaries"`
}

func RunUP230B() (UP230BResult,error) {
	holdout:=[]int{89,90,91,92,93,94,95,96,97,98,99,100,101,102,103,104,105,106,107,108,109,110,111,112,113,114,115,116,117,118,119,120}
	known:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	res:=UP230BResult{
		Schema:UP230BFarParitySchema,
		Experiment:"UP-230B-far-parity-extrapolation",
		SourceUP229BSeal:"4b83de60a51ebe02009ade0c7dd5b8b73eadb798",
		ReferencePhases:[]int{58,59},
		HoldoutPhases:holdout,
		HeldoutFittingUsed:false,
		ParityUsedForModeling:false,
		GateChanged:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFeatureUsed:false,
		LiveActivation:false,
	}
	for _,f:=range known{
		res.Families=append(res.Families,f.Name)
		evenRef:=up193bFeature(58,f.Canonical)
		oddRef:=up193bFeature(59,f.Canonical)
		sm:=UP230BFamilySummary{Family:f.Name}
		for _,ph:=range holdout{
			sm.HoldoutStates++
			x:=up193bFeature(ph,f.Canonical)
			same:=evenRef
			opp:=oddRef
			if ph%2!=0{same,opp=oddRef,evenRef}
			if x==same{sm.SameParityExactMatches++}
			if x==opp{sm.OppositeParityExactMatches++}
			if x!=same&&x!=opp{sm.UnmatchedHoldoutStates++}
		}
		res.Summaries=append(res.Summaries,sm)
	}
	return res,nil
}
