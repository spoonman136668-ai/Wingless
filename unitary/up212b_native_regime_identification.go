package unitary

import "math"

const UP212BRegimeSchema="wingless.up212b-native-regime-identification.v1"

type UP212BConfusion struct{
	Actual string `json:"actual"`
	Predicted string `json:"predicted"`
	Count int `json:"count"`
}
type UP212BRegimeSummary struct{
	Regime string `json:"regime"`
	Correct int `json:"correct"`
	Total int `json:"total"`
	Accuracy float64 `json:"accuracy"`
}
type UP212BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP211BSeal string `json:"source_up211b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Features []string `json:"features"`
	Classifier string `json:"classifier"`
	TrainingPoints int `json:"training_points"`
	EvaluationPoints int `json:"evaluation_points"`
	ChanceAccuracy float64 `json:"chance_accuracy"`
	Accuracy float64 `json:"accuracy"`
	MeanNearestDistance float64 `json:"mean_nearest_distance"`
	MeanConfidenceMargin float64 `json:"mean_confidence_margin"`
	EvaluationLabelFittingUsed bool `json:"evaluation_label_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	CalibrationTargetInputUsed bool `json:"calibration_target_input_used"`
	NonlinearClassifierUsed bool `json:"nonlinear_classifier_used"`
	Regimes []UP212BRegimeSummary `json:"regimes"`
	Confusion []UP212BConfusion `json:"confusion"`
}
func up212bDist(a,b [4]float64)float64{
	s:=0.0;for i:=0;i<4;i++{d:=a[i]-b[i];s+=d*d};return math.Sqrt(s)
}
func RunUP212B()(UP212BResult,error){
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{121,122,123,124,125,126,127,128,129,130,131,132}
	specs:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	type row struct{name string;x [4]float64}
	rows:=[]row{}
	for _,f:=range specs{for _,ph:=range train{rows=append(rows,row{name:f.Name,x:up193bFeature(ph,f.Canonical)})}}
	var mean,std [4]float64
	for _,r:=range rows{for j:=0;j<4;j++{mean[j]+=r.x[j]}}
	for j:=0;j<4;j++{mean[j]/=float64(len(rows))}
	for _,r:=range rows{for j:=0;j<4;j++{d:=r.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(rows)));if std[j]==0{std[j]=1}}
	cent:=map[string][4]float64{};counts:=map[string]int{}
	for _,r:=range rows{
		z:=[4]float64{}
		for j:=0;j<4;j++{z[j]=(r.x[j]-mean[j])/std[j]}
		c:=cent[r.name];for j:=0;j<4;j++{c[j]+=z[j]};cent[r.name]=c;counts[r.name]++
	}
	for name,c:=range cent{for j:=0;j<4;j++{c[j]/=float64(counts[name])};cent[name]=c}
	res:=UP212BResult{
		Schema:UP212BRegimeSchema,Experiment:"UP-212B-native-regime-identification",SourceUP211BSeal:"6a64c22843a6cd663e0e56f849fc5dfd4af9423e",
		TrainingPhases:train,EvaluationPhases:eval,
		Features:[]string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin"},
		Classifier:"standardized_nearest_centroid",TrainingPoints:len(rows),ChanceAccuracy:0.25,
		EvaluationLabelFittingUsed:false,AdaptiveFeatureSelectionUsed:false,PhaseInputUsed:false,ParityInputUsed:false,CalibrationTargetInputUsed:false,NonlinearClassifierUsed:false,
	}
	correctBy:=map[string]int{};totalBy:=map[string]int{};conf:=map[string]int{}
	sumNearest,sumMargin:=0.0,0.0
	for _,f:=range specs{for _,ph:=range eval{
		x:=up193bFeature(ph,f.Canonical);z:=[4]float64{}
		for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]}
		bestName:="";best,second:=math.Inf(1),math.Inf(1)
		for _,g:=range specs{
			d:=up212bDist(z,cent[g.Name])
			if d<best{second=best;best=d;bestName=g.Name}else if d<second{second=d}
		}
		res.EvaluationPoints++;totalBy[f.Name]++;sumNearest+=best;sumMargin+=second-best
		if bestName==f.Name{correctBy[f.Name]++}
		conf[f.Name+"|"+bestName]++
	}}
	totalCorrect:=0
	for _,f:=range specs{
		c:=correctBy[f.Name];t:=totalBy[f.Name];totalCorrect+=c
		s:=UP212BRegimeSummary{Regime:f.Name,Correct:c,Total:t}
		if t>0{s.Accuracy=float64(c)/float64(t)}
		res.Regimes=append(res.Regimes,s)
		for _,g:=range specs{res.Confusion=append(res.Confusion,UP212BConfusion{Actual:f.Name,Predicted:g.Name,Count:conf[f.Name+"|"+g.Name]})}
	}
	if res.EvaluationPoints>0{
		res.Accuracy=float64(totalCorrect)/float64(res.EvaluationPoints)
		res.MeanNearestDistance=sumNearest/float64(res.EvaluationPoints)
		res.MeanConfidenceMargin=sumMargin/float64(res.EvaluationPoints)
	}
	return res,nil
}
