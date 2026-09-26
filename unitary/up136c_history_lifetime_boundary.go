package unitary

const UP136CLifetimeSchema="wingless.up136c-history-lifetime-boundary.v1"

type UP136CPoint struct{
	Path string `json:"path"`
	Target int `json:"target"`
	Positions []int `json:"positions"`
	SecondPairGeneration int `json:"second_pair_generation"`
	AdmissionRate float64 `json:"admission_rate"`
	FinalAccuracy float64 `json:"final_accuracy"`
	HotAccuracy float64 `json:"hot_accuracy"`
	NonpersistentAdmissionRate float64 `json:"nonpersistent_admission_rate"`
	NonpersistentRetentionRate float64 `json:"nonpersistent_retention_rate"`
	OneShotFalseAdmissions int `json:"one_shot_false_admissions"`
	MaxCurrentTableEntries int `json:"max_current_table_entries"`
	MaxHistoryTableEntries int `json:"max_history_table_entries"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	PanicRate float64 `json:"panic_rate"`
}
type UP136CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP135CSeal string `json:"source_up135c_seal"`
	Generations int `json:"generations"`
	ScheduledWrites int `json:"scheduled_writes"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	MaxHistoryAge int `json:"max_history_age"`
	SingleRecurrenceTargetPerArm bool `json:"single_recurrence_target_per_arm"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	QueryPriorityUsed bool `json:"query_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP136CPoint `json:"points"`
}
type up136cEpisode struct{
	panicHit bool
	admit,hits int
	hotHits,hotTotal int
	nonAdmit,nonTrials,nonHits,nonTotal int
	fp,maxCurrent,maxHistory,recallEntries int
}
func up136cPositions(path string)[]int{
	switch path{
	case "live_pair_g1":return []int{15,16}
	case "rebuild_only":return []int{5,37,74,75}
	case "next_g4":return []int{5,37,74,75,106,107}
	case "after_one_blank_g5":return []int{5,37,74,75,138,139}
	case "after_two_blank_g6":return []int{5,37,74,75,170,171}
	default:return []int{5,37,74,75,202,203}
	}
}
func up136cSecondGen(path string)int{
	switch path{
	case "next_g4":return 4
	case "after_one_blank_g5":return 5
	case "after_two_blank_g6":return 6
	case "after_three_blank_g7":return 7
	default:return 0
	}
}
func up136cEpisodeRun(path string,target,base,ep int)(out up136cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	pathCode:=map[string]int{"live_pair_g1":0,"rebuild_only":1,"next_g4":2,"after_one_blank_g5":3,"after_two_blank_g6":4,"after_three_blank_g7":5}[path]
	rng:=newSQ0RNG(sq0Seed(base,7701+pathCode*317+target*5,ep))
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

	positions:=up136cPositions(path)
	for pos:=1;pos<=224;pos++{
		hit:=false
		for _,at:=range positions{
			if pos==at{
				up125cPresent(x,target,truth[target],1,&out.admit,&out.maxCurrent,&out.maxHistory)
				hit=true
				break
			}
		}
		if !hit{a,_:=up125cProcess(x,next,rng.intn(32),&out.maxCurrent,&out.maxHistory);next++;if a{out.fp++}}
	}

	second:=[]int{260,261,262,263}
	for _,k:=range second{v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++}
	if x.counter!=0{out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory)}
	for i:=0;i<32;i++{a,_:=up125cProcess(x,next,rng.intn(32),&out.maxCurrent,&out.maxHistory);next++;if a{out.fp++}}

	for j:=0;j<12288;j++{
		key:=1000+j;a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory);if a{out.fp++}
		if (j+1)%4==0{for k:=0;k<12;k++{x.mem.query(k)};for _,g:=range [][]int{targets,setA,setB,first,second}{for _,k:=range g{x.mem.query(k)}}}
	}
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;if ok&&got==truth[k]{out.hotHits++}}
	if got,ok:=x.mem.query(target);ok&&got==truth[target]{out.hits++}
	for _,g:=range [][]int{setA,setB,first,second}{for _,k:=range g{got,ok:=x.mem.query(k);out.nonTotal++;if ok&&got==truth[k]{out.nonHits++}}}
	out.recallEntries=x.mem.count
	return
}
func up136cRun(path string,target int,seeds []int)UP136CPoint{
	var panicN,total,admit,hits,hotHits,hotTotal,nonAdmit,nonTrials,nonHits,nonTotal,fp,maxCurrent,maxHistory,maxRecall int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up136cEpisodeRun(path,target,base,ep);total++;if e.panicHit{panicN++;continue}
		admit+=e.admit;hits+=e.hits;hotHits+=e.hotHits;hotTotal+=e.hotTotal;nonAdmit+=e.nonAdmit;nonTrials+=e.nonTrials;nonHits+=e.nonHits;nonTotal+=e.nonTotal;fp+=e.fp
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
	}}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	pos:=up136cPositions(path)
	return UP136CPoint{Path:path,Target:target,Positions:pos,SecondPairGeneration:up136cSecondGen(path),AdmissionRate:rate(admit,total),FinalAccuracy:rate(hits,total),HotAccuracy:rate(hotHits,hotTotal),NonpersistentAdmissionRate:rate(nonAdmit,nonTrials),NonpersistentRetentionRate:rate(nonHits,nonTotal),OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,RecallEntriesUsed:maxRecall,PanicRate:rate(panicN,total)}
}
func RunUP136C()(UP136CResult,error){
	res:=UP136CResult{Schema:UP136CLifetimeSchema,Experiment:"UP-136C-history-lifetime-boundary",SourceUP135CSeal:"6dad48c4b5522a2991a600b33757522dc9bd2ea1",Generations:7,ScheduledWrites:224,AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,MaxHistoryAge:2,SingleRecurrenceTargetPerArm:true,SemanticPriorityUsed:false,QueryPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	seeds:=[]int{259000000,260000000}
	paths:=[]string{"live_pair_g1","rebuild_only","next_g4","after_one_blank_g5","after_two_blank_g6","after_three_blank_g7"}
	for _,target:=range []int{100,103}{for _,path:=range paths{res.Points=append(res.Points,up136cRun(path,target,seeds))}}
	return res,nil
}
