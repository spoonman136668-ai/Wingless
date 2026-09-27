package unitary

const UP176CSecondGateSchema="wingless.up176c-second-action-gate.v1"

type UP176CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	SecondActionGate string `json:"second_action_gate"`
	Intervention string `json:"intervention"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	ActionsTaken int `json:"actions_taken"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP176CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP175CSeal string `json:"source_up175c_seal"`
	InitialHands []int `json:"initial_hands"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	SecondActionGates []string `json:"second_action_gates"`
	Interventions []string `json:"interventions"`
	FirstActionDelayIntervals int `json:"first_action_delay_intervals"`
	ForcedSecondDelayIntervals int `json:"forced_second_delay_intervals"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveGateUsed bool `json:"adaptive_gate_used"`
	AdaptiveActionBudgetUsed bool `json:"adaptive_action_budget_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP176CMetric `json:"metrics"`
}
func up176cForcedRun(x0 *up81cAging,endangered,cadence int,policy,intervention string)(lossStep,actions int){
	m:=*x0
	firstPending:=false
	secondPending:=false
	for start:=1;start<=64;start+=cadence{
		actedThisInterval:=false
		if secondPending&&actions==1{
			if up175cAct(&m,endangered,intervention){actions++;actedThisInterval=true}
			secondPending=false
		}
		if firstPending&&actions==0&&!actedThisInterval{
			if up175cAct(&m,endangered,intervention){actions++;actedThisInterval=true;secondPending=true}
			firstPending=false
		}
		if actions==0&&!firstPending&&!actedThisInterval{
			h:=up161cAdversarial(&m,endangered);if h<=cadence{firstPending=true}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1002000+step,step,policy){return step,actions}
		}
	}
	return 65,actions
}
func RunUP176C()(UP176CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4};gates:=[]string{"rewarning_gated","forced_second"};interventions:=[]string{"targeted_refresh","sham_refresh"}
	res:=UP176CResult{Schema:UP176CSecondGateSchema,Experiment:"UP-176C-second-action-gate",SourceUP175CSeal:"f878647f1722950e529d43b5189f5d665091964e",InitialHands:hands,Policies:policies,Cadences:cadences,SecondActionGates:gates,Interventions:interventions,FirstActionDelayIntervals:1,ForcedSecondDelayIntervals:1,MaxActionsPerArm:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveGateUsed:false,AdaptiveActionBudgetUsed:false,WarningThresholdChanged:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,gate:=range gates{for _,intervention:=range interventions{
		m:=UP176CMetric{Policy:policy,Cadence:cad,SecondActionGate:gate,Intervention:intervention};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,endangered,cad,policy);if base<=64{m.BaselineLosses++}
			treated,actions:=0,0
			if gate=="rewarning_gated"{treated,actions=up175cRun(x,endangered,cad,policy,intervention)}else{treated,actions=up176cForcedRun(x,endangered,cad,policy,intervention)}
			m.ActionsTaken+=actions
			if treated<=64{m.TreatedLosses++;delaySum+=treated-base;delayN++}else if base<=64{m.PreventedLosses++}
		}}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}}}
	return res,nil
}
