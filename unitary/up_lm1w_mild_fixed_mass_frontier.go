package unitary

const UPLM1WMildFrontierSchema="wingless.up-lm1w-mild-fixed-mass-frontier.v1"

type UPLM1WSummary struct{
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
type UPLM1WResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1VSeal string `json:"source_up_lm1v_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	RouterEpochs int `json:"router_epochs"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RecallCapChanged bool `json:"recall_cap_changed"`
	RouterModified bool `json:"router_modified"`
	AdaptiveWeightingUsed bool `json:"adaptive_weighting_used"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	ByteMetrics []UPLM1TFamilyMetric `json:"byte_metrics"`
	ByteSummaries []UPLM1WSummary `json:"byte_summaries"`
	RouterMetrics []UPLM1PRouterMetric `json:"router_metrics"`
	IntegratedMetrics []UPLM1VIntegratedMetric `json:"integrated_metrics"`
}

func uplm1wRates(arm string)(prior,fifth float64){
	switch arm{
	case "equal_mass": return 0.0800,0.0800
	case "fifth_1p125_mass": return 0.0775,0.0900
	case "fifth_1p25_mass": return 0.0750,0.1000
	case "fifth_1p375_mass": return 0.0725,0.1100
	case "fifth_1p5_mass": return 0.0700,0.1200
	}
	return 0,0
}
func uplm1wTrain(model *uplm0aModel,arm string,families [5][]uplm0fExample){
	prior,fifth:=uplm1wRates(arm)
	for epoch:=0;epoch<4;epoch++{
		for i:=range families[0]{
			for f:=0;f<4;f++{model.trainSentence(families[f][i].text,prior)}
			model.trainSentence(families[4][i].text,fifth)
		}
	}
}
func uplm1wSummary(arm string,metrics []UPLM1TFamilyMetric) UPLM1WSummary{
	priorLR,fifthLR:=uplm1wRates(arm)
	min,max,sum,priorSum:=1.0,-1.0,0.0,0.0
	fifthAcc,fifthPpl:=0.0,0.0
	count:=0
	for _,m:=range metrics{
		if m.Arm!=arm{continue}
		count++;sum+=m.Top1Accuracy
		if m.Top1Accuracy<min{min=m.Top1Accuracy}
		if m.Top1Accuracy>max{max=m.Top1Accuracy}
		if m.Family=="fifth"{fifthAcc=m.Top1Accuracy;fifthPpl=m.Perplexity}else{priorSum+=m.Top1Accuracy}
	}
	return UPLM1WSummary{
		Arm:arm,PriorLearningRate:priorLR,FifthLearningRate:fifthLR,
		TotalMassPerExample:4*priorLR+fifthLR,
		MinHeldoutAccuracy:min,MeanHeldoutAccuracy:sum/float64(count),
		HeldoutAccuracySpread:max-min,FifthHeldoutAccuracy:fifthAcc,
		FifthHeldoutPerplexity:fifthPpl,PriorMeanHeldoutAccuracy:priorSum/4,
	}
}
func RunUPLM1W()(UPLM1WResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil{return UPLM1WResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)

	d:=uplm1lStoreDirection()
	classifier:=uplm1pTrainClassifier(d)
	names:=uplm0gNames()

	res:=UPLM1WResult{
		Schema:UPLM1WMildFrontierSchema,Experiment:"UP-LM1W-mild-fixed-mass-frontier",
		SourceUPLM1VSeal:"c9349f4251421bf9395d0bff3e954463655d73b1",
		StateDimension:64,ExactRecallCap:16,RouterEpochs:20,AdaptationEpochs:4,TotalMassPerExample:0.40,
		RecurrentParametersTrained:false,RecallCapChanged:false,RouterModified:false,
		AdaptiveWeightingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	res.RouterMetrics=append(res.RouterMetrics,
		uplm1pRouterEval(classifier,d,"base","heldout",names[4:6],[]string{"stores","observes","reports"}),
		uplm1pRouterEval(classifier,d,"paraphrase","heldout",names[4:6],[]string{"saves","sees","recalls"}),
		uplm1pRouterEval(classifier,d,"third","heldout",names[4:6],[]string{"archives","notices","recounts"}),
		uplm1pRouterEval(classifier,d,"fourth","heldout",names[4:6],[]string{"retains","inspects","states"}),
		uplm1pRouterEval(classifier,d,"fifth","heldout",names[4:6],uplm1pFifthVerbs),
	)

	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	held:=[5][]uplm0fExample{baseHeld,paraHeld,thirdHeld,fourthHeld,fifthHeld}
	familyNames:=[5]string{"base","paraphrase","third","fourth","fifth"}
	arms:=[]string{"equal_mass","fifth_1p125_mass","fifth_1p25_mass","fifth_1p375_mass","fifth_1p5_mass"}

	for _,arm:=range arms{
		model:=uplm0oCloneModel(start)
		uplm1wTrain(model,arm,families)
		for f:=0;f<5;f++{
			m:=uplm0oByteEval(model,held[f],0,familyNames[f]+"_heldout")
			res.ByteMetrics=append(res.ByteMetrics,UPLM1TFamilyMetric{Arm:arm,Family:familyNames[f],Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity})
		}
		res.ByteSummaries=append(res.ByteSummaries,uplm1wSummary(arm,res.ByteMetrics))
		for f:=0;f<5;f++{
			for _,stream:=range []int{1,4}{
				res.IntegratedMetrics=append(res.IntegratedMetrics,UPLM1VIntegratedMetric{
					Arm:arm,Family:familyNames[f],Stream:stream,
					Metric:uplm1mEvaluate(model,classifier,d,held[f],stream,familyNames[f]+"_block_stream"+itoa(stream)),
				})
			}
		}
	}
	eq:=res.ByteSummaries[0]
	for i:=range res.ByteSummaries{
		res.ByteSummaries[i].FifthAccuracyDeltaVsEqual=res.ByteSummaries[i].FifthHeldoutAccuracy-eq.FifthHeldoutAccuracy
		res.ByteSummaries[i].PriorMeanDeltaVsEqual=res.ByteSummaries[i].PriorMeanHeldoutAccuracy-eq.PriorMeanHeldoutAccuracy
		res.ByteSummaries[i].MinAccuracyDeltaVsEqual=res.ByteSummaries[i].MinHeldoutAccuracy-eq.MinHeldoutAccuracy
	}
	return res,nil
}
