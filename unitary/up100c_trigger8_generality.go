package unitary

const UP100CTrigger8GeneralitySchema = "wingless.up100c-trigger8-generality.v1"

type UP100CPoint struct {
	Arm                           string  `json:"arm"`
	HotKeys                       int     `json:"hot_keys"`
	InitialOrder                  string  `json:"initial_order"`
	HotQueryAccuracy              float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy           float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy             float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed             int     `json:"recall_entries_used"`
	FalsePositiveAdmissions       int     `json:"false_positive_churn_admissions"`
	FirstFalsePositiveChurnIndex  int     `json:"first_false_positive_churn_index"`
	TriggerFiredEpisodes          int     `json:"trigger_fired_episodes"`
	MeanFilteredWriteIndexTrigger float64 `json:"mean_filtered_write_index_at_trigger"`
	PolicyMetadataBytes           int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes       int     `json:"total_bounded_memory_bytes"`
}

type UP100CTrigger8GeneralityResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP99CSeal             string        `json:"source_up99c_seal"`
	ActiveKeys                  int           `json:"active_keys"`
	FilterBits                  int           `json:"filter_bits"`
	HashCount                   int           `json:"hash_count"`
	GenerationInterval          int           `json:"generation_interval"`
	PostTriggerCounter          int           `json:"post_trigger_counter"`
	TriggerThreshold            int           `json:"trigger_threshold"`
	ChurnWrites                 int           `json:"churn_writes"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	EntryPayloadBytes           int           `json:"entry_payload_bytes"`
	ExplicitPhaseLabelUsed      bool          `json:"explicit_phase_label_used"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool          `json:"query_labels_used_for_admission"`
	Points                      []UP100CPoint `json:"points"`
}

func up100cInitialOrder(name string) []int {
	out:=make([]int,0,32)
	switch name {
	case "ascending":
		for k:=0;k<32;k++ { out=append(out,k) }
	case "reverse":
		for k:=31;k>=0;k-- { out=append(out,k) }
	case "evens_then_odds":
		for k:=0;k<32;k+=2 { out=append(out,k) }
		for k:=1;k<32;k+=2 { out=append(out,k) }
	}
	return out
}

func up100cRun(arm string,hotCount int,orderName string,seedBases []int) UP100CPoint {
	const active=32
	const churn=12288
	const width=1024
	threshold:=0
	if arm=="trigger8" { threshold=8 }

	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries:=0,0
	falsePositiveAdmissions:=0
	firstFPGlobal:=0
	triggerFiredEpisodes:=0
	triggerIndexSum:=0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3601+hotCount*73+len(orderName)*17+threshold,ep)
			rng:=newSQ0RNG(seed)
			m:=&up95cMemory{width:width,seen:make([]uint64,width/64)}
			truth:=make([]int,active)
			for k:=0;k<active;k++ { truth[k]=rng.intn(32) }

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

			for _,k:=range up100cInitialOrder(orderName) { process(k,truth[k]) }

			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount){continue}
				process(k,truth[k])
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ { if up85cIsHot(h,hotCount){m.query(h)} }
				}
			}
			if hotWritten%4!=0 {
				for h:=0;h<active;h++ { if up85cIsHot(h,hotCount){m.query(h)} }
			}

			firstEpisodeFP:=0
			for j:=0;j<churn;j++ {
				admitted,_:=process(active+j,rng.intn(32))
				if admitted {
					falsePositiveAdmissions++
					if firstEpisodeFP==0 { firstEpisodeFP=j+1 }
					if firstFPGlobal==0 { firstFPGlobal=j+1 }
				}
				if (j+1)%4==0 {
					for h:=0;h<active;h++ { if up85cIsHot(h,hotCount){m.query(h)} }
				}
			}
			if triggerFired { triggerFiredEpisodes++;triggerIndexSum+=triggerAt }

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
			_ = firstEpisodeFP
		}
	}

	meanTrigger:=0.0
	if triggerFiredEpisodes>0 { meanTrigger=float64(triggerIndexSum)/float64(triggerFiredEpisodes) }
	meta:=134
	if threshold>0 { meta=135 }
	return UP100CPoint{
		Arm:arm,HotKeys:hotCount,InitialOrder:orderName,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		FalsePositiveAdmissions:falsePositiveAdmissions,
		FirstFalsePositiveChurnIndex:firstFPGlobal,
		TriggerFiredEpisodes:triggerFiredEpisodes,
		MeanFilteredWriteIndexTrigger:meanTrigger,
		PolicyMetadataBytes:meta,
		TotalBoundedMemoryBytes:maxEntries*16+meta,
	}
}

func RunUP100C()(UP100CTrigger8GeneralityResult,error){
	result:=UP100CTrigger8GeneralityResult{
		Schema:UP100CTrigger8GeneralitySchema,
		Experiment:"UP-100C-trigger8-generality",
		SourceUP99CSeal:"816cc124614947ecdb2339f1f874b3c4cc88f355",
		ActiveKeys:32,FilterBits:1024,HashCount:2,GenerationInterval:32,PostTriggerCounter:8,TriggerThreshold:8,
		ChurnWrites:12288,ExactRecallCap:16,EntryPayloadBytes:16,
		ExplicitPhaseLabelUsed:false,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{189000000,190000000}
	for _,hot:=range []int{8,12,16} {
		for _,order:=range []string{"ascending","reverse","evens_then_odds"} {
			result.Points=append(result.Points,
				up100cRun("continuous_control",hot,order,seedBases),
				up100cRun("trigger8",hot,order,seedBases),
			)
		}
	}
	return result,nil
}
