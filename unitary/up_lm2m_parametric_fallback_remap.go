package unitary

import("math";"strings")

const UPLM2MFallbackSchema="wingless.up-lm2m-parametric-fallback-remap.v1"

type UPLM2MMetric struct{
	TrainingSchedule string `json:"training_schedule"`
	Family string `json:"family"`
	IdentityRotation int `json:"identity_rotation"`
	DeferredReportOrder string `json:"deferred_report_order"`
	ValueShift int `json:"value_shift"`
	RecallHitRate float64 `json:"recall_hit_rate"`
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
type UPLM2MResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2LSeal string `json:"source_up_lm2l_seal"`
	Allocation string `json:"allocation"`
	ExactRecallCap int `json:"exact_recall_cap"`
	EntityCount int `json:"entity_count"`
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	ValueShifts []int `json:"value_shifts"`
	IdentityRotations []int `json:"identity_rotations"`
	DeferredReportOrders []string `json:"deferred_report_orders"`
	TrainingSchedules int `json:"training_schedules"`
	Families int `json:"families"`
	CapacityChanged bool `json:"capacity_changed"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	RouterModified bool `json:"router_modified"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Metrics []UPLM2MMetric `json:"metrics"`
}
func uplm2mText(family string,rot,order,shift int)(string,map[int]bool){
	base:=uplm2iNames();values:=uplm0gValues();storeVerb,observeVerb,reportVerb:=uplm2dVerbs(family)
	nameAt:=func(pos int)string{return base[(pos+rot)%24]}
	valueAt:=func(pos int)string{return values[((pos+rot)+shift)%len(values)]}
	observeAt:=func(pos int)string{return values[((pos+rot)+shift+3)%len(values)]}
	s:="";targets:=map[int]bool{};deferN:=5
	reportOne:=func(pos int){s+=nameAt(pos)+" "+reportVerb+" ";targets[len(s)]=true;s+=valueAt(pos)+". "}
	for i:=0;i<12;i++{s+=nameAt(i)+" "+storeVerb+" "+valueAt(i)+". "}
	for i:=0;i<12;i++{s+=nameAt(i)+" "+observeVerb+" "+observeAt(i)+". "}
	early:=12-deferN
	for i:=0;i<early;i++{reportOne(i)}
	for i:=12;i<24;i++{s+=nameAt(i)+" "+storeVerb+" "+valueAt(i)+". "}
	for i:=12;i<24;i++{s+=nameAt(i)+" "+observeVerb+" "+observeAt(i)+". "}
	if order=="reverse"{for i:=11;i>=early;i--{reportOne(i)}}else{for i:=early;i<12;i++{reportOne(i)}}
	for i:=12;i<24;i++{reportOne(i)}
	return strings.TrimSpace(s)+"\n",targets
}
func uplm2mEval(model *uplm0aModel,classifier *uplm0jClassifier,d [64]float64,ts,family string,rot int,order string,shift int)UPLM2MMetric{
	s,targets:=uplm2mText(family,rot,order,shift);var h [64]float64;recall:=newUPLM0CRecall()
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
		if reportOverridePending&&queryName!=""{
			if value,ok:=recall.values[queryName];ok&&len(value)>0{
				memByte:=value[0];mi:=model.index[int(memByte)];if mi>=0{pred=mi};if memByte==targetByte{prob=1}else{prob=1e-12}
			}
			reportOverridePending=false
		}
		if prob<1e-12{prob=1e-12};total++;nll-=math.Log(prob);if okTarget&&pred==target{hits++}
		if isDep{depTotal++;if okTarget&&pred==target{depHits++}else{allReports=false}}
		if b=='.'{clean:=strings.TrimSuffix(strings.TrimSpace(clause),".");fields:=strings.Fields(clean);if len(fields)>=3&&trueClass>=0&&routeKnown&&routeClass==uplm0jStore{recall.write(fields[0],fields[2])};clause="";routeKnown=false;routeClass=-1;trueClass=-1;queryName="";trueQueryName="";reportOverridePending=false}
		if len(recall.order)>maxEntries{maxEntries=len(recall.order)}
	}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)};exact:=0.0;if allReports&&depTotal==24{exact=1}
	return UPLM2MMetric{TrainingSchedule:ts,Family:family,IdentityRotation:rot,DeferredReportOrder:order,ValueShift:shift,RecallHitRate:rate(recallHits,recallTotal),DependentFirstByteAccuracy:rate(depHits,depTotal),ReportSetExactAccuracy:exact,EventRoutingAccuracy:rate(eventHits,eventTotal),StoreRoutingAccuracy:rate(storeHits,storeTotal),ObserveRoutingAccuracy:rate(observeHits,observeTotal),ReportRoutingAccuracy:rate(reportHits,reportTotal),MaxRecallEntries:maxEntries,Top1Accuracy:rate(hits,total),Perplexity:math.Exp(nll/float64(total))}
}
func RunUPLM2M()(UPLM2MResult,error){
	start,baseTrain,_,paraTrain,_,thirdTrain,_,fourthTrain,_,err:=uplm1mStartModel();if err!=nil{return UPLM2MResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain};uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus();start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection();classifier:=uplm1pTrainClassifier(d)
	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain};familyNames:=[]string{"base","paraphrase","third","fourth","fifth"}
	schedules:=[]string{"canonical_prior","balanced_prior"};rots:=[]int{0,7};orders:=[]string{"forward","reverse"};shifts:=[]int{0,1,2,3}
	res:=UPLM2MResult{Schema:UPLM2MFallbackSchema,Experiment:"UP-LM2M-parametric-fallback-remap",SourceUPLM2LSeal:"c59dd2bce72ab784dec39ab356c00ccb1d379ee7",Allocation:"equal_mass",ExactRecallCap:16,EntityCount:24,DeferredFirstChunkReports:5,ValueShifts:shifts,IdentityRotations:rots,DeferredReportOrders:orders,TrainingSchedules:2,Families:5,CapacityChanged:false,ExtraTrainingUsed:false,RouterModified:false,AttentionUsed:false,FutureOracleUsed:false}
	for _,ts:=range schedules{
		model:=uplm0oCloneModel(start);uplm2cTrain(model,"equal_mass",ts,families)
		for _,f:=range familyNames{for _,rot:=range rots{for _,ord:=range orders{for _,shift:=range shifts{res.Metrics=append(res.Metrics,uplm2mEval(model,classifier,d,ts,f,rot,ord,shift))}}}}
	}
	return res,nil
}
