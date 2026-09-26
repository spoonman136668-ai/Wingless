package unitary

const UP143CLiveMixedAgeSchema="wingless.up143c-live-mixed-age-replacement.v1"

type UP143CPoint struct{
	OlderTarget int `json:"older_target"`
	YoungerTarget int `json:"younger_target"`
	OlderPresentAfterBoundary1Rate float64 `json:"older_present_after_boundary1_rate"`
	OlderMeanAgeAfterBoundary1 float64 `json:"older_mean_age_after_boundary1"`
	YoungerPresentAfterBoundary1Rate float64 `json:"younger_present_after_boundary1_rate"`
	YoungerMeanAgeAfterBoundary1 float64 `json:"younger_mean_age_after_boundary1"`
	MeanHistoryEntriesAfterBoundary1 float64 `json:"mean_history_entries_after_boundary1"`
	OlderPresentAfterBoundary2Rate float64 `json:"older_present_after_boundary2_rate"`
	OlderMeanAgeAfterBoundary2 float64 `json:"older_mean_age_after_boundary2"`
	YoungerPresentAfterBoundary2Rate float64 `json:"younger_present_after_boundary2_rate"`
	YoungerMeanAgeAfterBoundary2 float64 `json:"younger_mean_age_after_boundary2"`
	MeanHistoryEntriesAfterBoundary2 float64 `json:"mean_history_entries_after_boundary2"`
	Age0Evictions int `json:"age0_evictions"`
	Age1Evictions int `json:"age1_evictions"`
	Age2Evictions int `json:"age2_evictions"`
	OlderAdmissionRate float64 `json:"older_admission_rate"`
	YoungerAdmissionRate float64 `json:"younger_admission_rate"`
	HotAccuracy float64 `json:"hot_accuracy"`
	UnexpectedManipulationAdmissions int `json:"unexpected_manipulation_admissions"`
	PanicRate float64 `json:"panic_rate"`
}
type UP143CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP142CSeal string `json:"source_up142c_seal"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	MaxHistoryAge int `json:"max_history_age"`
	DirectStateFixture bool `json:"direct_state_fixture"`
	MemoryIncreased bool `json:"memory_increased"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Points []UP143CPoint `json:"points"`
}
type up143cEpisode struct{
	panicHit bool
	o1p,o1a,y1p,y1a,h1 int
	o2p,o2a,y2p,y2a,h2 int
	ev [3]int
	olderAdmit,youngerAdmit int
	hotHits,hotTotal int
	unexpected int
}
func up143cEpisodeRun(older,younger,base,ep int)(out up143cEpisode){
	defer func(){if recover()!=nil{out.panicHit=true}}()
	x,truth,_,rng,maxC,maxH:=up140cSetup(older,base,ep)
	startEv:=x.evictions
	// Window 1: younger + 15 fresh paired identities.
	for rep:=0;rep<2;rep++{a,_:=up125cProcess(x,younger,truth[younger],&maxC,&maxH);if a{out.unexpected++}}
	for i:=0;i<15;i++{key:=50000+i+older*100;for rep:=0;rep<2;rep++{a,_:=up125cProcess(x,key,rng.intn(32),&maxC,&maxH);if a{out.unexpected++}}}
	if x.counter!=0{panic("UP143C_WINDOW1_NOT_BOUNDARY")}
	out.o1p,out.o1a=up140cHist(x,older);out.y1p,out.y1a=up140cHist(x,younger);out.h1=up125cCount(&x.history)
	// Window 2: 16 fresh paired identities, forcing exactly one replacement at boundary.
	for i:=0;i<16;i++{key:=60000+i+older*100;for rep:=0;rep<2;rep++{a,_:=up125cProcess(x,key,rng.intn(32),&maxC,&maxH);if a{out.unexpected++}}}
	if x.counter!=0{panic("UP143C_WINDOW2_NOT_BOUNDARY")}
	out.o2p,out.o2a=up140cHist(x,older);out.y2p,out.y2a=up140cHist(x,younger);out.h2=up125cCount(&x.history)
	for i:=0;i<3;i++{out.ev[i]=x.evictions[i]-startEv[i]}
	for k:=0;k<12;k++{got,ok:=x.mem.query(k);out.hotTotal++;if ok&&got==truth[k]{out.hotHits++}}
	xOlder:=*x;xYounger:=*x
	up125cPresent(&xOlder,older,truth[older],2,&out.olderAdmit,&maxC,&maxH)
	up125cPresent(&xYounger,younger,truth[younger],2,&out.youngerAdmit,&maxC,&maxH)
	return
}
func up143cRun(older,younger int,seeds []int)UP143CPoint{
	var total,panicN,o1p,o1a,y1p,y1a,h1,o2p,o2a,y2p,y2a,h2,olderAdmit,youngerAdmit,hotHits,hotTotal,unexpected int
	var ev [3]int
	for _,base:=range seeds{for ep:=0;ep<32;ep++{
		e:=up143cEpisodeRun(older,younger,base,ep);total++;if e.panicHit{panicN++;continue}
		o1p+=e.o1p;o1a+=e.o1a;y1p+=e.y1p;y1a+=e.y1a;h1+=e.h1;o2p+=e.o2p;o2a+=e.o2a;y2p+=e.y2p;y2a+=e.y2a;h2+=e.h2
		olderAdmit+=e.olderAdmit;youngerAdmit+=e.youngerAdmit;hotHits+=e.hotHits;hotTotal+=e.hotTotal;unexpected+=e.unexpected
		for i:=0;i<3;i++{ev[i]+=e.ev[i]}
	}}
	n:=total-panicN;rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)};age:=func(sum,p int)float64{if p==0{return 0};return float64(sum)/float64(p)}
	return UP143CPoint{OlderTarget:older,YoungerTarget:younger,OlderPresentAfterBoundary1Rate:rate(o1p,n),OlderMeanAgeAfterBoundary1:age(o1a,o1p),YoungerPresentAfterBoundary1Rate:rate(y1p,n),YoungerMeanAgeAfterBoundary1:age(y1a,y1p),MeanHistoryEntriesAfterBoundary1:rate(h1,n),OlderPresentAfterBoundary2Rate:rate(o2p,n),OlderMeanAgeAfterBoundary2:age(o2a,o2p),YoungerPresentAfterBoundary2Rate:rate(y2p,n),YoungerMeanAgeAfterBoundary2:age(y2a,y2p),MeanHistoryEntriesAfterBoundary2:rate(h2,n),Age0Evictions:ev[0],Age1Evictions:ev[1],Age2Evictions:ev[2],OlderAdmissionRate:rate(olderAdmit,n),YoungerAdmissionRate:rate(youngerAdmit,n),HotAccuracy:rate(hotHits,hotTotal),UnexpectedManipulationAdmissions:unexpected,PanicRate:rate(panicN,total)}
}
func RunUP143C()(UP143CResult,error){
	res:=UP143CResult{Schema:UP143CLiveMixedAgeSchema,Experiment:"UP-143C-live-mixed-age-replacement",SourceUP142CSeal:"33d3928831081f204fecfdc28c58d2d213b85b87",AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,MaxHistoryAge:2,DirectStateFixture:false,MemoryIncreased:false,SemanticPriorityUsed:false,FutureOracleUsed:false}
	seeds:=[]int{271000000,272000000}
	for _,pair:=range [][2]int{{100,103},{103,100}}{res.Points=append(res.Points,up143cRun(pair[0],pair[1],seeds))}
	return res,nil
}
