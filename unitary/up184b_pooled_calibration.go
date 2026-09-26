package unitary

import "math"

const UP184BPooledSchema="wingless.up184b-pooled-calibration.v1"

type up184bModel struct{
	name string
	cells map[string]*up182bCell
	globalP float64
}
type UP184BMetric struct{
	Model string `json:"model"`
	Phase int `json:"phase"`
	ActualCrossings int `json:"actual_crossings"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	PredictedActualRatio float64 `json:"predicted_actual_ratio"`
	AbsoluteMassRatioError float64 `json:"absolute_mass_ratio_error"`
	BrierScore float64 `json:"brier_score"`
	ExpectedCalibrationError float64 `json:"expected_calibration_error"`
	AUROC float64 `json:"auroc"`
	AveragePrecision float64 `json:"average_precision"`
	TotalSlots int `json:"total_slots"`
}
type UP184BModelSummary struct{
	Model string `json:"model"`
	MeanAbsoluteMassRatioError float64 `json:"mean_absolute_mass_ratio_error"`
}
type UP184BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP183BSeal string `json:"source_up183b_seal"`
	SingleCalibrationPhases []int `json:"single_calibration_phases"`
	PooledCalibrationPhases []int `json:"pooled_calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	ModelsFrozenBeforeEvaluation bool `json:"models_frozen_before_evaluation"`
	EvaluationFittingUsed bool `json:"evaluation_fitting_used"`
	PerPhaseScalingUsed bool `json:"per_phase_scaling_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP184BMetric `json:"metrics"`
	Summaries []UP184BModelSummary `json:"summaries"`
}
func up184bBuild(name string,phases []int)up184bModel{
	cells:=map[string]*up182bCell{};pos,total:=0,0
	for _,phase:=range phases{
		for _,r:=range up182bRows(phase){
			key:=up182bKey(r);if cells[key]==nil{cells[key]=&up182bCell{}};cells[key].slots++;total++
			if r.positive{cells[key].cross++;pos++}
		}
	}
	return up184bModel{name:name,cells:cells,globalP:float64(pos+1)/float64(total+2)}
}
func up184bEval(m up184bModel,phase int)UP184BMetric{
	test:=up182bRows(phase);obs:=make([]up178bObs,0,len(test));brier,pred:=0.0,0.0;actual:=0
	type eceBin struct{n,pos int;psum float64};ece:=make([]eceBin,10)
	for _,r:=range test{
		p:=m.globalP;if c:=m.cells[up182bKey(r)];c!=nil{p=float64(c.cross+1)/float64(c.slots+2)}
		y:=0.0;if r.positive{y=1;actual++};diff:=p-y;brier+=diff*diff;pred+=p
		b:=int(p*10);if b>9{b=9};ece[b].n++;ece[b].psum+=p;if r.positive{ece[b].pos++}
		obs=append(obs,up178bObs{score:p,positive:r.positive})
	}
	eceV:=0.0;for _,b:=range ece{if b.n>0{ap:=b.psum/float64(b.n);fr:=float64(b.pos)/float64(b.n);eceV+=float64(b.n)/float64(len(test))*math.Abs(ap-fr)}}
	disc:=up178bDisc(m.name,obs);ratio:=0.0;if actual>0{ratio=pred/float64(actual)}
	return UP184BMetric{Model:m.name,Phase:phase,ActualCrossings:actual,PredictedCrossingsSum:pred,PredictedActualRatio:ratio,AbsoluteMassRatioError:math.Abs(ratio-1),BrierScore:brier/float64(len(test)),ExpectedCalibrationError:eceV,AUROC:disc.AUROC,AveragePrecision:disc.AveragePrecision,TotalSlots:len(test)}
}
func RunUP184B()(UP184BResult,error){
	single:=up184bBuild("single_phase_26",[]int{26})
	pooled:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	phases:=[]int{31,32,33}
	res:=UP184BResult{Schema:UP184BPooledSchema,Experiment:"UP-184B-pooled-calibration",SourceUP183BSeal:"d4e715241c94f32b46e52f595bf6a84b7c675b1d",SingleCalibrationPhases:[]int{26},PooledCalibrationPhases:[]int{26,27,28,29,30},EvaluationPhases:phases,ModelsFrozenBeforeEvaluation:true,EvaluationFittingUsed:false,PerPhaseScalingUsed:false,MaintenanceTriggered:false}
	models:=[]up184bModel{single,pooled};sums:=map[string]float64{}
	for _,m:=range models{for _,phase:=range phases{x:=up184bEval(m,phase);res.Metrics=append(res.Metrics,x);sums[m.name]+=x.AbsoluteMassRatioError}}
	for _,m:=range models{res.Summaries=append(res.Summaries,UP184BModelSummary{Model:m.name,MeanAbsoluteMassRatioError:sums[m.name]/float64(len(phases))})}
	return res,nil
}
