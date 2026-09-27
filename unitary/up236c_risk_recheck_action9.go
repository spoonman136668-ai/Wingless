package unitary

const UP236CRecheckSchema="wingless.up236c-risk-recheck-action9.v1"

type UP236CMetric struct{
 Schedule string `json:"schedule"`
 PhaseAdvanceWrites int `json:"phase_advance_writes"`
 Arms int `json:"arms"`
 Cap8Losses int `json:"cap8_losses"`
 SingleLosses int `json:"single_losses"`
 FirstRiskLosses int `json:"first_risk_losses"`
 RecheckRiskLosses int `json:"recheck_risk_losses"`
 SingleNinth int `json:"single_ninth"`
 FirstRiskNinth int `json:"first_risk_ninth"`
 RecheckRiskNinth int `json:"recheck_risk_ninth"`
 FirstRiskRescues int `json:"first_risk_rescues"`
 RecheckRiskRescues int `json:"recheck_risk_rescues"`
 FirstRiskUnnecessary int `json:"first_risk_unnecessary"`
 RecheckRiskUnnecessary int `json:"recheck_risk_unnecessary"`
 DelayedAuthorizations int `json:"delayed_authorizations"`
}
type UP236CResult struct{
 Schema string `json:"schema"`; Experiment string `json:"experiment"`; SourceUP235CSeal string `json:"source_up235c_seal"`
 MaxActions int `json:"max_actions"`; CounterfactualOnly bool `json:"counterfactual_only"`
 ClassifierChanged bool `json:"classifier_changed"`; PrototypeMutationUsed bool `json:"prototype_mutation_used"`
 ThresholdFittingUsed bool `json:"threshold_fitting_used"`; NewNativeFieldUsed bool `json:"new_native_field_used"`
 ActionBudgetIncreased bool `json:"action_budget_increased"`; LiveActivation bool `json:"live_activation"`
 Metrics []UP236CMetric `json:"metrics"`
}
func up236cRun(x0 *up81cAging,endangered,advance int,a,b,mode string)(loss,ninth int,delayed bool){
 m:=*x0;committed:=false;actions:=0
 steps:=make([]int,9);for i:=range steps{steps[i]=-1}
 action4:=-1;due5,due6,due7:=-1,-1,-1;eighth:=-1;seen:=false
 failEx:=up233cExemplars(up232cFailureSignatures);survEx:=up233cExemplars(up232cSurvivorSignatures)
 for start:=1;start<=80;start+=2{
  if !committed{
   if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;steps[0]=start;committed=true}
  }else if actions>=1&&actions<4{
   last:=steps[actions-1];if last>=0&&start>=last+8{m.query(endangered);actions++;steps[actions-1]=start;if actions==4{action4=start;due5=action4+8;due6=action4+16;due7=action4+24}}
  }else if actions==4{
   fire:=start>=due5;if start<due5&&up161cAdversarial(&m,endangered)<=2{fire=true};if fire{m.query(endangered);actions++;steps[4]=start}
  }else if actions==5{
   if start>=due6{m.query(endangered);actions++;steps[5]=start}
  }else if actions==6{
   if start>=due7{m.query(endangered);actions++;steps[6]=start}
  }else if actions==7{
   if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;steps[7]=start;eighth=start}
  }else if actions==8&&mode!="cap8"&&eighth>=0&&start>eighth&&up161cAdversarial(&m,endangered)<=2{
   fire:=false
   if mode=="single"{fire=true
   }else{
    snap:=up223cSnapshot(&m,endangered,start);risk,_,_:=up233cClass(true,snap,failEx,survEx)
    if risk=="high_risk"{fire=true;if seen{delayed=true}}
    if mode=="first"&&seen{fire=false}
   }
   if !seen{seen=true}
   if fire{m.query(endangered);actions++;steps[8]=start;ninth++}
  }
  for j:=0;j<2&&start+j<=80;j++{
   step:=start+j
   if up165cRealStep(&m,endangered,1056000+step,step,up210cPolicy(step,advance,a,b)){return step,ninth,delayed}
  }
 }
 return 81,ninth,delayed
}
func RunUP236C()(UP236CResult,error){
 cohorts:=[][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}}
 hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
 type sched struct{name,a,b string}
 schedules:=[]sched{{"no_hostile","no_refresh","hostile_shield"},{"hostile_no","hostile_shield","no_refresh"},{"alternating_fixed","alternating_shield","fixed_offset_refresh"},{"fixed_alternating","fixed_offset_refresh","alternating_shield"}}
 r:=UP236CResult{Schema:UP236CRecheckSchema,Experiment:"UP-236C-risk-recheck-action9",SourceUP235CSeal:"0868cccf7a9d0ab0d8242251eef8de5f1fd7c7aa",MaxActions:9,CounterfactualOnly:true,ClassifierChanged:false,PrototypeMutationUsed:false,ThresholdFittingUsed:false,NewNativeFieldUsed:false,ActionBudgetIncreased:false,LiveActivation:false}
 for _,s:=range schedules{for _,adv:=range []int{0,22}{
  m:=UP236CMetric{Schedule:s.name,PhaseAdvanceWrites:adv}
  for _,c:=range cohorts{for _,hand:=range hands{
   x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
   bLoss,_,_:=up236cRun(x,e,adv,s.a,s.b,"cap8")
   sLoss,sn,_:=up236cRun(x,e,adv,s.a,s.b,"single")
   fLoss,fn,_:=up236cRun(x,e,adv,s.a,s.b,"first")
   rLoss,rn,delayed:=up236cRun(x,e,adv,s.a,s.b,"recheck")
   bf,sf,ff,rf:=bLoss<=80,sLoss<=80,fLoss<=80,rLoss<=80
   if bf{m.Cap8Losses++};if sf{m.SingleLosses++};if ff{m.FirstRiskLosses++};if rf{m.RecheckRiskLosses++}
   m.SingleNinth+=sn;m.FirstRiskNinth+=fn;m.RecheckRiskNinth+=rn
   if bf&&!ff{m.FirstRiskRescues++};if bf&&!rf{m.RecheckRiskRescues++}
   if !bf&&fn>0{m.FirstRiskUnnecessary++};if !bf&&rn>0{m.RecheckRiskUnnecessary++}
   if delayed{m.DelayedAuthorizations++}
  }}
  r.Metrics=append(r.Metrics,m)
 }}
 return r,nil
}
