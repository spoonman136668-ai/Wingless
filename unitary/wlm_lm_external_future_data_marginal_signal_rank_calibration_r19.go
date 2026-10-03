package unitary

import "math"

type wlmLmFutureDataMarginalSignalRankCalibrationR19Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmR19Interval struct {
	signal float64
	eff float64
	positive bool
}

func RunWlmLmExternalFutureDataMarginalSignalRankCalibrationR19(code, structured, prose []byte) interface{} {
	type src struct{ d string; b []byte; h string; n int; budgets []int }
	ss:=[]src{
		{"code",code,"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",41453,[]int{0,436,582,727,872}},
		{"structured",structured,"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",14365,[]int{0,291,436,581,726}},
		{"technical_prose",prose,"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",1454,[]int{0,291,436,581,726}},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,"arm_count":3,
		"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,
		"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,
		"budget_point_count_total":15,"interval_count_total":12,
		"signal_weight_baseline_argmax_change":1,"signal_weight_motif_churn":3,"signal_weight_retained_best_change":0,
		"signal_evaluation_byte_use_count":0,"signal_evaluation_label_use_count":0,"signal_post_result_choice_count":0,
		"positive_interval_count":0,"zero_interval_count":0,"negative_interval_count":0,
		"positive_vs_zero_auc":0,"marginal_efficiency_pairwise_concordance":0,
		"mean_positive_rank_signal":0,"mean_zero_rank_signal":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,s:=range ss {
		m["total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h { m["source_identity_mismatch_count"]++ }
	}
	intervals:=[]wlmLmR19Interval{}
	posSignalSum:=0.0;zeroSignalSum:=0.0
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
		signals:=make([]float64,len(t.budgets)-1)
		for i:=1;i<len(t.budgets);i++ {
			lo,hi:=t.budgets[i-1],t.budgets[i]
			bc,en,ex,best:=wlmLmR18StateDelta(&bases[i-1],&bases[i],sels[i-1],sels[i])
			_ = best
			width:=hi-lo
			signal:=float64(bc+3*(en+ex))/float64(width)
			signals[i-1]=signal
			k:="arm_"+t.d+"_interval_"+itoaR17(lo)+"_"+itoaR17(hi)
			m[k+"_baseline_argmax_change_count"]=float64(bc)
			m[k+"_motif_entry_count"]=float64(en)
			m[k+"_motif_exit_count"]=float64(ex)
			m[k+"_retained_motif_best_change_count"]=float64(best)
			m[k+"_rank_signal"]=signal
		}
		hits:=make([]int,len(t.budgets))
		for bi:=range t.budgets {
			hits[bi]=wlmLmR10Hits(eval,&bases[bi],sels[bi])
			m["arm_"+t.d+"_budget_"+itoaR17(t.budgets[bi])+"_exact_hits"]=float64(hits[bi])
		}
		for i:=1;i<len(t.budgets);i++ {
			lo,hi:=t.budgets[i-1],t.budgets[i]
			width:=hi-lo
			gain:=hits[i]-hits[i-1]
			eff:=float64(gain)/float64(width)
			k:="arm_"+t.d+"_interval_"+itoaR17(lo)+"_"+itoaR17(hi)
			m[k+"_gain"]=float64(gain);m[k+"_marginal_efficiency"]=eff
			iv:=wlmLmR19Interval{signal:signals[i-1],eff:eff,positive:gain>0}
			intervals=append(intervals,iv)
			if gain>0 { m["positive_interval_count"]++;posSignalSum+=iv.signal
			} else if gain==0 { m["zero_interval_count"]++;zeroSignalSum+=iv.signal
			} else { m["negative_interval_count"]++ }
		}
	}
	if m["positive_interval_count"]>0 { m["mean_positive_rank_signal"]=posSignalSum/m["positive_interval_count"] }
	if m["zero_interval_count"]>0 { m["mean_zero_rank_signal"]=zeroSignalSum/m["zero_interval_count"] }
	pos:=[]wlmLmR19Interval{};zero:=[]wlmLmR19Interval{}
	for _,iv:=range intervals {
		if iv.positive { pos=append(pos,iv) } else if iv.eff==0 { zero=append(zero,iv) }
	}
	if len(pos)==0||len(zero)==0 { m["invalid_row_count"]++ } else {
		points:=0.0;total:=0.0
		for _,p:=range pos { for _,z:=range zero {
			total++
			if p.signal>z.signal { points++ } else if p.signal==z.signal { points+=0.5 }
		}}
		m["positive_vs_zero_auc"]=points/total
	}
	concordant:=0.0;comparable:=0.0
	for i:=0;i<len(intervals);i++ { for j:=i+1;j<len(intervals);j++ {
		a,b:=intervals[i],intervals[j]
		if a.eff==b.eff { continue }
		comparable++
		if a.signal==b.signal { concordant+=0.5
		} else if (a.signal>b.signal)==(a.eff>b.eff) { concordant++ }
	}}
	if comparable==0 { m["invalid_row_count"]++ } else { m["marginal_efficiency_pairwise_concordance"]=concordant/comparable }
	if len(intervals)!=12 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmFutureDataMarginalSignalRankCalibrationR19Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-MARGINAL-SIGNAL-RANK-CALIBRATION-R19",
		m,
	}
}
