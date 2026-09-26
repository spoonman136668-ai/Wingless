package unitary

const UP106BVocabularyGrowthSchema = "wingless.up106b-vocabulary-growth.v1"

type UP106BStagePoint struct {
	Stage                     int                `json:"stage"`
	AcquiredVerbs             int                `json:"acquired_verbs"`
	BaseSeenAccuracy          float64            `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64            `json:"acquired_aggregate_accuracy"`
	VerbAccuracy              map[string]float64 `json:"verb_accuracy"`
}

type UP106BVocabularyGrowthResult struct {
	Schema                 string             `json:"schema"`
	Experiment             string             `json:"experiment"`
	SourceUP105BSeal       string             `json:"source_up105b_seal"`
	StateDimension         int                `json:"state_dimension"`
	EpochsPerBatch         int                `json:"epochs_per_batch"`
	LearningRate           float64            `json:"learning_rate"`
	GroundingPerVerb       int                `json:"grounding_per_verb"`
	BaseReplaySubjects     int                `json:"base_replay_subjects"`
	PriorReplaySubjects    int                `json:"prior_replay_subjects"`
	Points                 []UP106BStagePoint `json:"points"`
}

type up106bVerbSpec struct {
	verb  string
	class int
}

var up106bBase = []up106bVerbSpec{
	{"stores", up97bStore},
	{"keeps", up97bStore},
	{"observes", up97bObserve},
	{"sees", up97bObserve},
	{"reports", up97bReport},
	{"recalls", up97bReport},
}

var up106bAcquire = []up106bVerbSpec{
	{"holds", up97bStore},
	{"notes", up97bObserve},
	{"tells", up97bReport},
	{"saves", up97bStore},
	{"watches", up97bObserve},
	{"remembers", up97bReport},
}

func up106bEncode(ni int, verb string) [64]float64 {
	return up95bEncode(up97bNames[ni] + " " + verb)
}

func up106bTrainStep(c *up97bClassifier, ni int, spec up106bVerbSpec) {
	h := up106bEncode(ni, spec.verb)
	p := c.probs(h)
	for k := 0; k < 3; k++ {
		g := p[k]
		if k == spec.class { g -= 1 }
		for i := 0; i < 64; i++ {
			c.w[k][i] -= 0.08 * g * h[i]
		}
		c.b[k] -= 0.08 * g
	}
}

func up106bTrainBase() *up97bClassifier {
	c := &up97bClassifier{}
	for epoch := 0; epoch < 20; epoch++ {
		for ni := 0; ni < 6; ni++ {
			for _, spec := range up106bBase {
				up106bTrainStep(c, ni, spec)
			}
		}
	}
	return c
}

func up106bAccuracy(c *up97bClassifier, spec up106bVerbSpec, start, end int) float64 {
	hits,total := 0,0
	for ni := start; ni < end; ni++ {
		pred := up97bArgmax(c.probs(up106bEncode(ni, spec.verb)))
		total++
		if pred == spec.class { hits++ }
	}
	return float64(hits)/float64(total)
}

func up106bBaseAccuracy(c *up97bClassifier) float64 {
	sum := 0.0
	for _,spec := range up106bBase {
		sum += up106bAccuracy(c,spec,0,6)
	}
	return sum/float64(len(up106bBase))
}

func up106bPoint(stage int,c *up97bClassifier) UP106BStagePoint {
	m := map[string]float64{}
	sum := 0.0
	for i,spec := range up106bAcquire {
		acc := up106bAccuracy(c,spec,4,6)
		m[spec.verb]=acc
		if i < stage { sum += acc }
	}
	agg := 0.0
	if stage > 0 { agg = sum/float64(stage) }
	return UP106BStagePoint{
		Stage:stage,AcquiredVerbs:stage,
		BaseSeenAccuracy:up106bBaseAccuracy(c),
		AcquiredAggregateAccuracy:agg,
		VerbAccuracy:m,
	}
}

func up106bAcquireOne(c *up97bClassifier,current up106bVerbSpec,previous []up106bVerbSpec) {
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ { up106bTrainStep(c,ni,current) }
		for _,spec:=range up106bBase { up106bTrainStep(c,0,spec) }
		for _,spec:=range previous {
			for ni:=0;ni<2;ni++ { up106bTrainStep(c,ni,spec) }
		}
	}
}

func RunUP106B()(UP106BVocabularyGrowthResult,error){
	c:=up106bTrainBase()
	result:=UP106BVocabularyGrowthResult{
		Schema:UP106BVocabularyGrowthSchema,
		Experiment:"UP-106B-vocabulary-growth",
		SourceUP105BSeal:"20ff005b26a6b5cd955dadfc947e66f43a7f606a",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,
		GroundingPerVerb:4,BaseReplaySubjects:1,PriorReplaySubjects:2,
	}
	result.Points=append(result.Points,up106bPoint(0,c))
	acquired:=[]up106bVerbSpec{}
	for stage,spec:=range up106bAcquire {
		up106bAcquireOne(c,spec,acquired)
		acquired=append(acquired,spec)
		result.Points=append(result.Points,up106bPoint(stage+1,c))
	}
	return result,nil
}
