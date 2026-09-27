package unitary

const UP228CCompositionSchema="wingless.up228c-bounded-correction-composition.v1"

type UP228CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	RetimeOnlyLosses int `json:"retime_only_losses"`
	Tail9OnlyLosses int `json:"tail9_only_losses"`
	CombinedLosses int `json:"combined_losses"`
	BaselineActions int `json:"baseline_actions"`
	RetimeOnlyActions int `json:"retime_only_actions"`
	Tail9OnlyActions int `json:"tail9_only_actions"`
	CombinedActions int `json:"combined_actions"`
	RetimeOnlyAdvanceEvents int `json:"retime_only_advance_events"`
	CombinedAdvanceEvents int `json:"combined_advance_events"`
	RetimeOnlyAdvanceWrites int `json:"retime_only_advance_writes"`
	CombinedAdvanceWrites int `json:"combined_advance_writes"`
	Tail9OnlyNinthActions int `json:"tail9_only_ninth_actions"`
	CombinedNinthActions int `json:"combined_ninth_actions"`
	CombinedRescueBeyondBaseline int `json:"combined_rescue_beyond_baseline"`
	CombinedRescueBeyondRetimeOnly int `json:"combined_rescue_beyond_retime_only"`
	CombinedRescueBeyondTail9Only int `json:"combined_rescue_beyond_tail9_only"`
}
type UP228CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP227CSeal string `json:"source_up227c_seal"`
	PolicyArms []string `json:"policy_arms"`
	PolicyPairs []string `json:"policy_pairs"`
	PhaseAdvances []int `json:"phase_advances"`
	Horizons []int `json:"horizons"`
	MonitoringCadence int `json:"monitoring_cadence"`
	MaxActions int `json:"max_actions"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	NewTriggerUsed bool `json:"new_trigger_used"`
	NewActionTypeUsed bool `json:"new_action_type_used"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ScheduleCompressionUsed bool `json:"schedule_compression_used"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	LiveActivation bool `json:"live_activation"`
	Metrics []UP228CMetric `json:"metrics"`
}
func up228cRun(x0 *up81cAging,endangered,advance,horizon int,a,b string,retime5,allow9 bool)(loss,actions int,advanced bool,advanceWrites,ninth int){
	m:=*x0;committed:=false
	steps:=make([]int,9);for i:=range steps{steps[i]=-1}
	action4:=-1;due5,due6,due7:=-1,-1,-1;eighthBoundary:=-1
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;steps[0]=start;committed=true
			}
		}else{
			if actions>=1&&actions<4{
				last:=steps[actions-1]
				if last>=0&&start>=last+8{
					m.query(endangered);actions++;steps[actions-1]=start
					if actions==4{action4=start;due5=action4+8;due6=action4+16;due7=action4+24}
				}
			}else if actions==4{
				fire:=start>=due5
				if retime5&&start<due5&&up161cAdversarial(&m,endangered)<=2{fire=true;advanced=true;advanceWrites=due5-start}
				if fire{m.query(endangered);actions++;steps[4]=start}
			}else if actions==5{
				if start>=due6{m.query(endangered);actions++;steps[5]=start}
			}else if actions==6{
				if start>=due7{m.query(endangered);actions++;steps[6]=start}
			}else if actions==7{
				if up161cAdversarial(&m,endangered)<=2{
					m.query(endangered);actions++;steps[7]=start;eighthBoundary=start
				}
			}else if actions==8&&allow9&&eighthBoundary>=0&&start>eighthBoundary{
				if up161cAdversarial(&m,endangered)<=2{
					m.query(endangered);actions++;steps[8]=start;ninth++
				}
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1051000+step,step,up210cPolicy(step,advance,a,b)){return step,actions,advanced,advanceWrites,ninth}
		}
	}
	return horizon+1,actions,advanced,advanceWrites,ninth
}
func RunUP228C()(UP228CResult,error){
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
	res:=UP228CResult{
		Schema:UP228CCompositionSchema,Experiment:"UP-228C-bounded-correction-composition",
		SourceUP227CSeal:"128e1909e0b4f460628de6806bca6f6da4ac4f04",
		PolicyArms:[]string{"baseline","retime_only","tail9_only","combined"},
		PolicyPairs:[]string{"no_hostile","hostile_no","alternating_fixed","fixed_alternating"},
		PhaseAdvances:advances,Horizons:horizons,MonitoringCadence:2,MaxActions:9,
		CounterfactualOnly:true,NewTriggerUsed:false,NewActionTypeUsed:false,
		ThresholdFittingUsed:false,ScheduleCompressionUsed:false,FuturePolicyInputUsed:false,LiveActivation:false,
	}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{
		m:=UP228CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			bLoss,bAct,_,_,_:=up228cRun(x,e,advance,h,s.a,s.b,false,false)
			rLoss,rAct,rAdv,rDelta,_:=up228cRun(x,e,advance,h,s.a,s.b,true,false)
			tLoss,tAct,_,_,tNinth:=up228cRun(x,e,advance,h,s.a,s.b,false,true)
			cLoss,cAct,cAdv,cDelta,cNinth:=up228cRun(x,e,advance,h,s.a,s.b,true,true)
			if bLoss<=h{m.BaselineLosses++};if rLoss<=h{m.RetimeOnlyLosses++};if tLoss<=h{m.Tail9OnlyLosses++};if cLoss<=h{m.CombinedLosses++}
			m.BaselineActions+=bAct;m.RetimeOnlyActions+=rAct;m.Tail9OnlyActions+=tAct;m.CombinedActions+=cAct
			if rAdv{m.RetimeOnlyAdvanceEvents++;m.RetimeOnlyAdvanceWrites+=rDelta}
			if cAdv{m.CombinedAdvanceEvents++;m.CombinedAdvanceWrites+=cDelta}
			m.Tail9OnlyNinthActions+=tNinth;m.CombinedNinthActions+=cNinth
			if bLoss<=h&&cLoss>h{m.CombinedRescueBeyondBaseline++}
			if rLoss<=h&&cLoss>h{m.CombinedRescueBeyondRetimeOnly++}
			if tLoss<=h&&cLoss>h{m.CombinedRescueBeyondTail9Only++}
		}}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
