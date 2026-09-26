package unitary

const UP115BBaseClassDiagnosticSchema = "wingless.up115b-base-class-diagnostic.v1"

type UP115BPoint struct {
	ClassOrder                string             `json:"class_order"`
	Policy                    string             `json:"policy"`
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	BaseStoreAccuracy         float64            `json:"base_store_accuracy"`
	BaseObserveAccuracy       float64            `json:"base_observe_accuracy"`
	BaseReportAccuracy        float64            `json:"base_report_accuracy"`
	BaseVerbAccuracy          map[string]float64 `json:"base_verb_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	AcquiredVerbAccuracy      map[string]float64 `json:"acquired_verb_accuracy"`
}

type UP115BBaseClassDiagnosticResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP114BSeal string        `json:"source_up114b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	FixedBudget      int           `json:"fixed_budget"`
	TrainingChanged  bool          `json:"training_changed"`
	Points           []UP115BPoint `json:"points"`
}

func up115bBaseDiagnostics(c *up97bClassifier)(float64,float64,float64,map[string]float64){
	perVerb:=map[string]float64{}
	classSum:=map[int]float64{up97bStore:0,up97bObserve:0,up97bReport:0}
	classN:=map[int]int{up97bStore:0,up97bObserve:0,up97bReport:0}
	for _,spec:=range up106bBase {
		acc:=up106bAccuracy(c,spec,0,6)
		perVerb[spec.verb]=acc
		classSum[spec.class]+=acc
		classN[spec.class]++
	}
	return classSum[up97bStore]/float64(classN[up97bStore]),
		classSum[up97bObserve]/float64(classN[up97bObserve]),
		classSum[up97bReport]/float64(classN[up97bReport]),
		perVerb
}

func up115bPoint(orderName,policy string,stage int,c *up97bClassifier,seq []up106bVerbSpec) UP115BPoint {
	store,observe,report,baseVerb:=up115bBaseDiagnostics(c)
	acquired:=map[string]float64{}
	sum:=0.0
	for i,spec:=range seq {
		acc:=up106bAccuracy(c,spec,4,6)
		acquired[spec.verb]=acc
		if i<stage { sum+=acc }
	}
	agg:=0.0
	if stage>0 { agg=sum/float64(stage) }
	return UP115BPoint{
		ClassOrder:orderName,Policy:policy,Stage:stage,AcquiredVerbs:stage,
		BaseSeenAccuracy:up106bBaseAccuracy(c),
		BaseStoreAccuracy:store,BaseObserveAccuracy:observe,BaseReportAccuracy:report,
		BaseVerbAccuracy:baseVerb,
		AcquiredAggregateAccuracy:agg,AcquiredVerbAccuracy:acquired,
	}
}

func RunUP115B()(UP115BBaseClassDiagnosticResult,error){
	result:=UP115BBaseClassDiagnosticResult{
		Schema:UP115BBaseClassDiagnosticSchema,
		Experiment:"UP-115B-base-class-diagnostic",
		SourceUP114BSeal:"3c64686896dd22342a1d5f230c1292f6f9c91755",
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
			result.Points=append(result.Points,up115bPoint(name,policy,0,c,seq))
			acquired:=[]up106bVerbSpec{}
			for idx,spec:=range seq {
				stage:=idx+1
				up113bAcquire(c,spec,acquired,policy,stage)
				acquired=append(acquired,spec)
				result.Points=append(result.Points,up115bPoint(name,policy,stage,c,seq))
			}
		}
	}
	return result,nil
}
