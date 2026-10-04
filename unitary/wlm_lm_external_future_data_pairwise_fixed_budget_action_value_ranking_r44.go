package unitary

import (
	"math"
	"strings"
)

type wlmLmR44PairCase struct {
	Manifest string `json:"manifest"`
	ActionA [3]int `json:"action_a"`
	ActionB [3]int `json:"action_b"`
	R43Choice [3]int `json:"r43_choice"`
	ActionOnlyChoice [3]int `json:"action_only_choice"`
	StateShuffledChoice [3]int `json:"state_shuffled_choice"`
	ActualChoice [3]int `json:"actual_choice"`
	ActualDelta float64 `json:"actual_delta_a_minus_b"`
	Tie bool `json:"tie"`
}

type wlmLmR44Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR44PairCase `json:"cases"`
}

func wlmLmR44Choice(a,b [3]int,sa,sb float64) [3]int {
	if sa>sb { return a }
	if sb>sa { return b }
	if wlmLmR27LexLess(a,b) { return a }
	return b
}

func wlmLmR44FitState(mans []wlmLmR27Manifest, shuffled bool, m map[string]float64) ([]float64,[3]float64,bool) {
	ex:=[]wlmLmR34Example{}
	if !shuffled {
		ex=wlmLmR35Examples(mans,false)
	} else {
		for i:=range mans {
			xs:=wlmLmR35Examples([]wlmLmR27Manifest{mans[i]},false)
			ys:=wlmLmR35Examples([]wlmLmR27Manifest{mans[(i+1)%len(mans)]},false)
			if len(xs)!=len(ys) { m["invalid_row_count"]++; return nil,[3]float64{},false }
			for k:=range xs { ex=append(ex,wlmLmR34Example{x:xs[k].x,y:ys[k].y}) }
		}
	}
	beta,ok:=wlmLmR34Fit(ex)
	if !ok { m["fit_failure_count"]++; return nil,[3]float64{},false }
	var sums [3]float64
	var ns [3]int
	for i,man:=range mans {
		actualMan:=man
		if shuffled { actualMan=mans[(i+1)%len(mans)] }
		for _,alloc:=range wlmLmR27Candidates {
			if alloc==wlmLmR27Equal { continue }
			d:=wlmLmR43StateIndex(alloc,m)
			pred:=wlmLmR35PolicyPred(man,beta,false,alloc)-wlmLmR35PolicyPred(man,beta,false,wlmLmR27Equal)
			actual:=wlmLmR27PolicyActual(actualMan,alloc)-wlmLmR27PolicyActual(actualMan,wlmLmR27Equal)
			sums[d]+=actual-pred;ns[d]++
		}
	}
	var offsets [3]float64
	for d:=0;d<3;d++ {
		if ns[d]==0 { m["state_stratum_insufficient_count"]++; return nil,[3]float64{},false }
		offsets[d]=sums[d]/float64(ns[d])
		if math.IsNaN(offsets[d])||math.IsInf(offsets[d],0) { m["invalid_row_count"]++;return nil,[3]float64{},false }
	}
	return beta,offsets,true
}

func wlmLmR44StateScore(man wlmLmR27Manifest,beta []float64,offsets [3]float64,alloc [3]int,m map[string]float64) float64 {
	d:=wlmLmR43StateIndex(alloc,m)
	return wlmLmR35PolicyPred(man,beta,false,alloc)-wlmLmR35PolicyPred(man,beta,false,wlmLmR27Equal)+offsets[d]
}

func RunWlmLmExternalFutureDataPairwiseFixedBudgetActionValueRankingR44(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse,
	sixthCode,sixthStructured,sixthProse,
	seventhCode,seventhStructured,seventhProse,
	eighthCode,eighthStructured,eighthProse []byte,
) interface{} {
	budCode:=[]int{0,436,582,727,872}
	budOther:=[]int{0,291,436,581,726}
	m:=map[string]float64{
		"historical_manifest_count":4,"unseen_manifest_count":3,"candidate_action_count":9,
		"candidate_pair_count":0,"non_tie_pair_count":0,"exact_tie_count":0,
		"r43_correct_pair_count":0,"action_only_correct_pair_count":0,"state_shuffled_correct_pair_count":0,
		"r43_pairwise_winner_accuracy":0,"action_only_pairwise_winner_accuracy":0,"state_shuffled_pairwise_winner_accuracy":0,
		"r43_minus_action_only_accuracy":0,"r43_mean_decision_regret":0,"action_only_mean_decision_regret":0,
		"state_shuffled_mean_decision_regret":0,"fit_failure_count":0,"state_stratum_insufficient_count":0,
		"state_definition_mismatch_count":0,"fixed_budget_mismatch_count":0,"heldout_outcome_use_before_pair_freeze":0,
		"post_result_pair_or_threshold_choice_count":0,"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	mk:=func(d string,b []byte,h string,n int,bud []int) wlmLmR27Source {
		return wlmLmR27Source{d:d,b:b,h:h,n:n,budgets:bud}
	}
	mkSelf:=func(d string,b []byte,n int,bud []int) wlmLmR27Source {
		if len(b)!=n { m["source_identity_mismatch_count"]++ }
		return wlmLmR27Source{d:d,b:b,h:wlmLmExternalRawRepPredR1SHA256(b),n:n,budgets:bud}
	}
	historicalInputs:=[]wlmLmR35ManifestInput{
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
	unseenInputs:=[]wlmLmR35ManifestInput{
		{"sixth",[]wlmLmR27Source{
			mkSelf("code",sixthCode,41453,budCode),mkSelf("structured",sixthStructured,14365,budOther),mkSelf("technical_prose",sixthProse,1454,budOther),
		}},
		{"seventh",[]wlmLmR27Source{
			mkSelf("code",seventhCode,41453,budCode),mkSelf("structured",seventhStructured,14365,budOther),mkSelf("technical_prose",seventhProse,1454,budOther),
		}},
		{"eighth",[]wlmLmR27Source{
			mkSelf("code",eighthCode,41453,budCode),mkSelf("structured",eighthStructured,14365,budOther),mkSelf("technical_prose",eighthProse,1454,budOther),
		}},
	}
	actions:=make([][3]int,0,9)
	for _,a:=range wlmLmR27Candidates {
		if a==wlmLmR27Equal { continue }
		if a[0]+a[1]+a[2]!=1744 { m["fixed_budget_mismatch_count"]++ }
		actions=append(actions,a)
	}
	if len(actions)!=9 { m["invalid_row_count"]++ }

	historical:=make([]wlmLmR27Manifest,0,4)
	for i,in:=range historicalInputs {
		man,ex:=wlmLmR27Build("r44_hist_"+in.name,in.sources,i,true,m)
		if len(ex)!=12 { m["invalid_row_count"]++ }
		historical=append(historical,man)
	}
	beta,offsets,ok:=wlmLmR44FitState(historical,false,m)
	shufBeta,shufOffsets,ok2:=wlmLmR44FitState(historical,true,m)
	if !ok||!ok2 { m["invalid_row_count"]++ }

	actionMean:=map[[3]int]float64{}
	for _,a:=range actions {
		total:=0.0
		for _,man:=range historical {
			total+=wlmLmR27PolicyActual(man,a)-wlmLmR27PolicyActual(man,wlmLmR27Equal)
		}
		actionMean[a]=total/float64(len(historical))
	}

	cases:=make([]wlmLmR44PairCase,0,108)
	r43Regret,actionRegret,shufRegret:=0.0,0.0,0.0
	for ui,in:=range unseenInputs {
		features,_:=wlmLmR27Build("r44_frozen_"+in.name,in.sources,4+ui,false,m)
		type score struct{r43,action,shuf float64}
		scores:=map[[3]int]score{}
		for _,a:=range actions {
			scores[a]=score{
				r43:wlmLmR44StateScore(features,beta,offsets,a,m),
				action:actionMean[a],
				shuf:wlmLmR44StateScore(features,shufBeta,shufOffsets,a,m),
			}
		}
		actualMan,_:=wlmLmR27Build("r44_actual_"+in.name,in.sources,4+ui,true,m)
		actual:=map[[3]int]float64{}
		for _,a:=range actions { actual[a]=wlmLmR27PolicyActual(actualMan,a)-wlmLmR27PolicyActual(actualMan,wlmLmR27Equal) }
		localNonTie,localR43,localAction:=0,0,0
		for i:=0;i<len(actions);i++ {
			for j:=i+1;j<len(actions);j++ {
				a,b:=actions[i],actions[j];m["candidate_pair_count"]++
				aa,ab:=actual[a],actual[b]
				r43c:=wlmLmR44Choice(a,b,scores[a].r43,scores[b].r43)
				actc:=wlmLmR44Choice(a,b,scores[a].action,scores[b].action)
				shufc:=wlmLmR44Choice(a,b,scores[a].shuf,scores[b].shuf)
				row:=wlmLmR44PairCase{Manifest:in.name,ActionA:a,ActionB:b,R43Choice:r43c,ActionOnlyChoice:actc,StateShuffledChoice:shufc,ActualDelta:aa-ab}
				if aa==ab {
					row.Tie=true;m["exact_tie_count"]++
				} else {
					m["non_tie_pair_count"]++;localNonTie++
					winner:=a;if ab>aa { winner=b };row.ActualChoice=winner
					gap:=math.Abs(aa-ab)
					if r43c==winner { m["r43_correct_pair_count"]++;localR43++ } else { r43Regret+=gap }
					if actc==winner { m["action_only_correct_pair_count"]++;localAction++ } else { actionRegret+=gap }
					if shufc==winner { m["state_shuffled_correct_pair_count"]++ } else { shufRegret+=gap }
				}
				cases=append(cases,row)
			}
		}
		if localNonTie>0 {
			racc:=float64(localR43)/float64(localNonTie)
			aacc:=float64(localAction)/float64(localNonTie)
			m[in.name+"_r43_pairwise_winner_accuracy"]=racc
			m[in.name+"_action_only_pairwise_winner_accuracy"]=aacc
			m[in.name+"_r43_minus_action_only_accuracy"]=racc-aacc
		} else { m["invalid_row_count"]++ }
	}
	n:=m["non_tie_pair_count"]
	if n>0 {
		m["r43_pairwise_winner_accuracy"]=m["r43_correct_pair_count"]/n
		m["action_only_pairwise_winner_accuracy"]=m["action_only_correct_pair_count"]/n
		m["state_shuffled_pairwise_winner_accuracy"]=m["state_shuffled_correct_pair_count"]/n
		m["r43_minus_action_only_accuracy"]=m["r43_pairwise_winner_accuracy"]-m["action_only_pairwise_winner_accuracy"]
		m["r43_mean_decision_regret"]=r43Regret/n
		m["action_only_mean_decision_regret"]=actionRegret/n
		m["state_shuffled_mean_decision_regret"]=shufRegret/n
	}
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" { m["source_identity_mismatch_count"]+=v }
	}
	if int(m["candidate_pair_count"])!=108 || len(cases)!=108 { m["invalid_row_count"]++ }
	if m["fit_failure_count"]!=0 || m["state_stratum_insufficient_count"]!=0 || m["state_definition_mismatch_count"]!=0 || m["fixed_budget_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR44Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-PAIRWISE-FIXED-BUDGET-ACTION-VALUE-RANKING-R44",Metrics:m,Cases:cases}
}
