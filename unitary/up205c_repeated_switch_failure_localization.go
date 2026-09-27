package unitary

type UP205CPoint struct{
	Schedule string `json:"schedule"`
	CohortIndex int `json:"cohort_index"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	BaselineLossStep int `json:"baseline_loss_step"`
	TreatedLossStep int `json:"treated_loss_step"`
	TriggerOnset int `json:"trigger_onset"`
	ActionCount int `json:"action_count"`
	ActionWrites []int `json:"action_writes"`
	LastActionWrite int `json:"last_action_write"`
	ActionCapReached bool `json:"action_cap_reached"`
	LossDelay int `json:"loss_delay"`
	SurvivesHorizon bool `json:"survives_horizon"`
}
type UP205CSummary struct{
	Schedule string `json:"schedule"`
	Arms int `json:"arms"`
	TreatedFailures int `json:"treated_failures"`
	FailuresAtActionCap int `json:"failures_at_action_cap"`
	MeanActionsAmongFailures float64 `json:"mean_actions_among_failures"`
	MeanLastActionToFailure float64 `json:"mean_last_action_to_failure"`
	MinTreatedFailureStep int `json:"min_treated_failure_step"`
	MaxTreatedFailureStep int `json:"max_treated_failure_step"`
}
type UP205CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP204CSeal string `json:"source_up204c_seal"`
	SwitchIntervalWrites int `json:"switch_interval_writes"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Horizon int `json:"horizon"`
	RuleChanged bool `json:"rule_changed"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Points []UP205CPoint `json:"points"`
	Summaries []UP205CSummary `json:"summaries"`
}
const UP205CLocalizeSchema="wingless.up205c-repeated-switch-failure-localization.v1"

func up205cTreated(x0 *up81cAging,endangered int,a,b string)(lossStep,onset int,actionWrites []int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=56;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actionWrites=append(actionWrites,start);onset=start;committed=true;since=0
			}
		}else if len(actionWrites)<8{
			since++
			if since>=4{m.query(endangered);actionWrites=append(actionWrites,start);since=0}
		}
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1028000+step,step,up204cPolicy(step,a,b)){return step,onset,actionWrites}
		}
	}
	return 57,onset,actionWrites
}
func RunUP205C()(UP205CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile_repeated","no_refresh","hostile_shield"},
		{"hostile_no_repeated","hostile_shield","no_refresh"},
		{"alternating_fixed_repeated","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating_repeated","fixed_offset_refresh","alternating_shield"},
	}
	res:=UP205CResult{Schema:UP205CLocalizeSchema,Experiment:"UP-205C-repeated-switch-failure-localization",SourceUP204CSeal:"922d8cdf5a6a06a5298e0e01472b5362c742ba05",SwitchIntervalWrites:8,Cadence:2,SpacingIntervals:4,ActionCap:8,Horizon:56,RuleChanged:false,CounterfactualOnly:true,LiveActivation:false}
	for _,s:=range schedules{
		sm:=UP205CSummary{Schedule:s.name,MinTreatedFailureStep:1<<30}
		sumActions,sumGap,nFail:=0,0,0
		for ci,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};sm.Arms++
			base:=up204cBaseline(x,e,s.a,s.b)
			treated,onset,writes:=up205cTreated(x,e,s.a,s.b)
			last:=0;if len(writes)>0{last=writes[len(writes)-1]}
			surv:=treated>56
			res.Points=append(res.Points,UP205CPoint{Schedule:s.name,CohortIndex:ci,InitialHand:hand,EndangeredKey:e,BaselineLossStep:base,TreatedLossStep:treated,TriggerOnset:onset,ActionCount:len(writes),ActionWrites:append([]int(nil),writes...),LastActionWrite:last,ActionCapReached:len(writes)==8,LossDelay:treated-base,SurvivesHorizon:surv})
			if !surv{
				sm.TreatedFailures++;nFail++;sumActions+=len(writes)
				if len(writes)==8{sm.FailuresAtActionCap++}
				if last>0{sumGap+=treated-last}
				if treated<sm.MinTreatedFailureStep{sm.MinTreatedFailureStep=treated}
				if treated>sm.MaxTreatedFailureStep{sm.MaxTreatedFailureStep=treated}
			}
		}}
		if nFail>0{sm.MeanActionsAmongFailures=float64(sumActions)/float64(nFail);sm.MeanLastActionToFailure=float64(sumGap)/float64(nFail)}
		if sm.MinTreatedFailureStep==1<<30{sm.MinTreatedFailureStep=0}
		res.Summaries=append(res.Summaries,sm)
	}
	return res,nil
}
