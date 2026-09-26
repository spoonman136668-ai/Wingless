package unitary

const UP113CMixedSetSchema = "wingless.up113c-exact-probation-mixedset.v1"

type UP113CPoint struct {
	RepeatGap                  int     `json:"repeat_gap"`
	ChurnWrites                int     `json:"churn_writes"`
	RecurringAdmissionRate     float64 `json:"recurring_admission_rate"`
	RecurringValueAccuracy     float64 `json:"recurring_value_accuracy"`
	HotQueryAccuracy           float64 `json:"hot_query_accuracy"`
	MixedSetExactAccuracy      float64 `json:"mixed_set_exact_accuracy"`
	OneShotFalseAdmissions     int     `json:"one_shot_false_admissions"`
	RecallEntriesUsed          int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode  float64 `json:"mean_checkpoints_per_episode"`
	PolicyMetadataBytes        int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes    int     `json:"total_bounded_memory_bytes"`
}

type UP113CMixedSetResult struct {
	Schema                     string        `json:"schema"`
	Experiment                 string        `json:"experiment"`
	SourceUP112CSeal           string        `json:"source_up112c_seal"`
	HotKeys                    int           `json:"hot_keys"`
	RecurringKeys              int           `json:"recurring_keys"`
	TargetWorkingSet           int           `json:"target_working_set"`
	GenerationInterval         int           `json:"generation_interval"`
	AdmissionMemoryBytes       int           `json:"admission_memory_bytes"`
	ExactRecallCap             int           `json:"exact_recall_cap"`
	EpisodesPerSeed            int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission bool       `json:"query_evidence_used_for_admission"`
	FutureOracleUsed           bool          `json:"future_oracle_used"`
	PhaseLabelUsed             bool          `json:"phase_label_used"`
	Points                     []UP113CPoint `json:"points"`
}

func up113cQueryWorkingSet(x *up111cExactMachine,recurring []int) {
	for k:=0;k<12;k++ { x.mem.query(k) }
	for _,k:=range recurring { x.mem.query(k) }
}

func up113cFillToGenerationBoundary(x *up111cExactMachine,next *int,rng *sq0RNG)(falseAdmissions int) {
	for x.counter!=0 {
		key:=*next
		*next++
		admitted,_:=x.process(key,rng.intn(32))
		if admitted { falseAdmissions++ }
	}
	return
}

func up113cRun(gap,churn int,seeds []int) UP113CPoint {
	admittedRecurring,totalRecurring:=0,0
	recHits,recTotal:=0,0
	hotHits,hotTotal:=0,0
	exactHits,episodes:=0,0
	falseAdmissions,maxEntries,totalCheckpoints:=0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,5101+gap*149+churn,ep))
			x:=&up111cExactMachine{}
			truth:=map[int]int{}

			// Deterministic original load. The first 16 become exact residents;
			// 16..31 are first-seen probation writes.
			for k:=0;k<32;k++ {
				v:=rng.intn(32)
				truth[k]=v
				admitted,_:=x.process(k,v)
				if k>=16 && admitted { falseAdmissions++ }
			}

			// Promote/protect the twelve persistent hot residents.
			for k:=0;k<12;k++ {
				x.process(k,truth[k])
				x.mem.query(k)
			}

			nextOneShot:=1000000+ep*200000
			falseAdmissions+=up113cFillToGenerationBoundary(x,&nextOneShot,rng)

			recurring:=[]int{100,101,102,103}
			for _,key:=range recurring {
				v:=rng.intn(32)
				truth[key]=v
				// Start each legitimate recurrence at a fresh generation.
				if x.counter!=0 { falseAdmissions+=up113cFillToGenerationBoundary(x,&nextOneShot,rng) }

				firstAdmit,_:=x.process(key,v)
				if firstAdmit { falseAdmissions++ }

				for i:=0;i<gap-1;i++ {
					k:=nextOneShot;nextOneShot++
					admitted,_:=x.process(k,rng.intn(32))
					if admitted { falseAdmissions++ }
				}

				admitted,_:=x.process(key,v)
				totalRecurring++
				if admitted {
					admittedRecurring++
					x.mem.query(key)
				}

				// Complete this generation before starting the next candidate.
				if x.counter!=0 { falseAdmissions+=up113cFillToGenerationBoundary(x,&nextOneShot,rng) }
				up113cQueryWorkingSet(x,recurring)
			}

			// Long unique one-shot churn.
			for j:=0;j<churn;j++ {
				key:=2000000+ep*200000+j
				admitted,_:=x.process(key,rng.intn(32))
				if admitted { falseAdmissions++ }
				if (j+1)%4==0 { up113cQueryWorkingSet(x,recurring) }
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=x.mem.query(k)
				hotTotal++
				if ok&&got==truth[k] { hotHits++ } else { exact=false }
			}
			for _,key:=range recurring {
				got,ok:=x.mem.query(key)
				recTotal++
				if ok&&got==truth[key] { recHits++ } else { exact=false }
			}
			if exact { exactHits++ }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			totalCheckpoints+=x.checkpoints
			episodes++
		}
	}

	return UP113CPoint{
		RepeatGap:gap,ChurnWrites:churn,
		RecurringAdmissionRate:float64(admittedRecurring)/float64(totalRecurring),
		RecurringValueAccuracy:float64(recHits)/float64(recTotal),
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		MixedSetExactAccuracy:float64(exactHits)/float64(episodes),
		OneShotFalseAdmissions:falseAdmissions,
		RecallEntriesUsed:maxEntries,
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		PolicyMetadataBytes:136,
		TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP113C()(UP113CMixedSetResult,error){
	result:=UP113CMixedSetResult{
		Schema:UP113CMixedSetSchema,
		Experiment:"UP-113C-exact-probation-mixedset",
		SourceUP112CSeal:"11d807704fc0eab679530497bfceba1d0e0f2c86",
		HotKeys:12,RecurringKeys:4,TargetWorkingSet:16,
		GenerationInterval:32,AdmissionMemoryBytes:128,ExactRecallCap:16,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,FutureOracleUsed:false,PhaseLabelUsed:false,
	}
	seeds:=[]int{213000000,214000000}
	for _,gap:=range []int{8,16,24,31} {
		for _,churn:=range []int{24576,49152,98304} {
			result.Points=append(result.Points,up113cRun(gap,churn,seeds))
		}
	}
	return result,nil
}
