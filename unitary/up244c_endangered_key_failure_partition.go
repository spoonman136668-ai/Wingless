package unitary

import "sort"

const UP244CEndangeredKeySchema = "wingless.up244c-endangered-key-failure-partition.v1"

type UP244CKeySummary struct {
	EndangeredKey int `json:"endangered_key"`
	States int `json:"states"`
	AFailures int `json:"a_failures"`
	BFailures int `json:"b_failures"`
	BothFail int `json:"both_fail"`
	NeitherFail int `json:"neither_fail"`
	AEvents int `json:"a_events"`
	BEvents int `json:"b_events"`
}

type UP244CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP243CSeal string `json:"source_up243c_seal"`
	ConditionA string `json:"condition_a"`
	ConditionB string `json:"condition_b"`
	PairedInitialStates int `json:"paired_initial_states"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
	KeySummaries []UP244CKeySummary `json:"key_summaries"`
}

func RunUP244C()(UP244CResult,error) {
	cohorts:=[][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res:=UP244CResult{
		Schema:UP244CEndangeredKeySchema,
		Experiment:"UP-244C-endangered-key-failure-partition",
		SourceUP243CSeal:"ecb6f850178fc6e866f1b63a0d32cd49a85c0807",
		ConditionA:"schedule2|advance22",
		ConditionB:"schedule3|advance0",
		InterventionChanged:false,
		ThresholdFittingUsed:false,
		ClassifierTrainingUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFieldUsed:false,
		LiveActivation:false,
	}
	agg:=map[int]*UP244CKeySummary{}
	for _,c:=range cohorts {
		for _,hand:=range hands {
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
			a,err:=up237cRun(x,e,22,"alternating_shield","fixed_offset_refresh");if err!=nil{return UP244CResult{},err}
			b,err:=up237cRun(x,e,0,"fixed_offset_refresh","alternating_shield");if err!=nil{return UP244CResult{},err}
			res.PairedInitialStates++
			sm:=agg[e]
			if sm==nil {sm=&UP244CKeySummary{EndangeredKey:e};agg[e]=sm}
			sm.States++
			af:=a.Outcome=="loss";bf:=b.Outcome=="loss"
			if af{sm.AFailures++};if bf{sm.BFailures++}
			if af&&bf{sm.BothFail++}
			if !af&&!bf{sm.NeitherFail++}
			if a.EventPresent{sm.AEvents++};if b.EventPresent{sm.BEvents++}
		}
	}
	keys:=make([]int,0,len(agg))
	for k:=range agg{keys=append(keys,k)}
	sort.Ints(keys)
	for _,k:=range keys{res.KeySummaries=append(res.KeySummaries,*agg[k])}
	return res,nil
}
