package unitary

import "fmt"

const UPLM2SClosureSchema="wingless.up-lm2s-bounded-early-closure.v1"

type UPLM2SPoint struct{
 DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
 IdentityRotation int `json:"identity_rotation"`
 ValueShift int `json:"value_shift"`
 BaselineUnreportedEvictions int `json:"baseline_unreported_evictions"`
 Interventions int `json:"interventions"`
 CounterfactualUnreportedEvictions int `json:"counterfactual_unreported_evictions"`
 PreventedEvictions int `json:"prevented_evictions"`
 UnnecessaryInterventions int `json:"unnecessary_interventions"`
 MovedReports int `json:"moved_reports"`
 RemainingDeferredReports int `json:"remaining_deferred_reports"`
}
type UPLM2SResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUPLM2RSeal string `json:"source_up_lm2r_seal"`
 ExactRecallCap int `json:"exact_recall_cap"`
 DeferredLevels []int `json:"deferred_levels"`
 IdentityRotations []int `json:"identity_rotations"`
 ValueShifts []int `json:"value_shifts"`
 CounterfactualOnly bool `json:"counterfactual_only"`
 LiveActivation bool `json:"live_activation"`
 CapacityChanged bool `json:"capacity_changed"`
 ExtraTrainingUsed bool `json:"extra_training_used"`
 BaselineHarmfulEvictions int `json:"baseline_harmful_evictions"`
 TotalInterventions int `json:"total_interventions"`
 CounterfactualHarmfulEvictions int `json:"counterfactual_harmful_evictions"`
 PreventedEvictions int `json:"prevented_evictions"`
 UnnecessaryInterventions int `json:"unnecessary_interventions"`
 PreventedFraction float64 `json:"prevented_fraction"`
 InterventionPrecision float64 `json:"intervention_precision"`
 InterventionRecall float64 `json:"intervention_recall"`
 Points []UPLM2SPoint `json:"points"`
}
func uplm2sRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func uplm2sInit(d,rot,shift int)(*uplm0cRecall,map[string]bool){
 r:=newUPLM0CRecall();reported:=map[string]bool{}
 for i:=0;i<12;i++{r.write(uplm2nName(i,rot),uplm2nValue(i,rot,shift))}
 for i:=0;i<12-d;i++{reported[uplm2nName(i,rot)]=true}
 return r,reported
}
func uplm2sRun(d,rot,shift int)UPLM2SPoint{
 base,baseReported:=uplm2sInit(d,rot,shift)
 baselineHarm:=0
 for j:=0;j<12;j++{
  if len(base.order)>=16&&!baseReported[base.order[0]]{baselineHarm++}
  base.write(fmt.Sprintf("base-%d-%d-%d-%d",d,rot,shift,j),"x")
 }

 cf,reported:=uplm2sInit(d,rot,shift)
 interventions,harm,unnecessary:=0,0,0
 moved:=map[string]bool{}
 for j:=0;j<12;j++{
  if len(cf.order)>=16{
   victim:=cf.order[0]
   if !reported[victim]{
    interventions++
    if _,ok:=cf.values[victim];!ok{unnecessary++}
    reported[victim]=true
    moved[victim]=true
   }
  }
  if len(cf.order)>=16&&!reported[cf.order[0]]{harm++}
  cf.write(fmt.Sprintf("cf-%d-%d-%d-%d",d,rot,shift,j),"x")
 }
 prevented:=baselineHarm-harm
 return UPLM2SPoint{
  DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,
  BaselineUnreportedEvictions:baselineHarm,Interventions:interventions,
  CounterfactualUnreportedEvictions:harm,PreventedEvictions:prevented,
  UnnecessaryInterventions:unnecessary,MovedReports:len(moved),
  RemainingDeferredReports:d-len(moved),
 }
}
func RunUPLM2S()(UPLM2SResult,error){
 levels:=[]int{4,5,6};rots:=[]int{0,7};shifts:=[]int{0,1,2,3}
 res:=UPLM2SResult{Schema:UPLM2SClosureSchema,Experiment:"UP-LM2S-bounded-early-closure",SourceUPLM2RSeal:"4e254762af5e380cbaf689002c884d5574512888",ExactRecallCap:16,DeferredLevels:levels,IdentityRotations:rots,ValueShifts:shifts,CounterfactualOnly:true,LiveActivation:false,CapacityChanged:false,ExtraTrainingUsed:false}
 for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{
  p:=uplm2sRun(d,rot,shift);res.Points=append(res.Points,p)
  res.BaselineHarmfulEvictions+=p.BaselineUnreportedEvictions
  res.TotalInterventions+=p.Interventions
  res.CounterfactualHarmfulEvictions+=p.CounterfactualUnreportedEvictions
  res.PreventedEvictions+=p.PreventedEvictions
  res.UnnecessaryInterventions+=p.UnnecessaryInterventions
 }}}
 res.PreventedFraction=uplm2sRate(res.PreventedEvictions,res.BaselineHarmfulEvictions)
 res.InterventionPrecision=uplm2sRate(res.PreventedEvictions,res.TotalInterventions)
 res.InterventionRecall=uplm2sRate(res.PreventedEvictions,res.BaselineHarmfulEvictions)
 return res,nil
}
