package unitary

const UP114BClassOrderSchema = "wingless.up114b-class-order-generalization.v1"

type UP114BPoint struct {
	ClassOrder                string             `json:"class_order"`
	Policy                    string             `json:"policy"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP114BClassOrderResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP113BSeal string        `json:"source_up113b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	Points           []UP114BPoint `json:"points"`
}

func up114bClassPair(class int) []up106bVerbSpec {
	switch class {
	case up97bStore:
		return []up106bVerbSpec{{verb:"holds",class:up97bStore},{verb:"saves",class:up97bStore}}
	case up97bObserve:
		return []up106bVerbSpec{{verb:"notes",class:up97bObserve},{verb:"watches",class:up97bObserve}}
	default:
		return []up106bVerbSpec{{verb:"tells",class:up97bReport},{verb:"remembers",class:up97bReport}}
	}
}

func up114bSequence(order []int) []up106bVerbSpec {
	out:=make([]up106bVerbSpec,0,6)
	for _,class:=range order { out=append(out,up114bClassPair(class)[0]) }
	for _,class:=range order { out=append(out,up114bClassPair(class)[1]) }
	return out
}

func up114bOrderName(order []int) string {
	b:=make([]byte,0,3)
	for _,class:=range order {
		switch class {
		case up97bStore: b=append(b,'S')
		case up97bObserve: b=append(b,'O')
		case up97bReport: b=append(b,'R')
		}
	}
	return string(b)
}

func up114bPoint(orderName,policy string,stage int,c *up97bClassifier,seq []up106bVerbSpec) UP114BPoint {
	m:=map[string]float64{}
	sum:=0.0
	for i,spec:=range seq {
		acc:=up106bAccuracy(c,spec,4,6)
		m[spec.verb]=acc
		if i<stage { sum+=acc }
	}
	agg:=0.0
	if stage>0 { agg=sum/float64(stage) }
	return UP114BPoint{
		ClassOrder:orderName,Policy:policy,Stage:stage,AcquiredVerbs:stage,
		BaseSeenAccuracy:up106bBaseAccuracy(c),AcquiredAggregateAccuracy:agg,VerbAccuracy:m,
	}
}

func RunUP114B()(UP114BClassOrderResult,error){
	result:=UP114BClassOrderResult{
		Schema:UP114BClassOrderSchema,Experiment:"UP-114B-class-order-generalization",
		SourceUP113BSeal:"bea5f6847499e5287788309bb66d85f73a6efbdd",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,FixedBudget:6,
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},
		{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},
		{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},
		{up97bReport,up97bObserve,up97bStore},
	}
	for _,order:=range orders {
		name:=up114bOrderName(order)
		seq:=up114bSequence(order)
		for _,policy:=range []string{"current_class_excluded","stage3_anchor2"} {
			c:=up106bTrainBase()
			result.Points=append(result.Points,up114bPoint(name,policy,0,c,seq))
			acquired:=[]up106bVerbSpec{}
			for idx,spec:=range seq {
				stage:=idx+1
				up113bAcquire(c,spec,acquired,policy,stage)
				acquired=append(acquired,spec)
				result.Points=append(result.Points,up114bPoint(name,policy,stage,c,seq))
			}
		}
	}
	return result,nil
}
