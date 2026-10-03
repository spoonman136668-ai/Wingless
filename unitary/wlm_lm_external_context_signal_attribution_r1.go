package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalContextSignalAttributionR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalContextSignalAttributionR1(code,structured,prose []byte) interface{} {
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
		"eligible_current_row_count":0,
		"ambiguity_median_margin":0,
		"supported_pair_query_count":0,
		"ambiguous_supported_query_count":0,
		"confident_supported_query_count":0,
		"ambiguous_supported_exact_top1_accuracy":0,
		"ambiguous_supported_current_top1_accuracy":0,
		"ambiguous_supported_exact_minus_current_accuracy":0,
		"confident_supported_exact_top1_accuracy":0,
		"confident_supported_current_top1_accuracy":0,
		"confident_supported_exact_minus_current_accuracy":0,
		"context_signal_concentration_gap":0,
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
		train=append(train,src.data[:split])
		evals=append(evals,src.data[split:])
	}

	_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(train,metrics)
	metrics["selected_motif_count"]=float64(len(selected))
	keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index:=make(map[[4]uint8]int,len(keys))
	for i,k:=range keys { index[k]=i }

	streams:=make([][]int,len(train))
	var current [512][512]uint32
	var unigram [512]uint32
	pairCounts:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		streams[si]=stream
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
	dedicated:=511
	if len(rows)<dedicated { dedicated=len(rows) }
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,dedicated)
	for i:=0;i<dedicated;i++ { exactID[rows[i].pair]=i }

	var exact [512][512]uint32
	for _,stream:=range streams {
		for j:=1;j+1<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			id:=511
			if v,ok:=exactID[p];ok{id=v}
			target:=stream[j+1]
			if exact[id][target]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			exact[id][target]++
		}
	}

	margins:=make([]float64,0,512)
	rowMargin:=make([]float64,512)
	rowEligible:=make([]bool,512)
	for i:=0;i<512;i++ {
		var total uint64
		var top1,top2 uint32
		for j:=0;j<512;j++ {
			v:=current[i][j]
			total+=uint64(v)
			if v>top1 { top2=top1; top1=v } else if v>top2 { top2=v }
		}
		if total>=4 {
			m:=(float64(top1)-float64(top2))/float64(total)
			rowMargin[i]=m
			rowEligible[i]=true
			margins=append(margins,m)
		}
	}
	sort.Float64s(margins)
	metrics["eligible_current_row_count"]=float64(len(margins))
	median:=0.0
	if len(margins)==0 {
		metrics["invalid_row_count"]++
	} else if len(margins)%2==1 {
		median=margins[len(margins)/2]
	} else {
		median=(margins[len(margins)/2-1]+margins[len(margins)/2])/2
	}
	metrics["ambiguity_median_margin"]=median

	globalBest:=0
	for i:=1;i<512;i++ { if unigram[i]>unigram[globalBest] { globalBest=i } }

	ambQ,ambExact,ambCurrent:=0,0,0
	confQ,confExact,confCurrent:=0,0,0
	for _,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		for j:=1;j+1<len(stream);j++ {
			prev,curr,target:=stream[j-1],stream[j],stream[j+1]
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:prev,curr:curr}
			id,ok:=exactID[p]
			if !ok || pairCounts[p]<3 || !rowEligible[curr] { continue }
			ep,et:=wlmLmExternalMotifRelationHoldoutR1Best(&exact[id]); if et==0 { ep=globalBest }
			cp,ct:=wlmLmExternalMotifRelationHoldoutR1Best(&current[curr]); if ct==0 { cp=globalBest }
			metrics["supported_pair_query_count"]++
			if rowMargin[curr]<=median {
				ambQ++
				if ep==target { ambExact++ }
				if cp==target { ambCurrent++ }
			} else {
				confQ++
				if ep==target { confExact++ }
				if cp==target { confCurrent++ }
			}
		}
	}
	metrics["ambiguous_supported_query_count"]=float64(ambQ)
	metrics["confident_supported_query_count"]=float64(confQ)
	if ambQ>0 {
		e:=float64(ambExact)/float64(ambQ); c:=float64(ambCurrent)/float64(ambQ)
		metrics["ambiguous_supported_exact_top1_accuracy"]=e
		metrics["ambiguous_supported_current_top1_accuracy"]=c
		metrics["ambiguous_supported_exact_minus_current_accuracy"]=e-c
	} else { metrics["invalid_row_count"]++ }
	if confQ>0 {
		e:=float64(confExact)/float64(confQ); c:=float64(confCurrent)/float64(confQ)
		metrics["confident_supported_exact_top1_accuracy"]=e
		metrics["confident_supported_current_top1_accuracy"]=c
		metrics["confident_supported_exact_minus_current_accuracy"]=e-c
	} else { metrics["invalid_row_count"]++ }
	metrics["context_signal_concentration_gap"]=metrics["ambiguous_supported_exact_minus_current_accuracy"]-metrics["confident_supported_exact_minus_current_accuracy"]

	for _,v:=range metrics { if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++} }
	return wlmLmExternalContextSignalAttributionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-CONTEXT-SIGNAL-ATTRIBUTION-R1",
		Metrics:metrics,
	}
}
