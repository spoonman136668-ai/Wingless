package unitary

import "math"

const UP202BTrendSchema="wingless.up202b-permutation-trend-transfer.v1"

type UP202BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	BaselinePredictedFactor float64 `json:"baseline_predicted_factor"`
	IdentityMeanPredictedFactor float64 `json:"identity_mean_predicted_factor"`
	IdentityTrendPredictedFactor float64 `json:"identity_trend_predicted_factor"`
	BaselineAbsoluteError float64 `json:"baseline_absolute_error"`
	IdentityMeanAbsoluteError float64 `json:"identity_mean_absolute_error"`
	IdentityTrendAbsoluteError float64 `json:"identity_trend_absolute_error"`
}

type UP202BSummary struct{
	Composition string `json:"composition"`
	EvaluationPoints int `json:"evaluation_points"`
	BaselineMeanAbsoluteError float64 `json:"baseline_mean_absolute_error"`
	IdentityMeanMeanAbsoluteError float64 `json:"identity_mean_mean_absolute_error"`
	IdentityTrendMeanAbsoluteError float64 `json:"identity_trend_mean_absolute_error"`
	BaselineMaxAbsoluteError float64 `json:"baseline_max_absolute_error"`
	IdentityMeanMaxAbsoluteError float64 `json:"identity_mean_max_absolute_error"`
	IdentityTrendMaxAbsoluteError float64 `json:"identity_trend_max_absolute_error"`
	TrendBetterThanIdentityMean int `json:"trend_better_than_identity_mean"`
	TrendBetterThanBaseline int `json:"trend_better_than_baseline"`
}

type UP202BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP201BSeal string `json:"source_up201b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	PermutationCountPerComposition int `json:"permutation_count_per_composition"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveTrendSelectionUsed bool `json:"adaptive_trend_selection_used"`
	NonlinearTrendUsed bool `json:"nonlinear_trend_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP202BMetric `json:"metrics"`
	Summaries []UP202BSummary `json:"summaries"`
}

type up202bLine struct{Mean,Slope,Center float64}

func up202bFitLine(phases []int, ys []float64) up202bLine{
	center:=0.0
	for _,p:=range phases{center+=float64(p)}
	center/=float64(len(phases))
	mean:=0.0
	for _,y:=range ys{mean+=y}
	mean/=float64(len(ys))
	num,den:=0.0,0.0
	for i,p:=range phases{
		x:=float64(p)-center
		num+=x*(ys[i]-mean)
		den+=x*x
	}
	slope:=0.0
	if den>0{slope=num/den}
	return up202bLine{Mean:mean,Slope:slope,Center:center}
}
func (l up202bLine) predict(phase int)float64{return l.Mean+l.Slope*(float64(phase)-l.Center)}

func RunUP202B()(UP202BResult,error){
	train:=[]int{55,56,57};eval:=[]int{73,74,75}
	comps:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP202BResult{Schema:UP202BTrendSchema,Experiment:"UP-202B-permutation-trend-transfer",SourceUP201BSeal:"6529a78c26a1ce8a86830918dd9a4673b94b3fe1",TrainingPhases:train,EvaluationPhases:eval,PermutationCountPerComposition:24,HeldoutFittingUsed:false,AdaptiveTrendSelectionUsed:false,NonlinearTrendUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		lines:=map[string]up202bLine{}
		base:=0.0
		for _,p:=range perms{
			ys:=make([]float64,0,len(train))
			for _,phase:=range train{
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_train",SurfaceIndices:p},phase)
				ys=append(ys,pt.RequiredMassFactor);base+=pt.RequiredMassFactor
			}
			lines[up201bKey(p)]=up202bFitLine(train,ys)
		}
		base/=float64(len(perms)*len(train))
		s:=UP202BSummary{Composition:c.Name}
		sumB,sumI,sumT:=0.0,0.0,0.0
		for _,phase:=range eval{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase)
			line:=lines[up201bKey(p)]
			ip:=line.Mean;tp:=line.predict(phase)
			be:=math.Abs(base-pt.RequiredMassFactor);ie:=math.Abs(ip-pt.RequiredMassFactor);te:=math.Abs(tp-pt.RequiredMassFactor)
			res.Metrics=append(res.Metrics,UP202BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,BaselinePredictedFactor:base,IdentityMeanPredictedFactor:ip,IdentityTrendPredictedFactor:tp,BaselineAbsoluteError:be,IdentityMeanAbsoluteError:ie,IdentityTrendAbsoluteError:te})
			s.EvaluationPoints++;sumB+=be;sumI+=ie;sumT+=te
			if be>s.BaselineMaxAbsoluteError{s.BaselineMaxAbsoluteError=be}
			if ie>s.IdentityMeanMaxAbsoluteError{s.IdentityMeanMaxAbsoluteError=ie}
			if te>s.IdentityTrendMaxAbsoluteError{s.IdentityTrendMaxAbsoluteError=te}
			if te<ie{s.TrendBetterThanIdentityMean++}
			if te<be{s.TrendBetterThanBaseline++}
		}}
		s.BaselineMeanAbsoluteError=sumB/float64(s.EvaluationPoints)
		s.IdentityMeanMeanAbsoluteError=sumI/float64(s.EvaluationPoints)
		s.IdentityTrendMeanAbsoluteError=sumT/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
