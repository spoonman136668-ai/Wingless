package unitary

const UP130BZeroShotLexicalSchema = "wingless.up130b-dualview-zero-shot-lexical.v1"

type UP130BSplitMetric struct {
	Split            string  `json:"split"`
	OverallAccuracy  float64 `json:"overall_accuracy"`
	StoreAccuracy    float64 `json:"store_accuracy"`
	ObserveAccuracy  float64 `json:"observe_accuracy"`
	ReportAccuracy   float64 `json:"report_accuracy"`
	StorePrecision   float64 `json:"store_precision"`
	StoreRecall      float64 `json:"store_recall"`
	WorstSurfaceAcc  float64 `json:"worst_surface_accuracy"`
	Examples         int     `json:"examples"`
}

type UP130BSurfaceMetric struct {
	Split              string  `json:"split"`
	Surface            string  `json:"surface"`
	TargetClass        string  `json:"target_class"`
	Accuracy           float64 `json:"accuracy"`
	StoreRate          float64 `json:"store_rate"`
	ObserveRate        float64 `json:"observe_rate"`
	ReportRate         float64 `json:"report_rate"`
	MeanStoreProb      float64 `json:"mean_store_probability"`
	MeanObserveProb    float64 `json:"mean_observe_probability"`
	MeanReportProb     float64 `json:"mean_report_probability"`
	MeanStoreRawLogit  float64 `json:"mean_store_raw_view_logit_contribution"`
	MeanStoreProjLogit float64 `json:"mean_store_projected_view_logit_contribution"`
	Examples           int     `json:"examples"`
}

type UP130BZeroShotLexicalResult struct {
	Schema                    string                `json:"schema"`
	Experiment                string                `json:"experiment"`
	SourceUP129BSeal          string                `json:"source_up129b_seal"`
	StateDimension            int                   `json:"state_dimension"`
	GateInputDimension        int                   `json:"gate_input_dimension"`
	GateEpochs                int                   `json:"gate_epochs"`
	LearningRate              float64               `json:"learning_rate"`
	NewSurfaceTrainingUsed    bool                  `json:"new_surface_training_used"`
	NewSubjectTrainingUsed    bool                  `json:"new_subject_training_used"`
	ProjectorRecomputed       bool                  `json:"projector_recomputed"`
	ExplicitClassAtInference  bool                  `json:"explicit_class_at_inference"`
	SplitMetrics              []UP130BSplitMetric   `json:"split_metrics"`
	SurfaceMetrics            []UP130BSurfaceMetric `json:"surface_metrics"`
}

var up130bStore=[]string{"lodges","stashes","caches","files"}
var up130bObserve=[]string{"scans","checks","views","monitors"}
var up130bReport=[]string{"relays","announces","cites","summarizes"}
var up130bNewNames=[]string{"mia","noah","opal","pax","quin","rue"}

func up130bSurfaces() []struct{verb string; class int}{
	out:=make([]struct{verb string;class int},0,12)
	for _,v:=range up130bStore{out=append(out,struct{verb string;class int}{v,up97bStore})}
	for _,v:=range up130bObserve{out=append(out,struct{verb string;class int}{v,up97bObserve})}
	for _,v:=range up130bReport{out=append(out,struct{verb string;class int}{v,up97bReport})}
	return out
}

func up130bSurface(g *up129bGate,split string,names []string,verb string,class int,o,r [64]float64) UP130BSurfaceMetric{
	hits:=0
	pred:=[3]int{}
	sum:=[3]float64{}
	rawStore,projStore:=0.0,0.0
	for _,name:=range names{
		raw,proj:=up129bViews(name,verb,o,r)
		p:=up129bProbs(g,raw,proj)
		k:=up129bArgmax(p)
		pred[k]++
		if k==class{hits++}
		for j:=0;j<3;j++{sum[j]+=p[j]}
		for i:=0;i<64;i++{
			rawStore+=g.w[up97bStore][i]*raw[i]
			projStore+=g.w[up97bStore][64+i]*proj[i]
		}
	}
	n:=float64(len(names))
	return UP130BSurfaceMetric{
		Split:split,Surface:verb,TargetClass:up97bClassName(class),Accuracy:float64(hits)/n,
		StoreRate:float64(pred[up97bStore])/n,ObserveRate:float64(pred[up97bObserve])/n,ReportRate:float64(pred[up97bReport])/n,
		MeanStoreProb:sum[up97bStore]/n,MeanObserveProb:sum[up97bObserve]/n,MeanReportProb:sum[up97bReport]/n,
		MeanStoreRawLogit:rawStore/n,MeanStoreProjLogit:projStore/n,Examples:len(names),
	}
}

func up130bSplit(g *up129bGate,split string,names []string,o,r [64]float64)(UP130BSplitMetric,[]UP130BSurfaceMetric){
	hits,total,tp,fp,fn:=0,0,0,0,0
	classHits:=[3]int{}
	classTotal:=[3]int{}
	worst:=1.0
	surfaces:=[]UP130BSurfaceMetric{}
	for _,s:=range up130bSurfaces(){
		m:=up130bSurface(g,split,names,s.verb,s.class,o,r)
		surfaces=append(surfaces,m)
		if m.Accuracy<worst{worst=m.Accuracy}
		for _,name:=range names{
			k:=up129bGateClass(g,name,s.verb,o,r)
			total++;classTotal[s.class]++
			if k==s.class{hits++;classHits[s.class]++}
			if k==up97bStore&&s.class==up97bStore{tp++}
			if k==up97bStore&&s.class!=up97bStore{fp++}
			if k!=up97bStore&&s.class==up97bStore{fn++}
		}
	}
	prec,rec:=1.0,1.0
	if tp+fp>0{prec=float64(tp)/float64(tp+fp)}
	if tp+fn>0{rec=float64(tp)/float64(tp+fn)}
	return UP130BSplitMetric{
		Split:split,OverallAccuracy:float64(hits)/float64(total),
		StoreAccuracy:float64(classHits[up97bStore])/float64(classTotal[up97bStore]),
		ObserveAccuracy:float64(classHits[up97bObserve])/float64(classTotal[up97bObserve]),
		ReportAccuracy:float64(classHits[up97bReport])/float64(classTotal[up97bReport]),
		StorePrecision:prec,StoreRecall:rec,WorstSurfaceAcc:worst,Examples:total,
	},surfaces
}

func RunUP130B()(UP130BZeroShotLexicalResult,error){
	o,r:=up124bCompetitorDirections()
	g:=up129bTrainGate(o,r)
	result:=UP130BZeroShotLexicalResult{
		Schema:UP130BZeroShotLexicalSchema,Experiment:"UP-130B-dualview-zero-shot-lexical",
		SourceUP129BSeal:"3a69812d126a79dc509b2eaf5825173012472e35",
		StateDimension:64,GateInputDimension:128,GateEpochs:20,LearningRate:0.08,
		NewSurfaceTrainingUsed:false,NewSubjectTrainingUsed:false,ProjectorRecomputed:false,ExplicitClassAtInference:false,
	}
	splits:=[]struct{name string;names []string}{
		{"original_heldout",up121bOriginalNames[4:6]},
		{"prior_unseen",up121bUnseenNames},
		{"new_unseen",up130bNewNames},
	}
	for _,sp:=range splits{
		m,s:=up130bSplit(g,sp.name,sp.names,o,r)
		result.SplitMetrics=append(result.SplitMetrics,m)
		result.SurfaceMetrics=append(result.SurfaceMetrics,s...)
	}
	return result,nil
}
