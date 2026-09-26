package unitary

const UP103CPromotionOrderSchema = "wingless.up103c-promotion-order.v1"

type UP103CPoint struct {
	InitialOrder                  string  `json:"initial_order"`
	SecondWriteOrder              string  `json:"second_write_order"`
	ResidentHotBeforeChurn        float64 `json:"mean_resident_hot_before_churn"`
	TriggerFiredEpisodes          int     `json:"trigger_fired_episodes"`
	MeanFilteredWriteIndexTrigger float64 `json:"mean_filtered_write_index_at_trigger"`
	FalsePositiveAdmissions       int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy         float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy      float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed             int     `json:"recall_entries_used"`
}

type UP103CPromotionOrderResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP102CSeal            string        `json:"source_up102c_seal"`
	HotKeys                     int           `json:"hot_keys"`
	TriggerThreshold            int           `json:"trigger_threshold"`
	FilterBits                  int           `json:"filter_bits"`
	GenerationInterval          int           `json:"generation_interval"`
	ChurnWrites                 int           `json:"churn_writes"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	ExtraDiagnosticReadsUsed    bool          `json:"extra_diagnostic_reads_used"`
	AdaptiveOrderingUsed        bool          `json:"adaptive_ordering_used"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool          `json:"query_labels_used_for_admission"`
	Points                      []UP103CPoint `json:"points"`
}

func up103cSecondWriteOrder(name string) []int {
	switch name {
	case "ascending":
		out:=make([]int,0,16)
		for k:=0;k<32;k+=2 { out=append(out,k) }
		return out
	case "reverse":
		out:=make([]int,0,16)
		for k:=30;k>=0;k-=2 { out=append(out,k) }
		return out
	case "alternating_ends":
		out:=make([]int,0,16)
		lo,hi:=0,30
		for lo<hi {
			out=append(out,lo,hi)
			lo+=2
			hi-=2
		}
		if lo==hi { out=append(out,lo) }
		return out
	default:
		return nil
	}
}

func up103cRun(initialOrder,secondOrder string,seedBases []int) UP103CPoint {
	const active=32
	const churn=12288
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	residentHotSum:=0
	episodes,maxEntries:=0,0
	triggerEpisodes,triggerIndexSum,falseAdmissions:=0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3901+len(initialOrder)*89+len(secondOrder)*23,ep)
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

			for _,k:=range up100cInitialOrder(initialOrder) { process(k,truth[k]) }

			written:=0
			for _,k:=range up103cSecondWriteOrder(secondOrder) {
				process(k,truth[k])
				written++
				if written%4==0 { up102cHotSweep(m) }
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
	return UP103CPoint{
		InitialOrder:initialOrder,SecondWriteOrder:secondOrder,
		ResidentHotBeforeChurn:float64(residentHotSum)/float64(episodes),
		TriggerFiredEpisodes:triggerEpisodes,
		MeanFilteredWriteIndexTrigger:meanTrigger,
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,
	}
}

func RunUP103C()(UP103CPromotionOrderResult,error){
	result:=UP103CPromotionOrderResult{
		Schema:UP103CPromotionOrderSchema,
		Experiment:"UP-103C-promotion-order",
		SourceUP102CSeal:"1eaba94eccfcb3cc0d1fd441d907077b00e5587d",
		HotKeys:16,TriggerThreshold:8,FilterBits:1024,GenerationInterval:32,ChurnWrites:12288,ExactRecallCap:16,
		ExtraDiagnosticReadsUsed:false,AdaptiveOrderingUsed:false,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{195000000,196000000}
	for _,initial:=range []string{"ascending","reverse","evens_then_odds"} {
		for _,second:=range []string{"ascending","reverse","alternating_ends"} {
			result.Points=append(result.Points,up103cRun(initial,second,seedBases))
		}
	}
	return result,nil
}
