package unitary

import "fmt"

const UPLM5MPolicySchema="wingless.up-lm5m-allocation-policy-transfer.v1"

type UPLM5MSummary struct{
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
type UPLM5MResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUPLM5LSeal string `json:"source_up_lm5l_seal"`
 Policies []string `json:"policies"`
 PooledConditionCells int `json:"pooled_condition_cells"`
 PointsPerPooledCell int `json:"points_per_pooled_cell"`
 CandidateCoordinates []string `json:"candidate_coordinates"`
 PerArmDeliveryInKeys bool `json:"per_arm_delivery_in_keys"`
 AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
 AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
 CounterfactualOnly bool `json:"counterfactual_only"`
 LiveActivation bool `json:"live_activation"`
 Summaries []UPLM5MSummary `json:"summaries"`
}
func uplm5mRun(rot int,perm,profile,policy string,budget,start,tp int,w UPLM4WWindow,bred,tred int)int{
 arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile);actions:=0
 for round:=0;round<6;round++{
  if round>=start{
   cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,bred)
   eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred)
   used:=map[int]bool{}
   for k:=0;k<eff&&actions<cap;k++{
    i:=uplm3cChoose(arms,policy,used);if i<0{break}
    if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
   }
  }
  for i:=range arms{arms[i].r.write(fmt.Sprintf("5m-%s-%s-%s-%d-%d-%d-%d",profile,policy,perm,budget,start,tp,round),"x")}
 }
 failed:=0;for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f};return failed
}
func RunUPLM5M()(UPLM5MResult,error){
 profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
 windows:=[]UPLM4WWindow{
  {BudgetCut:0,BudgetRestore:1,ThroughputCut:1,ThroughputRestore:3},
  {BudgetCut:1,BudgetRestore:3,ThroughputCut:0,ThroughputRestore:1},
  {BudgetCut:0,BudgetRestore:1,ThroughputCut:2,ThroughputRestore:4},
  {BudgetCut:2,BudgetRestore:4,ThroughputCut:0,ThroughputRestore:1},
  {BudgetCut:0,BudgetRestore:1,ThroughputCut:3,ThroughputRestore:5},
  {BudgetCut:3,BudgetRestore:5,ThroughputCut:0,ThroughputRestore:1},
 }
 policies:=[]string{"earliest_deadline","fixed_order"};breds:=[]int{1,2};treds:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
 budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5};coords:=[]string{"enriched","enriched+policy"}
 r:=UPLM5MResult{Schema:UPLM5MPolicySchema,Experiment:"UP-LM5M-allocation-policy-transfer",SourceUPLM5LSeal:"52a30bf88ae98e36fd668e7203968fcd93c31e41",Policies:policies,PooledConditionCells:72,PointsPerPooledCell:768,CandidateCoordinates:coords,PerArmDeliveryInKeys:false,AdaptivePolicySelectionUsed:false,AdaptiveCoordinateSearchUsed:false,CounterfactualOnly:true,LiveActivation:false}
 type point struct{r,tp,end,preE,postE,preL,postL,earlyR,lateR,onset,out int;earlyResource,policy string};type acc struct{min,max,n int}
 for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
  pts:=make([]point,0,768)
  for _,policy:=range policies{for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
   rr,end,preE,postE,preL,postL,earlyR,lateR,onset,earlyResource:=uplm5iSchedule(b,st,tp,w,br,tr)
   out:=uplm5mRun(rot,perm,profile,policy,b,st,tp,w,br,tr)
   pts=append(pts,point{r:rr,tp:tp,end:end,preE:preE,postE:postE,preL:preL,postL:postL,earlyR:earlyR,lateR:lateR,onset:onset,out:out,earlyResource:earlyResource,policy:policy})
  }}}}}
  for _,coord:=range coords{
   g:=map[string]*acc{}
   for _,p:=range pts{
    k:=uplm5iKey("frozen+deployment-onset",p.r,p.tp,p.end,p.preE,p.postE,p.preL,p.postL,p.earlyR,p.lateR,p.onset,p.earlyResource)
    if coord=="enriched+policy"{k+="|policy="+p.policy}
    a:=g[k];if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a};a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
   }
   s:=UPLM5MSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,PooledPoints:len(pts),Groups:len(g)};sum:=0
   for _,a:=range g{if a.n>1{s.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
   if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))};r.Summaries=append(r.Summaries,s)
  }
 }}}}}
 return r,nil
}
