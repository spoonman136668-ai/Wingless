package unitary

const UP91CContinuousGenerationSchema = "wingless.up91c-continuous-generation.v1"

type UP91CPoint struct {
	Arm                     string  `json:"arm"`
	Interval                int     `json:"interval"`
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

type UP91CContinuousGenerationResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP90CSeal            string       `json:"source_up90c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	HashCount                  int          `json:"hash_count"`
	ChurnWrites                int          `json:"churn_writes"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP91CPoint `json:"points"`
}

func up91cRun(arm string,interval int,seedBases []int) UP91CPoint {
	const active=32
	const hotCount=16
	const churn=192
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries:=0,0
	rejectedTotal,falsePositiveAdmissions,resets:=0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,2701+interval*41+len(arm),ep)
			rng:=newSQ0RNG(seed)
			m:=&up88cFilterMemory{width:256}
			truth:=make([]int,active)
			counter:=0

			process:=func(key,value int,countContinuous bool)(admitted,rejected bool){
				existing:=m.mem.find(key)>=0
				full:=m.mem.count>=16
				admitted,rejected=m.write(key,value)
				if arm!="phase_local_32" && countContinuous && !existing && full {
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
				v:=rng.intn(32);truth[k]=v
				process(k,v,true)
			}
			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount){continue}
				process(k,truth[k],true)
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ { if up85cIsHot(h,hotCount){m.query(h)} }
				}
			}
			for j:=0;j<churn;j++ {
				key:=active+j
				admitted,rejected:=process(key,rng.intn(32),true)
				if rejected { rejectedTotal++ }
				if admitted { falsePositiveAdmissions++ }
				if (j+1)%4==0 {
					for h:=0;h<active;h++ { if up85cIsHot(h,hotCount){m.query(h)} }
				}
				if arm=="phase_local_32" && (j+1)%32==0 {
					m.seen=[4]uint64{}
					resets++
				}
			}
			hotExact:=true
			for k:=0;k<active;k++ {
				got,ok:=m.query(k)
				if up85cIsHot(k,hotCount) {
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
	return UP91CPoint{
		Arm:arm,Interval:interval,
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

func RunUP91C()(UP91CContinuousGenerationResult,error){
	result:=UP91CContinuousGenerationResult{
		Schema:UP91CContinuousGenerationSchema,Experiment:"UP-91C-continuous-generation",
		SourceUP90CSeal:"68d924c4ddc53e4bcc3145cc5fa1eb5ef11d1539",
		ActiveKeys:32,HotKeys:16,FilterBits:256,HashCount:2,ChurnWrites:192,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{171000000,172000000}
	result.Points=append(result.Points,
		up91cRun("phase_local_32",32,seedBases),
		up91cRun("continuous_16",16,seedBases),
		up91cRun("continuous_32",32,seedBases),
		up91cRun("continuous_64",64,seedBases),
	)
	return result,nil
}
