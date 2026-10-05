package unitary

import "math"

type wlmLmR44Case struct {
	Manifest string `json:"manifest"`
	ActionA [3]int `json:"action_a"`
	ActionB [3]int `json:"action_b"`
	ActualA float64 `json:"actual_a"`
	ActualB float64 `json:"actual_b"`
	R43ScoreA float64 `json:"r43_score_a"`
	R43ScoreB float64 `json:"r43_score_b"`
	ActionOnlyScoreA float64 `json:"action_only_score_a"`
	ActionOnlyScoreB float64 `json:"action_only_score_b"`
	ShuffledScoreA float64 `json:"shuffled_score_a"`
	ShuffledScoreB float64 `json:"shuffled_score_b"`
	ModeScoreA float64 `json:"mode_score_a"`
	ModeScoreB float64 `json:"mode_score_b"`
	Tie bool `json:"tie"`
}

type wlmLmR44Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Mode string `json:"mode"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR44Case `json:"cases"`
}

type wlmLmR44ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

func wlmLmR44StateIndex(alloc [3]int) int {
	if alloc[0]==582 { return 0 }
	if alloc[0]==727 { return 1 }
	if alloc[0]==872 { return 2 }
	return -1
}

func wlmLmR44DynamicSources(code,structured,prose []byte) []wlmLmR27Source {
	bc:=[]int{0,436,582,727,872}; bo:=[]int{0,291,436,581,726}
	mk:=func(d string,b []byte,n int,bud []int) wlmLmR27Source {
		return wlmLmR27Source{d:d,b:b,h:wlmLmExternalRawRepPredR1SHA256(b),n:n,budgets:bud}
	}
	return []wlmLmR27Source{
		mk("code",code,41453,bc),mk("structured",structured,14365,bo),mk("technical_prose",prose,1454,bo),
	}
}

func wlmLmR44HistoricalInputs(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse []byte,
) []wlmLmR44ManifestInput {
	bc:=[]int{0,436,582,727,872}; bo:=[]int{0,291,436,581,726}
	mk:=func(d string,b []byte,h string,n int,bud []int) wlmLmR27Source {
		return wlmLmR27Source{d:d,b:b,h:h,n:n,budgets:bud}
	}
	return []wlmLmR44ManifestInput{
		{"transfer",[]wlmLmR27Source{
			mk("code",transferCode,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45",41453,bc),
			mk("structured",transferStructured,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932",14365,bo),
			mk("technical_prose",transferProse,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b",1454,bo)}},
		{"third",[]wlmLmR27Source{
			mk("code",thirdCode,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff",41453,bc),
			mk("structured",thirdStructured,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56",14365,bo),
			mk("technical_prose",thirdProse,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017",1454,bo)}},
		{"fourth",[]wlmLmR27Source{
			mk("code",fourthCode,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1",41453,bc),
			mk("structured",fourthStructured,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933",14365,bo),
			mk("technical_prose",fourthProse,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222",1454,bo)}},
		{"fifth",[]wlmLmR27Source{
			mk("code",fifthCode,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a",41453,bc),
			mk("structured",fifthStructured,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281",14365,bo),
			mk("technical_prose",fifthProse,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49",1454,bo)}},
	}
}

func wlmLmR44Vec(arm wlmLmR27Arm,hi int,mode string) []float64 {
	f:=arm.features[hi]
	if mode=="ablation" {
		return []float64{f[0],f[1],0,0,0,f[5]}
	}
	if mode=="causal-latent" {
		prev:=[6]float64{}
		for i:=1;i<len(arm.budgets);i++ {
			if arm.budgets[i]==hi && i>1 { prev=arm.features[arm.budgets[i-1]];break }
		}
		churn:=0.0
		for i:=0;i<6;i++ { churn+=math.Abs(f[i]-prev[i]) }
		return []float64{f[0],f[1],f[2]-f[3],f[4],f[5],churn}
	}
	return []float64{f[0],f[1],f[2],f[3],f[4],f[5]}
}

func wlmLmR44Examples(mans []wlmLmR27Manifest,mode string) []wlmLmR34Example {
	out:=[]wlmLmR34Example{}
	for _,man:=range mans {
		for _,arm:=range man.arms {
			for i:=1;i<len(arm.budgets);i++ {
				lo,hi:=arm.budgets[i-1],arm.budgets[i]
				out=append(out,wlmLmR34Example{x:wlmLmR44Vec(arm,hi,mode),y:float64(arm.hits[hi]-arm.hits[lo])})
			}
		}
	}
	return out
}

func wlmLmR44PolicyPred(man wlmLmR27Manifest,beta []float64,mode string,alloc [3]int) float64 {
	v:=0.0
	for ai,arm:=range man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			if hi<=alloc[ai] { v+=wlmLmR34Predict(beta,wlmLmR44Vec(arm,hi,mode)) }
		}
	}
	return v
}

func wlmLmR44Offsets(mans []wlmLmR27Manifest,beta []float64,mode string,m map[string]float64) [3]float64 {
	var sum [3]float64; var n [3]int
	for _,man:=range mans {
		for _,a:=range wlmLmR27Candidates {
			if a==wlmLmR27Equal { continue }
			k:=wlmLmR44StateIndex(a)
			if k<0 { m["state_definition_mismatch_count"]++;continue }
			p:=wlmLmR44PolicyPred(man,beta,mode,a)-wlmLmR44PolicyPred(man,beta,mode,wlmLmR27Equal)
			y:=wlmLmR27PolicyActual(man,a)-wlmLmR27PolicyActual(man,wlmLmR27Equal)
			sum[k]+=y-p;n[k]++
		}
	}
	var out [3]float64
	for i:=0;i<3;i++ {
		if n[i]==0 { m["invalid_row_count"]++;continue }
		out[i]=sum[i]/float64(n[i])
	}
	return out
}

func wlmLmR44Score(man wlmLmR27Manifest,beta []float64,offs [3]float64,mode string,a [3]int) float64 {
	k:=wlmLmR44StateIndex(a)
	v:=wlmLmR44PolicyPred(man,beta,mode,a)
	if k>=0 { v+=offs[k] }
	return v
}

func wlmLmR44Win(scoreA,scoreB,actualA,actualB float64)(bool,float64) {
	if actualA==actualB { return false,0 }
	chooseA:=scoreA>=scoreB
	actualChooseA:=actualA>actualB
	regret:=0.0
	if chooseA!=actualChooseA { regret=math.Abs(actualA-actualB) }
	return chooseA==actualChooseA,regret
}

func RunWlmLmExternalFutureDataPairwiseActionValueR44(
	mode string,
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse,
	e1Code,e1Structured,e1Prose,
	e2Code,e2Structured,e2Prose,
	e3Code,e3Structured,e3Prose []byte,
) interface{} {
	if mode!="r44"&&mode!="replication"&&mode!="ablation"&&mode!="causal-latent" { mode="invalid" }
	m:=map[string]float64{
		"historical_manifest_count":4,"evaluation_manifest_count":3,"candidate_action_count":9,
		"pair_count":0,"non_tie_pair_count":0,"exact_tie_count":0,
		"r43_pairwise_correct":0,"action_only_pairwise_correct":0,"state_shuffled_pairwise_correct":0,"mode_pairwise_correct":0,
		"r43_decision_regret_sum":0,"action_only_decision_regret_sum":0,"state_shuffled_decision_regret_sum":0,"mode_decision_regret_sum":0,
		"fit_failure_count":0,"state_definition_mismatch_count":0,
		"heldout_outcome_use_before_pair_freeze":0,"post_result_pair_or_threshold_choice_count":0,
		"source_identity_mismatch_count":0,"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"invalid_row_count":0,
	}
	hInputs:=wlmLmR44HistoricalInputs(
		transferCode,transferStructured,transferProse,thirdCode,thirdStructured,thirdProse,
		fourthCode,fourthStructured,fourthProse,fifthCode,fifthStructured,fifthProse)
	hist:=make([]wlmLmR27Manifest,0,4)
	for i,in:=range hInputs {
		man,_:=wlmLmR27Build("r44_hist_"+in.name,in.sources,i,true,m);hist=append(hist,man)
	}
	fullBeta,ok:=wlmLmR34Fit(wlmLmR44Examples(hist,"r43"))
	if !ok { m["fit_failure_count"]++;m["invalid_row_count"]++ }
	fullOff:=wlmLmR44Offsets(hist,fullBeta,"r43",m)

	activeMode:=mode
	if activeMode=="r44"||activeMode=="replication" { activeMode="r43" }
	modeBeta,ok2:=wlmLmR34Fit(wlmLmR44Examples(hist,activeMode))
	if !ok2 { m["fit_failure_count"]++;m["invalid_row_count"]++ }
	modeOff:=wlmLmR44Offsets(hist,modeBeta,activeMode,m)

	var actionMean=map[[3]int]float64{}
	for _,a:=range wlmLmR27Candidates {
		if a==wlmLmR27Equal { continue }
		sum:=0.0
		for _,man:=range hist { sum+=wlmLmR27PolicyActual(man,a) }
		actionMean[a]=sum/float64(len(hist))
	}
	evalInputs:=[]wlmLmR44ManifestInput{
		{"eval1",wlmLmR44DynamicSources(e1Code,e1Structured,e1Prose)},
		{"eval2",wlmLmR44DynamicSources(e2Code,e2Structured,e2Prose)},
		{"eval3",wlmLmR44DynamicSources(e3Code,e3Structured,e3Prose)},
	}
	candidates:=make([][3]int,0,9)
	for _,a:=range wlmLmR27Candidates { if a!=wlmLmR27Equal { candidates=append(candidates,a) } }
	cases:=make([]wlmLmR44Case,0,108)
	for mi,in:=range evalInputs {
		features,_:=wlmLmR27Build("r44_features_"+in.name,in.sources,10+mi,false,m)
		type sc struct{a [3]int;r43,ao,shuf,mode float64}
		scores:=make([]sc,0,9)
		for _,a:=range candidates {
			k:=wlmLmR44StateIndex(a)
			shuf:=wlmLmR44PolicyPred(features,fullBeta,"r43",a)
			if k>=0 { shuf+=fullOff[(k+1)%3] }
			scores=append(scores,sc{a:a,r43:wlmLmR44Score(features,fullBeta,fullOff,"r43",a),ao:actionMean[a],shuf:shuf,mode:wlmLmR44Score(features,modeBeta,modeOff,activeMode,a)})
		}
		actual,_:=wlmLmR27Build("r44_actual_"+in.name,in.sources,10+mi,true,m)
		manN:=0.0;manR:=0.0;manAO:=0.0
		for i:=0;i<len(scores);i++ {
			for j:=i+1;j<len(scores);j++ {
				a,b:=scores[i],scores[j]
				ya:=wlmLmR27PolicyActual(actual,a.a);yb:=wlmLmR27PolicyActual(actual,b.a)
				row:=wlmLmR44Case{Manifest:in.name,ActionA:a.a,ActionB:b.a,ActualA:ya,ActualB:yb,
					R43ScoreA:a.r43,R43ScoreB:b.r43,ActionOnlyScoreA:a.ao,ActionOnlyScoreB:b.ao,
					ShuffledScoreA:a.shuf,ShuffledScoreB:b.shuf,ModeScoreA:a.mode,ModeScoreB:b.mode,Tie:ya==yb}
				m["pair_count"]++
				if ya==yb { m["exact_tie_count"]++;cases=append(cases,row);continue }
				m["non_tie_pair_count"]++;manN++
				if ok,r:=wlmLmR44Win(a.r43,b.r43,ya,yb);ok { m["r43_pairwise_correct"]++;manR++ } else { m["r43_decision_regret_sum"]+=r }
				if ok,r:=wlmLmR44Win(a.ao,b.ao,ya,yb);ok { m["action_only_pairwise_correct"]++;manAO++ } else { m["action_only_decision_regret_sum"]+=r }
				if ok,r:=wlmLmR44Win(a.shuf,b.shuf,ya,yb);ok { m["state_shuffled_pairwise_correct"]++ } else { m["state_shuffled_decision_regret_sum"]+=r }
				if ok,r:=wlmLmR44Win(a.mode,b.mode,ya,yb);ok { m["mode_pairwise_correct"]++ } else { m["mode_decision_regret_sum"]+=r }
				cases=append(cases,row)
			}
		}
		if manN>0 {
			m[in.name+"_r43_accuracy"]=manR/manN
			m[in.name+"_action_only_accuracy"]=manAO/manN
			m[in.name+"_r43_minus_action_only_accuracy"]=(manR-manAO)/manN
		}
	}
	n:=m["non_tie_pair_count"]
	if n>0 {
		m["r43_pairwise_winner_accuracy"]=m["r43_pairwise_correct"]/n
		m["action_only_pairwise_winner_accuracy"]=m["action_only_pairwise_correct"]/n
		m["state_shuffled_pairwise_winner_accuracy"]=m["state_shuffled_pairwise_correct"]/n
		m["mode_pairwise_winner_accuracy"]=m["mode_pairwise_correct"]/n
		m["r43_minus_action_only_accuracy"]=m["r43_pairwise_winner_accuracy"]-m["action_only_pairwise_winner_accuracy"]
		m["mode_minus_r43_accuracy"]=m["mode_pairwise_winner_accuracy"]-m["r43_pairwise_winner_accuracy"]
		m["r43_mean_decision_regret"]=m["r43_decision_regret_sum"]/n
		m["action_only_mean_decision_regret"]=m["action_only_decision_regret_sum"]/n
		m["state_shuffled_mean_decision_regret"]=m["state_shuffled_decision_regret_sum"]/n
		m["mode_mean_decision_regret"]=m["mode_decision_regret_sum"]/n
	}
	if m["pair_count"]!=108 || len(cases)!=108 || len(candidates)!=9 { m["invalid_row_count"]++ }
	if mode=="invalid" { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	exp:="WLM-LM-EXTERNAL-FUTURE-DATA-PAIRWISE-FIXED-BUDGET-ACTION-VALUE-RANKING-R44"
	if mode=="replication" { exp="WLM-LM-EXTERNAL-FUTURE-DATA-PAIRWISE-RANKING-DISJOINT-REPLICATION-R45" }
	if mode=="ablation" { exp="WLM-LM-EXTERNAL-FUTURE-DATA-PAIRWISE-STATE-COMPONENT-ABLATION-R45" }
	if mode=="causal-latent" { exp="WLM-LM-EXTERNAL-FUTURE-DATA-CAUSAL-LATENT-LEARNING-STATE-REPRESENTATION-R45" }
	return wlmLmR44Result{Schema:"wingless.research-scientific-result.v1",Experiment:exp,Mode:mode,Metrics:m,Cases:cases}
}
