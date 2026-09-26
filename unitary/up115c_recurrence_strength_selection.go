package unitary

const UP115CStrengthSelectionSchema = "wingless.up115c-recurrence-strength-selection.v1"

type UP115CPoint struct {
	Arm                      string  `json:"arm"`
	AdmissionSightings       int     `json:"admission_sightings"`
	WeakAdmissionRate        float64 `json:"weak_admission_rate"`
	StrongAdmissionRate      float64 `json:"strong_admission_rate"`
	HotAccuracy              float64 `json:"hot_accuracy"`
	WeakAccuracy             float64 `json:"weak_accuracy"`
	StrongAccuracy           float64 `json:"strong_accuracy"`
	Target16Accuracy         float64 `json:"target16_accuracy"`
	Target16ExactAccuracy    float64 `json:"target16_exact_accuracy"`
	OneShotFalseAdmissions   int     `json:"one_shot_false_admissions"`
	RecallEntriesUsed        int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode float64 `json:"mean_checkpoints_per_episode"`
	AdmissionMemoryBytes     int     `json:"admission_memory_bytes"`
	TotalBoundedMemoryBytes  int     `json:"total_bounded_memory_bytes"`
}

type UP115CStrengthSelectionResult struct {
	Schema                     string        `json:"schema"`
	Experiment                 string        `json:"experiment"`
	SourceUP114CSeal           string        `json:"source_up114c_seal"`
	HotKeys                    int           `json:"hot_keys"`
	WeakCandidates             int           `json:"weak_candidates"`
	StrongCandidates           int           `json:"strong_candidates"`
	GenerationInterval         int           `json:"generation_interval"`
	AdmissionMemoryBytes       int           `json:"admission_memory_bytes"`
	ExactRecallCap             int           `json:"exact_recall_cap"`
	EpisodesPerSeed            int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission bool       `json:"query_evidence_used_for_admission"`
	SemanticPriorityUsed       bool          `json:"semantic_priority_used"`
	FutureOracleUsed           bool          `json:"future_oracle_used"`
	PhaseLabelUsed             bool          `json:"phase_label_used"`
	Points                     []UP115CPoint `json:"points"`
}

const up115cKeyMask uint32 = (1<<30)-1

type up115cStrengthMachine struct {
	mem         up81cAging
	table       [32]uint32
	counter     int
	reuseMask   uint16
	checkpoints int
	threshold   int
}

func up115cPack(key,count int) uint32 {
	return (uint32(count&3)<<30)|((uint32(key)+1)&up115cKeyMask)
}

func up115cKey(v uint32) int { return int((v&up115cKeyMask)-1) }
func up115cCount(v uint32) int { return int(v>>30) }

func (x *up115cStrengthMachine) clearTable(){for i:=range x.table{x.table[i]=0}}

func (x *up115cStrengthMachine) checkpoint(){
	x.clearTable()
	x.counter=8
	x.reuseMask=0
	x.checkpoints++
}

func (x *up115cStrengthMachine) find(key int) int {
	for i,v:=range x.table {
		if v!=0&&up115cKey(v)==key{return i}
	}
	return -1
}

func (x *up115cStrengthMachine) insert(key int){
	for i,v:=range x.table{
		if v==0{x.table[i]=up115cPack(key,1);return}
	}
	panic("UP115C probation table full before generation boundary")
}

func (x *up115cStrengthMachine) process(key,value int)(admitted,rejected bool){
	existing:=x.mem.find(key)
	if existing>=0{
		x.mem.write(key,value)
		x.reuseMask|=uint16(1)<<uint(existing)
		if x.reuseMask==0xffff{x.checkpoint()}
		return true,false
	}
	if x.mem.count<16{
		x.mem.write(key,value)
		return true,false
	}

	idx:=x.find(key)
	if idx<0{
		x.insert(key)
		rejected=true
	}else{
		count:=up115cCount(x.table[idx])+1
		if count>=x.threshold{
			x.mem.write(key,value)
			x.table[idx]=0
			admitted=true
		}else{
			x.table[idx]=up115cPack(key,count)
			rejected=true
		}
	}

	x.counter++
	if x.counter==32{
		x.clearTable()
		x.counter=0
	}

	if admitted{
		idx=x.mem.find(key)
		if idx>=0{
			x.reuseMask &^= uint16(1)<<uint(idx)
			x.reuseMask |= uint16(1)<<uint(idx)
		}
		if x.reuseMask==0xffff{x.checkpoint()}
	}
	return
}

func up115cQueryTargets(x *up115cStrengthMachine,weak,strong []int){
	for k:=0;k<12;k++{x.mem.query(k)}
	for _,k:=range weak{x.mem.query(k)}
	for _,k:=range strong{x.mem.query(k)}
}

func up115cFillBoundary(x *up115cStrengthMachine,next *int,rng *sq0RNG)(falseAdmissions int){
	for x.counter!=0{
		key:=*next;(*next)++
		admitted,_:=x.process(key,rng.intn(32))
		if admitted{falseAdmissions++}
	}
	return
}

func up115cPresent(x *up115cStrengthMachine,key,value,sightings int,next *int,rng *sq0RNG,weak,strong []int)(admitted bool,falseAdmissions int){
	if x.counter!=0{falseAdmissions+=up115cFillBoundary(x,next,rng)}
	first,_:=x.process(key,value)
	if first{falseAdmissions++}
	for sight:=2;sight<=sightings;sight++{
		for i:=0;i<7;i++{
			k:=*next;(*next)++
			a,_:=x.process(k,rng.intn(32))
			if a{falseAdmissions++}
			if i%4==3{up115cQueryTargets(x,weak,strong)}
		}
		a,_:=x.process(key,value)
		if a{admitted=true}
	}
	if x.counter!=0{falseAdmissions+=up115cFillBoundary(x,next,rng)}
	up115cQueryTargets(x,weak,strong)
	return
}

func up115cRun(arm string,threshold int,seeds []int) UP115CPoint{
	const churn=49152
	weakAdmit,strongAdmit:=0,0
	weakTotal,strongTotal:=0,0
	hotHits,hotTotal:=0,0
	weakHits,weakEval:=0,0
	strongHits,strongEval:=0,0
	targetHits,targetTotal:=0,0
	targetExact,episodes:=0,0
	falseAdmissions,maxEntries,totalCheckpoints:=0,0,0

	for _,base:=range seeds{
		for ep:=0;ep<32;ep++{
			rng:=newSQ0RNG(sq0Seed(base,5301+threshold*163,ep))
			x:=&up115cStrengthMachine{threshold:threshold}
			truth:=map[int]int{}
			for k:=0;k<32;k++{
				v:=rng.intn(32);truth[k]=v
				a,_:=x.process(k,v)
				if k>=16&&a{falseAdmissions++}
			}
			for k:=0;k<12;k++{x.process(k,truth[k]);x.mem.query(k)}

			weak:=[]int{100,101,102,103}
			strong:=[]int{200,201,202,203}
			next:=1000000+ep*300000
			if x.counter!=0{falseAdmissions+=up115cFillBoundary(x,&next,rng)}

			for _,key:=range weak{
				v:=rng.intn(32);truth[key]=v
				a,fp:=up115cPresent(x,key,v,2,&next,rng,weak,strong)
				falseAdmissions+=fp;weakTotal++;if a{weakAdmit++}
			}
			for _,key:=range strong{
				v:=rng.intn(32);truth[key]=v
				a,fp:=up115cPresent(x,key,v,3,&next,rng,weak,strong)
				falseAdmissions+=fp;strongTotal++;if a{strongAdmit++}
			}

			for j:=0;j<churn;j++{
				key:=2000000+ep*300000+j
				a,_:=x.process(key,rng.intn(32))
				if a{falseAdmissions++}
				if (j+1)%4==0{up115cQueryTargets(x,weak,strong)}
			}

			exact:=true
			for k:=0;k<12;k++{
				got,ok:=x.mem.query(k);hotTotal++;targetTotal++
				if ok&&got==truth[k]{hotHits++;targetHits++}else{exact=false}
			}
			for _,key:=range weak{
				got,ok:=x.mem.query(key);weakEval++
				if ok&&got==truth[key]{weakHits++}
			}
			for _,key:=range strong{
				got,ok:=x.mem.query(key);strongEval++;targetTotal++
				if ok&&got==truth[key]{strongHits++;targetHits++}else{exact=false}
			}
			if exact{targetExact++}
			if x.mem.count>maxEntries{maxEntries=x.mem.count}
			totalCheckpoints+=x.checkpoints
			episodes++
		}
	}

	return UP115CPoint{
		Arm:arm,AdmissionSightings:threshold,
		WeakAdmissionRate:float64(weakAdmit)/float64(weakTotal),
		StrongAdmissionRate:float64(strongAdmit)/float64(strongTotal),
		HotAccuracy:float64(hotHits)/float64(hotTotal),
		WeakAccuracy:float64(weakHits)/float64(weakEval),
		StrongAccuracy:float64(strongHits)/float64(strongEval),
		Target16Accuracy:float64(targetHits)/float64(targetTotal),
		Target16ExactAccuracy:float64(targetExact)/float64(episodes),
		OneShotFalseAdmissions:falseAdmissions,RecallEntriesUsed:maxEntries,
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		AdmissionMemoryBytes:128,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP115C()(UP115CStrengthSelectionResult,error){
	result:=UP115CStrengthSelectionResult{
		Schema:UP115CStrengthSelectionSchema,Experiment:"UP-115C-recurrence-strength-selection",
		SourceUP114CSeal:"775d35badbbd0a320f85739598575f1655fecb2b",
		HotKeys:12,WeakCandidates:4,StrongCandidates:4,GenerationInterval:32,AdmissionMemoryBytes:128,ExactRecallCap:16,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,PhaseLabelUsed:false,
	}
	seeds:=[]int{217000000,218000000}
	result.Points=append(result.Points,
		up115cRun("exact2_control",2,seeds),
		up115cRun("exact3_strength",3,seeds),
	)
	return result,nil
}
