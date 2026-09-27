package unitary

import (
	"fmt"
	"sort"
)

const UPLM3HHybridDeadlineSchema="wingless.up-lm3h-hybrid-deadline.v1"

type UPLM3HMetric struct{
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
type UPLM3HReachSummary struct{
	ReachableBudget int `json:"reachable_budget"`
	Observations int `json:"observations"`
	MinEarliestFailed int `json:"min_earliest_failed"`
	MaxEarliestFailed int `json:"max_earliest_failed"`
	FailureSpread int `json:"failure_spread"`
}
type UPLM3HResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3GSeal string `json:"source_up_lm3g_seal"`
	DeadlineProfile string `json:"deadline_profile"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Budgets []int `json:"budgets"`
	ActionStartRounds []int `json:"action_start_rounds"`
	Throughputs []int `json:"throughputs"`
	GlobalRounds int `json:"global_rounds"`
	DistinctArmPerRound bool `json:"distinct_arm_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveDeadlineUsed bool `json:"adaptive_deadline_used"`
	AdaptiveResourcesUsed bool `json:"adaptive_resources_used"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Metrics []UPLM3HMetric `json:"metrics"`
	ReachSummaries []UPLM3HReachSummary `json:"reach_summaries"`
}
func uplm3hDeferredTarget(a uplm2xArm)int{
	if a.d==4{return 1}
	if a.d==5{return 2}
	return 3
}
func uplm3hLayoutTarget(a uplm2xArm)int{
	if a.layout=="suffix_reported"{return 1}
	return 2
}
func uplm3hTarget(a uplm2xArm)int{
	d,l:=uplm3hDeferredTarget(a),uplm3hLayoutTarget(a)
	if d<l{return d}
	return l
}
func uplm3hPrepressure(arms []uplm2xArm,rot int){
	for i:=range arms{
		target:=uplm3hTarget(arms[i])
		for n:=0;n<32;n++{
			_,cd,ok:=uplm2xFirstPending(&arms[i]);if !ok||cd<=target{break}
			arms[i].r.write(fmt.Sprintf("3h-pre-%d-%d-%d",rot,i,n),"x")
		}
	}
}
func uplm3hRun(rot int,perm,policy string,budget,start,tp int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3hPrepressure(arms,rot)
	for round:=0;round<6;round++{
		if policy!="baseline"&&round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3h-%s-%s-%d-%d-%d-%d-%d",policy,perm,budget,start,tp,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM3H()(UPLM3HResult,error){
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3HResult{Schema:UPLM3HHybridDeadlineSchema,Experiment:"UP-LM3H-hybrid-deadline",SourceUPLM3GSeal:"afe5ac7f16beb36eac779ca77654c1cba5b8c911",DeadlineProfile:"hybrid_min",IdentityRotations:rots,Permutations:perms,Budgets:budgets,ActionStartRounds:starts,Throughputs:tps,GlobalRounds:6,DistinctArmPerRound:true,CounterfactualOnly:true,LiveActivation:false,AdaptiveDeadlineUsed:false,AdaptiveResourcesUsed:false,FutureScheduleUsed:false}
	type acc struct{n,min,max int};groups:=map[int]*acc{}
	for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		reach:=uplm3fReachable(budget,start,tp);m:=UPLM3HMetric{Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:reach}
		for _,rot:=range rots{for _,perm:=range perms{
			_,bf,_:=uplm3hRun(rot,perm,"baseline",budget,start,tp)
			_,ef,ea:=uplm3hRun(rot,perm,"earliest_deadline",budget,start,tp)
			_,ff,fa:=uplm3hRun(rot,perm,"fixed_order",budget,start,tp)
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed
		res.Metrics=append(res.Metrics,m)
		a:=groups[reach];if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};groups[reach]=a};a.n++;if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
	}}}
	keys:=make([]int,0,len(groups));for k:=range groups{keys=append(keys,k)};sort.Ints(keys)
	for _,k:=range keys{a:=groups[k];res.ReachSummaries=append(res.ReachSummaries,UPLM3HReachSummary{ReachableBudget:k,Observations:a.n,MinEarliestFailed:a.min,MaxEarliestFailed:a.max,FailureSpread:a.max-a.min})}
	return res,nil
}
