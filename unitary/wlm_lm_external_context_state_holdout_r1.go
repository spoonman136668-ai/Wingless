package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalContextStateHoldoutR1Source struct {
	domain string
	data []byte
	sha256 string
	bytes int
}

type wlmLmExternalContextStateHoldoutR1Pair struct {
	prev int
	curr int
}

type wlmLmExternalContextStateHoldoutR1PairRow struct {
	pair wlmLmExternalContextStateHoldoutR1Pair
	count uint32
}

type wlmLmExternalContextStateHoldoutR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalContextStateHoldoutR1(code,structured,prose []byte) interface{} {
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
		"context_state_count":512,
		"dedicated_pair_context_count":0,
		"context_predictive_counter_capacity":262144,
		"training_distinct_pair_count":0,
		"training_pair_coverage_captured_by_dedicated_contexts":0,
		"eval_source_count_with_triples":0,
		"pooled_context_top1_accuracy":0,
		"pooled_current_top1_accuracy":0,
		"pooled_context_minus_current_accuracy":0,
		"overflow_context_query_fraction":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
		"counter_overflow_rows":0,
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
	pairCounts:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	totalPairOccurrences:=uint64(0)
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		streams[si]=stream
		for j:=1;j<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			if pairCounts[p]==^uint32(0) { metrics["counter_overflow_count"]++;continue }
			pairCounts[p]++
			totalPairOccurrences++
		}
	}
	metrics["training_distinct_pair_count"]=float64(len(pairCounts))

	rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairCounts))
	for p,c:=range pairCounts { rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:p,count:c}) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count { return rows[i].count>rows[j].count }
		if rows[i].pair.prev!=rows[j].pair.prev { return rows[i].pair.prev<rows[j].pair.prev }
		return rows[i].pair.curr<rows[j].pair.curr
	})
	dedicated:=511
	if len(rows)<dedicated { dedicated=len(rows) }
	contextID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,dedicated)
	dedicatedOccurrences:=uint64(0)
	for i:=0;i<dedicated;i++ {
		contextID[rows[i].pair]=i
		dedicatedOccurrences+=uint64(rows[i].count)
	}
	metrics["dedicated_pair_context_count"]=float64(dedicated)
	if totalPairOccurrences>0 {
		metrics["training_pair_coverage_captured_by_dedicated_contexts"]=float64(dedicatedOccurrences)/float64(totalPairOccurrences)
	} else {
		metrics["invalid_row_count"]++
	}

	var context [512][512]uint32
	var current [512][512]uint32
	var unigram [512]uint32
	for _,stream:=range streams {
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if current[a][b]==^uint32(0)||unigram[b]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			current[a][b]++
			unigram[b]++
		}
		for j:=1;j+1<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			cid:=511
			if id,ok:=contextID[p];ok { cid=id }
			target:=stream[j+1]
			if context[cid][target]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			context[cid][target]++
		}
	}
	globalBest:=0
	for i:=1;i<512;i++ { if unigram[i]>unigram[globalBest] { globalBest=i } }

	pooledContextCorrect:=0
	pooledCurrentCorrect:=0
	pooledQueries:=0
	pooledOverflow:=0
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		if len(stream)<3 { continue }
		metrics["eval_source_count_with_triples"]++
		contextCorrect:=0
		currentCorrect:=0
		queries:=0
		overflowQueries:=0
		for j:=1;j+1<len(stream);j++ {
			pair:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			cid:=511
			if id,ok:=contextID[pair];ok { cid=id } else { overflowQueries++;pooledOverflow++ }
			target:=stream[j+1]
			cp,ctotal:=wlmLmExternalMotifRelationHoldoutR1Best(&context[cid])
			if ctotal==0 { cp=globalBest }
			rp,rtotal:=wlmLmExternalMotifRelationHoldoutR1Best(&current[stream[j]])
			if rtotal==0 { rp=globalBest }
			if cp==target { contextCorrect++;pooledContextCorrect++ }
			if rp==target { currentCorrect++;pooledCurrentCorrect++ }
			queries++;pooledQueries++
		}
		contextAcc:=float64(contextCorrect)/float64(queries)
		currentAcc:=float64(currentCorrect)/float64(queries)
		delta:=contextAcc-currentAcc
		metrics["eval_"+sources[si].domain+"_context_top1_accuracy"]=contextAcc
		metrics["eval_"+sources[si].domain+"_current_top1_accuracy"]=currentAcc
		metrics["eval_"+sources[si].domain+"_context_minus_current_accuracy"]=delta
		metrics["eval_"+sources[si].domain+"_dedicated_context_fraction"]=1-float64(overflowQueries)/float64(queries)
	}
	if pooledQueries>0 {
		metrics["pooled_context_top1_accuracy"]=float64(pooledContextCorrect)/float64(pooledQueries)
		metrics["pooled_current_top1_accuracy"]=float64(pooledCurrentCorrect)/float64(pooledQueries)
		metrics["pooled_context_minus_current_accuracy"]=metrics["pooled_context_top1_accuracy"]-metrics["pooled_current_top1_accuracy"]
		metrics["overflow_context_query_fraction"]=float64(pooledOverflow)/float64(pooledQueries)
	} else {
		metrics["invalid_row_count"]++
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++}
	}
	return wlmLmExternalContextStateHoldoutR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-CONTEXT-STATE-HOLDOUT-R1",
		Metrics:metrics,
	}
}
