package unitary

const UP230CHysteresisSchema="wingless.up230c-action9-hysteresis.v1"

type UP230CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Arms int `json:"arms"`
	Cap8Losses int `json:"cap8_losses"`
	SingleCriticalLosses int `json:"single_critical_losses"`
	HysteresisLosses int `json:"hysteresis_losses"`
	SingleNinthActions int `json:"single_ninth_actions"`
	HysteresisNinthActions int `json:"hysteresis_ninth_actions"`
	SingleRescues int `json:"single_rescues"`
	HysteresisRescues int `json:"hysteresis_rescues"`
	SingleUnnecessaryNinth int `json:"single_unnecessary_ninth"`
	HysteresisUnnecessaryNinth int `json:"hysteresis_unnecessary_ninth"`
	SingleIneffectiveNinth int `json:"single_ineffective_ninth"`
	HysteresisIneffectiveNinth int `json:"hysteresis_ineffective_ninth"`
	SingleRescuePerUnnecessary float64 `json:"single_rescue_per_unnecessary"`
	HysteresisRescuePerUnnecessary float64 `json:"hysteresis_rescue_per_unnecessary"`
}
type UP230CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP229CSeal string `json:"source_up229c_seal"`
	PolicyPairs []string `json:"policy_pairs"`
	PhaseAdvances []int `json:"phase_advances"`
	Horizon int `json:"horizon"`
	MonitoringCadence int `json:"monitoring_cadence"`
	MaxActions int `json:"max_actions"`
	CriticalThreshold int `json:"critical_threshold"`
	PersistenceCount int `json:"persistence_count"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	ThresholdChanged bool `json:"threshold_changed"`
	ActionBudgetIncreased bool `json:"action_budget_increased"`
	NewActionTypeUsed bool `json:"new_action_type_used"`
	CohortSpecificTuningUsed bool `json:"cohort_specific_tuning_used"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	LiveActivation bool `json:"live_activation"`
	Metrics []UP230CMetric `json:"metrics"`
}
func up230cRun(x0 *up81cAging,endangered,advance int,a,b,mode string)(loss,actions,ninth int){
	m:=*x0;committed:=false
	steps:=make([]int,9);for i:=range steps{steps[i]=-1}
	action4:=-1;due5,due6,due7:=-1,-1,-1;eighthBoundary:=-1;criticalStreak:=0
	for start:=1;start<=80;start+=2{
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
				if start<due5&&up161cAdversarial(&m,endangered)<=2{fire=true}
				if fire{m.query(endangered);actions++;steps[4]=start}
			}else if actions==5{
				if start>=due6{m.query(endangered);actions++;steps[5]=start}
			}else if actions==6{
				if start>=due7{m.query(endangered);actions++;steps[6]=start}
			}else if actions==7{
				if up161cAdversarial(&m,endangered)<=2{
					m.query(endangered);actions++;steps[7]=start;eighthBoundary=start;criticalStreak=0
				}
			}else if actions==8&&mode!="none"&&eighthBoundary>=0&&start>eighthBoundary{
				critical:=up161cAdversarial(&m,endangered)<=2
				if mode=="single_critical"{
					if critical{m.query(endangered);actions++;steps[8]=start;ninth++}
				}else{
					if critical{criticalStreak++}else{criticalStreak=0}
					if criticalStreak>=2{m.query(endangered);actions++;steps[8]=start;ninth++}
				}
			}
		}
		for j:=0;j<2&&start+j<=80;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1052000+step,step,up210cPolicy(step,advance,a,b)){return step,actions,ninth}
		}
	}
	return 81,actions,ninth
}
func RunUP230C()(UP230CResult,error){
	cohorts:=[][]int{{0,5,10,15},{1,6,11,12},{2,7,8,13},{3,4,9,14}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{8,12}
	res:=UP230CResult{
		Schema:UP230CHysteresisSchema,Experiment:"UP-230C-action9-hysteresis",
		SourceUP229CSeal:"f0c4678f08038b23942d5e6fe4d942378723f7e4",
		PolicyPairs:[]string{"no_hostile","hostile_no","alternating_fixed","fixed_alternating"},
		PhaseAdvances:advances,Horizon:80,MonitoringCadence:2,MaxActions:9,CriticalThreshold:2,PersistenceCount:2,
		CounterfactualOnly:true,ThresholdChanged:false,ActionBudgetIncreased:false,NewActionTypeUsed:false,
		CohortSpecificTuningUsed:false,FuturePolicyInputUsed:false,LiveActivation:false,
	}
	for _,s:=range schedules{for _,advance:=range advances{
		m:=UP230CMetric{Schedule:s.name,PhaseAdvanceWrites:advance}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base,_,_:=up230cRun(x,e,advance,s.a,s.b,"none")
			single,_,sn:=up230cRun(x,e,advance,s.a,s.b,"single_critical")
			hyst,_,hn:=up230cRun(x,e,advance,s.a,s.b,"two_consecutive_critical")
			bf:=base<=80;sf:=single<=80;hf:=hyst<=80
			if bf{m.Cap8Losses++};if sf{m.SingleCriticalLosses++};if hf{m.HysteresisLosses++}
			m.SingleNinthActions+=sn;m.HysteresisNinthActions+=hn
			if bf&&!sf{m.SingleRescues++};if bf&&!hf{m.HysteresisRescues++}
			if !bf&&sn>0{m.SingleUnnecessaryNinth++};if !bf&&hn>0{m.HysteresisUnnecessaryNinth++}
			if bf&&sf&&sn>0{m.SingleIneffectiveNinth++};if bf&&hf&&hn>0{m.HysteresisIneffectiveNinth++}
		}}
		if m.SingleUnnecessaryNinth>0{m.SingleRescuePerUnnecessary=float64(m.SingleRescues)/float64(m.SingleUnnecessaryNinth)}
		if m.HysteresisUnnecessaryNinth>0{m.HysteresisRescuePerUnnecessary=float64(m.HysteresisRescues)/float64(m.HysteresisUnnecessaryNinth)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
