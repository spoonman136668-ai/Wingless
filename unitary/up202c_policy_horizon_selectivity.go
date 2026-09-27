package unitary

const UP202CPolicySelectivitySchema="wingless.up202c-policy-horizon-selectivity.v1"

type UP202CMetric struct{
	Policy string `json:"policy"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	BaselineFailures int `json:"baseline_failures"`
	BaselineSurvivors int `json:"baseline_survivors"`
	TreatedFailures int `json:"treated_failures"`
	PreventedFailures int `json:"prevented_failures"`
	ActionsTaken int `json:"actions_taken"`
	ApparentUnnecessaryActionArms int `json:"apparent_unnecessary_action_arms"`
	NearFutureActionArms int `json:"near_future_action_arms"`
	TrueUnnecessaryActionArms int `json:"true_unnecessary_action_arms"`
	BaselineSurvivorsBeyondGrace int `json:"baseline_survivors_beyond_grace"`
	PreventedPerTrueUnnecessary float64 `json:"prevented_per_true_unnecessary"`
	TrueUnnecessaryFractionAmongBeyondGraceSurvivors float64 `json:"true_unnecessary_fraction_among_beyond_grace_survivors"`
}
type UP202CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP201CSeal string `json:"source_up201c_seal"`
	Policies []string `json:"policies"`
	Horizons []int `json:"horizons"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	GraceWrites int `json:"grace_writes"`
	CohortTopology string `json:"cohort_topology"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	PolicySpecificRuleUsed bool `json:"policy_specific_rule_used"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP202CMetric `json:"metrics"`
}
func RunUP202C()(UP202CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"}
	horizons:=[]int{8,12,16,20,24,28,32,40,48,56}
	res:=UP202CResult{Schema:UP202CPolicySelectivitySchema,Experiment:"UP-202C-policy-horizon-selectivity",SourceUP201CSeal:"eb199a3866190a78f20279f7cbcf892393ebdc9c",Policies:policies,Horizons:horizons,Cadence:2,SpacingIntervals:4,ActionCap:8,GraceWrites:8,CohortTopology:"contiguous_quartets",CounterfactualOnly:true,LiveActivation:false,PolicySpecificRuleUsed:false,AdaptiveSelectionUsed:false}
	for _,policy:=range policies{for _,horizon:=range horizons{
		m:=UP202CMetric{Policy:policy,Horizon:horizon}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,2,policy)
			treated,actions:=up201cRun(x,e,horizon,policy)
			if base<=horizon{m.BaselineFailures++}else{m.BaselineSurvivors++}
			if treated<=horizon{m.TreatedFailures++}
			if base<=horizon&&treated>horizon{m.PreventedFailures++}
			if base>horizon+8{m.BaselineSurvivorsBeyondGrace++}
			if base>horizon&&actions>0{
				m.ApparentUnnecessaryActionArms++
				if base<=horizon+8{m.NearFutureActionArms++}else{m.TrueUnnecessaryActionArms++}
			}
			m.ActionsTaken+=actions
		}}
		if m.TrueUnnecessaryActionArms>0{m.PreventedPerTrueUnnecessary=float64(m.PreventedFailures)/float64(m.TrueUnnecessaryActionArms)}
		if m.BaselineSurvivorsBeyondGrace>0{m.TrueUnnecessaryFractionAmongBeyondGraceSurvivors=float64(m.TrueUnnecessaryActionArms)/float64(m.BaselineSurvivorsBeyondGrace)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
