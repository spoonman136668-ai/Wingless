package unitary

import("fmt";"math")

const UP182BCalibrationSchema="wingless.up182b-continuous-risk-calibration.v1"

type UP182BCalibration struct{
	BrierScore float64 `json:"brier_score"`
	ExpectedCalibrationError float64 `json:"expected_calibration_error"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	ActualCrossings int `json:"actual_crossings"`
	TotalSlots int `json:"total_slots"`
}
type UP182BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP181BSeal string `json:"source_up181b_seal"`
	CalibrationPhase int `json:"calibration_phase"`
	EvaluationPhase int `json:"evaluation_phase"`
	ProbabilityEstimator string `json:"probability_estimator"`
	FeatureBinsFrozenBeforeEvaluation bool `json:"feature_bins_frozen_before_evaluation"`
	EvaluationFittingUsed bool `json:"evaluation_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Calibration UP182BCalibration `json:"calibration"`
	Metrics []UP178BMetric `json:"metrics"`
}
type up182bRow struct{cls,temp string;mb int;margin float64;positive bool}
type up182bCell struct{slots,cross int}
func up182bMarginBin(x float64)int{
	if x<0.01{return 0};if x<0.02{return 1};if x<0.03{return 2};if x<0.04{return 3};if x<0.05{return 4};if x<0.10{return 5};return 6
}
func up182bRows(epoch int)[]up182bRow{
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
	rows:=make([]up182bRow,0,17280)
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r)
		for _,path:=range paths{
			g:=*base;streak:=make([]int,120)
			for _,idx:=range path{
				before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls)
				for i:=range before{if up166bRank(before,i)<=band{streak[i]++}else{streak[i]=0}}
				items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
				for i:=range before{
					temp:=up181bTemporal(streak[i]);if temp==""{temp="out_of_band"}
					am:=math.Abs(before[i].margin)
					rows=append(rows,up182bRow{cls:cls,temp:temp,mb:up182bMarginBin(am),margin:am,positive:before[i].correct!=after[i].correct})
				}
			}
		}
	}
	return rows
}
func up182bKey(r up182bRow)string{return fmt.Sprintf("%s|%d|%s",r.cls,r.mb,r.temp)}
func RunUP182B()(UP182BResult,error){
	train:=up182bRows(26);test:=up182bRows(27)
	cells:=map[string]*up182bCell{};trainPos:=0
	for _,r:=range train{key:=up182bKey(r);if cells[key]==nil{cells[key]=&up182bCell{}};cells[key].slots++;if r.positive{cells[key].cross++;trainPos++}}
	globalP:=float64(trainPos+1)/float64(len(train)+2)
	riskObs:=make([]up178bObs,0,len(test));marginObs:=make([]up178bObs,0,len(test))
	brier,pred:=0.0,0.0;actual:=0
	type eceBin struct{n,pos int;psum float64};ece:=make([]eceBin,10)
	for _,r:=range test{
		p:=globalP
		if c:=cells[up182bKey(r)];c!=nil{p=float64(c.cross+1)/float64(c.slots+2)}
		y:=0.0;if r.positive{y=1;actual++}
		diff:=p-y;brier+=diff*diff;pred+=p
		b:=int(p*10);if b>9{b=9};ece[b].n++;ece[b].psum+=p;if r.positive{ece[b].pos++}
		riskObs=append(riskObs,up178bObs{score:p,positive:r.positive})
		marginObs=append(marginObs,up178bObs{score:-r.margin,positive:r.positive})
	}
	eceTotal:=0.0
	for _,b:=range ece{if b.n>0{avgP:=b.psum/float64(b.n);freq:=float64(b.pos)/float64(b.n);eceTotal+=float64(b.n)/float64(len(test))*math.Abs(avgP-freq)}}
	cal:=UP182BCalibration{BrierScore:brier/float64(len(test)),ExpectedCalibrationError:eceTotal,PredictedCrossingsSum:pred,ActualCrossings:actual,TotalSlots:len(test)}
	return UP182BResult{Schema:UP182BCalibrationSchema,Experiment:"UP-182B-continuous-risk-calibration",SourceUP181BSeal:"9a6d2dccea33afd67bd93a10d3199b307f3203e7",CalibrationPhase:26,EvaluationPhase:27,ProbabilityEstimator:"class_x_margin01_x_temporal_laplace",FeatureBinsFrozenBeforeEvaluation:true,EvaluationFittingUsed:false,MaintenanceTriggered:false,Calibration:cal,Metrics:[]UP178BMetric{up178bDisc("calibrated_margin_temporal",riskObs),up178bDisc("margin_only",marginObs)}},nil
}
