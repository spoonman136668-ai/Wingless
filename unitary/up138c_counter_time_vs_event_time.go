package unitary

const UP138CCounterTimeSchema="wingless.up138c-counter-time-vs-event-time.v1"

type UP138CPoint struct{
	Arm string `json:"arm"`
	Target int `json:"target"`
	RawOperations int `json:"raw_operations"`
	CandidateWrites int `json:"candidate_writes"`
	ReadOnlyQueries int `json:"read_only_queries"`
	CandidateBoundaries int `json:"candidate_boundaries"`
	HistoryPresentBeforePairRate float64 `json:"history_present_before_pair_rate"`
	MeanHistoryAgeBeforePair float64 `json:"mean_history_age_before_pair"`
	AdmissionRate float64 `json:"admission_rate"`
	FinalAccuracy float64 `json:"final_accuracy"`
	HotAccuracy float64 `json:"hot_accuracy"`
	OneShotFalseAdmissions int `json:"one_shot_false_admissions"`
	MaxCurrentTableEntries int `json:"max_current_table_entries"`
	MaxHistoryTableEntries int `json:"max_history_table_entries"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	PanicRate float64 `json:"panic_rate"`
}
type UP138CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP137CSeal string `json:"source_up137c_seal"`
	RawOperationsPerArm int `json:"raw_operations_per_arm"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	MaxHistoryAge int `json:"max_history_age"`
	QueryOperationsMutateAdmissionState bool `json:"query_operations_mutate_admission_state"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP138CPoint `json:"points"`
}
type up138cEpisode struct{
	panicHit bool
	admit,hits,hotHits,hotTotal int
	fp,maxCurrent,maxHistory,recallEntries int
	historyPresent int
	historyAge int
}
func up138cArmWrites(arm string)int{
	switch arm{
	case "write_0_query_96":return 0
	case "write_32_query_64":return 32
	case "write_64_query_32":return 64
	default:return 96
	}
}
func up138cIsWrite(op,writes int)bool{
	if writes==0{return false}
	if writes==96{return true}
	// Deterministic even spacing over 96 slots.
	return ((op+1)*writes)/96 > (op*writes)/96
}
func up138cEpisodeRun(arm string,target,base,ep int)(out up138cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	writes:=up138cArmWrites(arm)
	rng:=newSQ0RNG(sq0Seed(base,7901+writes*337+target*7,ep))
	x:=&up125cAgeEvictMachine{maxAge:2};truth:=map[int]int{};next:=300
	targets:=[]int{100,101,102,103}
	setA:=make([]int,16);for i:=range setA{setA[i]=200+i}
	setB:=make([]int,12);for i:=range setB{setB[i]=220+i}

	for k:=0;k<32;k++{v:=rng.intn(32);truth[k]=v;a,_:=up125cProcess(x,k,v,&out.maxCurrent,&out.maxHistory);if k>=16&&a{out.fp++}}
	for k:=0;k<12;k++{x.process(k,truth[k]);x.mem.query(k)}
	up125cObserve(x,&out.maxCurrent,&out.maxHistory)
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}

	for _,k:=range setA{v:=rng.intn(32);truth[k]=v;dummy:=0;up125cPresent(x,k,v,2,&dummy,&out.maxCurrent,&out.maxHistory)}
	for _,k:=range targets{v:=rng.intn(32);truth[k]=v;dummy:=0;up125cPresent(x,k,v,2,&dummy,&out.maxCurrent,&out.maxHistory)}
	for _,k:=range setB{v:=rng.intn(32);truth[k]=v;dummy:=0;up125cPresent(x,k,v,2,&dummy,&out.maxCurrent,&out.maxHistory)}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	first:=make([]int,16);for i:=0;i<16;i++{first[i]=240+i}
	for _,k:=range first{v:=rng.intn(32);truth[k]=v;dummy:=0;up125cPresent(x,k,v,2,&dummy,&out.maxCurrent,&out.maxHistory)}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}

	// G1/G2 isolated sightings, then a G3 pair rebuilds fresh history.
	for pos:=1;pos<=96;pos++{
		switch pos{
		case 5,37:
			dummy:=0;up125cPresent(x,target,truth[target],1,&dummy,&out.maxCurrent,&out.maxHistory)
		case 74,75:
			dummy:=0;up125cPresent(x,target,truth[target],1,&dummy,&out.maxCurrent,&out.maxHistory)
		default:
			a,_:=up125cProcess(x,next,rng.intn(32),&out.maxCurrent,&out.maxHistory);next++;if a{out.fp++}
		}
	}
	if x.counter!=0{panic("UP138C_SETUP_NOT_ON_BOUNDARY")}

	// Exactly 96 raw operations; only candidate writes advance admission time.
	for op:=0;op<96;op++{
		if up138cIsWrite(op,writes){
			a,_:=up125cProcess(x,next,rng.intn(32),&out.maxCurrent,&out.maxHistory);next++;if a{out.fp++}
		}else{
			x.mem.query(op%12)
		}
	}
	if hi:=up118cFind(&x.history,target);hi>=0{
		out.historyPresent=1
		out.historyAge=int(x.history[hi]>>14)
	}

	up125cPresent(x,target,truth[target],2,&out.admit,&out.maxCurrent,&out.maxHistory)

	for j:=0;j<12288;j++{
		key:=1000+j;a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory);if a{out.fp++}
		if (j+1)%4==0{for k:=0;k<12;k++{x.mem.query(k)}}
	}
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;if ok&&got==truth[k]{out.hotHits++}}
	if got,ok:=x.mem.query(target);ok&&got==truth[target]{out.hits++}
	out.recallEntries=x.mem.count
	return
}
func up138cRun(arm string,target int,seeds []int)UP138CPoint{
	var panicN,total,admit,hits,hotHits,hotTotal,fp,maxCurrent,maxHistory,maxRecall,histPresent,histAge int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up138cEpisodeRun(arm,target,base,ep);total++
		if e.panicHit{panicN++;continue}
		admit+=e.admit;hits+=e.hits;hotHits+=e.hotHits;hotTotal+=e.hotTotal;fp+=e.fp;histPresent+=e.historyPresent;histAge+=e.historyAge
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
	}}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	w:=up138cArmWrites(arm);q:=96-w
	meanAge:=0.0;if histPresent>0{meanAge=float64(histAge)/float64(histPresent)}
	return UP138CPoint{Arm:arm,Target:target,RawOperations:96,CandidateWrites:w,ReadOnlyQueries:q,CandidateBoundaries:w/32,HistoryPresentBeforePairRate:rate(histPresent,total-panicN),MeanHistoryAgeBeforePair:meanAge,AdmissionRate:rate(admit,total-panicN),FinalAccuracy:rate(hits,total-panicN),HotAccuracy:rate(hotHits,hotTotal),OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall,PanicRate:rate(panicN,total)}
}
func RunUP138C()(UP138CResult,error){
	res:=UP138CResult{Schema:UP138CCounterTimeSchema,Experiment:"UP-138C-counter-time-vs-event-time",SourceUP137CSeal:"85c0cd809c80ef71dd5b4475b8506f52feef0026",RawOperationsPerArm:96,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,MaxHistoryAge:2,QueryOperationsMutateAdmissionState:false,SemanticPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	seeds:=[]int{263000000,264000000}
	arms:=[]string{"write_0_query_96","write_32_query_64","write_64_query_32","write_96_query_0"}
	for _,target:=range []int{100,103}{for _,arm:=range arms{res.Points=append(res.Points,up138cRun(arm,target,seeds))}}
	return res,nil
}
