package unitary

import (
	"math"
	"strings"
)

type wlmLmR30Holdout struct {
	HoldoutManifest string `json:"holdout_manifest"`
	CorrectedSelection [3]int `json:"corrected_selection"`
	UncorrectedSelection [3]int `json:"uncorrected_selection"`
	CorrectedActualAdvantage float64 `json:"corrected_actual_advantage"`
	UncorrectedActualAdvantage float64 `json:"uncorrected_actual_advantage"`
	StableCandidateCount int `json:"stable_candidate_count"`
	UnstableCandidateCount int `json:"unstable_candidate_count"`
}

type wlmLmR30Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Holdouts []wlmLmR30Holdout `json:"holdouts"`
}

type wlmLmR30ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

func wlmLmR30Selection(held wlmLmR27Manifest, beta [7]float64, train []wlmLmR27Manifest, corrected bool) ([3]int,int,int) {
	best:=wlmLmR27Equal
	bestPred:=0.0
	stable,unstable:=0,0
	for _,alloc:=range wlmLmR27Candidates {
		if alloc==wlmLmR27Equal { continue }
		signs:=map[int]bool{}
		for _,tm:=range train {
			adv:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
			signs[wlmLmR29Sign(adv)]=true
		}
		stableNonZero:=false
		if len(signs)==1 {
			stableNonZero=signs[1] || signs[-1]
		}
		if stableNonZero { stable++ } else { unstable++ }
		if corrected && !stableNonZero { continue }
		pred:=wlmLmR27PolicyPred(held,beta,alloc)-wlmLmR27PolicyPred(held,beta,wlmLmR27Equal)
		if pred<=0 { continue }
		if best==wlmLmR27Equal || pred>bestPred || (pred==bestPred && wlmLmR27LexLess(alloc,best)) {
			best=alloc
			bestPred=pred
		}
	}
	return best,stable,unstable
}

func RunWlmLmExternalFutureDataHistoricalSignInstabilityCorrectionR30(
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
	inputs:=[]wlmLmR30ManifestInput{
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
		"outer_holdout_manifest_count":4,"case_count":0,
		"heldout_outcome_use_before_selection":0,"post_result_policy_choice_count":0,
		"stable_candidate_count":0,"unstable_candidate_count":0,
		"corrected_reallocation_count":0,"uncorrected_reallocation_count":0,"corrected_fallback_equal_count":0,
		"corrected_positive_holdout_count":0,"corrected_negative_holdout_count":0,"corrected_zero_holdout_count":0,
		"uncorrected_positive_holdout_count":0,"uncorrected_negative_holdout_count":0,"uncorrected_zero_holdout_count":0,
		"corrected_total_actual_advantage":0,"uncorrected_total_actual_advantage":0,"corrected_delta_vs_uncorrected":0,
		"source_identity_mismatch_count":0,"capacity_growth_event_count":0,
		"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0,
	}
	holdouts:=make([]wlmLmR30Holdout,0,4)
	for hold:=0;hold<4;hold++ {
		trainEx:=[]wlmLmR27Example{}
		trainMans:=make([]wlmLmR27Manifest,0,3)
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r30_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man)
			trainEx=append(trainEx,ex...)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		beta,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["invalid_row_count"]++ }
		heldFeatures,_:=wlmLmR27Build("r30_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		corrected,stable,unstable:=wlmLmR30Selection(heldFeatures,beta,trainMans,true)
		uncorrected,_,_:=wlmLmR30Selection(heldFeatures,beta,trainMans,false)
		m["stable_candidate_count"]+=float64(stable)
		m["unstable_candidate_count"]+=float64(unstable)
		if corrected==wlmLmR27Equal { m["corrected_fallback_equal_count"]++ } else { m["corrected_reallocation_count"]++ }
		if uncorrected!=wlmLmR27Equal { m["uncorrected_reallocation_count"]++ }
		if corrected[0]+corrected[1]+corrected[2]!=1744 || uncorrected[0]+uncorrected[1]+uncorrected[2]!=1744 {
			m["invalid_row_count"]++
		}
		// Both selectors are frozen above. Held-out outcomes are opened only below.
		heldActual,_:=wlmLmR27Build("r30_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		ca:=wlmLmR27PolicyActual(heldActual,corrected)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
		ua:=wlmLmR27PolicyActual(heldActual,uncorrected)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
		m["corrected_total_actual_advantage"]+=ca
		m["uncorrected_total_actual_advantage"]+=ua
		if ca>0 { m["corrected_positive_holdout_count"]++ } else if ca<0 { m["corrected_negative_holdout_count"]++ } else { m["corrected_zero_holdout_count"]++ }
		if ua>0 { m["uncorrected_positive_holdout_count"]++ } else if ua<0 { m["uncorrected_negative_holdout_count"]++ } else { m["uncorrected_zero_holdout_count"]++ }
		holdouts=append(holdouts,wlmLmR30Holdout{
			HoldoutManifest:inputs[hold].name,CorrectedSelection:corrected,UncorrectedSelection:uncorrected,
			CorrectedActualAdvantage:ca,UncorrectedActualAdvantage:ua,
			StableCandidateCount:stable,UnstableCandidateCount:unstable,
		})
	}
	m["case_count"]=float64(len(holdouts))
	m["corrected_delta_vs_uncorrected"]=m["corrected_total_actual_advantage"]-m["uncorrected_total_actual_advantage"]
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if len(holdouts)!=4 { m["invalid_row_count"]++ }
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmR30Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-HISTORICAL-SIGN-INSTABILITY-CORRECTION-R30",
		Metrics:m,Holdouts:holdouts,
	}
}
