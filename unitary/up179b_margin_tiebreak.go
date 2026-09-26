package unitary

import "math"

const UP179BMarginTieSchema="wingless.up179b-margin-tiebreak.v1"

type UP179BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP178BSeal string `json:"source_up178b_seal"`
	ContextPhase int `json:"context_phase"`
	TieBreakEpsilon float64 `json:"tie_break_epsilon"`
	PredictorsFrozenBeforeRun bool `json:"predictors_frozen_before_run"`
	ThresholdTuningUsed bool `json:"threshold_tuning_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP178BMetric `json:"metrics"`
}
func RunUP179B()(UP179BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=24;eps:=1e-6
	paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
	baseObs:=make([]up178bObs,0,17280);refinedObs:=make([]up178bObs,0,17280);marginObs:=make([]up178bObs,0,17280)
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r)
		for _,path:=range paths{
			g:=*base;streak:=make([]int,120)
			for _,idx:=range path{
				before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls);ranks:=make([]int,len(before))
				for i:=range before{ranks[i]=up166bRank(before,i);if ranks[i]<=band{streak[i]++}else{streak[i]=0}}
				items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
				for i:=range before{
					rs:=up175bRank(ranks[i]);ts:=up175bTemporal(streak[i]);p:=up175bP(cls,rs,ts);changed:=before[i].correct!=after[i].correct;am:=math.Abs(before[i].margin)
					baseObs=append(baseObs,up178bObs{score:p,positive:changed})
					refinedObs=append(refinedObs,up178bObs{score:p+eps/(1.0+am),positive:changed})
					marginObs=append(marginObs,up178bObs{score:-am,positive:changed})
				}
			}
		}
	}
	return UP179BResult{Schema:UP179BMarginTieSchema,Experiment:"UP-179B-margin-tiebreak",SourceUP178BSeal:"65fa1a872cbaa372dbd72df9c7c7df9ab0269e30",ContextPhase:24,TieBreakEpsilon:eps,PredictorsFrozenBeforeRun:true,ThresholdTuningUsed:false,MaintenanceTriggered:false,Metrics:[]UP178BMetric{up178bDisc("rank_plus_temporal",baseObs),up178bDisc("rank_plus_temporal_margin_tiebreak",refinedObs),up178bDisc("margin_only",marginObs)}},nil
}
