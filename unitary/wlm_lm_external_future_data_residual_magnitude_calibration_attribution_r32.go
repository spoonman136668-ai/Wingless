package unitary

import (
	"math"
	"strings"
)

type wlmLmR32Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	PredictedAdvantage float64 `json:"predicted_advantage"`
	TrainingNetAdvantage float64 `json:"training_net_advantage"`
	TrainingMaxAbsAdvantage float64 `json:"training_max_abs_advantage"`
	NetMagnitudeSignConflict bool `json:"net_magnitude_sign_conflict"`
	PredictionMagnitudeExtrapolation bool `json:"prediction_magnitude_extrapolation"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR32Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR32Case `json:"cases"`
}

type wlmLmR32ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

type wlmLmR32Agg struct {
	flagged int
	unflagged int
	flaggedSignErr int
	unflaggedSignErr int
	flaggedAbsErr float64
	unflaggedAbsErr float64
}

func wlmLmR32Update(a *wlmLmR32Agg, flag bool, pred, actual float64) {
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

func wlmLmR32Finalize(prefix string,a wlmLmR32Agg,total int,m map[string]float64) {
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

func RunWlmLmExternalFutureDataResidualMagnitudeCalibrationAttributionR32(
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
	inputs:=[]wlmLmR32ManifestInput{
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
		"heldout_outcome_use_before_flags":0,
		"post_result_mechanism_choice_count":0,
		"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"counter_overflow_count":0,
		"invalid_row_count":0,
	}
	var netAgg,predMagAgg wlmLmR32Agg
	cases:=make([]wlmLmR32Case,0,33)
	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainEx:=[]wlmLmR27Example{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r32_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man)
			trainEx=append(trainEx,ex...)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		beta,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["invalid_row_count"]++ }
		heldFeatures,_:=wlmLmR27Build("r32_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct {
			alloc [3]int
			pred float64
			net float64
			maxAbs float64
			netConflict bool
			predMagExtra bool
		}
		frozen:=make([]frozenCase,0,9)
		for _,alloc:=range wlmLmR27Candidates {
			if alloc==wlmLmR27Equal { continue }
			m["non_equal_case_count"]++
			signs:=map[int]bool{}
			net,maxAbs:=0.0,0.0
			for _,tm:=range trainMans {
				adv:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				signs[wlmLmR29Sign(adv)]=true
				net+=adv
				if math.Abs(adv)>maxAbs { maxAbs=math.Abs(adv) }
			}
			stableNonZero:=len(signs)==1 && (signs[1] || signs[-1])
			if stableNonZero {
				m["stable_nonzero_case_count"]++
				continue
			}
			m["excluded_case_count"]++
			pred:=wlmLmR27PolicyPred(heldFeatures,beta,alloc)-wlmLmR27PolicyPred(heldFeatures,beta,wlmLmR27Equal)
			netConflict:=wlmLmR29Sign(pred)!=wlmLmR29Sign(net)
			predMagExtra:=math.Abs(pred)>maxAbs
			frozen=append(frozen,frozenCase{alloc:alloc,pred:pred,net:net,maxAbs:maxAbs,netConflict:netConflict,predMagExtra:predMagExtra})
		}
		// All diagnostics and predictions above are frozen before held-out outcomes open.
		heldActual,_:=wlmLmR27Build("r32_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR32Update(&netAgg,fc.netConflict,fc.pred,actual)
			wlmLmR32Update(&predMagAgg,fc.predMagExtra,fc.pred,actual)
			cases=append(cases,wlmLmR32Case{
				HoldoutManifest:inputs[hold].name,
				Allocation:fc.alloc,
				PredictedAdvantage:fc.pred,
				TrainingNetAdvantage:fc.net,
				TrainingMaxAbsAdvantage:fc.maxAbs,
				NetMagnitudeSignConflict:fc.netConflict,
				PredictionMagnitudeExtrapolation:fc.predMagExtra,
				ActualAdvantage:actual,
			})
		}
	}
	wlmLmR32Finalize("net_magnitude_sign_conflict",netAgg,len(cases),m)
	wlmLmR32Finalize("prediction_magnitude_extrapolation",predMagAgg,len(cases),m)
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 || int(m["excluded_case_count"])!=33 || len(cases)!=33 {
		m["invalid_row_count"]++
	}
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmR32Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-RESIDUAL-MAGNITUDE-CALIBRATION-ATTRIBUTION-R32",
		Metrics:m,
		Cases:cases,
	}
}
