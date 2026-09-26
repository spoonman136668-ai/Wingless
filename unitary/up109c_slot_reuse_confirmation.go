package unitary

const UP109CSlotReuseSchema = "wingless.up109c-slot-reuse-confirmation.v1"

type UP109CPoint struct {
	Arm                       string  `json:"arm"`
	InitialOrder              string  `json:"initial_order"`
	MeanCheckpointsPerEpisode float64 `json:"mean_checkpoints_per_episode"`
	MeanFirstCheckpointIndex  float64 `json:"mean_first_checkpoint_filtered_write_index"`
	ResidentHotBeforeChurn    float64 `json:"mean_resident_hot_before_churn"`
	FalsePositiveAdmissions   int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy     float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy  float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed         int     `json:"recall_entries_used"`
	PolicyMetadataBytes       int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes   int     `json:"total_bounded_memory_bytes"`
}

type UP109CSlotReuseResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP108CSeal            string        `json:"source_up108c_seal"`
	HotKeys                     int           `json:"hot_keys"`
	FilterBits                  int           `json:"filter_bits"`
	GenerationInterval          int           `json:"generation_interval"`
	ChurnWrites                 int           `json:"churn_writes"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	ReuseMaskBytes              int           `json:"reuse_mask_bytes"`
	PhasefreeArmUsesPhaseLabel  bool          `json:"phasefree_arm_uses_phase_label"`
	QueryLabelsUsedForDetector  bool          `json:"query_labels_used_for_detector"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	Points                      []UP109CPoint `json:"points"`
}

type up109cMachine struct {
	core                 *up95cMemory
	reuseMask            uint16
	checkpoints          int
	firstCheckpointIndex int
	filteredWrites       int
	mode                 string
}

func (x *up109cMachine) checkpoint() {
	x.core.clearSeen()
	x.core.counter=8
	x.reuseMask=0
	x.checkpoints++
	if x.firstCheckpointIndex==0 { x.firstCheckpointIndex=x.filteredWrites }
}

func (x *up109cMachine) process(key,value int)(admitted,rejected bool){
	existingIdx:=x.core.mem.find(key)
	full:=x.core.mem.count>=16
	seenBefore:=false
	if existingIdx<0&&full {
		x.filteredWrites++
		seenBefore=x.core.seenBefore(key)
	}
	admitted,rejected=x.core.write(key,value)

	switch x.mode {
	case "resident_write_mask":
		if existingIdx>=0 {
			x.reuseMask|=uint16(1)<<uint(existingIdx)
			if x.reuseMask==0xffff { x.checkpoint() }
		} else if admitted {
			newIdx:=x.core.mem.find(key)
			if newIdx>=0 { x.reuseMask &^= uint16(1)<<uint(newIdx) }
		}
	case "slot_reuse_confirmation":
		if existingIdx>=0 {
			x.reuseMask|=uint16(1)<<uint(existingIdx)
		} else if admitted {
			newIdx:=x.core.mem.find(key)
			if newIdx>=0 {
				x.reuseMask &^= uint16(1)<<uint(newIdx)
				if full&&seenBefore { x.reuseMask|=uint16(1)<<uint(newIdx) }
			}
		}
		if x.reuseMask==0xffff { x.checkpoint() }
	}
	return
}

func up109cRun(arm,initialOrder string,seedBases []int) UP109CPoint {
	const active=32
	const churn=12288
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	residentSum,episodes,maxEntries:=0,0,0
	falseAdmissions,totalCheckpoints,firstCheckpointSum,firstCheckpointEpisodes:=0,0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,4501+len(initialOrder)*113+len(arm)*43,ep)
			rng:=newSQ0RNG(seed)
			x:=&up109cMachine{
				core:&up95cMemory{width:width,seen:make([]uint64,width/64)},
				mode:arm,
			}
			truth:=make([]int,active)
			for k:=0;k<active;k++ { truth[k]=rng.intn(32) }

			for _,k:=range up104cInitialOrder(initialOrder) { x.process(k,truth[k]) }
			for pass:=0;pass<2;pass++ {
				written:=0
				for _,k:=range up103cSecondWriteOrder("alternating_ends") {
					x.process(k,truth[k])
					written++
					if written%4==0 { up102cHotSweep(x.core) }
				}
			}

			if arm=="explicit_checkpoint" { x.checkpoint() }

			resHot:=0
			for k:=0;k<active;k++ {
				if up85cIsHot(k,16)&&x.core.mem.find(k)>=0 { resHot++ }
			}
			residentSum+=resHot

			for j:=0;j<churn;j++ {
				admitted,_:=x.process(active+j,rng.intn(32))
				if admitted { falseAdmissions++ }
				if (j+1)%4==0 { up102cHotSweep(x.core) }
			}

			exact:=true
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,16){continue}
				got,ok:=x.core.query(k)
				hotTotal++
				if ok&&got==truth[k]{hotHits++}else{exact=false}
			}
			if exact { hotExactHits++ }
			if x.core.mem.count>maxEntries { maxEntries=x.core.mem.count }
			totalCheckpoints+=x.checkpoints
			if x.firstCheckpointIndex>0 {
				firstCheckpointSum+=x.firstCheckpointIndex
				firstCheckpointEpisodes++
			}
			episodes++
		}
	}
	firstMean:=0.0
	if firstCheckpointEpisodes>0 { firstMean=float64(firstCheckpointSum)/float64(firstCheckpointEpisodes) }
	meta:=134
	if arm=="resident_write_mask"||arm=="slot_reuse_confirmation" { meta=136 }
	return UP109CPoint{
		Arm:arm,InitialOrder:initialOrder,
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		MeanFirstCheckpointIndex:firstMean,
		ResidentHotBeforeChurn:float64(residentSum)/float64(episodes),
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:meta,TotalBoundedMemoryBytes:maxEntries*16+meta,
	}
}

func RunUP109C()(UP109CSlotReuseResult,error){
	result:=UP109CSlotReuseResult{
		Schema:UP109CSlotReuseSchema,Experiment:"UP-109C-slot-reuse-confirmation",
		SourceUP108CSeal:"05e006e769ccb16b6f8ff78eac02edd29ae8dc77",
		HotKeys:16,FilterBits:1024,GenerationInterval:32,ChurnWrites:12288,ExactRecallCap:16,ReuseMaskBytes:2,
		PhasefreeArmUsesPhaseLabel:false,QueryLabelsUsedForDetector:false,FutureOracleUsed:false,
	}
	seedBases:=[]int{199000000,200000000}
	for _,initial:=range []string{"ascending","evens_then_odds","odds_then_evens","rotate8"} {
		for _,arm:=range []string{"explicit_checkpoint","resident_write_mask","slot_reuse_confirmation"} {
			result.Points=append(result.Points,up109cRun(arm,initial,seedBases))
		}
	}
	return result,nil
}
