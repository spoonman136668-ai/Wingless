package unitary

const UP128BThreeWayGateSchema = "wingless.up128b-threeway-store-gate.v1"

type UP128BSplitMetric struct {
	Split          string  `json:"split"`
	ClassAccuracy  float64 `json:"class_accuracy"`
	StorePrecision float64 `json:"store_gate_precision"`
	StoreRecall    float64 `json:"store_gate_recall"`
	Examples       int     `json:"examples"`
}

type UP128BSurfaceMetric struct {
	Split             string  `json:"split"`
	Surface           string  `json:"surface"`
	TargetClass       string  `json:"target_class"`
	Accuracy          float64 `json:"accuracy"`
	PredictedStoreRate float64 `json:"predicted_store_rate"`
	PredictedObserveRate float64 `json:"predicted_observe_rate"`
	PredictedReportRate float64 `json:"predicted_report_rate"`
	MeanStoreProb     float64 `json:"mean_store_probability"`
	MeanObserveProb   float64 `json:"mean_observe_probability"`
	MeanReportProb    float64 `json:"mean_report_probability"`
	Examples          int     `json:"examples"`
}

type UP128BThreeWayGateResult struct {
	Schema                    string                 `json:"schema"`
	Experiment                string                 `json:"experiment"`
	SourceUP127BSeal          string                 `json:"source_up127b_seal"`
	StateDimension            int                    `json:"state_dimension"`
	Epochs                    int                    `json:"epochs"`
	LearningRate              float64                `json:"learning_rate"`
	ExplicitClassAtInference  bool                   `json:"explicit_class_at_inference"`
	ClassWeightingUsed        bool                   `json:"class_weighting_used"`
	ThresholdUsed             bool                   `json:"threshold_used"`
	ControlBalancedBinary     []UP127BSplitMetric    `json:"control_balanced_binary"`
	SplitMetrics              []UP128BSplitMetric    `json:"split_metrics"`
	SurfaceMetrics            []UP128BSurfaceMetric  `json:"surface_metrics"`
}

var up128bStore=[]string{"stores","keeps","holds","saves","archives"}
var up128bObserve=[]string{"observes","sees","notes","watches","notices"}
var up128bReport=[]string{"reports","recalls","tells","remembers","recounts"}

func up128bTrainStep(c *up97bClassifier,name,verb string,class int) {
	h:=up95bEncode(name+" "+verb)
	p:=c.probs(h)
	for k:=0;k<3;k++ {
		g:=p[k]
		if k==class { g-=1 }
		for i:=0;i<64;i++ { c.w[k][i]-=0.08*g*h[i] }
		c.b[k]-=0.08*g
	}
}

func up128bTrain()*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ {
			name:=up121bOriginalNames[ni]
			for i:=0;i<5;i++ {
				up128bTrainStep(c,name,up128bStore[i],up97bStore)
				up128bTrainStep(c,name,up128bObserve[i],up97bObserve)
				up128bTrainStep(c,name,up128bReport[i],up97bReport)
			}
		}
	}
	return c
}

func up128bSurfaces() []struct{verb string; class int} {
	out:=make([]struct{verb string; class int},0,15)
	for _,v:=range up128bStore { out=append(out,struct{verb string;class int}{v,up97bStore}) }
	for _,v:=range up128bObserve { out=append(out,struct{verb string;class int}{v,up97bObserve}) }
	for _,v:=range up128bReport { out=append(out,struct{verb string;class int}{v,up97bReport}) }
	return out
}

func up128bSplit(c *up97bClassifier,split string,names []string) UP128BSplitMetric {
	hits,total,tp,fp,fn:=0,0,0,0,0
	for _,name:=range names {
		for _,s:=range up128bSurfaces() {
			pred:=up97bArgmax(c.probs(up95bEncode(name+" "+s.verb)))
			total++
			if pred==s.class { hits++ }
			if pred==up97bStore&&s.class==up97bStore { tp++ }
			if pred==up97bStore&&s.class!=up97bStore { fp++ }
			if pred!=up97bStore&&s.class==up97bStore { fn++ }
		}
	}
	prec,rec:=1.0,1.0
	if tp+fp>0 { prec=float64(tp)/float64(tp+fp) }
	if tp+fn>0 { rec=float64(tp)/float64(tp+fn) }
	return UP128BSplitMetric{Split:split,ClassAccuracy:float64(hits)/float64(total),StorePrecision:prec,StoreRecall:rec,Examples:total}
}

func up128bSurface(c *up97bClassifier,split string,names []string,verb string,class int) UP128BSurfaceMetric {
	hits:=0
	pred:=[3]int{}
	sum:=[3]float64{}
	for _,name:=range names {
		p:=c.probs(up95bEncode(name+" "+verb))
		k:=up97bArgmax(p)
		pred[k]++
		if k==class { hits++ }
		for j:=0;j<3;j++ { sum[j]+=p[j] }
	}
	n:=float64(len(names))
	return UP128BSurfaceMetric{
		Split:split,Surface:verb,TargetClass:up97bClassName(class),
		Accuracy:float64(hits)/n,
		PredictedStoreRate:float64(pred[up97bStore])/n,
		PredictedObserveRate:float64(pred[up97bObserve])/n,
		PredictedReportRate:float64(pred[up97bReport])/n,
		MeanStoreProb:sum[up97bStore]/n,MeanObserveProb:sum[up97bObserve]/n,MeanReportProb:sum[up97bReport]/n,
		Examples:len(names),
	}
}

func up97bClassName(class int) string {
	switch class {
	case up97bStore: return "STORE"
	case up97bObserve: return "OBSERVE"
	default: return "REPORT"
	}
}

func RunUP128B()(UP128BThreeWayGateResult,error){
	c:=up128bTrain()
	result:=UP128BThreeWayGateResult{
		Schema:UP128BThreeWayGateSchema,Experiment:"UP-128B-threeway-store-gate",
		SourceUP127BSeal:"d40cf39935ae917e3d1440032c0d5923c860c2c1",
		StateDimension:64,Epochs:20,LearningRate:0.08,
		ExplicitClassAtInference:false,ClassWeightingUsed:false,ThresholdUsed:false,
		ControlBalancedBinary:[]UP127BSplitMetric{
			{Arm:"balanced_positive_weight2",Split:"train",Accuracy:0.8833333333333333,Precision:0.7407407407407407,Recall:1,Examples:60},
			{Arm:"balanced_positive_weight2",Split:"heldout",Accuracy:0.9,Precision:0.8181818181818182,Recall:0.9,Examples:30},
			{Arm:"balanced_positive_weight2",Split:"unseen",Accuracy:0.9111111111111111,Precision:0.7894736842105263,Recall:1,Examples:90},
		},
	}
	splits:=[]struct{name string;names []string}{
		{"train",up121bOriginalNames[:4]},
		{"heldout",up121bOriginalNames[4:6]},
		{"unseen",up121bUnseenNames},
	}
	for _,sp:=range splits {
		result.SplitMetrics=append(result.SplitMetrics,up128bSplit(c,sp.name,sp.names))
		for _,s:=range up128bSurfaces() {
			result.SurfaceMetrics=append(result.SurfaceMetrics,up128bSurface(c,sp.name,sp.names,s.verb,s.class))
		}
	}
	return result,nil
}
