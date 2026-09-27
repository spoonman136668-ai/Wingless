package unitary

const UPLM3EReachabilitySchema="wingless.up-lm3e-reachability-envelope.v1"

type UPLM3EMetric struct{
	Profile string `json:"profile"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	RemainingActionRounds int `json:"remaining_action_rounds"`
	ReachableActionsPerScenario int `json:"reachable_actions_per_scenario"`
	BaselineFailed int `json:"baseline_failed"`
	EarliestFailed int `json:"earliest_failed"`
	FixedFailed int `json:"fixed_failed"`
	EarliestActions int `json:"earliest_actions"`
	FixedActions int `json:"fixed_actions"`
	EarliestPrevented int `json:"earliest_prevented"`
	FixedPrevented int `json:"fixed_prevented"`
	EarliestAdvantage int `json:"earliest_advantage"`
}
type UPLM3EResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3DSeal string `json:"source_up_lm3d_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Profiles []string `json:"profiles"`
	ActionStartRounds []int `json:"action_start_rounds"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Throughputs []int `json:"throughputs"`
	TotalBudget int `json:"total_budget"`
	GlobalRounds int `json:"global_rounds"`
	DistinctArmPerRound bool `json:"distinct_arm_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveOnsetUsed bool `json:"adaptive_onset_used"`
	AdaptiveBudgetUsed bool `json:"adaptive_budget_used"`
	AdaptiveThroughputUsed bool `json:"adaptive_throughput_used"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Metrics []UPLM3EMetric `json:"metrics"`
}
func uplm3eReachable(start,tp int)int{
	remaining:=6-start
	if remaining<0{remaining=0}
	n:=remaining*tp
	if n>6{n=6}
	return n
}
func RunUPLM3E()(UPLM3EResult,error){
	profiles:=[]string{"by_deferred_level","by_layout"};starts:=[]int{0,1,2,3,4};rots:=[]int{3,11};perms:=[]string{"identity","reverse"};tps:=[]int{1,2,3}
	res:=UPLM3EResult{Schema:UPLM3EReachabilitySchema,Experiment:"UP-LM3E-reachability-envelope",SourceUPLM3DSeal:"7b6759c78d0faceb7ba59f5f3adddf9f3b24186d",ExactRecallCap:16,Profiles:profiles,ActionStartRounds:starts,IdentityRotations:rots,Permutations:perms,Throughputs:tps,TotalBudget:6,GlobalRounds:6,DistinctArmPerRound:true,CounterfactualOnly:true,LiveActivation:false,AdaptiveOnsetUsed:false,AdaptiveBudgetUsed:false,AdaptiveThroughputUsed:false,FutureScheduleUsed:false}
	for _,profile:=range profiles{for _,start:=range starts{for _,tp:=range tps{
		m:=UPLM3EMetric{Profile:profile,ActionStartRound:start,Throughput:tp,RemainingActionRounds:6-start,ReachableActionsPerScenario:uplm3eReachable(start,tp)}
		for _,rot:=range rots{for _,perm:=range perms{
			_,bf,_:=uplm3dRun(rot,perm,profile,"baseline",start,tp)
			_,ef,ea:=uplm3dRun(rot,perm,profile,"earliest_deadline",start,tp)
			_,ff,fa:=uplm3dRun(rot,perm,profile,"fixed_order",start,tp)
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed;m.EarliestAdvantage=m.FixedFailed-m.EarliestFailed
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
