package unitary

import "fmt"

const UPLM2YReplicationSchema="wingless.up-lm2y-budget-allocation-replication.v1"

type UPLM2YPoint struct{
	IdentityRotation int `json:"identity_rotation"`
	Permutation string `json:"permutation"`
	Budget int `json:"budget"`
	Policy string `json:"policy"`
	Actions int `json:"actions"`
	Completed int `json:"completed"`
	Failed int `json:"failed"`
}
type UPLM2YBudgetMetric struct{
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
}
type UPLM2YResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2XSeal string `json:"source_up_lm2x_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Budgets []int `json:"budgets"`
	GlobalRounds int `json:"global_rounds"`
	MaxActionsPerRound int `json:"max_actions_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Points []UPLM2YPoint `json:"points"`
	Metrics []UPLM2YBudgetMetric `json:"metrics"`
}
func uplm2yPermute(in []uplm2xArm,perm string)[]uplm2xArm{
	n:=len(in);out:=make([]uplm2xArm,n)
	switch perm{
	case "reverse": for i:=0;i<n;i++{out[i]=in[n-1-i]}
	case "rotate2": for i:=0;i<n;i++{out[i]=in[(i+2)%n]}
	default: copy(out,in)
	}
	return out
}
func uplm2yRun(rot int,perm,policy string,budget int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	for round:=0;round<12;round++{
		if policy!="baseline"&&actions<budget{
			i:=uplm2xChoose(arms,policy)
			if i>=0{if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++}}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("y-%s-%s-%d-%d-%d-%d",policy,perm,budget,rot,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM2Y()(UPLM2YResult,error){
	rots:=[]int{3,11};perms:=[]string{"identity","reverse","rotate2"};budgets:=[]int{4,8,12}
	res:=UPLM2YResult{Schema:UPLM2YReplicationSchema,Experiment:"UP-LM2Y-budget-allocation-replication",SourceUPLM2XSeal:"79a4c00a527e037ac4209c81ecb7b9f49eb44cdc",ExactRecallCap:16,IdentityRotations:rots,Permutations:perms,Budgets:budgets,GlobalRounds:12,MaxActionsPerRound:1,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false}
	for _,budget:=range budgets{
		m:=UPLM2YBudgetMetric{Budget:budget}
		for _,rot:=range rots{for _,perm:=range perms{
			bc,bf,_:=uplm2yRun(rot,perm,"baseline",0)
			ec,ef,ea:=uplm2yRun(rot,perm,"earliest_deadline",budget)
			fc,ff,fa:=uplm2yRun(rot,perm,"fixed_order",budget)
			res.Points=append(res.Points,
				UPLM2YPoint{IdentityRotation:rot,Permutation:perm,Budget:budget,Policy:"baseline",Actions:0,Completed:bc,Failed:bf},
				UPLM2YPoint{IdentityRotation:rot,Permutation:perm,Budget:budget,Policy:"earliest_deadline",Actions:ea,Completed:ec,Failed:ef},
				UPLM2YPoint{IdentityRotation:rot,Permutation:perm,Budget:budget,Policy:"fixed_order",Actions:fa,Completed:fc,Failed:ff})
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed
		m.EarliestPreventedPerAction=uplm2sRate(m.EarliestPrevented,m.EarliestActions);m.FixedPreventedPerAction=uplm2sRate(m.FixedPrevented,m.FixedActions)
		m.EarliestAdvantage=m.FixedFailed-m.EarliestFailed
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
