package unitary

const UP229BParityHoldoutSchema = "wingless.up229b-parity-holdout-generalization.v1"

type UP229BFamilySummary struct {
	Family string `json:"family"`
	HoldoutStates int `json:"holdout_states"`
	SameParityExactMatches int `json:"same_parity_exact_matches"`
	OppositeParityExactMatches int `json:"opposite_parity_exact_matches"`
	UnmatchedHoldoutStates int `json:"unmatched_holdout_states"`
}

type UP229BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP228BSeal string `json:"source_up228b_seal"`
	TrainingPhases []int `json:"training_phases"`
	HoldoutPhases []int `json:"holdout_phases"`
	Families []string `json:"families"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	ParityUsedForModeling bool `json:"parity_used_for_modeling"`
	GateChanged bool `json:"gate_changed"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP229BFamilySummary `json:"summaries"`
}

func RunUP229B() (UP229BResult,error) {
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	holdout:=[]int{73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88}
	known:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	res:=UP229BResult{
		Schema:UP229BParityHoldoutSchema,
		Experiment:"UP-229B-parity-holdout-generalization",
		SourceUP228BSeal:"64e5835019c182e19192b89b37b8f9357f340f41",
		TrainingPhases:train,
		HoldoutPhases:holdout,
		HeldoutFittingUsed:false,
		ParityUsedForModeling:false,
		GateChanged:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFeatureUsed:false,
		LiveActivation:false,
	}
	for _,f:=range known {
		res.Families=append(res.Families,f.Name)
		evenRef:=up193bFeature(58,f.Canonical)
		oddRef:=up193bFeature(59,f.Canonical)
		sm:=UP229BFamilySummary{Family:f.Name}
		for _,ph:=range holdout {
			sm.HoldoutStates++
			x:=up193bFeature(ph,f.Canonical)
			same:=evenRef
			opp:=oddRef
			if ph%2!=0 {same,opp=oddRef,evenRef}
			if x==same {sm.SameParityExactMatches++}
			if x==opp {sm.OppositeParityExactMatches++}
			if x!=same && x!=opp {sm.UnmatchedHoldoutStates++}
		}
		res.Summaries=append(res.Summaries,sm)
	}
	return res,nil
}
