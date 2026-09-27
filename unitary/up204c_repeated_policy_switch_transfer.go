package unitary

const UP204CRepeatedSwitchSchema="wingless.up204c-repeated-policy-switch-transfer.v1"

type UP204CMetric struct{
	Schedule string `json:"schedule"`
	PolicyA string `json:"policy_a"`
	PolicyB string `json:"policy_b"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	ArmsWithAction int `json:"arms_with_action"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
}
type UP204CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP203CSeal string `json:"source_up203c_seal"`
	SwitchIntervalWrites int `json:"switch_interval_writes"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Horizon int `json:"horizon"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsedByWarning bool `json:"future_schedule_used_by_warning"`
	PolicySpecificRuleUsed bool `json:"policy_specific_rule_used"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP204CMetric `json:"metrics"`
}
func up204cPolicy(step int,a,b string)string{
	block:=(step-1)/8
	if block%2==0{return a}
	return b
}
func up204cBaseline(x0 *up81cAging,endangered int,a,b string)int{
	m:=*x0
	for start:=1;start<=56;start+=2{
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1026000+step,step,up204cPolicy(step,a,b)){return step}
		}
	}
	return 57
}
func up204cTreated(x0 *up81cAging,endangered int,a,b string)(lossStep,actions int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=56;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;committed=true;since=0
			}
		}else if actions<8{
			since++
			if since>=4{m.query(endangered);actions++;since=0}
		}
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1027000+step,step,up204cPolicy(step,a,b)){return step,actions}
		}
	}
	return 57,actions
}
func RunUP204C()(UP204CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile_repeated","no_refresh","hostile_shield"},
		{"hostile_no_repeated","hostile_shield","no_refresh"},
		{"alternating_fixed_repeated","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating_repeated","fixed_offset_refresh","alternating_shield"},
	}
	res:=UP204CResult{Schema:UP204CRepeatedSwitchSchema,Experiment:"UP-204C-repeated-policy-switch-transfer",SourceUP203CSeal:"dd4edb69ae2127725643f00a4564e7426ade8159",SwitchIntervalWrites:8,Cadence:2,SpacingIntervals:4,ActionCap:8,Horizon:56,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsedByWarning:false,PolicySpecificRuleUsed:false,AdaptiveSelectionUsed:false}
	for _,s:=range schedules{
		m:=UP204CMetric{Schedule:s.name,PolicyA:s.a,PolicyB:s.b};sumExt,nFail:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up204cBaseline(x,e,s.a,s.b)
			treated,actions:=up204cTreated(x,e,s.a,s.b)
			if base<=56{m.BaselineLosses++}
			if treated<=56{m.TreatedLosses++;sumExt+=treated-base;nFail++}
			if base<=56&&treated>56{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			if actions>0{m.ArmsWithAction++}
			m.ActionsTaken+=actions
		}}
		if nFail>0{m.MeanLossStepExtensionAmongFailures=float64(sumExt)/float64(nFail)}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
