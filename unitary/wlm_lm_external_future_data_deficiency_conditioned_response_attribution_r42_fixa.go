package unitary

import (
	"math"
	"strings"
)

type wlmLmR42FixACase struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	DeficiencyDomain string `json:"deficiency_domain"`
	GlobalPredictedAdvantage float64 `json:"global_predicted_advantage"`
	DeficiencyDomainOffset float64 `json:"deficiency_domain_offset"`
	DeficiencyConditionalPredictedAdvantage float64 `json:"deficiency_conditional_predicted_advantage"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR42FixAResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR42FixACase `json:"cases"`
}

func wlmLmR42FixADeficiencyDomain(alloc [3]int,m map[string]float64) int {
	best:=1
	bestDelta:=wlmLmR27Equal[1]-alloc[1]
	for i:=2;i<3;i++ {
		d:=wlmLmR27Equal[i]-alloc[i]
		if d>bestDelta {
			best=i
			bestDelta=d
		}
	}
	if bestDelta<=0 {
		m["deficiency_domain_definition_mismatch_count"]++
		return 1
	}
	return best
}

func wlmLmR42FixADomainName(i int) string {
	if i==0 { return "code" }
	if i==1 { return "structured" }
	return "technical_prose"
}

func RunWlmLmExternalFutureDataDeficiencyConditionalResponseAttributionR42FixA(
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
		"deficiency_conditional_sign_error_rate":0,"deficiency_conditional_mean_absolute_error":0,
		"deficiency_conditional_sign_error_improvement":0,"deficiency_conditional_mae_ratio":0,
		"code_case_count":0,"structured_case_count":0,"technical_prose_case_count":0,
		"code_offset":0,"structured_offset":0,"technical_prose_offset":0,
		"base6_identity_mismatch_count":0,"deficiency_domain_definition_mismatch_count":0,
		"target_stratum_insufficient_count":0,"fit_failure_count":0,
		"heldout_outcome_use_before_deficiency_condition_freeze":0,
		"post_result_deficiency_condition_choice_count":0,"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	var globalAgg,conditionalAgg wlmLmR34Agg
	var offsetTotals [3]float64
	cases:=make([]wlmLmR42FixACase,0,33)

	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainEx:=[]wlmLmR27Example{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r42_fixa_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man)
			trainEx=append(trainEx,ex...)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		baseR27,ok:=wlmLmR27Fit(trainEx,-1)
		baseBeta,ok0:=wlmLmR34Fit(wlmLmR35Examples(trainMans,false))
		if !ok || !ok0 {
			m["fit_failure_count"]++
			m["invalid_row_count"]++
		}

		var residualSum [3]float64
		var residualN [3]int
		for _,tm:=range trainMans {
			for _,alloc:=range wlmLmR27Candidates {
				if alloc==wlmLmR27Equal { continue }
				d:=wlmLmR42FixADeficiencyDomain(alloc,m)
				pred:=wlmLmR35PolicyPred(tm,baseBeta,false,alloc)-wlmLmR35PolicyPred(tm,baseBeta,false,wlmLmR27Equal)
				actual:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				residualSum[d]+=actual-pred
				residualN[d]++
			}
		}
		var offsets [3]float64
		totalN:=0
		for d:=1;d<3;d++ {
			totalN+=residualN[d]
			if residualN[d]<3 {
				m["target_stratum_insufficient_count"]++
				m["invalid_row_count"]++
				continue
			}
			offsets[d]=residualSum[d]/float64(residualN[d])
			if math.IsNaN(offsets[d])||math.IsInf(offsets[d],0) {
				m["target_stratum_insufficient_count"]++
				m["invalid_row_count"]++
			}
			offsetTotals[d]+=offsets[d]
		}
		if totalN!=27 {
			m["deficiency_domain_definition_mismatch_count"]++
			m["invalid_row_count"]++
		}

		heldFeatures,_:=wlmLmR27Build("r42_fixa_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct {
			alloc [3]int
			excluded bool
			domain int
			globalPred,conditionalPred,offset float64
		}
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

			globalPred:=wlmLmR35PolicyPred(heldFeatures,baseBeta,false,alloc)-wlmLmR35PolicyPred(heldFeatures,baseBeta,false,wlmLmR27Equal)
			r27:=wlmLmR27PolicyPred(heldFeatures,baseR27,alloc)-wlmLmR27PolicyPred(heldFeatures,baseR27,wlmLmR27Equal)
			if math.Abs(globalPred-r27)>1e-9 { m["base6_identity_mismatch_count"]++ }
			d:=wlmLmR42FixADeficiencyDomain(alloc,m)
			frozen=append(frozen,frozenCase{
				alloc:alloc,excluded:!stableNonZero,domain:d,
				globalPred:globalPred,conditionalPred:globalPred+offsets[d],offset:offsets[d],
			})
		}

		heldActual,_:=wlmLmR27Build("r42_fixa_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			if !fc.excluded { continue }
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR34Update(&globalAgg,fc.globalPred,actual)
			wlmLmR34Update(&conditionalAgg,fc.conditionalPred,actual)
			if fc.domain==0 { m["code_case_count"]++; m["deficiency_domain_definition_mismatch_count"]++ }
			if fc.domain==1 { m["structured_case_count"]++ }
			if fc.domain==2 { m["technical_prose_case_count"]++ }
			cases=append(cases,wlmLmR42FixACase{
				HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,
				DeficiencyDomain:wlmLmR42FixADomainName(fc.domain),
				GlobalPredictedAdvantage:fc.globalPred,
				DeficiencyDomainOffset:fc.offset,
				DeficiencyConditionalPredictedAdvantage:fc.conditionalPred,
				ActualAdvantage:actual,
			})
		}
	}

	wlmLmR34Finalize("global",globalAgg,m)
	wlmLmR34Finalize("deficiency_conditional",conditionalAgg,m)
	m["deficiency_conditional_sign_error_improvement"]=m["global_sign_error_rate"]-m["deficiency_conditional_sign_error_rate"]
	m["deficiency_conditional_mae_ratio"]=wlmLmR34Ratio(m["deficiency_conditional_mean_absolute_error"],m["global_mean_absolute_error"],m)
	m["code_offset"]=0
	m["structured_offset"]=offsetTotals[1]/4.0
	m["technical_prose_offset"]=offsetTotals[2]/4.0

	if math.Abs(m["global_sign_error_rate"]-0.5151515151515151)>1e-9 ||
		math.Abs(m["global_mean_absolute_error"]-10.731783167689903)>1e-9 {
		m["base6_identity_mismatch_count"]++
	}
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 ||
		int(m["excluded_case_count"])!=33 || len(cases)!=33 || globalAgg.n!=33 || conditionalAgg.n!=33 ||
		int(m["code_case_count"])!=0 || int(m["structured_case_count"]+m["technical_prose_case_count"])!=33 {
		m["invalid_row_count"]++
	}
	if m["base6_identity_mismatch_count"]!=0 || m["deficiency_domain_definition_mismatch_count"]!=0 ||
		m["target_stratum_insufficient_count"]!=0 || m["fit_failure_count"]!=0 {
		m["invalid_row_count"]++
	}
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmR42FixAResult{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-DEFICIENCY-CONDITIONED-RESPONSE-ATTRIBUTION-R42-FIXA",
		Metrics:m,Cases:cases,
	}
}
