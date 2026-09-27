package unitary

const UP203CPolicySwitchSchema="wingless.up203c-policy-switch-transfer.v1"

type UP203CMetric struct{
	Schedule string `json:"schedule"`
	PolicyBefore string `json:"policy_before"`
	PolicyAfter string `json:"policy_after"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	ArmsWithAction int `json:"arms_with_action"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
}
type UP203CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP202CSeal string `json:"source_up202c_seal"`
	SwitchWrite int `json:"switch_write"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Horizon int `json:"horizon"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsedByWarning bool `json:"future_schedule_used_by_warning"`
	PolicySpecificRuleUsed bool `json:"policy_specific_rule_used"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP203CMetric `json:"metrics"`
}
func up203cPolicy(step int,before,after string)string{
	if step<24{return before}
	return after
}
func up203cBaseline(x0 *up81cAging,endangered int,before,after string)int{
	m:=*x0
	for start:=1;start<=56;start+=2{
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1024000+step,step,up203cPolicy(step,before,after)){return step}
		}
	}
	return 57
}
func up203cTreated(x0 *up81cAging,endangered int,before,after string)(lossStep,actions int){
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
			if up165cRealStep(&m,endangered,1025000+step,step,up203cPolicy(step,before,after)){return step,actions}
		}
	}
	return 57,actions
}
func RunUP203C()(UP203CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,before,after string}
	schedules:=[]sched{
		{"no_to_hostile","no_refresh","hostile_shield"},
		{"hostile_to_no","hostile_shield","no_refresh"},
		{"alternating_to_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_to_alternating","fixed_offset_refresh","alternating_shield"},
	}
	res:=UP203CResult{Schema:UP203CPolicySwitchSchema,Experiment:"UP-203C-policy-switch-transfer",SourceUP202CSeal:"d53d84720a7706eba18e2ebfdb08c103f28e060c",SwitchWrite:24,Cadence:2,SpacingIntervals:4,ActionCap:8,Horizon:56,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsedByWarning:false,PolicySpecificRuleUsed:false,AdaptiveSelectionUsed:false}
	for _,s:=range schedules{
		m:=UP203CMetric{Schedule:s.name,PolicyBefore:s.before,PolicyAfter:s.after};sumExt,nFail:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up203cBaseline(x,e,s.before,s.after)
			treated,actions:=up203cTreated(x,e,s.before,s.after)
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
