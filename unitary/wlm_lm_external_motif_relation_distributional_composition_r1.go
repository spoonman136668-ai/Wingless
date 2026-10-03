package unitary

import "math"

type wlmLmExternalMotifRelationDistributionalR1Source struct {
	domain string
	data []byte
	sha256 string
	bytes int
}

type wlmLmExternalMotifRelationDistributionalR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalMotifRelationDistributionalR1Predict(
	start,depth int,
	relation *[512][512]uint32,
	rowTotals *[512]uint64,
	unigram *[512]uint32,
	unigramTotal uint64,
)(int,float64) {
	var dist [512]float64
	dist[start]=1.0
	for step:=0;step<depth;step++ {
		var next [512]float64
		for from,p:=range dist {
			if p==0 { continue }
			if rowTotals[from]>0 {
				den:=float64(rowTotals[from])
				for to,c:=range relation[from] {
					if c>0 { next[to]+=p*float64(c)/den }
				}
			} else if unigramTotal>0 {
				den:=float64(unigramTotal)
				for to,c:=range unigram {
					if c>0 { next[to]+=p*float64(c)/den }
				}
			}
		}
		dist=next
	}
	best:=0
	total:=0.0
	for i,p:=range dist {
		total+=p
		if p>dist[best] { best=i }
	}
	return best,math.Abs(1.0-total)
}

func RunWlmLmExternalMotifRelationDistributionalCompositionR1(code,structured,prose []byte) interface{} {
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
		"depth2_pooled_distributional_accuracy_gain":0,
		"depth3_pooled_distributional_accuracy_gain":0,
		"depth2_distributional_minus_greedy_accuracy":0,
		"depth3_distributional_minus_greedy_accuracy":0,
		"depth2_sources_with_positive_distributional_gain":0,
		"depth3_sources_with_positive_distributional_gain":0,
		"maximum_probability_mass_error":0,
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

	var relation [512][512]uint32
	var rowTotals [512]uint64
	var unigram [512]uint32
	var direct2 [512][512]uint32
	var direct3 [512][512]uint32
	var unigramTotal uint64
	for _,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if relation[a][b]==^uint32(0)||unigram[b]==^uint32(0) { metrics["counter_overflow_count"]++;continue }
			relation[a][b]++;rowTotals[a]++;unigram[b]++;unigramTotal++
		}
		for j:=0;j+2<len(stream);j++ {
			a,b:=stream[j],stream[j+2]
			if direct2[a][b]==^uint32(0){metrics["counter_overflow_count"]++}else{direct2[a][b]++}
		}
		for j:=0;j+3<len(stream);j++ {
			a,b:=stream[j],stream[j+3]
			if direct3[a][b]==^uint32(0){metrics["counter_overflow_count"]++}else{direct3[a][b]++}
		}
	}
	if unigramTotal==0 { metrics["invalid_row_count"]++ }
	globalBest:=0
	for i:=1;i<512;i++ { if unigram[i]>unigram[globalBest] { globalBest=i } }

	var pred2,pred3 [512]int
	for start:=0;start<512;start++ {
		p,e:=wlmLmExternalMotifRelationDistributionalR1Predict(start,2,&relation,&rowTotals,&unigram,unigramTotal)
		pred2[start]=p
		metrics["maximum_probability_mass_error"]=wlmLmRawRepPredFreshHoldoutR1Max(metrics["maximum_probability_mass_error"],e)
		p,e=wlmLmExternalMotifRelationDistributionalR1Predict(start,3,&relation,&rowTotals,&unigram,unigramTotal)
		pred3[start]=p
		metrics["maximum_probability_mass_error"]=wlmLmRawRepPredFreshHoldoutR1Max(metrics["maximum_probability_mass_error"],e)
	}

	var compCorrect,baseCorrect,directCorrect [4]int
	var queries [4]int
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		for _,depth:=range []int{2,3} {
			sc,sb,sd,sq:=0,0,0,0
			for j:=0;j+depth<len(stream);j++ {
				start,target:=stream[j],stream[j+depth]
				pred:=pred2[start]
				drow:=&direct2[start]
				if depth==3 { pred=pred3[start];drow=&direct3[start] }
				dpred,_:=wlmLmExternalMotifRelationCompositionR1Argmax(drow,globalBest)
				if pred==target { sc++;compCorrect[depth]++ }
				if globalBest==target { sb++;baseCorrect[depth]++ }
				if dpred==target { sd++;directCorrect[depth]++ }
				sq++;queries[depth]++
			}
			if sq==0 { metrics["invalid_row_count"]++;continue }
			ca:=float64(sc)/float64(sq)
			ba:=float64(sb)/float64(sq)
			da:=float64(sd)/float64(sq)
			gain:=ca-ba
			prefix:="depth2_";if depth==3{prefix="depth3_"}
			metrics[prefix+"eval_"+sources[si].domain+"_distributional_top1_accuracy"]=ca
			metrics[prefix+"eval_"+sources[si].domain+"_unigram_top1_accuracy"]=ba
			metrics[prefix+"eval_"+sources[si].domain+"_direct_lookup_top1_accuracy"]=da
			metrics[prefix+"eval_"+sources[si].domain+"_distributional_accuracy_gain"]=gain
			if gain>0 { metrics[prefix+"sources_with_positive_distributional_gain"]++ }
		}
	}
	for _,depth:=range []int{2,3} {
		if queries[depth]==0 { metrics["invalid_row_count"]++;continue }
		prefix:="depth2_";greedy:=0.323167469234885
		if depth==3 { prefix="depth3_";greedy=0.2508701472556894 }
		ca:=float64(compCorrect[depth])/float64(queries[depth])
		ba:=float64(baseCorrect[depth])/float64(queries[depth])
		da:=float64(directCorrect[depth])/float64(queries[depth])
		metrics[prefix+"query_count"]=float64(queries[depth])
		metrics[prefix+"pooled_distributional_top1_accuracy"]=ca
		metrics[prefix+"pooled_unigram_top1_accuracy"]=ba
		metrics[prefix+"pooled_direct_lookup_top1_accuracy"]=da
		metrics[prefix+"pooled_distributional_accuracy_gain"]=ca-ba
		metrics[prefix+"distributional_minus_greedy_accuracy"]=ca-greedy
	}
	if metrics["maximum_probability_mass_error"]>1e-9 { metrics["invalid_row_count"]++ }
	for _,v:=range metrics { if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++} }
	return wlmLmExternalMotifRelationDistributionalR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-MOTIF-RELATION-DISTRIBUTIONAL-COMPOSITION-R1",
		Metrics:metrics,
	}
}
