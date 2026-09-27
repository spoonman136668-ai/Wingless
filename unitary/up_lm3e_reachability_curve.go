package unitary

import "fmt"

const UPLM3EReachabilitySchema="wingless.up-lm3e-reachability-curve.v1"

type UPLM3EPoint struct{
	Profile string `json:"profile"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	IdentityRotation int `json:"identity_rotation"`
	Permutation string `json:"permutation"`
	Policy string `json:"policy"`
	Actions int `json:"actions"`
	Completed int `json:"completed"`
	Failed int `json:"failed"`
}
type UPLM3EMetric struct{
	Profile string `json:"profile"`
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
	EarliestAdvantage int `json:"earliest_advantage"`
}
type UPLM3EResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3DSeal string `json:"source_up_lm3d_seal"`
	Profiles []string `json:"profiles"`
	ActionStartRounds []int `json:"action_start_rounds"`
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
	Points []UPLM3EPoint `json:"points"`
	Metrics []UPLM3EMetric `json:"metrics"`
}
func uplm3eReachable(start,tp int)int{n:=tp*(6-start);if n>6{n=6};if n<0{n=0};return n}
func uplm3eRun(rot int,perm,profile,policy string,start,tp int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3bPrepressure(arms,profile,rot);budget:=6
	for round:=0;round<6;round++{
		if policy!="baseline"&&round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3e-%s-%s-%s-%d-%d-%d-%d",profile,policy,perm,start,tp,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM3E()(UPLM3EResult,error){
	profiles:=[]string{"by_deferred_level","by_layout"};starts:=[]int{0,1,2,3,4};tps:=[]int{1,2,3};rots:=[]int{3,11};perms:=[]string{"identity","reverse"}
	res:=UPLM3EResult{Schema:UPLM3EReachabilitySchema,Experiment:"UP-LM3E-reachability-curve",SourceUPLM3DSeal:"7b6759c78d0faceb7ba59f5f3adddf9f3b24186d",Profiles:profiles,ActionStartRounds:starts,Throughputs:tps,TotalBudget:6,GlobalRounds:6,DistinctArmPerRound:true,CounterfactualOnly:true,LiveActivation:false,AdaptiveOnsetUsed:false,AdaptiveBudgetUsed:false,AdaptiveThroughputUsed:false,FutureScheduleUsed:false}
	for _,profile:=range profiles{for _,start:=range starts{for _,tp:=range tps{
		reach:=uplm3eReachable(start,tp);m:=UPLM3EMetric{Profile:profile,ActionStartRound:start,Throughput:tp,ReachableBudget:reach}
		for _,rot:=range rots{for _,perm:=range perms{
			bc,bf,_:=uplm3eRun(rot,perm,profile,"baseline",start,tp)
			ec,ef,ea:=uplm3eRun(rot,perm,profile,"earliest_deadline",start,tp)
			fc,ff,fa:=uplm3eRun(rot,perm,profile,"fixed_order",start,tp)
			res.Points=append(res.Points,
				UPLM3EPoint{Profile:profile,ActionStartRound:start,Throughput:tp,ReachableBudget:reach,IdentityRotation:rot,Permutation:perm,Policy:"baseline",Completed:bc,Failed:bf},
				UPLM3EPoint{Profile:profile,ActionStartRound:start,Throughput:tp,ReachableBudget:reach,IdentityRotation:rot,Permutation:perm,Policy:"earliest_deadline",Actions:ea,Completed:ec,Failed:ef},
				UPLM3EPoint{Profile:profile,ActionStartRound:start,Throughput:tp,ReachableBudget:reach,IdentityRotation:rot,Permutation:perm,Policy:"fixed_order",Actions:fa,Completed:fc,Failed:ff})
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed;m.EarliestAdvantage=m.FixedFailed-m.EarliestFailed
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
