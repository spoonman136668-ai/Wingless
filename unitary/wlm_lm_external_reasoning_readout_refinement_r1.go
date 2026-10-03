package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalReasoningReadoutRefinementR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalReasoningReadoutRefinementR1Top1(dist *[512]float64) int {
	best:=0
	for i:=1;i<512;i++ {
		if dist[i]>dist[best] { best=i }
	}
	return best
}

func wlmLmExternalReasoningReadoutRefinementR1Class(dist *[512]float64, classes *[512]int) int {
	var mass [8]float64
	for i,p:=range dist { mass[classes[i]]+=p }
	best:=0
	for i:=1;i<8;i++ { if mass[i]>mass[best] { best=i } }
	return best
}

func RunWlmLmExternalReasoningReadoutRefinementR1(code,structured,prose []byte) interface{} {
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
		"refined_reasoning_class_count":8,
		"motif_class_assignment_count":512,
		"minimum_refined_class_size":512,
		"maximum_refined_class_size":0,
		"depth2_query_count":0,
		"depth3_query_count":0,
		"depth2_refined_baseline_reasoning_accuracy":0,
		"depth2_refined_repaired_reasoning_accuracy":0,
		"depth2_refined_repaired_top1_projection_accuracy":0,
		"depth2_refined_repaired_reasoning_minus_baseline_accuracy":0,
		"depth2_refined_repaired_reasoning_minus_top1_projection_accuracy":0,
		"depth3_refined_baseline_reasoning_accuracy":0,
		"depth3_refined_repaired_reasoning_accuracy":0,
		"depth3_refined_repaired_top1_projection_accuracy":0,
		"depth3_refined_repaired_reasoning_minus_baseline_accuracy":0,
		"depth3_refined_repaired_reasoning_minus_top1_projection_accuracy":0,
		"depth2_baseline_exact_top1_accuracy":0,
		"depth2_repaired_exact_top1_accuracy":0,
		"depth2_repaired_exact_top1_minus_baseline_exact_top1_accuracy":0,
		"depth3_baseline_exact_top1_accuracy":0,
		"depth3_repaired_exact_top1_accuracy":0,
		"depth3_repaired_exact_top1_minus_baseline_exact_top1_accuracy":0,
		"depth3_exact_repair_gain_query_count":0,
		"depth3_exact_gain_hidden_by_refined_readout_count":0,
		"depth3_exact_gain_hidden_by_refined_readout_fraction":0,
		"minimum_source_depth3_refined_repaired_reasoning_minus_baseline_accuracy":math.Inf(1),
		"maximum_probability_mass_error":0,
		"minimum_repaired_retained_mass":1,
		"heldout_selection_use_count":0,
		"relation_capacity_growth_event_count":0,
		"adaptive_readout_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	train:=make([][]byte,0,len(sources));evals:=make([][]byte,0,len(sources))
	for _,src:=range sources {
		metrics["total_source_bytes"]+=float64(len(src.data))
		if len(src.data)!=src.bytes||wlmLmExternalRawRepPredR1SHA256(src.data)!=src.sha256{metrics["source_identity_mismatch_count"]++}
		split:=len(src.data)*3/5
		if split<5||len(src.data)-split<5{metrics["invalid_row_count"]++}
		train=append(train,src.data[:split]);evals=append(evals,src.data[split:])
	}
	_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(train,metrics)
	metrics["selected_motif_count"]=float64(len(selected))
	keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index:=make(map[[4]uint8]int,len(keys));for i,k:=range keys{index[k]=i}

	streams:=make([][]int,len(train))
	var trainMotifCounts [512]uint64
	var relation [512][512]uint32
	var rowTotals [512]uint64
	var unigram [512]uint32
	var unigramTotal uint64
	pairCounts:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si,data:=range train {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index);streams[si]=stream
		for _,m:=range stream{trainMotifCounts[m]++}
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
	if unigramTotal==0{metrics["invalid_row_count"]++}

	ids:=make([]int,512);for i:=0;i<512;i++{ids[i]=i}
	sort.Slice(ids,func(i,j int)bool{
		if trainMotifCounts[ids[i]]!=trainMotifCounts[ids[j]]{return trainMotifCounts[ids[i]]>trainMotifCounts[ids[j]]}
		return ids[i]<ids[j]
	})
	var classes [512]int
	var classSizes [8]int
	for rank,id:=range ids {
		classID:=rank/64
		classes[id]=classID
		classSizes[classID]++
	}
	for _,n:=range classSizes {
		if float64(n)<metrics["minimum_refined_class_size"]{metrics["minimum_refined_class_size"]=float64(n)}
		if float64(n)>metrics["maximum_refined_class_size"]{metrics["maximum_refined_class_size"]=float64(n)}
	}

	rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairCounts))
	for p,c:=range pairCounts{rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:p,count:c})}
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count{return rows[i].count>rows[j].count}
		if rows[i].pair.prev!=rows[j].pair.prev{return rows[i].pair.prev<rows[j].pair.prev}
		return rows[i].pair.curr<rows[j].pair.curr
	})
	dedicated:=511;if len(rows)<dedicated{dedicated=len(rows)}
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,dedicated)
	for i:=0;i<dedicated;i++{exactID[rows[i].pair]=i}
	var exact [512][512]uint32
	var exactTotals [512]uint64
	for _,stream:=range streams{
		for j:=1;j+1<len(stream);j++{
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
	for i:=0;i<512;i++{
		var total uint64;var top1,top2 uint32
		for _,v:=range relation[i]{
			total+=uint64(v)
			if v>top1{top2=top1;top1=v}else if v>top2{top2=v}
		}
		if total>=4{
			rowEligible[i]=true
			rowMargin[i]=(float64(top1)-float64(top2))/float64(total)
			margins=append(margins,rowMargin[i])
		}
	}
	sort.Float64s(margins)
	median:=0.0
	if len(margins)==0{metrics["invalid_row_count"]++}else if len(margins)%2==1{median=margins[len(margins)/2]}else{median=(margins[len(margins)/2-1]+margins[len(margins)/2])/2}

	var baseHits,repairHits,top1Hits [4]int
	var baseExactHits,repairExactHits [4]int
	var queries [4]int
	for si,data:=range evals {
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		sourceBase,sourceRepair,sourceQueries:=0,0,0
		for _,depth:=range []int{2,3}{
			for j:=0;j+1+depth<len(stream);j++{
				prev,curr,target:=stream[j],stream[j+1],stream[j+1+depth]
				bd,be,_:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(prev,curr,depth,2048,false,&relation,&rowTotals,&exact,&exactTotals,exactID,pairCounts,&rowEligible,&rowMargin,median,&unigram,unigramTotal)
				rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(prev,curr,depth,2048,true,&relation,&rowTotals,&exact,&exactTotals,exactID,pairCounts,&rowEligible,&rowMargin,median,&unigram,unigramTotal)
				if be>metrics["maximum_probability_mass_error"]{metrics["maximum_probability_mass_error"]=be}
				if re>metrics["maximum_probability_mass_error"]{metrics["maximum_probability_mass_error"]=re}
				if rr<metrics["minimum_repaired_retained_mass"]{metrics["minimum_repaired_retained_mass"]=rr}
				truth:=classes[target]
				btop:=wlmLmExternalReasoningReadoutRefinementR1Top1(&bd)
				rtop:=wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
				bc:=wlmLmExternalReasoningReadoutRefinementR1Class(&bd,&classes)
				rc:=wlmLmExternalReasoningReadoutRefinementR1Class(&rd,&classes)
				tc:=classes[rtop]
				if btop==target{baseExactHits[depth]++}
				if rtop==target{repairExactHits[depth]++}
				if depth==3 && btop!=target && rtop==target{
					metrics["depth3_exact_repair_gain_query_count"]++
					if classes[btop]==classes[rtop] || bc==truth{
						metrics["depth3_exact_gain_hidden_by_refined_readout_count"]++
					}
				}
				if bc==truth{baseHits[depth]++}
				if rc==truth{repairHits[depth]++}
				if tc==truth{top1Hits[depth]++}
				queries[depth]++
				if depth==3{
					sourceQueries++
					if bc==truth{sourceBase++}
					if rc==truth{sourceRepair++}
				}
			}
		}
		if sourceQueries>0{
			gain:=float64(sourceRepair-sourceBase)/float64(sourceQueries)
			metrics["eval_"+sources[si].domain+"_depth3_refined_repaired_reasoning_minus_baseline_accuracy"]=gain
			if gain<metrics["minimum_source_depth3_refined_repaired_reasoning_minus_baseline_accuracy"]{metrics["minimum_source_depth3_refined_repaired_reasoning_minus_baseline_accuracy"]=gain}
		}
	}
	for _,depth:=range []int{2,3}{
		if queries[depth]==0{metrics["invalid_row_count"]++;continue}
		p:="depth2_";if depth==3{p="depth3_"}
		b:=float64(baseHits[depth])/float64(queries[depth])
		r:=float64(repairHits[depth])/float64(queries[depth])
		t:=float64(top1Hits[depth])/float64(queries[depth])
		metrics[p+"query_count"]=float64(queries[depth])
		metrics[p+"refined_baseline_reasoning_accuracy"]=b
		metrics[p+"refined_repaired_reasoning_accuracy"]=r
		metrics[p+"refined_repaired_top1_projection_accuracy"]=t
		metrics[p+"refined_repaired_reasoning_minus_baseline_accuracy"]=r-b
		metrics[p+"refined_repaired_reasoning_minus_top1_projection_accuracy"]=r-t
		beExact:=float64(baseExactHits[depth])/float64(queries[depth])
		reExact:=float64(repairExactHits[depth])/float64(queries[depth])
		metrics[p+"baseline_exact_top1_accuracy"]=beExact
		metrics[p+"repaired_exact_top1_accuracy"]=reExact
		metrics[p+"repaired_exact_top1_minus_baseline_exact_top1_accuracy"]=reExact-beExact
	}
	if metrics["depth3_exact_repair_gain_query_count"]>0{
		metrics["depth3_exact_gain_hidden_by_refined_readout_fraction"]=metrics["depth3_exact_gain_hidden_by_refined_readout_count"]/metrics["depth3_exact_repair_gain_query_count"]
	}
	if math.IsInf(metrics["minimum_source_depth3_refined_repaired_reasoning_minus_baseline_accuracy"],1){metrics["minimum_source_depth3_refined_repaired_reasoning_minus_baseline_accuracy"]=0;metrics["invalid_row_count"]++}
	if metrics["maximum_probability_mass_error"]>1e-9{metrics["invalid_row_count"]++}
	for _,v:=range metrics{if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++}}

	return wlmLmExternalReasoningReadoutRefinementR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-REASONING-READOUT-REFINEMENT-R1",
		Metrics:metrics,
	}
}
