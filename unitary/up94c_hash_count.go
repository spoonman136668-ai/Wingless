package unitary

const UP94CHashCountSchema = "wingless.up94c-hash-count.v1"

type UP94CPoint struct {
	Arm                     string  `json:"arm"`
	HashCount               int     `json:"hash_count"`
	HotQueryAccuracy        float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy     float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy       float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed       int     `json:"recall_entries_used"`
	RejectedOneShotWrites   int     `json:"rejected_one_shot_writes"`
	FalsePositiveAdmissions int     `json:"false_positive_churn_admissions"`
	FilterResets            int     `json:"filter_resets"`
	PolicyMetadataBytes     int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes int     `json:"total_bounded_memory_bytes"`
}

type UP94CHashCountResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP93CSeal            string       `json:"source_up93c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	GenerationInterval         int          `json:"generation_interval"`
	ChurnWrites                int          `json:"churn_writes"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP94CPoint `json:"points"`
}

type up94cMemory struct {
	mem      up81cAging
	seen     [4]uint64
	counter  int
	hashes   int
	legacy   bool
	rejected int
	resets   int
}

func up94cBits(key,hashes int,legacy bool) []int {
	if legacy {
		a,b:=up88cBits(key,256)
		return []int{a,b}
	}
	out:=make([]int,0,hashes)
	x:=uint64(uint32(key))
	for i:=0;i<hashes;i++ {
		z:=sq0Mix64(x ^ ((uint64(i)+1)*0x9e3779b97f4a7c15) ^ 0xd1b54a32d192ed03)
		out=append(out,int(z&255))
	}
	return out
}

func (m *up94cMemory) bitSet(bit int) bool {
	word:=bit/64
	off:=uint(bit%64)
	return (m.seen[word]&(uint64(1)<<off))!=0
}
func (m *up94cMemory) setBit(bit int){
	word:=bit/64
	off:=uint(bit%64)
	m.seen[word]|=uint64(1)<<off
}
func (m *up94cMemory) seenBefore(key int) bool {
	for _,bit:=range up94cBits(key,m.hashes,m.legacy) {
		if !m.bitSet(bit) { return false }
	}
	return true
}
func (m *up94cMemory) markSeen(key int){
	for _,bit:=range up94cBits(key,m.hashes,m.legacy) { m.setBit(bit) }
}

func (m *up94cMemory) write(key,value int)(admitted,rejected bool){
	if m.mem.find(key)>=0 {
		m.mem.write(key,value)
		return true,false
	}
	if m.mem.count<16 {
		m.mem.write(key,value)
		return true,false
	}
	if m.seenBefore(key) {
		m.mem.write(key,value)
		admitted=true
	}else{
		m.markSeen(key)
		m.rejected++
		rejected=true
	}
	m.counter++
	if m.counter==32 {
		m.seen=[4]uint64{}
		m.counter=0
		m.resets++
	}
	return admitted,rejected
}
func (m *up94cMemory) query(key int)(int,bool){return m.mem.query(key)}

func up94cRun(arm string,hashes int,legacy bool,seedBases []int) UP94CPoint {
	const active=32
	const hotCount=16
	const churn=1536
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries:=0,0
	rejectedTotal,falsePositiveAdmissions,resets:=0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3001+hashes*53+len(arm),ep)
			rng:=newSQ0RNG(seed)
			m:=&up94cMemory{hashes:hashes,legacy:legacy}
			truth:=make([]int,active)
			for k:=0;k<active;k++ {
				v:=rng.intn(32);truth[k]=v;m.write(k,v)
			}
			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount){continue}
				m.write(k,truth[k])
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ {if up85cIsHot(h,hotCount){m.query(h)}}
				}
			}
			for j:=0;j<churn;j++ {
				key:=active+j
				admitted,rejected:=m.write(key,rng.intn(32))
				if rejected{rejectedTotal++}
				if admitted{falsePositiveAdmissions++}
				if (j+1)%4==0 {
					for h:=0;h<active;h++ {if up85cIsHot(h,hotCount){m.query(h)}}
				}
			}
			hotExact:=true
			for k:=0;k<active;k++ {
				got,ok:=m.query(k)
				if up85cIsHot(k,hotCount){
					hotTotal++;if ok&&got==truth[k]{hotHits++}else{hotExact=false}
				}else{
					coldTotal++;if ok&&got==truth[k]{coldHits++}
				}
			}
			if hotExact{hotExactHits++}
			episodes++
			if m.mem.count>maxEntries{maxEntries=m.mem.count}
			resets+=m.resets
		}
	}
	return UP94CPoint{
		Arm:arm,HashCount:hashes,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,RejectedOneShotWrites:rejectedTotal,
		FalsePositiveAdmissions:falsePositiveAdmissions,FilterResets:resets,
		PolicyMetadataBytes:38,TotalBoundedMemoryBytes:maxEntries*16+38,
	}
}

func RunUP94C()(UP94CHashCountResult,error){
	result:=UP94CHashCountResult{
		Schema:UP94CHashCountSchema,Experiment:"UP-94C-hash-count",
		SourceUP93CSeal:"37acecd7c314f2c86e40ecff6d687382db864864",
		ActiveKeys:32,HotKeys:16,FilterBits:256,GenerationInterval:32,ChurnWrites:1536,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{177000000,178000000}
	result.Points=append(result.Points,
		up94cRun("legacy2",2,true,seedBases),
		up94cRun("multi2",2,false,seedBases),
		up94cRun("multi4",4,false,seedBases),
		up94cRun("multi6",6,false,seedBases),
	)
	return result,nil
}
