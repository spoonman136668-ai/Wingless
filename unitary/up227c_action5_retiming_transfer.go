package unitary

const UP227CRetimingTransferSchema="wingless.up227c-action5-retiming-transfer.v1"

type UP227CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	ComparatorLosses int `json:"comparator_losses"`
	PrimaryLosses int `json:"primary_losses"`
	AdvancedAction5Arms int `json:"advanced_action5_arms"`
	RescuedLosses int `json:"rescued_losses"`
	UnnecessaryAdvances int `json:"unnecessary_advances"`
	NecessaryAdvanceAttempts int `json:"necessary_advance_attempts"`
	IneffectiveAdvanceAttempts int `json:"ineffective_advance_attempts"`
	TotalAdvanceWrites int `json:"total_advance_writes"`
}
type UP227CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP226CSeal string `json:"source_up226c_seal"`
	PolicyPairs []string `json:"policy_pairs"`
	PhaseAdvances []int `json:"phase_advances"`
	Horizons []int `json:"horizons"`
	MonitoringCadence int `json:"monitoring_cadence"`
	MaxActions int `json:"max_actions"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ScheduleCompressionUsed bool `json:"schedule_compression_used"`
	ActionBudgetIncreased bool `json:"action_budget_increased"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	LiveActivation bool `json:"live_activation"`
	Metrics []UP227CMetric `json:"metrics"`
}
func up227cRun(x0 *up81cAging,endangered,advance,horizon int,a,b string,retime bool)(loss int,advanced bool,advanceWrites int){
	m:=*x0;committed:=false;actions:=0
	steps:=make([]int,8);for i:=range steps{steps[i]=-1}
	action4:=-1;due5,due6,due7:=-1,-1,-1
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;steps[0]=start;committed=true}
		}else{
			if actions>=1&&actions<4{
				last:=steps[actions-1]
				if last>=0&&start>=last+8{
					m.query(endangered);actions++;steps[actions-1]=start
					if actions==4{action4=start;due5=action4+8;due6=action4+16;due7=action4+24}
				}
			}else if actions==4{
				fire:=start>=due5
				if retime&&start<due5&&up161cAdversarial(&m,endangered)<=2{fire=true;advanced=true;advanceWrites=due5-start}
				if fire{m.query(endangered);actions++;steps[4]=start}
			}else if actions==5{
				if start>=due6{m.query(endangered);actions++;steps[5]=start}
			}else if actions==6{
				if start>=due7{m.query(endangered);actions++;steps[6]=start}
			}else if actions==7{
				if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;steps[7]=start}
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1050000+step,step,up210cPolicy(step,advance,a,b)){return step,advanced,advanceWrites}
		}
	}
	return horizon+1,advanced,advanceWrites
}
func RunUP227C()(UP227CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{8,12};horizons:=[]int{74,80}
	res:=UP227CResult{
		Schema:UP227CRetimingTransferSchema,Experiment:"UP-227C-action5-retiming-transfer",
		SourceUP226CSeal:"16da4b7a1cf7e4c12dd63a8b84d423973ab0c52a",
		PolicyPairs:[]string{"no_hostile","hostile_no","alternating_fixed","fixed_alternating"},
		PhaseAdvances:advances,Horizons:horizons,MonitoringCadence:2,MaxActions:8,
		CounterfactualOnly:true,ThresholdFittingUsed:false,ScheduleCompressionUsed:false,
		ActionBudgetIncreased:false,FuturePolicyInputUsed:false,LiveActivation:false,
	}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{
		m:=UP227CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base,_,_:=up227cRun(x,e,advance,h,s.a,s.b,false)
			primary,advanced,delta:=up227cRun(x,e,advance,h,s.a,s.b,true)
			if base<=h{m.ComparatorLosses++};if primary<=h{m.PrimaryLosses++}
			if advanced{
				m.AdvancedAction5Arms++;m.TotalAdvanceWrites+=delta
				if base<=h{m.NecessaryAdvanceAttempts++}else{m.UnnecessaryAdvances++}
				if base<=h&&primary>h{m.RescuedLosses++}
				if base<=h&&primary<=h{m.IneffectiveAdvanceAttempts++}
			}
		}}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
