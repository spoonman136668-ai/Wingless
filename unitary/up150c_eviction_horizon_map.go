package unitary

const UP150CEvictionSchema="wingless.up150c-eviction-horizon-map.v1"

type UP150CPoint struct{
	InitialHand int `json:"initial_hand"`
	Admission int `json:"admission"`
	NewKey int `json:"new_key"`
	NewlyEvictedOriginalKey int `json:"newly_evicted_original_key"`
	NewlyEvictedOriginalSlot int `json:"newly_evicted_original_slot"`
	OriginalFactsRemaining int `json:"original_facts_remaining"`
	CurrentHand int `json:"current_hand"`
	RecallEntriesUsed int `json:"recall_entries_used"`
}
type UP150CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP149CSeal string `json:"source_up149c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	InitialHands []int `json:"initial_hands"`
	AdmissionsPerArm int `json:"admissions_per_arm"`
	DurableQueriesAfterInitialization int `json:"durable_queries_after_initialization"`
	MemoryPolicyChanged bool `json:"memory_policy_changed"`
	CapacityIncreased bool `json:"capacity_increased"`
	Points []UP150CPoint `json:"points"`
}
func up150cRun(initialHand int)[]UP150CPoint{
	x:=&up81cAging{};durable:=up145cDurable()
	for _,d:=range durable{x.write(d.key,d.class)}
	x.hand=initialHand
	present:=make([]bool,16);for k:=0;k<16;k++{present[k]=x.find(k)>=0}
	out:=[]UP150CPoint{}
	for j:=0;j<16;j++{
		x.write(1000+j,j%3)
		evictedKey,evictedSlot:=-1,-1
		remaining:=0
		for k:=0;k<16;k++{
			now:=x.find(k)>=0
			if present[k]&&!now{evictedKey=k;evictedSlot=k}
			present[k]=now
			if now{remaining++}
		}
		out=append(out,UP150CPoint{InitialHand:initialHand,Admission:j+1,NewKey:1000+j,NewlyEvictedOriginalKey:evictedKey,NewlyEvictedOriginalSlot:evictedSlot,OriginalFactsRemaining:remaining,CurrentHand:x.hand,RecallEntriesUsed:x.count})
	}
	return out
}
func RunUP150C()(UP150CResult,error){
	hands:=[]int{0,4,8,12}
	res:=UP150CResult{Schema:UP150CEvictionSchema,Experiment:"UP-150C-eviction-horizon-map",SourceUP149CSeal:"d1e1d0794501e7820d24fed83589ed26f47baa8d",ExactRecallCap:16,InitialHands:hands,AdmissionsPerArm:16,DurableQueriesAfterInitialization:0,MemoryPolicyChanged:false,CapacityIncreased:false}
	for _,h:=range hands{res.Points=append(res.Points,up150cRun(h)...)}
	return res,nil
}
