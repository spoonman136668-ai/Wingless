package unitary

const UP214CNativeTailSchema="wingless.up214c-native-near-tail-trigger.v1"

type UP214CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Mode string `json:"mode"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	EighthActions int `json:"eighth_actions"`
}
type UP214CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP213CSeal string `json:"source_up213c_seal"`
	Modes []string `json:"modes"`
	Horizons []int `json:"horizons"`
	NearRiskHorizon int `json:"near_risk_horizon"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	Metrics []UP214CMetric `json:"metrics"`
}
func up214cTreated(x0 *up81cAging,endangered,advance,horizon int,a,b,mode string)(loss,actions,eighth int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{m.query(endangered);actions++;since=0}
			if actions==7{
				fire:=false
				if mode=="nominal"&&since>=4{fire=true}
				if mode=="native_near"&&up161cAdversarial(&m,endangered)<=4{fire=true}
				if fire{m.query(endangered);actions++;eighth++;since=0}
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1039000+step,step,up210cPolicy(step,advance,a,b)){return step,actions,eighth}
		}
	}
	return horizon+1,actions,eighth
}
func RunUP214C()(UP214CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{{"no_hostile","no_refresh","hostile_shield"},{"hostile_no","hostile_shield","no_refresh"},{"alternating_fixed","alternating_shield","fixed_offset_refresh"},{"fixed_alternating","fixed_offset_refresh","alternating_shield"}}
	advances:=[]int{0,4};horizons:=[]int{66,68,70,72};modes:=[]string{"none","nominal","native_near"}
	res:=UP214CResult{Schema:UP214CNativeTailSchema,Experiment:"UP-214C-native-near-tail-trigger",SourceUP213CSeal:"3f042545dd0e2077aa54e1a6be06a8c5fee8b8cc",Modes:modes,Horizons:horizons,NearRiskHorizon:4,CounterfactualOnly:true,LiveActivation:false,ThresholdFittingUsed:false,FuturePolicyInputUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{for _,mode:=range modes{
		m:=UP214CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h,Mode:mode}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
			base:=up211cBaseline(x,e,advance,h,s.a,s.b)
			treated,actions,eighth:=up214cTreated(x,e,advance,h,s.a,s.b,mode)
			if base<=h{m.BaselineLosses++};if treated<=h{m.TreatedLosses++};if base<=h&&treated>h{m.PreventedLosses++};if treated<base{m.AcceleratedLosses++}
			m.ActionsTaken+=actions;m.EighthActions+=eighth
		}}
		res.Metrics=append(res.Metrics,m)
	}}}}
	return res,nil
}
