package unitary

const UP146CDurableReuseSchema="wingless.up146c-durable-reuse-pressure.v1"

type UP146CCandidate struct{
	Subject string `json:"subject"`
	Verb string `json:"verb"`
	Class int `json:"class"`
	Role string `json:"role"`
	AdmissionCount int `json:"admission_count"`
	FinalRetained bool `json:"final_retained"`
}
type UP146CPoint struct{
	Assignment string `json:"assignment"`
	ProtectionArm string `json:"protection_arm"`
	QuinStore4Accuracy float64 `json:"quin_store4_accuracy"`
	RueStore4Accuracy float64 `json:"rue_store4_accuracy"`
	ProtectedGroupAccuracy float64 `json:"protected_group_accuracy"`
	UnprotectedMatchedGroupAccuracy float64 `json:"unprotected_matched_group_accuracy"`
	BackgroundQuinAccuracy float64 `json:"background_quin_accuracy"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	MaxCurrentTableEntries int `json:"max_current_table_entries"`
	MaxHistoryTableEntries int `json:"max_history_table_entries"`
	OneShotFalseAdmissions int `json:"one_shot_false_admissions"`
	PanicRate float64 `json:"panic_rate"`
	Candidates []UP146CCandidate `json:"candidates"`
}
type UP146CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP145CSeal string `json:"source_up145c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	CandidateCalls int `json:"candidate_calls"`
	QueriesPerRefresh int `json:"queries_per_refresh"`
	RefreshEveryCandidateCalls int `json:"refresh_every_candidate_calls"`
	ProtectionArms int `json:"protection_arms"`
	Assignments int `json:"assignments"`
	MemoryPolicyChanged bool `json:"memory_policy_changed"`
	AdaptiveQuerySelection bool `json:"adaptive_query_selection"`
	MemoryIncreased bool `json:"memory_increased"`
	Points []UP146CPoint `json:"points"`
}

func up146cProtectedKeys(arm string)[]int{
	if arm=="protect_rue_store4"{return []int{12,13,14,15}}
	return []int{0,1,2,3}
}
func up146cGroupAccuracy(x *up125cAgeEvictMachine,truth map[int]int,keys []int)float64{
	h:=0
	for _,k:=range keys{if v,ok:=x.mem.query(k);ok&&v==truth[k]{h++}}
	return float64(h)/float64(len(keys))
}
func up146cRun(assignment,arm string)UP146CPoint{
	x:=&up125cAgeEvictMachine{maxAge:2}
	durable:=up145cDurable();cands:=up145cCandidates();roles:=up145cRoles(assignment)
	truth:=map[int]int{}
	for _,d:=range durable{truth[d.key]=d.class;x.mem.write(d.key,d.class)}
	for _,c:=range cands{truth[c.key]=c.class}
	protected:=up146cProtectedKeys(arm)
	admit:=map[int]int{};maxC,maxH,fp:=0,0,0;panicHit:=false
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
			if pos%4==0{for _,k:=range protected{x.mem.query(k)}}
		}
	}()
	p:=UP146CPoint{Assignment:assignment,ProtectionArm:arm,RecallEntriesUsed:x.mem.count,MaxCurrentTableEntries:maxC,MaxHistoryTableEntries:maxH,OneShotFalseAdmissions:fp}
	if panicHit{p.PanicRate=1;return p}
	quin:=[]int{0,1,2,3};rue:=[]int{12,13,14,15};bg:=[]int{4,5,6,7,8,9,10,11}
	p.QuinStore4Accuracy=up146cGroupAccuracy(x,truth,quin)
	p.RueStore4Accuracy=up146cGroupAccuracy(x,truth,rue)
	p.BackgroundQuinAccuracy=up146cGroupAccuracy(x,truth,bg)
	if arm=="protect_rue_store4"{p.ProtectedGroupAccuracy=p.RueStore4Accuracy;p.UnprotectedMatchedGroupAccuracy=p.QuinStore4Accuracy}else{p.ProtectedGroupAccuracy=p.QuinStore4Accuracy;p.UnprotectedMatchedGroupAccuracy=p.RueStore4Accuracy}
	for _,c:=range cands{
		_,ok:=x.mem.query(c.key)
		p.Candidates=append(p.Candidates,UP146CCandidate{Subject:c.subject,Verb:c.verb,Class:c.class,Role:roles[c.key],AdmissionCount:admit[c.key],FinalRetained:ok})
	}
	return p
}
func RunUP146C()(UP146CResult,error){
	res:=UP146CResult{Schema:UP146CDurableReuseSchema,Experiment:"UP-146C-durable-reuse-pressure",SourceUP145CSeal:"365893cb14058ecded93d22719185878ffd08316",ExactRecallCap:16,CandidateCalls:96,QueriesPerRefresh:4,RefreshEveryCandidateCalls:4,ProtectionArms:2,Assignments:2,MemoryPolicyChanged:false,AdaptiveQuerySelection:false,MemoryIncreased:false}
	for _,a:=range []string{"A","B"}{for _,arm:=range []string{"protect_quin_store4","protect_rue_store4"}{res.Points=append(res.Points,up146cRun(a,arm))}}
	return res,nil
}
