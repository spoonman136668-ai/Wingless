package unitary

const UP98CPromotionChurnSeparationSchema = "wingless.up98c-promotion-churn-separation.v1"

type UP98CPoint struct {
	Arm                           string  `json:"arm"`
	PhaseOffset                   int     `json:"phase_offset"`
	ClearedAfterPromotion         bool    `json:"cleared_after_promotion"`
	HotQueryAccuracy              float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy           float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy             float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed             int     `json:"recall_entries_used"`
	RejectedOneShotChurnWrites    int     `json:"rejected_one_shot_churn_writes"`
	FalsePositiveAdmissions       int     `json:"false_positive_churn_admissions"`
	FirstFalsePositiveChurnIndex  int     `json:"first_false_positive_churn_index"`
	EpisodesWithFalsePositive     int     `json:"episodes_with_false_positive"`
	AutomaticFilterResets         int     `json:"automatic_filter_resets"`
	PromotionFilterClears         int     `json:"promotion_filter_clears"`
	PolicyMetadataBytes           int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes       int     `json:"total_bounded_memory_bytes"`
}

type UP98CPromotionChurnSeparationResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP97CSeal            string       `json:"source_up97c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	HashCount                  int          `json:"hash_count"`
	GenerationInterval         int          `json:"generation_interval"`
	ChurnWrites                int          `json:"churn_writes"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	DiagnosticPhaseLabelUsed   bool         `json:"diagnostic_phase_label_used"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP98CPoint `json:"points"`
}

func up98cRun(arm string,offset int,clearAfterPromotion bool,seedBases []int) UP98CPoint {
	const active=32
	const hotCount=16
	const churn=12288
	const width=1024

	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries:=0,0
	rejectedChurn,falsePositiveAdmissions,episodesWithFP:=0,0,0
	firstFPGlobal:=0
	autoResets,promotionClears:=0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3401+offset*67+len(arm),ep)
			rng:=newSQ0RNG(seed)
			m:=&up95cMemory{width:width,seen:make([]uint64,width/64)}
			truth:=make([]int,active)

			for k:=0;k<active;k++ {
				v:=rng.intn(32)
				truth[k]=v
				m.write(k,v)
			}

			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount){continue}
				m.write(k,truth[k])
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ {
						if up85cIsHot(h,hotCount){m.query(h)}
					}
				}
			}

			if clearAfterPromotion {
				m.clearSeen()
				m.counter=offset
				promotionClears++
			}

			firstEpisodeFP:=0
			for j:=0;j<churn;j++ {
				key:=active+j
				admitted,rejected:=m.write(key,rng.intn(32))
				if rejected { rejectedChurn++ }
				if admitted {
					falsePositiveAdmissions++
					if firstEpisodeFP==0 { firstEpisodeFP=j+1 }
					if firstFPGlobal==0 { firstFPGlobal=j+1 }
				}
				if (j+1)%4==0 {
					for h:=0;h<active;h++ {
						if up85cIsHot(h,hotCount){m.query(h)}
					}
				}
			}
			if firstEpisodeFP>0 { episodesWithFP++ }

			hotExact:=true
			for k:=0;k<active;k++ {
				got,ok:=m.query(k)
				if up85cIsHot(k,hotCount) {
					hotTotal++
					if ok&&got==truth[k]{hotHits++}else{hotExact=false}
				}else{
					coldTotal++
					if ok&&got==truth[k]{coldHits++}
				}
			}
			if hotExact { hotExactHits++ }
			episodes++
			if m.mem.count>maxEntries { maxEntries=m.mem.count }
			autoResets+=m.resets
		}
	}

	return UP98CPoint{
		Arm:arm,PhaseOffset:offset,ClearedAfterPromotion:clearAfterPromotion,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		RejectedOneShotChurnWrites:rejectedChurn,
		FalsePositiveAdmissions:falsePositiveAdmissions,
		FirstFalsePositiveChurnIndex:firstFPGlobal,
		EpisodesWithFalsePositive:episodesWithFP,
		AutomaticFilterResets:autoResets,
		PromotionFilterClears:promotionClears,
		PolicyMetadataBytes:134,
		TotalBoundedMemoryBytes:maxEntries*16+134,
	}
}

func RunUP98C()(UP98CPromotionChurnSeparationResult,error){
	result:=UP98CPromotionChurnSeparationResult{
		Schema:UP98CPromotionChurnSeparationSchema,
		Experiment:"UP-98C-promotion-churn-separation",
		SourceUP97CSeal:"a52f15992f6e0476f2f36eb88d222822694a9772",
		ActiveKeys:32,HotKeys:16,FilterBits:1024,HashCount:2,GenerationInterval:32,ChurnWrites:12288,
		ExactRecallCap:16,EntryPayloadBytes:16,DiagnosticPhaseLabelUsed:true,
		FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{185000000,186000000}
	result.Points=append(result.Points,
		up98cRun("continuous_control",0,false,seedBases),
		up98cRun("clear_after_promotion_phase0",0,true,seedBases),
		up98cRun("clear_after_promotion_phase8",8,true,seedBases),
		up98cRun("clear_after_promotion_phase16",16,true,seedBases),
		up98cRun("clear_after_promotion_phase24",24,true,seedBases),
	)
	return result,nil
}
