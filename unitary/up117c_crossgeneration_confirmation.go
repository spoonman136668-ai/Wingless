package unitary

const UP117CCrossGenerationSchema = "wingless.up117c-crossgeneration-confirmation.v1"

type UP117CPoint struct {
	Arm                       string  `json:"arm"`
	TargetAdmissionRate       float64 `json:"target_admission_rate"`
	BurstAdmissionRate        float64 `json:"burst_admission_rate"`
	HotAccuracy               float64 `json:"hot_accuracy"`
	TargetAccuracy            float64 `json:"target_accuracy"`
	BurstAccuracy             float64 `json:"burst_accuracy"`
	Target16Accuracy          float64 `json:"target16_accuracy"`
	Target16ExactAccuracy     float64 `json:"target16_exact_accuracy"`
	OneShotFalseAdmissions    int     `json:"one_shot_false_admissions"`
	RecallEntriesUsed         int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode float64 `json:"mean_checkpoints_per_episode"`
	AdmissionMemoryBytes      int     `json:"admission_memory_bytes"`
	PolicyMetadataBytes       int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes   int     `json:"total_bounded_memory_bytes"`
}

type UP117CCrossGenerationResult struct {
	Schema                         string        `json:"schema"`
	Experiment                     string        `json:"experiment"`
	SourceUP116CSeal               string        `json:"source_up116c_seal"`
	HotKeys                        int           `json:"hot_keys"`
	TargetCandidates               int           `json:"target_candidates"`
	BurstDistractors               int           `json:"burst_distractors"`
	GenerationInterval             int           `json:"generation_interval"`
	AdmissionMemoryBytes           int           `json:"admission_memory_bytes"`
	ExactRecallCap                 int           `json:"exact_recall_cap"`
	EpisodesPerSeed                int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission  bool          `json:"query_evidence_used_for_admission"`
	SemanticPriorityUsed           bool          `json:"semantic_priority_used"`
	FutureOracleUsed               bool          `json:"future_oracle_used"`
	PhaseLabelUsed                 bool          `json:"phase_label_used"`
	Points                         []UP117CPoint `json:"points"`
}

const up117cQualified uint16 = 1 << 15
const up117cKeyMask uint16 = up117cQualified - 1

type up117cCrossGenMachine struct {
	mem          up81cAging
	current      [32]uint16
	previous     [32]uint16
	counter      int
	reuseMask    uint16
	checkpoints  int
}

func up117cCode(key int) uint16 { return uint16(key+1) & up117cKeyMask }
func up117cKey(v uint16) int { return int((v & up117cKeyMask)-1) }

func up117cFind(table *[32]uint16,key int) int {
	code:=up117cCode(key)
	for i,v:=range table {
		if v&up117cKeyMask==code && v&up117cKeyMask!=0 { return i }
	}
	return -1
}

func up117cInsert(table *[32]uint16,key int,qualified bool) {
	code:=up117cCode(key)
	if qualified { code|=up117cQualified }
	for i,v:=range table {
		if v==0 {
			table[i]=code
			return
		}
	}
	panic("UP117C exact admission table full before generation boundary")
}

func up117cClear(table *[32]uint16) { for i:=range table { table[i]=0 } }

func (x *up117cCrossGenMachine) generationBoundary() {
	up117cClear(&x.previous)
	for _,v:=range x.current {
		if v&up117cQualified!=0 {
			up117cInsert(&x.previous,up117cKey(v),false)
		}
	}
	up117cClear(&x.current)
	x.counter=0
}

func (x *up117cCrossGenMachine) checkpoint() {
	up117cClear(&x.current)
	up117cClear(&x.previous)
	x.counter=8
	x.reuseMask=0
	x.checkpoints++
}

func (x *up117cCrossGenMachine) process(key,value int)(admitted,rejected bool) {
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

	idx:=up117cFind(&x.current,key)
	if idx<0 {
		up117cInsert(&x.current,key,false)
		rejected=true
	}else if x.current[idx]&up117cQualified==0 {
		x.current[idx]|=up117cQualified
		if up117cFind(&x.previous,key)>=0 {
			x.mem.write(key,value)
			admitted=true
			x.current[idx]=0
			if pi:=up117cFind(&x.previous,key);pi>=0 { x.previous[pi]=0 }
		}else{
			rejected=true
		}
	}else{
		rejected=true
	}

	x.counter++
	if x.counter==32 { x.generationBoundary() }

	if admitted {
		idx=x.mem.find(key)
		if idx>=0 {
			x.reuseMask &^= uint16(1)<<uint(idx)
			x.reuseMask |= uint16(1)<<uint(idx)
		}
		if x.reuseMask==0xffff { x.checkpoint() }
	}
	return
}

func up117cQuery(mem *up81cAging,targets,bursts []int) {
	for k:=0;k<12;k++ { mem.query(k) }
	for _,k:=range targets { mem.query(k) }
	for _,k:=range bursts { mem.query(k) }
}

func up117cFillExact3(x *up115cStrengthMachine,next *int,rng *sq0RNG)(falseAdmissions int) {
	return up115cFillBoundary(x,next,rng)
}

func up117cFillCross(x *up117cCrossGenMachine,next *int,rng *sq0RNG)(falseAdmissions int) {
	for x.counter!=0 {
		key:=*next
		(*next)++
		a,_:=x.process(key,rng.intn(32))
		if a { falseAdmissions++ }
	}
	return
}

func up117cRunExact3(seeds []int) UP117CPoint {
	const churn=24576
	targets:=[]int{100,101,102,103}
	bursts:=[]int{200,201,202,203}
	targetAdmit,burstAdmit,targetTotal,burstTotal:=0,0,0,0
	hotHits,hotTotal,targetHits,targetEval,burstHits,burstEval:=0,0,0,0,0,0
	targetSetHits,targetSetTotal,targetExact:=0,0,0
	falseAdmissions,maxEntries,totalCheckpoints,episodes:=0,0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,5501,ep))
			x:=&up115cStrengthMachine{threshold:3}
			truth:=map[int]int{}
			for k:=0;k<32;k++ {
				v:=rng.intn(32);truth[k]=v
				a,_:=x.process(k,v)
				if k>=16&&a { falseAdmissions++ }
			}
			for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }

			next:=1000
			if x.counter!=0 { falseAdmissions+=up117cFillExact3(x,&next,rng) }

			// Generation 1: target 2x, burst 4x, ascending key order.
			for _,key:=range targets {
				v:=rng.intn(32);truth[key]=v
				for s:=0;s<2;s++ {
					was:=x.mem.find(key)>=0
					a,_:=x.process(key,v)
					if !was&&a { targetAdmit++ }
				}
				targetTotal++
			}
			for _,key:=range bursts {
				v:=rng.intn(32);truth[key]=v
				for s:=0;s<4;s++ {
					was:=x.mem.find(key)>=0
					a,_:=x.process(key,v)
					if !was&&a { burstAdmit++ }
				}
				burstTotal++
			}
			if x.counter!=0 { falseAdmissions+=up117cFillExact3(x,&next,rng) }

			// Generation 2: targets recur 2x; bursts do not.
			for _,key:=range targets {
				for s:=0;s<2;s++ {
					was:=x.mem.find(key)>=0
					a,_:=x.process(key,truth[key])
					if !was&&a { targetAdmit++ }
				}
			}
			if x.counter!=0 { falseAdmissions+=up117cFillExact3(x,&next,rng) }

			for j:=0;j<churn;j++ {
				key:=5000+j
				a,_:=x.process(key,rng.intn(32))
				if a { falseAdmissions++ }
				if (j+1)%4==0 { up117cQuery(&x.mem,targets,bursts) }
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=x.mem.query(k);hotTotal++;targetSetTotal++
				if ok&&got==truth[k] { hotHits++;targetSetHits++ } else { exact=false }
			}
			for _,key:=range targets {
				got,ok:=x.mem.query(key);targetEval++;targetSetTotal++
				if ok&&got==truth[key] { targetHits++;targetSetHits++ } else { exact=false }
			}
			for _,key:=range bursts {
				got,ok:=x.mem.query(key);burstEval++
				if ok&&got==truth[key] { burstHits++ }
			}
			if exact { targetExact++ }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			totalCheckpoints+=x.checkpoints
			episodes++
		}
	}
	return UP117CPoint{
		Arm:"exact3_control",
		TargetAdmissionRate:float64(targetAdmit)/float64(targetTotal),
		BurstAdmissionRate:float64(burstAdmit)/float64(burstTotal),
		HotAccuracy:float64(hotHits)/float64(hotTotal),
		TargetAccuracy:float64(targetHits)/float64(targetEval),
		BurstAccuracy:float64(burstHits)/float64(burstEval),
		Target16Accuracy:float64(targetSetHits)/float64(targetSetTotal),
		Target16ExactAccuracy:float64(targetExact)/float64(episodes),
		OneShotFalseAdmissions:falseAdmissions,RecallEntriesUsed:maxEntries,
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		AdmissionMemoryBytes:128,PolicyMetadataBytes:136,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func up117cRunCross(seeds []int) UP117CPoint {
	const churn=24576
	targets:=[]int{100,101,102,103}
	bursts:=[]int{200,201,202,203}
	targetAdmit,burstAdmit,targetTotal,burstTotal:=0,0,0,0
	hotHits,hotTotal,targetHits,targetEval,burstHits,burstEval:=0,0,0,0,0,0
	targetSetHits,targetSetTotal,targetExact:=0,0,0
	falseAdmissions,maxEntries,totalCheckpoints,episodes:=0,0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,5601,ep))
			x:=&up117cCrossGenMachine{}
			truth:=map[int]int{}
			for k:=0;k<32;k++ {
				v:=rng.intn(32);truth[k]=v
				a,_:=x.process(k,v)
				if k>=16&&a { falseAdmissions++ }
			}
			for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }

			next:=1000
			if x.counter!=0 { falseAdmissions+=up117cFillCross(x,&next,rng) }

			for _,key:=range targets {
				v:=rng.intn(32);truth[key]=v
				for s:=0;s<2;s++ {
					was:=x.mem.find(key)>=0
					a,_:=x.process(key,v)
					if !was&&a { targetAdmit++ }
				}
				targetTotal++
			}
			for _,key:=range bursts {
				v:=rng.intn(32);truth[key]=v
				for s:=0;s<4;s++ {
					was:=x.mem.find(key)>=0
					a,_:=x.process(key,v)
					if !was&&a { burstAdmit++ }
				}
				burstTotal++
			}
			if x.counter!=0 { falseAdmissions+=up117cFillCross(x,&next,rng) }

			for _,key:=range targets {
				for s:=0;s<2;s++ {
					was:=x.mem.find(key)>=0
					a,_:=x.process(key,truth[key])
					if !was&&a { targetAdmit++ }
				}
			}
			if x.counter!=0 { falseAdmissions+=up117cFillCross(x,&next,rng) }

			for j:=0;j<churn;j++ {
				key:=5000+j
				a,_:=x.process(key,rng.intn(32))
				if a { falseAdmissions++ }
				if (j+1)%4==0 { up117cQuery(&x.mem,targets,bursts) }
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=x.mem.query(k);hotTotal++;targetSetTotal++
				if ok&&got==truth[k] { hotHits++;targetSetHits++ } else { exact=false }
			}
			for _,key:=range targets {
				got,ok:=x.mem.query(key);targetEval++;targetSetTotal++
				if ok&&got==truth[key] { targetHits++;targetSetHits++ } else { exact=false }
			}
			for _,key:=range bursts {
				got,ok:=x.mem.query(key);burstEval++
				if ok&&got==truth[key] { burstHits++ }
			}
			if exact { targetExact++ }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			totalCheckpoints+=x.checkpoints
			episodes++
		}
	}
	return UP117CPoint{
		Arm:"crossgen_2x2",
		TargetAdmissionRate:float64(targetAdmit)/float64(targetTotal),
		BurstAdmissionRate:float64(burstAdmit)/float64(burstTotal),
		HotAccuracy:float64(hotHits)/float64(hotTotal),
		TargetAccuracy:float64(targetHits)/float64(targetEval),
		BurstAccuracy:float64(burstHits)/float64(burstEval),
		Target16Accuracy:float64(targetSetHits)/float64(targetSetTotal),
		Target16ExactAccuracy:float64(targetExact)/float64(episodes),
		OneShotFalseAdmissions:falseAdmissions,RecallEntriesUsed:maxEntries,
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		AdmissionMemoryBytes:128,PolicyMetadataBytes:136,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP117C()(UP117CCrossGenerationResult,error) {
	result:=UP117CCrossGenerationResult{
		Schema:UP117CCrossGenerationSchema,Experiment:"UP-117C-crossgeneration-confirmation",
		SourceUP116CSeal:"f8669aaa748a610c85847512cb96d6af850b91a1",
		HotKeys:12,TargetCandidates:4,BurstDistractors:4,GenerationInterval:32,AdmissionMemoryBytes:128,ExactRecallCap:16,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,PhaseLabelUsed:false,
	}
	seeds:=[]int{221000000,222000000}
	result.Points=append(result.Points,up117cRunExact3(seeds),up117cRunCross(seeds))
	return result,nil
}
