package unitary

import "math"

const UP183BDriftSchema="wingless.up183b-calibration-drift.v1"

type UP183BPhaseMetric struct{
	Phase int `json:"phase"`
	ActualCrossings int `json:"actual_crossings"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	PredictedActualRatio float64 `json:"predicted_actual_ratio"`
	BrierScore float64 `json:"brier_score"`
	ExpectedCalibrationError float64 `json:"expected_calibration_error"`
	AUROC float64 `json:"auroc"`
	AveragePrecision float64 `json:"average_precision"`
	TotalSlots int `json:"total_slots"`
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
	Metrics []UP183BPhaseMetric `json:"metrics"`
}
func RunUP183B()(UP183BResult,error){
	train:=up182bRows(26);cells:=map[string]*up182bCell{};trainPos:=0
	for _,r:=range train{key:=up182bKey(r);if cells[key]==nil{cells[key]=&up182bCell{}};cells[key].slots++;if r.positive{cells[key].cross++;trainPos++}}
	globalP:=float64(trainPos+1)/float64(len(train)+2)
	phases:=[]int{28,29,30}
	res:=UP183BResult{Schema:UP183BDriftSchema,Experiment:"UP-183B-calibration-drift",SourceUP182BSeal:"5d5978457f6b1a25fc493f66fffed448f3941d6c",CalibrationPhase:26,EvaluationPhases:phases,ModelFrozenBeforeEvaluation:true,EvaluationFittingUsed:false,MaintenanceTriggered:false}
	for _,phase:=range phases{
		test:=up182bRows(phase);obs:=make([]up178bObs,0,len(test));brier,pred:=0.0,0.0;actual:=0
		type eceBin struct{n,pos int;psum float64};ece:=make([]eceBin,10)
		for _,r:=range test{
			p:=globalP;if c:=cells[up182bKey(r)];c!=nil{p=float64(c.cross+1)/float64(c.slots+2)}
			y:=0.0;if r.positive{y=1;actual++};diff:=p-y;brier+=diff*diff;pred+=p
			b:=int(p*10);if b>9{b=9};ece[b].n++;ece[b].psum+=p;if r.positive{ece[b].pos++}
			obs=append(obs,up178bObs{score:p,positive:r.positive})
		}
		eceV:=0.0;for _,b:=range ece{if b.n>0{ap:=b.psum/float64(b.n);fr:=float64(b.pos)/float64(b.n);eceV+=float64(b.n)/float64(len(test))*math.Abs(ap-fr)}}
		disc:=up178bDisc("frozen_calibrated",obs);ratio:=0.0;if actual>0{ratio=pred/float64(actual)}
		res.Metrics=append(res.Metrics,UP183BPhaseMetric{Phase:phase,ActualCrossings:actual,PredictedCrossingsSum:pred,PredictedActualRatio:ratio,BrierScore:brier/float64(len(test)),ExpectedCalibrationError:eceV,AUROC:disc.AUROC,AveragePrecision:disc.AveragePrecision,TotalSlots:len(test)})
	}
	return res,nil
}
