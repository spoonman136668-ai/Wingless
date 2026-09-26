package unitary

const UP131CMixedSchema="wingless.up131c-mixed-phase-selectivity.v1"

type UP131CPoint struct{
	Arm string `json:"arm"`
	SameTargets []int `json:"same_targets"`
	CrossTargets []int `json:"cross_targets"`
	SameAdmissionRate float64 `json:"same_admission_rate"`
	CrossAdmissionRate float64 `json:"cross_admission_rate"`
	SameAccuracy float64 `json:"same_accuracy"`
	CrossAccuracy float64 `json:"cross_accuracy"`
	HotAccuracy float64 `json:"hot_accuracy"`
	Target16Accuracy float64 `json:"target16_accuracy"`
	Target16ExactAccuracy float64 `json:"target16_exact_accuracy"`
	NonpersistentAdmissionRate float64 `json:"nonpersistent_admission_rate"`
	NonpersistentRetentionRate float64 `json:"nonpersistent_retention_rate"`
	OneShotFalseAdmissions int `json:"one_shot_false_admissions"`
	MaxCurrentTableEntries int `json:"max_current_table_entries"`
	MaxHistoryTableEntries int `json:"max_history_table_entries"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	PanicRate float64 `json:"panic_rate"`
}
type UP131CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP130CSeal string `json:"source_up130c_seal"`
	FirstWaveQualified int `json:"first_wave_qualified"`
	SecondWaveQualified int `json:"second_wave_qualified"`
	PreOverflowWrites int `json:"pre_overflow_writes"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	QueryPriorityUsed bool `json:"query_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP131CPoint `json:"points"`
}

type up131cEpisode struct{
	panicHit bool
	sameAdmit,crossAdmit int
	sameHits,crossHits int
	hotHits,hotTotal,targetHits,targetTotal int
	exact bool
	nonAdmit,nonTrials,nonHits,nonTotal int
	fp,maxCurrent,maxHistory,recallEntries int
}
func up131cOneShots(x *up125cAgeEvictMachine,n int,next *int,rng *sq0RNG,out *up131cEpisode){
	for i:=0;i<n;i++{a,_:=up125cProcess(x,*next,rng.intn(32),&out.maxCurrent,&out.maxHistory);(*next)++;if a{out.fp++}}
}
func up131cEpisodeRun(arm string,base,ep int)(out up131cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	code:=0;if arm=="assignment_b"{code=1}
	rng:=newSQ0RNG(sq0Seed(base,7201+code*283,ep))
	x:=&up125cAgeEvictMachine{maxAge:2};truth:=map[int]int{};next:=300
	targets:=[]int{100,101,102,103}
	same:=[]int{100,101};cross:=[]int{102,103}
	if arm=="assignment_b"{same,cross=cross,same}
	setA:=make([]int,16);for i:=range setA{setA[i]=200+i}
	setB:=make([]int,12);for i:=range setB{setB[i]=220+i}

	for k:=0;k<32;k++{v:=rng.intn(32);truth[k]=v;a,_:=up125cProcess(x,k,v,&out.maxCurrent,&out.maxHistory);if k>=16&&a{out.fp++}}
	for k:=0;k<12;k++{x.process(k,truth[k]);x.mem.query(k)}
	up125cObserve(x,&out.maxCurrent,&out.maxHistory);if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	for _,k:=range setA{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	for _,k:=range targets{v:=rng.intn(32);truth[k]=v;dummy:=0;up125cPresent(x,k,v,2,&dummy,&out.maxCurrent,&out.maxHistory)}
	for _,k:=range setB{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	first:=make([]int,16);for i:=0;i<16;i++{first[i]=240+i}
	for _,k:=range first{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}

	for _,k:=range same{up125cPresent(x,k,truth[k],1,&out.sameAdmit,&out.maxCurrent,&out.maxHistory)}
	up131cOneShots(x,24,&next,rng,&out)
	for _,k:=range same{up125cPresent(x,k,truth[k],1,&out.sameAdmit,&out.maxCurrent,&out.maxHistory)}
	for _,k:=range cross{up125cPresent(x,k,truth[k],1,&out.crossAdmit,&out.maxCurrent,&out.maxHistory)}
	up131cOneShots(x,2,&next,rng,&out)
	for _,k:=range cross{up125cPresent(x,k,truth[k],1,&out.crossAdmit,&out.maxCurrent,&out.maxHistory)}
	up131cOneShots(x,30,&next,rng,&out)

	second:=[]int{260,261,262,263}
	for _,k:=range second{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	up131cOneShots(x,32,&next,rng,&out)

	for j:=0;j<12288;j++{
		key:=1000+j;a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory);if a{out.fp++}
		if (j+1)%4==0{for k:=0;k<12;k++{x.mem.query(k)};for _,g:=range [][]int{targets,setA,setB,first,second}{for _,k:=range g{x.mem.query(k)}}}
	}
	out.exact=true
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;out.targetTotal++;if ok&&got==truth[k]{out.hotHits++;out.targetHits++}else{out.exact=false}}
	for _,k:=range same{got,ok:=x.mem.query(k);out.targetTotal++;if ok&&got==truth[k]{out.sameHits++;out.targetHits++}else{out.exact=false}}
	for _,k:=range cross{got,ok:=x.mem.query(k);out.targetTotal++;if ok&&got==truth[k]{out.crossHits++;out.targetHits++}else{out.exact=false}}
	for _,g:=range [][]int{setA,setB,first,second}{for _,k:=range g{got,ok:=x.mem.query(k);out.nonTotal++;if ok&&got==truth[k]{out.nonHits++}}}
	out.recallEntries=x.mem.count
	return
}
func up131cRun(arm string,seeds []int)UP131CPoint{
	var panicN,total,sameAdmit,crossAdmit,sameHits,crossHits,hotHits,hotTotal,targetHits,targetTotal,exactHits int
	var nonAdmit,nonTrials,nonHits,nonTotal,fp,maxCurrent,maxHistory,maxRecall int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up131cEpisodeRun(arm,base,ep);total++;if e.panicHit{panicN++;continue}
		sameAdmit+=e.sameAdmit;crossAdmit+=e.crossAdmit;sameHits+=e.sameHits;crossHits+=e.crossHits
		hotHits+=e.hotHits;hotTotal+=e.hotTotal;targetHits+=e.targetHits;targetTotal+=e.targetTotal;if e.exact{exactHits++}
		nonAdmit+=e.nonAdmit;nonTrials+=e.nonTrials;nonHits+=e.nonHits;nonTotal+=e.nonTotal;fp+=e.fp
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
	}}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	same:=[]int{100,101};cross:=[]int{102,103};if arm=="assignment_b"{same,cross=cross,same}
	return UP131CPoint{Arm:arm,SameTargets:same,CrossTargets:cross,SameAdmissionRate:rate(sameAdmit,total*2),CrossAdmissionRate:rate(crossAdmit,total*2),SameAccuracy:rate(sameHits,total*2),CrossAccuracy:rate(crossHits,total*2),HotAccuracy:rate(hotHits,hotTotal),Target16Accuracy:rate(targetHits,targetTotal),Target16ExactAccuracy:rate(exactHits,total-panicN),NonpersistentAdmissionRate:rate(nonAdmit,nonTrials),NonpersistentRetentionRate:rate(nonHits,nonTotal),OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall,PanicRate:rate(panicN,total)}
}
func RunUP131C()(UP131CResult,error){
	res:=UP131CResult{Schema:UP131CMixedSchema,Experiment:"UP-131C-mixed-phase-selectivity",SourceUP130CSeal:"639b5607ea32406c092b7731ee565752bb322025",FirstWaveQualified:16,SecondWaveQualified:4,PreOverflowWrites:64,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,SemanticPriorityUsed:false,QueryPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	seeds:=[]int{249000000,250000000}
	for _,arm:=range []string{"assignment_a","assignment_b"}{res.Points=append(res.Points,up131cRun(arm,seeds))}
	return res,nil
}
