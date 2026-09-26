package unitary

const UPLM2BPositionSchema="wingless.up-lm2b-prior-position-map.v1"

type UPLM2BMetric struct{
	Allocation string `json:"allocation"`
	Rotation int `json:"rotation"`
	Family string `json:"family"`
	PriorPosition int `json:"prior_position"`
	Stream int `json:"stream"`
	Metric UPLM0JMetric `json:"metric"`
}
type UPLM2BFamilySummary struct{
	Allocation string `json:"allocation"`
	Rotation int `json:"rotation"`
	Family string `json:"family"`
	PriorPosition int `json:"prior_position"`
	IntegratedMeanTop1Accuracy float64 `json:"integrated_mean_top1_accuracy"`
	IntegratedMinTop1Accuracy float64 `json:"integrated_min_top1_accuracy"`
	IntegratedMeanPerplexity float64 `json:"integrated_mean_perplexity"`
	StructuralExactAllCells bool `json:"structural_exact_all_cells"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM2BRotationSummary struct{
	Allocation string `json:"allocation"`
	Rotation int `json:"rotation"`
	PriorOrder []string `json:"prior_order"`
	FifthIntegratedMean float64 `json:"fifth_integrated_mean"`
	PriorIntegratedMean float64 `json:"prior_integrated_mean"`
	AllFamilyIntegratedMean float64 `json:"all_family_integrated_mean"`
	MinFamilyIntegratedMean float64 `json:"min_family_integrated_mean"`
	StructuralExactAllFamilies bool `json:"structural_exact_all_families"`
}
type UPLM2BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2ASeal string `json:"source_up_lm2a_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	Allocations int `json:"allocations"`
	Rotations int `json:"rotations"`
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
	Metrics []UPLM2BMetric `json:"metrics"`
	FamilySummaries []UPLM2BFamilySummary `json:"family_summaries"`
	RotationSummaries []UPLM2BRotationSummary `json:"rotation_summaries"`
}

func uplm2bOrder(rotation int)[4]int{
	return [4]int{rotation%4,(rotation+1)%4,(rotation+2)%4,(rotation+3)%4}
}
func uplm2bPosition(rotation,family int)int{
	order:=uplm2bOrder(rotation)
	for i,f:=range order{if f==family{return i+1}}
	return 5
}
func uplm2bTrain(model *uplm0aModel,allocation string,rotation int,families [5][]uplm0fExample){
	prior,fifth:=uplm1wRates(allocation)
	order:=uplm2bOrder(rotation)
	for epoch:=0;epoch<4;epoch++{
		for i:=range families[0]{
			for _,f:=range order{model.trainSentence(families[f][i].text,prior)}
			model.trainSentence(families[4][i].text,fifth)
		}
	}
}
func uplm2bFamilySummary(allocation string,rotation int,family string,pos int,metrics []UPLM2BMetric)UPLM2BFamilySummary{
	sumAcc,sumPpl:=0.0,0.0
	minAcc:=1.0
	count,maxRecall:=0,0
	exact:=true
	for _,x:=range metrics{
		if x.Allocation!=allocation||x.Rotation!=rotation||x.Family!=family{continue}
		count++;sumAcc+=x.Metric.Top1Accuracy;sumPpl+=x.Metric.Perplexity
		if x.Metric.Top1Accuracy<minAcc{minAcc=x.Metric.Top1Accuracy}
		if x.Metric.MaxRecallEntries>maxRecall{maxRecall=x.Metric.MaxRecallEntries}
		if !uplm1wStructuralExact(x.Metric){exact=false}
	}
	return UPLM2BFamilySummary{Allocation:allocation,Rotation:rotation,Family:family,PriorPosition:pos,IntegratedMeanTop1Accuracy:sumAcc/float64(count),IntegratedMinTop1Accuracy:minAcc,IntegratedMeanPerplexity:sumPpl/float64(count),StructuralExactAllCells:exact,MaxRecallEntries:maxRecall}
}
func uplm2bRotationSummary(allocation string,rotation int,names []string,s []UPLM2BFamilySummary)UPLM2BRotationSummary{
	sum,prior,fifth,min:=0.0,0.0,0.0,1.0
	count,priorN:=0,0
	exact:=true
	for _,x:=range s{
		if x.Allocation!=allocation||x.Rotation!=rotation{continue}
		count++;sum+=x.IntegratedMeanTop1Accuracy
		if x.IntegratedMeanTop1Accuracy<min{min=x.IntegratedMeanTop1Accuracy}
		if x.Family=="fifth"{fifth=x.IntegratedMeanTop1Accuracy}else{prior+=x.IntegratedMeanTop1Accuracy;priorN++}
		if !x.StructuralExactAllCells{exact=false}
	}
	order:=uplm2bOrder(rotation)
	labels:=[]string{names[order[0]],names[order[1]],names[order[2]],names[order[3]]}
	return UPLM2BRotationSummary{Allocation:allocation,Rotation:rotation,PriorOrder:labels,FifthIntegratedMean:fifth,PriorIntegratedMean:prior/float64(priorN),AllFamilyIntegratedMean:sum/float64(count),MinFamilyIntegratedMean:min,StructuralExactAllFamilies:exact}
}

func RunUPLM2B()(UPLM2BResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil{return UPLM2BResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection()
	classifier:=uplm1pTrainClassifier(d)
	names:=uplm0gNames()

	res:=UPLM2BResult{
		Schema:UPLM2BPositionSchema,Experiment:"UP-LM2B-prior-position-map",
		SourceUPLM2ASeal:"ccdbc414cc87ccf331d7504a370368a7c550cf17",
		StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,
		Allocations:5,Rotations:4,Families:5,StreamDepths:[]int{1,4},FifthAlwaysLast:true,
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
		for rotation:=0;rotation<4;rotation++{
			model:=uplm0oCloneModel(start)
			uplm2bTrain(model,allocation,rotation,families)
			for f,family:=range familyNames{
				pos:=5;if f<4{pos=uplm2bPosition(rotation,f)}
				for _,stream:=range []int{1,4}{
					res.Metrics=append(res.Metrics,UPLM2BMetric{Allocation:allocation,Rotation:rotation,Family:family,PriorPosition:pos,Stream:stream,Metric:uplm1mEvaluate(model,classifier,d,held[f],stream,family+"_block_stream"+itoa(stream))})
				}
				res.FamilySummaries=append(res.FamilySummaries,uplm2bFamilySummary(allocation,rotation,family,pos,res.Metrics))
			}
			res.RotationSummaries=append(res.RotationSummaries,uplm2bRotationSummary(allocation,rotation,familyNames,res.FamilySummaries))
		}
	}
	return res,nil
}
