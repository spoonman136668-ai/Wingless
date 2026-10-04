package unitary

import (
	"math"
	"sort"
)

type wlmLmR33Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	FullPredictedAdvantage float64 `json:"full_predicted_advantage"`
	JackknifePredictedAdvantages [3]float64 `json:"jackknife_predicted_advantages"`
	JackknifeSignInstability bool `json:"jackknife_sign_instability"`
	JackknifeSpread float64 `json:"jackknife_spread"`
	HoldoutMedianJackknifeSpread float64 `json:"holdout_median_jackknife_spread"`
	HighRelativeJackknifeSpread bool `json:"high_relative_jackknife_spread"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR33Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR33Case `json:"cases"`
}

type wlmLmR33ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

type wlmLmR33Agg struct {
	flagged int
	unflagged int
	flaggedSignErr int
	unflaggedSignErr int
	flaggedAbsErr float64
	unflaggedAbsErr float64
}

func wlmLmR33Update(a *wlmLmR33Agg, flag bool, pred, actual float64) {
	signErr:=wlmLmR29Sign(pred)!=wlmLmR29Sign(actual)
	absErr:=math.Abs(pred-actual)
	if flag {
		a.flagged++
		if signErr { a.flaggedSignErr++ }
		a.flaggedAbsErr+=absErr
	} else {
		a.unflagged++
		if signErr { a.unflaggedSignErr++ }
		a.unflaggedAbsErr+=absErr
	}
}

func wlmLmR33Finalize(prefix string,a wlmLmR33Agg,total int,m map[string]float64) {
	m[prefix+"_flagged_count"]=float64(a.flagged)
	m[prefix+"_unflagged_count"]=float64(a.unflagged)
	if total>0 { m[prefix+"_coverage"]=float64(a.flagged)/float64(total) }
	flaggedRate,unflaggedRate:=0.0,0.0
	flaggedMAE,unflaggedMAE:=0.0,0.0
	if a.flagged>0 {
		flaggedRate=float64(a.flaggedSignErr)/float64(a.flagged)
		flaggedMAE=a.flaggedAbsErr/float64(a.flagged)
	}
	if a.unflagged>0 {
		unflaggedRate=float64(a.unflaggedSignErr)/float64(a.unflagged)
		unflaggedMAE=a.unflaggedAbsErr/float64(a.unflagged)
	}
	m[prefix+"_flagged_sign_error_rate"]=flaggedRate
	m[prefix+"_unflagged_sign_error_rate"]=unflaggedRate
	m[prefix+"_flagged_mean_absolute_error"]=flaggedMAE
	m[prefix+"_unflagged_mean_absolute_error"]=unflaggedMAE
	m[prefix+"_sign_error_rate_delta"]=flaggedRate-unflaggedRate
	if unflaggedMAE>0 { m[prefix+"_mae_ratio"]=flaggedMAE/unflaggedMAE }
}

func RunWlmLmExternalFutureDataCalibrationResidualStructureAttributionR33(
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
	inputs:=[]wlmLmR33ManifestInput{
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
		"jackknife_fit_failure_count":0,
		"heldout_outcome_use_before_flags":0,
		"post_result_mechanism_choice_count":0,
		"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"counter_overflow_count":0,
		"invalid_row_count":0,
	}
	var signAgg,spreadAgg wlmLmR33Agg
	cases:=make([]wlmLmR33Case,0,33)
	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainEx:=[]wlmLmR27Example{}
		trainIDs:=[]int{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r33_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man)
			trainEx=append(trainEx,ex...)
			trainIDs=append(trainIDs,j)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		full,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["invalid_row_count"]++ }
		var jack [3][7]float64
		for i,id:=range trainIDs {
			beta,jok:=wlmLmR27Fit(trainEx,id)
			if !jok {
				m["jackknife_fit_failure_count"]++
				m["invalid_row_count"]++
			}
			jack[i]=beta
		}
		heldFeatures,_:=wlmLmR27Build("r33_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct {
			alloc [3]int
			excluded bool
			fullPred float64
			jackPred [3]float64
			signInstability bool
			spread float64
			highSpread bool
		}
		frozen:=make([]frozenCase,0,9)
		spreads:=make([]float64,0,9)
		for _,alloc:=range wlmLmR27Candidates {
			if alloc==wlmLmR27Equal { continue }
			m["non_equal_case_count"]++
			signs:=map[int]bool{}
			for _,tm:=range trainMans {
				adv:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				signs[wlmLmR29Sign(adv)]=true
			}
			stableNonZero:=len(signs)==1 && (signs[1] || signs[-1])
			if stableNonZero { m["stable_nonzero_case_count"]++ } else { m["excluded_case_count"]++ }
			fullPred:=wlmLmR27PolicyPred(heldFeatures,full,alloc)-wlmLmR27PolicyPred(heldFeatures,full,wlmLmR27Equal)
			var jp [3]float64
			minP,maxP:=math.Inf(1),math.Inf(-1)
			for i:=0;i<3;i++ {
				jp[i]=wlmLmR27PolicyPred(heldFeatures,jack[i],alloc)-wlmLmR27PolicyPred(heldFeatures,jack[i],wlmLmR27Equal)
				if jp[i]<minP { minP=jp[i] }
				if jp[i]>maxP { maxP=jp[i] }
			}
			s0:=wlmLmR29Sign(jp[0])
			unanimousNonZero:=s0!=0
			for i:=1;i<3;i++ {
				if wlmLmR29Sign(jp[i])!=s0 { unanimousNonZero=false }
			}
			spread:=maxP-minP
			if math.IsNaN(spread)||math.IsInf(spread,0) { m["invalid_row_count"]++ }
			spreads=append(spreads,spread)
			frozen=append(frozen,frozenCase{
				alloc:alloc,excluded:!stableNonZero,fullPred:fullPred,jackPred:jp,
				signInstability:!unanimousNonZero,spread:spread,
			})
		}
		if len(spreads)!=9 { m["invalid_row_count"]++ }
		sorted:=append([]float64(nil),spreads...)
		sort.Float64s(sorted)
		median:=0.0
		if len(sorted)==9 { median=sorted[4] }
		for i:=range frozen { frozen[i].highSpread=frozen[i].spread>median }

		heldActual,_:=wlmLmR27Build("r33_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			if !fc.excluded { continue }
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR33Update(&signAgg,fc.signInstability,fc.fullPred,actual)
			wlmLmR33Update(&spreadAgg,fc.highSpread,fc.fullPred,actual)
			cases=append(cases,wlmLmR33Case{
				HoldoutManifest:inputs[hold].name,
				Allocation:fc.alloc,
				FullPredictedAdvantage:fc.fullPred,
				JackknifePredictedAdvantages:fc.jackPred,
				JackknifeSignInstability:fc.signInstability,
				JackknifeSpread:fc.spread,
				HoldoutMedianJackknifeSpread:median,
				HighRelativeJackknifeSpread:fc.highSpread,
				ActualAdvantage:actual,
			})
		}
	}
	wlmLmR33Finalize("jackknife_sign_instability",signAgg,len(cases),m)
	wlmLmR33Finalize("high_relative_jackknife_spread",spreadAgg,len(cases),m)
	for k,v:=range m {
		if len(k)>31 && k[len(k)-31:]=="_source_identity_mismatch_count" && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 || int(m["excluded_case_count"])!=33 || len(cases)!=33 {
		m["invalid_row_count"]++
	}
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmR33Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-CALIBRATION-RESIDUAL-STRUCTURE-ATTRIBUTION-R33",
		Metrics:m,
		Cases:cases,
	}
}
