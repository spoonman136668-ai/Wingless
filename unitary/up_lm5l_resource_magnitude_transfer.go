package unitary

const UPLM5LMagnitudeSchema="wingless.up-lm5l-resource-magnitude-transfer.v1"

type UPLM5LSummary struct{
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
type UPLM5LResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUPLM5KSeal string `json:"source_up_lm5k_seal"`
 PooledConditionCells int `json:"pooled_condition_cells"`
 PointsPerPooledCell int `json:"points_per_pooled_cell"`
 CandidateCoordinates []string `json:"candidate_coordinates"`
 Budgets []int `json:"budgets"`
 Starts []int `json:"starts"`
 Throughputs []int `json:"throughputs"`
 RawActionStartInKeys bool `json:"raw_action_start_in_keys"`
 RawBudgetInKeys bool `json:"raw_budget_in_keys"`
 CutRoundLabelsInKeys bool `json:"cut_round_labels_in_keys"`
 AdaptiveGridSelectionUsed bool `json:"adaptive_grid_selection_used"`
 AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
 CounterfactualOnly bool `json:"counterfactual_only"`
 LiveActivation bool `json:"live_activation"`
 Summaries []UPLM5LSummary `json:"summaries"`
}
func RunUPLM5L()(UPLM5LResult,error){
 profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
 windows:=[]UPLM4WWindow{
  {BudgetCut:0,BudgetRestore:1,ThroughputCut:1,ThroughputRestore:3},
  {BudgetCut:1,BudgetRestore:3,ThroughputCut:0,ThroughputRestore:1},
  {BudgetCut:0,BudgetRestore:1,ThroughputCut:2,ThroughputRestore:4},
  {BudgetCut:2,BudgetRestore:4,ThroughputCut:0,ThroughputRestore:1},
  {BudgetCut:0,BudgetRestore:1,ThroughputCut:3,ThroughputRestore:5},
  {BudgetCut:3,BudgetRestore:5,ThroughputCut:0,ThroughputRestore:1},
 }
 breds:=[]int{1,2};treds:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
 budgets:=[]int{2,3,8,9};starts:=[]int{1,2,5,6};tps:=[]int{1,2,5,6};coords:=[]string{"frozen","frozen+deployment-onset"}
 r:=UPLM5LResult{Schema:UPLM5LMagnitudeSchema,Experiment:"UP-LM5L-resource-magnitude-transfer",SourceUPLM5KSeal:"72fa0ef77f4c3e482bc3308aa533f62ec41c85ea",PooledConditionCells:72,PointsPerPooledCell:384,CandidateCoordinates:coords,Budgets:budgets,Starts:starts,Throughputs:tps,RawActionStartInKeys:false,RawBudgetInKeys:false,CutRoundLabelsInKeys:false,AdaptiveGridSelectionUsed:false,AdaptiveCoordinateSearchUsed:false,CounterfactualOnly:true,LiveActivation:false}
 type point struct{r,tp,end,preE,postE,preL,postL,earlyR,lateR,onset,out int;earlyResource string};type acc struct{min,max,n int}
 for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
  pts:=make([]point,0,384)
  for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
   rr,end,preE,postE,preL,postL,earlyR,lateR,onset,earlyResource:=uplm5iSchedule(b,st,tp,w,br,tr)
   out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
   pts=append(pts,point{r:rr,tp:tp,end:end,preE:preE,postE:postE,preL:preL,postL:postL,earlyR:earlyR,lateR:lateR,onset:onset,out:out,earlyResource:earlyResource})
  }}}}
  for _,coord:=range coords{
   g:=map[string]*acc{}
   for _,p:=range pts{
    k:=uplm5iKey(coord,p.r,p.tp,p.end,p.preE,p.postE,p.preL,p.postL,p.earlyR,p.lateR,p.onset,p.earlyResource);a:=g[k]
    if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a};a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
   }
   s:=UPLM5LSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,PooledPoints:len(pts),Groups:len(g)};sum:=0
   for _,a:=range g{if a.n>1{s.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
   if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))};r.Summaries=append(r.Summaries,s)
  }
 }}}}}
 return r,nil
}
