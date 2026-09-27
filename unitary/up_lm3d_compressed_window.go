package unitary

import "fmt"

const UPLM3DCompressedWindowSchema="wingless.up-lm3d-compressed-window.v1"

type UPLM3DPoint struct{
	Profile string `json:"profile"`
	ActionStartRound int `json:"action_start_round"`
	IdentityRotation int `json:"identity_rotation"`
	Permutation string `json:"permutation"`
	Throughput int `json:"throughput"`
	Policy string `json:"policy"`
	Actions int `json:"actions"`
	Completed int `json:"completed"`
	Failed int `json:"failed"`
}
type UPLM3DMetric struct{
	Profile string `json:"profile"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	BaselineFailed int `json:"baseline_failed"`
	EarliestFailed int `json:"earliest_failed"`
	FixedFailed int `json:"fixed_failed"`
	EarliestActions int `json:"earliest_actions"`
	FixedActions int `json:"fixed_actions"`
	EarliestPrevented int `json:"earliest_prevented"`
	FixedPrevented int `json:"fixed_prevented"`
	EarliestAdvantage int `json:"earliest_advantage"`
}
type UPLM3DResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3CSeal string `json:"source_up_lm3c_seal"`
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
	Points []UPLM3DPoint `json:"points"`
	Metrics []UPLM3DMetric `json:"metrics"`
}
func uplm3dRun(rot int,perm,profile,policy string,startRound,throughput int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3bPrepressure(arms,profile,rot)
	budget:=6
	for round:=0;round<6;round++{
		if policy!="baseline"&&round>=startRound{
			used:=map[int]bool{}
			for k:=0;k<throughput&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3d-%s-%s-%s-%d-%d-%d-%d-%d",profile,policy,perm,startRound,throughput,rot,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM3D()(UPLM3DResult,error){
	profiles:=[]string{"by_deferred_level","by_layout"};starts:=[]int{1,2};rots:=[]int{3,11};perms:=[]string{"identity","reverse"};tps:=[]int{1,2,3}
	res:=UPLM3DResult{Schema:UPLM3DCompressedWindowSchema,Experiment:"UP-LM3D-compressed-window",SourceUPLM3CSeal:"bb29b01bf50424f94cc40412fd395059d3b0cc48",ExactRecallCap:16,Profiles:profiles,ActionStartRounds:starts,IdentityRotations:rots,Permutations:perms,Throughputs:tps,TotalBudget:6,GlobalRounds:6,DistinctArmPerRound:true,CounterfactualOnly:true,LiveActivation:false,AdaptiveOnsetUsed:false,AdaptiveBudgetUsed:false,AdaptiveThroughputUsed:false,FutureScheduleUsed:false}
	for _,profile:=range profiles{for _,start:=range starts{for _,tp:=range tps{
		m:=UPLM3DMetric{Profile:profile,ActionStartRound:start,Throughput:tp}
		for _,rot:=range rots{for _,perm:=range perms{
			bc,bf,_:=uplm3dRun(rot,perm,profile,"baseline",start,tp)
			ec,ef,ea:=uplm3dRun(rot,perm,profile,"earliest_deadline",start,tp)
			fc,ff,fa:=uplm3dRun(rot,perm,profile,"fixed_order",start,tp)
			res.Points=append(res.Points,
				UPLM3DPoint{Profile:profile,ActionStartRound:start,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"baseline",Completed:bc,Failed:bf},
				UPLM3DPoint{Profile:profile,ActionStartRound:start,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"earliest_deadline",Actions:ea,Completed:ec,Failed:ef},
				UPLM3DPoint{Profile:profile,ActionStartRound:start,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"fixed_order",Actions:fa,Completed:fc,Failed:ff})
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed;m.EarliestAdvantage=m.FixedFailed-m.EarliestFailed
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
