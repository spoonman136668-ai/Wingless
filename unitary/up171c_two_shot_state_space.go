package unitary

const UP171CStateSpaceSchema="wingless.up171c-two-shot-state-space.v1"

type UP171CMetric struct{
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
type UP171CSubset struct{
	Subset string `json:"subset"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	ActionsTaken int `json:"actions_taken"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
}
type UP171CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP170CSeal string `json:"source_up170c_seal"`
	InitialHands []int `json:"initial_hands"`
	OriginalHands []int `json:"original_hands"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	Interventions []string `json:"interventions"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveHandSelectionUsed bool `json:"adaptive_hand_selection_used"`
	AdaptiveActionBudgetUsed bool `json:"adaptive_action_budget_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP171CMetric `json:"metrics"`
	NoRefreshCadence2TargetedSubsets []UP171CSubset `json:"no_refresh_cadence2_targeted_subsets"`
}
func up171cOriginalHand(hand int)bool{return hand==0||hand==4||hand==8||hand==12}
func RunUP171C()(UP171CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15};orig:=[]int{0,4,8,12}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4};interventions:=[]string{"targeted_refresh","sham_refresh"}
	res:=UP171CResult{Schema:UP171CStateSpaceSchema,Experiment:"UP-171C-two-shot-state-space",SourceUP170CSeal:"81cc048a2a5a54fcc5c064b632859b2ad5e65c7c",InitialHands:hands,OriginalHands:orig,Policies:policies,Cadences:cadences,Interventions:interventions,MaxActionsPerArm:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveHandSelectionUsed:false,AdaptiveActionBudgetUsed:false,WarningThresholdChanged:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,intervention:=range interventions{
		m:=UP171CMetric{Policy:policy,Cadence:cad,Intervention:intervention};delaySum,delayN:=0,0
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
	for _,name:=range []string{"original_hands","expanded_hands"}{
		s:=UP171CSubset{Subset:name}
		for _,c:=range cohorts{for _,hand:=range hands{
			if (name=="original_hands")!=up171cOriginalHand(hand){continue}
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};s.Arms++
			base:=up169cBaselineLossStep(x,endangered,2,"no_refresh");if base<=64{s.BaselineLosses++}
			treated,actions:=up170cRun(x,endangered,2,"no_refresh","targeted_refresh");s.ActionsTaken+=actions
			if treated<=64{s.TreatedLosses++}else if base<=64{s.PreventedLosses++}
		}}
		res.NoRefreshCadence2TargetedSubsets=append(res.NoRefreshCadence2TargetedSubsets,s)
	}
	return res,nil
}
