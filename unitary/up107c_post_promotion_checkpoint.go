package unitary

const UP107CPostPromotionCheckpointSchema = "wingless.up107c-post-promotion-checkpoint.v1"

type UP107CPoint struct {
	InitialOrder             string  `json:"initial_order"`
	Arm                      string  `json:"arm"`
	ResidentHotBeforeChurn   float64 `json:"mean_resident_hot_before_churn"`
	TriggerFiredEpisodes     int     `json:"trigger_fired_episodes"`
	FalsePositiveAdmissions  int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy    float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed        int     `json:"recall_entries_used"`
}

type UP107CPostPromotionCheckpointResult struct {
	Schema                   string        `json:"schema"`
	Experiment               string        `json:"experiment"`
	SourceUP106CSeal         string        `json:"source_up106c_seal"`
	HotKeys                  int           `json:"hot_keys"`
	TriggerThreshold         int           `json:"trigger_threshold"`
	FilterBits               int           `json:"filter_bits"`
	GenerationInterval       int           `json:"generation_interval"`
	ChurnWrites              int           `json:"churn_writes"`
	ExactRecallCap           int           `json:"exact_recall_cap"`
	ExplicitPhaseLabelUsed   bool          `json:"explicit_phase_label_used"`
	FutureOracleUsed         bool          `json:"future_oracle_used"`
	Points                   []UP107CPoint `json:"points"`
}

func up107cRun(initialOrder,arm string,seedBases []int) UP107CPoint {
	const active=32
	const churn=12288
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	residentSum:=0
	episodes,maxEntries:=0,0
	triggerEpisodes,falseAdmissions:=0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,4301+len(initialOrder)*107+len(arm)*37,ep)
			rng:=newSQ0RNG(seed)
			m:=&up95cMemory{width:width,seen:make([]uint64,width/64)}
			truth:=make([]int,active)
			for k:=0;k<active;k++ { truth[k]=rng.intn(32) }

			repeatAdmissions:=0
			triggerFired:=false
			process:=func(key,value int)(admitted,rejected bool){
				existing:=m.mem.find(key)>=0
				full:=m.mem.count>=16
				seenBefore:=false
				if !existing&&full { seenBefore=m.seenBefore(key) }
				admitted,rejected=m.write(key,value)
				if !existing&&full&&seenBefore&&admitted {
					repeatAdmissions++
					if !triggerFired&&repeatAdmissions==8 {
						m.clearSeen()
						m.counter=8
						triggerFired=true
					}
				}
				return
			}

			for _,k:=range up104cInitialOrder(initialOrder) { process(k,truth[k]) }

			for pass:=0;pass<2;pass++ {
				written:=0
				for _,k:=range up103cSecondWriteOrder("alternating_ends") {
					process(k,truth[k])
					written++
					if written%4==0 { up102cHotSweep(m) }
				}
			}

			if arm=="post_promotion_filter_clear" {
				m.clearSeen();m.counter=0
			}
			if arm=="post_promotion_filter_clear_counter8" {
				m.clearSeen();m.counter=8
			}

			resHot:=0
			for k:=0;k<active;k++ { if up85cIsHot(k,16)&&m.mem.find(k)>=0 { resHot++ } }
			residentSum+=resHot

			for j:=0;j<churn;j++ {
				admitted,_:=process(active+j,rng.intn(32))
				if admitted { falseAdmissions++ }
				if (j+1)%4==0 { up102cHotSweep(m) }
			}

			exact:=true
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,16){continue}
				got,ok:=m.query(k)
				hotTotal++
				if ok&&got==truth[k]{hotHits++}else{exact=false}
			}
			if exact{hotExactHits++}
			if triggerFired{triggerEpisodes++}
			if m.mem.count>maxEntries{maxEntries=m.mem.count}
			episodes++
		}
	}
	return UP107CPoint{
		InitialOrder:initialOrder,Arm:arm,
		ResidentHotBeforeChurn:float64(residentSum)/float64(episodes),
		TriggerFiredEpisodes:triggerEpisodes,
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,
	}
}

func RunUP107C()(UP107CPostPromotionCheckpointResult,error){
	result:=UP107CPostPromotionCheckpointResult{
		Schema:UP107CPostPromotionCheckpointSchema,
		Experiment:"UP-107C-post-promotion-checkpoint",
		SourceUP106CSeal:"84ddc7a82455820c20efa64c9647f332c92c395a",
		HotKeys:16,TriggerThreshold:8,FilterBits:1024,GenerationInterval:32,ChurnWrites:12288,ExactRecallCap:16,
		ExplicitPhaseLabelUsed:true,FutureOracleUsed:false,
	}
	seedBases:=[]int{203000000,204000000}
	for _,initial:=range []string{"odds_then_evens","rotate8","evens_then_odds","ascending"} {
		for _,arm:=range []string{"two_pass_control","post_promotion_filter_clear","post_promotion_filter_clear_counter8"} {
			result.Points=append(result.Points,up107cRun(initial,arm,seedBases))
		}
	}
	return result,nil
}
