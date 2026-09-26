package unitary

const UP134CMultiGenerationSchema="wingless.up134c-multigeneration-reset.v1"

type UP134CPathMetric struct{
	Path string `json:"path"`
	Target int `json:"target"`
	Positions []int `json:"positions"`
	TotalSightings int `json:"total_sightings"`
	AdmissionRate float64 `json:"admission_rate"`
	FinalAccuracy float64 `json:"final_accuracy"`
}
type UP134CPoint struct{
	Arm string `json:"arm"`
	PathMetrics []UP134CPathMetric `json:"path_metrics"`
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
type UP134CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP133CSeal string `json:"source_up133c_seal"`
	FirstWaveQualified int `json:"first_wave_qualified"`
	SecondWaveQualified int `json:"second_wave_qualified"`
	PreOverflowWrites int `json:"pre_overflow_writes"`
	Generations int `json:"generations"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	QueryPriorityUsed bool `json:"query_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP134CPoint `json:"points"`
}
type up134cPath struct{name string;target int;positions []int}
type up134cEpisode struct{
	panicHit bool
	admit [4]int
	hits [4]int
	hotHits,hotTotal,targetHits,targetTotal int
	exact bool
	nonAdmit,nonTrials,nonHits,nonTotal int
	fp,maxCurrent,maxHistory,recallEntries int
}
func up134cPaths(arm string)[]up134cPath{
	targets:=[]int{100,101,102,103}
	if arm=="assignment_b"{targets=[]int{102,103,100,101}}
	return []up134cPath{
		{"isolated_4",targets[0],[]int{5,37,69,101}},
		{"late_pair_g3",targets[1],[]int{10,42,74,75}},
		{"early_pair_then_singles",targets[2],[]int{15,16,50,82}},
		{"pair_after_three_isolated",targets[3],[]int{20,52,84,116,117}},
	}
}
func up134cEpisodeRun(arm string,base,ep int)(out up134cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	code:=0;if arm=="assignment_b"{code=1}
	rng:=newSQ0RNG(sq0Seed(base,7501+code*311,ep))
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

	paths:=up134cPaths(arm)
	for pos:=1;pos<=128;pos++{
		handled:=false
		for i,p:=range paths{
			for _,at:=range p.positions{
				if pos==at{
					up125cPresent(x,p.target,truth[p.target],1,&out.admit[i],&out.maxCurrent,&out.maxHistory)
					handled=true
					break
				}
			}
			if handled{break}
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
func up134cRun(arm string,seeds []int)UP134CPoint{
	var panicN,total,hotHits,hotTotal,targetHits,targetTotal,exactHits int
	var admit,hits [4]int
	var nonAdmit,nonTrials,nonHits,nonTotal,fp,maxCurrent,maxHistory,maxRecall int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up134cEpisodeRun(arm,base,ep);total++;if e.panicHit{panicN++;continue}
		for i:=0;i<4;i++{admit[i]+=e.admit[i];hits[i]+=e.hits[i]}
		hotHits+=e.hotHits;hotTotal+=e.hotTotal;targetHits+=e.targetHits;targetTotal+=e.targetTotal;if e.exact{exactHits++}
		nonAdmit+=e.nonAdmit;nonTrials+=e.nonTrials;nonHits+=e.nonHits;nonTotal+=e.nonTotal;fp+=e.fp
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
	}}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	paths:=up134cPaths(arm)
	pm:=make([]UP134CPathMetric,0,4)
	for i,p:=range paths{pm=append(pm,UP134CPathMetric{Path:p.name,Target:p.target,Positions:append([]int(nil),p.positions...),TotalSightings:len(p.positions),AdmissionRate:rate(admit[i],total),FinalAccuracy:rate(hits[i],total)})}
	return UP134CPoint{Arm:arm,PathMetrics:pm,HotAccuracy:rate(hotHits,hotTotal),Target16Accuracy:rate(targetHits,targetTotal),Target16ExactAccuracy:rate(exactHits,total-panicN),NonpersistentAdmissionRate:rate(nonAdmit,nonTrials),NonpersistentRetentionRate:rate(nonHits,nonTotal),OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall,PanicRate:rate(panicN,total)}
}
func RunUP134C()(UP134CResult,error){
	res:=UP134CResult{Schema:UP134CMultiGenerationSchema,Experiment:"UP-134C-multigeneration-reset",SourceUP133CSeal:"e84d69e7daa206134e909ce111ceeee5e9297a7c",FirstWaveQualified:16,SecondWaveQualified:4,PreOverflowWrites:128,Generations:4,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,SemanticPriorityUsed:false,QueryPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	seeds:=[]int{255000000,256000000}
	for _,arm:=range []string{"assignment_a","assignment_b"}{res.Points=append(res.Points,up134cRun(arm,seeds))}
	return res,nil
}
