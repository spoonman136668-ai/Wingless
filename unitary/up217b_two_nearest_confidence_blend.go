package unitary

import "math"

const UP217BBlendSchema="wingless.up217b-two-nearest-confidence-blend.v1"

type UP217BMetric struct{
	TrueRegime string `json:"true_regime"`
	PredictedRegime string `json:"predicted_regime"`
	SecondRegime string `json:"second_regime"`
	Phase int `json:"phase"`
	Permutation []int `json:"permutation"`
	ConfidenceMargin float64 `json:"confidence_margin"`
	BlendUsed bool `json:"blend_used"`
	StaticAbsoluteError float64 `json:"static_absolute_error"`
	RoutedAbsoluteError float64 `json:"routed_absolute_error"`
	BlendAbsoluteError float64 `json:"blend_absolute_error"`
	OracleFamilyAbsoluteError float64 `json:"oracle_family_absolute_error"`
	OracleAnchorAbsoluteError float64 `json:"oracle_anchor_absolute_error"`
}
type UP217BSummary struct{
	Regime string `json:"regime"`
	RouteCorrect int `json:"route_correct"`
	RouteTotal int `json:"route_total"`
	BlendDecisions int `json:"blend_decisions"`
	EvaluationPoints int `json:"evaluation_points"`
	StaticMeanAbsoluteError float64 `json:"static_mean_absolute_error"`
	RoutedMeanAbsoluteError float64 `json:"routed_mean_absolute_error"`
	BlendMeanAbsoluteError float64 `json:"blend_mean_absolute_error"`
	OracleFamilyMeanAbsoluteError float64 `json:"oracle_family_mean_absolute_error"`
	OracleAnchorMeanAbsoluteError float64 `json:"oracle_anchor_mean_absolute_error"`
	BlendBetterThanRouted int `json:"blend_better_than_routed"`
	BlendBetterThanStatic int `json:"blend_better_than_static"`
}
type UP217BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP216BSeal string `json:"source_up216b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	ThresholdRule string `json:"threshold_rule"`
	BlendRule string `json:"blend_rule"`
	TrainingDerivedThreshold float64 `json:"training_derived_threshold"`
	RouteDecisions int `json:"route_decisions"`
	RouteCorrect int `json:"route_correct"`
	RouteAccuracy float64 `json:"route_accuracy"`
	BlendDecisions int `json:"blend_decisions"`
	EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`
	AdaptiveWeightingUsed bool `json:"adaptive_weighting_used"`
	RetrainingUsed bool `json:"retraining_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	NonlinearModelUsed bool `json:"nonlinear_model_used"`
	Metrics []UP217BMetric `json:"metrics"`
	Summaries []UP217BSummary `json:"summaries"`
}
func RunUP217B()(UP217BResult,error){
	template:=[]int{55,56,57};train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{196,197,198,199,200,201,202,203,204,205,206,207}
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
		ff:=up213bFamily{name:f.Name,canonical:f.Canonical,perms:perms,means:means,canon:canon};n:=0
		for _,r:=range raws{if r.name==f.Name{for j:=0;j<4;j++{ff.centroid[j]+=(r.x[j]-mean[j])/std[j]};n++}}
		for j:=0;j<4;j++{ff.centroid[j]/=float64(n)}
		var a [5][5]float64;var b [5]float64
		for _,ph:=range train{
			x:=up193bFeature(ph,f.Canonical);v:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
			pt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train",SurfaceIndices:f.Canonical},ph);y:=pt.RequiredMassFactor-means[canon]
			for i:=0;i<5;i++{b[i]+=v[i]*y;for j:=0;j<5;j++{a[i][j]+=v[i]*v[j]}}
		}
		for i:=1;i<5;i++{a[i][i]+=1e-6}
		c,ok:=up193bSolve(a,b);if !ok{return UP217BResult{},nil};ff.coef=c;fams=append(fams,ff)
	}
	threshold:=math.Inf(1)
	for fi,f:=range fams{for _,ph:=range train{
		x:=up193bFeature(ph,f.canonical);z:=[4]float64{}
		for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]}
		bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
		for gi,g:=range fams{
			d:=up212bDist(z,g.centroid)
			if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}
		}
		if bestIdx==fi{m:=second-best;if m<threshold{threshold=m}}
	}}
	if math.IsInf(threshold,1){threshold=0}
	res:=UP217BResult{
		Schema:UP217BBlendSchema,Experiment:"UP-217B-two-nearest-confidence-blend",SourceUP216BSeal:"1de95869087acfed0f23858a0108c456332c8c67",
		TrainingPhases:train,EvaluationPhases:eval,ThresholdRule:"minimum_correct_training_confidence_margin",BlendRule:"equal_mean_of_two_nearest_frozen_calibrators",
		TrainingDerivedThreshold:threshold,EvaluationDerivedThresholdUsed:false,AdaptiveWeightingUsed:false,RetrainingUsed:false,PhaseInputUsed:false,ParityInputUsed:false,NonlinearModelUsed:false,
	}
	for fi,f:=range fams{
		sm:=UP217BSummary{Regime:f.name};sumS,sumR,sumB,sumF,sumA:=0.0,0.0,0.0,0.0,0.0
		for _,ph:=range eval{
			x:=up193bFeature(ph,f.canonical);z4:=[4]float64{};v:=[5]float64{1}
			for j:=0;j<4;j++{z4[j]=(x[j]-mean[j])/std[j];v[j+1]=z4[j]}
			bestIdx,secondIdx:=-1,-1;best,second:=math.Inf(1),math.Inf(1)
			for gi,g:=range fams{
				d:=up212bDist(z4,g.centroid)
				if d<best{second=best;secondIdx=bestIdx;best=d;bestIdx=gi}else if d<second{second=d;secondIdx=gi}
			}
			margin:=second-best;blend:=margin<threshold
			res.RouteDecisions++;sm.RouteTotal++;if bestIdx==fi{res.RouteCorrect++;sm.RouteCorrect++}
			if blend{res.BlendDecisions++;sm.BlendDecisions++}
			routedShift,secondShift,oracleFamilyShift:=0.0,0.0,0.0
			for i:=0;i<5;i++{
				routedShift+=fams[bestIdx].coef[i]*v[i]
				secondShift+=fams[secondIdx].coef[i]*v[i]
				oracleFamilyShift+=f.coef[i]*v[i]
			}
			blendShift:=routedShift;if blend{blendShift=(routedShift+secondShift)/2}
			canonPt:=up191bPoint(model,UP191BProfile{Name:f.name+"_eval",SurfaceIndices:f.canonical},ph);anchorShift:=canonPt.RequiredMassFactor-f.means[f.canon]
			for i,p:=range f.perms{
				if i==f.canon{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:f.name+"_eval",SurfaceIndices:p},ph)
				sp:=f.means[i];rp:=sp+routedShift;bp:=sp+blendShift;fp:=sp+oracleFamilyShift;ap:=sp+anchorShift
				se:=math.Abs(sp-pt.RequiredMassFactor);re:=math.Abs(rp-pt.RequiredMassFactor);be:=math.Abs(bp-pt.RequiredMassFactor);fe:=math.Abs(fp-pt.RequiredMassFactor);ae:=math.Abs(ap-pt.RequiredMassFactor)
				res.Metrics=append(res.Metrics,UP217BMetric{TrueRegime:f.name,PredictedRegime:fams[bestIdx].name,SecondRegime:fams[secondIdx].name,Phase:ph,Permutation:append([]int(nil),p...),ConfidenceMargin:margin,BlendUsed:blend,StaticAbsoluteError:se,RoutedAbsoluteError:re,BlendAbsoluteError:be,OracleFamilyAbsoluteError:fe,OracleAnchorAbsoluteError:ae})
				sm.EvaluationPoints++;sumS+=se;sumR+=re;sumB+=be;sumF+=fe;sumA+=ae
				if be<re{sm.BlendBetterThanRouted++};if be<se{sm.BlendBetterThanStatic++}
			}
		}
		n:=float64(sm.EvaluationPoints);sm.StaticMeanAbsoluteError=sumS/n;sm.RoutedMeanAbsoluteError=sumR/n;sm.BlendMeanAbsoluteError=sumB/n;sm.OracleFamilyMeanAbsoluteError=sumF/n;sm.OracleAnchorMeanAbsoluteError=sumA/n
		res.Summaries=append(res.Summaries,sm)
	}
	if res.RouteDecisions>0{res.RouteAccuracy=float64(res.RouteCorrect)/float64(res.RouteDecisions)}
	return res,nil
}
