package unitary

import "fmt"

const UPLM5DHeldoutSchema="wingless.up-lm5d-heldout-cut-window-transfer.v1"

type UPLM5DSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	BudgetReduction int `json:"budget_reduction"`
	ThroughputReduction int `json:"throughput_reduction"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	PooledPoints int `json:"pooled_points"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM5DResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5CSeal string `json:"source_up_lm5c_seal"`
	HeldoutScheduleCount int `json:"heldout_schedule_count"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	PointsPerPooledCell int `json:"points_per_pooled_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CutRoundLabelsInKeys bool `json:"cut_round_labels_in_keys"`
	ResourceIdentityInKeys bool `json:"resource_identity_in_keys"`
	FullScheduleIdentityInKeys bool `json:"full_schedule_identity_in_keys"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveScheduleSelectionUsed bool `json:"adaptive_schedule_selection_used"`
	Summaries []UPLM5DSummary `json:"summaries"`
}
func uplm5dKey(coord string,r,tp,end,pe,preE,preL,earlyR,lateR int)string{
	if coord=="ordered-path"{
		return fmt.Sprintf("r=%d|tp=%d|end=%d|preE=%d|postE=%d|preL=%d",r,tp,end,preE,pe,preL)
	}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|preE=%d|postE=%d|preL=%d|earlyR=%d|lateR=%d",r,tp,end,preE,pe,preL,earlyR,lateR)
}
func RunUPLM5D()(UPLM5DResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{
		{BudgetCut:2,BudgetRestore:3,ThroughputCut:1,ThroughputRestore:4},
		{BudgetCut:1,BudgetRestore:4,ThroughputCut:2,ThroughputRestore:3},
		{BudgetCut:2,BudgetRestore:4,ThroughputCut:3,ThroughputRestore:5},
		{BudgetCut:3,BudgetRestore:5,ThroughputCut:2,ThroughputRestore:4},
		{BudgetCut:2,BudgetRestore:3,ThroughputCut:1,ThroughputRestore:5},
		{BudgetCut:1,BudgetRestore:5,ThroughputCut:2,ThroughputRestore:3},
	}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"ordered-path","ordered-path+restore-rounds"}
	res:=UPLM5DResult{Schema:UPLM5DHeldoutSchema,Experiment:"UP-LM5D-heldout-cut-window-transfer",SourceUPLM5CSeal:"0e2c683764c234b28289be43dc0a859fc00e14fb",HeldoutScheduleCount:len(windows),PooledConditionCells:72,PointsPerPooledCell:384,CandidateCoordinates:coords,CutRoundLabelsInKeys:false,ResourceIdentityInKeys:false,FullScheduleIdentityInKeys:false,CounterfactualOnly:true,LiveActivation:false,AdaptiveScheduleSelectionUsed:false}
	type point struct{r,tp,end,pe,preE,preL,earlyR,lateR,out int};type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pb,pt,pe:=uplm4ySchedule(b,st,tp,w,br,tr)
			preE,preL:=pb,pt;earlyR,lateR:=w.BudgetRestore,w.ThroughputRestore
			if w.ThroughputRestore<w.BudgetRestore{preE,preL=pt,pb;earlyR,lateR=w.ThroughputRestore,w.BudgetRestore}
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
			pts=append(pts,point{r:r,tp:tp,end:end,pe:pe,preE:preE,preL:preL,earlyR:earlyR,lateR:lateR,out:out})
		}}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm5dKey(coord,p.r,p.tp,p.end,p.pe,p.preE,p.preL,p.earlyR,p.lateR);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM5DSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,PooledPoints:len(pts),Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
