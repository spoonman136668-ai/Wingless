package unitary

const UP109BReplayTimingSchema = "wingless.up109b-replay-timing.v1"

type UP109BPoint struct {
	Arm                       string             `json:"arm"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP109BReplayTimingResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP108BSeal string        `json:"source_up108b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	Points           []UP109BPoint `json:"points"`
}

type up109bReplayExample struct {
	spec up106bVerbSpec
	ni   int
}

func up109bReplaySequence(previous []up106bVerbSpec, epoch int) []up109bReplayExample {
	if len(previous)==0 { return nil }
	pairs:=make([]up109bReplayExample,0,len(previous)*2)
	for _,spec:=range previous {
		pairs=append(pairs,up109bReplayExample{spec:spec,ni:0},up109bReplayExample{spec:spec,ni:1})
	}
	start:=(epoch*6)%len(pairs)
	out:=make([]up109bReplayExample,0,6)
	for i:=0;i<6;i++ { out=append(out,pairs[(start+i)%len(pairs)]) }
	return out
}

func up109bTrainReplay(c *up97bClassifier, seq []up109bReplayExample) {
	for _,x:=range seq { up106bTrainStep(c,x.ni,x.spec) }
}

func up109bGroundCurrent(c *up97bClassifier,current up106bVerbSpec) {
	for ni:=0;ni<4;ni++ { up106bTrainStep(c,ni,current) }
}

func up109bBaseReplay(c *up97bClassifier) {
	for _,spec:=range up106bBase { up106bTrainStep(c,0,spec) }
}

func up109bAcquire(c *up97bClassifier,current up106bVerbSpec,previous []up106bVerbSpec,arm string) {
	for epoch:=0;epoch<20;epoch++ {
		seq:=up109bReplaySequence(previous,epoch)
		switch arm {
		case "after":
			up109bGroundCurrent(c,current)
			up109bBaseReplay(c)
			up109bTrainReplay(c,seq)
		case "before":
			up109bTrainReplay(c,seq)
			up109bGroundCurrent(c,current)
			up109bBaseReplay(c)
		case "split_3_3":
			if len(seq)>0 { up109bTrainReplay(c,seq[:3]) }
			up109bGroundCurrent(c,current)
			up109bBaseReplay(c)
			if len(seq)>0 { up109bTrainReplay(c,seq[3:]) }
		}
	}
}

func up109bPoint(arm string,stage int,c *up97bClassifier) UP109BPoint {
	p:=up106bPoint(stage,c)
	return UP109BPoint{
		Arm:arm,Stage:p.Stage,AcquiredVerbs:p.AcquiredVerbs,
		BaseSeenAccuracy:p.BaseSeenAccuracy,AcquiredAggregateAccuracy:p.AcquiredAggregateAccuracy,
		VerbAccuracy:p.VerbAccuracy,
	}
}

func RunUP109B()(UP109BReplayTimingResult,error){
	result:=UP109BReplayTimingResult{
		Schema:UP109BReplayTimingSchema,Experiment:"UP-109B-replay-timing",
		SourceUP108BSeal:"095c70f4ad4cafa1d3a92e3ebbe9d0dded53583d",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,FixedBudget:6,
	}
	for _,arm:=range []string{"after","before","split_3_3"} {
		c:=up106bTrainBase()
		result.Points=append(result.Points,up109bPoint(arm,0,c))
		acquired:=[]up106bVerbSpec{}
		for stage,spec:=range up106bAcquire {
			up109bAcquire(c,spec,acquired,arm)
			acquired=append(acquired,spec)
			result.Points=append(result.Points,up109bPoint(arm,stage+1,c))
		}
	}
	return result,nil
}
