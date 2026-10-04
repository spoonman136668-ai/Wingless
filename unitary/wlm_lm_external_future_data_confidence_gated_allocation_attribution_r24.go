package unitary

import "math"

type wlmLmFutureDataConfidenceGatedAllocationAttributionR24Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalFutureDataConfidenceGatedAllocationAttributionR24(code,structured,prose []byte) interface{} {
	type src struct{d string;b []byte;h string;n int}
	ss:=[]src{
		{"code",code,"283073d9f6c0dd868c39a913364bce6744ff1e29c038f6920197c0d33e0c2ac1",41453},
		{"structured",structured,"a46fcfb7d862b03b750b61a5f667d4ac25a064df9ccb746e395db7e864068933",14365},
		{"technical_prose",prose,"5d0c2efd139bd6094098bc893ed746020f03e0860a25f278348f43f47c236222",1454},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"fourth_manifest_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,
		"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,
		"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,
		"baseline_argmax_weight":1,"motif_churn_weight":3,
		"target_interval_count":4,
		"predicted_added_code_contribution":0,"predicted_omitted_structured_contribution":0,
		"actual_added_code_gain":0,"actual_omitted_structured_gain":0,
		"predicted_net":0,"actual_net":0,
		"selection_evaluation_byte_use_count":0,"selection_evaluation_label_use_count":0,"allocation_post_result_choice_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,s:=range ss {
		m["total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h {
			m["source_identity_mismatch_count"]++
			m["fourth_manifest_identity_mismatch_count"]++
		}
	}
	type target struct{idx,lo,hi int;name,side string}
	targets:=[]target{
		{0,582,727,"code_582_727","added"},
		{0,727,872,"code_727_872","added"},
		{1,291,436,"structured_291_436","omitted"},
		{1,436,581,"structured_436_581","omitted"},
	}
	for _,q:=range targets {
		t:=ss[q.idx];idx:=[]int{}
		for i:=range ss {if i!=q.idx {idx=append(idx,i)}}
		a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b)
		if len(a)+len(b)!=15819||len(t.b)<1454 {m["invalid_row_count"]++;continue}
		eval,res:=t.b[:582],t.b[582:1454]
		trainLo:=[][]byte{a,b};if q.lo>0 {trainLo=append(trainLo,res[:q.lo])}
		trainHi:=[][]byte{a,b};if q.hi>0 {trainHi=append(trainHi,res[:q.hi])}
		baseLo,selLo:=wlmLmR10Train256(trainLo,m)
		baseHi,selHi:=wlmLmR10Train256(trainHi,m)
		if len(selLo)!=256||len(selHi)!=256 {m["invalid_row_count"]++}
		bc,en,ex,best:=wlmLmR18StateDelta(&baseLo,&baseHi,selLo,selHi);_ = best
		marginSum:=uint64(0)
		for p:=0;p<256;p++ {
			if wlmLmR22Argmax(&baseLo[p])!=wlmLmR22Argmax(&baseHi[p]) {
				marginSum+=uint64(wlmLmR22TopMargin(&baseHi[p]))
			}
		}
		gatedBC:=0
		if marginSum>0 {gatedBC=bc}
		width:=q.hi-q.lo
		signal:=float64(gatedBC+3*(en+ex))/float64(width)
		pred:=signal*float64(width)
		loHits:=wlmLmR10Hits(eval,&baseLo,selLo)
		hiHits:=wlmLmR10Hits(eval,&baseHi,selHi)
		gain:=float64(hiHits-loHits)
		k:="interval_"+q.name
		m[k+"_baseline_argmax_change_count"]=float64(bc)
		m[k+"_upper_margin_sum"]=float64(marginSum)
		m[k+"_gated_baseline_count"]=float64(gatedBC)
		m[k+"_motif_entry_count"]=float64(en)
		m[k+"_motif_exit_count"]=float64(ex)
		m[k+"_confidence_gated_signal"]=signal
		m[k+"_predicted_contribution"]=pred
		m[k+"_lower_exact_hits"]=float64(loHits)
		m[k+"_upper_exact_hits"]=float64(hiHits)
		m[k+"_actual_gain"]=gain
		if q.side=="added" {
			m["predicted_added_code_contribution"]+=pred
			m["actual_added_code_gain"]+=gain
		} else {
			m["predicted_omitted_structured_contribution"]+=pred
			m["actual_omitted_structured_gain"]+=gain
		}
	}
	m["predicted_net"]=m["predicted_added_code_contribution"]-m["predicted_omitted_structured_contribution"]
	m["actual_net"]=m["actual_added_code_gain"]-m["actual_omitted_structured_gain"]
	for _,v:=range m {if math.IsNaN(v)||math.IsInf(v,0) {m["invalid_row_count"]++}}
	return wlmLmFutureDataConfidenceGatedAllocationAttributionR24Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-CONFIDENCE-GATED-ALLOCATION-ATTRIBUTION-R24",
		m,
	}
}
