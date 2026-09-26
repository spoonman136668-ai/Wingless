package unitary

import "math"

const UPLM1LStatesOrthogonalSchema = "wingless.up-lm1l-states-store-orthogonal.v1"

func uplm1lDot(a,b [64]float64) float64 {
	s:=0.0
	for i:=0;i<64;i++ { s+=a[i]*b[i] }
	return s
}

func uplm1lNorm(a [64]float64) float64 {
	return math.Sqrt(uplm1lDot(a,a))
}

type UPLM1LFamilyMetric struct {
	Arm           string  `json:"arm"`
	Family        string  `json:"family"`
	Split         string  `json:"split"`
	Accuracy      float64 `json:"accuracy"`
	StoreRecall   float64 `json:"store_recall"`
	ObserveRecall float64 `json:"observe_recall"`
	ReportRecall  float64 `json:"report_recall"`
	Examples      int     `json:"examples"`
}

type UPLM1LSurfaceMetric struct {
	Arm                 string  `json:"arm"`
	Family              string  `json:"family"`
	Split               string  `json:"split"`
	Surface             string  `json:"surface"`
	TargetClass         string  `json:"target_class"`
	Accuracy            float64 `json:"accuracy"`
	PredictedStoreRate  float64 `json:"predicted_store_rate"`
	PredictedObserveRate float64 `json:"predicted_observe_rate"`
	PredictedReportRate float64 `json:"predicted_report_rate"`
	MeanStoreProb       float64 `json:"mean_store_probability"`
	MeanObserveProb     float64 `json:"mean_observe_probability"`
	MeanReportProb      float64 `json:"mean_report_probability"`
	MinTargetMargin     float64 `json:"min_target_margin"`
	Examples            int     `json:"examples"`
}

type UPLM1LStatesOrthogonalResult struct {
	Schema                    string                 `json:"schema"`
	Experiment                string                 `json:"experiment"`
	SourceUPLM1KSeal          string                 `json:"source_up_lm1k_seal"`
	StateDimension            int                    `json:"state_dimension"`
	FourthGroundingEpochs     int                    `json:"fourth_grounding_epochs"`
	LearningRate              float64                `json:"learning_rate"`
	CorrectionRecomputed      bool                   `json:"correction_recomputed"`
	ThresholdChanged          bool                   `json:"threshold_changed"`
	ByteModelUsed             bool                   `json:"byte_model_used"`
	StoreDirectionNorm        float64                `json:"store_direction_norm"`
	FamilyMetrics             []UPLM1LFamilyMetric   `json:"family_metrics"`
	SurfaceMetrics            []UPLM1LSurfaceMetric  `json:"surface_metrics"`
}

func uplm1lStoreDirection() [64]float64 {
	names:=uplm0gNames()
	surfaces:=[]string{"stores","saves","archives","retains"}
	var d [64]float64
	count:=0.0
	for ni:=0;ni<4;ni++{
		for _,surface:=range surfaces{
			h:=uplm0fEncode(names[ni]+" "+surface)
			for i:=0;i<64;i++{d[i]+=h[i]}
			count++
		}
	}
	for i:=0;i<64;i++{d[i]/=count}
	n:=uplm1lNorm(d)
	if n>1e-12{for i:=0;i<64;i++{d[i]/=n}}
	return d
}

func uplm1lEncode(arm,name,surface string,d [64]float64)[64]float64{
	h:=uplm0fEncode(name+" "+surface)
	if arm!="states_store_orthogonal" || surface!="states" { return h }
	dot:=uplm1lDot(h,d)
	var residual [64]float64
	for i:=0;i<64;i++{residual[i]=h[i]-dot*d[i]}
	nh,nr:=uplm1lNorm(h),uplm1lNorm(residual)
	if nr>1e-12{
		s:=nh/nr
		for i:=0;i<64;i++{residual[i]*=s}
	}
	return residual
}

func uplm1lTrainStep(c *uplm0jClassifier,arm,name,surface string,class int,d [64]float64){
	h:=uplm1lEncode(arm,name,surface,d)
	p:=c.probs(h)
	for k:=0;k<3;k++{
		g:=p[k]
		if k==class{g-=1}
		for i:=0;i<64;i++{c.w[k][i]-=0.08*g*h[i]}
		c.b[k]-=0.08*g
	}
}

func uplm1lTrainClassifier(arm string,d [64]float64)*uplm0jClassifier{
	if arm=="raw_control" { return uplm1jTrainClassifier() }
	c:=uplm1cTrainClassifier()
	names:=uplm0gNames()
	prior:=[][]string{
		{"stores","observes","reports"},
		{"saves","sees","recalls"},
		{"archives","notices","recounts"},
	}
	fourth:=[]string{"retains","inspects","states"}
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<4;ni++{
			for class,surface:=range fourth{
				uplm1lTrainStep(c,arm,names[ni],surface,class,d)
			}
		}
		for _,family:=range prior{
			for class,surface:=range family{
				uplm0nTrainStep(c,uplm0nAnchor(names[0],surface),class)
			}
		}
	}
	return c
}

func uplm1lProbs(c *uplm0jClassifier,arm,name,surface string,d [64]float64)[3]float64{
	return c.probs(uplm1lEncode(arm,name,surface,d))
}

func uplm1lFamily(c *uplm0jClassifier,arm,family,split string,names,surfaces []string,d [64]float64) UPLM1LFamilyMetric {
	hits,total:=0,0
	classHits:=[3]int{}
	classTotal:=[3]int{}
	for _,name:=range names{
		for class,surface:=range surfaces{
			pred:=uplm0jArgmax(uplm1lProbs(c,arm,name,surface,d))
			total++;classTotal[class]++
			if pred==class{hits++;classHits[class]++}
		}
	}
	return UPLM1LFamilyMetric{
		Arm:arm,Family:family,Split:split,Accuracy:float64(hits)/float64(total),
		StoreRecall:float64(classHits[0])/float64(classTotal[0]),
		ObserveRecall:float64(classHits[1])/float64(classTotal[1]),
		ReportRecall:float64(classHits[2])/float64(classTotal[2]),Examples:total,
	}
}

func uplm1lSurface(c *uplm0jClassifier,arm,family,split string,names []string,surface string,class int,d [64]float64) UPLM1LSurfaceMetric {
	hits:=0
	pred:=[3]int{}
	sum:=[3]float64{}
	minMargin:=math.Inf(1)
	for _,name:=range names{
		p:=uplm1lProbs(c,arm,name,surface,d)
		k:=uplm0jArgmax(p)
		pred[k]++
		if k==class{hits++}
		for j:=0;j<3;j++{sum[j]+=p[j]}
		bestOther:=math.Inf(-1)
		for j:=0;j<3;j++{if j!=class&&p[j]>bestOther{bestOther=p[j]}}
		margin:=p[class]-bestOther
		if margin<minMargin{minMargin=margin}
	}
	n:=float64(len(names))
	return UPLM1LSurfaceMetric{
		Arm:arm,Family:family,Split:split,Surface:surface,TargetClass:uplm0jClassLabel(class),
		Accuracy:float64(hits)/n,
		PredictedStoreRate:float64(pred[uplm0jStore])/n,
		PredictedObserveRate:float64(pred[uplm0jObserve])/n,
		PredictedReportRate:float64(pred[uplm0jReport])/n,
		MeanStoreProb:sum[uplm0jStore]/n,MeanObserveProb:sum[uplm0jObserve]/n,MeanReportProb:sum[uplm0jReport]/n,
		MinTargetMargin:minMargin,Examples:len(names),
	}
}

func RunUPLM1L()(UPLM1LStatesOrthogonalResult,error){
	d:=uplm1lStoreDirection()
	result:=UPLM1LStatesOrthogonalResult{
		Schema:UPLM1LStatesOrthogonalSchema,Experiment:"UP-LM1L-states-store-orthogonal",
		SourceUPLM1KSeal:"e5944e3a6a52a1e4bd98c5c880c5eca964075fd5",
		StateDimension:64,FourthGroundingEpochs:20,LearningRate:0.08,
		CorrectionRecomputed:false,ThresholdChanged:false,ByteModelUsed:false,StoreDirectionNorm:uplm1lNorm(d),
	}
	names:=uplm0gNames()
	type cell struct{family,split string;names,surfaces []string}
	cells:=[]cell{
		{"fourth","train",names[:4],[]string{"retains","inspects","states"}},
		{"fourth","heldout",names[4:6],[]string{"retains","inspects","states"}},
		{"fourth","unseen",uplm1kUnseenNames,[]string{"retains","inspects","states"}},
		{"base","heldout",names[4:6],[]string{"stores","observes","reports"}},
		{"paraphrase","heldout",names[4:6],[]string{"saves","sees","recalls"}},
		{"third","heldout",names[4:6],[]string{"archives","notices","recounts"}},
	}
	for _,arm:=range []string{"raw_control","states_store_orthogonal"}{
		c:=uplm1lTrainClassifier(arm,d)
		for _,x:=range cells{
			result.FamilyMetrics=append(result.FamilyMetrics,uplm1lFamily(c,arm,x.family,x.split,x.names,x.surfaces,d))
			for class,surface:=range x.surfaces{
				result.SurfaceMetrics=append(result.SurfaceMetrics,uplm1lSurface(c,arm,x.family,x.split,x.names,surface,class,d))
			}
		}
	}
	return result,nil
}
