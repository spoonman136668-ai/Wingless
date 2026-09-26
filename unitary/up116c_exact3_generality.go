package unitary

const UP116CExact3GeneralitySchema = "wingless.up116c-exact3-generality.v1"

type UP116CPoint struct {
	Scenario                    string  `json:"scenario"`
	WeakCandidates              int     `json:"weak_candidates"`
	StrongCandidates            int     `json:"strong_candidates"`
	DistractorCandidates        int     `json:"distractor_candidates"`
	WeakAdmissionRate           float64 `json:"weak_admission_rate"`
	StrongAdmissionRate         float64 `json:"strong_admission_rate"`
	DistractorAdmissionRate     float64 `json:"distractor_admission_rate"`
	HotAccuracy                 float64 `json:"hot_accuracy"`
	WeakAccuracy                float64 `json:"weak_accuracy"`
	StrongAccuracy              float64 `json:"strong_accuracy"`
	DistractorAccuracy          float64 `json:"distractor_accuracy"`
	Target16Accuracy            float64 `json:"target16_accuracy"`
	Target16ExactAccuracy       float64 `json:"target16_exact_accuracy"`
	AllCandidateAccuracy        float64 `json:"all_candidate_accuracy"`
	OneShotFalseAdmissions      int     `json:"one_shot_false_admissions"`
	RecallEntriesUsed           int     `json:"recall_entries_used"`
	MeanCheckpointsPerEpisode   float64 `json:"mean_checkpoints_per_episode"`
	AdmissionMemoryBytes        int     `json:"admission_memory_bytes"`
	TotalBoundedMemoryBytes     int     `json:"total_bounded_memory_bytes"`
}

type UP116CExact3GeneralityResult struct {
	Schema                     string        `json:"schema"`
	Experiment                 string        `json:"experiment"`
	SourceUP115CSeal           string        `json:"source_up115c_seal"`
	HotKeys                    int           `json:"hot_keys"`
	GenerationInterval         int           `json:"generation_interval"`
	AdmissionSightings         int           `json:"admission_sightings"`
	AdmissionMemoryBytes       int           `json:"admission_memory_bytes"`
	ExactRecallCap             int           `json:"exact_recall_cap"`
	EpisodesPerSeed            int           `json:"episodes_per_seed"`
	QueryEvidenceUsedForAdmission bool       `json:"query_evidence_used_for_admission"`
	SemanticPriorityUsed       bool          `json:"semantic_priority_used"`
	FutureOracleUsed           bool          `json:"future_oracle_used"`
	PhaseLabelUsed             bool          `json:"phase_label_used"`
	Points                     []UP116CPoint `json:"points"`
}

func up116cQuery(x *up115cStrengthMachine,weak,strong,distractors []int) {
	for k:=0;k<12;k++ { x.mem.query(k) }
	for _,k:=range weak { x.mem.query(k) }
	for _,k:=range strong { x.mem.query(k) }
	for _,k:=range distractors { x.mem.query(k) }
}

func up116cPresentGaps(x *up115cStrengthMachine,key,value int,gaps []int,next *int,rng *sq0RNG,weak,strong,distractors []int)(admitted bool,falseAdmissions int) {
	if x.counter!=0 { falseAdmissions+=up115cFillBoundary(x,next,rng) }

	first,_:=x.process(key,value)
	if first { falseAdmissions++ }

	for _,gap:=range gaps {
		for i:=0;i<gap-1;i++ {
			k:=*next
			(*next)++
			a,_:=x.process(k,rng.intn(32))
			if a { falseAdmissions++ }
			if i%4==3 { up116cQuery(x,weak,strong,distractors) }
		}
		a,_:=x.process(key,value)
		if a { admitted=true }
	}

	if x.counter!=0 { falseAdmissions+=up115cFillBoundary(x,next,rng) }
	up116cQuery(x,weak,strong,distractors)
	return
}

func up116cRun(scenario string,seeds []int) UP116CPoint {
	const churn=49152

	weak:=[]int{}
	strong:=[]int{200,201,202,203}
	distractors:=[]int{}

	if scenario=="mixed_gap_clean" {
		weak=[]int{100,101,102,103}
	}
	if scenario=="four_hit_distractor" {
		distractors=[]int{300,301,302,303}
	}

	weakAdmit,strongAdmit,distractorAdmit:=0,0,0
	weakTotal,strongTotal,distractorTotal:=0,0,0

	hotHits,hotTotal:=0,0
	weakHits,weakEval:=0,0
	strongHits,strongEval:=0,0
	distractorHits,distractorEval:=0,0
	targetHits,targetTotal:=0,0
	candidateHits,candidateTotal:=0,0
	targetExact,episodes:=0,0
	falseAdmissions,maxEntries,totalCheckpoints:=0,0,0

	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			rng:=newSQ0RNG(sq0Seed(base,5401+len(scenario)*173,ep))
			x:=&up115cStrengthMachine{threshold:3}
			truth:=map[int]int{}

			for k:=0;k<32;k++ {
				v:=rng.intn(32)
				truth[k]=v
				a,_:=x.process(k,v)
				if k>=16&&a { falseAdmissions++ }
			}
			for k:=0;k<12;k++ {
				x.process(k,truth[k])
				x.mem.query(k)
			}

			next:=1000000+ep*400000
			if x.counter!=0 { falseAdmissions+=up115cFillBoundary(x,&next,rng) }

			if scenario=="mixed_gap_clean" {
				weakGaps:=[]int{4,12,20,28}
				for i,key:=range weak {
					v:=rng.intn(32)
					truth[key]=v
					a,fp:=up116cPresentGaps(x,key,v,[]int{weakGaps[i]},&next,rng,weak,strong,distractors)
					falseAdmissions+=fp
					weakTotal++
					if a { weakAdmit++ }
				}
				firstGaps:=[]int{5,9,13,17}
				secondGaps:=[]int{6,10,14,12}
				for i,key:=range strong {
					v:=rng.intn(32)
					truth[key]=v
					a,fp:=up116cPresentGaps(x,key,v,[]int{firstGaps[i],secondGaps[i]},&next,rng,weak,strong,distractors)
					falseAdmissions+=fp
					strongTotal++
					if a { strongAdmit++ }
				}
			}else{
				for _,key:=range strong {
					v:=rng.intn(32)
					truth[key]=v
					a,fp:=up116cPresentGaps(x,key,v,[]int{8,8},&next,rng,weak,strong,distractors)
					falseAdmissions+=fp
					strongTotal++
					if a { strongAdmit++ }
				}
				for _,key:=range distractors {
					v:=rng.intn(32)
					truth[key]=v
					a,fp:=up116cPresentGaps(x,key,v,[]int{4,4,4},&next,rng,weak,strong,distractors)
					falseAdmissions+=fp
					distractorTotal++
					if a { distractorAdmit++ }
				}
			}

			for j:=0;j<churn;j++ {
				key:=2000000+ep*400000+j
				a,_:=x.process(key,rng.intn(32))
				if a { falseAdmissions++ }
				if (j+1)%4==0 { up116cQuery(x,weak,strong,distractors) }
			}

			exact:=true
			for k:=0;k<12;k++ {
				got,ok:=x.mem.query(k)
				hotTotal++;targetTotal++
				if ok&&got==truth[k] {
					hotHits++;targetHits++
				}else{
					exact=false
				}
			}
			for _,key:=range weak {
				got,ok:=x.mem.query(key)
				weakEval++;candidateTotal++
				if ok&&got==truth[key] { weakHits++;candidateHits++ }
			}
			for _,key:=range strong {
				got,ok:=x.mem.query(key)
				strongEval++;candidateTotal++;targetTotal++
				if ok&&got==truth[key] {
					strongHits++;candidateHits++;targetHits++
				}else{
					exact=false
				}
			}
			for _,key:=range distractors {
				got,ok:=x.mem.query(key)
				distractorEval++;candidateTotal++
				if ok&&got==truth[key] { distractorHits++;candidateHits++ }
			}

			if exact { targetExact++ }
			if x.mem.count>maxEntries { maxEntries=x.mem.count }
			totalCheckpoints+=x.checkpoints
			episodes++
		}
	}

	rate:=func(hits,total int)float64{
		if total==0 { return 0 }
		return float64(hits)/float64(total)
	}

	return UP116CPoint{
		Scenario:scenario,
		WeakCandidates:len(weak),StrongCandidates:len(strong),DistractorCandidates:len(distractors),
		WeakAdmissionRate:rate(weakAdmit,weakTotal),
		StrongAdmissionRate:rate(strongAdmit,strongTotal),
		DistractorAdmissionRate:rate(distractorAdmit,distractorTotal),
		HotAccuracy:rate(hotHits,hotTotal),
		WeakAccuracy:rate(weakHits,weakEval),
		StrongAccuracy:rate(strongHits,strongEval),
		DistractorAccuracy:rate(distractorHits,distractorEval),
		Target16Accuracy:rate(targetHits,targetTotal),
		Target16ExactAccuracy:rate(targetExact,episodes),
		AllCandidateAccuracy:rate(candidateHits,candidateTotal),
		OneShotFalseAdmissions:falseAdmissions,
		RecallEntriesUsed:maxEntries,
		MeanCheckpointsPerEpisode:float64(totalCheckpoints)/float64(episodes),
		AdmissionMemoryBytes:128,
		TotalBoundedMemoryBytes:maxEntries*16+136,
	}
}

func RunUP116C()(UP116CExact3GeneralityResult,error) {
	result:=UP116CExact3GeneralityResult{
		Schema:UP116CExact3GeneralitySchema,
		Experiment:"UP-116C-exact3-generality",
		SourceUP115CSeal:"b1623072a56658734db37d149a0b50824ae17fc2",
		HotKeys:12,GenerationInterval:32,AdmissionSightings:3,AdmissionMemoryBytes:128,ExactRecallCap:16,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,PhaseLabelUsed:false,
	}
	seeds:=[]int{219000000,220000000}
	result.Points=append(result.Points,
		up116cRun("mixed_gap_clean",seeds),
		up116cRun("four_hit_distractor",seeds),
	)
	return result,nil
}
