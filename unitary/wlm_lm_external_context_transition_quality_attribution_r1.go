package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalContextTransitionQualityAttributionR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalContextTransitionQualityAttributionR1(code,structured,prose []byte) interface{} {
	sources:=[]wlmLmExternalMotifRelationDistributionalR1Source{
		{domain:"code",data:code,sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
		{domain:"structured",data:structured,sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
		{domain:"technical_prose",data:prose,sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
	}
	metrics:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":float64(len(sources)),
		"total_source_bytes":0,
		"selected_motif_count":0,
		"depth2_query_count":0,
		"depth3_query_count":0,
		"depth2_error_count":0,
		"depth3_error_count":0,
		"total_error_count":0,
		"stratified_error_count_sum":0,
		"heldout_selection_use_count":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
		"ambiguity_median_margin":0,
		"depth3_low_support_query_share":0,
		"depth3_low_support_error_share":0,
		"depth3_low_support_error_excess":0,
		"depth3_high_support_query_share":0,
		"depth3_high_support_error_share":0,
		"depth3_high_support_error_excess":0,
		"depth3_high_ambiguity_query_share":0,
		"depth3_high_ambiguity_error_share":0,
		"depth3_high_ambiguity_error_excess":0,
		"depth3_low_ambiguity_query_share":0,
		"depth3_low_ambiguity_error_share":0,
		"depth3_low_ambiguity_error_excess":0,
		"depth3_primary_attribution_excess":0,
	}

	train:=make([][]byte,0,len(sources))
	evals:=make([][]byte,0,len(sources))
	for _,src:=range sources {
		metrics["total_source_bytes"]+=float64(len(src.data))
		if len(src.data)!=src.bytes || wlmLmExternalRawRepPredR1SHA256(src.data)!=src.sha256 {
			metrics["source_identity_mismatch_count"]++
		}
		split:=len(src.data)*3/5
		if split<5 || len(src.data)-split<5 { metrics["invalid_row_count"]++ }
		train=append(train,src.data[:split])
		evals=append(evals,src.data[split:])
	}

	_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(train,metrics)
	metrics["selected_motif_count"]=float64(len(selected))
	keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index:=make(map[[4]uint8]int,len(keys))
	for i,k:=range keys { index[k]=i }

	streams:=make([][]int,len(train))
	var relation [512][512]uint32
	var rowTotals [512]uint64
	var unigram [512]uint32
	var unigramTotal uint64
	pairCounts:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		streams[si]=stream
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if relation[a][b]==^uint32(0)||unigram[b]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			relation[a][b]++;rowTotals[a]++;unigram[b]++;unigramTotal++
		}
		for j:=1;j<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			if pairCounts[p]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			pairCounts[p]++
		}
	}
	if unigramTotal==0 { metrics["invalid_row_count"]++ }

	rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairCounts))
	for p,c:=range pairCounts { rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:p,count:c}) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count{return rows[i].count>rows[j].count}
		if rows[i].pair.prev!=rows[j].pair.prev{return rows[i].pair.prev<rows[j].pair.prev}
		return rows[i].pair.curr<rows[j].pair.curr
	})
	dedicated:=511
	if len(rows)<dedicated { dedicated=len(rows) }
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,dedicated)
	for i:=0;i<dedicated;i++ { exactID[rows[i].pair]=i }
	var exact [512][512]uint32
	var exactTotals [512]uint64
	for _,stream:=range streams {
		for j:=1;j+1<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			id:=511
			if v,ok:=exactID[p];ok{id=v}
			target:=stream[j+1]
			if exact[id][target]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			exact[id][target]++;exactTotals[id]++
		}
	}

	var rowEligible [512]bool
	var rowMargin [512]float64
	margins:=make([]float64,0,512)
	for i:=0;i<512;i++ {
		var total uint64
		var top1,top2 uint32
		for _,v:=range relation[i] {
			total+=uint64(v)
			if v>top1 { top2=top1;top1=v } else if v>top2 { top2=v }
		}
		if total>=4 {
			rowEligible[i]=true
			rowMargin[i]=(float64(top1)-float64(top2))/float64(total)
			margins=append(margins,rowMargin[i])
		}
	}
	sort.Float64s(margins)
	median:=0.0
	if len(margins)==0 {
		metrics["invalid_row_count"]++
	} else if len(margins)%2==1 {
		median=margins[len(margins)/2]
	} else {
		median=(margins[len(margins)/2-1]+margins[len(margins)/2])/2
	}
	metrics["ambiguity_median_margin"]=median

	type counts struct{ q,e int }
	var d2LowSupport,d2HighSupport,d2HighAmb,d2LowAmb counts
	var d3LowSupport,d3HighSupport,d3HighAmb,d3LowAmb counts

	for _,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		for _,depth:=range []int{2,3} {
			for j:=0;j+1+depth<len(stream);j++ {
				prev,curr,target:=stream[j],stream[j+1],stream[j+1+depth]
				pred,_,_:=wlmLmExternalCompositionErrorAttributionR2Predict(
					prev,curr,depth,2048,&relation,&rowTotals,&exact,&exactTotals,
					exactID,pairCounts,&rowEligible,&rowMargin,median,&unigram,unigramTotal,
				)
				p:=wlmLmExternalContextStateHoldoutR1Pair{prev:prev,curr:curr}
				lowSupport:=pairCounts[p]<3
				highAmbiguity:=!rowEligible[curr] || rowMargin[curr]<median
				err:=pred!=target

				if depth==2 {
					metrics["depth2_query_count"]++
					if err { metrics["depth2_error_count"]++;metrics["total_error_count"]++ }
					if lowSupport { d2LowSupport.q++;if err{d2LowSupport.e++} } else { d2HighSupport.q++;if err{d2HighSupport.e++} }
					if highAmbiguity { d2HighAmb.q++;if err{d2HighAmb.e++} } else { d2LowAmb.q++;if err{d2LowAmb.e++} }
				} else {
					metrics["depth3_query_count"]++
					if err { metrics["depth3_error_count"]++;metrics["total_error_count"]++ }
					if lowSupport { d3LowSupport.q++;if err{d3LowSupport.e++} } else { d3HighSupport.q++;if err{d3HighSupport.e++} }
					if highAmbiguity { d3HighAmb.q++;if err{d3HighAmb.e++} } else { d3LowAmb.q++;if err{d3LowAmb.e++} }
				}
			}
		}
	}

	metrics["stratified_error_count_sum"]=float64(d2LowSupport.e+d2HighSupport.e+d3LowSupport.e+d3HighSupport.e)
	setShares:=func(prefix string,a,b counts,totalQ,totalE int){
		if totalQ<=0 { metrics["invalid_row_count"]++;return }
		aq:=float64(a.q)/float64(totalQ)
		bq:=float64(b.q)/float64(totalQ)
		ae,be:=0.0,0.0
		if totalE>0 {
			ae=float64(a.e)/float64(totalE)
			be=float64(b.e)/float64(totalE)
		}
		metrics[prefix+"_a_query_share"]=aq
		metrics[prefix+"_a_error_share"]=ae
		metrics[prefix+"_a_error_excess"]=ae-aq
		metrics[prefix+"_b_query_share"]=bq
		metrics[prefix+"_b_error_share"]=be
		metrics[prefix+"_b_error_excess"]=be-bq
	}
	setShares("depth2_support",d2LowSupport,d2HighSupport,d2LowSupport.q+d2HighSupport.q,d2LowSupport.e+d2HighSupport.e)
	setShares("depth2_ambiguity",d2HighAmb,d2LowAmb,d2HighAmb.q+d2LowAmb.q,d2HighAmb.e+d2LowAmb.e)
	setShares("depth3_support",d3LowSupport,d3HighSupport,d3LowSupport.q+d3HighSupport.q,d3LowSupport.e+d3HighSupport.e)
	setShares("depth3_ambiguity",d3HighAmb,d3LowAmb,d3HighAmb.q+d3LowAmb.q,d3HighAmb.e+d3LowAmb.e)

	metrics["depth3_low_support_query_share"]=metrics["depth3_support_a_query_share"]
	metrics["depth3_low_support_error_share"]=metrics["depth3_support_a_error_share"]
	metrics["depth3_low_support_error_excess"]=metrics["depth3_support_a_error_excess"]
	metrics["depth3_high_support_query_share"]=metrics["depth3_support_b_query_share"]
	metrics["depth3_high_support_error_share"]=metrics["depth3_support_b_error_share"]
	metrics["depth3_high_support_error_excess"]=metrics["depth3_support_b_error_excess"]
	metrics["depth3_high_ambiguity_query_share"]=metrics["depth3_ambiguity_a_query_share"]
	metrics["depth3_high_ambiguity_error_share"]=metrics["depth3_ambiguity_a_error_share"]
	metrics["depth3_high_ambiguity_error_excess"]=metrics["depth3_ambiguity_a_error_excess"]
	metrics["depth3_low_ambiguity_query_share"]=metrics["depth3_ambiguity_b_query_share"]
	metrics["depth3_low_ambiguity_error_share"]=metrics["depth3_ambiguity_b_error_share"]
	metrics["depth3_low_ambiguity_error_excess"]=metrics["depth3_ambiguity_b_error_excess"]
	primary:=metrics["depth3_low_support_error_excess"]
	if metrics["depth3_high_ambiguity_error_excess"]>primary { primary=metrics["depth3_high_ambiguity_error_excess"] }
	metrics["depth3_primary_attribution_excess"]=primary

	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++}
	}
	return wlmLmExternalContextTransitionQualityAttributionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-CONTEXT-TRANSITION-QUALITY-ATTRIBUTION-R1",
		Metrics:metrics,
	}
}
