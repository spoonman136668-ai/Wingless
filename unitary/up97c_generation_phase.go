package unitary

const UP97CGenerationPhaseSchema = "wingless.up97c-generation-phase.v1"

type UP97CPoint struct {
	PhaseOffset                  int     `json:"phase_offset"`
	HotQueryAccuracy             float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy          float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy            float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed            int     `json:"recall_entries_used"`
	RejectedOneShotWrites        int     `json:"rejected_one_shot_writes"`
	FalsePositiveAdmissions      int     `json:"false_positive_churn_admissions"`
	FirstFalsePositiveChurnIndex int     `json:"first_false_positive_churn_index"`
	EpisodesWithFalsePositive    int     `json:"episodes_with_false_positive"`
	FilterResets                 int     `json:"filter_resets"`
	PolicyMetadataBytes          int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes      int     `json:"total_bounded_memory_bytes"`
}

type UP97CGenerationPhaseResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP96CSeal            string       `json:"source_up96c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	HashCount                  int          `json:"hash_count"`
	GenerationInterval         int          `json:"generation_interval"`
	ChurnWrites                int          `json:"churn_writes"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP97CPoint `json:"points"`
}

func up97cRun(offset int,seedBases []int) UP97CPoint {
	const active=32
	const hotCount=16
	const width=1024
	const churn=12288
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries,maxBytes:=0,0,0
	rejectedTotal,falsePositiveAdmissions,resets:=0,0,0
	firstFalse:=0
	episodesWithFalse:=0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3301+offset*67,ep)
			rng:=newSQ0RNG(seed)
			m:=&up95cMemory{width:width,seen:make([]uint64,width/64),counter:offset}
			truth:=make([]int,active)

			for k:=0;k<active;k++ {
				v:=rng.intn(32);truth[k]=v;m.write(k,v)
			}
			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount){continue}
				m.write(k,truth[k]);hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ {if up85cIsHot(h,hotCount){m.query(h)}}
				}
			}

			episodeFalse:=false
			for j:=0;j<churn;j++ {
				key:=active+j
				admitted,rejected:=m.write(key,rng.intn(32))
				if rejected{rejectedTotal++}
				if admitted{
					falsePositiveAdmissions++
					if !episodeFalse {
						episodeFalse=true
						episodesWithFalse++
						if firstFalse==0 { firstFalse=j+1 }
					}
				}
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

	return UP97CPoint{
		PhaseOffset:offset,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		RejectedOneShotWrites:rejectedTotal,
		FalsePositiveAdmissions:falsePositiveAdmissions,
		FirstFalsePositiveChurnIndex:firstFalse,
		EpisodesWithFalsePositive:episodesWithFalse,
		FilterResets:resets,
		PolicyMetadataBytes:134,
		TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP97C()(UP97CGenerationPhaseResult,error){
	result:=UP97CGenerationPhaseResult{
		Schema:UP97CGenerationPhaseSchema,
		Experiment:"UP-97C-generation-phase",
		SourceUP96CSeal:"1a9d4ff092ab58db8741d5a526ea6b352867956e",
		ActiveKeys:32,HotKeys:16,FilterBits:1024,HashCount:2,GenerationInterval:32,ChurnWrites:12288,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{183000000,184000000}
	for _,offset:=range []int{0,8,16,24} {
		result.Points=append(result.Points,up97cRun(offset,seedBases))
	}
	return result,nil
}
