package unitary

const UP146CAccessRecencySchema="wingless.up146c-durable-access-recency.v1"

type UP146CDurableMetric struct{
	Key int `json:"key"`
	Subject string `json:"subject"`
	Verb string `json:"verb"`
	Class int `json:"class"`
	Queried bool `json:"queried"`
	Retained bool `json:"retained"`
}
type UP146CPoint struct{
	Assignment string `json:"assignment"`
	ColdBlock string `json:"cold_block"`
	ColdIndices []int `json:"cold_indices"`
	HotDurableAccuracy float64 `json:"hot_durable_accuracy"`
	ColdDurableAccuracy float64 `json:"cold_durable_accuracy"`
	Durable []UP146CDurableMetric `json:"durable"`
	Candidates []UP145CCandidateMetric `json:"candidates"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	MaxCurrentTableEntries int `json:"max_current_table_entries"`
	MaxHistoryTableEntries int `json:"max_history_table_entries"`
	OneShotFalseAdmissions int `json:"one_shot_false_admissions"`
	PanicRate float64 `json:"panic_rate"`
}
type UP146CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP145CSeal string `json:"source_up145c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	ColdFactsPerArm int `json:"cold_facts_per_arm"`
	QueriedFactsPerArm int `json:"queried_facts_per_arm"`
	Assignments int `json:"assignments"`
	ColdBlocks int `json:"cold_blocks"`
	MemoryPolicyChanged bool `json:"memory_policy_changed"`
	AdaptiveRefreshUsed bool `json:"adaptive_refresh_used"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	QueryPriorityUsed bool `json:"query_priority_used"`
	Points []UP146CPoint `json:"points"`
}
func up146cColdName(start int)string{
	switch start{case 0:return "cold_0_3";case 4:return "cold_4_7";case 8:return "cold_8_11";default:return "cold_12_15"}
}
func up146cIsCold(k,start int)bool{return k>=start&&k<start+4}
func up146cRun(assignment string,coldStart int)UP146CPoint{
	x:=&up125cAgeEvictMachine{maxAge:2};durable:=up145cDurable();cands:=up145cCandidates();roles:=up145cRoles(assignment)
	truth:=map[int]int{};for _,d:=range durable{truth[d.key]=d.class;x.mem.write(d.key,d.class)};for _,c:=range cands{truth[c.key]=c.class}
	admit:=map[int]int{};maxC,maxH,fp:=0,0,0;panicHit:=false
	func(){
		defer func(){if recover()!=nil{panicHit=true}}()
		next:=1000
		for pos:=1;pos<=96;pos++{
			if c,ok:=up145cAt(pos,cands,roles);ok{
				was:=x.mem.find(c.key)>=0;a,_:=up125cProcess(x,c.key,c.class,&maxC,&maxH);if !was&&a{admit[c.key]++}
			}else{
				a,_:=up125cProcess(x,next,(pos-1)%3,&maxC,&maxH);next++;if a{fp++}
			}
			if pos%4==0{for k:=0;k<16;k++{if !up146cIsCold(k,coldStart){x.mem.query(k)}}}
		}
	}()
	p:=UP146CPoint{Assignment:assignment,ColdBlock:up146cColdName(coldStart),ColdIndices:[]int{coldStart,coldStart+1,coldStart+2,coldStart+3},RecallEntriesUsed:x.mem.count,MaxCurrentTableEntries:maxC,MaxHistoryTableEntries:maxH,OneShotFalseAdmissions:fp}
	if panicHit{p.PanicRate=1;return p}
	hotHits,coldHits:=0,0
	for _,d:=range durable{
		v,ok:=x.mem.query(d.key);retained:=ok&&v==truth[d.key];queried:=!up146cIsCold(d.key,coldStart)
		if queried{if retained{hotHits++}}else{if retained{coldHits++}}
		p.Durable=append(p.Durable,UP146CDurableMetric{Key:d.key,Subject:d.subject,Verb:d.verb,Class:d.class,Queried:queried,Retained:retained})
	}
	p.HotDurableAccuracy=float64(hotHits)/12.0;p.ColdDurableAccuracy=float64(coldHits)/4.0
	for _,c:=range cands{_,ok:=x.mem.query(c.key);p.Candidates=append(p.Candidates,UP145CCandidateMetric{Subject:c.subject,Verb:c.verb,Class:c.class,Role:roles[c.key],Positions:append([]int(nil),up145cPositions(roles[c.key])...),AdmissionCount:admit[c.key],FinalRetained:ok})}
	return p
}
func RunUP146C()(UP146CResult,error){
	res:=UP146CResult{Schema:UP146CAccessRecencySchema,Experiment:"UP-146C-durable-access-recency",SourceUP145CSeal:"365893cb14058ecded93d22719185878ffd08316",ExactRecallCap:16,HistoryEntries:32,ColdFactsPerArm:4,QueriedFactsPerArm:12,Assignments:2,ColdBlocks:4,MemoryPolicyChanged:false,AdaptiveRefreshUsed:false,SemanticPriorityUsed:false,QueryPriorityUsed:false}
	for _,a:=range []string{"A","B"}{for _,start:=range []int{0,4,8,12}{res.Points=append(res.Points,up146cRun(a,start))}}
	return res,nil
}
