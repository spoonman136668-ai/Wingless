package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalContextStateAttributionR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalContextStateAttributionR1(code,structured,prose []byte) interface{} {
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
		"dedicated_query_count":0,
		"overflow_query_count":0,
		"partition_query_count_mismatch":0,
		"pooled_dedicated_context_top1_accuracy":0,
		"pooled_dedicated_current_top1_accuracy":0,
		"pooled_dedicated_context_minus_current_accuracy":0,
		"pooled_overflow_context_top1_accuracy":0,
		"pooled_overflow_current_top1_accuracy":0,
		"pooled_overflow_context_minus_current_accuracy":0,
		"support_1_2_query_count":0,
		"support_1_2_context_minus_current_accuracy":0,
		"support_3_7_query_count":0,
		"support_3_7_context_minus_current_accuracy":0,
		"support_8_plus_query_count":0,
		"support_8_plus_context_minus_current_accuracy":0,
		"support_0_query_count":0,
		"training_pair_coverage_captured_by_dedicated_contexts":0,
		"evaluation_dedicated_context_fraction":0,
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
	} else { metrics["invalid_row_count"]++ }

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

	type tally struct{ q,context,current int }
	var dedicatedT,overflowT tally
	supportTallies:=map[string]*tally{
		"support_1_2":{},
		"support_3_7":{},
		"support_8_plus":{},
	}
	totalQueries:=0
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		if len(stream)<3 { continue }
		var sd,so tally
		for j:=1;j+1<len(stream);j++ {
			pair:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			cid:=511
			isDedicated:=false
			if id,ok:=contextID[pair];ok { cid=id;isDedicated=true }
			target:=stream[j+1]
			cp,ctotal:=wlmLmExternalMotifRelationHoldoutR1Best(&context[cid])
			if ctotal==0 { cp=globalBest }
			rp,rtotal:=wlmLmExternalMotifRelationHoldoutR1Best(&current[stream[j]])
			if rtotal==0 { rp=globalBest }

			var dst *tally
			if isDedicated { dst=&sd } else { dst=&so }
			dst.q++
			if cp==target { dst.context++ }
			if rp==target { dst.current++ }

			if isDedicated {
				dedicatedT.q++
				if cp==target { dedicatedT.context++ }
				if rp==target { dedicatedT.current++ }
			} else {
				overflowT.q++
				if cp==target { overflowT.context++ }
				if rp==target { overflowT.current++ }
			}
			support:=pairCounts[pair]
			var bin string
			switch {
			case support==0:
				metrics["support_0_query_count"]++
			case support<=2:
				bin="support_1_2"
			case support<=7:
				bin="support_3_7"
			default:
				bin="support_8_plus"
			}
			if bin!="" {
				t:=supportTallies[bin]
				t.q++
				if cp==target { t.context++ }
				if rp==target { t.current++ }
			}
			totalQueries++
		}
		if sd.q>0 {
			ca:=float64(sd.context)/float64(sd.q)
			ra:=float64(sd.current)/float64(sd.q)
			metrics["eval_"+sources[si].domain+"_dedicated_context_top1_accuracy"]=ca
			metrics["eval_"+sources[si].domain+"_dedicated_current_top1_accuracy"]=ra
			metrics["eval_"+sources[si].domain+"_dedicated_context_minus_current_accuracy"]=ca-ra
			metrics["eval_"+sources[si].domain+"_dedicated_query_count"]=float64(sd.q)
		}
		if so.q>0 {
			ca:=float64(so.context)/float64(so.q)
			ra:=float64(so.current)/float64(so.q)
			metrics["eval_"+sources[si].domain+"_overflow_context_top1_accuracy"]=ca
			metrics["eval_"+sources[si].domain+"_overflow_current_top1_accuracy"]=ra
			metrics["eval_"+sources[si].domain+"_overflow_context_minus_current_accuracy"]=ca-ra
			metrics["eval_"+sources[si].domain+"_overflow_query_count"]=float64(so.q)
		}
	}

	metrics["dedicated_query_count"]=float64(dedicatedT.q)
	metrics["overflow_query_count"]=float64(overflowT.q)
	if dedicatedT.q>0 {
		ca:=float64(dedicatedT.context)/float64(dedicatedT.q)
		ra:=float64(dedicatedT.current)/float64(dedicatedT.q)
		metrics["pooled_dedicated_context_top1_accuracy"]=ca
		metrics["pooled_dedicated_current_top1_accuracy"]=ra
		metrics["pooled_dedicated_context_minus_current_accuracy"]=ca-ra
	}
	if overflowT.q>0 {
		ca:=float64(overflowT.context)/float64(overflowT.q)
		ra:=float64(overflowT.current)/float64(overflowT.q)
		metrics["pooled_overflow_context_top1_accuracy"]=ca
		metrics["pooled_overflow_current_top1_accuracy"]=ra
		metrics["pooled_overflow_context_minus_current_accuracy"]=ca-ra
	}
	if totalQueries>0 {
		metrics["evaluation_dedicated_context_fraction"]=float64(dedicatedT.q)/float64(totalQueries)
		if dedicatedT.q+overflowT.q!=totalQueries { metrics["partition_query_count_mismatch"]++ }
	} else { metrics["invalid_row_count"]++ }
	for name,t:=range supportTallies {
		metrics[name+"_query_count"]=float64(t.q)
		if t.q>0 {
			ca:=float64(t.context)/float64(t.q)
			ra:=float64(t.current)/float64(t.q)
			metrics[name+"_context_top1_accuracy"]=ca
			metrics[name+"_current_top1_accuracy"]=ra
			metrics[name+"_context_minus_current_accuracy"]=ca-ra
		}
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++}
	}
	return wlmLmExternalContextStateAttributionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-CONTEXT-STATE-ATTRIBUTION-R1",
		Metrics:metrics,
	}
}
