package unitary

const UP95CFilterWidthSchema = "wingless.up95c-filter-width-generational.v1"

type UP95CPoint struct {
	FilterBits               int     `json:"filter_bits"`
	HotQueryAccuracy         float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy      float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy        float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed        int     `json:"recall_entries_used"`
	RejectedOneShotWrites    int     `json:"rejected_one_shot_writes"`
	FalsePositiveAdmissions  int     `json:"false_positive_churn_admissions"`
	FilterResets             int     `json:"filter_resets"`
	PolicyMetadataBytes      int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes  int     `json:"total_bounded_memory_bytes"`
}

type UP95CFilterWidthResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP94CSeal            string       `json:"source_up94c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	HashCount                  int          `json:"hash_count"`
	GenerationInterval         int          `json:"generation_interval"`
	ChurnWrites                int          `json:"churn_writes"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP95CPoint `json:"points"`
}

type up95cMemory struct {
	mem      up81cAging
	seen     []uint64
	width    int
	counter  int
	rejected int
	resets   int
}

func up95cHashBits(key,width int)(int,int){
	x:=uint64(uint32(key))
	bits:=8
	for (1<<bits)<width { bits++ }
	shift:=uint(64-bits)
	h1:=int((x*0x9e3779b97f4a7c15)>>shift)&(width-1)
	y:=x^0xbf58476d1ce4e5b9
	h2:=int((y*0x94d049bb133111eb)>>shift)&(width-1)
	return h1,h2
}
func (m *up95cMemory) bitSet(bit int) bool {
	word:=bit/64
	off:=uint(bit%64)
	return (m.seen[word]&(uint64(1)<<off))!=0
}
func (m *up95cMemory) setBit(bit int){
	word:=bit/64
	off:=uint(bit%64)
	m.seen[word]|=uint64(1)<<off
}
func (m *up95cMemory) seenBefore(key int) bool {
	a,b:=up95cHashBits(key,m.width)
	return m.bitSet(a)&&m.bitSet(b)
}
func (m *up95cMemory) markSeen(key int){
	a,b:=up95cHashBits(key,m.width)
	m.setBit(a);m.setBit(b)
}
func (m *up95cMemory) clearSeen(){for i:=range m.seen{m.seen[i]=0}}

func (m *up95cMemory) write(key,value int)(admitted,rejected bool){
	if m.mem.find(key)>=0 {
		m.mem.write(key,value);return true,false
	}
	if m.mem.count<16 {
		m.mem.write(key,value);return true,false
	}
	if m.seenBefore(key) {
		m.mem.write(key,value);admitted=true
	}else{
		m.markSeen(key);m.rejected++;rejected=true
	}
	m.counter++
	if m.counter==32 {
		m.clearSeen();m.counter=0;m.resets++
	}
	return admitted,rejected
}
func (m *up95cMemory) query(key int)(int,bool){return m.mem.query(key)}
func (m *up95cMemory) metadataBytes() int {return 6+m.width/8}
func (m *up95cMemory) totalBytes() int {return m.mem.count*16+m.metadataBytes()}

func up95cRun(width int,seedBases []int) UP95CPoint {
	const active=32
	const hotCount=16
	const churn=1536
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries,maxBytes:=0,0,0
	rejectedTotal,falsePositiveAdmissions,resets:=0,0,0
	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3101+width*59,ep)
			rng:=newSQ0RNG(seed)
			m:=&up95cMemory{width:width,seen:make([]uint64,width/64)}
			truth:=make([]int,active)
			for k:=0;k<active;k++ {v:=rng.intn(32);truth[k]=v;m.write(k,v)}
			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount){continue}
				m.write(k,truth[k]);hotWritten++
				if hotWritten%4==0 {for h:=0;h<active;h++{if up85cIsHot(h,hotCount){m.query(h)}}}
			}
			for j:=0;j<churn;j++ {
				key:=active+j
				admitted,rejected:=m.write(key,rng.intn(32))
				if rejected{rejectedTotal++}
				if admitted{falsePositiveAdmissions++}
				if (j+1)%4==0 {for h:=0;h<active;h++{if up85cIsHot(h,hotCount){m.query(h)}}}
			}
			hotExact:=true
			for k:=0;k<active;k++ {
				got,ok:=m.query(k)
				if up85cIsHot(k,hotCount){hotTotal++;if ok&&got==truth[k]{hotHits++}else{hotExact=false}}else{coldTotal++;if ok&&got==truth[k]{coldHits++}}
			}
			if hotExact{hotExactHits++}
			episodes++
			if m.mem.count>maxEntries{maxEntries=m.mem.count}
			if m.totalBytes()>maxBytes{maxBytes=m.totalBytes()}
			resets+=m.resets
		}
	}
	return UP95CPoint{
		FilterBits:width,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,RejectedOneShotWrites:rejectedTotal,
		FalsePositiveAdmissions:falsePositiveAdmissions,FilterResets:resets,
		PolicyMetadataBytes:6+width/8,TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP95C()(UP95CFilterWidthResult,error){
	result:=UP95CFilterWidthResult{
		Schema:UP95CFilterWidthSchema,Experiment:"UP-95C-filter-width-generational",
		SourceUP94CSeal:"43727974dc4d24e3244e8eec0e069a1e8c16ae48",
		ActiveKeys:32,HotKeys:16,HashCount:2,GenerationInterval:32,ChurnWrites:1536,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{179000000,180000000}
	for _,width:=range []int{256,512,1024} {result.Points=append(result.Points,up95cRun(width,seedBases))}
	return result,nil
}
