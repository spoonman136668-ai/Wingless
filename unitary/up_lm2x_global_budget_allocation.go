package unitary

import "fmt"

const UPLM2XAllocationSchema="wingless.up-lm2x-global-budget-allocation.v1"

type uplm2xArm struct{
	d int
	layout string
	r *uplm0cRecall
	reported map[string]bool
}
type UPLM2XPoint struct{
	IdentityRotation int `json:"identity_rotation"`
	Policy string `json:"policy"`
	Actions int `json:"actions"`
	Completed int `json:"completed"`
	Failed int `json:"failed"`
}
type UPLM2XResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2WSeal string `json:"source_up_lm2w_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	IdentityRotations []int `json:"identity_rotations"`
	ArmsPerScenario int `json:"arms_per_scenario"`
	GlobalRounds int `json:"global_rounds"`
	MaxActionsPerRound int `json:"max_actions_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Points []UPLM2XPoint `json:"points"`
}
func uplm2xArms(rot int)[]uplm2xArm{
	out:=[]uplm2xArm{}
	for _,d:=range []int{4,5,6}{for _,layout:=range []string{"suffix_reported","alternating_reported"}{
		r,rep:=uplm2vInit(d,rot,0,layout);out=append(out,uplm2xArm{d:d,layout:layout,r:r,reported:rep})
	}}
	return out
}
func uplm2xFirstPending(a *uplm2xArm)(string,int,bool){
	free:=16-len(a.r.order)
	for i,n:=range a.r.order{if !a.reported[n]{return n,free+i+1,true}}
	return "",0,false
}
func uplm2xChoose(arms []uplm2xArm,policy string)int{
	best:=-1
	if policy=="fixed_order"{for i:=range arms{if _,_,ok:=uplm2xFirstPending(&arms[i]);ok{return i}};return -1}
	bestCd:=1<<30
	for i:=range arms{_,cd,ok:=uplm2xFirstPending(&arms[i]);if ok&&(cd<bestCd||(cd==bestCd&&(best<0||i<best))){best=i;bestCd=cd}}
	return best
}
func uplm2xRun(rot int,policy string)(completed,failed,actions int){
	arms:=uplm2xArms(rot)
	for round:=0;round<12;round++{
		if policy!="baseline"{
			i:=uplm2xChoose(arms,policy)
			if i>=0{if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++}}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("x-%s-%d-%d-%d",policy,rot,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM2X()(UPLM2XResult,error){
	rots:=[]int{3,11};policies:=[]string{"baseline","earliest_deadline","fixed_order"}
	res:=UPLM2XResult{Schema:UPLM2XAllocationSchema,Experiment:"UP-LM2X-global-budget-allocation",SourceUPLM2WSeal:"4941bb92181dd47a4917e67efa1ebc92180ae3d1",ExactRecallCap:16,IdentityRotations:rots,ArmsPerScenario:6,GlobalRounds:12,MaxActionsPerRound:1,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false}
	for _,rot:=range rots{for _,policy:=range policies{c,f,a:=uplm2xRun(rot,policy);res.Points=append(res.Points,UPLM2XPoint{IdentityRotation:rot,Policy:policy,Actions:a,Completed:c,Failed:f})}}
	return res,nil
}
