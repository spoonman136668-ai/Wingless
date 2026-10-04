package unitary

import "math"

type wlmLmFutureDataMotifChurnProspectiveR27Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalFutureDataMotifChurnFactorialProspectiveValidationR27(code,structured,prose []byte) interface{} {
	type src struct{d string;b []byte;h string;n int;budgets []int}
	ss:=[]src{
		{"code",code,"1c93e8c1d4554e0e1dd6d2a9210ad97b9882bb1a1251ead1215891e7119b70d4",41453,[]int{0,436,582,727,872}},
		{"structured",structured,"9ddfc53d69bd42f120a6d8b64eb9178dc28e53d0332293e8bf294a3f2b6386cc",14365,[]int{0,291,436,581,726}},
		{"technical_prose",prose,"527e21110a7f84a1939ccf6063fdfeb77905d9ec180e5fde18488d21b90c4f49",1454,[]int{0,291,436,581,726}},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"fifth_manifest_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,
		"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,
		"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,
		"churn_interval_count":0,"productive_churn_interval_count":0,"zero_gain_churn_interval_count":0,"negative_gain_churn_interval_count":0,
		"productive_signal_sum":0,"zero_gain_signal_sum":0,"productive_signal_mean":0,"zero_gain_signal_mean":0,
		"productive_vs_zero_pairwise_auc":0,
		"control_evaluation_byte_use_count":0,"control_evaluation_label_use_count":0,"post_result_signal_change_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,s:=range ss{
		m["total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h{
			m["source_identity_mismatch_count"]++;m["fifth_manifest_identity_mismatch_count"]++
		}
	}
	type rec struct{signal,gain float64}
	recs:=[]rec{}
	for ti,t:=range ss{
		idx:=[]int{}
		for i:=range ss{if i!=ti{idx=append(idx,i)}}
		a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b)
		if len(a)+len(b)!=15819||len(t.b)<1454{m["invalid_row_count"]++;continue}
		eval,res:=t.b[:582],t.b[582:1454]
		bases:=map[int][256][256]uint32{}
		sels:=map[int]map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif{}
		for _,bud:=range t.budgets{
			train:=[][]byte{a,b};if bud>0{train=append(train,res[:bud])}
			base,sel:=wlmLmR10Train256(train,m);if len(sel)!=256{m["invalid_row_count"]++}
			bases[bud]=base;sels[bud]=sel
		}
		for i:=1;i<len(t.budgets);i++{
			lo,hi:=t.budgets[i-1],t.budgets[i]
			ls,hs:=sels[lo],sels[hi]
			entries,exits:=0,0
			for k:=range hs{if _,ok:=ls[k];!ok{entries++}}
			for k:=range ls{if _,ok:=hs[k];!ok{exits++}}
			churn:=entries+exits
			if churn==0{continue}
			signal:=float64(churn)/float64(hi-lo)
			lb,hb:=bases[lo],bases[hi]
			loHits:=wlmLmR10Hits(eval,&lb,ls);hiHits:=wlmLmR10Hits(eval,&hb,hs)
			gain:=float64(hiHits-loHits)
			k:="arm_"+t.d+"_interval_"+itoaR17(lo)+"_"+itoaR17(hi)
			m[k+"_motif_entry_count"]=float64(entries)
			m[k+"_motif_exit_count"]=float64(exits)
			m[k+"_raw_churn_signal"]=signal
			m[k+"_actual_gain"]=gain
			m["churn_interval_count"]++
			if gain>0{
				m["productive_churn_interval_count"]++;m["productive_signal_sum"]+=signal
			}else if gain==0{
				m["zero_gain_churn_interval_count"]++;m["zero_gain_signal_sum"]+=signal
			}else{
				m["negative_gain_churn_interval_count"]++
			}
			recs=append(recs,rec{signal,gain})
		}
	}
	if m["productive_churn_interval_count"]>0{
		m["productive_signal_mean"]=m["productive_signal_sum"]/m["productive_churn_interval_count"]
	}
	if m["zero_gain_churn_interval_count"]>0{
		m["zero_gain_signal_mean"]=m["zero_gain_signal_sum"]/m["zero_gain_churn_interval_count"]
	}
	prod:=[]float64{};zero:=[]float64{}
	for _,x:=range recs{
		if x.gain>0{prod=append(prod,x.signal)}else if x.gain==0{zero=append(zero,x.signal)}
	}
	if len(prod)>0&&len(zero)>0{
		s:=0.0
		for _,p:=range prod{for _,z:=range zero{if p>z{s+=1}else if p==z{s+=0.5}}}
		m["productive_vs_zero_pairwise_auc"]=s/float64(len(prod)*len(zero))
	}
	for _,v:=range m{if math.IsNaN(v)||math.IsInf(v,0){m["invalid_row_count"]++}}
	return wlmLmFutureDataMotifChurnProspectiveR27Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-MOTIF-CHURN-FACTORIAL-PROSPECTIVE-VALIDATION-R27",
		m,
	}
}
