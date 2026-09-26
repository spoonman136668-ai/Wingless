package unitary

import "math"

const UP120BCompetitorOrthogonalSchema = "wingless.up120b-competitor-orthogonal.v1"

type UP120BStaticMetric struct {
	Arm               string  `json:"arm"`
	StoreCentroidCos  float64 `json:"store_centroid_cosine"`
	ObserveCentroidCos float64 `json:"observe_centroid_cosine"`
	ReportCentroidCos float64 `json:"report_centroid_cosine"`
}

type UP120BPoint struct {
	RepresentationArm          string  `json:"representation_arm"`
	ClassOrder                 string  `json:"class_order"`
	ReplayPolicy               string  `json:"replay_policy"`
	Stage                      int     `json:"stage"`
	StoresAccuracy             float64 `json:"stores_accuracy"`
	KeepsAccuracy              float64 `json:"keeps_accuracy"`
	BaseLexiconAccuracy        float64 `json:"base_lexicon_accuracy"`
	AcquiredAggregateAccuracy  float64 `json:"acquired_aggregate_accuracy"`
	MeanStoreObserveMargin     float64 `json:"mean_store_observe_margin"`
	MinStoreObserveMargin      float64 `json:"min_store_observe_margin"`
	MeanStoreReportMargin      float64 `json:"mean_store_report_margin"`
	MinStoreReportMargin       float64 `json:"min_store_report_margin"`
	WrongObserveCount          int     `json:"wrong_observe_count"`
	WrongReportCount           int     `json:"wrong_report_count"`
}

type UP120BCompetitorOrthogonalResult struct {
	Schema           string               `json:"schema"`
	Experiment       string               `json:"experiment"`
	SourceUP119BSeal string               `json:"source_up119b_seal"`
	StateDimension   int                  `json:"state_dimension"`
	LearningRate     float64              `json:"learning_rate"`
	EpochsPerBatch   int                  `json:"epochs_per_batch"`
	GroundingPerVerb int                  `json:"grounding_per_verb"`
	FixedBudget      int                  `json:"fixed_budget"`
	AdaptiveGeometry bool                 `json:"adaptive_geometry"`
	StaticMetrics    []UP120BStaticMetric `json:"static_metrics"`
	Points           []UP120BPoint        `json:"points"`
}

func up120bDot(a,b [64]float64)float64{
	s:=0.0
	for i:=0;i<64;i++{s+=a[i]*b[i]}
	return s
}

func up120bNorm(a [64]float64)float64{return math.Sqrt(up120bDot(a,a))}

func up120bOrthogonalDelta()[64]float64{
	s:=up117bCentroid(up97bStore,"stores")
	o:=up117bCentroid(up97bObserve,"stores")
	r:=up117bCentroid(up97bReport,"stores")
	oo,rr,or:=up120bDot(o,o),up120bDot(r,r),up120bDot(o,r)
	so,sr:=up120bDot(s,o),up120bDot(s,r)
	det:=oo*rr-or*or
	a,b:=0.0,0.0
	if math.Abs(det)>1e-12{
		a=(so*rr-sr*or)/det
		b=(sr*oo-so*or)/det
	}
	var residual [64]float64
	for i:=0;i<64;i++{residual[i]=s[i]-a*o[i]-b*r[i]}
	ns,nr:=up120bNorm(s),up120bNorm(residual)
	if nr>1e-12{
		scale:=ns/nr
		for i:=0;i<64;i++{residual[i]*=scale}
	}
	stores:=up117bMeanVector("stores")
	var d [64]float64
	for i:=0;i<64;i++{d[i]=residual[i]-stores[i]}
	return d
}

func up120bDeltaForArm(arm string,centroid,orth [64]float64)[64]float64{
	switch arm{
	case "stores_centroid_shift": return centroid
	case "stores_competitor_orthogonal": return orth
	default: return [64]float64{}
	}
}

func up120bEncode(arm string,ni int,verb string,centroid,orth [64]float64)[64]float64{
	h:=up106bEncode(ni,verb)
	if verb=="stores"{
		d:=up120bDeltaForArm(arm,centroid,orth)
		for i:=0;i<64;i++{h[i]+=d[i]}
	}
	return h
}

func up120bMeanVector(arm,verb string,centroid,orth [64]float64)[64]float64{
	var out [64]float64
	for ni:=0;ni<6;ni++{
		h:=up120bEncode(arm,ni,verb,centroid,orth)
		for i:=0;i<64;i++{out[i]+=h[i]/6.0}
	}
	return out
}

func up120bTrainStep(c *up97bClassifier,arm string,ni int,spec up106bVerbSpec,centroid,orth [64]float64){
	h:=up120bEncode(arm,ni,spec.verb,centroid,orth)
	p:=c.probs(h)
	for k:=0;k<3;k++{
		g:=p[k]
		if k==spec.class{g-=1}
		for i:=0;i<64;i++{c.w[k][i]-=0.08*g*h[i]}
		c.b[k]-=0.08*g
	}
}

func up120bTrainBase(arm string,centroid,orth [64]float64)*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<6;ni++{for _,spec:=range up106bBase{up120bTrainStep(c,arm,ni,spec,centroid,orth)}}
	}
	return c
}

func up120bReplay(c *up97bClassifier,arm string,seq []up109bReplayExample,centroid,orth [64]float64){
	for _,x:=range seq{up120bTrainStep(c,arm,x.ni,x.spec,centroid,orth)}
}

func up120bAcquire(c *up97bClassifier,arm string,current up106bVerbSpec,previous []up106bVerbSpec,policy string,stage int,centroid,orth [64]float64){
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<4;ni++{up120bTrainStep(c,arm,ni,current,centroid,orth)}
		for _,spec:=range up106bBase{up120bTrainStep(c,arm,0,spec,centroid,orth)}
		var seq []up109bReplayExample
		if policy=="stage3_anchor2"&&stage==3{seq=up112bExtraReplay(previous,current.class,epoch,2)}else{seq=up111bCurrentExcluded(previous,current.class,epoch)}
		up120bReplay(c,arm,seq,centroid,orth)
	}
}

func up120bVerbAccuracy(c *up97bClassifier,arm string,spec up106bVerbSpec,start,end int,centroid,orth [64]float64)float64{
	hits,total:=0,0
	for ni:=start;ni<end;ni++{
		p:=up97bArgmax(c.probs(up120bEncode(arm,ni,spec.verb,centroid,orth)))
		total++
		if p==spec.class{hits++}
	}
	return float64(hits)/float64(total)
}

func up120bBaseAccuracy(c *up97bClassifier,arm string,centroid,orth [64]float64)float64{
	s:=0.0
	for _,spec:=range up106bBase{s+=up120bVerbAccuracy(c,arm,spec,0,6,centroid,orth)}
	return s/float64(len(up106bBase))
}

func up120bAcquiredAccuracy(c *up97bClassifier,arm string,acquired []up106bVerbSpec,centroid,orth [64]float64)float64{
	if len(acquired)==0{return 0}
	s:=0.0
	for _,spec:=range acquired{s+=up120bVerbAccuracy(c,arm,spec,4,6,centroid,orth)}
	return s/float64(len(acquired))
}

func up120bPoint(arm,order,policy string,stage int,c *up97bClassifier,acquired []up106bVerbSpec,centroid,orth [64]float64)UP120BPoint{
	sumSO,sumSR:=0.0,0.0
	minSO,minSR:=math.Inf(1),math.Inf(1)
	hits,wrongO,wrongR:=0,0,0
	for ni:=0;ni<6;ni++{
		z:=up117bLogits(c,up120bEncode(arm,ni,"stores",centroid,orth))
		p:=0;if z[1]>z[p]{p=1};if z[2]>z[p]{p=2}
		if p==up97bStore{hits++};if p==up97bObserve{wrongO++};if p==up97bReport{wrongR++}
		so:=z[up97bStore]-z[up97bObserve];sr:=z[up97bStore]-z[up97bReport]
		sumSO+=so;sumSR+=sr;if so<minSO{minSO=so};if sr<minSR{minSR=sr}
	}
	keeps:=up106bVerbSpec{verb:"keeps",class:up97bStore}
	return UP120BPoint{
		RepresentationArm:arm,ClassOrder:order,ReplayPolicy:policy,Stage:stage,
		StoresAccuracy:float64(hits)/6.0,KeepsAccuracy:up120bVerbAccuracy(c,arm,keeps,0,6,centroid,orth),
		BaseLexiconAccuracy:up120bBaseAccuracy(c,arm,centroid,orth),
		AcquiredAggregateAccuracy:up120bAcquiredAccuracy(c,arm,acquired,centroid,orth),
		MeanStoreObserveMargin:sumSO/6,MinStoreObserveMargin:minSO,
		MeanStoreReportMargin:sumSR/6,MinStoreReportMargin:minSR,
		WrongObserveCount:wrongO,WrongReportCount:wrongR,
	}
}

func RunUP120B()(UP120BCompetitorOrthogonalResult,error){
	centroid:=up118bDelta()
	orth:=up120bOrthogonalDelta()
	result:=UP120BCompetitorOrthogonalResult{
		Schema:UP120BCompetitorOrthogonalSchema,Experiment:"UP-120B-competitor-orthogonal",
		SourceUP119BSeal:"fc178cbf08c5ab6cc223e11a723cf8c683ce3ec0",
		StateDimension:64,LearningRate:0.08,EpochsPerBatch:20,GroundingPerVerb:4,FixedBudget:6,AdaptiveGeometry:false,
	}
	for _,arm:=range []string{"baseline","stores_centroid_shift","stores_competitor_orthogonal"}{
		v:=up120bMeanVector(arm,"stores",centroid,orth)
		result.StaticMetrics=append(result.StaticMetrics,UP120BStaticMetric{
			Arm:arm,
			StoreCentroidCos:up117bCos(v,up117bCentroid(up97bStore,"stores")),
			ObserveCentroidCos:up117bCos(v,up117bCentroid(up97bObserve,"stores")),
			ReportCentroidCos:up117bCos(v,up117bCentroid(up97bReport,"stores")),
		})
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,arm:=range []string{"baseline","stores_centroid_shift","stores_competitor_orthogonal"}{
		for _,order:=range orders{
			name:=up114bOrderName(order);seq:=up114bSequence(order)
			for _,policy:=range []string{"current_class_excluded","stage3_anchor2"}{
				c:=up120bTrainBase(arm,centroid,orth);acquired:=[]up106bVerbSpec{}
				result.Points=append(result.Points,up120bPoint(arm,name,policy,0,c,acquired,centroid,orth))
				for idx,spec:=range seq{
					stage:=idx+1;up120bAcquire(c,arm,spec,acquired,policy,stage,centroid,orth);acquired=append(acquired,spec)
					result.Points=append(result.Points,up120bPoint(arm,name,policy,stage,c,acquired,centroid,orth))
				}
			}
		}
	}
	return result,nil
}
