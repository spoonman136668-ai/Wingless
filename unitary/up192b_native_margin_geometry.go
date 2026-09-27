package unitary

import "math"

const UP192BNativeMarginGeometrySchema="wingless.up192b-native-margin-geometry.v1"

type UP192BPoint struct{
	Profile string `json:"profile"`
	Phase int `json:"phase"`
	NativeCorrectCount int `json:"native_correct_count"`
	MeanAbsoluteMargin float64 `json:"mean_absolute_margin"`
	NearZeroMarginCount int `json:"near_zero_margin_count"`
	MinAbsoluteMargin float64 `json:"min_absolute_margin"`
	RequiredMassFactor float64 `json:"required_mass_factor"`
}
type UP192BCorrelations struct{
	NativeCorrectCount float64 `json:"native_correct_count"`
	MeanAbsoluteMargin float64 `json:"mean_absolute_margin"`
	NearZeroMarginCount float64 `json:"near_zero_margin_count"`
	MinAbsoluteMargin float64 `json:"min_absolute_margin"`
}
type UP192BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP191BSeal string `json:"source_up191b_seal"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Profiles []UP191BProfile `json:"profiles"`
	NearZeroThreshold float64 `json:"near_zero_threshold"`
	FittedCorrectionUsed bool `json:"fitted_correction_used"`
	HeldoutTuningUsed bool `json:"heldout_tuning_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	AdaptiveProfileSearchUsed bool `json:"adaptive_profile_search_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Correlations UP192BCorrelations `json:"correlations"`
	Points []UP192BPoint `json:"points"`
}
func up192bGeometry(phase int,indices []int)(correct int,meanAbs float64,nearZero int,minAbs float64){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);minAbs=math.Inf(1);n:=0
	for _,subject:=range up130bNewNames{
		var g *up129bGate
		if len(indices)==0{g=up169bPostReport(common,subject,phase,o,r)}else{g=up190bShiftBase(common,subject,phase,indices,o,r)}
		for _,e:=range up165bMargins(g,o,r){
			if e.correct{correct++};a:=math.Abs(e.margin);meanAbs+=a;n++;if a<0.01{nearZero++};if a<minAbs{minAbs=a}
		}
	}
	if n>0{meanAbs/=float64(n)};if math.IsInf(minAbs,1){minAbs=0};return
}
func up192bCorr(points []UP192BPoint,which int)float64{
	if len(points)==0{return 0};mx,my:=0.0,0.0
	xf:=func(p UP192BPoint)float64{switch which{case 0:return float64(p.NativeCorrectCount);case 1:return p.MeanAbsoluteMargin;case 2:return float64(p.NearZeroMarginCount);default:return p.MinAbsoluteMargin}}
	for _,p:=range points{mx+=xf(p);my+=p.RequiredMassFactor};mx/=float64(len(points));my/=float64(len(points))
	num,dx,dy:=0.0,0.0,0.0
	for _,p:=range points{x:=xf(p)-mx;y:=p.RequiredMassFactor-my;num+=x*y;dx+=x*x;dy+=y*y}
	if dx==0||dy==0{return 0};return num/math.Sqrt(dx*dy)
}
func RunUP192B()(UP192BResult,error){
	eval:=[]int{43,44,45}
	profiles:=[]UP191BProfile{
		{Name:"control"},{Name:"store1",SurfaceIndices:[]int{0}},{Name:"store2",SurfaceIndices:[]int{0,1}},{Name:"store4",SurfaceIndices:[]int{0,1,2,3}},
		{Name:"observe1",SurfaceIndices:[]int{5}},{Name:"observe2",SurfaceIndices:[]int{5,6}},{Name:"observe4",SurfaceIndices:[]int{5,6,7,8}},
		{Name:"mixed2",SurfaceIndices:[]int{0,5}},{Name:"mixed4",SurfaceIndices:[]int{0,5,1,6}},
	}
	pooled:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP192BResult{Schema:UP192BNativeMarginGeometrySchema,Experiment:"UP-192B-native-margin-geometry",SourceUP191BSeal:"ae9469a8672adc110f6ec1c9382432fef9a9230b",EvaluationPhases:eval,Profiles:profiles,NearZeroThreshold:0.01,FittedCorrectionUsed:false,HeldoutTuningUsed:false,AdaptiveFeatureSelectionUsed:false,AdaptiveProfileSearchUsed:false,MaintenanceTriggered:false}
	for _,p:=range profiles{for _,phase:=range eval{
		base:=up191bPoint(pooled,p,phase);correct,meanAbs,nearZero,minAbs:=up192bGeometry(phase,p.SurfaceIndices)
		res.Points=append(res.Points,UP192BPoint{Profile:p.Name,Phase:phase,NativeCorrectCount:correct,MeanAbsoluteMargin:meanAbs,NearZeroMarginCount:nearZero,MinAbsoluteMargin:minAbs,RequiredMassFactor:base.RequiredMassFactor})
	}}
	res.Correlations=UP192BCorrelations{NativeCorrectCount:up192bCorr(res.Points,0),MeanAbsoluteMargin:up192bCorr(res.Points,1),NearZeroMarginCount:up192bCorr(res.Points,2),MinAbsoluteMargin:up192bCorr(res.Points,3)}
	return res,nil
}
