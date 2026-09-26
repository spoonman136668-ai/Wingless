package unitary

import "math"

const UPLM1KRouterTopologySchema = "wingless.up-lm1k-fourth-router-topology.v1"

type UPLM1KFamilyMetric struct {
	Family        string  `json:"family"`
	Split         string  `json:"split"`
	Accuracy      float64 `json:"accuracy"`
	StoreRecall   float64 `json:"store_recall"`
	ObserveRecall float64 `json:"observe_recall"`
	ReportRecall  float64 `json:"report_recall"`
	Examples      int     `json:"examples"`
}

type UPLM1KSurfaceMetric struct {
	Family               string  `json:"family"`
	Split                string  `json:"split"`
	Surface              string  `json:"surface"`
	TargetClass          string  `json:"target_class"`
	Accuracy             float64 `json:"accuracy"`
	PredictedStoreRate   float64 `json:"predicted_store_rate"`
	PredictedObserveRate float64 `json:"predicted_observe_rate"`
	PredictedReportRate  float64 `json:"predicted_report_rate"`
	MeanStoreProb        float64 `json:"mean_store_probability"`
	MeanObserveProb      float64 `json:"mean_observe_probability"`
	MeanReportProb       float64 `json:"mean_report_probability"`
	MinTargetMargin      float64 `json:"min_target_margin"`
	Examples             int     `json:"examples"`
}

type UPLM1KRouterTopologyResult struct {
	Schema                  string                 `json:"schema"`
	Experiment              string                 `json:"experiment"`
	SourceUPLM1JSeal        string                 `json:"source_up_lm1j_seal"`
	StateDimension          int                    `json:"state_dimension"`
	FourthGroundingEpochs   int                    `json:"fourth_grounding_epochs"`
	LearningRate            float64                `json:"learning_rate"`
	RouterReconstructed     bool                   `json:"router_reconstructed"`
	RouterModifiedAfterBuild bool                  `json:"router_modified_after_build"`
	ByteModelUsed           bool                   `json:"byte_model_used"`
	ThresholdChanged        bool                   `json:"threshold_changed"`
	FamilyMetrics           []UPLM1KFamilyMetric  `json:"family_metrics"`
	SurfaceMetrics          []UPLM1KSurfaceMetric `json:"surface_metrics"`
}

var uplm1kUnseenNames=[]string{"gia","hal","ivy","jon","kia","leo"}

func uplm1kSurface(c *uplm0jClassifier,family,split string,names []string,surface string,class int) UPLM1KSurfaceMetric {
	hits:=0
	pred:=[3]int{}
	sum:=[3]float64{}
	minMargin:=math.Inf(1)
	for _,name:=range names {
		p:=c.probs(uplm0fEncode(name+" "+surface))
		k:=uplm0jArgmax(p)
		pred[k]++
		if k==class { hits++ }
		for j:=0;j<3;j++ { sum[j]+=p[j] }
		bestOther:=math.Inf(-1)
		for j:=0;j<3;j++ {
			if j!=class && p[j]>bestOther { bestOther=p[j] }
		}
		margin:=p[class]-bestOther
		if margin<minMargin { minMargin=margin }
	}
	n:=float64(len(names))
	return UPLM1KSurfaceMetric{
		Family:family,Split:split,Surface:surface,TargetClass:uplm0jClassLabel(class),
		Accuracy:float64(hits)/n,
		PredictedStoreRate:float64(pred[uplm0jStore])/n,
		PredictedObserveRate:float64(pred[uplm0jObserve])/n,
		PredictedReportRate:float64(pred[uplm0jReport])/n,
		MeanStoreProb:sum[uplm0jStore]/n,MeanObserveProb:sum[uplm0jObserve]/n,MeanReportProb:sum[uplm0jReport]/n,
		MinTargetMargin:minMargin,Examples:len(names),
	}
}

func uplm1kFamily(c *uplm0jClassifier,family,split string,names,surfaces []string) UPLM1KFamilyMetric {
	hits,total:=0,0
	classHits:=[3]int{}
	classTotal:=[3]int{}
	for _,name:=range names {
		for class,surface:=range surfaces {
			pred:=uplm0jArgmax(c.probs(uplm0fEncode(name+" "+surface)))
			total++;classTotal[class]++
			if pred==class { hits++;classHits[class]++ }
		}
	}
	return UPLM1KFamilyMetric{
		Family:family,Split:split,Accuracy:float64(hits)/float64(total),
		StoreRecall:float64(classHits[0])/float64(classTotal[0]),
		ObserveRecall:float64(classHits[1])/float64(classTotal[1]),
		ReportRecall:float64(classHits[2])/float64(classTotal[2]),
		Examples:total,
	}
}

func RunUPLM1K()(UPLM1KRouterTopologyResult,error){
	c:=uplm1jTrainClassifier()
	result:=UPLM1KRouterTopologyResult{
		Schema:UPLM1KRouterTopologySchema,Experiment:"UP-LM1K-fourth-router-topology",
		SourceUPLM1JSeal:"077930c37f583288be2c56ea157d797ad41207a7",
		StateDimension:64,FourthGroundingEpochs:20,LearningRate:0.08,
		RouterReconstructed:true,RouterModifiedAfterBuild:false,ByteModelUsed:false,ThresholdChanged:false,
	}
	type cell struct{family,split string;names,surfaces []string}
	names:=uplm0gNames()
	cells:=[]cell{
		{"fourth","train",names[:4],[]string{"retains","inspects","states"}},
		{"fourth","heldout",names[4:6],[]string{"retains","inspects","states"}},
		{"fourth","unseen",uplm1kUnseenNames,[]string{"retains","inspects","states"}},
		{"base","heldout",names[4:6],[]string{"stores","observes","reports"}},
		{"paraphrase","heldout",names[4:6],[]string{"saves","sees","recalls"}},
		{"third","heldout",names[4:6],[]string{"archives","notices","recounts"}},
	}
	for _,x:=range cells {
		result.FamilyMetrics=append(result.FamilyMetrics,uplm1kFamily(c,x.family,x.split,x.names,x.surfaces))
		for class,surface:=range x.surfaces {
			result.SurfaceMetrics=append(result.SurfaceMetrics,uplm1kSurface(c,x.family,x.split,x.names,surface,class))
		}
	}
	return result,nil
}
