package unitary

import "fmt"

const UPLM2QDynamicSchema="wingless.up-lm2q-dynamic-pressure-countdown.v1"

type UPLM2QPoint struct{
 DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
 IdentityRotation int `json:"identity_rotation"`
 ValueShift int `json:"value_shift"`
 Snapshot string `json:"snapshot"`
 RecallEntries int `json:"recall_entries"`
 FreeSlots int `json:"free_slots"`
 TargetName string `json:"target_name"`
 TargetFIFOIndex int `json:"target_fifo_index"`
 PredictedUniqueStores int `json:"predicted_unique_stores"`
 ActualUniqueStores int `json:"actual_unique_stores"`
 Exact bool `json:"exact"`
}
type UPLM2QResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUPLM2PSeal string `json:"source_up_lm2p_seal"`
 ExactRecallCap int `json:"exact_recall_cap"`
 DeferredLevels []int `json:"deferred_levels"`
 IdentityRotations []int `json:"identity_rotations"`
 ValueShifts []int `json:"value_shifts"`
 Snapshots []string `json:"snapshots"`
 FutureOracleUsed bool `json:"future_oracle_used"`
 InterventionTriggered bool `json:"intervention_triggered"`
 ExactPredictions int `json:"exact_predictions"`
 TotalPredictions int `json:"total_predictions"`
 ExactPredictionRate float64 `json:"exact_prediction_rate"`
 Points []UPLM2QPoint `json:"points"`
}
func uplm2qClone(r *uplm0cRecall)*uplm0cRecall{
 c:=newUPLM0CRecall()
 c.order=append([]string(nil),r.order...)
 for k,v:=range r.values{c.values[k]=v}
 return c
}
func uplm2qTarget(r *uplm0cRecall,reported map[string]bool)(string,int){
 for i,n:=range r.order{if !reported[n]{return n,i}}
 return "",-1
}
func uplm2qActual(r *uplm0cRecall,reported map[string]bool,tag string)int{
 c:=uplm2qClone(r)
 for j:=1;j<=64;j++{
  if len(c.order)>=16{
   victim:=c.order[0]
   if !reported[victim]{return j}
  }
  c.write(fmt.Sprintf("future-%s-%d",tag,j),"x")
 }
 return -1
}
func uplm2qSnapshot(d,rot,shift int,label string,r *uplm0cRecall,reported map[string]bool)UPLM2QPoint{
 target,idx:=uplm2qTarget(r,reported);free:=16-len(r.order);pred:=-1
 if idx>=0{pred=free+idx+1}
 actual:=uplm2qActual(r,reported,fmt.Sprintf("%d-%d-%d-%s",d,rot,shift,label))
 return UPLM2QPoint{DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,Snapshot:label,RecallEntries:len(r.order),FreeSlots:free,TargetName:target,TargetFIFOIndex:idx,PredictedUniqueStores:pred,ActualUniqueStores:actual,Exact:pred==actual}
}
func uplm2qRun(d,rot,shift int)[]UPLM2QPoint{
 r:=newUPLM0CRecall();reported:=map[string]bool{}
 for i:=0;i<12;i++{r.write(uplm2nName(i,rot),uplm2nValue(i,rot,shift))}
 for i:=0;i<12-d;i++{reported[uplm2nName(i,rot)]=true}
 firstPending:=uplm2nName(12-d,rot);secondPending:=uplm2nName(13-d,rot)
 out:=[]UPLM2QPoint{uplm2qSnapshot(d,rot,shift,"initial",r,reported)}
 for j:=1;j<=4;j++{r.write(fmt.Sprintf("pressure-a-%d-%d-%d-%d",d,rot,shift,j),"x")}
 reported[firstPending]=true
 out=append(out,uplm2qSnapshot(d,rot,shift,"after_store4_close1",r,reported))
 for j:=5;j<=8;j++{r.write(fmt.Sprintf("pressure-b-%d-%d-%d-%d",d,rot,shift,j),"x")}
 reported[secondPending]=true
 out=append(out,uplm2qSnapshot(d,rot,shift,"after_store8_close2",r,reported))
 return out
}
func RunUPLM2Q()(UPLM2QResult,error){
 levels:=[]int{4,5,6};rots:=[]int{0,7};shifts:=[]int{0,1,2,3};snaps:=[]string{"initial","after_store4_close1","after_store8_close2"}
 res:=UPLM2QResult{Schema:UPLM2QDynamicSchema,Experiment:"UP-LM2Q-dynamic-pressure-countdown",SourceUPLM2PSeal:"bcd8b5ea2a5c6c8c6ed9e5bcaad081e561227a03",ExactRecallCap:16,DeferredLevels:levels,IdentityRotations:rots,ValueShifts:shifts,Snapshots:snaps,FutureOracleUsed:false,InterventionTriggered:false}
 for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{for _,p:=range uplm2qRun(d,rot,shift){res.Points=append(res.Points,p);if p.Exact{res.ExactPredictions++}}}}}
 res.TotalPredictions=len(res.Points);res.ExactPredictionRate=float64(res.ExactPredictions)/float64(res.TotalPredictions)
 return res,nil
}
