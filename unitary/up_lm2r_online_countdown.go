package unitary

import "fmt"

const UPLM2ROnlineSchema="wingless.up-lm2r-online-countdown.v1"

type UPLM2RPoint struct{
 DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
 IdentityRotation int `json:"identity_rotation"`
 ValueShift int `json:"value_shift"`
 EventOrdinal int `json:"event_ordinal"`
 EventKind string `json:"event_kind"`
 PredictedUniqueStores int `json:"predicted_unique_stores"`
 ActualUniqueStores int `json:"actual_unique_stores"`
 Exact bool `json:"exact"`
 NearTermWarning bool `json:"near_term_warning"`
 WarningRetraction bool `json:"warning_retraction"`
}
type UPLM2RResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUPLM2QSeal string `json:"source_up_lm2q_seal"`
 ExactRecallCap int `json:"exact_recall_cap"`
 DeferredLevels []int `json:"deferred_levels"`
 IdentityRotations []int `json:"identity_rotations"`
 ValueShifts []int `json:"value_shifts"`
 EventsPerArm int `json:"events_per_arm"`
 NearTermThreshold int `json:"near_term_threshold"`
 FutureOracleUsed bool `json:"future_oracle_used"`
 InterventionTriggered bool `json:"intervention_triggered"`
 ExactPredictions int `json:"exact_predictions"`
 TotalPredictions int `json:"total_predictions"`
 ExactPredictionRate float64 `json:"exact_prediction_rate"`
 NearTermWarnings int `json:"near_term_warnings"`
 WarningRetractions int `json:"warning_retractions"`
 Points []UPLM2RPoint `json:"points"`
}
func uplm2rCountdown(r *uplm0cRecall,reported map[string]bool,tag string)(int,int){
 target,idx:=uplm2qTarget(r,reported);_ = target
 pred:=-1;if idx>=0{pred=16-len(r.order)+idx+1}
 actual:=uplm2qActual(r,reported,tag)
 return pred,actual
}
func uplm2rRun(d,rot,shift int)[]UPLM2RPoint{
 r:=newUPLM0CRecall();reported:=map[string]bool{}
 for i:=0;i<12;i++{r.write(uplm2nName(i,rot),uplm2nValue(i,rot,shift))}
 for i:=0;i<12-d;i++{reported[uplm2nName(i,rot)]=true}
 pending:=[]string{}
 for i:=12-d;i<12;i++{pending=append(pending,uplm2nName(i,rot))}
 storeCount:=0;closeCount:=0;prevWarn:=false
 out:=[]UPLM2RPoint{}
 kinds:=[]string{"STORE","REPORT","STORE","STORE","REPORT","STORE","STORE","REPORT","STORE","STORE","STORE","REPORT","STORE","STORE","STORE","STORE"}
 for e,kind:=range kinds{
  if kind=="STORE"{
   storeCount++;r.write(fmt.Sprintf("online-%d-%d-%d-%d",d,rot,shift,storeCount),"x")
  }else{
   if closeCount<len(pending){reported[pending[closeCount]]=true};closeCount++
  }
  pred,actual:=uplm2rCountdown(r,reported,fmt.Sprintf("%d-%d-%d-%d",d,rot,shift,e+1))
  warn:=pred>=1&&pred<=4
  retract:=kind=="REPORT"&&prevWarn&&!warn
  out=append(out,UPLM2RPoint{DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,EventOrdinal:e+1,EventKind:kind,PredictedUniqueStores:pred,ActualUniqueStores:actual,Exact:pred==actual,NearTermWarning:warn,WarningRetraction:retract})
  prevWarn=warn
 }
 return out
}
func RunUPLM2R()(UPLM2RResult,error){
 levels:=[]int{4,5,6};rots:=[]int{0,7};shifts:=[]int{0,1,2,3}
 res:=UPLM2RResult{Schema:UPLM2ROnlineSchema,Experiment:"UP-LM2R-online-countdown",SourceUPLM2QSeal:"9a7ac068b8a685dbfd44b087771e3cf5ce46fd28",ExactRecallCap:16,DeferredLevels:levels,IdentityRotations:rots,ValueShifts:shifts,EventsPerArm:16,NearTermThreshold:4,FutureOracleUsed:false,InterventionTriggered:false}
 for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{for _,p:=range uplm2rRun(d,rot,shift){res.Points=append(res.Points,p);if p.Exact{res.ExactPredictions++};if p.NearTermWarning{res.NearTermWarnings++};if p.WarningRetraction{res.WarningRetractions++}}}}}
 res.TotalPredictions=len(res.Points);res.ExactPredictionRate=float64(res.ExactPredictions)/float64(res.TotalPredictions)
 return res,nil
}
