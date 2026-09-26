package unitary

import "fmt"

const UP124CHistoryOverflowSchema = "wingless.up124c-history-overflow-boundary.v1"

type UP124CPoint struct {
	OverflowQualifiedCandidates int     `json:"overflow_qualified_candidates"`
	OverflowEpisodeRate         float64 `json:"overflow_episode_rate"`
	OverflowEpisodes            int     `json:"overflow_episodes"`
	TotalEpisodes               int     `json:"total_episodes"`
	OverflowMessage             string  `json:"overflow_message"`
	ValidEvaluatedEpisodes      int     `json:"valid_evaluated_episodes"`
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
	RecallEntriesUsed           int     `json:"recall_entries_used"`
	MeanCheckpointsPerEvaluated float64 `json:"mean_checkpoints_per_evaluated_episode"`
	AdmissionMemoryBytes        int     `json:"admission_memory_bytes"`
	TotalBoundedMemoryBytes     int     `json:"total_bounded_memory_bytes"`
}

type UP124CHistoryOverflowResult struct {
	Schema                        string        `json:"schema"`
	Experiment                    string        `json:"experiment"`
	SourceUP123CSeal              string        `json:"source_up123c_seal"`
	HotKeys                       int           `json:"hot_keys"`
	ValidTargets                  int           `json:"valid_targets"`
	SetACount                     int           `json:"set_a_count"`
	SetBCount                     int           `json:"set_b_count"`
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
	ReplacementRuleAdded          bool          `json:"replacement_rule_added"`
	Points                        []UP124CPoint `json:"points"`
}

type up124cEpisode struct {
	overflow bool
	message string
	validAdmit int
	validTrials int
	validHits int
	validTotal int
	hotHits int
	hotTotal int
	targetHits int
	targetTotal int
	exact bool
	nonAdmit int
	nonTrials int
	nonHits int
	nonTotal int
	fp int
	maxCurrent int
	maxHistory int
	recallEntries int
	checkpoints int
}

func up124cEpisodeRun(q,base,ep int)(out up124cEpisode) {
	defer func(){
		if v:=recover();v!=nil {
			out.overflow=true
			out.message=fmt.Sprint(v)
		}
	}()

	rng:=newSQ0RNG(sq0Seed(base,6501+q*241,ep))
	x:=&up120cExtendedMachine{maxAge:2}
	truth:=map[int]int{}
	next:=300
	targets:=[]int{100,101,102,103}
	setA:=make([]int,16)
	setB:=make([]int,12)
	for i:=range setA { setA[i]=200+i }
	for i:=range setB { setB[i]=220+i }

	for k:=0;k<32;k++ {
		v:=rng.intn(32);truth[k]=v
		a,_:=up122cProcess(x,k,v,&out.maxCurrent,&out.maxHistory)
		if k>=16&&a { out.fp++ }
	}
	for k:=0;k<12;k++ { x.process(k,truth[k]);x.mem.query(k) }
	up122cObserve(x,&out.maxCurrent,&out.maxHistory)
	if x.counter!=0 { out.fp+=up122cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory) }

	// Generation 1: fill history with 16 age-0 qualified set-A records.
	for _,k:=range setA {
		v:=rng.intn(32);truth[k]=v
		up122cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory)
		out.nonTrials++
	}

	// Generation 2: add four targets + 12 set-B records, yielding exact 32/32 history occupancy.
	for _,k:=range targets {
		v:=rng.intn(32);truth[k]=v
		up122cPresent(x,k,v,2,&out.validAdmit,&out.maxCurrent,&out.maxHistory)
		out.validTrials++
	}
	for _,k:=range setB {
		v:=rng.intn(32);truth[k]=v
		up122cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory)
		out.nonTrials++
	}
	if x.counter!=0 { out.fp+=up122cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory) }

	// Overflow generation: qualified newcomers are scientifically allowed to trigger
	// the unchanged mechanism's existing full-history panic at the generation boundary.
	overflowKeys:=make([]int,q)
	for i:=0;i<q;i++ { overflowKeys[i]=240+i }
	for _,k:=range overflowKeys {
		v:=rng.intn(32);truth[k]=v
		up122cPresent(x,k,v,2,&out.nonAdmit,&out.maxCurrent,&out.maxHistory)
		out.nonTrials++
	}
	if x.counter!=0 { out.fp+=up122cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory) }

	// Only reachable when the unchanged mechanism survives the overflow generation.
	for _,k:=range targets {
		up122cPresent(x,k,truth[k],2,&out.validAdmit,&out.maxCurrent,&out.maxHistory)
	}
	if x.counter!=0 { out.fp+=up122cFill(x,&next,rng,&out.maxCurrent,&out.maxHistory) }

	for j:=0;j<12288;j++ {
		key:=1000+j
		a,_:=up122cProcess(x,key,rng.intn(32),&out.maxCurrent,&out.maxHistory)
		if a { out.fp++ }
		if (j+1)%4==0 {
			for k:=0;k<12;k++ { x.mem.query(k) }
			for _,k:=range targets { x.mem.query(k) }
			for _,k:=range setA { x.mem.query(k) }
			for _,k:=range setB { x.mem.query(k) }
			for _,k:=range overflowKeys { x.mem.query(k) }
		}
	}

	out.exact=true
	for k:=0;k<12;k++ {
		got,ok:=x.mem.query(k);out.hotTotal++;out.targetTotal++
		if ok&&got==truth[k] { out.hotHits++;out.targetHits++ } else { out.exact=false }
	}
	for _,k:=range targets {
		got,ok:=x.mem.query(k);out.validTotal++;out.targetTotal++
		if ok&&got==truth[k] { out.validHits++;out.targetHits++ } else { out.exact=false }
	}
	for _,group:=range [][]int{setA,setB,overflowKeys} {
		for _,k:=range group {
			got,ok:=x.mem.query(k);out.nonTotal++
			if ok&&got==truth[k] { out.nonHits++ }
		}
	}
	out.recallEntries=x.mem.count
	out.checkpoints=x.checkpoints
	return
}

func up124cRun(q int,seeds []int) UP124CPoint {
	var overflow,total,evaluated int
	var message string
	var validAdmit,validTrials,validHits,validTotal int
	var hotHits,hotTotal,targetHits,targetTotal,exactHits int
	var nonAdmit,nonTrials,nonHits,nonTotal,fp int
	var maxCurrent,maxHistory,maxRecall,checkpoints int
	for _,base:=range seeds {
		for ep:=0;ep<32;ep++ {
			e:=up124cEpisodeRun(q,base,ep)
			total++
			if e.maxCurrent>maxCurrent { maxCurrent=e.maxCurrent }
			if e.maxHistory>maxHistory { maxHistory=e.maxHistory }
			if e.recallEntries>maxRecall { maxRecall=e.recallEntries }
			fp+=e.fp;nonAdmit+=e.nonAdmit;nonTrials+=e.nonTrials
			if e.overflow {
				overflow++
				if message=="" { message=e.message }
				continue
			}
			evaluated++
			validAdmit+=e.validAdmit;validTrials+=e.validTrials
			validHits+=e.validHits;validTotal+=e.validTotal
			hotHits+=e.hotHits;hotTotal+=e.hotTotal
			targetHits+=e.targetHits;targetTotal+=e.targetTotal
			nonHits+=e.nonHits;nonTotal+=e.nonTotal
			if e.exact { exactHits++ }
			checkpoints+=e.checkpoints
		}
	}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	meanCheck:=0.0
	if evaluated>0 { meanCheck=float64(checkpoints)/float64(evaluated) }
	return UP124CPoint{
		OverflowQualifiedCandidates:q,OverflowEpisodeRate:rate(overflow,total),
		OverflowEpisodes:overflow,TotalEpisodes:total,OverflowMessage:message,
		ValidEvaluatedEpisodes:evaluated,ValidAdmissionRate:rate(validAdmit,validTrials),
		ValidAccuracy:rate(validHits,validTotal),HotAccuracy:rate(hotHits,hotTotal),
		Target16Accuracy:rate(targetHits,targetTotal),Target16ExactAccuracy:rate(exactHits,evaluated),
		NonpersistentAdmissionRate:rate(nonAdmit,nonTrials),NonpersistentRetentionRate:rate(nonHits,nonTotal),
		OneShotFalseAdmissions:fp,MaxCurrentTableEntries:maxCurrent,MaxHistoryTableEntries:maxHistory,
		RecallEntriesUsed:maxRecall,MeanCheckpointsPerEvaluated:meanCheck,
		AdmissionMemoryBytes:128,TotalBoundedMemoryBytes:392,
	}
}

func RunUP124C()(UP124CHistoryOverflowResult,error) {
	res:=UP124CHistoryOverflowResult{
		Schema:UP124CHistoryOverflowSchema,Experiment:"UP-124C-history-overflow-boundary",
		SourceUP123CSeal:"dfa1f6e4176b2a35213b91c97b82cfa8ae7af6f4",
		HotKeys:12,ValidTargets:4,SetACount:16,SetBCount:12,GenerationInterval:32,
		AdmissionMemoryBytes:128,ExactRecallCap:16,HistoryEntries:32,KeyBits:14,
		ChurnWrites:12288,EpisodesPerSeed:32,
		QueryEvidenceUsedForAdmission:false,SemanticPriorityUsed:false,FutureOracleUsed:false,
		ReplacementRuleAdded:false,
	}
	seeds:=[]int{235000000,236000000}
	for _,q:=range []int{0,1,4,8} { res.Points=append(res.Points,up124cRun(q,seeds)) }
	return res,nil
}
