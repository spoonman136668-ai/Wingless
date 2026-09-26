package unitary

const UP93CGenerationSaltedSchema = "wingless.up93c-generation-salted-hash.v1"

type UP93CPoint struct {
	Arm                     string  `json:"arm"`
	ChurnWrites             int     `json:"churn_writes"`
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

type UP93CGenerationSaltedResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP92CSeal            string       `json:"source_up92c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	HashCount                  int          `json:"hash_count"`
	GenerationInterval         int          `json:"generation_interval"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP93CPoint `json:"points"`
}

type up93cMemory struct {
	mem        up81cAging
	seen       [4]uint64
	counter    int
	generation uint8
	salted     bool
	rejected   int
	resets     int
}

func up93cBits(key int,generation uint8,salted bool)(int,int){
	x:=uint64(uint32(key))
	if salted {
		salt:=sq0Mix64((uint64(generation)+1)*0x9e3779b97f4a7c15)
		x^=salt
	}
	h1:=int((x*0x9e3779b97f4a7c15)>>56)
	y:=x^0xbf58476d1ce4e5b9
	h2:=int((y*0x94d049bb133111eb)>>56)
	return h1,h2
}

func (m *up93cMemory) bitSet(bit int) bool {
	word:=bit/64
	off:=uint(bit%64)
	return (m.seen[word]&(uint64(1)<<off))!=0
}
func (m *up93cMemory) setBit(bit int){
	word:=bit/64
	off:=uint(bit%64)
	m.seen[word]|=uint64(1)<<off
}
func (m *up93cMemory) seenBefore(key int) bool {
	a,b:=up93cBits(key,m.generation,m.salted)
	return m.bitSet(a)&&m.bitSet(b)
}
func (m *up93cMemory) markSeen(key int){
	a,b:=up93cBits(key,m.generation,m.salted)
	m.setBit(a);m.setBit(b)
}

func (m *up93cMemory) write(key,value int)(admitted,rejected bool){
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
		if m.salted { m.generation++ }
	}
	return admitted,rejected
}
func (m *up93cMemory) query(key int)(int,bool){return m.mem.query(key)}
func (m *up93cMemory) metadataBytes() int {
	if m.salted { return 39 }
	return 38
}
func (m *up93cMemory) totalBytes() int { return m.mem.count*16+m.metadataBytes() }

func up93cRun(salted bool,churn int,seedBases []int) UP93CPoint {
	const active=32
	const hotCount=16
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries,maxBytes:=0,0,0
	rejectedTotal,falsePositiveAdmissions,resets:=0,0,0
	arm:="static_hash"
	if salted { arm="generation_salted_hash" }

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,2901+churn*47+boolInt(salted),ep)
			rng:=newSQ0RNG(seed)
			m:=&up93cMemory{salted:salted}
			truth:=make([]int,active)

			for k:=0;k<active;k++ {
				v:=rng.intn(32);truth[k]=v
				m.write(k,v)
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
			if m.totalBytes()>maxBytes{maxBytes=m.totalBytes()}
			resets+=m.resets
		}
	}
	meta:=38
	if salted{meta=39}
	return UP93CPoint{
		Arm:arm,ChurnWrites:churn,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,RejectedOneShotWrites:rejectedTotal,
		FalsePositiveAdmissions:falsePositiveAdmissions,FilterResets:resets,
		PolicyMetadataBytes:meta,TotalBoundedMemoryBytes:maxBytes,
	}
}

func boolInt(v bool) int { if v{return 1};return 0 }

func RunUP93C()(UP93CGenerationSaltedResult,error){
	result:=UP93CGenerationSaltedResult{
		Schema:UP93CGenerationSaltedSchema,Experiment:"UP-93C-generation-salted-hash",
		SourceUP92CSeal:"ce595308c7c181f5a43a1f8e6e7a5b9c8014a810",
		ActiveKeys:32,HotKeys:16,FilterBits:256,HashCount:2,GenerationInterval:32,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{175000000,176000000}
	for _,salted:=range []bool{false,true} {
		for _,churn:=range []int{768,1536,3072} {
			result.Points=append(result.Points,up93cRun(salted,churn,seedBases))
		}
	}
	return result,nil
}
