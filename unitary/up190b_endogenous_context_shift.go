package unitary

import "math"

const UP190BContextShiftSchema="wingless.up190b-endogenous-context-shift.v1"

type UP190BProfile struct{
	Name string `json:"name"`
	SurfaceIndices []int `json:"surface_indices"`
}
type UP190BMetric struct{
	Model string `json:"model"`
	Profile string `json:"profile"`
	Phase int `json:"phase"`
	Context string `json:"context"`
	NativeCorrectCount int `json:"native_correct_count"`
	OutsideKnownNativeCounts bool `json:"outside_known_native_counts"`
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
type UP190BSummary struct{
	Model string `json:"model"`
	MeanAbsoluteMassRatioError float64 `json:"mean_absolute_mass_ratio_error"`
	MaxAbsoluteMassRatioError float64 `json:"max_absolute_mass_ratio_error"`
}
type UP190BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP189BSeal string `json:"source_up189b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	ContextThreshold int `json:"context_threshold"`
	KnownNativeCorrectCounts []int `json:"known_native_correct_counts"`
	ContextAFactor float64 `json:"context_a_factor"`
	ContextBFactor float64 `json:"context_b_factor"`
	Profiles []UP190BProfile `json:"profiles"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveShiftSearchUsed bool `json:"adaptive_shift_search_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP190BMetric `json:"metrics"`
	Summaries []UP190BSummary `json:"summaries"`
}

func up190bShiftBase(common *up129bGate,subject string,phase int,indices []int,o,r [64]float64)*up129bGate{
	g:=up169bPostReport(common,subject,phase,o,r)
	items:=up158bItems(phase,indices)
	up156bApplyOld(g,items,0,len(items),o,r)
	return g
}
func up190bNativeCorrectCount(phase int,indices []int)int{
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);n:=0
	for _,subject:=range up130bNewNames{
		g:=up190bShiftBase(common,subject,phase,indices,o,r)
		for _,e:=range up165bMargins(g,o,r){if e.correct{n++}}
	}
	return n
}
func up190bRows(phase int,indices []int)[]up182bRow{
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
	rows:=make([]up182bRow,0,17280)
	for _,subject:=range up130bNewNames{
		base:=up190bShiftBase(common,subject,phase,indices,o,r)
		for _,path:=range paths{
			g:=*base;streak:=make([]int,120)
			for _,idx:=range path{
				before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls)
				for i:=range before{if up166bRank(before,i)<=band{streak[i]++}else{streak[i]=0}}
				items:=up158bItems(phase,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
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
func up190bEval(m up184bModel,factors map[string]float64,profile string,phase int,indices []int,correct bool)UP190BMetric{
	nc:=up190bNativeCorrectCount(phase,indices);ctx:="context_B";if nc<=621{ctx="context_A"}
	test:=up190bRows(phase,indices);obs:=make([]up178bObs,0,len(test));brier,pred:=0.0,0.0;actual:=0
	type eceBin struct{n,pos int;psum float64};ece:=make([]eceBin,10)
	for _,row:=range test{
		p:=up188bP(m,row)
		if correct{p*=factors[ctx];if p>1{p=1};if p<0{p=0}}
		y:=0.0;if row.positive{y=1;actual++};d:=p-y;brier+=d*d;pred+=p
		b:=int(p*10);if b>9{b=9};ece[b].n++;ece[b].psum+=p;if row.positive{ece[b].pos++}
		obs=append(obs,up178bObs{score:p,positive:row.positive})
	}
	eceV:=0.0;for _,b:=range ece{if b.n>0{ap:=b.psum/float64(b.n);fr:=float64(b.pos)/float64(b.n);eceV+=float64(b.n)/float64(len(test))*math.Abs(ap-fr)}}
	disc:=up178bDisc("up190b",obs);ratio:=0.0;if actual>0{ratio=pred/float64(actual)}
	name:="pooled_26_30";if correct{name="pooled_plus_native_mass_correction"}
	outside:=nc!=620&&nc!=623
	return UP190BMetric{Model:name,Profile:profile,Phase:phase,Context:ctx,NativeCorrectCount:nc,OutsideKnownNativeCounts:outside,ActualCrossings:actual,PredictedCrossingsSum:pred,PredictedActualRatio:ratio,AbsoluteMassRatioError:math.Abs(ratio-1),BrierScore:brier/float64(len(test)),ExpectedCalibrationError:eceV,AUROC:disc.AUROC,AveragePrecision:disc.AveragePrecision,TotalSlots:len(test)}
}
func RunUP190B()(UP190BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{43,44,45}
	profiles:=[]UP190BProfile{{Name:"store_shift",SurfaceIndices:[]int{0,1,2,3}},{Name:"observe_shift",SurfaceIndices:[]int{5,6,7,8}},{Name:"mixed_shift",SurfaceIndices:[]int{0,5,1,6}}}
	pooled:=up184bBuild("pooled_26_30",cal)
	f:=map[string]float64{"context_A":1.0103904235372925,"context_B":0.8983526543771945}
	res:=UP190BResult{Schema:UP190BContextShiftSchema,Experiment:"UP-190B-endogenous-context-shift",SourceUP189BSeal:"62429a84c67c974773efe5ddb51d7324be880507",CalibrationPhases:cal,EvaluationPhases:eval,ContextThreshold:621,KnownNativeCorrectCounts:[]int{620,623},ContextAFactor:f["context_A"],ContextBFactor:f["context_B"],Profiles:profiles,PhaseInputUsed:false,ParityInputUsed:false,HeldoutFittingUsed:false,AdaptiveShiftSearchUsed:false,MaintenanceTriggered:false}
	sum:=map[string]float64{};maxv:=map[string]float64{};count:=map[string]int{}
	for _,p:=range profiles{for _,phase:=range eval{
		a:=up190bEval(pooled,f,p.Name,phase,p.SurfaceIndices,false);b:=up190bEval(pooled,f,p.Name,phase,p.SurfaceIndices,true)
		res.Metrics=append(res.Metrics,a,b)
		for _,m:=range []UP190BMetric{a,b}{sum[m.Model]+=m.AbsoluteMassRatioError;count[m.Model]++;if m.AbsoluteMassRatioError>maxv[m.Model]{maxv[m.Model]=m.AbsoluteMassRatioError}}
	}}
	for _,name:=range []string{"pooled_26_30","pooled_plus_native_mass_correction"}{res.Summaries=append(res.Summaries,UP190BSummary{Model:name,MeanAbsoluteMassRatioError:sum[name]/float64(count[name]),MaxAbsoluteMassRatioError:maxv[name]})}
	return res,nil
}
