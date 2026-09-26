package unitary

import "math"

const UP183BDriftSchema="wingless.up183b-calibration-drift-sweep.v1"

type UP183BPoint struct{
	EvaluationPhase int `json:"evaluation_phase"`
	ActualCrossings int `json:"actual_crossings"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	PredictedToActualRatio float64 `json:"predicted_to_actual_ratio"`
	BrierScore float64 `json:"brier_score"`
	ExpectedCalibrationError float64 `json:"expected_calibration_error"`
	AUROC float64 `json:"auroc"`
	AveragePrecision float64 `json:"average_precision"`
}
type UP183BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP182BSeal string `json:"source_up182b_seal"`
	CalibrationPhase int `json:"calibration_phase"`
	EvaluationPhases []int `json:"evaluation_phases"`
	ModelFrozenBeforeEvaluation bool `json:"model_frozen_before_evaluation"`
	EvaluationFittingUsed bool `json:"evaluation_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	MeanPredictedToActualRatio float64 `json:"mean_predicted_to_actual_ratio"`
	MinPredictedToActualRatio float64 `json:"min_predicted_to_actual_ratio"`
	MaxPredictedToActualRatio float64 `json:"max_predicted_to_actual_ratio"`
	Points []UP183BPoint `json:"points"`
}
type up183bECEBin struct{n,pos int;psum float64}
func up183bEval(test []up182bRow,cells map[string]*up182bCell,globalP float64,phase int)UP183BPoint{
	riskObs:=make([]up178bObs,0,len(test));brier,pred:=0.0,0.0;actual:=0;ece:=make([]up183bECEBin,10)
	for _,r:=range test{
		p:=globalP;if c:=cells[up182bKey(r)];c!=nil{p=float64(c.cross+1)/float64(c.slots+2)}
		y:=0.0;if r.positive{y=1;actual++};d:=p-y;brier+=d*d;pred+=p
		b:=int(p*10);if b>9{b=9};ece[b].n++;ece[b].psum+=p;if r.positive{ece[b].pos++}
		riskObs=append(riskObs,up178bObs{score:p,positive:r.positive})
	}
	eceTotal:=0.0;for _,b:=range ece{if b.n>0{avg:=b.psum/float64(b.n);freq:=float64(b.pos)/float64(b.n);eceTotal+=float64(b.n)/float64(len(test))*math.Abs(avg-freq)}}
	m:=up178bDisc("frozen_risk",riskObs);ratio:=0.0;if actual>0{ratio=pred/float64(actual)}
	return UP183BPoint{EvaluationPhase:phase,ActualCrossings:actual,PredictedCrossingsSum:pred,PredictedToActualRatio:ratio,BrierScore:brier/float64(len(test)),ExpectedCalibrationError:eceTotal,AUROC:m.AUROC,AveragePrecision:m.AveragePrecision}
}
func RunUP183B()(UP183BResult,error){
	train:=up182bRows(26);cells:=map[string]*up182bCell{};trainPos:=0
	for _,r:=range train{key:=up182bKey(r);if cells[key]==nil{cells[key]=&up182bCell{}};cells[key].slots++;if r.positive{cells[key].cross++;trainPos++}}
	globalP:=float64(trainPos+1)/float64(len(train)+2);phases:=[]int{27,28,29,30}
	res:=UP183BResult{Schema:UP183BDriftSchema,Experiment:"UP-183B-calibration-drift-sweep",SourceUP182BSeal:"5d5978457f6b1a25fc493f66fffed448f3941d6c",CalibrationPhase:26,EvaluationPhases:phases,ModelFrozenBeforeEvaluation:true,EvaluationFittingUsed:false,MaintenanceTriggered:false}
	sum:=0.0;min,max:=0.0,0.0
	for i,p:=range phases{pt:=up183bEval(up182bRows(p),cells,globalP,p);res.Points=append(res.Points,pt);sum+=pt.PredictedToActualRatio;if i==0||pt.PredictedToActualRatio<min{min=pt.PredictedToActualRatio};if i==0||pt.PredictedToActualRatio>max{max=pt.PredictedToActualRatio}}
	res.MeanPredictedToActualRatio=sum/float64(len(phases));res.MinPredictedToActualRatio=min;res.MaxPredictedToActualRatio=max
	return res,nil
}
