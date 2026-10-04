package unitary

import (
	"math"
	"strings"
)

type wlmLmR38Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	RawBaselinePredictedAdvantage float64 `json:"raw_baseline_predicted_advantage"`
	RawBaselineSign int `json:"raw_baseline_sign"`
	TrainingPositiveVotes int `json:"training_positive_votes"`
	TrainingNegativeVotes int `json:"training_negative_votes"`
	MajorityCovered bool `json:"majority_covered"`
	MajoritySign int `json:"majority_sign"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR38Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR38Case `json:"cases"`
}

type wlmLmR38SignAgg struct {
	n int
	err int
}

func wlmLmR38Update(a *wlmLmR38SignAgg,predSign int,actual float64) {
	a.n++
	if predSign!=wlmLmR29Sign(actual) { a.err++ }
}

func wlmLmR38Rate(a wlmLmR38SignAgg,m map[string]float64,key string) {
	if a.n==0 { m[key]=0;return }
	m[key]=float64(a.err)/float64(a.n)
}

func RunWlmLmExternalFutureDataCalibrationLabelStabilityAttributionR38(
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
		"majority_covered_case_count":0,"majority_uncovered_case_count":0,"majority_coverage":0,
		"raw_baseline_sign_error_rate":0,"raw_covered_sign_error_rate":0,"raw_uncovered_sign_error_rate":0,
		"majority_sign_error_rate":0,"majority_sign_error_improvement":0,
		"raw_baseline_identity_mismatch_count":0,"consensus_definition_mismatch_count":0,
		"fit_failure_count":0,"heldout_outcome_use_before_consensus_freeze":0,
		"post_result_consensus_choice_count":0,"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	var rawAll,rawCovered,rawUncovered,majority wlmLmR38SignAgg
	cases:=make([]wlmLmR38Case,0,33)
	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainEx:=[]wlmLmR27Example{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r38_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man);trainEx=append(trainEx,ex...)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		baseR27,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		rawBeta,ok0:=wlmLmR34Fit(wlmLmR35Examples(trainMans,false))
		if !ok0 { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		heldFeatures,_:=wlmLmR27Build("r38_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)

		type frozenCase struct{
			alloc [3]int
			excluded bool
			rawPred float64
			posVotes,negVotes int
			covered bool
			majoritySign int
		}
		frozen:=make([]frozenCase,0,9)
		for _,alloc:=range wlmLmR27Candidates {
			if alloc==wlmLmR27Equal { continue }
			m["non_equal_case_count"]++
			signs:=map[int]bool{}
			pos,neg:=0,0
			for _,tm:=range trainMans {
				adv:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				s:=wlmLmR29Sign(adv)
				signs[s]=true
				if s>0 { pos++ } else if s<0 { neg++ }
			}
			stableNonZero:=len(signs)==1 && (signs[1]||signs[-1])
			if stableNonZero { m["stable_nonzero_case_count"]++ } else { m["excluded_case_count"]++ }

			covered:=false;maj:=0
			if pos>=2 && neg>=2 { m["consensus_definition_mismatch_count"]++ }
			if pos>=2 { covered=true;maj=1 }
			if neg>=2 {
				if covered { m["consensus_definition_mismatch_count"]++ }
				covered=true;maj=-1
			}
			if covered && maj==0 { m["consensus_definition_mismatch_count"]++ }
			if !covered && maj!=0 { m["consensus_definition_mismatch_count"]++ }

			rawPred:=wlmLmR35PolicyPred(heldFeatures,rawBeta,false,alloc)-wlmLmR35PolicyPred(heldFeatures,rawBeta,false,wlmLmR27Equal)
			r27:=wlmLmR27PolicyPred(heldFeatures,baseR27,alloc)-wlmLmR27PolicyPred(heldFeatures,baseR27,wlmLmR27Equal)
			if math.Abs(rawPred-r27)>1e-9 { m["raw_baseline_identity_mismatch_count"]++ }
			frozen=append(frozen,frozenCase{alloc:alloc,excluded:!stableNonZero,rawPred:rawPred,posVotes:pos,negVotes:neg,covered:covered,majoritySign:maj})
		}

		heldActual,_:=wlmLmR27Build("r38_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			if !fc.excluded { continue }
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			rawSign:=wlmLmR29Sign(fc.rawPred)
			wlmLmR38Update(&rawAll,rawSign,actual)
			if fc.covered {
				m["majority_covered_case_count"]++
				wlmLmR38Update(&rawCovered,rawSign,actual)
				wlmLmR38Update(&majority,fc.majoritySign,actual)
			}else{
				m["majority_uncovered_case_count"]++
				wlmLmR38Update(&rawUncovered,rawSign,actual)
			}
			cases=append(cases,wlmLmR38Case{
				HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,
				RawBaselinePredictedAdvantage:fc.rawPred,RawBaselineSign:rawSign,
				TrainingPositiveVotes:fc.posVotes,TrainingNegativeVotes:fc.negVotes,
				MajorityCovered:fc.covered,MajoritySign:fc.majoritySign,ActualAdvantage:actual,
			})
		}
	}
	wlmLmR38Rate(rawAll,m,"raw_baseline_sign_error_rate")
	wlmLmR38Rate(rawCovered,m,"raw_covered_sign_error_rate")
	wlmLmR38Rate(rawUncovered,m,"raw_uncovered_sign_error_rate")
	wlmLmR38Rate(majority,m,"majority_sign_error_rate")
	if m["excluded_case_count"]>0 { m["majority_coverage"]=m["majority_covered_case_count"]/m["excluded_case_count"] }
	if m["majority_covered_case_count"]>0 {
		m["majority_sign_error_improvement"]=m["raw_covered_sign_error_rate"]-m["majority_sign_error_rate"]
	}
	if math.Abs(m["raw_baseline_sign_error_rate"]-0.5151515151515151)>1e-9 {
		m["raw_baseline_identity_mismatch_count"]++
	}
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 || int(m["excluded_case_count"])!=33 || len(cases)!=33 {
		m["invalid_row_count"]++
	}
	if int(m["majority_covered_case_count"]+m["majority_uncovered_case_count"])!=33 { m["invalid_row_count"]++ }
	if m["raw_baseline_identity_mismatch_count"]!=0 || m["consensus_definition_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR38Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-CALIBRATION-LABEL-STABILITY-ATTRIBUTION-R38",
		Metrics:m,Cases:cases,
	}
}
