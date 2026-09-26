package unitary

const UPLM2HDepthSchema="wingless.up-lm2h-combined-locality-stream-depth.v1"

type UPLM2HDepthSummary struct{
	Allocation string `json:"allocation"`
	TrainingSchedule string `json:"training_schedule"`
	Variant string `json:"variant"`
	Stream int `json:"stream"`
	FifthIntegratedMean float64 `json:"fifth_integrated_mean"`
	PriorIntegratedMean float64 `json:"prior_integrated_mean"`
	AllFamilyIntegratedMean float64 `json:"all_family_integrated_mean"`
	MinFamilyIntegratedMean float64 `json:"min_family_integrated_mean"`
	StructuralExactAllFamilies bool `json:"structural_exact_all_families"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM2HResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2GSeal string `json:"source_up_lm2g_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	Allocations int `json:"allocations"`
	TrainingSchedules int `json:"training_schedules"`
	Variants int `json:"variants"`
	Families int `json:"families"`
	StreamDepths []int `json:"stream_depths"`
	FifthAlwaysLast bool `json:"fifth_always_last"`
	StoreObserveGlobalAllVariants bool `json:"store_observe_global_all_variants"`
	EvaluationRetrainingUsed bool `json:"evaluation_retraining_used"`
	RecallCapChanged bool `json:"recall_cap_changed"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RouterModified bool `json:"router_modified"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	RouterMetrics []UPLM1PRouterMetric `json:"router_metrics"`
	Metrics []UPLM2EMetric `json:"metrics"`
	DepthSummaries []UPLM2HDepthSummary `json:"depth_summaries"`
}
func uplm2hDepthSummary(allocation,schedule,variant string,stream int,metrics []UPLM2EMetric)UPLM2HDepthSummary{
	sum,prior,fifth,min:=0.0,0.0,0.0,1.0
	count,priorN,maxRecall:=0,0,0
	exact:=true
	for _,x:=range metrics{
		if x.Allocation!=allocation||x.TrainingSchedule!=schedule||x.LocalityVariant!=variant||x.Stream!=stream{continue}
		acc:=x.Metric.Top1Accuracy;sum+=acc;count++
		if acc<min{min=acc}
		if x.Family=="fifth"{fifth=acc}else{prior+=acc;priorN++}
		if x.Metric.MaxRecallEntries>maxRecall{maxRecall=x.Metric.MaxRecallEntries}
		if !uplm1wStructuralExact(x.Metric){exact=false}
	}
	return UPLM2HDepthSummary{Allocation:allocation,TrainingSchedule:schedule,Variant:variant,Stream:stream,FifthIntegratedMean:fifth,PriorIntegratedMean:prior/float64(priorN),AllFamilyIntegratedMean:sum/float64(count),MinFamilyIntegratedMean:min,StructuralExactAllFamilies:exact,MaxRecallEntries:maxRecall}
}
func RunUPLM2H()(UPLM2HResult,error){
	start,baseTrain,_,paraTrain,_,thirdTrain,_,fourthTrain,_,err:=uplm1mStartModel()
	if err!=nil{return UPLM2HResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain};uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus();start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection();classifier:=uplm1pTrainClassifier(d);names:=uplm0gNames()
	res:=UPLM2HResult{Schema:UPLM2HDepthSchema,Experiment:"UP-LM2H-combined-locality-stream-depth",SourceUPLM2GSeal:"be861515e9020fe1be8d8e3a1193a77582e2621d",StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,Allocations:5,TrainingSchedules:2,Variants:2,Families:5,StreamDepths:[]int{1,4,8,16,48},FifthAlwaysLast:true,StoreObserveGlobalAllVariants:true,EvaluationRetrainingUsed:false,RecallCapChanged:false,RecurrentParametersTrained:false,RouterModified:false,AttentionUsed:false,FutureOracleUsed:false}
	res.RouterMetrics=append(res.RouterMetrics,
		uplm1pRouterEval(classifier,d,"base","heldout",names[4:6],[]string{"stores","observes","reports"}),
		uplm1pRouterEval(classifier,d,"paraphrase","heldout",names[4:6],[]string{"saves","sees","recalls"}),
		uplm1pRouterEval(classifier,d,"third","heldout",names[4:6],[]string{"archives","notices","recounts"}),
		uplm1pRouterEval(classifier,d,"fourth","heldout",names[4:6],[]string{"retains","inspects","states"}),
		uplm1pRouterEval(classifier,d,"fifth","heldout",names[4:6],uplm1pFifthVerbs),
	)
	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	familyNames:=[]string{"base","paraphrase","third","fourth","fifth"}
	allocations:=[]string{"equal_mass","fifth_1p125_mass","fifth_1p25_mass","fifth_1p375_mass","fifth_1p5_mass"}
	schedules:=[]string{"canonical_prior","balanced_prior"}
	variants:=[]string{"forward_tail","reverse_report_tail"}
	streams:=[]int{1,4,8,16,48}
	for _,allocation:=range allocations{for _,schedule:=range schedules{
		model:=uplm0oCloneModel(start);uplm2cTrain(model,allocation,schedule,families)
		for _,variant:=range variants{for _,family:=range familyNames{
			held:=uplm2gHeldout(family,variant)
			for _,stream:=range streams{res.Metrics=append(res.Metrics,UPLM2EMetric{Allocation:allocation,TrainingSchedule:schedule,LocalityVariant:variant,Family:family,Stream:stream,Metric:uplm1mEvaluate(model,classifier,d,held,stream,family+"_"+variant+"_stream"+itoa(stream))})}
		}}
		for _,variant:=range variants{for _,stream:=range streams{res.DepthSummaries=append(res.DepthSummaries,uplm2hDepthSummary(allocation,schedule,variant,stream,res.Metrics))}}
	}}
	return res,nil
}
