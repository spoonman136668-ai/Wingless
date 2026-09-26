package unitary

const UP121CMixedSharingSchema = "wingless.up121c-mixed-age3-sharing.v1"

type UP121CPoint struct {
	Arm                         string  `json:"arm"`
	AdjacentAdmissionRate       float64 `json:"adjacent_admission_rate"`
	SkipOneAdmissionRate        float64 `json:"skip_one_admission_rate"`
	SkipTwoAdmissionRate        float64 `json:"skip_two_admission_rate"`
	BurstAdmissionRate          float64 `json:"burst_admission_rate"`
	SkipThreeAdmissionRate      float64 `json:"skip_three_admission_rate"`
	HotAccuracy                 float64 `json:"hot_accuracy"`
	AdjacentAccuracy            float64 `json:"adjacent_accuracy"`
	SkipOneAccuracy             float64 `json:"skip_one_accuracy"`
	SkipTwoAccuracy             float64 `json:"skip_two_accuracy"`
	BurstAccuracy               float64 `json:"burst_accuracy"`
	SkipThreeAccuracy           float64 `json:"skip_three_accuracy"`
	Target16Accuracy            float64 `json:"target16_accuracy"`
	Target16ExactAccuracy       float64 `json:"target16_exact_accuracy"`
	NegativeRetentionRate       float64 `json:"negative_retention_rate"`
	OneShotFalseAdmissions      int     `json:"one_shot_false_admissions"`
	RecallEntriesUsed           int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode   float64 `json:"mean_checkpoints_per_episode"`
	AdmissionMemoryBytes        int     `json:"admission_memory_bytes"`
	TotalBoundedMemoryBytes     int     `json:"total_bounded_memory_bytes"`
}

type UP121CMixedSharingResult struct {
	Schema                         string        `json:"schema"`
	Experiment                     string        `json:"experiment"`
	SourceUP120CSeal               string        `json:"source_up120c_seal"`
	HotKeys                        int           `json:"hot_keys"`
	ValidRecurringCandidates       int           `json:"valid_recurring_candidates"`
	BurstDistractors               int           `json:"burst_distractors"`
	SkipThreeNegatives             int           `json:"skip_three_negatives"`
	GenerationInterval             int           `json:"generation_interval"`
	AdmissionMemoryBytes           int           `json:"admission_memory_bytes"`
	ExactRecallCap                 int           `json:"exact_recall_cap"`
	KeyBits                        int           `json:"key_bits"`
	ChurnWrites                    int           `json:"churn_writes"`
	EpisodesPerSeed                int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission  bool          `json:"query_evidence_used_for_admission"`
	SemanticPriorityUsed           bool          `json:"semantic_priority_used"`
	FutureOracleUsed               bool          `json:"future_oracle_used"`
	AdaptiveHorizonUsed            bool          `json:"adaptive_horizon_used"`
	Points                         []UP121CPoint `json:"points"`
}

func up121cRate(h,t int) float64 {
	if t==0 { return 0 }
	return float64(h)/float64(t)
}

func up121cQuery(mem *up81cAging,groups ...[]int) {
	for k:=0;k<12;k++ { mem.query(k) }
	for _,g:=range groups {
		for _,k:=range g { mem.query(k) }
	}
}

func up121cPresentAge1(x *up118cAge2Machine,key,value,sightings int,admit *int) {
	for i:=0;i<sightings;i++ {
		was:=x.mem.find(key)>=0
		a,_:=x.process(key,value)
		if !was&&a { (*admit)++ }
	}
}

func up121cPresentAge2(x *up120cExtendedMachine,key,value,sightings int,admit *int) {
	for i:=0;i<sightings;i++ {
		was:=x.mem.find(key)>=0
		a,_:=x.process(key,value)
		if !was&&a { (*admit)++ }
	}
}

func up121cRunAge1(seeds []int) UP121CPoint {
	const churn=12288
	adj:=[]int{100,101}
	skip1:=[]int{102}
	skip2:=[]int{103}
	bursts:=[]int{200,201}
	skip3:=[]int{202,203}

	adjAdmit,skip1Admit,skip2Admit,burstAdmit,skip3Admit:=0,0,0,0,0
	adjTrials,skip1Trials,skip2Trials,burstTrials,skip3Trials:=0,0,0,0,0
	hotHits,hotTotal:=0,0
	adjHits,adjTotal,skip1Hits,skip1Total,skip2Hits,skip2Total:=0,0,0,0,0,0
	burstHits,burstTotal,skip3Hits,skip3Total:=0,0,0,0
	targetHits,targetTotal,exactHits,episodes:=0,0,0,0
	fp,maxEntries,checkpoints:=0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,6101,ep))
			x:=&up118cAge2Machine{}
			truth:=map[int]int{}
			next:=300

			for k:=0;k<32;k++ {
				v:=rng.intn(32);truth[k]=v
				a,_:=x.process(k,v)
				if k>=16&&a { fp++ }
			}
			for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }
			if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

			for _,k:=range adj {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge1(x,k,v,2,&adjAdmit);adjTrials++
			}
			for _,k:=range skip1 {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge1(x,k,v,2,&skip1Admit);skip1Trials++
			}
			for _,k:=range skip2 {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge1(x,k,v,2,&skip2Admit);skip2Trials++
			}
			for _,k:=range bursts {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge1(x,k,v,4,&burstAdmit);burstTrials++
			}
			for _,k:=range skip3 {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge1(x,k,v,2,&skip3Admit);skip3Trials++
			}
			if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

			for _,k:=range adj { up121cPresentAge1(x,k,truth[k],2,&adjAdmit) }
			if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

			for _,k:=range skip1 { up121cPresentAge1(x,k,truth[k],2,&skip1Admit) }
			if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

			for _,k:=range skip2 { up121cPresentAge1(x,k,truth[k],2,&skip2Admit) }
			if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

			for _,k:=range skip3 { up121cPresentAge1(x,k,truth[k],2,&skip3Admit) }
			if x.counter!=0 { fp+=up118cFillAge2(x,&next,rng) }

			for j:=0;j<churn;j++ {
				key:=1000+j
				a,_:=x.process(key,rng.intn(32))
				if a { fp++ }
				if (j+1)%4==0 { up121cQuery(&x.mem,adj,skip1,skip2,bursts,skip3) }
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=x.mem.query(k);hotTotal++;targetTotal++
				if ok&&got==truth[k] { hotHits++;targetHits++ } else { exact=false }
			}
			for _,k:=range adj {
				got,ok:=x.mem.query(k);adjTotal++;targetTotal++
				if ok&&got==truth[k] { adjHits++;targetHits++ } else { exact=false }
			}
			for _,k:=range skip1 {
				got,ok:=x.mem.query(k);skip1Total++;targetTotal++
				if ok&&got==truth[k] { skip1Hits++;targetHits++ } else { exact=false }
			}
			for _,k:=range skip2 {
				got,ok:=x.mem.query(k);skip2Total++;targetTotal++
				if ok&&got==truth[k] { skip2Hits++;targetHits++ } else { exact=false }
			}
			for _,k:=range bursts {
				got,ok:=x.mem.query(k);burstTotal++
				if ok&&got==truth[k] { burstHits++ }
			}
			for _,k:=range skip3 {
				got,ok:=x.mem.query(k);skip3Total++
				if ok&&got==truth[k] { skip3Hits++ }
			}
			if exact { exactHits++ }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			checkpoints+=x.checkpoints
			episodes++
		}
	}
	negHits:=burstHits+skip3Hits
	negTotal:=burstTotal+skip3Total
	return UP121CPoint{
		Arm:"age1_control",
		AdjacentAdmissionRate:up121cRate(adjAdmit,adjTrials),
		SkipOneAdmissionRate:up121cRate(skip1Admit,skip1Trials),
		SkipTwoAdmissionRate:up121cRate(skip2Admit,skip2Trials),
		BurstAdmissionRate:up121cRate(burstAdmit,burstTrials),
		SkipThreeAdmissionRate:up121cRate(skip3Admit,skip3Trials),
		HotAccuracy:up121cRate(hotHits,hotTotal),AdjacentAccuracy:up121cRate(adjHits,adjTotal),
		SkipOneAccuracy:up121cRate(skip1Hits,skip1Total),SkipTwoAccuracy:up121cRate(skip2Hits,skip2Total),
		BurstAccuracy:up121cRate(burstHits,burstTotal),SkipThreeAccuracy:up121cRate(skip3Hits,skip3Total),
		Target16Accuracy:up121cRate(targetHits,targetTotal),Target16ExactAccuracy:up121cRate(exactHits,episodes),
		NegativeRetentionRate:up121cRate(negHits,negTotal),OneShotFalseAdmissions:fp,
		RecallEntriesUsed:maxEntries,MeanCheckpointsPerEpisode:float64(checkpoints)/float64(episodes),
		AdmissionMemoryBytes:128,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func up121cRunAge2(seeds []int) UP121CPoint {
	const churn=12288
	adj:=[]int{100,101}
	skip1:=[]int{102}
	skip2:=[]int{103}
	bursts:=[]int{200,201}
	skip3:=[]int{202,203}

	adjAdmit,skip1Admit,skip2Admit,burstAdmit,skip3Admit:=0,0,0,0,0
	adjTrials,skip1Trials,skip2Trials,burstTrials,skip3Trials:=0,0,0,0,0
	hotHits,hotTotal:=0,0
	adjHits,adjTotal,skip1Hits,skip1Total,skip2Hits,skip2Total:=0,0,0,0,0,0
	burstHits,burstTotal,skip3Hits,skip3Total:=0,0,0,0
	targetHits,targetTotal,exactHits,episodes:=0,0,0,0
	fp,maxEntries,checkpoints:=0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,6201,ep))
			x:=&up120cExtendedMachine{maxAge:2}
			truth:=map[int]int{}
			next:=300

			for k:=0;k<32;k++ {
				v:=rng.intn(32);truth[k]=v
				a,_:=x.process(k,v)
				if k>=16&&a { fp++ }
			}
			for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }
			if x.counter!=0 { fp+=x.fill(&next,rng) }

			for _,k:=range adj {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge2(x,k,v,2,&adjAdmit);adjTrials++
			}
			for _,k:=range skip1 {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge2(x,k,v,2,&skip1Admit);skip1Trials++
			}
			for _,k:=range skip2 {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge2(x,k,v,2,&skip2Admit);skip2Trials++
			}
			for _,k:=range bursts {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge2(x,k,v,4,&burstAdmit);burstTrials++
			}
			for _,k:=range skip3 {
				v:=rng.intn(32);truth[k]=v;up121cPresentAge2(x,k,v,2,&skip3Admit);skip3Trials++
			}
			if x.counter!=0 { fp+=x.fill(&next,rng) }

			for _,k:=range adj { up121cPresentAge2(x,k,truth[k],2,&adjAdmit) }
			if x.counter!=0 { fp+=x.fill(&next,rng) }

			for _,k:=range skip1 { up121cPresentAge2(x,k,truth[k],2,&skip1Admit) }
			if x.counter!=0 { fp+=x.fill(&next,rng) }

			for _,k:=range skip2 { up121cPresentAge2(x,k,truth[k],2,&skip2Admit) }
			if x.counter!=0 { fp+=x.fill(&next,rng) }

			for _,k:=range skip3 { up121cPresentAge2(x,k,truth[k],2,&skip3Admit) }
			if x.counter!=0 { fp+=x.fill(&next,rng) }

			for j:=0;j<churn;j++ {
				key:=1000+j
				a,_:=x.process(key,rng.intn(32))
				if a { fp++ }
				if (j+1)%4==0 { up121cQuery(&x.mem,adj,skip1,skip2,bursts,skip3) }
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=x.mem.query(k);hotTotal++;targetTotal++
				if ok&&got==truth[k] { hotHits++;targetHits++ } else { exact=false }
			}
			for _,k:=range adj {
				got,ok:=x.mem.query(k);adjTotal++;targetTotal++
				if ok&&got==truth[k] { adjHits++;targetHits++ } else { exact=false }
			}
			for _,k:=range skip1 {
				got,ok:=x.mem.query(k);skip1Total++;targetTotal++
				if ok&&got==truth[k] { skip1Hits++;targetHits++ } else { exact=false }
			}
			for _,k:=range skip2 {
				got,ok:=x.mem.query(k);skip2Total++;targetTotal++
				if ok&&got==truth[k] { skip2Hits++;targetHits++ } else { exact=false }
			}
			for _,k:=range bursts {
				got,ok:=x.mem.query(k);burstTotal++
				if ok&&got==truth[k] { burstHits++ }
			}
			for _,k:=range skip3 {
				got,ok:=x.mem.query(k);skip3Total++
				if ok&&got==truth[k] { skip3Hits++ }
			}
			if exact { exactHits++ }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			checkpoints+=x.checkpoints
			episodes++
		}
	}
	negHits:=burstHits+skip3Hits
	negTotal:=burstTotal+skip3Total
	return UP121CPoint{
		Arm:"age2_extended",
		AdjacentAdmissionRate:up121cRate(adjAdmit,adjTrials),
		SkipOneAdmissionRate:up121cRate(skip1Admit,skip1Trials),
		SkipTwoAdmissionRate:up121cRate(skip2Admit,skip2Trials),
		BurstAdmissionRate:up121cRate(burstAdmit,burstTrials),
		SkipThreeAdmissionRate:up121cRate(skip3Admit,skip3Trials),
		HotAccuracy:up121cRate(hotHits,hotTotal),AdjacentAccuracy:up121cRate(adjHits,adjTotal),
		SkipOneAccuracy:up121cRate(skip1Hits,skip1Total),SkipTwoAccuracy:up121cRate(skip2Hits,skip2Total),
		BurstAccuracy:up121cRate(burstHits,burstTotal),SkipThreeAccuracy:up121cRate(skip3Hits,skip3Total),
		Target16Accuracy:up121cRate(targetHits,targetTotal),Target16ExactAccuracy:up121cRate(exactHits,episodes),
		NegativeRetentionRate:up121cRate(negHits,negTotal),OneShotFalseAdmissions:fp,
		RecallEntriesUsed:maxEntries,MeanCheckpointsPerEpisode:float64(checkpoints)/float64(episodes),
		AdmissionMemoryBytes:128,TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP121C()(UP121CMixedSharingResult,error) {
	result:=UP121CMixedSharingResult{
		Schema:UP121CMixedSharingSchema,Experiment:"UP-121C-mixed-age3-sharing",
		SourceUP120CSeal:"d457b42196b2ee87cb701bed1da63f759701e840",
		HotKeys:12,ValidRecurringCandidates:4,BurstDistractors:2,SkipThreeNegatives:2,
		GenerationInterval:32,AdmissionMemoryBytes:128,ExactRecallCap:16,KeyBits:14,
		ChurnWrites:12288,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,AdaptiveHorizonUsed:false,
	}
	seeds:=[]int{229000000,230000000}
	result.Points=append(result.Points,up121cRunAge1(seeds),up121cRunAge2(seeds))
	return result,nil
}
