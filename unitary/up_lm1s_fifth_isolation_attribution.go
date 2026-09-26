package unitary

const UPLM1SIsolationSchema="wingless.up-lm1s-fifth-isolation-attribution.v1"

type UPLM1SFamilyMetric struct {
	Arm string `json:"arm"`
	Family string `json:"family"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity float64 `json:"perplexity"`
}
type UPLM1SSummary struct {
	Arm string `json:"arm"`
	MinHeldoutAccuracy float64 `json:"min_heldout_accuracy"`
	MeanHeldoutAccuracy float64 `json:"mean_heldout_accuracy"`
	HeldoutAccuracySpread float64 `json:"heldout_accuracy_spread"`
	FifthHeldoutAccuracy float64 `json:"fifth_heldout_accuracy"`
	FifthHeldoutPerplexity float64 `json:"fifth_heldout_perplexity"`
	PriorMeanHeldoutAccuracy float64 `json:"prior_mean_heldout_accuracy"`
}
type UPLM1SIsolationResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1RSeal string `json:"source_up_lm1r_seal"`
	StateDimension int `json:"state_dimension"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	LearningRate float64 `json:"learning_rate"`
	FifthUpdatesPerEpoch int `json:"fifth_updates_per_epoch"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RouterUsed bool `json:"router_used"`
	ExactRecallUsed bool `json:"exact_recall_used"`
	AdaptiveScheduleUsed bool `json:"adaptive_schedule_used"`
	SixthFamilyUsed bool `json:"sixth_family_used"`
	Metrics []UPLM1SFamilyMetric `json:"metrics"`
	Summaries []UPLM1SSummary `json:"summaries"`
	FifthOnlyMinusCyclicFifthAccuracy float64 `json:"fifth_only_minus_cyclic_fifth_accuracy"`
	FifthOnlyMinusCyclicPriorMeanAccuracy float64 `json:"fifth_only_minus_cyclic_prior_mean_accuracy"`
	FifthOnlyGainOverNoAdaptation float64 `json:"fifth_only_gain_over_no_adaptation"`
}

func uplm1sTrain(model *uplm0aModel,arm string,families [5][]uplm0fExample){
	for epoch:=0;epoch<4;epoch++{
		for i:=range families[0]{
			switch arm{
			case "fifth_only":
				model.trainSentence(families[4][i].text,0.08)
			case "five_family_cyclic":
				for f:=0;f<5;f++{model.trainSentence(families[f][i].text,0.08)}
			}
		}
	}
}
func uplm1sSummary(arm string,metrics []UPLM1SFamilyMetric) UPLM1SSummary{
	min,max,sum:=1.0,-1.0,0.0
	prior:=0.0
	count:=0
	fifthAcc,fifthPpl:=0.0,0.0
	for _,m:=range metrics{
		if m.Arm!=arm{continue}
		count++;sum+=m.Top1Accuracy
		if m.Top1Accuracy<min{min=m.Top1Accuracy}
		if m.Top1Accuracy>max{max=m.Top1Accuracy}
		if m.Family=="fifth"{fifthAcc=m.Top1Accuracy;fifthPpl=m.Perplexity}else{prior+=m.Top1Accuracy}
	}
	return UPLM1SSummary{
		Arm:arm,MinHeldoutAccuracy:min,MeanHeldoutAccuracy:sum/float64(count),HeldoutAccuracySpread:max-min,
		FifthHeldoutAccuracy:fifthAcc,FifthHeldoutPerplexity:fifthPpl,PriorMeanHeldoutAccuracy:prior/4,
	}
}
func RunUPLM1S()(UPLM1SIsolationResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil{return UPLM1SIsolationResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	held:=[5][]uplm0fExample{baseHeld,paraHeld,thirdHeld,fourthHeld,fifthHeld}
	res:=UPLM1SIsolationResult{
		Schema:UPLM1SIsolationSchema,Experiment:"UP-LM1S-fifth-isolation-attribution",
		SourceUPLM1RSeal:"d11cd78832db0685d7e5485103863266b4792da9",
		StateDimension:64,AdaptationEpochs:4,LearningRate:0.08,FifthUpdatesPerEpoch:len(fifthTrain),
		RecurrentParametersTrained:false,RouterUsed:false,ExactRecallUsed:false,AdaptiveScheduleUsed:false,SixthFamilyUsed:false,
	}
	arms:=[]string{"no_adaptation","fifth_only","five_family_cyclic"}
	for _,arm:=range arms{
		model:=uplm0oCloneModel(start)
		uplm1sTrain(model,arm,families)
		for f:=0;f<5;f++{
			m:=uplm0oByteEval(model,held[f],0,uplm1qFamilyNames[f]+"_heldout")
			res.Metrics=append(res.Metrics,UPLM1SFamilyMetric{Arm:arm,Family:uplm1qFamilyNames[f],Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity})
		}
		res.Summaries=append(res.Summaries,uplm1sSummary(arm,res.Metrics))
	}
	var no,iso,cyc UPLM1SSummary
	for _,s:=range res.Summaries{
		switch s.Arm{case "no_adaptation":no=s;case "fifth_only":iso=s;case "five_family_cyclic":cyc=s}
	}
	res.FifthOnlyMinusCyclicFifthAccuracy=iso.FifthHeldoutAccuracy-cyc.FifthHeldoutAccuracy
	res.FifthOnlyMinusCyclicPriorMeanAccuracy=iso.PriorMeanHeldoutAccuracy-cyc.PriorMeanHeldoutAccuracy
	res.FifthOnlyGainOverNoAdaptation=iso.FifthHeldoutAccuracy-no.FifthHeldoutAccuracy
	return res,nil
}
