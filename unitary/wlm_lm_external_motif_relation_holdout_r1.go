package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalMotifRelationHoldoutR1Source struct {
	domain string
	data []byte
	sha256 string
	bytes int
}

type wlmLmExternalMotifRelationHoldoutR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif) [][4]uint8 {
	keys:=make([][4]uint8,0,len(selected))
	for k:=range selected { keys=append(keys,k) }
	sort.Slice(keys,func(i,j int)bool{return wlmLmRawRepPredFreshHoldoutR1Less(keys[i],keys[j])})
	return keys
}

func wlmLmExternalMotifRelationHoldoutR1Decode(data []byte,index map[[4]uint8]int) []int {
	out:=make([]int,0,len(data)/2)
	for i:=0;i<len(data); {
		if i+4<=len(data) {
			k:=[4]uint8{data[i],data[i+1],data[i+2],data[i+3]}
			if id,ok:=index[k];ok {
				out=append(out,id)
				i+=4
				continue
			}
		}
		i++
	}
	return out
}

func wlmLmExternalMotifRelationHoldoutR1Best(row *[512]uint32)(int,uint64) {
	best:=0
	bestCount:=row[0]
	total:=uint64(row[0])
	for i:=1;i<512;i++ {
		total+=uint64(row[i])
		if row[i]>bestCount {
			best=i
			bestCount=row[i]
		}
	}
	return best,total
}

func RunWlmLmExternalMotifRelationHoldoutR1(code,structured,prose []byte) interface{} {
	sources:=[]wlmLmExternalMotifRelationHoldoutR1Source{
		{domain:"code",data:code,sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
		{domain:"structured",data:structured,sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
		{domain:"technical_prose",data:prose,sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
	}
	metrics:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":float64(len(sources)),
		"total_source_bytes":0,
		"selected_motif_count":0,
		"train_relation_transition_count":0,
		"eval_relation_transition_count":0,
		"eval_source_count_with_transitions":0,
		"sources_with_positive_relation_accuracy_gain":0,
		"pooled_relation_top1_accuracy_gain":0,
		"minimum_source_relation_accuracy_gain":1,
		"unseen_current_motif_query_fraction":0,
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

	var relation [512][512]uint32
	var unigram [512]uint32
	for i,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		metrics["train_"+sources[i].domain+"_decoded_motif_count"]=float64(len(stream))
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if relation[a][b]==^uint32(0) || unigram[b]==^uint32(0) {
				metrics["counter_overflow_count"]++
				continue
			}
			relation[a][b]++
			unigram[b]++
			metrics["train_relation_transition_count"]++
		}
	}
	globalBest:=0
	for i:=1;i<512;i++ {
		if unigram[i]>unigram[globalBest] { globalBest=i }
	}

	pooledRelationCorrect:=0
	pooledBaselineCorrect:=0
	pooledQueries:=0
	unseenQueries:=0
	for i,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		metrics["eval_"+sources[i].domain+"_decoded_motif_count"]=float64(len(stream))
		if len(stream)<2 { continue }
		metrics["eval_source_count_with_transitions"]++
		relationCorrect:=0
		baselineCorrect:=0
		queries:=0
		sourceUnseen:=0
		for j:=0;j+1<len(stream);j++ {
			a,target:=stream[j],stream[j+1]
			pred,total:=wlmLmExternalMotifRelationHoldoutR1Best(&relation[a])
			if total==0 {
				pred=globalBest
				sourceUnseen++
				unseenQueries++
			}
			if pred==target { relationCorrect++;pooledRelationCorrect++ }
			if globalBest==target { baselineCorrect++;pooledBaselineCorrect++ }
			queries++
			pooledQueries++
			metrics["eval_relation_transition_count"]++
		}
		relationAcc:=float64(relationCorrect)/float64(queries)
		baselineAcc:=float64(baselineCorrect)/float64(queries)
		gain:=relationAcc-baselineAcc
		metrics["eval_"+sources[i].domain+"_relation_transition_count"]=float64(queries)
		metrics["eval_"+sources[i].domain+"_conditional_top1_accuracy"]=relationAcc
		metrics["eval_"+sources[i].domain+"_unigram_top1_accuracy"]=baselineAcc
		metrics["eval_"+sources[i].domain+"_relation_accuracy_gain"]=gain
		metrics["eval_"+sources[i].domain+"_unseen_current_motif_query_fraction"]=float64(sourceUnseen)/float64(queries)
		if gain>0 { metrics["sources_with_positive_relation_accuracy_gain"]++ }
		metrics["minimum_source_relation_accuracy_gain"]=wlmLmRawRepPredFreshHoldoutR1Min(metrics["minimum_source_relation_accuracy_gain"],gain)
	}
	if pooledQueries>0 {
		metrics["pooled_conditional_top1_accuracy"]=float64(pooledRelationCorrect)/float64(pooledQueries)
		metrics["pooled_unigram_top1_accuracy"]=float64(pooledBaselineCorrect)/float64(pooledQueries)
		metrics["pooled_relation_top1_accuracy_gain"]=metrics["pooled_conditional_top1_accuracy"]-metrics["pooled_unigram_top1_accuracy"]
		metrics["unseen_current_motif_query_fraction"]=float64(unseenQueries)/float64(pooledQueries)
	} else {
		metrics["invalid_row_count"]++
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) { metrics["invalid_row_count"]++ }
	}
	return wlmLmExternalMotifRelationHoldoutR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-MOTIF-RELATION-HOLDOUT-R1",
		Metrics:metrics,
	}
}
