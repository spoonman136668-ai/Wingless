package unitary

const UP139CNoveltyClockSchema="wingless.up139c-novelty-pressure-clock.v1"

type UP139CPoint struct{
	Arm string `json:"arm"`
	Target int `json:"target"`
	ProcessCalls int `json:"process_calls"`
	NovelCandidateCalls int `json:"novel_candidate_calls"`
	ExactMemoryUpdateCalls int `json:"exact_memory_update_calls"`
	CandidateBoundaries int `json:"candidate_boundaries"`
	Checkpoints int `json:"checkpoints"`
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
type UP139CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP138CSeal string `json:"source_up138c_seal"`
	ProcessCallsPerArm int `json:"process_calls_per_arm"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	MaxHistoryAge int `json:"max_history_age"`
	ExactUpdateKeyCount int `json:"exact_update_key_count"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP139CPoint `json:"points"`
}
type up139cEpisode struct{
	panicHit bool
	admit,hits,hotHits,hotTotal int
	fp,maxCurrent,maxHistory,recallEntries int
	historyPresent,historyAge,checkpoints int
}
func up139cNovelCount(arm string)int{
	switch arm{
	case "novel_0_exact_96":return 0
	case "novel_32_exact_64":return 32
	case "novel_64_exact_32":return 64
	default:return 96
	}
}
func up139cEpisodeRun(arm string,target,base,ep int)(out up139cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	novel:=up139cNovelCount(arm)
	rng:=newSQ0RNG(sq0Seed(base,8001+novel*347+target*11,ep))
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
	if x.counter!=0{panic("UP139C_SETUP_NOT_ON_BOUNDARY")}
	checkBefore:=x.checkpoints

	for op:=0;op<96;op++{
		if up138cIsWrite(op,novel){
			a,_:=up125cProcess(x,next,rng.intn(32),&out.maxCurrent,&out.maxHistory);next++;if a{out.fp++}
		}else{
			a,_:=up125cProcess(x,0,truth[0],&out.maxCurrent,&out.maxHistory)
			if !a{panic("UP139C_EXACT_UPDATE_NOT_ADMITTED")}
		}
	}
	out.checkpoints=x.checkpoints-checkBefore
	if hi:=up118cFind(&x.history,target);hi>=0{out.historyPresent=1;out.historyAge=int(x.history[hi]>>14)}
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
func up139cRun(arm string,target int,seeds []int)UP139CPoint{
	var panicN,total,admit,hits,hotHits,hotTotal,fp,maxCurrent,maxHistory,maxRecall,histPresent,histAge,checkpoints int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up139cEpisodeRun(arm,target,base,ep);total++
		if e.panicHit{panicN++;continue}
		admit+=e.admit;hits+=e.hits;hotHits+=e.hotHits;hotTotal+=e.hotTotal;fp+=e.fp;histPresent+=e.historyPresent;histAge+=e.historyAge;checkpoints+=e.checkpoints
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
	}}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	novel:=up139cNovelCount(arm);exact:=96-novel;meanAge:=0.0;if histPresent>0{meanAge=float64(histAge)/float64(histPresent)}
	return UP139CPoint{Arm:arm,Target:target,ProcessCalls:96,NovelCandidateCalls:novel,ExactMemoryUpdateCalls:exact,CandidateBoundaries:novel/32,Checkpoints:checkpoints,HistoryPresentBeforePairRate:rate(histPresent,total-panicN),MeanHistoryAgeBeforePair:meanAge,AdmissionRate:rate(admit,total-panicN),FinalAccuracy:rate(hits,total-panicN),HotAccuracy:rate(hotHits,hotTotal),OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall,PanicRate:rate(panicN,total)}
}
func RunUP139C()(UP139CResult,error){
	res:=UP139CResult{Schema:UP139CNoveltyClockSchema,Experiment:"UP-139C-novelty-pressure-clock",SourceUP138CSeal:"5ce5a9dd0776e3054dc16c90a322db9e5ce16caf",ProcessCallsPerArm:96,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,MaxHistoryAge:2,ExactUpdateKeyCount:1,SemanticPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	seeds:=[]int{265000000,266000000};arms:=[]string{"novel_0_exact_96","novel_32_exact_64","novel_64_exact_32","novel_96_exact_0"}
	for _,target:=range []int{100,103}{for _,arm:=range arms{res.Points=append(res.Points,up139cRun(arm,target,seeds))}}
	return res,nil
}
