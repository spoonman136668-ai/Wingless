package unitary

const UP216CCriticalNearSchema="wingless.up216c-critical-vs-near-tail-trigger.v1"

type UP216CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Trigger string `json:"trigger"`
	Arms int `json:"arms"`
	Cap7Losses int `json:"cap7_losses"`
	TreatedLosses int `json:"treated_losses"`
	EighthActions int `json:"eighth_actions"`
	TailRescues int `json:"tail_rescues"`
	UnnecessaryEighth int `json:"unnecessary_eighth"`
	NecessaryEighthAttempts int `json:"necessary_eighth_attempts"`
	IneffectiveEighth int `json:"ineffective_eighth"`
	NoEighthSurvivors int `json:"no_eighth_survivors"`
	RescuePerUnnecessary float64 `json:"rescue_per_unnecessary"`
}
type UP216CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP215CSeal string `json:"source_up215c_seal"`
	Triggers []string `json:"triggers"`
	CriticalHorizon int `json:"critical_horizon"`
	NearHorizon int `json:"near_horizon"`
	Horizons []int `json:"horizons"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	Metrics []UP216CMetric `json:"metrics"`
}
func up216cTreated(x0 *up81cAging,endangered,advance,horizon int,a,b string,threshold int)(loss,actions,eighth int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{m.query(endangered);actions++;since=0}
			if actions==7&&up161cAdversarial(&m,endangered)<=threshold{
				m.query(endangered);actions++;eighth++;since=0
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1040000+step,step,up210cPolicy(step,advance,a,b)){return step,actions,eighth}
		}
	}
	return horizon+1,actions,eighth
}
func RunUP216C()(UP216CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{{"no_hostile","no_refresh","hostile_shield"},{"hostile_no","hostile_shield","no_refresh"},{"alternating_fixed","alternating_shield","fixed_offset_refresh"},{"fixed_alternating","fixed_offset_refresh","alternating_shield"}}
	advances:=[]int{0,4};horizons:=[]int{66,68,70,72};triggers:=[]string{"critical","near"}
	res:=UP216CResult{Schema:UP216CCriticalNearSchema,Experiment:"UP-216C-critical-vs-near-tail-trigger",SourceUP215CSeal:"a4a74a4ff3c45f76466236c2350bb3b6f5b3ef1e",Triggers:triggers,CriticalHorizon:2,NearHorizon:4,Horizons:horizons,CounterfactualOnly:true,LiveActivation:false,ThresholdFittingUsed:false,FuturePolicyInputUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{for _,trigger:=range triggers{
		threshold:=2;if trigger=="near"{threshold=4}
		m:=UP216CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h,Trigger:trigger}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			cap7,_,_:=up214cTreated(x,e,advance,h,s.a,s.b,"none")
			treated,_,eighth:=up216cTreated(x,e,advance,h,s.a,s.b,threshold)
			capLoss:=cap7<=h;treatedLoss:=treated<=h
			if capLoss{m.Cap7Losses++};if treatedLoss{m.TreatedLosses++};m.EighthActions+=eighth
			if capLoss&&!treatedLoss{m.TailRescues++}
			if !capLoss&&eighth>0{m.UnnecessaryEighth++}
			if capLoss&&eighth>0{m.NecessaryEighthAttempts++}
			if capLoss&&treatedLoss&&eighth>0{m.IneffectiveEighth++}
			if !capLoss&&eighth==0{m.NoEighthSurvivors++}
		}}
		if m.UnnecessaryEighth>0{m.RescuePerUnnecessary=float64(m.TailRescues)/float64(m.UnnecessaryEighth)}
		res.Metrics=append(res.Metrics,m)
	}}}}
	return res,nil
}
