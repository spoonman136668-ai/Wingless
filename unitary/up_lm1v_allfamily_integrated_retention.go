package unitary

const UPLM1VAllFamilySchema="wingless.up-lm1v-allfamily-integrated-retention.v1"

type UPLM1VIntegratedMetric struct{
	Arm string `json:"arm"`
	Family string `json:"family"`
	Stream int `json:"stream"`
	Metric UPLM0JMetric `json:"metric"`
}
type UPLM1VResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1USeal string `json:"source_up_lm1u_seal"`
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
	ByteSummaries []UPLM1TSummary `json:"byte_summaries"`
	RouterMetrics []UPLM1PRouterMetric `json:"router_metrics"`
	IntegratedMetrics []UPLM1VIntegratedMetric `json:"integrated_metrics"`
}

func RunUPLM1V()(UPLM1VResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil{return UPLM1VResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)

	d:=uplm1lStoreDirection()
	classifier:=uplm1pTrainClassifier(d)
	names:=uplm0gNames()

	res:=UPLM1VResult{
		Schema:UPLM1VAllFamilySchema,Experiment:"UP-LM1V-allfamily-integrated-retention",
		SourceUPLM1USeal:"d47b0f9b03e98b819d31e574f49bcf22d63f2b66",
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
	arms:=[]string{"equal_mass","fifth_1p5_mass","fifth_2x_mass"}
	for _,arm:=range arms{
		model:=uplm0oCloneModel(start)
		uplm1tTrain(model,arm,families)
		for f:=0;f<5;f++{
			m:=uplm0oByteEval(model,held[f],0,familyNames[f]+"_heldout")
			res.ByteMetrics=append(res.ByteMetrics,UPLM1TFamilyMetric{Arm:arm,Family:familyNames[f],Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity})
		}
		res.ByteSummaries=append(res.ByteSummaries,uplm1tSummary(arm,res.ByteMetrics))
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
