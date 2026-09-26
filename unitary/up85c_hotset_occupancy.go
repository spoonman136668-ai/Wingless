package unitary

const UP85CHotsetOccupancySchema = "wingless.up85c-hotset-occupancy.v1"

type UP85CPoint struct {
	Policy                  string  `json:"policy"`
	HotKeys                 int     `json:"hot_keys"`
	HotQueryAccuracy        float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy     float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy       float64 `json:"cold_query_accuracy"`
	ColdSetExactAccuracy    float64 `json:"cold_set_exact_accuracy"`
	RecallEntriesUsed       int     `json:"recall_entries_used"`
	PolicyMetadataBytes     int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes int     `json:"total_bounded_memory_bytes"`
}

type UP85CHotsetOccupancyResult struct {
	Schema            string       `json:"schema"`
	Experiment        string       `json:"experiment"`
	SourceUP84CSeal   string       `json:"source_up84c_seal"`
	ActiveKeys        int          `json:"active_keys"`
	ChurnWrites       int          `json:"churn_writes"`
	ExactRecallCap    int          `json:"exact_recall_cap"`
	EntryPayloadBytes int          `json:"entry_payload_bytes"`
	FutureOracleUsed  bool         `json:"future_oracle_used"`
	Points            []UP85CPoint `json:"points"`
}

func up85cIsHot(key, hotCount int) bool {
	if key < 0 || key >= 32 || key%2 != 0 { return false }
	return key/2 < hotCount
}

func up85cRun(policy string, hotCount int, seedBases []int) UP85CPoint {
	const active=32
	const churn=24
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal,coldExactHits:=0,0,0
	episodes,maxEntries,maxBytes:=0,0,0
	meta:=16

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,2201+hotCount*23,ep)
			rng:=newSQ0RNG(seed)

			var write func(int,int)
			var query func(int)(int,bool)
			var used func()int
			var bytesUsed func()int

			switch policy {
			case "fifo","lru":
				m:=&up79cMemory{policy:policy}
				write=m.write; query=m.query
				used=func()int{return m.count}; bytesUsed=m.totalBytes
				if policy=="lru" { meta=136 } else { meta=16 }
			case "two_bit_aging":
				m:=&up81cAging{}
				write=m.write; query=m.query
				used=func()int{return m.count}; bytesUsed=m.totalBytes
				meta=5
			}

			truth:=make([]int,active)
			for k:=0;k<active;k++ {
				v:=rng.intn(32); truth[k]=v; write(k,v)
				if (k+1)%4==0 {
					for h:=0;h<active;h++ {
						if h<=k&&up85cIsHot(h,hotCount) { query(h) }
					}
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
		}
	}

	return UP85CPoint{
		Policy:policy,HotKeys:hotCount,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		ColdSetExactAccuracy:float64(coldExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:meta,TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP85C()(UP85CHotsetOccupancyResult,error){
	result:=UP85CHotsetOccupancyResult{
		Schema:UP85CHotsetOccupancySchema,
		Experiment:"UP-85C-hotset-occupancy",
		SourceUP84CSeal:"51805b22b80b26e927e12dd8b3c66e1282e3cb55",
		ActiveKeys:32,ChurnWrites:24,ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,
	}
	seedBases:=[]int{163000000,164000000}
	for _,policy:=range []string{"fifo","lru","two_bit_aging"} {
		for _,hot:=range []int{4,8,12,16} {
			result.Points=append(result.Points,up85cRun(policy,hot,seedBases))
		}
	}
	return result,nil
}
