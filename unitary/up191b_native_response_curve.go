package unitary

import "math"

const UP191BNativeResponseSchema="wingless.up191b-native-response-curve.v1"

type UP191BProfile struct{
	Name string `json:"name"`
	SurfaceIndices []int `json:"surface_indices"`
}
type UP191BPoint struct{
	Profile string `json:"profile"`
	ShiftLength int `json:"shift_length"`
	Phase int `json:"phase"`
	NativeCorrectCount int `json:"native_correct_count"`
	ActualCrossings int `json:"actual_crossings"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	RequiredMassFactor float64 `json:"required_mass_factor"`
	PooledAbsoluteMassRatioError float64 `json:"pooled_absolute_mass_ratio_error"`
}
type UP191BCountGroup struct{
	NativeCorrectCount int `json:"native_correct_count"`
	Observations int `json:"observations"`
	MeanRequiredMassFactor float64 `json:"mean_required_mass_factor"`
	MinRequiredMassFactor float64 `json:"min_required_mass_factor"`
	MaxRequiredMassFactor float64 `json:"max_required_mass_factor"`
	FactorSpread float64 `json:"factor_spread"`
}
type UP191BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP190BSeal string `json:"source_up190b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Profiles []UP191BProfile `json:"profiles"`
	FittedCorrectionUsed bool `json:"fitted_correction_used"`
	HeldoutTuningUsed bool `json:"heldout_tuning_used"`
	AdaptiveProfileSearchUsed bool `json:"adaptive_profile_search_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	PearsonNativeCountRequiredFactor float64 `json:"pearson_native_count_required_factor"`
	MaxWithinCountFactorSpread float64 `json:"max_within_count_factor_spread"`
	Points []UP191BPoint `json:"points"`
	CountGroups []UP191BCountGroup `json:"count_groups"`
}
func up191bRowsAndCount(phase int,indices []int)([]up182bRow,int){
	if len(indices)==0{return up182bRows(phase),up187bCorrectCount(phase)}
	return up190bRows(phase,indices),up190bNativeCorrectCount(phase,indices)
}
func up191bPoint(m up184bModel,p UP191BProfile,phase int)UP191BPoint{
	rows,nc:=up191bRowsAndCount(phase,p.SurfaceIndices);actual:=0;pred:=0.0
	for _,row:=range rows{pred+=up188bP(m,row);if row.positive{actual++}}
	factor:=0.0;ratio:=0.0
	if pred>0{factor=float64(actual)/pred}
	if actual>0{ratio=pred/float64(actual)}
	return UP191BPoint{Profile:p.Name,ShiftLength:len(p.SurfaceIndices),Phase:phase,NativeCorrectCount:nc,ActualCrossings:actual,PredictedCrossingsSum:pred,RequiredMassFactor:factor,PooledAbsoluteMassRatioError:math.Abs(ratio-1)}
}
func up191bPearson(points []UP191BPoint)float64{
	if len(points)==0{return 0};mx,my:=0.0,0.0
	for _,p:=range points{mx+=float64(p.NativeCorrectCount);my+=p.RequiredMassFactor}
	mx/=float64(len(points));my/=float64(len(points));num,dx,dy:=0.0,0.0,0.0
	for _,p:=range points{x:=float64(p.NativeCorrectCount)-mx;y:=p.RequiredMassFactor-my;num+=x*y;dx+=x*x;dy+=y*y}
	if dx==0||dy==0{return 0};return num/math.Sqrt(dx*dy)
}
func RunUP191B()(UP191BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{43,44,45}
	profiles:=[]UP191BProfile{
		{Name:"control"},
		{Name:"store1",SurfaceIndices:[]int{0}},
		{Name:"store2",SurfaceIndices:[]int{0,1}},
		{Name:"store4",SurfaceIndices:[]int{0,1,2,3}},
		{Name:"observe1",SurfaceIndices:[]int{5}},
		{Name:"observe2",SurfaceIndices:[]int{5,6}},
		{Name:"observe4",SurfaceIndices:[]int{5,6,7,8}},
		{Name:"mixed2",SurfaceIndices:[]int{0,5}},
		{Name:"mixed4",SurfaceIndices:[]int{0,5,1,6}},
	}
	res:=UP191BResult{Schema:UP191BNativeResponseSchema,Experiment:"UP-191B-native-response-curve",SourceUP190BSeal:"bae5008f361cca459694cffc82fd2302162920fb",CalibrationPhases:cal,EvaluationPhases:eval,Profiles:profiles,FittedCorrectionUsed:false,HeldoutTuningUsed:false,AdaptiveProfileSearchUsed:false,PhaseInputUsed:false,ParityInputUsed:false,MaintenanceTriggered:false}
	pooled:=up184bBuild("pooled_26_30",cal)
	for _,p:=range profiles{for _,phase:=range eval{res.Points=append(res.Points,up191bPoint(pooled,p,phase))}}
	res.PearsonNativeCountRequiredFactor=up191bPearson(res.Points)
	counts:=[]int{};seen:=map[int]bool{}
	for _,p:=range res.Points{if !seen[p.NativeCorrectCount]{seen[p.NativeCorrectCount]=true;counts=append(counts,p.NativeCorrectCount)}}
	for _,count:=range counts{
		n:=0;sum:=0.0;minv,maxv:=0.0,0.0
		for _,p:=range res.Points{if p.NativeCorrectCount!=count{continue};if n==0{minv=p.RequiredMassFactor;maxv=p.RequiredMassFactor};n++;sum+=p.RequiredMassFactor;if p.RequiredMassFactor<minv{minv=p.RequiredMassFactor};if p.RequiredMassFactor>maxv{maxv=p.RequiredMassFactor}}
		spread:=maxv-minv;if spread>res.MaxWithinCountFactorSpread{res.MaxWithinCountFactorSpread=spread}
		res.CountGroups=append(res.CountGroups,UP191BCountGroup{NativeCorrectCount:count,Observations:n,MeanRequiredMassFactor:sum/float64(n),MinRequiredMassFactor:minv,MaxRequiredMassFactor:maxv,FactorSpread:spread})
	}
	return res,nil
}
