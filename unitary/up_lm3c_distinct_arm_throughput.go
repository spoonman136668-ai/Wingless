package unitary

import "fmt"

const UPLM3CDistinctArmSchema="wingless.up-lm3c-distinct-arm-throughput.v1"

type UPLM3CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3BSeal string `json:"source_up_lm3b_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Profiles []string `json:"profiles"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Throughputs []int `json:"throughputs"`
	TotalBudget int `json:"total_budget"`
	GlobalRounds int `json:"global_rounds"`
	OneActionPerArmPerRound bool `json:"one_action_per_arm_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Points []UPLM3BPoint `json:"points"`
	Metrics []UPLM3BMetric `json:"metrics"`
}

func uplm3cChoose(arms []uplm2xArm,policy string,used map[int]bool)int{
	if policy=="fixed_order"{
		for i:=range arms{
			if used[i]{continue}
			if _,_,ok:=uplm2xFirstPending(&arms[i]);ok{return i}
		}
		return -1
	}
	best:=-1
	bestCd:=1<<30
	for i:=range arms{
		if used[i]{continue}
		_,cd,ok:=uplm2xFirstPending(&arms[i])
		if ok&&(cd<bestCd||(cd==bestCd&&(best<0||i<best))){best=i;bestCd=cd}
	}
	return best
}

func uplm3cRun(rot int,perm,profile,policy string,throughput int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	uplm3bPrepressure(arms,profile,rot)
	budget:=6
	for round:=0;round<6;round++{
		if policy!="baseline"{
			used:=map[int]bool{}
			for k:=0;k<throughput&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used)
				if i<0{break}
				used[i]=true
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3c-%s-%s-%s-%d-%d-%d-%d",profile,policy,perm,throughput,rot,round,i),"x")}
	}
	for i:=range arms{
		c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot)
		completed+=c
		failed+=f
	}
	return
}

func RunUPLM3C()(UPLM3CResult,error){
	profiles:=[]string{"by_deferred_level","by_layout"}
	rots:=[]int{3,11}
	perms:=[]string{"identity","reverse"}
	tps:=[]int{1,2,3}
	res:=UPLM3CResult{
		Schema:UPLM3CDistinctArmSchema,
		Experiment:"UP-LM3C-distinct-arm-throughput",
		SourceUPLM3BSeal:"eb17353a337835e67ff68a28ed5d4c3af7f7e9a3",
		ExactRecallCap:16,
		Profiles:profiles,
		IdentityRotations:rots,
		Permutations:perms,
		Throughputs:tps,
		TotalBudget:6,
		GlobalRounds:6,
		OneActionPerArmPerRound:true,
		CounterfactualOnly:true,
		LiveActivation:false,
	}
	for _,profile:=range profiles{
		for _,tp:=range tps{
			m:=UPLM3BMetric{Profile:profile,Throughput:tp}
			for _,rot:=range rots{
				for _,perm:=range perms{
					bc,bf,_:=uplm3cRun(rot,perm,profile,"baseline",tp)
					ec,ef,ea:=uplm3cRun(rot,perm,profile,"earliest_deadline",tp)
					fc,ff,fa:=uplm3cRun(rot,perm,profile,"fixed_order",tp)
					res.Points=append(res.Points,
						UPLM3BPoint{Profile:profile,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"baseline",Completed:bc,Failed:bf},
						UPLM3BPoint{Profile:profile,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"earliest_deadline",Actions:ea,Completed:ec,Failed:ef},
						UPLM3BPoint{Profile:profile,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"fixed_order",Actions:fa,Completed:fc,Failed:ff})
					m.BaselineFailed+=bf
					m.EarliestFailed+=ef
					m.FixedFailed+=ff
					m.EarliestActions+=ea
					m.FixedActions+=fa
				}
			}
			m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed
			m.FixedPrevented=m.BaselineFailed-m.FixedFailed
			m.EarliestAdvantage=m.FixedFailed-m.EarliestFailed
			res.Metrics=append(res.Metrics,m)
		}
	}
	return res,nil
}
