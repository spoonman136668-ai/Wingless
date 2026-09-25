package unitary

const UP80CSecondChanceSchema = "wingless.up80c-second-chance-eviction.v1"

type UP80CPoint struct {
	Policy string `json:"policy"`
	ActiveKeys int `json:"active_keys"`
	ChurnWrites int `json:"churn_writes"`
	HotQueryAccuracy float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy float64 `json:"cold_query_accuracy"`
	ColdSetExactAccuracy float64 `json:"cold_set_exact_accuracy"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	PolicyMetadataBytes int `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes int `json:"total_bounded_memory_bytes"`
}

type UP80CSecondChanceResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP79CSeal string `json:"source_up79c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	EntryPayloadBytes int `json:"entry_payload_bytes"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Points []UP80CPoint `json:"points"`
}

type up80cClock struct {
	entries [16]up79cEntry
	ref [16]bool
	count int
	hand int
}
func (m *up80cClock) find(key int) int {
	for i:=0;i<16;i++{if m.entries[i].used&&m.entries[i].key==key{return i}}
	return -1
}
func (m *up80cClock) chooseSlot() int {
	if m.count<16{
		for i:=0;i<16;i++{
			idx:=(m.hand+i)%16
			if !m.entries[idx].used{m.hand=(idx+1)%16;return idx}
		}
	}
	for{
		if !m.ref[m.hand]{
			idx:=m.hand;m.hand=(m.hand+1)%16;return idx
		}
		m.ref[m.hand]=false
		m.hand=(m.hand+1)%16
	}
}
func (m *up80cClock) write(key,value int){
	if i:=m.find(key);i>=0{m.entries[i].value=value;m.ref[i]=true;return}
	i:=m.chooseSlot()
	if !m.entries[i].used{m.count++}
	m.entries[i]=up79cEntry{key:key,value:value,used:true}
	m.ref[i]=true
}
func (m *up80cClock) query(key int)(int,bool){
	i:=m.find(key);if i<0{return 0,false}
	m.ref[i]=true
	return m.entries[i].value,true
}
func (m *up80cClock) totalBytes() int {return m.count*16+3}

func up80cRunClock(active,churn int,seedBases []int) UP80CPoint {
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal,coldExactHits:=0,0,0
	episodes:=0;maxEntries:=0;maxBytes:=0
	for _,base:=range seedBases{
		for ep:=0;ep<64;ep++{
			seed:=sq0Seed(base,1501+active*11+churn,ep);rng:=newSQ0RNG(seed)
			m:=&up80cClock{};truth:=make([]int,active)
			for k:=0;k<active;k++{v:=rng.intn(32);truth[k]=v;m.write(k,v)}
			for k:=0;k<active;k+=2{m.query(k)}
			for j:=0;j<churn;j++{m.write(active+j,rng.intn(32))}
			hotExact:=true;coldExact:=true
			for k:=0;k<active;k++{
				got,ok:=m.query(k)
				if k%2==0{hotTotal++;if ok&&got==truth[k]{hotHits++}else{hotExact=false}}else{coldTotal++;if ok&&got==truth[k]{coldHits++}else{coldExact=false}}
			}
			if hotExact{hotExactHits++};if coldExact{coldExactHits++};episodes++
			if m.count>maxEntries{maxEntries=m.count};if m.totalBytes()>maxBytes{maxBytes=m.totalBytes()}
		}
	}
	return UP80CPoint{Policy:"clock_second_chance",ActiveKeys:active,ChurnWrites:churn,HotQueryAccuracy:float64(hotHits)/float64(hotTotal),HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),ColdSetExactAccuracy:float64(coldExactHits)/float64(episodes),RecallEntriesUsed:maxEntries,PolicyMetadataBytes:3,TotalBoundedMemoryBytes:maxBytes}
}

func RunUP80C()(UP80CSecondChanceResult,error){
	result:=UP80CSecondChanceResult{
		Schema:UP80CSecondChanceSchema,Experiment:"UP-80C-second-chance-eviction",
		SourceUP79CSeal:"525487c6011abc396948bee12f032398f201eb89",
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,
	}
	seedBases:=[]int{149000000,150000000}
	for _,policy:=range []string{"fifo","lru"}{
		for _,active:=range []int{24,32}{
			for _,churn:=range []int{8,16}{
				p:=up79cRun(policy,active,churn,seedBases)
				result.Points=append(result.Points,UP80CPoint{
					Policy:p.Policy,ActiveKeys:p.ActiveKeys,ChurnWrites:p.ChurnWrites,
					HotQueryAccuracy:p.HotQueryAccuracy,HotSetExactAccuracy:p.HotSetExactAccuracy,
					ColdQueryAccuracy:p.ColdQueryAccuracy,ColdSetExactAccuracy:p.ColdSetExactAccuracy,
					RecallEntriesUsed:p.RecallEntriesUsed,PolicyMetadataBytes:p.PolicyMetadataBytes,TotalBoundedMemoryBytes:p.TotalBoundedMemoryBytes,
				})
			}
		}
	}
	for _,active:=range []int{24,32}{
		for _,churn:=range []int{8,16}{result.Points=append(result.Points,up80cRunClock(active,churn,seedBases))}
	}
	return result,nil
}
