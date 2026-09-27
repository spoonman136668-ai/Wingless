package unitary

const UP243CFailurePartitionSchema = "wingless.up243c-initial-state-failure-partition.v1"

type UP243CPoint struct {
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	AOutcome string `json:"a_outcome"`
	BOutcome string `json:"b_outcome"`
	AEventPresent bool `json:"a_event_present"`
	BEventPresent bool `json:"b_event_present"`
}

type UP243CCohortSummary struct {
	Cohort int `json:"cohort"`
	States int `json:"states"`
	AFailures int `json:"a_failures"`
	BFailures int `json:"b_failures"`
	BothFail int `json:"both_fail"`
	NeitherFail int `json:"neither_fail"`
}

type UP243CHandSummary struct {
	InitialHand int `json:"initial_hand"`
	States int `json:"states"`
	AFailures int `json:"a_failures"`
	BFailures int `json:"b_failures"`
	BothFail int `json:"both_fail"`
	NeitherFail int `json:"neither_fail"`
}

type UP243CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP242CSeal string `json:"source_up242c_seal"`
	ConditionA string `json:"condition_a"`
	ConditionB string `json:"condition_b"`
	PairedInitialStates int `json:"paired_initial_states"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
	CohortSummaries []UP243CCohortSummary `json:"cohort_summaries"`
	HandSummaries []UP243CHandSummary `json:"hand_summaries"`
	Points []UP243CPoint `json:"points"`
}

func RunUP243C() (UP243CResult,error) {
	cohorts := [][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}}
	hands := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res := UP243CResult{
		Schema:UP243CFailurePartitionSchema,
		Experiment:"UP-243C-initial-state-failure-partition",
		SourceUP242CSeal:"cd2e47edf7d0b1cfb8b1463f53c7fbf2bd9379dc",
		ConditionA:"schedule2|advance22",
		ConditionB:"schedule3|advance0",
		InterventionChanged:false,
		ThresholdFittingUsed:false,
		ClassifierTrainingUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFieldUsed:false,
		LiveActivation:false,
	}
	cohortAgg:=make([]UP243CCohortSummary,len(cohorts))
	for i:=range cohortAgg {cohortAgg[i].Cohort=i}
	handAgg:=make([]UP243CHandSummary,len(hands))
	for i,h:=range hands {handAgg[i].InitialHand=h}
	for ci,c := range cohorts {
		for _,hand := range hands {
			x,e,ok := up161cTriggerState(hand,c); if !ok {continue}
			a,err := up237cRun(x,e,22,"alternating_shield","fixed_offset_refresh"); if err!=nil{return UP243CResult{},err}
			b,err := up237cRun(x,e,0,"fixed_offset_refresh","alternating_shield"); if err!=nil{return UP243CResult{},err}
			res.PairedInitialStates++
			p:=UP243CPoint{Cohort:ci,InitialHand:hand,EndangeredKey:e,AOutcome:a.Outcome,BOutcome:b.Outcome,AEventPresent:a.EventPresent,BEventPresent:b.EventPresent}
			res.Points=append(res.Points,p)
			af:=a.Outcome=="loss"; bf:=b.Outcome=="loss"
			cohortAgg[ci].States++; handAgg[hand].States++
			if af {cohortAgg[ci].AFailures++;handAgg[hand].AFailures++}
			if bf {cohortAgg[ci].BFailures++;handAgg[hand].BFailures++}
			if af&&bf {cohortAgg[ci].BothFail++;handAgg[hand].BothFail++}
			if !af&&!bf {cohortAgg[ci].NeitherFail++;handAgg[hand].NeitherFail++}
		}
	}
	res.CohortSummaries=cohortAgg
	res.HandSummaries=handAgg
	return res,nil
}
