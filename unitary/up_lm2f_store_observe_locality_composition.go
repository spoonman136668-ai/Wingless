package unitary

const UPLM2FLocalityCompositionSchema="wingless.up-lm2f-store-observe-locality-composition.v1"

type UPLM2FResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2ESeal string `json:"source_up_lm2e_seal"`
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

func uplm2fCombinedExample(n,v,p int,family string)uplm0fExample{
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
	for i:=0;i<4;i++{storeInitial(i)}
	for i:=0;i<4;i++{observe(i)}
	for i:=0;i<4;i++{update(i);report(i,i==3)}
	s+="\\n"
	return uplm0fExample{text:s,targetPos:targets,updateCount:uc}
}
func uplm2fHeldout(family,variant string)[]uplm0fExample{
	if variant!="stores_observes_global"{return uplm2eHeldout(family,variant)}
	out:=[]uplm0fExample{}
	for n:=0;n<6;n++{for v:=0;v<6;v++{for p:=0;p<4;p++{
		if (n+2*v+p)%3!=2{continue}
		out=append(out,uplm2fCombinedExample(n,v,p,family))
	}}}
	return out
}
func RunUPLM2F()(UPLM2FResult,error){
	start,baseTrain,_,paraTrain,_,thirdTrain,_,fourthTrain,_,err:=uplm1mStartModel()
	if err!=nil{return UPLM2FResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus();start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection();classifier:=uplm1pTrainClassifier(d);names:=uplm0gNames()
	res:=UPLM2FResult{
		Schema:UPLM2FLocalityCompositionSchema,Experiment:"UP-LM2F-store-observe-locality-composition",SourceUPLM2ESeal:"a58f81e197cbd9f1928970a30fa59fb29a8e5b15",
		StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,
		Allocations:5,TrainingSchedules:2,LocalityVariants:4,Families:5,StreamDepths:[]int{1,4},FifthAlwaysLast:true,
		EvaluationRetrainingUsed:false,ContentChangedAcrossVariants:false,RecurrentParametersTrained:false,RecallCapChanged:false,RouterModified:false,AttentionUsed:false,FutureOracleUsed:false,
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
	variants:=[]string{"per_name","stores_global","observes_global","stores_observes_global"}
	for _,allocation:=range allocations{
		for _,schedule:=range schedules{
			model:=uplm0oCloneModel(start);uplm2cTrain(model,allocation,schedule,families)
			for _,variant:=range variants{
				for _,family:=range familyNames{
					held:=uplm2fHeldout(family,variant)
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
