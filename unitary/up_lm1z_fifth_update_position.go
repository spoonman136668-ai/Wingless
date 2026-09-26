package unitary

const UPLM1ZFifthPositionSchema="wingless.up-lm1z-fifth-update-position.v1"

type UPLM1ZMetric struct{
	Allocation string `json:"allocation"`
	FifthPosition int `json:"fifth_position"`
	Family string `json:"family"`
	Stream int `json:"stream"`
	Metric UPLM0JMetric `json:"metric"`
}
type UPLM1ZFamilySummary struct{
	Allocation string `json:"allocation"`
	FifthPosition int `json:"fifth_position"`
	Family string `json:"family"`
	IntegratedMeanTop1Accuracy float64 `json:"integrated_mean_top1_accuracy"`
	IntegratedMinTop1Accuracy float64 `json:"integrated_min_top1_accuracy"`
	IntegratedMeanPerplexity float64 `json:"integrated_mean_perplexity"`
	StructuralExactAllCells bool `json:"structural_exact_all_cells"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM1ZPositionSummary struct{
	Allocation string `json:"allocation"`
	FifthPosition int `json:"fifth_position"`
	FifthIntegratedMean float64 `json:"fifth_integrated_mean"`
	PriorIntegratedMean float64 `json:"prior_integrated_mean"`
	AllFamilyIntegratedMean float64 `json:"all_family_integrated_mean"`
	MinFamilyIntegratedMean float64 `json:"min_family_integrated_mean"`
	StructuralExactAllFamilies bool `json:"structural_exact_all_families"`
}
type UPLM1ZResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1YSeal string `json:"source_up_lm1y_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	Allocations int `json:"allocations"`
	FifthPositions int `json:"fifth_positions"`
	Families int `json:"families"`
	StreamDepths []int `json:"stream_depths"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RecallCapChanged bool `json:"recall_cap_changed"`
	RouterModified bool `json:"router_modified"`
	AdaptiveWeightingUsed bool `json:"adaptive_weighting_used"`
	AdaptiveOrderUsed bool `json:"adaptive_order_used"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	RouterMetrics []UPLM1PRouterMetric `json:"router_metrics"`
	Metrics []UPLM1ZMetric `json:"metrics"`
	FamilySummaries []UPLM1ZFamilySummary `json:"family_summaries"`
	PositionSummaries []UPLM1ZPositionSummary `json:"position_summaries"`
}

func uplm1zFamilyOrder(fifthPosition int)[5]int{
	priors:=[4]int{0,1,2,3}
	var out [5]int
	pi:=0
	for pos:=1;pos<=5;pos++{
		if pos==fifthPosition{out[pos-1]=4}else{out[pos-1]=priors[pi];pi++}
	}
	return out
}
func uplm1zTrain(model *uplm0aModel,allocation string,fifthPosition int,families [5][]uplm0fExample){
	prior,fifth:=uplm1wRates(allocation)
	order:=uplm1zFamilyOrder(fifthPosition)
	for epoch:=0;epoch<4;epoch++{
		for i:=range families[0]{
			for _,f:=range order{
				lr:=prior
				if f==4{lr=fifth}
				model.trainSentence(families[f][i].text,lr)
			}
		}
	}
}
func uplm1zFamilySummary(allocation string,pos int,family string,metrics []UPLM1ZMetric)UPLM1ZFamilySummary{
	sumAcc,sumPpl:=0.0,0.0
	minAcc:=1.0
	count,maxRecall:=0,0
	exact:=true
	for _,x:=range metrics{
		if x.Allocation!=allocation||x.FifthPosition!=pos||x.Family!=family{continue}
		count++;sumAcc+=x.Metric.Top1Accuracy;sumPpl+=x.Metric.Perplexity
		if x.Metric.Top1Accuracy<minAcc{minAcc=x.Metric.Top1Accuracy}
		if x.Metric.MaxRecallEntries>maxRecall{maxRecall=x.Metric.MaxRecallEntries}
		if !uplm1wStructuralExact(x.Metric){exact=false}
	}
	return UPLM1ZFamilySummary{Allocation:allocation,FifthPosition:pos,Family:family,IntegratedMeanTop1Accuracy:sumAcc/float64(count),IntegratedMinTop1Accuracy:minAcc,IntegratedMeanPerplexity:sumPpl/float64(count),StructuralExactAllCells:exact,MaxRecallEntries:maxRecall}
}
func uplm1zPositionSummary(allocation string,pos int,s []UPLM1ZFamilySummary)UPLM1ZPositionSummary{
	sum,prior,fifth,min:=0.0,0.0,0.0,1.0
	count,priorN:=0,0
	exact:=true
	for _,x:=range s{
		if x.Allocation!=allocation||x.FifthPosition!=pos{continue}
		count++;sum+=x.IntegratedMeanTop1Accuracy
		if x.IntegratedMeanTop1Accuracy<min{min=x.IntegratedMeanTop1Accuracy}
		if x.Family=="fifth"{fifth=x.IntegratedMeanTop1Accuracy}else{prior+=x.IntegratedMeanTop1Accuracy;priorN++}
		if !x.StructuralExactAllCells{exact=false}
	}
	return UPLM1ZPositionSummary{Allocation:allocation,FifthPosition:pos,FifthIntegratedMean:fifth,PriorIntegratedMean:prior/float64(priorN),AllFamilyIntegratedMean:sum/float64(count),MinFamilyIntegratedMean:min,StructuralExactAllFamilies:exact}
}

func RunUPLM1Z()(UPLM1ZResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil{return UPLM1ZResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection()
	classifier:=uplm1pTrainClassifier(d)
	names:=uplm0gNames()

	res:=UPLM1ZResult{
		Schema:UPLM1ZFifthPositionSchema,Experiment:"UP-LM1Z-fifth-update-position",
		SourceUPLM1YSeal:"feeca16b52f60cb0dfdac58de46ac02ff061b0c5",
		StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,
		Allocations:5,FifthPositions:5,Families:5,StreamDepths:[]int{1,4},
		RecurrentParametersTrained:false,RecallCapChanged:false,RouterModified:false,
		AdaptiveWeightingUsed:false,AdaptiveOrderUsed:false,AttentionUsed:false,FutureOracleUsed:false,
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
	familyNames:=[]string{"base","paraphrase","third","fourth","fifth"}
	allocations:=[]string{"equal_mass","fifth_1p125_mass","fifth_1p25_mass","fifth_1p375_mass","fifth_1p5_mass"}

	for _,allocation:=range allocations{
		for pos:=1;pos<=5;pos++{
			model:=uplm0oCloneModel(start)
			uplm1zTrain(model,allocation,pos,families)
			for f,family:=range familyNames{
				for _,stream:=range []int{1,4}{
					res.Metrics=append(res.Metrics,UPLM1ZMetric{
						Allocation:allocation,FifthPosition:pos,Family:family,Stream:stream,
						Metric:uplm1mEvaluate(model,classifier,d,held[f],stream,family+"_block_stream"+itoa(stream)),
					})
				}
				res.FamilySummaries=append(res.FamilySummaries,uplm1zFamilySummary(allocation,pos,family,res.Metrics))
			}
			res.PositionSummaries=append(res.PositionSummaries,uplm1zPositionSummary(allocation,pos,res.FamilySummaries))
		}
	}
	return res,nil
}
