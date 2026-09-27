package unitary

const UP177CActionCountSchema="wingless.up177c-delayed-action-count.v1"

type UP177CMetric struct{
	Cadence int `json:"cadence"`
	MaxActions int `json:"max_actions"`
	Intervention string `json:"intervention"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	ActionsTaken int `json:"actions_taken"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP177CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP176CSeal string `json:"source_up176c_seal"`
	InitialHands []int `json:"initial_hands"`
	Cadences []int `json:"cadences"`
	ActionCaps []int `json:"action_caps"`
	Interventions []string `json:"interventions"`
	DelayIntervals int `json:"delay_intervals"`
	Policy string `json:"policy"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveActionCapUsed bool `json:"adaptive_action_cap_used"`
	AdaptiveDelayUsed bool `json:"adaptive_delay_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP177CMetric `json:"metrics"`
}
func up177cRun(x0 *up81cAging,endangered,cadence,maxActions int,intervention string)(lossStep,actions int){
	m:=*x0;pending:=false
	for start:=1;start<=64;start+=cadence{
		actedThis:=false
		if pending&&actions<maxActions{
			if up175cAct(&m,endangered,intervention){actions++;actedThis=true}
			pending=false
		}
		if !actedThis&&!pending&&actions<maxActions{
			h:=up161cAdversarial(&m,endangered)
			if h<=cadence{pending=true}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1004000+step,step,"no_refresh"){return step,actions}
		}
	}
	return 65,actions
}
func RunUP177C()(UP177CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15};cadences:=[]int{2,4};caps:=[]int{1,2};interventions:=[]string{"targeted_refresh","sham_refresh"}
	res:=UP177CResult{Schema:UP177CActionCountSchema,Experiment:"UP-177C-delayed-action-count",SourceUP176CSeal:"8b79af488a6d5864078c1e6bc8ab065fd64ad7b3",InitialHands:hands,Cadences:cadences,ActionCaps:caps,Interventions:interventions,DelayIntervals:1,Policy:"no_refresh",CounterfactualOnly:true,LiveActivation:false,AdaptiveActionCapUsed:false,AdaptiveDelayUsed:false,WarningThresholdChanged:false}
	for _,cad:=range cadences{for _,cap:=range caps{for _,intervention:=range interventions{
		m:=UP177CMetric{Cadence:cad,MaxActions:cap,Intervention:intervention};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,endangered,cad,"no_refresh");if base<=64{m.BaselineLosses++}
			treated,actions:=up177cRun(x,endangered,cad,cap,intervention);m.ActionsTaken+=actions
			if treated<=64{m.TreatedLosses++;delaySum+=treated-base;delayN++}else if base<=64{m.PreventedLosses++}
		}}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
