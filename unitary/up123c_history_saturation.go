package unitary

const UP123CHistorySaturationSchema = "wingless.up123c-history-saturation.v1"

type UP123CPoint struct {
	SetBCount                   int     `json:"set_b_count"`
	ExpectedHistoryEntries      int     `json:"expected_history_entries"`
	ValidAdmissionRate          float64 `json:"valid_admission_rate"`
	ValidAccuracy               float64 `json:"valid_accuracy"`
	HotAccuracy                 float64 `json:"hot_accuracy"`
	Target16Accuracy            float64 `json:"target16_accuracy"`
	Target16ExactAccuracy       float64 `json:"target16_exact_accuracy"`
	NonpersistentAdmissionRate  float64 `json:"nonpersistent_admission_rate"`
	NonpersistentRetentionRate  float64 `json:"nonpersistent_retention_rate"`
	OneShotFalseAdmissions      int     `json:"one_shot_false_admissions"`
	MaxCurrentTableEntries      int     `json:"max_current_table_entries"`
	MaxHistoryTableEntries      int     `json:"max_history_table_entries"`
	HistoryReachedFull          bool    `json:"history_reached_full"`
	RecallEntriesUsed           int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode   float64 `json:"mean_checkpoints_per_episode"`
	AdmissionMemoryBytes        int     `json:"admission_memory_bytes"`
	TotalBoundedMemoryBytes     int     `json:"total_bounded_memory_bytes"`
}

type UP123CHistorySaturationResult struct {
	Schema                        string        `json:"schema"`
	Experiment                    string        `json:"experiment"`
	SourceUP122CSeal              string        `json:"source_up122c_seal"`
	HotKeys                       int           `json:"hot_keys"`
	ValidTargets                  int           `json:"valid_targets"`
	SetACount                     int           `json:"set_a_count"`
	GenerationInterval            int           `json:"generation_interval"`
	AdmissionMemoryBytes          int           `json:"admission_memory_bytes"`
	ExactRecallCap                int           `json:"exact_recall_cap"`
	HistoryEntries                int           `json:"history_entries"`
	KeyBits                       int           `json:"key_bits"`
	ChurnWrites                   int           `json:"churn_writes"`
	EpisodesPerSeed               int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission bool          `json:"query_evidence_used_for_admission"`
	SemanticPriorityUsed          bool          `json:"semantic_priority_used"`
	FutureOracleUsed              bool          `json:"future_oracle_used"`
	AdaptiveHorizonUsed           bool          `json:"adaptive_horizon_used"`
	Points                        []UP123CPoint `json:"points"`
}

func up123cRun(setBCount int,seeds []int) UP123CPoint {
	const churn=12288
	targets:=[]int{100,101,102,103}
	setA:=make([]int,16)
	for i:=range setA { setA[i]=200+i }

	validAdmit,validTrials:=0,0
	nonAdmit,nonTrials:=0,0
	hotHits,hotTotal,validHits,validTotal:=0,0,0,0
	nonHits,nonTotal,targetHits,targetTotal,exactHits:=0,0,0,0,0
	fp,maxCurrent,maxHistory,maxEntries,checkpoints,episodes:=0,0,0,0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,6401+setBCount*239,ep))
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

			// Generation 1: exactly 16 qualified nonpersistent set-A keys.
			for _,k:=range setA {
				v:=rng.intn(32);truth[k]=v
				up122cPresent(x,k,v,2,&nonAdmit,&localCurrent,&localHistory)
				nonTrials++
			}

			// Generation 2: four valid targets + P qualified nonpersistent set-B keys,
			// then unique one-shot writes to the exact 32-write boundary.
			for _,k:=range targets {
				v:=rng.intn(32);truth[k]=v
				up122cPresent(x,k,v,2,&validAdmit,&localCurrent,&localHistory)
				validTrials++
			}
			setB:=make([]int,setBCount)
			for i:=0;i<setBCount;i++ { setB[i]=220+i }
			for _,k:=range setB {
				v:=rng.intn(32);truth[k]=v
				up122cPresent(x,k,v,2,&nonAdmit,&localCurrent,&localHistory)
				nonTrials++
			}
			if x.counter!=0 { fp+=up122cFill(x,&next,rng,&localCurrent,&localHistory) }

			// Generation 3: only the valid targets recur.
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
				got,ok:=x.mem.query(k);nonTotal++
				if ok&&got==truth[k] { nonHits++ }
			}
			for _,k:=range setB {
				got,ok:=x.mem.query(k);nonTotal++
				if ok&&got==truth[k] { nonHits++ }
			}
			if exact { exactHits++ }
			if localCurrent>maxCurrent { maxCurrent=localCurrent }
			if localHistory>maxHistory { maxHistory=localHistory }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			checkpoints+=x.checkpoints
			episodes++
		}
	}

	return UP123CPoint{
		SetBCount:setBCount,ExpectedHistoryEntries:20+setBCount,
		ValidAdmissionRate:up121cRate(validAdmit,validTrials),
		ValidAccuracy:up121cRate(validHits,validTotal),
		HotAccuracy:up121cRate(hotHits,hotTotal),
		Target16Accuracy:up121cRate(targetHits,targetTotal),
		Target16ExactAccuracy:up121cRate(exactHits,episodes),
		NonpersistentAdmissionRate:up121cRate(nonAdmit,nonTrials),
		NonpersistentRetentionRate:up121cRate(nonHits,nonTotal),
		OneShotFalseAdmissions:fp,
		MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,
		HistoryReachedFull:maxHistory==32,
		RecallEntriesUsed:maxEntries,MeanCheckpointsPerEpisode:float64(checkpoints)/float64(episodes),
		AdmissionMemoryBytes:128,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP123C()(UP123CHistorySaturationResult,error) {
	result:=UP123CHistorySaturationResult{
		Schema:UP123CHistorySaturationSchema,Experiment:"UP-123C-history-saturation",
		SourceUP122CSeal:"6cac73d570c9f795bcb7879cc6d57d1f4f8211dd",
		HotKeys:12,ValidTargets:4,SetACount:16,GenerationInterval:32,
		AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,KeyBits:14,
		ChurnWrites:12288,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,AdaptiveHorizonUsed:false,
	}
	seeds:=[]int{233000000,234000000}
	for _,p:=range []int{4,8,12} { result.Points=append(result.Points,up123cRun(p,seeds)) }
	return result,nil
}
