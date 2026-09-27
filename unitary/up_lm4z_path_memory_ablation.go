package unitary

import "fmt"

const UPLM4ZPathMemorySchema="wingless.up-lm4z-path-memory-ablation.v1"

type UPLM4ZSummary struct{
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
type UPLM4ZResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4YSeal string `json:"source_up_lm4y_seal"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	Summaries []UPLM4ZSummary `json:"summaries"`
}
func uplm4zKey(coord string,r,tp,end,pb,pt,postEarly,preLate int)string{
	switch coord{
	case "final+post-early":
		return fmt.Sprintf("r=%d|tp=%d|end=%d|pe=%d",r,tp,end,postEarly)
	case "final+post-early+pre-late":
		return fmt.Sprintf("r=%d|tp=%d|end=%d|pe=%d|pl=%d",r,tp,end,postEarly,preLate)
	default:
		return fmt.Sprintf("r=%d|tp=%d|end=%d|pb=%d|pt=%d|pe=%d",r,tp,end,pb,pt,postEarly)
	}
}
func RunUPLM4Z()(UPLM4ZResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{{1,3,2,4},{2,4,1,3},{1,4,2,5},{2,5,1,4},{1,3,3,5},{3,5,1,3}}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"final+post-early","final+post-early+pre-late","full"}
	res:=UPLM4ZResult{Schema:UPLM4ZPathMemorySchema,Experiment:"UP-LM4Z-path-memory-ablation",SourceUPLM4YSeal:"0724b89e0630e14e832d56d3754736449610e2cd",TopologyConditionCells:432,ResourceConfigurationsPerCell:64,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCoordinateSearchUsed:false}
	type point struct{r,tp,end,pb,pt,pe,pl,out int};type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,w:=range windows{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pb,pt,pe:=uplm4ySchedule(b,st,tp,w,br,tr)
			pl:=pb;if w.ThroughputRestore>w.BudgetRestore{pl=pt}
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
			pts=append(pts,point{r,tp,end,pb,pt,pe,pl,out})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm4zKey(coord,p.r,p.tp,p.end,p.pb,p.pt,p.pe,p.pl);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4ZSummary{DeadlineProfile:profile,BudgetCut:w.BudgetCut,BudgetRestore:w.BudgetRestore,ThroughputCut:w.ThroughputCut,ThroughputRestore:w.ThroughputRestore,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}}
	return res,nil
}
