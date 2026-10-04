package unitary

import (
	"math"
	"strings"
)

type wlmLmR29Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	ReferencePrediction float64 `json:"reference_prediction"`
	ActualAdvantage float64 `json:"actual_advantage"`
	FeatureExtrapolation bool `json:"feature_extrapolation"`
	FeatureExcursionCount int `json:"feature_excursion_count"`
	FeatureExcursionMagnitude float64 `json:"feature_excursion_magnitude"`
	HistoricalSignInstability bool `json:"historical_sign_instability"`
	DomainContributionConflict bool `json:"domain_contribution_conflict"`
	DomainPredictedDeltas [3]float64 `json:"domain_predicted_deltas"`
}

type wlmLmR29Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR29Case `json:"cases"`
}

type wlmLmR29ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

type wlmLmR29Agg struct {
	flagged int
	unflagged int
	flaggedErr int
	unflaggedErr int
	flaggedAbs float64
	unflaggedAbs float64
	flaggedMag float64
}

func wlmLmR29Sign(v float64) int {
	if v>0 { return 1 }
	if v<0 { return -1 }
	return 0
}

func wlmLmR29PolicyVector(man wlmLmR27Manifest, alloc [3]int) [6]float64 {
	var out [6]float64
	for ai,arm:=range man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			if hi<=alloc[ai] {
				x:=arm.features[hi]
				for k:=0;k<6;k++ { out[k]+=x[k] }
			}
		}
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			if hi<=wlmLmR27Equal[ai] {
				x:=arm.features[hi]
				for k:=0;k<6;k++ { out[k]-=x[k] }
			}
		}
	}
	return out
}

func wlmLmR29DomainPred(man wlmLmR27Manifest,beta [7]float64,alloc [3]int,ai int) float64 {
	arm:=man.arms[ai]
	total:=0.0
	for i:=1;i<len(arm.budgets);i++ {
		hi:=arm.budgets[i]
		if hi<=alloc[ai] { total+=wlmLmR27Predict(beta,arm.features[hi]) }
		if hi<=wlmLmR27Equal[ai] { total-=wlmLmR27Predict(beta,arm.features[hi]) }
	}
	return total
}

func wlmLmR29FeatureExtrapolation(hold [6]float64, train [][6]float64) (bool,int,float64) {
	count:=0
	mag:=0.0
	for k:=0;k<6;k++ {
		lo,hi:=train[0][k],train[0][k]
		for i:=1;i<len(train);i++ {
			if train[i][k]<lo { lo=train[i][k] }
			if train[i][k]>hi { hi=train[i][k] }
		}
		rng:=hi-lo
		if rng<1e-12 { rng=1e-12 }
		if hold[k]<lo {
			count++
			mag+=(lo-hold[k])/rng
		} else if hold[k]>hi {
			count++
			mag+=(hold[k]-hi)/rng
		}
	}
	return count>0,count,mag
}

func wlmLmR29Update(a *wlmLmR29Agg,flag bool,pred,actual,mag float64) {
	err:=wlmLmR29Sign(pred)!=wlmLmR29Sign(actual)
	ae:=math.Abs(pred-actual)
	if flag {
		a.flagged++
		if err { a.flaggedErr++ }
		a.flaggedAbs+=ae
		a.flaggedMag+=mag
	} else {
		a.unflagged++
		if err { a.unflaggedErr++ }
		a.unflaggedAbs+=ae
	}
}

func wlmLmR29Finalize(prefix string,a wlmLmR29Agg,total int,m map[string]float64) {
	m[prefix+"_flagged_case_count"]=float64(a.flagged)
	m[prefix+"_unflagged_case_count"]=float64(a.unflagged)
	if total>0 { m[prefix+"_coverage"]=float64(a.flagged)/float64(total) }
	var fer,uer,fmae,umae float64
	if a.flagged>0 {
		fer=float64(a.flaggedErr)/float64(a.flagged)
		fmae=a.flaggedAbs/float64(a.flagged)
		m[prefix+"_feature_extrapolation_mean_magnitude"]=a.flaggedMag/float64(a.flagged)
	}
	if a.unflagged>0 {
		uer=float64(a.unflaggedErr)/float64(a.unflagged)
		umae=a.unflaggedAbs/float64(a.unflagged)
	}
	m[prefix+"_flagged_sign_error_rate"]=fer
	m[prefix+"_unflagged_sign_error_rate"]=uer
	m[prefix+"_flagged_mean_absolute_error"]=fmae
	m[prefix+"_unflagged_mean_absolute_error"]=umae
	m[prefix+"_sign_error_rate_delta"]=fer-uer
	if umae>0 { m[prefix+"_mae_ratio"]=fmae/umae }
}

func RunWlmLmExternalFutureDataCalibrationOutcomeMechanismAttributionR29(
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
	inputs:=[]wlmLmR29ManifestInput{
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
		"outer_holdout_manifest_count":4,"policies_per_holdout":10,"case_count":0,
		"heldout_outcome_use_before_flags":0,"post_result_mechanism_choice_count":0,
		"source_identity_mismatch_count":0,"capacity_growth_event_count":0,
		"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,
		"invalid_row_count":0,
	}
	cases:=make([]wlmLmR29Case,0,40)
	var exAgg,instAgg,confAgg wlmLmR29Agg
	for hold:=0;hold<4;hold++ {
		trainEx:=[]wlmLmR27Example{}
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainIDs:=make([]int,0,3)
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r29_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man)
			trainEx=append(trainEx,ex...)
			trainIDs=append(trainIDs,j)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		beta,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["invalid_row_count"]++ }
		heldFeatures,_:=wlmLmR27Build("r29_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct {
			alloc [3]int
			pred float64
			extrap bool
			exCount int
			exMag float64
			instability bool
			conflict bool
			domain [3]float64
		}
		frozen:=make([]frozenCase,0,10)
		for _,alloc:=range wlmLmR27Candidates {
			holdVec:=wlmLmR29PolicyVector(heldFeatures,alloc)
			trainVecs:=make([][6]float64,0,3)
			signs:=map[int]bool{}
			for _,tm:=range trainMans {
				trainVecs=append(trainVecs,wlmLmR29PolicyVector(tm,alloc))
				adv:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				signs[wlmLmR29Sign(adv)]=true
			}
			exflag,excount,exmag:=wlmLmR29FeatureExtrapolation(holdVec,trainVecs)
			inst:=len(signs)>1
			var domain [3]float64
			pos,neg:=false,false
			for ai:=0;ai<3;ai++ {
				domain[ai]=wlmLmR29DomainPred(heldFeatures,beta,alloc,ai)
				if domain[ai]>0 { pos=true }
				if domain[ai]<0 { neg=true }
			}
			conf:=pos&&neg
			pred:=wlmLmR27PolicyPred(heldFeatures,beta,alloc)-wlmLmR27PolicyPred(heldFeatures,beta,wlmLmR27Equal)
			frozen=append(frozen,frozenCase{alloc:alloc,pred:pred,extrap:exflag,exCount:excount,exMag:exmag,instability:inst,conflict:conf,domain:domain})
		}
		// Held-out outcomes are constructed only after all mechanism flags for this holdout are frozen.
		heldActual,_:=wlmLmR27Build("r29_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR29Update(&exAgg,fc.extrap,fc.pred,actual,fc.exMag)
			wlmLmR29Update(&instAgg,fc.instability,fc.pred,actual,0)
			wlmLmR29Update(&confAgg,fc.conflict,fc.pred,actual,0)
			cases=append(cases,wlmLmR29Case{
				HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,
				ReferencePrediction:fc.pred,ActualAdvantage:actual,
				FeatureExtrapolation:fc.extrap,FeatureExcursionCount:fc.exCount,FeatureExcursionMagnitude:fc.exMag,
				HistoricalSignInstability:fc.instability,DomainContributionConflict:fc.conflict,DomainPredictedDeltas:fc.domain,
			})
		}
		_ = trainIDs
	}
	m["case_count"]=float64(len(cases))
	wlmLmR29Finalize("feature_extrapolation",exAgg,40,m)
	wlmLmR29Finalize("historical_sign_instability",instAgg,40,m)
	wlmLmR29Finalize("domain_contribution_conflict",confAgg,40,m)
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if len(cases)!=40 { m["invalid_row_count"]++ }
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmR29Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-CALIBRATION-OUTCOME-MECHANISM-ATTRIBUTION-R29",
		Metrics:m,Cases:cases,
	}
}
