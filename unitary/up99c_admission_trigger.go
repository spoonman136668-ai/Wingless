package unitary

const UP99CAdmissionTriggerSchema = "wingless.up99c-admission-trigger.v1"

type UP99CPoint struct {
	Arm                           string  `json:"arm"`
	TriggerThreshold              int     `json:"trigger_threshold"`
	HotQueryAccuracy              float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy           float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy             float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed             int     `json:"recall_entries_used"`
	RejectedChurnWrites           int     `json:"rejected_churn_writes"`
	FalsePositiveAdmissions       int     `json:"false_positive_churn_admissions"`
	FirstFalsePositiveChurnIndex  int     `json:"first_false_positive_churn_index"`
	EpisodesWithFalsePositive     int     `json:"episodes_with_false_positive"`
	TriggerFiredEpisodes          int     `json:"trigger_fired_episodes"`
	MeanFilteredWriteIndexTrigger float64 `json:"mean_filtered_write_index_at_trigger"`
	AutomaticFilterResets         int     `json:"automatic_filter_resets"`
	PolicyMetadataBytes           int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes       int     `json:"total_bounded_memory_bytes"`
}

type UP99CAdmissionTriggerResult struct {
	Schema                      string       `json:"schema"`
	Experiment                  string       `json:"experiment"`
	SourceUP98CSeal             string       `json:"source_up98c_seal"`
	ActiveKeys                  int          `json:"active_keys"`
	HotKeys                     int          `json:"hot_keys"`
	FilterBits                  int          `json:"filter_bits"`
	HashCount                   int          `json:"hash_count"`
	GenerationInterval          int          `json:"generation_interval"`
	PostTriggerCounter          int          `json:"post_trigger_counter"`
	ChurnWrites                 int          `json:"churn_writes"`
	ExactRecallCap              int          `json:"exact_recall_cap"`
	EntryPayloadBytes           int          `json:"entry_payload_bytes"`
	ExplicitPhaseLabelUsed      bool         `json:"explicit_phase_label_used"`
	FutureOracleUsed            bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool         `json:"query_labels_used_for_admission"`
	Points                      []UP99CPoint `json:"points"`
}

func up99cRun(arm string,threshold int,seedBases []int) UP99CPoint {
	const active=32
	const hotCount=16
	const churn=12288
	const width=1024

	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries:=0,0
	rejectedChurn,falsePositiveAdmissions,episodesWithFP:=0,0,0
	firstFPGlobal:=0
	triggerFiredEpisodes:=0
	triggerIndexSum:=0
	autoResets:=0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3501+threshold*71+len(arm),ep)
			rng:=newSQ0RNG(seed)
			m:=&up95cMemory{width:width,seen:make([]uint64,width/64)}
			truth:=make([]int,active)
			repeatAdmissions:=0
			triggerFired:=false
			filteredWriteIndex:=0
			triggerAt:=0

			process:=func(key,value int)(admitted,rejected bool){
				existing:=m.mem.find(key)>=0
				full:=m.mem.count>=16
				seenBefore:=false
				if !existing && full {
					filteredWriteIndex++
					seenBefore=m.seenBefore(key)
				}
				admitted,rejected=m.write(key,value)
				if threshold>0 && !triggerFired && !existing && full && seenBefore && admitted {
					repeatAdmissions++
					if repeatAdmissions==threshold {
						m.clearSeen()
						m.counter=8
						triggerFired=true
						triggerAt=filteredWriteIndex
					}
				}
				return
			}

			for k:=0;k<active;k++ {
				v:=rng.intn(32);truth[k]=v;process(k,v)
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

			firstEpisodeFP:=0
			for j:=0;j<churn;j++ {
				admitted,rejected:=process(active+j,rng.intn(32))
				if rejected { rejectedChurn++ }
				if admitted {
					falsePositiveAdmissions++
					if firstEpisodeFP==0 { firstEpisodeFP=j+1 }
					if firstFPGlobal==0 { firstFPGlobal=j+1 }
				}
				if (j+1)%4==0 {
					for h:=0;h<active;h++ { if up85cIsHot(h,hotCount){m.query(h)} }
				}
			}
			if firstEpisodeFP>0 { episodesWithFP++ }
			if triggerFired { triggerFiredEpisodes++;triggerIndexSum+=triggerAt }

			hotExact:=true
			for k:=0;k<active;k++ {
				got,ok:=m.query(k)
				if up85cIsHot(k,hotCount) {
					hotTotal++;if ok&&got==truth[k]{hotHits++}else{hotExact=false}
				}else{
					coldTotal++;if ok&&got==truth[k]{coldHits++}
				}
			}
			if hotExact { hotExactHits++ }
			episodes++
			if m.mem.count>maxEntries { maxEntries=m.mem.count }
			autoResets+=m.resets
		}
	}

	meanTrigger:=0.0
	if triggerFiredEpisodes>0 { meanTrigger=float64(triggerIndexSum)/float64(triggerFiredEpisodes) }
	meta:=134
	if threshold>0 { meta=135 }
	return UP99CPoint{
		Arm:arm,TriggerThreshold:threshold,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		RejectedChurnWrites:rejectedChurn,
		FalsePositiveAdmissions:falsePositiveAdmissions,
		FirstFalsePositiveChurnIndex:firstFPGlobal,
		EpisodesWithFalsePositive:episodesWithFP,
		TriggerFiredEpisodes:triggerFiredEpisodes,
		MeanFilteredWriteIndexTrigger:meanTrigger,
		AutomaticFilterResets:autoResets,
		PolicyMetadataBytes:meta,
		TotalBoundedMemoryBytes:maxEntries*16+meta,
	}
}

func RunUP99C()(UP99CAdmissionTriggerResult,error){
	result:=UP99CAdmissionTriggerResult{
		Schema:UP99CAdmissionTriggerSchema,
		Experiment:"UP-99C-admission-trigger",
		SourceUP98CSeal:"89f6fd97bff95723c248a1a34a07ca43e163d864",
		ActiveKeys:32,HotKeys:16,FilterBits:1024,HashCount:2,GenerationInterval:32,PostTriggerCounter:8,
		ChurnWrites:12288,ExactRecallCap:16,EntryPayloadBytes:16,
		ExplicitPhaseLabelUsed:false,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{187000000,188000000}
	result.Points=append(result.Points,
		up99cRun("continuous_control",0,seedBases),
		up99cRun("trigger4",4,seedBases),
		up99cRun("trigger8",8,seedBases),
		up99cRun("trigger12",12,seedBases),
	)
	return result,nil
}
