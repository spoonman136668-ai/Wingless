package unitary

import (
	"math"
	"strings"
)

type wlmLmR40Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	PolicyDeltaBase6 [6]float64 `json:"policy_delta_base6"`
	Region string `json:"region"`
	Base6PredictedAdvantage float64 `json:"base6_predicted_advantage"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR40Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR40Case `json:"cases"`
}

func wlmLmR40PolicyDelta(man wlmLmR27Manifest,alloc [3]int) [6]float64 {
	var out [6]float64
	for ai,arm:=range man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			coef:=0.0
			if hi<=alloc[ai] { coef+=1 }
			if hi<=wlmLmR27Equal[ai] { coef-=1 }
			if coef==0 { continue }
			fv:=arm.features[hi]
			for k:=0;k<6;k++ { out[k]+=coef*fv[k] }
		}
	}
	return out
}

func wlmLmR40Envelope(mans []wlmLmR27Manifest,m map[string]float64) ([6]float64,[6]float64) {
	var lo,hi [6]float64
	first:=true
	count:=0
	for _,man:=range mans {
		for _,alloc:=range wlmLmR27Candidates {
			if alloc==wlmLmR27Equal { continue }
			v:=wlmLmR40PolicyDelta(man,alloc)
			if first {
				lo=v;hi=v;first=false
			} else {
				for k:=0;k<6;k++ {
					if v[k]<lo[k] { lo[k]=v[k] }
					if v[k]>hi[k] { hi[k]=v[k] }
				}
			}
			count++
		}
	}
	if first || count!=27 { m["region_definition_mismatch_count"]++ }
	for k:=0;k<6;k++ {
		if math.IsNaN(lo[k])||math.IsInf(lo[k],0)||math.IsNaN(hi[k])||math.IsInf(hi[k],0)||lo[k]>hi[k] {
			m["region_definition_mismatch_count"]++
		}
	}
	return lo,hi
}

func wlmLmR40InEnvelope(v,lo,hi [6]float64) bool {
	for k:=0;k<6;k++ {
		if v[k]<lo[k] || v[k]>hi[k] { return false }
	}
	return true
}

func wlmLmR40Finalize(prefix string,a wlmLmR34Agg,m map[string]float64) {
	if a.n==0 {
		m[prefix+"_case_count"]=0
		m[prefix+"_sign_error_rate"]=0
		m[prefix+"_mean_absolute_error"]=0
		return
	}
	m[prefix+"_case_count"]=float64(a.n)
	m[prefix+"_sign_error_rate"]=float64(a.signErr)/float64(a.n)
	m[prefix+"_mean_absolute_error"]=a.absErr/float64(a.n)
}

func RunWlmLmExternalFutureDataProspectiveApplicabilityRegionAttributionR40(
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
	inputs:=[]wlmLmR35ManifestInput{
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
		"non_equal_case_count":0,"stable_nonzero_case_count":0,"excluded_case_count":0,
		"global_sign_error_rate":0,"global_mean_absolute_error":0,
		"in_envelope_case_count":0,"out_of_envelope_case_count":0,"in_envelope_coverage":0,
		"in_envelope_sign_error_rate":0,"in_envelope_mean_absolute_error":0,
		"out_of_envelope_sign_error_rate":0,"out_of_envelope_mean_absolute_error":0,
		"in_envelope_sign_error_improvement":0,"out_minus_in_sign_error_gap":0,
		"base6_identity_mismatch_count":0,"region_definition_mismatch_count":0,
		"fit_failure_count":0,"heldout_outcome_use_before_region_freeze":0,
		"post_result_region_choice_count":0,"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	var globalAgg,inAgg,outAgg wlmLmR34Agg
	cases:=make([]wlmLmR40Case,0,33)
	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainEx:=[]wlmLmR27Example{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r40_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man);trainEx=append(trainEx,ex...)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		baseR27,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		baseExamples:=wlmLmR35Examples(trainMans,false)
		baseBeta,ok0:=wlmLmR34Fit(baseExamples)
		if !ok0 { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		lo,hi:=wlmLmR40Envelope(trainMans,m)

		heldFeatures,_:=wlmLmR27Build("r40_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct{alloc [3]int;excluded,in bool;pred float64;vec [6]float64}
		frozen:=make([]frozenCase,0,9)
		for _,alloc:=range wlmLmR27Candidates {
			if alloc==wlmLmR27Equal { continue }
			m["non_equal_case_count"]++
			signs:=map[int]bool{}
			for _,tm:=range trainMans {
				adv:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				signs[wlmLmR29Sign(adv)]=true
			}
			stableNonZero:=len(signs)==1 && (signs[1]||signs[-1])
			if stableNonZero { m["stable_nonzero_case_count"]++ } else { m["excluded_case_count"]++ }

			pred:=wlmLmR35PolicyPred(heldFeatures,baseBeta,false,alloc)-wlmLmR35PolicyPred(heldFeatures,baseBeta,false,wlmLmR27Equal)
			r27:=wlmLmR27PolicyPred(heldFeatures,baseR27,alloc)-wlmLmR27PolicyPred(heldFeatures,baseR27,wlmLmR27Equal)
			if math.Abs(pred-r27)>1e-9 { m["base6_identity_mismatch_count"]++ }
			vec:=wlmLmR40PolicyDelta(heldFeatures,alloc)
			frozen=append(frozen,frozenCase{alloc:alloc,excluded:!stableNonZero,in:wlmLmR40InEnvelope(vec,lo,hi),pred:pred,vec:vec})
		}

		heldActual,_:=wlmLmR27Build("r40_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			if !fc.excluded { continue }
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR34Update(&globalAgg,fc.pred,actual)
			region:="OUT_OF_ENVELOPE"
			if fc.in {
				region="IN_ENVELOPE";wlmLmR34Update(&inAgg,fc.pred,actual)
			} else {
				wlmLmR34Update(&outAgg,fc.pred,actual)
			}
			cases=append(cases,wlmLmR40Case{HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,
				PolicyDeltaBase6:fc.vec,Region:region,Base6PredictedAdvantage:fc.pred,ActualAdvantage:actual})
		}
	}
	wlmLmR34Finalize("global",globalAgg,m)
	wlmLmR40Finalize("in_envelope",inAgg,m)
	wlmLmR40Finalize("out_of_envelope",outAgg,m)
	if globalAgg.n>0 { m["in_envelope_coverage"]=float64(inAgg.n)/float64(globalAgg.n) }
	m["in_envelope_sign_error_improvement"]=m["global_sign_error_rate"]-m["in_envelope_sign_error_rate"]
	m["out_minus_in_sign_error_gap"]=m["out_of_envelope_sign_error_rate"]-m["in_envelope_sign_error_rate"]
	if math.Abs(m["global_sign_error_rate"]-0.5151515151515151)>1e-9 ||
		math.Abs(m["global_mean_absolute_error"]-10.731783167689903)>1e-9 {
		m["base6_identity_mismatch_count"]++
	}
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 || int(m["excluded_case_count"])!=33 ||
		len(cases)!=33 || globalAgg.n!=33 || inAgg.n+outAgg.n!=33 {
		m["invalid_row_count"]++
	}
	if m["base6_identity_mismatch_count"]!=0 || m["region_definition_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR40Result{Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-PROSPECTIVE-APPLICABILITY-REGION-ATTRIBUTION-R40",
		Metrics:m,Cases:cases}
}
