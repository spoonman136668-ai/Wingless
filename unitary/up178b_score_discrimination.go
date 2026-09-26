package unitary

import "sort"

const UP178BDiscriminationSchema="wingless.up178b-score-discrimination.v1"

type UP178BMetric struct{
	Predictor string `json:"predictor"`
	AUROC float64 `json:"auroc"`
	AveragePrecision float64 `json:"average_precision"`
	BaseCrossingRate float64 `json:"base_crossing_rate"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	ActualCrossings int `json:"actual_crossings"`
	TotalSlots int `json:"total_slots"`
}
type UP178BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP177BSeal string `json:"source_up177b_seal"`
	ContextPhase int `json:"context_phase"`
	PredictorsFrozenBeforeRun bool `json:"predictors_frozen_before_run"`
	ThresholdTuningUsed bool `json:"threshold_tuning_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP178BMetric `json:"metrics"`
}
type up178bObs struct{score float64;positive bool}
type up178bBin struct{pos,neg int}
func up178bDisc(name string,obs []up178bObs)UP178BMetric{
	bins:=map[float64]*up178bBin{};pos,neg:=0,0;pred:=0.0
	for _,o:=range obs{
		if bins[o.score]==nil{bins[o.score]=&up178bBin{}}
		if o.positive{bins[o.score].pos++;pos++}else{bins[o.score].neg++;neg++}
		pred+=o.score
	}
	scores:=make([]float64,0,len(bins));for s:=range bins{scores=append(scores,s)}
	sort.Float64s(scores)
	cumNeg:=0;u:=0.0
	for _,s:=range scores{b:=bins[s];u+=float64(b.pos*cumNeg)+0.5*float64(b.pos*b.neg);cumNeg+=b.neg}
	auc:=0.0;if pos>0&&neg>0{auc=u/float64(pos*neg)}
	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))
	cumPos,cumAll:=0,0;prevRecall,ap:=0.0,0.0
	for _,s:=range scores{
		b:=bins[s];cumPos+=b.pos;cumAll+=b.pos+b.neg
		recall:=0.0;if pos>0{recall=float64(cumPos)/float64(pos)}
		precision:=0.0;if cumAll>0{precision=float64(cumPos)/float64(cumAll)}
		ap+=(recall-prevRecall)*precision;prevRecall=recall
	}
	base:=0.0;if len(obs)>0{base=float64(pos)/float64(len(obs))}
	return UP178BMetric{Predictor:name,AUROC:auc,AveragePrecision:ap,BaseCrossingRate:base,PredictedCrossingsSum:pred,ActualCrossings:pos,TotalSlots:len(obs)}
}
func RunUP178B()(UP178BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=23
	paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
	composite:=make([]up178bObs,0,17280);rankOnly:=make([]up178bObs,0,17280)
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r)
		for _,path:=range paths{
			g:=*base;streak:=make([]int,120)
			for _,idx:=range path{
				before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls);ranks:=make([]int,len(before))
				for i:=range before{ranks[i]=up166bRank(before,i);if ranks[i]<=band{streak[i]++}else{streak[i]=0}}
				items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
				for i:=range before{
					rs:=up175bRank(ranks[i]);ts:=up175bTemporal(streak[i]);changed:=before[i].correct!=after[i].correct
					composite=append(composite,up178bObs{score:up175bP(cls,rs,ts),positive:changed})
					rankOnly=append(rankOnly,up178bObs{score:up176bRankOnlyP(cls,rs),positive:changed})
				}
			}
		}
	}
	return UP178BResult{Schema:UP178BDiscriminationSchema,Experiment:"UP-178B-score-discrimination",SourceUP177BSeal:"c49343cffe447e23f4932845d7ee792d2e28b0fd",ContextPhase:23,PredictorsFrozenBeforeRun:true,ThresholdTuningUsed:false,MaintenanceTriggered:false,Metrics:[]UP178BMetric{up178bDisc("rank_plus_temporal",composite),up178bDisc("rank_only",rankOnly)}},nil
}
