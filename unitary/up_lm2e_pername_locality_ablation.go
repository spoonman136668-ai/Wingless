package unitary

const UPLM2ELocalitySchema="wingless.up-lm2e-pername-locality-ablation.v1"

type UPLM2EMetric struct{
	Allocation string `json:"allocation"`
	TrainingSchedule string `json:"training_schedule"`
	LocalityVariant string `json:"locality_variant"`
	Family string `json:"family"`
	Stream int `json:"stream"`
	Metric UPLM0JMetric `json:"metric"`
}
type UPLM2EFamilySummary struct{
	Allocation string `json:"allocation"`
	TrainingSchedule string `json:"training_schedule"`
	LocalityVariant string `json:"locality_variant"`
	Family string `json:"family"`
	IntegratedMeanTop1Accuracy float64 `json:"integrated_mean_top1_accuracy"`
	IntegratedMinTop1Accuracy float64 `json:"integrated_min_top1_accuracy"`
	IntegratedMeanPerplexity float64 `json:"integrated_mean_perplexity"`
	StructuralExactAllCells bool `json:"structural_exact_all_cells"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM2EVariantSummary struct{
	Allocation string `json:"allocation"`
	TrainingSchedule string `json:"training_schedule"`
	LocalityVariant string `json:"locality_variant"`
	FifthIntegratedMean float64 `json:"fifth_integrated_mean"`
	PriorIntegratedMean float64 `json:"prior_integrated_mean"`
	AllFamilyIntegratedMean float64 `json:"all_family_integrated_mean"`
	MinFamilyIntegratedMean float64 `json:"min_family_integrated_mean"`
	StructuralExactAllFamilies bool `json:"structural_exact_all_families"`
}
type UPLM2EResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2DSeal string `json:"source_up_lm2d_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	Allocations int `json:"allocations"`
	TrainingSchedules int `json:"training_schedules"`
	LocalityVariants int `json:"locality_variants"`
	Families int `json:"families"`
	StreamDepths []int `json:"stream_depths"`
	FifthAlwaysLast bool `json:"fifth_always_last"`
	EvaluationRetrainingUsed bool `json:"evaluation_retraining_used"`
	ContentChangedAcrossVariants bool `json:"content_changed_across_variants"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RecallCapChanged bool `json:"recall_cap_changed"`
	RouterModified bool `json:"router_modified"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	RouterMetrics []UPLM1PRouterMetric `json:"router_metrics"`
	Metrics []UPLM2EMetric `json:"metrics"`
	FamilySummaries []UPLM2EFamilySummary `json:"family_summaries"`
	VariantSummaries []UPLM2EVariantSummary `json:"variant_summaries"`
}

func uplm2eExample(n,v,p int,family,variant string)uplm0fExample{
	names:=uplm0gNames();values:=uplm0gValues()
	ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
	initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
	latest:=initial;uc:=uplm0eUpdateCount(p)
	storeVerb,observeVerb,reportVerb:=uplm2dVerbs(family)
	s:="";var targets [4]int
	storeInitial:=func(i int){s+=ns[i]+" "+storeVerb+" "+initial[i]+". "}
	observe:=func(i int){obs:=values[(v+i+3)%6];s+=ns[(i+p)%4]+" "+observeVerb+" "+obs+". "}
	update:=func(i int){if i>=uc{return};latest[i]=values[(v+i+2)%6];s+=ns[i]+" "+storeVerb+" "+latest[i]+". "}
	report:=func(i int,last bool){s+=ns[i]+" "+reportVerb+" ";targets[i]=len(s);s+=latest[i]+".";if !last{s+=" "}}
	switch variant{
	case "per_name":
		for i:=0;i<4;i++{storeInitial(i);observe(i);update(i);report(i,i==3)}
	case "stores_global":
		for i:=0;i<4;i++{storeInitial(i)}
		for i:=0;i<4;i++{observe(i);update(i);report(i,i==3)}
	case "observes_global":
		for i:=0;i<4;i++{observe(i)}
		for i:=0;i<4;i++{storeInitial(i);update(i);report(i,i==3)}
	default: // reports_global
		for i:=0;i<4;i++{storeInitial(i);observe(i);update(i)}
		for i:=0;i<4;i++{report(i,i==3)}
	}
	s+="\n"
	return uplm0fExample{text:s,targetPos:targets,updateCount:uc}
}
func uplm2eHeldout(family,variant string)[]uplm0fExample{
	out:=[]uplm0fExample{}
	for n:=0;n<6;n++{for v:=0;v<6;v++{for p:=0;p<4;p++{
		if (n+2*v+p)%3!=2{continue}
		out=append(out,uplm2eExample(n,v,p,family,variant))
	}}}
	return out
}
func uplm2eFamilySummary(allocation,schedule,variant,family string,metrics []UPLM2EMetric)UPLM2EFamilySummary{
	sumAcc,sumPpl:=0.0,0.0;minAcc:=1.0;count,maxRecall:=0,0;exact:=true
	for _,x:=range metrics{
		if x.Allocation!=allocation||x.TrainingSchedule!=schedule||x.LocalityVariant!=variant||x.Family!=family{continue}
		count++;sumAcc+=x.Metric.Top1Accuracy;sumPpl+=x.Metric.Perplexity
		if x.Metric.Top1Accuracy<minAcc{minAcc=x.Metric.Top1Accuracy}
		if x.Metric.MaxRecallEntries>maxRecall{maxRecall=x.Metric.MaxRecallEntries}
		if !uplm1wStructuralExact(x.Metric){exact=false}
	}
	return UPLM2EFamilySummary{Allocation:allocation,TrainingSchedule:schedule,LocalityVariant:variant,Family:family,IntegratedMeanTop1Accuracy:sumAcc/float64(count),IntegratedMinTop1Accuracy:minAcc,IntegratedMeanPerplexity:sumPpl/float64(count),StructuralExactAllCells:exact,MaxRecallEntries:maxRecall}
}
func uplm2eVariantSummary(allocation,schedule,variant string,s []UPLM2EFamilySummary)UPLM2EVariantSummary{
	sum,prior,fifth,min:=0.0,0.0,0.0,1.0;count,priorN:=0,0;exact:=true
	for _,x:=range s{
		if x.Allocation!=allocation||x.TrainingSchedule!=schedule||x.LocalityVariant!=variant{continue}
		count++;sum+=x.IntegratedMeanTop1Accuracy
		if x.IntegratedMeanTop1Accuracy<min{min=x.IntegratedMeanTop1Accuracy}
		if x.Family=="fifth"{fifth=x.IntegratedMeanTop1Accuracy}else{prior+=x.IntegratedMeanTop1Accuracy;priorN++}
		if !x.StructuralExactAllCells{exact=false}
	}
	return UPLM2EVariantSummary{Allocation:allocation,TrainingSchedule:schedule,LocalityVariant:variant,FifthIntegratedMean:fifth,PriorIntegratedMean:prior/float64(priorN),AllFamilyIntegratedMean:sum/float64(count),MinFamilyIntegratedMean:min,StructuralExactAllFamilies:exact}
}
func RunUPLM2E()(UPLM2EResult,error){
	start,baseTrain,_,paraTrain,_,thirdTrain,_,fourthTrain,_,err:=uplm1mStartModel()
	if err!=nil{return UPLM2EResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus();start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection();classifier:=uplm1pTrainClassifier(d);names:=uplm0gNames()
	res:=UPLM2EResult{
		Schema:UPLM2ELocalitySchema,Experiment:"UP-LM2E-pername-locality-ablation",SourceUPLM2DSeal:"12374b0c3bfc3227b9e00b335ca4c5ca6af307fb",
		StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,
		Allocations:5,TrainingSchedules:2,LocalityVariants:4,Families:5,StreamDepths:[]int{1,4},FifthAlwaysLast:true,
		EvaluationRetrainingUsed:false,ContentChangedAcrossVariants:false,
		RecurrentParametersTrained:false,RecallCapChanged:false,RouterModified:false,AttentionUsed:false,FutureOracleUsed:false,
	}
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
	variants:=[]string{"per_name","stores_global","observes_global","reports_global"}
	for _,allocation:=range allocations{
		for _,schedule:=range schedules{
			model:=uplm0oCloneModel(start);uplm2cTrain(model,allocation,schedule,families)
			for _,variant:=range variants{
				for _,family:=range familyNames{
					held:=uplm2eHeldout(family,variant)
					for _,stream:=range []int{1,4}{
						res.Metrics=append(res.Metrics,UPLM2EMetric{Allocation:allocation,TrainingSchedule:schedule,LocalityVariant:variant,Family:family,Stream:stream,Metric:uplm1mEvaluate(model,classifier,d,held,stream,family+"_"+variant+"_stream"+itoa(stream))})
					}
					res.FamilySummaries=append(res.FamilySummaries,uplm2eFamilySummary(allocation,schedule,variant,family,res.Metrics))
				}
				res.VariantSummaries=append(res.VariantSummaries,uplm2eVariantSummary(allocation,schedule,variant,res.FamilySummaries))
			}
		}
	}
	return res,nil
}
