package unitary

import "fmt"

const UPLM4VCombinedSchema="wingless.up-lm4v-budget-throughput-revoke-restore.v1"

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
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4VSummary `json:"summaries"`
}
func uplm4vCap(budget,round,cut,restore,reduce int)int{
	if round<cut||round>=restore{return budget}
	x:=budget-reduce;if x<0{return 0};return x
}
func uplm4vTP(tp,round,cut,restore,reduce int)int{
	if round<cut||round>=restore{return tp}
	x:=tp-reduce;if x<0{return 0};return x
}
func uplm4vSchedule(budget,start,tp,cut,restore,bred,tred int)(reach,end,preRestore int){
	actions:=0;end=start-1
	for round:=0;round<6;round++{
		if round==restore{preRestore=actions}
		if round<start{continue}
		cap:=uplm4vCap(budget,round,cut,restore,bred)
		eff:=uplm4vTP(tp,round,cut,restore,tred)
		before:=actions
		for k:=0;k<eff&&actions<cap;k++{actions++}
		if actions>before{end=round}
	}
	return actions,end,preRestore
}
func uplm4vRun(rot int,perm,profile string,budget,start,tp,cut,restore,bred,tred int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			cap:=uplm4vCap(budget,round,cut,restore,bred)
			eff:=uplm4vTP(tp,round,cut,restore,tred)
			used:=map[int]bool{}
			for k:=0;k<eff&&actions<cap;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4v-%s-%d-%d-%d-%d-%s-%d-%d-%d-%d",profile,cut,restore,bred,tred,perm,budget,start,tp,round),"x")}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4vKey(coord string,r,tp,end,pre int)string{
	if coord=="final"{return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|prerestore=%d",r,tp,end,pre)
}
func RunUPLM4V()(UPLM4VResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	schedules:=[]UPLM4NSchedule{{1,3},{2,4},{3,5}}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"final","final+prerestore"}
	res:=UPLM4VResult{Schema:UPLM4VCombinedSchema,Experiment:"UP-LM4V-budget-throughput-revoke-restore",TopologyConditionCells:216,ResourceConfigurationsPerCell:64,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{r,tp,end,pre,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,sch:=range schedules{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pre:=uplm4vSchedule(b,st,tp,sch.CutRound,sch.RestoreRound,br,tr)
			pts=append(pts,point{r,tp,end,pre,uplm4vRun(rot,perm,profile,b,st,tp,sch.CutRound,sch.RestoreRound,br,tr)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm4vKey(coord,p.r,p.tp,p.end,p.pre);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4VSummary{DeadlineProfile:profile,CutRound:sch.CutRound,RestoreRound:sch.RestoreRound,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}}
	return res,nil
}
