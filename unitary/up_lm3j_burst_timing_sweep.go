package unitary

import (
	"fmt"
	"sort"
)

const UPLM3JBurstTimingSchema="wingless.up-lm3j-burst-timing-sweep.v1"

type UPLM3JMetric struct{
	BurstRound int `json:"burst_round"`
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
type UPLM3JReachSummary struct{
	BurstRound int `json:"burst_round"`
	ReachableBudget int `json:"reachable_budget"`
	Observations int `json:"observations"`
	MinEarliestFailed int `json:"min_earliest_failed"`
	MaxEarliestFailed int `json:"max_earliest_failed"`
	FailureSpread int `json:"failure_spread"`
}
type UPLM3JResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3ISeal string `json:"source_up_lm3i_seal"`
	DeadlineProfile string `json:"deadline_profile"`
	BurstRounds []int `json:"burst_rounds"`
	ExtraWritesPerArm int `json:"extra_writes_per_arm"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Budgets []int `json:"budgets"`
	ActionStartRounds []int `json:"action_start_rounds"`
	Throughputs []int `json:"throughputs"`
	GlobalRounds int `json:"global_rounds"`
	DistinctArmPerRound bool `json:"distinct_arm_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveBurstUsed bool `json:"adaptive_burst_used"`
	AdaptiveResourcesUsed bool `json:"adaptive_resources_used"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Metrics []UPLM3JMetric `json:"metrics"`
	ReachSummaries []UPLM3JReachSummary `json:"reach_summaries"`
}
func uplm3jRun(rot int,perm,policy string,budget,start,tp,burstRound int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3hPrepressure(arms,rot)
	for round:=0;round<6;round++{
		if policy!="baseline"&&round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{
			arms[i].r.write(fmt.Sprintf("3j-%d-%s-%s-%d-%d-%d-%d-%d",burstRound,policy,perm,budget,start,tp,round,i),"x")
			if round==burstRound{arms[i].r.write(fmt.Sprintf("3j-burst-%d-%s-%s-%d-%d-%d-%d",burstRound,policy,perm,budget,start,tp,i),"x")}
		}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM3J()(UPLM3JResult,error){
	bursts:=[]int{0,1,2,3,4,5};budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3JResult{Schema:UPLM3JBurstTimingSchema,Experiment:"UP-LM3J-burst-timing-sweep",SourceUPLM3ISeal:"4fdaccbd98457f3cf9543ad50b65710ad8335e4b",DeadlineProfile:"hybrid_min",BurstRounds:bursts,ExtraWritesPerArm:1,IdentityRotations:rots,Permutations:perms,Budgets:budgets,ActionStartRounds:starts,Throughputs:tps,GlobalRounds:6,DistinctArmPerRound:true,CounterfactualOnly:true,LiveActivation:false,AdaptiveBurstUsed:false,AdaptiveResourcesUsed:false,FutureScheduleUsed:false}
	type key struct{burst,reach int};type acc struct{n,min,max int};groups:=map[key]*acc{}
	for _,burst:=range bursts{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		reach:=uplm3fReachable(budget,start,tp);m:=UPLM3JMetric{BurstRound:burst,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:reach}
		for _,rot:=range rots{for _,perm:=range perms{
			_,bf,_:=uplm3jRun(rot,perm,"baseline",budget,start,tp,burst)
			_,ef,ea:=uplm3jRun(rot,perm,"earliest_deadline",budget,start,tp,burst)
			_,ff,fa:=uplm3jRun(rot,perm,"fixed_order",budget,start,tp,burst)
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed
		res.Metrics=append(res.Metrics,m)
		k:=key{burst:burst,reach:reach};a:=groups[k];if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};groups[k]=a};a.n++;if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
	}}}}
	keys:=make([]key,0,len(groups));for k:=range groups{keys=append(keys,k)}
	sort.Slice(keys,func(i,j int)bool{if keys[i].burst!=keys[j].burst{return keys[i].burst<keys[j].burst};return keys[i].reach<keys[j].reach})
	for _,k:=range keys{a:=groups[k];res.ReachSummaries=append(res.ReachSummaries,UPLM3JReachSummary{BurstRound:k.burst,ReachableBudget:k.reach,Observations:a.n,MinEarliestFailed:a.min,MaxEarliestFailed:a.max,FailureSpread:a.max-a.min})}
	return res,nil
}
