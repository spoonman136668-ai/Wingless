package unitary

import "math"

const UP239BHeldoutCompositionParitySchema = "wingless.up239b-heldout-composition-parity.v1"

type UP239BFamilySummary struct {
	Family         string  `json:"family"`
	FarStates      int     `json:"far_states"`
	FarCorrect     int     `json:"far_correct"`
	FarAccuracy    float64 `json:"far_accuracy"`
	ExtremeStates  int     `json:"extreme_states"`
	ExtremeCorrect int     `json:"extreme_correct"`
	ExtremeAccuracy float64 `json:"extreme_accuracy"`
}

type UP239BResult struct {
	Schema                       string                `json:"schema"`
	Experiment                   string                `json:"experiment"`
	SourceUP238BSeal             string                `json:"source_up238b_seal"`
	ParentClassification         string                `json:"parent_classification"`
	ParentExtremeAccuracy        float64               `json:"parent_extreme_accuracy"`
	TrainingPhases               []int                 `json:"training_phases"`
	FarEvaluationPhases          []int                 `json:"far_evaluation_phases"`
	ExtremeEvaluationPhases      []int                 `json:"extreme_evaluation_phases"`
	Feature                      string                `json:"feature"`
	HeldoutFamilies              []UP197BComposition   `json:"heldout_families"`
	UnlabeledFamilyCenteringUsed bool                  `json:"unlabeled_family_centering_used"`
	EvaluationLabelFittingUsed   bool                  `json:"evaluation_label_fitting_used"`
	PhaseInputAtInferenceUsed    bool                  `json:"phase_input_at_inference_used"`
	AdaptiveFeatureSelectionUsed bool                  `json:"adaptive_feature_selection_used"`
	NonlinearClassifierUsed      bool                  `json:"nonlinear_classifier_used"`
	LiveActivation               bool                  `json:"live_activation"`
	FarAccuracy                  float64               `json:"far_accuracy"`
	ExtremeAccuracy              float64               `json:"extreme_accuracy"`
	Summaries                    []UP239BFamilySummary `json:"summaries"`
	Classification               string                `json:"classification"`
}

func RunUP239B() (UP239BResult, error) {
	parent, err := RunUP238B()
	if err != nil {
		return UP239BResult{}, err
	}
	train := []int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	far := []int{185,186,187,188,189,190,191,192,193,194,195,196,197,198,199,200,201,202,203,204,205,206,207,208,209,210,211,212,213,214,215,216}
	extreme := []int{249,250,251,252,253,254,255,256,257,258,259,260,261,262,263,264,265,266,267,268,269,270,271,272,273,274,275,276,277,278,279,280}
	trainFamilies := []UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	heldout := []UP197BComposition{
		{Name:"store3",Canonical:[]int{0,1,2}},
		{Name:"observe3",Canonical:[]int{5,6,7}},
		{Name:"mixed3",Canonical:[]int{0,5,1}},
		{Name:"report2",Canonical:[]int{13,14}},
	}
	feat := func(ph int, f UP197BComposition) float64 { return up193bFeature(ph, f.Canonical)[2] }
	parity := func(ph int) string { if ph%2==0 { return "even" }; return "odd" }

	trainMean := map[string]float64{}
	trainN := map[string]int{}
	type sample struct{name string; ph int; x float64}
	var samples []sample
	for _,f := range trainFamilies {
		for _,ph := range train {
			x:=feat(ph,f)
			samples=append(samples,sample{name:f.Name,ph:ph,x:x})
			trainMean[f.Name]+=x
			trainN[f.Name]++
		}
	}
	for _,f:=range trainFamilies { trainMean[f.Name]/=float64(trainN[f.Name]) }
	std:=0.0
	for _,s:=range samples { d:=s.x-trainMean[s.name]; std+=d*d }
	std=math.Sqrt(std/float64(len(samples))); if std==0 { std=1 }
	sum:=map[string]float64{"even":0,"odd":0}; n:=map[string]int{"even":0,"odd":0}
	for _,s:=range samples { q:=parity(s.ph); sum[q]+=(s.x-trainMean[s.name])/std; n[q]++ }
	cent:=map[string]float64{"even":sum["even"]/float64(n["even"]),"odd":sum["odd"]/float64(n["odd"])}

	heldoutMean:=map[string]float64{}
	for _,f:=range heldout {
		for _,ph:=range train { heldoutMean[f.Name]+=feat(ph,f) }
		heldoutMean[f.Name]/=float64(len(train))
	}

	r:=UP239BResult{
		Schema:UP239BHeldoutCompositionParitySchema,
		Experiment:"UP-239B-heldout-composition-parity",
		SourceUP238BSeal:"6ade177b8895ce06170b70330c4ef79b4becad59",
		ParentClassification:parent.Classification,
		ParentExtremeAccuracy:parent.ExtremeAccuracy,
		TrainingPhases:train,
		FarEvaluationPhases:far,
		ExtremeEvaluationPhases:extreme,
		Feature:"near_zero_margin_count",
		HeldoutFamilies:heldout,
		UnlabeledFamilyCenteringUsed:true,
	}
	farCorrect,farTotal,extCorrect,extTotal:=0,0,0,0
	classify:=func(ph int,f UP197BComposition) string {
		z:=(feat(ph,f)-heldoutMean[f.Name])/std
		p:="even"
		if math.Abs(z-cent["odd"])<math.Abs(z-cent["even"]) { p="odd" }
		return p
	}
	for _,f:=range heldout {
		sm:=UP239BFamilySummary{Family:f.Name}
		for _,ph:=range far {
			sm.FarStates++; farTotal++
			if classify(ph,f)==parity(ph) { sm.FarCorrect++; farCorrect++ }
		}
		for _,ph:=range extreme {
			sm.ExtremeStates++; extTotal++
			if classify(ph,f)==parity(ph) { sm.ExtremeCorrect++; extCorrect++ }
		}
		sm.FarAccuracy=float64(sm.FarCorrect)/float64(sm.FarStates)
		sm.ExtremeAccuracy=float64(sm.ExtremeCorrect)/float64(sm.ExtremeStates)
		r.Summaries=append(r.Summaries,sm)
	}
	r.FarAccuracy=float64(farCorrect)/float64(farTotal)
	r.ExtremeAccuracy=float64(extCorrect)/float64(extTotal)
	if parent.Classification!="EXTREME_TRANSFER_PERFECT" || parent.ExtremeAccuracy!=1.0 {
		r.Classification="ANCHOR_NOT_REPRODUCED"
	} else if r.FarAccuracy==1.0 && r.ExtremeAccuracy==1.0 {
		r.Classification="HELDOUT_COMPOSITION_TRANSFER_PERFECT"
	} else {
		r.Classification="HELDOUT_COMPOSITION_TRANSFER_DECAY"
	}
	return r,nil
}
