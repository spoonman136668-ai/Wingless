package unitary

import "fmt"

const UPLM5GPostLateSchema="wingless.up-lm5g-post-late-restore-checkpoint.v1"

type UPLM5GSummary struct{
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
type UPLM5GResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5FSeal string `json:"source_up_lm5f_seal"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	PointsPerPooledCell int `json:"points_per_pooled_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CutRoundLabelsInKeys bool `json:"cut_round_labels_in_keys"`
	RawActionStartInKeys bool `json:"raw_action_start_in_keys"`
	RawBudgetInKeys bool `json:"raw_budget_in_keys"`
	FullScheduleIdentityInKeys bool `json:"full_schedule_identity_in_keys"`
	AdaptiveCheckpointSearchUsed bool `json:"adaptive_checkpoint_search_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5GSummary `json:"summaries"`
}
func uplm5gSchedule(budget,start,tp int,w UPLM4WWindow,bred,tred int)(r,end,preE,postE,preL,postL,earlyR,lateR int,earlyResource string){
	actions:=0;end=start-1
	earlyR=w.BudgetRestore;lateR=w.ThroughputRestore;earlyResource="budget_first"
	if w.ThroughputRestore<w.BudgetRestore{earlyR=w.ThroughputRestore;lateR=w.BudgetRestore;earlyResource="throughput_first"}
	for round:=0;round<6;round++{
		if round==earlyR{preE=actions}
		if round==lateR{preL=actions}
		if round>=start{
			cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,bred)
			eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred)
			before:=actions
			for k:=0;k<eff&&actions<cap;k++{actions++}
			if actions>before{end=round}
		}
		if round==earlyR{postE=actions}
		if round==lateR{postL=actions}
	}
	return actions,end,preE,postE,preL,postL,earlyR,lateR,earlyResource
}
func uplm5gKey(coord string,r,tp,end,preE,postE,preL,postL,earlyR,lateR int,earlyResource string)string{
	base:=fmt.Sprintf("r=%d|tp=%d|end=%d|preE=%d|postE=%d|preL=%d|earlyR=%d|lateR=%d|earlyResource=%s",r,tp,end,preE,postE,preL,earlyR,lateR,earlyResource)
	if coord=="timing-aware+early-resource"{return base}
	return fmt.Sprintf("%s|postL=%d",base,postL)
}
func RunUPLM5G()(UPLM5GResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:1,ThroughputRestore:3},
		{BudgetCut:1,BudgetRestore:3,ThroughputCut:0,ThroughputRestore:2},
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:2,ThroughputRestore:4},
		{BudgetCut:2,BudgetRestore:4,ThroughputCut:0,ThroughputRestore:2},
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:3,ThroughputRestore:5},
		{BudgetCut:3,BudgetRestore:5,ThroughputCut:0,ThroughputRestore:2},
	}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"timing-aware+early-resource","timing-aware+early-resource+post-late"}
	res:=UPLM5GResult{
		Schema:UPLM5GPostLateSchema,Experiment:"UP-LM5G-post-late-restore-checkpoint",
		SourceUPLM5FSeal:"1ccefa94b81ce0ed4e0eab711ee066e5f985721a",
		PooledConditionCells:72,PointsPerPooledCell:384,CandidateCoordinates:coords,
		CutRoundLabelsInKeys:false,RawActionStartInKeys:false,RawBudgetInKeys:false,
		FullScheduleIdentityInKeys:false,AdaptiveCheckpointSearchUsed:false,
		CounterfactualOnly:true,LiveActivation:false,
	}
	type point struct{r,tp,end,preE,postE,preL,postL,earlyR,lateR,out int;earlyResource string}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,preE,postE,preL,postL,earlyR,lateR,earlyResource:=uplm5gSchedule(b,st,tp,w,br,tr)
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
			pts=append(pts,point{r:r,tp:tp,end:end,preE:preE,postE:postE,preL:preL,postL:postL,earlyR:earlyR,lateR:lateR,out:out,earlyResource:earlyResource})
		}}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm5gKey(coord,p.r,p.tp,p.end,p.preE,p.postE,p.preL,p.postL,p.earlyR,p.lateR,p.earlyResource)
				a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM5GSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,PooledPoints:len(pts),Groups:len(g)}
			sum:=0
			for _,a:=range g{
				if a.n>1{sm.MultiMemberGroups++}
				d:=a.max-a.min;sum+=d
				if d>0{sm.NonzeroSpreadGroups++}
				if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}
			}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
