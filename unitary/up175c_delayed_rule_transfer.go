package unitary

const UP175CDelayedTransferSchema="wingless.up175c-delayed-rule-transfer.v1"

type UP175CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Intervention string `json:"intervention"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	ActionsTaken int `json:"actions_taken"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP175CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP174CSeal string `json:"source_up174c_seal"`
	InitialHands []int `json:"initial_hands"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	Interventions []string `json:"interventions"`
	DelayIntervals int `json:"delay_intervals"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveDelayUsed bool `json:"adaptive_delay_used"`
	AdaptiveActionBudgetUsed bool `json:"adaptive_action_budget_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	FutureRealPolicyScheduleUsedByWarning bool `json:"future_real_policy_schedule_used_by_warning"`
	Metrics []UP175CMetric `json:"metrics"`
}
func up175cAct(m *up81cAging,endangered int,intervention string)bool{
	if intervention=="targeted_refresh"{m.query(endangered);return true}
	if intervention=="sham_refresh"{return up168cShamQuery(m,endangered)}
	return false
}
func up175cRun(x0 *up81cAging,endangered,cadence int,policy,intervention string)(lossStep,actions int){
	m:=*x0;pending:=false
	for start:=1;start<=64;start+=cadence{
		actedThisInterval:=false
		if pending&&actions<2{
			if up175cAct(&m,endangered,intervention){actions++;actedThisInterval=true}
			pending=false
		}
		if !actedThisInterval&&!pending&&actions<2{
			h:=up161cAdversarial(&m,endangered);warn:=h<=cadence
			if warn{pending=true}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1000000+step,step,policy){return step,actions}
		}
	}
	return 65,actions
}
func RunUP175C()(UP175CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4};interventions:=[]string{"targeted_refresh","sham_refresh"}
	res:=UP175CResult{Schema:UP175CDelayedTransferSchema,Experiment:"UP-175C-delayed-rule-transfer",SourceUP174CSeal:"53b69b7677df06f83b1e701d0dbe8b1af53ee7ea",InitialHands:hands,Policies:policies,Cadences:cadences,Interventions:interventions,DelayIntervals:1,MaxActionsPerArm:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveDelayUsed:false,AdaptiveActionBudgetUsed:false,WarningThresholdChanged:false,FutureRealPolicyScheduleUsedByWarning:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,intervention:=range interventions{
		m:=UP175CMetric{Policy:policy,Cadence:cad,Intervention:intervention};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,endangered,cad,policy);if base<=64{m.BaselineLosses++}
			treated,actions:=up175cRun(x,endangered,cad,policy,intervention);m.ActionsTaken+=actions
			if treated<=64{m.TreatedLosses++;delaySum+=treated-base;delayN++}else if base<=64{m.PreventedLosses++}
		}}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
