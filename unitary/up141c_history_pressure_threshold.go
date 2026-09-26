package unitary

const UP141CHistoryPressureSchema="wingless.up141c-history-pressure-threshold.v1"

type UP141CPoint struct{
	QualifiedPerWindow int `json:"qualified_per_window"`
	Target int `json:"target"`
	UnresolvedCalls int `json:"unresolved_calls"`
	HistoryPresentAfterBoundary1Rate float64 `json:"history_present_after_boundary1_rate"`
	MeanHistoryAgeAfterBoundary1 float64 `json:"mean_history_age_after_boundary1"`
	MeanHistoryEntriesAfterBoundary1 float64 `json:"mean_history_entries_after_boundary1"`
	HistoryPresentAfterBoundary2Rate float64 `json:"history_present_after_boundary2_rate"`
	MeanHistoryAgeAfterBoundary2 float64 `json:"mean_history_age_after_boundary2"`
	MeanHistoryEntriesAfterBoundary2 float64 `json:"mean_history_entries_after_boundary2"`
	Age0Evictions int `json:"age0_evictions"`
	Age1Evictions int `json:"age1_evictions"`
	Age2Evictions int `json:"age2_evictions"`
	UnexpectedManipulationAdmissions int `json:"unexpected_manipulation_admissions"`
	AdmissionRate float64 `json:"admission_rate"`
	FinalAccuracy float64 `json:"final_accuracy"`
	HotAccuracy float64 `json:"hot_accuracy"`
	MaxCurrentTableEntries int `json:"max_current_table_entries"`
	MaxHistoryTableEntries int `json:"max_history_table_entries"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	PanicRate float64 `json:"panic_rate"`
}
type UP141CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP140CSeal string `json:"source_up140c_seal"`
	UnresolvedCallsPerArm int `json:"unresolved_calls_per_arm"`
	CandidateBoundariesPerArm int `json:"candidate_boundaries_per_arm"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	MaxHistoryAge int `json:"max_history_age"`
	QualifiedLevels []int `json:"qualified_levels"`
	MemoryIncreased bool `json:"memory_increased"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Points []UP141CPoint `json:"points"`
}
type up141cEpisode struct{
	panicHit bool
	admit,hits,hotHits,hotTotal int
	maxCurrent,maxHistory,recallEntries int
	h1Present,h1Age,h1Entries int
	h2Present,h2Age,h2Entries int
	ev [3]int
	unexpected int
}
func up141cEpisodeRun(q,target,base,ep int)(out up141cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	maxC,maxH:=0,0
	x,truth,next,rng:=up140cSetup(target,base,ep,&maxC,&maxH)
	out.maxCurrent=maxC;out.maxHistory=maxH
	startEv:=x.evictions
	for window:=0;window<2;window++{
		keyBase:=30000+window*1000+target*100
		// q paired identities -> q qualified current entries.
		for i:=0;i<q;i++{
			key:=keyBase+i
			for rep:=0;rep<2;rep++{
				a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory)
				if a{out.unexpected++}
			}
		}
		// Remaining calls are one-shot unique identities.
		for i:=0;i<32-2*q;i++{
			key:=keyBase+100+i
			a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory)
			if a{out.unexpected++}
		}
		if x.counter!=0{panic("UP141C_WINDOW_NOT_ON_BOUNDARY")}
		p,a:=up140cHist(x,target)
		if window==0{out.h1Present=p;out.h1Age=a;out.h1Entries=up125cCount(&x.history)}else{out.h2Present=p;out.h2Age=a;out.h2Entries=up125cCount(&x.history)}
	}
	for i:=0;i<3;i++{out.ev[i]=x.evictions[i]-startEv[i]}
	up125cPresent(x,target,truth[target],2,&out.admit,&out.maxCurrent,&out.maxHistory)
	for j:=0;j<12288;j++{
		key:=40000+j;up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory)
		if (j+1)%4==0{for k:=0;k<12;k++{x.mem.query(k)}}
	}
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;if ok&&got==truth[k]{out.hotHits++}}
	if got,ok:=x.mem.query(target);ok&&got==truth[target]{out.hits++}
	out.recallEntries=x.mem.count
	_=next
	return
}
func up141cRun(q,target int,seeds []int)UP141CPoint{
	var panicN,total,admit,hits,hotHits,hotTotal,maxCurrent,maxHistory,maxRecall int
	var h1p,h1a,h1e,h2p,h2a,h2e,unexpected int
	var ev [3]int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up141cEpisodeRun(q,target,base,ep);total++
		if e.panicHit{panicN++;continue}
		admit+=e.admit;hits+=e.hits;hotHits+=e.hotHits;hotTotal+=e.hotTotal
		h1p+=e.h1Present;h1a+=e.h1Age;h1e+=e.h1Entries;h2p+=e.h2Present;h2a+=e.h2Age;h2e+=e.h2Entries
		unexpected+=e.unexpected;for i:=0;i<3;i++{ev[i]+=e.ev[i]}
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
	}}
	n:=total-panicN
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	ageMean:=func(sum,present int)float64{if present==0{return 0};return float64(sum)/float64(present)}
	return UP141CPoint{QualifiedPerWindow:q,Target:target,UnresolvedCalls:64,
		HistoryPresentAfterBoundary1Rate:rate(h1p,n),MeanHistoryAgeAfterBoundary1:ageMean(h1a,h1p),MeanHistoryEntriesAfterBoundary1:rate(h1e,n),
		HistoryPresentAfterBoundary2Rate:rate(h2p,n),MeanHistoryAgeAfterBoundary2:ageMean(h2a,h2p),MeanHistoryEntriesAfterBoundary2:rate(h2e,n),
		Age0Evictions:ev[0],Age1Evictions:ev[1],Age2Evictions:ev[2],UnexpectedManipulationAdmissions:unexpected,
		AdmissionRate:rate(admit,n),FinalAccuracy:rate(hits,n),HotAccuracy:rate(hotHits,hotTotal),
		MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall,PanicRate:rate(panicN,total)}
}
func RunUP141C()(UP141CResult,error){
	levels:=[]int{0,12,14,15,16}
	res:=UP141CResult{Schema:UP141CHistoryPressureSchema,Experiment:"UP-141C-history-pressure-threshold",SourceUP140CSeal:"1ef8af85e7ee3d459e2147a3fb4611bffa78701e",UnresolvedCallsPerArm:64,CandidateBoundariesPerArm:2,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,MaxHistoryAge:2,QualifiedLevels:append([]int(nil),levels...),MemoryIncreased:false,SemanticPriorityUsed:false,FutureOracleUsed:false}
	seeds:=[]int{269000000,270000000}
	for _,target:=range []int{100,103}{for _,q:=range levels{res.Points=append(res.Points,up141cRun(q,target,seeds))}}
	return res,nil
}
