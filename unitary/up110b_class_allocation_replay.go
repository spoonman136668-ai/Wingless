package unitary

const UP110BClassAllocationSchema = "wingless.up110b-class-allocation-replay.v1"

type UP110BPoint struct {
	Arm                       string             `json:"arm"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP110BClassAllocationResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP109BSeal string        `json:"source_up109b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	Points           []UP110BPoint `json:"points"`
}

func up110bClassMembers(previous []up106bVerbSpec, class int) []up106bVerbSpec {
	out:=[]up106bVerbSpec{}
	for _,s:=range previous {
		if s.class==class { out=append(out,s) }
	}
	return out
}

func up110bEqualPrior(previous []up106bVerbSpec, epoch int) []up109bReplayExample {
	if len(previous)==0 { return nil }
	classes:=[]int{up97bStore,up97bStore,up97bObserve,up97bObserve,up97bReport,up97bReport}
	out:=make([]up109bReplayExample,0,6)
	for slot,class:=range classes {
		group:=up110bClassMembers(previous,class)
		var spec up106bVerbSpec
		if len(group)>0 {
			spec=group[(epoch+slot)%len(group)]
		}else{
			spec=previous[(epoch+slot)%len(previous)]
		}
		out=append(out,up109bReplayExample{spec:spec,ni:(epoch+slot)&1})
	}
	return out
}

func up110bExposureBalanced(previous []up106bVerbSpec,currentClass,epoch int) []up109bReplayExample {
	if len(previous)==0 { return nil }
	groups:=map[int][]up106bVerbSpec{
		up97bStore:up110bClassMembers(previous,up97bStore),
		up97bObserve:up110bClassMembers(previous,up97bObserve),
		up97bReport:up110bClassMembers(previous,up97bReport),
	}
	if len(groups[up97bStore])==0||len(groups[up97bObserve])==0||len(groups[up97bReport])==0 {
		return up110bEqualPrior(previous,epoch)
	}
	out:=make([]up109bReplayExample,0,6)
	slot:=0
	for _,class:=range []int{up97bStore,up97bObserve,up97bReport} {
		if class==currentClass { continue }
		group:=groups[class]
		for j:=0;j<3;j++ {
			spec:=group[(epoch+j)%len(group)]
			out=append(out,up109bReplayExample{spec:spec,ni:(epoch+slot)&1})
			slot++
		}
	}
	return out
}

func up110bAcquire(c *up97bClassifier,current up106bVerbSpec,previous []up106bVerbSpec,arm string) {
	for epoch:=0;epoch<20;epoch++ {
		up109bGroundCurrent(c,current)
		up109bBaseReplay(c)
		var seq []up109bReplayExample
		switch arm {
		case "pair_major":
			seq=up109bReplaySequence(previous,epoch)
		case "equal_prior_class":
			seq=up110bEqualPrior(previous,epoch)
		case "exposure_balanced":
			seq=up110bExposureBalanced(previous,current.class,epoch)
		}
		up109bTrainReplay(c,seq)
	}
}

func up110bPoint(arm string,stage int,c *up97bClassifier) UP110BPoint {
	p:=up106bPoint(stage,c)
	return UP110BPoint{
		Arm:arm,Stage:p.Stage,AcquiredVerbs:p.AcquiredVerbs,
		BaseSeenAccuracy:p.BaseSeenAccuracy,
		AcquiredAggregateAccuracy:p.AcquiredAggregateAccuracy,
		VerbAccuracy:p.VerbAccuracy,
	}
}

func RunUP110B()(UP110BClassAllocationResult,error){
	result:=UP110BClassAllocationResult{
		Schema:UP110BClassAllocationSchema,Experiment:"UP-110B-class-allocation-replay",
		SourceUP109BSeal:"734f49c0da903da1505353d44e41b9311d0fc550",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,FixedBudget:6,
	}
	for _,arm:=range []string{"pair_major","equal_prior_class","exposure_balanced"} {
		c:=up106bTrainBase()
		result.Points=append(result.Points,up110bPoint(arm,0,c))
		acquired:=[]up106bVerbSpec{}
		for stage,spec:=range up106bAcquire {
			up110bAcquire(c,spec,acquired,arm)
			acquired=append(acquired,spec)
			result.Points=append(result.Points,up110bPoint(arm,stage+1,c))
		}
	}
	return result,nil
}
