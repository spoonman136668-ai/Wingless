package unitary

import "math"

const UP188BMassCorrectionSchema="wingless.up188b-native-mass-correction.v1"

type UP188BFactor struct{
	Context string `json:"context"`
	ActualCrossings int `json:"actual_crossings"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	Factor float64 `json:"factor"`
}
type UP188BMetric struct{
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
type UP188BSummary struct{
	Model string `json:"model"`
	MeanAbsoluteMassRatioError float64 `json:"mean_absolute_mass_ratio_error"`
}
type UP188BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP187BSeal string `json:"source_up187b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	ContextThreshold int `json:"context_threshold"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Factors []UP188BFactor `json:"factors"`
	Metrics []UP188BMetric `json:"metrics"`
	Summaries []UP188BSummary `json:"summaries"`
}
func up188bP(m up184bModel,r up182bRow)float64{
	p:=m.globalP
	if c:=m.cells[up182bKey(r)];c!=nil{p=float64(c.cross+1)/float64(c.slots+2)}
	return p
}
func up188bFactors(m up184bModel,phases []int)map[string]float64{
	type acc struct{actual int;pred float64}
	a:=map[string]*acc{}
	for _,phase:=range phases{
		ctx,_:=up187bContext(phase)
		if a[ctx]==nil{a[ctx]=&acc{}}
		for _,r:=range up182bRows(phase){
			a[ctx].pred+=up188bP(m,r)
			if r.positive{a[ctx].actual++}
		}
	}
	out:=map[string]float64{}
	for ctx,x:=range a{out[ctx]=float64(x.actual)/x.pred}
	return out
}
func up188bEval(m up184bModel,factors map[string]float64,phase int,correct bool)UP188BMetric{
	ctx,nc:=up187bContext(phase);test:=up182bRows(phase);obs:=make([]up178bObs,0,len(test))
	brier,pred:=0.0,0.0;actual:=0
	type eceBin struct{n,pos int;psum float64};ece:=make([]eceBin,10)
	for _,r:=range test{
		p:=up188bP(m,r)
		if correct{p*=factors[ctx];if p>1{p=1};if p<0{p=0}}
		y:=0.0;if r.positive{y=1;actual++};d:=p-y;brier+=d*d;pred+=p
		b:=int(p*10);if b>9{b=9};ece[b].n++;ece[b].psum+=p;if r.positive{ece[b].pos++}
		obs=append(obs,up178bObs{score:p,positive:r.positive})
	}
	eceV:=0.0;for _,b:=range ece{if b.n>0{ap:=b.psum/float64(b.n);fr:=float64(b.pos)/float64(b.n);eceV+=float64(b.n)/float64(len(test))*math.Abs(ap-fr)}}
	disc:=up178bDisc("x",obs);ratio:=0.0;if actual>0{ratio=pred/float64(actual)}
	name:="pooled_26_30";if correct{name="pooled_plus_native_mass_correction"}
	return UP188BMetric{Model:name,Phase:phase,Context:ctx,NativeCorrectCount:nc,ActualCrossings:actual,PredictedCrossingsSum:pred,PredictedActualRatio:ratio,AbsoluteMassRatioError:math.Abs(ratio-1),BrierScore:brier/float64(len(test)),ExpectedCalibrationError:eceV,AUROC:disc.AUROC,AveragePrecision:disc.AveragePrecision,TotalSlots:len(test)}
}
func RunUP188B()(UP188BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{34,35,36};pooled:=up184bBuild("pooled_26_30",cal);f:=up188bFactors(pooled,cal)
	res:=UP188BResult{Schema:UP188BMassCorrectionSchema,Experiment:"UP-188B-native-mass-correction",SourceUP187BSeal:"5647d1c20194b0e5d133055c2e25f90c7bd73518",CalibrationPhases:cal,EvaluationPhases:eval,ContextThreshold:621,PhaseInputUsed:false,ParityInputUsed:false,HeldoutFittingUsed:false,MaintenanceTriggered:false}
	for _,ctx:=range []string{"context_A","context_B"}{
		actual,pred:=0,0.0
		for _,phase:=range cal{c,_:=up187bContext(phase);if c!=ctx{continue};for _,r:=range up182bRows(phase){pred+=up188bP(pooled,r);if r.positive{actual++}}}
		res.Factors=append(res.Factors,UP188BFactor{Context:ctx,ActualCrossings:actual,PredictedCrossingsSum:pred,Factor:f[ctx]})
	}
	sums:=map[string]float64{}
	for _,phase:=range eval{
		a:=up188bEval(pooled,f,phase,false);b:=up188bEval(pooled,f,phase,true)
		res.Metrics=append(res.Metrics,a,b);sums[a.Model]+=a.AbsoluteMassRatioError;sums[b.Model]+=b.AbsoluteMassRatioError
	}
	for _,name:=range []string{"pooled_26_30","pooled_plus_native_mass_correction"}{res.Summaries=append(res.Summaries,UP188BSummary{Model:name,MeanAbsoluteMassRatioError:sums[name]/3})}
	return res,nil
}
