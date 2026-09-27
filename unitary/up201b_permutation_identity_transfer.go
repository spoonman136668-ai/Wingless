package unitary

import (
	"fmt"
	"math"
)

const UP201BIdentitySchema="wingless.up201b-permutation-identity-transfer.v1"

type UP201BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	BaselinePredictedFactor float64 `json:"baseline_predicted_factor"`
	IdentityPredictedFactor float64 `json:"identity_predicted_factor"`
	BaselineAbsoluteError float64 `json:"baseline_absolute_error"`
	IdentityAbsoluteError float64 `json:"identity_absolute_error"`
}
type UP201BSummary struct{
	Composition string `json:"composition"`
	EvaluationPoints int `json:"evaluation_points"`
	BaselineMeanAbsoluteError float64 `json:"baseline_mean_absolute_error"`
	IdentityMeanAbsoluteError float64 `json:"identity_mean_absolute_error"`
	BaselineMaxAbsoluteError float64 `json:"baseline_max_absolute_error"`
	IdentityMaxAbsoluteError float64 `json:"identity_max_absolute_error"`
	IdentityBetterThanBaseline int `json:"identity_better_than_baseline"`
}
type UP201BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP200BSeal string `json:"source_up200b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	PermutationCountPerComposition int `json:"permutation_count_per_composition"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	FeatureExtractionUsed bool `json:"feature_extraction_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP201BMetric `json:"metrics"`
	Summaries []UP201BSummary `json:"summaries"`
}
func up201bKey(p []int)string{return fmt.Sprint(p)}
func RunUP201B()(UP201BResult,error){
	train:=[]int{55,56,57};eval:=[]int{70,71,72}
	comps:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP201BResult{Schema:UP201BIdentitySchema,Experiment:"UP-201B-permutation-identity-transfer",SourceUP200BSeal:"56e9730fa7863308566b4c2489f8dd9baca8d8d0",TrainingPhases:train,EvaluationPhases:eval,PermutationCountPerComposition:24,HeldoutFittingUsed:false,FeatureExtractionUsed:false,PhaseInputUsed:false,ParityInputUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical);means:=map[string]float64{};base:=0.0
		for _,p:=range perms{
			s:=0.0
			for _,phase:=range train{
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_train",SurfaceIndices:p},phase)
				s+=pt.RequiredMassFactor;base+=pt.RequiredMassFactor
			}
			means[up201bKey(p)]=s/float64(len(train))
		}
		base/=float64(len(perms)*len(train))
		sumB,sumI:=0.0,0.0;s:=UP201BSummary{Composition:c.Name}
		for _,phase:=range eval{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase)
			ip:=means[up201bKey(p)];be:=math.Abs(base-pt.RequiredMassFactor);ie:=math.Abs(ip-pt.RequiredMassFactor)
			res.Metrics=append(res.Metrics,UP201BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,BaselinePredictedFactor:base,IdentityPredictedFactor:ip,BaselineAbsoluteError:be,IdentityAbsoluteError:ie})
			s.EvaluationPoints++;sumB+=be;sumI+=ie
			if be>s.BaselineMaxAbsoluteError{s.BaselineMaxAbsoluteError=be};if ie>s.IdentityMaxAbsoluteError{s.IdentityMaxAbsoluteError=ie};if ie<be{s.IdentityBetterThanBaseline++}
		}}
		s.BaselineMeanAbsoluteError=sumB/float64(s.EvaluationPoints);s.IdentityMeanAbsoluteError=sumI/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
