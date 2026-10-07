package unitary

import (
	"math"
	"strings"
)

type wlmLmR49Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	BestPairwiseReset string `json:"best_pairwise_reset"`
	Arms []wlmLmR48ArmSummary `json:"arms"`
	Baseline wlmLmR48ArmSummary `json:"baseline"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalFutureDataLearnedStateFormerInteractionAttributionR49(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse,
	twentyFirstCode,twentyFirstStructured,twentyFirstProse,
	twentySecondCode,twentySecondStructured,twentySecondProse,
	twentyThirdCode,twentyThirdStructured,twentyThirdProse []byte,
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
	for i,in:=range histInputs { hist=append(hist,wlmLmR47PreparedFrom("r49_hist_"+in.name,in.sources,i,true,m)) }
	unseen:=[]wlmLmR35ManifestInput{
		{"twenty-first",[]wlmLmR27Source{mk("code",twentyFirstCode,"3ae0150ba9a53b5dfc24480ace2d7439a9f72107e7da57eb05d95060e1cf3cca",41453,budCode),mk("structured",twentyFirstStructured,"1b53e5f36896d34022e0930b7d0c0e8d81c0a6346f248324250f961170e15ba8",14365,budOther),mk("technical_prose",twentyFirstProse,"5c152abd1237f9ab3b1b56f9c6738e1f7f6f0d0ac980efa1daa291f35c3bc6cc",1454,budOther)}},
		{"twenty-second",[]wlmLmR27Source{mk("code",twentySecondCode,"685394ef1dee84b20445ce8241177e95438ffaaa65c72df9ba7791db66abcc91",41453,budCode),mk("structured",twentySecondStructured,"86af1e7a250f29b2a82d7350b77923c62a18cc6decca294f4e949c50feb512f8",14365,budOther),mk("technical_prose",twentySecondProse,"a40eb8b08f758ca9667789e86b1d6d77d47effe565cd3c93496e0c721d0b3ee0",1454,budOther)}},
		{"twenty-third",[]wlmLmR27Source{mk("code",twentyThirdCode,"c339ef0f8ac223d4691cefe9bcb5e64191026c454642f4aec921d1b4749a37d1",41453,budCode),mk("structured",twentyThirdStructured,"5a49822316595b7f8288ed775f5ed8d5c7e6e29fda13f2eed193464db722578a",14365,budOther),mk("technical_prose",twentyThirdProse,"bed50a24e24c98133d08517bc7f9975f6434d4b77ccf486b580087370b0a1654",1454,budOther)}},
	}
	armDefs:=[]struct{id string;alpha [6]float64}{
		{"full-selected",[6]float64{0.75,0.75,0.125,0.125,1,1}},
		{"reset-distribution-character",[6]float64{0.5,0.5,0.5,0.5,1,1}},
		{"reset-character-layout-repetition",[6]float64{0.75,0.75,0.5,0.5,0.5,0.5}},
		{"reset-distribution-layout-repetition",[6]float64{0.5,0.5,0.125,0.125,0.5,0.5}},
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
	m["best_pairwise_reset_accuracy_delta"]=bestDelta
	m["best_vs_second_margin"]=bestDelta-secondDelta
	m["baseline_minus_full_accuracy"]=baseline.PairwiseWinnerAccuracy-full.PairwiseWinnerAccuracy
	for _,a:=range arms { if a.ID==bestReset { m["best_reset_positive_manifest_count"]=float64(a.PositiveManifestCountVsFull) } }
	m["total_eval_pair_count"]=float64((len(arms)+1)*108)
	for _,a:=range arms { if a.NonTiePairCount<72 || len(a.ManifestAccuracy)!=3 { m["invalid_row_count"]++ } }
	if baseline.NonTiePairCount<72 || len(baseline.ManifestAccuracy)!=3 { m["invalid_row_count"]++ }
	for k,v:=range m { if strings.HasSuffix(k,"_source_identity_mismatch_count")&&k!="source_identity_mismatch_count" { m["source_identity_mismatch_count"]+=v } }
	if len(arms)!=4||m["state_dimension"]!=6||m["readout_parameter_count"]!=7||m["total_adaptation_budget"]!=1744 { m["invalid_row_count"]++ }
	if m["fixed_budget_mismatch_count"]!=0||m["fit_failure_count"]!=0||m["raw_state_source_mismatch_count"]!=0||m["source_identity_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR49Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-INTERACTION-ATTRIBUTION-R49",BestPairwiseReset:bestReset,Arms:arms,Baseline:baseline,Metrics:m}
}
