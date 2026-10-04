package unitary

import "math"

type wlmLmR27Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Decision string `json:"decision"`
	SelectedAllocation [3]int `json:"selected_allocation"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmR27Example struct {
	manifest int
	x [6]float64
	y float64
}

type wlmLmR27Arm struct {
	d string
	eval []byte
	budgets []int
	bases map[int][256][256]uint32
	sels map[int]map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif
	features map[int][6]float64
	hits map[int]int
}

type wlmLmR27Manifest struct {
	name string
	arms []wlmLmR27Arm
}

type wlmLmR27Source struct {
	d string
	b []byte
	h string
	n int
	budgets []int
}

var wlmLmR27Candidates = [][3]int{
	{582,436,726},{582,581,581},{582,726,436},{727,291,726},{727,436,581},
	{727,581,436},{727,726,291},{872,291,581},{872,436,436},{872,581,291},
}

var wlmLmR27Equal = [3]int{582,581,581}

func wlmLmR27LexLess(a,b [3]int) bool {
	if a[0]!=b[0] { return a[0]<b[0] }
	if a[1]!=b[1] { return a[1]<b[1] }
	return a[2]<b[2]
}

func wlmLmR27Build(name string, ss []wlmLmR27Source, manifestIndex int, openHistoricalEval bool, m map[string]float64) (wlmLmR27Manifest, []wlmLmR27Example) {
	out:=wlmLmR27Manifest{name:name,arms:make([]wlmLmR27Arm,0,3)}
	ex:=[]wlmLmR27Example{}
	for _,s:=range ss {
		m[name+"_total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n || wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h {
			m[name+"_source_identity_mismatch_count"]++
		}
	}
	for ti,t:=range ss {
		idx:=[]int{}
		for i:=range ss { if i!=ti { idx=append(idx,i) } }
		a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b)
		if len(a)+len(b)!=15819 || len(t.b)<1454 { m["invalid_row_count"]++; continue }
		eval,res:=t.b[:582],t.b[582:1454]
		arm:=wlmLmR27Arm{
			d:t.d,eval:eval,budgets:t.budgets,
			bases:map[int][256][256]uint32{},
			sels:map[int]map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif{},
			features:map[int][6]float64{},hits:map[int]int{},
		}
		for _,bud:=range t.budgets {
			train:=[][]byte{a,b}
			if bud>0 { train=append(train,res[:bud]) }
			base,sel:=wlmLmR10Train256(train,m)
			if len(sel)!=256 { m["invalid_row_count"]++ }
			arm.bases[bud]=base
			arm.sels[bud]=sel
		}
		for i:=1;i<len(t.budgets);i++ {
			lo,hi:=t.budgets[i-1],t.budgets[i]
			lb,hb:=arm.bases[lo],arm.bases[hi]
			ls,hs:=arm.sels[lo],arm.sels[hi]
			bc,en,exit,best:=wlmLmR18StateDelta(&lb,&hb,ls,hs)
			marginSum:=uint64(0)
			for p:=0;p<256;p++ {
				if wlmLmR22Argmax(&lb[p])!=wlmLmR22Argmax(&hb[p]) {
					marginSum+=uint64(wlmLmR22TopMargin(&hb[p]))
				}
			}
			active:=0
			for k,motif:=range hs {
				if _,ok:=ls[k]; ok { continue }
				if motif.best!=wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(&hb,k[3]) { active++ }
			}
			for k,motif:=range ls {
				if _,ok:=hs[k]; ok { continue }
				if motif.best!=wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(&lb,k[3]) { active++ }
			}
			w:=float64(hi-lo)
			arm.features[hi]=[6]float64{
				float64(bc)/w,
				float64(marginSum)/w,
				float64(en)/w,
				float64(exit)/w,
				float64(active)/w,
				float64(best)/w,
			}
		}
		if openHistoricalEval {
			for _,bud:=range t.budgets {
				bb:=arm.bases[bud]
				arm.hits[bud]=wlmLmR10Hits(eval,&bb,arm.sels[bud])
			}
			for i:=1;i<len(t.budgets);i++ {
				lo,hi:=t.budgets[i-1],t.budgets[i]
				ex=append(ex,wlmLmR27Example{manifest:manifestIndex,x:arm.features[hi],y:float64(arm.hits[hi]-arm.hits[lo])})
			}
		}
		out.arms=append(out.arms,arm)
	}
	if len(out.arms)!=3 { m["invalid_row_count"]++ }
	return out,ex
}

func wlmLmR27Fit(ex []wlmLmR27Example, omit int) ([7]float64,bool) {
	var a [7][8]float64
	count:=0
	for _,e:=range ex {
		if e.manifest==omit { continue }
		count++
		var z [7]float64
		z[0]=1
		for i:=0;i<6;i++ { z[i+1]=e.x[i] }
		for i:=0;i<7;i++ {
			for j:=0;j<7;j++ { a[i][j]+=z[i]*z[j] }
			a[i][7]+=z[i]*e.y
		}
	}
	if count==0 { return [7]float64{},false }
	for i:=1;i<7;i++ { a[i][i]+=1.0 }
	for col:=0;col<7;col++ {
		pivot:=col
		best:=math.Abs(a[col][col])
		for r:=col+1;r<7;r++ {
			v:=math.Abs(a[r][col])
			if v>best { best=v;pivot=r }
		}
		if best<1e-12 || math.IsNaN(best) || math.IsInf(best,0) { return [7]float64{},false }
		if pivot!=col { a[col],a[pivot]=a[pivot],a[col] }
		pv:=a[col][col]
		for j:=col;j<8;j++ { a[col][j]/=pv }
		for r:=0;r<7;r++ {
			if r==col { continue }
			f:=a[r][col]
			if f==0 { continue }
			for j:=col;j<8;j++ { a[r][j]-=f*a[col][j] }
		}
	}
	var beta [7]float64
	for i:=0;i<7;i++ {
		beta[i]=a[i][7]
		if math.IsNaN(beta[i])||math.IsInf(beta[i],0) { return [7]float64{},false }
	}
	return beta,true
}

func wlmLmR27Predict(beta [7]float64,x [6]float64) float64 {
	v:=beta[0]
	for i:=0;i<6;i++ { v+=beta[i+1]*x[i] }
	return v
}

func wlmLmR27PolicyPred(man wlmLmR27Manifest,beta [7]float64,alloc [3]int) float64 {
	total:=0.0
	for ai,arm:=range man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			if hi<=alloc[ai] { total+=wlmLmR27Predict(beta,arm.features[hi]) }
		}
	}
	return total
}

func wlmLmR27PolicyActual(man wlmLmR27Manifest,alloc [3]int) float64 {
	total:=0.0
	for ai,arm:=range man.arms { total+=float64(arm.hits[alloc[ai]]) }
	return total
}

func RunWlmLmExternalFutureDataProspectiveCalibrationR27(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse []byte,
) interface{} {
	budCode:=[]int{0,436,582,727,872}
	budOther:=[]int{0,291,436,581,726}
	mk:=func(d string,b []byte,h string,n int,bud []int) wlmLmR27Source { return wlmLmR27Source{d:d,b:b,h:h,n:n,budgets:bud} }
	transfer:=[]wlmLmR27Source{
		mk("code",transferCode,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45",41453,budCode),
		mk("structured",transferStructured,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932",14365,budOther),
		mk("technical_prose",transferProse,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b",1454,budOther),
	}
	third:=[]wlmLmR27Source{
		mk("code",thirdCode,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff",41453,budCode),
		mk("structured",thirdStructured,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56",14365,budOther),
		mk("technical_prose",thirdProse,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017",1454,budOther),
	}
	fourth:=[]wlmLmR27Source{
		mk("code",fourthCode,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1",41453,budCode),
		mk("structured",fourthStructured,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933",14365,budOther),
		mk("technical_prose",fourthProse,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222",1454,budOther),
	}
	fifth:=[]wlmLmR27Source{
		mk("code",fifthCode,"504b68653b5478b88216f6342a74bacc5982549657005fb486dd00d753b4ea9a",41453,budCode),
		mk("structured",fifthStructured,"5c0f3a215ba35b7fbcaae212d27a33ba5109e16d89809987d16fa89072534281",14365,budOther),
		mk("technical_prose",fifthProse,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49",1454,budOther),
	}
	m:=map[string]float64{
		"historical_manifest_count":3,"historical_example_count":0,"historical_uncertainty_radius":0,
		"historical_policy_advantage_rmse":0,"historical_policy_sign_accuracy":0,
		"fifth_manifest_policy_evaluation_byte_use_count":0,"fifth_manifest_policy_evaluation_label_use_count":0,
		"post_result_policy_choice_count":0,"candidate_allocation_count":10,"total_adaptation_budget_equal":1744,
		"total_adaptation_budget_selected":0,"selected_predicted_advantage_full":0,"selected_predicted_advantage_min_loo":0,
		"selected_code_budget":0,"selected_structured_budget":0,"selected_technical_prose_budget":0,
		"decision_reallocate":0,"decision_fallback_equal":0,
		"packet0_aggregate_exact_hits":0,"equal_aggregate_exact_hits":0,"selected_aggregate_exact_hits":0,
		"actual_aggregate_exact_hit_delta_vs_equal":0,"selected_domain_below_packet0_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,
		"invalid_row_count":0,
	}
	m1,e1:=wlmLmR27Build("historical_transfer",transfer,0,true,m)
	m2,e2:=wlmLmR27Build("historical_third",third,1,true,m)
	m3,e3:=wlmLmR27Build("historical_fourth",fourth,2,true,m)
	ex:=append(append(e1,e2...),e3...)
	m["historical_example_count"]=float64(len(ex))
	if len(ex)!=36 { m["invalid_row_count"]++ }
	historical:=[]wlmLmR27Manifest{m1,m2,m3}
	loo:=make([][7]float64,3)
	uncertainty:=0.0
	sq:=0.0;signOK:=0.0;signN:=0.0
	for hold:=0;hold<3;hold++ {
		beta,ok:=wlmLmR27Fit(ex,hold)
		if !ok { m["invalid_row_count"]++;continue }
		loo[hold]=beta
		for _,c:=range wlmLmR27Candidates {
			pred:=wlmLmR27PolicyPred(historical[hold],beta,c)-wlmLmR27PolicyPred(historical[hold],beta,wlmLmR27Equal)
			actual:=wlmLmR27PolicyActual(historical[hold],c)-wlmLmR27PolicyActual(historical[hold],wlmLmR27Equal)
			over:=pred-actual
			if over>uncertainty { uncertainty=over }
			d:=pred-actual;sq+=d*d;signN++
			if (pred>0&&actual>0)||(pred<0&&actual<0)||(pred==0&&actual==0) { signOK++ }
		}
	}
	if uncertainty<0 { uncertainty=0 }
	m["historical_uncertainty_radius"]=uncertainty
	if signN>0 {
		m["historical_policy_advantage_rmse"]=math.Sqrt(sq/signN)
		m["historical_policy_sign_accuracy"]=signOK/signN
	}
	full,ok:=wlmLmR27Fit(ex,-1)
	if !ok { m["invalid_row_count"]++ }
	fifthMan,_:=wlmLmR27Build("fifth",fifth,3,false,m)
	best:=wlmLmR27Equal
	bestAdv:=math.Inf(-1)
	for _,c:=range wlmLmR27Candidates {
		adv:=wlmLmR27PolicyPred(fifthMan,full,c)-wlmLmR27PolicyPred(fifthMan,full,wlmLmR27Equal)
		if adv>bestAdv || (adv==bestAdv && wlmLmR27LexLess(c,best)) { bestAdv=adv;best=c }
	}
	minLOO:=math.Inf(1)
	for i:=0;i<3;i++ {
		adv:=wlmLmR27PolicyPred(fifthMan,loo[i],best)-wlmLmR27PolicyPred(fifthMan,loo[i],wlmLmR27Equal)
		if adv<minLOO { minLOO=adv }
	}
	m["selected_predicted_advantage_full"]=bestAdv
	m["selected_predicted_advantage_min_loo"]=minLOO
	decision:="FALLBACK_EQUAL"
	selected:=wlmLmR27Equal
	if best!=wlmLmR27Equal && bestAdv>uncertainty && minLOO>uncertainty {
		decision="REALLOCATE";selected=best;m["decision_reallocate"]=1
	} else { m["decision_fallback_equal"]=1 }
	m["selected_code_budget"]=float64(selected[0]);m["selected_structured_budget"]=float64(selected[1]);m["selected_technical_prose_budget"]=float64(selected[2])
	m["total_adaptation_budget_selected"]=float64(selected[0]+selected[1]+selected[2])
	if m["total_adaptation_budget_selected"]!=1744 { m["invalid_row_count"]++ }
	// The prospective policy is frozen above. Evaluation begins only here.
	for ai,arm:=range fifthMan.arms {
		p0b:=arm.bases[0];eqb:=arm.bases[wlmLmR27Equal[ai]];sb:=arm.bases[selected[ai]]
		p0:=wlmLmR10Hits(arm.eval,&p0b,arm.sels[0])
		eq:=wlmLmR10Hits(arm.eval,&eqb,arm.sels[wlmLmR27Equal[ai]])
		sel:=wlmLmR10Hits(arm.eval,&sb,arm.sels[selected[ai]])
		m["arm_"+arm.d+"_packet0_exact_hits"]=float64(p0)
		m["arm_"+arm.d+"_equal_exact_hits"]=float64(eq)
		m["arm_"+arm.d+"_selected_exact_hits"]=float64(sel)
		m["arm_"+arm.d+"_selected_delta_vs_equal"]=float64(sel-eq)
		m["packet0_aggregate_exact_hits"]+=float64(p0)
		m["equal_aggregate_exact_hits"]+=float64(eq)
		m["selected_aggregate_exact_hits"]+=float64(sel)
		if sel<p0 { m["selected_domain_below_packet0_count"]++ }
	}
	m["actual_aggregate_exact_hit_delta_vs_equal"]=m["selected_aggregate_exact_hits"]-m["equal_aggregate_exact_hits"]
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR27Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-PROSPECTIVE-CALIBRATION-R27",
		Decision:decision,SelectedAllocation:selected,Metrics:m,
	}
}
