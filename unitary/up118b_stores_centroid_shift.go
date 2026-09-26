package unitary

import "math"

const UP118BStoresCentroidShiftSchema = "wingless.up118b-stores-centroid-shift.v1"

type UP118BStaticMetric struct {
	Arm              string  `json:"arm"`
	StoreCentroidCos float64 `json:"store_centroid_cosine"`
	ObserveCentroidCos float64 `json:"observe_centroid_cosine"`
	ReportCentroidCos float64 `json:"report_centroid_cosine"`
}

type UP118BPoint struct {
	RepresentationArm          string  `json:"representation_arm"`
	ClassOrder                 string  `json:"class_order"`
	ReplayPolicy               string  `json:"replay_policy"`
	Stage                      int     `json:"stage"`
	StoresAccuracy             float64 `json:"stores_accuracy"`
	KeepsAccuracy              float64 `json:"keeps_accuracy"`
	BaseLexiconAccuracy        float64 `json:"base_lexicon_accuracy"`
	AcquiredAggregateAccuracy  float64 `json:"acquired_aggregate_accuracy"`
	StoresMeanStoreObserveMargin float64 `json:"stores_mean_store_observe_margin"`
	StoresMinStoreObserveMargin  float64 `json:"stores_min_store_observe_margin"`
}

type UP118BStoresCentroidShiftResult struct {
	Schema            string              `json:"schema"`
	Experiment        string              `json:"experiment"`
	SourceUP117BSeal  string              `json:"source_up117b_seal"`
	StateDimension    int                 `json:"state_dimension"`
	LearningRate      float64             `json:"learning_rate"`
	EpochsPerBatch    int                 `json:"epochs_per_batch"`
	GroundingPerVerb  int                 `json:"grounding_per_verb"`
	FixedBudget       int                 `json:"fixed_budget"`
	AdaptiveGeometry  bool                `json:"adaptive_geometry"`
	StaticMetrics     []UP118BStaticMetric `json:"static_metrics"`
	Points            []UP118BPoint        `json:"points"`
}

func up118bDelta() [64]float64 {
	stores:=up117bMeanVector("stores")
	var sib [64]float64
	for _,verb:=range []string{"keeps","holds","saves"} {
		v:=up117bMeanVector(verb)
		for i:=0;i<64;i++ { sib[i]+=v[i]/3.0 }
	}
	var d [64]float64
	for i:=0;i<64;i++ { d[i]=sib[i]-stores[i] }
	return d
}

func up118bEncode(arm string,ni int,verb string,delta [64]float64)[64]float64{
	h:=up106bEncode(ni,verb)
	if arm=="stores_centroid_shift" && verb=="stores" {
		for i:=0;i<64;i++ { h[i]+=delta[i] }
	}
	return h
}

func up118bMeanVector(arm,verb string,delta [64]float64)[64]float64{
	var out [64]float64
	for ni:=0;ni<6;ni++ {
		h:=up118bEncode(arm,ni,verb,delta)
		for i:=0;i<64;i++ { out[i]+=h[i]/6.0 }
	}
	return out
}

func up118bTrainStep(c *up97bClassifier,arm string,ni int,spec up106bVerbSpec,delta [64]float64){
	h:=up118bEncode(arm,ni,spec.verb,delta)
	p:=c.probs(h)
	for k:=0;k<3;k++ {
		g:=p[k]
		if k==spec.class { g-=1 }
		for i:=0;i<64;i++ { c.w[k][i]-=0.08*g*h[i] }
		c.b[k]-=0.08*g
	}
}

func up118bTrainBase(arm string,delta [64]float64)*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<6;ni++ {
			for _,spec:=range up106bBase { up118bTrainStep(c,arm,ni,spec,delta) }
		}
	}
	return c
}

func up118bGround(c *up97bClassifier,arm string,current up106bVerbSpec,delta [64]float64){
	for ni:=0;ni<4;ni++ { up118bTrainStep(c,arm,ni,current,delta) }
}

func up118bBaseReplay(c *up97bClassifier,arm string,delta [64]float64){
	for _,spec:=range up106bBase { up118bTrainStep(c,arm,0,spec,delta) }
}

func up118bReplay(c *up97bClassifier,arm string,seq []up109bReplayExample,delta [64]float64){
	for _,x:=range seq { up118bTrainStep(c,arm,x.ni,x.spec,delta) }
}

func up118bAcquire(c *up97bClassifier,arm string,current up106bVerbSpec,previous []up106bVerbSpec,policy string,stage int,delta [64]float64){
	for epoch:=0;epoch<20;epoch++ {
		up118bGround(c,arm,current,delta)
		up118bBaseReplay(c,arm,delta)
		var seq []up109bReplayExample
		if policy=="stage3_anchor2" && stage==3 {
			seq=up112bExtraReplay(previous,current.class,epoch,2)
		}else{
			seq=up111bCurrentExcluded(previous,current.class,epoch)
		}
		up118bReplay(c,arm,seq,delta)
	}
}

func up118bVerbAccuracy(c *up97bClassifier,arm string,spec up106bVerbSpec,start,end int,delta [64]float64)float64{
	hits,total:=0,0
	for ni:=start;ni<end;ni++ {
		pred:=up97bArgmax(c.probs(up118bEncode(arm,ni,spec.verb,delta)))
		total++
		if pred==spec.class { hits++ }
	}
	return float64(hits)/float64(total)
}

func up118bBaseAccuracy(c *up97bClassifier,arm string,delta [64]float64)float64{
	sum:=0.0
	for _,spec:=range up106bBase { sum+=up118bVerbAccuracy(c,arm,spec,0,6,delta) }
	return sum/float64(len(up106bBase))
}

func up118bAcquiredAccuracy(c *up97bClassifier,arm string,acquired []up106bVerbSpec,delta [64]float64)float64{
	if len(acquired)==0{return 0}
	sum:=0.0
	for _,spec:=range acquired { sum+=up118bVerbAccuracy(c,arm,spec,4,6,delta) }
	return sum/float64(len(acquired))
}

func up118bStoresMargin(c *up97bClassifier,arm string,delta [64]float64)(mean,min float64){
	sum:=0.0
	min=math.Inf(1)
	for ni:=0;ni<6;ni++ {
		z:=up117bLogits(c,up118bEncode(arm,ni,"stores",delta))
		m:=z[up97bStore]-z[up97bObserve]
		sum+=m
		if m<min{min=m}
	}
	return sum/6.0,min
}

func up118bPoint(arm,order,policy string,stage int,c *up97bClassifier,acquired []up106bVerbSpec,delta [64]float64)UP118BPoint{
	stores:=up106bVerbSpec{verb:"stores",class:up97bStore}
	keeps:=up106bVerbSpec{verb:"keeps",class:up97bStore}
	mean,min:=up118bStoresMargin(c,arm,delta)
	return UP118BPoint{
		RepresentationArm:arm,ClassOrder:order,ReplayPolicy:policy,Stage:stage,
		StoresAccuracy:up118bVerbAccuracy(c,arm,stores,0,6,delta),
		KeepsAccuracy:up118bVerbAccuracy(c,arm,keeps,0,6,delta),
		BaseLexiconAccuracy:up118bBaseAccuracy(c,arm,delta),
		AcquiredAggregateAccuracy:up118bAcquiredAccuracy(c,arm,acquired,delta),
		StoresMeanStoreObserveMargin:mean,StoresMinStoreObserveMargin:min,
	}
}

func RunUP118B()(UP118BStoresCentroidShiftResult,error){
	delta:=up118bDelta()
	result:=UP118BStoresCentroidShiftResult{
		Schema:UP118BStoresCentroidShiftSchema,Experiment:"UP-118B-stores-centroid-shift",
		SourceUP117BSeal:"fa5add5fe53f609913041f6a12b50c545dabc2cd",
		StateDimension:64,LearningRate:0.08,EpochsPerBatch:20,GroundingPerVerb:4,FixedBudget:6,AdaptiveGeometry:false,
	}
	for _,arm:=range []string{"baseline","stores_centroid_shift"} {
		sv:=up118bMeanVector(arm,"stores",delta)
		result.StaticMetrics=append(result.StaticMetrics,UP118BStaticMetric{
			Arm:arm,
			StoreCentroidCos:up117bCos(sv,up117bCentroid(up97bStore,"stores")),
			ObserveCentroidCos:up117bCos(sv,up117bCentroid(up97bObserve,"stores")),
			ReportCentroidCos:up117bCos(sv,up117bCentroid(up97bReport,"stores")),
		})
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,arm:=range []string{"baseline","stores_centroid_shift"} {
		for _,order:=range orders {
			orderName:=up114bOrderName(order)
			seq:=up114bSequence(order)
			for _,policy:=range []string{"current_class_excluded","stage3_anchor2"} {
				c:=up118bTrainBase(arm,delta)
				acquired:=[]up106bVerbSpec{}
				result.Points=append(result.Points,up118bPoint(arm,orderName,policy,0,c,acquired,delta))
				for idx,spec:=range seq {
					stage:=idx+1
					up118bAcquire(c,arm,spec,acquired,policy,stage,delta)
					acquired=append(acquired,spec)
					result.Points=append(result.Points,up118bPoint(arm,orderName,policy,stage,c,acquired,delta))
				}
			}
		}
	}
	return result,nil
}
