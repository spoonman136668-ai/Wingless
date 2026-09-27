package unitary

import "fmt"

const UPLM3BBurstySchema="wingless.up-lm3b-bursty-throughput.v1"

type UPLM3BPoint struct{
	Profile string `json:"profile"`
	IdentityRotation int `json:"identity_rotation"`
	Permutation string `json:"permutation"`
	Throughput int `json:"throughput"`
	Policy string `json:"policy"`
	Actions int `json:"actions"`
	Completed int `json:"completed"`
	Failed int `json:"failed"`
}
type UPLM3BMetric struct{
	Profile string `json:"profile"`
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
type UPLM3BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3ASeal string `json:"source_up_lm3a_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Profiles []string `json:"profiles"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Throughputs []int `json:"throughputs"`
	TotalBudget int `json:"total_budget"`
	GlobalRounds int `json:"global_rounds"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Points []UPLM3BPoint `json:"points"`
	Metrics []UPLM3BMetric `json:"metrics"`
}
func uplm3bTarget(a uplm2xArm,profile string)int{
	if profile=="by_deferred_level"{
		if a.d==4{return 1};if a.d==5{return 2};return 3
	}
	if a.layout=="suffix_reported"{return 1}
	return 2
}
func uplm3bPrepressure(arms []uplm2xArm,profile string,rot int){
	for i:=range arms{
		target:=uplm3bTarget(arms[i],profile)
		for n:=0;n<32;n++{
			_,cd,ok:=uplm2xFirstPending(&arms[i]);if !ok||cd<=target{break}
			arms[i].r.write(fmt.Sprintf("3b-pre-%s-%d-%d-%d",profile,rot,i,n),"x")
		}
	}
}
func uplm3bRun(rot int,perm,profile,policy string,throughput int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3bPrepressure(arms,profile,rot)
	budget:=6
	for round:=0;round<6;round++{
		if policy!="baseline"{
			for k:=0;k<throughput&&actions<budget;k++{
				i:=uplm2xChoose(arms,policy);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3b-%s-%s-%s-%d-%d-%d-%d",profile,policy,perm,throughput,rot,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM3B()(UPLM3BResult,error){
	profiles:=[]string{"by_deferred_level","by_layout"};rots:=[]int{3,11};perms:=[]string{"identity","reverse"};tps:=[]int{1,2,3}
	res:=UPLM3BResult{Schema:UPLM3BBurstySchema,Experiment:"UP-LM3B-bursty-throughput",SourceUPLM3ASeal:"3798bd40e2ad514f0be6b62d7236b3ecbcaf178b",ExactRecallCap:16,Profiles:profiles,IdentityRotations:rots,Permutations:perms,Throughputs:tps,TotalBudget:6,GlobalRounds:6,CounterfactualOnly:true,LiveActivation:false}
	for _,profile:=range profiles{for _,tp:=range tps{
		m:=UPLM3BMetric{Profile:profile,Throughput:tp}
		for _,rot:=range rots{for _,perm:=range perms{
			bc,bf,_:=uplm3bRun(rot,perm,profile,"baseline",tp)
			ec,ef,ea:=uplm3bRun(rot,perm,profile,"earliest_deadline",tp)
			fc,ff,fa:=uplm3bRun(rot,perm,profile,"fixed_order",tp)
			res.Points=append(res.Points,
				UPLM3BPoint{Profile:profile,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"baseline",Completed:bc,Failed:bf},
				UPLM3BPoint{Profile:profile,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"earliest_deadline",Actions:ea,Completed:ec,Failed:ef},
				UPLM3BPoint{Profile:profile,IdentityRotation:rot,Permutation:perm,Throughput:tp,Policy:"fixed_order",Actions:fa,Completed:fc,Failed:ff})
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		m.EarliestPrevented=m.BaselineFailed-m.EarliestFailed;m.FixedPrevented=m.BaselineFailed-m.FixedFailed;m.EarliestAdvantage=m.FixedFailed-m.EarliestFailed
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
