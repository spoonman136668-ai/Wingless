package unitary

const UP177BFrontierSchema="wingless.up177b-threshold-frontier.v1"

type UP177BMetric struct{
 Threshold float64 `json:"threshold"`
 TruePositive int `json:"true_positive"`
 FalsePositive int `json:"false_positive"`
 FalseNegative int `json:"false_negative"`
 TrueNegative int `json:"true_negative"`
 WarningCount int `json:"warning_count"`
 Precision float64 `json:"precision"`
 Recall float64 `json:"recall"`
}
type UP177BResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP176BSeal string `json:"source_up176b_seal"`
 ContextPhase int `json:"context_phase"`
 Thresholds []float64 `json:"thresholds"`
 TotalSlots int `json:"total_slots"`
 ActualCrossings int `json:"actual_crossings"`
 PredictorFrozenBeforeRun bool `json:"predictor_frozen_before_run"`
 AdaptiveThresholdUsed bool `json:"adaptive_threshold_used"`
 MaintenanceTriggered bool `json:"maintenance_triggered"`
 Metrics []UP177BMetric `json:"metrics"`
}
type up177bAgg struct{tp,fp,fn,tn int}
func up177bRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func RunUP177B()(UP177BResult,error){
 o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=22
 thresholds:=[]float64{0.20,0.30,0.40};aggs:=make([]up177bAgg,len(thresholds));total,cross:=0,0
 paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
 for _,subject:=range up130bNewNames{
  base:=up169bPostReport(common,subject,epoch,o,r)
  for _,path:=range paths{
   g:=*base;streak:=make([]int,120)
   for _,idx:=range path{
    before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls);ranks:=make([]int,len(before))
    for i:=range before{ranks[i]=up166bRank(before,i);if ranks[i]<=band{streak[i]++}else{streak[i]=0}}
    items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
    for i:=range before{
     rs:=up175bRank(ranks[i]);ts:=up175bTemporal(streak[i]);p:=up175bP(cls,rs,ts);changed:=before[i].correct!=after[i].correct;total++;if changed{cross++}
     for ti,t:=range thresholds{
      warn:=p>=t;a:=&aggs[ti]
      if warn&&changed{a.tp++}else if warn&&!changed{a.fp++}else if !warn&&changed{a.fn++}else{a.tn++}
     }
    }
   }
  }
 }
 res:=UP177BResult{Schema:UP177BFrontierSchema,Experiment:"UP-177B-threshold-frontier",SourceUP176BSeal:"acfe77e44d2082f0db0d85ba9d7deba6f91f84f5",ContextPhase:22,Thresholds:thresholds,TotalSlots:total,ActualCrossings:cross,PredictorFrozenBeforeRun:true,AdaptiveThresholdUsed:false,MaintenanceTriggered:false}
 for i,t:=range thresholds{a:=aggs[i];res.Metrics=append(res.Metrics,UP177BMetric{Threshold:t,TruePositive:a.tp,FalsePositive:a.fp,FalseNegative:a.fn,TrueNegative:a.tn,WarningCount:a.tp+a.fp,Precision:up177bRate(a.tp,a.tp+a.fp),Recall:up177bRate(a.tp,a.tp+a.fn)})}
 return res,nil
}
