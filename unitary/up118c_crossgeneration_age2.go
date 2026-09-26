package unitary

const UP118CAge2Schema = "wingless.up118c-crossgeneration-age2.v1"

type UP118CPoint struct {
	Arm                       string  `json:"arm"`
	Workload                  string  `json:"workload"`
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

type UP118CAge2Result struct {
	Schema                        string        `json:"schema"`
	Experiment                    string        `json:"experiment"`
	SourceUP117CSeal              string        `json:"source_up117c_seal"`
	HotKeys                       int           `json:"hot_keys"`
	TargetCandidates              int           `json:"target_candidates"`
	BurstDistractors              int           `json:"burst_distractors"`
	GenerationInterval            int           `json:"generation_interval"`
	AdmissionMemoryBytes          int           `json:"admission_memory_bytes"`
	ExactRecallCap                int           `json:"exact_recall_cap"`
	KeyBits                       int           `json:"key_bits"`
	ChurnWrites                   int           `json:"churn_writes"`
	EpisodesPerSeed               int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission bool          `json:"query_evidence_used_for_admission"`
	SemanticPriorityUsed          bool          `json:"semantic_priority_used"`
	FutureOracleUsed              bool          `json:"future_oracle_used"`
	PhaseLabelUsed                bool          `json:"phase_label_used"`
	Points                        []UP118CPoint `json:"points"`
}

const up118cKeyMask uint16 = (1<<14)-1
const up118cCurrentQualified uint16 = 1<<14

type up118cAge2Machine struct {
	mem         up81cAging
	current     [32]uint16
	history     [32]uint16
	counter     int
	reuseMask   uint16
	checkpoints int
}

func up118cCode(key int) uint16 { return uint16(key+1)&up118cKeyMask }
func up118cKey(v uint16) int { return int((v&up118cKeyMask)-1) }

func up118cFind(table *[32]uint16,key int) int {
	code:=up118cCode(key)
	for i,v:=range table {
		if v&up118cKeyMask==code && v&up118cKeyMask!=0 { return i }
	}
	return -1
}

func up118cInsertCurrent(table *[32]uint16,key int,qualified bool) {
	code:=up118cCode(key)
	if qualified { code|=up118cCurrentQualified }
	for i,v:=range table { if v==0 { table[i]=code; return } }
	panic("UP118C current table full")
}

func up118cInsertHistory(table *[32]uint16,key,age int) {
	code:=up118cCode(key)|(uint16(age&3)<<14)
	if i:=up118cFind(table,key);i>=0 { table[i]=code; return }
	for i,v:=range table { if v==0 { table[i]=code; return } }
	panic("UP118C history table full")
}

func up118cClear(table *[32]uint16){ for i:=range table { table[i]=0 } }

func (x *up118cAge2Machine) boundary() {
	for i,v:=range x.history {
		if v==0 { continue }
		age:=int(v>>14)
		if age>=1 { x.history[i]=0 } else { x.history[i]=(v&up118cKeyMask)|(uint16(1)<<14) }
	}
	for _,v:=range x.current {
		if v&up118cCurrentQualified!=0 { up118cInsertHistory(&x.history,up118cKey(v),0) }
	}
	up118cClear(&x.current)
	x.counter=0
}

func (x *up118cAge2Machine) checkpoint(){
	up118cClear(&x.current);up118cClear(&x.history)
	x.counter=8;x.reuseMask=0;x.checkpoints++
}

func (x *up118cAge2Machine) process(key,value int)(admitted,rejected bool){
	existing:=x.mem.find(key)
	if existing>=0{
		x.mem.write(key,value)
		x.reuseMask|=uint16(1)<<uint(existing)
		if x.reuseMask==0xffff{x.checkpoint()}
		return true,false
	}
	if x.mem.count<16{x.mem.write(key,value);return true,false}

	idx:=up118cFind(&x.current,key)
	if idx<0{
		up118cInsertCurrent(&x.current,key,false);rejected=true
	}else if x.current[idx]&up118cCurrentQualified==0{
		x.current[idx]|=up118cCurrentQualified
		if hi:=up118cFind(&x.history,key);hi>=0{
			x.mem.write(key,value);admitted=true;x.current[idx]=0;x.history[hi]=0
		}else{rejected=true}
	}else{rejected=true}

	x.counter++
	if x.counter==32{x.boundary()}
	if admitted{
		idx=x.mem.find(key)
		if idx>=0{x.reuseMask|=uint16(1)<<uint(idx)}
		if x.reuseMask==0xffff{x.checkpoint()}
	}
	return
}

func up118cFillControl(x *up117cCrossGenMachine,next *int,rng *sq0RNG)(fp int){
	for x.counter!=0{key:=*next;(*next)++;a,_:=x.process(key,rng.intn(32));if a{fp++}}
	return
}
func up118cFillAge2(x *up118cAge2Machine,next *int,rng *sq0RNG)(fp int){
	for x.counter!=0{key:=*next;(*next)++;a,_:=x.process(key,rng.intn(32));if a{fp++}}
	return
}

func up118cEval(mem *up81cAging,truth map[int]int,targets,bursts []int)(hotHits,hotTotal,targetHits,targetTotal,burstHits,burstTotal,targetSetHits,targetSetTotal int,exact bool){
	exact=true
	for k:=0;k<12;k++{got,ok:=mem.query(k);hotTotal++;targetSetTotal++;if ok&&got==truth[k]{hotHits++;targetSetHits++}else{exact=false}}
	for _,k:=range targets{got,ok:=mem.query(k);targetTotal++;targetSetTotal++;if ok&&got==truth[k]{targetHits++;targetSetHits++}else{exact=false}}
	for _,k:=range bursts{got,ok:=mem.query(k);burstTotal++;if ok&&got==truth[k]{burstHits++}}
	return
}

func up118cRunControl(workload string,seeds []int) UP118CPoint{
	const churn=12288
	targets:=[]int{100,101,102,103};bursts:=[]int{200,201,202,203}
	targetAdmit,burstAdmit,targetTrials,burstTrials:=0,0,0,0
	hotHits,hotTotal,targetHits,targetEval,burstHits,burstEval:=0,0,0,0,0,0
	setHits,setTotal,exactHits,episodes,fp,maxEntries,checkpoints:=0,0,0,0,0,0,0
	for _,base:=range seeds{
		for ep:=0;ep<32;ep++{
			rng:=newSQ0RNG(sq0Seed(base,5701+len(workload)*181,ep))
			x:=&up117cCrossGenMachine{};truth:=map[int]int{}
			for k:=0;k<32;k++{v:=rng.intn(32);truth[k]=v;a,_:=x.process(k,v);if k>=16&&a{fp++}}
			for k:=0;k<12;k++{x.process(k,truth[k]);x.mem.query(k)}
			next:=300
			if x.counter!=0{fp+=up118cFillControl(x,&next,rng)}

			for _,k:=range targets{v:=rng.intn(32);truth[k]=v;for s:=0;s<2;s++{was:=x.mem.find(k)>=0;a,_:=x.process(k,v);if !was&&a{targetAdmit++}};targetTrials++}
			for _,k:=range bursts{v:=rng.intn(32);truth[k]=v;for s:=0;s<4;s++{was:=x.mem.find(k)>=0;a,_:=x.process(k,v);if !was&&a{burstAdmit++}};burstTrials++}
			if x.counter!=0{fp+=up118cFillControl(x,&next,rng)}

			if workload=="skip_one_generation"{if x.counter!=0{fp+=up118cFillControl(x,&next,rng)};for x.counter==0{for i:=0;i<32;i++{a,_:=x.process(next,rng.intn(32));next++;if a{fp++}};break}}
			for _,k:=range targets{for s:=0;s<2;s++{was:=x.mem.find(k)>=0;a,_:=x.process(k,truth[k]);if !was&&a{targetAdmit++}}}
			if x.counter!=0{fp+=up118cFillControl(x,&next,rng)}

			for j:=0;j<churn;j++{key:=2000+j;a,_:=x.process(key,rng.intn(32));if a{fp++};if (j+1)%4==0{up117cQuery(&x.mem,targets,bursts)}}
			hh,ht,th,tt,bh,bt,sh,st,exact:=up118cEval(&x.mem,truth,targets,bursts)
			hotHits+=hh;hotTotal+=ht;targetHits+=th;targetEval+=tt;burstHits+=bh;burstEval+=bt;setHits+=sh;setTotal+=st;if exact{exactHits++}
			if x.mem.count>maxEntries{maxEntries=x.mem.count};checkpoints+=x.checkpoints;episodes++
		}
	}
	return UP118CPoint{Arm:"previous_generation_control",Workload:workload,TargetAdmissionRate:float64(targetAdmit)/float64(targetTrials),BurstAdmissionRate:float64(burstAdmit)/float64(burstTrials),HotAccuracy:float64(hotHits)/float64(hotTotal),TargetAccuracy:float64(targetHits)/float64(targetEval),BurstAccuracy:float64(burstHits)/float64(burstEval),Target16Accuracy:float64(setHits)/float64(setTotal),Target16ExactAccuracy:float64(exactHits)/float64(episodes),OneShotFalseAdmissions:fp,RecallEntriesUsed:maxEntries,MeanCheckpointsPerEpisode:float64(checkpoints)/float64(episodes),AdmissionMemoryBytes:128,PolicyMetadataBytes:136,TotalBoundedMemoryBytes:maxEntries*16+136}
}

func up118cRunAge2(workload string,seeds []int) UP118CPoint{
	const churn=12288
	targets:=[]int{100,101,102,103};bursts:=[]int{200,201,202,203}
	targetAdmit,burstAdmit,targetTrials,burstTrials:=0,0,0,0
	hotHits,hotTotal,targetHits,targetEval,burstHits,burstEval:=0,0,0,0,0,0
	setHits,setTotal,exactHits,episodes,fp,maxEntries,checkpoints:=0,0,0,0,0,0,0
	for _,base:=range seeds{
		for ep:=0;ep<32;ep++{
			rng:=newSQ0RNG(sq0Seed(base,5801+len(workload)*191,ep))
			x:=&up118cAge2Machine{};truth:=map[int]int{}
			for k:=0;k<32;k++{v:=rng.intn(32);truth[k]=v;a,_:=x.process(k,v);if k>=16&&a{fp++}}
			for k:=0;k<12;k++{x.process(k,truth[k]);x.mem.query(k)}
			next:=300
			if x.counter!=0{fp+=up118cFillAge2(x,&next,rng)}

			for _,k:=range targets{v:=rng.intn(32);truth[k]=v;for s:=0;s<2;s++{was:=x.mem.find(k)>=0;a,_:=x.process(k,v);if !was&&a{targetAdmit++}};targetTrials++}
			for _,k:=range bursts{v:=rng.intn(32);truth[k]=v;for s:=0;s<4;s++{was:=x.mem.find(k)>=0;a,_:=x.process(k,v);if !was&&a{burstAdmit++}};burstTrials++}
			if x.counter!=0{fp+=up118cFillAge2(x,&next,rng)}

			if workload=="skip_one_generation"{for i:=0;i<32;i++{a,_:=x.process(next,rng.intn(32));next++;if a{fp++}}}
			for _,k:=range targets{for s:=0;s<2;s++{was:=x.mem.find(k)>=0;a,_:=x.process(k,truth[k]);if !was&&a{targetAdmit++}}}
			if x.counter!=0{fp+=up118cFillAge2(x,&next,rng)}

			for j:=0;j<churn;j++{key:=2000+j;a,_:=x.process(key,rng.intn(32));if a{fp++};if (j+1)%4==0{up117cQuery(&x.mem,targets,bursts)}}
			hh,ht,th,tt,bh,bt,sh,st,exact:=up118cEval(&x.mem,truth,targets,bursts)
			hotHits+=hh;hotTotal+=ht;targetHits+=th;targetEval+=tt;burstHits+=bh;burstEval+=bt;setHits+=sh;setTotal+=st;if exact{exactHits++}
			if x.mem.count>maxEntries{maxEntries=x.mem.count};checkpoints+=x.checkpoints;episodes++
		}
	}
	return UP118CPoint{Arm:"age2_history",Workload:workload,TargetAdmissionRate:float64(targetAdmit)/float64(targetTrials),BurstAdmissionRate:float64(burstAdmit)/float64(burstTrials),HotAccuracy:float64(hotHits)/float64(hotTotal),TargetAccuracy:float64(targetHits)/float64(targetEval),BurstAccuracy:float64(burstHits)/float64(burstEval),Target16Accuracy:float64(setHits)/float64(setTotal),Target16ExactAccuracy:float64(exactHits)/float64(episodes),OneShotFalseAdmissions:fp,RecallEntriesUsed:maxEntries,MeanCheckpointsPerEpisode:float64(checkpoints)/float64(episodes),AdmissionMemoryBytes:128,PolicyMetadataBytes:136,TotalBoundedMemoryBytes:maxEntries*16+136}
}

func RunUP118C()(UP118CAge2Result,error){
	result:=UP118CAge2Result{Schema:UP118CAge2Schema,Experiment:"UP-118C-crossgeneration-age2",SourceUP117CSeal:"c2df738f6ea8974036f7138696fddedfa1a39f78",HotKeys:12,TargetCandidates:4,BurstDistractors:4,GenerationInterval:32,AdmissionMemoryBytes:128,ExactRecallCap:16,KeyBits:14,ChurnWrites:12288,EpisodesPerSeed:32,QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,PhaseLabelUsed:false}
	seeds:=[]int{223000000,224000000}
	for _,workload:=range []string{"adjacent_recurrence","skip_one_generation"}{
		result.Points=append(result.Points,up118cRunControl(workload,seeds),up118cRunAge2(workload,seeds))
	}
	return result,nil
}
