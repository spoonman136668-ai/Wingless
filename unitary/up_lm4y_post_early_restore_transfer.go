package unitary

import "fmt"

const UPLM4YPostEarlySchema="wingless.up-lm4y-post-early-restore-transfer.v1"

type UPLM4YSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	BudgetCut int `json:"budget_cut"`
	BudgetRestore int `json:"budget_restore"`
	ThroughputCut int `json:"throughput_cut"`
	ThroughputRestore int `json:"throughput_restore"`
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
type UPLM4YResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4XSeal string `json:"source_up_lm4x_seal"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCheckpointSearchUsed bool `json:"adaptive_checkpoint_search_used"`
	Summaries []UPLM4YSummary `json:"summaries"`
}
func uplm4ySchedule(budget,start,tp int,w UPLM4WWindow,bred,tred int)(r,end,pb,pt,postEarly int){
	actions:=0;end=start-1
	early:=w.BudgetRestore;if w.ThroughputRestore<early{early=w.ThroughputRestore}
	for round:=0;round<6;round++{
		if round==w.BudgetRestore{pb=actions}
		if round==w.ThroughputRestore{pt=actions}
		if round>=start{
			cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,bred)
			eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred)
			before:=actions
			for k:=0;k<eff&&actions<cap;k++{actions++}
			if actions>before{end=round}
		}
		if round==early{postEarly=actions}
	}
	return actions,end,pb,pt,postEarly
}
func uplm4yKey(coord string,r,tp,end,pb,pt,pe int)string{
	if coord=="both-restores"{return fmt.Sprintf("r=%d|tp=%d|end=%d|pb=%d|pt=%d",r,tp,end,pb,pt)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|pb=%d|pt=%d|pe=%d",r,tp,end,pb,pt,pe)
}
func RunUPLM4Y()(UPLM4YResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{{1,3,2,4},{2,4,1,3},{1,4,2,5},{2,5,1,4},{1,3,3,5},{3,5,1,3}}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"both-restores","both-restores+post-early-restore"}
	res:=UPLM4YResult{Schema:UPLM4YPostEarlySchema,Experiment:"UP-LM4Y-post-early-restore-transfer",SourceUPLM4XSeal:"acacd69acdea3a684b23a8f8403fdf68211cfcf7",TopologyConditionCells:432,ResourceConfigurationsPerCell:64,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCheckpointSearchUsed:false}
	type point struct{r,tp,end,pb,pt,pe,out int};type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,w:=range windows{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pb,pt,pe:=uplm4ySchedule(b,st,tp,w,br,tr)
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
			pts=append(pts,point{r,tp,end,pb,pt,pe,out})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm4yKey(coord,p.r,p.tp,p.end,p.pb,p.pt,p.pe);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4YSummary{DeadlineProfile:profile,BudgetCut:w.BudgetCut,BudgetRestore:w.BudgetRestore,ThroughputCut:w.ThroughputCut,ThroughputRestore:w.ThroughputRestore,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}}
	return res,nil
}
