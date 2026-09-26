package unitary

import "math"

const UP125BLearnedStoreGateSchema = "wingless.up125b-learned-store-gate.v1"

type up125bGate struct {
	w [64]float64
	b float64
}

type UP125BGateMetric struct {
	Split     string  `json:"split"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Examples  int     `json:"examples"`
}

type UP125BPoint struct {
	Arm                  string              `json:"arm"`
	ClassOrder           string              `json:"class_order"`
	ReplayPolicy         string              `json:"replay_policy"`
	Stage                int                 `json:"stage"`
	BaseOriginalAccuracy float64             `json:"base_original_accuracy"`
	BaseUnseenAccuracy   float64             `json:"base_unseen_accuracy"`
	AcquiredAggregateAcc float64             `json:"acquired_aggregate_accuracy"`
	StoreMetrics         []UP124BStoreMetric `json:"store_metrics"`
	NonStoreErrorRate    float64             `json:"non_store_error_rate"`
}

type UP125BLearnedStoreGateResult struct {
	Schema                    string             `json:"schema"`
	Experiment                string             `json:"experiment"`
	SourceUP124BSeal          string             `json:"source_up124b_seal"`
	StateDimension            int                `json:"state_dimension"`
	GateEpochs                int                `json:"gate_epochs"`
	GateLearningRate          float64            `json:"gate_learning_rate"`
	GateThreshold             float64            `json:"gate_threshold"`
	ExplicitClassAtInference  bool               `json:"explicit_class_at_inference"`
	ProjectorRetrained        bool               `json:"projector_retrained"`
	AdaptiveGeometry          bool               `json:"adaptive_geometry"`
	GateMetrics               []UP125BGateMetric `json:"gate_metrics"`
	Points                    []UP125BPoint      `json:"points"`
}

var up125bStoreSurfaces=[]string{"stores","keeps","holds","saves","archives"}
var up125bNonStoreSurfaces=[]string{"observes","sees","notes","watches","notices","reports","recalls","tells","remembers","recounts"}

func up125bSigmoid(x float64) float64 {
	if x>=0 { z:=math.Exp(-x); return 1/(1+z) }
	z:=math.Exp(x); return z/(1+z)
}

func up125bGateProb(g *up125bGate,h [64]float64) float64 {
	z:=g.b
	for i:=0;i<64;i++ { z+=g.w[i]*h[i] }
	return up125bSigmoid(z)
}

func up125bGateStep(g *up125bGate,name,verb string,target float64) {
	h:=up95bEncode(name+" "+verb)
	p:=up125bGateProb(g,h)
	d:=p-target
	for i:=0;i<64;i++ { g.w[i]-=0.08*d*h[i] }
	g.b-=0.08*d
}

func up125bTrainGate()*up125bGate{
	g:=&up125bGate{}
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ {
			name:=up121bOriginalNames[ni]
			for _,verb:=range up125bStoreSurfaces { up125bGateStep(g,name,verb,1) }
			for _,verb:=range up125bNonStoreSurfaces { up125bGateStep(g,name,verb,0) }
		}
	}
	return g
}

func up125bGatePredict(g *up125bGate,name,verb string) bool {
	return up125bGateProb(g,up95bEncode(name+" "+verb))>=0.5
}

func up125bGateEval(g *up125bGate,split string,names []string) UP125BGateMetric {
	tp,fp,fn,hits,total:=0,0,0,0,0
	for _,name:=range names {
		for _,verb:=range up125bStoreSurfaces {
			p:=up125bGatePredict(g,name,verb); total++
			if p { hits++;tp++ } else { fn++ }
		}
		for _,verb:=range up125bNonStoreSurfaces {
			p:=up125bGatePredict(g,name,verb); total++
			if !p { hits++ } else { fp++ }
		}
	}
	prec,rec:=1.0,1.0
	if tp+fp>0 { prec=float64(tp)/float64(tp+fp) }
	if tp+fn>0 { rec=float64(tp)/float64(tp+fn) }
	return UP125BGateMetric{Split:split,Accuracy:float64(hits)/float64(total),Precision:prec,Recall:rec,Examples:total}
}

func up125bEncode(arm string,g *up125bGate,name,verb string,o,r [64]float64)[64]float64{
	h:=up95bEncode(name+" "+verb)
	apply:=up124bIsStoreVerb(verb)
	if arm=="learned_store_gate" { apply=up125bGatePredict(g,name,verb) }
	if apply { h=up124bProject(h,o,r) }
	return h
}

func up125bTrainStep(c *up97bClassifier,arm string,g *up125bGate,name string,spec up106bVerbSpec,o,r [64]float64){
	h:=up125bEncode(arm,g,name,spec.verb,o,r)
	p:=c.probs(h)
	for k:=0;k<3;k++ {
		d:=p[k]
		if k==spec.class { d-=1 }
		for i:=0;i<64;i++ { c.w[k][i]-=0.08*d*h[i] }
		c.b[k]-=0.08*d
	}
}

func up125bTrainBase(arm string,g *up125bGate,o,r [64]float64)*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++ {
		for _,name:=range up121bOriginalNames {
			for _,spec:=range up106bBase { up125bTrainStep(c,arm,g,name,spec,o,r) }
		}
	}
	return c
}

func up125bAcquire(c *up97bClassifier,arm string,g *up125bGate,current up106bVerbSpec,previous []up106bVerbSpec,policy string,stage int,o,r [64]float64){
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ { up125bTrainStep(c,arm,g,up121bOriginalNames[ni],current,o,r) }
		for _,spec:=range up106bBase { up125bTrainStep(c,arm,g,up121bOriginalNames[0],spec,o,r) }
		var seq []up109bReplayExample
		if policy=="stage3_anchor2"&&stage==3 { seq=up112bExtraReplay(previous,current.class,epoch,2) } else { seq=up111bCurrentExcluded(previous,current.class,epoch) }
		for _,x:=range seq { up125bTrainStep(c,arm,g,up121bOriginalNames[x.ni],x.spec,o,r) }
	}
}

func up125bVerbAcc(c *up97bClassifier,arm string,g *up125bGate,spec up106bVerbSpec,names []string,o,r [64]float64)float64{
	hits:=0
	for _,name:=range names {
		if up97bArgmax(c.probs(up125bEncode(arm,g,name,spec.verb,o,r)))==spec.class { hits++ }
	}
	return float64(hits)/float64(len(names))
}

func up125bMargin(c *up97bClassifier,arm string,g *up125bGate,verb string,names []string,o,r [64]float64)float64{
	min:=math.Inf(1)
	for _,name:=range names {
		z:=up117bLogits(c,up125bEncode(arm,g,name,verb,o,r))
		m:=z[up97bStore]-math.Max(z[up97bObserve],z[up97bReport])
		if m<min { min=m }
	}
	return min
}

func up125bBaseAcc(c *up97bClassifier,arm string,g *up125bGate,names []string,o,r [64]float64)float64{
	s:=0.0
	for _,spec:=range up106bBase { s+=up125bVerbAcc(c,arm,g,spec,names,o,r) }
	return s/float64(len(up106bBase))
}

func up125bNonStoreError(g *up125bGate,names []string)float64{
	errs,total:=0,0
	for _,name:=range names {
		for _,verb:=range up125bNonStoreSurfaces {
			total++
			if up125bGatePredict(g,name,verb) { errs++ }
		}
	}
	return float64(errs)/float64(total)
}

func up125bPoint(arm,order,policy string,stage int,c *up97bClassifier,g *up125bGate,acquired []up106bVerbSpec,o,r [64]float64) UP125BPoint {
	p:=UP125BPoint{
		Arm:arm,ClassOrder:order,ReplayPolicy:policy,Stage:stage,
		BaseOriginalAccuracy:up125bBaseAcc(c,arm,g,up121bOriginalNames,o,r),
		BaseUnseenAccuracy:up125bBaseAcc(c,arm,g,up121bUnseenNames,o,r),
		NonStoreErrorRate:up125bNonStoreError(g,up121bUnseenNames),
	}
	if len(acquired)>0 {
		s:=0.0
		for _,spec:=range acquired { s+=up125bVerbAcc(c,arm,g,spec,up121bUnseenNames[4:6],o,r) }
		p.AcquiredAggregateAcc=s/float64(len(acquired))
	}
	for _,verb:=range []string{"stores","keeps","keepsafe","archives"} {
		spec:=up106bVerbSpec{verb:verb,class:up97bStore}
		p.StoreMetrics=append(p.StoreMetrics,UP124BStoreMetric{
			Verb:verb,
			OriginalAcc:up125bVerbAcc(c,arm,g,spec,up121bOriginalNames,o,r),
			UnseenAcc:up125bVerbAcc(c,arm,g,spec,up121bUnseenNames,o,r),
			OriginalMargin:up125bMargin(c,arm,g,verb,up121bOriginalNames,o,r),
			UnseenMargin:up125bMargin(c,arm,g,verb,up121bUnseenNames,o,r),
		})
	}
	return p
}

func RunUP125B()(UP125BLearnedStoreGateResult,error){
	o,r:=up124bCompetitorDirections()
	g:=up125bTrainGate()
	result:=UP125BLearnedStoreGateResult{
		Schema:UP125BLearnedStoreGateSchema,Experiment:"UP-125B-learned-store-gate",
		SourceUP124BSeal:"adf0ead18133ee4390a309ed114440e27847ddd4",
		StateDimension:64,GateEpochs:20,GateLearningRate:0.08,GateThreshold:0.5,
		ExplicitClassAtInference:false,ProjectorRetrained:false,AdaptiveGeometry:false,
	}
	result.GateMetrics=[]UP125BGateMetric{
		up125bGateEval(g,"heldout_subjects",up121bOriginalNames[4:6]),
		up125bGateEval(g,"unseen_subjects",up121bUnseenNames),
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,arm:=range []string{"explicit_store_gate","learned_store_gate"} {
		for _,order:=range orders {
			name:=up114bOrderName(order)
			seq:=up121bSequence(order,"alternate")
			for _,policy:=range []string{"current_class_excluded","stage3_anchor2"} {
				c:=up125bTrainBase(arm,g,o,r)
				acquired:=[]up106bVerbSpec{}
				result.Points=append(result.Points,up125bPoint(arm,name,policy,0,c,g,acquired,o,r))
				for idx,spec:=range seq {
					stage:=idx+1
					up125bAcquire(c,arm,g,spec,acquired,policy,stage,o,r)
					acquired=append(acquired,spec)
					result.Points=append(result.Points,up125bPoint(arm,name,policy,stage,c,g,acquired,o,r))
				}
			}
		}
	}
	return result,nil
}
