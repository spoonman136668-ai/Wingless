package unitary

import "math"

type wlmLmExternalReasoningTransferCapacityFeasibilityR2Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmTransferCapacityBudgetR2(a,b []byte)([]byte,[]byte){
	const budget=15819
	total:=len(a)+len(b)
	if total<budget || total==0 { return nil,nil }
	na:=budget*len(a)/total
	nb:=budget-na
	if na>len(a) { na=len(a); nb=budget-na }
	if nb>len(b) { nb=len(b); na=budget-nb }
	if na<0 || nb<0 || na>len(a) || nb>len(b) || na+nb!=budget { return nil,nil }
	return a[:na],b[:nb]
}

func RunWlmLmExternalReasoningTransferCapacityFeasibilityR2(code,structured,prose []byte) interface{} {
	sources:=[]wlmLmExternalMotifRelationDistributionalR1Source{
		{domain:"code",data:code,sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
		{domain:"structured",data:structured,sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
		{domain:"technical_prose",data:prose,sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":3,
		"total_source_bytes":0,
		"arm_count":3,
		"training_byte_budget_per_arm":15819,
		"target_byte_use_count":0,
		"target_label_use_count":0,
		"candidate_capacity_count":4,
		"minimum_eligible_motif_count_across_arms":math.Inf(1),
		"selected_common_capacity":0,
		"capacity_selection_post_target_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"counter_overflow_count":0,
		"invalid_row_count":0,
	}
	for _,s:=range sources {
		m["total_source_bytes"]+=float64(len(s.data))
		if len(s.data)!=s.bytes || wlmLmExternalRawRepPredR1SHA256(s.data)!=s.sha256 { m["source_identity_mismatch_count"]++ }
	}
	for targetIdx,target:=range sources {
		idx:=make([]int,0,2)
		for i:=range sources { if i!=targetIdx { idx=append(idx,i) } }
		a,b:=wlmLmTransferCapacityBudgetR2(sources[idx[0]].data,sources[idx[1]].data)
		if len(a)+len(b)!=15819 {
			m["invalid_row_count"]++
			continue
		}
		local:=map[string]float64{"counter_overflow_rows":0}
		_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel([][]byte{a,b},local)
		m["counter_overflow_count"]+=local["counter_overflow_rows"]
		n:=float64(len(selected))
		m["arm_"+target.domain+"_eligible_motif_count"]=n
		m["arm_"+target.domain+"_training_bytes"]=float64(len(a)+len(b))
		if n<m["minimum_eligible_motif_count_across_arms"] { m["minimum_eligible_motif_count_across_arms"]=n }
		for _,capv:=range []int{128,256,384,512} {
			key:="arm_"+target.domain+"_supports_capacity_"+string(rune('0'+capv/128))
			if len(selected)>=capv { m[key]=1 } else { m[key]=0 }
		}
	}
	minCount:=m["minimum_eligible_motif_count_across_arms"]
	for _,capv:=range []int{128,256,384,512} {
		if minCount>=float64(capv) { m["selected_common_capacity"]=float64(capv) }
	}
	if math.IsInf(m["minimum_eligible_motif_count_across_arms"],0) {
		m["minimum_eligible_motif_count_across_arms"]=0
		m["invalid_row_count"]++
	}
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmExternalReasoningTransferCapacityFeasibilityR2Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-REASONING-TRANSFER-CAPACITY-FEASIBILITY-R2",
		Metrics:m,
	}
}
