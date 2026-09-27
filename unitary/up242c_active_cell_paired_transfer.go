package unitary

const UP242CActiveCellPairedSchema = "wingless.up242c-active-cell-paired-transfer.v1"

type UP242CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP241CSeal string `json:"source_up241c_seal"`
	ConditionA string `json:"condition_a"`
	ConditionB string `json:"condition_b"`
	PairedInitialStates int `json:"paired_initial_states"`
	AFailures int `json:"a_failures"`
	BFailures int `json:"b_failures"`
	BothFail int `json:"both_fail"`
	AOnlyFail int `json:"a_only_fail"`
	BOnlyFail int `json:"b_only_fail"`
	BothSurvive int `json:"both_survive"`
	BothEventPresent int `json:"both_event_present"`
	AOnlyEventPresent int `json:"a_only_event_present"`
	BOnlyEventPresent int `json:"b_only_event_present"`
	NeitherEventPresent int `json:"neither_event_present"`
	ExactCurrentSignatureMatches int `json:"exact_current_signature_matches"`
	ExactTrajectorySignatureMatches int `json:"exact_trajectory_signature_matches"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
}

func RunUP242C() (UP242CResult,error) {
	cohorts := [][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}}
	hands := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res := UP242CResult{
		Schema:UP242CActiveCellPairedSchema,
		Experiment:"UP-242C-active-cell-paired-transfer",
		SourceUP241CSeal:"095315dcec5859806ad5e7c89276447b435f1581",
		ConditionA:"schedule2|advance22",
		ConditionB:"schedule3|advance0",
		InterventionChanged:false,
		ThresholdFittingUsed:false,
		ClassifierTrainingUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFieldUsed:false,
		LiveActivation:false,
	}
	for _,c := range cohorts {
		for _,hand := range hands {
			x,e,ok := up161cTriggerState(hand,c); if !ok {continue}
			a,err := up237cRun(x,e,22,"alternating_shield","fixed_offset_refresh"); if err!=nil{return UP242CResult{},err}
			b,err := up237cRun(x,e,0,"fixed_offset_refresh","alternating_shield"); if err!=nil{return UP242CResult{},err}
			res.PairedInitialStates++
			af:=a.Outcome=="loss"; bf:=b.Outcome=="loss"
			if af {res.AFailures++}; if bf {res.BFailures++}
			switch {
			case af&&bf: res.BothFail++
			case af&&!bf: res.AOnlyFail++
			case !af&&bf: res.BOnlyFail++
			default: res.BothSurvive++
			}
			switch {
			case a.EventPresent&&b.EventPresent:
				res.BothEventPresent++
				if a.CurrentSignature!="" && a.CurrentSignature==b.CurrentSignature {res.ExactCurrentSignatureMatches++}
				if a.TrajectorySignature!="" && a.TrajectorySignature==b.TrajectorySignature {res.ExactTrajectorySignatureMatches++}
			case a.EventPresent&&!b.EventPresent: res.AOnlyEventPresent++
			case !a.EventPresent&&b.EventPresent: res.BOnlyEventPresent++
			default: res.NeitherEventPresent++
			}
		}
	}
	return res,nil
}
