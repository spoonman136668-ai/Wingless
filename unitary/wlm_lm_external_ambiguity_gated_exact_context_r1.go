package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalAmbiguityGatedExactContextR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalAmbiguityGatedExactContextR1(code,structured,prose []byte) interface{} {
	sources:=[]wlmLmExternalContextStateHoldoutR1Source{
		{domain:"code",data:code,sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
		{domain:"structured",data:structured,sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
		{domain:"technical_prose",data:prose,sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
	}
	metrics:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":float64(len(sources)),
		"total_source_bytes":0,
		"selected_motif_count":0,
		"ambiguity_gated_query_count":0,
		"selector_rule_violation_count":0,
		"ambiguity_median_margin":0,
		"pooled_hybrid_top1_accuracy":0,
		"pooled_current_top1_accuracy":0,
		"pooled_support_backoff_top1_accuracy":0,
		"pooled_hybrid_minus_current_accuracy":0,
		"pooled_hybrid_minus_support_backoff_accuracy":0,
		"minimum_source_hybrid_minus_current_accuracy":math.Inf(1),
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}
	train:=make([][]byte,0,len(sources))
	evals:=make([][]byte,0,len(sources))
	for _,src:=range sources {
		metrics["total_source_bytes"]+=float64(len(src.data))
		if len(src.data)!=src.bytes || wlmLmExternalRawRepPredR1SHA256(src.data)!=src.sha256 { metrics["source_identity_mismatch_count"]++ }
		split:=len(src.data)*3/5
		if split<5 || len(src.data)-split<5 { metrics["invalid_row_count"]++ }
		train=append(train,src.data[:split]); evals=append(evals,src.data[split:])
	}
	_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(train,metrics)
	metrics["selected_motif_count"]=float64(len(selected))
	keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index:=make(map[[4]uint8]int,len(keys)); for i,k:=range keys { index[k]=i }

	streams:=make([][]int,len(train))
	var current [512][512]uint32
	var unigram [512]uint32
	pairCounts:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index); streams[si]=stream
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if current[a][b]==^uint32(0)||unigram[b]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			current[a][b]++; unigram[b]++
		}
		for j:=1;j<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			if pairCounts[p]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			pairCounts[p]++
		}
	}
	rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairCounts))
	for p,c:=range pairCounts { rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:p,count:c}) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count{return rows[i].count>rows[j].count}
		if rows[i].pair.prev!=rows[j].pair.prev{return rows[i].pair.prev<rows[j].pair.prev}
		return rows[i].pair.curr<rows[j].pair.curr
	})
	dedicated:=511; if len(rows)<dedicated { dedicated=len(rows) }
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,dedicated)
	for i:=0;i<dedicated;i++ { exactID[rows[i].pair]=i }
	var exact [512][512]uint32
	for _,stream:=range streams {
		for j:=1;j+1<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			id:=511; if v,ok:=exactID[p];ok{id=v}
			target:=stream[j+1]
			if exact[id][target]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			exact[id][target]++
		}
	}

	margins:=make([]float64,0,512)
	rowMargin:=make([]float64,512)
	rowEligible:=make([]bool,512)
	for i:=0;i<512;i++ {
		var total uint64; var top1,top2 uint32
		for j:=0;j<512;j++ {
			v:=current[i][j]; total+=uint64(v)
			if v>top1 { top2=top1; top1=v } else if v>top2 { top2=v }
		}
		if total>=4 {
			m:=(float64(top1)-float64(top2))/float64(total)
			rowMargin[i]=m; rowEligible[i]=true; margins=append(margins,m)
		}
	}
	sort.Float64s(margins)
	median:=0.0
	if len(margins)==0 { metrics["invalid_row_count"]++ } else if len(margins)%2==1 { median=margins[len(margins)/2] } else { median=(margins[len(margins)/2-1]+margins[len(margins)/2])/2 }
	metrics["ambiguity_median_margin"]=median
	globalBest:=0; for i:=1;i<512;i++ { if unigram[i]>unigram[globalBest] { globalBest=i } }

	total:=0; hybridHits,currentHits,backoffHits:=0,0,0
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		q:=0; hh,ch,bh:=0,0,0
		for j:=1;j+1<len(stream);j++ {
			prev,curr,target:=stream[j-1],stream[j],stream[j+1]
			cp,ct:=wlmLmExternalMotifRelationHoldoutR1Best(&current[curr]); if ct==0 {cp=globalBest}
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:prev,curr:curr}
			id,isDedicated:=exactID[p]; if !isDedicated {id=511}
			ep,et:=wlmLmExternalMotifRelationHoldoutR1Best(&exact[id]); if et==0 {ep=globalBest}
			useBackoff:=isDedicated && pairCounts[p]>=3
			bp:=cp; if useBackoff {bp=ep}
			useHybrid:=useBackoff && rowEligible[curr] && rowMargin[curr]<=median
			hp:=cp; if useHybrid { hp=ep; metrics["ambiguity_gated_query_count"]++ }
			if useHybrid && (!isDedicated || pairCounts[p]<3 || !rowEligible[curr] || rowMargin[curr]>median) { metrics["selector_rule_violation_count"]++ }
			if hp==target {hh++;hybridHits++}
			if cp==target {ch++;currentHits++}
			if bp==target {bh++;backoffHits++}
			q++;total++
		}
		if q>0 {
			ha:=float64(hh)/float64(q); ca:=float64(ch)/float64(q); ba:=float64(bh)/float64(q)
			metrics["eval_"+sources[si].domain+"_hybrid_top1_accuracy"]=ha
			metrics["eval_"+sources[si].domain+"_current_top1_accuracy"]=ca
			metrics["eval_"+sources[si].domain+"_support_backoff_top1_accuracy"]=ba
			metrics["eval_"+sources[si].domain+"_hybrid_minus_current_accuracy"]=ha-ca
			if ha-ca<metrics["minimum_source_hybrid_minus_current_accuracy"] { metrics["minimum_source_hybrid_minus_current_accuracy"]=ha-ca }
		}
	}
	if total>0 {
		metrics["pooled_hybrid_top1_accuracy"]=float64(hybridHits)/float64(total)
		metrics["pooled_current_top1_accuracy"]=float64(currentHits)/float64(total)
		metrics["pooled_support_backoff_top1_accuracy"]=float64(backoffHits)/float64(total)
		metrics["pooled_hybrid_minus_current_accuracy"]=metrics["pooled_hybrid_top1_accuracy"]-metrics["pooled_current_top1_accuracy"]
		metrics["pooled_hybrid_minus_support_backoff_accuracy"]=metrics["pooled_hybrid_top1_accuracy"]-metrics["pooled_support_backoff_top1_accuracy"]
	} else { metrics["invalid_row_count"]++ }
	if math.IsInf(metrics["minimum_source_hybrid_minus_current_accuracy"],1) { metrics["minimum_source_hybrid_minus_current_accuracy"]=0; metrics["invalid_row_count"]++ }
	for _,v:=range metrics { if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++} }
	return wlmLmExternalAmbiguityGatedExactContextR1Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-AMBIGUITY-GATED-EXACT-CONTEXT-R1",Metrics:metrics}
}
