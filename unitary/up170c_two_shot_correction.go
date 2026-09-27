package unitary

const UP170CTwoShotSchema="wingless.up170c-two-shot-correction.v1"

type UP170CMetric struct{
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
type UP170CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP169CSeal string `json:"source_up169c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	Interventions []string `json:"interventions"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveActionBudgetUsed bool `json:"adaptive_action_budget_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	FutureRealPolicyScheduleUsedByWarning bool `json:"future_real_policy_schedule_used_by_warning"`
	Metrics []UP170CMetric `json:"metrics"`
}
func up170cRun(x0 *up81cAging,endangered,cadence int,policy,intervention string)(lossStep,actions int){
	m:=*x0
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered);warn:=h<=cadence
		if warn&&actions<2{
			acted:=false
			if intervention=="targeted_refresh"{m.query(endangered);acted=true}
			if intervention=="sham_refresh"{acted=up168cShamQuery(&m,endangered)}
			if acted{actions++}
		}
		for j:=0;j<cadence&&start+j<=64;j++{step:=start+j;if up165cRealStep(&m,endangered,992000+step,step,policy){return step,actions}}
	}
	return 65,actions
}
func RunUP170C()(UP170CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4};interventions:=[]string{"targeted_refresh","sham_refresh"}
	res:=UP170CResult{Schema:UP170CTwoShotSchema,Experiment:"UP-170C-two-shot-correction",SourceUP169CSeal:"dc3131bf32e6860f34a53ab9f76d688e2ffefba6",Policies:policies,Cadences:cadences,Interventions:interventions,MaxActionsPerArm:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveActionBudgetUsed:false,WarningThresholdChanged:false,FutureRealPolicyScheduleUsedByWarning:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,intervention:=range interventions{
		m:=UP170CMetric{Policy:policy,Cadence:cad,Intervention:intervention};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,endangered,cad,policy);if base<=64{m.BaselineLosses++}
			treated,actions:=up170cRun(x,endangered,cad,policy,intervention);m.ActionsTaken+=actions
			if treated<=64{m.TreatedLosses++;delaySum+=treated-base;delayN++}else if base<=64{m.PreventedLosses++}
		}}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
