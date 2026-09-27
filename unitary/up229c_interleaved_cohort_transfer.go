package unitary

const UP229CInterleavedSchema="wingless.up229c-interleaved-cohort-transfer.v1"

type UP229CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	CombinedLosses int `json:"combined_losses"`
	BaselineActions int `json:"baseline_actions"`
	CombinedActions int `json:"combined_actions"`
	AdvanceEvents int `json:"advance_events"`
	AdvanceWrites int `json:"advance_writes"`
	NinthActions int `json:"ninth_actions"`
	RescuedBaselineFailures int `json:"rescued_baseline_failures"`
	InterventionBearingBaselineSurvivors int `json:"intervention_bearing_baseline_survivors"`
}
type UP229CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP228CSeal string `json:"source_up228c_seal"`
	Cohorts [][]int `json:"cohorts"`
	PolicyPairs []string `json:"policy_pairs"`
	PhaseAdvances []int `json:"phase_advances"`
	Horizons []int `json:"horizons"`
	MonitoringCadence int `json:"monitoring_cadence"`
	MaxActions int `json:"max_actions"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	CohortSpecificTuningUsed bool `json:"cohort_specific_tuning_used"`
	NewTriggerUsed bool `json:"new_trigger_used"`
	NewActionTypeUsed bool `json:"new_action_type_used"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ScheduleCompressionUsed bool `json:"schedule_compression_used"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	LiveActivation bool `json:"live_activation"`
	Metrics []UP229CMetric `json:"metrics"`
}
func RunUP229C()(UP229CResult,error){
	cohorts:=[][]int{{0,5,10,15},{1,6,11,12},{2,7,8,13},{3,4,9,14}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{8,12};horizons:=[]int{74,80}
	res:=UP229CResult{
		Schema:UP229CInterleavedSchema,Experiment:"UP-229C-interleaved-cohort-transfer",
		SourceUP228CSeal:"f6e4274d318708d735e7e5544a79d4e8dbb26b19",
		Cohorts:cohorts,PolicyPairs:[]string{"no_hostile","hostile_no","alternating_fixed","fixed_alternating"},
		PhaseAdvances:advances,Horizons:horizons,MonitoringCadence:2,MaxActions:9,
		CounterfactualOnly:true,CohortSpecificTuningUsed:false,NewTriggerUsed:false,
		NewActionTypeUsed:false,ThresholdFittingUsed:false,ScheduleCompressionUsed:false,
		FuturePolicyInputUsed:false,LiveActivation:false,
	}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{
		m:=UP229CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			baseLoss,baseActions,_,_,_:=up228cRun(x,e,advance,h,s.a,s.b,false,false)
			combinedLoss,combinedActions,advanced,delta,ninth:=up228cRun(x,e,advance,h,s.a,s.b,true,true)
			baseFailed:=baseLoss<=h;combinedFailed:=combinedLoss<=h
			if baseFailed{m.BaselineLosses++};if combinedFailed{m.CombinedLosses++}
			m.BaselineActions+=baseActions;m.CombinedActions+=combinedActions
			if advanced{m.AdvanceEvents++;m.AdvanceWrites+=delta}
			m.NinthActions+=ninth
			if baseFailed&&!combinedFailed{m.RescuedBaselineFailures++}
			if !baseFailed&&(advanced||ninth>0){m.InterventionBearingBaselineSurvivors++}
		}}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
