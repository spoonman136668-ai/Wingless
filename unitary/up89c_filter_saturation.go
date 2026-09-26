package unitary

const UP89CFilterSaturationSchema = "wingless.up89c-filter-saturation.v1"

type UP89CPoint struct {
	ChurnWrites               int     `json:"churn_writes"`
	HotQueryAccuracy          float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy       float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy         float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed         int     `json:"recall_entries_used"`
	RejectedOneShotWrites     int     `json:"rejected_one_shot_writes"`
	FalsePositiveAdmissions   int     `json:"false_positive_churn_admissions"`
	PolicyMetadataBytes       int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes   int     `json:"total_bounded_memory_bytes"`
}

type UP89CFilterSaturationResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP88CSeal            string       `json:"source_up88c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	HashCount                  int          `json:"hash_count"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP89CPoint `json:"points"`
}

func up89cRun(churn int,seedBases []int) UP89CPoint {
	const active=32
	const hotCount=16
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries,maxBytes:=0,0,0
	rejectedTotal,falsePositiveAdmissions:=0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,2501+churn*31,ep)
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

	return UP89CPoint{
		ChurnWrites:churn,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		RejectedOneShotWrites:rejectedTotal,
		FalsePositiveAdmissions:falsePositiveAdmissions,
		PolicyMetadataBytes:37,
		TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP89C()(UP89CFilterSaturationResult,error){
	result:=UP89CFilterSaturationResult{
		Schema:UP89CFilterSaturationSchema,
		Experiment:"UP-89C-filter-saturation",
		SourceUP88CSeal:"db9f1c6d6478cd5c2f14268163d7dd9e996c01da",
		ActiveKeys:32,HotKeys:16,FilterBits:256,HashCount:2,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{167000000,168000000}
	for _,churn:=range []int{24,48,96,192} {
		result.Points=append(result.Points,up89cRun(churn,seedBases))
	}
	return result,nil
}
