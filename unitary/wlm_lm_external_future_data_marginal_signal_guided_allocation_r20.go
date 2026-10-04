package unitary

import "math"

type wlmLmFutureDataMarginalSignalGuidedAllocationR20Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmR20Arm struct {
	d string
	eval []byte
	budgets []int
	bases map[int][256][256]uint32
	sels map[int]map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif
	pred map[int]float64
}

func RunWlmLmExternalFutureDataMarginalSignalGuidedAllocationR20(code,structured,prose []byte) interface{} {
	type src struct{d string;b []byte;h string;n int;budgets []int}
	ss:=[]src{
		{"code",code,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff",41453,[]int{0,436,582,727,872}},
		{"structured",structured,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56",14365,[]int{0,291,436,581,726}},
		{"technical_prose",prose,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017",1454,[]int{0,291,436,581,726}},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"third_manifest_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,
		"arm_count":3,"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,
		"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,
		"signal_weight_baseline_argmax_change":1,"signal_weight_motif_churn":3,"signal_weight_retained_best_change":0,
		"candidate_allocation_count":10,"total_adaptation_budget_equal":1744,"total_adaptation_budget_selected":0,
		"selection_evaluation_byte_use_count":0,"selection_evaluation_label_use_count":0,"allocation_post_result_choice_count":0,
		"selected_code_budget":0,"selected_structured_budget":0,"selected_technical_prose_budget":0,
		"selected_predicted_value":0,"selected_differs_equal":0,
		"packet0_aggregate_exact_hits":0,"equal_aggregate_exact_hits":0,"selected_aggregate_exact_hits":0,
		"selected_aggregate_exact_hit_delta_vs_equal":0,"selected_domain_below_packet0_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,s:=range ss {
		m["total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h {m["source_identity_mismatch_count"]++;m["third_manifest_identity_mismatch_count"]++}
	}
	arms:=make([]wlmLmR20Arm,0,3)
	for ti,t:=range ss {
		idx:=[]int{}
		for i:=range ss {if i!=ti {idx=append(idx,i)}}
		a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b)
		if len(a)+len(b)!=15819||len(t.b)<1454 {m["invalid_row_count"]++;continue}
		eval,res:=t.b[:582],t.b[582:1454]
		arm:=wlmLmR20Arm{d:t.d,eval:eval,budgets:t.budgets,bases:map[int][256][256]uint32{},sels:map[int]map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif{},pred:map[int]float64{0:0}}
		for _,bud:=range t.budgets {
			train:=[][]byte{a,b};if bud>0 {train=append(train,res[:bud])}
			base,sel:=wlmLmR10Train256(train,m);if len(sel)!=256 {m["invalid_row_count"]++}
			arm.bases[bud]=base;arm.sels[bud]=sel
		}
		cum:=0.0
		for i:=1;i<len(t.budgets);i++ {
			lo,hi:=t.budgets[i-1],t.budgets[i]
			lb,hb:=arm.bases[lo],arm.bases[hi]
			bc,en,ex,best:=wlmLmR18StateDelta(&lb,&hb,arm.sels[lo],arm.sels[hi]);_ = best
			width:=hi-lo;signal:=float64(bc+3*(en+ex))/float64(width)
			cum+=signal*float64(width);arm.pred[hi]=cum
			k:="arm_"+t.d+"_interval_"+itoaR17(lo)+"_"+itoaR17(hi)
			m[k+"_baseline_argmax_change_count"]=float64(bc);m[k+"_motif_entry_count"]=float64(en);m[k+"_motif_exit_count"]=float64(ex);m[k+"_rank_signal"]=signal
		}
		arms=append(arms,arm)
	}
	if len(arms)!=3 {m["invalid_row_count"]++}
	candidates:=[][3]int{{582,436,726},{582,581,581},{582,726,436},{727,291,726},{727,436,581},{727,581,436},{727,726,291},{872,291,581},{872,436,436},{872,581,291}}
	equal:=[3]int{582,581,581}
	bestScore:=math.Inf(-1);best:=[3]int{}
	for _,c:=range candidates {
		if len(arms)!=3 {break}
		score:=arms[0].pred[c[0]]+arms[1].pred[c[1]]+arms[2].pred[c[2]]
		if score>bestScore {bestScore=score;best=c}
	}
	if math.IsInf(bestScore,-1) {m["invalid_row_count"]++} else {
		m["selected_code_budget"]=float64(best[0]);m["selected_structured_budget"]=float64(best[1]);m["selected_technical_prose_budget"]=float64(best[2])
		m["selected_predicted_value"]=bestScore;m["total_adaptation_budget_selected"]=float64(best[0]+best[1]+best[2])
		if best!=equal {m["selected_differs_equal"]=1}
	}
	if m["total_adaptation_budget_selected"]!=1744 {m["invalid_row_count"]++}
	if len(arms)==3 {
		for i,arm:=range arms {
			p0b:=arm.bases[0];eqb:=arm.bases[equal[i]];sb:=arm.bases[best[i]]
			p0:=wlmLmR10Hits(arm.eval,&p0b,arm.sels[0])
			eh:=wlmLmR10Hits(arm.eval,&eqb,arm.sels[equal[i]])
			sh:=wlmLmR10Hits(arm.eval,&sb,arm.sels[best[i]])
			p:="arm_"+arm.d+"_"
			m[p+"packet0_exact_hits"]=float64(p0);m[p+"equal_exact_hits"]=float64(eh);m[p+"selected_exact_hits"]=float64(sh);m[p+"selected_delta_vs_equal"]=float64(sh-eh)
			m["packet0_aggregate_exact_hits"]+=float64(p0);m["equal_aggregate_exact_hits"]+=float64(eh);m["selected_aggregate_exact_hits"]+=float64(sh)
			if sh<p0 {m["selected_domain_below_packet0_count"]++}
		}
		m["selected_aggregate_exact_hit_delta_vs_equal"]=m["selected_aggregate_exact_hits"]-m["equal_aggregate_exact_hits"]
	}
	for _,v:=range m {if math.IsNaN(v)||math.IsInf(v,0) {m["invalid_row_count"]++}}
	return wlmLmFutureDataMarginalSignalGuidedAllocationR20Result{"wingless.research-scientific-result.v1","WLM-LM-EXTERNAL-FUTURE-DATA-MARGINAL-SIGNAL-GUIDED-ALLOCATION-R20",m}
}
