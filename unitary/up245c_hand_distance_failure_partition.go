package unitary

import "sort"

const UP245CHandDistanceSchema = "wingless.up245c-hand-distance-failure-partition.v1"

type UP245CDistanceSummary struct {
	HandDistance int `json:"hand_distance"`
	States int `json:"states"`
	AFailures int `json:"a_failures"`
	BFailures int `json:"b_failures"`
	BothFail int `json:"both_fail"`
	NeitherFail int `json:"neither_fail"`
	AEvents int `json:"a_events"`
	BEvents int `json:"b_events"`
}

type UP245CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP244CSeal string `json:"source_up244c_seal"`
	ConditionA string `json:"condition_a"`
	ConditionB string `json:"condition_b"`
	PairedInitialStates int `json:"paired_initial_states"`
	DistanceDerivedFromNativeState bool `json:"distance_derived_from_native_state"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
	DistanceSummaries []UP245CDistanceSummary `json:"distance_summaries"`
}

func RunUP245C()(UP245CResult,error){
	cohorts:=[][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res:=UP245CResult{
		Schema:UP245CHandDistanceSchema,
		Experiment:"UP-245C-hand-distance-failure-partition",
		SourceUP244CSeal:"abfde6073a84cbaae2fc63c16506a209a619d911",
		ConditionA:"schedule2|advance22",
		ConditionB:"schedule3|advance0",
		DistanceDerivedFromNativeState:true,
		InterventionChanged:false,
		ThresholdFittingUsed:false,
		ClassifierTrainingUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFieldUsed:false,
		LiveActivation:false,
	}
	agg:=map[int]*UP245CDistanceSummary{}
	for _,cohort:=range cohorts{
		for _,initialHand:=range hands{
			x,e,ok:=up161cTriggerState(initialHand,cohort);if !ok{continue}
			slot:=x.find(e)
			if slot<0{continue}
			dist:=(slot-x.hand+16)%16
			a,err:=up237cRun(x,e,22,"alternating_shield","fixed_offset_refresh");if err!=nil{return UP245CResult{},err}
			b,err:=up237cRun(x,e,0,"fixed_offset_refresh","alternating_shield");if err!=nil{return UP245CResult{},err}
			res.PairedInitialStates++
			sm:=agg[dist]
			if sm==nil{sm=&UP245CDistanceSummary{HandDistance:dist};agg[dist]=sm}
			sm.States++
			af:=a.Outcome=="loss";bf:=b.Outcome=="loss"
			if af{sm.AFailures++};if bf{sm.BFailures++}
			if af&&bf{sm.BothFail++}
			if !af&&!bf{sm.NeitherFail++}
			if a.EventPresent{sm.AEvents++};if b.EventPresent{sm.BEvents++}
		}
	}
	ds:=make([]int,0,len(agg))
	for d:=range agg{ds=append(ds,d)}
	sort.Ints(ds)
	for _,d:=range ds{res.DistanceSummaries=append(res.DistanceSummaries,*agg[d])}
	return res,nil
}
