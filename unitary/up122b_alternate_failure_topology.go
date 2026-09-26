package unitary

import "math"

const UP122BAlternateFailureTopologySchema = "wingless.up122b-alternate-failure-topology.v1"

type UP122BControlPoint struct {
	ClassOrder             string  `json:"class_order"`
	ReplayPolicy           string  `json:"replay_policy"`
	Stage                  int     `json:"stage"`
	StoresOriginalAccuracy float64 `json:"stores_original_accuracy"`
	StoresUnseenAccuracy   float64 `json:"stores_unseen_accuracy"`
	BaseOriginalAccuracy   float64 `json:"base_original_accuracy"`
	BaseUnseenAccuracy     float64 `json:"base_unseen_accuracy"`
}

type UP122BVerbDiagnostic struct {
	ClassOrder       string  `json:"class_order"`
	ReplayPolicy     string  `json:"replay_policy"`
	Stage            int     `json:"stage"`
	Verb             string  `json:"verb"`
	Class            int     `json:"class"`
	SubjectFamily    string  `json:"subject_family"`
	Acquired         bool    `json:"acquired"`
	Accuracy         float64 `json:"accuracy"`
	MeanClassMargin  float64 `json:"mean_correct_class_margin"`
	MinClassMargin   float64 `json:"min_correct_class_margin"`
	WrongStoreCount  int     `json:"wrong_store_count"`
	WrongObserveCount int    `json:"wrong_observe_count"`
	WrongReportCount int     `json:"wrong_report_count"`
}

type UP122BAlternateFailureTopologyResult struct {
	Schema           string                 `json:"schema"`
	Experiment       string                 `json:"experiment"`
	SourceUP121BSeal string                 `json:"source_up121b_seal"`
	StateDimension   int                    `json:"state_dimension"`
	LearningRate     float64                `json:"learning_rate"`
	EpochsPerBatch   int                    `json:"epochs_per_batch"`
	GroundingPerVerb int                    `json:"grounding_per_verb"`
	FixedBudget      int                    `json:"fixed_budget"`
	TrainingChanged  bool                   `json:"training_changed"`
	AdaptiveGeometry bool                   `json:"adaptive_geometry"`
	Controls         []UP122BControlPoint   `json:"controls"`
	Diagnostics      []UP122BVerbDiagnostic `json:"diagnostics"`
}

func up122bDiagnostic(c *up97bClassifier,order,policy string,stage int,spec up106bVerbSpec,names []string,family string,acquired bool,orth [64]float64) UP122BVerbDiagnostic {
	hits,total:=0,0
	sumMargin:=0.0
	minMargin:=math.Inf(1)
	wrongStore,wrongObserve,wrongReport:=0,0,0
	for _,name:=range names {
		z:=up117bLogits(c,up121bEncodeName("stores_competitor_orthogonal",name,spec.verb,orth))
		pred:=0
		if z[1]>z[pred]{pred=1}
		if z[2]>z[pred]{pred=2}
		total++
		if pred==spec.class { hits++ } else {
			switch pred {
			case up97bStore: wrongStore++
			case up97bObserve: wrongObserve++
			case up97bReport: wrongReport++
			}
		}
		bestOther:=math.Inf(-1)
		for k:=0;k<3;k++ {
			if k==spec.class { continue }
			if z[k]>bestOther { bestOther=z[k] }
		}
		margin:=z[spec.class]-bestOther
		sumMargin+=margin
		if margin<minMargin { minMargin=margin }
	}
	return UP122BVerbDiagnostic{
		ClassOrder:order,ReplayPolicy:policy,Stage:stage,Verb:spec.verb,Class:spec.class,
		SubjectFamily:family,Acquired:acquired,
		Accuracy:float64(hits)/float64(total),
		MeanClassMargin:sumMargin/float64(total),MinClassMargin:minMargin,
		WrongStoreCount:wrongStore,WrongObserveCount:wrongObserve,WrongReportCount:wrongReport,
	}
}

func RunUP122B()(UP122BAlternateFailureTopologyResult,error){
	orth:=up120bOrthogonalDelta()
	result:=UP122BAlternateFailureTopologyResult{
		Schema:UP122BAlternateFailureTopologySchema,
		Experiment:"UP-122B-alternate-failure-topology",
		SourceUP121BSeal:"b887eb467212f16540ebc1af20feec74bfa263ec",
		StateDimension:64,LearningRate:0.08,EpochsPerBatch:20,GroundingPerVerb:4,FixedBudget:6,
		TrainingChanged:false,AdaptiveGeometry:false,
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	stores:=up106bVerbSpec{verb:"stores",class:up97bStore}
	for _,order:=range orders {
		orderName:=up114bOrderName(order)
		seq:=up121bSequence(order,"alternate")
		for _,policy:=range []string{"current_class_excluded","stage3_anchor2"} {
			c:=up121bTrainBase("stores_competitor_orthogonal",orth)
			acquired:=[]up106bVerbSpec{}
			for stage:=0;stage<=6;stage++ {
				if stage>0 {
					spec:=seq[stage-1]
					up121bAcquire(c,"stores_competitor_orthogonal",spec,acquired,policy,stage,orth)
					acquired=append(acquired,spec)
				}
				result.Controls=append(result.Controls,UP122BControlPoint{
					ClassOrder:orderName,ReplayPolicy:policy,Stage:stage,
					StoresOriginalAccuracy:up121bVerbAccuracyNames(c,"stores_competitor_orthogonal",stores,up121bOriginalNames,orth),
					StoresUnseenAccuracy:up121bVerbAccuracyNames(c,"stores_competitor_orthogonal",stores,up121bUnseenNames,orth),
					BaseOriginalAccuracy:up121bBaseAccuracyNames(c,"stores_competitor_orthogonal",up121bOriginalNames,orth),
					BaseUnseenAccuracy:up121bBaseAccuracyNames(c,"stores_competitor_orthogonal",up121bUnseenNames,orth),
				})
				for i,spec:=range seq {
					isAcquired:=i<stage
					result.Diagnostics=append(result.Diagnostics,
						up122bDiagnostic(c,orderName,policy,stage,spec,up121bOriginalNames[4:6],"original_heldout",isAcquired,orth),
						up122bDiagnostic(c,orderName,policy,stage,spec,up121bUnseenNames[4:6],"unseen_heldout",isAcquired,orth),
					)
				}
			}
		}
	}
	return result,nil
}
