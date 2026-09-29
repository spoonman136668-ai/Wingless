package unitary

import "math"

const UP235BFeatureSubsetSchema = "wingless.up235b-centered-parity-feature-subsets.v1"

type UP235BFamilySummary struct {
	Family  string `json:"family"`
	States  int    `json:"states"`
	Correct int    `json:"correct"`
}

type UP235BSubsetSummary struct {
	Mask        int                  `json:"mask"`
	Features    []string             `json:"features"`
	Cardinality int                  `json:"cardinality"`
	Accuracy    float64              `json:"accuracy"`
	Families    []UP235BFamilySummary `json:"families"`
}

type UP235BResult struct {
	Schema                       string                `json:"schema"`
	Experiment                   string                `json:"experiment"`
	SourceUP234BSeal             string                `json:"source_up234b_seal"`
	ParentAccuracy               float64               `json:"parent_family_centered_pooled_accuracy"`
	FeatureNames                 []string              `json:"feature_names"`
	SubsetCount                  int                   `json:"subset_count"`
	MinimumPerfectCardinality    int                   `json:"minimum_perfect_cardinality"`
	MinimumPerfectMasks          []int                 `json:"minimum_perfect_masks"`
	TrainingParityLabelsUsed     bool                  `json:"training_parity_labels_used"`
	FamilyCenteringUsed          bool                  `json:"family_centering_used"`
	PhaseInputAtInferenceUsed    bool                  `json:"phase_input_at_inference_used"`
	EvaluationLabelFittingUsed   bool                  `json:"evaluation_label_fitting_used"`
	AdaptiveFeatureSelectionUsed bool                  `json:"adaptive_feature_selection_used"`
	NonlinearClassifierUsed      bool                  `json:"nonlinear_classifier_used"`
	LiveActivation               bool                  `json:"live_activation"`
	Subsets                      []UP235BSubsetSummary `json:"subsets"`
}

func RunUP235B() (UP235BResult, error) {
	parent, err := RunUP234B()
	if err != nil {
		return UP235BResult{}, err
	}
	train := []int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval := []int{121,122,123,124,125,126,127,128,129,130,131,132,133,134,135,136,137,138,139,140,141,142,143,144,145,146,147,148,149,150,151,152}
	specs := []UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	featureNames := []string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin"}
	type row struct{name string; phase int; x [4]float64}
	rows:=[]row{}
	familySum:=map[string][4]float64{}
	familyN:=map[string]int{}
	for _,f:=range specs{
		for _,ph:=range train{
			x:=up193bFeature(ph,f.Canonical)
			rows=append(rows,row{f.Name,ph,x})
			a:=familySum[f.Name]
			for j:=0;j<4;j++{a[j]+=x[j]}
			familySum[f.Name]=a
			familyN[f.Name]++
		}
	}
	familyMean:=map[string][4]float64{}
	for _,f:=range specs{
		a:=familySum[f.Name]
		for j:=0;j<4;j++{a[j]/=float64(familyN[f.Name])}
		familyMean[f.Name]=a
	}
	parity:=func(ph int)string{if ph%2==0{return "even"};return "odd"}
	res:=UP235BResult{
		Schema:UP235BFeatureSubsetSchema,
		Experiment:"UP-235B-centered-parity-feature-subsets",
		SourceUP234BSeal:"cbc29927d9aa9e1d008e053578728baba44bdeaa",
		ParentAccuracy:parent.FamilyCenteredPooledAccuracy,
		FeatureNames:featureNames,
		MinimumPerfectCardinality:99,
		TrainingParityLabelsUsed:true,
		FamilyCenteringUsed:true,
		PhaseInputAtInferenceUsed:false,
		EvaluationLabelFittingUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		NonlinearClassifierUsed:false,
		LiveActivation:false,
	}
	for mask:=1;mask<16;mask++{
		selected:=[]int{}
		features:=[]string{}
		for j:=0;j<4;j++{if mask&(1<<j)!=0{selected=append(selected,j);features=append(features,featureNames[j])}}
		std:=make([]float64,len(selected))
		for _,r:=range rows{
			m:=familyMean[r.name]
			for q,j:=range selected{d:=r.x[j]-m[j];std[q]+=d*d}
		}
		for q:=range std{std[q]=math.Sqrt(std[q]/float64(len(rows)));if std[q]==0{std[q]=1}}
		type acc struct{sum []float64;n int}
		pooled:=map[string]acc{"even":{sum:make([]float64,len(selected))},"odd":{sum:make([]float64,len(selected))}}
		for _,r:=range rows{
			p:=parity(r.phase);a:=pooled[p];m:=familyMean[r.name]
			for q,j:=range selected{a.sum[q]+=(r.x[j]-m[j])/std[q]}
			a.n++;pooled[p]=a
		}
		cent:=map[string][]float64{"even":make([]float64,len(selected)),"odd":make([]float64,len(selected))}
		for _,p:=range []string{"even","odd"}{a:=pooled[p];for q:=range selected{cent[p][q]=a.sum[q]/float64(a.n)}}
		sm:=UP235BSubsetSummary{Mask:mask,Features:features,Cardinality:len(selected)}
		correct:=0
		for _,f:=range specs{
			fs:=UP235BFamilySummary{Family:f.Name}
			m:=familyMean[f.Name]
			for _,ph:=range eval{
				x:=up193bFeature(ph,f.Canonical)
				best:="";bestD:=math.Inf(1)
				for _,p:=range []string{"even","odd"}{
					d:=0.0
					for q,j:=range selected{z:=(x[j]-m[j])/std[q];e:=z-cent[p][q];d+=e*e}
					if d<bestD{bestD=d;best=p}
				}
				fs.States++
				if best==parity(ph){fs.Correct++;correct++}
			}
			sm.Families=append(sm.Families,fs)
		}
		sm.Accuracy=float64(correct)/128.0
		if sm.Accuracy==1.0{
			if sm.Cardinality<res.MinimumPerfectCardinality{res.MinimumPerfectCardinality=sm.Cardinality;res.MinimumPerfectMasks=[]int{mask}
			}else if sm.Cardinality==res.MinimumPerfectCardinality{res.MinimumPerfectMasks=append(res.MinimumPerfectMasks,mask)}
		}
		res.Subsets=append(res.Subsets,sm)
	}
	res.SubsetCount=len(res.Subsets)
	if res.MinimumPerfectCardinality==99{res.MinimumPerfectCardinality=0}
	return res,nil
}
