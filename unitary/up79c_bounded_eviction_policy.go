package unitary

const UP79CBoundedEvictionSchema = "wingless.up79c-bounded-eviction-policy.v1"

type UP79CPoint struct {
	Policy                  string  `json:"policy"`
	ActiveKeys              int     `json:"active_keys"`
	ChurnWrites             int     `json:"churn_writes"`
	HotQueryAccuracy        float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy     float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy       float64 `json:"cold_query_accuracy"`
	ColdSetExactAccuracy    float64 `json:"cold_set_exact_accuracy"`
	RecallEntriesUsed       int     `json:"recall_entries_used"`
	PolicyMetadataBytes     int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes int     `json:"total_bounded_memory_bytes"`
}

type UP79CBoundedEvictionResult struct {
	Schema            string       `json:"schema"`
	Experiment        string       `json:"experiment"`
	SourceUP78CSeal   string       `json:"source_up78c_seal"`
	ExactRecallCap    int          `json:"exact_recall_cap"`
	EntryPayloadBytes int          `json:"entry_payload_bytes"`
	FutureOracleUsed  bool         `json:"future_oracle_used"`
	Points            []UP79CPoint `json:"points"`
}

type up79cEntry struct {
	key int
	value int
	used bool
	born uint64
	touch uint64
}

type up79cMemory struct {
	policy string
	entries [16]up79cEntry
	count int
	clock uint64
}

func (m *up79cMemory) tick() uint64 {m.clock++;return m.clock}
func (m *up79cMemory) find(key int) int {
	for i:=0;i<16;i++{if m.entries[i].used&&m.entries[i].key==key{return i}}
	return -1
}
func (m *up79cMemory) chooseVictim() int {
	best:=-1
	var stamp uint64
	for i:=0;i<16;i++{
		if !m.entries[i].used{return i}
		s:=m.entries[i].born
		if m.policy=="lru"{s=m.entries[i].touch}
		if best<0||s<stamp{best=i;stamp=s}
	}
	return best
}
func (m *up79cMemory) write(key,value int){
	now:=m.tick()
	if i:=m.find(key);i>=0{
		m.entries[i].value=value
		if m.policy=="lru"{m.entries[i].touch=now}
		return
	}
	i:=m.chooseVictim()
	if !m.entries[i].used{m.count++}
	m.entries[i]=up79cEntry{key:key,value:value,used:true,born:now,touch:now}
}
func (m *up79cMemory) query(key int)(int,bool){
	i:=m.find(key)
	if i<0{return 0,false}
	if m.policy=="lru"{m.entries[i].touch=m.tick()}
	return m.entries[i].value,true
}
func (m *up79cMemory) metadataBytes() int {if m.policy=="lru"{return 136};return 16}
func (m *up79cMemory) totalBytes() int {return m.count*16+m.metadataBytes()}

func up79cRun(policy string,active,churn int,seedBases []int) UP79CPoint {
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal,coldExactHits:=0,0,0
	episodes:=0;maxEntries:=0;maxBytes:=0
	for _,base:=range seedBases{
		for ep:=0;ep<64;ep++{
			seed:=sq0Seed(base,1301+active*11+churn,ep)
			rng:=newSQ0RNG(seed)
			m:=&up79cMemory{policy:policy}
			truth:=make([]int,active)
			for k:=0;k<active;k++{v:=rng.intn(32);truth[k]=v;m.write(k,v)}
			// Rehearse deterministic hot set: even original keys.
			for k:=0;k<active;k+=2{m.query(k)}
			for j:=0;j<churn;j++{m.write(active+j,rng.intn(32))}
			hotExact:=true;coldExact:=true
			for k:=0;k<active;k++{
				got,ok:=m.query(k)
				if k%2==0{
					hotTotal++;if ok&&got==truth[k]{hotHits++}else{hotExact=false}
				}else{
					coldTotal++;if ok&&got==truth[k]{coldHits++}else{coldExact=false}
				}
			}
			if hotExact{hotExactHits++};if coldExact{coldExactHits++};episodes++
			if m.count>maxEntries{maxEntries=m.count};if m.totalBytes()>maxBytes{maxBytes=m.totalBytes()}
		}
	}
	meta:=16;if policy=="lru"{meta=136}
	return UP79CPoint{
		Policy:policy,ActiveKeys:active,ChurnWrites:churn,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),ColdSetExactAccuracy:float64(coldExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:meta,TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP79C()(UP79CBoundedEvictionResult,error){
	result:=UP79CBoundedEvictionResult{
		Schema:UP79CBoundedEvictionSchema,Experiment:"UP-79C-bounded-eviction-policy",
		SourceUP78CSeal:"4644c0ba295515a9d9212b7184ea0b28ff2341c7",
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,
	}
	seedBases:=[]int{145000000,146000000}
	for _,policy:=range []string{"fifo","lru"}{
		for _,active:=range []int{24,32}{
			for _,churn:=range []int{8,16}{
				result.Points=append(result.Points,up79cRun(policy,active,churn,seedBases))
			}
		}
	}
	return result,nil
}
