package unitary

import (
	"math"
	"strings"
)

type wlmLmR52Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Arms []wlmLmR50ArmSummary `json:"arms"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmR52EvalArm(id string, alpha [6]float64, groups [6]int, groupCount int, hist []wlmLmR47Prepared, unseen []wlmLmR35ManifestInput, m map[string]float64) wlmLmR50ArmSummary {
	examples:=[]wlmLmR47Example{}
	for i,p:=range hist {
		st:=wlmLmR47States(p,alpha,m)
		examples=append(examples,wlmLmR47Examples(p,st,i)...)
	}
	beta,ok:=wlmLmR50Fit(examples,groups,groupCount)
	if !ok { m["fit_failure_count"]++ }
	actions:=wlmLmR47Actions(m)
	correct,total:=0,0
	regret:=0.0
	per:=map[string]float64{}
	for ui,in:=range unseen {
		pred:=wlmLmR47PreparedFrom("r52_pred_"+id+"_"+in.name,in.sources,4+ui,false,m)
		st:=wlmLmR47States(pred,alpha,m)
		scores:=map[[3]int]float64{}
		for _,a:=range actions { scores[a]=wlmLmR47PolicyPred(pred,st,beta,a) }
		actual:=wlmLmR47PreparedFrom("r52_actual_"+id+"_"+in.name,in.sources,4+ui,true,m)
		localCorrect,localN:=0,0
		for i:=0;i<len(actions);i++ {
			for j:=i+1;j<len(actions);j++ {
				a,b:=actions[i],actions[j]
				aa,ab:=wlmLmR27PolicyActual(actual.man,a),wlmLmR27PolicyActual(actual.man,b)
				if aa==ab { continue }
				localN++;total++
				want:=a;if ab>aa { want=b }
				got:=wlmLmR44Choice(a,b,scores[a],scores[b])
				if got==want { correct++;localCorrect++ } else { regret+=math.Abs(aa-ab) }
			}
		}
		if localN==0 { m["invalid_row_count"]++ } else { per[in.name]=float64(localCorrect)/float64(localN) }
	}
	acc,mr:=0.0,0.0
	if total==0 { m["invalid_row_count"]++ } else { acc=float64(correct)/float64(total);mr=regret/float64(total) }
	return wlmLmR50ArmSummary{ID:id,Groups:groups,EffectiveParameters:groupCount+1,PairwiseWinnerAccuracy:acc,MeanDecisionRegret:mr,NonTiePairCount:total,CorrectPairCount:correct,ManifestAccuracy:per}
}

func RunWlmLmExternalFutureDataLearnedStateFormerRepresentationReadoutCouplingDiagnosisR52(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse,
	twentyFourthCode,twentyFourthStructured,twentyFourthProse,
	twentyFifthCode,twentyFifthStructured,twentyFifthProse,
	twentySixthCode,twentySixthStructured,twentySixthProse []byte,
) interface{} {
	budCode:=[]int{0,436,582,727,872};budOther:=[]int{0,291,436,581,726}
	m:=map[string]float64{
		"historical_manifest_count":4,"unseen_manifest_count":3,"representation_arm_count":2,"readout_arm_count":2,"coupling_arm_count":4,
		"state_dimension":6,"max_readout_parameter_count":7,"total_adaptation_budget":1744,
		"candidate_pair_count_per_arm":108,"heldout_outcome_use_before_arm_freeze":0,"post_result_arm_or_threshold_choice_count":0,
		"capacity_growth_event_count":0,"external_model_call_count":0,"tokenizer_use_count":0,"fixed_budget_mismatch_count":0,
		"raw_state_source_mismatch_count":0,"source_identity_mismatch_count":0,"fit_failure_count":0,"invalid_row_count":0,
	}
	mk:=func(d string,b []byte,h string,n int,bud []int)wlmLmR27Source{return wlmLmR27Source{d:d,b:b,h:h,n:n,budgets:bud}}
	histInputs:=[]wlmLmR35ManifestInput{
		{"transfer",[]wlmLmR27Source{mk("code",transferCode,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45",41453,budCode),mk("structured",transferStructured,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932",14365,budOther),mk("technical_prose",transferProse,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b",1454,budOther)}},
		{"third",[]wlmLmR27Source{mk("code",thirdCode,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff",41453,budCode),mk("structured",thirdStructured,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56",14365,budOther),mk("technical_prose",thirdProse,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017",1454,budOther)}},
		{"fourth",[]wlmLmR27Source{mk("code",fourthCode,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1",41453,budCode),mk("structured",fourthStructured,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933",14365,budOther),mk("technical_prose",fourthProse,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222",1454,budOther)}},
		{"fifth",[]wlmLmR27Source{mk("code",fifthCode,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a",41453,budCode),mk("structured",fifthStructured,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281",14365,budOther),mk("technical_prose",fifthProse,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49",1454,budOther)}},
	}
	hist:=make([]wlmLmR47Prepared,0,4);for i,in:=range histInputs{hist=append(hist,wlmLmR47PreparedFrom("r52_hist_"+in.name,in.sources,i,true,m))}
	unseen:=[]wlmLmR35ManifestInput{
		{"twenty-fourth",[]wlmLmR27Source{mk("code",twentyFourthCode,"807bdfd41ab37f46c268ad4cc19c5d8d9fb80834dfcb5fa687ba8b1034f57d29",41453,budCode),mk("structured",twentyFourthStructured,"0f929cc80665958eac496aaa81913ecba99b08d1d2c90841f11c73faf8ec4d41",14365,budOther),mk("technical_prose",twentyFourthProse,"855c28784bbccb8ffc659055a6ed590e2a4ad71498bbce4fe902018b5f53deb2",1454,budOther)}},
		{"twenty-fifth",[]wlmLmR27Source{mk("code",twentyFifthCode,"ad758be3460bc6c0bf132e513756c7adf9cf469a00002a684cfde10089b81a79",41453,budCode),mk("structured",twentyFifthStructured,"d3d256ceffe672f2ad3a070c1f55f4ce0629ecec6b9d81e3e6d67be94ad89c89",14365,budOther),mk("technical_prose",twentyFifthProse,"7989e9b02507217afa6db6a06277c7d204abfd427f4f710292e626c2eefea24f",1454,budOther)}},
		{"twenty-sixth",[]wlmLmR27Source{mk("code",twentySixthCode,"717dbb2865eedcf185507a5fcf21fdb3f67d363dec58bfaf28de7dd5934c5fc1",41453,budCode),mk("structured",twentySixthStructured,"019edf98b26422a780a50535d9b3dc6b3ac184b71c397e0ee782b273ec074727",14365,budOther),mk("technical_prose",twentySixthProse,"f424f61fb7309a87fc4ccf092f959b39ac04795974aa41c1f1828bdf3751daf3",1454,budOther)}},
	}

	selected:=[6]float64{0.75,0.75,0.125,0.125,1,1}
	r46:=[6]float64{0.5,0.5,0.5,0.5,0.5,0.5}
	fullGroups:=[6]int{0,1,2,3,4,5}
	tiedGroups:=[6]int{0,0,1,1,1,1}
	arms:=[]wlmLmR50ArmSummary{
		wlmLmR52EvalArm("selected-full6",selected,fullGroups,6,hist,unseen,m),
		wlmLmR52EvalArm("selected-two-group2",selected,tiedGroups,2,hist,unseen,m),
		wlmLmR52EvalArm("r46-full6",r46,fullGroups,6,hist,unseen,m),
		wlmLmR52EvalArm("r46-two-group2",r46,tiedGroups,2,hist,unseen,m),
	}
	pos:=func(tied,full wlmLmR50ArmSummary)int{n:=0;for name,acc:=range tied.ManifestAccuracy{if acc>full.ManifestAccuracy[name]{n++}};return n}
	arms[1].AccuracyDeltaVsFull=arms[1].PairwiseWinnerAccuracy-arms[0].PairwiseWinnerAccuracy
	arms[1].PositiveManifestCountVsFull=pos(arms[1],arms[0])
	arms[3].AccuracyDeltaVsFull=arms[3].PairwiseWinnerAccuracy-arms[2].PairwiseWinnerAccuracy
	arms[3].PositiveManifestCountVsFull=pos(arms[3],arms[2])
	selectedDelta:=arms[1].AccuracyDeltaVsFull
	r46Delta:=arms[3].AccuracyDeltaVsFull
	m["selected_tied_accuracy_delta"]=selectedDelta
	m["selected_tied_positive_manifest_count"]=float64(arms[1].PositiveManifestCountVsFull)
	m["selected_tied_mean_regret"]=arms[1].MeanDecisionRegret
	m["selected_full_mean_regret"]=arms[0].MeanDecisionRegret
	m["r46_tied_accuracy_delta"]=r46Delta
	m["r46_tied_positive_manifest_count"]=float64(arms[3].PositiveManifestCountVsFull)
	m["r46_tied_mean_regret"]=arms[3].MeanDecisionRegret
	m["r46_full_mean_regret"]=arms[2].MeanDecisionRegret
	m["coupling_interaction_accuracy_delta"]=r46Delta-selectedDelta
	totalPairs:=0;for _,a:=range arms{totalPairs+=a.NonTiePairCount;if a.NonTiePairCount<72||len(a.ManifestAccuracy)!=3{m["invalid_row_count"]++}}
	m["total_eval_pair_count"]=float64(totalPairs)
	for k,v:=range m{if strings.HasSuffix(k,"_source_identity_mismatch_count")&&k!="source_identity_mismatch_count"{m["source_identity_mismatch_count"]+=v}}
	if len(arms)!=4||m["representation_arm_count"]!=2||m["readout_arm_count"]!=2||m["state_dimension"]!=6||m["max_readout_parameter_count"]!=7||m["total_adaptation_budget"]!=1744{m["invalid_row_count"]++}
	if m["fixed_budget_mismatch_count"]!=0||m["fit_failure_count"]!=0||m["raw_state_source_mismatch_count"]!=0||m["source_identity_mismatch_count"]!=0{m["invalid_row_count"]++}
	for _,v:=range m{if math.IsNaN(v)||math.IsInf(v,0){m["invalid_row_count"]++}}
	return wlmLmR52Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-REPRESENTATION-READOUT-COUPLING-DIAGNOSIS-R52",Arms:arms,Metrics:m}
}
