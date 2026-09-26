package unitary

const UP110CSlotReuseGeneralitySchema = "wingless.up110c-slot-reuse-generality.v1"

type UP110CPoint struct {
	Arm                       string  `json:"arm"`
	InitialOrder              string  `json:"initial_order"`
	ChurnWrites               int     `json:"churn_writes"`
	ResidentHotBeforeChurn    float64 `json:"mean_resident_hot_before_churn"`
	MeanCheckpointsPerEpisode float64 `json:"mean_checkpoints_per_episode"`
	MeanFirstCheckpointIndex  float64 `json:"mean_first_checkpoint_filtered_write_index"`
	FalsePositiveAdmissions   int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy     float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy  float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed         int     `json:"recall_entries_used"`
	PolicyMetadataBytes       int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes   int     `json:"total_bounded_memory_bytes"`
}

type UP110CSlotReuseGeneralityResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP109CSeal            string        `json:"source_up109c_seal"`
	HotKeys                     int           `json:"hot_keys"`
	FilterBits                  int           `json:"filter_bits"`
	GenerationInterval          int           `json:"generation_interval"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	ReuseMaskBytes              int           `json:"reuse_mask_bytes"`
	EpisodesPerSeed             int           `json:"episodes_per_seed"`
	PhasefreeArmUsesPhaseLabel  bool          `json:"phasefree_arm_uses_phase_label"`
	QueryLabelsUsedForDetector  bool          `json:"query_labels_used_for_detector"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	Points                      []UP110CPoint `json:"points"`
}

func up110cRun(arm,initialOrder string,churn int,seedBases []int) UP110CPoint {
	const active=32
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	residentSum,episodes,maxEntries:=0,0,0
	falseAdmissions,totalCheckpoints,firstCheckpointSum,firstCheckpointEpisodes:=0,0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<32;ep++ {
			seed:=sq0Seed(base,4601+len(initialOrder)*127+len(arm)*47+churn,ep)
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
	if arm=="slot_reuse_confirmation" { meta=136 }
	return UP110CPoint{
		Arm:arm,InitialOrder:initialOrder,ChurnWrites:churn,
		ResidentHotBeforeChurn:float64(residentSum)/float64(episodes),
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		MeanFirstCheckpointIndex:firstMean,
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:meta,TotalBoundedMemoryBytes:maxEntries*16+meta,
	}
}

func RunUP110C()(UP110CSlotReuseGeneralityResult,error){
	result:=UP110CSlotReuseGeneralityResult{
		Schema:UP110CSlotReuseGeneralitySchema,Experiment:"UP-110C-slot-reuse-generality",
		SourceUP109CSeal:"5a328c4115d82061704a01b6431fb6b6bcb50dcb",
		HotKeys:16,FilterBits:1024,GenerationInterval:32,ExactRecallCap:16,ReuseMaskBytes:2,EpisodesPerSeed:32,
		PhasefreeArmUsesPhaseLabel:false,QueryLabelsUsedForDetector:false,FutureOracleUsed:false,
	}
	seedBases:=[]int{207000000,208000000}
	for _,initial:=range []string{"ascending","reverse","evens_then_odds","odds_then_evens","interleaved_low_high","rotate8"} {
		for _,churn:=range []int{12288,24576,49152} {
			for _,arm:=range []string{"explicit_checkpoint","slot_reuse_confirmation"} {
				result.Points=append(result.Points,up110cRun(arm,initial,churn,seedBases))
			}
		}
	}
	return result,nil
}
