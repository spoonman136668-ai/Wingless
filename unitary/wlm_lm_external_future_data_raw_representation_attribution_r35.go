package unitary

import (
	"math"
	"strings"
)

type wlmLmR35Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	Base6PredictedAdvantage float64 `json:"base6_predicted_advantage"`
	RawBytePredictedAdvantage float64 `json:"rawbyte_predicted_advantage"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR35Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR35Case `json:"cases"`
}

type wlmLmR35ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

func wlmLmR35Raw(eval []byte) [5]float64 {
	var out [5]float64
	if len(eval)==0 { return out }
	var counts [256]int
	ws,letters,digits,punct:=0,0,0,0
	for _,b:=range eval {
		counts[int(b)]++
		if b==9||b==10||b==13||b==32 { ws++ }
		if (b>='A'&&b<='Z')||(b>='a'&&b<='z') { letters++ }
		if b>='0'&&b<='9' { digits++ }
		if strings.ContainsRune("{}[]():,.;=_-+/\\<>\"'\x60@#$%^&*|!?~",rune(b)) { punct++ }
	}
	n:=float64(len(eval))
	h:=0.0
	for _,c:=range counts {
		if c==0 { continue }
		p:=float64(c)/n
		h-=p*math.Log2(p)
	}
	out[0]=float64(ws)/n
	out[1]=float64(letters)/n
	out[2]=float64(digits)/n
	out[3]=float64(punct)/n
	out[4]=h/8.0
	return out
}

func wlmLmR35Vec(arm wlmLmR27Arm, hi int, raw bool) []float64 {
	base:=arm.features[hi]
	x:=make([]float64,0,11)
	for i:=0;i<6;i++ { x=append(x,base[i]) }
	if raw {
		r:=wlmLmR35Raw(arm.eval)
		for i:=0;i<5;i++ { x=append(x,r[i]) }
	}
	return x
}

func wlmLmR35Examples(mans []wlmLmR27Manifest, raw bool) []wlmLmR34Example {
	out:=[]wlmLmR34Example{}
	for _,man:=range mans {
		for _,arm:=range man.arms {
			for i:=1;i<len(arm.budgets);i++ {
				lo,hi:=arm.budgets[i-1],arm.budgets[i]
				out=append(out,wlmLmR34Example{x:wlmLmR35Vec(arm,hi,raw),y:float64(arm.hits[hi]-arm.hits[lo])})
			}
		}
	}
	return out
}

func wlmLmR35PolicyPred(man wlmLmR27Manifest,beta []float64,raw bool,alloc [3]int) float64 {
	total:=0.0
	for ai,arm:=range man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			if hi<=alloc[ai] { total+=wlmLmR34Predict(beta,wlmLmR35Vec(arm,hi,raw)) }
		}
	}
	return total
}

func RunWlmLmExternalFutureDataRawRepresentationAttributionR35(
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
	inputs:=[]wlmLmR35ManifestInput{
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
		"non_equal_case_count":0,"stable_nonzero_case_count":0,"excluded_case_count":0,
		"fit_failure_count":0,"base6_identity_mismatch_count":0,
		"heldout_outcome_use_before_representation_freeze":0,
		"post_result_representation_choice_count":0,"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	var baseAgg,rawAgg wlmLmR34Agg
	cases:=make([]wlmLmR35Case,0,33)
	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainEx:=[]wlmLmR27Example{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r35_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man);trainEx=append(trainEx,ex...)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		baseR27,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		base,ok0:=wlmLmR34Fit(wlmLmR35Examples(trainMans,false))
		raw,ok1:=wlmLmR34Fit(wlmLmR35Examples(trainMans,true))
		if !ok0||!ok1 { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		heldFeatures,_:=wlmLmR27Build("r35_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct{alloc [3]int;excluded bool;p0,p1 float64}
		frozen:=make([]frozenCase,0,9)
		for _,alloc:=range wlmLmR27Candidates {
			if alloc==wlmLmR27Equal { continue }
			m["non_equal_case_count"]++
			signs:=map[int]bool{}
			for _,tm:=range trainMans {
				adv:=wlmLmR27PolicyActual(tm,alloc)-wlmLmR27PolicyActual(tm,wlmLmR27Equal)
				signs[wlmLmR29Sign(adv)]=true
			}
			stableNonZero:=len(signs)==1 && (signs[1]||signs[-1])
			if stableNonZero { m["stable_nonzero_case_count"]++ } else { m["excluded_case_count"]++ }
			p0:=wlmLmR35PolicyPred(heldFeatures,base,false,alloc)-wlmLmR35PolicyPred(heldFeatures,base,false,wlmLmR27Equal)
			r27:=wlmLmR27PolicyPred(heldFeatures,baseR27,alloc)-wlmLmR27PolicyPred(heldFeatures,baseR27,wlmLmR27Equal)
			if math.Abs(p0-r27)>1e-9 { m["base6_identity_mismatch_count"]++ }
			p1:=wlmLmR35PolicyPred(heldFeatures,raw,true,alloc)-wlmLmR35PolicyPred(heldFeatures,raw,true,wlmLmR27Equal)
			frozen=append(frozen,frozenCase{alloc:alloc,excluded:!stableNonZero,p0:p0,p1:p1})
		}
		heldActual,_:=wlmLmR27Build("r35_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			if !fc.excluded { continue }
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR34Update(&baseAgg,fc.p0,actual);wlmLmR34Update(&rawAgg,fc.p1,actual)
			cases=append(cases,wlmLmR35Case{
				HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,
				Base6PredictedAdvantage:fc.p0,RawBytePredictedAdvantage:fc.p1,ActualAdvantage:actual,
			})
		}
	}
	wlmLmR34Finalize("base6",baseAgg,m);wlmLmR34Finalize("rawbyte",rawAgg,m)
	m["rawbyte_sign_error_improvement"]=m["base6_sign_error_rate"]-m["rawbyte_sign_error_rate"]
	m["rawbyte_mae_ratio"]=wlmLmR34Ratio(m["rawbyte_mean_absolute_error"],m["base6_mean_absolute_error"],m)
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 || int(m["excluded_case_count"])!=33 || len(cases)!=33 {
		m["invalid_row_count"]++
	}
	if m["base6_identity_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR35Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-RAW-REPRESENTATION-ATTRIBUTION-R35",
		Metrics:m,Cases:cases,
	}
}
