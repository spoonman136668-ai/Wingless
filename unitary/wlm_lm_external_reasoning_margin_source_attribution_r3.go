package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalReasoningMarginSourceAttributionR3Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmMarginSourceMedianR3(xs []float64) float64 {
	sort.Float64s(xs)
	if len(xs)==0 { return 0 }
	n:=len(xs)
	if n%2==1 { return xs[n/2] }
	return (xs[n/2-1]+xs[n/2])/2
}

func RunWlmLmExternalReasoningMarginSourceAttributionR3(code, structured, prose []byte) interface{} {
	sources:=[]wlmLmExternalMotifRelationDistributionalR1Source{
		{domain:"code",data:code,sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
		{domain:"structured",data:structured,sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
		{domain:"technical_prose",data:prose,sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,
		"selected_motif_count":0,"training_margin_sample_count":0,"training_margin_threshold":0,
		"depth3_query_count":0,"same_class_hidden_gain_count":0,"same_class_hidden_below_threshold_count":0,
		"same_class_source_code_count":0,"same_class_source_structured_count":0,"same_class_source_technical_prose_count":0,
		"source_training_margin_code_sample_count":0,"source_training_margin_structured_sample_count":0,"source_training_margin_technical_prose_sample_count":0,
		"source_training_margin_code_median":0,"source_training_margin_structured_median":0,"source_training_margin_technical_prose_median":0,
		"code_structured_training_margin_median_gap":0,
		"source_below_threshold_code_query_count":0,"source_below_threshold_structured_query_count":0,"source_below_threshold_technical_prose_query_count":0,
		"source_below_threshold_code_exact_hit_count":0,"source_below_threshold_structured_exact_hit_count":0,"source_below_threshold_technical_prose_exact_hit_count":0,
		"source_below_threshold_code_exact_precision":0,"source_below_threshold_structured_exact_precision":0,"source_below_threshold_technical_prose_exact_precision":0,
		"code_structured_below_threshold_precision_gap":0,"source_attribution_accounting_error_count":0,
		"maximum_probability_mass_error":0,"minimum_repaired_retained_mass":1,
		"heldout_selection_use_count":0,"relation_capacity_growth_event_count":0,"adaptive_readout_growth_event_count":0,
		"tokenizer_use_count":0,"external_model_call_count":0,"invalid_row_count":0,"counter_overflow_count":0,
	}
	train:=make([][]byte,0,3); evals:=make([][]byte,0,3)
	for _,s:=range sources {
		m["total_source_bytes"]+=float64(len(s.data))
		if len(s.data)!=s.bytes || wlmLmExternalRawRepPredR1SHA256(s.data)!=s.sha256 { m["source_identity_mismatch_count"]++ }
		q:=len(s.data)*3/5
		if q<5 || len(s.data)-q<5 { m["invalid_row_count"]++ }
		train=append(train,s.data[:q]); evals=append(evals,s.data[q:])
	}
	_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(train,m); m["selected_motif_count"]=float64(len(selected))
	keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index:=make(map[[4]uint8]int,len(keys)); for i,k:=range keys { index[k]=i }

	streams:=make([][]int,len(train))
	var counts [512]uint64
	var rel [512][512]uint32
	var totals [512]uint64
	var uni [512]uint32
	var unitotal uint64
	pairs:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si,data:=range train {
		st:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index); streams[si]=st
		for _,x:=range st { counts[x]++ }
		for j:=0;j+1<len(st);j++ {
			a,b:=st[j],st[j+1]
			if rel[a][b]==^uint32(0) || uni[b]==^uint32(0) { m["counter_overflow_count"]++; continue }
			rel[a][b]++; totals[a]++; uni[b]++; unitotal++
		}
		for j:=1;j<len(st);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:st[j-1],curr:st[j]}
			if pairs[p]==^uint32(0) { m["counter_overflow_count"]++; continue }
			pairs[p]++
		}
	}
	if unitotal==0 { m["invalid_row_count"]++ }

	ids:=make([]int,512); for i:=range ids { ids[i]=i }
	sort.Slice(ids,func(i,j int)bool {
		if counts[ids[i]]!=counts[ids[j]] { return counts[ids[i]]>counts[ids[j]] }
		return ids[i]<ids[j]
	})
	var classes [512]int
	for rank,id:=range ids { classes[id]=rank/64 }

	rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairs))
	for p,c:=range pairs { rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:p,count:c}) }
	sort.Slice(rows,func(i,j int)bool {
		if rows[i].count!=rows[j].count { return rows[i].count>rows[j].count }
		if rows[i].pair.prev!=rows[j].pair.prev { return rows[i].pair.prev<rows[j].pair.prev }
		return rows[i].pair.curr<rows[j].pair.curr
	})
	n:=511; if len(rows)<n { n=len(rows) }
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,n)
	for i:=0;i<n;i++ { exactID[rows[i].pair]=i }
	var exact [512][512]uint32
	var exactTotals [512]uint64
	for _,st:=range streams {
		for j:=1;j+1<len(st);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:st[j-1],curr:st[j]}
			id:=511; if v,ok:=exactID[p];ok { id=v }
			t:=st[j+1]
			if exact[id][t]==^uint32(0) { m["counter_overflow_count"]++; continue }
			exact[id][t]++; exactTotals[id]++
		}
	}
	var eligible [512]bool
	var marginsByRow [512]float64
	rowMargins:=make([]float64,0,512)
	for i:=0;i<512;i++ {
		var total uint64; var a,b uint32
		for _,v:=range rel[i] { total+=uint64(v); if v>a { b=a;a=v } else if v>b { b=v } }
		if total>=4 { eligible[i]=true; marginsByRow[i]=(float64(a)-float64(b))/float64(total); rowMargins=append(rowMargins,marginsByRow[i]) }
	}
	median:=wlmLmMarginSourceMedianR3(rowMargins)
	if len(rowMargins)==0 { m["invalid_row_count"]++ }

	gateMargin:=func(dist *[512]float64) float64 {
		top:=wlmLmExternalReasoningReadoutRefinementR1Top1(dist)
		cid:=classes[top]; mass,runner:=0.0,0.0
		for k,p:=range dist {
			if classes[k]!=cid { continue }
			mass+=p
			if k!=top && p>runner { runner=p }
		}
		if mass<=0 { m["invalid_row_count"]++; return 0 }
		return (dist[top]-runner)/mass
	}

	allTrainMargins:=make([]float64,0)
	sourceTrainMargins:=make([][]float64,3)
	for si,st:=range streams {
		for j:=0;j+4<len(st);j++ {
			rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(st[j],st[j+1],3,2048,true,&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal)
			if re>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=re }
			if rr<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=rr }
			gm:=gateMargin(&rd)
			allTrainMargins=append(allTrainMargins,gm)
			sourceTrainMargins[si]=append(sourceTrainMargins[si],gm)
		}
	}
	global:=wlmLmMarginSourceMedianR3(allTrainMargins)
	m["training_margin_sample_count"]=float64(len(allTrainMargins)); m["training_margin_threshold"]=global
	domains:=[]string{"code","structured","technical_prose"}
	for i,d:=range domains {
		m["source_training_margin_"+d+"_sample_count"]=float64(len(sourceTrainMargins[i]))
		if len(sourceTrainMargins[i])==0 { m["invalid_row_count"]++; continue }
		m["source_training_margin_"+d+"_median"]=wlmLmMarginSourceMedianR3(sourceTrainMargins[i])
	}
	m["code_structured_training_margin_median_gap"]=math.Abs(m["source_training_margin_code_median"]-m["source_training_margin_structured_median"])

	for si,data:=range evals {
		st:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index); d:=sources[si].domain
		for j:=0;j+4<len(st);j++ {
			target:=st[j+4]
			bd,be,_:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(st[j],st[j+1],3,2048,false,&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal)
			rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(st[j],st[j+1],3,2048,true,&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal)
			if be>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=be }
			if re>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=re }
			if rr<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=rr }
			m["depth3_query_count"]++
			gm:=gateMargin(&rd)
			btop:=wlmLmExternalReasoningReadoutRefinementR1Top1(&bd)
			rtop:=wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
			if gm<global {
				m["source_below_threshold_"+d+"_query_count"]++
				if rtop==target { m["source_below_threshold_"+d+"_exact_hit_count"]++ }
			}
			if btop==target || rtop!=target { continue }
			truth:=classes[target]
			bc:=wlmLmExternalReasoningReadoutRefinementR1Class(&bd,&classes)
			if classes[btop]!=classes[rtop] && bc!=truth { continue }
			if classes[btop]==classes[rtop] {
				m["same_class_hidden_gain_count"]++
				m["same_class_source_"+d+"_count"]++
				if gm<global { m["same_class_hidden_below_threshold_count"]++ }
			}
		}
	}
	for _,d:=range domains {
		q:=m["source_below_threshold_"+d+"_query_count"]
		if q>0 { m["source_below_threshold_"+d+"_exact_precision"]=m["source_below_threshold_"+d+"_exact_hit_count"]/q }
	}
	m["code_structured_below_threshold_precision_gap"]=math.Abs(m["source_below_threshold_code_exact_precision"]-m["source_below_threshold_structured_exact_precision"])
	if m["source_training_margin_code_sample_count"]+m["source_training_margin_structured_sample_count"]+m["source_training_margin_technical_prose_sample_count"]!=m["training_margin_sample_count"] { m["source_attribution_accounting_error_count"]++ }
	if m["same_class_source_code_count"]+m["same_class_source_structured_count"]+m["same_class_source_technical_prose_count"]!=m["same_class_hidden_gain_count"] { m["source_attribution_accounting_error_count"]++ }
	if m["maximum_probability_mass_error"]>1e-9 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmExternalReasoningMarginSourceAttributionR3Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-REASONING-MARGIN-SOURCE-ATTRIBUTION-R3",Metrics:m}
}
