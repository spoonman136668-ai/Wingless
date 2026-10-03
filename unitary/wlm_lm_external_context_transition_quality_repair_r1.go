package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalContextTransitionQualityRepairR1PairMass struct {
	prev int
	curr int
	mass float64
}

type wlmLmExternalContextTransitionQualityRepairR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalContextTransitionQualityRepairR1Predict(
	prev,curr,depth,capStates int,
	repair bool,
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
	active:=[]wlmLmExternalContextTransitionQualityRepairR1PairMass{{prev:prev,curr:curr,mass:1.0}}
	maxMassError:=0.0
	minRetained:=1.0
	for step:=0;step<depth;step++ {
		next:=make(map[[2]int]float64)
		for _,state:=range active {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:state.prev,curr:state.curr}
			support:=pairCounts[p]
			id,dedicated:=exactID[p]
			if repair && support<3 {
				currentTotal:=rowTotals[state.curr]
				overflowTotal:=exactTotals[511]
				switch {
				case currentTotal>0 && overflowTotal>0:
					cd:=float64(currentTotal)
					od:=float64(overflowTotal)
					for to:=0;to<512;to++ {
						prob:=0.5*float64(relation[state.curr][to])/cd + 0.5*float64(exact[511][to])/od
						if prob>0 { next[[2]int{state.curr,to}]+=state.mass*prob }
					}
				case currentTotal>0:
					cd:=float64(currentTotal)
					for to:=0;to<512;to++ {
						c:=relation[state.curr][to]
						if c>0 { next[[2]int{state.curr,to}]+=state.mass*float64(c)/cd }
					}
				case overflowTotal>0:
					od:=float64(overflowTotal)
					for to:=0;to<512;to++ {
						c:=exact[511][to]
						if c>0 { next[[2]int{state.curr,to}]+=state.mass*float64(c)/od }
					}
				case unigramTotal>0:
					ud:=float64(unigramTotal)
					for to:=0;to<512;to++ {
						c:=unigram[to]
						if c>0 { next[[2]int{state.curr,to}]+=state.mass*float64(c)/ud }
					}
				}
				continue
			}

			useExact:=dedicated && support>=3 && rowEligible[state.curr] && rowMargin[state.curr]<=median
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
		entries:=make([]wlmLmExternalContextTransitionQualityRepairR1PairMass,0,len(next))
		for p,m:=range next {
			entries=append(entries,wlmLmExternalContextTransitionQualityRepairR1PairMass{prev:p[0],curr:p[1],mass:m})
		}
		sort.Slice(entries,func(i,j int)bool{
			if entries[i].mass!=entries[j].mass{return entries[i].mass>entries[j].mass}
			if entries[i].prev!=entries[j].prev{return entries[i].prev<entries[j].prev}
			return entries[i].curr<entries[j].curr
		})
		total:=0.0
		for _,e:=range entries { total+=e.mass }
		if err:=math.Abs(1.0-total);err>maxMassError{maxMassError=err}
		if total==0{return 0,maxMassError,0}
		if len(entries)>capStates{entries=entries[:capStates]}
		kept:=0.0
		for _,e:=range entries{kept+=e.mass}
		if kept==0{return 0,maxMassError,0}
		retained:=kept/total
		if retained<minRetained{minRetained=retained}
		for i:=range entries{entries[i].mass/=kept}
		active=entries
	}
	var final [512]float64
	total:=0.0
	for _,e:=range active{final[e.curr]+=e.mass;total+=e.mass}
	if err:=math.Abs(1.0-total);err>maxMassError{maxMassError=err}
	best:=0
	for i:=1;i<512;i++{if final[i]>final[best]{best=i}}
	return best,maxMassError,minRetained
}

func RunWlmLmExternalContextTransitionQualityRepairR1(code,structured,prose []byte) interface{} {
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
		"depth3_low_support_query_count":0,
		"depth2_baseline_pooled_top1_accuracy":0,
		"depth2_repair_pooled_top1_accuracy":0,
		"depth2_pooled_repair_minus_baseline_accuracy":0,
		"depth3_baseline_pooled_top1_accuracy":0,
		"depth3_repair_pooled_top1_accuracy":0,
		"depth3_pooled_repair_minus_baseline_accuracy":0,
		"depth3_low_support_baseline_top1_accuracy":0,
		"depth3_low_support_repair_top1_accuracy":0,
		"depth3_low_support_repair_minus_baseline_accuracy":0,
		"minimum_source_depth3_repair_minus_baseline_accuracy":math.Inf(1),
		"maximum_probability_mass_error":0,
		"minimum_repair_retained_mass":1,
		"heldout_selection_use_count":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	train:=make([][]byte,0,len(sources))
	evals:=make([][]byte,0,len(sources))
	for _,src:=range sources{
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
	var relation [512][512]uint32
	var rowTotals [512]uint64
	var unigram [512]uint32
	var unigramTotal uint64
	pairCounts:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si,data:=range train{
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index);streams[si]=stream
		for j:=0;j+1<len(stream);j++{
			a,b:=stream[j],stream[j+1]
			if relation[a][b]==^uint32(0)||unigram[b]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			relation[a][b]++;rowTotals[a]++;unigram[b]++;unigramTotal++
		}
		for j:=1;j<len(stream);j++{
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:stream[j-1],curr:stream[j]}
			if pairCounts[p]==^uint32(0){metrics["counter_overflow_count"]++;continue}
			pairCounts[p]++
		}
	}
	if unigramTotal==0{metrics["invalid_row_count"]++}

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
	metrics["ambiguity_median_margin"]=median

	var depthBase,depthRepair [4]int
	var depthQueries [4]int
	lowBase,lowRepair,lowQueries:=0,0,0
	for si,data:=range evals{
		stream:=wlmLmExternalMotifRelationHoldoutR1Decode(data,index)
		sourceBase,sourceRepair,sourceQueries:=0,0,0
		for _,depth:=range []int{2,3}{
			for j:=0;j+1+depth<len(stream);j++{
				prev,curr,target:=stream[j],stream[j+1],stream[j+1+depth]
				bp,bErr,_:=wlmLmExternalContextTransitionQualityRepairR1Predict(prev,curr,depth,2048,false,&relation,&rowTotals,&exact,&exactTotals,exactID,pairCounts,&rowEligible,&rowMargin,median,&unigram,unigramTotal)
				rp,rErr,rRet:=wlmLmExternalContextTransitionQualityRepairR1Predict(prev,curr,depth,2048,true,&relation,&rowTotals,&exact,&exactTotals,exactID,pairCounts,&rowEligible,&rowMargin,median,&unigram,unigramTotal)
				if bErr>metrics["maximum_probability_mass_error"]{metrics["maximum_probability_mass_error"]=bErr}
				if rErr>metrics["maximum_probability_mass_error"]{metrics["maximum_probability_mass_error"]=rErr}
				if rRet<metrics["minimum_repair_retained_mass"]{metrics["minimum_repair_retained_mass"]=rRet}
				if bp==target{depthBase[depth]++}
				if rp==target{depthRepair[depth]++}
				depthQueries[depth]++
				if depth==3{
					sourceQueries++
					if bp==target{sourceBase++}
					if rp==target{sourceRepair++}
					p:=wlmLmExternalContextStateHoldoutR1Pair{prev:prev,curr:curr}
					if pairCounts[p]<3{
						lowQueries++
						if bp==target{lowBase++}
						if rp==target{lowRepair++}
					}
				}
			}
		}
		if sourceQueries>0{
			gain:=float64(sourceRepair-sourceBase)/float64(sourceQueries)
			metrics["eval_"+sources[si].domain+"_depth3_repair_minus_baseline_accuracy"]=gain
			if gain<metrics["minimum_source_depth3_repair_minus_baseline_accuracy"]{metrics["minimum_source_depth3_repair_minus_baseline_accuracy"]=gain}
		}
	}
	for _,depth:=range []int{2,3}{
		if depthQueries[depth]==0{metrics["invalid_row_count"]++;continue}
		prefix:="depth2_";if depth==3{prefix="depth3_"}
		b:=float64(depthBase[depth])/float64(depthQueries[depth])
		r:=float64(depthRepair[depth])/float64(depthQueries[depth])
		metrics[prefix+"query_count"]=float64(depthQueries[depth])
		metrics[prefix+"baseline_pooled_top1_accuracy"]=b
		metrics[prefix+"repair_pooled_top1_accuracy"]=r
		metrics[prefix+"pooled_repair_minus_baseline_accuracy"]=r-b
	}
	if lowQueries>0{
		metrics["depth3_low_support_query_count"]=float64(lowQueries)
		metrics["depth3_low_support_baseline_top1_accuracy"]=float64(lowBase)/float64(lowQueries)
		metrics["depth3_low_support_repair_top1_accuracy"]=float64(lowRepair)/float64(lowQueries)
		metrics["depth3_low_support_repair_minus_baseline_accuracy"]=metrics["depth3_low_support_repair_top1_accuracy"]-metrics["depth3_low_support_baseline_top1_accuracy"]
	}else{metrics["invalid_row_count"]++}
	if math.IsInf(metrics["minimum_source_depth3_repair_minus_baseline_accuracy"],1){metrics["minimum_source_depth3_repair_minus_baseline_accuracy"]=0;metrics["invalid_row_count"]++}
	if metrics["maximum_probability_mass_error"]>1e-9{metrics["invalid_row_count"]++}
	for _,v:=range metrics{if math.IsNaN(v)||math.IsInf(v,0){metrics["invalid_row_count"]++}}

	return wlmLmExternalContextTransitionQualityRepairR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-CONTEXT-TRANSITION-QUALITY-REPAIR-R1",
		Metrics:metrics,
	}
}
