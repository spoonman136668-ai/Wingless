package unitary

import "sort"

const UPLM3GHeldoutTopologySchema="wingless.up-lm3g-heldout-topology.v1"

type UPLM3GMetric struct{
	Profile string `json:"profile"`
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	BaselineFailed int `json:"baseline_failed"`
	EarliestFailed int `json:"earliest_failed"`
	FixedFailed int `json:"fixed_failed"`
	EarliestActions int `json:"earliest_actions"`
	FixedActions int `json:"fixed_actions"`
	EarliestPrevented int `json:"earliest_prevented"`
	FixedPrevented int `json:"fixed_prevented"`
}
type UPLM3GReachSummary struct{
	Profile string `json:"profile"`
	ReachableBudget int `json:"reachable_budget"`
	Observations int `json:"observations"`
	MinEarliestFailed int `json:"min_earliest_failed"`
	MaxEarliestFailed int `json:"max_earliest_failed"`
	FailureSpread int `json:"failure_spread"`
}
type UPLM3GResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3FSeal string `json:"source_up_lm3f_seal"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Profiles []string `json:"profiles"`
	Budgets []int `json:"budgets"`
	ActionStartRounds []int `json:"action_start_rounds"`
	Throughputs []int `json:"throughputs"`
	GlobalRounds int `json:"global_rounds"`
	DistinctArmPerRound bool `json:"distinct_arm_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveResourcesUsed bool `json:"adaptive_resources_used"`
	AdaptiveTopologySearchUsed bool `json:"adaptive_topology_search_used"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Metrics []UPLM3GMetric `json:"metrics"`
	ReachSummaries []UPLM3GReachSummary `json:"reach_summaries"`
}
func RunUPLM3G()(UPLM3GResult,error){
	profiles:=[]string{"by_deferred_level","by_layout"};budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3GResult{Schema:UPLM3GHeldoutTopologySchema,Experiment:"UP-LM3G-heldout-topology",SourceUPLM3FSeal:"062aa92cfa59ab09ce77b336ca788d96e76ed0dc",IdentityRotations:rots,Permutations:perms,Profiles:profiles,Budgets:budgets,ActionStartRounds:starts,Throughputs:tps,GlobalRounds:6,DistinctArmPerRound:true,CounterfactualOnly:true,LiveActivation:false,AdaptiveResourcesUsed:false,AdaptiveTopologySearchUsed:false,FutureScheduleUsed:false}
	type key struct{profile string;reach int};type acc struct{n,min,max int}
	groups:=map[key]*acc{}
	for _,profile:=range profiles{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		reach:=uplm3fReachable(budget,start,tp);m:=UPLM3GMetric{Profile:profile,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:reach}
		for _,rot:=range rots{for _,perm:=range perms{
			_,bf,_:=uplm3fRun(rot,perm,profile,"baseline",budget,start,tp)
			_,ef,ea:=uplm3fRun(rot,perm,profile,"earliest_deadline",budget,start,tp)
			_,ff,fa:=uplm3fRun(rot,perm,profile,"fixed_order",budget,start,tp)
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed
		res.Metrics=append(res.Metrics,m)
		k:=key{profile,reach};a:=groups[k];if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};groups[k]=a};a.n++;if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
	}}}}
	keys:=make([]key,0,len(groups));for k:=range groups{keys=append(keys,k)}
	sort.Slice(keys,func(i,j int)bool{if keys[i].profile!=keys[j].profile{return keys[i].profile<keys[j].profile};return keys[i].reach<keys[j].reach})
	for _,k:=range keys{a:=groups[k];res.ReachSummaries=append(res.ReachSummaries,UPLM3GReachSummary{Profile:k.profile,ReachableBudget:k.reach,Observations:a.n,MinEarliestFailed:a.min,MaxEarliestFailed:a.max,FailureSpread:a.max-a.min})}
	return res,nil
}
