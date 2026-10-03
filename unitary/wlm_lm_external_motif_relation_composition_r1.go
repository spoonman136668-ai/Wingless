package unitary

import "math"

type wlmLmExternalMotifRelationCompositionR1Source struct {
	domain string
	data []byte
	sha256 string
	bytes int
}

type wlmLmExternalMotifRelationCompositionR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalMotifRelationCompositionR1Argmax(row *[512]uint32, fallback int)(int,bool) {
	best,total:=wlmLmExternalMotifRelationHoldoutR1Best(row)
	if total==0 { return fallback,true }
	return best,false
}

func RunWlmLmExternalMotifRelationCompositionR1(code,structured,prose []byte) interface{} {
	sources:=[]wlmLmExternalMotifRelationCompositionR1Source{
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
		"depth2_sources_with_positive_composed_gain":0,
		"depth3_sources_with_positive_composed_gain":0,
		"depth2_pooled_composed_accuracy_gain":0,
		"depth3_pooled_composed_accuracy_gain":0,
		"maximum_chain_fallback_fraction":0,
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
	var direct2 [512][512]uint32
	var direct3 [512][512]uint32
	for _,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if relation[a][b]==^uint32(0) || unigram[b]==^uint32(0) {
				metrics["counter_overflow_count"]++
			} else {
				relation[a][b]++
				unigram[b]++
			}
		}
		for j:=0;j+2<len(stream);j++ {
			a,b:=stream[j],stream[j+2]
			if direct2[a][b]==^uint32(0) { metrics["counter_overflow_count"]++ } else { direct2[a][b]++ }
		}
		for j:=0;j+3<len(stream);j++ {
			a,b:=stream[j],stream[j+3]
			if direct3[a][b]==^uint32(0) { metrics["counter_overflow_count"]++ } else { direct3[a][b]++ }
		}
	}
	globalBest:=0
	for i:=1;i<512;i++ { if unigram[i]>unigram[globalBest] { globalBest=i } }

	pooledCompCorrect:=[4]int{}
	pooledBaseCorrect:=[4]int{}
	pooledDirectCorrect:=[4]int{}
	pooledQueries:=[4]int{}

	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		for _,depth:=range []int{2,3} {
			compCorrect:=0
			baseCorrect:=0
			directCorrect:=0
			queries:=0
			fallbackQueries:=0
			for j:=0;j+depth<len(stream);j++ {
				current:=stream[j]
				target:=stream[j+depth]
				pred:=current
				hadFallback:=false
				for step:=0;step<depth;step++ {
					var fallback bool
					pred,fallback=wlmLmExternalMotifRelationCompositionR1Argmax(&relation[pred],globalBest)
					hadFallback=hadFallback||fallback
				}
				var directPred int
				if depth==2 {
					directPred,_=wlmLmExternalMotifRelationCompositionR1Argmax(&direct2[current],globalBest)
				} else {
					directPred,_=wlmLmExternalMotifRelationCompositionR1Argmax(&direct3[current],globalBest)
				}
				if pred==target { compCorrect++;pooledCompCorrect[depth]++ }
				if globalBest==target { baseCorrect++;pooledBaseCorrect[depth]++ }
				if directPred==target { directCorrect++;pooledDirectCorrect[depth]++ }
				if hadFallback { fallbackQueries++ }
				queries++
				pooledQueries[depth]++
			}
			if queries==0 { metrics["invalid_row_count"]++; continue }
			compAcc:=float64(compCorrect)/float64(queries)
			baseAcc:=float64(baseCorrect)/float64(queries)
			directAcc:=float64(directCorrect)/float64(queries)
			gain:=compAcc-baseAcc
			fallbackFrac:=float64(fallbackQueries)/float64(queries)
			prefix:="depth2_"
			if depth==3 { prefix="depth3_" }
			metrics[prefix+"eval_"+sources[si].domain+"_composed_top1_accuracy"]=compAcc
			metrics[prefix+"eval_"+sources[si].domain+"_unigram_top1_accuracy"]=baseAcc
			metrics[prefix+"eval_"+sources[si].domain+"_direct_lookup_top1_accuracy"]=directAcc
			metrics[prefix+"eval_"+sources[si].domain+"_composed_accuracy_gain"]=gain
			metrics[prefix+"eval_"+sources[si].domain+"_chain_fallback_fraction"]=fallbackFrac
			if gain>0 { metrics[prefix+"sources_with_positive_composed_gain"]++ }
			metrics["maximum_chain_fallback_fraction"]=wlmLmRawRepPredFreshHoldoutR1Max(metrics["maximum_chain_fallback_fraction"],fallbackFrac)
		}
	}
	for _,depth:=range []int{2,3} {
		if pooledQueries[depth]==0 { metrics["invalid_row_count"]++; continue }
		prefix:="depth2_"
		if depth==3 { prefix="depth3_" }
		metrics[prefix+"query_count"]=float64(pooledQueries[depth])
		compAcc:=float64(pooledCompCorrect[depth])/float64(pooledQueries[depth])
		baseAcc:=float64(pooledBaseCorrect[depth])/float64(pooledQueries[depth])
		directAcc:=float64(pooledDirectCorrect[depth])/float64(pooledQueries[depth])
		metrics[prefix+"pooled_composed_top1_accuracy"]=compAcc
		metrics[prefix+"pooled_unigram_top1_accuracy"]=baseAcc
		metrics[prefix+"pooled_direct_lookup_top1_accuracy"]=directAcc
		metrics[prefix+"pooled_composed_accuracy_gain"]=compAcc-baseAcc
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) { metrics["invalid_row_count"]++ }
	}
	return wlmLmExternalMotifRelationCompositionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-MOTIF-RELATION-COMPOSITION-R1",
		Metrics:metrics,
	}
}
