package unitary

const UP241CEventSupportSchema = "wingless.up241c-event-support-coverage.v1"

type UP241CSummary struct {
	Schedule int `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Arms int `json:"arms"`
	Failures int `json:"failures"`
	Survivors int `json:"survivors"`
	EventArms int `json:"event_arms"`
	EventFailures int `json:"event_failures"`
	EventSurvivors int `json:"event_survivors"`
	SilentFailures int `json:"silent_failures"`
}

type UP241CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP240CSeal string `json:"source_up240c_seal"`
	ArmsTotal int `json:"arms_total"`
	Schedules int `json:"schedules"`
	AdvanceConditions []int `json:"advance_conditions"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP241CSummary `json:"summaries"`
}

func RunUP241C() (UP241CResult,error) {
	cohorts := [][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}}
	hands := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{a,b string}
	schedules := []sched{
		{"no_refresh","hostile_shield"},
		{"hostile_shield","no_refresh"},
		{"alternating_shield","fixed_offset_refresh"},
		{"fixed_offset_refresh","alternating_shield"},
	}
	advances := []int{0,22}
	res := UP241CResult{
		Schema:UP241CEventSupportSchema,
		Experiment:"UP-241C-event-support-coverage",
		SourceUP240CSeal:"428a28addb7b9be8de9c55ccc327e3c168babae9",
		Schedules:len(schedules),
		AdvanceConditions:advances,
		InterventionChanged:false,
		ThresholdFittingUsed:false,
		ClassifierTrainingUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFieldUsed:false,
		LiveActivation:false,
	}
	for si,s := range schedules {
		for _,adv := range advances {
			sm := UP241CSummary{Schedule:si,PhaseAdvanceWrites:adv}
			for _,c := range cohorts {
				for _,hand := range hands {
					x,e,ok := up161cTriggerState(hand,c); if !ok {continue}
					a,err := up237cRun(x,e,adv,s.a,s.b); if err != nil {return UP241CResult{},err}
					sm.Arms++; res.ArmsTotal++
					if a.Outcome=="loss" {sm.Failures++} else {sm.Survivors++}
					if a.EventPresent {
						sm.EventArms++
						if a.Outcome=="loss" {sm.EventFailures++} else {sm.EventSurvivors++}
					} else if a.Outcome=="loss" {
						sm.SilentFailures++
					}
				}
			}
			res.Summaries = append(res.Summaries,sm)
		}
	}
	return res,nil
}
