package unitary

const UP127BBalancedStoreGateSchema = "wingless.up127b-balanced-store-gate.v1"

type UP127BSplitMetric struct {
	Arm       string  `json:"arm"`
	Split     string  `json:"split"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Examples  int     `json:"examples"`
}

type UP127BSurfaceMetric struct {
	Arm             string  `json:"arm"`
	Split           string  `json:"split"`
	Surface         string  `json:"surface"`
	TargetStore     bool    `json:"target_store"`
	MeanProbability float64 `json:"mean_probability"`
	MinProbability  float64 `json:"min_probability"`
	MaxProbability  float64 `json:"max_probability"`
	PositiveRate    float64 `json:"positive_rate"`
	Examples        int     `json:"examples"`
}

type UP127BBalancedStoreGateResult struct {
	Schema               string                `json:"schema"`
	Experiment           string                `json:"experiment"`
	SourceUP126BSeal     string                `json:"source_up126b_seal"`
	StateDimension       int                   `json:"state_dimension"`
	GateEpochs           int                   `json:"gate_epochs"`
	GateLearningRate     float64               `json:"gate_learning_rate"`
	GateThreshold        float64               `json:"gate_threshold"`
	PositiveWeight       float64               `json:"positive_weight"`
	ThresholdChanged     bool                  `json:"threshold_changed"`
	EncoderChanged       bool                  `json:"encoder_changed"`
	SplitMetrics         []UP127BSplitMetric   `json:"split_metrics"`
	SurfaceMetrics       []UP127BSurfaceMetric `json:"surface_metrics"`
}

func up127bGateStepWeighted(g *up125bGate,name,verb string,target,weight float64) {
	h:=up95bEncode(name+" "+verb)
	p:=up125bGateProb(g,h)
	d:=(p-target)*weight
	for i:=0;i<64;i++ { g.w[i]-=0.08*d*h[i] }
	g.b-=0.08*d
}

func up127bTrainBalanced()*up125bGate{
	g:=&up125bGate{}
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ {
			name:=up121bOriginalNames[ni]
			for _,verb:=range up125bStoreSurfaces { up127bGateStepWeighted(g,name,verb,1,2) }
			for _,verb:=range up125bNonStoreSurfaces { up127bGateStepWeighted(g,name,verb,0,1) }
		}
	}
	return g
}

func up127bSplitMetric(arm string,g *up125bGate,split string,names []string) UP127BSplitMetric {
	m:=up125bGateEval(g,split,names)
	return UP127BSplitMetric{Arm:arm,Split:split,Accuracy:m.Accuracy,Precision:m.Precision,Recall:m.Recall,Examples:m.Examples}
}

func up127bSurfaceMetric(arm string,g *up125bGate,split,surface string,target bool,names []string) UP127BSurfaceMetric {
	m:=up126bSurfaceMetric(g,split,surface,target,names)
	return UP127BSurfaceMetric{
		Arm:arm,Split:split,Surface:surface,TargetStore:target,
		MeanProbability:m.MeanProbability,MinProbability:m.MinProbability,MaxProbability:m.MaxProbability,
		PositiveRate:m.PositiveRate,Examples:m.Examples,
	}
}

func RunUP127B()(UP127BBalancedStoreGateResult,error){
	result:=UP127BBalancedStoreGateResult{
		Schema:UP127BBalancedStoreGateSchema,Experiment:"UP-127B-balanced-store-gate",
		SourceUP126BSeal:"ad44f3b77c0914f60198de80f22cd8094f967366",
		StateDimension:64,GateEpochs:20,GateLearningRate:0.08,GateThreshold:0.5,PositiveWeight:2,
		ThresholdChanged:false,EncoderChanged:false,
	}
	arms:=[]struct{name string; gate *up125bGate}{
		{"original_unbalanced",up125bTrainGate()},
		{"balanced_positive_weight2",up127bTrainBalanced()},
	}
	splits:=[]struct{name string; names []string}{
		{"train",up121bOriginalNames[:4]},
		{"heldout",up121bOriginalNames[4:6]},
		{"unseen",up121bUnseenNames},
	}
	for _,arm:=range arms {
		for _,s:=range splits {
			result.SplitMetrics=append(result.SplitMetrics,up127bSplitMetric(arm.name,arm.gate,s.name,s.names))
			for _,surface:=range up125bStoreSurfaces {
				result.SurfaceMetrics=append(result.SurfaceMetrics,up127bSurfaceMetric(arm.name,arm.gate,s.name,surface,true,s.names))
			}
			for _,surface:=range up125bNonStoreSurfaces {
				result.SurfaceMetrics=append(result.SurfaceMetrics,up127bSurfaceMetric(arm.name,arm.gate,s.name,surface,false,s.names))
			}
		}
	}
	return result,nil
}
