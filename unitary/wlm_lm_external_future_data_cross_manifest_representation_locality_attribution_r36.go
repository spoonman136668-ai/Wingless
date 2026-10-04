package unitary

import (
	"math"
	"sort"
	"strings"
)

type wlmLmR36Case struct {
	HoldoutManifest string `json:"holdout_manifest"`
	Allocation [3]int `json:"allocation"`
	LocalTrainingManifests [2]string `json:"local_training_manifests"`
	Global3PredictedAdvantage float64 `json:"global3_predicted_advantage"`
	Local2PredictedAdvantage float64 `json:"local2_predicted_advantage"`
	ActualAdvantage float64 `json:"actual_advantage"`
}

type wlmLmR36Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR36Case `json:"cases"`
}

type wlmLmR36Distance struct {
	pos int
	manifestIndex int
	d float64
}

func wlmLmR36ManifestDistance(a,b wlmLmR27Manifest) float64 {
	if len(a.arms)!=len(b.arms) { return math.Inf(1) }
	total:=0.0
	for i:=range a.arms {
		ra:=wlmLmR35Raw(a.arms[i].eval)
		rb:=wlmLmR35Raw(b.arms[i].eval)
		for j:=0;j<5;j++ {
			x:=ra[j]-rb[j]
			total+=x*x
		}
	}
	return total
}

func RunWlmLmExternalFutureDataCrossManifestRepresentationLocalityAttributionR36(
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
		"locality_selection_mismatch_count":0,
		"heldout_outcome_use_before_locality_freeze":0,
		"post_result_locality_choice_count":0,"source_identity_mismatch_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	var globalAgg,localAgg wlmLmR34Agg
	cases:=make([]wlmLmR36Case,0,33)
	for hold:=0;hold<4;hold++ {
		trainMans:=make([]wlmLmR27Manifest,0,3)
		trainIdx:=make([]int,0,3)
		trainEx:=[]wlmLmR27Example{}
		for j:=0;j<4;j++ {
			if j==hold { continue }
			man,ex:=wlmLmR27Build("r36_train_"+inputs[j].name+"_for_"+inputs[hold].name,inputs[j].sources,j,true,m)
			trainMans=append(trainMans,man);trainIdx=append(trainIdx,j);trainEx=append(trainEx,ex...)
		}
		if len(trainEx)!=36 { m["invalid_row_count"]++ }
		baseR27,ok:=wlmLmR27Fit(trainEx,-1)
		if !ok { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		global,ok0:=wlmLmR34Fit(wlmLmR35Examples(trainMans,false))
		if !ok0 { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		heldFeatures,_:=wlmLmR27Build("r36_held_features_"+inputs[hold].name,inputs[hold].sources,hold,false,m)
		ds:=make([]wlmLmR36Distance,0,3)
		for pos,man:=range trainMans {
			ds=append(ds,wlmLmR36Distance{pos:pos,manifestIndex:trainIdx[pos],d:wlmLmR36ManifestDistance(heldFeatures,man)})
		}
		sort.Slice(ds,func(i,j int)bool{
			if ds[i].d!=ds[j].d { return ds[i].d<ds[j].d }
			return ds[i].manifestIndex<ds[j].manifestIndex
		})
		if len(ds)!=3 || ds[0].pos==ds[1].pos || ds[0].manifestIndex==hold || ds[1].manifestIndex==hold {
			m["locality_selection_mismatch_count"]++
		}
		localMans:=[]wlmLmR27Manifest{trainMans[ds[0].pos],trainMans[ds[1].pos]}
		local,ok1:=wlmLmR34Fit(wlmLmR35Examples(localMans,false))
		if !ok1 { m["fit_failure_count"]++;m["invalid_row_count"]++ }
		localNames:=[2]string{inputs[ds[0].manifestIndex].name,inputs[ds[1].manifestIndex].name}
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
			p0:=wlmLmR35PolicyPred(heldFeatures,global,false,alloc)-wlmLmR35PolicyPred(heldFeatures,global,false,wlmLmR27Equal)
			r27:=wlmLmR27PolicyPred(heldFeatures,baseR27,alloc)-wlmLmR27PolicyPred(heldFeatures,baseR27,wlmLmR27Equal)
			if math.Abs(p0-r27)>1e-9 { m["base6_identity_mismatch_count"]++ }
			p1:=wlmLmR35PolicyPred(heldFeatures,local,false,alloc)-wlmLmR35PolicyPred(heldFeatures,local,false,wlmLmR27Equal)
			frozen=append(frozen,frozenCase{alloc:alloc,excluded:!stableNonZero,p0:p0,p1:p1})
		}
		heldActual,_:=wlmLmR27Build("r36_held_actual_"+inputs[hold].name,inputs[hold].sources,hold,true,m)
		for _,fc:=range frozen {
			if !fc.excluded { continue }
			actual:=wlmLmR27PolicyActual(heldActual,fc.alloc)-wlmLmR27PolicyActual(heldActual,wlmLmR27Equal)
			wlmLmR34Update(&globalAgg,fc.p0,actual);wlmLmR34Update(&localAgg,fc.p1,actual)
			cases=append(cases,wlmLmR36Case{
				HoldoutManifest:inputs[hold].name,Allocation:fc.alloc,LocalTrainingManifests:localNames,
				Global3PredictedAdvantage:fc.p0,Local2PredictedAdvantage:fc.p1,ActualAdvantage:actual,
			})
		}
	}
	wlmLmR34Finalize("global3",globalAgg,m);wlmLmR34Finalize("local2",localAgg,m)
	m["local2_sign_error_improvement"]=m["global3_sign_error_rate"]-m["local2_sign_error_rate"]
	m["local2_mae_ratio"]=wlmLmR34Ratio(m["local2_mean_absolute_error"],m["global3_mean_absolute_error"],m)
	for k,v:=range m {
		if strings.HasSuffix(k,"_source_identity_mismatch_count") && k!="source_identity_mismatch_count" {
			m["source_identity_mismatch_count"]+=v
		}
	}
	if int(m["non_equal_case_count"])!=36 || int(m["stable_nonzero_case_count"])!=3 || int(m["excluded_case_count"])!=33 || len(cases)!=33 {
		m["invalid_row_count"]++
	}
	if m["base6_identity_mismatch_count"]!=0 || m["locality_selection_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR36Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-CROSS-MANIFEST-REPRESENTATION-LOCALITY-ATTRIBUTION-R36",
		Metrics:m,Cases:cases,
	}
}
