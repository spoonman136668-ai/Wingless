package unitary

const UPLM2ATerminalFifthSchema="wingless.up-lm2a-terminal-fifth-prior-order.v1"

type UPLM2AMetric struct{
	Allocation string `json:"allocation"`
	PriorOrder string `json:"prior_order"`
	Family string `json:"family"`
	Stream int `json:"stream"`
	Metric UPLM0JMetric `json:"metric"`
}
type UPLM2AFamilySummary struct{
	Allocation string `json:"allocation"`
	PriorOrder string `json:"prior_order"`
	Family string `json:"family"`
	IntegratedMeanTop1Accuracy float64 `json:"integrated_mean_top1_accuracy"`
	IntegratedMinTop1Accuracy float64 `json:"integrated_min_top1_accuracy"`
	IntegratedMeanPerplexity float64 `json:"integrated_mean_perplexity"`
	StructuralExactAllCells bool `json:"structural_exact_all_cells"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM2AOrderSummary struct{
	Allocation string `json:"allocation"`
	PriorOrder string `json:"prior_order"`
	FifthIntegratedMean float64 `json:"fifth_integrated_mean"`
	PriorIntegratedMean float64 `json:"prior_integrated_mean"`
	AllFamilyIntegratedMean float64 `json:"all_family_integrated_mean"`
	MinFamilyIntegratedMean float64 `json:"min_family_integrated_mean"`
	StructuralExactAllFamilies bool `json:"structural_exact_all_families"`
}
type UPLM2AResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1ZSeal string `json:"source_up_lm1z_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	Allocations int `json:"allocations"`
	PriorOrders int `json:"prior_orders"`
	Families int `json:"families"`
	StreamDepths []int `json:"stream_depths"`
	FifthAlwaysLast bool `json:"fifth_always_last"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RecallCapChanged bool `json:"recall_cap_changed"`
	RouterModified bool `json:"router_modified"`
	AdaptiveWeightingUsed bool `json:"adaptive_weighting_used"`
	AdaptiveOrderUsed bool `json:"adaptive_order_used"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	RouterMetrics []UPLM1PRouterMetric `json:"router_metrics"`
	Metrics []UPLM2AMetric `json:"metrics"`
	FamilySummaries []UPLM2AFamilySummary `json:"family_summaries"`
	OrderSummaries []UPLM2AOrderSummary `json:"order_summaries"`
}

func uplm2aPriorOrder(order string,epoch,i int)[4]int{
	switch order{
	case "canonical_prior":
		return [4]int{0,1,2,3}
	case "reverse_prior":
		return [4]int{3,2,1,0}
	default:
		start:=(epoch+i)%4
		return [4]int{start,(start+1)%4,(start+2)%4,(start+3)%4}
	}
}
func uplm2aTrain(model *uplm0aModel,allocation,order string,families [5][]uplm0fExample){
	prior,fifth:=uplm1wRates(allocation)
	for epoch:=0;epoch<4;epoch++{
		for i:=range families[0]{
			seq:=uplm2aPriorOrder(order,epoch,i)
			for _,f:=range seq{model.trainSentence(families[f][i].text,prior)}
			model.trainSentence(families[4][i].text,fifth)
		}
	}
}
func uplm2aFamilySummary(allocation,order,family string,metrics []UPLM2AMetric)UPLM2AFamilySummary{
	sumAcc,sumPpl:=0.0,0.0
	minAcc:=1.0
	count,maxRecall:=0,0
	exact:=true
	for _,x:=range metrics{
		if x.Allocation!=allocation||x.PriorOrder!=order||x.Family!=family{continue}
		count++;sumAcc+=x.Metric.Top1Accuracy;sumPpl+=x.Metric.Perplexity
		if x.Metric.Top1Accuracy<minAcc{minAcc=x.Metric.Top1Accuracy}
		if x.Metric.MaxRecallEntries>maxRecall{maxRecall=x.Metric.MaxRecallEntries}
		if !uplm1wStructuralExact(x.Metric){exact=false}
	}
	return UPLM2AFamilySummary{Allocation:allocation,PriorOrder:order,Family:family,IntegratedMeanTop1Accuracy:sumAcc/float64(count),IntegratedMinTop1Accuracy:minAcc,IntegratedMeanPerplexity:sumPpl/float64(count),StructuralExactAllCells:exact,MaxRecallEntries:maxRecall}
}
func uplm2aOrderSummary(allocation,order string,s []UPLM2AFamilySummary)UPLM2AOrderSummary{
	sum,prior,fifth,min:=0.0,0.0,0.0,1.0
	count,priorN:=0,0
	exact:=true
	for _,x:=range s{
		if x.Allocation!=allocation||x.PriorOrder!=order{continue}
		count++;sum+=x.IntegratedMeanTop1Accuracy
		if x.IntegratedMeanTop1Accuracy<min{min=x.IntegratedMeanTop1Accuracy}
		if x.Family=="fifth"{fifth=x.IntegratedMeanTop1Accuracy}else{prior+=x.IntegratedMeanTop1Accuracy;priorN++}
		if !x.StructuralExactAllCells{exact=false}
	}
	return UPLM2AOrderSummary{Allocation:allocation,PriorOrder:order,FifthIntegratedMean:fifth,PriorIntegratedMean:prior/float64(priorN),AllFamilyIntegratedMean:sum/float64(count),MinFamilyIntegratedMean:min,StructuralExactAllFamilies:exact}
}

func RunUPLM2A()(UPLM2AResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil{return UPLM2AResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection()
	classifier:=uplm1pTrainClassifier(d)
	names:=uplm0gNames()

	res:=UPLM2AResult{
		Schema:UPLM2ATerminalFifthSchema,Experiment:"UP-LM2A-terminal-fifth-prior-order",
		SourceUPLM1ZSeal:"89c18dbebdc5c2d26cf763235e97dfb77c6cfcf4",
		StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,
		Allocations:5,PriorOrders:3,Families:5,StreamDepths:[]int{1,4},FifthAlwaysLast:true,
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
	orders:=[]string{"canonical_prior","reverse_prior","rotating_prior"}

	for _,allocation:=range allocations{
		for _,order:=range orders{
			model:=uplm0oCloneModel(start)
			uplm2aTrain(model,allocation,order,families)
			for f,family:=range familyNames{
				for _,stream:=range []int{1,4}{
					res.Metrics=append(res.Metrics,UPLM2AMetric{
						Allocation:allocation,PriorOrder:order,Family:family,Stream:stream,
						Metric:uplm1mEvaluate(model,classifier,d,held[f],stream,family+"_block_stream"+itoa(stream)),
					})
				}
				res.FamilySummaries=append(res.FamilySummaries,uplm2aFamilySummary(allocation,order,family,res.Metrics))
			}
			res.OrderSummaries=append(res.OrderSummaries,uplm2aOrderSummary(allocation,order,res.FamilySummaries))
		}
	}
	return res,nil
}
