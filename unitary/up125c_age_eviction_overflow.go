package unitary

import "fmt"

const UP125CAgeEvictionSchema = "wingless.up125c-age-eviction-overflow.v1"

type up125cAgeEvictMachine struct {
	mem up81cAging
	current [32]uint16
	history [32]uint16
	counter int
	reuseMask uint16
	checkpoints int
	maxAge int
	evictions [3]int
}

func (x *up125cAgeEvictMachine) insertHistory(key,age int) {
	code:=up118cCode(key)|(uint16(age&3)<<14)
	if i:=up118cFind(&x.history,key);i>=0 { x.history[i]=code;return }
	for i,v:=range x.history {
		if v==0 { x.history[i]=code;return }
	}
	bestIdx,bestAge:=-1,-1
	for i,v:=range x.history {
		a:=int(v>>14)
		if a>bestAge { bestIdx=i;bestAge=a }
	}
	if bestIdx<0 { panic("UP125C age eviction found no victim") }
	if bestAge>=0&&bestAge<3 { x.evictions[bestAge]++ }
	x.history[bestIdx]=code
}

func (x *up125cAgeEvictMachine) boundary() {
	for i,v:=range x.history {
		if v==0 { continue }
		age:=int(v>>14)
		if age>=x.maxAge { x.history[i]=0 } else { x.history[i]=(v&up118cKeyMask)|(uint16(age+1)<<14) }
	}
	for _,v:=range x.current {
		if v&up118cCurrentQualified!=0 { x.insertHistory(up118cKey(v),0) }
	}
	up118cClear(&x.current)
	x.counter=0
}
func (x *up125cAgeEvictMachine) checkpoint() {
	up118cClear(&x.current);up118cClear(&x.history)
	x.counter=8;x.reuseMask=0;x.checkpoints++
}
func (x *up125cAgeEvictMachine) process(key,value int)(admitted,rejected bool) {
	existing:=x.mem.find(key)
	if existing>=0 {
		x.mem.write(key,value)
		x.reuseMask|=uint16(1)<<uint(existing)
		if x.reuseMask==0xffff { x.checkpoint() }
		return true,false
	}
	if x.mem.count<16 { x.mem.write(key,value);return true,false }

	idx:=up118cFind(&x.current,key)
	if idx<0 {
		up118cInsertCurrent(&x.current,key,false);rejected=true
	} else if x.current[idx]&up118cCurrentQualified==0 {
		x.current[idx]|=up118cCurrentQualified
		if hi:=up118cFind(&x.history,key);hi>=0 {
			x.mem.write(key,value);admitted=true;x.current[idx]=0;x.history[hi]=0
		} else { rejected=true }
	} else { rejected=true }

	x.counter++
	if x.counter==32 { x.boundary() }
	if admitted {
		idx=x.mem.find(key)
		if idx>=0 { x.reuseMask|=uint16(1)<<uint(idx) }
		if x.reuseMask==0xffff { x.checkpoint() }
	}
	return
}
func up125cCount(table *[32]uint16) int {
	n:=0;for _,v:=range table { if v!=0 { n++ } };return n
}
func up125cObserve(x *up125cAgeEvictMachine,maxCurrent,maxHistory *int) {
	c:=up125cCount(&x.current);h:=up125cCount(&x.history)
	if c>*maxCurrent { *maxCurrent=c };if h>*maxHistory { *maxHistory=h }
}
func up125cProcess(x *up125cAgeEvictMachine,key,value int,maxCurrent,maxHistory *int)(bool,bool) {
	a,r:=x.process(key,value);up125cObserve(x,maxCurrent,maxHistory);return a,r
}
func up125cFill(x *up125cAgeEvictMachine,next *int,rng *sq0RNG,maxCurrent,maxHistory *int)(fp int) {
	for x.counter!=0 {
		k:=*next;(*next)++;a,_:=up125cProcess(x,k,rng.intn(32),maxCurrent,maxHistory);if a{fp++}
	}
	return
}
func up125cPresent(x *up125cAgeEvictMachine,key,value,sightings int,admit *int,maxCurrent,maxHistory *int) {
	for i:=0;i<sightings;i++ {
		was:=x.mem.find(key)>=0;a,_:=up125cProcess(x,key,value,maxCurrent,maxHistory)
		if !was&&a { (*admit)++ }
	}
}

type up125cEpisode struct {
	panicHit bool
	panicMessage string
	validAdmit,validTrials,validHits,validTotal int
	hotHits,hotTotal,targetHits,targetTotal,exactHits int
	nonAdmit,nonTrials,nonHits,nonTotal int
	fp,maxCurrent,maxHistory,recallEntries,checkpoints int
	evictions [3]int
}

func up125cEpisodeRun(q,base,ep int)(out up125cEpisode) {
	defer func(){if v:=recover();v!=nil{out.panicHit=true;out.panicMessage=fmt.Sprint(v)}}()
	rng:=newSQ0RNG(sq0Seed(base,6601+q*251,ep))
	x:=&up125cAgeEvictMachine{maxAge:2}
	truth:=map[int]int{}
	next:=300
	targets:=[]int{100,101,102,103}
	setA:=make([]int,16);for i:=range setA{setA[i]=200+i}
	setB:=make([]int,12);for i:=range setB{setB[i]=220+i}

	for k:=0;k<32;k++ {
		v:=rng.intn(32);truth[k]=v;a,_:=up125cProcess(x,k,v,&out.maxCurrent,&out.maxHistory);if k>=16&&a{out.fp++}
	}
	for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }
	up125cObserve(x,&out.maxCurrent,&out.maxHistory)
	if x.counter!=0 { out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory) }

	for _,k:=range setA {
		v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++
	}
	for _,k:=range targets {
		v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.validAdmit,&out.maxCurrent,&out.maxHistory);out.validTrials++
	}
	for _,k:=range setB {
		v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++
	}
	if x.counter!=0 { out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory) }

	overflow:=make([]int,q);for i:=0;i<q;i++{overflow[i]=240+i}
	for _,k:=range overflow {
		v:=rng.intn(32);truth[k]=v;up125cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory);out.nonTrials++
	}
	if x.counter!=0 { out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory) }

	for _,k:=range targets { up125cPresent(x,k,truth[k],2,&out.validAdmit,&out.maxCurrent,&out.maxHistory) }
	if x.counter!=0 { out.fp+=up125cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory) }

	for j:=0;j<12288;j++ {
		key:=1000+j;a,_:=up125cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory);if a{out.fp++}
		if (j+1)%4==0 {
			for k:=0;k<12;k++ { x.mem.query(k) }
			for _,g:=range [][]int{targets,setA,setB,overflow} { for _,k:=range g { x.mem.query(k) } }
		}
	}
	exact:=true
	for k:=0;k<12;k++ {
		got,ok:=x.mem.query(k);out.hotTotal++;out.targetTotal++
		if ok&&got==truth[k] { out.hotHits++;out.targetHits++ } else { exact=false }
	}
	for _,k:=range targets {
		got,ok:=x.mem.query(k);out.validTotal++;out.targetTotal++
		if ok&&got==truth[k] { out.validHits++;out.targetHits++ } else { exact=false }
	}
	for _,g:=range [][]int{setA,setB,overflow} {
		for _,k:=range g {
			got,ok:=x.mem.query(k);out.nonTotal++;if ok&&got==truth[k]{out.nonHits++}
		}
	}
	if exact { out.exactHits=1 }
	out.recallEntries=x.mem.count;out.checkpoints=x.checkpoints;out.evictions=x.evictions
	return
}

type UP125CPoint struct {
	OverflowQualifiedCandidates int `json:"overflow_qualified_candidates"`
	PanicRate float64 `json:"panic_rate"`
	PanicMessage string `json:"panic_message"`
	TotalEvictions int `json:"total_evictions"`
	Age0Evictions int `json:"age0_evictions"`
	Age1Evictions int `json:"age1_evictions"`
	Age2Evictions int `json:"age2_evictions"`
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
	MeanCheckpointsPerEpisode float64 `json:"mean_checkpoints_per_episode"`
}
type UP125CAgeEvictionResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP124CSeal string `json:"source_up124c_seal"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	EpisodesPerSeed int `json:"episodes_per_seed"`
	ReplacementSignal string `json:"replacement_signal"`
	TieBreak string `json:"tie_break"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	QueryPriorityUsed bool `json:"query_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP125CPoint `json:"points"`
}
func up125cRun(q int,seeds []int) UP125CPoint {
	var panicN,total int;msg:=""
	var validAdmit,validTrials,validHits,validTotal,hotHits,hotTotal,targetHits,targetTotal,exactHits int
	var nonAdmit,nonTrials,nonHits,nonTotal,fp,maxCurrent,maxHistory,maxRecall,checkpoints int
	var ev [3]int
	for _,base:=range seeds { for ep:=0;ep<32;ep++ {
		e:=up125cEpisodeRun(q,base,ep);total++
		if e.maxCurrent>maxCurrent{maxCurrent=e.maxCurrent};if e.maxHistory>maxHistory{maxHistory=e.maxHistory};if e.recallEntries>maxRecall{maxRecall=e.recallEntries}
		fp+=e.fp;nonAdmit+=e.nonAdmit;nonTrials+=e.nonTrials;for i:=0;i<3;i++{ev[i]+=e.evictions[i]}
		if e.panicHit { panicN++;if msg==""{msg=e.panicMessage};continue }
		validAdmit+=e.validAdmit;validTrials+=e.validTrials;validHits+=e.validHits;validTotal+=e.validTotal
		hotHits+=e.hotHits;hotTotal+=e.hotTotal;targetHits+=e.targetHits;targetTotal+=e.targetTotal;exactHits+=e.exactHits
		nonHits+=e.nonHits;nonTotal+=e.nonTotal;checkpoints+=e.checkpoints
	}}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	return UP125CPoint{
		OverflowQualifiedCandidates:q,PanicRate:rate(panicN,total),PanicMessage:msg,
		TotalEvictions:ev[0]+ev[1]+ev[2],Age0Evictions:ev[0],Age1Evictions:ev[1],Age2Evictions:ev[2],
		ValidAdmissionRate:rate(validAdmit,validTrials),ValidAccuracy:rate(validHits,validTotal),HotAccuracy:rate(hotHits,hotTotal),
		Target16Accuracy:rate(targetHits,targetTotal),Target16ExactAccuracy:rate(exactHits,total-panicN),
		NonpersistentAdmissionRate:rate(nonAdmit,nonTrials),NonpersistentRetentionRate:rate(nonHits,nonTotal),
		OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,
		RecallEntriesUsed:maxRecall,MeanCheckpointsPerEpisode:rate(checkpoints,total-panicN),
	}
}
func RunUP125C()(UP125CAgeEvictionResult,error) {
	res:=UP125CAgeEvictionResult{
		Schema:UP125CAgeEvictionSchema,Experiment:"UP-125C-age-eviction-overflow",
		SourceUP124CSeal:"b78728a9f3247f5c09f88f8b5f42b20c33b465a6",
		AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,EpisodesPerSeed:32,
		ReplacementSignal:"maximum_history_age_only",TieBreak:"lowest_table_index",
		SemanticPriorityUsed:false,QueryPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false,
	}
	seeds:=[]int{237000000,238000000}
	for _,q:=range []int{0,1,8,16}{res.Points=append(res.Points,up125cRun(q,seeds))}
	return res,nil
}
