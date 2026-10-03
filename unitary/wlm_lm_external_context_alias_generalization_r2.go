package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalContextAliasGeneralizationR2Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalContextAliasGeneralizationR2Distance(a,b *[512]uint32,ta,tb uint64) float64 {
	if ta==0 && tb==0 { return 0 }
	if ta==0 || tb==0 { return 1 }
	var d float64
	da:=float64(ta)
	db:=float64(tb)
	for i:=0;i<512;i++ {
		d += math.Abs(float64(a[i])/da-float64(b[i])/db)
	}
	return d
}

func wlmLmExternalContextAliasGeneralizationR2Aliases(rows *[512][512]uint32, totals *[512]uint64, k int) ([]int,int) {
	medoids:=make([]int,0,k)
	first:=0
	for i:=1;i<512;i++ {
		if totals[i]>totals[first] || (totals[i]==totals[first] && i<first) { first=i }
	}
	medoids=append(medoids,first)
	chosen:=map[int]bool{first:true}
	for len(medoids)<k {
		best:=-1
		bestDist:=-1.0
		for i:=0;i<512;i++ {
			if chosen[i] { continue }
			minDist:=math.Inf(1)
			for _,m:=range medoids {
				d:=wlmLmExternalContextAliasGeneralizationR2Distance(&rows[i],&rows[m],totals[i],totals[m])
				if d<minDist { minDist=d }
			}
			if best<0 || minDist>bestDist || (minDist==bestDist && (totals[i]>totals[best] || (totals[i]==totals[best] && i<best))) {
				best=i
				bestDist=minDist
			}
		}
		if best<0 { break }
		chosen[best]=true
		medoids=append(medoids,best)
	}
	aliases:=make([]int,512)
	pop:=make([]int,len(medoids))
	for i:=0;i<512;i++ {
		bestAlias:=0
		bestDist:=wlmLmExternalContextAliasGeneralizationR2Distance(&rows[i],&rows[medoids[0]],totals[i],totals[medoids[0]])
		for a:=1;a<len(medoids);a++ {
			d:=wlmLmExternalContextAliasGeneralizationR2Distance(&rows[i],&rows[medoids[a]],totals[i],totals[medoids[a]])
			if d<bestDist {
				bestAlias=a
				bestDist=d
			}
		}
		aliases[i]=bestAlias
		pop[bestAlias]++
	}
	nonempty:=0
	for _,n:=range pop { if n>0 { nonempty++ } }
	return aliases,nonempty
}

func RunWlmLmExternalContextAliasGeneralizationR2(code,structured,prose []byte) interface{} {
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
		"alias_state_count":512,
		"alias_predictive_counter_capacity":262144,
		"nonempty_previous_alias_count":0,
		"nonempty_current_alias_count":0,
		"previous_role_zero_vector_count":0,
		"current_role_zero_vector_count":0,
		"overflow_unseen_query_count":0,
		"pooled_alias_top1_accuracy":0,
		"pooled_current_top1_accuracy":0,
		"pooled_support_backoff_top1_accuracy":0,
		"pooled_alias_minus_current_accuracy":0,
		"minimum_source_alias_minus_current_accuracy":math.Inf(1),
		"overflow_unseen_alias_top1_accuracy":0,
		"overflow_unseen_exact_backoff_top1_accuracy":0,
		"overflow_unseen_alias_minus_exact_backoff_accuracy":0,
		"alias_state_training_occupancy_fraction":0,
		"alias_state_evaluation_coverage_fraction":0,
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
	var outgoing [512][512]uint32
	var incoming [512][512]uint32
	var outgoingTotals [512]uint64
	var incomingTotals [512]uint64
	pairCounts:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		streams[si]=stream
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if outgoing[a][b]==^uint32(0) || incoming[b][a]==^uint32(0) {
				metrics["counter_overflow_count"]++
				continue
			}
			outgoing[a][b]++
			incoming[b][a]++
			outgoingTotals[a]++
			incomingTotals[b]++
		}
		for j:=1;j<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			if pairCounts[p]==^uint32(0) { metrics["counter_overflow_count"]++; continue }
			pairCounts[p]++
		}
	}
	for i:=0;i<512;i++ {
		if outgoingTotals[i]==0 { metrics["previous_role_zero_vector_count"]++ }
		if incomingTotals[i]==0 { metrics["current_role_zero_vector_count"]++ }
	}
	prevAlias,prevNonempty:=wlmLmExternalContextAliasGeneralizationR2Aliases(&outgoing,&outgoingTotals,16)
	currAlias,currNonempty:=wlmLmExternalContextAliasGeneralizationR2Aliases(&incoming,&incomingTotals,32)
	metrics["nonempty_previous_alias_count"]=float64(prevNonempty)
	metrics["nonempty_current_alias_count"]=float64(currNonempty)

	rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairCounts))
	for p,c:=range pairCounts { rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:p,count:c}) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count { return rows[i].count>rows[j].count }
		if rows[i].pair.prev!=rows[j].pair.prev { return rows[i].pair.prev<rows[j].pair.prev }
		return rows[i].pair.curr<rows[j].pair.curr
	})
	dedicated:=511
	if len(rows)<dedicated { dedicated=len(rows) }
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,dedicated)
	for i:=0;i<dedicated;i++ { exactID[rows[i].pair]=i }

	var aliasCounts [512][512]uint32
	var exactCounts [512][512]uint32
	var current [512][512]uint32
	var unigram [512]uint32
	aliasRowUsed:=make([]bool,512)
	for _,stream:=range streams {
		for j:=0;j+1<len(stream);j++ {
			a,b:=stream[j],stream[j+1]
			if current[a][b]==^uint32(0)||unigram[b]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			current[a][b]++
			unigram[b]++
		}
		for j:=1;j+1<len(stream);j++ {
			prev,curr,target:=stream[j-1],stream[j],stream[j+1]
			state:=prevAlias[prev]*32+currAlias[curr]
			if state<0||state>=512 { metrics["invalid_row_count"]++; continue }
			if aliasCounts[state][target]==^uint32(0){metrics["counter_overflow_count"]++} else {
				aliasCounts[state][target]++
				aliasRowUsed[state]=true
			}
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:prev,curr:curr}
			eid:=511
			if id,ok:=exactID[p];ok { eid=id }
			if exactCounts[eid][target]==^uint32(0){metrics["counter_overflow_count"]++} else { exactCounts[eid][target]++ }
		}
	}
	occupied:=0
	for _,v:=range aliasRowUsed { if v { occupied++ } }
	metrics["alias_state_training_occupancy_fraction"]=float64(occupied)/512.0
	globalBest:=0
	for i:=1;i<512;i++ { if unigram[i]>unigram[globalBest] { globalBest=i } }

	pooledQ:=0
	aliasCorrect:=0
	currentCorrect:=0
	backoffCorrect:=0
	coveredQ:=0
	overflowQ:=0
	overflowAliasCorrect:=0
	overflowBackoffCorrect:=0
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		if len(stream)<3 { continue }
		q:=0
		sa,sc,sb:=0,0,0
		for j:=1;j+1<len(stream);j++ {
			prev,curr,target:=stream[j-1],stream[j],stream[j+1]
			state:=prevAlias[prev]*32+currAlias[curr]
			ap,atotal:=wlmLmExternalMotifRelationHoldoutR1Best(&aliasCounts[state])
			if atotal==0 { ap=globalBest } else { coveredQ++ }
			cp,ctotal:=wlmLmExternalMotifRelationHoldoutR1Best(&current[curr])
			if ctotal==0 { cp=globalBest }

			pair:=wlmLmExternalContextStateHoldoutR1Pair{prev:prev,curr:curr}
			eid,isDedicated:=exactID[pair]
			if !isDedicated { eid=511 }
			ep,etotal:=wlmLmExternalMotifRelationHoldoutR1Best(&exactCounts[eid])
			if etotal==0 { ep=globalBest }
			useExact:=isDedicated && pairCounts[pair]>=3
			bp:=cp
			if useExact { bp=ep }

			if ap==target { sa++;aliasCorrect++ }
			if cp==target { sc++;currentCorrect++ }
			if bp==target { sb++;backoffCorrect++ }
			if !isDedicated {
				overflowQ++
				if ap==target { overflowAliasCorrect++ }
				if bp==target { overflowBackoffCorrect++ }
			}
			q++;pooledQ++
		}
		if q>0 {
			aa:=float64(sa)/float64(q)
			ca:=float64(sc)/float64(q)
			ba:=float64(sb)/float64(q)
			delta:=aa-ca
			metrics["eval_"+sources[si].domain+"_alias_top1_accuracy"]=aa
			metrics["eval_"+sources[si].domain+"_current_top1_accuracy"]=ca
			metrics["eval_"+sources[si].domain+"_support_backoff_top1_accuracy"]=ba
			metrics["eval_"+sources[si].domain+"_alias_minus_current_accuracy"]=delta
			if delta<metrics["minimum_source_alias_minus_current_accuracy"] { metrics["minimum_source_alias_minus_current_accuracy"]=delta }
		}
	}
	metrics["overflow_unseen_query_count"]=float64(overflowQ)
	if pooledQ>0 {
		metrics["pooled_alias_top1_accuracy"]=float64(aliasCorrect)/float64(pooledQ)
		metrics["pooled_current_top1_accuracy"]=float64(currentCorrect)/float64(pooledQ)
		metrics["pooled_support_backoff_top1_accuracy"]=float64(backoffCorrect)/float64(pooledQ)
		metrics["pooled_alias_minus_current_accuracy"]=metrics["pooled_alias_top1_accuracy"]-metrics["pooled_current_top1_accuracy"]
		metrics["alias_state_evaluation_coverage_fraction"]=float64(coveredQ)/float64(pooledQ)
	} else { metrics["invalid_row_count"]++ }
	if overflowQ>0 {
		metrics["overflow_unseen_alias_top1_accuracy"]=float64(overflowAliasCorrect)/float64(overflowQ)
		metrics["overflow_unseen_exact_backoff_top1_accuracy"]=float64(overflowBackoffCorrect)/float64(overflowQ)
		metrics["overflow_unseen_alias_minus_exact_backoff_accuracy"]=metrics["overflow_unseen_alias_top1_accuracy"]-metrics["overflow_unseen_exact_backoff_top1_accuracy"]
	} else { metrics["invalid_row_count"]++ }
	if math.IsInf(metrics["minimum_source_alias_minus_current_accuracy"],1) {
		metrics["minimum_source_alias_minus_current_accuracy"]=0
		metrics["invalid_row_count"]++
	}
	for _,v:=range metrics { if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++} }
	return wlmLmExternalContextAliasGeneralizationR2Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-CONTEXT-ALIAS-GENERALIZATION-R2",
		Metrics:metrics,
	}
}
