package unitary

import "math"

const UP129BDualViewGateSchema = "wingless.up129b-dualview-store-gate.v1"

type up129bGate struct {
	w [3][128]float64
	b [3]float64
}

type UP129BSplitMetric struct {
	Split          string  `json:"split"`
	ClassAccuracy  float64 `json:"class_accuracy"`
	StorePrecision float64 `json:"store_precision"`
	StoreRecall    float64 `json:"store_recall"`
	Examples       int     `json:"examples"`
}

type UP129BSurfaceMetric struct {
	Split               string  `json:"split"`
	Surface             string  `json:"surface"`
	TargetClass         string  `json:"target_class"`
	Accuracy            float64 `json:"accuracy"`
	StoreRate           float64 `json:"store_rate"`
	ObserveRate         float64 `json:"observe_rate"`
	ReportRate          float64 `json:"report_rate"`
	MeanStoreProb       float64 `json:"mean_store_probability"`
	MeanObserveProb     float64 `json:"mean_observe_probability"`
	MeanReportProb      float64 `json:"mean_report_probability"`
	MeanStoreRawLogit   float64 `json:"mean_store_raw_view_logit_contribution"`
	MeanStoreProjLogit  float64 `json:"mean_store_projected_view_logit_contribution"`
	Examples            int     `json:"examples"`
}

type UP129BPoint struct {
	ClassOrder            string              `json:"class_order"`
	ReplayPolicy          string              `json:"replay_policy"`
	Stage                 int                 `json:"stage"`
	BaseOriginalAccuracy  float64             `json:"base_original_accuracy"`
	BaseUnseenAccuracy    float64             `json:"base_unseen_accuracy"`
	AcquiredAggregateAcc  float64             `json:"acquired_aggregate_accuracy"`
	NonStoreProjectionRate float64            `json:"non_store_projection_rate"`
	StoreMetrics          []UP124BStoreMetric `json:"store_metrics"`
}

type UP129BDualViewGateResult struct {
	Schema                   string                `json:"schema"`
	Experiment               string                `json:"experiment"`
	SourceUP128BSeal         string                `json:"source_up128b_seal"`
	StateDimension           int                   `json:"state_dimension"`
	GateInputDimension       int                   `json:"gate_input_dimension"`
	Epochs                   int                   `json:"epochs"`
	LearningRate             float64               `json:"learning_rate"`
	ExplicitClassAtInference bool                  `json:"explicit_class_at_inference"`
	ProjectorRetrained       bool                  `json:"projector_retrained"`
	AdaptiveGeometry         bool                  `json:"adaptive_geometry"`
	SplitMetrics             []UP129BSplitMetric   `json:"split_metrics"`
	SurfaceMetrics           []UP129BSurfaceMetric `json:"surface_metrics"`
	Points                   []UP129BPoint         `json:"points"`
}

func up129bViews(name,verb string,o,r [64]float64)(raw,proj [64]float64){
	raw=up95bEncode(name+" "+verb)
	proj=up124bProject(raw,o,r)
	return
}

func up129bProbs(g *up129bGate,raw,proj [64]float64)[3]float64{
	z:=[3]float64{}
	maxz:=math.Inf(-1)
	for k:=0;k<3;k++{
		z[k]=g.b[k]
		for i:=0;i<64;i++{
			z[k]+=g.w[k][i]*raw[i]
			z[k]+=g.w[k][64+i]*proj[i]
		}
		if z[k]>maxz{maxz=z[k]}
	}
	sum:=0.0
	for k:=0;k<3;k++{z[k]=math.Exp(z[k]-maxz);sum+=z[k]}
	for k:=0;k<3;k++{z[k]/=sum}
	return z
}

func up129bArgmax(p [3]float64)int{
	best:=0
	for i:=1;i<3;i++{if p[i]>p[best]{best=i}}
	return best
}

func up129bGateStep(g *up129bGate,name,verb string,class int,o,r [64]float64){
	raw,proj:=up129bViews(name,verb,o,r)
	p:=up129bProbs(g,raw,proj)
	for k:=0;k<3;k++{
		d:=p[k]
		if k==class{d-=1}
		for i:=0;i<64;i++{
			g.w[k][i]-=0.08*d*raw[i]
			g.w[k][64+i]-=0.08*d*proj[i]
		}
		g.b[k]-=0.08*d
	}
}

func up129bTrainGate(o,r [64]float64)*up129bGate{
	g:=&up129bGate{}
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<4;ni++{
			name:=up121bOriginalNames[ni]
			for i:=0;i<5;i++{
				up129bGateStep(g,name,up128bStore[i],up97bStore,o,r)
				up129bGateStep(g,name,up128bObserve[i],up97bObserve,o,r)
				up129bGateStep(g,name,up128bReport[i],up97bReport,o,r)
			}
		}
	}
	return g
}

func up129bGateClass(g *up129bGate,name,verb string,o,r [64]float64)int{
	raw,proj:=up129bViews(name,verb,o,r)
	return up129bArgmax(up129bProbs(g,raw,proj))
}

func up129bSplit(g *up129bGate,split string,names []string,o,r [64]float64) UP129BSplitMetric{
	hits,total,tp,fp,fn:=0,0,0,0,0
	for _,name:=range names{
		for _,s:=range up128bSurfaces(){
			pred:=up129bGateClass(g,name,s.verb,o,r)
			total++
			if pred==s.class{hits++}
			if pred==up97bStore&&s.class==up97bStore{tp++}
			if pred==up97bStore&&s.class!=up97bStore{fp++}
			if pred!=up97bStore&&s.class==up97bStore{fn++}
		}
	}
	prec,rec:=1.0,1.0
	if tp+fp>0{prec=float64(tp)/float64(tp+fp)}
	if tp+fn>0{rec=float64(tp)/float64(tp+fn)}
	return UP129BSplitMetric{Split:split,ClassAccuracy:float64(hits)/float64(total),StorePrecision:prec,StoreRecall:rec,Examples:total}
}

func up129bSurface(g *up129bGate,split string,names []string,verb string,class int,o,r [64]float64) UP129BSurfaceMetric{
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
	return UP129BSurfaceMetric{
		Split:split,Surface:verb,TargetClass:up97bClassName(class),Accuracy:float64(hits)/n,
		StoreRate:float64(pred[up97bStore])/n,ObserveRate:float64(pred[up97bObserve])/n,ReportRate:float64(pred[up97bReport])/n,
		MeanStoreProb:sum[up97bStore]/n,MeanObserveProb:sum[up97bObserve]/n,MeanReportProb:sum[up97bReport]/n,
		MeanStoreRawLogit:rawStore/n,MeanStoreProjLogit:projStore/n,Examples:len(names),
	}
}

func up129bEncode(g *up129bGate,name,verb string,o,r [64]float64)[64]float64{
	raw,proj:=up129bViews(name,verb,o,r)
	if up129bArgmax(up129bProbs(g,raw,proj))==up97bStore{return proj}
	return raw
}

func up129bTrainStep(c *up97bClassifier,g *up129bGate,name string,spec up106bVerbSpec,o,r [64]float64){
	h:=up129bEncode(g,name,spec.verb,o,r)
	p:=c.probs(h)
	for k:=0;k<3;k++{
		d:=p[k]
		if k==spec.class{d-=1}
		for i:=0;i<64;i++{c.w[k][i]-=0.08*d*h[i]}
		c.b[k]-=0.08*d
	}
}

func up129bTrainBase(g *up129bGate,o,r [64]float64)*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++{
		for _,name:=range up121bOriginalNames{
			for _,spec:=range up106bBase{up129bTrainStep(c,g,name,spec,o,r)}
		}
	}
	return c
}

func up129bAcquire(c *up97bClassifier,g *up129bGate,current up106bVerbSpec,previous []up106bVerbSpec,policy string,stage int,o,r [64]float64){
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<4;ni++{up129bTrainStep(c,g,up121bOriginalNames[ni],current,o,r)}
		for _,spec:=range up106bBase{up129bTrainStep(c,g,up121bOriginalNames[0],spec,o,r)}
		var seq []up109bReplayExample
		if policy=="stage3_anchor2"&&stage==3{seq=up112bExtraReplay(previous,current.class,epoch,2)}else{seq=up111bCurrentExcluded(previous,current.class,epoch)}
		for _,x:=range seq{up129bTrainStep(c,g,up121bOriginalNames[x.ni],x.spec,o,r)}
	}
}

func up129bVerbAcc(c *up97bClassifier,g *up129bGate,spec up106bVerbSpec,names []string,o,r [64]float64)float64{
	hits:=0
	for _,name:=range names{
		if up97bArgmax(c.probs(up129bEncode(g,name,spec.verb,o,r)))==spec.class{hits++}
	}
	return float64(hits)/float64(len(names))
}

func up129bBaseAcc(c *up97bClassifier,g *up129bGate,names []string,o,r [64]float64)float64{
	sum:=0.0
	for _,spec:=range up106bBase{sum+=up129bVerbAcc(c,g,spec,names,o,r)}
	return sum/float64(len(up106bBase))
}

func up129bMargin(c *up97bClassifier,g *up129bGate,verb string,names []string,o,r [64]float64)float64{
	min:=math.Inf(1)
	for _,name:=range names{
		z:=up117bLogits(c,up129bEncode(g,name,verb,o,r))
		m:=z[up97bStore]-math.Max(z[up97bObserve],z[up97bReport])
		if m<min{min=m}
	}
	return min
}

func up129bNonStoreProjectionRate(g *up129bGate,names []string,o,r [64]float64)float64{
	projected,total:=0,0
	for _,name:=range names{
		for _,verb:=range append(append([]string{},up128bObserve...),up128bReport...){
			total++
			if up129bGateClass(g,name,verb,o,r)==up97bStore{projected++}
		}
	}
	return float64(projected)/float64(total)
}

func up129bPoint(order,policy string,stage int,c *up97bClassifier,g *up129bGate,acquired []up106bVerbSpec,o,r [64]float64) UP129BPoint{
	p:=UP129BPoint{
		ClassOrder:order,ReplayPolicy:policy,Stage:stage,
		BaseOriginalAccuracy:up129bBaseAcc(c,g,up121bOriginalNames,o,r),
		BaseUnseenAccuracy:up129bBaseAcc(c,g,up121bUnseenNames,o,r),
		NonStoreProjectionRate:up129bNonStoreProjectionRate(g,up121bUnseenNames,o,r),
	}
	if len(acquired)>0{
		s:=0.0
		for _,spec:=range acquired{s+=up129bVerbAcc(c,g,spec,up121bUnseenNames[4:6],o,r)}
		p.AcquiredAggregateAcc=s/float64(len(acquired))
	}
	for _,verb:=range []string{"stores","archives","keepsafe"}{
		spec:=up106bVerbSpec{verb:verb,class:up97bStore}
		p.StoreMetrics=append(p.StoreMetrics,UP124BStoreMetric{
			Verb:verb,
			OriginalAcc:up129bVerbAcc(c,g,spec,up121bOriginalNames,o,r),
			UnseenAcc:up129bVerbAcc(c,g,spec,up121bUnseenNames,o,r),
			OriginalMargin:up129bMargin(c,g,verb,up121bOriginalNames,o,r),
			UnseenMargin:up129bMargin(c,g,verb,up121bUnseenNames,o,r),
		})
	}
	return p
}

func RunUP129B()(UP129BDualViewGateResult,error){
	o,r:=up124bCompetitorDirections()
	g:=up129bTrainGate(o,r)
	result:=UP129BDualViewGateResult{
		Schema:UP129BDualViewGateSchema,Experiment:"UP-129B-dualview-store-gate",
		SourceUP128BSeal:"34240b7e3da4c1de6e8fb9011bb635c6aa52a8dc",
		StateDimension:64,GateInputDimension:128,Epochs:20,LearningRate:0.08,
		ExplicitClassAtInference:false,ProjectorRetrained:false,AdaptiveGeometry:false,
	}
	splits:=[]struct{name string;names []string}{
		{"train",up121bOriginalNames[:4]},{"heldout",up121bOriginalNames[4:6]},{"unseen",up121bUnseenNames},
	}
	for _,sp:=range splits{
		result.SplitMetrics=append(result.SplitMetrics,up129bSplit(g,sp.name,sp.names,o,r))
		for _,s:=range up128bSurfaces(){result.SurfaceMetrics=append(result.SurfaceMetrics,up129bSurface(g,sp.name,sp.names,s.verb,s.class,o,r))}
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,order:=range orders{
		name:=up114bOrderName(order)
		seq:=up121bSequence(order,"alternate")
		for _,policy:=range []string{"current_class_excluded","stage3_anchor2"}{
			c:=up129bTrainBase(g,o,r)
			acquired:=[]up106bVerbSpec{}
			result.Points=append(result.Points,up129bPoint(name,policy,0,c,g,acquired,o,r))
			for idx,spec:=range seq{
				stage:=idx+1
				up129bAcquire(c,g,spec,acquired,policy,stage,o,r)
				acquired=append(acquired,spec)
				result.Points=append(result.Points,up129bPoint(name,policy,stage,c,g,acquired,o,r))
			}
		}
	}
	return result,nil
}
