package unitary

import (
	"fmt"
	"sort"
)

const UPLM3KPressurePlacementSchema="wingless.up-lm3k-pressure-placement-control.v1"

type UPLM3KMetric struct{
	PressureMode string `json:"pressure_mode"`
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	BaselineFailed int `json:"baseline_failed"`
	EarliestFailed int `json:"earliest_failed"`
	FixedFailed int `json:"fixed_failed"`
	EarliestActions int `json:"earliest_actions"`
	FixedActions int `json:"fixed_actions"`
}
type UPLM3KSummary struct{
	PressureMode string `json:"pressure_mode"`
	ReachableBudget int `json:"reachable_budget"`
	Observations int `json:"observations"`
	MinEarliestFailed int `json:"min_earliest_failed"`
	MaxEarliestFailed int `json:"max_earliest_failed"`
	FailureSpread int `json:"failure_spread"`
}
type UPLM3KResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3JSeal string `json:"source_up_lm3j_seal"`
	PressureModes []string `json:"pressure_modes"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Budgets []int `json:"budgets"`
	ActionStartRounds []int `json:"action_start_rounds"`
	Throughputs []int `json:"throughputs"`
	GlobalRounds int `json:"global_rounds"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptivePressureUsed bool `json:"adaptive_pressure_used"`
	AdaptiveResourcesUsed bool `json:"adaptive_resources_used"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Metrics []UPLM3KMetric `json:"metrics"`
	Summaries []UPLM3KSummary `json:"summaries"`
}
func uplm3kRun(rot int,perm,policy,mode string,budget,start,tp int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3hPrepressure(arms,rot)
	if mode=="static_pre"{for i:=range arms{arms[i].r.write(fmt.Sprintf("3k-static-%s-%s-%d-%d-%d-%d",policy,perm,budget,start,tp,i),"x")}}
	for round:=0;round<6;round++{
		if policy!="baseline"&&round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{
			arms[i].r.write(fmt.Sprintf("3k-%s-%s-%s-%d-%d-%d-%d-%d",mode,policy,perm,budget,start,tp,round,i),"x")
			if (mode=="dynamic_r3"&&round==3)||(mode=="dynamic_r5"&&round==5){arms[i].r.write(fmt.Sprintf("3k-extra-%s-%s-%s-%d-%d-%d-%d",mode,policy,perm,budget,start,tp,i),"x")}
		}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func RunUPLM3K()(UPLM3KResult,error){
	modes:=[]string{"none","static_pre","dynamic_r3","dynamic_r5"};budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3KResult{Schema:UPLM3KPressurePlacementSchema,Experiment:"UP-LM3K-pressure-placement-control",SourceUPLM3JSeal:"2a8db4e75c0bc5c49271eb7760cbb8cb085c78c0",PressureModes:modes,IdentityRotations:rots,Permutations:perms,Budgets:budgets,ActionStartRounds:starts,Throughputs:tps,GlobalRounds:6,CounterfactualOnly:true,LiveActivation:false,AdaptivePressureUsed:false,AdaptiveResourcesUsed:false,FutureScheduleUsed:false}
	type key struct{mode string;reach int};type acc struct{n,min,max int};groups:=map[key]*acc{}
	for _,mode:=range modes{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		reach:=uplm3fReachable(budget,start,tp);m:=UPLM3KMetric{PressureMode:mode,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:reach}
		for _,rot:=range rots{for _,perm:=range perms{
			_,bf,_:=uplm3kRun(rot,perm,"baseline",mode,budget,start,tp)
			_,ef,ea:=uplm3kRun(rot,perm,"earliest_deadline",mode,budget,start,tp)
			_,ff,fa:=uplm3kRun(rot,perm,"fixed_order",mode,budget,start,tp)
			m.BaselineFailed+=bf;m.EarliestFailed+=ef;m.FixedFailed+=ff;m.EarliestActions+=ea;m.FixedActions+=fa
		}}
		res.Metrics=append(res.Metrics,m)
		k:=key{mode:mode,reach:reach};a:=groups[k];if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};groups[k]=a};a.n++;if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
	}}}}
	keys:=make([]key,0,len(groups));for k:=range groups{keys=append(keys,k)}
	sort.Slice(keys,func(i,j int)bool{if keys[i].mode!=keys[j].mode{return keys[i].mode<keys[j].mode};return keys[i].reach<keys[j].reach})
	for _,k:=range keys{a:=groups[k];res.Summaries=append(res.Summaries,UPLM3KSummary{PressureMode:k.mode,ReachableBudget:k.reach,Observations:a.n,MinEarliestFailed:a.min,MaxEarliestFailed:a.max,FailureSpread:a.max-a.min})}
	return res,nil
}
