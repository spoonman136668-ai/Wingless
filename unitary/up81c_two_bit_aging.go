package unitary

const UP81CTwoBitAgingSchema = "wingless.up81c-two-bit-aging.v1"

type UP81CPoint struct {
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

type UP81CTwoBitAgingResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP80CSeal string `json:"source_up80c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	EntryPayloadBytes int `json:"entry_payload_bytes"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Points []UP81CPoint `json:"points"`
}

type up81cAging struct {
	entries [16]up79cEntry
	age [16]uint8
	count int
	hand int
}

func (m *up81cAging) find(key int) int {
	for i:=0;i<16;i++{if m.entries[i].used&&m.entries[i].key==key{return i}}
	return -1
}

func (m *up81cAging) chooseSlot() int {
	if m.count<16{
		for i:=0;i<16;i++{
			idx:=(m.hand+i)%16
			if !m.entries[idx].used{
				m.hand=(idx+1)%16
				return idx
			}
		}
	}
	for{
		if m.age[m.hand]==0{
			idx:=m.hand
			m.hand=(m.hand+1)%16
			return idx
		}
		m.age[m.hand]--
		m.hand=(m.hand+1)%16
	}
}

func (m *up81cAging) write(key,value int){
	if i:=m.find(key);i>=0{
		m.entries[i].value=value
		m.age[i]=3
		return
	}
	i:=m.chooseSlot()
	if !m.entries[i].used{m.count++}
	m.entries[i]=up79cEntry{key:key,value:value,used:true}
	m.age[i]=1
}

func (m *up81cAging) query(key int)(int,bool){
	i:=m.find(key)
	if i<0{return 0,false}
	m.age[i]=3
	return m.entries[i].value,true
}

func (m *up81cAging) totalBytes() int {return m.count*16+5}

func up81cRunAging(active,churn int,seedBases []int) UP81CPoint {
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal,coldExactHits:=0,0,0
	episodes:=0
	maxEntries,maxBytes:=0,0
	for _,base:=range seedBases{
		for ep:=0;ep<64;ep++{
			seed:=sq0Seed(base,1601+active*11+churn,ep)
			rng:=newSQ0RNG(seed)
			m:=&up81cAging{}
			truth:=make([]int,active)
			for k:=0;k<active;k++{
				v:=rng.intn(32)
				truth[k]=v
				m.write(k,v)
			}
			for k:=0;k<active;k+=2{m.query(k)}
			for j:=0;j<churn;j++{m.write(active+j,rng.intn(32))}
			hotExact,coldExact:=true,true
			for k:=0;k<active;k++{
				got,ok:=m.query(k)
				if k%2==0{
					hotTotal++
					if ok&&got==truth[k]{hotHits++}else{hotExact=false}
				}else{
					coldTotal++
					if ok&&got==truth[k]{coldHits++}else{coldExact=false}
				}
			}
			if hotExact{hotExactHits++}
			if coldExact{coldExactHits++}
			episodes++
			if m.count>maxEntries{maxEntries=m.count}
			if m.totalBytes()>maxBytes{maxBytes=m.totalBytes()}
		}
	}
	return UP81CPoint{
		Policy:"two_bit_aging",ActiveKeys:active,ChurnWrites:churn,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		ColdSetExactAccuracy:float64(coldExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:5,TotalBoundedMemoryBytes:maxBytes,
	}
}

func up81cFrom79(p UP79CPoint) UP81CPoint {
	return UP81CPoint{
		Policy:p.Policy,ActiveKeys:p.ActiveKeys,ChurnWrites:p.ChurnWrites,
		HotQueryAccuracy:p.HotQueryAccuracy,HotSetExactAccuracy:p.HotSetExactAccuracy,
		ColdQueryAccuracy:p.ColdQueryAccuracy,ColdSetExactAccuracy:p.ColdSetExactAccuracy,
		RecallEntriesUsed:p.RecallEntriesUsed,PolicyMetadataBytes:p.PolicyMetadataBytes,
		TotalBoundedMemoryBytes:p.TotalBoundedMemoryBytes,
	}
}

func up81cFrom80(p UP80CPoint) UP81CPoint {
	return UP81CPoint{
		Policy:p.Policy,ActiveKeys:p.ActiveKeys,ChurnWrites:p.ChurnWrites,
		HotQueryAccuracy:p.HotQueryAccuracy,HotSetExactAccuracy:p.HotSetExactAccuracy,
		ColdQueryAccuracy:p.ColdQueryAccuracy,ColdSetExactAccuracy:p.ColdSetExactAccuracy,
		RecallEntriesUsed:p.RecallEntriesUsed,PolicyMetadataBytes:p.PolicyMetadataBytes,
		TotalBoundedMemoryBytes:p.TotalBoundedMemoryBytes,
	}
}

func RunUP81C()(UP81CTwoBitAgingResult,error){
	result:=UP81CTwoBitAgingResult{
		Schema:UP81CTwoBitAgingSchema,Experiment:"UP-81C-two-bit-aging",
		SourceUP80CSeal:"b100e17cbac6354c44b42e9dad25101dd97c40be",
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,
	}
	seedBases:=[]int{151000000,152000000}
	for _,policy:=range []string{"fifo","lru"}{
		for _,active:=range []int{24,32}{
			for _,churn:=range []int{8,16}{
				result.Points=append(result.Points,up81cFrom79(up79cRun(policy,active,churn,seedBases)))
			}
		}
	}
	for _,active:=range []int{24,32}{
		for _,churn:=range []int{8,16}{
			result.Points=append(result.Points,up81cFrom80(up80cRunClock(active,churn,seedBases)))
		}
	}
	for _,active:=range []int{24,32}{
		for _,churn:=range []int{8,16}{
			result.Points=append(result.Points,up81cRunAging(active,churn,seedBases))
		}
	}
	return result,nil
}
