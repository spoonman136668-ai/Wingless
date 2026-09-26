package unitary

const UPLM1TFixedMassSchema = "wingless.up-lm1t-fixed-mass-plasticity-map.v1"

type UPLM1TFamilyMetric struct {
	Arm string `json:"arm"`
	Family string `json:"family"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity float64 `json:"perplexity"`
}

type UPLM1TSummary struct {
	Arm string `json:"arm"`
	PriorLearningRate float64 `json:"prior_learning_rate"`
	FifthLearningRate float64 `json:"fifth_learning_rate"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	MinHeldoutAccuracy float64 `json:"min_heldout_accuracy"`
	MeanHeldoutAccuracy float64 `json:"mean_heldout_accuracy"`
	HeldoutAccuracySpread float64 `json:"heldout_accuracy_spread"`
	FifthHeldoutAccuracy float64 `json:"fifth_heldout_accuracy"`
	FifthHeldoutPerplexity float64 `json:"fifth_heldout_perplexity"`
	PriorMeanHeldoutAccuracy float64 `json:"prior_mean_heldout_accuracy"`
	FifthAccuracyDeltaVsEqual float64 `json:"fifth_accuracy_delta_vs_equal"`
	PriorMeanDeltaVsEqual float64 `json:"prior_mean_delta_vs_equal"`
	MinAccuracyDeltaVsEqual float64 `json:"min_accuracy_delta_vs_equal"`
}

type UPLM1TFixedMassResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1SSeal string `json:"source_up_lm1s_seal"`
	StateDimension int `json:"state_dimension"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	UpdatesPerFamilyPerExample int `json:"updates_per_family_per_example"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RouterUsed bool `json:"router_used"`
	ExactRecallUsed bool `json:"exact_recall_used"`
	AdaptiveWeightingUsed bool `json:"adaptive_weighting_used"`
	SixthFamilyUsed bool `json:"sixth_family_used"`
	Metrics []UPLM1TFamilyMetric `json:"metrics"`
	Summaries []UPLM1TSummary `json:"summaries"`
}

func uplm1tRates(arm string)(prior,fifth float64) {
	switch arm {
	case "equal_mass": return 0.08,0.08
	case "fifth_1p5_mass": return 0.07,0.12
	case "fifth_2x_mass": return 0.06,0.16
	}
	return 0,0
}

func uplm1tTrain(model *uplm0aModel,arm string,families [5][]uplm0fExample) {
	prior,fifth:=uplm1tRates(arm)
	for epoch:=0;epoch<4;epoch++ {
		for i:=range families[0] {
			for f:=0;f<4;f++ { model.trainSentence(families[f][i].text,prior) }
			model.trainSentence(families[4][i].text,fifth)
		}
	}
}

func uplm1tSummary(arm string,metrics []UPLM1TFamilyMetric) UPLM1TSummary {
	priorLR,fifthLR:=uplm1tRates(arm)
	min,max,sum,priorSum:=1.0,-1.0,0.0,0.0
	fifthAcc,fifthPpl:=0.0,0.0
	count:=0
	for _,m:=range metrics {
		if m.Arm!=arm { continue }
		count++;sum+=m.Top1Accuracy
		if m.Top1Accuracy<min { min=m.Top1Accuracy }
		if m.Top1Accuracy>max { max=m.Top1Accuracy }
		if m.Family=="fifth" { fifthAcc=m.Top1Accuracy;fifthPpl=m.Perplexity } else { priorSum+=m.Top1Accuracy }
	}
	return UPLM1TSummary{
		Arm:arm,PriorLearningRate:priorLR,FifthLearningRate:fifthLR,
		TotalMassPerExample:4*priorLR+fifthLR,
		MinHeldoutAccuracy:min,MeanHeldoutAccuracy:sum/float64(count),
		HeldoutAccuracySpread:max-min,FifthHeldoutAccuracy:fifthAcc,
		FifthHeldoutPerplexity:fifthPpl,PriorMeanHeldoutAccuracy:priorSum/4,
	}
}

func RunUPLM1T()(UPLM1TFixedMassResult,error) {
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil { return UPLM1TFixedMassResult{},err }
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)

	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	held:=[5][]uplm0fExample{baseHeld,paraHeld,thirdHeld,fourthHeld,fifthHeld}

	res:=UPLM1TFixedMassResult{
		Schema:UPLM1TFixedMassSchema,Experiment:"UP-LM1T-fixed-mass-plasticity-map",
		SourceUPLM1SSeal:"4f572229717f454cbebe0f369c4fe8b5c6f6080e",
		StateDimension:64,AdaptationEpochs:4,UpdatesPerFamilyPerExample:1,TotalMassPerExample:0.40,
		RecurrentParametersTrained:false,RouterUsed:false,ExactRecallUsed:false,
		AdaptiveWeightingUsed:false,SixthFamilyUsed:false,
	}
	arms:=[]string{"equal_mass","fifth_1p5_mass","fifth_2x_mass"}
	for _,arm:=range arms {
		model:=uplm0oCloneModel(start)
		uplm1tTrain(model,arm,families)
		for f:=0;f<5;f++ {
			m:=uplm0oByteEval(model,held[f],0,uplm1qFamilyNames[f]+"_heldout")
			res.Metrics=append(res.Metrics,UPLM1TFamilyMetric{
				Arm:arm,Family:uplm1qFamilyNames[f],Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity,
			})
		}
		res.Summaries=append(res.Summaries,uplm1tSummary(arm,res.Metrics))
	}
	equal:=res.Summaries[0]
	for i:=range res.Summaries {
		res.Summaries[i].FifthAccuracyDeltaVsEqual=res.Summaries[i].FifthHeldoutAccuracy-equal.FifthHeldoutAccuracy
		res.Summaries[i].PriorMeanDeltaVsEqual=res.Summaries[i].PriorMeanHeldoutAccuracy-equal.PriorMeanHeldoutAccuracy
		res.Summaries[i].MinAccuracyDeltaVsEqual=res.Summaries[i].MinHeldoutAccuracy-equal.MinHeldoutAccuracy
	}
	return res,nil
}
