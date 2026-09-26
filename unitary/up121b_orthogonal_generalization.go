package unitary

import "math"

const UP121BOrthogonalGeneralizationSchema = "wingless.up121b-orthogonal-generalization.v1"

type UP121BPoint struct {
	RepresentationArm            string  `json:"representation_arm"`
	LexicalSet                   string  `json:"lexical_set"`
	ClassOrder                   string  `json:"class_order"`
	ReplayPolicy                 string  `json:"replay_policy"`
	Stage                        int     `json:"stage"`
	StoresOriginalAccuracy       float64 `json:"stores_original_accuracy"`
	StoresUnseenAccuracy         float64 `json:"stores_unseen_accuracy"`
	KeepsOriginalAccuracy        float64 `json:"keeps_original_accuracy"`
	KeepsUnseenAccuracy          float64 `json:"keeps_unseen_accuracy"`
	BaseOriginalAccuracy         float64 `json:"base_original_accuracy"`
	BaseUnseenAccuracy           float64 `json:"base_unseen_accuracy"`
	AcquiredAggregateAccuracy    float64 `json:"acquired_aggregate_accuracy"`
	MeanStoreObserveMarginOrig   float64 `json:"mean_store_observe_margin_original"`
	MinStoreObserveMarginOrig    float64 `json:"min_store_observe_margin_original"`
	MeanStoreReportMarginOrig    float64 `json:"mean_store_report_margin_original"`
	MinStoreReportMarginOrig     float64 `json:"min_store_report_margin_original"`
	MeanStoreObserveMarginUnseen float64 `json:"mean_store_observe_margin_unseen"`
	MinStoreObserveMarginUnseen  float64 `json:"min_store_observe_margin_unseen"`
	MeanStoreReportMarginUnseen  float64 `json:"mean_store_report_margin_unseen"`
	MinStoreReportMarginUnseen   float64 `json:"min_store_report_margin_unseen"`
}

type UP121BOrthogonalGeneralizationResult struct {
	Schema                 string        `json:"schema"`
	Experiment             string        `json:"experiment"`
	SourceUP120BSeal       string        `json:"source_up120b_seal"`
	StateDimension         int           `json:"state_dimension"`
	LearningRate           float64       `json:"learning_rate"`
	EpochsPerBatch         int           `json:"epochs_per_batch"`
	GroundingPerVerb       int           `json:"grounding_per_verb"`
	FixedBudget            int           `json:"fixed_budget"`
	CorrectionRecomputed   bool          `json:"correction_recomputed"`
	AdaptiveGeometry       bool          `json:"adaptive_geometry"`
	OriginalSubjects       []string      `json:"original_subjects"`
	UnseenSubjects         []string      `json:"unseen_subjects"`
	Points                 []UP121BPoint `json:"points"`
}

var up121bOriginalNames=[]string{"ada","ben","cy","dee","eli","fay"}
var up121bUnseenNames=[]string{"gia","hal","ivy","jon","kia","leo"}

func up121bAltPair(class int) []up106bVerbSpec {
	switch class {
	case up97bStore:
		return []up106bVerbSpec{{verb:"keepsafe",class:up97bStore},{verb:"archives",class:up97bStore}}
	case up97bObserve:
		return []up106bVerbSpec{{verb:"notices",class:up97bObserve},{verb:"spots",class:up97bObserve}}
	default:
		return []up106bVerbSpec{{verb:"recounts",class:up97bReport},{verb:"retells",class:up97bReport}}
	}
}

func up121bSequence(order []int,lexical string) []up106bVerbSpec {
	if lexical=="original" { return up114bSequence(order) }
	out:=make([]up106bVerbSpec,0,6)
	for _,class:=range order { out=append(out,up121bAltPair(class)[0]) }
	for _,class:=range order { out=append(out,up121bAltPair(class)[1]) }
	return out
}

func up121bEncodeName(arm,name,verb string,orth [64]float64)[64]float64{
	h:=up95bEncode(name+" "+verb)
	if arm=="stores_competitor_orthogonal" && verb=="stores" {
		for i:=0;i<64;i++ { h[i]+=orth[i] }
	}
	return h
}

func up121bTrainStep(c *up97bClassifier,arm,name string,spec up106bVerbSpec,orth [64]float64){
	h:=up121bEncodeName(arm,name,spec.verb,orth)
	p:=c.probs(h)
	for k:=0;k<3;k++{
		g:=p[k]
		if k==spec.class { g-=1 }
		for i:=0;i<64;i++ { c.w[k][i]-=0.08*g*h[i] }
		c.b[k]-=0.08*g
	}
}

func up121bTrainBase(arm string,orth [64]float64)*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++{
		for _,name:=range up121bOriginalNames {
			for _,spec:=range up106bBase { up121bTrainStep(c,arm,name,spec,orth) }
		}
	}
	return c
}

func up121bReplay(c *up97bClassifier,arm string,seq []up109bReplayExample,orth [64]float64){
	for _,x:=range seq {
		up121bTrainStep(c,arm,up121bOriginalNames[x.ni],x.spec,orth)
	}
}

func up121bAcquire(c *up97bClassifier,arm string,current up106bVerbSpec,previous []up106bVerbSpec,policy string,stage int,orth [64]float64){
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<4;ni++ { up121bTrainStep(c,arm,up121bOriginalNames[ni],current,orth) }
		for _,spec:=range up106bBase { up121bTrainStep(c,arm,up121bOriginalNames[0],spec,orth) }
		var seq []up109bReplayExample
		if policy=="stage3_anchor2" && stage==3 {
			seq=up112bExtraReplay(previous,current.class,epoch,2)
		}else{
			seq=up111bCurrentExcluded(previous,current.class,epoch)
		}
		up121bReplay(c,arm,seq,orth)
	}
}

func up121bVerbAccuracyNames(c *up97bClassifier,arm string,spec up106bVerbSpec,names []string,orth [64]float64)float64{
	hits:=0
	for _,name:=range names {
		if up97bArgmax(c.probs(up121bEncodeName(arm,name,spec.verb,orth)))==spec.class { hits++ }
	}
	return float64(hits)/float64(len(names))
}

func up121bBaseAccuracyNames(c *up97bClassifier,arm string,names []string,orth [64]float64)float64{
	s:=0.0
	for _,spec:=range up106bBase { s+=up121bVerbAccuracyNames(c,arm,spec,names,orth) }
	return s/float64(len(up106bBase))
}

func up121bAcquiredAccuracy(c *up97bClassifier,arm,lexical string,acquired []up106bVerbSpec,orth [64]float64)float64{
	if len(acquired)==0 { return 0 }
	names:=up121bOriginalNames[4:6]
	if lexical=="alternate" { names=up121bUnseenNames[4:6] }
	s:=0.0
	for _,spec:=range acquired { s+=up121bVerbAccuracyNames(c,arm,spec,names,orth) }
	return s/float64(len(acquired))
}

func up121bMargins(c *up97bClassifier,arm string,names []string,orth [64]float64)(float64,float64,float64,float64){
	sumSO,sumSR:=0.0,0.0
	minSO,minSR:=math.Inf(1),math.Inf(1)
	for _,name:=range names {
		z:=up117bLogits(c,up121bEncodeName(arm,name,"stores",orth))
		so:=z[up97bStore]-z[up97bObserve]
		sr:=z[up97bStore]-z[up97bReport]
		sumSO+=so;sumSR+=sr
		if so<minSO{minSO=so}
		if sr<minSR{minSR=sr}
	}
	return sumSO/float64(len(names)),minSO,sumSR/float64(len(names)),minSR
}

func up121bPoint(arm,lexical,order,policy string,stage int,c *up97bClassifier,acquired []up106bVerbSpec,orth [64]float64)UP121BPoint{
	stores:=up106bVerbSpec{verb:"stores",class:up97bStore}
	keeps:=up106bVerbSpec{verb:"keeps",class:up97bStore}
	msoO,minsoO,msrO,minsrO:=up121bMargins(c,arm,up121bOriginalNames,orth)
	msoU,minsoU,msrU,minsrU:=up121bMargins(c,arm,up121bUnseenNames,orth)
	return UP121BPoint{
		RepresentationArm:arm,LexicalSet:lexical,ClassOrder:order,ReplayPolicy:policy,Stage:stage,
		StoresOriginalAccuracy:up121bVerbAccuracyNames(c,arm,stores,up121bOriginalNames,orth),
		StoresUnseenAccuracy:up121bVerbAccuracyNames(c,arm,stores,up121bUnseenNames,orth),
		KeepsOriginalAccuracy:up121bVerbAccuracyNames(c,arm,keeps,up121bOriginalNames,orth),
		KeepsUnseenAccuracy:up121bVerbAccuracyNames(c,arm,keeps,up121bUnseenNames,orth),
		BaseOriginalAccuracy:up121bBaseAccuracyNames(c,arm,up121bOriginalNames,orth),
		BaseUnseenAccuracy:up121bBaseAccuracyNames(c,arm,up121bUnseenNames,orth),
		AcquiredAggregateAccuracy:up121bAcquiredAccuracy(c,arm,lexical,acquired,orth),
		MeanStoreObserveMarginOrig:msoO,MinStoreObserveMarginOrig:minsoO,
		MeanStoreReportMarginOrig:msrO,MinStoreReportMarginOrig:minsrO,
		MeanStoreObserveMarginUnseen:msoU,MinStoreObserveMarginUnseen:minsoU,
		MeanStoreReportMarginUnseen:msrU,MinStoreReportMarginUnseen:minsrU,
	}
}

func RunUP121B()(UP121BOrthogonalGeneralizationResult,error){
	orth:=up120bOrthogonalDelta()
	result:=UP121BOrthogonalGeneralizationResult{
		Schema:UP121BOrthogonalGeneralizationSchema,
		Experiment:"UP-121B-orthogonal-generalization",
		SourceUP120BSeal:"17bb8323a0df4583e8b41807e19b3d473b3a1043",
		StateDimension:64,LearningRate:0.08,EpochsPerBatch:20,GroundingPerVerb:4,FixedBudget:6,
		CorrectionRecomputed:false,AdaptiveGeometry:false,
		OriginalSubjects:append([]string(nil),up121bOriginalNames...),
		UnseenSubjects:append([]string(nil),up121bUnseenNames...),
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,arm:=range []string{"baseline","stores_competitor_orthogonal"} {
		for _,lexical:=range []string{"original","alternate"} {
			for _,order:=range orders {
				orderName:=up114bOrderName(order)
				seq:=up121bSequence(order,lexical)
				for _,policy:=range []string{"current_class_excluded","stage3_anchor2"} {
					c:=up121bTrainBase(arm,orth)
					acquired:=[]up106bVerbSpec{}
					result.Points=append(result.Points,up121bPoint(arm,lexical,orderName,policy,0,c,acquired,orth))
					for idx,spec:=range seq {
						stage:=idx+1
						up121bAcquire(c,arm,spec,acquired,policy,stage,orth)
						acquired=append(acquired,spec)
						result.Points=append(result.Points,up121bPoint(arm,lexical,orderName,policy,stage,c,acquired,orth))
					}
				}
			}
		}
	}
	return result,nil
}
