package unitary

import "math"

const UP197BTransitionSchema="wingless.up197b-directed-transition-calibration.v1"

type UP197BComposition struct{
	Name string `json:"name"`
	Canonical []int `json:"canonical"`
}
type UP197BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	TransitionPredictedFactor float64 `json:"transition_predicted_factor"`
	BaselinePredictedFactor float64 `json:"baseline_predicted_factor"`
	TransitionAbsoluteError float64 `json:"transition_absolute_error"`
	BaselineAbsoluteError float64 `json:"baseline_absolute_error"`
}
type UP197BSummary struct{
	Composition string `json:"composition"`
	TrainingMeanFactor float64 `json:"training_mean_factor"`
	TransitionMeanAbsoluteError float64 `json:"transition_mean_absolute_error"`
	BaselineMeanAbsoluteError float64 `json:"baseline_mean_absolute_error"`
	TransitionMaxAbsoluteError float64 `json:"transition_max_absolute_error"`
	BaselineMaxAbsoluteError float64 `json:"baseline_max_absolute_error"`
	TransitionBetterPoints int `json:"transition_better_points"`
	EvaluationPoints int `json:"evaluation_points"`
}
type UP197BModel struct{
	Composition string `json:"composition"`
	Coefficients []float64 `json:"coefficients"`
}
type UP197BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP196BSeal string `json:"source_up196b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Compositions []UP197BComposition `json:"compositions"`
	FeatureCount int `json:"feature_count"`
	RidgeLambda float64 `json:"ridge_lambda"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Models []UP197BModel `json:"models"`
	Metrics []UP197BMetric `json:"metrics"`
	Summaries []UP197BSummary `json:"summaries"`
}
func up197bPairIndex(a,b int)int{
	k:=0
	for i:=0;i<4;i++{for j:=0;j<4;j++{if i==j{continue};if i==a&&j==b{return k};k++}}
	return -1
}
func up197bFeatures(canonical,perm []int)[12]float64{
	pos:=map[int]int{};for i,v:=range canonical{pos[v]=i}
	var x [12]float64
	for i:=0;i<len(perm)-1;i++{a,b:=pos[perm[i]],pos[perm[i+1]];idx:=up197bPairIndex(a,b);if idx>=0{x[idx]=1}}
	return x
}
func up197bSolve(a [13][13]float64,b [13]float64)([13]float64,bool){
	for i:=0;i<13;i++{
		p:=i;best:=math.Abs(a[i][i]);for r:=i+1;r<13;r++{if v:=math.Abs(a[r][i]);v>best{best=v;p=r}}
		if best<1e-12{return [13]float64{},false}
		if p!=i{a[i],a[p]=a[p],a[i];b[i],b[p]=b[p],b[i]}
		d:=a[i][i];for c:=i;c<13;c++{a[i][c]/=d};b[i]/=d
		for r:=0;r<13;r++{if r==i{continue};f:=a[r][i];for c:=i;c<13;c++{a[r][c]-=f*a[i][c]};b[r]-=f*b[i]}
	}
	return b,true
}
func up197bPredict(coef [13]float64,x [12]float64)float64{
	y:=coef[0];for i:=0;i<12;i++{y+=coef[i+1]*x[i]};return y
}
func RunUP197B()(UP197BResult,error){
	train:=[]int{55,56,57};eval:=[]int{58,59,60}
	comps:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP197BResult{Schema:UP197BTransitionSchema,Experiment:"UP-197B-directed-transition-calibration",SourceUP196BSeal:"6a4ed7ef52e69cf835d143402d670e13800793c6",TrainingPhases:train,EvaluationPhases:eval,Compositions:comps,FeatureCount:12,RidgeLambda:1e-6,HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,PhaseInputUsed:false,ParityInputUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		var a [13][13]float64;var b [13]float64;meanY:=0.0;nTrain:=0
		for _,phase:=range train{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_train",SurfaceIndices:p},phase);x:=up197bFeatures(c.Canonical,p);v:=[13]float64{};v[0]=1;for i:=0;i<12;i++{v[i+1]=x[i]}
			for i:=0;i<13;i++{b[i]+=v[i]*pt.RequiredMassFactor;for j:=0;j<13;j++{a[i][j]+=v[i]*v[j]}}
			meanY+=pt.RequiredMassFactor;nTrain++
		}}
		meanY/=float64(nTrain);for i:=1;i<13;i++{a[i][i]+=1e-6}
		coef,ok:=up197bSolve(a,b);if !ok{return UP197BResult{},nil}
		res.Models=append(res.Models,UP197BModel{Composition:c.Name,Coefficients:append([]float64(nil),coef[:]...)})
		s:=UP197BSummary{Composition:c.Name,TrainingMeanFactor:meanY};sumT,sumB:=0.0,0.0
		for _,phase:=range eval{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase);pred:=up197bPredict(coef,up197bFeatures(c.Canonical,p));te:=math.Abs(pred-pt.RequiredMassFactor);be:=math.Abs(meanY-pt.RequiredMassFactor)
			res.Metrics=append(res.Metrics,UP197BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,TransitionPredictedFactor:pred,BaselinePredictedFactor:meanY,TransitionAbsoluteError:te,BaselineAbsoluteError:be})
			sumT+=te;sumB+=be;s.EvaluationPoints++;if te<be{s.TransitionBetterPoints++};if te>s.TransitionMaxAbsoluteError{s.TransitionMaxAbsoluteError=te};if be>s.BaselineMaxAbsoluteError{s.BaselineMaxAbsoluteError=be}
		}}
		s.TransitionMeanAbsoluteError=sumT/float64(s.EvaluationPoints);s.BaselineMeanAbsoluteError=sumB/float64(s.EvaluationPoints);res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
