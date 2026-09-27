package unitary

const UP172CWarningTimingSchema="wingless.up172c-warning-timing-control.v1"

type UP172CMetric struct{
	Intervention string `json:"intervention"`
	Arms int `json:"arms"`
	ActionsTaken int `json:"actions_taken"`
	Losses int `json:"losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP172CPaired struct{
	WarningOnlyPrevents int `json:"warning_only_prevents"`
	FixedOnlyPrevents int `json:"fixed_only_prevents"`
	BothPrevent int `json:"both_prevent"`
	NeitherPrevent int `json:"neither_prevent"`
}
type UP172CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP171CSeal string `json:"source_up171c_seal"`
	InitialHands []int `json:"initial_hands"`
	Cadence int `json:"cadence"`
	Policy string `json:"policy"`
	FixedEarlyWriteStarts []int `json:"fixed_early_write_starts"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveTimingUsed bool `json:"adaptive_timing_used"`
	AdaptiveActionBudgetUsed bool `json:"adaptive_action_budget_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP172CMetric `json:"metrics"`
	Paired UP172CPaired `json:"paired"`
}
func up172cFixedEarly(x0 *up81cAging,endangered int)(lossStep,actions int){
	m:=*x0
	for start:=1;start<=64;start+=2{
		if (start==1||start==3)&&actions<2{m.query(endangered);actions++}
		for j:=0;j<2&&start+j<=64;j++{step:=start+j;if up165cRealStep(&m,endangered,994000+step,step,"no_refresh"){return step,actions}}
	}
	return 65,actions
}
func up172cNoAction(x0 *up81cAging,endangered int)int{
	return up169cBaselineLossStep(x0,endangered,2,"no_refresh")
}
func RunUP172C()(UP172CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	modes:=[]string{"warning_targeted","fixed_early_targeted","warning_sham","no_action"}
	res:=UP172CResult{Schema:UP172CWarningTimingSchema,Experiment:"UP-172C-warning-timing-control",SourceUP171CSeal:"a82baf0cb475c054621f289774d7d563dfc642e1",InitialHands:hands,Cadence:2,Policy:"no_refresh",FixedEarlyWriteStarts:[]int{1,3},MaxActionsPerArm:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveTimingUsed:false,AdaptiveActionBudgetUsed:false,WarningThresholdChanged:false}
	type arm struct{x *up81cAging;endangered int}
	arms:=[]arm{}
	for _,c:=range cohorts{for _,hand:=range hands{x,e,ok:=up161cTriggerState(hand,c);if ok{arms=append(arms,arm{x:x,endangered:e})}}}
	base:=make([]int,len(arms));warn:=make([]int,len(arms));fixed:=make([]int,len(arms))
	for i,a:=range arms{base[i]=up172cNoAction(a.x,a.endangered);warn[i],_=up170cRun(a.x,a.endangered,2,"no_refresh","targeted_refresh");fixed[i],_=up172cFixedEarly(a.x,a.endangered)}
	for _,mode:=range modes{
		m:=UP172CMetric{Intervention:mode,Arms:len(arms)};delaySum,delayN:=0,0
		for i,a:=range arms{
			b:=base[i];treated,actions:=b,0
			switch mode{
			case "warning_targeted": treated,actions=up170cRun(a.x,a.endangered,2,"no_refresh","targeted_refresh")
			case "fixed_early_targeted": treated,actions=up172cFixedEarly(a.x,a.endangered)
			case "warning_sham": treated,actions=up170cRun(a.x,a.endangered,2,"no_refresh","sham_refresh")
			case "no_action":
			}
			m.ActionsTaken+=actions
			if treated<=64{m.Losses++;delaySum+=treated-b;delayN++}else if b<=64{m.PreventedLosses++}
		}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}
	for i:=range arms{
		wp:=base[i]<=64&&warn[i]>64;fp:=base[i]<=64&&fixed[i]>64
		switch{case wp&&fp:res.Paired.BothPrevent++;case wp:res.Paired.WarningOnlyPrevents++;case fp:res.Paired.FixedOnlyPrevents++;default:res.Paired.NeitherPrevent++}
	}
	return res,nil
}
