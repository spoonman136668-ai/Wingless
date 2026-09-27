package unitary

const UP226CAdvance5Schema="wingless.up226c-critical-advanced-action5.v1"

type UP226CArm struct{
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	ComparatorLossStep int `json:"comparator_loss_step"`
	PrimaryLossStep int `json:"primary_loss_step"`
	ComparatorActionSteps []int `json:"comparator_action_steps"`
	PrimaryActionSteps []int `json:"primary_action_steps"`
	Action5Advanced bool `json:"action5_advanced"`
	Action5AdvanceWrites int `json:"action5_advance_writes"`
}
type UP226CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP225CSeal string `json:"source_up225c_seal"`
	Horizon int `json:"horizon"`
	MonitoringCadence int `json:"monitoring_cadence"`
	MaxActions int `json:"max_actions"`
	ComparatorLosses int `json:"comparator_losses"`
	PrimaryLosses int `json:"primary_losses"`
	AdvancedAction5Arms int `json:"advanced_action5_arms"`
	RescuedLosses int `json:"rescued_losses"`
	UnnecessaryAdvances int `json:"unnecessary_advances"`
	NecessaryAdvanceAttempts int `json:"necessary_advance_attempts"`
	IneffectiveAdvanceAttempts int `json:"ineffective_advance_attempts"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ScheduleCompressionUsed bool `json:"schedule_compression_used"`
	ActionBudgetIncreased bool `json:"action_budget_increased"`
	FuturePolicyInputUsed bool `json:"future_policy_input_used"`
	LiveActivation bool `json:"live_activation"`
	Arms []UP226CArm `json:"arms"`
}

func up226cRun(x0 *up81cAging,endangered int,advance5 bool)(loss int,steps []int,advanced bool,advanceWrites int){
	m:=*x0
	steps=make([]int,8);for i:=range steps{steps[i]=-1}
	committed:=false;actions:=0
	// Fixed comparator due boundaries are derived from the actual action-4 step.
	action4Step:=-1
	due5,due6,due7:=-1,-1,-1
	for start:=1;start<=80;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;steps[0]=start;committed=true
			}
		}else{
			if actions>=1&&actions<4{
				last:=steps[actions-1]
				if last>=0&&start>=last+8{
					m.query(endangered);actions++;steps[actions-1]=start
					if actions==4{
						action4Step=start
						due5=action4Step+8;due6=action4Step+16;due7=action4Step+24
					}
				}
			}else if actions==4{
				fire:=start>=due5
				if advance5&&start<due5&&up161cAdversarial(&m,endangered)<=2{fire=true;advanced=true;advanceWrites=due5-start}
				if fire{m.query(endangered);actions++;steps[4]=start}
			}else if actions==5{
				if start>=due6{m.query(endangered);actions++;steps[5]=start}
			}else if actions==6{
				if start>=due7{m.query(endangered);actions++;steps[6]=start}
			}else if actions==7{
				if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;steps[7]=start}
			}
		}
		for j:=0;j<2&&start+j<=80;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1049000+step,step,up210cPolicy(step,8,"alternating_shield","fixed_offset_refresh")){
				return step,steps,advanced,advanceWrites
			}
		}
	}
	return 81,steps,advanced,advanceWrites
}
func RunUP226C()(UP226CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res:=UP226CResult{
		Schema:UP226CAdvance5Schema,Experiment:"UP-226C-critical-advanced-action5",
		SourceUP225CSeal:"9f0985eecb8d0580f28d5399e930d37d0cc3cd1a",
		Horizon:80,MonitoringCadence:2,MaxActions:8,
		CounterfactualOnly:true,ThresholdFittingUsed:false,ScheduleCompressionUsed:false,
		ActionBudgetIncreased:false,FuturePolicyInputUsed:false,LiveActivation:false,
	}
	for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
		base,baseSteps,_,_:=up226cRun(x,e,false)
		primary,primarySteps,advanced,delta:=up226cRun(x,e,true)
		a:=UP226CArm{Cohort:ci,InitialHand:hand,EndangeredKey:e,ComparatorLossStep:base,PrimaryLossStep:primary,ComparatorActionSteps:baseSteps,PrimaryActionSteps:primarySteps,Action5Advanced:advanced,Action5AdvanceWrites:delta}
		if base<=80{res.ComparatorLosses++};if primary<=80{res.PrimaryLosses++}
		if advanced{
			res.AdvancedAction5Arms++
			if base<=80{res.NecessaryAdvanceAttempts++}else{res.UnnecessaryAdvances++}
			if base<=80&&primary>80{res.RescuedLosses++}
			if base<=80&&primary<=80{res.IneffectiveAdvanceAttempts++}
		}
		res.Arms=append(res.Arms,a)
	}}
	return res,nil
}
