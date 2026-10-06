package unitary

import (
	"math"
	"strings"
)

type wlmLmR48ArmSummary struct {
	ID string `json:"id"`
	Alphas [6]float64 `json:"alphas"`
	PairwiseWinnerAccuracy float64 `json:"pairwise_winner_accuracy"`
	MeanDecisionRegret float64 `json:"mean_decision_regret"`
	NonTiePairCount int `json:"non_tie_pair_count"`
	CorrectPairCount int `json:"correct_pair_count"`
	ManifestAccuracy map[string]float64 `json:"manifest_accuracy"`
	AccuracyDeltaVsFull float64 `json:"accuracy_delta_vs_full"`
	PositiveManifestCountVsFull int `json:"positive_manifest_count_vs_full"`
}

type wlmLmR48Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	BestReset string `json:"best_reset"`
	Arms []wlmLmR48ArmSummary `json:"arms"`
	Baseline wlmLmR48ArmSummary `json:"baseline"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmR48EvalArm(id string, alpha [6]float64, hist []wlmLmR47Prepared, unseen []wlmLmR35ManifestInput, m map[string]float64) wlmLmR48ArmSummary {
	examples:=[]wlmLmR47Example{}
	for i,p:=range hist {
		st:=wlmLmR47States(p,alpha,m)
		examples=append(examples,wlmLmR47Examples(p,st,i)...)
	}
	beta,ok:=wlmLmR47Fit(examples,-1)
	if !ok { m["fit_failure_count"]++ }
	actions:=wlmLmR47Actions(m)
	correct,total:=0,0
	regret:=0.0
	per:=map[string]float64{}
	for ui,in:=range unseen {
		pred:=wlmLmR47PreparedFrom("r48_pred_"+id+"_"+in.name,in.sources,4+ui,false,m)
		st:=wlmLmR47States(pred,alpha,m)
		scores:=map[[3]int]float64{}
		for _,a:=range actions { scores[a]=wlmLmR47PolicyPred(pred,st,beta,a) }
		actual:=wlmLmR47PreparedFrom("r48_actual_"+id+"_"+in.name,in.sources,4+ui,true,m)
		localCorrect,localN:=0,0
		for i:=0;i<len(actions);i++ {
			for j:=i+1;j<len(actions);j++ {
				a,b:=actions[i],actions[j]
				aa,ab:=wlmLmR27PolicyActual(actual.man,a),wlmLmR27PolicyActual(actual.man,b)
				if aa==ab { continue }
				localN++;total++
				want:=a
				if ab>aa { want=b }
				got:=wlmLmR44Choice(a,b,scores[a],scores[b])
				if got==want { correct++;localCorrect++ } else { regret+=math.Abs(aa-ab) }
			}
		}
		if localN==0 { m["invalid_row_count"]++ } else { per[in.name]=float64(localCorrect)/float64(localN) }
	}
	acc:=0.0;mr:=0.0
	if total==0 { m["invalid_row_count"]++ } else { acc=float64(correct)/float64(total);mr=regret/float64(total) }
	return wlmLmR48ArmSummary{ID:id,Alphas:alpha,PairwiseWinnerAccuracy:acc,MeanDecisionRegret:mr,NonTiePairCount:total,CorrectPairCount:correct,ManifestAccuracy:per}
}

func RunWlmLmExternalFutureDataLearnedStateFormerFeatureAttributionR48(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse,
	eighteenthCode,eighteenthStructured,eighteenthProse,
	nineteenthCode,nineteenthStructured,nineteenthProse,
	twentiethCode,twentiethStructured,twentiethProse []byte,
) interface{} {
	budCode:=[]int{0,436,582,727,872}
	budOther:=[]int{0,291,436,581,726}
	m:=map[string]float64{
		"historical_manifest_count":4,
		"unseen_manifest_count":3,
		"attribution_arm_count":4,
		"baseline_anchor_count":1,
		"state_dimension":6,
		"readout_parameter_count":7,
		"total_adaptation_budget":1744,
		"candidate_pair_count_per_arm":108,
		"heldout_outcome_use_before_arm_freeze":0,
		"post_result_arm_or_threshold_choice_count":0,
		"capacity_growth_event_count":0,
		"external_model_call_count":0,
		"tokenizer_use_count":0,
		"fixed_budget_mismatch_count":0,
		"raw_state_source_mismatch_count":0,
		"source_identity_mismatch_count":0,
		"fit_failure_count":0,
		"invalid_row_count":0,
	}
	mk:=func(d string,b []byte,h string,n int,bud []int)wlmLmR27Source{return wlmLmR27Source{d:d,b:b,h:h,n:n,budgets:bud}}
	histInputs:=[]wlmLmR35ManifestInput{
		{"transfer",[]wlmLmR27Source{mk("code",transferCode,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45",41453,budCode),mk("structured",transferStructured,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932",14365,budOther),mk("technical_prose",transferProse,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b",1454,budOther)}},
		{"third",[]wlmLmR27Source{mk("code",thirdCode,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff",41453,budCode),mk("structured",thirdStructured,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56",14365,budOther),mk("technical_prose",thirdProse,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017",1454,budOther)}},
		{"fourth",[]wlmLmR27Source{mk("code",fourthCode,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1",41453,budCode),mk("structured",fourthStructured,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933",14365,budOther),mk("technical_prose",fourthProse,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222",1454,budOther)}},
		{"fifth",[]wlmLmR27Source{mk("code",fifthCode,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a",41453,budCode),mk("structured",fifthStructured,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281",14365,budOther),mk("technical_prose",fifthProse,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49",1454,budOther)}},
	}
	hist:=make([]wlmLmR47Prepared,0,4)
	for i,in:=range histInputs { hist=append(hist,wlmLmR47PreparedFrom("r48_hist_"+in.name,in.sources,i,true,m)) }
	unseen:=[]wlmLmR35ManifestInput{
		{"eighteenth",[]wlmLmR27Source{mk("code",eighteenthCode,"2d1317ffe8676e0ce80220112dc43e633d4638b21dd70df4a548cfca3233c2c5",41453,budCode),mk("structured",eighteenthStructured,"92065652c4a6412e229bd0a7254aff2ff654daa6282e112521b8795e4c4c61f7",14365,budOther),mk("technical_prose",eighteenthProse,"69f42897f882ddfa2afac098e2f32f30f4e93ef4fdcbfab7aff459f9b6946d52",1454,budOther)}},
		{"nineteenth",[]wlmLmR27Source{mk("code",nineteenthCode,"59fe6feae1755891c443a3ec58b36251428c47debf918bdb5d00ef841b884cde",41453,budCode),mk("structured",nineteenthStructured,"c674adc46f03be25284f1296698d52a3b749f61d7513831ef9987a07111b8a00",14365,budOther),mk("technical_prose",nineteenthProse,"1be40e2cee5a5e2dbb1636b822309f8d0e18d60689454267627b8a0ab73261b1",1454,budOther)}},
		{"twentieth",[]wlmLmR27Source{mk("code",twentiethCode,"09d797fcc5d55c2611781588c2c2b594e429062664d3dac1b86f5033fc5eec32",41453,budCode),mk("structured",twentiethStructured,"695e6385b752759f21ebccd45ec9217137b8fee42212a81218f91698e6d17bee",14365,budOther),mk("technical_prose",twentiethProse,"c47c609d3766e4eb0e95cec4aa3978271d00cbd9677aa44bdc0f165d1ab09c21",1454,budOther)}},
	}
	armDefs:=[]struct{id string;alpha [6]float64}{
		{"full-selected",[6]float64{0.75,0.75,0.125,0.125,1,1}},
		{"reset-distribution",[6]float64{0.5,0.5,0.125,0.125,1,1}},
		{"reset-character",[6]float64{0.75,0.75,0.5,0.5,1,1}},
		{"reset-layout-repetition",[6]float64{0.75,0.75,0.125,0.125,0.5,0.5}},
	}
	arms:=make([]wlmLmR48ArmSummary,0,4)
	for _,d:=range armDefs { arms=append(arms,wlmLmR48EvalArm(d.id,d.alpha,hist,unseen,m)) }
	baseline:=wlmLmR48EvalArm("r46-alpha-0.5",[6]float64{0.5,0.5,0.5,0.5,0.5,0.5},hist,unseen,m)
	full:=arms[0]
	bestReset:=""
	bestDelta:=math.Inf(-1)
	secondDelta:=math.Inf(-1)
	for i:=1;i<len(arms);i++ {
		arms[i].AccuracyDeltaVsFull=arms[i].PairwiseWinnerAccuracy-full.PairwiseWinnerAccuracy
		pos:=0
		for name,acc:=range arms[i].ManifestAccuracy { if acc>full.ManifestAccuracy[name] { pos++ } }
		arms[i].PositiveManifestCountVsFull=pos
		d:=arms[i].AccuracyDeltaVsFull
		if d>bestDelta || (d==bestDelta && (bestReset=="" || arms[i].ID<bestReset)) {
			secondDelta=bestDelta;bestDelta=d;bestReset=arms[i].ID
		} else if d>secondDelta { secondDelta=d }
	}
	if math.IsInf(secondDelta,-1) { secondDelta=bestDelta }
	m["best_reset_accuracy_delta"]=bestDelta
	m["best_vs_second_margin"]=bestDelta-secondDelta
	m["baseline_minus_full_accuracy"]=baseline.PairwiseWinnerAccuracy-full.PairwiseWinnerAccuracy
	for _,a:=range arms { if a.ID==bestReset { m["best_reset_positive_manifest_count"]=float64(a.PositiveManifestCountVsFull) } }
	m["total_eval_pair_count"]=float64((len(arms)+1)*108)
	for _,a:=range arms {
		if a.NonTiePairCount<72 || len(a.ManifestAccuracy)!=3 { m["invalid_row_count"]++ }
	}
	if baseline.NonTiePairCount<72 || len(baseline.ManifestAccuracy)!=3 { m["invalid_row_count"]++ }
	for k,v:=range m { if strings.HasSuffix(k,"_source_identity_mismatch_count")&&k!="source_identity_mismatch_count" { m["source_identity_mismatch_count"]+=v } }
	if len(arms)!=4||m["state_dimension"]!=6||m["readout_parameter_count"]!=7||m["total_adaptation_budget"]!=1744 { m["invalid_row_count"]++ }
	if m["fixed_budget_mismatch_count"]!=0||m["fit_failure_count"]!=0||m["raw_state_source_mismatch_count"]!=0||m["source_identity_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR48Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-FEATURE-ATTRIBUTION-R48",BestReset:bestReset,Arms:arms,Baseline:baseline,Metrics:m}
}
