package unitary

func up234cL1(a,b up233cVec)int{
 d:=0
 for i:=0;i<9;i++{x:=a[i]-b[i];if x<0{x=-x};d+=x}
 return d
}
func up234cMin(x up233cVec,ps []up233cVec,metric string)int{
 m:=1<<30
 for _,p:=range ps{
  d:=up233cDist(x,p);if metric=="manhattan"{d=up234cL1(x,p)}
  if d<m{m=d}
 }
 return m
}
type UP234CSummary struct{
 Metric string `json:"metric"`
 Arms int `json:"arms"`; Failures int `json:"failures"`; Survivors int `json:"survivors"`
 EventFailures int `json:"event_failures"`; EventSurvivors int `json:"event_survivors"`
 HighRisk int `json:"high_risk"`; LowRisk int `json:"low_risk"`; Unknown int `json:"unknown"`; NoEvent int `json:"no_event"`
 HighRiskFailures int `json:"high_risk_failures"`; HighRiskSurvivors int `json:"high_risk_survivors"`
 LowRiskFailures int `json:"low_risk_failures"`; LowRiskSurvivors int `json:"low_risk_survivors"`
 UnknownFailures int `json:"unknown_failures"`; UnknownSurvivors int `json:"unknown_survivors"`
 HighRiskFailureRecall float64 `json:"high_risk_failure_recall"`; HighRiskPrecision float64 `json:"high_risk_precision"`
}
type UP234CResult struct{
 Schema string `json:"schema"`; Experiment string `json:"experiment"`; SourceUP233CSeal string `json:"source_up233c_seal"`
 HeldoutFittingUsed bool `json:"heldout_fitting_used"`; NewNativeFieldUsed bool `json:"new_native_field_used"`
 LearnedScalingUsed bool `json:"learned_scaling_used"`; LearnedWeightsUsed bool `json:"learned_weights_used"`
 ThresholdFittingUsed bool `json:"threshold_fitting_used"`; InterventionChanged bool `json:"intervention_changed"`; LiveActivation bool `json:"live_activation"`
 Summaries []UP234CSummary `json:"summaries"`
}
func RunUP234C()(UP234CResult,error){
 failEx:=up233cExemplars(up232cFailureSignatures);survEx:=up233cExemplars(up232cSurvivorSignatures)
 cohorts:=[][]int{{0,4,10,14},{1,5,11,15},{2,6,8,12},{3,7,9,13}}
 hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
 type sched struct{a,b string};schedules:=[]sched{{"no_refresh","hostile_shield"},{"hostile_shield","no_refresh"},{"alternating_shield","fixed_offset_refresh"},{"fixed_offset_refresh","alternating_shield"}}
 metrics:=[]string{"hamming","manhattan"}
 r:=UP234CResult{Schema:"wingless.up234c-distance-metric-transfer.v1",Experiment:"UP-234C-distance-metric-transfer",SourceUP233CSeal:"016fe588ae331a4456b35bc094b8c56c5d147d17",HeldoutFittingUsed:false,NewNativeFieldUsed:false,LearnedScalingUsed:false,LearnedWeightsUsed:false,ThresholdFittingUsed:false,InterventionChanged:false,LiveActivation:false}
 for _,metric:=range metrics{
  s:=UP234CSummary{Metric:metric}
  for _,pol:=range schedules{for _,adv:=range []int{4,18}{for _,c:=range cohorts{for _,hand:=range hands{
   x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
   a,err:=up231cRun(x,e,adv,pol.a,pol.b);if err!=nil{return UP234CResult{},err}
   fail:=a.Outcome=="loss";s.Arms++;if fail{s.Failures++}else{s.Survivors++}
   if !a.EventPresent{s.NoEvent++;continue}
   if fail{s.EventFailures++}else{s.EventSurvivors++}
   z:=up233cFromSnapshot(a.Snapshot);fd:=up234cMin(z,failEx,metric);sd:=up234cMin(z,survEx,metric)
   if fd<sd{s.HighRisk++;if fail{s.HighRiskFailures++}else{s.HighRiskSurvivors++}
   }else if sd<fd{s.LowRisk++;if fail{s.LowRiskFailures++}else{s.LowRiskSurvivors++}
   }else{s.Unknown++;if fail{s.UnknownFailures++}else{s.UnknownSurvivors++}}
  }}}}
  if s.Failures>0{s.HighRiskFailureRecall=float64(s.HighRiskFailures)/float64(s.Failures)}
  if s.HighRisk>0{s.HighRiskPrecision=float64(s.HighRiskFailures)/float64(s.HighRisk)}
  r.Summaries=append(r.Summaries,s)
 }
 return r,nil
}
