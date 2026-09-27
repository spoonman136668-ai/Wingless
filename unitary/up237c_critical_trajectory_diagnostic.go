package unitary

import "fmt"

const UP237CTrajectorySchema="wingless.up237c-critical-trajectory-diagnostic.v1"

type UP237CArm struct{
 Outcome string `json:"outcome"`
 EventPresent bool `json:"event_present"`
 CurrentSignature string `json:"current_signature,omitempty"`
 TrajectorySignature string `json:"trajectory_signature,omitempty"`
}
type UP237CResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP236CSeal string `json:"source_up236c_seal"`
 ArmsTotal int `json:"arms_total"`
 Failures int `json:"failures"`
 Survivors int `json:"survivors"`
 EventFailures int `json:"event_failures"`
 EventSurvivors int `json:"event_survivors"`
 SilentFailures int `json:"silent_failures"`
 CurrentFailureSignatures int `json:"current_failure_signatures"`
 CurrentSurvivorSignatures int `json:"current_survivor_signatures"`
 CurrentSharedSignatures int `json:"current_shared_signatures"`
 TrajectoryFailureSignatures int `json:"trajectory_failure_signatures"`
 TrajectorySurvivorSignatures int `json:"trajectory_survivor_signatures"`
 TrajectorySharedSignatures int `json:"trajectory_shared_signatures"`
 InterventionChanged bool `json:"intervention_changed"`
 ThresholdFittingUsed bool `json:"threshold_fitting_used"`
 ClassifierTrainingUsed bool `json:"classifier_training_used"`
 AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
 NewNativeFieldUsed bool `json:"new_native_field_used"`
 LiveActivation bool `json:"live_activation"`
 Arms []UP237CArm `json:"arms"`
}
func up237cDelta(a,b UP223CSnapshot)string{
 pred:=0;if a.PredictedIsEndangered!=b.PredictedIsEndangered{if a.PredictedIsEndangered{pred=1}else{pred=-1}}
 return fmt.Sprintf("dage=%d|ddist=%d|dpred=%d|da0=%d|da1=%d|da2=%d|da3=%d|dadv=%d|dnoq=%d",
  a.EndangeredAge-b.EndangeredAge,a.HandDistance-b.HandDistance,pred,
  a.Age0-b.Age0,a.Age1-b.Age1,a.Age2-b.Age2,a.Age3-b.Age3,
  a.AdversarialHorizon-b.AdversarialHorizon,a.NoQueryHorizon-b.NoQueryHorizon)
}
func up237cRun(x0 *up81cAging,endangered,advance int,a,b string)(UP237CArm,error){
 m:=*x0;committed:=false;actions:=0
 steps:=make([]int,8);for i:=range steps{steps[i]=-1}
 action4:=-1;due5,due6,due7:=-1,-1,-1;eighth:=-1
 var prev UP223CSnapshot;prevValid:=false
 out:=UP237CArm{Outcome:"survive"}
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
   if up161cAdversarial(&m,endangered)<=2{
    m.query(endangered);actions++;steps[7]=start;eighth=start
    prev=up223cSnapshot(&m,endangered,start);prevValid=true
   }
  }else if actions==8&&eighth>=0&&start>eighth&&!out.EventPresent{
   cur:=up223cSnapshot(&m,endangered,start)
   if up161cAdversarial(&m,endangered)<=2{
    out.EventPresent=true
    out.CurrentSignature=up231cSignature(cur)
    if prevValid{out.TrajectorySignature=out.CurrentSignature+"|"+up237cDelta(cur,prev)}
   }else{prev=cur;prevValid=true}
  }
  for j:=0;j<2&&start+j<=80;j++{
   step:=start+j
   if up165cRealStep(&m,endangered,1057000+step,step,up210cPolicy(step,advance,a,b)){out.Outcome="loss";return out,nil}
  }
 }
 return out,nil
}
func RunUP237C()(UP237CResult,error){
 cohorts:=[][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}}
 hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
 type sched struct{a,b string};schedules:=[]sched{{"no_refresh","hostile_shield"},{"hostile_shield","no_refresh"},{"alternating_shield","fixed_offset_refresh"},{"fixed_offset_refresh","alternating_shield"}}
 r:=UP237CResult{Schema:UP237CTrajectorySchema,Experiment:"UP-237C-critical-trajectory-diagnostic",SourceUP236CSeal:"397f7a0b3e75d7c8519e9575b5459d4d7eeec9d3",InterventionChanged:false,ThresholdFittingUsed:false,ClassifierTrainingUsed:false,AdaptiveFeatureSelectionUsed:false,NewNativeFieldUsed:false,LiveActivation:false}
 cf,cs,tf,ts:=map[string]bool{},map[string]bool{},map[string]bool{},map[string]bool{}
 for _,s:=range schedules{for _,adv:=range []int{0,22}{for _,c:=range cohorts{for _,hand:=range hands{
  x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
  a,err:=up237cRun(x,e,adv,s.a,s.b);if err!=nil{return UP237CResult{},err};r.ArmsTotal++
  fail:=a.Outcome=="loss";if fail{r.Failures++}else{r.Survivors++}
  if a.EventPresent{
   if fail{r.EventFailures++;cf[a.CurrentSignature]=true;if a.TrajectorySignature!=""{tf[a.TrajectorySignature]=true}
   }else{r.EventSurvivors++;cs[a.CurrentSignature]=true;if a.TrajectorySignature!=""{ts[a.TrajectorySignature]=true}}
  }else if fail{r.SilentFailures++}
  r.Arms=append(r.Arms,a)
 }}}}
 r.CurrentFailureSignatures=len(cf);r.CurrentSurvivorSignatures=len(cs);for k:=range cf{if cs[k]{r.CurrentSharedSignatures++}}
 r.TrajectoryFailureSignatures=len(tf);r.TrajectorySurvivorSignatures=len(ts);for k:=range tf{if ts[k]{r.TrajectorySharedSignatures++}}
 return r,nil
}
