package unitary

import (
	"math"
	"sort"
	"strings"
)

type wlmLmR28Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	Predictions [4]float64 `json:"predictions"`
	ReferencePrediction float64 `json:"reference_prediction"`
	ActualAdvantage float64 `json:"actual_advantage"`
	ConsensusOnly bool `json:"consensus_only"`
	ConsensusRelativeSpread bool `json:"consensus_relative_spread"`
	RobustLowerBound bool `json:"robust_lower_bound"`
}

type wlmLmR28Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR28Case `json:"cases"`
}

type wlmLmR28ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

type wlmLmR28Agg struct {
	covered int
	signCorrect int
	absErr float64
	positiveReco int
	falsePositive int
}

func wlmLmR28Sign(v float64) int {
	if v>0 { return 1 }
	if v<0 { return -1 }
	return 0
}

func wlmLmR28Median4(v [4]float64) float64 {
	x:=[]float64{v[0],v[1],v[2],v[3]}
	sort.Float64s(x)
	return (x[1]+x[2])/2
}

func wlmLmR28Flags(p [4]float64) (bool,bool,bool) {
	s0:=wlmLmR28Sign(p[0])
	consensus:=s0!=0
	for i:=1;i<4;i++ {
		if wlmLmR28Sign(p[i])!=s0 { consensus=false }
	}
	lo,hi:=p[0],p[0]
	allPos:=p[0]>0
	for i:=1;i<4;i++ {
		if p[i]<lo { lo=p[i] }
		if p[i]>hi { hi=p[i] }
		if p[i]<=0 { allPos=false }
	}
	rng:=hi-lo
	med:=wlmLmR28Median4(p)
	rel:=consensus && rng<=math.Abs(med)
	robust:=allPos && lo>rng
	return consensus,rel,robust
}

func wlmLmR28Update(a *wlmLmR28Agg, pred, actual float64) {
	a.covered++
	if wlmLmR28Sign(pred)==wlmLmR28Sign(actual) { a.signCorrect++ }
	a.absErr+=math.Abs(pred-actual)
	if pred>0 {
		a.positiveReco++
		if actual<=0 { a.falsePositive++ }
	}
}

func wlmLmR28Finalize(prefix string,a wlmLmR28Agg,total int,m map[string]float64) {
	m[prefix+"_covered_case_count"]=float64(a.covered)
	if total>0 { m[prefix+"_coverage"]=float64(a.covered)/float64(total) }
	if a.covered>0 {
		m[prefix+"_sign_accuracy"]=float64(a.signCorrect)/float64(a.covered)
		m[prefix+"_mean_absolute_error"]=a.absErr/float64(a.covered)
	}
	m[prefix+"_positive_recommendation_count"]=float64(a.positiveReco)
	if a.positiveReco>0 {
		m[prefix+"_false_positive_reallocation_rate"]=float64(a.falsePositive)/float64(a.positiveReco)
	}
	m[prefix+"_false_positive_count"]=float64(a.falsePositive)
}

func RunWlmLmExternalFutureDataCalibrationConfidenceAttributionR28(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse []byte,
) interface{} {
	budCode:=[]int{0,436,582,727,872}
	budOther:=[]int{0,291,436,581,726}
	mk:=func(d string,b []byte,h string,n int,bud []int) wlmLmR27Source {
		return wlmLmR27Source{d:d,b:b,h:h,n:n,budgets:bud}
	}
	inputs:=[]wlmLmR28ManifestInput{
		{"transfer",[]wlmLmR27Source{
			mk("code",transferCode,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45",41453,budCode),
			mk("structured",transferStructured,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932",14365,budOther),
			mk("technical_prose",transferProse,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b",1454,budOther),
		}},
		{"third",[]wlmLmR27Source{
			mk("code",thirdCode,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff",41453,budCode),
			mk("structured",thirdStructured,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56",14365,budOther),
			mk("technical_prose",thirdProse,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017",1454,budOther),
		}},
		{"fourth",[]wlmLmR27Source{
			mk("code",fourthCode,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1",41453,budCode),
			mk("structured",fourthStructured,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933",14365,budOther),
			mk("technical_prose",fourthProse,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222",1454,budOther),
		}},
		{"fifth",[]wlmLmR27Source{
			mk("code",fifthCode,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a",41453,budCode),
			mk("structured",fifthStructured,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281",14365,budOther),
			mk("technical_prose",fifthProse,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49",1454,budOther),
		}},
	}
	m:=map[string]float64{
		"outer_holdout_manifest_count":4,
		"policies_per_holdout":10,
		"prospective_case_count":0,
		"models_per_case":4,
		"heldout_outcome_use_before_flags":0,
		"post_result_predicate_choice_count":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"counter_overflow_count":0,
		"invalid_row_count":0,
		"source_identity_mismatch_count":0,
	}
	cases:=make([]wlmLmR28Case,0,40)
	var ungated,consensus,relative,robust wlmLmR28Agg
	for hold:=0;hold<4;hold++ {
		trainEx:=[]wlmLmR27Example{}
		trainIDs:=[]int{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			_,ex:=wlmLmR27Build("r28_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainEx=append(trainEx,ex...)
			trainIDs=append(trainIDs,j)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		mainBeta,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["invalid_row_count"]++ }
		var sub [3][7]float64
		for si,id:=range trainIDs {
			b,ok:=wlmLmR27Fit(trainEx,id)
			if !ok { m["invalid_row_count"]++ }
			sub[si]=b
		}
		heldFeatures,_:=wlmLmR27Build("r28_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct {
			alloc [3]int
			pred [4]float64
			cons bool
			rel bool
			rob bool
		}
		frozen:=make([]frozenCase,0,10)
		for _,alloc:=range wlmLmR27Candidates {
			var p [4]float64
			p[0]=wlmLmR27PolicyPred(heldFeatures,mainBeta,alloc)-wlmLmR27PolicyPred(heldFeatures,mainBeta,wlmLmR27Equal)
			for i:=0;i<3;i++ {
				p[i+1]=wlmLmR27PolicyPred(heldFeatures,sub[i],alloc)-wlmLmR27PolicyPred(heldFeatures,sub[i],wlmLmR27Equal)
			}
			co,re,ro:=wlmLmR28Flags(p)
			frozen=append(frozen,frozenCase{alloc:alloc,pred:p,cons:co,rel:re,rob:ro})
		}
		// Held-out outcomes are deliberately constructed only after all 10 policy reliability flags are frozen.
		heldActual,_:=wlmLmR27Build("r28_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			ref:=fc.pred[0]
			wlmLmR28Update(&ungated,ref,actual)
			if fc.cons { wlmLmR28Update(&consensus,ref,actual) }
			if fc.rel { wlmLmR28Update(&relative,ref,actual) }
			if fc.rob { wlmLmR28Update(&robust,ref,actual) }
			cases=append(cases,wlmLmR28Case{
				HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,Predictions:fc.pred,
				ReferencePrediction:ref,ActualAdvantage:actual,
				ConsensusOnly:fc.cons,ConsensusRelativeSpread:fc.rel,RobustLowerBound:fc.rob,
			})
		}
	}
	m["prospective_case_count"]=float64(len(cases))
	wlmLmR28Finalize("ungated",ungated,40,m)
	wlmLmR28Finalize("consensus_only",consensus,40,m)
	wlmLmR28Finalize("consensus_relative_spread",relative,40,m)
	wlmLmR28Finalize("robust_lower_bound",robust,40,m)
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if len(cases)!=40 { m["invalid_row_count"]++ }
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmR28Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-CALIBRATION-CONFIDENCE-ATTRIBUTION-R28",
		Metrics:m,Cases:cases,
	}
}
