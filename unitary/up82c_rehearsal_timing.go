package unitary

const UP82CRehearsalTimingSchema = "wingless.up82c-rehearsal-timing.v1"

type UP82CPoint struct {
	Policy string `json:"policy"`
	ActiveKeys int `json:"active_keys"`
	ChurnWrites int `json:"churn_writes"`
	RehearsalTiming string `json:"rehearsal_timing"`
	HotQueryAccuracy float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy float64 `json:"cold_query_accuracy"`
	ColdSetExactAccuracy float64 `json:"cold_set_exact_accuracy"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	PolicyMetadataBytes int `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes int `json:"total_bounded_memory_bytes"`
}

type UP82CRehearsalTimingResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP81CSeal string `json:"source_up81c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	EntryPayloadBytes int `json:"entry_payload_bytes"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Points []UP82CPoint `json:"points"`
}

func up82cRun(policy string,active int,timing string,seedBases []int) UP82CPoint {
	const churn=16
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal,coldExactHits:=0,0,0
	episodes,maxEntries,maxBytes:=0,0,0
	meta:=16
	for _,base:=range seedBases{
		for ep:=0;ep<64;ep++{
			seed:=sq0Seed(base,1701+active*13+ep,ep)
			rng:=newSQ0RNG(seed)

			var write func(int,int)
			var query func(int)(int,bool)
			var used func()int
			var bytes func()int

			switch policy{
			case "fifo","lru":
				m:=&up79cMemory{policy:policy}
				write=m.write;query=m.query
				used=func()int{return m.count}
				bytes=m.totalBytes
				if policy=="lru"{meta=136}else{meta=16}
			case "two_bit_aging":
				m:=&up81cAging{}
				write=m.write;query=m.query
				used=func()int{return m.count}
				bytes=m.totalBytes
				meta=5
			}

			truth:=make([]int,active)
			for k:=0;k<active;k++{
				v:=rng.intn(32)
				truth[k]=v
				write(k,v)
			}

			rehearse:=func(){for k:=0;k<active;k+=2{query(k)}}
			rehearsalAfter:=0
			switch timing{
			case "early":rehearsalAfter=0
			case "mid":rehearsalAfter=8
			case "late":rehearsalAfter=16
			}
			if rehearsalAfter==0{rehearse()}
			for j:=0;j<churn;j++{
				if j==rehearsalAfter&&j>0{rehearse()}
				write(active+j,rng.intn(32))
			}
			if rehearsalAfter==churn{rehearse()}

			hotExact,coldExact:=true,true
			for k:=0;k<active;k++{
				got,ok:=query(k)
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
			if used()>maxEntries{maxEntries=used()}
			if bytes()>maxBytes{maxBytes=bytes()}
		}
	}
	return UP82CPoint{
		Policy:policy,ActiveKeys:active,ChurnWrites:churn,RehearsalTiming:timing,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		ColdSetExactAccuracy:float64(coldExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:meta,TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP82C()(UP82CRehearsalTimingResult,error){
	result:=UP82CRehearsalTimingResult{
		Schema:UP82CRehearsalTimingSchema,Experiment:"UP-82C-rehearsal-timing",
		SourceUP81CSeal:"54d90fc5c2a95a37b9afb12898005afa7d7f6614",
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,
	}
	seedBases:=[]int{153000000,154000000}
	for _,policy:=range []string{"fifo","lru","two_bit_aging"}{
		for _,active:=range []int{24,32}{
			for _,timing:=range []string{"early","mid","late"}{
				result.Points=append(result.Points,up82cRun(policy,active,timing,seedBases))
			}
		}
	}
	return result,nil
}
