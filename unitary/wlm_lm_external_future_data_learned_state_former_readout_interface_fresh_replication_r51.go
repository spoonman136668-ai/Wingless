package unitary

import "strings"

type wlmLmR51Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Arms []wlmLmR50ArmSummary `json:"arms"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalFutureDataLearnedStateFormerReadoutInterfaceFreshReplicationR51(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse,
	twentyFourthCode,twentyFourthStructured,twentyFourthProse,
	twentyFifthCode,twentyFifthStructured,twentyFifthProse,
	twentySixthCode,twentySixthStructured,twentySixthProse []byte,
) interface{} {
	budCode:=[]int{0,436,582,727,872};budOther:=[]int{0,291,436,581,726}
	m:=map[string]float64{"historical_manifest_count":4,"unseen_manifest_count":3,"readout_arm_count":2,"state_dimension":6,"max_readout_parameter_count":7,"selected_readout_parameter_count":3,"total_adaptation_budget":1744,"candidate_pair_count_per_arm":108,"heldout_outcome_use_before_arm_freeze":0,"post_result_arm_or_threshold_choice_count":0,"capacity_growth_event_count":0,"external_model_call_count":0,"tokenizer_use_count":0,"fixed_budget_mismatch_count":0,"raw_state_source_mismatch_count":0,"source_identity_mismatch_count":0,"fit_failure_count":0,"invalid_row_count":0}
	mk:=func(d string,b []byte,h string,n int,bud []int)wlmLmR27Source{return wlmLmR27Source{d:d,b:b,h:h,n:n,budgets:bud}}
	histInputs:=[]wlmLmR35ManifestInput{
		{"transfer",[]wlmLmR27Source{mk("code",transferCode,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45",41453,budCode),mk("structured",transferStructured,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932",14365,budOther),mk("technical_prose",transferProse,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b",1454,budOther)}},
		{"third",[]wlmLmR27Source{mk("code",thirdCode,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff",41453,budCode),mk("structured",thirdStructured,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56",14365,budOther),mk("technical_prose",thirdProse,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017",1454,budOther)}},
		{"fourth",[]wlmLmR27Source{mk("code",fourthCode,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1",41453,budCode),mk("structured",fourthStructured,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933",14365,budOther),mk("technical_prose",fourthProse,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222",1454,budOther)}},
		{"fifth",[]wlmLmR27Source{mk("code",fifthCode,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a",41453,budCode),mk("structured",fifthStructured,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281",14365,budOther),mk("technical_prose",fifthProse,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49",1454,budOther)}},
	}
	hist:=make([]wlmLmR47Prepared,0,4);for i,in:=range histInputs{hist=append(hist,wlmLmR47PreparedFrom("r51_hist_"+in.name,in.sources,i,true,m))}
	unseen:=[]wlmLmR35ManifestInput{
		{"twenty-fourth",[]wlmLmR27Source{mk("code",twentyFourthCode,"807bdfd41ab37f46c268ad4cc19c5d8d9fb80834dfcb5fa687ba8b1034f57d29",41453,budCode),mk("structured",twentyFourthStructured,"0f929cc80665958eac496aaa81913ecba99b08d1d2c90841f11c73faf8ec4d41",14365,budOther),mk("technical_prose",twentyFourthProse,"855c28784bbccb8ffc659055a6ed590e2a4ad71498bbce4fe902018b5f53deb2",1454,budOther)}},
		{"twenty-fifth",[]wlmLmR27Source{mk("code",twentyFifthCode,"ad758be3460bc6c0bf132e513756c7adf9cf469a00002a684cfde10089b81a79",41453,budCode),mk("structured",twentyFifthStructured,"d3d256ceffe672f2ad3a070c1f55f4ce0629ecec6b9d81e3e6d67be94ad89c89",14365,budOther),mk("technical_prose",twentyFifthProse,"7989e9b02507217afa6db6a06277c7d204abfd427f4f710292e626c2eefea24f",1454,budOther)}},
		{"twenty-sixth",[]wlmLmR27Source{mk("code",twentySixthCode,"717dbb2865eedcf185507a5fcf21fdb3f67d363dec58bfaf28de7dd5934c5fc1",41453,budCode),mk("structured",twentySixthStructured,"019edf98b26422a780a50535d9b3dc6b3ac184b71c397e0ee782b273ec074727",14365,budOther),mk("technical_prose",twentySixthProse,"f424f61fb7309a87fc4ccf092f959b39ac04795974aa41c1f1828bdf3751daf3",1454,budOther)}},
	}
	full:=wlmLmR50EvalArm("full-6",[6]int{0,1,2,3,4,5},6,hist,unseen,m)
	selected:=wlmLmR50EvalArm("two-group-2",[6]int{0,0,1,1,1,1},2,hist,unseen,m)
	selected.AccuracyDeltaVsFull=selected.PairwiseWinnerAccuracy-full.PairwiseWinnerAccuracy
	pos:=0;for name,acc:=range selected.ManifestAccuracy{if acc>full.ManifestAccuracy[name]{pos++}}
	selected.PositiveManifestCountVsFull=pos;m["selected_accuracy_delta"]=selected.AccuracyDeltaVsFull;m["selected_positive_manifest_count"]=float64(pos);m["selected_mean_regret"]=selected.MeanDecisionRegret;m["full_mean_regret"]=full.MeanDecisionRegret;m["total_eval_pair_count"]=float64(full.NonTiePairCount+selected.NonTiePairCount)
	for _,a:=range []wlmLmR50ArmSummary{full,selected}{if a.NonTiePairCount<72||len(a.ManifestAccuracy)!=3{m["invalid_row_count"]++}}
	for k,v:=range m{if strings.HasSuffix(k,"_source_identity_mismatch_count")&&k!="source_identity_mismatch_count"{m["source_identity_mismatch_count"]+=v}}
	if m["fixed_budget_mismatch_count"]!=0||m["fit_failure_count"]!=0||m["raw_state_source_mismatch_count"]!=0||m["source_identity_mismatch_count"]!=0{m["invalid_row_count"]++}
	return wlmLmR51Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-READOUT-INTERFACE-FRESH-REPLICATION-R51",Arms:[]wlmLmR50ArmSummary{full,selected},Metrics:m}
}
