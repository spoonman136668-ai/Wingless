package unitary

const UP87CSeenOnceAdmissionSchema = "wingless.up87c-seenonce-admission.v1"

type UP87CPoint struct {
	Policy                  string  `json:"policy"`
	HotKeys                 int     `json:"hot_keys"`
	HotQueryAccuracy        float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy     float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy       float64 `json:"cold_query_accuracy"`
	ColdSetExactAccuracy    float64 `json:"cold_set_exact_accuracy"`
	RecallEntriesUsed       int     `json:"recall_entries_used"`
	PolicyMetadataBytes     int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes int     `json:"total_bounded_memory_bytes"`
	RejectedOneShotWrites   int     `json:"rejected_one_shot_writes"`
	SecondWriteAdmissions   int     `json:"second_write_admissions"`
}

type UP87CSeenOnceAdmissionResult struct {
	Schema            string       `json:"schema"`
	Experiment        string       `json:"experiment"`
	SourceUP86CSeal   string       `json:"source_up86c_seal"`
	ActiveKeys        int          `json:"active_keys"`
	ChurnWrites       int          `json:"churn_writes"`
	ExactRecallCap    int          `json:"exact_recall_cap"`
	EntryPayloadBytes int          `json:"entry_payload_bytes"`
	FutureOracleUsed  bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool `json:"query_labels_used_for_admission"`
	Points            []UP87CPoint `json:"points"`
}

type up87cSeenOnceMemory struct {
	mem       up81cAging
	seen      uint64
	rejected  int
	admitted2 int
}

func up87cBits(key int)(uint,uint){
	x:=uint64(uint32(key))
	h1:=(x*0x9e3779b97f4a7c15)>>58
	y:=x^0xbf58476d1ce4e5b9
	h2:=(y*0x94d049bb133111eb)>>58
	return uint(h1),uint(h2)
}

func (m *up87cSeenOnceMemory) seenBefore(key int) bool {
	a,b:=up87cBits(key)
	return (m.seen&(uint64(1)<<a))!=0 && (m.seen&(uint64(1)<<b))!=0
}

func (m *up87cSeenOnceMemory) markSeen(key int) {
	a,b:=up87cBits(key)
	m.seen|=uint64(1)<<a
	m.seen|=uint64(1)<<b
}

func (m *up87cSeenOnceMemory) write(key,value int) {
	if m.mem.find(key)>=0 {
		m.mem.write(key,value)
		return
	}
	if m.mem.count<16 {
		m.mem.write(key,value)
		return
	}
	if m.seenBefore(key) {
		m.mem.write(key,value)
		m.admitted2++
		return
	}
	m.markSeen(key)
	m.rejected++
}

func (m *up87cSeenOnceMemory) query(key int)(int,bool){
	return m.mem.query(key)
}

func (m *up87cSeenOnceMemory) totalBytes() int {
	return m.mem.count*16+13
}

func up87cRun(policy string,hotCount int,seedBases []int) UP87CPoint {
	const active=32
	const churn=24
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal,coldExactHits:=0,0,0
	episodes,maxEntries,maxBytes:=0,0,0
	rejected,secondAdmissions:=0,0
	meta:=5

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,2301+hotCount*29,ep)
			rng:=newSQ0RNG(seed)

			var write func(int,int)
			var query func(int)(int,bool)
			var used func()int
			var bytesUsed func()int
			var rejectCount func()int
			var admit2Count func()int

			if policy=="two_bit_aging" {
				m:=&up81cAging{}
				write=m.write; query=m.query
				used=func()int{return m.count}; bytesUsed=m.totalBytes
				rejectCount=func()int{return 0}; admit2Count=func()int{return 0}
				meta=5
			} else {
				m:=&up87cSeenOnceMemory{}
				write=m.write; query=m.query
				used=func()int{return m.mem.count}; bytesUsed=m.totalBytes
				rejectCount=func()int{return m.rejected}; admit2Count=func()int{return m.admitted2}
				meta=13
			}

			truth:=make([]int,active)
			for k:=0;k<active;k++ {
				v:=rng.intn(32)
				truth[k]=v
				write(k,v)
			}

			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount) { continue }
				write(k,truth[k])
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ {
						if up85cIsHot(h,hotCount) { query(h) }
					}
				}
			}
			if hotWritten%4!=0 {
				for h:=0;h<active;h++ {
					if up85cIsHot(h,hotCount) { query(h) }
				}
			}

			for j:=0;j<churn;j++ {
				write(active+j,rng.intn(32))
				if (j+1)%4==0 {
					for h:=0;h<active;h++ {
						if up85cIsHot(h,hotCount) { query(h) }
					}
				}
			}

			hotExact,coldExact:=true,true
			for k:=0;k<active;k++ {
				got,ok:=query(k)
				if up85cIsHot(k,hotCount) {
					hotTotal++
					if ok&&got==truth[k] { hotHits++ } else { hotExact=false }
				} else {
					coldTotal++
					if ok&&got==truth[k] { coldHits++ } else { coldExact=false }
				}
			}
			if hotExact { hotExactHits++ }
			if coldExact { coldExactHits++ }
			episodes++
			if used()>maxEntries { maxEntries=used() }
			if bytesUsed()>maxBytes { maxBytes=bytesUsed() }
			rejected+=rejectCount()
			secondAdmissions+=admit2Count()
		}
	}

	return UP87CPoint{
		Policy:policy,HotKeys:hotCount,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		ColdSetExactAccuracy:float64(coldExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:meta,TotalBoundedMemoryBytes:maxBytes,
		RejectedOneShotWrites:rejected,SecondWriteAdmissions:secondAdmissions,
	}
}

func RunUP87C()(UP87CSeenOnceAdmissionResult,error){
	result:=UP87CSeenOnceAdmissionResult{
		Schema:UP87CSeenOnceAdmissionSchema,
		Experiment:"UP-87C-seenonce-admission",
		SourceUP86CSeal:"58cfc6b5179fdc23a972866250ab380943f7601f",
		ActiveKeys:32,ChurnWrites:24,ExactRecallCap:16,EntryPayloadBytes:16,
		FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{165000000,166000000}
	for _,policy:=range []string{"two_bit_aging","seen_once_gate_plus_aging"} {
		for _,hot:=range []int{13,16} {
			result.Points=append(result.Points,up87cRun(policy,hot,seedBases))
		}
	}
	return result,nil
}
