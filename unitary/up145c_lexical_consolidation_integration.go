package unitary

const UP145CLexicalIntegrationSchema="wingless.up145c-lexical-consolidation-integration.v1"

type up145cLexeme struct{key int;subject,verb string;class int}
type UP145CCandidateMetric struct{
	Subject string `json:"subject"`
	Verb string `json:"verb"`
	Class int `json:"class"`
	Role string `json:"role"`
	Positions []int `json:"positions"`
	AdmissionCount int `json:"admission_count"`
	FinalRetained bool `json:"final_retained"`
}
type UP145CPoint struct{
	Assignment string `json:"assignment"`
	Candidates []UP145CCandidateMetric `json:"candidates"`
	HotDurableAccuracy float64 `json:"hot_durable_accuracy"`
	ColdDurableAccuracy float64 `json:"cold_durable_accuracy"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	MaxCurrentTableEntries int `json:"max_current_table_entries"`
	MaxHistoryTableEntries int `json:"max_history_table_entries"`
	OneShotFalseAdmissions int `json:"one_shot_false_admissions"`
	PanicRate float64 `json:"panic_rate"`
}
type UP145CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP144CSeal string `json:"source_up144c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	HistoryEntries int `json:"history_entries"`
	AdmissionMemoryBytes int `json:"admission_memory_bytes"`
	CandidateGenerations int `json:"candidate_generations"`
	CandidateCalls int `json:"candidate_calls"`
	HotDurableFacts int `json:"hot_durable_facts"`
	ColdDurableFacts int `json:"cold_durable_facts"`
	LexicalClasses []string `json:"lexical_classes"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	QueryPriorityUsed bool `json:"query_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP145CPoint `json:"points"`
}
func up145cDurable()[]up145cLexeme{
	return []up145cLexeme{
		{0,"quin","lodges",0},{1,"quin","stashes",0},{2,"quin","caches",0},{3,"quin","files",0},
		{4,"quin","scans",1},{5,"quin","checks",1},{6,"quin","views",1},{7,"quin","monitors",1},
		{8,"quin","relays",2},{9,"quin","announces",2},{10,"quin","cites",2},{11,"quin","summarizes",2},
		{12,"rue","lodges",0},{13,"rue","stashes",0},{14,"rue","caches",0},{15,"rue","files",0},
	}
}
func up145cCandidates()[]up145cLexeme{
	return []up145cLexeme{{100,"mia","caches",0},{101,"noah","scans",1},{102,"opal","relays",2},{103,"pax","files",0}}
}
func up145cRoles(assignment string)map[int]string{
	if assignment=="B"{return map[int]string{102:"two_stage",103:"cross_boundary",100:"delayed_two_stage",101:"isolated"}}
	return map[int]string{100:"two_stage",101:"cross_boundary",102:"delayed_two_stage",103:"isolated"}
}
func up145cPositions(role string)[]int{
	switch role{
	case "two_stage":return []int{5,6,37,38}
	case "cross_boundary":return []int{31,33}
	case "delayed_two_stage":return []int{15,16,79,80}
	default:return []int{20,52,84}
	}
}
func up145cAt(pos int,cands []up145cLexeme,roles map[int]string)(up145cLexeme,bool){
	for _,c:=range cands{for _,p:=range up145cPositions(roles[c.key]){if p==pos{return c,true}}}
	return up145cLexeme{},false
}
func up145cRun(assignment string)UP145CPoint{
	x:=&up125cAgeEvictMachine{maxAge:2}
	durable:=up145cDurable();cands:=up145cCandidates();roles:=up145cRoles(assignment)
	truth:=map[int]int{}
	for _,d:=range durable{truth[d.key]=d.class;x.mem.write(d.key,d.class)}
	for _,c:=range cands{truth[c.key]=c.class}
	admit:=map[int]int{};maxC,maxH,fp:=0,0,0
	panicHit:=false
	func(){
		defer func(){if recover()!=nil{panicHit=true}}()
		next:=1000
		for pos:=1;pos<=96;pos++{
			if c,ok:=up145cAt(pos,cands,roles);ok{
				was:=x.mem.find(c.key)>=0
				a,_:=up125cProcess(x,c.key,c.class,&maxC,&maxH)
				if !was&&a{admit[c.key]++}
			}else{
				a,_:=up125cProcess(x,next,(pos-1)%3,&maxC,&maxH);next++
				if a{fp++}
			}
			if pos%4==0{for k:=0;k<12;k++{x.mem.query(k)}}
		}
	}()
	point:=UP145CPoint{Assignment:assignment,RecallEntriesUsed:x.mem.count,MaxCurrentTableEntries:maxC,MaxHistoryTableEntries:maxH,OneShotFalseAdmissions:fp}
	if panicHit{point.PanicRate=1;return point}
	hotHits,coldHits:=0,0
	for k:=0;k<12;k++{if v,ok:=x.mem.query(k);ok&&v==truth[k]{hotHits++}}
	for k:=12;k<16;k++{if v,ok:=x.mem.query(k);ok&&v==truth[k]{coldHits++}}
	point.HotDurableAccuracy=float64(hotHits)/12.0;point.ColdDurableAccuracy=float64(coldHits)/4.0
	for _,c:=range cands{
		_,ok:=x.mem.query(c.key)
		point.Candidates=append(point.Candidates,UP145CCandidateMetric{Subject:c.subject,Verb:c.verb,Class:c.class,Role:roles[c.key],Positions:append([]int(nil),up145cPositions(roles[c.key])...),AdmissionCount:admit[c.key],FinalRetained:ok})
	}
	return point
}
func RunUP145C()(UP145CResult,error){
	res:=UP145CResult{Schema:UP145CLexicalIntegrationSchema,Experiment:"UP-145C-lexical-consolidation-integration",SourceUP144CSeal:"7bd71e089b99a4e24e958466f3d06603613c0972",ExactRecallCap:16,HistoryEntries:32,AdmissionMemoryBytes:128,CandidateGenerations:3,CandidateCalls:96,HotDurableFacts:12,ColdDurableFacts:4,LexicalClasses:[]string{"STORE","OBSERVE","REPORT"},SemanticPriorityUsed:false,QueryPriorityUsed:false,FutureOracleUsed:false,MemoryIncreased:false}
	for _,a:=range []string{"A","B"}{res.Points=append(res.Points,up145cRun(a))}
	return res,nil
}
