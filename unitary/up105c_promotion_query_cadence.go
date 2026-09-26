package unitary

const UP105CPromotionQueryCadenceSchema = "wingless.up105c-promotion-query-cadence.v1"

type UP105CPoint struct {
	InitialOrder                  string  `json:"initial_order"`
	PromotionQueryCadence         string  `json:"promotion_query_cadence"`
	ResidentHotBeforeChurn        float64 `json:"mean_resident_hot_before_churn"`
	TriggerFiredEpisodes          int     `json:"trigger_fired_episodes"`
	MeanFilteredWriteIndexTrigger float64 `json:"mean_filtered_write_index_at_trigger"`
	FalsePositiveAdmissions       int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy         float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy      float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed             int     `json:"recall_entries_used"`
}

type UP105CPromotionQueryCadenceResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP104CSeal            string        `json:"source_up104c_seal"`
	HotKeys                     int           `json:"hot_keys"`
	TriggerThreshold            int           `json:"trigger_threshold"`
	FilterBits                  int           `json:"filter_bits"`
	GenerationInterval          int           `json:"generation_interval"`
	ChurnWrites                 int           `json:"churn_writes"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	SecondWriteOrder            string        `json:"second_write_order"`
	AdaptiveOrderingUsed        bool          `json:"adaptive_ordering_used"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool          `json:"query_labels_used_for_admission"`
	Points                      []UP105CPoint `json:"points"`
}

func up105cRun(initialOrder,cadence string,seedBases []int) UP105CPoint {
	const active=32
	const churn=12288
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	residentHotSum:=0
	episodes,maxEntries:=0,0
	triggerEpisodes,triggerIndexSum,falseAdmissions:=0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,4101+len(initialOrder)*101+len(cadence)*29,ep)
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

			for _,k:=range up104cInitialOrder(initialOrder) { process(k,truth[k]) }

			written:=0
			for _,k:=range up103cSecondWriteOrder("alternating_ends") {
				process(k,truth[k])
				written++
				switch cadence {
				case "query_every4":
					if written%4==0 { up102cHotSweep(m) }
				case "query_every8":
					if written%8==0 { up102cHotSweep(m) }
				case "query_after16":
					if written==16 { up102cHotSweep(m) }
				case "no_promotion_queries":
				}
			}

			resHot:=0
			for k:=0;k<active;k++ {
				if up85cIsHot(k,16)&&m.mem.find(k)>=0 { resHot++ }
			}
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
	if triggerEpisodes>0 { meanTrigger=float64(triggerIndexSum)/float64(triggerEpisodes) }
	return UP105CPoint{
		InitialOrder:initialOrder,PromotionQueryCadence:cadence,
		ResidentHotBeforeChurn:float64(residentHotSum)/float64(episodes),
		TriggerFiredEpisodes:triggerEpisodes,
		MeanFilteredWriteIndexTrigger:meanTrigger,
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,
	}
}

func RunUP105C()(UP105CPromotionQueryCadenceResult,error){
	result:=UP105CPromotionQueryCadenceResult{
		Schema:UP105CPromotionQueryCadenceSchema,
		Experiment:"UP-105C-promotion-query-cadence",
		SourceUP104CSeal:"d34cdb9c75d5b0209969b9dd7244f734ee153ae8",
		HotKeys:16,TriggerThreshold:8,FilterBits:1024,GenerationInterval:32,ChurnWrites:12288,ExactRecallCap:16,
		SecondWriteOrder:"alternating_ends",AdaptiveOrderingUsed:false,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{199000000,200000000}
	for _,initial:=range []string{"ascending","reverse","evens_then_odds","odds_then_evens","rotate8"} {
		for _,cadence:=range []string{"query_every4","query_every8","query_after16","no_promotion_queries"} {
			result.Points=append(result.Points,up105cRun(initial,cadence,seedBases))
		}
	}
	return result,nil
}
