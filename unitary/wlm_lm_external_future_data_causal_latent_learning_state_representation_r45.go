package unitary

import (
	"math"
	"strings"
)

type wlmLmR45PairCase struct {
	Manifest string `json:"manifest"`
	ActionA [3]int `json:"action_a"`
	ActionB [3]int `json:"action_b"`
	TrajectoryChoice [3]int `json:"trajectory_choice"`
	Base6Choice [3]int `json:"base6_choice"`
	ActionOnlyChoice [3]int `json:"action_only_choice"`
	ActualChoice [3]int `json:"actual_choice"`
	ActualDelta float64 `json:"actual_delta_a_minus_b"`
	Tie bool `json:"tie"`
}

type wlmLmR45Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR45PairCase `json:"cases"`
}

func wlmLmR45Channels(f [6]float64) [3]float64 {
	return [3]float64{
		f[0] + f[5],
		f[1],
		f[2] + f[3] + f[4],
	}
}

func wlmLmR45TrajectoryVec(arm wlmLmR27Arm, hi int, m map[string]float64) []float64 {
	k:=-1
	for i,b:=range arm.budgets {
		if b==hi { k=i; break }
	}
	if k<1 {
		m["trajectory_definition_mismatch_count"]++
		return make([]float64,6)
	}
	var weighted [3]float64
	var first,last [3]float64
	ws:=0.0
	for j:=1;j<=k;j++ {
		f,ok:=arm.features[arm.budgets[j]]
		if !ok {
			m["trajectory_definition_mismatch_count"]++
			return make([]float64,6)
		}
		ch:=wlmLmR45Channels(f)
		if j==1 { first=ch }
		last=ch
		w:=float64(j)
		ws+=w
		for q:=0;q<3;q++ { weighted[q]+=w*ch[q] }
	}
	if ws<=0 {
		m["trajectory_definition_mismatch_count"]++
		return make([]float64,6)
	}
	x:=[]float64{
		weighted[0]/ws,last[0]-first[0],
		weighted[1]/ws,last[1]-first[1],
		weighted[2]/ws,last[2]-first[2],
	}
	for _,v:=range x {
		if math.IsNaN(v)||math.IsInf(v,0) { m["trajectory_definition_mismatch_count"]++ }
	}
	return x
}

func wlmLmR45Examples(mans []wlmLmR27Manifest,m map[string]float64) []wlmLmR34Example {
	out:=[]wlmLmR34Example{}
	for _,man:=range mans {
		for _,arm:=range man.arms {
			for i:=1;i<len(arm.budgets);i++ {
				lo,hi:=arm.budgets[i-1],arm.budgets[i]
				x:=wlmLmR45TrajectoryVec(arm,hi,m)
				if len(x)!=6 { m["trajectory_definition_mismatch_count"]++ }
				out=append(out,wlmLmR34Example{x:x,y:float64(arm.hits[hi]-arm.hits[lo])})
			}
		}
	}
	return out
}

func wlmLmR45PolicyPred(man wlmLmR27Manifest,beta []float64,alloc [3]int,m map[string]float64) float64 {
	total:=0.0
	for ai,arm:=range man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			if hi<=alloc[ai] {
				total+=wlmLmR34Predict(beta,wlmLmR45TrajectoryVec(arm,hi,m))
			}
		}
	}
	return total
}

func RunWlmLmExternalFutureDataCausalLatentLearningStateRepresentationR45(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse,
	ninthCode,ninthStructured,ninthProse,
	tenthCode,tenthStructured,tenthProse,
	eleventhCode,eleventhStructured,eleventhProse []byte,
) interface{} {
	budCode:=[]int{0,436,582,727,872}
	budOther:=[]int{0,291,436,581,726}
	m:=map[string]float64{
		"historical_manifest_count":4,"unseen_manifest_count":3,"candidate_action_count":9,
		"candidate_pair_count":0,"non_tie_pair_count":0,"exact_tie_count":0,
		"trajectory_correct_pair_count":0,"base6_correct_pair_count":0,"action_only_correct_pair_count":0,
		"trajectory_pairwise_winner_accuracy":0,"base6_pairwise_winner_accuracy":0,"action_only_pairwise_winner_accuracy":0,
		"trajectory_minus_base6_accuracy":0,"trajectory_minus_action_only_accuracy":0,
		"trajectory_mean_decision_regret":0,"base6_mean_decision_regret":0,"action_only_mean_decision_regret":0,
		"positive_manifest_count":0,
		"representation_dimension":6,"trajectory_readout_parameter_count":0,"base6_readout_parameter_count":0,
		"fit_failure_count":0,"trajectory_definition_mismatch_count":0,"readout_parameter_count_mismatch_count":0,
		"fixed_budget_mismatch_count":0,"heldout_outcome_use_before_pair_freeze":0,
		"post_result_representation_or_threshold_choice_count":0,"source_identity_mismatch_count":0,
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
		{"ninth",[]wlmLmR27Source{mkSelf("code",ninthCode,41453,budCode),mkSelf("structured",ninthStructured,14365,budOther),mkSelf("technical_prose",ninthProse,1454,budOther)}},
		{"tenth",[]wlmLmR27Source{mkSelf("code",tenthCode,41453,budCode),mkSelf("structured",tenthStructured,14365,budOther),mkSelf("technical_prose",tenthProse,1454,budOther)}},
		{"eleventh",[]wlmLmR27Source{mkSelf("code",eleventhCode,41453,budCode),mkSelf("structured",eleventhStructured,14365,budOther),mkSelf("technical_prose",eleventhProse,1454,budOther)}},
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
		man,ex:=wlmLmR27Build("r45_hist_"+in.name,in.sources,i,true,m)
		if len(ex)!=12 { m["invalid_row_count"]++ }
		historical=append(historical,man)
	}
	baseBeta,ok0:=wlmLmR34Fit(wlmLmR35Examples(historical,false))
	trajBeta,ok1:=wlmLmR34Fit(wlmLmR45Examples(historical,m))
	if !ok0||!ok1 { m["fit_failure_count"]++;m["invalid_row_count"]++ }
	m["base6_readout_parameter_count"]=float64(len(baseBeta))
	m["trajectory_readout_parameter_count"]=float64(len(trajBeta))
	if len(baseBeta)!=7||len(trajBeta)!=7 {
		m["readout_parameter_count_mismatch_count"]++
		m["invalid_row_count"]++
	}

	actionMean:=map[[3]int]float64{}
	for _,a:=range actions {
		total:=0.0
		for _,man:=range historical {
			total+=wlmLmR27PolicyActual(man,a)-wlmLmR27PolicyActual(man,wlmLmR27Equal)
		}
		actionMean[a]=total/float64(len(historical))
	}

	cases:=make([]wlmLmR45PairCase,0,108)
	trajRegret,baseRegret,actionRegret:=0.0,0.0,0.0
	for ui,in:=range unseenInputs {
		features,_:=wlmLmR27Build("r45_frozen_"+in.name,in.sources,4+ui,false,m)
		type score struct{traj,base,action float64}
		scores:=map[[3]int]score{}
		for _,a:=range actions {
			scores[a]=score{
				traj:wlmLmR45PolicyPred(features,trajBeta,a,m)-wlmLmR45PolicyPred(features,trajBeta,wlmLmR27Equal,m),
				base:wlmLmR35PolicyPred(features,baseBeta,false,a)-wlmLmR35PolicyPred(features,baseBeta,false,wlmLmR27Equal),
				action:actionMean[a],
			}
		}
		actualMan,_:=wlmLmR27Build("r45_actual_"+in.name,in.sources,4+ui,true,m)
		actual:=map[[3]int]float64{}
		for _,a:=range actions { actual[a]=wlmLmR27PolicyActual(actualMan,a)-wlmLmR27PolicyActual(actualMan,wlmLmR27Equal) }
		localNonTie,localTraj,localBase:=0,0,0
		for i:=0;i<len(actions);i++ {
			for j:=i+1;j<len(actions);j++ {
				a,b:=actions[i],actions[j]
				m["candidate_pair_count"]++
				aa,ab:=actual[a],actual[b]
				tc:=wlmLmR44Choice(a,b,scores[a].traj,scores[b].traj)
				bc:=wlmLmR44Choice(a,b,scores[a].base,scores[b].base)
				ac:=wlmLmR44Choice(a,b,scores[a].action,scores[b].action)
				row:=wlmLmR45PairCase{Manifest:in.name,ActionA:a,ActionB:b,TrajectoryChoice:tc,Base6Choice:bc,ActionOnlyChoice:ac,ActualDelta:aa-ab}
				if aa==ab {
					row.Tie=true;m["exact_tie_count"]++
				} else {
					m["non_tie_pair_count"]++;localNonTie++
					winner:=a
					if ab>aa { winner=b }
					row.ActualChoice=winner
					gap:=math.Abs(aa-ab)
					if tc==winner { m["trajectory_correct_pair_count"]++;localTraj++ } else { trajRegret+=gap }
					if bc==winner { m["base6_correct_pair_count"]++;localBase++ } else { baseRegret+=gap }
					if ac==winner { m["action_only_correct_pair_count"]++ } else { actionRegret+=gap }
				}
				cases=append(cases,row)
			}
		}
		if localNonTie<=0 {
			m["invalid_row_count"]++
		} else {
			ta:=float64(localTraj)/float64(localNonTie)
			ba:=float64(localBase)/float64(localNonTie)
			m[in.name+"_trajectory_pairwise_winner_accuracy"]=ta
			m[in.name+"_base6_pairwise_winner_accuracy"]=ba
			m[in.name+"_trajectory_minus_base6_accuracy"]=ta-ba
			if ta>ba { m["positive_manifest_count"]++ }
		}
	}
	n:=m["non_tie_pair_count"]
	if n<=0 {
		m["invalid_row_count"]++
	} else {
		m["trajectory_pairwise_winner_accuracy"]=m["trajectory_correct_pair_count"]/n
		m["base6_pairwise_winner_accuracy"]=m["base6_correct_pair_count"]/n
		m["action_only_pairwise_winner_accuracy"]=m["action_only_correct_pair_count"]/n
		m["trajectory_minus_base6_accuracy"]=m["trajectory_pairwise_winner_accuracy"]-m["base6_pairwise_winner_accuracy"]
		m["trajectory_minus_action_only_accuracy"]=m["trajectory_pairwise_winner_accuracy"]-m["action_only_pairwise_winner_accuracy"]
		m["trajectory_mean_decision_regret"]=trajRegret/n
		m["base6_mean_decision_regret"]=baseRegret/n
		m["action_only_mean_decision_regret"]=actionRegret/n
	}
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" { m["source_identity_mismatch_count"]+=v }
	}
	if int(m["candidate_pair_count"])!=108||len(cases)!=108 { m["invalid_row_count"]++ }
	if m["fit_failure_count"]!=0||m["trajectory_definition_mismatch_count"]!=0||m["readout_parameter_count_mismatch_count"]!=0||m["fixed_budget_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR45Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-CAUSAL-LATENT-LEARNING-STATE-REPRESENTATION-R45",Metrics:m,Cases:cases}
}
