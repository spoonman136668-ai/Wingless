package unitary

import (
	"math"
	"strings"
)

type wlmLmR43Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	AllocationState string `json:"allocation_state"`
	GlobalPredictedAdvantage float64 `json:"global_predicted_advantage"`
	StateOffset float64 `json:"state_offset"`
	StateConditionalPredictedAdvantage float64 `json:"state_conditional_predicted_advantage"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR43Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR43Case `json:"cases"`
}

func wlmLmR43StateIndex(alloc [3]int,m map[string]float64) int {
	switch alloc[0] {
	case 582:
		return 0
	case 727:
		return 1
	case 872:
		return 2
	default:
		m["state_definition_mismatch_count"]++
		return 0
	}
}

func wlmLmR43StateName(i int) string {
	if i==0 { return "code_budget_582" }
	if i==1 { return "code_budget_727" }
	return "code_budget_872"
}

func RunWlmLmExternalFutureDataStateConditionalResponseAttributionR43(
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
		"state_conditional_sign_error_rate":0,"state_conditional_mean_absolute_error":0,
		"state_conditional_sign_error_improvement":0,"state_conditional_mae_ratio":0,
		"state_582_case_count":0,"state_727_case_count":0,"state_872_case_count":0,
		"state_582_offset":0,"state_727_offset":0,"state_872_offset":0,
		"base6_identity_mismatch_count":0,"state_definition_mismatch_count":0,
		"state_stratum_insufficient_count":0,"fit_failure_count":0,
		"heldout_outcome_use_before_state_condition_freeze":0,
		"post_result_state_condition_choice_count":0,"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	var globalAgg,conditionalAgg wlmLmR34Agg
	var offsetTotals [3]float64
	cases:=make([]wlmLmR43Case,0,33)

	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainEx:=[]wlmLmR27Example{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r43_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
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
				d:=wlmLmR43StateIndex(alloc,m)
				pred:=wlmLmR35PolicyPred(tm,baseBeta,false,alloc)-wlmLmR35PolicyPred(tm,baseBeta,false,wlmLmR27Equal)
				actual:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				residualSum[d]+=actual-pred
				residualN[d]++
			}
		}
		var offsets [3]float64
		totalN:=0
		for d:=0;d<3;d++ {
			totalN+=residualN[d]
			if residualN[d]<3 {
				m["state_stratum_insufficient_count"]++
				m["invalid_row_count"]++
				continue
			}
			offsets[d]=residualSum[d]/float64(residualN[d])
			if math.IsNaN(offsets[d])||math.IsInf(offsets[d],0) {
				m["state_stratum_insufficient_count"]++
				m["invalid_row_count"]++
			}
			offsetTotals[d]+=offsets[d]
		}
		if totalN!=27 {
			m["state_definition_mismatch_count"]++
			m["invalid_row_count"]++
		}

		heldFeatures,_:=wlmLmR27Build("r43_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct {
			alloc [3]int
			excluded bool
			state int
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
			d:=wlmLmR43StateIndex(alloc,m)
			frozen=append(frozen,frozenCase{
				alloc:alloc,excluded:!stableNonZero,state:d,
				globalPred:globalPred,conditionalPred:globalPred+offsets[d],offset:offsets[d],
			})
		}

		heldActual,_:=wlmLmR27Build("r43_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			if !fc.excluded { continue }
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR34Update(&globalAgg,fc.globalPred,actual)
			wlmLmR34Update(&conditionalAgg,fc.conditionalPred,actual)
			if fc.state==0 { m["state_582_case_count"]++ }
			if fc.state==1 { m["state_727_case_count"]++ }
			if fc.state==2 { m["state_872_case_count"]++ }
			cases=append(cases,wlmLmR43Case{
				HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,
				AllocationState:wlmLmR43StateName(fc.state),
				GlobalPredictedAdvantage:fc.globalPred,
				StateOffset:fc.offset,
				StateConditionalPredictedAdvantage:fc.conditionalPred,
				ActualAdvantage:actual,
			})
		}
	}

	wlmLmR34Finalize("global",globalAgg,m)
	wlmLmR34Finalize("state_conditional",conditionalAgg,m)
	m["state_conditional_sign_error_improvement"]=m["global_sign_error_rate"]-m["state_conditional_sign_error_rate"]
	m["state_conditional_mae_ratio"]=wlmLmR34Ratio(m["state_conditional_mean_absolute_error"],m["global_mean_absolute_error"],m)
	m["state_582_offset"]=offsetTotals[0]/4.0
	m["state_727_offset"]=offsetTotals[1]/4.0
	m["state_872_offset"]=offsetTotals[2]/4.0

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
		int(m["state_582_case_count"]+m["state_727_case_count"]+m["state_872_case_count"])!=33 {
		m["invalid_row_count"]++
	}
	if m["base6_identity_mismatch_count"]!=0 || m["state_definition_mismatch_count"]!=0 ||
		m["state_stratum_insufficient_count"]!=0 || m["fit_failure_count"]!=0 {
		m["invalid_row_count"]++
	}
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmR43Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-STATE-CONDITIONED-RESPONSE-ATTRIBUTION-R43",
		Metrics:m,Cases:cases,
	}
}
