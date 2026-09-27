package unitary

import "fmt"

const UPLM4XOverlapSchema="wingless.up-lm4x-overlap-transition-checkpoint.v1"

type UPLM4XSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
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
type UPLM4XResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4WSeal string `json:"source_up_lm4w_seal"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCheckpointSearchUsed bool `json:"adaptive_checkpoint_search_used"`
	Summaries []UPLM4XSummary `json:"summaries"`
}
func uplm4xSchedule(budget,start,tp,tred int)(r,end,pb,pt,post3 int){
	w:=UPLM4WWindow{BudgetCut:3,BudgetRestore:5,ThroughputCut:1,ThroughputRestore:3}
	actions:=0;end=start-1
	for round:=0;round<6;round++{
		if round==w.BudgetRestore{pb=actions}
		if round==w.ThroughputRestore{pt=actions}
		if round>=start{
			cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,2)
			eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred)
			before:=actions
			for k:=0;k<eff&&actions<cap;k++{actions++}
			if actions>before{end=round}
		}
		if round==3{post3=actions}
	}
	return actions,end,pb,pt,post3
}
func uplm4xKey(coord string,r,tp,end,pb,pt,post3 int)string{
	if coord=="both-restores"{return fmt.Sprintf("r=%d|tp=%d|end=%d|pb=%d|pt=%d",r,tp,end,pb,pt)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|pb=%d|pt=%d|post3=%d",r,tp,end,pb,pt,post3)
}
func RunUPLM4X()(UPLM4XResult,error){
	profiles:=[]string{"layout_only","hybrid_min"};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"both-restores","both-restores+post-round3"}
	res:=UPLM4XResult{Schema:UPLM4XOverlapSchema,Experiment:"UP-LM4X-overlap-transition-checkpoint",SourceUPLM4WSeal:"417dbcd932142c34f8e3a9f4f17ca67d25807636",TopologyConditionCells:24,ResourceConfigurationsPerCell:64,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCheckpointSearchUsed:false}
	type point struct{r,tp,end,pb,pt,p3,out int};type acc struct{min,max,n int}
	w:=UPLM4WWindow{BudgetCut:3,BudgetRestore:5,ThroughputCut:1,ThroughputRestore:3}
	for _,profile:=range profiles{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pb,pt,p3:=uplm4xSchedule(b,st,tp,tr)
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,2,tr)
			pts=append(pts,point{r,tp,end,pb,pt,p3,out})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm4xKey(coord,p.r,p.tp,p.end,p.pb,p.pt,p.p3);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4XSummary{DeadlineProfile:profile,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}
	return res,nil
}
