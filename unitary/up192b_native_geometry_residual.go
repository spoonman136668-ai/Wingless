package unitary

import "math"

const UP192BNativeGeometrySchema="wingless.up192b-native-geometry-residual.v1"

type UP192BProfile struct{
	Name string `json:"name"`
	SurfaceIndices []int `json:"surface_indices"`
}
type UP192BMetric struct{
	Profile string `json:"profile"`
	Phase int `json:"phase"`
	NativeCorrectCount int `json:"native_correct_count"`
	MeanSignedMargin float64 `json:"mean_signed_margin"`
	MeanAbsoluteMargin float64 `json:"mean_absolute_margin"`
	MinimumMargin float64 `json:"minimum_margin"`
	NearZeroMarginCount int `json:"near_zero_margin_count"`
	NegativeMarginMass float64 `json:"negative_margin_mass"`
	RequiredCalibrationFactor float64 `json:"required_calibration_factor"`
}
type UP192BCorrelation struct{
	Feature string `json:"feature"`
	PearsonRequiredFactor float64 `json:"pearson_required_factor"`
}
type UP192BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP191BSeal string `json:"source_up191b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Profiles []UP192BProfile `json:"profiles"`
	NearZeroThreshold float64 `json:"near_zero_threshold"`
	NewCorrectionFit bool `json:"new_correction_fit"`
	CorrectionApplied bool `json:"correction_applied"`
	AdaptiveFeatureSearchUsed bool `json:"adaptive_feature_search_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP192BMetric `json:"metrics"`
	Correlations []UP192BCorrelation `json:"correlations"`
}

func up192bGeometry(phase int,indices []int)(correct int,mean,meanAbs,minMargin float64,nearZero int,negMass float64){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	minMargin=math.Inf(1);n:=0
	for _,subject:=range up130bNewNames{
		g:=up190bShiftBase(common,subject,phase,indices,o,r)
		for _,e:=range up165bMargins(g,o,r){
			n++;if e.correct{correct++}
			mean+=e.margin;am:=math.Abs(e.margin);meanAbs+=am
			if e.margin<minMargin{minMargin=e.margin}
			if am<0.01{nearZero++}
			if e.margin<0{negMass+=-e.margin}
		}
	}
	if n>0{mean/=float64(n);meanAbs/=float64(n)}
	return
}
func up192bPearson(metrics []UP192BMetric,feature func(UP192BMetric)float64)float64{
	if len(metrics)<2{return 0}
	mx,my:=0.0,0.0
	for _,m:=range metrics{mx+=feature(m);my+=m.RequiredCalibrationFactor}
	n:=float64(len(metrics));mx/=n;my/=n
	num,dx,dy:=0.0,0.0,0.0
	for _,m:=range metrics{x:=feature(m)-mx;y:=m.RequiredCalibrationFactor-my;num+=x*y;dx+=x*x;dy+=y*y}
	if dx==0||dy==0{return 0};return num/math.Sqrt(dx*dy)
}
func RunUP192B()(UP192BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{43,44,45}
	profiles:=[]UP192BProfile{
		{Name:"control",SurfaceIndices:[]int{}},
		{Name:"store1",SurfaceIndices:[]int{0}},
		{Name:"store2",SurfaceIndices:[]int{0,1}},
		{Name:"store_forward",SurfaceIndices:[]int{0,1,2,3}},
		{Name:"store_reverse",SurfaceIndices:[]int{3,2,1,0}},
		{Name:"observe1",SurfaceIndices:[]int{5}},
		{Name:"observe2",SurfaceIndices:[]int{5,6}},
		{Name:"observe_forward",SurfaceIndices:[]int{5,6,7,8}},
		{Name:"observe_reverse",SurfaceIndices:[]int{8,7,6,5}},
		{Name:"mixed2",SurfaceIndices:[]int{0,5}},
		{Name:"mixed_forward",SurfaceIndices:[]int{0,5,1,6}},
		{Name:"mixed_reverse",SurfaceIndices:[]int{6,1,5,0}},
	}
	model:=up184bBuild("pooled_26_30",cal)
	res:=UP192BResult{Schema:UP192BNativeGeometrySchema,Experiment:"UP-192B-native-geometry-residual",SourceUP191BSeal:"e6cf06e3bc44ca603df9cf9dc8f54b87dc2d36ba",CalibrationPhases:cal,EvaluationPhases:eval,Profiles:profiles,NearZeroThreshold:0.01,NewCorrectionFit:false,CorrectionApplied:false,AdaptiveFeatureSearchUsed:false,MaintenanceTriggered:false}
	for _,p:=range profiles{for _,phase:=range eval{
		base:=up191bMetric(model,p.Name,phase,p.SurfaceIndices)
		c,mean,meanAbs,minMargin,nz,neg:=up192bGeometry(phase,p.SurfaceIndices)
		if c!=base.NativeCorrectCount{panic("UP192B_NATIVE_COUNT_MISMATCH")}
		res.Metrics=append(res.Metrics,UP192BMetric{Profile:p.Name,Phase:phase,NativeCorrectCount:c,MeanSignedMargin:mean,MeanAbsoluteMargin:meanAbs,MinimumMargin:minMargin,NearZeroMarginCount:nz,NegativeMarginMass:neg,RequiredCalibrationFactor:base.RequiredCalibrationFactor})
	}}
	features:=[]struct{name string;f func(UP192BMetric)float64}{
		{"native_correct_count",func(m UP192BMetric)float64{return float64(m.NativeCorrectCount)}},
		{"mean_signed_margin",func(m UP192BMetric)float64{return m.MeanSignedMargin}},
		{"mean_absolute_margin",func(m UP192BMetric)float64{return m.MeanAbsoluteMargin}},
		{"minimum_margin",func(m UP192BMetric)float64{return m.MinimumMargin}},
		{"near_zero_margin_count",func(m UP192BMetric)float64{return float64(m.NearZeroMarginCount)}},
		{"negative_margin_mass",func(m UP192BMetric)float64{return m.NegativeMarginMass}},
	}
	for _,x:=range features{res.Correlations=append(res.Correlations,UP192BCorrelation{Feature:x.name,PearsonRequiredFactor:up192bPearson(res.Metrics,x.f)})}
	return res,nil
}
