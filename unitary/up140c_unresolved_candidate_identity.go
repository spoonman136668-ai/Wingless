package unitary

const UP140CIdentitySchema="wingless.up140c-unresolved-candidate-identity.v1"

type UP140CPoint struct{
	Arm string `json:"arm"`
	Target int `json:"target"`
	UnresolvedCalls int `json:"unresolved_calls"`
	UniqueIdentitiesTotal int `json:"unique_identities_total"`
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
type UP140CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP139CSeal string `json:"source_up139c_seal"`
	UnresolvedCallsPerArm int `json:"unresolved_calls_per_arm"`
	CandidateBoundariesPerArm int `json:"candidate_boundaries_per_arm"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	MaxHistoryAge int `json:"max_history_age"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP140CPoint `json:"points"`
}
type up140cEpisode struct{
	panicHit bool
	admit,hits,hotHits,hotTotal int
	maxCurrent,maxHistory,recallEntries int
	h1Present,h1Age,h1Entries int
	h2Present,h2Age,h2Entries int
	ev [3]int
	unexpectedAdmissions int
}

func up140cGroupSize(arm string)int{
	switch arm{
	case "unique_32":return 1
	case "paired_16":return 2
	case "quartet_8":return 4
	default:return 32
	}
}
func up140cUniquePerWindow(arm string)int{return 32/up140cGroupSize(arm)}
func up140cHist(x *up125cAgeEvictMachine,target int)(present,age int){
	if hi:=up118cFind(&x.history,target);hi>=0{return 1,int(x.history[hi]>>14)}
	return 0,0
}
func up140cSetup(target,base,ep int)(x *up125cAgeEvictMachine,truth map[int]int,next int,rng *sq0RNG,outMaxCurrent,outMaxHistory *int){
	rng=newSQ0RNG(sq0Seed(base,8101+target*13,ep))
	x=&up125cAgeEvictMachine{maxAge:2};truth=map[int]int{};next=300
	targets:=[]int{100,101,102,103}
	setA:=make([]int,16);for i:=range setA{setA[i]=200+i}
	setB:=make([]int,12);for i:=range setB{setB[i]=220+i}
	maxC,maxH:=0,0
	for k:=0;k<32;k++{v:=rng.intn(32);truth[k]=v;up125cProcess(x,k,v,&maxC,&maxH)}
	for k:=0;k<12;k++{x.process(k,truth[k]);x.mem.query(k)}
	up125cObserve(x,&maxC,&maxH);if x.counter!=0{up125cFill(x,&next,rng,&maxC,&maxH)}
	for _,k:=range setA{v:=rng.intn(32);truth[k]=v;d:=0;up125cPresent(x,k,v,2,&d,&maxC,&maxH)}
	for _,k:=range targets{v:=rng.intn(32);truth[k]=v;d:=0;up125cPresent(x,k,v,2,&d,&maxC,&maxH)}
	for _,k:=range setB{v:=rng.intn(32);truth[k]=v;d:=0;up125cPresent(x,k,v,2,&d,&maxC,&maxH)}
	if x.counter!=0{up125cFill(x,&next,rng,&maxC,&maxH)}
	first:=make([]int,16);for i:=range first{first[i]=240+i}
	for _,k:=range first{v:=rng.intn(32);truth[k]=v;d:=0;up125cPresent(x,k,v,2,&d,&maxC,&maxH)}
	if x.counter!=0{up125cFill(x,&next,rng,&maxC,&maxH)}
	for pos:=1;pos<=96;pos++{
		switch pos{
		case 5,37:
			d:=0;up125cPresent(x,target,truth[target],1,&d,&maxC,&maxH)
		case 74,75:
			d:=0;up125cPresent(x,target,truth[target],1,&d,&maxC,&maxH)
		default:
			up125cProcess(x,next,rng.intn(32),&maxC,&maxH);next++
		}
	}
	if x.counter!=0{panic("UP140C_SETUP_NOT_ON_BOUNDARY")}
	*outMaxCurrent=maxC;*outMaxHistory=maxH
	return
}
func up140cEpisodeRun(arm string,target,base,ep int)(out up140cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	maxC,maxH:=0,0
	x,truth,next,rng:=up140cSetup(target,base,ep,&maxC,&maxH)
	out.maxCurrent=maxC;out.maxHistory=maxH
	startEv:=x.evictions
	group:=up140cGroupSize(arm)
	for window:=0;window<2;window++{
		keyBase:=10000+window*1000+target*100
		for op:=0;op<32;op++{
			key:=keyBase+(op/group)
			a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory)
			if a{out.unexpectedAdmissions++}
		}
		if x.counter!=0{panic("UP140C_WINDOW_NOT_ON_BOUNDARY")}
		p,a:=up140cHist(x,target)
		if window==0{
			out.h1Present=p;out.h1Age=a;out.h1Entries=up125cCount(&x.history)
		}else{
			out.h2Present=p;out.h2Age=a;out.h2Entries=up125cCount(&x.history)
		}
	}
	for i:=0;i<3;i++{out.ev[i]=x.evictions[i]-startEv[i]}
	up125cPresent(x,target,truth[target],2,&out.admit,&out.maxCurrent,&out.maxHistory)
	for j:=0;j<12288;j++{
		key:=20000+j;up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory)
		if (j+1)%4==0{for k:=0;k<12;k++{x.mem.query(k)}}
	}
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;if ok&&got==truth[k]{out.hotHits++}}
	if got,ok:=x.mem.query(target);ok&&got==truth[target]{out.hits++}
	out.recallEntries=x.mem.count
	_ = next
	return
}
func up140cRun(arm string,target int,seeds []int)UP140CPoint{
	var panicN,total,admit,hits,hotHits,hotTotal,maxCurrent,maxHistory,maxRecall int
	var h1p,h1a,h1e,h2p,h2a,h2e,unexpected int
	var ev [3]int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up140cEpisodeRun(arm,target,base,ep);total++
		if e.panicHit{panicN++;continue}
		admit+=e.admit;hits+=e.hits;hotHits+=e.hotHits;hotTotal+=e.hotTotal
		h1p+=e.h1Present;h1a+=e.h1Age;h1e+=e.h1Entries;h2p+=e.h2Present;h2a+=e.h2Age;h2e+=e.h2Entries
		unexpected+=e.unexpectedAdmissions;for i:=0;i<3;i++{ev[i]+=e.ev[i]}
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
	}}
	n:=total-panicN
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	ageMean:=func(sum,present int)float64{if present==0{return 0};return float64(sum)/float64(present)}
	return UP140CPoint{
		Arm:arm,Target:target,UnresolvedCalls:64,UniqueIdentitiesTotal:2*up140cUniquePerWindow(arm),
		HistoryPresentAfterBoundary1Rate:rate(h1p,n),MeanHistoryAgeAfterBoundary1:ageMean(h1a,h1p),MeanHistoryEntriesAfterBoundary1:rate(h1e,n),
		HistoryPresentAfterBoundary2Rate:rate(h2p,n),MeanHistoryAgeAfterBoundary2:ageMean(h2a,h2p),MeanHistoryEntriesAfterBoundary2:rate(h2e,n),
		Age0Evictions:ev[0],Age1Evictions:ev[1],Age2Evictions:ev[2],UnexpectedManipulationAdmissions:unexpected,
		AdmissionRate:rate(admit,n),FinalAccuracy:rate(hits,n),HotAccuracy:rate(hotHits,hotTotal),
		MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall,PanicRate:rate(panicN,total),
	}
}
func RunUP140C()(UP140CResult,error){
	res:=UP140CResult{Schema:UP140CIdentitySchema,Experiment:"UP-140C-unresolved-candidate-identity",SourceUP139CSeal:"d0b747eada76bb3ca82b140b4888e033620e3a6e",UnresolvedCallsPerArm:64,CandidateBoundariesPerArm:2,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,MaxHistoryAge:2,SemanticPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	seeds:=[]int{267000000,268000000};arms:=[]string{"unique_32","paired_16","quartet_8","repeat_1"}
	for _,target:=range []int{100,103}{for _,arm:=range arms{res.Points=append(res.Points,up140cRun(arm,target,seeds))}}
	return res,nil
}
