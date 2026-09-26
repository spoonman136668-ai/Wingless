package unitary

import "fmt"

const UP130CEdgeSchema="wingless.up130c-generation-edge-causality.v1"

type UP130CPoint struct{
	Arm string `json:"arm"`
	BoundaryBetweenSightings bool `json:"boundary_between_sightings"`
	PreOverflowTargetSightings int `json:"pre_overflow_target_sightings"`
	PreOverflowWrites int `json:"pre_overflow_writes"`
	PanicRate float64 `json:"panic_rate"`
	PanicMessage string `json:"panic_message"`
	PostFirstWaveEvictions int `json:"post_first_wave_evictions"`
	ValidAdmissionRate float64 `json:"valid_admission_rate"`
	ValidAccuracy float64 `json:"valid_accuracy"`
	HotAccuracy float64 `json:"hot_accuracy"`
	Target16Accuracy float64 `json:"target16_accuracy"`
	Target16ExactAccuracy float64 `json:"target16_exact_accuracy"`
	NonpersistentAdmissionRate float64 `json:"nonpersistent_admission_rate"`
	NonpersistentRetentionRate float64 `json:"nonpersistent_retention_rate"`
	OneShotFalseAdmissions int `json:"one_shot_false_admissions"`
	MaxCurrentTableEntries int `json:"max_current_table_entries"`
	MaxHistoryTableEntries int `json:"max_history_table_entries"`
	RecallEntriesUsed int `json:"recall_entries_used"`
}
type UP130CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP129CSeal string `json:"source_up129c_seal"`
	FirstWaveQualified int `json:"first_wave_qualified"`
	SecondWaveQualified int `json:"second_wave_qualified"`
	PreOverflowTargetSightings int `json:"pre_overflow_target_sightings"`
	PreOverflowWrites int `json:"pre_overflow_writes"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	ReplacementSignal string `json:"replacement_signal"`
	TieBreak string `json:"tie_break"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	QueryPriorityUsed bool `json:"query_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP130CPoint `json:"points"`
}
type up130cEpisode struct{
	panicHit bool;panicMessage string
	validAdmit,validTrials,validHits,validTotal int
	hotHits,hotTotal,targetHits,targetTotal int
	exact bool
	nonAdmit,nonTrials,nonHits,nonTotal int
	fp,maxCurrent,maxHistory,recallEntries,postEvictions int
}
func up130cOneShots(x *up125cAgeEvictMachine,n int,next *int,rng *sq0RNG,out *up130cEpisode){
	for i:=0;i<n;i++{a,_:=up125cProcess(x,*next,rng.intn(32),&out.maxCurrent,&out.maxHistory);(*next)++;if a{out.fp++}}
}
func up130cEpisodeRun(arm string,base,ep int)(out up130cEpisode){
	defer func(){if v:=recover();v!=nil{out.panicHit=true;out.panicMessage=fmt.Sprint(v)}}()
	armCode:=0;if arm=="edge_same_generation"{armCode=1}else if arm=="edge_cross_boundary"{armCode=2}
	rng:=newSQ0RNG(sq0Seed(base,7101+armCode*281,ep))
	x:=&up125cAgeEvictMachine{maxAge:2};truth:=map[int]int{};next:=300
	targets:=[]int{100,101,102,103}
	setA:=make([]int,16);for i:=range setA{setA[i]=200+i}
	setB:=make([]int,12);for i:=range setB{setB[i]=220+i}

	for k:=0;k<32;k++{v:=rng.intn(32);truth[k]=v;a,_:=up125cProcess(x,k,v,&out.maxCurrent,&out.maxHistory);if k>=16&&a{out.fp++}}
	for k:=0;k<12;k++{x.process(k,truth[k]);x.mem.query(k)}
	up125cObserve(x,&out.maxCurrent,&out.maxHistory);if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	for _,k:=range setA{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	for _,k:=range targets{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.validAdmit,&out.maxCurrent,&out.maxHistory);out.validTrials++}
	for _,k:=range setB{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	first:=make([]int,16);for i:=0;i<16;i++{first[i]=240+i}
	for _,k:=range first{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	evBefore:=x.evictions[0]+x.evictions[1]+x.evictions[2]

	switch arm{
	case "early_same_generation":
		for _,k:=range targets{up125cPresent(x,k,truth[k],1,&out.validAdmit,&out.maxCurrent,&out.maxHistory)}
		up130cOneShots(x,24,&next,rng,&out)
		for _,k:=range targets{up125cPresent(x,k,truth[k],1,&out.validAdmit,&out.maxCurrent,&out.maxHistory)}
		up130cOneShots(x,32,&next,rng,&out)
	case "edge_same_generation":
		up130cOneShots(x,24,&next,rng,&out)
		for _,k:=range targets{up125cPresent(x,k,truth[k],1,&out.validAdmit,&out.maxCurrent,&out.maxHistory)}
		for _,k:=range targets{up125cPresent(x,k,truth[k],1,&out.validAdmit,&out.maxCurrent,&out.maxHistory)}
		up130cOneShots(x,32,&next,rng,&out)
	case "edge_cross_boundary":
		up130cOneShots(x,28,&next,rng,&out)
		for _,k:=range targets{up125cPresent(x,k,truth[k],1,&out.validAdmit,&out.maxCurrent,&out.maxHistory)}
		for _,k:=range targets{up125cPresent(x,k,truth[k],1,&out.validAdmit,&out.maxCurrent,&out.maxHistory)}
		up130cOneShots(x,28,&next,rng,&out)
	}

	second:=[]int{260,261,262,263}
	for _,k:=range second{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	up130cOneShots(x,32,&next,rng,&out)
	out.postEvictions=x.evictions[0]+x.evictions[1]+x.evictions[2]-evBefore

	for j:=0;j<12288;j++{
		key:=1000+j;a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory);if a{out.fp++}
		if (j+1)%4==0{for k:=0;k<12;k++{x.mem.query(k)};for _,g:=range [][]int{targets,setA,setB,first,second}{for _,k:=range g{x.mem.query(k)}}}
	}
	out.exact=true
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;out.targetTotal++;if ok&&got==truth[k]{out.hotHits++;out.targetHits++}else{out.exact=false}}
	for _,k:=range targets{got,ok:=x.mem.query(k);out.validTotal++;out.targetTotal++;if ok&&got==truth[k]{out.validHits++;out.targetHits++}else{out.exact=false}}
	for _,g:=range [][]int{setA,setB,first,second}{for _,k:=range g{got,ok:=x.mem.query(k);out.nonTotal++;if ok&&got==truth[k]{out.nonHits++}}}
	out.recallEntries=x.mem.count
	return
}
func up130cRun(arm string,seeds []int)UP130CPoint{
	var panicN,total,postE int;msg:=""
	var validAdmit,validTrials,validHits,validTotal,hotHits,hotTotal,targetHits,targetTotal,exactHits int
	var nonAdmit,nonTrials,nonHits,nonTotal,fp,maxCurrent,maxHistory,maxRecall int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up130cEpisodeRun(arm,base,ep);total++;postE+=e.postEvictions
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
		fp+=e.fp;nonAdmit+=e.nonAdmit;nonTrials+=e.nonTrials
		if e.panicHit{panicN++;if msg==""{msg=e.panicMessage};continue}
		validAdmit+=e.validAdmit;validTrials+=e.validTrials;validHits+=e.validHits;validTotal+=e.validTotal;hotHits+=e.hotHits;hotTotal+=e.hotTotal;targetHits+=e.targetHits;targetTotal+=e.targetTotal
		if e.exact{exactHits++};nonHits+=e.nonHits;nonTotal+=e.nonTotal
	}}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	return UP130CPoint{Arm:arm,BoundaryBetweenSightings:arm=="edge_cross_boundary",PreOverflowTargetSightings:2,PreOverflowWrites:64,PanicRate:rate(panicN,total),PanicMessage:msg,PostFirstWaveEvictions:postE,ValidAdmissionRate:rate(validAdmit,validTrials),ValidAccuracy:rate(validHits,validTotal),HotAccuracy:rate(hotHits,hotTotal),Target16Accuracy:rate(targetHits,targetTotal),Target16ExactAccuracy:rate(exactHits,total-panicN),NonpersistentAdmissionRate:rate(nonAdmit,nonTrials),NonpersistentRetentionRate:rate(nonHits,nonTotal),OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall}
}
func RunUP130C()(UP130CResult,error){
	res:=UP130CResult{Schema:UP130CEdgeSchema,Experiment:"UP-130C-generation-edge-causality",SourceUP129CSeal:"608584f4c506cf8780a348fe00d4322690466e9f",FirstWaveQualified:16,SecondWaveQualified:4,PreOverflowTargetSightings:2,PreOverflowWrites:64,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,ReplacementSignal:"maximum_history_age_only",TieBreak:"lowest_table_index",SemanticPriorityUsed:false,QueryPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	seeds:=[]int{247000000,248000000}
	for _,arm:=range []string{"early_same_generation","edge_same_generation","edge_cross_boundary"}{res.Points=append(res.Points,up130cRun(arm,seeds))}
	return res,nil
}
