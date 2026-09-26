package unitary

const UP102CQueryAgingTimingSchema = "wingless.up102c-query-aging-timing.v1"

type UP102CPoint struct {
	InitialOrder                  string  `json:"initial_order"`
	ReadTiming                    string  `json:"read_timing"`
	ResidentHotBeforeChurn        float64 `json:"mean_resident_hot_before_churn"`
	TriggerFiredEpisodes          int     `json:"trigger_fired_episodes"`
	MeanFilteredWriteIndexTrigger float64 `json:"mean_filtered_write_index_at_trigger"`
	FalsePositiveAdmissions       int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy         float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy      float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed             int     `json:"recall_entries_used"`
}

type UP102CQueryAgingTimingResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP101CSeal            string        `json:"source_up101c_seal"`
	HotKeys                     int           `json:"hot_keys"`
	TriggerThreshold            int           `json:"trigger_threshold"`
	FilterBits                  int           `json:"filter_bits"`
	GenerationInterval          int           `json:"generation_interval"`
	ChurnWrites                 int           `json:"churn_writes"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	AdaptiveReadsUsed           bool          `json:"adaptive_reads_used"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool          `json:"query_labels_used_for_admission"`
	Points                      []UP102CPoint `json:"points"`
}

func up102cHotSweep(m *up95cMemory) {
	for h:=0;h<32;h++ { if up85cIsHot(h,16){ m.query(h) } }
}

func up102cRun(orderName,readTiming string,seedBases []int) UP102CPoint {
	const active=32
	const churn=12288
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	residentHotSum:=0
	episodes,maxEntries:=0,0
	triggerEpisodes,triggerIndexSum,falseAdmissions:=0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3801+len(orderName)*83+len(readTiming)*19,ep)
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
				if !existing&&full {
					filteredWriteIndex++
					seenBefore=m.seenBefore(key)
				}
				admitted,rejected=m.write(key,value)
				if !existing&&full&&seenBefore&&admitted {
					repeatAdmissions++
					if !triggerFired&&repeatAdmissions==8 {
						m.clearSeen()
						m.counter=8
						triggerFired=true
						triggerAt=filteredWriteIndex
					}
				}
				return
			}

			for _,k:=range up100cInitialOrder(orderName){ process(k,truth[k]) }
			if readTiming=="after_initial"||readTiming=="both"{ up102cHotSweep(m) }

			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,16){continue}
				process(k,truth[k])
				hotWritten++
				if hotWritten%4==0 { up102cHotSweep(m) }
			}
			if readTiming=="after_promotion"||readTiming=="both"{ up102cHotSweep(m) }

			resHot:=0
			for k:=0;k<active;k++ { if up85cIsHot(k,16)&&m.mem.find(k)>=0 { resHot++ } }
			residentHotSum+=resHot

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
			if triggerFired{triggerEpisodes++;triggerIndexSum+=triggerAt}
			if m.mem.count>maxEntries{maxEntries=m.mem.count}
			episodes++
		}
	}
	meanTrigger:=0.0
	if triggerEpisodes>0{meanTrigger=float64(triggerIndexSum)/float64(triggerEpisodes)}
	return UP102CPoint{
		InitialOrder:orderName,ReadTiming:readTiming,
		ResidentHotBeforeChurn:float64(residentHotSum)/float64(episodes),
		TriggerFiredEpisodes:triggerEpisodes,MeanFilteredWriteIndexTrigger:meanTrigger,
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,
	}
}

func RunUP102C()(UP102CQueryAgingTimingResult,error){
	result:=UP102CQueryAgingTimingResult{
		Schema:UP102CQueryAgingTimingSchema,Experiment:"UP-102C-query-aging-timing",
		SourceUP101CSeal:"0deaa954944277a63e1b96c01d2769bba82b4b73",
		HotKeys:16,TriggerThreshold:8,FilterBits:1024,GenerationInterval:32,ChurnWrites:12288,ExactRecallCap:16,
		AdaptiveReadsUsed:false,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{193000000,194000000}
	for _,order:=range []string{"ascending","reverse","evens_then_odds"} {
		for _,timing:=range []string{"no_extra_reads","after_initial","after_promotion","both"} {
			result.Points=append(result.Points,up102cRun(order,timing,seedBases))
		}
	}
	return result,nil
}
