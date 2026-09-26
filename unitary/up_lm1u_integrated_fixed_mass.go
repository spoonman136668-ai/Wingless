package unitary

const UPLM1UIntegratedSchema="wingless.up-lm1u-integrated-fixed-mass.v1"

type UPLM1UResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1TSeal string `json:"source_up_lm1t_seal"`
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
	IntegratedMetrics []UPLM1PIntegratedMetric `json:"integrated_metrics"`
}

func RunUPLM1U()(UPLM1UResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil{return UPLM1UResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)

	d:=uplm1lStoreDirection()
	classifier:=uplm1pTrainClassifier(d)
	names:=uplm0gNames()

	res:=UPLM1UResult{
		Schema:UPLM1UIntegratedSchema,Experiment:"UP-LM1U-integrated-fixed-mass",
		SourceUPLM1TSeal:"2dd43d89cdd8480eec3c7be779b2de5f48596159",
		StateDimension:64,ExactRecallCap:16,RouterEpochs:20,AdaptationEpochs:4,TotalMassPerExample:0.40,
		RecurrentParametersTrained:false,RecallCapChanged:false,RouterModified:false,
		AdaptiveWeightingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	res.RouterMetrics=append(res.RouterMetrics,
		uplm1pRouterEval(classifier,d,"base","heldout",names[4:6],[]string{"stores","observes","reports"}),
		uplm1pRouterEval(classifier,d,"paraphrase","heldout",names[4:6],[]string{"saves","sees","recalls"}),
		uplm1pRouterEval(classifier,d,"third","heldout",names[4:6],[]string{"archives","notices","recounts"}),
		uplm1pRouterEval(classifier,d,"fourth","heldout",names[4:6],[]string{"retains","inspects","states"}),
		uplm1pRouterEval(classifier,d,"fifth","train",names[:4],uplm1pFifthVerbs),
		uplm1pRouterEval(classifier,d,"fifth","heldout",names[4:6],uplm1pFifthVerbs),
		uplm1pRouterEval(classifier,d,"fifth","unseen",uplm1kUnseenNames,uplm1pFifthVerbs),
	)

	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	held:=[5][]uplm0fExample{baseHeld,paraHeld,thirdHeld,fourthHeld,fifthHeld}
	arms:=[]string{"equal_mass","fifth_1p5_mass","fifth_2x_mass"}
	for _,arm:=range arms{
		model:=uplm0oCloneModel(start)
		uplm1tTrain(model,arm,families)
		for f:=0;f<5;f++{
			m:=uplm0oByteEval(model,held[f],0,uplm1qFamilyNames[f]+"_heldout")
			res.ByteMetrics=append(res.ByteMetrics,UPLM1TFamilyMetric{Arm:arm,Family:uplm1qFamilyNames[f],Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity})
		}
		s:=uplm1tSummary(arm,res.ByteMetrics)
		res.ByteSummaries=append(res.ByteSummaries,s)

		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"}{
			h:=uplm1pHeldout(order)
			for _,stream:=range []int{1,4}{
				res.IntegratedMetrics=append(res.IntegratedMetrics,UPLM1PIntegratedMetric{
					Arm:arm,Metric:uplm1mEvaluate(model,classifier,d,h,stream,"fifth_"+order+"_stream"+itoa(stream)),
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
