package unitary

const UP144CAgeStackSchema="wingless.up144c-live-age-stack-replacement.v1"

type UP144CTargetState struct{
	Target int `json:"target"`
	Present bool `json:"present"`
	Age int `json:"age"`
}
type UP144CPoint struct{
	Assignment string `json:"assignment"`
	Oldest1 int `json:"oldest1"`
	Oldest2 int `json:"oldest2"`
	Younger1 int `json:"younger1"`
	Younger2 int `json:"younger2"`
	AfterBoundary1 []UP144CTargetState `json:"after_boundary1"`
	AfterBoundary2 []UP144CTargetState `json:"after_boundary2"`
	AfterBoundary3 []UP144CTargetState `json:"after_boundary3"`
	MeanHistoryEntriesBoundary1 float64 `json:"mean_history_entries_boundary1"`
	MeanHistoryEntriesBoundary2 float64 `json:"mean_history_entries_boundary2"`
	MeanHistoryEntriesBoundary3 float64 `json:"mean_history_entries_boundary3"`
	Age0Evictions int `json:"age0_evictions"`
	Age1Evictions int `json:"age1_evictions"`
	Age2Evictions int `json:"age2_evictions"`
	Oldest1AdmissionRate float64 `json:"oldest1_admission_rate"`
	Oldest2AdmissionRate float64 `json:"oldest2_admission_rate"`
	Younger1AdmissionRate float64 `json:"younger1_admission_rate"`
	Younger2AdmissionRate float64 `json:"younger2_admission_rate"`
	HotAccuracy float64 `json:"hot_accuracy"`
	UnexpectedManipulationAdmissions int `json:"unexpected_manipulation_admissions"`
	PanicRate float64 `json:"panic_rate"`
}
type UP144CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP143CSeal string `json:"source_up143c_seal"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	MaxHistoryAge int `json:"max_history_age"`
	DirectStateFixture bool `json:"direct_state_fixture"`
	MemoryIncreased bool `json:"memory_increased"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Points []UP144CPoint `json:"points"`
}
type up144cEpisode struct{
	panicHit bool
	b1,b2,b3 [4]UP144CTargetState
	h1,h2,h3 int
	ev [3]int
	admit [4]int
	hotHits,hotTotal,unexpected int
}
func up144cStates(x *up125cAgeEvictMachine,targets [4]int)[4]UP144CTargetState{
	var out [4]UP144CTargetState
	for i,t:=range targets{p,a:=up140cHist(x,t);out[i]=UP144CTargetState{Target:t,Present:p==1,Age:a}}
	return out
}
func up144cEpisodeRun(targets [4]int,base,ep int)(out up144cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	o1,o2,y1,y2:=targets[0],targets[1],targets[2],targets[3]
	x,truth,_,rng,maxC,maxH:=up140cSetup(o1,base,ep);startEv:=x.evictions
	// W1: oldest2 first, then 15 paired distractors.
	for rep:=0;rep<2;rep++{a,_:=up125cProcess(x,o2,truth[o2],&maxC,&maxH);if a{out.unexpected++}}
	for i:=0;i<15;i++{key:=70000+i+o1*100;for rep:=0;rep<2;rep++{a,_:=up125cProcess(x,key,rng.intn(32),&maxC,&maxH);if a{out.unexpected++}}}
	if x.counter!=0{panic("UP144C_W1_NOT_BOUNDARY")};out.b1=up144cStates(x,targets);out.h1=up125cCount(&x.history)
	// W2: younger1, younger2, then 14 paired distractors.
	for _,t:=range []int{y1,y2}{for rep:=0;rep<2;rep++{a,_:=up125cProcess(x,t,truth[t],&maxC,&maxH);if a{out.unexpected++}}}
	for i:=0;i<14;i++{key:=80000+i+o1*100;for rep:=0;rep<2;rep++{a,_:=up125cProcess(x,key,rng.intn(32),&maxC,&maxH);if a{out.unexpected++}}}
	if x.counter!=0{panic("UP144C_W2_NOT_BOUNDARY")};out.b2=up144cStates(x,targets);out.h2=up125cCount(&x.history)
	// W3: one qualified fresh key + 30 singletons.
	key:=90000+o1*100
	for rep:=0;rep<2;rep++{a,_:=up125cProcess(x,key,rng.intn(32),&maxC,&maxH);if a{out.unexpected++}}
	for i:=0;i<30;i++{k:=91000+i+o1*100;a,_:=up125cProcess(x,k,rng.intn(32),&maxC,&maxH);if a{out.unexpected++}}
	if x.counter!=0{panic("UP144C_W3_NOT_BOUNDARY")};out.b3=up144cStates(x,targets);out.h3=up125cCount(&x.history)
	for i:=0;i<3;i++{out.ev[i]=x.evictions[i]-startEv[i]}
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;if ok&&got==truth[k]{out.hotHits++}}
	for i,t:=range targets{clone:=*x;dummyC,dummyH:=maxC,maxH;up125cPresent(&clone,t,truth[t],2,&out.admit[i],&dummyC,&dummyH)}
	return
}
func up144cAggregateState(episodes []up144cEpisode,which int,idx int,target int)UP144CTargetState{
	present,age,n:=0,0,0
	for _,e:=range episodes{if e.panicHit{continue};var s UP144CTargetState;if which==1{s=e.b1[idx]}else if which==2{s=e.b2[idx]}else{s=e.b3[idx]};if s.Present{present++;age+=s.Age};n++}
	a:=0;if present>0{a=age/present};return UP144CTargetState{Target:target,Present:present==n,Age:a}
}
func up144cRun(name string,targets [4]int,seeds []int)UP144CPoint{
	episodes:=[]up144cEpisode{};total,panicN:=0,0
	for _,base:=range seeds{for ep:=0;ep<32;ep++{e:=up144cEpisodeRun(targets,base,ep);episodes=append(episodes,e);total++;if e.panicHit{panicN++}}}
	n:=total-panicN;rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	var h1,h2,h3,hotHits,hotTotal,unexpected int;var ev [3]int;var admit [4]int
	for _,e:=range episodes{if e.panicHit{continue};h1+=e.h1;h2+=e.h2;h3+=e.h3;hotHits+=e.hotHits;hotTotal+=e.hotTotal;unexpected+=e.unexpected;for i:=0;i<3;i++{ev[i]+=e.ev[i]};for i:=0;i<4;i++{admit[i]+=e.admit[i]}}
	b1,b2,b3:=[]UP144CTargetState{},[]UP144CTargetState{},[]UP144CTargetState{}
	for i,t:=range targets{b1=append(b1,up144cAggregateState(episodes,1,i,t));b2=append(b2,up144cAggregateState(episodes,2,i,t));b3=append(b3,up144cAggregateState(episodes,3,i,t))}
	return UP144CPoint{Assignment:name,Oldest1:targets[0],Oldest2:targets[1],Younger1:targets[2],Younger2:targets[3],AfterBoundary1:b1,AfterBoundary2:b2,AfterBoundary3:b3,MeanHistoryEntriesBoundary1:rate(h1,n),MeanHistoryEntriesBoundary2:rate(h2,n),MeanHistoryEntriesBoundary3:rate(h3,n),Age0Evictions:ev[0],Age1Evictions:ev[1],Age2Evictions:ev[2],Oldest1AdmissionRate:rate(admit[0],n),Oldest2AdmissionRate:rate(admit[1],n),Younger1AdmissionRate:rate(admit[2],n),Younger2AdmissionRate:rate(admit[3],n),HotAccuracy:rate(hotHits,hotTotal),UnexpectedManipulationAdmissions:unexpected,PanicRate:rate(panicN,total)}
}
func RunUP144C()(UP144CResult,error){
	res:=UP144CResult{Schema:UP144CAgeStackSchema,Experiment:"UP-144C-live-age-stack-replacement",SourceUP143CSeal:"35572f35815a5acbafd86b02fba9b5cae930e383",AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,MaxHistoryAge:2,DirectStateFixture:false,MemoryIncreased:false,FutureOracleUsed:false}
	seeds:=[]int{273000000,274000000}
	res.Points=append(res.Points,up144cRun("A",[4]int{100,101,102,103},seeds),up144cRun("B",[4]int{102,103,100,101},seeds))
	return res,nil
}
