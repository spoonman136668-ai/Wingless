package unitary

const UP189CCorrectionPolicySchema="wingless.up189c-correction-policy-transfer.v1"

type UP189CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	StageOneActions int `json:"stage_one_actions"`
	StageTwoActions int `json:"stage_two_actions"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossStepChangeAmongFailures float64 `json:"mean_loss_step_change_among_failures"`
}
type UP189CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP188CSeal string `json:"source_up188c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	MaxActions int `json:"max_actions"`
	StageOneDelayIntervals int `json:"stage_one_delay_intervals"`
	StageTwoDelayIntervals int `json:"stage_two_delay_intervals"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	AdaptiveDelayUsed bool `json:"adaptive_delay_used"`
	PolicySpecificRuleUsed bool `json:"policy_specific_rule_used"`
	FuturePolicyScheduleUsedByWarning bool `json:"future_policy_schedule_used_by_warning"`
	Metrics []UP189CMetric `json:"metrics"`
}
func up189cRun(x0 *up81cAging,endangered,cadence int,policy string)(lossStep,a1,a2 int){
	m:=*x0
	stage:=1
	pendingSecond:=false
	for start:=1;start<=64;start+=cadence{
		actedThis:=false
		if pendingSecond&&stage==2{
			m.query(endangered);a2++;stage++;pendingSecond=false;actedThis=true
		}
		if !actedThis&&stage<=2{
			h:=up161cAdversarial(&m,endangered)
			if h<=cadence{
				if stage==1{
					m.query(endangered);a1++;stage++;actedThis=true
				}else if stage==2{
					pendingSecond=true
				}
			}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1014000+step,step,policy){return step,a1,a2}
		}
	}
	return 65,a1,a2
}
func RunUP189C()(UP189CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"}
	cadences:=[]int{2,4}
	res:=UP189CResult{
		Schema:UP189CCorrectionPolicySchema,Experiment:"UP-189C-correction-policy-transfer",
		SourceUP188CSeal:"a64adbb6d9098d810eb908a0aa0987d25db96216",
		Policies:policies,Cadences:cadences,MaxActions:2,StageOneDelayIntervals:0,StageTwoDelayIntervals:1,
		CounterfactualOnly:true,LiveActivation:false,WarningThresholdChanged:false,AdaptiveDelayUsed:false,PolicySpecificRuleUsed:false,FuturePolicyScheduleUsedByWarning:false,
	}
	for _,policy:=range policies{for _,cad:=range cadences{
		m:=UP189CMetric{Policy:policy,Cadence:cad};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,policy)
			treated,a1,a2:=up189cRun(x,e,cad,policy)
			m.StageOneActions+=a1;m.StageTwoActions+=a2
			if base<=64{m.BaselineLosses++}
			if treated<=64{m.TreatedLosses++}
			if base<=64&&treated>64{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			if treated<=64{delaySum+=treated-base;delayN++}
		}}
		total:=m.StageOneActions+m.StageTwoActions
		if total>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(total)}
		if delayN>0{m.MeanLossStepChangeAmongFailures=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
