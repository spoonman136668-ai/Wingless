package unitary

import "math"

const UP214BLateRoutedSchema="wingless.up214b-late-routed-offset-transfer.v1"

type UP214BMetric struct{
	TrueRegime string `json:"true_regime"`
	PredictedRegime string `json:"predicted_regime"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	StaticPredictedFactor float64 `json:"static_predicted_factor"`
	RoutedPredictedFactor float64 `json:"routed_predicted_factor"`
	OracleFamilyPredictedFactor float64 `json:"oracle_family_predicted_factor"`
	OracleAnchorPredictedFactor float64 `json:"oracle_anchor_predicted_factor"`
	StaticAbsoluteError float64 `json:"static_absolute_error"`
	RoutedAbsoluteError float64 `json:"routed_absolute_error"`
	OracleFamilyAbsoluteError float64 `json:"oracle_family_absolute_error"`
	OracleAnchorAbsoluteError float64 `json:"oracle_anchor_absolute_error"`
}
type UP214BSummary struct{
	Regime string `json:"regime"`
	RouteCorrect int `json:"route_correct"`
	RouteTotal int `json:"route_total"`
	EvaluationPoints int `json:"evaluation_points"`
	StaticMeanAbsoluteError float64 `json:"static_mean_absolute_error"`
	RoutedMeanAbsoluteError float64 `json:"routed_mean_absolute_error"`
	OracleFamilyMeanAbsoluteError float64 `json:"oracle_family_mean_absolute_error"`
	OracleAnchorMeanAbsoluteError float64 `json:"oracle_anchor_mean_absolute_error"`
	StaticMaxAbsoluteError float64 `json:"static_max_absolute_error"`
	RoutedMaxAbsoluteError float64 `json:"routed_max_absolute_error"`
	OracleFamilyMaxAbsoluteError float64 `json:"oracle_family_max_absolute_error"`
	OracleAnchorMaxAbsoluteError float64 `json:"oracle_anchor_max_absolute_error"`
	RoutedBetterThanStatic int `json:"routed_better_than_static"`
}
type UP214BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP213BSeal string `json:"source_up213b_seal"`
	TemplatePhases []int `json:"template_phases"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Features []string `json:"features"`
	Router string `json:"router"`
	RidgeLambda float64 `json:"ridge_lambda"`
	RouteDecisions int `json:"route_decisions"`
	RouteCorrect int `json:"route_correct"`
	RouteAccuracy float64 `json:"route_accuracy"`
	EvaluationLabelFittingUsed bool `json:"evaluation_label_fitting_used"`
	AdaptiveModelSelectionUsed bool `json:"adaptive_model_selection_used"`
	RetrainingUsed bool `json:"retraining_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	AnchorMeasurementUsedByRouted bool `json:"anchor_measurement_used_by_routed"`
	AnchorMeasurementUsedByOracleFamily bool `json:"anchor_measurement_used_by_oracle_family"`
	NonlinearModelUsed bool `json:"nonlinear_model_used"`
	Metrics []UP214BMetric `json:"metrics"`
	Summaries []UP214BSummary `json:"summaries"`
}
func RunUP214B()(UP214BResult,error){
	template:=[]int{55,56,57};train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{136,137,138,139,140,141,142,143,144,145,146,147,148,149,150,151,152,153,154,155,156,157,158,159}
	specs:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	type raw struct{name string;x [4]float64}
	raws:=[]raw{}
	for _,f:=range specs{for _,ph:=range train{raws=append(raws,raw{name:f.Name,x:up193bFeature(ph,f.Canonical)})}}
	var mean,std [4]float64
	for _,r:=range raws{for j:=0;j<4;j++{mean[j]+=r.x[j]}}
	for j:=0;j<4;j++{mean[j]/=float64(len(raws))}
	for _,r:=range raws{for j:=0;j<4;j++{d:=r.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(raws)));if std[j]==0{std[j]=1}}
	fams:=[]up213bFamily{}
	for _,f:=range specs{
		perms:=up196bPermutations(f.Canonical);means:=make([]float64,len(perms));canon:=-1
		for i,p:=range perms{
			if up201bKey(p)==up201bKey(f.Canonical){canon=i}
			s:=0.0;for _,ph:=range template{s+=up191bPoint(model,UP191BProfile{Name:f.Name+"_template",SurfaceIndices:p},ph).RequiredMassFactor}
			means[i]=s/float64(len(template))
		}
		ff:=up213bFamily{name:f.Name,canonical:f.Canonical,perms:perms,means:means,canon:canon}
		n:=0
		for _,r:=range raws{if r.name==f.Name{for j:=0;j<4;j++{ff.centroid[j]+=(r.x[j]-mean[j])/std[j]};n++}}
		for j:=0;j<4;j++{ff.centroid[j]/=float64(n)}
		var a [5][5]float64;var b [5]float64
		for _,ph:=range train{
			x:=up193bFeature(ph,f.Canonical);v:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
			pt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train",SurfaceIndices:f.Canonical},ph)
			y:=pt.RequiredMassFactor-means[canon]
			for i:=0;i<5;i++{b[i]+=v[i]*y;for j:=0;j<5;j++{a[i][j]+=v[i]*v[j]}}
		}
		for i:=1;i<5;i++{a[i][i]+=1e-6}
		c,ok:=up193bSolve(a,b);if !ok{return UP214BResult{},nil};ff.coef=c
		fams=append(fams,ff)
	}
	res:=UP214BResult{Schema:UP214BLateRoutedSchema,Experiment:"UP-214B-late-routed-offset-transfer",SourceUP213BSeal:"e39f9e0c9cc3c453d205d0fac4f13237fa530d28",TemplatePhases:template,TrainingPhases:train,EvaluationPhases:eval,Features:[]string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin"},Router:"standardized_nearest_centroid",RidgeLambda:1e-6,EvaluationLabelFittingUsed:false,AdaptiveModelSelectionUsed:false,RetrainingUsed:false,PhaseInputUsed:false,ParityInputUsed:false,AnchorMeasurementUsedByRouted:false,AnchorMeasurementUsedByOracleFamily:false,NonlinearModelUsed:false}
	for fi,f:=range fams{
		sm:=UP214BSummary{Regime:f.name};sumS,sumR,sumF,sumA:=0.0,0.0,0.0,0.0
		for _,ph:=range eval{
			x:=up193bFeature(ph,f.canonical);z4:=[4]float64{};v:=[5]float64{1}
			for j:=0;j<4;j++{z4[j]=(x[j]-mean[j])/std[j];v[j+1]=z4[j]}
			predIdx:=0;best:=math.Inf(1)
			for gi,g:=range fams{d:=up212bDist(z4,g.centroid);if d<best{best=d;predIdx=gi}}
			res.RouteDecisions++;sm.RouteTotal++;if predIdx==fi{res.RouteCorrect++;sm.RouteCorrect++}
			routedShift:=0.0;oracleFamilyShift:=0.0
			for i:=0;i<5;i++{routedShift+=fams[predIdx].coef[i]*v[i];oracleFamilyShift+=f.coef[i]*v[i]}
			canonPt:=up191bPoint(model,UP191BProfile{Name:f.name+"_eval",SurfaceIndices:f.canonical},ph)
			anchorShift:=canonPt.RequiredMassFactor-f.means[f.canon]
			for i,p:=range f.perms{
				if i==f.canon{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:f.name+"_eval",SurfaceIndices:p},ph)
				sp:=f.means[i];rp:=sp+routedShift;fp:=sp+oracleFamilyShift;ap:=sp+anchorShift
				se:=math.Abs(sp-pt.RequiredMassFactor);re:=math.Abs(rp-pt.RequiredMassFactor);fe:=math.Abs(fp-pt.RequiredMassFactor);ae:=math.Abs(ap-pt.RequiredMassFactor)
				res.Metrics=append(res.Metrics,UP214BMetric{TrueRegime:f.name,PredictedRegime:fams[predIdx].name,Permutation:append([]int(nil),p...),Phase:ph,ActualRequiredFactor:pt.RequiredMassFactor,StaticPredictedFactor:sp,RoutedPredictedFactor:rp,OracleFamilyPredictedFactor:fp,OracleAnchorPredictedFactor:ap,StaticAbsoluteError:se,RoutedAbsoluteError:re,OracleFamilyAbsoluteError:fe,OracleAnchorAbsoluteError:ae})
				sm.EvaluationPoints++;sumS+=se;sumR+=re;sumF+=fe;sumA+=ae
				if se>sm.StaticMaxAbsoluteError{sm.StaticMaxAbsoluteError=se};if re>sm.RoutedMaxAbsoluteError{sm.RoutedMaxAbsoluteError=re};if fe>sm.OracleFamilyMaxAbsoluteError{sm.OracleFamilyMaxAbsoluteError=fe};if ae>sm.OracleAnchorMaxAbsoluteError{sm.OracleAnchorMaxAbsoluteError=ae}
				if re<se{sm.RoutedBetterThanStatic++}
			}
		}
		sm.StaticMeanAbsoluteError=sumS/float64(sm.EvaluationPoints);sm.RoutedMeanAbsoluteError=sumR/float64(sm.EvaluationPoints);sm.OracleFamilyMeanAbsoluteError=sumF/float64(sm.EvaluationPoints);sm.OracleAnchorMeanAbsoluteError=sumA/float64(sm.EvaluationPoints)
		res.Summaries=append(res.Summaries,sm)
	}
	if res.RouteDecisions>0{res.RouteAccuracy=float64(res.RouteCorrect)/float64(res.RouteDecisions)}
	return res,nil
}
