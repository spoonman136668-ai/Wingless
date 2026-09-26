package unitary

const UP88CFilterWidthSchema = "wingless.up88c-filter-width.v1"

type UP88CPoint struct {
	FilterBits               int     `json:"filter_bits"`
	HotQueryAccuracy         float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy      float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy        float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed        int     `json:"recall_entries_used"`
	PolicyMetadataBytes      int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes  int     `json:"total_bounded_memory_bytes"`
	RejectedOneShotWrites    int     `json:"rejected_one_shot_writes"`
	SecondWriteAdmissions    int     `json:"second_write_admissions"`
}

type UP88CFilterWidthResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP87CSeal            string       `json:"source_up87c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	ChurnWrites                int          `json:"churn_writes"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	HashCount                  int          `json:"hash_count"`
	Points                     []UP88CPoint `json:"points"`
}

type up88cFilterMemory struct {
	mem      up81cAging
	seen     [4]uint64
	width    int
	rejected int
}

func up88cShiftForWidth(width int) uint {
	switch width {
	case 64:
		return 58
	case 128:
		return 57
	default:
		return 56
	}
}

func up88cBits(key, width int) (int, int) {
	x := uint64(uint32(key))
	shift := up88cShiftForWidth(width)
	h1 := int((x * 0x9e3779b97f4a7c15) >> shift)
	y := x ^ 0xbf58476d1ce4e5b9
	h2 := int((y * 0x94d049bb133111eb) >> shift)
	return h1, h2
}

func (m *up88cFilterMemory) bitSet(bit int) bool {
	word := bit / 64
	offset := uint(bit % 64)
	return (m.seen[word] & (uint64(1) << offset)) != 0
}

func (m *up88cFilterMemory) setBit(bit int) {
	word := bit / 64
	offset := uint(bit % 64)
	m.seen[word] |= uint64(1) << offset
}

func (m *up88cFilterMemory) seenBefore(key int) bool {
	a,b := up88cBits(key,m.width)
	return m.bitSet(a) && m.bitSet(b)
}

func (m *up88cFilterMemory) markSeen(key int) {
	a,b := up88cBits(key,m.width)
	m.setBit(a); m.setBit(b)
}

func (m *up88cFilterMemory) write(key,value int) (admitted bool, rejected bool) {
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
		return true,false
	}
	m.markSeen(key)
	m.rejected++
	return false,true
}

func (m *up88cFilterMemory) query(key int)(int,bool){ return m.mem.query(key) }
func (m *up88cFilterMemory) totalBytes() int { return m.mem.count*16 + 5 + m.width/8 }

func up88cRun(width int, seedBases []int) UP88CPoint {
	const active=32
	const hotCount=16
	const churn=24
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries,maxBytes:=0,0,0
	rejectedTotal,secondAdmissions:=0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,2401+width,ep)
			rng:=newSQ0RNG(seed)
			m:=&up88cFilterMemory{width:width}
			truth:=make([]int,active)

			for k:=0;k<active;k++ {
				v:=rng.intn(32)
				truth[k]=v
				m.write(k,v)
			}

			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount) { continue }
				before:=m.mem.find(k)>=0
				admitted,_:=m.write(k,truth[k])
				after:=m.mem.find(k)>=0
				if !before && admitted && after { secondAdmissions++ }
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ {
						if up85cIsHot(h,hotCount) { m.query(h) }
					}
				}
			}

			for j:=0;j<churn;j++ {
				_,rej:=m.write(active+j,rng.intn(32))
				if rej { rejectedTotal++ }
				if (j+1)%4==0 {
					for h:=0;h<active;h++ {
						if up85cIsHot(h,hotCount) { m.query(h) }
					}
				}
			}

			hotExact:=true
			for k:=0;k<active;k++ {
				got,ok:=m.query(k)
				if up85cIsHot(k,hotCount) {
					hotTotal++
					if ok&&got==truth[k] { hotHits++ } else { hotExact=false }
				} else {
					coldTotal++
					if ok&&got==truth[k] { coldHits++ }
				}
			}
			if hotExact { hotExactHits++ }
			episodes++
			if m.mem.count>maxEntries { maxEntries=m.mem.count }
			if m.totalBytes()>maxBytes { maxBytes=m.totalBytes() }
		}
	}

	return UP88CPoint{
		FilterBits:width,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		PolicyMetadataBytes:5+width/8,
		TotalBoundedMemoryBytes:maxBytes,
		RejectedOneShotWrites:rejectedTotal,
		SecondWriteAdmissions:secondAdmissions,
	}
}

func RunUP88C()(UP88CFilterWidthResult,error){
	result:=UP88CFilterWidthResult{
		Schema:UP88CFilterWidthSchema,
		Experiment:"UP-88C-filter-width",
		SourceUP87CSeal:"f28e7790d8808ce5e61c2dd3be59d60f4541cb47",
		ActiveKeys:32,HotKeys:16,ChurnWrites:24,ExactRecallCap:16,EntryPayloadBytes:16,
		FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,HashCount:2,
	}
	seedBases:=[]int{165000000,166000000}
	for _,width:=range []int{64,128,256} {
		result.Points=append(result.Points,up88cRun(width,seedBases))
	}
	return result,nil
}
