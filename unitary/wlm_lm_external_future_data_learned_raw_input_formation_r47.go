package unitary

import (
	"math"
	"strings"
)

type wlmLmR47CandidateSummary struct {
	ID string `json:"id"`
	LearnedScalars int `json:"learned_scalars"`
	Alphas []float64 `json:"alphas"`
	HistoricalPairwiseAccuracy float64 `json:"historical_pairwise_accuracy"`
	ComplexityPenalty float64 `json:"complexity_penalty"`
	SearchScore float64 `json:"search_score"`
	ResourceEligible bool `json:"resource_eligible"`
}

type wlmLmR47PairCase struct {
	Manifest string `json:"manifest"`
	ActionA [3]int `json:"action_a"`
	ActionB [3]int `json:"action_b"`
	LearnedChoice [3]int `json:"learned_choice"`
	BaselineChoice [3]int `json:"baseline_choice"`
	ActualChoice [3]int `json:"actual_choice"`
	ActualDelta float64 `json:"actual_delta_a_minus_b"`
	Tie bool `json:"tie"`
}

type wlmLmR47Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SelectedCandidate string `json:"selected_candidate"`
	SelectedAlphas []float64 `json:"selected_alphas"`
	CandidateHistory []wlmLmR47CandidateSummary `json:"candidate_history"`
	Metrics map[string]float64 `json:"metrics"`
	Cases []wlmLmR47PairCase `json:"cases"`
}

type wlmLmR47Prepared struct {
	man wlmLmR27Manifest
	sources []wlmLmR27Source
}

type wlmLmR47Example struct {
	manifest int
	x [6]float64
	y float64
}

type wlmLmR47CandidateSpec struct {
	id string
	groups [6]int
	scalars int
}

var wlmLmR47Grid=[]float64{0.125,0.25,0.5,0.75,1.0}

var wlmLmR47Specs=[]wlmLmR47CandidateSpec{
	{id:"shared-alpha",groups:[6]int{0,0,0,0,0,0},scalars:1},
	{id:"two-group-alpha",groups:[6]int{0,0,1,1,1,1},scalars:2},
	{id:"three-group-alpha",groups:[6]int{0,0,1,1,2,2},scalars:3},
	{id:"fieldwise-alpha",groups:[6]int{0,1,2,3,4,5},scalars:6},
}

func wlmLmR47RawStats(b []byte)[6]float64 {
	if len(b)==0 { return [6]float64{} }
	sum:=0.0
	printable,digits,space,repeats:=0,0,0,0
	for i,v:=range b {
		x:=float64(v)
		sum+=x
		if v>=32&&v<=126 { printable++ }
		if v>='0'&&v<='9' { digits++ }
		if v==' '||v=='\t'||v=='\r'||v=='\n' { space++ }
		if i>0&&v==b[i-1] { repeats++ }
	}
	mean:=sum/float64(len(b))
	varsq:=0.0
	for _,v:=range b { d:=float64(v)-mean;varsq+=d*d }
	varsq/=float64(len(b))
	den:=float64(len(b))
	repDen:=float64(len(b)-1);if repDen<1 { repDen=1 }
	return [6]float64{
		mean/255.0,
		varsq/(255.0*255.0),
		float64(printable)/den,
		float64(digits)/den,
		float64(space)/den,
		float64(repeats)/repDen,
	}
}

func wlmLmR47Alpha(spec wlmLmR47CandidateSpec, vals []float64)[6]float64 {
	var out [6]float64
	for i:=0;i<6;i++ { out[i]=vals[spec.groups[i]] }
	return out
}

func wlmLmR47States(p wlmLmR47Prepared, alpha [6]float64, m map[string]float64)[3]map[int][6]float64 {
	var out [3]map[int][6]float64
	if len(p.man.arms)!=3||len(p.sources)!=3 { m["invalid_row_count"]++;return out }
	for ai,arm:=range p.man.arms {
		out[ai]=map[int][6]float64{}
		src:=p.sources[ai].b
		if len(src)<1454 { m["raw_state_source_mismatch_count"]++;continue }
		res:=src[582:1454]
		var state [6]float64
		for i:=1;i<len(arm.budgets);i++ {
			lo,hi:=arm.budgets[i-1],arm.budgets[i]
			if lo<0||hi>len(res)||lo>=hi { m["raw_state_source_mismatch_count"]++;continue }
			x:=wlmLmR47RawStats(res[lo:hi])
			for q:=0;q<6;q++ { state[q]=(1-alpha[q])*state[q]+alpha[q]*x[q] }
			out[ai][hi]=state
		}
	}
	return out
}

func wlmLmR47Examples(p wlmLmR47Prepared, states [3]map[int][6]float64, manifest int)[]wlmLmR47Example {
	out:=[]wlmLmR47Example{}
	for ai,arm:=range p.man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			lo,hi:=arm.budgets[i-1],arm.budgets[i]
			x,ok:=states[ai][hi];if !ok { continue }
			out=append(out,wlmLmR47Example{manifest:manifest,x:x,y:float64(arm.hits[hi]-arm.hits[lo])})
		}
	}
	return out
}

func wlmLmR47Fit(ex []wlmLmR47Example, omit int)([7]float64,bool) {
	var a [7][8]float64
	count:=0
	for _,e:=range ex {
		if e.manifest==omit { continue }
		count++
		var z [7]float64;z[0]=1
		for i:=0;i<6;i++ { z[i+1]=e.x[i] }
		for i:=0;i<7;i++ {
			for j:=0;j<7;j++ { a[i][j]+=z[i]*z[j] }
			a[i][7]+=z[i]*e.y
		}
	}
	if count==0 { return [7]float64{},false }
	for i:=1;i<7;i++ { a[i][i]+=1.0 }
	for col:=0;col<7;col++ {
		pivot:=col;best:=math.Abs(a[col][col])
		for r:=col+1;r<7;r++ { if v:=math.Abs(a[r][col]);v>best { best=v;pivot=r } }
		if best<1e-12||math.IsNaN(best)||math.IsInf(best,0) { return [7]float64{},false }
		if pivot!=col { a[col],a[pivot]=a[pivot],a[col] }
		pv:=a[col][col]
		for j:=col;j<8;j++ { a[col][j]/=pv }
		for r:=0;r<7;r++ {
			if r==col { continue }
			f:=a[r][col];if f==0 { continue }
			for j:=col;j<8;j++ { a[r][j]-=f*a[col][j] }
		}
	}
	var beta [7]float64
	for i:=0;i<7;i++ { beta[i]=a[i][7];if math.IsNaN(beta[i])||math.IsInf(beta[i],0){return [7]float64{},false} }
	return beta,true
}

func wlmLmR47Predict(beta [7]float64,x [6]float64)float64 {
	v:=beta[0];for i:=0;i<6;i++ { v+=beta[i+1]*x[i] };return v
}

func wlmLmR47PolicyPred(p wlmLmR47Prepared, states [3]map[int][6]float64, beta [7]float64, alloc [3]int)float64 {
	total:=0.0
	for ai,arm:=range p.man.arms {
		for i:=1;i<len(arm.budgets);i++ {
			hi:=arm.budgets[i]
			if hi<=alloc[ai] { total+=wlmLmR47Predict(beta,states[ai][hi]) }
		}
	}
	return total
}

func wlmLmR47Actions(m map[string]float64)[][3]int {
	out:=make([][3]int,0,9)
	for _,a:=range wlmLmR27Candidates {
		if a==wlmLmR27Equal { continue }
		if a[0]+a[1]+a[2]!=1744 { m["fixed_budget_mismatch_count"]++ }
		out=append(out,a)
	}
	if len(out)!=9 { m["invalid_row_count"]++ }
	return out
}

func wlmLmR47HistAccuracy(hist []wlmLmR47Prepared, alpha [6]float64, m map[string]float64)float64 {
	states:=make([][3]map[int][6]float64,len(hist))
	ex:=[]wlmLmR47Example{}
	for i,p:=range hist {
		states[i]=wlmLmR47States(p,alpha,m)
		ex=append(ex,wlmLmR47Examples(p,states[i],i)...)
	}
	actions:=wlmLmR47Actions(m)
	correct,total:=0,0
	for hold,p:=range hist {
		beta,ok:=wlmLmR47Fit(ex,hold);if !ok { m["fit_failure_count"]++;continue }
		scores:=map[[3]int]float64{}
		for _,a:=range actions { scores[a]=wlmLmR47PolicyPred(p,states[hold],beta,a) }
		for i:=0;i<len(actions);i++ { for j:=i+1;j<len(actions);j++ {
			a,b:=actions[i],actions[j]
			aa,ab:=wlmLmR27PolicyActual(p.man,a),wlmLmR27PolicyActual(p.man,b)
			if aa==ab { continue }
			total++
			want:=a;if ab>aa { want=b }
			got:=wlmLmR44Choice(a,b,scores[a],scores[b])
			if got==want { correct++ }
		}}
	}
	if total==0 { m["invalid_row_count"]++;return 0 }
	return float64(correct)/float64(total)
}

func wlmLmR47LexAlphaLess(a,b []float64)bool {
	for i:=0;i<len(a)&&i<len(b);i++ { if a[i]!=b[i] { return a[i]<b[i] } }
	return len(a)<len(b)
}

func wlmLmR47SearchSpec(spec wlmLmR47CandidateSpec,hist []wlmLmR47Prepared,m map[string]float64)wlmLmR47CandidateSummary {
	best:=wlmLmR47CandidateSummary{ID:spec.id,LearnedScalars:spec.scalars,SearchScore:math.Inf(-1)}
	vals:=make([]float64,spec.scalars)
	var walk func(int)
	walk=func(g int){
		if g==spec.scalars {
			alpha:=wlmLmR47Alpha(spec,vals)
			acc:=wlmLmR47HistAccuracy(hist,alpha,m)
			pen:=0.002*float64(spec.scalars-1)
			score:=acc-pen
			cp:=append([]float64(nil),vals...)
			if score>best.SearchScore || (score==best.SearchScore&&wlmLmR47LexAlphaLess(cp,best.Alphas)) {
				best.Alphas=cp;best.HistoricalPairwiseAccuracy=acc;best.ComplexityPenalty=pen;best.SearchScore=score
			}
			return
		}
		for _,a:=range wlmLmR47Grid { vals[g]=a;walk(g+1) }
	}
	walk(0)
	best.ResourceEligible=spec.scalars<=1
	return best
}

func wlmLmR47PreparedFrom(name string,ss []wlmLmR27Source,index int,open bool,m map[string]float64)wlmLmR47Prepared {
	man,_:=wlmLmR27Build(name,ss,index,open,m)
	return wlmLmR47Prepared{man:man,sources:ss}
}

func RunWlmLmExternalFutureDataLearnedRawInputFormationR47(
	transferCode,transferStructured,transferProse,
	thirdCode,thirdStructured,thirdProse,
	fourthCode,fourthStructured,fourthProse,
	fifthCode,fifthStructured,fifthProse,
	fifteenthCode,fifteenthStructured,fifteenthProse,
	sixteenthCode,sixteenthStructured,sixteenthProse,
	seventeenthCode,seventeenthStructured,seventeenthProse []byte,
) interface{} {
	budCode:=[]int{0,436,582,727,872};budOther:=[]int{0,291,436,581,726}
	m:=map[string]float64{
		"historical_manifest_count":4,"unseen_manifest_count":3,"candidate_mechanism_count":4,
		"state_dimension":6,"readout_parameter_count":7,"baseline_effective_state_former_scalars":1,
		"selected_effective_state_former_scalars":0,"selected_resource_eligible":0,
		"candidate_pair_count":0,"non_tie_pair_count":0,"exact_tie_count":0,
		"learned_correct_pair_count":0,"baseline_correct_pair_count":0,
		"learned_pairwise_winner_accuracy":0,"baseline_pairwise_winner_accuracy":0,
		"learned_minus_baseline_accuracy":0,"learned_mean_decision_regret":0,"baseline_mean_decision_regret":0,
		"positive_manifest_count":0,"heldout_outcome_use_before_candidate_freeze":0,
		"post_result_candidate_or_threshold_choice_count":0,"capacity_growth_event_count":0,
		"external_model_call_count":0,"tokenizer_use_count":0,"fixed_budget_mismatch_count":0,
		"raw_state_source_mismatch_count":0,"source_identity_mismatch_count":0,"fit_failure_count":0,
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
	for i,in:=range histInputs { hist=append(hist,wlmLmR47PreparedFrom("r47_hist_"+in.name,in.sources,i,true,m)) }

	history:=make([]wlmLmR47CandidateSummary,0,4)
	for _,spec:=range wlmLmR47Specs { history=append(history,wlmLmR47SearchSpec(spec,hist,m)) }
	best:=history[0]
	for _,c:=range history[1:] {
		if c.SearchScore>best.SearchScore || (c.SearchScore==best.SearchScore&&c.ID<best.ID) { best=c }
	}
	var bestSpec wlmLmR47CandidateSpec
	for _,s:=range wlmLmR47Specs { if s.id==best.ID { bestSpec=s } }
	bestAlpha:=wlmLmR47Alpha(bestSpec,best.Alphas)
	m["selected_effective_state_former_scalars"]=float64(best.LearnedScalars)
	if best.ResourceEligible { m["selected_resource_eligible"]=1 }

	bestStates:=make([][3]map[int][6]float64,len(hist));allEx:=[]wlmLmR47Example{}
	mans:=make([]wlmLmR27Manifest,0,len(hist))
	for i,p:=range hist { bestStates[i]=wlmLmR47States(p,bestAlpha,m);allEx=append(allEx,wlmLmR47Examples(p,bestStates[i],i)...);mans=append(mans,p.man) }
	learnedBeta,ok:=wlmLmR47Fit(allEx,-1);if !ok { m["fit_failure_count"]++ }
	baselineBeta,okb:=wlmLmR34Fit(wlmLmR46Examples(mans,m));if !okb||len(baselineBeta)!=7 { m["fit_failure_count"]++ }

	unseen:=[]wlmLmR35ManifestInput{
		{"fifteenth",[]wlmLmR27Source{mk("code",fifteenthCode,"d9d63e15e7bbd585765661f4864db73d4c32503a7c7d286350fc00b5548f5b49",41453,budCode),mk("structured",fifteenthStructured,"667d1b6953fb1fc6b452a312bbb2fa7658dff68dab7bc41d60769771b1768252",14365,budOther),mk("technical_prose",fifteenthProse,"a493b2f056a0cbfee8dae55551cc91aef8661e19b0aaa20a04c46493dc30e4bb",1454,budOther)}},
		{"sixteenth",[]wlmLmR27Source{mk("code",sixteenthCode,"2d07e4331a31d7f2804ffab378a577a2aeb8f3c79ae64884845644c00b93485e",41453,budCode),mk("structured",sixteenthStructured,"5bb9dd4fd3dcc51b77340d819ab80a04f06621112f6696d8bbc9874d8a6acfb7",14365,budOther),mk("technical_prose",sixteenthProse,"5e8a418aa51278d2e6a029786ca59dd5ca24da104093d9959c00bbb3242bd450",1454,budOther)}},
		{"seventeenth",[]wlmLmR27Source{mk("code",seventeenthCode,"417d03c0172e12eab882d82093ff1369a0f60b7123a1802038f30e1329ba28bb",41453,budCode),mk("structured",seventeenthStructured,"b6561ecbbbf6683e8ab7307e8180612b55d1d8bb38636d81c1042ac52a088d81",14365,budOther),mk("technical_prose",seventeenthProse,"855c28784bbccb8ffc659055a6ed590e2a4ad71498bbce4fe902018b5f53deb2",1454,budOther)}},
	}
	actions:=wlmLmR47Actions(m);cases:=[]wlmLmR47PairCase{}
	lReg,bReg:=0.0,0.0
	for ui,in:=range unseen {
		pred:=wlmLmR47PreparedFrom("r47_pred_"+in.name,in.sources,4+ui,false,m)
		st:=wlmLmR47States(pred,bestAlpha,m)
		type sc struct{l,b float64};scores:=map[[3]int]sc{}
		for _,a:=range actions { scores[a]=sc{l:wlmLmR47PolicyPred(pred,st,learnedBeta,a),b:wlmLmR46PolicyPred(pred.man,baselineBeta,a,m)} }
		actual:=wlmLmR47PreparedFrom("r47_actual_"+in.name,in.sources,4+ui,true,m)
		localN,localL,localB:=0,0,0
		for i:=0;i<len(actions);i++ { for j:=i+1;j<len(actions);j++ {
			a,b:=actions[i],actions[j];aa,ab:=wlmLmR27PolicyActual(actual.man,a),wlmLmR27PolicyActual(actual.man,b)
			m["candidate_pair_count"]++
			lc:=wlmLmR44Choice(a,b,scores[a].l,scores[b].l);bc:=wlmLmR44Choice(a,b,scores[a].b,scores[b].b)
			row:=wlmLmR47PairCase{Manifest:in.name,ActionA:a,ActionB:b,LearnedChoice:lc,BaselineChoice:bc,ActualDelta:aa-ab}
			if aa==ab { row.Tie=true;m["exact_tie_count"]++ } else {
				m["non_tie_pair_count"]++;localN++;want:=a;if ab>aa { want=b };row.ActualChoice=want;gap:=math.Abs(aa-ab)
				if lc==want { m["learned_correct_pair_count"]++;localL++ } else { lReg+=gap }
				if bc==want { m["baseline_correct_pair_count"]++;localB++ } else { bReg+=gap }
			}
			cases=append(cases,row)
		}}
		if localN==0 { m["invalid_row_count"]++ } else {
			la:=float64(localL)/float64(localN);ba:=float64(localB)/float64(localN)
			m[in.name+"_learned_pairwise_winner_accuracy"]=la;m[in.name+"_baseline_pairwise_winner_accuracy"]=ba
			if la>ba { m["positive_manifest_count"]++ }
		}
	}
	n:=m["non_tie_pair_count"];if n<=0 { m["invalid_row_count"]++ } else {
		m["learned_pairwise_winner_accuracy"]=m["learned_correct_pair_count"]/n;m["baseline_pairwise_winner_accuracy"]=m["baseline_correct_pair_count"]/n
		m["learned_minus_baseline_accuracy"]=m["learned_pairwise_winner_accuracy"]-m["baseline_pairwise_winner_accuracy"]
		m["learned_mean_decision_regret"]=lReg/n;m["baseline_mean_decision_regret"]=bReg/n
	}
	for k,v:=range m { if strings.HasSuffix(k,"_source_identity_mismatch_count")&&k!="source_identity_mismatch_count" { m["source_identity_mismatch_count"]+=v } }
	if m["candidate_mechanism_count"]!=4||len(history)!=4||m["state_dimension"]!=6||m["readout_parameter_count"]!=7||m["candidate_pair_count"]!=108 { m["invalid_row_count"]++ }
	if m["fixed_budget_mismatch_count"]!=0||m["fit_failure_count"]!=0||m["raw_state_source_mismatch_count"]!=0||m["source_identity_mismatch_count"]!=0 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmR47Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-RAW-INPUT-FORMATION-R47",SelectedCandidate:best.ID,SelectedAlphas:best.Alphas,CandidateHistory:history,Metrics:m,Cases:cases}
}
