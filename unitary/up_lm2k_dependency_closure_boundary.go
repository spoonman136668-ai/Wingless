package unitary

import("math";"strings")

const UPLM2KClosureSchema="wingless.up-lm2k-dependency-closure-boundary.v1"

type UPLM2KMetric struct{
	Allocation string `json:"allocation"`
	TrainingSchedule string `json:"training_schedule"`
	Family string `json:"family"`
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity float64 `json:"perplexity"`
	DependentFirstByteAccuracy float64 `json:"dependent_first_byte_accuracy"`
	ReportSetExactAccuracy float64 `json:"report_set_exact_accuracy"`
	RecallHitRate float64 `json:"recall_hit_rate"`
	ExpectedRecallHitRate float64 `json:"expected_recall_hit_rate"`
	EventRoutingAccuracy float64 `json:"event_routing_accuracy"`
	StoreRoutingAccuracy float64 `json:"store_routing_accuracy"`
	ObserveRoutingAccuracy float64 `json:"observe_routing_accuracy"`
	ReportRoutingAccuracy float64 `json:"report_routing_accuracy"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM2KResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2JSeal string `json:"source_up_lm2j_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	EntityCount int `json:"entity_count"`
	ChunkSizes []int `json:"chunk_sizes"`
	DeferredLevels []int `json:"deferred_levels"`
	Allocations int `json:"allocations"`
	TrainingSchedules int `json:"training_schedules"`
	Families int `json:"families"`
	EventMultisetIdentical bool `json:"event_multiset_identical"`
	CapacityChanged bool `json:"capacity_changed"`
	AdaptiveChunkingUsed bool `json:"adaptive_chunking_used"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	RouterModified bool `json:"router_modified"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Metrics []UPLM2KMetric `json:"metrics"`
}
func uplm2kExpected(d int)float64{if d<=4{return 1};return float64(24-(d-4))/24.0}
func uplm2kText(family string,deferN int)(string,map[int]bool){
	names:=uplm2iNames();values:=uplm0gValues();storeVerb,observeVerb,reportVerb:=uplm2dVerbs(family)
	s:="";targets:=map[int]bool{}
	reportOne:=func(i int){s+=names[i]+" "+reportVerb+" ";targets[len(s)]=true;s+=values[i%len(values)]+". "}
	for i:=0;i<12;i++{s+=names[i]+" "+storeVerb+" "+values[i%len(values)]+". "}
	for i:=0;i<12;i++{s+=names[i]+" "+observeVerb+" "+values[(i+3)%len(values)]+". "}
	early:=12-deferN
	for i:=0;i<early;i++{reportOne(i)}
	for i:=12;i<24;i++{s+=names[i]+" "+storeVerb+" "+values[i%len(values)]+". "}
	for i:=12;i<24;i++{s+=names[i]+" "+observeVerb+" "+values[(i+3)%len(values)]+". "}
	for i:=early;i<12;i++{reportOne(i)}
	for i:=12;i<24;i++{reportOne(i)}
	s=strings.TrimSpace(s)+"\n";return s,targets
}
func uplm2kEval(model *uplm0aModel,classifier *uplm0jClassifier,d [64]float64,allocation,trainSchedule,family string,deferN int)UPLM2KMetric{
	s,targets:=uplm2kText(family,deferN);var h [64]float64;recall:=newUPLM0CRecall()
	hits,total,depHits,depTotal:=0,0,0,0;nll:=0.0;maxEntries:=0
	eventHits,eventTotal:=0,0;storeHits,storeTotal,observeHits,observeTotal,reportHits,reportTotal:=0,0,0,0,0,0
	recallHits,recallTotal:=0,0;allReports:=true
	clause:="";routeKnown:=false;routeClass,trueClass:=-1,-1;queryName:="";trueQueryName:="";reportOverridePending:=false
	for t:=0;t<len(s)-1;t++{
		b:=s[t];h=uplm0aStep(h,b);clause+=string(b);trim:=strings.TrimSpace(clause)
		if b==' '&&!routeKnown{
			fields:=strings.Fields(trim)
			if len(fields)==2{
				trueClass=uplm0nClassForVerb(fields[1])
				if trueClass>=0{
					routeClass=uplm1mPredict(classifier,fields[0],fields[1],d);routeKnown=true;eventTotal++
					if routeClass==trueClass{eventHits++}
					switch trueClass{
					case uplm0jStore:storeTotal++;if routeClass==trueClass{storeHits++}
					case uplm0jObserve:observeTotal++;if routeClass==trueClass{observeHits++}
					case uplm0jReport:reportTotal++;trueQueryName=fields[0];if routeClass==trueClass{reportHits++}
					}
					if routeClass==uplm0jReport{queryName=fields[0];reportOverridePending=true}
				}
			}
		}
		targetByte:=s[t+1];target:=model.index[int(targetByte)];okTarget:=target>=0;if !okTarget{target=0}
		p:=model.probs(h);pred:=uplm0aArgmax(p);prob:=p[target];isDep:=targets[t+1]
		if isDep{recallTotal++;if value,ok:=recall.values[trueQueryName];ok&&len(value)>0{recallHits++}}
		if reportOverridePending&&queryName!=""{
			if value,ok:=recall.values[queryName];ok&&len(value)>0{
				memByte:=value[0];mi:=model.index[int(memByte)];if mi>=0{pred=mi};if memByte==targetByte{prob=1}else{prob=1e-12}
			}
			reportOverridePending=false
		}
		if prob<1e-12{prob=1e-12};total++;nll-=math.Log(prob);if okTarget&&pred==target{hits++}
		if isDep{depTotal++;if okTarget&&pred==target{depHits++}else{allReports=false}}
		if b=='.'{
			clean:=strings.TrimSuffix(strings.TrimSpace(clause),".");fields:=strings.Fields(clean)
			if len(fields)>=3&&trueClass>=0&&routeKnown&&routeClass==uplm0jStore{recall.write(fields[0],fields[2])}
			clause="";routeKnown=false;routeClass=-1;trueClass=-1;queryName="";trueQueryName="";reportOverridePending=false
		}
		if len(recall.order)>maxEntries{maxEntries=len(recall.order)}
	}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	exact:=0.0;if allReports&&depTotal==24{exact=1}
	return UPLM2KMetric{Allocation:allocation,TrainingSchedule:trainSchedule,Family:family,DeferredFirstChunkReports:deferN,Top1Accuracy:rate(hits,total),Perplexity:math.Exp(nll/float64(total)),DependentFirstByteAccuracy:rate(depHits,depTotal),ReportSetExactAccuracy:exact,RecallHitRate:rate(recallHits,recallTotal),ExpectedRecallHitRate:uplm2kExpected(deferN),EventRoutingAccuracy:rate(eventHits,eventTotal),StoreRoutingAccuracy:rate(storeHits,storeTotal),ObserveRoutingAccuracy:rate(observeHits,observeTotal),ReportRoutingAccuracy:rate(reportHits,reportTotal),MaxRecallEntries:maxEntries}
}
func RunUPLM2K()(UPLM2KResult,error){
	start,baseTrain,_,paraTrain,_,thirdTrain,_,fourthTrain,_,err:=uplm1mStartModel();if err!=nil{return UPLM2KResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain};uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus();start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection();classifier:=uplm1pTrainClassifier(d)
	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	familyNames:=[]string{"base","paraphrase","third","fourth","fifth"}
	allocations:=[]string{"equal_mass","fifth_1p125_mass","fifth_1p25_mass","fifth_1p375_mass","fifth_1p5_mass"}
	trainingSchedules:=[]string{"canonical_prior","balanced_prior"};levels:=[]int{0,2,4,6,8,12}
	res:=UPLM2KResult{Schema:UPLM2KClosureSchema,Experiment:"UP-LM2K-dependency-closure-boundary",SourceUPLM2JSeal:"f7d186efd394989b9cc8b98a06af0cc394b0e49b",StateDimension:64,ExactRecallCap:16,EntityCount:24,ChunkSizes:[]int{12,12},DeferredLevels:append([]int(nil),levels...),Allocations:5,TrainingSchedules:2,Families:5,EventMultisetIdentical:true,CapacityChanged:false,AdaptiveChunkingUsed:false,ExtraTrainingUsed:false,RouterModified:false,AttentionUsed:false,FutureOracleUsed:false}
	for _,allocation:=range allocations{for _,ts:=range trainingSchedules{
		model:=uplm0oCloneModel(start);uplm2cTrain(model,allocation,ts,families)
		for _,family:=range familyNames{for _,level:=range levels{res.Metrics=append(res.Metrics,uplm2kEval(model,classifier,d,allocation,ts,family,level))}}
	}}
	return res,nil
}
