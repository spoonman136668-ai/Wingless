package unitary

const UP217CCriticalLateSchema="wingless.up217c-critical-tail-late-transfer.v1"

type UP217CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	Cap7Losses int `json:"cap7_losses"`
	CriticalTailLosses int `json:"critical_tail_losses"`
	EighthActions int `json:"eighth_actions"`
	TailRescues int `json:"tail_rescues"`
	UnnecessaryEighth int `json:"unnecessary_eighth"`
	NecessaryEighthAttempts int `json:"necessary_eighth_attempts"`
	IneffectiveEighth int `json:"ineffective_eighth"`
	NoEighthSurvivors int `json:"no_eighth_survivors"`
	RescuePerUnnecessary float64 `json:"rescue_per_unnecessary"`
}
type UP217CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP216CSeal string `json:"source_up216c_seal"`
	CriticalHorizon int `json:"critical_horizon"`
	PhaseAdvances []int `json:"phase_advances"`
	Horizons []int `json:"horizons"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ScheduleRetuningUsed bool `json:"schedule_retuning_used"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	Metrics []UP217CMetric `json:"metrics"`
}
func RunUP217C()(UP217CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{{"no_hostile","no_refresh","hostile_shield"},{"hostile_no","hostile_shield","no_refresh"},{"alternating_fixed","alternating_shield","fixed_offset_refresh"},{"fixed_alternating","fixed_offset_refresh","alternating_shield"}}
	advances:=[]int{8,12};horizons:=[]int{74,76,78,80}
	res:=UP217CResult{Schema:UP217CCriticalLateSchema,Experiment:"UP-217C-critical-tail-late-transfer",SourceUP216CSeal:"14b37dbfb429e4264bb3d279f1ff3488f79efbd9",CriticalHorizon:2,PhaseAdvances:advances,Horizons:horizons,CounterfactualOnly:true,LiveActivation:false,ThresholdFittingUsed:false,ScheduleRetuningUsed:false,FuturePolicyInputUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{
		m:=UP217CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			cap7,_,_:=up214cTreated(x,e,advance,h,s.a,s.b,"none")
			treated,_,eighth:=up216cTreated(x,e,advance,h,s.a,s.b,2)
			capLoss:=cap7<=h;treatedLoss:=treated<=h
			if capLoss{m.Cap7Losses++};if treatedLoss{m.CriticalTailLosses++};m.EighthActions+=eighth
			if capLoss&&!treatedLoss{m.TailRescues++}
			if !capLoss&&eighth>0{m.UnnecessaryEighth++}
			if capLoss&&eighth>0{m.NecessaryEighthAttempts++}
			if capLoss&&treatedLoss&&eighth>0{m.IneffectiveEighth++}
			if !capLoss&&eighth==0{m.NoEighthSurvivors++}
		}}
		if m.UnnecessaryEighth>0{m.RescuePerUnnecessary=float64(m.TailRescues)/float64(m.UnnecessaryEighth)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
