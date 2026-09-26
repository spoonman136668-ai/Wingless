package unitary

import (
	"math"
	"strings"
)

const UPLM2IConcurrencySchema="wingless.up-lm2i-entity-concurrency-capacity.v1"

type UPLM2IMetric struct{
	Allocation string `json:"allocation"`
	TrainingSchedule string `json:"training_schedule"`
	Family string `json:"family"`
	ReportOrder string `json:"report_order"`
	Entities int `json:"entities"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity float64 `json:"perplexity"`
	DependentFirstByteAccuracy float64 `json:"dependent_first_byte_accuracy"`
	ReportSetExactAccuracy float64 `json:"report_set_exact_accuracy"`
	RecallHitRate float64 `json:"recall_hit_rate"`
	ExpectedCapacityHitRate float64 `json:"expected_capacity_hit_rate"`
	EventRoutingAccuracy float64 `json:"event_routing_accuracy"`
	StoreRoutingAccuracy float64 `json:"store_routing_accuracy"`
	ObserveRoutingAccuracy float64 `json:"observe_routing_accuracy"`
	ReportRoutingAccuracy float64 `json:"report_routing_accuracy"`
	MaxRecallEntries int `json:"max_recall_entries"`
}
type UPLM2IResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2HSeal string `json:"source_up_lm2h_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	AdaptationEpochs int `json:"adaptation_epochs"`
	TotalMassPerExample float64 `json:"total_mass_per_example"`
	Allocations int `json:"allocations"`
	TrainingSchedules int `json:"training_schedules"`
	Families int `json:"families"`
	ReportOrders int `json:"report_orders"`
	ConcurrencyLevels []int `json:"concurrency_levels"`
	EntityNameCount int `json:"entity_name_count"`
	EntityNamesUseExistingAlphabet bool `json:"entity_names_use_existing_alphabet"`
	RecallCapChanged bool `json:"recall_cap_changed"`
	RecurrentParametersTrained bool `json:"recurrent_parameters_trained"`
	RouterModified bool `json:"router_modified"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Metrics []UPLM2IMetric `json:"metrics"`
}
func uplm2iNames()[]string{
	return []string{"ada","ben","cy","dee","eli","fay","adaben","adacy","adadee","adaeli","adafay","benada","bency","bendee","beneli","benfay","cyada","cyben","cydee","cyeli","cyfay","deeada","deeben","deecy"}
}
func uplm2iExpected(n int)float64{if n<=16{return 1};return 16.0/float64(n)}
func uplm2iText(family string,n int,order string)(string,map[int]bool){
	names:=uplm2iNames()[:n];values:=uplm0gValues();storeVerb,observeVerb,reportVerb:=uplm2dVerbs(family)
	s:="";targets:=map[int]bool{}
	for i,name:=range names{s+=name+" "+storeVerb+" "+values[i%len(values)]+". "}
	for i,name:=range names{s+=name+" "+observeVerb+" "+values[(i+3)%len(values)]+". "}
	reportOne:=func(i int,last bool){s+=names[i]+" "+reportVerb+" ";targets[len(s)]=true;s+=values[i%len(values)]+".";if !last{s+=" "}}
	if order=="reverse_report"{for j:=n-1;j>=0;j--{reportOne(j,j==0)}}else{for i:=0;i<n;i++{reportOne(i,i==n-1)}}
	s+="\n";return s,targets
}
func uplm2iEval(model *uplm0aModel,classifier *uplm0jClassifier,d [64]float64,allocation,schedule,family,order string,n int)UPLM2IMetric{
	s,targets:=uplm2iText(family,n,order);var h [64]float64;recall:=newUPLM0CRecall()
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
		p:=model.probs(h);pred:=uplm0aArgmax(p);prob:=p[target]
		isDep:=targets[t+1]
		if isDep{
			recallTotal++
			if value,ok:=recall.values[trueQueryName];ok&&len(value)>0{recallHits++}
		}
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
	exact:=0.0;if allReports&&depTotal==n{exact=1}
	return UPLM2IMetric{Allocation:allocation,TrainingSchedule:schedule,Family:family,ReportOrder:order,Entities:n,Top1Accuracy:rate(hits,total),Perplexity:math.Exp(nll/float64(total)),DependentFirstByteAccuracy:rate(depHits,depTotal),ReportSetExactAccuracy:exact,RecallHitRate:rate(recallHits,recallTotal),ExpectedCapacityHitRate:uplm2iExpected(n),EventRoutingAccuracy:rate(eventHits,eventTotal),StoreRoutingAccuracy:rate(storeHits,storeTotal),ObserveRoutingAccuracy:rate(observeHits,observeTotal),ReportRoutingAccuracy:rate(reportHits,reportTotal),MaxRecallEntries:maxEntries}
}
func RunUPLM2I()(UPLM2IResult,error){
	start,baseTrain,_,paraTrain,_,thirdTrain,_,fourthTrain,_,err:=uplm1mStartModel();if err!=nil{return UPLM2IResult{},err}
	four:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain};uplm1nTrainArm(start,"rotating_palindromic_split",four)
	fifthTrain,fifthHeld:=uplm1pCorpus();start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)
	d:=uplm1lStoreDirection();classifier:=uplm1pTrainClassifier(d)
	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	familyNames:=[]string{"base","paraphrase","third","fourth","fifth"}
	allocations:=[]string{"equal_mass","fifth_1p125_mass","fifth_1p25_mass","fifth_1p375_mass","fifth_1p5_mass"}
	schedules:=[]string{"canonical_prior","balanced_prior"};orders:=[]string{"forward_report","reverse_report"};levels:=[]int{8,12,16,20,24}
	res:=UPLM2IResult{Schema:UPLM2IConcurrencySchema,Experiment:"UP-LM2I-entity-concurrency-capacity",SourceUPLM2HSeal:"0e79f05ed3bc3ac46f290de9bdfc09ff0a0f0239",StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,TotalMassPerExample:0.40,Allocations:5,TrainingSchedules:2,Families:5,ReportOrders:2,ConcurrencyLevels:append([]int(nil),levels...),EntityNameCount:24,EntityNamesUseExistingAlphabet:true,RecallCapChanged:false,RecurrentParametersTrained:false,RouterModified:false,AttentionUsed:false,FutureOracleUsed:false}
	for _,allocation:=range allocations{for _,schedule:=range schedules{
		model:=uplm0oCloneModel(start);uplm2cTrain(model,allocation,schedule,families)
		for _,family:=range familyNames{for _,order:=range orders{for _,n:=range levels{res.Metrics=append(res.Metrics,uplm2iEval(model,classifier,d,allocation,schedule,family,order,n))}}}
	}}
	return res,nil
}
