package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalReasoningSelfDiagnosisProspectiveR3Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmProspectiveR3Select256(files [][]byte, metrics map[string]float64) map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif {
	_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(files,metrics)
	rows:=make([]wlmLmRawRepPredFreshHoldoutR1Motif,0,len(selected))
	for _,row:=range selected { rows=append(rows,row) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].bestCount!=rows[j].bestCount { return rows[i].bestCount>rows[j].bestCount }
		if rows[i].consistency!=rows[j].consistency { return rows[i].consistency>rows[j].consistency }
		return wlmLmRawRepPredFreshHoldoutR1Less(rows[i].key,rows[j].key)
	})
	if len(rows)>256 { rows=rows[:256] }
	out:=make(map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif,len(rows))
	for _,row:=range rows { out[row.key]=row }
	return out
}

func wlmLmProspectiveR3Median(xs []float64) float64 {
	sort.Float64s(xs)
	if len(xs)==0 { return 0 }
	n:=len(xs)
	if n%2==1 { return xs[n/2] }
	return (xs[n/2-1]+xs[n/2])/2
}

func RunWlmLmExternalReasoningSelfDiagnosisProspectiveR3(code,structured,prose []byte) interface{} {
	m:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":3,
		"total_source_bytes":float64(len(code)+len(structured)+len(prose)),
		"training_byte_budget":15819,
		"target_eval_byte_budget":1454,
		"selected_motif_count":0,
		"reasoning_class_count":8,
		"minimum_reasoning_class_size":256,
		"maximum_reasoning_class_size":0,
		"training_margin_sample_count":0,
		"zero_training_margin_class_count":0,
		"preflight_unevaluable_flag":0,
		"preflight_target_byte_use_count":0,
		"preflight_target_label_use_count":0,
		"preflight_decision_frozen_flag":0,
		"target_query_count":0,
		"unsupported_target_query_count":0,
		"unsupported_target_exact_hit_count":0,
		"maximum_probability_mass_error":0,
		"minimum_repaired_retained_mass":1,
		"relation_capacity_growth_event_count":0,
		"adaptive_readout_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"counter_overflow_count":0,
		"counter_overflow_rows":0,
		"invalid_row_count":0,
	}
	if len(code)!=41453 || wlmLmExternalRawRepPredR1SHA256(code)!="7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6" { m["source_identity_mismatch_count"]++ }
	if len(structured)!=14365 || wlmLmExternalRawRepPredR1SHA256(structured)!="4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192" { m["source_identity_mismatch_count"]++ }
	if len(prose)!=1454 || wlmLmExternalRawRepPredR1SHA256(prose)!="8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3" { m["source_identity_mismatch_count"]++ }
	if len(structured)+len(prose)!=15819 { m["invalid_row_count"]++ }

	train:=[][]byte{structured,prose}
	selected:=wlmLmProspectiveR3Select256(train,m)
	m["counter_overflow_count"]+=m["counter_overflow_rows"]
	delete(m,"counter_overflow_rows")
	m["selected_motif_count"]=float64(len(selected))
	if len(selected)!=256 { m["invalid_row_count"]++ }

	keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index:=make(map[[4]uint8]int,len(keys))
	for i,k:=range keys { index[k]=i }
	streams:=[][]int{
		wlmLmExternalMotifRelationHoldoutR1Decode(structured,index),
		wlmLmExternalMotifRelationHoldoutR1Decode(prose,index),
	}
	var counts [512]uint64
	var rel [512][512]uint32
	var totals [512]uint64
	var uni [512]uint32
	var unitotal uint64
	pairs:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for _,st:=range streams {
		for _,x:=range st {
			if x<0||x>=256 { m["invalid_row_count"]++; continue }
			counts[x]++
		}
		for j:=0;j+1<len(st);j++ {
			a,b:=st[j],st[j+1]
			if rel[a][b]==^uint32(0)||uni[b]==^uint32(0){ m["counter_overflow_count"]++; continue }
			rel[a][b]++; totals[a]++; uni[b]++; unitotal++
		}
		for j:=1;j<len(st);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:st[j-1],curr:st[j]}
			if pairs[p]==^uint32(0){ m["counter_overflow_count"]++; continue }
			pairs[p]++
		}
	}
	if unitotal==0 { m["invalid_row_count"]++ }

	ids:=make([]int,256)
	for i:=range ids { ids[i]=i }
	sort.Slice(ids,func(i,j int)bool{
		if counts[ids[i]]!=counts[ids[j]] { return counts[ids[i]]>counts[ids[j]] }
		return ids[i]<ids[j]
	})
	var classes [512]int
	var classSizes [8]int
	for rank,id:=range ids {
		c:=rank/32
		classes[id]=c
		classSizes[c]++
	}
	for i:=0;i<8;i++ {
		n:=float64(classSizes[i])
		if n<m["minimum_reasoning_class_size"] { m["minimum_reasoning_class_size"]=n }
		if n>m["maximum_reasoning_class_size"] { m["maximum_reasoning_class_size"]=n }
	}

	rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairs))
	for p,c:=range pairs { rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:p,count:c}) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count { return rows[i].count>rows[j].count }
		if rows[i].pair.prev!=rows[j].pair.prev { return rows[i].pair.prev<rows[j].pair.prev }
		return rows[i].pair.curr<rows[j].pair.curr
	})
	n:=511
	if len(rows)<n { n=len(rows) }
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,n)
	for i:=0;i<n;i++ { exactID[rows[i].pair]=i }
	var exact [512][512]uint32
	var exactTotals [512]uint64
	for _,st:=range streams {
		for j:=1;j+1<len(st);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:st[j-1],curr:st[j]}
			id:=511
			if v,ok:=exactID[p];ok { id=v }
			t:=st[j+1]
			if exact[id][t]==^uint32(0){ m["counter_overflow_count"]++; continue }
			exact[id][t]++; exactTotals[id]++
		}
	}
	var eligible [512]bool
	var marginsByRow [512]float64
	rowMargins:=make([]float64,0,256)
	for i:=0;i<256;i++ {
		var total uint64
		var top,runner uint32
		for j:=0;j<256;j++ {
			v:=rel[i][j]
			total+=uint64(v)
			if v>top { runner=top; top=v } else if v>runner { runner=v }
		}
		if total>=4 {
			eligible[i]=true
			marginsByRow[i]=(float64(top)-float64(runner))/float64(total)
			rowMargins=append(rowMargins,marginsByRow[i])
		}
	}
	if len(rowMargins)==0 { m["invalid_row_count"]++ }
	median:=wlmLmProspectiveR3Median(rowMargins)

	gateMargin:=func(dist *[512]float64) float64 {
		top:=wlmLmExternalReasoningReadoutRefinementR1Top1(dist)
		if top<0||top>=256 { m["invalid_row_count"]++; return 0 }
		cid:=classes[top]
		mass,runner:=0.0,0.0
		for k:=0;k<256;k++ {
			if classes[k]!=cid { continue }
			v:=dist[k]
			mass+=v
			if k!=top && v>runner { runner=v }
		}
		if mass<=0 { m["invalid_row_count"]++; return 0 }
		return (dist[top]-runner)/mass
	}

	classTrainMargins:=make([][]float64,8)
	for _,st:=range streams {
		for j:=0;j+4<len(st);j++ {
			rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(
				st[j],st[j+1],3,2048,true,
				&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal,
			)
			if re>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=re }
			if rr<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=rr }
			_ = gateMargin(&rd)
			top:=wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
			if top<0||top>=256 { m["invalid_row_count"]++; continue }
			classTrainMargins[classes[top]]=append(classTrainMargins[classes[top]],gateMargin(&rd))
			m["training_margin_sample_count"]++
		}
	}
	zeroSet:=make(map[int]bool)
	for i:=0;i<8;i++ {
		n:=len(classTrainMargins[i])
		m["class_"+string(rune('0'+i))+"_training_margin_sample_count"]=float64(n)
		if n==0 {
			m["zero_training_margin_class_count"]++
			zeroSet[i]=true
		}
	}
	if len(zeroSet)>0 { m["preflight_unevaluable_flag"]=1 }
	m["preflight_decision_frozen_flag"]=1

	// Target bytes are used scientifically only after the training-only preflight is frozen.
	target:=code[:1454]
	targetStream:=wlmLmExternalMotifRelationHoldoutR1Decode(target,index)
	for j:=0;j+4<len(targetStream);j++ {
		truth:=targetStream[j+4]
		rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(
			targetStream[j],targetStream[j+1],3,2048,true,
			&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal,
		)
		if re>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=re }
		if rr<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=rr }
		top:=wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
		if top<0||top>=256 { m["invalid_row_count"]++; continue }
		m["target_query_count"]++
		if zeroSet[classes[top]] {
			m["unsupported_target_query_count"]++
			if top==truth { m["unsupported_target_exact_hit_count"]++ }
		}
	}
	if m["maximum_probability_mass_error"]>1e-9 { m["invalid_row_count"]++ }
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmExternalReasoningSelfDiagnosisProspectiveR3Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-REASONING-SELF-DIAGNOSIS-PROSPECTIVE-R3",
		Metrics:m,
	}
}
