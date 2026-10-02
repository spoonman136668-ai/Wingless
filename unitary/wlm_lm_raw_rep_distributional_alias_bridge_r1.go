package unitary

import (
	"math"
	"sort"
)

type wlmLmRawRepDistributionalAliasBridgeR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmRawRepDistributionalAliasBridgeR1Transition struct {
	counts [5][4][4][16]uint32
}

func (t *wlmLmRawRepDistributionalAliasBridgeR1Transition) observe(op,left,right,outLeft,outRight int,metrics map[string]float64) {
	if op<0||op>=5||left<0||left>=4||right<0||right>=4||outLeft<0||outLeft>=4||outRight<0||outRight>=4 {
		metrics["invalid_row_count"]++
		return
	}
	idx:=outLeft*4+outRight
	if t.counts[op][left][right][idx]==^uint32(0) {
		metrics["counter_overflow_count"]++
		return
	}
	t.counts[op][left][right][idx]++
}

func (t *wlmLmRawRepDistributionalAliasBridgeR1Transition) predict(op,left,right int)(int,int) {
	if op<0||op>=5||left<0||left>=4||right<0||right>=4 {
		return 0,0
	}
	best:=0
	bestCount:=t.counts[op][left][right][0]
	for i:=1;i<16;i++ {
		if t.counts[op][left][right][i]>bestCount {
			best=i
			bestCount=t.counts[op][left][right][i]
		}
	}
	return best/4,best%4
}

func wlmLmRawRepDistributionalAliasBridgeR1CalibrationRecords() [][]uint8 {
	valueCounts:=[]int{17,23,31,41}
	opCounts:=[]int{53,67,83,101,127}
	total:=0
	for _,c:=range valueCounts { total+=2*c }
	for _,c:=range opCounts { total+=2*c }
	out:=make([][]uint8,0,total)
	tag:=0
	for family:=0;family<2;family++ {
		for semantic,count:=range valueCounts {
			key:=wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family,semantic)
			for i:=0;i<count;i++ {
				state:=uint32(0x9e3779b9 ^ uint32((family+1)*100003+(semantic+1)*8191+(i+1)*131))
				raw:=make([]uint8,0,20)
				raw=append(raw,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,5+(tag%3))...)
				raw=append(raw,key[:]...)
				raw=append(raw,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,7+(tag%5))...)
				out=append(out,raw)
				tag++
			}
		}
		for semantic,count:=range opCounts {
			key:=wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(family,semantic)
			for i:=0;i<count;i++ {
				state:=uint32(0x85ebca6b ^ uint32((family+1)*70001+(semantic+1)*12289+(i+1)*313))
				raw:=make([]uint8,0,20)
				raw=append(raw,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,6+(tag%4))...)
				raw=append(raw,key[:]...)
				raw=append(raw,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,8+(tag%6))...)
				out=append(out,raw)
				tag++
			}
		}
	}
	return out
}

func wlmLmRawRepDistributionalAliasBridgeR1SignatureCounts(
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	records [][]uint8,
) map[int]int {
	index:=make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,len(reps))
	for i,k:=range reps { index[k]=i }
	counts:=make(map[int]int,len(reps))
	for _,raw:=range records {
		for i:=0;i+4<=len(raw);i++ {
			k:=wlmSiRawRepAliasInvarianceFalsificationR1Key{raw[i],raw[i+1],raw[i+2],raw[i+3]}
			if id,ok:=index[k];ok {
				counts[id]++
			}
		}
	}
	return counts
}

func wlmLmRawRepDistributionalAliasBridgeR1Pair(
	ids []int,
	signatures map[int]int,
)(map[int]int,int,int,int) {
	groups:=make(map[int][]int)
	zeroIDs:=make([]int,0)
	for _,id:=range ids {
		sig:=signatures[id]
		if sig<=0 {
			zeroIDs=append(zeroIDs,id)
			continue
		}
		groups[sig]=append(groups[sig],id)
	}
	sigs:=make([]int,0,len(groups))
	for sig:=range groups { sigs=append(sigs,sig) }
	sort.Ints(sigs)
	canonical:=make(map[int]int,len(ids))
	pairs:=0
	ambiguous:=0
	unpaired:=len(zeroIDs)
	next:=0
	for _,sig:=range sigs {
		group:=groups[sig]
		sort.Ints(group)
		if len(group)!=2 {
			ambiguous++
			unpaired+=len(group)
			continue
		}
		for _,id:=range group { canonical[id]=next }
		next++
		pairs++
	}
	return canonical,pairs,ambiguous,unpaired
}

func wlmLmRawRepDistributionalAliasBridgeR1SemanticOfCanonicalValues(
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	valueCanonical map[int]int,
	metrics map[string]float64,
) map[int]int {
	out:=make(map[int]int)
	for repID,canonical:=range valueCanonical {
		semantic:=wlmSiRawRepAliasInvarianceFalsificationR1SemanticOfValueKey(reps[repID])
		if semantic<0 {
			metrics["invalid_row_count"]++
			continue
		}
		if prior,ok:=out[canonical];ok && prior!=semantic {
			metrics["invalid_row_count"]++
			continue
		}
		out[canonical]=semantic
	}
	return out
}

func wlmLmRawRepDistributionalAliasBridgeR1Evaluate(
	valueFamily,opFamily int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	valueSurfaceMap,opSurfaceMap map[int]int,
	valueCanonical,opCanonical map[int]int,
	semanticOfSurface,semanticOfCanonical map[int]int,
	unbridged *wlmSiRawRepAliasInvarianceFalsificationR1Transition,
	bridged *wlmLmRawRepDistributionalAliasBridgeR1Transition,
	metrics map[string]float64,
)(int,int,int) {
	bridgedHits:=0
	unbridgedHits:=0
	total:=0
	for opSemantic:=0;opSemantic<5;opSemantic++ {
		for leftSemantic:=0;leftSemantic<4;leftSemantic++ {
			for rightSemantic:=0;rightSemantic<4;rightSemantic++ {
				tag:=3000+valueFamily*1000+opFamily*500+opSemantic*50+leftSemantic*5+rightSemantic
				raw:=wlmSiRawRepAliasInvarianceFalsificationR1EvalRecord(valueFamily,opFamily,leftSemantic,rightSemantic,opSemantic,tag)
				d:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
				total++
				if len(d)!=3 {
					metrics["invalid_row_count"]++
					continue
				}

				sLeft,ok0:=valueSurfaceMap[d[0]]
				sRight,ok1:=valueSurfaceMap[d[1]]
				sOp,ok2:=opSurfaceMap[d[2]]
				cLeft,ok3:=valueCanonical[d[0]]
				cRight,ok4:=valueCanonical[d[1]]
				cOp,ok5:=opCanonical[d[2]]
				if !ok0||!ok1||!ok2||!ok3||!ok4||!ok5 {
					metrics["invalid_row_count"]++
					continue
				}

				uLeft,uRight:=unbridged.predict(sOp,sLeft,sRight)
				uSemLeft,okUL:=semanticOfSurface[uLeft]
				uSemRight,okUR:=semanticOfSurface[uRight]

				bLeft,bRight:=bridged.predict(cOp,cLeft,cRight)
				bSemLeft,okBL:=semanticOfCanonical[bLeft]
				bSemRight,okBR:=semanticOfCanonical[bRight]

				tLeft,tRight:=wlmSiRawRepAliasInvarianceFalsificationR1Apply(leftSemantic,rightSemantic,opSemantic)
				if okUL&&okUR&&uSemLeft==tLeft&&uSemRight==tRight { unbridgedHits++ }
				if okBL&&okBR&&bSemLeft==tLeft&&bSemRight==tRight { bridgedHits++ }
			}
		}
	}
	return bridgedHits,unbridgedHits,total
}

// RunWlmLmRawRepDistributionalAliasBridgeR1 tests a distributional surface-alias canonicalization bridge.
func RunWlmLmRawRepDistributionalAliasBridgeR1() interface{} {
	metrics:=map[string]float64{
		"transition_training_record_count":0,
		"selected_representation_count":0,
		"selected_true_motif_match_count":0,
		"value_alias_pair_count":0,
		"operation_alias_pair_count":0,
		"ambiguous_signature_group_count":0,
		"unpaired_selected_motif_count":0,
		"training_decode_failure_count":0,
		"same_family_evaluation_count":0,
		"bridged_same_family_accuracy":0,
		"unbridged_same_family_accuracy":0,
		"cross_family_evaluation_count":0,
		"bridged_cross_family_accuracy":0,
		"unbridged_cross_family_accuracy":0,
		"bridged_minus_unbridged_cross_family_accuracy":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	samples:=make([][]uint8,0,3200)
	for repeat:=0;repeat<20;repeat++ {
		for family:=0;family<2;family++ {
			for op:=0;op<5;op++ {
				for left:=0;left<4;left++ {
					for right:=0;right<4;right++ {
						raw:=wlmSiRawRepAliasInvarianceFalsificationR1TrainingRecord(family,left,right,op,repeat)
						samples=append(samples,raw)
						metrics["transition_training_record_count"]++
					}
				}
			}
		}
	}

	reps:=wlmSiRawRepAliasInvarianceFalsificationR1LearnRepresentation(samples)
	metrics["selected_representation_count"]=float64(len(reps))
	metrics["selected_true_motif_match_count"]=float64(wlmSiRawRepAliasInvarianceFalsificationR1TrueMotifCount(reps))

	valueSet:=make(map[int]bool)
	opSet:=make(map[int]bool)
	for _,raw:=range samples {
		d:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
		if len(d)!=5 { continue }
		valueSet[d[0]]=true
		valueSet[d[1]]=true
		opSet[d[2]]=true
		valueSet[d[3]]=true
		valueSet[d[4]]=true
	}
	valueIDs:=make([]int,0,len(valueSet))
	for id:=range valueSet { valueIDs=append(valueIDs,id) }
	opIDs:=make([]int,0,len(opSet))
	for id:=range opSet { opIDs=append(opIDs,id) }
	sort.Ints(valueIDs)
	sort.Ints(opIDs)

	calibration:=wlmLmRawRepDistributionalAliasBridgeR1CalibrationRecords()
	signatures:=wlmLmRawRepDistributionalAliasBridgeR1SignatureCounts(reps,calibration)
	valueCanonical,valuePairs,valueAmbiguous,valueUnpaired:=wlmLmRawRepDistributionalAliasBridgeR1Pair(valueIDs,signatures)
	opCanonical,opPairs,opAmbiguous,opUnpaired:=wlmLmRawRepDistributionalAliasBridgeR1Pair(opIDs,signatures)
	metrics["value_alias_pair_count"]=float64(valuePairs)
	metrics["operation_alias_pair_count"]=float64(opPairs)
	metrics["ambiguous_signature_group_count"]=float64(valueAmbiguous+opAmbiguous)
	metrics["unpaired_selected_motif_count"]=float64(valueUnpaired+opUnpaired)

	valueSurfaceMap:=wlmSiRawRepAliasInvarianceFalsificationR1IndexMap(valueIDs)
	opSurfaceMap:=wlmSiRawRepAliasInvarianceFalsificationR1IndexMap(opIDs)

	semanticOfSurface:=make(map[int]int,len(valueIDs))
	for surface,repID:=range valueIDs {
		semanticOfSurface[surface]=wlmSiRawRepAliasInvarianceFalsificationR1SemanticOfValueKey(reps[repID])
	}
	semanticOfCanonical:=wlmLmRawRepDistributionalAliasBridgeR1SemanticOfCanonicalValues(reps,valueCanonical,metrics)

	var unbridged wlmSiRawRepAliasInvarianceFalsificationR1Transition
	var bridged wlmLmRawRepDistributionalAliasBridgeR1Transition

	for _,raw:=range samples {
		d:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
		if len(d)!=5 {
			metrics["training_decode_failure_count"]++
			continue
		}
		sLeft,ok0:=valueSurfaceMap[d[0]]
		sRight,ok1:=valueSurfaceMap[d[1]]
		sOp,ok2:=opSurfaceMap[d[2]]
		sOutLeft,ok3:=valueSurfaceMap[d[3]]
		sOutRight,ok4:=valueSurfaceMap[d[4]]
		cLeft,ok5:=valueCanonical[d[0]]
		cRight,ok6:=valueCanonical[d[1]]
		cOp,ok7:=opCanonical[d[2]]
		cOutLeft,ok8:=valueCanonical[d[3]]
		cOutRight,ok9:=valueCanonical[d[4]]
		if !ok0||!ok1||!ok2||!ok3||!ok4||!ok5||!ok6||!ok7||!ok8||!ok9 {
			metrics["training_decode_failure_count"]++
			continue
		}
		unbridged.observe(sOp,sLeft,sRight,sOutLeft,sOutRight,metrics)
		bridged.observe(cOp,cLeft,cRight,cOutLeft,cOutRight,metrics)
	}

	bridgedSame:=0
	unbridgedSame:=0
	sameTotal:=0
	for family:=0;family<2;family++ {
		b,u,n:=wlmLmRawRepDistributionalAliasBridgeR1Evaluate(
			family,family,reps,valueSurfaceMap,opSurfaceMap,valueCanonical,opCanonical,
			semanticOfSurface,semanticOfCanonical,&unbridged,&bridged,metrics,
		)
		bridgedSame+=b
		unbridgedSame+=u
		sameTotal+=n
	}
	metrics["same_family_evaluation_count"]=float64(sameTotal)
	if sameTotal>0 {
		metrics["bridged_same_family_accuracy"]=float64(bridgedSame)/float64(sameTotal)
		metrics["unbridged_same_family_accuracy"]=float64(unbridgedSame)/float64(sameTotal)
	}

	bridgedCross:=0
	unbridgedCross:=0
	crossTotal:=0
	for valueFamily:=0;valueFamily<2;valueFamily++ {
		opFamily:=1-valueFamily
		b,u,n:=wlmLmRawRepDistributionalAliasBridgeR1Evaluate(
			valueFamily,opFamily,reps,valueSurfaceMap,opSurfaceMap,valueCanonical,opCanonical,
			semanticOfSurface,semanticOfCanonical,&unbridged,&bridged,metrics,
		)
		bridgedCross+=b
		unbridgedCross+=u
		crossTotal+=n
	}
	metrics["cross_family_evaluation_count"]=float64(crossTotal)
	if crossTotal>0 {
		metrics["bridged_cross_family_accuracy"]=float64(bridgedCross)/float64(crossTotal)
		metrics["unbridged_cross_family_accuracy"]=float64(unbridgedCross)/float64(crossTotal)
		metrics["bridged_minus_unbridged_cross_family_accuracy"]=metrics["bridged_cross_family_accuracy"]-metrics["unbridged_cross_family_accuracy"]
	}

	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) {
			metrics["invalid_row_count"]++
		}
	}
	return wlmLmRawRepDistributionalAliasBridgeR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-REP-DISTRIBUTIONAL-ALIAS-BRIDGE-R1",
		Metrics:metrics,
	}
}
