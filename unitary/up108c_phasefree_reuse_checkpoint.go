package unitary

const UP108CPhasefreeReuseSchema = "wingless.up108c-phasefree-reuse-checkpoint.v1"

type UP108CPoint struct {
	Arm                         string  `json:"arm"`
	InitialOrder                string  `json:"initial_order"`
	MeanCheckpointsPerEpisode   float64 `json:"mean_checkpoints_per_episode"`
	MeanFirstCheckpointIndex    float64 `json:"mean_first_checkpoint_filtered_write_index"`
	ResidentHotBeforeChurn      float64 `json:"mean_resident_hot_before_churn"`
	FalsePositiveAdmissions     int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy       float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy    float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed           int     `json:"recall_entries_used"`
	PolicyMetadataBytes         int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes     int     `json:"total_bounded_memory_bytes"`
}

type UP108CPhasefreeReuseResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP107CSeal            string        `json:"source_up107c_seal"`
	HotKeys                     int           `json:"hot_keys"`
	FilterBits                  int           `json:"filter_bits"`
	GenerationInterval          int           `json:"generation_interval"`
	ChurnWrites                 int           `json:"churn_writes"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	ReuseMaskBytes              int           `json:"reuse_mask_bytes"`
	PhasefreeArmUsesPhaseLabel  bool          `json:"phasefree_arm_uses_phase_label"`
	QueryLabelsUsedForDetector  bool          `json:"query_labels_used_for_detector"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	Points                      []UP108CPoint `json:"points"`
}

type up108cMachine struct {
	core                  *up95cMemory
	reuseMask             uint16
	checkpoints           int
	firstCheckpointIndex  int
	filteredWrites        int
	repeatAdmissions      int
	triggerFired          bool
	phasefree             bool
}

func (x *up108cMachine) checkpoint() {
	x.core.clearSeen()
	x.core.counter=8
	x.reuseMask=0
	x.checkpoints++
	if x.firstCheckpointIndex==0 { x.firstCheckpointIndex=x.filteredWrites }
}

func (x *up108cMachine) process(key,value int)(admitted,rejected bool){
	existingIdx:=x.core.mem.find(key)
	full:=x.core.mem.count>=16
	seenBefore:=false
	if existingIdx<0&&full {
		x.filteredWrites++
		seenBefore=x.core.seenBefore(key)
	}
	admitted,rejected=x.core.write(key,value)

	if existingIdx>=0 {
		if x.phasefree {
			x.reuseMask|=uint16(1)<<uint(existingIdx)
			if x.reuseMask==0xffff { x.checkpoint() }
		}
		return
	}

	if admitted {
		newIdx:=x.core.mem.find(key)
		if newIdx>=0 && x.phasefree {
			x.reuseMask &^= uint16(1)<<uint(newIdx)
		}
	}

	if full&&seenBefore&&admitted {
		x.repeatAdmissions++
		if !x.triggerFired&&x.repeatAdmissions==8 {
			x.core.clearSeen()
			x.core.counter=8
			x.triggerFired=true
		}
	}
	return
}

func up108cRun(arm,initialOrder string,seedBases []int) UP108CPoint {
	const active=32
	const churn=12288
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	residentSum,episodes,maxEntries:=0,0,0
	falseAdmissions,totalCheckpoints,firstCheckpointSum,firstCheckpointEpisodes:=0,0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,4401+len(initialOrder)*109+len(arm)*41,ep)
			rng:=newSQ0RNG(seed)
			x:=&up108cMachine{
				core:&up95cMemory{width:width,seen:make([]uint64,width/64)},
				phasefree:arm=="phasefree_reuse_checkpoint",
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

			if arm=="explicit_checkpoint" {
				x.core.clearSeen()
				x.core.counter=8
				x.checkpoints++
				if x.firstCheckpointIndex==0 { x.firstCheckpointIndex=x.filteredWrites }
			}

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
			if exact{hotExactHits++}
			if x.core.mem.count>maxEntries{maxEntries=x.core.mem.count}
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
	if arm=="phasefree_reuse_checkpoint" { meta=136 }
	return UP108CPoint{
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

func RunUP108C()(UP108CPhasefreeReuseResult,error){
	result:=UP108CPhasefreeReuseResult{
		Schema:UP108CPhasefreeReuseSchema,Experiment:"UP-108C-phasefree-reuse-checkpoint",
		SourceUP107CSeal:"e96eec6518d83079735325a1728d3db4d191cc52",
		HotKeys:16,FilterBits:1024,GenerationInterval:32,ChurnWrites:12288,ExactRecallCap:16,ReuseMaskBytes:2,
		PhasefreeArmUsesPhaseLabel:false,QueryLabelsUsedForDetector:false,FutureOracleUsed:false,
	}
	seedBases:=[]int{205000000,206000000}
	for _,initial:=range []string{"odds_then_evens","rotate8","evens_then_odds","ascending"} {
		for _,arm:=range []string{"explicit_checkpoint","phasefree_reuse_checkpoint","no_checkpoint"} {
			result.Points=append(result.Points,up108cRun(arm,initial,seedBases))
		}
	}
	return result,nil
}
