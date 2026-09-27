package unitary

const UP218CNinthSchema="wingless.up218c-fresh-critical-ninth-action.v1"

type UP218CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	Cap8Losses int `json:"cap8_losses"`
	Cap9Losses int `json:"cap9_losses"`
	NinthActions int `json:"ninth_actions"`
	NinthRescues int `json:"ninth_rescues"`
	UnnecessaryNinth int `json:"unnecessary_ninth"`
	NecessaryNinthAttempts int `json:"necessary_ninth_attempts"`
	IneffectiveNinth int `json:"ineffective_ninth"`
	NoNinthSurvivors int `json:"no_ninth_survivors"`
	RescuePerUnnecessary float64 `json:"rescue_per_unnecessary"`
}
type UP218CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP217CSeal string `json:"source_up217c_seal"`
	CriticalHorizon int `json:"critical_horizon"`
	MaxActions int `json:"max_actions"`
	PhaseAdvances []int `json:"phase_advances"`
	Horizons []int `json:"horizons"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ScheduleRetuningUsed bool `json:"schedule_retuning_used"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	Metrics []UP218CMetric `json:"metrics"`
}
func up218cTreated(x0 *up81cAging,endangered,advance,horizon int,a,b string,maxActions int)(loss,actions,ninth int){
	m:=*x0;committed:=false;since:=0;eighthFired:=false;eighthBoundary:=-1
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{m.query(endangered);actions++;since=0}
			if actions==7&&up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;since=0;eighthFired=true;eighthBoundary=start
			}else if maxActions>=9&&actions==8&&eighthFired&&start>eighthBoundary&&up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;ninth++;since=0
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1041000+step,step,up210cPolicy(step,advance,a,b)){return step,actions,ninth}
		}
	}
	return horizon+1,actions,ninth
}
func RunUP218C()(UP218CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{{"no_hostile","no_refresh","hostile_shield"},{"hostile_no","hostile_shield","no_refresh"},{"alternating_fixed","alternating_shield","fixed_offset_refresh"},{"fixed_alternating","fixed_offset_refresh","alternating_shield"}}
	advances:=[]int{8,12};horizons:=[]int{74,76,78,80}
	res:=UP218CResult{Schema:UP218CNinthSchema,Experiment:"UP-218C-fresh-critical-ninth-action",SourceUP217CSeal:"fd3336e688a2e0a61a4150e1cfbe59a540d6568c",CriticalHorizon:2,MaxActions:9,PhaseAdvances:advances,Horizons:horizons,CounterfactualOnly:true,LiveActivation:false,ThresholdFittingUsed:false,ScheduleRetuningUsed:false,FuturePolicyInputUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{
		m:=UP218CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			cap8,_,_:=up218cTreated(x,e,advance,h,s.a,s.b,8)
			cap9,_,ninth:=up218cTreated(x,e,advance,h,s.a,s.b,9)
			l8:=cap8<=h;l9:=cap9<=h
			if l8{m.Cap8Losses++};if l9{m.Cap9Losses++};m.NinthActions+=ninth
			if l8&&!l9{m.NinthRescues++}
			if !l8&&ninth>0{m.UnnecessaryNinth++}
			if l8&&ninth>0{m.NecessaryNinthAttempts++}
			if l8&&l9&&ninth>0{m.IneffectiveNinth++}
			if !l8&&ninth==0{m.NoNinthSurvivors++}
		}}
		if m.UnnecessaryNinth>0{m.RescuePerUnnecessary=float64(m.NinthRescues)/float64(m.UnnecessaryNinth)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
