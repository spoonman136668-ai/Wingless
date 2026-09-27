package unitary

const UPLM2ZOnsetSchema="wingless.up-lm2z-budget-onset.v1"

type UPLM2ZScenario struct{
	Budget int `json:"budget"`
	IdentityRotation int `json:"identity_rotation"`
	Permutation string `json:"permutation"`
	EarliestFailed int `json:"earliest_failed"`
	FixedFailed int `json:"fixed_failed"`
	EarliestAdvantage int `json:"earliest_advantage"`
}
type UPLM2ZMetric struct{
	Budget int `json:"budget"`
	BaselineFailed int `json:"baseline_failed"`
	EarliestFailed int `json:"earliest_failed"`
	FixedFailed int `json:"fixed_failed"`
	EarliestActions int `json:"earliest_actions"`
	FixedActions int `json:"fixed_actions"`
	EarliestPrevented int `json:"earliest_prevented"`
	FixedPrevented int `json:"fixed_prevented"`
	EarliestPreventedPerAction float64 `json:"earliest_prevented_per_action"`
	FixedPreventedPerAction float64 `json:"fixed_prevented_per_action"`
	EarliestAdvantage int `json:"earliest_advantage"`
	ScenariosEarliestStrictlyBetter int `json:"scenarios_earliest_strictly_better"`
	ScenariosTied int `json:"scenarios_tied"`
}
type UPLM2ZResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2YSeal string `json:"source_up_lm2y_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Budgets []int `json:"budgets"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	GlobalRounds int `json:"global_rounds"`
	MaxActionsPerRound int `json:"max_actions_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	FirstAggregateAdvantageBudget int `json:"first_aggregate_advantage_budget"`
	Scenarios []UPLM2ZScenario `json:"scenarios"`
	Metrics []UPLM2ZMetric `json:"metrics"`
}
func RunUPLM2Z()(UPLM2ZResult,error){
	rots:=[]int{3,11};perms:=[]string{"identity","reverse","rotate2"};budgets:=[]int{4,5,6,7,8}
	res:=UPLM2ZResult{Schema:UPLM2ZOnsetSchema,Experiment:"UP-LM2Z-budget-onset",SourceUPLM2YSeal:"e955419fdb6a38573cf6007555ecd4583f80889b",ExactRecallCap:16,Budgets:budgets,IdentityRotations:rots,Permutations:perms,GlobalRounds:12,MaxActionsPerRound:1,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false,FirstAggregateAdvantageBudget:-1}
	for _,budget:=range budgets{
		m:=UPLM2ZMetric{Budget:budget}
		for _,rot:=range rots{for _,perm:=range perms{
			_,bf,_:=uplm2yRun(rot,perm,"baseline",0)
			_,ef,ea:=uplm2yRun(rot,perm,"earliest_deadline",budget)
			_,ff,fa:=uplm2yRun(rot,perm,"fixed_order",budget)
			adv:=ff-ef
			res.Scenarios=append(res.Scenarios,UPLM2ZScenario{Budget:budget,IdentityRotation:rot,Permutation:perm,EarliestFailed:ef,FixedFailed:ff,EarliestAdvantage:adv})
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
			if adv>0{m.ScenariosEarliestStrictlyBetter++}else if adv==0{m.ScenariosTied++}
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed
		m.EarliestPreventedPerAction=uplm2sRate(m.EarliestPrevented,m.EarliestActions);m.FixedPreventedPerAction=uplm2sRate(m.FixedPrevented,m.FixedActions)
		m.EarliestAdvantage=m.FixedFailed-m.EarliestFailed
		if res.FirstAggregateAdvantageBudget<0&&m.EarliestAdvantage>0{res.FirstAggregateAdvantageBudget=budget}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
