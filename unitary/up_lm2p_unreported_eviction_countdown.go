package unitary

import "fmt"

const UPLM2PCountdownSchema="wingless.up-lm2p-unreported-eviction-countdown.v1"

type UPLM2PPoint struct{
 DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
 IdentityRotation int `json:"identity_rotation"`
 ValueShift int `json:"value_shift"`
 TargetName string `json:"target_name"`
 TargetFIFOIndex int `json:"target_fifo_index"`
 FreeSlots int `json:"free_slots"`
 PredictedUniqueStores int `json:"predicted_unique_stores"`
 ActualUniqueStores int `json:"actual_unique_stores"`
 Exact bool `json:"exact"`
}
type UPLM2PResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUPLM2OSeal string `json:"source_up_lm2o_seal"`
 ExactRecallCap int `json:"exact_recall_cap"`
 DeferredLevels []int `json:"deferred_levels"`
 IdentityRotations []int `json:"identity_rotations"`
 ValueShifts []int `json:"value_shifts"`
 FutureOracleUsed bool `json:"future_oracle_used"`
 InterventionTriggered bool `json:"intervention_triggered"`
 ExactPredictions int `json:"exact_predictions"`
 TotalPredictions int `json:"total_predictions"`
 ExactPredictionRate float64 `json:"exact_prediction_rate"`
 Points []UPLM2PPoint `json:"points"`
}
func uplm2pRun(d,rot,shift int)UPLM2PPoint{
 recall:=newUPLM0CRecall();reported:=map[string]bool{}
 for i:=0;i<12;i++{recall.write(uplm2nName(i,rot),uplm2nValue(i,rot,shift))}
 for i:=0;i<12-d;i++{reported[uplm2nName(i,rot)]=true}
 target:="";idx:=-1
 for i,n:=range recall.order{if !reported[n]{target=n;idx=i;break}}
 free:=16-len(recall.order);pred:=free+idx+1;actual:=-1
 for j:=1;j<=13;j++{
  recall.write(fmt.Sprintf("probe-%d-%d-%d-%d",d,rot,shift,j),"x")
  if _,ok:=recall.values[target];!ok{actual=j;break}
 }
 return UPLM2PPoint{DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,TargetName:target,TargetFIFOIndex:idx,FreeSlots:free,PredictedUniqueStores:pred,ActualUniqueStores:actual,Exact:pred==actual}
}
func RunUPLM2P()(UPLM2PResult,error){
 levels:=[]int{4,5,6};rots:=[]int{0,7};shifts:=[]int{0,1,2,3}
 res:=UPLM2PResult{Schema:UPLM2PCountdownSchema,Experiment:"UP-LM2P-unreported-eviction-countdown",SourceUPLM2OSeal:"a10aeb370020820e12636dd5182e91d48287288e",ExactRecallCap:16,DeferredLevels:levels,IdentityRotations:rots,ValueShifts:shifts,FutureOracleUsed:false,InterventionTriggered:false}
 for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{p:=uplm2pRun(d,rot,shift);res.Points=append(res.Points,p);if p.Exact{res.ExactPredictions++}}}}
 res.TotalPredictions=len(res.Points);res.ExactPredictionRate=float64(res.ExactPredictions)/float64(res.TotalPredictions)
 return res,nil
}
