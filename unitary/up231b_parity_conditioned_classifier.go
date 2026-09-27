package unitary

import "math"

const UP231BParityClassifierSchema = "wingless.up231b-parity-conditioned-classifier.v1"

type UP231BFamilySummary struct {
	Family string `json:"family"`
	EvaluationStates int `json:"evaluation_states"`
	BaselineCorrect int `json:"baseline_correct"`
	ParityConditionedCorrect int `json:"parity_conditioned_correct"`
}

type UP231BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP230BSeal string `json:"source_up230b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Features []string `json:"features"`
	BaselineClassifier string `json:"baseline_classifier"`
	ParityClassifier string `json:"parity_classifier"`
	TrainingPoints int `json:"training_points"`
	EvaluationPoints int `json:"evaluation_points"`
	BaselineAccuracy float64 `json:"baseline_accuracy"`
	ParityConditionedAccuracy float64 `json:"parity_conditioned_accuracy"`
	EvaluationLabelFittingUsed bool `json:"evaluation_label_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	NonlinearClassifierUsed bool `json:"nonlinear_classifier_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP231BFamilySummary `json:"summaries"`
}

func RunUP231B()(UP231BResult,error){
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{121,122,123,124,125,126,127,128,129,130,131,132,133,134,135,136,137,138,139,140,141,142,143,144,145,146,147,148,149,150,151,152}
	specs:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	type row struct{name string;phase int;x [4]float64}
	rows:=[]row{}
	for _,f:=range specs{for _,ph:=range train{rows=append(rows,row{name:f.Name,phase:ph,x:up193bFeature(ph,f.Canonical)})}}
	var mean,std [4]float64
	for _,r:=range rows{for j:=0;j<4;j++{mean[j]+=r.x[j]}}
	for j:=0;j<4;j++{mean[j]/=float64(len(rows))}
	for _,r:=range rows{for j:=0;j<4;j++{d:=r.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(rows)));if std[j]==0{std[j]=1}}
	type accum struct{sum [4]float64;n int}
	base:=map[string]accum{}
	par:=map[string]accum{}
	for _,r:=range rows{
		z:=[4]float64{}
		for j:=0;j<4;j++{z[j]=(r.x[j]-mean[j])/std[j]}
		a:=base[r.name];for j:=0;j<4;j++{a.sum[j]+=z[j]};a.n++;base[r.name]=a
		k:=r.name
		if r.phase%2==0{k+="|even"}else{k+="|odd"}
		p:=par[k];for j:=0;j<4;j++{p.sum[j]+=z[j]};p.n++;par[k]=p
	}
	centroid:=func(a accum)(c [4]float64){for j:=0;j<4;j++{c[j]=a.sum[j]/float64(a.n)};return}
	baseCent:=map[string][4]float64{}
	parCent:=map[string][4]float64{}
	for k,a:=range base{baseCent[k]=centroid(a)}
	for k,a:=range par{parCent[k]=centroid(a)}
	res:=UP231BResult{
		Schema:UP231BParityClassifierSchema,
		Experiment:"UP-231B-parity-conditioned-classifier",
		SourceUP230BSeal:"8242983e1aac1449ed70ca17726c5a5c05b2ff1d",
		TrainingPhases:train,EvaluationPhases:eval,
		Features:[]string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin"},
		BaselineClassifier:"standardized_nearest_centroid",
		ParityClassifier:"standardized_nearest_centroid_by_fixed_phase_parity",
		TrainingPoints:len(rows),
		EvaluationLabelFittingUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		ParityInputUsed:true,
		NonlinearClassifierUsed:false,
		LiveActivation:false,
	}
	baseTotal,parTotal:=0,0
	for _,f:=range specs{
		s:=UP231BFamilySummary{Family:f.Name}
		for _,ph:=range eval{
			x:=up193bFeature(ph,f.Canonical);z:=[4]float64{}
			for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]}
			bestBase:="";bestBaseD:=math.Inf(1)
			bestPar:="";bestParD:=math.Inf(1)
			parity:="|odd";if ph%2==0{parity="|even"}
			for _,g:=range specs{
				if d:=up212bDist(z,baseCent[g.Name]);d<bestBaseD{bestBaseD=d;bestBase=g.Name}
				if d:=up212bDist(z,parCent[g.Name+parity]);d<bestParD{bestParD=d;bestPar=g.Name}
			}
			s.EvaluationStates++;res.EvaluationPoints++
			if bestBase==f.Name{s.BaselineCorrect++;baseTotal++}
			if bestPar==f.Name{s.ParityConditionedCorrect++;parTotal++}
		}
		res.Summaries=append(res.Summaries,s)
	}
	if res.EvaluationPoints>0{
		res.BaselineAccuracy=float64(baseTotal)/float64(res.EvaluationPoints)
		res.ParityConditionedAccuracy=float64(parTotal)/float64(res.EvaluationPoints)
	}
	return res,nil
}
