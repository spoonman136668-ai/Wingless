package unitary

import "math"

type wlmLmFutureDataMarginalSignalMotifChurnControlR25Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalFutureDataMarginalSignalMotifChurnControlR25(code,structured,prose []byte) interface{} {
	type src struct{d string;b []byte;h string;n int;budgets []int}
	ss:=[]src{
		{"code",code,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1",41453,[]int{0,436,582,727,872}},
		{"structured",structured,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933",14365,[]int{0,291,436,581,726}},
		{"technical_prose",prose,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222",1454,[]int{0,291,436,581,726}},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"fourth_manifest_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,
		"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,
		"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,
		"churn_interval_count":0,"productive_churn_interval_count":0,"zero_gain_churn_interval_count":0,"negative_gain_churn_interval_count":0,
		"productive_control_signal_sum":0,"zero_gain_control_signal_sum":0,"productive_control_signal_mean":0,"zero_gain_control_signal_mean":0,
		"productive_vs_zero_pairwise_auc":0,
		"control_evaluation_byte_use_count":0,"control_evaluation_label_use_count":0,"control_post_result_choice_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,s:=range ss {
		m["total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h {
			m["source_identity_mismatch_count"]++;m["fourth_manifest_identity_mismatch_count"]++
		}
	}
	type obs struct{control,gain float64}
	productive:=[]obs{};zero:=[]obs{}
	for ti,t:=range ss {
		idx:=[]int{}
		for i:=range ss {if i!=ti {idx=append(idx,i)}}
		a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b)
		if len(a)+len(b)!=15819||len(t.b)<1454 {m["invalid_row_count"]++;continue}
		eval,res:=t.b[:582],t.b[582:1454]
		bases:=map[int][256][256]uint32{}
		sels:=map[int]map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif{}
		for _,bud:=range t.budgets {
			train:=[][]byte{a,b};if bud>0 {train=append(train,res[:bud])}
			base,sel:=wlmLmR10Train256(train,m);if len(sel)!=256 {m["invalid_row_count"]++}
			bases[bud]=base;sels[bud]=sel
		}
		for i:=1;i<len(t.budgets);i++ {
			lo,hi:=t.budgets[i-1],t.budgets[i]
			lb,hb:=bases[lo],bases[hi];ls,hs:=sels[lo],sels[hi]
			entries,exits:=0,0;activeEntries,activeExits:=0,0
			for k,motif:=range hs {
				if _,ok:=ls[k];ok {continue}
				entries++
				if motif.best!=wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(&hb,k[3]) {activeEntries++}
			}
			for k,motif:=range ls {
				if _,ok:=hs[k];ok {continue}
				exits++
				if motif.best!=wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(&lb,k[3]) {activeExits++}
			}
			churn:=entries+exits
			if churn==0 {continue}
			width:=hi-lo;active:=activeEntries+activeExits
			control:=float64(active)/float64(width)
			loHits:=wlmLmR10Hits(eval,&lb,ls);hiHits:=wlmLmR10Hits(eval,&hb,hs)
			gain:=float64(hiHits-loHits)
			k:="arm_"+t.d+"_interval_"+itoaR17(lo)+"_"+itoaR17(hi)
			m[k+"_motif_entry_count"]=float64(entries)
			m[k+"_motif_exit_count"]=float64(exits)
			m[k+"_prediction_active_entry_count"]=float64(activeEntries)
			m[k+"_prediction_active_exit_count"]=float64(activeExits)
			m[k+"_prediction_active_churn_count"]=float64(active)
			m[k+"_control_signal"]=control
			m[k+"_actual_gain"]=gain
			m["churn_interval_count"]++
			if gain>0 {
				m["productive_churn_interval_count"]++
				m["productive_control_signal_sum"]+=control
				productive=append(productive,obs{control,gain})
			} else if gain==0 {
				m["zero_gain_churn_interval_count"]++
				m["zero_gain_control_signal_sum"]+=control
				zero=append(zero,obs{control,gain})
			} else {
				m["negative_gain_churn_interval_count"]++
			}
		}
	}
	if len(productive)>0 {m["productive_control_signal_mean"]=m["productive_control_signal_sum"]/float64(len(productive))}
	if len(zero)>0 {m["zero_gain_control_signal_mean"]=m["zero_gain_control_signal_sum"]/float64(len(zero))}
	if len(productive)>0&&len(zero)>0 {
		score:=0.0
		for _,p:=range productive {
			for _,z:=range zero {
				if p.control>z.control {score+=1} else if p.control==z.control {score+=0.5}
			}
		}
		m["productive_vs_zero_pairwise_auc"]=score/float64(len(productive)*len(zero))
	}
	for _,v:=range m {if math.IsNaN(v)||math.IsInf(v,0) {m["invalid_row_count"]++}}
	return wlmLmFutureDataMarginalSignalMotifChurnControlR25Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-MARGINAL-SIGNAL-MOTIF-CHURN-CONTROL-R25",
		m,
	}
}
