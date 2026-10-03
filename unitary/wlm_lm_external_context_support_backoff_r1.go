package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalContextSupportBackoffR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalContextSupportBackoffR1(code,structured,prose []byte) interface{} {
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
		"current_predictive_counter_capacity":262144,
		"selector_context_query_count":0,
		"selector_backoff_query_count":0,
		"selector_support_violation_count":0,
		"selector_overflow_or_unseen_backoff_query_count":0,
		"selector_support_1_2_backoff_query_count":0,
		"pooled_hybrid_top1_accuracy":0,
		"pooled_current_top1_accuracy":0,
		"pooled_ungated_context_top1_accuracy":0,
		"pooled_hybrid_minus_current_accuracy":0,
		"pooled_hybrid_minus_ungated_context_accuracy":0,
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
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		streams[si]=stream
		for j:=1;j<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			if pairCounts[p]==^uint32(0) { metrics["counter_overflow_count"]++; continue }
			pairCounts[p]++
		}
	}

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
	for i:=0;i<dedicated;i++ { contextID[rows[i].pair]=i }
	metrics["dedicated_pair_context_count"]=float64(dedicated)

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

	pooledQ:=0
	pooledHybrid:=0
	pooledCurrent:=0
	pooledContext:=0
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		if len(stream)<3 { continue }
		q:=0
		hybridCorrect:=0
		currentCorrect:=0
		contextCorrect:=0
		contextQueries:=0
		backoffQueries:=0
		for j:=1;j+1<len(stream);j++ {
			pair:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			cid:=511
			id,isDedicated:=contextID[pair]
			if isDedicated { cid=id }
			support:=pairCounts[pair]
			useContext:=isDedicated && support>=3
			if useContext {
				contextQueries++
				metrics["selector_context_query_count"]++
				if !isDedicated || support<3 { metrics["selector_support_violation_count"]++ }
			} else {
				backoffQueries++
				metrics["selector_backoff_query_count"]++
				if !isDedicated || support==0 {
					metrics["selector_overflow_or_unseen_backoff_query_count"]++
				} else if support<=2 {
					metrics["selector_support_1_2_backoff_query_count"]++
				}
			}

			target:=stream[j+1]
			cp,ctotal:=wlmLmExternalMotifRelationHoldoutR1Best(&context[cid])
			if ctotal==0 { cp=globalBest }
			rp,rtotal:=wlmLmExternalMotifRelationHoldoutR1Best(&current[stream[j]])
			if rtotal==0 { rp=globalBest }
			hp:=rp
			if useContext { hp=cp }

			if hp==target { hybridCorrect++; pooledHybrid++ }
			if rp==target { currentCorrect++; pooledCurrent++ }
			if cp==target { contextCorrect++; pooledContext++ }
			q++; pooledQ++
		}
		if q>0 {
			ha:=float64(hybridCorrect)/float64(q)
			ra:=float64(currentCorrect)/float64(q)
			ca:=float64(contextCorrect)/float64(q)
			delta:=ha-ra
			metrics["eval_"+sources[si].domain+"_hybrid_top1_accuracy"]=ha
			metrics["eval_"+sources[si].domain+"_current_top1_accuracy"]=ra
			metrics["eval_"+sources[si].domain+"_ungated_context_top1_accuracy"]=ca
			metrics["eval_"+sources[si].domain+"_hybrid_minus_current_accuracy"]=delta
			metrics["eval_"+sources[si].domain+"_hybrid_minus_ungated_context_accuracy"]=ha-ca
			metrics["eval_"+sources[si].domain+"_context_use_fraction"]=float64(contextQueries)/float64(q)
			metrics["eval_"+sources[si].domain+"_backoff_fraction"]=float64(backoffQueries)/float64(q)
			if delta<metrics["minimum_source_hybrid_minus_current_accuracy"] {
				metrics["minimum_source_hybrid_minus_current_accuracy"]=delta
			}
		}
	}
	if pooledQ>0 {
		metrics["pooled_hybrid_top1_accuracy"]=float64(pooledHybrid)/float64(pooledQ)
		metrics["pooled_current_top1_accuracy"]=float64(pooledCurrent)/float64(pooledQ)
		metrics["pooled_ungated_context_top1_accuracy"]=float64(pooledContext)/float64(pooledQ)
		metrics["pooled_hybrid_minus_current_accuracy"]=metrics["pooled_hybrid_top1_accuracy"]-metrics["pooled_current_top1_accuracy"]
		metrics["pooled_hybrid_minus_ungated_context_accuracy"]=metrics["pooled_hybrid_top1_accuracy"]-metrics["pooled_ungated_context_top1_accuracy"]
	} else {
		metrics["invalid_row_count"]++
	}
	if math.IsInf(metrics["minimum_source_hybrid_minus_current_accuracy"],1) {
		metrics["minimum_source_hybrid_minus_current_accuracy"]=0
		metrics["invalid_row_count"]++
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++}
	}
	return wlmLmExternalContextSupportBackoffR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-CONTEXT-SUPPORT-BACKOFF-R1",
		Metrics:metrics,
	}
}
