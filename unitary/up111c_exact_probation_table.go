package unitary

const UP111CExactProbationSchema = "wingless.up111c-exact-probation-table.v1"

type UP111CPoint struct {
	Arm                       string  `json:"arm"`
	InitialOrder              string  `json:"initial_order"`
	ChurnWrites               int     `json:"churn_writes"`
	ResidentHotBeforeChurn    float64 `json:"mean_resident_hot_before_churn"`
	MeanCheckpointsPerEpisode float64 `json:"mean_checkpoints_per_episode"`
	FalsePositiveAdmissions   int     `json:"false_positive_churn_admissions"`
	FinalHotQueryAccuracy     float64 `json:"final_hot_query_accuracy"`
	FinalHotSetExactAccuracy  float64 `json:"final_hot_set_exact_accuracy"`
	RecallEntriesUsed         int     `json:"recall_entries_used"`
	PolicyMetadataBytes       int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes   int     `json:"total_bounded_memory_bytes"`
}

type UP111CExactProbationResult struct {
	Schema                      string        `json:"schema"`
	Experiment                  string        `json:"experiment"`
	SourceUP110CSeal            string        `json:"source_up110c_seal"`
	HotKeys                     int           `json:"hot_keys"`
	AdmissionMemoryBytes        int           `json:"admission_memory_bytes"`
	GenerationInterval          int           `json:"generation_interval"`
	ExactRecallCap              int           `json:"exact_recall_cap"`
	ReuseMaskBytes              int           `json:"reuse_mask_bytes"`
	EpisodesPerSeed             int           `json:"episodes_per_seed"`
	PhaseLabelUsed              bool          `json:"phase_label_used"`
	QueryLabelsUsedForDetector  bool          `json:"query_labels_used_for_detector"`
	FutureOracleUsed            bool          `json:"future_oracle_used"`
	Points                      []UP111CPoint `json:"points"`
}

type up111cExactMachine struct {
	mem          up81cAging
	table        [32]uint32
	counter      int
	reuseMask    uint16
	checkpoints  int
}

func (x *up111cExactMachine) clearTable() {
	for i:=range x.table { x.table[i]=0 }
}

func (x *up111cExactMachine) checkpoint() {
	x.clearTable()
	x.counter=8
	x.reuseMask=0
	x.checkpoints++
}

func (x *up111cExactMachine) tableSeen(key int) bool {
	code:=uint32(key)+1
	for _,v:=range x.table {
		if v==code { return true }
	}
	return false
}

func (x *up111cExactMachine) tableInsert(key int) {
	code:=uint32(key)+1
	for i,v:=range x.table {
		if v==0 {
			x.table[i]=code
			return
		}
	}
	panic("UP111C probation table unexpectedly full before generation boundary")
}

func (x *up111cExactMachine) process(key,value int)(admitted,rejected bool) {
	existing:=x.mem.find(key)
	if existing>=0 {
		x.mem.write(key,value)
		x.reuseMask|=uint16(1)<<uint(existing)
		if x.reuseMask==0xffff { x.checkpoint() }
		return true,false
	}
	if x.mem.count<16 {
		x.mem.write(key,value)
		return true,false
	}

	seen:=x.tableSeen(key)
	if seen {
		x.mem.write(key,value)
		admitted=true
	}else{
		x.tableInsert(key)
		rejected=true
	}

	x.counter++
	if x.counter==32 {
		x.clearTable()
		x.counter=0
	}

	if admitted {
		idx:=x.mem.find(key)
		if idx>=0 {
			x.reuseMask &^= uint16(1)<<uint(idx)
			x.reuseMask |= uint16(1)<<uint(idx)
		}
		if x.reuseMask==0xffff { x.checkpoint() }
	}
	return
}

func up111cHotSweepExact(x *up111cExactMachine) {
	for h:=0;h<32;h++ {
		if up85cIsHot(h,16) { x.mem.query(h) }
	}
}

func up111cRunBloom(initial string,churn int,seedBases []int) UP111CPoint {
	hotHits,hotTotal,exactHits:=0,0,0
	residentSum,episodes,maxEntries:=0,0,0
	falseAdmissions,totalCheckpoints:=0,0
	for _,base:=range seedBases {
		for ep:=0;ep<16;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,4701+len(initial)*131+churn,ep))
			x:=&up109cMachine{core:&up95cMemory{width:1024,seen:make([]uint64,16)},mode:"slot_reuse_confirmation"}
			truth:=make([]int,32)
			for k:=0;k<32;k++ { truth[k]=rng.intn(32) }
			for _,k:=range up104cInitialOrder(initial) { x.process(k,truth[k]) }
			for pass:=0;pass<2;pass++ {
				w:=0
				for _,k:=range up103cSecondWriteOrder("alternating_ends") {
					x.process(k,truth[k]);w++
					if w%4==0 { up102cHotSweep(x.core) }
				}
			}
			res:=0
			for k:=0;k<32;k++ { if up85cIsHot(k,16)&&x.core.mem.find(k)>=0 { res++ } }
			residentSum+=res
			for j:=0;j<churn;j++ {
				admitted,_:=x.process(32+j,rng.intn(32))
				if admitted { falseAdmissions++ }
				if (j+1)%4==0 { up102cHotSweep(x.core) }
			}
			exact:=true
			for k:=0;k<32;k++ {
				if !up85cIsHot(k,16){continue}
				got,ok:=x.core.query(k);hotTotal++
				if ok&&got==truth[k]{hotHits++}else{exact=false}
			}
			if exact { exactHits++ }
			if x.core.mem.count>maxEntries { maxEntries=x.core.mem.count }
			totalCheckpoints+=x.checkpoints
			episodes++
		}
	}
	return UP111CPoint{
		Arm:"bloom1024",InitialOrder:initial,ChurnWrites:churn,
		ResidentHotBeforeChurn:float64(residentSum)/float64(episodes),
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(exactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:136,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func up111cRunExact(initial string,churn int,seedBases []int) UP111CPoint {
	hotHits,hotTotal,exactHits:=0,0,0
	residentSum,episodes,maxEntries:=0,0,0
	falseAdmissions,totalCheckpoints:=0,0
	for _,base:=range seedBases {
		for ep:=0;ep<16;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,4801+len(initial)*137+churn,ep))
			x:=&up111cExactMachine{}
			truth:=make([]int,32)
			for k:=0;k<32;k++ { truth[k]=rng.intn(32) }
			for _,k:=range up104cInitialOrder(initial) { x.process(k,truth[k]) }
			for pass:=0;pass<2;pass++ {
				w:=0
				for _,k:=range up103cSecondWriteOrder("alternating_ends") {
					x.process(k,truth[k]);w++
					if w%4==0 { up111cHotSweepExact(x) }
				}
			}
			res:=0
			for k:=0;k<32;k++ { if up85cIsHot(k,16)&&x.mem.find(k)>=0 { res++ } }
			residentSum+=res
			for j:=0;j<churn;j++ {
				admitted,_:=x.process(32+j,rng.intn(32))
				if admitted { falseAdmissions++ }
				if (j+1)%4==0 { up111cHotSweepExact(x) }
			}
			exact:=true
			for k:=0;k<32;k++ {
				if !up85cIsHot(k,16){continue}
				got,ok:=x.mem.query(k);hotTotal++
				if ok&&got==truth[k]{hotHits++}else{exact=false}
			}
			if exact { exactHits++ }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			totalCheckpoints+=x.checkpoints
			episodes++
		}
	}
	return UP111CPoint{
		Arm:"exact32",InitialOrder:initial,ChurnWrites:churn,
		ResidentHotBeforeChurn:float64(residentSum)/float64(episodes),
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		FalsePositiveAdmissions:falseAdmissions,
		FinalHotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		FinalHotSetExactAccuracy:float64(exactHits)/float64(episodes),
		RecallEntriesUsed:maxEntries,PolicyMetadataBytes:136,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP111C()(UP111CExactProbationResult,error){
	result:=UP111CExactProbationResult{
		Schema:UP111CExactProbationSchema,Experiment:"UP-111C-exact-probation-table",
		SourceUP110CSeal:"9954d1813002fbed244db9645572edce699cb256",
		HotKeys:16,AdmissionMemoryBytes:128,GenerationInterval:32,ExactRecallCap:16,ReuseMaskBytes:2,EpisodesPerSeed:16,
		PhaseLabelUsed:false,QueryLabelsUsedForDetector:false,FutureOracleUsed:false,
	}
	seeds:=[]int{209000000,210000000}
	for _,initial:=range []string{"ascending","reverse","evens_then_odds","odds_then_evens","interleaved_low_high","rotate8"} {
		for _,churn:=range []int{24576,49152,98304} {
			result.Points=append(result.Points,
				up111cRunBloom(initial,churn,seeds),
				up111cRunExact(initial,churn,seeds),
			)
		}
	}
	return result,nil
}
