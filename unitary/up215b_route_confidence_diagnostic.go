package unitary

import "math"

const UP215BRouteConfidenceSchema="wingless.up215b-route-confidence-diagnostic.v1"

type UP215BPoint struct{
	TrueRegime string `json:"true_regime"`
	PredictedRegime string `json:"predicted_regime"`
	Phase int `json:"phase"`
	Correct bool `json:"correct"`
	NearestDistance float64 `json:"nearest_distance"`
	SecondNearestDistance float64 `json:"second_nearest_distance"`
	ConfidenceMargin float64 `json:"confidence_margin"`
}
type UP215BConfusion struct{
	TrueRegime string `json:"true_regime"`
	PredictedRegime string `json:"predicted_regime"`
	Count int `json:"count"`
}
type UP215BSummary struct{
	Regime string `json:"regime"`
	Correct int `json:"correct"`
	Total int `json:"total"`
	MeanNearestCorrect float64 `json:"mean_nearest_correct"`
	MeanNearestIncorrect float64 `json:"mean_nearest_incorrect"`
	MeanMarginCorrect float64 `json:"mean_margin_correct"`
	MeanMarginIncorrect float64 `json:"mean_margin_incorrect"`
	MinMargin float64 `json:"min_margin"`
	MaxMargin float64 `json:"max_margin"`
}
type UP215BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP214BSeal string `json:"source_up214b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Features []string `json:"features"`
	Router string `json:"router"`
	RouteDecisions int `json:"route_decisions"`
	RouteCorrect int `json:"route_correct"`
	RouteAccuracy float64 `json:"route_accuracy"`
	EvaluationLabelFittingUsed bool `json:"evaluation_label_fitting_used"`
	RetrainingUsed bool `json:"retraining_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	ConfidenceThresholdUsed bool `json:"confidence_threshold_used"`
	FallbackRouteUsed bool `json:"fallback_route_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	NonlinearClassifierUsed bool `json:"nonlinear_classifier_used"`
	Points []UP215BPoint `json:"points"`
	Confusion []UP215BConfusion `json:"confusion"`
	Summaries []UP215BSummary `json:"summaries"`
}
func RunUP215B()(UP215BResult,error){
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{160,161,162,163,164,165,166,167,168,169,170,171,172,173,174,175,176,177,178,179,180,181,182,183}
	specs:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	type raw struct{name string;x [4]float64}
	raws:=[]raw{}
	for _,f:=range specs{for _,ph:=range train{raws=append(raws,raw{name:f.Name,x:up193bFeature(ph,f.Canonical)})}}
	var mean,std [4]float64
	for _,r:=range raws{for j:=0;j<4;j++{mean[j]+=r.x[j]}}
	for j:=0;j<4;j++{mean[j]/=float64(len(raws))}
	for _,r:=range raws{for j:=0;j<4;j++{d:=r.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(raws)));if std[j]==0{std[j]=1}}
	type fam struct{name string;canonical []int;centroid [4]float64}
	fams:=[]fam{}
	for _,f:=range specs{
		ff:=fam{name:f.Name,canonical:f.Canonical};n:=0
		for _,r:=range raws{
			if r.name!=f.Name{continue}
			for j:=0;j<4;j++{ff.centroid[j]+=(r.x[j]-mean[j])/std[j]}
			n++
		}
		for j:=0;j<4;j++{ff.centroid[j]/=float64(n)}
		fams=append(fams,ff)
	}
	res:=UP215BResult{
		Schema:UP215BRouteConfidenceSchema,Experiment:"UP-215B-route-confidence-diagnostic",SourceUP214BSeal:"3664d805bbb1a9acf70228edb8286a8df9336838",
		TrainingPhases:train,EvaluationPhases:eval,Features:[]string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin"},
		Router:"standardized_nearest_centroid",
		EvaluationLabelFittingUsed:false,RetrainingUsed:false,AdaptiveFeatureSelectionUsed:false,ConfidenceThresholdUsed:false,FallbackRouteUsed:false,PhaseInputUsed:false,ParityInputUsed:false,NonlinearClassifierUsed:false,
	}
	conf:=map[string]int{}
	for fi,f:=range fams{
		sm:=UP215BSummary{Regime:f.name,MinMargin:math.Inf(1)}
		sumNC,sumNI,sumMC,sumMI:=0.0,0.0,0.0,0.0;nC,nI:=0,0
		for _,ph:=range eval{
			x:=up193bFeature(ph,f.canonical);z:=[4]float64{}
			for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]}
			bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
			for gi,g:=range fams{
				d:=up212bDist(z,g.centroid)
				if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}
			}
			margin:=second-best;correct:=bestIdx==fi
			res.RouteDecisions++;sm.Total++
			if correct{res.RouteCorrect++;sm.Correct++;sumNC+=best;sumMC+=margin;nC++}else{sumNI+=best;sumMI+=margin;nI++}
			if margin<sm.MinMargin{sm.MinMargin=margin};if margin>sm.MaxMargin{sm.MaxMargin=margin}
			p:=UP215BPoint{TrueRegime:f.name,PredictedRegime:fams[bestIdx].name,Phase:ph,Correct:correct,NearestDistance:best,SecondNearestDistance:second,ConfidenceMargin:margin}
			res.Points=append(res.Points,p)
			conf[f.name+"|"+fams[bestIdx].name]++
		}
		if nC>0{sm.MeanNearestCorrect=sumNC/float64(nC);sm.MeanMarginCorrect=sumMC/float64(nC)}
		if nI>0{sm.MeanNearestIncorrect=sumNI/float64(nI);sm.MeanMarginIncorrect=sumMI/float64(nI)}
		if math.IsInf(sm.MinMargin,1){sm.MinMargin=0}
		res.Summaries=append(res.Summaries,sm)
	}
	for _,t:=range fams{for _,p:=range fams{res.Confusion=append(res.Confusion,UP215BConfusion{TrueRegime:t.name,PredictedRegime:p.name,Count:conf[t.name+"|"+p.name]})}}
	if res.RouteDecisions>0{res.RouteAccuracy=float64(res.RouteCorrect)/float64(res.RouteDecisions)}
	return res,nil
}
