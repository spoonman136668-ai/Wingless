package unitary

import "math"

const UP194BOrderHistorySchema="wingless.up194b-order-history-effect.v1"

type UP194BPair struct{
	Name string `json:"name"`
	Forward []int `json:"forward"`
	Reverse []int `json:"reverse"`
}
type UP194BMetric struct{
	Pair string `json:"pair"`
	Phase int `json:"phase"`
	ForwardNativeCorrectCount int `json:"forward_native_correct_count"`
	ReverseNativeCorrectCount int `json:"reverse_native_correct_count"`
	ForwardRequiredFactor float64 `json:"forward_required_factor"`
	ReverseRequiredFactor float64 `json:"reverse_required_factor"`
	AbsoluteFactorDifference float64 `json:"absolute_factor_difference"`
	AbsoluteNativeCountDifference int `json:"absolute_native_count_difference"`
	SameNativeCorrectCount bool `json:"same_native_correct_count"`
}
type UP194BSummary struct{
	MeanAbsoluteFactorDifference float64 `json:"mean_absolute_factor_difference"`
	MaxAbsoluteFactorDifference float64 `json:"max_absolute_factor_difference"`
	SameCountPairs int `json:"same_count_pairs"`
	SameCountNonzeroFactorPairs int `json:"same_count_nonzero_factor_pairs"`
	MeanSameCountAbsoluteFactorDifference float64 `json:"mean_same_count_absolute_factor_difference"`
}
type UP194BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP193BSeal string `json:"source_up193b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Pairs []UP194BPair `json:"pairs"`
	NewCorrectionFit bool `json:"new_correction_fit"`
	CorrectionApplied bool `json:"correction_applied"`
	AdaptivePairSearchUsed bool `json:"adaptive_pair_search_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP194BMetric `json:"metrics"`
	Summary UP194BSummary `json:"summary"`
}
func RunUP194B()(UP194BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{49,50,51}
	pairs:=[]UP194BPair{
		{Name:"store4",Forward:[]int{0,1,2,3},Reverse:[]int{3,2,1,0}},
		{Name:"observe4",Forward:[]int{5,6,7,8},Reverse:[]int{8,7,6,5}},
		{Name:"mixed4",Forward:[]int{0,5,1,6},Reverse:[]int{6,1,5,0}},
		{Name:"cross3",Forward:[]int{0,5,13},Reverse:[]int{13,5,0}},
	}
	model:=up184bBuild("pooled_26_30",cal)
	res:=UP194BResult{Schema:UP194BOrderHistorySchema,Experiment:"UP-194B-order-history-effect",SourceUP193BSeal:"50aaa41f36c10283ca677124f5f07cb2b1be294f",CalibrationPhases:cal,EvaluationPhases:eval,Pairs:pairs,NewCorrectionFit:false,CorrectionApplied:false,AdaptivePairSearchUsed:false,MaintenanceTriggered:false}
	sum,maxv,sameSum:=0.0,0.0,0.0;sameN,sameNonzero:=0,0
	for _,p:=range pairs{for _,phase:=range eval{
		f:=up191bPoint(model,UP191BProfile{Name:p.Name+"_forward",SurfaceIndices:p.Forward},phase)
		r:=up191bPoint(model,UP191BProfile{Name:p.Name+"_reverse",SurfaceIndices:p.Reverse},phase)
		d:=math.Abs(f.RequiredMassFactor-r.RequiredMassFactor);cd:=f.NativeCorrectCount-r.NativeCorrectCount;if cd<0{cd=-cd};same:=cd==0
		res.Metrics=append(res.Metrics,UP194BMetric{Pair:p.Name,Phase:phase,ForwardNativeCorrectCount:f.NativeCorrectCount,ReverseNativeCorrectCount:r.NativeCorrectCount,ForwardRequiredFactor:f.RequiredMassFactor,ReverseRequiredFactor:r.RequiredMassFactor,AbsoluteFactorDifference:d,AbsoluteNativeCountDifference:cd,SameNativeCorrectCount:same})
		sum+=d;if d>maxv{maxv=d}
		if same{sameN++;sameSum+=d;if d>1e-12{sameNonzero++}}
	}}
	res.Summary.MeanAbsoluteFactorDifference=sum/float64(len(res.Metrics));res.Summary.MaxAbsoluteFactorDifference=maxv;res.Summary.SameCountPairs=sameN;res.Summary.SameCountNonzeroFactorPairs=sameNonzero
	if sameN>0{res.Summary.MeanSameCountAbsoluteFactorDifference=sameSum/float64(sameN)}
	return res,nil
}
