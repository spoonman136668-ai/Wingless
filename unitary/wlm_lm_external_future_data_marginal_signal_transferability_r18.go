package unitary

import "math"

type wlmLmFutureDataMarginalSignalTransferabilityR18Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmR18StateDelta(
	loBase, hiBase *[256][256]uint32,
	loSel, hiSel map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif,
) (baselineChanges, entries, exits, bestChanges int) {
	for i:=0;i<256;i++ {
		if wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(loBase,byte(i)) != wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(hiBase,byte(i)) {
			baselineChanges++
		}
	}
	for k,lo:=range loSel {
		hi,ok:=hiSel[k]
		if !ok { exits++; continue }
		if lo.best!=hi.best { bestChanges++ }
	}
	for k:=range hiSel {
		if _,ok:=loSel[k];!ok { entries++ }
	}
	return
}

func RunWlmLmExternalFutureDataMarginalSignalTransferabilityR18(code, structured, prose []byte) interface{} {
	type src struct{ d string; b []byte; h string; n int; budgets []int }
	ss:=[]src{
		{"code",code,"66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45",41453,[]int{0,436,582,727,872}},
		{"structured",structured,"95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932",14365,[]int{0,291,436,581,726}},
		{"technical_prose",prose,"48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b",1454,[]int{0,291,436,581,726}},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"transfer_manifest_identity_mismatch_count":0,
		"source_count":3,"total_source_bytes":0,"arm_count":3,
		"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,
		"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,
		"budget_point_count_total":15,"interval_count_total":12,
		"signal_evaluation_byte_use_count":0,"signal_evaluation_label_use_count":0,"signal_post_result_choice_count":0,
		"positive_gain_interval_count":0,"zero_gain_interval_count":0,"negative_gain_interval_count":0,
		"positive_gain_with_zero_state_delta_count":0,"zero_gain_with_zero_state_delta_count":0,
		"positive_gain_with_positive_state_delta_count":0,"zero_gain_with_positive_state_delta_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,s:=range ss {
		m["total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h { m["source_identity_mismatch_count"]++ }
	}
	for ti,t:=range ss {
		idx:=[]int{}
		for i:=range ss { if i!=ti { idx=append(idx,i) } }
		a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b)
		if len(a)+len(b)!=15819||len(t.b)<1454 { m["invalid_row_count"]++; continue }
		eval,res:=t.b[:582],t.b[582:1454]
		bases:=make([][256][256]uint32,len(t.budgets))
		sels:=make([]map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif,len(t.budgets))
		for bi,bud:=range t.budgets {
			train:=[][]byte{a,b}
			if bud>0 { train=append(train,res[:bud]) }
			base,sel:=wlmLmR10Train256(train,m)
			if len(sel)!=256 { m["invalid_row_count"]++ }
			bases[bi]=base;sels[bi]=sel
		}
		for i:=1;i<len(t.budgets);i++ {
			lo,hi:=t.budgets[i-1],t.budgets[i]
			bc,en,ex,best:=wlmLmR18StateDelta(&bases[i-1],&bases[i],sels[i-1],sels[i])
			total:=bc+en+ex+best
			k:="arm_"+t.d+"_interval_"+itoaR17(lo)+"_"+itoaR17(hi)
			m[k+"_baseline_argmax_change_count"]=float64(bc)
			m[k+"_motif_entry_count"]=float64(en)
			m[k+"_motif_exit_count"]=float64(ex)
			m[k+"_retained_motif_best_change_count"]=float64(best)
			m[k+"_state_delta"]=float64(total)
		}
		hits:=make([]int,len(t.budgets))
		for bi:=range t.budgets {
			hits[bi]=wlmLmR10Hits(eval,&bases[bi],sels[bi])
			m["arm_"+t.d+"_budget_"+itoaR17(t.budgets[bi])+"_exact_hits"]=float64(hits[bi])
		}
		for i:=1;i<len(t.budgets);i++ {
			lo,hi:=t.budgets[i-1],t.budgets[i]
			k:="arm_"+t.d+"_interval_"+itoaR17(lo)+"_"+itoaR17(hi)
			gain:=hits[i]-hits[i-1]
			m[k+"_gain"]=float64(gain)
			delta:=m[k+"_state_delta"]
			if gain>0 {
				m["positive_gain_interval_count"]++
				if delta==0 { m["positive_gain_with_zero_state_delta_count"]++ } else { m["positive_gain_with_positive_state_delta_count"]++ }
			} else if gain==0 {
				m["zero_gain_interval_count"]++
				if delta==0 { m["zero_gain_with_zero_state_delta_count"]++ } else { m["zero_gain_with_positive_state_delta_count"]++ }
			} else {
				m["negative_gain_interval_count"]++
			}
		}
	}
	if m["positive_gain_interval_count"]+m["zero_gain_interval_count"]+m["negative_gain_interval_count"]!=12 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmFutureDataMarginalSignalTransferabilityR18Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-MARGINAL-SIGNAL-TRANSFERABILITY-R18",
		m,
	}
}
