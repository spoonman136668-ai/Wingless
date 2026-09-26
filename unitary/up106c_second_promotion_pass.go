package unitary

const UP106CSecondPromotionPassSchema = "wingless.up106c-second-promotion-pass.v1"

type UP106CPoint struct {
	InitialOrder             string  `json:"initial_order"`
	Arm                      string  `json:"arm"`
	ResidentHotAfterPass1    float64 `json:"mean_resident_hot_after_pass1"`
	ResidentHotAfterPromotion float64 `json:"mean_resident_hot_after_promotion"`
	TriggerFiredEpisodes     int     `json:"trigger_fired_episodes"`
	FalsePositiveAdmissions  int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy    float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed        int     `json:"recall_entries_used"`
}

type UP106CSecondPromotionPassResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP105CSeal            string        `json:"source_up105c_seal"`
	HotKeys                     int           `json:"hot_keys"`
	TriggerThreshold            int           `json:"trigger_threshold"`
	FilterBits                  int           `json:"filter_bits"`
	GenerationInterval          int           `json:"generation_interval"`
	ChurnWrites                 int           `json:"churn_writes"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	SecondWriteOrder            string        `json:"second_write_order"`
	PromotionQueryCadence       string        `json:"promotion_query_cadence"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool          `json:"query_labels_used_for_admission"`
	Points                      []UP106CPoint `json:"points"`
}

func up106cRun(initialOrder,arm string,seedBases []int) UP106CPoint {
	const active=32
	const churn=12288
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	resident1Sum,residentFinalSum:=0,0
	episodes,maxEntries:=0,0
	triggerEpisodes,falseAdmissions:=0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,4201+len(initialOrder)*103+len(arm)*31,ep)
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

			promote:=func(){
				written:=0
				for _,k:=range up103cSecondWriteOrder("alternating_ends") {
					process(k,truth[k])
					written++
					if written%4==0 { up102cHotSweep(m) }
				}
			}
			countResident:=func()int{
				n:=0
				for k:=0;k<active;k++ { if up85cIsHot(k,16)&&m.mem.find(k)>=0 { n++ } }
				return n
			}

			promote()
			resident1Sum+=countResident()
			if arm=="two_pass" { promote() }
			residentFinalSum+=countResident()

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

	return UP106CPoint{
		InitialOrder:initialOrder,Arm:arm,
		ResidentHotAfterPass1:float64(resident1Sum)/float64(episodes),
		ResidentHotAfterPromotion:float64(residentFinalSum)/float64(episodes),
		TriggerFiredEpisodes:triggerEpisodes,
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,
	}
}

func RunUP106C()(UP106CSecondPromotionPassResult,error){
	result:=UP106CSecondPromotionPassResult{
		Schema:UP106CSecondPromotionPassSchema,
		Experiment:"UP-106C-second-promotion-pass",
		SourceUP105CSeal:"8534d020165f5f666e924db88f8a54dd2ce67f16",
		HotKeys:16,TriggerThreshold:8,FilterBits:1024,GenerationInterval:32,ChurnWrites:12288,ExactRecallCap:16,
		SecondWriteOrder:"alternating_ends",PromotionQueryCadence:"query_every4",
		FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{201000000,202000000}
	for _,initial:=range []string{"ascending","reverse","evens_then_odds","odds_then_evens","rotate8"} {
		for _,arm:=range []string{"one_pass","two_pass"} {
			result.Points=append(result.Points,up106cRun(initial,arm,seedBases))
		}
	}
	return result,nil
}
