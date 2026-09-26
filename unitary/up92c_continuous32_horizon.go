package unitary

const UP92CContinuous32HorizonSchema = "wingless.up92c-continuous32-horizon.v1"

type UP92CPoint struct {
	ChurnWrites               int     `json:"churn_writes"`
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

type UP92CContinuous32HorizonResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP91CSeal            string       `json:"source_up91c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	HashCount                  int          `json:"hash_count"`
	GenerationInterval         int          `json:"generation_interval"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP92CPoint `json:"points"`
}

func up92cRun(churn int,seedBases []int) UP92CPoint {
	const active=32
	const hotCount=16
	const interval=32
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries:=0,0
	rejectedTotal,falsePositiveAdmissions,resets:=0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,2801+churn*43,ep)
			rng:=newSQ0RNG(seed)
			m:=&up88cFilterMemory{width:256}
			truth:=make([]int,active)
			counter:=0

			process:=func(key,value int)(admitted,rejected bool){
				existing:=m.mem.find(key)>=0
				full:=m.mem.count>=16
				admitted,rejected=m.write(key,value)
				if !existing && full {
					counter++
					if counter==interval {
						m.seen=[4]uint64{}
						counter=0
						resets++
					}
				}
				return
			}

			for k:=0;k<active;k++ {
				v:=rng.intn(32)
				truth[k]=v
				process(k,v)
			}
			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount){continue}
				process(k,truth[k])
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ { if up85cIsHot(h,hotCount){m.query(h)} }
				}
			}
			for j:=0;j<churn;j++ {
				key:=active+j
				admitted,rejected:=process(key,rng.intn(32))
				if rejected{rejectedTotal++}
				if admitted{falsePositiveAdmissions++}
				if (j+1)%4==0 {
					for h:=0;h<active;h++ { if up85cIsHot(h,hotCount){m.query(h)} }
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
		}
	}
	return UP92CPoint{
		ChurnWrites:churn,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		RejectedOneShotWrites:rejectedTotal,
		FalsePositiveAdmissions:falsePositiveAdmissions,
		FilterResets:resets,
		PolicyMetadataBytes:38,
		TotalBoundedMemoryBytes:maxEntries*16+38,
	}
}

func RunUP92C()(UP92CContinuous32HorizonResult,error){
	result:=UP92CContinuous32HorizonResult{
		Schema:UP92CContinuous32HorizonSchema,Experiment:"UP-92C-continuous32-horizon",
		SourceUP91CSeal:"db52aa6b7ede98945d7b22e065b1b2a23aec4b53",
		ActiveKeys:32,HotKeys:16,FilterBits:256,HashCount:2,GenerationInterval:32,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{173000000,174000000}
	for _,churn:=range []int{192,384,768,1536} {
		result.Points=append(result.Points,up92cRun(churn,seedBases))
	}
	return result,nil
}
