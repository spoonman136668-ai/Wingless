package unitary

const UP235CRiskGateSchema="wingless.up235c-risk-gated-action9.v1"

type UP235CMetric struct{
 Schedule string `json:"schedule"`
 PhaseAdvanceWrites int `json:"phase_advance_writes"`
 Arms int `json:"arms"`
 Cap8Losses int `json:"cap8_losses"`
 SingleLosses int `json:"single_losses"`
 RiskLosses int `json:"risk_losses"`
 SingleNinthActions int `json:"single_ninth_actions"`
 RiskNinthActions int `json:"risk_ninth_actions"`
 SingleRescues int `json:"single_rescues"`
 RiskRescues int `json:"risk_rescues"`
 SingleUnnecessary int `json:"single_unnecessary"`
 RiskUnnecessary int `json:"risk_unnecessary"`
 SingleIneffective int `json:"single_ineffective"`
 RiskIneffective int `json:"risk_ineffective"`
 SingleRescuePerUnnecessary float64 `json:"single_rescue_per_unnecessary"`
 RiskRescuePerUnnecessary float64 `json:"risk_rescue_per_unnecessary"`
}
type UP235CResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP234CSeal string `json:"source_up234c_seal"`
 Horizon int `json:"horizon"`
 MonitoringCadence int `json:"monitoring_cadence"`
 MaxActions int `json:"max_actions"`
 CounterfactualOnly bool `json:"counterfactual_only"`
 ThresholdFittingUsed bool `json:"threshold_fitting_used"`
 PrototypeMutationUsed bool `json:"prototype_mutation_used"`
 LearnedWeightsUsed bool `json:"learned_weights_used"`
 NewNativeFieldUsed bool `json:"new_native_field_used"`
 RetryAfterRejectedFirstEvent bool `json:"retry_after_rejected_first_event"`
 ActionBudgetIncreased bool `json:"action_budget_increased"`
 NewActionTypeUsed bool `json:"new_action_type_used"`
 LiveActivation bool `json:"live_activation"`
 Metrics []UP235CMetric `json:"metrics"`
}
func up235cRun(x0 *up81cAging,endangered,advance int,a,b,mode string)(loss,actions,ninth int){
 m:=*x0;committed:=false
 steps:=make([]int,9);for i:=range steps{steps[i]=-1}
 action4:=-1;due5,due6,due7:=-1,-1,-1;eighthBoundary:=-1;firstEventSeen:=false
 failEx:=up233cExemplars(up232cFailureSignatures);survEx:=up233cExemplars(up232cSurvivorSignatures)
 for start:=1;start<=80;start+=2{
  if !committed{
   if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;steps[0]=start;committed=true}
  }else{
   if actions>=1&&actions<4{
    last:=steps[actions-1]
    if last>=0&&start>=last+8{
     m.query(endangered);actions++;steps[actions-1]=start
     if actions==4{action4=start;due5=action4+8;due6=action4+16;due7=action4+24}
    }
   }else if actions==4{
    fire:=start>=due5
    if start<due5&&up161cAdversarial(&m,endangered)<=2{fire=true}
    if fire{m.query(endangered);actions++;steps[4]=start}
   }else if actions==5{
    if start>=due6{m.query(endangered);actions++;steps[5]=start}
   }else if actions==6{
    if start>=due7{m.query(endangered);actions++;steps[6]=start}
   }else if actions==7{
    if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;steps[7]=start;eighthBoundary=start}
   }else if actions==8&&mode!="cap8"&&eighthBoundary>=0&&start>eighthBoundary&&!firstEventSeen{
    if up161cAdversarial(&m,endangered)<=2{
     firstEventSeen=true;fire:=mode=="single_critical"
     if mode=="hamming_risk"{
      snap:=up223cSnapshot(&m,endangered,start)
      risk,_,_:=up233cClass(true,snap,failEx,survEx)
      fire=risk=="high_risk"
     }
     if fire{m.query(endangered);actions++;steps[8]=start;ninth++}
    }
   }
  }
  for j:=0;j<2&&start+j<=80;j++{
   step:=start+j
   if up165cRealStep(&m,endangered,1055000+step,step,up210cPolicy(step,advance,a,b)){return step,actions,ninth}
  }
 }
 return 81,actions,ninth
}
func RunUP235C()(UP235CResult,error){
 cohorts:=[][]int{{0,8,9,15},{1,4,10,14},{2,5,11,12},{3,6,7,13}}
 hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
 type sched struct{name,a,b string}
 schedules:=[]sched{{"no_hostile","no_refresh","hostile_shield"},{"hostile_no","hostile_shield","no_refresh"},{"alternating_fixed","alternating_shield","fixed_offset_refresh"},{"fixed_alternating","fixed_offset_refresh","alternating_shield"}}
 r:=UP235CResult{Schema:UP235CRiskGateSchema,Experiment:"UP-235C-risk-gated-action9",SourceUP234CSeal:"ec5e25131dc99cd85a2bccec16b324a3a6118e7c",Horizon:80,MonitoringCadence:2,MaxActions:9,CounterfactualOnly:true,ThresholdFittingUsed:false,PrototypeMutationUsed:false,LearnedWeightsUsed:false,NewNativeFieldUsed:false,RetryAfterRejectedFirstEvent:false,ActionBudgetIncreased:false,NewActionTypeUsed:false,LiveActivation:false}
 for _,s:=range schedules{for _,adv:=range []int{2,20}{
  m:=UP235CMetric{Schedule:s.name,PhaseAdvanceWrites:adv}
  for _,c:=range cohorts{for _,hand:=range hands{
   x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
   bLoss,_,_:=up235cRun(x,e,adv,s.a,s.b,"cap8")
   sLoss,_,sn:=up235cRun(x,e,adv,s.a,s.b,"single_critical")
   rLoss,_,rn:=up235cRun(x,e,adv,s.a,s.b,"hamming_risk")
   bf,sf,rf:=bLoss<=80,sLoss<=80,rLoss<=80
   if bf{m.Cap8Losses++};if sf{m.SingleLosses++};if rf{m.RiskLosses++}
   m.SingleNinthActions+=sn;m.RiskNinthActions+=rn
   if bf&&!sf{m.SingleRescues++};if bf&&!rf{m.RiskRescues++}
   if !bf&&sn>0{m.SingleUnnecessary++};if !bf&&rn>0{m.RiskUnnecessary++}
   if bf&&sf&&sn>0{m.SingleIneffective++};if bf&&rf&&rn>0{m.RiskIneffective++}
  }}
  if m.SingleUnnecessary>0{m.SingleRescuePerUnnecessary=float64(m.SingleRescues)/float64(m.SingleUnnecessary)}
  if m.RiskUnnecessary>0{m.RiskRescuePerUnnecessary=float64(m.RiskRescues)/float64(m.RiskUnnecessary)}
  r.Metrics=append(r.Metrics,m)
 }}
 return r,nil
}
