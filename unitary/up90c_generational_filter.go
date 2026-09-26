package unitary

const UP90CGenerationalFilterSchema = "wingless.up90c-generational-filter.v1"

type UP90CPoint struct {
	ResetInterval             int     `json:"reset_interval"`
	HotQueryAccuracy          float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy       float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy         float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed         int     `json:"recall_entries_used"`
	RejectedOneShotWrites     int     `json:"rejected_one_shot_writes"`
	FalsePositiveAdmissions   int     `json:"false_positive_churn_admissions"`
	FilterResets              int     `json:"filter_resets"`
	PolicyMetadataBytes       int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes   int     `json:"total_bounded_memory_bytes"`
}

type UP90CGenerationalFilterResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP89CSeal            string       `json:"source_up89c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	HashCount                  int          `json:"hash_count"`
	ChurnWrites                int          `json:"churn_writes"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP90CPoint `json:"points"`
}

func up90cRun(resetInterval int,seedBases []int) UP90CPoint {
	const active=32
	const hotCount=16
	const churn=192
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries,maxBytes:=0,0,0
	rejectedTotal,falsePositiveAdmissions,resets:=0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,2601+resetInterval*37,ep)
			rng:=newSQ0RNG(seed)
			m:=&up88cFilterMemory{width:256}
			truth:=make([]int,active)

			for k:=0;k<active;k++ {
				v:=rng.intn(32)
				truth[k]=v
				m.write(k,v)
			}

			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount) { continue }
				m.write(k,truth[k])
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ {
						if up85cIsHot(h,hotCount) { m.query(h) }
					}
				}
			}

			for j:=0;j<churn;j++ {
				key:=active+j
				admitted,rejected:=m.write(key,rng.intn(32))
				if rejected { rejectedTotal++ }
				if admitted { falsePositiveAdmissions++ }
				if (j+1)%4==0 {
					for h:=0;h<active;h++ {
						if up85cIsHot(h,hotCount) { m.query(h) }
					}
				}
				if resetInterval>0 && (j+1)%resetInterval==0 {
					m.seen=[4]uint64{}
					resets++
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
			bytes:=m.mem.count*16+37
			if resetInterval>0 { bytes++ }
			if bytes>maxBytes { maxBytes=bytes }
		}
	}

	meta:=37
	if resetInterval>0 { meta=38 }
	return UP90CPoint{
		ResetInterval:resetInterval,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		RejectedOneShotWrites:rejectedTotal,
		FalsePositiveAdmissions:falsePositiveAdmissions,
		FilterResets:resets,
		PolicyMetadataBytes:meta,
		TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP90C()(UP90CGenerationalFilterResult,error){
	result:=UP90CGenerationalFilterResult{
		Schema:UP90CGenerationalFilterSchema,
		Experiment:"UP-90C-generational-filter",
		SourceUP89CSeal:"4834165acfdd79552d12517014d59cf30a2e75e0",
		ActiveKeys:32,HotKeys:16,FilterBits:256,HashCount:2,ChurnWrites:192,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{169000000,170000000}
	for _,interval:=range []int{0,32,64,96} {
		result.Points=append(result.Points,up90cRun(interval,seedBases))
	}
	return result,nil
}
