package unitary

const UP174CWarningDelaySchema="wingless.up174c-warning-delay-tolerance.v1"

type UP174CMetric struct{
	DelayIntervals int `json:"delay_intervals"`
	Arms int `json:"arms"`
	ActionsTaken int `json:"actions_taken"`
	Losses int `json:"losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP174CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP173CSeal string `json:"source_up173c_seal"`
	InitialHands []int `json:"initial_hands"`
	Cadence int `json:"cadence"`
	Policy string `json:"policy"`
	DelayIntervals []int `json:"delay_intervals"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveDelayUsed bool `json:"adaptive_delay_used"`
	AdaptiveActionBudgetUsed bool `json:"adaptive_action_budget_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP174CMetric `json:"metrics"`
}
func up174cRun(x0 *up81cAging,endangered,delay int)(lossStep,actions int){
	m:=*x0
	pending:=0
	for start:=1;start<=64;start+=2{
		actedThisInterval:=false
		if pending>0{
			pending--
			if pending==0&&actions<2{
				m.query(endangered);actions++;actedThisInterval=true
			}
		}
		if !actedThisInterval&&pending==0&&actions<2{
			h:=up161cAdversarial(&m,endangered);warn:=h<=2
			if warn{
				if delay==0{m.query(endangered);actions++;actedThisInterval=true}else{pending=delay}
			}
		}
		for j:=0;j<2&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,998000+step,step,"no_refresh"){return step,actions}
		}
	}
	return 65,actions
}
func RunUP174C()(UP174CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15};delays:=[]int{0,1,2}
	res:=UP174CResult{Schema:UP174CWarningDelaySchema,Experiment:"UP-174C-warning-delay-tolerance",SourceUP173CSeal:"e5a89e17a3f717d2f0d917479fa0417c35ffe850",InitialHands:hands,Cadence:2,Policy:"no_refresh",DelayIntervals:delays,MaxActionsPerArm:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveDelayUsed:false,AdaptiveActionBudgetUsed:false,WarningThresholdChanged:false}
	type arm struct{x *up81cAging;endangered int}
	arms:=[]arm{}
	for _,c:=range cohorts{for _,hand:=range hands{x,e,ok:=up161cTriggerState(hand,c);if ok{arms=append(arms,arm{x:x,endangered:e})}}}
	base:=make([]int,len(arms));for i,a:=range arms{base[i]=up169cBaselineLossStep(a.x,a.endangered,2,"no_refresh")}
	for _,delay:=range delays{
		m:=UP174CMetric{DelayIntervals:delay,Arms:len(arms)};delaySum,delayN:=0,0
		for i,a:=range arms{
			treated,actions:=up174cRun(a.x,a.endangered,delay);m.ActionsTaken+=actions
			if treated<=64{m.Losses++;delaySum+=treated-base[i];delayN++}else if base[i]<=64{m.PreventedLosses++}
		}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
