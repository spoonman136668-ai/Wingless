package unitary

import "math"

const UP126BStoreGateTopologySchema = "wingless.up126b-store-gate-topology.v1"

type UP126BSurfaceMetric struct {
	Split          string  `json:"split"`
	Surface        string  `json:"surface"`
	TargetStore    bool    `json:"target_store"`
	MeanProbability float64 `json:"mean_probability"`
	MinProbability  float64 `json:"min_probability"`
	MaxProbability  float64 `json:"max_probability"`
	PositiveRate     float64 `json:"positive_rate"`
	Examples         int     `json:"examples"`
}

type UP126BSplitMetric struct {
	Split     string  `json:"split"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Examples  int     `json:"examples"`
}

type UP126BStoreGateTopologyResult struct {
	Schema                  string                `json:"schema"`
	Experiment              string                `json:"experiment"`
	SourceUP125BSeal        string                `json:"source_up125b_seal"`
	StateDimension          int                   `json:"state_dimension"`
	GateEpochs              int                   `json:"gate_epochs"`
	GateLearningRate        float64               `json:"gate_learning_rate"`
	GateThreshold           float64               `json:"gate_threshold"`
	GateRetrained           bool                  `json:"gate_retrained"`
	ThresholdChanged        bool                  `json:"threshold_changed"`
	ProjectorUsed           bool                  `json:"projector_used"`
	SurfaceMetrics          []UP126BSurfaceMetric `json:"surface_metrics"`
	SplitMetrics            []UP126BSplitMetric   `json:"split_metrics"`
	MinStoreMarginToThreshold float64             `json:"min_store_margin_to_threshold"`
	MaxNonStoreProbability  float64               `json:"max_non_store_probability"`
}

func up126bSurfaceMetric(g *up125bGate,split,surface string,target bool,names []string) UP126BSurfaceMetric {
	sum:=0.0
	minp:=math.Inf(1)
	maxp:=math.Inf(-1)
	pos:=0
	for _,name:=range names {
		p:=up125bGateProb(g,up95bEncode(name+" "+surface))
		sum+=p
		if p<minp { minp=p }
		if p>maxp { maxp=p }
		if p>=0.5 { pos++ }
	}
	return UP126BSurfaceMetric{
		Split:split,Surface:surface,TargetStore:target,
		MeanProbability:sum/float64(len(names)),
		MinProbability:minp,MaxProbability:maxp,
		PositiveRate:float64(pos)/float64(len(names)),
		Examples:len(names),
	}
}

func RunUP126B()(UP126BStoreGateTopologyResult,error){
	g:=up125bTrainGate()
	result:=UP126BStoreGateTopologyResult{
		Schema:UP126BStoreGateTopologySchema,
		Experiment:"UP-126B-store-gate-topology",
		SourceUP125BSeal:"2e86cdcd20f88a1d6480ff647349b3baf90a8758",
		StateDimension:64,GateEpochs:20,GateLearningRate:0.08,GateThreshold:0.5,
		GateRetrained:false,ThresholdChanged:false,ProjectorUsed:false,
		MinStoreMarginToThreshold:math.Inf(1),
		MaxNonStoreProbability:math.Inf(-1),
	}
	splits:=[]struct{name string; names []string}{
		{"train",up121bOriginalNames[:4]},
		{"heldout",up121bOriginalNames[4:6]},
		{"unseen",up121bUnseenNames},
	}
	for _,s:=range splits {
		result.SplitMetrics=append(result.SplitMetrics,UP126BSplitMetric{
			Split:s.name,
			Accuracy:up125bGateEval(g,s.name,s.names).Accuracy,
			Precision:up125bGateEval(g,s.name,s.names).Precision,
			Recall:up125bGateEval(g,s.name,s.names).Recall,
			Examples:up125bGateEval(g,s.name,s.names).Examples,
		})
		for _,surface:=range up125bStoreSurfaces {
			m:=up126bSurfaceMetric(g,s.name,surface,true,s.names)
			result.SurfaceMetrics=append(result.SurfaceMetrics,m)
			margin:=m.MinProbability-0.5
			if margin<result.MinStoreMarginToThreshold { result.MinStoreMarginToThreshold=margin }
		}
		for _,surface:=range up125bNonStoreSurfaces {
			m:=up126bSurfaceMetric(g,s.name,surface,false,s.names)
			result.SurfaceMetrics=append(result.SurfaceMetrics,m)
			if m.MaxProbability>result.MaxNonStoreProbability { result.MaxNonStoreProbability=m.MaxProbability }
		}
	}
	return result,nil
}
