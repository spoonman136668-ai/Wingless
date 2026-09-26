package unitary

const UP111BCurrentClassExcludedSchema = "wingless.up111b-current-class-excluded.v1"

type UP111BPoint struct {
	Arm                       string             `json:"arm"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP111BCurrentClassExcludedResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP110BSeal string        `json:"source_up110b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	Points           []UP111BPoint `json:"points"`
}

func up111bCurrentExcluded(previous []up106bVerbSpec,currentClass,epoch int) []up109bReplayExample {
	if len(previous)==0 { return nil }
	eligibleClasses:=[]int{}
	for _,class:=range []int{up97bStore,up97bObserve,up97bReport} {
		if class==currentClass { continue }
		if len(up110bClassMembers(previous,class))>0 { eligibleClasses=append(eligibleClasses,class) }
	}
	if len(eligibleClasses)==0 {
		return up109bReplaySequence(previous,epoch)
	}
	out:=make([]up109bReplayExample,0,6)
	for slot:=0;slot<6;slot++ {
		class:=eligibleClasses[slot%len(eligibleClasses)]
		group:=up110bClassMembers(previous,class)
		spec:=group[(epoch+slot/len(eligibleClasses))%len(group)]
		out=append(out,up109bReplayExample{spec:spec,ni:(epoch+slot)&1})
	}
	return out
}

func up111bAcquire(c *up97bClassifier,current up106bVerbSpec,previous []up106bVerbSpec,arm string) {
	for epoch:=0;epoch<20;epoch++ {
		up109bGroundCurrent(c,current)
		up109bBaseReplay(c)
		var seq []up109bReplayExample
		switch arm {
		case "pair_major":
			seq=up109bReplaySequence(previous,epoch)
		case "exposure_balanced_control":
			seq=up110bExposureBalanced(previous,current.class,epoch)
		case "current_class_excluded":
			seq=up111bCurrentExcluded(previous,current.class,epoch)
		}
		up109bTrainReplay(c,seq)
	}
}

func up111bPoint(arm string,stage int,c *up97bClassifier) UP111BPoint {
	p:=up106bPoint(stage,c)
	return UP111BPoint{
		Arm:arm,Stage:p.Stage,AcquiredVerbs:p.AcquiredVerbs,
		BaseSeenAccuracy:p.BaseSeenAccuracy,AcquiredAggregateAccuracy:p.AcquiredAggregateAccuracy,
		VerbAccuracy:p.VerbAccuracy,
	}
}

func RunUP111B()(UP111BCurrentClassExcludedResult,error){
	result:=UP111BCurrentClassExcludedResult{
		Schema:UP111BCurrentClassExcludedSchema,Experiment:"UP-111B-current-class-excluded",
		SourceUP110BSeal:"892e78c7ea1e97b2a7c1cd7e635cfa47d049d471",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,FixedBudget:6,
	}
	for _,arm:=range []string{"pair_major","exposure_balanced_control","current_class_excluded"} {
		c:=up106bTrainBase()
		result.Points=append(result.Points,up111bPoint(arm,0,c))
		acquired:=[]up106bVerbSpec{}
		for stage,spec:=range up106bAcquire {
			up111bAcquire(c,spec,acquired,arm)
			acquired=append(acquired,spec)
			result.Points=append(result.Points,up111bPoint(arm,stage+1,c))
		}
	}
	return result,nil
}
