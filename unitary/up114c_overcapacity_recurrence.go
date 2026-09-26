package unitary

const UP114COvercapacitySchema = "wingless.up114c-overcapacity-recurrence.v1"

type UP114CPoint struct {
	RecurringDemand            int     `json:"recurring_demand"`
	TotalTargetDemand          int     `json:"total_target_demand"`
	ChurnWrites                int     `json:"churn_writes"`
	RecurringAdmissionRate     float64 `json:"recurring_admission_rate"`
	HotQueryAccuracy           float64 `json:"hot_query_accuracy"`
	RecurringAccuracy          float64 `json:"recurring_accuracy"`
	MeanRetainedRecurringCount float64 `json:"mean_retained_recurring_count"`
	TargetSetAccuracy          float64 `json:"target_set_accuracy"`
	WholeTargetExactAccuracy   float64 `json:"whole_target_exact_accuracy"`
	OneShotFalseAdmissions     int     `json:"one_shot_false_admissions"`
	RecallEntriesUsed          int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode  float64 `json:"mean_checkpoints_per_episode"`
	PolicyMetadataBytes        int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes    int     `json:"total_bounded_memory_bytes"`
}

type UP114COvercapacityResult struct {
	Schema                     string        `json:"schema"`
	Experiment                 string        `json:"experiment"`
	SourceUP113CSeal           string        `json:"source_up113c_seal"`
	HotKeys                    int           `json:"hot_keys"`
	GenerationInterval         int           `json:"generation_interval"`
	RepeatGap                  int           `json:"repeat_gap"`
	AdmissionMemoryBytes       int           `json:"admission_memory_bytes"`
	ExactRecallCap             int           `json:"exact_recall_cap"`
	EpisodesPerSeed            int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission bool       `json:"query_evidence_used_for_admission"`
	PriorityLabelsUsed         bool          `json:"priority_labels_used"`
	FutureOracleUsed           bool          `json:"future_oracle_used"`
	PhaseLabelUsed             bool          `json:"phase_label_used"`
	Points                     []UP114CPoint `json:"points"`
}

func up114cQueryTargets(x *up111cExactMachine, recurring []int) {
	for k:=0;k<12;k++ { x.mem.query(k) }
	for _,k:=range recurring { x.mem.query(k) }
}

func up114cRun(recurringCount int,seeds []int) UP114CPoint {
	const churn=49152
	const gap=16
	admittedRecurring,totalRecurring:=0,0
	hotHits,hotTotal:=0,0
	recHits,recTotal:=0,0
	targetHits,targetTotal:=0,0
	retainedRecurringSum:=0
	exactHits,episodes:=0,0
	falseAdmissions,maxEntries,totalCheckpoints:=0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,5201+recurringCount*157,ep))
			x:=&up111cExactMachine{}
			truth:=map[int]int{}

			for k:=0;k<32;k++ {
				v:=rng.intn(32)
				truth[k]=v
				admitted,_:=x.process(k,v)
				if k>=16&&admitted{falseAdmissions++}
			}

			for k:=0;k<12;k++ {
				x.process(k,truth[k])
				x.mem.query(k)
			}

			nextOneShot:=1000000+ep*300000
			if x.counter!=0 { falseAdmissions+=up113cFillToGenerationBoundary(x,&nextOneShot,rng) }

			recurring:=make([]int,recurringCount)
			for i:=0;i<recurringCount;i++ { recurring[i]=100+i }

			for _,key:=range recurring {
				v:=rng.intn(32)
				truth[key]=v
				if x.counter!=0 { falseAdmissions+=up113cFillToGenerationBoundary(x,&nextOneShot,rng) }

				firstAdmit,_:=x.process(key,v)
				if firstAdmit { falseAdmissions++ }

				for i:=0;i<gap-1;i++ {
					k:=nextOneShot
					nextOneShot++
					admitted,_:=x.process(k,rng.intn(32))
					if admitted { falseAdmissions++ }
				}

				admitted,_:=x.process(key,v)
				totalRecurring++
				if admitted { admittedRecurring++ }
				if x.counter!=0 { falseAdmissions+=up113cFillToGenerationBoundary(x,&nextOneShot,rng) }
				up114cQueryTargets(x,recurring)
			}

			for j:=0;j<churn;j++ {
				key:=2000000+ep*300000+j
				admitted,_:=x.process(key,rng.intn(32))
				if admitted { falseAdmissions++ }
				if (j+1)%4==0 { up114cQueryTargets(x,recurring) }
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=x.mem.query(k)
				hotTotal++;targetTotal++
				if ok&&got==truth[k] { hotHits++;targetHits++ } else { exact=false }
			}
			retained:=0
			for _,key:=range recurring {
				got,ok:=x.mem.query(key)
				recTotal++;targetTotal++
				if ok&&got==truth[key] {
					recHits++;targetHits++;retained++
				} else {
					exact=false
				}
			}
			retainedRecurringSum+=retained
			if exact { exactHits++ }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			totalCheckpoints+=x.checkpoints
			episodes++
		}
	}

	return UP114CPoint{
		RecurringDemand:recurringCount,TotalTargetDemand:12+recurringCount,ChurnWrites:churn,
		RecurringAdmissionRate:float64(admittedRecurring)/float64(totalRecurring),
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		RecurringAccuracy:float64(recHits)/float64(recTotal),
		MeanRetainedRecurringCount:float64(retainedRecurringSum)/float64(episodes),
		TargetSetAccuracy:float64(targetHits)/float64(targetTotal),
		WholeTargetExactAccuracy:float64(exactHits)/float64(episodes),
		OneShotFalseAdmissions:falseAdmissions,RecallEntriesUsed:maxEntries,
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		PolicyMetadataBytes:136,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP114C()(UP114COvercapacityResult,error){
	result:=UP114COvercapacityResult{
		Schema:UP114COvercapacitySchema,Experiment:"UP-114C-overcapacity-recurrence",
		SourceUP113CSeal:"f425d2410d413ee53086cf794efc397e66b3660f",
		HotKeys:12,GenerationInterval:32,RepeatGap:16,AdmissionMemoryBytes:128,ExactRecallCap:16,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,PriorityLabelsUsed:false,FutureOracleUsed:false,PhaseLabelUsed:false,
	}
	seeds:=[]int{215000000,216000000}
	for _,n:=range []int{4,5,8} { result.Points=append(result.Points,up114cRun(n,seeds)) }
	return result,nil
}
