package unitary

const UPLM5KRestore1Schema="wingless.up-lm5k-restore-round1-extrapolation.v1"

type UPLM5KSummary struct{
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
type UPLM5KResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5JSeal string `json:"source_up_lm5j_seal"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	PointsPerPooledCell int `json:"points_per_pooled_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	IntroducesRestoreRound1 bool `json:"introduces_restore_round1"`
	RawActionStartInKeys bool `json:"raw_action_start_in_keys"`
	RawBudgetInKeys bool `json:"raw_budget_in_keys"`
	CutRoundLabelsInKeys bool `json:"cut_round_labels_in_keys"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5KSummary `json:"summaries"`
}
func RunUPLM5K()(UPLM5KResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{
		{BudgetCut:0,BudgetRestore:1,ThroughputCut:1,ThroughputRestore:3},
		{BudgetCut:1,BudgetRestore:3,ThroughputCut:0,ThroughputRestore:1},
		{BudgetCut:0,BudgetRestore:1,ThroughputCut:2,ThroughputRestore:4},
		{BudgetCut:2,BudgetRestore:4,ThroughputCut:0,ThroughputRestore:1},
		{BudgetCut:0,BudgetRestore:1,ThroughputCut:3,ThroughputRestore:5},
		{BudgetCut:3,BudgetRestore:5,ThroughputCut:0,ThroughputRestore:1},
	}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"frozen","frozen+deployment-onset"}
	res:=UPLM5KResult{Schema:UPLM5KRestore1Schema,Experiment:"UP-LM5K-restore-round1-extrapolation",SourceUPLM5JSeal:"bd3c13a67abe238062a3a2be93369617c3c0b864",PooledConditionCells:72,PointsPerPooledCell:384,CandidateCoordinates:coords,IntroducesRestoreRound1:true,RawActionStartInKeys:false,RawBudgetInKeys:false,CutRoundLabelsInKeys:false,AdaptiveCoordinateSearchUsed:false,CounterfactualOnly:true,LiveActivation:false}
	type point struct{r,tp,end,preE,postE,preL,postL,earlyR,lateR,onset,out int;earlyResource string}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		pts:=make([]point,0,384)
		for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,preE,postE,preL,postL,earlyR,lateR,onset,earlyResource:=uplm5iSchedule(b,st,tp,w,br,tr)
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
			pts=append(pts,point{r:r,tp:tp,end:end,preE:preE,postE:postE,preL:preL,postL:postL,earlyR:earlyR,lateR:lateR,onset:onset,out:out,earlyResource:earlyResource})
		}}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm5iKey(coord,p.r,p.tp,p.end,p.preE,p.postE,p.preL,p.postL,p.earlyR,p.lateR,p.onset,p.earlyResource)
				a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			s:=UPLM5KSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,PooledPoints:len(pts),Groups:len(g)}
			sum:=0
			for _,a:=range g{if a.n>1{s.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
			if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,s)
		}
	}}}}}
	return res,nil
}
