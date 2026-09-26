package unitary

const UP142CAgeRankSchema="wingless.up142c-age-rank-replacement.v1"

type UP142CPoint struct{
	Target int `json:"target"`
	TargetAge int `json:"target_age"`
	TargetIndex int `json:"target_index"`
	TargetPresentBefore bool `json:"target_present_before"`
	TargetPresentAfter bool `json:"target_present_after"`
	IncomingPresentAfter bool `json:"incoming_present_after"`
	EvictedKey int `json:"evicted_key"`
	EvictedIndex int `json:"evicted_index"`
	Age0Evictions int `json:"age0_evictions"`
	Age1Evictions int `json:"age1_evictions"`
	Age2Evictions int `json:"age2_evictions"`
	FinalHistoryEntries int `json:"final_history_entries"`
}
type UP142CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP141CSeal string `json:"source_up141c_seal"`
	HistoryEntries int `json:"history_entries"`
	MaxHistoryAge int `json:"max_history_age"`
	ReplacementSignal string `json:"replacement_signal"`
	TieBreak string `json:"tie_break"`
	DirectStateFixture bool `json:"direct_state_fixture"`
	MemoryIncreased bool `json:"memory_increased"`
	SemanticPriorityUsed bool `json:"semantic_priority_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Points []UP142CPoint `json:"points"`
}
func up142cKeyAt(v uint16)int{if v==0{return -1};return up118cKey(v)}
func up142cRun(target,age,index int)UP142CPoint{
	x:=&up125cAgeEvictMachine{maxAge:2}
	filler:=10000
	for i:=0;i<32;i++{
		key:=filler+i
		if i==index{key=target}
		a:=0;if i==index{a=age}
		x.history[i]=up118cCode(key)|(uint16(a&3)<<14)
	}
	before:=up118cFind(&x.history,target)>=0
	snapshot:=x.history
	incoming:=20000+target+age*100+index
	x.insertHistory(incoming,0)
	after:=up118cFind(&x.history,target)>=0
	incomingPresent:=up118cFind(&x.history,incoming)>=0
	evictedKey,evictedIndex:=-1,-1
	for i:=0;i<32;i++{
		if up142cKeyAt(snapshot[i])!=up142cKeyAt(x.history[i]){
			evictedKey=up142cKeyAt(snapshot[i]);evictedIndex=i;break
		}
	}
	return UP142CPoint{Target:target,TargetAge:age,TargetIndex:index,TargetPresentBefore:before,TargetPresentAfter:after,IncomingPresentAfter:incomingPresent,EvictedKey:evictedKey,EvictedIndex:evictedIndex,Age0Evictions:x.evictions[0],Age1Evictions:x.evictions[1],Age2Evictions:x.evictions[2],FinalHistoryEntries:up125cCount(&x.history)}
}
func RunUP142C()(UP142CResult,error){
	res:=UP142CResult{Schema:UP142CAgeRankSchema,Experiment:"UP-142C-age-rank-replacement",SourceUP141CSeal:"349da133458274f0d25cbde00135928ce4b79494",HistoryEntries:32,MaxHistoryAge:2,ReplacementSignal:"maximum_history_age_only",TieBreak:"lowest_table_index",DirectStateFixture:true,MemoryIncreased:false,SemanticPriorityUsed:false,FutureOracleUsed:false}
	for _,target:=range []int{100,103}{for _,age:=range []int{0,1,2}{for _,idx:=range []int{0,31}{res.Points=append(res.Points,up142cRun(target,age,idx))}}}
	return res,nil
}
