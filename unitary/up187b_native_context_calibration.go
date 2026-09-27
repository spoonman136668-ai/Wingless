package unitary

import "math"

const UP187BNativeCalibrationSchema="wingless.up187b-native-context-calibration.v1"

type up187bModel struct{
	name string
	cells map[string]*up182bCell
	globalP float64
}
type UP187BMetric struct{
	Model string `json:"model"`
	Phase int `json:"phase"`
	Context string `json:"context"`
	NativeCorrectCount int `json:"native_correct_count"`
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
type UP187BSummary struct{
	Model string `json:"model"`
	MeanAbsoluteMassRatioError float64 `json:"mean_absolute_mass_ratio_error"`
}
type UP187BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP186BSeal string `json:"source_up186b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	ContextThreshold int `json:"context_threshold"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	EvaluationFittingUsed bool `json:"evaluation_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP187BMetric `json:"metrics"`
	Summaries []UP187BSummary `json:"summaries"`
}
func up187bCorrectCount(phase int)int{
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);n:=0
	for _,subject:=range up130bNewNames{
		g:=up169bPostReport(common,subject,phase,o,r)
		for _,e:=range up165bMargins(g,o,r){if e.correct{n++}}
	}
	return n
}
func up187bContext(phase int)(string,int){
	n:=up187bCorrectCount(phase)
	if n<=621{return "context_A",n}
	return "context_B",n
}
func up187bBuild(name string,phases []int,condition bool)up187bModel{
	cells:=map[string]*up182bCell{};pos,total:=0,0
	for _,phase:=range phases{
		ctx,_:=up187bContext(phase)
		for _,r:=range up182bRows(phase){
			key:=up182bKey(r);if condition{key=ctx+"|"+key}
			if cells[key]==nil{cells[key]=&up182bCell{}}
			cells[key].slots++;total++;if r.positive{cells[key].cross++;pos++}
		}
	}
	return up187bModel{name:name,cells:cells,globalP:float64(pos+1)/float64(total+2)}
}
func up187bEval(m up187bModel,phase int,condition bool)UP187BMetric{
	ctx,correct:=up187bContext(phase);test:=up182bRows(phase);obs:=make([]up178bObs,0,len(test))
	brier,pred:=0.0,0.0;actual:=0
	type eceBin struct{n,pos int;psum float64};ece:=make([]eceBin,10)
	for _,r:=range test{
		key:=up182bKey(r);if condition{key=ctx+"|"+key}
		p:=m.globalP;if c:=m.cells[key];c!=nil{p=float64(c.cross+1)/float64(c.slots+2)}
		y:=0.0;if r.positive{y=1;actual++};d:=p-y;brier+=d*d;pred+=p
		b:=int(p*10);if b>9{b=9};ece[b].n++;ece[b].psum+=p;if r.positive{ece[b].pos++}
		obs=append(obs,up178bObs{score:p,positive:r.positive})
	}
	eceV:=0.0;for _,b:=range ece{if b.n>0{ap:=b.psum/float64(b.n);fr:=float64(b.pos)/float64(b.n);eceV+=float64(b.n)/float64(len(test))*math.Abs(ap-fr)}}
	disc:=up178bDisc(m.name,obs);ratio:=0.0;if actual>0{ratio=pred/float64(actual)}
	return UP187BMetric{Model:m.name,Phase:phase,Context:ctx,NativeCorrectCount:correct,ActualCrossings:actual,PredictedCrossingsSum:pred,PredictedActualRatio:ratio,AbsoluteMassRatioError:math.Abs(ratio-1),BrierScore:brier/float64(len(test)),ExpectedCalibrationError:eceV,AUROC:disc.AUROC,AveragePrecision:disc.AveragePrecision,TotalSlots:len(test)}
}
func RunUP187B()(UP187BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{31,32,33}
	pooled:=up187bBuild("pooled_26_30",cal,false)
	native:=up187bBuild("native_context_26_30",cal,true)
	res:=UP187BResult{Schema:UP187BNativeCalibrationSchema,Experiment:"UP-187B-native-context-calibration",SourceUP186BSeal:"b95c60d4a41f91ab089b0983fda82157a797885d",CalibrationPhases:cal,EvaluationPhases:eval,ContextThreshold:621,PhaseInputUsed:false,ParityInputUsed:false,EvaluationFittingUsed:false,MaintenanceTriggered:false}
	sums:=map[string]float64{}
	for _,phase:=range eval{
		a:=up187bEval(pooled,phase,false);b:=up187bEval(native,phase,true)
		res.Metrics=append(res.Metrics,a,b);sums[a.Model]+=a.AbsoluteMassRatioError;sums[b.Model]+=b.AbsoluteMassRatioError
	}
	for _,name:=range []string{"pooled_26_30","native_context_26_30"}{res.Summaries=append(res.Summaries,UP187BSummary{Model:name,MeanAbsoluteMassRatioError:sums[name]/3.0})}
	return res,nil
}
