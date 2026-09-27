package unitary

import "math"

const UP211BInstabilitySchema="wingless.up211b-family-holdout-instability-magnitude.v1"

type UP211BMetric struct{
	HeldoutFamily string `json:"heldout_family"`
	Phase int `json:"phase"`
	ActualMagnitude float64 `json:"actual_magnitude"`
	TrainingMeanPredictedMagnitude float64 `json:"training_mean_predicted_magnitude"`
	NativePredictedMagnitude float64 `json:"native_predicted_magnitude"`
	ZeroAbsoluteError float64 `json:"zero_absolute_error"`
	TrainingMeanAbsoluteError float64 `json:"training_mean_absolute_error"`
	NativeAbsoluteError float64 `json:"native_absolute_error"`
}
type UP211BSummary struct{
	HeldoutFamily string `json:"heldout_family"`
	TrainingPoints int `json:"training_points"`
	EvaluationPoints int `json:"evaluation_points"`
	ZeroMeanAbsoluteError float64 `json:"zero_mean_absolute_error"`
	TrainingMeanAbsoluteError float64 `json:"training_mean_absolute_error"`
	NativeMeanAbsoluteError float64 `json:"native_mean_absolute_error"`
	ZeroMaxAbsoluteError float64 `json:"zero_max_absolute_error"`
	TrainingMeanMaxAbsoluteError float64 `json:"training_mean_max_absolute_error"`
	NativeMaxAbsoluteError float64 `json:"native_max_absolute_error"`
	NativeBetterThanTrainingMean int `json:"native_better_than_training_mean"`
	PearsonCorrelation float64 `json:"pearson_correlation"`
}
type UP211BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP210BSeal string `json:"source_up210b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Features []string `json:"features"`
	RidgeLambda float64 `json:"ridge_lambda"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	FamilyIdentityInputUsed bool `json:"family_identity_input_used"`
	SignedDriftInputUsed bool `json:"signed_drift_input_used"`
	AbsoluteStateInputUsed bool `json:"absolute_state_input_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	FutureStateInputUsed bool `json:"future_state_input_used"`
	Metrics []UP211BMetric `json:"metrics"`
	Summaries []UP211BSummary `json:"summaries"`
}
func up211bAbsDeltaFeature(phase int,indices []int)[4]float64{
	x:=up210bDeltaFeature(phase,indices)
	return [4]float64{math.Abs(x[0]),math.Abs(x[1]),math.Abs(x[2]),math.Abs(x[3])}
}
func up211bCorr(xs,ys []float64)float64{
	if len(xs)==0||len(xs)!=len(ys){return 0}
	mx,my:=0.0,0.0
	for i:=range xs{mx+=xs[i];my+=ys[i]}
	mx/=float64(len(xs));my/=float64(len(ys))
	num,dx,dy:=0.0,0.0,0.0
	for i:=range xs{a:=xs[i]-mx;b:=ys[i]-my;num+=a*b;dx+=a*a;dy+=b*b}
	if dx==0||dy==0{return 0}
	return num/math.Sqrt(dx*dy)
}
func RunUP211B()(UP211BResult,error){
	train:=[]int{59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{109,110,111,112,113,114,115,116,117,118,119,120}
	specs:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP211BResult{Schema:UP211BInstabilitySchema,Experiment:"UP-211B-family-holdout-instability-magnitude",SourceUP210BSeal:"7842b7338be06a648f326fad3b21a5ffd828304e",TrainingPhases:train,EvaluationPhases:eval,Features:[]string{"abs_delta_native_correct_count","abs_delta_mean_absolute_margin","abs_delta_near_zero_margin_count","abs_delta_min_absolute_margin"},RidgeLambda:1e-6,HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,FamilyIdentityInputUsed:false,SignedDriftInputUsed:false,AbsoluteStateInputUsed:false,PhaseInputUsed:false,ParityInputUsed:false,FutureStateInputUsed:false}
	for hi,h:=range specs{
		samples:=[]up210bSample{};trainMean:=0.0
		for fi,f:=range specs{
			if fi==hi{continue}
			for _,ph:=range train{
				now:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train",SurfaceIndices:f.Canonical},ph).RequiredMassFactor
				prev:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train_prev",SurfaceIndices:f.Canonical},ph-1).RequiredMassFactor
				y:=math.Abs(now-prev);samples=append(samples,up210bSample{x:up211bAbsDeltaFeature(ph,f.Canonical),y:y});trainMean+=y
			}
		}
		trainMean/=float64(len(samples))
		mean,std,coef,ok:=up210bFit(samples);if !ok{return UP211BResult{},nil}
		sm:=UP211BSummary{HeldoutFamily:h.Name,TrainingPoints:len(samples)}
		sumZ,sumT,sumN:=0.0,0.0,0.0;actuals:=[]float64{};preds:=[]float64{}
		for _,ph:=range eval{
			now:=up191bPoint(model,UP191BProfile{Name:h.Name+"_eval",SurfaceIndices:h.Canonical},ph).RequiredMassFactor
			prev:=up191bPoint(model,UP191BProfile{Name:h.Name+"_eval_prev",SurfaceIndices:h.Canonical},ph-1).RequiredMassFactor
			actual:=math.Abs(now-prev);x:=up211bAbsDeltaFeature(ph,h.Canonical)
			z:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
			pred:=0.0;for i:=0;i<5;i++{pred+=coef[i]*z[i]};if pred<0{pred=0}
			ze:=actual;te:=math.Abs(trainMean-actual);ne:=math.Abs(pred-actual)
			res.Metrics=append(res.Metrics,UP211BMetric{HeldoutFamily:h.Name,Phase:ph,ActualMagnitude:actual,TrainingMeanPredictedMagnitude:trainMean,NativePredictedMagnitude:pred,ZeroAbsoluteError:ze,TrainingMeanAbsoluteError:te,NativeAbsoluteError:ne})
			sm.EvaluationPoints++;sumZ+=ze;sumT+=te;sumN+=ne;actuals=append(actuals,actual);preds=append(preds,pred)
			if ze>sm.ZeroMaxAbsoluteError{sm.ZeroMaxAbsoluteError=ze};if te>sm.TrainingMeanMaxAbsoluteError{sm.TrainingMeanMaxAbsoluteError=te};if ne>sm.NativeMaxAbsoluteError{sm.NativeMaxAbsoluteError=ne}
			if ne<te{sm.NativeBetterThanTrainingMean++}
		}
		sm.ZeroMeanAbsoluteError=sumZ/float64(sm.EvaluationPoints);sm.TrainingMeanAbsoluteError=sumT/float64(sm.EvaluationPoints);sm.NativeMeanAbsoluteError=sumN/float64(sm.EvaluationPoints);sm.PearsonCorrelation=up211bCorr(preds,actuals)
		res.Summaries=append(res.Summaries,sm)
	}
	return res,nil
}
