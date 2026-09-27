package unitary

const UP201CPolicyTransferSchema="wingless.up201c-sparse-schedule-policy-transfer.v1"

type UP201CMetric struct{
	Policy string `json:"policy"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	ArmsWithAction int `json:"arms_with_action"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
}
type UP201CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP200CSeal string `json:"source_up200c_seal"`
	Policies []string `json:"policies"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Horizon int `json:"horizon"`
	CohortTopology string `json:"cohort_topology"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	PolicySpecificRuleUsed bool `json:"policy_specific_rule_used"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP201CMetric `json:"metrics"`
}
func up201cRun(x0 *up81cAging,endangered,horizon int,policy string)(lossStep,actions int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;committed=true;since=0
			}
		}else if actions<8{
			since++
			if since>=4{m.query(endangered);actions++;since=0}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1023000+step,step,policy){return step,actions}
		}
	}
	return horizon+1,actions
}
func RunUP201C()(UP201CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"}
	const horizon=56
	res:=UP201CResult{Schema:UP201CPolicyTransferSchema,Experiment:"UP-201C-sparse-schedule-policy-transfer",SourceUP200CSeal:"70915a195815f578ebae2f62ed2e2e1823ee1195",Policies:policies,Cadence:2,SpacingIntervals:4,ActionCap:8,Horizon:horizon,CohortTopology:"contiguous_quartets",CounterfactualOnly:true,LiveActivation:false,PolicySpecificRuleUsed:false,AdaptiveSelectionUsed:false}
	for _,policy:=range policies{
		m:=UP201CMetric{Policy:policy};sumExt,nFail:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,2,policy)
			treated,actions:=up201cRun(x,e,horizon,policy)
			if base<=horizon{m.BaselineLosses++}
			if treated<=horizon{m.TreatedLosses++;sumExt+=treated-base;nFail++}
			if base<=horizon&&treated>horizon{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			if actions>0{m.ArmsWithAction++}
			m.ActionsTaken+=actions
		}}
		if nFail>0{m.MeanLossStepExtensionAmongFailures=float64(sumExt)/float64(nFail)}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
