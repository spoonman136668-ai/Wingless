package unitary

import "fmt"

const UPLM3FBudgetReachabilitySchema="wingless.up-lm3f-budget-reachability.v1"

type UPLM3FPoint struct{
	Profile string `json:"profile"`
	Budget int `json:"budget"`
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
type UPLM3FMetric struct{
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
type UPLM3FResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3ESeal string `json:"source_up_lm3e_seal"`
	Profiles []string `json:"profiles"`
	Budgets []int `json:"budgets"`
	ActionStartRounds []int `json:"action_start_rounds"`
	Throughputs []int `json:"throughputs"`
	GlobalRounds int `json:"global_rounds"`
	DistinctArmPerRound bool `json:"distinct_arm_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveResourcesUsed bool `json:"adaptive_resources_used"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Points []UPLM3FPoint `json:"points"`
	Metrics []UPLM3FMetric `json:"metrics"`
}
func uplm3fReachable(budget,start,tp int)int{n:=tp*(6-start);if n>budget{n=budget};if n<0{n=0};return n}
func uplm3fRun(rot int,perm,profile,policy string,budget,start,tp int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3bPrepressure(arms,profile,rot)
	for round:=0;round<6;round++{
		if policy!="baseline"&&round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3f-%s-%s-%s-%d-%d-%d-%d-%d",profile,policy,perm,budget,start,tp,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM3F()(UPLM3FResult,error){
	profiles:=[]string{"by_deferred_level","by_layout"};budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{3,11};perms:=[]string{"identity","reverse"}
	res:=UPLM3FResult{Schema:UPLM3FBudgetReachabilitySchema,Experiment:"UP-LM3F-budget-reachability",SourceUPLM3ESeal:"5f3267b1b9451d4a785d232fcc05cfb4dc90a10f",Profiles:profiles,Budgets:budgets,ActionStartRounds:starts,Throughputs:tps,GlobalRounds:6,DistinctArmPerRound:true,CounterfactualOnly:true,LiveActivation:false,AdaptiveResourcesUsed:false,FutureScheduleUsed:false}
	for _,profile:=range profiles{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		reach:=uplm3fReachable(budget,start,tp);m:=UPLM3FMetric{Profile:profile,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:reach}
		for _,rot:=range rots{for _,perm:=range perms{
			bc,bf,_:=uplm3fRun(rot,perm,profile,"baseline",budget,start,tp)
			ec,ef,ea:=uplm3fRun(rot,perm,profile,"earliest_deadline",budget,start,tp)
			fc,ff,fa:=uplm3fRun(rot,perm,profile,"fixed_order",budget,start,tp)
			res.Points=append(res.Points,
				UPLM3FPoint{Profile:profile,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:reach,IdentityRotation:rot,Permutation:perm,Policy:"baseline",Completed:bc,Failed:bf},
				UPLM3FPoint{Profile:profile,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:reach,IdentityRotation:rot,Permutation:perm,Policy:"earliest_deadline",Actions:ea,Completed:ec,Failed:ef},
				UPLM3FPoint{Profile:profile,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:reach,IdentityRotation:rot,Permutation:perm,Policy:"fixed_order",Actions:fa,Completed:fc,Failed:ff})
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed
		res.Metrics=append(res.Metrics,m)
	}}}}
	return res,nil
}
