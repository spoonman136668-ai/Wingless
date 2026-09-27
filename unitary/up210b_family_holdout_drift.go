package unitary

import "math"

const UP210BDriftSchema="wingless.up210b-family-holdout-drift.v1"

type UP210BMetric struct{
	HeldoutFamily string `json:"heldout_family"`
	Phase int `json:"phase"`
	ActualDrift float64 `json:"actual_drift"`
	NativePredictedDrift float64 `json:"native_predicted_drift"`
	ZeroDriftAbsoluteError float64 `json:"zero_drift_absolute_error"`
	NativeDriftAbsoluteError float64 `json:"native_drift_absolute_error"`
}
type UP210BSummary struct{
	HeldoutFamily string `json:"heldout_family"`
	TrainingFamilies []string `json:"training_families"`
	TrainingPoints int `json:"training_points"`
	EvaluationPoints int `json:"evaluation_points"`
	ZeroDriftMeanAbsoluteError float64 `json:"zero_drift_mean_absolute_error"`
	NativeDriftMeanAbsoluteError float64 `json:"native_drift_mean_absolute_error"`
	ZeroDriftMaxAbsoluteError float64 `json:"zero_drift_max_absolute_error"`
	NativeDriftMaxAbsoluteError float64 `json:"native_drift_max_absolute_error"`
	NativeBetterThanZero int `json:"native_better_than_zero"`
	NonzeroActualDriftPoints int `json:"nonzero_actual_drift_points"`
	CorrectDriftDirection int `json:"correct_drift_direction"`
}
type UP210BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP209BSeal string `json:"source_up209b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Features []string `json:"features"`
	RidgeLambda float64 `json:"ridge_lambda"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	FamilyIdentityInputUsed bool `json:"family_identity_input_used"`
	AbsoluteStateInputUsed bool `json:"absolute_state_input_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	FutureStateInputUsed bool `json:"future_state_input_used"`
	Metrics []UP210BMetric `json:"metrics"`
	Summaries []UP210BSummary `json:"summaries"`
}
type up210bSample struct{x [4]float64;y float64}
func up210bDeltaFeature(phase int,indices []int)[4]float64{
	c:=up193bFeature(phase,indices);p:=up193bFeature(phase-1,indices)
	return [4]float64{c[0]-p[0],c[1]-p[1],c[2]-p[2],c[3]-p[3]}
}
func up210bFit(samples []up210bSample)(mean,std [4]float64,coef [5]float64,ok bool){
	for _,s:=range samples{for j:=0;j<4;j++{mean[j]+=s.x[j]}}
	for j:=0;j<4;j++{mean[j]/=float64(len(samples))}
	for _,s:=range samples{for j:=0;j<4;j++{d:=s.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(samples)));if std[j]==0{std[j]=1}}
	var a [5][5]float64;var b [5]float64
	for _,s:=range samples{
		v:=[5]float64{1,(s.x[0]-mean[0])/std[0],(s.x[1]-mean[1])/std[1],(s.x[2]-mean[2])/std[2],(s.x[3]-mean[3])/std[3]}
		for i:=0;i<5;i++{b[i]+=v[i]*s.y;for j:=0;j<5;j++{a[i][j]+=v[i]*v[j]}}
	}
	for i:=1;i<5;i++{a[i][i]+=1e-6}
	c,o:=up193bSolve(a,b);return mean,std,c,o
}
func up210bSign(v float64)int{if v>1e-12{return 1};if v< -1e-12{return -1};return 0}
func RunUP210B()(UP210BResult,error){
	train:=[]int{59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{97,98,99,100,101,102,103,104,105,106,107,108}
	specs:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP210BResult{
		Schema:UP210BDriftSchema,Experiment:"UP-210B-family-holdout-drift",SourceUP209BSeal:"3d286792ae090cbf748314d2a678532b73a77424",
		TrainingPhases:train,EvaluationPhases:eval,
		Features:[]string{"delta_native_correct_count","delta_mean_absolute_margin","delta_near_zero_margin_count","delta_min_absolute_margin"},RidgeLambda:1e-6,
		HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,FamilyIdentityInputUsed:false,AbsoluteStateInputUsed:false,PhaseInputUsed:false,ParityInputUsed:false,FutureStateInputUsed:false,
	}
	for hi,h:=range specs{
		samples:=[]up210bSample{};trainNames:=[]string{}
		for fi,f:=range specs{
			if fi==hi{continue};trainNames=append(trainNames,f.Name)
			for _,ph:=range train{
				now:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train",SurfaceIndices:f.Canonical},ph).RequiredMassFactor
				prev:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train_prev",SurfaceIndices:f.Canonical},ph-1).RequiredMassFactor
				samples=append(samples,up210bSample{x:up210bDeltaFeature(ph,f.Canonical),y:now-prev})
			}
		}
		mean,std,coef,ok:=up210bFit(samples);if !ok{return UP210BResult{},nil}
		s:=UP210BSummary{HeldoutFamily:h.Name,TrainingFamilies:trainNames,TrainingPoints:len(samples)}
		sumZ,sumN:=0.0,0.0
		for _,ph:=range eval{
			now:=up191bPoint(model,UP191BProfile{Name:h.Name+"_eval",SurfaceIndices:h.Canonical},ph).RequiredMassFactor
			prev:=up191bPoint(model,UP191BProfile{Name:h.Name+"_eval_prev",SurfaceIndices:h.Canonical},ph-1).RequiredMassFactor
			actual:=now-prev;x:=up210bDeltaFeature(ph,h.Canonical)
			z:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
			pred:=0.0;for i:=0;i<5;i++{pred+=coef[i]*z[i]}
			ze:=math.Abs(actual);ne:=math.Abs(pred-actual)
			res.Metrics=append(res.Metrics,UP210BMetric{HeldoutFamily:h.Name,Phase:ph,ActualDrift:actual,NativePredictedDrift:pred,ZeroDriftAbsoluteError:ze,NativeDriftAbsoluteError:ne})
			s.EvaluationPoints++;sumZ+=ze;sumN+=ne
			if ze>s.ZeroDriftMaxAbsoluteError{s.ZeroDriftMaxAbsoluteError=ze};if ne>s.NativeDriftMaxAbsoluteError{s.NativeDriftMaxAbsoluteError=ne}
			if ne<ze{s.NativeBetterThanZero++}
			if up210bSign(actual)!=0{s.NonzeroActualDriftPoints++;if up210bSign(pred)==up210bSign(actual){s.CorrectDriftDirection++}}
		}
		s.ZeroDriftMeanAbsoluteError=sumZ/float64(s.EvaluationPoints);s.NativeDriftMeanAbsoluteError=sumN/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
