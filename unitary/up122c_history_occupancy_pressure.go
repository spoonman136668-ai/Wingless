package unitary

const UP122CHistoryPressureSchema = "wingless.up122c-history-occupancy-pressure.v1"

type UP122CPoint struct {
	EphemeralPerGeneration     int     `json:"ephemeral_per_generation"`
	ValidAdmissionRate         float64 `json:"valid_admission_rate"`
	ValidAccuracy              float64 `json:"valid_accuracy"`
	HotAccuracy                float64 `json:"hot_accuracy"`
	Target16Accuracy           float64 `json:"target16_accuracy"`
	Target16ExactAccuracy      float64 `json:"target16_exact_accuracy"`
	EphemeralAdmissionRate     float64 `json:"ephemeral_admission_rate"`
	EphemeralRetentionRate     float64 `json:"ephemeral_retention_rate"`
	OneShotFalseAdmissions     int     `json:"one_shot_false_admissions"`
	MaxCurrentTableEntries     int     `json:"max_current_table_entries"`
	MaxHistoryTableEntries     int     `json:"max_history_table_entries"`
	RecallEntriesUsed          int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode  float64 `json:"mean_checkpoints_per_episode"`
	AdmissionMemoryBytes       int     `json:"admission_memory_bytes"`
	TotalBoundedMemoryBytes    int     `json:"total_bounded_memory_bytes"`
}

type UP122CHistoryPressureResult struct {
	Schema                         string        `json:"schema"`
	Experiment                     string        `json:"experiment"`
	SourceUP121CSeal               string        `json:"source_up121c_seal"`
	HotKeys                        int           `json:"hot_keys"`
	ValidTargets                   int           `json:"valid_targets"`
	GenerationInterval             int           `json:"generation_interval"`
	AdmissionMemoryBytes           int           `json:"admission_memory_bytes"`
	ExactRecallCap                 int           `json:"exact_recall_cap"`
	HistoryEntries                 int           `json:"history_entries"`
	KeyBits                        int           `json:"key_bits"`
	ChurnWrites                    int           `json:"churn_writes"`
	EpisodesPerSeed                int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission  bool          `json:"query_evidence_used_for_admission"`
	SemanticPriorityUsed           bool          `json:"semantic_priority_used"`
	FutureOracleUsed               bool          `json:"future_oracle_used"`
	AdaptiveHorizonUsed            bool          `json:"adaptive_horizon_used"`
	Points                         []UP122CPoint `json:"points"`
}

func up122cTableCount(table *[32]uint16) int {
	n:=0
	for _,v:=range table { if v!=0 { n++ } }
	return n
}

func up122cObserve(x *up120cExtendedMachine,maxCurrent,maxHistory *int) {
	c:=up122cTableCount(&x.current)
	h:=up122cTableCount(&x.history)
	if c>*maxCurrent { *maxCurrent=c }
	if h>*maxHistory { *maxHistory=h }
}

func up122cProcess(x *up120cExtendedMachine,key,value int,maxCurrent,maxHistory *int)(bool,bool) {
	a,r:=x.process(key,value)
	up122cObserve(x,maxCurrent,maxHistory)
	return a,r
}

func up122cFill(x *up120cExtendedMachine,next *int,rng *sq0RNG,maxCurrent,maxHistory *int)(fp int) {
	for x.counter!=0 {
		key:=*next
		(*next)++
		a,_:=up122cProcess(x,key,rng.intn(32),maxCurrent,maxHistory)
		if a { fp++ }
	}
	return
}

func up122cPresent(x *up120cExtendedMachine,key,value,sightings int,admit *int,maxCurrent,maxHistory *int) {
	for i:=0;i<sightings;i++ {
		was:=x.mem.find(key)>=0
		a,_:=up122cProcess(x,key,value,maxCurrent,maxHistory)
		if !was&&a { (*admit)++ }
	}
}

func up122cRun(pressure int,seeds []int) UP122CPoint {
	const churn=12288
	targets:=[]int{100,101,102,103}
	validAdmit,validTrials:=0,0
	ephAdmit,ephTrials:=0,0
	hotHits,hotTotal,validHits,validTotal:=0,0,0,0
	ephHits,ephTotal,targetHits,targetTotal,exactHits:=0,0,0,0,0
	fp,maxCurrent,maxHistory,maxEntries,checkpoints,episodes:=0,0,0,0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,6301+pressure*233,ep))
			x:=&up120cExtendedMachine{maxAge:2}
			truth:=map[int]int{}
			next:=300
			localCurrent,localHistory:=0,0

			for k:=0;k<32;k++ {
				v:=rng.intn(32);truth[k]=v
				a,_:=up122cProcess(x,k,v,&localCurrent,&localHistory)
				if k>=16&&a { fp++ }
			}
			for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }
			up122cObserve(x,&localCurrent,&localHistory)
			if x.counter!=0 { fp+=up122cFill(x,&next,rng,&localCurrent,&localHistory) }

			setA:=make([]int,pressure)
			setB:=make([]int,pressure)
			for i:=0;i<pressure;i++ { setA[i]=200+i;setB[i]=220+i }

			for _,k:=range targets {
				v:=rng.intn(32);truth[k]=v
				up122cPresent(x,k,v,2,&validAdmit,&localCurrent,&localHistory)
				validTrials++
			}
			for _,k:=range setA {
				v:=rng.intn(32);truth[k]=v
				up122cPresent(x,k,v,2,&ephAdmit,&localCurrent,&localHistory)
				ephTrials++
			}
			if x.counter!=0 { fp+=up122cFill(x,&next,rng,&localCurrent,&localHistory) }

			for _,k:=range setB {
				v:=rng.intn(32);truth[k]=v
				up122cPresent(x,k,v,2,&ephAdmit,&localCurrent,&localHistory)
				ephTrials++
			}
			if x.counter!=0 { fp+=up122cFill(x,&next,rng,&localCurrent,&localHistory) }

			for _,k:=range targets {
				up122cPresent(x,k,truth[k],2,&validAdmit,&localCurrent,&localHistory)
			}
			if x.counter!=0 { fp+=up122cFill(x,&next,rng,&localCurrent,&localHistory) }

			for j:=0;j<churn;j++ {
				key:=1000+j
				a,_:=up122cProcess(x,key,rng.intn(32),&localCurrent,&localHistory)
				if a { fp++ }
				if (j+1)%4==0 {
					for k:=0;k<12;k++ { x.mem.query(k) }
					for _,k:=range targets { x.mem.query(k) }
					for _,k:=range setA { x.mem.query(k) }
					for _,k:=range setB { x.mem.query(k) }
				}
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=x.mem.query(k);hotTotal++;targetTotal++
				if ok&&got==truth[k] { hotHits++;targetHits++ } else { exact=false }
			}
			for _,k:=range targets {
				got,ok:=x.mem.query(k);validTotal++;targetTotal++
				if ok&&got==truth[k] { validHits++;targetHits++ } else { exact=false }
			}
			for _,k:=range setA {
				got,ok:=x.mem.query(k);ephTotal++
				if ok&&got==truth[k] { ephHits++ }
			}
			for _,k:=range setB {
				got,ok:=x.mem.query(k);ephTotal++
				if ok&&got==truth[k] { ephHits++ }
			}
			if exact { exactHits++ }
			if localCurrent>maxCurrent { maxCurrent=localCurrent }
			if localHistory>maxHistory { maxHistory=localHistory }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			checkpoints+=x.checkpoints
			episodes++
		}
	}

	return UP122CPoint{
		EphemeralPerGeneration:pressure,
		ValidAdmissionRate:up121cRate(validAdmit,validTrials),
		ValidAccuracy:up121cRate(validHits,validTotal),
		HotAccuracy:up121cRate(hotHits,hotTotal),
		Target16Accuracy:up121cRate(targetHits,targetTotal),
		Target16ExactAccuracy:up121cRate(exactHits,episodes),
		EphemeralAdmissionRate:up121cRate(ephAdmit,ephTrials),
		EphemeralRetentionRate:up121cRate(ephHits,ephTotal),
		OneShotFalseAdmissions:fp,
		MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,
		RecallEntriesUsed:maxEntries,MeanCheckpointsPerEpisode:float64(checkpoints)/float64(episodes),
		AdmissionMemoryBytes:128,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP122C()(UP122CHistoryPressureResult,error) {
	result:=UP122CHistoryPressureResult{
		Schema:UP122CHistoryPressureSchema,Experiment:"UP-122C-history-occupancy-pressure",
		SourceUP121CSeal:"33516f1cc1a939e8859c74d1cef17c41fa6cf9a1",
		HotKeys:12,ValidTargets:4,GenerationInterval:32,AdmissionMemoryBytes:128,ExactRecallCap:16,
		HistoryEntries:32,KeyBits:14,ChurnWrites:12288,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,AdaptiveHorizonUsed:false,
	}
	seeds:=[]int{231000000,232000000}
	for _,pressure:=range []int{0,4,8,12} { result.Points=append(result.Points,up122cRun(pressure,seeds)) }
	return result,nil
}
