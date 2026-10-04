package unitary

import (
	"math"
	"strings"
)

type wlmLmR31Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	TrainingPositiveCount int `json:"training_positive_count"`
	TrainingNegativeCount int `json:"training_negative_count"`
	TrainingZeroCount int `json:"training_zero_count"`
	R30Excluded bool `json:"r30_excluded"`
	MajorityPositiveExcluded bool `json:"majority_positive_excluded"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR31Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR31Case `json:"cases"`
}

type wlmLmR31ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

func RunWlmLmExternalFutureDataResidualCalibrationAttributionR31(
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
	inputs:=[]wlmLmR31ManifestInput{
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
		"non_equal_case_count":0,
		"stable_nonzero_case_count":0,
		"excluded_case_count":0,
		"majority_positive_excluded_case_count":0,
		"other_excluded_case_count":0,
		"majority_positive_excluded_coverage":0,
		"majority_positive_excluded_positive_rate":0,
		"other_excluded_positive_rate":0,
		"majority_positive_excluded_mean_actual_advantage":0,
		"other_excluded_mean_actual_advantage":0,
		"majority_positive_positive_rate_delta":0,
		"majority_positive_mean_advantage_delta":0,
		"heldout_outcome_use_before_flags":0,
		"post_result_mechanism_choice_count":0,
		"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"counter_overflow_count":0,
		"invalid_row_count":0,
	}
	cases:=make([]wlmLmR31Case,0,33)
	majorityPositiveHeldoutPositive:=0
	otherHeldoutPositive:=0
	majorityPositiveActualSum:=0.0
	otherActualSum:=0.0
	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,_:=wlmLmR27Build("r31_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man)
		}
		type frozenCase struct {
			alloc [3]int
			pos int
			neg int
			zero int
			majorityPositive bool
		}
		frozen:=make([]frozenCase,0,9)
		for _,alloc:=range wlmLmR27Candidates {
			if alloc==wlmLmR27Equal { continue }
			m["non_equal_case_count"]++
			signs:=map[int]bool{}
			pos,neg,zero:=0,0,0
			for _,tm:=range trainMans {
				adv:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				sign:=wlmLmR29Sign(adv)
				signs[sign]=true
				if sign>0 { pos++ } else if sign<0 { neg++ } else { zero++ }
			}
			stableNonZero:=len(signs)==1 && (signs[1] || signs[-1])
			if stableNonZero {
				m["stable_nonzero_case_count"]++
				continue
			}
			m["excluded_case_count"]++
			majorityPositive:=pos>=2
			if majorityPositive { m["majority_positive_excluded_case_count"]++ } else { m["other_excluded_case_count"]++ }
			frozen=append(frozen,frozenCase{alloc:alloc,pos:pos,neg:neg,zero:zero,majorityPositive:majorityPositive})
		}
		// All R30-equivalent exclusion flags and majority-positive subgroup flags are frozen above.
		// Held-out outcomes are opened only below.
		heldActual,_:=wlmLmR27Build("r31_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			if fc.majorityPositive {
				majorityPositiveActualSum+=actual
				if actual>0 { majorityPositiveHeldoutPositive++ }
			} else {
				otherActualSum+=actual
				if actual>0 { otherHeldoutPositive++ }
			}
			cases=append(cases,wlmLmR31Case{
				HoldoutManifest:inputs[hold].name,
				Allocation:fc.alloc,
				TrainingPositiveCount:fc.pos,
				TrainingNegativeCount:fc.neg,
				TrainingZeroCount:fc.zero,
				R30Excluded:true,
				MajorityPositiveExcluded:fc.majorityPositive,
				ActualAdvantage:actual,
			})
		}
	}
	excluded:=int(m["excluded_case_count"])
	maj:=int(m["majority_positive_excluded_case_count"])
	other:=int(m["other_excluded_case_count"])
	if excluded>0 { m["majority_positive_excluded_coverage"]=float64(maj)/float64(excluded) }
	if maj>0 {
		m["majority_positive_excluded_positive_rate"]=float64(majorityPositiveHeldoutPositive)/float64(maj)
		m["majority_positive_excluded_mean_actual_advantage"]=majorityPositiveActualSum/float64(maj)
	}
	if other>0 {
		m["other_excluded_positive_rate"]=float64(otherHeldoutPositive)/float64(other)
		m["other_excluded_mean_actual_advantage"]=otherActualSum/float64(other)
	}
	m["majority_positive_positive_rate_delta"]=m["majority_positive_excluded_positive_rate"]-m["other_excluded_positive_rate"]
	m["majority_positive_mean_advantage_delta"]=m["majority_positive_excluded_mean_actual_advantage"]-m["other_excluded_mean_actual_advantage"]
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 || excluded!=33 || len(cases)!=33 {
		m["invalid_row_count"]++
	}
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmR31Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-RESIDUAL-CALIBRATION-ATTRIBUTION-R31",
		Metrics:m,
		Cases:cases,
	}
}
