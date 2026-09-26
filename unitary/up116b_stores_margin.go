package unitary

const UP116BStoresMarginSchema = "wingless.up116b-stores-margin.v1"

type UP116BVerbDecision struct {
	Accuracy          float64 `json:"accuracy"`
	MeanTargetMargin  float64 `json:"mean_target_margin"`
	MinTargetMargin   float64 `json:"min_target_margin"`
	PredictedStore    int     `json:"predicted_store"`
	PredictedObserve  int     `json:"predicted_observe"`
	PredictedReport   int     `json:"predicted_report"`
}

type UP116BPoint struct {
	ClassOrder                string              `json:"class_order"`
	Policy                    string              `json:"policy"`
	Stage                     int                 `json:"stage"`
	BaseSeenAccuracy          float64             `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64             `json:"acquired_aggregate_accuracy"`
	Stores                    UP116BVerbDecision  `json:"stores"`
	Keeps                     UP116BVerbDecision  `json:"keeps"`
}

type UP116BStoresMarginResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP115BSeal string        `json:"source_up115b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	TrainingChanged  bool          `json:"training_changed"`
	Points           []UP116BPoint `json:"points"`
}

func up116bDecision(c *up97bClassifier,spec up106bVerbSpec) UP116BVerbDecision {
	hits:=0
	marginSum:=0.0
	minMargin:=1e9
	counts:=[3]int{}
	for ni:=0;ni<6;ni++ {
		p:=c.probs(up106bEncode(ni,spec.verb))
		pred:=up97bArgmax(p)
		counts[pred]++
		if pred==spec.class { hits++ }
		maxOther:=-1.0
		for k:=0;k<3;k++ {
			if k==spec.class { continue }
			if p[k]>maxOther { maxOther=p[k] }
		}
		margin:=p[spec.class]-maxOther
		marginSum+=margin
		if margin<minMargin { minMargin=margin }
	}
	return UP116BVerbDecision{
		Accuracy:float64(hits)/6.0,
		MeanTargetMargin:marginSum/6.0,
		MinTargetMargin:minMargin,
		PredictedStore:counts[up97bStore],
		PredictedObserve:counts[up97bObserve],
		PredictedReport:counts[up97bReport],
	}
}

func up116bPoint(orderName,policy string,stage int,c *up97bClassifier,seq []up106bVerbSpec) UP116BPoint {
	sum:=0.0
	for i,spec:=range seq {
		if i>=stage { break }
		sum+=up106bAccuracy(c,spec,4,6)
	}
	agg:=0.0
	if stage>0 { agg=sum/float64(stage) }
	return UP116BPoint{
		ClassOrder:orderName,Policy:policy,Stage:stage,
		BaseSeenAccuracy:up106bBaseAccuracy(c),AcquiredAggregateAccuracy:agg,
		Stores:up116bDecision(c,up106bVerbSpec{verb:"stores",class:up97bStore}),
		Keeps:up116bDecision(c,up106bVerbSpec{verb:"keeps",class:up97bStore}),
	}
}

func RunUP116B()(UP116BStoresMarginResult,error){
	result:=UP116BStoresMarginResult{
		Schema:UP116BStoresMarginSchema,
		Experiment:"UP-116B-stores-margin",
		SourceUP115BSeal:"96e7968598df235267e07d099f0d77171a183f89",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,FixedBudget:6,
		TrainingChanged:false,
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
			result.Points=append(result.Points,up116bPoint(name,policy,0,c,seq))
			acquired:=[]up106bVerbSpec{}
			for idx,spec:=range seq {
				stage:=idx+1
				up113bAcquire(c,spec,acquired,policy,stage)
				acquired=append(acquired,spec)
				result.Points=append(result.Points,up116bPoint(name,policy,stage,c,seq))
			}
		}
	}
	return result,nil
}
