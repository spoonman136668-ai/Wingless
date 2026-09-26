package unitary

import("math";"strings")

const UPLM2LBoundarySchema="wingless.up-lm2l-five-dependency-boundary.v1"

type UPLM2LMetric struct{
	Allocation string `json:"allocation"`
	TrainingSchedule string `json:"training_schedule"`
	Family string `json:"family"`
	IdentityRotation int `json:"identity_rotation"`
	DeferredReportOrder string `json:"deferred_report_order"`
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	RecallHitRate float64 `json:"recall_hit_rate"`
	ExpectedRecallHitRate float64 `json:"expected_recall_hit_rate"`
	DependentFirstByteAccuracy float64 `json:"dependent_first_byte_accuracy"`
	ReportSetExactAccuracy float64 `json:"report_set_exact_accuracy"`
	EventRoutingAccuracy float64 `json:"event_routing_accuracy"`
	StoreRoutingAccuracy float64 `json:"store_routing_accuracy"`
	ObserveRoutingAccuracy float64 `json:"observe_routing_accuracy"`
	ReportRoutingAccuracy float64 `json:"report_routing_accuracy"`
	MaxRecallEntries int `json:"max_recall_entries"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity float64 `json:"perplexity"`
}
type UPLM2LResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2KSeal string `json:"source_up_lm2k_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	EntityCount int `json:"entity_count"`
	DeferredLevels []int `json:"deferred_levels"`
	IdentityRotations []int `json:"identity_rotations"`
	DeferredReportOrders []string `json:"deferred_report_orders"`
	Allocations int `json:"allocations"`
	TrainingSchedules int `json:"training_schedules"`
	Families int `json:"families"`
	EventMultisetIdentical bool `json:"event_multiset_identical"`
	CapacityChanged bool `json:"capacity_changed"`
	AdaptiveChunkingUsed bool `json:"adaptive_chunking_used"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	Metrics []UPLM2LMetric `json:"metrics"`
}
func uplm2lExpected(d int)float64{return float64(24-(d-4))/24.0}
func uplm2lText(family string,deferN,rot int,order string)(string,map[int]bool){
	base:=uplm2iNames();values:=uplm0gValues();storeVerb,observeVerb,reportVerb:=uplm2dVerbs(family)
	nameAt:=func(pos int)string{return base[(pos+rot)%24]}
	valueAt:=func(pos int)string{return values[((pos+rot)%24)%len(values)]}
	s:="";targets:=map[int]bool{}
	reportOne:=func(pos int){s+=nameAt(pos)+" "+reportVerb+" ";targets[len(s)]=true;s+=valueAt(pos)+". "}
	for i:=0;i<12;i++{s+=nameAt(i)+" "+storeVerb+" "+valueAt(i)+". "}
	for i:=0;i<12;i++{s+=nameAt(i)+" "+observeVerb+" "+values[((i+rot)+3)%len(values)]+". "}
	early:=12-deferN
	for i:=0;i<early;i++{reportOne(i)}
	for i:=12;i<24;i++{s+=nameAt(i)+" "+storeVerb+" "+valueAt(i)+". "}
	for i:=12;i<24;i++{s+=nameAt(i)+" "+observeVerb+" "+values[((i+rot)+3)%len(values)]+". "}
	if order=="reverse"{for i:=11;i>=early;i--{reportOne(i)}}else{for i:=early;i<12;i++{reportOne(i)}}
	for i:=12;i<24;i++{reportOne(i)}
	return strings.TrimSpace(s)+"\n",targets
}
func uplm2lEval(model *uplm0aModel,classifier *uplm0jClassifier,d [64]float64,allocation,ts,family string,level,rot int,order string)UPLM2LMetric{
	s,targets:=uplm2lText(family,level,rot,order);var h [64]float64;recall:=newUPLM0CRecall()
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
					switch trueClass{case uplm0jStore:storeTotal++;if routeClass==trueClass{storeHits++};case uplm0jObserve:observeTotal++;if routeClass==trueClass{observeHits++};case uplm0jReport:reportTotal++;trueQueryName=fields[0];if routeClass==trueClass{reportHits++}}
					if routeClass==uplm0jReport{queryName=fields[0];reportOverridePending=true}
				}
			}
		}
		targetByte:=s[t+1];target:=model.index[int(targetByte)];okTarget:=target>=0;if !okTarget{target=0}
		p:=model.probs(h);pred:=uplm0aArgmax(p);prob:=p[target];isDep:=targets[t+1]
		if isDep{recallTotal++;if value,ok:=recall.values[trueQueryName];ok&&len(value)>0{recallHits++}}
		if reportOverridePending&&queryName!=""{if value,ok:=recall.values[queryName];ok&&len(value)>0{memByte:=value[0];mi:=model.index[int(memByte)];if mi>=0{pred=mi};if memByte==targetByte{prob=1}else{prob=1e-12}};reportOverridePending=false}
		if prob<1e-12{prob=1e-12};total++;nll-=math.Log(prob);if okTarget&&pred==target{hits++}
		if isDep{depTotal++;if okTarget&&pred==target{depHits++}else{allReports=false}}
		if b=='.'{clean:=strings.TrimSuffix(strings.TrimSpace(clause),".");fields:=strings.Fields(clean);if len(fields)>=3&&trueClass>=0&&routeKnown&&routeClass==uplm0jStore{recall.write(fields[0],fields[2])};clause="";routeKnown=false;routeClass=-1;trueClass=-1;queryName="";trueQueryName="";reportOverridePending=false}
		if len(recall.order)>maxEntries{maxEntries=len(recall.order)}
	}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)};exact:=0.0;if allReports&&depTotal==24{exact=1}
	return UPLM2LMetric{Allocation:allocation,TrainingSchedule:ts,Family:family,IdentityRotation:rot,DeferredReportOrder:order,DeferredFirstChunkReports:level,RecallHitRate:rate(recallHits,recallTotal),ExpectedRecallHitRate:uplm2lExpected(level),DependentFirstByteAccuracy:rate(depHits,depTotal),ReportSetExactAccuracy:exact,EventRoutingAccuracy:rate(eventHits,eventTotal),StoreRoutingAccuracy:rate(storeHits,storeTotal),ObserveRoutingAccuracy:rate(observeHits,observeTotal),ReportRoutingAccuracy:rate(reportHits,reportTotal),MaxRecallEntries:maxEntries,Top1Accuracy:rate(hits,total),Perplexity:math.Exp(nll/float64(total))}
}
func RunUPLM2L()(UPLM2LResult,error){
	start,baseTrain,_,paraTrain,_,thirdTrain,_,fourthTrain,_,err:=uplm1mStartModel();if err!=nil{return UPLM2LResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain};uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus();start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection();classifier:=uplm1pTrainClassifier(d)
	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain};familyNames:=[]string{"base","paraphrase","third","fourth","fifth"}
	allocs:=[]string{"equal_mass","fifth_1p125_mass","fifth_1p25_mass","fifth_1p375_mass","fifth_1p5_mass"};schedules:=[]string{"canonical_prior","balanced_prior"};levels:=[]int{4,5,6};rots:=[]int{0,7};orders:=[]string{"forward","reverse"}
	res:=UPLM2LResult{Schema:UPLM2LBoundarySchema,Experiment:"UP-LM2L-five-dependency-boundary",SourceUPLM2KSeal:"bc6469d94e4b71d79ceaa85a55aeb6ab905b8464",ExactRecallCap:16,EntityCount:24,DeferredLevels:levels,IdentityRotations:rots,DeferredReportOrders:orders,Allocations:5,TrainingSchedules:2,Families:5,EventMultisetIdentical:true,CapacityChanged:false,AdaptiveChunkingUsed:false,ExtraTrainingUsed:false}
	for _,a:=range allocs{for _,ts:=range schedules{model:=uplm0oCloneModel(start);uplm2cTrain(model,a,ts,families);for _,f:=range familyNames{for _,rot:=range rots{for _,ord:=range orders{for _,level:=range levels{res.Metrics=append(res.Metrics,uplm2lEval(model,classifier,d,a,ts,f,level,rot,ord))}}}}}}
	return res,nil
}
