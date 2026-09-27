package unitary

const UP215CTailSelectivitySchema="wingless.up215c-native-tail-selectivity.v1"

type UP215CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	Cap7Losses int `json:"cap7_losses"`
	NativeTailLosses int `json:"native_tail_losses"`
	EighthActions int `json:"eighth_actions"`
	TailRescues int `json:"tail_rescues"`
	UnnecessaryEighth int `json:"unnecessary_eighth"`
	NecessaryEighthAttempts int `json:"necessary_eighth_attempts"`
	IneffectiveEighth int `json:"ineffective_eighth"`
	NoEighthSurvivors int `json:"no_eighth_survivors"`
	RescuePerUnnecessary float64 `json:"rescue_per_unnecessary"`
}
type UP215CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP214CSeal string `json:"source_up214c_seal"`
	NearRiskHorizon int `json:"near_risk_horizon"`
	Horizons []int `json:"horizons"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	Metrics []UP215CMetric `json:"metrics"`
}
func RunUP215C()(UP215CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{{"no_hostile","no_refresh","hostile_shield"},{"hostile_no","hostile_shield","no_refresh"},{"alternating_fixed","alternating_shield","fixed_offset_refresh"},{"fixed_alternating","fixed_offset_refresh","alternating_shield"}}
	advances:=[]int{0,4};horizons:=[]int{66,68,70,72}
	res:=UP215CResult{Schema:UP215CTailSelectivitySchema,Experiment:"UP-215C-native-tail-selectivity",SourceUP214CSeal:"1bbd8fc636474781a998a65015865e2be9ed7d28",NearRiskHorizon:4,Horizons:horizons,CounterfactualOnly:true,LiveActivation:false,ThresholdFittingUsed:false,FuturePolicyInputUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{
		m:=UP215CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			cap7,_,_:=up214cTreated(x,e,advance,h,s.a,s.b,"none")
			native,_,eighth:=up214cTreated(x,e,advance,h,s.a,s.b,"native_near")
			cap7Loss:=cap7<=h;nativeLoss:=native<=h
			if cap7Loss{m.Cap7Losses++};if nativeLoss{m.NativeTailLosses++};m.EighthActions+=eighth
			if cap7Loss&&!nativeLoss{m.TailRescues++}
			if !cap7Loss&&eighth>0{m.UnnecessaryEighth++}
			if cap7Loss&&eighth>0{m.NecessaryEighthAttempts++}
			if cap7Loss&&nativeLoss&&eighth>0{m.IneffectiveEighth++}
			if !cap7Loss&&eighth==0{m.NoEighthSurvivors++}
		}}
		if m.UnnecessaryEighth>0{m.RescuePerUnnecessary=float64(m.TailRescues)/float64(m.UnnecessaryEighth)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
