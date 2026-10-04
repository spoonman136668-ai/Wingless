package unitary

import "math"

type wlmLmFutureDataMotifChurnControlFactorialR26Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmR26Obs struct {
	scores [8]float64
	gain float64
}

func wlmLmR26CellIndex(gate,support,consistency bool) int {
	i:=0
	if gate {i+=4}
	if support {i+=2}
	if consistency {i++}
	return i
}

func RunWlmLmExternalFutureDataMotifChurnControlFactorialR26(code,structured,prose []byte) interface{} {
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
		"factorial_cell_count":8,"churn_interval_count":0,"productive_interval_count":0,"zero_gain_interval_count":0,"negative_interval_count":0,
		"heldout_feature_selection_count":0,"post_result_cell_install_count":0,"supported_cell_count":0,"mixed_cell_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,s:=range ss {
		m["total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h {
			m["source_identity_mismatch_count"]++;m["fourth_manifest_identity_mismatch_count"]++
		}
	}
	obs:=[]wlmLmR26Obs{}
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
			type motifView struct{motif wlmLmRawRepPredFreshHoldoutR1Motif;disagree bool}
			views:=[]motifView{}
			entries,exits:=0,0
			for k,motif:=range hs {
				if _,ok:=ls[k];ok {continue}
				entries++
				dis:=motif.best!=wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(&hb,k[3])
				views=append(views,motifView{motif,dis})
			}
			for k,motif:=range ls {
				if _,ok:=hs[k];ok {continue}
				exits++
				dis:=motif.best!=wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(&lb,k[3])
				views=append(views,motifView{motif,dis})
			}
			if entries+exits==0 {continue}
			width:=hi-lo
			var scores [8]float64
			for _,v:=range views {
				for _,gate:=range []bool{false,true} {
					for _,support:=range []bool{false,true} {
						for _,consistency:=range []bool{false,true} {
							ci:=wlmLmR26CellIndex(gate,support,consistency)
							if gate&&!v.disagree {continue}
							contrib:=1.0
							if support {contrib*=float64(v.motif.bestCount)}
							if consistency {contrib*=v.motif.consistency}
							scores[ci]+=contrib
						}
					}
				}
			}
			for ci:=0;ci<8;ci++ {scores[ci]/=float64(width)}
			loHits:=wlmLmR10Hits(eval,&lb,ls);hiHits:=wlmLmR10Hits(eval,&hb,hs);gain:=float64(hiHits-loHits)
			k:="arm_"+t.d+"_interval_"+itoaR17(lo)+"_"+itoaR17(hi)
			m[k+"_motif_entry_count"]=float64(entries);m[k+"_motif_exit_count"]=float64(exits);m[k+"_actual_gain"]=gain
			for ci:=0;ci<8;ci++ {m[k+"_cell_"+itoaR17(ci)+"_score"]=scores[ci]}
			m["churn_interval_count"]++
			if gain>0 {m["productive_interval_count"]++} else if gain==0 {m["zero_gain_interval_count"]++} else {m["negative_interval_count"]++}
			obs=append(obs,wlmLmR26Obs{scores:scores,gain:gain})
		}
	}
	for ci:=0;ci<8;ci++ {
		ps:=[]float64{};zs:=[]float64{}
		for _,o:=range obs {
			if o.gain>0 {ps=append(ps,o.scores[ci])} else if o.gain==0 {zs=append(zs,o.scores[ci])}
		}
		pmean,zmean:=0.0,0.0
		for _,v:=range ps {pmean+=v};if len(ps)>0 {pmean/=float64(len(ps))}
		for _,v:=range zs {zmean+=v};if len(zs)>0 {zmean/=float64(len(zs))}
		auc:=0.0
		if len(ps)>0&&len(zs)>0 {
			for _,p:=range ps {for _,z:=range zs {if p>z {auc+=1} else if p==z {auc+=0.5}}}
			auc/=float64(len(ps)*len(zs))
		}
		prefix:="cell_"+itoaR17(ci)
		m[prefix+"_productive_mean"]=pmean;m[prefix+"_zero_gain_mean"]=zmean;m[prefix+"_auc"]=auc
		if auc>=0.75&&pmean>zmean {m["supported_cell_count"]++} else if auc>0.5&&pmean>zmean {m["mixed_cell_count"]++}
	}
	for _,v:=range m {if math.IsNaN(v)||math.IsInf(v,0) {m["invalid_row_count"]++}}
	return wlmLmFutureDataMotifChurnControlFactorialR26Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-MOTIF-CHURN-CONTROL-FACTORIAL-R26",
		m,
	}
}
