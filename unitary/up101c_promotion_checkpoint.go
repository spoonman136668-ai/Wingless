package unitary

const UP101CPromotionCheckpointSchema = "wingless.up101c-promotion-checkpoint.v1"

type UP101CCheckpoint struct {
	Name                  string  `json:"name"`
	MeanResidentHotKeys   float64 `json:"mean_resident_hot_keys"`
	HotQueryAccuracy      float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy   float64 `json:"hot_set_exact_accuracy"`
	MeanResidentColdKeys  float64 `json:"mean_resident_cold_keys"`
}

type UP101CPoint struct {
	Arm                           string               `json:"arm"`
	InitialOrder                  string               `json:"initial_order"`
	Checkpoints                   []UP101CCheckpoint   `json:"checkpoints"`
	TotalRepeatAdmissions         int                  `json:"total_repeat_admissions"`
	MeanRepeatAdmissions          float64              `json:"mean_repeat_admissions"`
	TriggerFiredEpisodes          int                  `json:"trigger_fired_episodes"`
	MeanFilteredWriteIndexTrigger float64              `json:"mean_filtered_write_index_at_trigger"`
	FalsePositiveAdmissions       int                  `json:"false_positive_churn_admissions"`
	RecallEntriesUsed             int                  `json:"recall_entries_used"`
}

type UP101CPromotionCheckpointResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP100CSeal            string        `json:"source_up100c_seal"`
	ActiveKeys                  int           `json:"active_keys"`
	HotKeys                     int           `json:"hot_keys"`
	FilterBits                  int           `json:"filter_bits"`
	HashCount                   int           `json:"hash_count"`
	GenerationInterval          int           `json:"generation_interval"`
	PostTriggerCounter          int           `json:"post_trigger_counter"`
	TriggerThreshold            int           `json:"trigger_threshold"`
	ChurnWrites                 int           `json:"churn_writes"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	ExplicitPhaseLabelUsed      bool          `json:"explicit_phase_label_used"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool          `json:"query_labels_used_for_admission"`
	Points                      []UP101CPoint `json:"points"`
}

type up101cAccum struct {
	residentHot int
	hotHits int
	hotTotal int
	hotExact int
	residentCold int
	episodes int
}

func up101cCheckpoint(m *up95cMemory,truth []int,a *up101cAccum) {
	resHot,resCold:=0,0
	for k:=0;k<32;k++ {
		if m.mem.find(k)>=0 {
			if up85cIsHot(k,16){resHot++}else{resCold++}
		}
	}
	exact:=true
	for k:=0;k<32;k++ {
		if !up85cIsHot(k,16){continue}
		got,ok:=m.query(k)
		a.hotTotal++
		if ok&&got==truth[k]{a.hotHits++}else{exact=false}
	}
	a.residentHot+=resHot
	a.residentCold+=resCold
	if exact{a.hotExact++}
	a.episodes++
}

func up101cFinalCheckpoint(name string,a up101cAccum)UP101CCheckpoint{
	return UP101CCheckpoint{
		Name:name,
		MeanResidentHotKeys:float64(a.residentHot)/float64(a.episodes),
		HotQueryAccuracy:float64(a.hotHits)/float64(a.hotTotal),
		HotSetExactAccuracy:float64(a.hotExact)/float64(a.episodes),
		MeanResidentColdKeys:float64(a.residentCold)/float64(a.episodes),
	}
}

func up101cRun(arm,orderName string,seedBases []int) UP101CPoint {
	const active=32
	const churn=12288
	const width=1024
	threshold:=0
	if arm=="trigger8"{threshold=8}
	var initialAcc,promotionAcc,finalAcc up101cAccum
	totalRepeat,triggerEpisodes,triggerIndexSum,falseAdmissions,maxEntries:=0,0,0,0,0
	episodes:=0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3701+len(orderName)*79+threshold,ep)
			rng:=newSQ0RNG(seed)
			m:=&up95cMemory{width:width,seen:make([]uint64,width/64)}
			truth:=make([]int,active)
			for k:=0;k<active;k++{truth[k]=rng.intn(32)}
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
					if threshold>0&&!triggerFired&&repeatAdmissions==threshold {
						m.clearSeen()
						m.counter=8
						triggerFired=true
						triggerAt=filteredWriteIndex
					}
				}
				return
			}

			for _,k:=range up100cInitialOrder(orderName){process(k,truth[k])}
			up101cCheckpoint(m,truth,&initialAcc)

			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,16){continue}
				process(k,truth[k])
				hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++{if up85cIsHot(h,16){m.query(h)}}
				}
			}
			up101cCheckpoint(m,truth,&promotionAcc)

			for j:=0;j<churn;j++ {
				admitted,_:=process(active+j,rng.intn(32))
				if admitted{falseAdmissions++}
				if (j+1)%4==0 {
					for h:=0;h<active;h++{if up85cIsHot(h,16){m.query(h)}}
				}
			}
			up101cCheckpoint(m,truth,&finalAcc)

			totalRepeat+=repeatAdmissions
			if triggerFired{triggerEpisodes++;triggerIndexSum+=triggerAt}
			if m.mem.count>maxEntries{maxEntries=m.mem.count}
			episodes++
		}
	}
	meanTrigger:=0.0
	if triggerEpisodes>0{meanTrigger=float64(triggerIndexSum)/float64(triggerEpisodes)}
	return UP101CPoint{
		Arm:arm,InitialOrder:orderName,
		Checkpoints:[]UP101CCheckpoint{
			up101cFinalCheckpoint("after_initial",initialAcc),
			up101cFinalCheckpoint("after_promotion",promotionAcc),
			up101cFinalCheckpoint("after_churn",finalAcc),
		},
		TotalRepeatAdmissions:totalRepeat,
		MeanRepeatAdmissions:float64(totalRepeat)/float64(episodes),
		TriggerFiredEpisodes:triggerEpisodes,
		MeanFilteredWriteIndexTrigger:meanTrigger,
		FalsePositiveAdmissions:falseAdmissions,
		RecallEntriesUsed:maxEntries,
	}
}

func RunUP101C()(UP101CPromotionCheckpointResult,error){
	result:=UP101CPromotionCheckpointResult{
		Schema:UP101CPromotionCheckpointSchema,Experiment:"UP-101C-promotion-checkpoint",
		SourceUP100CSeal:"648c353adc5eed3e0b1ec43e28846f71405843a6",
		ActiveKeys:32,HotKeys:16,FilterBits:1024,HashCount:2,GenerationInterval:32,PostTriggerCounter:8,TriggerThreshold:8,
		ChurnWrites:12288,ExactRecallCap:16,ExplicitPhaseLabelUsed:false,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{191000000,192000000}
	for _,order:=range []string{"ascending","reverse","evens_then_odds"} {
		result.Points=append(result.Points,
			up101cRun("continuous_control",order,seedBases),
			up101cRun("trigger8",order,seedBases),
		)
	}
	return result,nil
}
