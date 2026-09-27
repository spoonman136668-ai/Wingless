package unitary

import "fmt"

const UPLM4VCombinedSchema="wingless.up-lm4v-combined-revoke-restore-transfer.v1"

type UPLM4VReduction struct{Budget,Throughput int}
type UPLM4VSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	CutRound int `json:"cut_round"`
	RestoreRound int `json:"restore_round"`
	BudgetReduction int `json:"budget_reduction"`
	ThroughputReduction int `json:"throughput_reduction"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4VResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4USeal string `json:"source_up_lm4u_seal"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4VSummary `json:"summaries"`
}
func uplm4vCaps(budget,tp,round,cut,restore int,red UPLM4VReduction)(int,int){
	if round<cut||round>=restore{return budget,tp}
	b:=budget-red.Budget;if b<0{b=0}
	t:=tp-red.Throughput;if t<0{t=0}
	return b,t
}
func uplm4vSchedule(budget,start,tp,cut,restore int,red UPLM4VReduction)(reach,end,preRestore int){
	actions:=0;end=start-1
	for round:=0;round<6;round++{
		if round==restore{preRestore=actions}
		if round<start{continue}
		cap,eff:=uplm4vCaps(budget,tp,round,cut,restore,red);before:=actions
		for k:=0;k<eff&&actions<cap;k++{actions++}
		if actions>before{end=round}
	}
	return actions,end,preRestore
}
func uplm4vRun(rot int,perm,profile string,budget,start,tp,cut,restore int,red UPLM4VReduction)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile);actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			cap,eff:=uplm4vCaps(budget,tp,round,cut,restore,red);used:=map[int]bool{}
			for k:=0;k<eff&&actions<cap;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4v-%s-%d-%d-%d-%d-%s-%d-%d-%d-%d",profile,cut,restore,red.Budget,red.Throughput,perm,budget,start,tp,round),"x")}
	}
	failed:=0;for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f};return failed
}
func RunUPLM4V()(UPLM4VResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	schedules:=[]UPLM4NSchedule{{1,3},{2,4},{3,5}}
	reductions:=[]UPLM4VReduction{{1,1},{1,2},{2,1},{2,2}}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"final","final+prerestore"}
	res:=UPLM4VResult{Schema:UPLM4VCombinedSchema,Experiment:"UP-LM4V-combined-revoke-restore-transfer",SourceUPLM4USeal:"8200c546c26e4b35b6e963df49f7adf08e907b07",TopologyConditionCells:216,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{r,tp,end,pre,out int};type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,sch:=range schedules{for _,red:=range reductions{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pre:=uplm4vSchedule(b,st,tp,sch.CutRound,sch.RestoreRound,red)
			pts=append(pts,point{r,tp,end,pre,uplm4vRun(rot,perm,profile,b,st,tp,sch.CutRound,sch.RestoreRound,red)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{k:=uplm4uKey(coord,p.r,p.tp,p.end,p.pre);a:=g[k];if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a};a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}}
			sm:=UPLM4VSummary{DeadlineProfile:profile,CutRound:sch.CutRound,RestoreRound:sch.RestoreRound,BudgetReduction:red.Budget,ThroughputReduction:red.Throughput,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))};res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
