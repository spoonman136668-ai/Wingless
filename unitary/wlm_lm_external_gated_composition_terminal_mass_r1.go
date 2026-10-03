package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalGatedCompositionTerminalMassR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmExternalGatedCompositionTerminalMassR1PairMass struct {
	prev int
	curr int
	mass float64
}

func wlmLmExternalGatedCompositionTerminalMassR1Predict(
	prev,curr,depth int,
	relation *[512][512]uint32,
	rowTotals *[512]uint64,
	exact *[512][512]uint32,
	exactTotals *[512]uint64,
	exactID map[wlmLmExternalContextStateHoldoutR1Pair]int,
	pairCounts map[wlmLmExternalContextStateHoldoutR1Pair]uint32,
	rowEligible *[512]bool,
	rowMargin *[512]float64,
	median float64,
	unigram *[512]uint32,
	unigramTotal uint64,
)(int,float64,float64) {
	active:=[]wlmLmExternalGatedCompositionTerminalMassR1PairMass{{prev:prev,curr:curr,mass:1.0}}
	minRetained:=1.0
	maxMassError:=0.0
	for step:=0;step<depth;step++ {
		next:=make(map[[2]int]float64)
		for _,state:=range active {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:state.prev,curr:state.curr}
			id,dedicated:=exactID[p]
			useExact:=dedicated && pairCounts[p]>=3 && rowEligible[state.curr] && rowMargin[state.curr]<=median
			if useExact && exactTotals[id]>0 {
				den:=float64(exactTotals[id])
				for to:=0;to<512;to++ {
					c:=exact[id][to]
					if c>0 { next[[2]int{state.curr,to}]+=state.mass*float64(c)/den }
				}
			} else if rowTotals[state.curr]>0 {
				den:=float64(rowTotals[state.curr])
				for to:=0;to<512;to++ {
					c:=relation[state.curr][to]
					if c>0 { next[[2]int{state.curr,to}]+=state.mass*float64(c)/den }
				}
			} else if unigramTotal>0 {
				den:=float64(unigramTotal)
				for to:=0;to<512;to++ {
					c:=unigram[to]
					if c>0 { next[[2]int{state.curr,to}]+=state.mass*float64(c)/den }
				}
			}
		}

		entries:=make([]wlmLmExternalGatedCompositionTerminalMassR1PairMass,0,len(next))
		for p,m:=range next { entries=append(entries,wlmLmExternalGatedCompositionTerminalMassR1PairMass{prev:p[0],curr:p[1],mass:m}) }
		sort.Slice(entries,func(i,j int)bool{
			if entries[i].prev!=entries[j].prev { return entries[i].prev<entries[j].prev }
			return entries[i].curr<entries[j].curr
		})
		total:=0.0
		for _,e:=range entries { total+=e.mass }
		if err:=math.Abs(1.0-total);err>maxMassError { maxMassError=err }
		if total==0 { return 0,maxMassError,0 }

		if step==depth-1 {
			var final [512]float64
			for _,e:=range entries { final[e.curr]+=e.mass }
			best:=0
			finalTotal:=0.0
			for i,m:=range final {
				finalTotal+=m
				if m>final[best] { best=i }
			}
			if err:=math.Abs(1.0-finalTotal);err>maxMassError { maxMassError=err }
			return best,maxMassError,minRetained
		}

		sort.Slice(entries,func(i,j int)bool{
			if entries[i].mass!=entries[j].mass { return entries[i].mass>entries[j].mass }
			if entries[i].prev!=entries[j].prev { return entries[i].prev<entries[j].prev }
			return entries[i].curr<entries[j].curr
		})
		if len(entries)>512 { entries=entries[:512] }
		kept:=0.0
		for _,e:=range entries { kept+=e.mass }
		retained:=kept/total
		if retained<minRetained { minRetained=retained }
		if kept==0 { return 0,maxMassError,0 }
		for i:=range entries { entries[i].mass/=kept }
		active=entries
	}
	return 0,maxMassError,minRetained
}

func RunWlmLmExternalGatedCompositionTerminalMassR1(code,structured,prose []byte) interface{} {
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
		"intermediate_runtime_pair_state_cap":512,
		"terminal_prune_event_count":0,
		"depth2_query_count":0,
		"depth3_query_count":0,
		"depth2_pooled_gated_minus_first_order_accuracy":0,
		"depth3_pooled_gated_minus_first_order_accuracy":0,
		"depth2_minimum_source_gated_minus_first_order_accuracy":math.Inf(1),
		"depth3_minimum_source_gated_minus_first_order_accuracy":math.Inf(1),
		"minimum_intermediate_beam_retained_mass":1,
		"maximum_probability_mass_error":0,
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
		train=append(train,src.data[:split]); evals=append(evals,src.data[split:])
	}
	_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(train,metrics)
	metrics["selected_motif_count"]=float64(len(selected))
	keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index:=make(map[[4]uint8]int,len(keys));for i,k:=range keys{index[k]=i}

	streams:=make([][]int,len(train))
	var relation [512][512]uint32
	var rowTotals [512]uint64
	var unigram [512]uint32
	var unigramTotal uint64
	pairCounts:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index);streams[si]=stream
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
	dedicated:=511;if len(rows)<dedicated{dedicated=len(rows)}
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,dedicated)
	for i:=0;i<dedicated;i++ { exactID[rows[i].pair]=i }
	var exact [512][512]uint32
	var exactTotals [512]uint64
	for _,stream:=range streams {
		for j:=1;j+1<len(stream);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			id:=511;if v,ok:=exactID[p];ok{id=v}
			target:=stream[j+1]
			if exact[id][target]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			exact[id][target]++;exactTotals[id]++
		}
	}

	var rowEligible [512]bool
	var rowMargin [512]float64
	margins:=make([]float64,0,512)
	for i:=0;i<512;i++ {
		var total uint64;var top1,top2 uint32
		for j:=0;j<512;j++ {
			v:=relation[i][j];total+=uint64(v)
			if v>top1 {top2=top1;top1=v} else if v>top2 {top2=v}
		}
		if total>=4 {
			rowEligible[i]=true
			rowMargin[i]=(float64(top1)-float64(top2))/float64(total)
			margins=append(margins,rowMargin[i])
		}
	}
	sort.Float64s(margins)
	median:=0.0
	if len(margins)==0 {metrics["invalid_row_count"]++} else if len(margins)%2==1 {median=margins[len(margins)/2]} else {median=(margins[len(margins)/2-1]+margins[len(margins)/2])/2}
	metrics["ambiguity_median_margin"]=median

	var gatedCorrect,firstCorrect [4]int
	var queries [4]int
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		for _,depth:=range []int{2,3} {
			sg,sf,sq:=0,0,0
			for j:=0;j+1+depth<len(stream);j++ {
				prev,curr,target:=stream[j],stream[j+1],stream[j+1+depth]
				gpred,merr,retained:=wlmLmExternalGatedCompositionTerminalMassR1Predict(prev,curr,depth,&relation,&rowTotals,&exact,&exactTotals,exactID,pairCounts,&rowEligible,&rowMargin,median,&unigram,unigramTotal)
				fpred,ferr:=wlmLmExternalMotifRelationDistributionalR1Predict(curr,depth,&relation,&rowTotals,&unigram,unigramTotal)
				if merr>metrics["maximum_probability_mass_error"]{metrics["maximum_probability_mass_error"]=merr}
				if ferr>metrics["maximum_probability_mass_error"]{metrics["maximum_probability_mass_error"]=ferr}
				if retained<metrics["minimum_intermediate_beam_retained_mass"]{metrics["minimum_intermediate_beam_retained_mass"]=retained}
				if gpred==target{sg++;gatedCorrect[depth]++}
				if fpred==target{sf++;firstCorrect[depth]++}
				sq++;queries[depth]++
			}
			if sq==0 {metrics["invalid_row_count"]++;continue}
			ga:=float64(sg)/float64(sq);fa:=float64(sf)/float64(sq);delta:=ga-fa
			prefix:="depth2_";if depth==3{prefix="depth3_"}
			metrics[prefix+"eval_"+sources[si].domain+"_gated_top1_accuracy"]=ga
			metrics[prefix+"eval_"+sources[si].domain+"_first_order_top1_accuracy"]=fa
			metrics[prefix+"eval_"+sources[si].domain+"_gated_minus_first_order_accuracy"]=delta
			k:=prefix+"minimum_source_gated_minus_first_order_accuracy"
			if delta<metrics[k]{metrics[k]=delta}
		}
	}
	for _,depth:=range []int{2,3} {
		prefix:="depth2_";if depth==3{prefix="depth3_"}
		if queries[depth]==0 {metrics["invalid_row_count"]++;continue}
		ga:=float64(gatedCorrect[depth])/float64(queries[depth]);fa:=float64(firstCorrect[depth])/float64(queries[depth])
		metrics[prefix+"query_count"]=float64(queries[depth])
		metrics[prefix+"pooled_gated_top1_accuracy"]=ga
		metrics[prefix+"pooled_first_order_top1_accuracy"]=fa
		metrics[prefix+"pooled_gated_minus_first_order_accuracy"]=ga-fa
		if math.IsInf(metrics[prefix+"minimum_source_gated_minus_first_order_accuracy"],1){metrics[prefix+"minimum_source_gated_minus_first_order_accuracy"]=0;metrics["invalid_row_count"]++}
	}
	if metrics["maximum_probability_mass_error"]>1e-9{metrics["invalid_row_count"]++}
	for _,v:=range metrics{if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++}}
	return wlmLmExternalGatedCompositionTerminalMassR1Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-EXTERNAL-GATED-COMPOSITION-TERMINAL-MASS-R1",Metrics:metrics}
}
