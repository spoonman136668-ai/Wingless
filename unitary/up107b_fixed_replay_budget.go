package unitary

const UP107BFixedReplaySchema = "wingless.up107b-fixed-replay-budget.v1"

type UP107BPoint struct {
	Arm                       string             `json:"arm"`
	FixedReplayBudget         int                `json:"fixed_replay_budget"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP107BFixedReplayResult struct {
	Schema             string        `json:"schema"`
	Experiment         string        `json:"experiment"`
	SourceUP106BSeal   string        `json:"source_up106b_seal"`
	StateDimension     int           `json:"state_dimension"`
	EpochsPerBatch     int           `json:"epochs_per_batch"`
	LearningRate       float64       `json:"learning_rate"`
	GroundingPerVerb   int           `json:"grounding_per_verb"`
	BaseReplaySubjects int           `json:"base_replay_subjects"`
	Points             []UP107BPoint `json:"points"`
}

func up107bPoint(arm string,budget,stage int,c *up97bClassifier) UP107BPoint {
	p:=up106bPoint(stage,c)
	return UP107BPoint{
		Arm:arm,FixedReplayBudget:budget,
		Stage:p.Stage,AcquiredVerbs:p.AcquiredVerbs,
		BaseSeenAccuracy:p.BaseSeenAccuracy,
		AcquiredAggregateAccuracy:p.AcquiredAggregateAccuracy,
		VerbAccuracy:p.VerbAccuracy,
	}
}

func up107bReplayFull(c *up97bClassifier,previous []up106bVerbSpec) {
	for _,spec:=range previous {
		up106bTrainStep(c,0,spec)
		up106bTrainStep(c,1,spec)
	}
}

func up107bReplayFixed(c *up97bClassifier,previous []up106bVerbSpec,budget,epoch int) {
	if len(previous)==0 || budget<=0 { return }
	type pair struct{ spec up106bVerbSpec; ni int }
	pairs:=make([]pair,0,len(previous)*2)
	for _,spec:=range previous {
		pairs=append(pairs,pair{spec:spec,ni:0},pair{spec:spec,ni:1})
	}
	start:=(epoch*budget)%len(pairs)
	for i:=0;i<budget;i++ {
		p:=pairs[(start+i)%len(pairs)]
		up106bTrainStep(c,p.ni,p.spec)
	}
}

func up107bAcquire(c *up97bClassifier,current up106bVerbSpec,previous []up106bVerbSpec,arm string,budget int) {
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ { up106bTrainStep(c,ni,current) }
		for _,spec:=range up106bBase { up106bTrainStep(c,0,spec) }
		if arm=="full_two_per_prior" {
			up107bReplayFull(c,previous)
		} else {
			up107bReplayFixed(c,previous,budget,epoch)
		}
	}
}

func RunUP107B()(UP107BFixedReplayResult,error){
	result:=UP107BFixedReplayResult{
		Schema:UP107BFixedReplaySchema,
		Experiment:"UP-107B-fixed-replay-budget",
		SourceUP106BSeal:"0918523a305c817e156a7326f0888445874d7f8d",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,
		GroundingPerVerb:4,BaseReplaySubjects:1,
	}
	arms:=[]struct{ name string; budget int }{
		{"full_two_per_prior",0},
		{"fixed_2",2},
		{"fixed_4",4},
		{"fixed_6",6},
	}
	for _,arm:=range arms {
		c:=up106bTrainBase()
		result.Points=append(result.Points,up107bPoint(arm.name,arm.budget,0,c))
		acquired:=[]up106bVerbSpec{}
		for stage,spec:=range up106bAcquire {
			up107bAcquire(c,spec,acquired,arm.name,arm.budget)
			acquired=append(acquired,spec)
			result.Points=append(result.Points,up107bPoint(arm.name,arm.budget,stage+1,c))
		}
	}
	return result,nil
}
