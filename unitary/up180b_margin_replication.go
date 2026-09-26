package unitary

import "math"

const UP180BMarginReplicationSchema="wingless.up180b-margin-replication.v1"

type UP180BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP179BSeal string `json:"source_up179b_seal"`
	ContextPhase int `json:"context_phase"`
	PredictorsFrozenBeforeRun bool `json:"predictors_frozen_before_run"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP178BMetric `json:"metrics"`
}
func RunUP180B()(UP180BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=25;eps:=1e-6
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
					refinedObs=append(refinedObs,up178bObs{score:p+eps/(1+am),positive:changed})
					marginObs=append(marginObs,up178bObs{score:-am,positive:changed})
				}
			}
		}
	}
	return UP180BResult{Schema:UP180BMarginReplicationSchema,Experiment:"UP-180B-margin-replication",SourceUP179BSeal:"3add97cb27d37b2818a56312ca3b1807d5b3f59a",ContextPhase:25,PredictorsFrozenBeforeRun:true,MaintenanceTriggered:false,Metrics:[]UP178BMetric{up178bDisc("rank_plus_temporal",baseObs),up178bDisc("rank_plus_temporal_margin_tiebreak",refinedObs),up178bDisc("margin_only",marginObs)}},nil
}
