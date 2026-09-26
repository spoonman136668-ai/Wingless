package unitary

const UP120CAge3Schema = "wingless.up120c-crossgeneration-age3.v1"

type UP120CPoint struct {
	Arm                       string  `json:"arm"`
	Workload                  string  `json:"workload"`
	SkippedGenerations        int     `json:"skipped_generations"`
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
	TotalBoundedMemoryBytes   int     `json:"total_bounded_memory_bytes"`
}

type UP120CAge3Result struct {
	Schema                        string        `json:"schema"`
	Experiment                    string        `json:"experiment"`
	SourceUP119CSeal              string        `json:"source_up119c_seal"`
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
	AdaptiveHorizonUsed           bool          `json:"adaptive_horizon_used"`
	Points                        []UP120CPoint `json:"points"`
}

type up120cExtendedMachine struct {
	mem         up81cAging
	current     [32]uint16
	history     [32]uint16
	counter     int
	reuseMask   uint16
	checkpoints int
	maxAge      int
}

func (x *up120cExtendedMachine) boundary() {
	for i,v:=range x.history {
		if v==0 { continue }
		age:=int(v>>14)
		if age>=x.maxAge {
			x.history[i]=0
		}else{
			x.history[i]=(v&up118cKeyMask)|(uint16(age+1)<<14)
		}
	}
	for _,v:=range x.current {
		if v&up118cCurrentQualified!=0 {
			up118cInsertHistory(&x.history,up118cKey(v),0)
		}
	}
	up118cClear(&x.current)
	x.counter=0
}

func (x *up120cExtendedMachine) checkpoint() {
	up118cClear(&x.current)
	up118cClear(&x.history)
	x.counter=8
	x.reuseMask=0
	x.checkpoints++
}

func (x *up120cExtendedMachine) process(key,value int)(admitted,rejected bool) {
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

	idx:=up118cFind(&x.current,key)
	if idx<0 {
		up118cInsertCurrent(&x.current,key,false)
		rejected=true
	}else if x.current[idx]&up118cCurrentQualified==0 {
		x.current[idx]|=up118cCurrentQualified
		if hi:=up118cFind(&x.history,key);hi>=0 {
			x.mem.write(key,value)
			admitted=true
			x.current[idx]=0
			x.history[hi]=0
		}else{
			rejected=true
		}
	}else{
		rejected=true
	}

	x.counter++
	if x.counter==32 { x.boundary() }

	if admitted {
		idx=x.mem.find(key)
		if idx>=0 { x.reuseMask|=uint16(1)<<uint(idx) }
		if x.reuseMask==0xffff { x.checkpoint() }
	}
	return
}

func (x *up120cExtendedMachine) fill(next *int,rng *sq0RNG)(fp int) {
	for x.counter!=0 {
		key:=*next
		(*next)++
		a,_:=x.process(key,rng.intn(32))
		if a { fp++ }
	}
	return
}

func up120cRun(arm,workload string,skip int,seeds []int) UP120CPoint {
	const churn=12288
	targets:=[]int{100,101,102,103}
	bursts:=[]int{200,201,202,203}

	targetAdmit,burstAdmit,targetTrials,burstTrials:=0,0,0,0
	hotHits,hotTotal,targetHits,targetEval,burstHits,burstEval:=0,0,0,0,0,0
	setHits,setTotal,exactHits,episodes:=0,0,0,0
	fp,maxEntries,checkpoints:=0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,6001+skip*223+len(arm)*227,ep))
			truth:=map[int]int{}
			next:=300

			if arm=="age1_control" {
				x:=&up118cAge2Machine{}
				for k:=0;k<32;k++ {
					v:=rng.intn(32);truth[k]=v
					a,_:=x.process(k,v);if k>=16&&a{fp++}
				}
				for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }
				if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

				for _,k:=range targets {
					v:=rng.intn(32);truth[k]=v
					for s:=0;s<2;s++ { was:=x.mem.find(k)>=0;a,_:=x.process(k,v);if !was&&a{targetAdmit++} }
					targetTrials++
				}
				for _,k:=range bursts {
					v:=rng.intn(32);truth[k]=v
					for s:=0;s<4;s++ { was:=x.mem.find(k)>=0;a,_:=x.process(k,v);if !was&&a{burstAdmit++} }
					burstTrials++
				}
				if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

				for g:=0;g<skip;g++ {
					for i:=0;i<32;i++ { a,_:=x.process(next,rng.intn(32));next++;if a{fp++} }
					if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }
				}
				for _,k:=range targets {
					for s:=0;s<2;s++ { was:=x.mem.find(k)>=0;a,_:=x.process(k,truth[k]);if !was&&a{targetAdmit++} }
				}
				if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

				for j:=0;j<churn;j++ {
					key:=1000+j;a,_:=x.process(key,rng.intn(32));if a{fp++}
					if (j+1)%4==0 { up117cQuery(&x.mem,targets,bursts) }
				}
				exact:=true
				for k:=0;k<12;k++ { got,ok:=x.mem.query(k);hotTotal++;setTotal++;if ok&&got==truth[k]{hotHits++;setHits++}else{exact=false} }
				for _,k:=range targets { got,ok:=x.mem.query(k);targetEval++;setTotal++;if ok&&got==truth[k]{targetHits++;setHits++}else{exact=false} }
				for _,k:=range bursts { got,ok:=x.mem.query(k);burstEval++;if ok&&got==truth[k]{burstHits++} }
				if exact{exactHits++};if x.mem.count>maxEntries{maxEntries=x.mem.count};checkpoints+=x.checkpoints
			}else{
				x:=&up120cExtendedMachine{maxAge:2}
				for k:=0;k<32;k++ {
					v:=rng.intn(32);truth[k]=v
					a,_:=x.process(k,v);if k>=16&&a{fp++}
				}
				for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }
				if x.counter!=0 { fp+=x.fill(&next,rng) }

				for _,k:=range targets {
					v:=rng.intn(32);truth[k]=v
					for s:=0;s<2;s++ { was:=x.mem.find(k)>=0;a,_:=x.process(k,v);if !was&&a{targetAdmit++} }
					targetTrials++
				}
				for _,k:=range bursts {
					v:=rng.intn(32);truth[k]=v
					for s:=0;s<4;s++ { was:=x.mem.find(k)>=0;a,_:=x.process(k,v);if !was&&a{burstAdmit++} }
					burstTrials++
				}
				if x.counter!=0 { fp+=x.fill(&next,rng) }

				for g:=0;g<skip;g++ {
					for i:=0;i<32;i++ { a,_:=x.process(next,rng.intn(32));next++;if a{fp++} }
					if x.counter!=0 { fp+=x.fill(&next,rng) }
				}
				for _,k:=range targets {
					for s:=0;s<2;s++ { was:=x.mem.find(k)>=0;a,_:=x.process(k,truth[k]);if !was&&a{targetAdmit++} }
				}
				if x.counter!=0 { fp+=x.fill(&next,rng) }

				for j:=0;j<churn;j++ {
					key:=1000+j;a,_:=x.process(key,rng.intn(32));if a{fp++}
					if (j+1)%4==0 { up117cQuery(&x.mem,targets,bursts) }
				}
				exact:=true
				for k:=0;k<12;k++ { got,ok:=x.mem.query(k);hotTotal++;setTotal++;if ok&&got==truth[k]{hotHits++;setHits++}else{exact=false} }
				for _,k:=range targets { got,ok:=x.mem.query(k);targetEval++;setTotal++;if ok&&got==truth[k]{targetHits++;setHits++}else{exact=false} }
				for _,k:=range bursts { got,ok:=x.mem.query(k);burstEval++;if ok&&got==truth[k]{burstHits++} }
				if exact{exactHits++};if x.mem.count>maxEntries{maxEntries=x.mem.count};checkpoints+=x.checkpoints
			}
			episodes++
		}
	}

	rate:=func(h,t int)float64{if t==0{return 0};return float64(h)/float64(t)}
	return UP120CPoint{
		Arm:arm,Workload:workload,SkippedGenerations:skip,
		TargetAdmissionRate:rate(targetAdmit,targetTrials),BurstAdmissionRate:rate(burstAdmit,burstTrials),
		HotAccuracy:rate(hotHits,hotTotal),TargetAccuracy:rate(targetHits,targetEval),BurstAccuracy:rate(burstHits,burstEval),
		Target16Accuracy:rate(setHits,setTotal),Target16ExactAccuracy:rate(exactHits,episodes),
		OneShotFalseAdmissions:fp,RecallEntriesUsed:maxEntries,MeanCheckpointsPerEpisode:float64(checkpoints)/float64(episodes),
		AdmissionMemoryBytes:128,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP120C()(UP120CAge3Result,error) {
	result:=UP120CAge3Result{
		Schema:UP120CAge3Schema,Experiment:"UP-120C-crossgeneration-age3",
		SourceUP119CSeal:"b49b8ecc045b71c2788fe3935bbbd3dc761576f4",
		HotKeys:12,TargetCandidates:4,BurstDistractors:4,GenerationInterval:32,
		AdmissionMemoryBytes:128,ExactRecallCap:16,KeyBits:14,ChurnWrites:12288,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,AdaptiveHorizonUsed:false,
	}
	seeds:=[]int{227000000,228000000}
	workloads:=[]struct{name string;skip int}{
		{"adjacent",0},{"skip_one",1},{"skip_two",2},{"skip_three_boundary",3},
	}
	for _,w:=range workloads {
		result.Points=append(result.Points,
			up120cRun("age1_control",w.name,w.skip,seeds),
			up120cRun("age2_extended",w.name,w.skip,seeds),
		)
	}
	return result,nil
}
