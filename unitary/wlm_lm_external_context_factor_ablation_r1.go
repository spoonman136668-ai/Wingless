package unitary

import "math"

type wlmLmExternalContextFactorAblationR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalContextFactorAblationR1(code,structured,prose []byte) interface{} {
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
		"previous_alias_count":16,
		"current_alias_count":32,
		"previous_predictive_counter_capacity":8192,
		"current_alias_predictive_counter_capacity":16384,
		"cartesian_alias_predictive_counter_capacity":262144,
		"previous_factor_training_occupancy_fraction":0,
		"current_alias_factor_training_occupancy_fraction":0,
		"pooled_previous_factor_accuracy":0,
		"pooled_current_alias_factor_accuracy":0,
		"pooled_current_exact_accuracy":0,
		"cartesian_alias_reference_accuracy":0,
		"previous_minus_cartesian_accuracy":0,
		"current_alias_minus_cartesian_accuracy":0,
		"maximum_factor_minus_cartesian_accuracy":0,
		"factor_dominance_gap":0,
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
	var outgoing [512][512]uint32
	var incoming [512][512]uint32
	var outgoingTotals [512]uint64
	var incomingTotals [512]uint64
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		streams[si]=stream
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if outgoing[a][b]==^uint32(0)||incoming[b][a]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			outgoing[a][b]++; incoming[b][a]++; outgoingTotals[a]++; incomingTotals[b]++
		}
	}
	prevAlias,_:=wlmLmExternalContextAliasGeneralizationR2Aliases(&outgoing,&outgoingTotals,16)
	currAlias,_:=wlmLmExternalContextAliasGeneralizationR2Aliases(&incoming,&incomingTotals,32)

	var previous [16][512]uint32
	var currentAlias [32][512]uint32
	var cartesian [512][512]uint32
	var currentExact [512][512]uint32
	var unigram [512]uint32
	prevUsed:=make([]bool,16)
	currUsed:=make([]bool,32)
	for _,stream:=range streams {
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if currentExact[a][b]==^uint32(0)||unigram[b]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			currentExact[a][b]++; unigram[b]++
		}
		for j:=1;j+1<len(stream);j++ {
			prev,curr,target:=stream[j-1],stream[j],stream[j+1]
			pa,ca:=prevAlias[prev],currAlias[curr]
			if previous[pa][target]==^uint32(0)||currentAlias[ca][target]==^uint32(0)||cartesian[pa*32+ca][target]==^uint32(0){
				metrics["counter_overflow_count"]++;continue
			}
			previous[pa][target]++
			currentAlias[ca][target]++
			cartesian[pa*32+ca][target]++
			prevUsed[pa]=true; currUsed[ca]=true
		}
	}
	prevOcc:=0; for _,v:=range prevUsed { if v { prevOcc++ } }
	currOcc:=0; for _,v:=range currUsed { if v { currOcc++ } }
	metrics["previous_factor_training_occupancy_fraction"]=float64(prevOcc)/16.0
	metrics["current_alias_factor_training_occupancy_fraction"]=float64(currOcc)/32.0

	globalBest:=0
	for i:=1;i<512;i++ { if unigram[i]>unigram[globalBest] { globalBest=i } }

	total:=0
	prevHits,currAliasHits,currentHits,cartHits:=0,0,0,0
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		if len(stream)<3 { continue }
		q:=0
		ph,ch,eh,ah:=0,0,0,0
		for j:=1;j+1<len(stream);j++ {
			prev,curr,target:=stream[j-1],stream[j],stream[j+1]
			pa,ca:=prevAlias[prev],currAlias[curr]
			pp,pt:=wlmLmExternalMotifRelationHoldoutR1Best(&previous[pa]); if pt==0 { pp=globalBest }
			cp,ct:=wlmLmExternalMotifRelationHoldoutR1Best(&currentAlias[ca]); if ct==0 { cp=globalBest }
			ep,et:=wlmLmExternalMotifRelationHoldoutR1Best(&currentExact[curr]); if et==0 { ep=globalBest }
			ap,at:=wlmLmExternalMotifRelationHoldoutR1Best(&cartesian[pa*32+ca]); if at==0 { ap=globalBest }
			if pp==target { ph++;prevHits++ }
			if cp==target { ch++;currAliasHits++ }
			if ep==target { eh++;currentHits++ }
			if ap==target { ah++;cartHits++ }
			q++;total++
		}
		if q>0 {
			metrics["eval_"+sources[si].domain+"_previous_factor_accuracy"]=float64(ph)/float64(q)
			metrics["eval_"+sources[si].domain+"_current_alias_factor_accuracy"]=float64(ch)/float64(q)
			metrics["eval_"+sources[si].domain+"_current_exact_accuracy"]=float64(eh)/float64(q)
			metrics["eval_"+sources[si].domain+"_cartesian_alias_accuracy"]=float64(ah)/float64(q)
		}
	}
	if total>0 {
		p:=float64(prevHits)/float64(total)
		c:=float64(currAliasHits)/float64(total)
		e:=float64(currentHits)/float64(total)
		a:=float64(cartHits)/float64(total)
		metrics["pooled_previous_factor_accuracy"]=p
		metrics["pooled_current_alias_factor_accuracy"]=c
		metrics["pooled_current_exact_accuracy"]=e
		metrics["cartesian_alias_reference_accuracy"]=a
		metrics["previous_minus_cartesian_accuracy"]=p-a
		metrics["current_alias_minus_cartesian_accuracy"]=c-a
		if p-a>c-a { metrics["maximum_factor_minus_cartesian_accuracy"]=p-a } else { metrics["maximum_factor_minus_cartesian_accuracy"]=c-a }
		metrics["factor_dominance_gap"]=math.Abs(p-c)
	} else { metrics["invalid_row_count"]++ }
	for _,v:=range metrics { if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++} }
	return wlmLmExternalContextFactorAblationR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-CONTEXT-FACTOR-ABLATION-R1",
		Metrics:metrics,
	}
}
