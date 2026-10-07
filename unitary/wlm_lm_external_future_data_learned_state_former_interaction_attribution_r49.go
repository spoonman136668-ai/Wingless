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
		{"twenty-first",[]wlmLmR27Source{mk("code",twentyFirstCode,"b82bc20dda9a4978ede68aa2db7f232a2b7cbaf7d3a152d4da9d7b5662be151a",41453,budCode),mk("structured",twentyFirstStructured,"21c287422a0f83bacdef4551ffa66435871f28ea62808736dd3b81db3d6b0b22",14365,budOther),mk("technical_prose",twentyFirstProse,"4357c8a5a220903e647e3e51d9bea46c71ed7bdb80db2986dbf4c6466b36832d",1454,budOther)}},
		{"twenty-second",[]wlmLmR27Source{mk("code",twentySecondCode,"02f3f7a21c41dc7c86e2ef0f783bb17ec4ff812f08dc0d1c5c7e26213ed7f499",41453,budCode),mk("structured",twentySecondStructured,"cd838e87e2ea93be1355b7ebb8cb05f427170918c8b5fa3b9e2083d508ee44f2",14365,budOther),mk("technical_prose",twentySecondProse,"c29df0165ddad99a1b4f7d798d6f7fd14a1acf1103eef89c1b21247ed596ecd6",1454,budOther)}},
		{"twenty-third",[]wlmLmR27Source{mk("code",twentyThirdCode,"3fcc427cc5e83aa450827ab1aea945663217b2ba867d907b49e81a44a7dd381f",41453,budCode),mk("structured",twentyThirdStructured,"1ae0bd65400663f591048f73212e20cebc172620df4519d0c9f3bacf2a240525",14365,budOther),mk("technical_prose",twentyThirdProse,"694a480a0dfb9e30f86b3f0d61f6343efb8795cbe56bcd9f406f0960695ce8a4",1454,budOther)}},
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
