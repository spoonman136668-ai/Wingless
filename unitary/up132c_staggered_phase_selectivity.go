package unitary

const UP132CStaggeredSchema="wingless.up132c-staggered-phase-selectivity.v1"

type UP132CPathMetric struct{
	Path string `json:"path"`
	Target int `json:"target"`
	FirstPosition int `json:"first_position"`
	SecondPosition int `json:"second_position"`
	CrossesBoundary bool `json:"crosses_boundary"`
	AdmissionRate float64 `json:"admission_rate"`
	FinalAccuracy float64 `json:"final_accuracy"`
}
type UP132CPoint struct{
	Arm string `json:"arm"`
	PathMetrics []UP132CPathMetric `json:"path_metrics"`
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
type UP132CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP131CSeal string `json:"source_up131c_seal"`
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
	Points []UP132CPoint `json:"points"`
}
type up132cPath struct{name string;target,first,second int}
type up132cEpisode struct{
	panicHit bool
	admit [4]int
	hits [4]int
	hotHits,hotTotal,targetHits,targetTotal int
	exact bool
	nonAdmit,nonTrials,nonHits,nonTotal int
	fp,maxCurrent,maxHistory,recallEntries int
}
func up132cPaths(arm string)[]up132cPath{
	targets:=[]int{100,101,102,103}
	if arm=="assignment_b"{targets=[]int{102,103,100,101}}
	return []up132cPath{
		{"same_early",targets[0],1,16},
		{"same_late",targets[1],17,32},
		{"cross_edge",targets[2],31,33},
		{"cross_wide",targets[3],20,45},
	}
}
func up132cEpisodeRun(arm string,base,ep int)(out up132cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	code:=0;if arm=="assignment_b"{code=1}
	rng:=newSQ0RNG(sq0Seed(base,7301+code*293,ep))
	x:=&up125cAgeEvictMachine{maxAge:2};truth:=map[int]int{};next:=300
	targets:=[]int{100,101,102,103}
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

	paths:=up132cPaths(arm)
	for pos:=1;pos<=64;pos++{
		handled:=false
		for i,p:=range paths{
			if pos==p.first||pos==p.second{
				up125cPresent(x,p.target,truth[p.target],1,&out.admit[i],&out.maxCurrent,&out.maxHistory)
				handled=true
				break
			}
		}
		if !handled{
			a,_:=up125cProcess(x,next,rng.intn(32),&out.maxCurrent,&out.maxHistory);next++;if a{out.fp++}
		}
	}

	second:=[]int{260,261,262,263}
	for _,k:=range second{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	for i:=0;i<32;i++{a,_:=up125cProcess(x,next,rng.intn(32),&out.maxCurrent,&out.maxHistory);next++;if a{out.fp++}}

	for j:=0;j<12288;j++{
		key:=1000+j;a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory);if a{out.fp++}
		if (j+1)%4==0{for k:=0;k<12;k++{x.mem.query(k)};for _,g:=range [][]int{targets,setA,setB,first,second}{for _,k:=range g{x.mem.query(k)}}}
	}
	out.exact=true
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;out.targetTotal++;if ok&&got==truth[k]{out.hotHits++;out.targetHits++}else{out.exact=false}}
	for i,p:=range paths{got,ok:=x.mem.query(p.target);out.targetTotal++;if ok&&got==truth[p.target]{out.hits[i]++;out.targetHits++}else{out.exact=false}}
	for _,g:=range [][]int{setA,setB,first,second}{for _,k:=range g{got,ok:=x.mem.query(k);out.nonTotal++;if ok&&got==truth[k]{out.nonHits++}}}
	out.recallEntries=x.mem.count
	return
}
func up132cRun(arm string,seeds []int)UP132CPoint{
	var panicN,total,hotHits,hotTotal,targetHits,targetTotal,exactHits int
	var admit,hits [4]int
	var nonAdmit,nonTrials,nonHits,nonTotal,fp,maxCurrent,maxHistory,maxRecall int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up132cEpisodeRun(arm,base,ep);total++;if e.panicHit{panicN++;continue}
		for i:=0;i<4;i++{admit[i]+=e.admit[i];hits[i]+=e.hits[i]}
		hotHits+=e.hotHits;hotTotal+=e.hotTotal;targetHits+=e.targetHits;targetTotal+=e.targetTotal;if e.exact{exactHits++}
		nonAdmit+=e.nonAdmit;nonTrials+=e.nonTrials;nonHits+=e.nonHits;nonTotal+=e.nonTotal;fp+=e.fp
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
	}}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	paths:=up132cPaths(arm)
	pm:=make([]UP132CPathMetric,0,4)
	for i,p:=range paths{pm=append(pm,UP132CPathMetric{Path:p.name,Target:p.target,FirstPosition:p.first,SecondPosition:p.second,CrossesBoundary:p.first<=32&&p.second>32,AdmissionRate:rate(admit[i],total),FinalAccuracy:rate(hits[i],total)})}
	return UP132CPoint{Arm:arm,PathMetrics:pm,HotAccuracy:rate(hotHits,hotTotal),Target16Accuracy:rate(targetHits,targetTotal),Target16ExactAccuracy:rate(exactHits,total-panicN),NonpersistentAdmissionRate:rate(nonAdmit,nonTrials),NonpersistentRetentionRate:rate(nonHits,nonTotal),OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall,PanicRate:rate(panicN,total)}
}
func RunUP132C()(UP132CResult,error){
	res:=UP132CResult{Schema:UP132CStaggeredSchema,Experiment:"UP-132C-staggered-phase-selectivity",SourceUP131CSeal:"0a778b520a912ab099a054c5b103489478b67fee",FirstWaveQualified:16,SecondWaveQualified:4,PreOverflowWrites:64,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,SemanticPriorityUsed:false,QueryPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	seeds:=[]int{251000000,252000000}
	for _,arm:=range []string{"assignment_a","assignment_b"}{res.Points=append(res.Points,up132cRun(arm,seeds))}
	return res,nil
}
