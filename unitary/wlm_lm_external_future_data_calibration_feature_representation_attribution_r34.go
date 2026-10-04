package unitary

import "math"

type wlmLmR34Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	Base6PredictedAdvantage float64 `json:"base6_predicted_advantage"`
	DomainPredictedAdvantage float64 `json:"domain_predicted_advantage"`
	BudgetPredictedAdvantage float64 `json:"budget_predicted_advantage"`
	DomainBudgetPredictedAdvantage float64 `json:"domain_budget_predicted_advantage"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR34Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR34Case `json:"cases"`
}

type wlmLmR34ManifestInput struct {
	name string
	sources []wlmLmR27Source
}

type wlmLmR34Example struct {
	x []float64
	y float64
}

type wlmLmR34Agg struct {
	n int
	signErr int
	absErr float64
}

func wlmLmR34Vec(arm wlmLmR27Arm, hi int, domain, budget bool) []float64 {
	base:=arm.features[hi]
	x:=make([]float64,0,9)
	for i:=0;i<6;i++ { x=append(x,base[i]) }
	if domain {
		code,structured:=0.0,0.0
		if arm.d=="code" { code=1 }
		if arm.d=="structured" { structured=1 }
		x=append(x,code,structured)
	}
	if budget {
		maxBudget:=arm.budgets[len(arm.budgets)-1]
		if maxBudget<=0 { x=append(x,0) } else { x=append(x,float64(hi)/float64(maxBudget)) }
	}
	return x
}

func wlmLmR34Examples(mans []wlmLmR27Manifest, domain, budget bool) []wlmLmR34Example {
	out:=[]wlmLmR34Example{}
	for _,man:=range mans {
		for _,arm:=range man.arms {
			for i:=1;i<len(arm.budgets);i++ {
				lo,hi:=arm.budgets[i-1],arm.budgets[i]
				out=append(out,wlmLmR34Example{x:wlmLmR34Vec(arm,hi,domain,budget),y:float64(arm.hits[hi]-arm.hits[lo])})
			}
		}
	}
	return out
}

func wlmLmR34Fit(ex []wlmLmR34Example) ([]float64,bool) {
	if len(ex)==0 || len(ex[0].x)==0 { return nil,false }
	p:=len(ex[0].x)+1
	a:=make([][]float64,p)
	for i:=range a { a[i]=make([]float64,p+1) }
	for _,e:=range ex {
		if len(e.x)+1!=p { return nil,false }
		z:=make([]float64,p);z[0]=1
		copy(z[1:],e.x)
		for i:=0;i<p;i++ {
			for j:=0;j<p;j++ { a[i][j]+=z[i]*z[j] }
			a[i][p]+=z[i]*e.y
		}
	}
	for i:=1;i<p;i++ { a[i][i]+=1.0 }
	for col:=0;col<p;col++ {
		pivot:=col;best:=math.Abs(a[col][col])
		for r:=col+1;r<p;r++ {
			v:=math.Abs(a[r][col]);if v>best { best=v;pivot=r }
		}
		if best<1e-12 || math.IsNaN(best) || math.IsInf(best,0) { return nil,false }
		if pivot!=col { a[col],a[pivot]=a[pivot],a[col] }
		pv:=a[col][col]
		for j:=col;j<=p;j++ { a[col][j]/=pv }
		for r:=0;r<p;r++ {
			if r==col { continue }
			f:=a[r][col];if f==0 { continue }
			for j:=col;j<=p;j++ { a[r][j]-=f*a[col][j] }
		}
	}
	beta:=make([]float64,p)
	for i:=0;i<p;i++ {
		beta[i]=a[i][p]
		if math.IsNaN(beta[i])||math.IsInf(beta[i],0) { return nil,false }
	}
	return beta,true
}

func wlmLmR34Predict(beta,x []float64) float64 {
	if len(beta)!=len(x)+1 { return math.NaN() }
	v:=beta[0]
	for i:=range x { v+=beta[i+1]*x[i] }
	return v
}

func wlmLmR34PolicyPred(man wlmLmR27Manifest,beta []float64,domain,budget bool,alloc [3]int) float64 {
	total:=0.0
	for ai,arm:=range man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			if hi<=alloc[ai] { total+=wlmLmR34Predict(beta,wlmLmR34Vec(arm,hi,domain,budget)) }
		}
	}
	return total
}

func wlmLmR34Update(a *wlmLmR34Agg,pred,actual float64) {
	a.n++
	if wlmLmR29Sign(pred)!=wlmLmR29Sign(actual) { a.signErr++ }
	a.absErr+=math.Abs(pred-actual)
}

func wlmLmR34Finalize(prefix string,a wlmLmR34Agg,m map[string]float64) {
	if a.n<=0 { m["invalid_row_count"]++;return }
	m[prefix+"_sign_error_rate"]=float64(a.signErr)/float64(a.n)
	m[prefix+"_mean_absolute_error"]=a.absErr/float64(a.n)
}

func wlmLmR34Ratio(num,den float64,m map[string]float64) float64 {
	if den<=0 { m["invalid_row_count"]++;return 0 }
	return num/den
}

func RunWlmLmExternalFutureDataCalibrationFeatureRepresentationAttributionR34(
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
	inputs:=[]wlmLmR34ManifestInput{
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
	var baseAgg,domainAgg,budgetAgg,bothAgg wlmLmR34Agg
	cases:=make([]wlmLmR34Case,0,33)
	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainEx:=[]wlmLmR27Example{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r34_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man);trainEx=append(trainEx,ex...)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		baseR27,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		base,ok0:=wlmLmR34Fit(wlmLmR34Examples(trainMans,false,false))
		domain,ok1:=wlmLmR34Fit(wlmLmR34Examples(trainMans,true,false))
		budget,ok2:=wlmLmR34Fit(wlmLmR34Examples(trainMans,false,true))
		both,ok3:=wlmLmR34Fit(wlmLmR34Examples(trainMans,true,true))
		if !ok0||!ok1||!ok2||!ok3 { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		heldFeatures,_:=wlmLmR27Build("r34_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		type frozenCase struct{ alloc [3]int;excluded bool;p0,p1,p2,p3 float64 }
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
			p0:=wlmLmR34PolicyPred(heldFeatures,base,false,false,alloc)-wlmLmR34PolicyPred(heldFeatures,base,false,false,wlmLmR27Equal)
			r27:=wlmLmR27PolicyPred(heldFeatures,baseR27,alloc)-wlmLmR27PolicyPred(heldFeatures,baseR27,wlmLmR27Equal)
			if math.Abs(p0-r27)>1e-9 { m["base6_identity_mismatch_count"]++ }
			p1:=wlmLmR34PolicyPred(heldFeatures,domain,true,false,alloc)-wlmLmR34PolicyPred(heldFeatures,domain,true,false,wlmLmR27Equal)
			p2:=wlmLmR34PolicyPred(heldFeatures,budget,false,true,alloc)-wlmLmR34PolicyPred(heldFeatures,budget,false,true,wlmLmR27Equal)
			p3:=wlmLmR34PolicyPred(heldFeatures,both,true,true,alloc)-wlmLmR34PolicyPred(heldFeatures,both,true,true,wlmLmR27Equal)
			frozen=append(frozen,frozenCase{alloc:alloc,excluded:!stableNonZero,p0:p0,p1:p1,p2:p2,p3:p3})
		}
		heldActual,_:=wlmLmR27Build("r34_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			if !fc.excluded { continue }
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR34Update(&baseAgg,fc.p0,actual);wlmLmR34Update(&domainAgg,fc.p1,actual)
			wlmLmR34Update(&budgetAgg,fc.p2,actual);wlmLmR34Update(&bothAgg,fc.p3,actual)
			cases=append(cases,wlmLmR34Case{
				HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,
				Base6PredictedAdvantage:fc.p0,DomainPredictedAdvantage:fc.p1,
				BudgetPredictedAdvantage:fc.p2,DomainBudgetPredictedAdvantage:fc.p3,
				ActualAdvantage:actual,
			})
		}
	}
	wlmLmR34Finalize("base6",baseAgg,m);wlmLmR34Finalize("domain",domainAgg,m)
	wlmLmR34Finalize("budget",budgetAgg,m);wlmLmR34Finalize("domain_budget",bothAgg,m)
	m["domain_main_sign_error_improvement_min"]=math.Min(m["base6_sign_error_rate"]-m["domain_sign_error_rate"],m["budget_sign_error_rate"]-m["domain_budget_sign_error_rate"])
	m["domain_main_mae_ratio_max"]=math.Max(wlmLmR34Ratio(m["domain_mean_absolute_error"],m["base6_mean_absolute_error"],m),wlmLmR34Ratio(m["domain_budget_mean_absolute_error"],m["budget_mean_absolute_error"],m))
	m["budget_main_sign_error_improvement_min"]=math.Min(m["base6_sign_error_rate"]-m["budget_sign_error_rate"],m["domain_sign_error_rate"]-m["domain_budget_sign_error_rate"])
	m["budget_main_mae_ratio_max"]=math.Max(wlmLmR34Ratio(m["budget_mean_absolute_error"],m["base6_mean_absolute_error"],m),wlmLmR34Ratio(m["domain_budget_mean_absolute_error"],m["domain_mean_absolute_error"],m))
	m["interaction_sign_error_improvement"]=m["base6_sign_error_rate"]-m["domain_budget_sign_error_rate"]
	m["interaction_mae_ratio"]=wlmLmR34Ratio(m["domain_budget_mean_absolute_error"],m["base6_mean_absolute_error"],m)
	for k,v:=range m {
		if len(k)>31 && k[len(k)-31:]=="_source_identity_mismatch_count" && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 || int(m["excluded_case_count"])!=33 || len(cases)!=33 {
		m["invalid_row_count"]++
	}
	if m["base6_identity_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR34Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-CALIBRATION-FEATURE-REPRESENTATION-ATTRIBUTION-R34",Metrics:m,Cases:cases}
}
