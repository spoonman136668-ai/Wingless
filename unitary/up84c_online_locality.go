package unitary

const UP84COnlineLocalitySchema = "wingless.up84c-online-locality.v1"

type UP84CPoint struct {
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

type UP84COnlineLocalityResult struct {
	Schema            string      `json:"schema"`
	Experiment        string      `json:"experiment"`
	SourceUP83CSeal   string      `json:"source_up83c_seal"`
	ExactRecallCap    int         `json:"exact_recall_cap"`
	EntryPayloadBytes int         `json:"entry_payload_bytes"`
	FutureOracleUsed  bool        `json:"future_oracle_used"`
	Points            []UP84CPoint `json:"points"`
}

func up84cRun(policy string, active, churn int, seedBases []int) UP84CPoint {
	hotHits, hotTotal, hotExactHits := 0, 0, 0
	coldHits, coldTotal, coldExactHits := 0, 0, 0
	episodes, maxEntries, maxBytes := 0, 0, 0
	meta := 16

	for _, base := range seedBases {
		for ep := 0; ep < 64; ep++ {
			seed := sq0Seed(base, 2101+active*19+churn, ep)
			rng := newSQ0RNG(seed)

			var write func(int,int)
			var query func(int)(int,bool)
			var used func()int
			var bytesUsed func()int

			switch policy {
			case "fifo","lru":
				m:=&up79cMemory{policy:policy}
				write=m.write; query=m.query
				used=func()int{return m.count}
				bytesUsed=m.totalBytes
				if policy=="lru" { meta=136 } else { meta=16 }
			case "two_bit_aging":
				m:=&up81cAging{}
				write=m.write; query=m.query
				used=func()int{return m.count}
				bytesUsed=m.totalBytes
				meta=5
			}

			truth:=make([]int,active)
			for k:=0;k<active;k++ {
				v:=rng.intn(32)
				truth[k]=v
				write(k,v)
				if (k+1)%4==0 {
					for h:=0;h<=k;h+=2 { query(h) }
				}
			}

			for j:=0;j<churn;j++ {
				write(active+j,rng.intn(32))
				if (j+1)%4==0 {
					for h:=0;h<active;h+=2 { query(h) }
				}
			}

			hotExact,coldExact:=true,true
			for k:=0;k<active;k++ {
				got,ok:=query(k)
				if k%2==0 {
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

	return UP84CPoint{
		Policy:policy,ActiveKeys:active,ChurnWrites:churn,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		ColdSetExactAccuracy:float64(coldExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:meta,TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP84C()(UP84COnlineLocalityResult,error){
	result:=UP84COnlineLocalityResult{
		Schema:UP84COnlineLocalitySchema,
		Experiment:"UP-84C-online-locality",
		SourceUP83CSeal:"3f6899cdbe85a666377e0724475b52205078a743",
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,
	}
	seedBases:=[]int{161000000,162000000}
	for _,policy:=range []string{"fifo","lru","two_bit_aging"} {
		for _,active:=range []int{24,32} {
			for _,churn:=range []int{16,24} {
				result.Points=append(result.Points,up84cRun(policy,active,churn,seedBases))
			}
		}
	}
	return result,nil
}
