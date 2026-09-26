package unitary

const UP103BRehearsalGroundingSchema = "wingless.up103b-rehearsal-grounding.v1"

type UP103BArmPoint struct {
	Arm               string  `json:"arm"`
	HeldoutAccuracy   float64 `json:"heldout_accuracy"`
	StorePrecision    float64 `json:"store_precision"`
	StoreRecall       float64 `json:"store_recall"`
	ObservePrecision  float64 `json:"observe_precision"`
	ObserveRecall     float64 `json:"observe_recall"`
	ReportPrecision   float64 `json:"report_precision"`
	ReportRecall      float64 `json:"report_recall"`
	SeenVerbAccuracy  float64 `json:"seen_verb_accuracy"`
}

type UP103BRehearsalGroundingResult struct {
	Schema            string           `json:"schema"`
	Experiment        string           `json:"experiment"`
	SourceUP102BSeal  string           `json:"source_up102b_seal"`
	StateDimension    int              `json:"state_dimension"`
	Epochs            int              `json:"epochs"`
	LearningRate      float64          `json:"learning_rate"`
	GroundingPerVerb  int              `json:"grounding_per_verb"`
	Points            []UP103BArmPoint `json:"points"`
}

func up103bTrainStep(c *up97bClassifier, ni, vi int) {
	h:=up99bPrefixEncoder(ni,vi,0)
	p:=c.probs(h)
	target:=up97bVerbClass(vi)
	for k:=0;k<3;k++ {
		g:=p[k]
		if k==target { g-=1 }
		for i:=0;i<64;i++ { c.w[k][i]-=0.08*g*h[i] }
		c.b[k]-=0.08*g
	}
}

func up103bReplayAdapt(base *up97bClassifier) *up97bClassifier {
	copyModel:=*base
	c:=&copyModel
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ {
			for _,vi:=range up101bHeldoutVerbs { up103bTrainStep(c,ni,vi) }
		}
		for _,vi:=range up101bTrainVerbs { up103bTrainStep(c,0,vi) }
	}
	return c
}

func up103bPoint(arm string,c *up97bClassifier) UP103BArmPoint {
	acc,p,r:=up102bHeldoutEval(c)
	return UP103BArmPoint{
		Arm:arm,HeldoutAccuracy:acc,
		StorePrecision:p[up97bStore],StoreRecall:r[up97bStore],
		ObservePrecision:p[up97bObserve],ObserveRecall:r[up97bObserve],
		ReportPrecision:p[up97bReport],ReportRecall:r[up97bReport],
		SeenVerbAccuracy:up102bSeenAccuracy(c),
	}
}

func RunUP103B()(UP103BRehearsalGroundingResult,error){
	base:=up101bTrain()
	adaptOnly:=up102bAdapt(base,4)
	replay:=up103bReplayAdapt(base)
	return UP103BRehearsalGroundingResult{
		Schema:UP103BRehearsalGroundingSchema,
		Experiment:"UP-103B-rehearsal-grounding",
		SourceUP102BSeal:"1fdb08115705ef9455e3c219fa5810e6f4f1b31e",
		StateDimension:64,Epochs:20,LearningRate:0.08,GroundingPerVerb:4,
		Points:[]UP103BArmPoint{
			up103bPoint("adaptation_only",adaptOnly),
			up103bPoint("rehearsal_stabilized",replay),
		},
	},nil
}
