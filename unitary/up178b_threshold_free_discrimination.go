package unitary

import "sort"

const UP178BDiscriminationSchema="wingless.up178b-threshold-free-discrimination.v1"

type up178bObs struct{score float64; y int}
type UP178BMetric struct{
 Predictor string `json:"predictor"`
 AUROC float64 `json:"auroc"`
 AUPRC float64 `json:"auprc"`
 BrierScore float64 `json:"brier_score"`
 BaseCrossingRate float64 `json:"base_crossing_rate"`
 MaxScore float64 `json:"max_score"`
 MaxScoreSlots int `json:"max_score_slots"`
 MaxScoreCrossings int `json:"max_score_crossings"`
 MaxScorePositiveRate float64 `json:"max_score_positive_rate"`
 MaxScoreEnrichment float64 `json:"max_score_enrichment"`
}
type UP178BResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP177BSeal string `json:"source_up177b_seal"`
 ContextPhase int `json:"context_phase"`
 TotalSlots int `json:"total_slots"`
 ActualCrossings int `json:"actual_crossings"`
 PredictorsFrozenBeforeRun bool `json:"predictors_frozen_before_run"`
 InterventionTriggered bool `json:"intervention_triggered"`
 Metrics []UP178BMetric `json:"metrics"`
}
func up178bMetric(name string,obs []up178bObs)UP178BMetric{
 sort.Slice(obs,func(i,j int)bool{return obs[i].score>obs[j].score})
 pos:=0;brier:=0.0;for _,o:=range obs{pos+=o.y;d:=o.score-float64(o.y);brier+=d*d}
 neg:=len(obs)-pos;base:=float64(pos)/float64(len(obs))
 tp,fp:=0,0;prevTPR,prevFPR:=0.0,0.0;auroc,auprc,prevRecall:=0.0,0.0,0.0
 maxScore:=obs[0].score;maxSlots,maxCross:=0,0
 for i:=0;i<len(obs);{
  j:=i;gPos:=0
  score:=obs[i].score
  for j<len(obs)&&obs[j].score==score{gPos+=obs[j].y;j++}
  gN:=j-i;tp+=gPos;fp+=gN-gPos
  tpr,fpr:=0.0,0.0;if pos>0{tpr=float64(tp)/float64(pos)};if neg>0{fpr=float64(fp)/float64(neg)}
  auroc+=(fpr-prevFPR)*(tpr+prevTPR)/2.0
  precision:=0.0;if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
  auprc+=(tpr-prevRecall)*precision
  prevTPR,prevFPR,prevRecall=tpr,fpr,tpr
  if score==maxScore{maxSlots=gN;maxCross=gPos}
  i=j
 }
 maxRate:=0.0;if maxSlots>0{maxRate=float64(maxCross)/float64(maxSlots)}
 enrich:=0.0;if base>0{enrich=maxRate/base}
 return UP178BMetric{Predictor:name,AUROC:auroc,AUPRC:auprc,BrierScore:brier/float64(len(obs)),BaseCrossingRate:base,MaxScore:maxScore,MaxScoreSlots:maxSlots,MaxScoreCrossings:maxCross,MaxScorePositiveRate:maxRate,MaxScoreEnrichment:enrich}
}
func RunUP178B()(UP178BResult,error){
 o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=23
 paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
 comp:=[]up178bObs{};rank:=[]up178bObs{};total,cross:=0,0
 for _,subject:=range up130bNewNames{
  base:=up169bPostReport(common,subject,epoch,o,r)
  for _,path:=range paths{
   g:=*base;streak:=make([]int,120)
   for _,idx:=range path{
    before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls);ranks:=make([]int,len(before))
    for i:=range before{ranks[i]=up166bRank(before,i);if ranks[i]<=band{streak[i]++}else{streak[i]=0}}
    items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
    for i:=range before{
     rs:=up175bRank(ranks[i]);ts:=up175bTemporal(streak[i]);y:=0;if before[i].correct!=after[i].correct{y=1;cross++};total++
     comp=append(comp,up178bObs{score:up175bP(cls,rs,ts),y:y})
     rank=append(rank,up178bObs{score:up176bRankOnlyP(cls,rs),y:y})
    }
   }
  }
 }
 return UP178BResult{Schema:UP178BDiscriminationSchema,Experiment:"UP-178B-threshold-free-discrimination",SourceUP177BSeal:"c49343cffe447e23f4932845d7ee792d2e28b0fd",ContextPhase:23,TotalSlots:total,ActualCrossings:cross,PredictorsFrozenBeforeRun:true,InterventionTriggered:false,Metrics:[]UP178BMetric{up178bMetric("rank_plus_temporal",comp),up178bMetric("rank_only",rank)}},nil
}
