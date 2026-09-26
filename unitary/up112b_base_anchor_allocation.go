package unitary

const UP112BBaseAnchorSchema = "wingless.up112b-base-anchor-allocation.v1"

type UP112BPoint struct {
	Arm                       string             `json:"arm"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP112BBaseAnchorResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP111BSeal string        `json:"source_up111b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	Points           []UP112BPoint `json:"points"`
}

func up112bBaseClassSpecs(class int) []up106bVerbSpec {
	out:=[]up106bVerbSpec{}
	for _,s:=range up106bBase {
		if s.class==class { out=append(out,s) }
	}
	return out
}

func up112bNonCurrentPrior(previous []up106bVerbSpec,currentClass,count,epoch int) []up109bReplayExample {
	if count<=0||len(previous)==0 { return nil }
	classes:=[]int{}
	for _,class:=range []int{up97bStore,up97bObserve,up97bReport} {
		if class==currentClass { continue }
		if len(up110bClassMembers(previous,class))>0 { classes=append(classes,class) }
	}
	if len(classes)==0 {
		seq:=up109bReplaySequence(previous,epoch)
		if count>len(seq){count=len(seq)}
		return seq[:count]
	}
	out:=make([]up109bReplayExample,0,count)
	for slot:=0;slot<count;slot++ {
		class:=classes[slot%len(classes)]
		group:=up110bClassMembers(previous,class)
		spec:=group[(epoch+slot/len(classes))%len(group)]
		out=append(out,up109bReplayExample{spec:spec,ni:(epoch+slot)&1})
	}
	return out
}

func up112bExtraReplay(previous []up106bVerbSpec,currentClass,epoch,anchorCount int) []up109bReplayExample {
	if len(previous)==0 { return nil }
	priorCount:=6-anchorCount
	out:=up112bNonCurrentPrior(previous,currentClass,priorCount,epoch)
	if anchorCount==0 { return out }
	base:=up112bBaseClassSpecs(currentClass)
	if len(base)!=2 { return out }
	if anchorCount>=2 {
		out=append(out,
			up109bReplayExample{spec:base[0],ni:1},
			up109bReplayExample{spec:base[1],ni:1},
		)
	}
	if anchorCount>=4 {
		out=append(out,
			up109bReplayExample{spec:base[0],ni:2},
			up109bReplayExample{spec:base[1],ni:2},
		)
	}
	return out
}

func up112bAcquire(c *up97bClassifier,current up106bVerbSpec,previous []up106bVerbSpec,arm string) {
	for epoch:=0;epoch<20;epoch++ {
		up109bGroundCurrent(c,current)
		up109bBaseReplay(c)
		var seq []up109bReplayExample
		switch arm {
		case "current_class_excluded":
			seq=up111bCurrentExcluded(previous,current.class,epoch)
		case "base_anchor_2":
			seq=up112bExtraReplay(previous,current.class,epoch,2)
		case "base_anchor_4":
			seq=up112bExtraReplay(previous,current.class,epoch,4)
		}
		up109bTrainReplay(c,seq)
	}
}

func up112bPoint(arm string,stage int,c *up97bClassifier) UP112BPoint {
	p:=up106bPoint(stage,c)
	return UP112BPoint{
		Arm:arm,Stage:p.Stage,AcquiredVerbs:p.AcquiredVerbs,
		BaseSeenAccuracy:p.BaseSeenAccuracy,AcquiredAggregateAccuracy:p.AcquiredAggregateAccuracy,
		VerbAccuracy:p.VerbAccuracy,
	}
}

func RunUP112B()(UP112BBaseAnchorResult,error){
	result:=UP112BBaseAnchorResult{
		Schema:UP112BBaseAnchorSchema,Experiment:"UP-112B-base-anchor-allocation",
		SourceUP111BSeal:"c9d4ba5a1a53b300c5fc435a41c3d965a4e1ac43",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,FixedBudget:6,
	}
	for _,arm:=range []string{"current_class_excluded","base_anchor_2","base_anchor_4"} {
		c:=up106bTrainBase()
		result.Points=append(result.Points,up112bPoint(arm,0,c))
		acquired:=[]up106bVerbSpec{}
		for stage,spec:=range up106bAcquire {
			up112bAcquire(c,spec,acquired,arm)
			acquired=append(acquired,spec)
			result.Points=append(result.Points,up112bPoint(arm,stage+1,c))
		}
	}
	return result,nil
}
