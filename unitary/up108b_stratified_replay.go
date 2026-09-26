package unitary

const UP108BStratifiedReplaySchema = "wingless.up108b-stratified-replay.v1"

type UP108BPoint struct {
	Arm                       string             `json:"arm"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP108BStratifiedReplayResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP107BSeal string        `json:"source_up107b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	Points           []UP108BPoint `json:"points"`
}

func up108bPoint(arm string,stage int,c *up97bClassifier) UP108BPoint {
	p:=up106bPoint(stage,c)
	return UP108BPoint{
		Arm:arm,Stage:p.Stage,AcquiredVerbs:p.AcquiredVerbs,
		BaseSeenAccuracy:p.BaseSeenAccuracy,AcquiredAggregateAccuracy:p.AcquiredAggregateAccuracy,
		VerbAccuracy:p.VerbAccuracy,
	}
}

func up108bReplayStratified(c *up97bClassifier,previous []up106bVerbSpec,epoch int) {
	if len(previous)==0 { return }
	type pair struct{ spec up106bVerbSpec; ni int }
	candidates:=make([]pair,0,len(previous)*2)
	for i,spec:=range previous {
		first:=(epoch+i)&1
		candidates=append(candidates,pair{spec:spec,ni:first})
	}
	for i,spec:=range previous {
		first:=(epoch+i)&1
		candidates=append(candidates,pair{spec:spec,ni:1-first})
	}
	for i:=0;i<6;i++ {
		p:=candidates[i%len(candidates)]
		up106bTrainStep(c,p.ni,p.spec)
	}
}

func up108bAcquire(c *up97bClassifier,current up106bVerbSpec,previous []up106bVerbSpec,arm string) {
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ { up106bTrainStep(c,ni,current) }
		for _,spec:=range up106bBase { up106bTrainStep(c,0,spec) }
		switch arm {
		case "full_two_per_prior":
			up107bReplayFull(c,previous)
		case "pair_major_fixed6":
			up107bReplayFixed(c,previous,6,epoch)
		case "verb_stratified_fixed6":
			up108bReplayStratified(c,previous,epoch)
		}
	}
}

func RunUP108B()(UP108BStratifiedReplayResult,error){
	result:=UP108BStratifiedReplayResult{
		Schema:UP108BStratifiedReplaySchema,Experiment:"UP-108B-stratified-replay",
		SourceUP107BSeal:"b368ef8385e231d69ac2379e8f8b9a39d31bd4e5",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,FixedBudget:6,
	}
	for _,arm:=range []string{"full_two_per_prior","pair_major_fixed6","verb_stratified_fixed6"} {
		c:=up106bTrainBase()
		result.Points=append(result.Points,up108bPoint(arm,0,c))
		acquired:=[]up106bVerbSpec{}
		for stage,spec:=range up106bAcquire {
			up108bAcquire(c,spec,acquired,arm)
			acquired=append(acquired,spec)
			result.Points=append(result.Points,up108bPoint(arm,stage+1,c))
		}
	}
	return result,nil
}
