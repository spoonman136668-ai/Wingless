package unitary

const UP113BStageLimitedAnchorSchema = "wingless.up113b-stage-limited-anchor.v1"

type UP113BPoint struct {
	Arm                       string             `json:"arm"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP113BStageLimitedAnchorResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP112BSeal string        `json:"source_up112b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	Points           []UP113BPoint `json:"points"`
}

func up113bAcquire(c *up97bClassifier,current up106bVerbSpec,previous []up106bVerbSpec,arm string,stage int) {
	for epoch:=0;epoch<20;epoch++ {
		up109bGroundCurrent(c,current)
		up109bBaseReplay(c)
		anchor:=false
		switch arm {
		case "stage3_anchor2":
			anchor=stage==3
		case "stages2_3_anchor2":
			anchor=stage==2||stage==3
		}
		var seq []up109bReplayExample
		if anchor {
			seq=up112bExtraReplay(previous,current.class,epoch,2)
		}else{
			seq=up111bCurrentExcluded(previous,current.class,epoch)
		}
		up109bTrainReplay(c,seq)
	}
}

func up113bPoint(arm string,stage int,c *up97bClassifier) UP113BPoint {
	p:=up106bPoint(stage,c)
	return UP113BPoint{
		Arm:arm,Stage:p.Stage,AcquiredVerbs:p.AcquiredVerbs,
		BaseSeenAccuracy:p.BaseSeenAccuracy,AcquiredAggregateAccuracy:p.AcquiredAggregateAccuracy,
		VerbAccuracy:p.VerbAccuracy,
	}
}

func RunUP113B()(UP113BStageLimitedAnchorResult,error){
	result:=UP113BStageLimitedAnchorResult{
		Schema:UP113BStageLimitedAnchorSchema,
		Experiment:"UP-113B-stage-limited-anchor",
		SourceUP112BSeal:"21c8645d492af224fdae7f734af62fce31af6f88",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,FixedBudget:6,
	}
	for _,arm:=range []string{"current_class_excluded","stage3_anchor2","stages2_3_anchor2"} {
		c:=up106bTrainBase()
		result.Points=append(result.Points,up113bPoint(arm,0,c))
		acquired:=[]up106bVerbSpec{}
		for idx,spec:=range up106bAcquire {
			stage:=idx+1
			up113bAcquire(c,spec,acquired,arm,stage)
			acquired=append(acquired,spec)
			result.Points=append(result.Points,up113bPoint(arm,stage,c))
		}
	}
	return result,nil
}
