package unitary

const UP176BAblationSchema="wingless.up176b-temporal-predictor-ablation.v1"

type UP176BPredictorMetric struct{
 Predictor string `json:"predictor"`
 PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
 BrierScore float64 `json:"brier_score"`
 TruePositive int `json:"true_positive"`
 FalsePositive int `json:"false_positive"`
 FalseNegative int `json:"false_negative"`
 TrueNegative int `json:"true_negative"`
 WarningPrecision float64 `json:"warning_precision"`
 WarningRecall float64 `json:"warning_recall"`
}
type UP176BResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP175BSeal string `json:"source_up175b_seal"`
 ContextPhase int `json:"context_phase"`
 WarningThreshold float64 `json:"warning_threshold"`
 TotalSlots int `json:"total_slots"`
 ActualCrossings int `json:"actual_crossings"`
 PredictorProbabilitiesFrozenBeforeRun bool `json:"predictor_probabilities_frozen_before_run"`
 AdaptiveThresholdUsed bool `json:"adaptive_threshold_used"`
 MaintenanceTriggered bool `json:"maintenance_triggered"`
 Metrics []UP176BPredictorMetric `json:"metrics"`
}
type up176bAgg struct{pred,brier float64;tp,fp,fn,tn int}
func up176bRankOnlyP(cls,rankState string)float64{
 bg:=0.00043793793793793793
 if rankState=="background"{return bg}
 switch cls{
 case "STORE":if rankState=="rank_1_4"{return 80.0/240.0};return 24.0/360.0
 case "OBSERVE":if rankState=="rank_1_4"{return 67.0/240.0};return 28.0/360.0
 case "REPORT":if rankState=="rank_1_4"{return 16.0/96.0}
 }
 return bg
}
func up176bUpdate(a *up176bAgg,p float64,changed bool,threshold float64){
 y:=0.0;if changed{y=1};a.pred+=p;d:=p-y;a.brier+=d*d;warn:=p>=threshold
 if warn&&changed{a.tp++}else if warn&&!changed{a.fp++}else if !warn&&changed{a.fn++}else{a.tn++}
}
func up176bMetric(name string,a up176bAgg,total int)UP176BPredictorMetric{
 rate:=func(x,y int)float64{if y==0{return 0};return float64(x)/float64(y)}
 return UP176BPredictorMetric{Predictor:name,PredictedCrossingsSum:a.pred,BrierScore:a.brier/float64(total),TruePositive:a.tp,FalsePositive:a.fp,FalseNegative:a.fn,TrueNegative:a.tn,WarningPrecision:rate(a.tp,a.tp+a.fp),WarningRecall:rate(a.tp,a.tp+a.fn)}
}
func RunUP176B()(UP176BResult,error){
 o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=21;threshold:=0.20
 paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
 var composite,rankOnly up176bAgg;total,cross:=0,0
 for _,subject:=range up130bNewNames{
  base:=up169bPostReport(common,subject,epoch,o,r)
  for _,path:=range paths{
   g:=*base;streak:=make([]int,120)
   for _,idx:=range path{
    before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls);ranks:=make([]int,len(before))
    for i:=range before{ranks[i]=up166bRank(before,i);if ranks[i]<=band{streak[i]++}else{streak[i]=0}}
    items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
    for i:=range before{
     rs:=up175bRank(ranks[i]);ts:=up175bTemporal(streak[i]);changed:=before[i].correct!=after[i].correct;total++;if changed{cross++}
     up176bUpdate(&composite,up175bP(cls,rs,ts),changed,threshold)
     up176bUpdate(&rankOnly,up176bRankOnlyP(cls,rs),changed,threshold)
    }
   }
  }
 }
 return UP176BResult{Schema:UP176BAblationSchema,Experiment:"UP-176B-temporal-predictor-ablation",SourceUP175BSeal:"6fc3ee7e5d2535654db6c75b2806b82e60bdfe90",ContextPhase:21,WarningThreshold:threshold,TotalSlots:total,ActualCrossings:cross,PredictorProbabilitiesFrozenBeforeRun:true,AdaptiveThresholdUsed:false,MaintenanceTriggered:false,Metrics:[]UP176BPredictorMetric{up176bMetric("rank_plus_temporal",composite,total),up176bMetric("rank_only",rankOnly,total)}},nil
}
