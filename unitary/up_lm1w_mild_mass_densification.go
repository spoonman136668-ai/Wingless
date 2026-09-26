package unitary

const UPLM1WMildMassSchema="wingless.up-lm1w-mild-mass-densification.v1"

type UPLM1WFamilySummary struct {
	Arm string `json:"arm"`
	Family string `json:"family"`
	PriorLearningRate float64 `json:"prior_learning_rate"`
	FifthLearningRate float64 `json:"fifth_learning_rate"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	ByteTop1Accuracy float64 `json:"byte_top1_accuracy"`
	BytePerplexity float64 `json:"byte_perplexity"`
	IntegratedMeanTop1Accuracy float64 `json:"integrated_mean_top1_accuracy"`
	IntegratedMinTop1Accuracy float64 `json:"integrated_min_top1_accuracy"`
	IntegratedMeanPerplexity float64 `json:"integrated_mean_perplexity"`
	StructuralExactAllCells bool `json:"structural_exact_all_cells"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM1WResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1VSeal string `json:"source_up_lm1v_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	Arms int `json:"arms"`
	Families int `json:"families"`
	StreamDepths []int `json:"stream_depths"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RecallCapChanged bool `json:"recall_cap_changed"`
	RouterModified bool `json:"router_modified"`
	AdaptiveWeightingUsed bool `json:"adaptive_weighting_used"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	RouterMetrics []UPLM1PRouterMetric `json:"router_metrics"`
	ByteMetrics []UPLM1TFamilyMetric `json:"byte_metrics"`
	IntegratedMetrics []UPLM1VIntegratedMetric `json:"integrated_metrics"`
	FamilySummaries []UPLM1WFamilySummary `json:"family_summaries"`
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
func uplm1wStructuralExact(m UPLM0JMetric) bool{
	return m.DependentFirstByteAccuracy==1 &&
		m.QuerySetExactAccuracy==1 &&
		m.Update0ExactAccuracy==1 &&
		m.Update1ExactAccuracy==1 &&
		m.Update2ExactAccuracy==1 &&
		m.Update4ExactAccuracy==1 &&
		m.AdmissionPrecision==1 &&
		m.AdmissionRecall==1 &&
		m.EventRoutingAccuracy==1 &&
		m.ReportRoutingAccuracy==1
}
func uplm1wSummary(arm,family string,byteMetric UPLM1TFamilyMetric,metrics []UPLM1VIntegratedMetric) UPLM1WFamilySummary{
	prior,fifth:=uplm1wRates(arm)
	sumAcc,sumPpl:=0.0,0.0
	minAcc:=1.0
	count,maxRecall:=0,0
	exact:=true
	for _,x:=range metrics{
		if x.Arm!=arm||x.Family!=family{continue}
		count++
		if x.Metric.Top1Accuracy<minAcc{minAcc=x.Metric.Top1Accuracy}
		sumAcc+=x.Metric.Top1Accuracy
		sumPpl+=x.Metric.Perplexity
		if x.Metric.MaxRecallEntries>maxRecall{maxRecall=x.Metric.MaxRecallEntries}
		if !uplm1wStructuralExact(x.Metric){exact=false}
	}
	return UPLM1WFamilySummary{
		Arm:arm,Family:family,PriorLearningRate:prior,FifthLearningRate:fifth,
		TotalMassPerExample:4*prior+fifth,
		ByteTop1Accuracy:byteMetric.Top1Accuracy,BytePerplexity:byteMetric.Perplexity,
		IntegratedMeanTop1Accuracy:sumAcc/float64(count),IntegratedMinTop1Accuracy:minAcc,
		IntegratedMeanPerplexity:sumPpl/float64(count),StructuralExactAllCells:exact,MaxRecallEntries:maxRecall,
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
	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	held:=[5][]uplm0fExample{baseHeld,paraHeld,thirdHeld,fourthHeld,fifthHeld}
	familyNames:=[5]string{"base","paraphrase","third","fourth","fifth"}

	res:=UPLM1WResult{
		Schema:UPLM1WMildMassSchema,Experiment:"UP-LM1W-mild-mass-densification",
		SourceUPLM1VSeal:"c9349f4251421bf9395d0bff3e954463655d73b1",
		StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,
		Arms:5,Families:5,StreamDepths:[]int{1,4},
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

	arms:=[]string{"equal_mass","fifth_1p125_mass","fifth_1p25_mass","fifth_1p375_mass","fifth_1p5_mass"}
	for _,arm:=range arms{
		model:=uplm0oCloneModel(start)
		uplm1wTrain(model,arm,families)
		armBytes:=make([]UPLM1TFamilyMetric,5)
		for f:=0;f<5;f++{
			m:=uplm0oByteEval(model,held[f],0,familyNames[f]+"_heldout")
			bm:=UPLM1TFamilyMetric{Arm:arm,Family:familyNames[f],Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity}
			armBytes[f]=bm
			res.ByteMetrics=append(res.ByteMetrics,bm)
			for _,stream:=range []int{1,4}{
				res.IntegratedMetrics=append(res.IntegratedMetrics,UPLM1VIntegratedMetric{
					Arm:arm,Family:familyNames[f],Stream:stream,
					Metric:uplm1mEvaluate(model,classifier,d,held[f],stream,familyNames[f]+"_block_stream"+itoa(stream)),
				})
			}
		}
		for f:=0;f<5;f++{
			res.FamilySummaries=append(res.FamilySummaries,uplm1wSummary(arm,familyNames[f],armBytes[f],res.IntegratedMetrics))
		}
	}
	return res,nil
}
