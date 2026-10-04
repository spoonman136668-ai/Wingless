package unitary

import "math"

type wlmLmFutureDataMarginalSignalBaselineFalsePositiveControlR22Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmR22Argmax(row *[256]uint32) int {
	best:=0
	for i:=1;i<256;i++ { if row[i]>row[best] { best=i } }
	return best
}

func wlmLmR22TopMargin(row *[256]uint32) uint32 {
	best,second:=uint32(0),uint32(0)
	for i:=0;i<256;i++ {
		v:=row[i]
		if v>best { second=best;best=v
		} else if v>second { second=v }
	}
	return best-second
}

func RunWlmLmExternalFutureDataMarginalSignalBaselineFalsePositiveControlR22(code,structured,prose []byte) interface{} {
	type src struct{d string;b []byte;h string;n int}
	ss:=[]src{
		{"code",code,"50744a9e70d67d62c97f3f434f4f05788b6b8514c6bce46cf7dafbaf49e2abff",41453},
		{"structured",structured,"2a97ba02bc5e479b1738f6f0c3e09318bb5a255350c84de014ddbcebea46af56",14365},
		{"technical_prose",prose,"23c002a1984ed065abfdbafa82100ed54d6bf6276a947676e710a30c75d96017",1454},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"third_manifest_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,
		"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,"target_evaluation_bytes_per_arm":582,
		"target_adaptation_reservoir_bytes_per_arm":872,"analyzed_interval_count":3,
		"control_evaluation_byte_use_count":0,"control_evaluation_label_use_count":0,"control_post_result_choice_count":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,s:=range ss {
		m["total_source_bytes"]+=float64(len(s.b))
		if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h {m["source_identity_mismatch_count"]++;m["third_manifest_identity_mismatch_count"]++}
	}
	type target struct{idx,lo,hi int;name string}
	targets:=[]target{
		{0,582,727,"code_582_727"},
		{2,0,291,"technical_prose_0_291"},
		{2,581,726,"technical_prose_581_726"},
	}
	for _,q:=range targets {
		t:=ss[q.idx];idx:=[]int{}
		for i:=range ss {if i!=q.idx {idx=append(idx,i)}}
		a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b)
		if len(a)+len(b)!=15819||len(t.b)<1454 {m["invalid_row_count"]++;continue}
		eval,res:=t.b[:582],t.b[582:1454]
		trainLo:=[][]byte{a,b};if q.lo>0 {trainLo=append(trainLo,res[:q.lo])}
		trainHi:=[][]byte{a,b};if q.hi>0 {trainHi=append(trainHi,res[:q.hi])}
		baseLo,selLo:=wlmLmR10Train256(trainLo,m);baseHi,selHi:=wlmLmR10Train256(trainHi,m)
		if len(selLo)!=256||len(selHi)!=256 {m["invalid_row_count"]++}
		changed:=0;marginSum:=uint64(0)
		for p:=0;p<256;p++ {
			if wlmLmR22Argmax(&baseLo[p])!=wlmLmR22Argmax(&baseHi[p]) {
				changed++
				marginSum+=uint64(wlmLmR22TopMargin(&baseHi[p]))
			}
		}
		width:=q.hi-q.lo
		control:=float64(marginSum)/float64(width)
		m[q.name+"_baseline_argmax_change_count"]=float64(changed)
		m[q.name+"_upper_margin_sum"]=float64(marginSum)
		m[q.name+"_control_signal"]=control
		loHits:=wlmLmR10Hits(eval,&baseLo,selLo);hiHits:=wlmLmR10Hits(eval,&baseHi,selHi)
		m[q.name+"_actual_gain"]=float64(hiHits-loHits)
	}
	for _,v:=range m {if math.IsNaN(v)||math.IsInf(v,0) {m["invalid_row_count"]++}}
	return wlmLmFutureDataMarginalSignalBaselineFalsePositiveControlR22Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-MARGINAL-SIGNAL-BASELINE-FALSE-POSITIVE-CONTROL-R22",
		m,
	}
}
