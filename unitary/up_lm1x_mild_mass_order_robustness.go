package unitary

const UPLM1XOrderRobustnessSchema="wingless.up-lm1x-mild-mass-order-robustness.v1"

type UPLM1XIntegratedMetric struct{
	Arm string `json:"arm"`
	Family string `json:"family"`
	Order string `json:"order"`
	Stream int `json:"stream"`
	Metric UPLM0JMetric `json:"metric"`
}
type UPLM1XFamilyOrderSummary struct{
	Arm string `json:"arm"`
	Family string `json:"family"`
	Order string `json:"order"`
	IntegratedMeanTop1Accuracy float64 `json:"integrated_mean_top1_accuracy"`
	IntegratedMinTop1Accuracy float64 `json:"integrated_min_top1_accuracy"`
	IntegratedMeanPerplexity float64 `json:"integrated_mean_perplexity"`
	StructuralExactAllCells bool `json:"structural_exact_all_cells"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM1XOrderSummary struct{
	Arm string `json:"arm"`
	Order string `json:"order"`
	FifthIntegratedMean float64 `json:"fifth_integrated_mean"`
	PriorIntegratedMean float64 `json:"prior_integrated_mean"`
	AllFamilyIntegratedMean float64 `json:"all_family_integrated_mean"`
	MinFamilyIntegratedMean float64 `json:"min_family_integrated_mean"`
	StructuralExactAllFamilies bool `json:"structural_exact_all_families"`
}
type UPLM1XResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM1WSeal string `json:"source_up_lm1w_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	Arms int `json:"arms"`
	Families int `json:"families"`
	Orders int `json:"orders"`
	StreamDepths []int `json:"stream_depths"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RecallCapChanged bool `json:"recall_cap_changed"`
	RouterModified bool `json:"router_modified"`
	AdaptiveWeightingUsed bool `json:"adaptive_weighting_used"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	RouterMetrics []UPLM1PRouterMetric `json:"router_metrics"`
	IntegratedMetrics []UPLM1XIntegratedMetric `json:"integrated_metrics"`
	FamilyOrderSummaries []UPLM1XFamilyOrderSummary `json:"family_order_summaries"`
	OrderSummaries []UPLM1XOrderSummary `json:"order_summaries"`
}

func uplm1xHeldout(family,order string,block [5][]uplm0fExample) []uplm0fExample{
	if order=="block"{
		switch family{
		case "base": return block[0]
		case "paraphrase": return block[1]
		case "third": return block[2]
		case "fourth": return block[3]
		default: return block[4]
		}
	}
	switch family{
	case "base": return uplm0mHeldout(order)
	case "paraphrase": return uplm0nHeldout(order)
	case "third": return uplm1cHeldout(order)
	case "fourth": return uplm1jHeldout(order)
	default: return uplm1pHeldout(order)
	}
}
func uplm1xFamilySummary(arm,family,order string,metrics []UPLM1XIntegratedMetric)UPLM1XFamilyOrderSummary{
	sumAcc,sumPpl:=0.0,0.0
	minAcc:=1.0
	count,maxRecall:=0,0
	exact:=true
	for _,x:=range metrics{
		if x.Arm!=arm||x.Family!=family||x.Order!=order{continue}
		count++;sumAcc+=x.Metric.Top1Accuracy;sumPpl+=x.Metric.Perplexity
		if x.Metric.Top1Accuracy<minAcc{minAcc=x.Metric.Top1Accuracy}
		if x.Metric.MaxRecallEntries>maxRecall{maxRecall=x.Metric.MaxRecallEntries}
		if !uplm1wStructuralExact(x.Metric){exact=false}
	}
	return UPLM1XFamilyOrderSummary{Arm:arm,Family:family,Order:order,IntegratedMeanTop1Accuracy:sumAcc/float64(count),IntegratedMinTop1Accuracy:minAcc,IntegratedMeanPerplexity:sumPpl/float64(count),StructuralExactAllCells:exact,MaxRecallEntries:maxRecall}
}
func uplm1xOrderSummary(arm,order string,s []UPLM1XFamilyOrderSummary)UPLM1XOrderSummary{
	sum,prior,fifth,min:=0.0,0.0,0.0,1.0
	count,priorN:=0,0
	exact:=true
	for _,x:=range s{
		if x.Arm!=arm||x.Order!=order{continue}
		count++;sum+=x.IntegratedMeanTop1Accuracy
		if x.IntegratedMeanTop1Accuracy<min{min=x.IntegratedMeanTop1Accuracy}
		if x.Family=="fifth"{fifth=x.IntegratedMeanTop1Accuracy}else{prior+=x.IntegratedMeanTop1Accuracy;priorN++}
		if !x.StructuralExactAllCells{exact=false}
	}
	return UPLM1XOrderSummary{Arm:arm,Order:order,FifthIntegratedMean:fifth,PriorIntegratedMean:prior/float64(priorN),AllFamilyIntegratedMean:sum/float64(count),MinFamilyIntegratedMean:min,StructuralExactAllFamilies:exact}
}

func RunUPLM1X()(UPLM1XResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil{return UPLM1XResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection()
	classifier:=uplm1pTrainClassifier(d)
	names:=uplm0gNames()

	res:=UPLM1XResult{
		Schema:UPLM1XOrderRobustnessSchema,Experiment:"UP-LM1X-mild-mass-order-robustness",
		SourceUPLM1WSeal:"8564bb800df3d0373efd3d361475037eb731aaf2",
		StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,
		Arms:5,Families:5,Orders:5,StreamDepths:[]int{1,4},
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
	block:=[5][]uplm0fExample{baseHeld,paraHeld,thirdHeld,fourthHeld,fifthHeld}
	familyNames:=[]string{"base","paraphrase","third","fourth","fifth"}
	orders:=[]string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"}
	arms:=[]string{"equal_mass","fifth_1p125_mass","fifth_1p25_mass","fifth_1p375_mass","fifth_1p5_mass"}

	for _,arm:=range arms{
		model:=uplm0oCloneModel(start)
		uplm1wTrain(model,arm,families)
		for _,family:=range familyNames{
			for _,order:=range orders{
				held:=uplm1xHeldout(family,order,block)
				for _,stream:=range []int{1,4}{
					res.IntegratedMetrics=append(res.IntegratedMetrics,UPLM1XIntegratedMetric{
						Arm:arm,Family:family,Order:order,Stream:stream,
						Metric:uplm1mEvaluate(model,classifier,d,held,stream,family+"_"+order+"_stream"+itoa(stream)),
					})
				}
				res.FamilyOrderSummaries=append(res.FamilyOrderSummaries,uplm1xFamilySummary(arm,family,order,res.IntegratedMetrics))
			}
		}
		for _,order:=range orders{res.OrderSummaries=append(res.OrderSummaries,uplm1xOrderSummary(arm,order,res.FamilyOrderSummaries))}
	}
	return res,nil
}
