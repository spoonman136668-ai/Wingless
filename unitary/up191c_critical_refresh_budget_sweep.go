package unitary

const UP191CCriticalBudgetSchema="wingless.up191c-critical-refresh-budget-sweep.v1"

type UP191CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	ActionCap int `json:"action_cap"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossStepChangeAmongFailures float64 `json:"mean_loss_step_change_among_failures"`
}
type UP191CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP190CSeal string `json:"source_up190c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	ActionCaps []int `json:"action_caps"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCapUsed bool `json:"adaptive_cap_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	PolicySpecificRuleUsed bool `json:"policy_specific_rule_used"`
	NewActionTypeUsed bool `json:"new_action_type_used"`
	Metrics []UP191CMetric `json:"metrics"`
}
func up191cRun(x0 *up81cAging,endangered,cadence int,policy string,cap int)(lossStep,actions int){
	m:=*x0
	stage:=1
	pendingSecond:=false
	for start:=1;start<=64;start+=cadence{
		actedThis:=false
		if pendingSecond&&stage==2&&actions<cap{
			m.query(endangered);actions++;stage=3;pendingSecond=false;actedThis=true
		}
		if !actedThis&&actions<cap{
			h:=up161cAdversarial(&m,endangered)
			if h<=cadence{
				if stage==1{
					m.query(endangered);actions++;stage=2;actedThis=true
				}else if stage==2{
					if cap>=2{pendingSecond=true}
				}else{
					m.query(endangered);actions++;actedThis=true
				}
			}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1015000+step,step,policy){return step,actions}
		}
	}
	return 65,actions
}
func RunUP191C()(UP191CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"}
	cadences:=[]int{2,4}
	caps:=[]int{1,2,3,4,6,8,12,16}
	res:=UP191CResult{
		Schema:UP191CCriticalBudgetSchema,Experiment:"UP-191C-critical-refresh-budget-sweep",
		SourceUP190CSeal:"c77d611c31413454c4c8f0772fbec26ce495060c",
		Policies:policies,Cadences:cadences,ActionCaps:caps,
		CounterfactualOnly:true,LiveActivation:false,AdaptiveCapUsed:false,WarningThresholdChanged:false,PolicySpecificRuleUsed:false,NewActionTypeUsed:false,
	}
	for _,policy:=range policies{for _,cad:=range cadences{for _,cap:=range caps{
		m:=UP191CMetric{Policy:policy,Cadence:cad,ActionCap:cap};sumDelay,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,policy)
			treated,actions:=up191cRun(x,e,cad,policy,cap);m.ActionsTaken+=actions
			if base<=64{m.BaselineLosses++}
			if treated<=64{m.TreatedLosses++}
			if base<=64&&treated>64{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			if treated<=64{sumDelay+=treated-base;delayN++}
		}}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossStepChangeAmongFailures=float64(sumDelay)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
