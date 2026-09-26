package unitary

import "math"

const UP124BClasswideProjectorSchema = "wingless.up124b-classwide-store-projector.v1"

type UP124BStoreMetric struct {
	Verb           string  `json:"verb"`
	OriginalAcc    float64 `json:"original_accuracy"`
	UnseenAcc      float64 `json:"unseen_accuracy"`
	OriginalMargin float64 `json:"original_min_margin"`
	UnseenMargin   float64 `json:"unseen_min_margin"`
}

type UP124BPoint struct {
	Arm                   string              `json:"arm"`
	ClassOrder            string              `json:"class_order"`
	ReplayPolicy          string              `json:"replay_policy"`
	Stage                 int                 `json:"stage"`
	BaseOriginalAccuracy  float64             `json:"base_original_accuracy"`
	BaseUnseenAccuracy    float64             `json:"base_unseen_accuracy"`
	AcquiredAggregateAcc  float64             `json:"acquired_aggregate_accuracy"`
	StoreMetrics          []UP124BStoreMetric `json:"store_metrics"`
}

type UP124BClasswideProjectorResult struct {
	Schema              string        `json:"schema"`
	Experiment          string        `json:"experiment"`
	SourceUP123BSeal    string        `json:"source_up123b_seal"`
	StateDimension      int           `json:"state_dimension"`
	LearningRate        float64       `json:"learning_rate"`
	EpochsPerBatch      int           `json:"epochs_per_batch"`
	GroundingPerVerb    int           `json:"grounding_per_verb"`
	FixedBudget         int           `json:"fixed_budget"`
	ClasswideProjector  bool          `json:"classwide_projector"`
	PerVerbRecompute    bool          `json:"per_verb_recompute"`
	AdaptiveGeometry    bool          `json:"adaptive_geometry"`
	Points              []UP124BPoint `json:"points"`
}

func up124bCompetitorDirections()(o,r [64]float64){
	o=up123bMeanOf("observes","sees","notes","watches","notices","spots")
	r=up123bMeanOf("reports","recalls","tells","remembers","recounts","retells")
	return
}

func up124bProject(h,o,r [64]float64)[64]float64{
	oo,rr,or:=up120bDot(o,o),up120bDot(r,r),up120bDot(o,r)
	ho,hr:=up120bDot(h,o),up120bDot(h,r)
	det:=oo*rr-or*or
	a,b:=0.0,0.0
	if math.Abs(det)>1e-12{
		a=(ho*rr-hr*or)/det
		b=(hr*oo-ho*or)/det
	}
	var residual [64]float64
	for i:=0;i<64;i++{residual[i]=h[i]-a*o[i]-b*r[i]}
	nh,nr:=up120bNorm(h),up120bNorm(residual)
	if nr>1e-12{
		scale:=nh/nr
		for i:=0;i<64;i++{residual[i]*=scale}
	}
	return residual
}

func up124bIsStoreVerb(verb string)bool{
	switch verb{
	case "stores","keeps","keepsafe","archives": return true
	default:return false
	}
}

func up124bEncode(arm,name,verb string,storesOrth,archivesDelta,o,r [64]float64)[64]float64{
	if arm=="local_corrections"{
		return up123bEncode("stores_plus_archives",name,verb,storesOrth,archivesDelta)
	}
	h:=up95bEncode(name+" "+verb)
	if up124bIsStoreVerb(verb){h=up124bProject(h,o,r)}
	return h
}

func up124bTrainStep(c *up97bClassifier,arm,name string,spec up106bVerbSpec,storesOrth,archivesDelta,o,r [64]float64){
	h:=up124bEncode(arm,name,spec.verb,storesOrth,archivesDelta,o,r)
	p:=c.probs(h)
	for k:=0;k<3;k++{
		g:=p[k];if k==spec.class{g-=1}
		for i:=0;i<64;i++{c.w[k][i]-=0.08*g*h[i]}
		c.b[k]-=0.08*g
	}
}

func up124bTrainBase(arm string,storesOrth,archivesDelta,o,r [64]float64)*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++{
		for _,name:=range up121bOriginalNames{
			for _,spec:=range up106bBase{up124bTrainStep(c,arm,name,spec,storesOrth,archivesDelta,o,r)}
		}
	}
	return c
}

func up124bAcquire(c *up97bClassifier,arm string,current up106bVerbSpec,previous []up106bVerbSpec,policy string,stage int,storesOrth,archivesDelta,o,r [64]float64){
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<4;ni++{up124bTrainStep(c,arm,up121bOriginalNames[ni],current,storesOrth,archivesDelta,o,r)}
		for _,spec:=range up106bBase{up124bTrainStep(c,arm,up121bOriginalNames[0],spec,storesOrth,archivesDelta,o,r)}
		var seq []up109bReplayExample
		if policy=="stage3_anchor2"&&stage==3{seq=up112bExtraReplay(previous,current.class,epoch,2)}else{seq=up111bCurrentExcluded(previous,current.class,epoch)}
		for _,x:=range seq{up124bTrainStep(c,arm,up121bOriginalNames[x.ni],x.spec,storesOrth,archivesDelta,o,r)}
	}
}

func up124bVerbAcc(c *up97bClassifier,arm string,spec up106bVerbSpec,names []string,storesOrth,archivesDelta,o,r [64]float64)float64{
	hits:=0
	for _,name:=range names{
		if up97bArgmax(c.probs(up124bEncode(arm,name,spec.verb,storesOrth,archivesDelta,o,r)))==spec.class{hits++}
	}
	return float64(hits)/float64(len(names))
}

func up124bMargin(c *up97bClassifier,arm,verb string,names []string,storesOrth,archivesDelta,o,r [64]float64)float64{
	min:=math.Inf(1)
	for _,name:=range names{
		z:=up117bLogits(c,up124bEncode(arm,name,verb,storesOrth,archivesDelta,o,r))
		m:=z[up97bStore]-math.Max(z[up97bObserve],z[up97bReport])
		if m<min{min=m}
	}
	return min
}

func up124bBaseAcc(c *up97bClassifier,arm string,names []string,storesOrth,archivesDelta,o,r [64]float64)float64{
	s:=0.0
	for _,spec:=range up106bBase{s+=up124bVerbAcc(c,arm,spec,names,storesOrth,archivesDelta,o,r)}
	return s/float64(len(up106bBase))
}

func up124bAcquiredAcc(c *up97bClassifier,arm string,acquired []up106bVerbSpec,storesOrth,archivesDelta,o,r [64]float64)float64{
	if len(acquired)==0{return 0}
	s:=0.0
	for _,spec:=range acquired{s+=up124bVerbAcc(c,arm,spec,up121bUnseenNames[4:6],storesOrth,archivesDelta,o,r)}
	return s/float64(len(acquired))
}

func up124bPoint(arm,order,policy string,stage int,c *up97bClassifier,acquired []up106bVerbSpec,storesOrth,archivesDelta,o,r [64]float64)UP124BPoint{
	p:=UP124BPoint{
		Arm:arm,ClassOrder:order,ReplayPolicy:policy,Stage:stage,
		BaseOriginalAccuracy:up124bBaseAcc(c,arm,up121bOriginalNames,storesOrth,archivesDelta,o,r),
		BaseUnseenAccuracy:up124bBaseAcc(c,arm,up121bUnseenNames,storesOrth,archivesDelta,o,r),
		AcquiredAggregateAcc:up124bAcquiredAcc(c,arm,acquired,storesOrth,archivesDelta,o,r),
	}
	for _,verb:=range []string{"stores","keeps","keepsafe","archives"}{
		spec:=up106bVerbSpec{verb:verb,class:up97bStore}
		p.StoreMetrics=append(p.StoreMetrics,UP124BStoreMetric{
			Verb:verb,
			OriginalAcc:up124bVerbAcc(c,arm,spec,up121bOriginalNames,storesOrth,archivesDelta,o,r),
			UnseenAcc:up124bVerbAcc(c,arm,spec,up121bUnseenNames,storesOrth,archivesDelta,o,r),
			OriginalMargin:up124bMargin(c,arm,verb,up121bOriginalNames,storesOrth,archivesDelta,o,r),
			UnseenMargin:up124bMargin(c,arm,verb,up121bUnseenNames,storesOrth,archivesDelta,o,r),
		})
	}
	return p
}

func RunUP124B()(UP124BClasswideProjectorResult,error){
	storesOrth:=up120bOrthogonalDelta()
	archivesDelta:=up123bArchivesDelta()
	o,r:=up124bCompetitorDirections()
	result:=UP124BClasswideProjectorResult{
		Schema:UP124BClasswideProjectorSchema,Experiment:"UP-124B-classwide-store-projector",
		SourceUP123BSeal:"b3f33eae8b7898c296cde750ac42553eac2bd834",
		StateDimension:64,LearningRate:0.08,EpochsPerBatch:20,GroundingPerVerb:4,FixedBudget:6,
		ClasswideProjector:true,PerVerbRecompute:false,AdaptiveGeometry:false,
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,arm:=range []string{"local_corrections","classwide_store_projector"}{
		for _,order:=range orders{
			name:=up114bOrderName(order)
			seq:=up121bSequence(order,"alternate")
			for _,policy:=range []string{"current_class_excluded","stage3_anchor2"}{
				c:=up124bTrainBase(arm,storesOrth,archivesDelta,o,r)
				acquired:=[]up106bVerbSpec{}
				result.Points=append(result.Points,up124bPoint(arm,name,policy,0,c,acquired,storesOrth,archivesDelta,o,r))
				for idx,spec:=range seq{
					stage:=idx+1
					up124bAcquire(c,arm,spec,acquired,policy,stage,storesOrth,archivesDelta,o,r)
					acquired=append(acquired,spec)
					result.Points=append(result.Points,up124bPoint(arm,name,policy,stage,c,acquired,storesOrth,archivesDelta,o,r))
				}
			}
		}
	}
	return result,nil
}
