package unitary

import (
	"math"
	"sort"
)

const UP196BPermutationSchema="wingless.up196b-permutation-distance.v1"

type UP196BComposition struct{
	Name string `json:"name"`
	Canonical []int `json:"canonical"`
}
type UP196BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	KendallDistance int `json:"kendall_distance"`
	Phase int `json:"phase"`
	NativeCorrectCount int `json:"native_correct_count"`
	CanonicalNativeCorrectCount int `json:"canonical_native_correct_count"`
	RequiredFactor float64 `json:"required_factor"`
	CanonicalRequiredFactor float64 `json:"canonical_required_factor"`
	AbsoluteFactorDifference float64 `json:"absolute_factor_difference"`
	SameNativeCorrectCount bool `json:"same_native_correct_count"`
}
type UP196BDistanceSummary struct{
	Composition string `json:"composition"`
	KendallDistance int `json:"kendall_distance"`
	Observations int `json:"observations"`
	MeanAbsoluteFactorDifference float64 `json:"mean_absolute_factor_difference"`
	MaxAbsoluteFactorDifference float64 `json:"max_absolute_factor_difference"`
}
type UP196BCompositionSummary struct{
	Composition string `json:"composition"`
	PearsonDistanceFactorDifference float64 `json:"pearson_distance_factor_difference"`
	SameCountObservations int `json:"same_count_observations"`
	SameCountNonzeroFactorObservations int `json:"same_count_nonzero_factor_observations"`
}
type UP196BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP195BSeal string `json:"source_up195b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Compositions []UP196BComposition `json:"compositions"`
	DistanceMetric string `json:"distance_metric"`
	NewCorrectionFit bool `json:"new_correction_fit"`
	AdaptivePermutationSelectionUsed bool `json:"adaptive_permutation_selection_used"`
	AdaptiveDistanceSearchUsed bool `json:"adaptive_distance_search_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP196BMetric `json:"metrics"`
	DistanceSummaries []UP196BDistanceSummary `json:"distance_summaries"`
	CompositionSummaries []UP196BCompositionSummary `json:"composition_summaries"`
}
func up196bPermutations(in []int)[][]int{
	out:=[][]int{}
	a:=append([]int(nil),in...)
	var rec func(int)
	rec=func(k int){
		if k==len(a){out=append(out,append([]int(nil),a...));return}
		for i:=k;i<len(a);i++{a[k],a[i]=a[i],a[k];rec(k+1);a[k],a[i]=a[i],a[k]}
	}
	rec(0);return out
}
func up196bDistance(canonical,perm []int)int{
	pos:=map[int]int{};for i,v:=range canonical{pos[v]=i}
	d:=0
	for i:=0;i<len(perm);i++{for j:=i+1;j<len(perm);j++{if pos[perm[i]]>pos[perm[j]]{d++}}}
	return d
}
func up196bPearson(ms []UP196BMetric)float64{
	if len(ms)<2{return 0}
	mx,my:=0.0,0.0
	for _,m:=range ms{mx+=float64(m.KendallDistance);my+=m.AbsoluteFactorDifference}
	n:=float64(len(ms));mx/=n;my/=n
	num,dx,dy:=0.0,0.0,0.0
	for _,m:=range ms{x:=float64(m.KendallDistance)-mx;y:=m.AbsoluteFactorDifference-my;num+=x*y;dx+=x*x;dy+=y*y}
	if dx==0||dy==0{return 0};return num/math.Sqrt(dx*dy)
}
func RunUP196B()(UP196BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{55,56,57}
	comps:=[]UP196BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
	}
	model:=up184bBuild("pooled_26_30",cal)
	res:=UP196BResult{Schema:UP196BPermutationSchema,Experiment:"UP-196B-permutation-distance",SourceUP195BSeal:"44faf848f60913208fc6089932ba977f9f97e5c8",CalibrationPhases:cal,EvaluationPhases:eval,Compositions:comps,DistanceMetric:"kendall_inversion_count",NewCorrectionFit:false,AdaptivePermutationSelectionUsed:false,AdaptiveDistanceSearchUsed:false,MaintenanceTriggered:false}
	type agg struct{n int;sum,max float64}
	groups:=map[string]map[int]*agg{}
	byComp:=map[string][]UP196BMetric{}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		for _,phase:=range eval{
			base:=up191bPoint(model,UP191BProfile{Name:c.Name+"_canonical",SurfaceIndices:c.Canonical},phase)
			for _,p:=range perms{
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_perm",SurfaceIndices:p},phase)
				d:=up196bDistance(c.Canonical,p);fd:=math.Abs(pt.RequiredMassFactor-base.RequiredMassFactor);same:=pt.NativeCorrectCount==base.NativeCorrectCount
				m:=UP196BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),KendallDistance:d,Phase:phase,NativeCorrectCount:pt.NativeCorrectCount,CanonicalNativeCorrectCount:base.NativeCorrectCount,RequiredFactor:pt.RequiredMassFactor,CanonicalRequiredFactor:base.RequiredMassFactor,AbsoluteFactorDifference:fd,SameNativeCorrectCount:same}
				res.Metrics=append(res.Metrics,m);byComp[c.Name]=append(byComp[c.Name],m)
				if groups[c.Name]==nil{groups[c.Name]=map[int]*agg{}}
				a:=groups[c.Name][d];if a==nil{a=&agg{};groups[c.Name][d]=a};a.n++;a.sum+=fd;if fd>a.max{a.max=fd}
			}
		}
	}
	for _,c:=range comps{
		distances:=make([]int,0,len(groups[c.Name]));for d:=range groups[c.Name]{distances=append(distances,d)};sort.Ints(distances)
		for _,d:=range distances{a:=groups[c.Name][d];res.DistanceSummaries=append(res.DistanceSummaries,UP196BDistanceSummary{Composition:c.Name,KendallDistance:d,Observations:a.n,MeanAbsoluteFactorDifference:a.sum/float64(a.n),MaxAbsoluteFactorDifference:a.max})}
		same,nonzero:=0,0;for _,m:=range byComp[c.Name]{if m.SameNativeCorrectCount{same++;if m.AbsoluteFactorDifference>1e-12{nonzero++}}}
		res.CompositionSummaries=append(res.CompositionSummaries,UP196BCompositionSummary{Composition:c.Name,PearsonDistanceFactorDifference:up196bPearson(byComp[c.Name]),SameCountObservations:same,SameCountNonzeroFactorObservations:nonzero})
	}
	return res,nil
}
