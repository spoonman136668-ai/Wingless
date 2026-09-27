package unitary

import "fmt"

const UPLM3AThroughputSchema="wingless.up-lm3a-throughput-scarcity.v1"

type UPLM3APoint struct{
	Budget int `json:"budget"`
	Throughput int `json:"throughput"`
	IdentityRotation int `json:"identity_rotation"`
	Permutation string `json:"permutation"`
	Policy string `json:"policy"`
	Actions int `json:"actions"`
	Completed int `json:"completed"`
	Failed int `json:"failed"`
}
type UPLM3AMetric struct{
	Budget int `json:"budget"`
	Throughput int `json:"throughput"`
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
}
type UPLM3AResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2ZSeal string `json:"source_up_lm2z_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Budgets []int `json:"budgets"`
	Throughputs []int `json:"throughputs"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	GlobalRounds int `json:"global_rounds"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	AdaptiveBudgetUsed bool `json:"adaptive_budget_used"`
	AdaptiveThroughputUsed bool `json:"adaptive_throughput_used"`
	Points []UPLM3APoint `json:"points"`
	Metrics []UPLM3AMetric `json:"metrics"`
}
func uplm3aRun(rot int,perm,policy string,budget,throughput int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	for round:=0;round<12;round++{
		if policy!="baseline"{
			for k:=0;k<throughput && actions<budget;k++{
				i:=uplm2xChoose(arms,policy)
				if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3a-%s-%s-%d-%d-%d-%d-%d",policy,perm,budget,throughput,rot,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM3A()(UPLM3AResult,error){
	rots:=[]int{3,11};perms:=[]string{"identity","reverse","rotate2"};budgets:=[]int{5,6};throughputs:=[]int{1,2,3}
	res:=UPLM3AResult{Schema:UPLM3AThroughputSchema,Experiment:"UP-LM3A-throughput-scarcity",SourceUPLM2ZSeal:"8f94f6a8a8a387115ae2c9da7f29182989457e09",ExactRecallCap:16,Budgets:budgets,Throughputs:throughputs,IdentityRotations:rots,Permutations:perms,GlobalRounds:12,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false,AdaptiveBudgetUsed:false,AdaptiveThroughputUsed:false}
	for _,budget:=range budgets{for _,tp:=range throughputs{
		m:=UPLM3AMetric{Budget:budget,Throughput:tp}
		for _,rot:=range rots{for _,perm:=range perms{
			bc,bf,_:=uplm3aRun(rot,perm,"baseline",0,tp)
			ec,ef,ea:=uplm3aRun(rot,perm,"earliest_deadline",budget,tp)
			fc,ff,fa:=uplm3aRun(rot,perm,"fixed_order",budget,tp)
			res.Points=append(res.Points,
				UPLM3APoint{Budget:budget,Throughput:tp,IdentityRotation:rot,Permutation:perm,Policy:"baseline",Actions:0,Completed:bc,Failed:bf},
				UPLM3APoint{Budget:budget,Throughput:tp,IdentityRotation:rot,Permutation:perm,Policy:"earliest_deadline",Actions:ea,Completed:ec,Failed:ef},
				UPLM3APoint{Budget:budget,Throughput:tp,IdentityRotation:rot,Permutation:perm,Policy:"fixed_order",Actions:fa,Completed:fc,Failed:ff})
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed
		m.EarliestPreventedPerAction=uplm2sRate(m.EarliestPrevented,m.EarliestActions);m.FixedPreventedPerAction=uplm2sRate(m.FixedPrevented,m.FixedActions)
		m.EarliestAdvantage=m.FixedFailed-m.EarliestFailed
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
