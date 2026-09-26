package unitary

const UP119CMixedPersistenceSchema = "wingless.up119c-mixed-persistence-horizon.v1"

type UP119CPoint struct {
	Arm                        string  `json:"arm"`
	Workload                   string  `json:"workload"`
	AdjacentAdmissionRate      float64 `json:"adjacent_admission_rate"`
	SkippedAdmissionRate       float64 `json:"skipped_admission_rate"`
	TargetAdmissionRate        float64 `json:"target_admission_rate"`
	BurstAdmissionRate         float64 `json:"burst_admission_rate"`
	HotAccuracy                float64 `json:"hot_accuracy"`
	AdjacentAccuracy           float64 `json:"adjacent_accuracy"`
	SkippedAccuracy            float64 `json:"skipped_accuracy"`
	TargetAccuracy             float64 `json:"target_accuracy"`
	BurstAccuracy              float64 `json:"burst_accuracy"`
	Target16Accuracy           float64 `json:"target16_accuracy"`
	Target16ExactAccuracy      float64 `json:"target16_exact_accuracy"`
	OneShotFalseAdmissions     int     `json:"one_shot_false_admissions"`
	RecallEntriesUsed          int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode  float64 `json:"mean_checkpoints_per_episode"`
	AdmissionMemoryBytes       int     `json:"admission_memory_bytes"`
	TotalBoundedMemoryBytes    int     `json:"total_bounded_memory_bytes"`
}

type UP119CMixedPersistenceResult struct {
	Schema                        string        `json:"schema"`
	Experiment                    string        `json:"experiment"`
	SourceUP118CSeal              string        `json:"source_up118c_seal"`
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
	Points                        []UP119CPoint `json:"points"`
}

type up119cMachine interface {
	process(key,value int)(bool,bool)
	fill(next *int,rng *sq0RNG) int
	memory() *up81cAging
	checkpointCount() int
	counterValue() int
}

type up119cControl struct{ x *up117cCrossGenMachine }
func (m *up119cControl) process(k,v int)(bool,bool){ return m.x.process(k,v) }
func (m *up119cControl) fill(next *int,rng *sq0RNG) int { return up118cFillControl(m.x,next,rng) }
func (m *up119cControl) memory()*up81cAging { return &m.x.mem }
func (m *up119cControl) checkpointCount() int { return m.x.checkpoints }
func (m *up119cControl) counterValue() int { return m.x.counter }

type up119cAge2 struct{ x *up118cAge2Machine }
func (m *up119cAge2) process(k,v int)(bool,bool){ return m.x.process(k,v) }
func (m *up119cAge2) fill(next *int,rng *sq0RNG) int { return up118cFillAge2(m.x,next,rng) }
func (m *up119cAge2) memory()*up81cAging { return &m.x.mem }
func (m *up119cAge2) checkpointCount() int { return m.x.checkpoints }
func (m *up119cAge2) counterValue() int { return m.x.counter }

func up119cQuery(mem *up81cAging,adjacent,skipped,bursts []int) {
	for k:=0;k<12;k++ { mem.query(k) }
	for _,k:=range adjacent { mem.query(k) }
	for _,k:=range skipped { mem.query(k) }
	for _,k:=range bursts { mem.query(k) }
}

func up119cFillGeneration(m up119cMachine,next *int,rng *sq0RNG)(fp int) {
	if m.counterValue()!=0 { fp+=m.fill(next,rng) }
	return
}

func up119cBlankGeneration(m up119cMachine,next *int,rng *sq0RNG)(fp int) {
	if m.counterValue()!=0 { fp+=m.fill(next,rng) }
	for i:=0;i<32;i++ {
		key:=*next
		(*next)++
		a,_:=m.process(key,rng.intn(32))
		if a { fp++ }
	}
	if m.counterValue()!=0 { fp+=m.fill(next,rng) }
	return
}

func up119cPresent(m up119cMachine,key,value,sightings int,admitCounter *int) {
	for i:=0;i<sightings;i++ {
		was:=m.memory().find(key)>=0
		a,_:=m.process(key,value)
		if !was&&a { (*admitCounter)++ }
	}
}

func up119cRun(arm,workload string,seeds []int) UP119CPoint {
	const churn=12288
	adjacent:=[]int{}
	skipped:=[]int{}
	if workload=="mixed_valid" {
		adjacent=[]int{100,101}
		skipped=[]int{102,103}
	}else{
		skipped=[]int{100,101,102,103}
	}
	bursts:=[]int{200,201,202,203}

	adjAdmit,skipAdmit,burstAdmit:=0,0,0
	adjTrials,skipTrials,burstTrials:=0,0,0
	hotHits,hotTotal:=0,0
	adjHits,adjTotal,skipHits,skipTotal:=0,0,0,0
	targetHits,targetTotal,burstHits,burstTotal:=0,0,0,0
	setHits,setTotal,exactHits:=0,0,0
	episodes,fp,maxEntries,checkpoints:=0,0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,5901+len(arm)*199+len(workload)*211,ep))
			var m up119cMachine
			if arm=="previous_generation_control" {
				m=&up119cControl{x:&up117cCrossGenMachine{}}
			}else{
				m=&up119cAge2{x:&up118cAge2Machine{}}
			}
			truth:=map[int]int{}

			for k:=0;k<32;k++ {
				v:=rng.intn(32);truth[k]=v
				a,_:=m.process(k,v)
				if k>=16&&a { fp++ }
			}
			for k:=0;k<12;k++ {
				m.process(k,truth[k])
				m.memory().query(k)
			}

			next:=300
			fp+=up119cFillGeneration(m,&next,rng)

			for _,k:=range adjacent {
				v:=rng.intn(32);truth[k]=v
				up119cPresent(m,k,v,2,&adjAdmit)
				adjTrials++
			}
			for _,k:=range skipped {
				v:=rng.intn(32);truth[k]=v
				up119cPresent(m,k,v,2,&skipAdmit)
				skipTrials++
			}
			for _,k:=range bursts {
				v:=rng.intn(32);truth[k]=v
				up119cPresent(m,k,v,4,&burstAdmit)
				burstTrials++
			}
			fp+=up119cFillGeneration(m,&next,rng)
			up119cQuery(m.memory(),adjacent,skipped,bursts)

			if workload=="mixed_valid" {
				for _,k:=range adjacent { up119cPresent(m,k,truth[k],2,&adjAdmit) }
				fp+=up119cFillGeneration(m,&next,rng)
				up119cQuery(m.memory(),adjacent,skipped,bursts)

				for _,k:=range skipped { up119cPresent(m,k,truth[k],2,&skipAdmit) }
				fp+=up119cFillGeneration(m,&next,rng)
			}else{
				fp+=up119cBlankGeneration(m,&next,rng)
				fp+=up119cBlankGeneration(m,&next,rng)
				for _,k:=range skipped { up119cPresent(m,k,truth[k],2,&skipAdmit) }
				fp+=up119cFillGeneration(m,&next,rng)
			}
			up119cQuery(m.memory(),adjacent,skipped,bursts)

			for j:=0;j<churn;j++ {
				key:=1000+j
				a,_:=m.process(key,rng.intn(32))
				if a { fp++ }
				if (j+1)%4==0 { up119cQuery(m.memory(),adjacent,skipped,bursts) }
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=m.memory().query(k)
				hotTotal++;targetTotal++;setTotal++
				if ok&&got==truth[k] { hotHits++;targetHits++;setHits++ } else { exact=false }
			}
			for _,k:=range adjacent {
				got,ok:=m.memory().query(k)
				adjTotal++;targetTotal++;setTotal++
				if ok&&got==truth[k] { adjHits++;targetHits++;setHits++ } else { exact=false }
			}
			for _,k:=range skipped {
				got,ok:=m.memory().query(k)
				skipTotal++;targetTotal++;setTotal++
				if ok&&got==truth[k] { skipHits++;targetHits++;setHits++ } else { exact=false }
			}
			for _,k:=range bursts {
				got,ok:=m.memory().query(k);burstTotal++
				if ok&&got==truth[k] { burstHits++ }
			}
			if exact { exactHits++ }
			if m.memory().count>maxEntries { maxEntries=m.memory().count }
			checkpoints+=m.checkpointCount()
			episodes++
		}
	}

	rate:=func(h,t int)float64{if t==0{return 0};return float64(h)/float64(t)}
	return UP119CPoint{
		Arm:arm,Workload:workload,
		AdjacentAdmissionRate:rate(adjAdmit,adjTrials),
		SkippedAdmissionRate:rate(skipAdmit,skipTrials),
		TargetAdmissionRate:rate(adjAdmit+skipAdmit,adjTrials+skipTrials),
		BurstAdmissionRate:rate(burstAdmit,burstTrials),
		HotAccuracy:rate(hotHits,hotTotal),
		AdjacentAccuracy:rate(adjHits,adjTotal),
		SkippedAccuracy:rate(skipHits,skipTotal),
		TargetAccuracy:rate(targetHits- hotHits,targetTotal-hotTotal),
		BurstAccuracy:rate(burstHits,burstTotal),
		Target16Accuracy:rate(setHits,setTotal),
		Target16ExactAccuracy:rate(exactHits,episodes),
		OneShotFalseAdmissions:fp,RecallEntriesUsed:maxEntries,
		MeanCheckpointsPerEpisode:float64(checkpoints)/float64(episodes),
		AdmissionMemoryBytes:128,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP119C()(UP119CMixedPersistenceResult,error) {
	result:=UP119CMixedPersistenceResult{
		Schema:UP119CMixedPersistenceSchema,Experiment:"UP-119C-mixed-persistence-horizon",
		SourceUP118CSeal:"c04c533451f752810dce1cc718ae6428ced42eb3",
		HotKeys:12,TargetCandidates:4,BurstDistractors:4,GenerationInterval:32,
		AdmissionMemoryBytes:128,ExactRecallCap:16,KeyBits:14,ChurnWrites:12288,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,AdaptiveHorizonUsed:false,
	}
	seeds:=[]int{225000000,226000000}
	for _,workload:=range []string{"mixed_valid","two_skip_boundary"} {
		result.Points=append(result.Points,
			up119cRun("previous_generation_control",workload,seeds),
			up119cRun("age2_history",workload,seeds),
		)
	}
	return result,nil
}
