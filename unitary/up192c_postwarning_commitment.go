package unitary

const UP192CCommitmentSchema="wingless.up192c-postwarning-commitment.v1"

type UP192CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	ActionCap int `json:"action_cap"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	FreshGatedPrevented int `json:"fresh_gated_prevented"`
	FreshGatedActions int `json:"fresh_gated_actions"`
	CommitmentPrevented int `json:"commitment_prevented"`
	CommitmentActions int `json:"commitment_actions"`
	CommitmentAccelerated int `json:"commitment_accelerated"`
	CommitmentMeanLossStepChangeAmongFailures float64 `json:"commitment_mean_loss_step_change_among_failures"`
}
type UP192CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP191CSeal string `json:"source_up191c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	ActionCaps []int `json:"action_caps"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	PolicySpecificRuleUsed bool `json:"policy_specific_rule_used"`
	FuturePolicyScheduleUsedByFirstWarning bool `json:"future_policy_schedule_used_by_first_warning"`
	Metrics []UP192CMetric `json:"metrics"`
}
func up192cRun(x0 *up81cAging,endangered,cadence int,policy string,cap int)(lossStep,actions int){
	m:=*x0
	committed:=false
	for start:=1;start<=64;start+=cadence{
		if actions<cap{
			if !committed{
				h:=up161cAdversarial(&m,endangered)
				if h<=cadence{
					m.query(endangered);actions++;committed=true
				}
			}else{
				m.query(endangered);actions++
			}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1016000+step,step,policy){return step,actions}
		}
	}
	return 65,actions
}
func RunUP192C()(UP192CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"}
	cadences:=[]int{2,4}
	caps:=[]int{2,4,6,8}
	res:=UP192CResult{Schema:UP192CCommitmentSchema,Experiment:"UP-192C-postwarning-commitment",SourceUP191CSeal:"7cd014c6e0eab81cd9fc2f73d4cd14941a2edcc6",Policies:policies,Cadences:cadences,ActionCaps:caps,CounterfactualOnly:true,LiveActivation:false,WarningThresholdChanged:false,PolicySpecificRuleUsed:false,FuturePolicyScheduleUsedByFirstWarning:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,cap:=range caps{
		m:=UP192CMetric{Policy:policy,Cadence:cad,ActionCap:cap};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,policy)
			fresh,freshActions:=up191cRun(x,e,cad,policy,cap)
			commit,commitActions:=up192cRun(x,e,cad,policy,cap)
			if base<=64{m.BaselineLosses++}
			if base<=64&&fresh>64{m.FreshGatedPrevented++}
			if base<=64&&commit>64{m.CommitmentPrevented++}
			m.FreshGatedActions+=freshActions;m.CommitmentActions+=commitActions
			if commit<base{m.CommitmentAccelerated++}
			if commit<=64{delaySum+=commit-base;delayN++}
		}}
		if delayN>0{m.CommitmentMeanLossStepChangeAmongFailures=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
