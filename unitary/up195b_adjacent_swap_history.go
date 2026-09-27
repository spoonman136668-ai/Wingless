package unitary

import (
	"math"
	"sort"
)

const UP195BAdjacentSwapSchema="wingless.up195b-adjacent-swap-history.v1"

type UP195BComposition struct{
	Name string `json:"name"`
	Canonical []int `json:"canonical"`
}
type UP195BMetric struct{
	Composition string `json:"composition"`
	SwapPosition int `json:"swap_position"`
	Phase int `json:"phase"`
	CanonicalNativeCorrectCount int `json:"canonical_native_correct_count"`
	SwappedNativeCorrectCount int `json:"swapped_native_correct_count"`
	CanonicalRequiredFactor float64 `json:"canonical_required_factor"`
	SwappedRequiredFactor float64 `json:"swapped_required_factor"`
	AbsoluteFactorDifference float64 `json:"absolute_factor_difference"`
	AbsoluteNativeCountDifference int `json:"absolute_native_count_difference"`
	SameNativeCorrectCount bool `json:"same_native_correct_count"`
}
type UP195BGroupSummary struct{
	Group string `json:"group"`
	Value string `json:"value"`
	Observations int `json:"observations"`
	MeanAbsoluteFactorDifference float64 `json:"mean_absolute_factor_difference"`
	MaxAbsoluteFactorDifference float64 `json:"max_absolute_factor_difference"`
}
type UP195BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP194BSeal string `json:"source_up194b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Compositions []UP195BComposition `json:"compositions"`
	NewCorrectionFit bool `json:"new_correction_fit"`
	AdaptiveSwapSearchUsed bool `json:"adaptive_swap_search_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP195BMetric `json:"metrics"`
	GroupSummaries []UP195BGroupSummary `json:"group_summaries"`
	SameCountSwapObservations int `json:"same_count_swap_observations"`
	SameCountNonzeroFactorObservations int `json:"same_count_nonzero_factor_observations"`
}
func up195bSwap(in []int,pos int)[]int{
	out:=append([]int(nil),in...)
	out[pos],out[pos+1]=out[pos+1],out[pos]
	return out
}
func RunUP195B()(UP195BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{52,53,54}
	comps:=[]UP195BComposition{
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	model:=up184bBuild("pooled_26_30",cal)
	res:=UP195BResult{Schema:UP195BAdjacentSwapSchema,Experiment:"UP-195B-adjacent-swap-history",SourceUP194BSeal:"4e2636ed3871b6ec2adc08318ede8567c302a320",CalibrationPhases:cal,EvaluationPhases:eval,Compositions:comps,NewCorrectionFit:false,AdaptiveSwapSearchUsed:false,MaintenanceTriggered:false}
	type agg struct{n int;sum,max float64}
	byComp:=map[string]*agg{};byPos:=map[int]*agg{}
	for _,c:=range comps{
		for pos:=0;pos<len(c.Canonical)-1;pos++{
			sw:=up195bSwap(c.Canonical,pos)
			for _,phase:=range eval{
				base:=up191bPoint(model,UP191BProfile{Name:c.Name+"_canonical",SurfaceIndices:c.Canonical},phase)
				alt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_swap",SurfaceIndices:sw},phase)
				d:=math.Abs(base.RequiredMassFactor-alt.RequiredMassFactor);cd:=base.NativeCorrectCount-alt.NativeCorrectCount;if cd<0{cd=-cd};same:=cd==0
				res.Metrics=append(res.Metrics,UP195BMetric{Composition:c.Name,SwapPosition:pos,Phase:phase,CanonicalNativeCorrectCount:base.NativeCorrectCount,SwappedNativeCorrectCount:alt.NativeCorrectCount,CanonicalRequiredFactor:base.RequiredMassFactor,SwappedRequiredFactor:alt.RequiredMassFactor,AbsoluteFactorDifference:d,AbsoluteNativeCountDifference:cd,SameNativeCorrectCount:same})
				if same{res.SameCountSwapObservations++;if d>1e-12{res.SameCountNonzeroFactorObservations++}}
				a:=byComp[c.Name];if a==nil{a=&agg{};byComp[c.Name]=a};a.n++;a.sum+=d;if d>a.max{a.max=d}
				p:=byPos[pos];if p==nil{p=&agg{};byPos[pos]=p};p.n++;p.sum+=d;if d>p.max{p.max=d}
			}
		}
	}
	names:=make([]string,0,len(byComp));for k:=range byComp{names=append(names,k)};sort.Strings(names)
	for _,k:=range names{a:=byComp[k];res.GroupSummaries=append(res.GroupSummaries,UP195BGroupSummary{Group:"composition",Value:k,Observations:a.n,MeanAbsoluteFactorDifference:a.sum/float64(a.n),MaxAbsoluteFactorDifference:a.max})}
	poses:=make([]int,0,len(byPos));for k:=range byPos{poses=append(poses,k)};sort.Ints(poses)
	for _,k:=range poses{a:=byPos[k];res.GroupSummaries=append(res.GroupSummaries,UP195BGroupSummary{Group:"swap_position",Value:string(rune('0'+k)),Observations:a.n,MeanAbsoluteFactorDifference:a.sum/float64(a.n),MaxAbsoluteFactorDifference:a.max})}
	return res,nil
}
