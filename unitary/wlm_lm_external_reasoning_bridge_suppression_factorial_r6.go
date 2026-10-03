package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalReasoningBridgeSuppressionFactorialR6Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`

}

func wlmLmExternalReasoningBridgeSuppressionFactorialR6Terminal(
	prev,curr,depth,capStates int,
	lowSupportFallback,exactContext bool,
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
)([512]float64,float64,float64) {
	var zero [512]float64
	active:=[]wlmLmExternalRepairedRelationReasoningBridgeR1PairMass{{prev:prev,curr:curr,mass:1}}
	maxMassError:=0.0
	minRetained:=1.0
	for step:=0;step<depth;step++ {
		next:=make(map[[2]int]float64)
		for _,state:=range active {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:state.prev,curr:state.curr}
			support:=pairCounts[p]
			id,dedicated:=exactID[p]
			if lowSupportFallback && support<3 {
				currentTotal:=rowTotals[state.curr]
				overflowTotal:=exactTotals[511]
				switch {
				case currentTotal>0 && overflowTotal>0:
					cd:=float64(currentTotal); od:=float64(overflowTotal)
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
			useExact:=exactContext && dedicated && support>=3 && rowEligible[state.curr] && rowMargin[state.curr]<=median
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
		entries:=make([]wlmLmExternalRepairedRelationReasoningBridgeR1PairMass,0,len(next))
		for p,mass:=range next {
			entries=append(entries,wlmLmExternalRepairedRelationReasoningBridgeR1PairMass{prev:p[0],curr:p[1],mass:mass})
		}
		sort.Slice(entries,func(i,j int)bool{
			if entries[i].mass!=entries[j].mass{return entries[i].mass>entries[j].mass}
			if entries[i].prev!=entries[j].prev{return entries[i].prev<entries[j].prev}
			return entries[i].curr<entries[j].curr
		})
		total:=0.0
		for _,e:=range entries{total+=e.mass}
		if err:=math.Abs(1-total);err>maxMassError{maxMassError=err}
		if total==0{return zero,maxMassError,0}
		if len(entries)>capStates{entries=entries[:capStates]}
		kept:=0.0
		for _,e:=range entries{kept+=e.mass}
		if kept==0{return zero,maxMassError,0}
		if r:=kept/total;r<minRetained{minRetained=r}
		for i:=range entries{entries[i].mass/=kept}
		active=entries
	}
	var final [512]float64
	total:=0.0
	for _,e:=range active{final[e.curr]+=e.mass;total+=e.mass}
	if err:=math.Abs(1-total);err>maxMassError{maxMassError=err}
	return final,maxMassError,minRetained
}

func RunWlmLmExternalReasoningBridgeSuppressionFactorialR6(code,structured,prose []byte) interface{} {
	m:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":3,
		"total_source_bytes":float64(len(code)+len(structured)+len(prose)),
		"training_byte_budget":15819,
		"target_eval_byte_budget":1454,
		"selected_motif_count":0,
		"reasoning_class_count":8,
		"minimum_reasoning_class_size":256,
		"maximum_reasoning_class_size":0,
		"training_margin_sample_count":0,
		"zero_training_margin_class_count":0,
		"preflight_unevaluable_flag":0,
		"preflight_target_byte_use_count":0,
		"preflight_target_label_use_count":0,
		"preflight_decision_frozen_flag":0,
		"target_query_count":0,
		"unsupported_target_query_count":0,
		"unsupported_target_exact_hit_count":0,
		"zero_support_class_count":0,
		"zero_support_class_identity_mismatch_count":0,
		"zero_support_target_top1_count":0,
		"target_representation_absence_count":0,
		"bridge_suppression_count":0,
		"diagnostic_active_count":0,
		"zero_support_attribution_record_count":0,
		"attribution_accounting_error_count":0,
		"suppressed_class_count":2,
		"suppression_attribution_record_count":0,
		"suppression_attribution_accounting_error_count":0,
		"base_relation_absence_count":0,
		"exact_context_suppression_count":0,
		"repair_selection_suppression_count":0,
		"final_competition_count":0,
		"class5_repaired_top1_count":0,
		"class6_repaired_top1_count":0,
		"selected_query_count":0,
		"factorial_record_count":0,
		"factorial_accounting_error_count":0,
		"low_support_main_count":0,
		"exact_context_main_count":0,
		"additive_shared_count":0,
		"interaction_only_count":0,
		"full_no_suppression_count":0,
		"unresolved_count":0,
		"B11_repaired_top1_class5_count":0,
		"B11_repaired_top1_class6_count":0,
		"B11_identity_mismatch_count":0,
		"selected_query_base_class_mass_sum":0,
		"selected_query_repaired_class_mass_sum":0,
		"maximum_probability_mass_error":0,
		"minimum_repaired_retained_mass":1,
		"relation_capacity_growth_event_count":0,
		"adaptive_readout_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"counter_overflow_count":0,
		"counter_overflow_rows":0,
		"invalid_row_count":0,
	}
	if len(code)!=41453 || wlmLmExternalRawRepPredR1SHA256(code)!="7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6" { m["source_identity_mismatch_count"]++ }
	if len(structured)!=14365 || wlmLmExternalRawRepPredR1SHA256(structured)!="4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192" { m["source_identity_mismatch_count"]++ }
	if len(prose)!=1454 || wlmLmExternalRawRepPredR1SHA256(prose)!="8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3" { m["source_identity_mismatch_count"]++ }
	if len(structured)+len(prose)!=15819 { m["invalid_row_count"]++ }

	train:=[][]byte{structured,prose}
	selected:=wlmLmProspectiveR3Select256(train,m)
	m["counter_overflow_count"]+=m["counter_overflow_rows"]
	delete(m,"counter_overflow_rows")
	m["selected_motif_count"]=float64(len(selected))
	if len(selected)!=256 { m["invalid_row_count"]++ }

	keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index:=make(map[[4]uint8]int,len(keys))
	for i,k:=range keys { index[k]=i }
	streams:=[][]int{
		wlmLmExternalMotifRelationHoldoutR1Decode(structured,index),
		wlmLmExternalMotifRelationHoldoutR1Decode(prose,index),
	}
	var counts [512]uint64
	var rel [512][512]uint32
	var totals [512]uint64
	var uni [512]uint32
	var unitotal uint64
	pairs:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for _,st:=range streams {
		for _,x:=range st {
			if x<0||x>=256 { m["invalid_row_count"]++; continue }
			counts[x]++
		}
		for j:=0;j+1<len(st);j++ {
			a,b:=st[j],st[j+1]
			if rel[a][b]==^uint32(0)||uni[b]==^uint32(0){ m["counter_overflow_count"]++; continue }
			rel[a][b]++; totals[a]++; uni[b]++; unitotal++
		}
		for j:=1;j<len(st);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:st[j-1],curr:st[j]}
			if pairs[p]==^uint32(0){ m["counter_overflow_count"]++; continue }
			pairs[p]++
		}
	}
	if unitotal==0 { m["invalid_row_count"]++ }

	ids:=make([]int,256)
	for i:=range ids { ids[i]=i }
	sort.Slice(ids,func(i,j int)bool{
		if counts[ids[i]]!=counts[ids[j]] { return counts[ids[i]]>counts[ids[j]] }
		return ids[i]<ids[j]
	})
	var classes [512]int
	var classSizes [8]int
	for rank,id:=range ids {
		c:=rank/32
		classes[id]=c
		classSizes[c]++
	}
	for i:=0;i<8;i++ {
		n:=float64(classSizes[i])
		if n<m["minimum_reasoning_class_size"] { m["minimum_reasoning_class_size"]=n }
		if n>m["maximum_reasoning_class_size"] { m["maximum_reasoning_class_size"]=n }
	}

	rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairs))
	for p,c:=range pairs { rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:p,count:c}) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count { return rows[i].count>rows[j].count }
		if rows[i].pair.prev!=rows[j].pair.prev { return rows[i].pair.prev<rows[j].pair.prev }
		return rows[i].pair.curr<rows[j].pair.curr
	})
	n:=511
	if len(rows)<n { n=len(rows) }
	exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,n)
	for i:=0;i<n;i++ { exactID[rows[i].pair]=i }
	var exact [512][512]uint32
	var exactTotals [512]uint64
	for _,st:=range streams {
		for j:=1;j+1<len(st);j++ {
			p:=wlmLmExternalContextStateHoldoutR1Pair{prev:st[j-1],curr:st[j]}
			id:=511
			if v,ok:=exactID[p];ok { id=v }
			t:=st[j+1]
			if exact[id][t]==^uint32(0){ m["counter_overflow_count"]++; continue }
			exact[id][t]++; exactTotals[id]++
		}
	}
	var eligible [512]bool
	var marginsByRow [512]float64
	rowMargins:=make([]float64,0,256)
	for i:=0;i<256;i++ {
		var total uint64
		var top,runner uint32
		for j:=0;j<256;j++ {
			v:=rel[i][j]
			total+=uint64(v)
			if v>top { runner=top; top=v } else if v>runner { runner=v }
		}
		if total>=4 {
			eligible[i]=true
			marginsByRow[i]=(float64(top)-float64(runner))/float64(total)
			rowMargins=append(rowMargins,marginsByRow[i])
		}
	}
	if len(rowMargins)==0 { m["invalid_row_count"]++ }
	median:=wlmLmProspectiveR3Median(rowMargins)

	gateMargin:=func(dist *[512]float64) float64 {
		top:=wlmLmExternalReasoningReadoutRefinementR1Top1(dist)
		if top<0||top>=256 { m["invalid_row_count"]++; return 0 }
		cid:=classes[top]
		mass,runner:=0.0,0.0
		for k:=0;k<256;k++ {
			if classes[k]!=cid { continue }
			v:=dist[k]
			mass+=v
			if k!=top && v>runner { runner=v }
		}
		if mass<=0 { m["invalid_row_count"]++; return 0 }
		return (dist[top]-runner)/mass
	}

	classTrainMargins:=make([][]float64,8)
	for _,st:=range streams {
		for j:=0;j+4<len(st);j++ {
			rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(
				st[j],st[j+1],3,2048,true,
				&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal,
			)
			if re>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=re }
			if rr<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=rr }
			_ = gateMargin(&rd)
			top:=wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
			if top<0||top>=256 { m["invalid_row_count"]++; continue }
			classTrainMargins[classes[top]]=append(classTrainMargins[classes[top]],gateMargin(&rd))
			m["training_margin_sample_count"]++
		}
	}
	zeroSet:=make(map[int]bool)
	for i:=0;i<8;i++ {
		n:=len(classTrainMargins[i])
		m["class_"+string(rune('0'+i))+"_training_margin_sample_count"]=float64(n)
		if n==0 {
			m["zero_training_margin_class_count"]++
			zeroSet[i]=true
		}
	}
	if len(zeroSet)>0 { m["preflight_unevaluable_flag"]=1 }
	m["preflight_decision_frozen_flag"]=1

	// Target bytes are used scientifically only after the training-only preflight is frozen.
	target:=code[:1454]
	targetStream:=wlmLmExternalMotifRelationHoldoutR1Decode(target,index)
	var rawClassCount [8]float64
	var topClassCount [8]float64
	var repairedMassByClass [8]float64
	for _,id:=range targetStream {
		if id<0||id>=256 { m["invalid_row_count"]++; continue }
		rawClassCount[classes[id]]++
	}
	for j:=0;j+4<len(targetStream);j++ {
		truth:=targetStream[j+4]
		rawFocus:=map[int]bool{}
		for q:=j;q<=j+4;q++ {
			id:=targetStream[q]
			if id>=0 && id<256 && (classes[id]==5 || classes[id]==6) { rawFocus[classes[id]]=true }
		}
		bd,be,br:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(
			targetStream[j],targetStream[j+1],3,2048,false,
			&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal,
		)
		if be>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=be }
		if br<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=br }
		rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(
			targetStream[j],targetStream[j+1],3,2048,true,
			&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal,
		)
		b10,e10,r10:=wlmLmExternalReasoningBridgeSuppressionFactorialR6Terminal(
			targetStream[j],targetStream[j+1],3,2048,true,false,
			&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal,
		)
		b01,e01,r01:=wlmLmExternalReasoningBridgeSuppressionFactorialR6Terminal(
			targetStream[j],targetStream[j+1],3,2048,false,true,
			&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal,
		)
		b11,e11,r11:=wlmLmExternalReasoningBridgeSuppressionFactorialR6Terminal(
			targetStream[j],targetStream[j+1],3,2048,true,true,
			&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal,
		)
		for _,errv:=range []float64{e10,e01,e11} { if errv>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=errv } }
		for _,ret:=range []float64{r10,r01,r11} { if ret<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=ret } }
		mismatch:=false
		for k:=0;k<256;k++ { if math.Abs(b11[k]-rd[k])>1e-15 { mismatch=true; break } }
		if mismatch { m["B11_identity_mismatch_count"]++ }
		if re>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=re }
		if rr<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=rr }
		top:=wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
		if top<0||top>=256 { m["invalid_row_count"]++; continue }
		m["target_query_count"]++
		if classes[top]==5 { m["class5_repaired_top1_count"]++; m["B11_repaired_top1_class5_count"]++ }
		if classes[top]==6 { m["class6_repaired_top1_count"]++; m["B11_repaired_top1_class6_count"]++ }
		if len(rawFocus)>0 {
			m["selected_query_count"]++
			m["factorial_record_count"]++
			type rp struct{ id int; p float64 }
			bestRank:=func(dist *[512]float64)(int,float64){
				rankv:=make([]rp,0,256)
				mass:=0.0
				for k:=0;k<256;k++ {
					rankv=append(rankv,rp{k,dist[k]})
					if rawFocus[classes[k]] { mass+=dist[k] }
				}
				sort.Slice(rankv,func(a,b int)bool{if rankv[a].p!=rankv[b].p{return rankv[a].p>rankv[b].p};return rankv[a].id<rankv[b].id})
				best:=513
				for rank,x:=range rankv { if rawFocus[classes[x.id]] { best=rank+1; break } }
				return best,mass
			}
			r00,mass00:=bestRank(&bd)
			r10,mass10:=bestRank(&b10)
			r01,mass01:=bestRank(&b01)
			r11,mass11:=bestRank(&b11)
			m["selected_query_base_class_mass_sum"]+=mass00
			m["selected_query_repaired_class_mass_sum"]+=mass11
			worse:=func(rank,base int,mass,baseMass float64)bool{
				if baseMass>0 && mass==0 { return true }
				return rank>base && rank>1
			}
			low:=worse(r10,r00,mass10,mass00)
			exactOnly:=worse(r01,r00,mass01,mass00)
			full:=worse(r11,r00,mass11,mass00)
			category:="unresolved"
			switch {
			case !full:
				category="full_no_suppression"
			case low && !exactOnly:
				category="low_support_main"
			case !low && exactOnly:
				category="exact_context_main"
			case low && exactOnly:
				category="additive_shared"
			case !low && !exactOnly:
				category="interaction_only"
			}
			switch category {
			case "low_support_main": m["low_support_main_count"]++
			case "exact_context_main": m["exact_context_main_count"]++
			case "additive_shared": m["additive_shared_count"]++
			case "interaction_only": m["interaction_only_count"]++
			case "full_no_suppression": m["full_no_suppression_count"]++
			default:
				m["unresolved_count"]++
				m["factorial_accounting_error_count"]++
			}
			_ = r10; _ = r01
		}
		topClass:=classes[top]
		topClassCount[topClass]++
		for k:=0;k<256;k++ {
			repairedMassByClass[classes[k]]+=rd[k]
		}
		if zeroSet[topClass] {
			m["unsupported_target_query_count"]++
			if top==truth { m["unsupported_target_exact_hit_count"]++ }
		}
	}
	expectedZero:=map[int]bool{5:true,6:true,7:true}
	m["zero_support_class_count"]=float64(len(zeroSet))
	for i:=0;i<8;i++ {
		if zeroSet[i]!=expectedZero[i] { m["zero_support_class_identity_mismatch_count"]++ }
		m["class_"+string(rune('0'+i))+"_target_raw_decoded_motif_count"]=rawClassCount[i]
		m["class_"+string(rune('0'+i))+"_target_repaired_top1_count"]=topClassCount[i]
		m["class_"+string(rune('0'+i))+"_target_repaired_probability_mass"]=repairedMassByClass[i]
	}
	for _,cid:=range []int{5,6,7} {
		m["zero_support_attribution_record_count"]++
		m["zero_support_target_top1_count"]+=topClassCount[cid]
		if rawClassCount[cid]==0 {
			m["target_representation_absence_count"]++
		} else if topClassCount[cid]==0 {
			m["bridge_suppression_count"]++
		} else {
			m["diagnostic_active_count"]++
		}
	}
	if m["target_representation_absence_count"]+m["bridge_suppression_count"]+m["diagnostic_active_count"]!=m["zero_support_attribution_record_count"] {
		m["attribution_accounting_error_count"]++
	}
	if m["base_relation_absence_count"]+m["exact_context_suppression_count"]+m["repair_selection_suppression_count"]+m["final_competition_count"]!=m["suppression_attribution_record_count"] {
		m["suppression_attribution_accounting_error_count"]++
	}
	if m["suppression_attribution_record_count"]!=m["selected_query_count"] { m["suppression_attribution_accounting_error_count"]++ }
	if m["low_support_main_count"]+m["exact_context_main_count"]+m["additive_shared_count"]+m["interaction_only_count"]+m["full_no_suppression_count"]+m["unresolved_count"]!=m["factorial_record_count"] {
		m["factorial_accounting_error_count"]++
	}
	if m["factorial_record_count"]!=m["selected_query_count"] { m["factorial_accounting_error_count"]++ }
	if m["maximum_probability_mass_error"]>1e-9 { m["invalid_row_count"]++ }
	for _,v:=range m {
		if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ }
	}
	return wlmLmExternalReasoningBridgeSuppressionFactorialR6Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-REASONING-BRIDGE-SUPPRESSION-FACTORIAL-R6",
		Metrics:m,
	}
}
