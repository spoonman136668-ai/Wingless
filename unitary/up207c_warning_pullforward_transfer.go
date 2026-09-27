package unitary

const UP207CPullForwardSchema="wingless.up207c-warning-pullforward-transfer.v1"

type UP207CMetric struct{
	Schedule string `json:"schedule"`
	Mode string `json:"mode"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	PullForwardActions int `json:"pull_forward_actions"`
	ArmsWithAction int `json:"arms_with_action"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
}
type UP207CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP206CSeal string `json:"source_up206c_seal"`
	PhaseShiftWrites int `json:"phase_shift_writes"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Horizon int `json:"horizon"`
	Modes []string `json:"modes"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ExtraBudgetUsed bool `json:"extra_budget_used"`
	PolicySpecificRuleUsed bool `json:"policy_specific_rule_used"`
	FutureScheduleUsedByWarning bool `json:"future_schedule_used_by_warning"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP207CMetric `json:"metrics"`
}
func up207cPolicy(step int,a,b string)string{
	if step<=4{return a}
	block:=(step-5)/8
	if block%2==0{return b}
	return a
}
func up207cBaseline(x0 *up81cAging,endangered int,a,b string)int{
	m:=*x0
	for start:=1;start<=56;start+=2{
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1030000+step,step,up207cPolicy(step,a,b)){return step}
		}
	}
	return 57
}
func up207cTreated(x0 *up81cAging,endangered int,a,b,mode string)(loss,actions,pull int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=56;start+=2{
		preH:=up161cAdversarial(&m,endangered)
		if !committed{
			if preH<=2{
				m.query(endangered);actions++;committed=true;since=0
			}
		}else if actions<8{
			since++
			act:=since>=4
			if mode=="pull_forward"&&preH<=2&&since<4{act=true;pull++}
			if act{m.query(endangered);actions++;since=0}
		}
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1031000+step,step,up207cPolicy(step,a,b)){return step,actions,pull}
		}
	}
	return 57,actions,pull
}
func RunUP207C()(UP207CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile_shifted","no_refresh","hostile_shield"},
		{"hostile_no_shifted","hostile_shield","no_refresh"},
		{"alternating_fixed_shifted","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating_shifted","fixed_offset_refresh","alternating_shield"},
	}
	modes:=[]string{"committed","pull_forward"}
	res:=UP207CResult{Schema:UP207CPullForwardSchema,Experiment:"UP-207C-warning-pullforward-transfer",SourceUP206CSeal:"e8a269efb3bf247d954617416d1fd783bba5faa4",PhaseShiftWrites:4,Cadence:2,SpacingIntervals:4,ActionCap:8,Horizon:56,Modes:modes,CounterfactualOnly:true,LiveActivation:false,ExtraBudgetUsed:false,PolicySpecificRuleUsed:false,FutureScheduleUsedByWarning:false,AdaptiveSelectionUsed:false}
	for _,s:=range schedules{for _,mode:=range modes{
		m:=UP207CMetric{Schedule:s.name,Mode:mode};sumExt,nFail:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up207cBaseline(x,e,s.a,s.b)
			treated,actions,pull:=up207cTreated(x,e,s.a,s.b,mode)
			if base<=56{m.BaselineLosses++}
			if treated<=56{m.TreatedLosses++;sumExt+=treated-base;nFail++}
			if base<=56&&treated>56{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			if actions>0{m.ArmsWithAction++}
			m.ActionsTaken+=actions;m.PullForwardActions+=pull
		}}
		if nFail>0{m.MeanLossStepExtensionAmongFailures=float64(sumExt)/float64(nFail)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
