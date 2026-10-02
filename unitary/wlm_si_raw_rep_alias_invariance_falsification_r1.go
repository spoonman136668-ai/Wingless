package unitary

import (
	"math"
	"sort"
)

type wlmSiRawRepAliasInvarianceFalsificationR1Key [4]uint8

type wlmSiRawRepAliasInvarianceFalsificationR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmSiRawRepAliasInvarianceFalsificationR1Sample struct {
	raw []uint8
}

type wlmSiRawRepAliasInvarianceFalsificationR1Transition struct {
	counts [10][8][8][64]uint32
}

func wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family,id int) wlmSiRawRepAliasInvarianceFalsificationR1Key {
	base:=uint8(16+8*id)
	if family==1 {
		base=uint8(56+8*id)
	}
	return wlmSiRawRepAliasInvarianceFalsificationR1Key{base,base+1,base+2,base+3}
}

func wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(family,id int) wlmSiRawRepAliasInvarianceFalsificationR1Key {
	base:=uint8(96+8*id)
	if family==1 {
		base=uint8(136+8*id)
	}
	return wlmSiRawRepAliasInvarianceFalsificationR1Key{base,base+1,base+2,base+3}
}

func wlmSiRawRepAliasInvarianceFalsificationR1Noise(state *uint32,n int) []uint8 {
	out:=make([]uint8,n)
	for i:=0;i<n;i++ {
		*state=*state*1664525+1013904223
		out[i]=uint8(192+((*state>>16)%64))
	}
	return out
}

func wlmSiRawRepAliasInvarianceFalsificationR1Apply(left,right,op int)(int,int) {
	if op==4 {
		return right,left
	}
	a:=op
	b:=(op+1)%4
	switch left {
	case a:
		left=b
	case b:
		left=a
	}
	return left,right
}

func wlmSiRawRepAliasInvarianceFalsificationR1Append(out []uint8,k wlmSiRawRepAliasInvarianceFalsificationR1Key) []uint8 {
	return append(out,k[:]...)
}

func wlmSiRawRepAliasInvarianceFalsificationR1TrainingRecord(family,left,right,op,repeat int) []uint8 {
	outLeft,outRight:=wlmSiRawRepAliasInvarianceFalsificationR1Apply(left,right,op)
	state:=uint32(0x27d4eb2d ^ uint32((family+1)*100003+(left+1)*131+(right+1)*977+(op+1)*8191+(repeat+1)*65537))
	out:=make([]uint8,0,90)
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,5+(repeat+family+left)%5)...)
	out=wlmSiRawRepAliasInvarianceFalsificationR1Append(out,wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family,left))
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,7+(repeat+right)%5)...)
	out=wlmSiRawRepAliasInvarianceFalsificationR1Append(out,wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family,right))
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,9+(repeat+op)%5)...)
	out=wlmSiRawRepAliasInvarianceFalsificationR1Append(out,wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(family,op))
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,11+(repeat+left+right)%5)...)
	out=wlmSiRawRepAliasInvarianceFalsificationR1Append(out,wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family,outLeft))
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,13+(repeat+op+right)%5)...)
	out=wlmSiRawRepAliasInvarianceFalsificationR1Append(out,wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family,outRight))
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,5+(repeat+family+op)%5)...)
	return out
}

func wlmSiRawRepAliasInvarianceFalsificationR1EvalRecord(valueFamily,opFamily,left,right,op,tag int) []uint8 {
	state:=uint32(0x165667b1 ^ uint32((valueFamily+1)*1009+(opFamily+1)*9176+(left+1)*131+(right+1)*977+(op+1)*8191+(tag+1)*65537))
	out:=make([]uint8,0,60)
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,6+(left+tag)%5)...)
	out=wlmSiRawRepAliasInvarianceFalsificationR1Append(out,wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(valueFamily,left))
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,8+(right+tag)%5)...)
	out=wlmSiRawRepAliasInvarianceFalsificationR1Append(out,wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(valueFamily,right))
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,10+(op+tag)%5)...)
	out=wlmSiRawRepAliasInvarianceFalsificationR1Append(out,wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(opFamily,op))
	out=append(out,wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state,7+(tag+left+right)%5)...)
	return out
}

func wlmSiRawRepAliasInvarianceFalsificationR1Less(a,b wlmSiRawRepAliasInvarianceFalsificationR1Key) bool {
	for i:=0;i<4;i++ {
		if a[i]!=b[i] {
			return a[i]<b[i]
		}
	}
	return false
}

func wlmSiRawRepAliasInvarianceFalsificationR1LearnRepresentation(records [][]uint8) []wlmSiRawRepAliasInvarianceFalsificationR1Key {
	counts:=make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]uint32)
	for _,data:=range records {
		for i:=0;i+4<=len(data);i++ {
			k:=wlmSiRawRepAliasInvarianceFalsificationR1Key{data[i],data[i+1],data[i+2],data[i+3]}
			counts[k]++
		}
	}
	type row struct{key wlmSiRawRepAliasInvarianceFalsificationR1Key;count uint32}
	rows:=make([]row,0,len(counts))
	for k,c:=range counts {
		rows=append(rows,row{k,c})
	}
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count {
			return rows[i].count>rows[j].count
		}
		return wlmSiRawRepAliasInvarianceFalsificationR1Less(rows[i].key,rows[j].key)
	})
	if len(rows)>18 {
		rows=rows[:18]
	}
	out:=make([]wlmSiRawRepAliasInvarianceFalsificationR1Key,len(rows))
	for i,r:=range rows {
		out[i]=r.key
	}
	sort.Slice(out,func(i,j int)bool{return wlmSiRawRepAliasInvarianceFalsificationR1Less(out[i],out[j])})
	return out
}

func wlmSiRawRepAliasInvarianceFalsificationR1Decode(data []uint8,reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) []int {
	index:=make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,len(reps))
	for i,k:=range reps {
		index[k]=i
	}
	out:=make([]int,0,5)
	for i:=0;i+4<=len(data); {
		k:=wlmSiRawRepAliasInvarianceFalsificationR1Key{data[i],data[i+1],data[i+2],data[i+3]}
		if id,ok:=index[k];ok {
			out=append(out,id)
			i+=4
			continue
		}
		i++
	}
	return out
}

func wlmSiRawRepAliasInvarianceFalsificationR1IndexMap(ids []int) map[int]int {
	out:=make(map[int]int,len(ids))
	for i,id:=range ids {
		out[id]=i
	}
	return out
}

func wlmSiRawRepAliasInvarianceFalsificationR1PairIndex(left,right int) int {
	return left*8+right
}

func wlmSiRawRepAliasInvarianceFalsificationR1PairDecode(index int)(int,int) {
	return index/8,index%8
}

func (t *wlmSiRawRepAliasInvarianceFalsificationR1Transition) observe(op,left,right,outLeft,outRight int,metrics map[string]float64) {
	if op<0||op>=10||left<0||left>=8||right<0||right>=8||outLeft<0||outLeft>=8||outRight<0||outRight>=8 {
		metrics["invalid_row_count"]++
		return
	}
	idx:=wlmSiRawRepAliasInvarianceFalsificationR1PairIndex(outLeft,outRight)
	if t.counts[op][left][right][idx]==^uint32(0) {
		metrics["counter_overflow_count"]++
		return
	}
	t.counts[op][left][right][idx]++
}

func (t *wlmSiRawRepAliasInvarianceFalsificationR1Transition) has(op,left,right int) bool {
	if op<0||op>=10||left<0||left>=8||right<0||right>=8 {
		return false
	}
	var total uint64
	for i:=0;i<64;i++ {
		total+=uint64(t.counts[op][left][right][i])
	}
	return total>0
}

func (t *wlmSiRawRepAliasInvarianceFalsificationR1Transition) predict(op,left,right int)(int,int) {
	if op<0||op>=10||left<0||left>=8||right<0||right>=8 {
		return 0,0
	}
	best:=0
	bestCount:=t.counts[op][left][right][0]
	for i:=1;i<64;i++ {
		if t.counts[op][left][right][i]>bestCount {
			best=i
			bestCount=t.counts[op][left][right][i]
		}
	}
	return wlmSiRawRepAliasInvarianceFalsificationR1PairDecode(best)
}

func wlmSiRawRepAliasInvarianceFalsificationR1TrueMotifCount(reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) int {
	count:=0
	for _,r:=range reps {
		matched:=false
		for family:=0;family<2;family++ {
			for i:=0;i<4;i++ {
				if r==wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family,i) {
					matched=true
				}
			}
			for i:=0;i<5;i++ {
				if r==wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(family,i) {
					matched=true
				}
			}
		}
		if matched {
			count++
		}
	}
	return count
}

func wlmSiRawRepAliasInvarianceFalsificationR1SemanticOfValueKey(k wlmSiRawRepAliasInvarianceFalsificationR1Key) int {
	for family:=0;family<2;family++ {
		for semantic:=0;semantic<4;semantic++ {
			if k==wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family,semantic) {
				return semantic
			}
		}
	}
	return -1
}

// RunWlmSiRawRepAliasInvarianceFalsificationR1 tests whether exact raw representations are surface-alias invariant.
func RunWlmSiRawRepAliasInvarianceFalsificationR1() interface{} {
	metrics:=map[string]float64{
		"training_record_count":0,
		"selected_representation_count":0,
		"selected_true_motif_match_count":0,
		"surface_value_representation_count":0,
		"surface_operation_representation_count":0,
		"training_decode_failure_count":0,
		"same_family_evaluation_count":0,
		"same_family_accuracy":0,
		"cross_family_evaluation_count":0,
		"cross_family_accuracy":0,
		"cross_family_missing_transition_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	samples:=make([]wlmSiRawRepAliasInvarianceFalsificationR1Sample,0,3200)
	records:=make([][]uint8,0,3200)
	for repeat:=0;repeat<20;repeat++ {
		for family:=0;family<2;family++ {
			for op:=0;op<5;op++ {
				for left:=0;left<4;left++ {
					for right:=0;right<4;right++ {
						raw:=wlmSiRawRepAliasInvarianceFalsificationR1TrainingRecord(family,left,right,op,repeat)
						samples=append(samples,wlmSiRawRepAliasInvarianceFalsificationR1Sample{raw:raw})
						records=append(records,raw)
						metrics["training_record_count"]++
					}
				}
			}
		}
	}

	reps:=wlmSiRawRepAliasInvarianceFalsificationR1LearnRepresentation(records)
	metrics["selected_representation_count"]=float64(len(reps))
	metrics["selected_true_motif_match_count"]=float64(wlmSiRawRepAliasInvarianceFalsificationR1TrueMotifCount(reps))

	valueSet:=make(map[int]bool)
	opSet:=make(map[int]bool)
	for _,sample:=range samples {
		d:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(sample.raw,reps)
		if len(d)!=5 {
			continue
		}
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
	metrics["surface_value_representation_count"]=float64(len(valueIDs))
	metrics["surface_operation_representation_count"]=float64(len(opIDs))
	valueMap:=wlmSiRawRepAliasInvarianceFalsificationR1IndexMap(valueIDs)
	opMap:=wlmSiRawRepAliasInvarianceFalsificationR1IndexMap(opIDs)

	semanticOfSurface:=make(map[int]int,len(valueIDs))
	for surface,id:=range valueIDs {
		semanticOfSurface[surface]=wlmSiRawRepAliasInvarianceFalsificationR1SemanticOfValueKey(reps[id])
	}

	var transition wlmSiRawRepAliasInvarianceFalsificationR1Transition
	for _,sample:=range samples {
		d:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(sample.raw,reps)
		if len(d)!=5 {
			metrics["training_decode_failure_count"]++
			continue
		}
		left,ok0:=valueMap[d[0]]
		right,ok1:=valueMap[d[1]]
		op,ok2:=opMap[d[2]]
		outLeft,ok3:=valueMap[d[3]]
		outRight,ok4:=valueMap[d[4]]
		if !ok0||!ok1||!ok2||!ok3||!ok4 {
			metrics["training_decode_failure_count"]++
			continue
		}
		transition.observe(op,left,right,outLeft,outRight,metrics)
	}

	sameHits:=0
	for family:=0;family<2;family++ {
		for opSemantic:=0;opSemantic<5;opSemantic++ {
			for leftSemantic:=0;leftSemantic<4;leftSemantic++ {
				for rightSemantic:=0;rightSemantic<4;rightSemantic++ {
					raw:=wlmSiRawRepAliasInvarianceFalsificationR1EvalRecord(family,family,leftSemantic,rightSemantic,opSemantic,1000+family*100+opSemantic*20+leftSemantic*4+rightSemantic)
					d:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
					metrics["same_family_evaluation_count"]++
					if len(d)!=3 {
						metrics["invalid_row_count"]++
						continue
					}
					left,ok0:=valueMap[d[0]]
					right,ok1:=valueMap[d[1]]
					op,ok2:=opMap[d[2]]
					if !ok0||!ok1||!ok2 {
						metrics["invalid_row_count"]++
						continue
					}
					pl,pr:=transition.predict(op,left,right)
					psl,okL:=semanticOfSurface[pl]
					psr,okR:=semanticOfSurface[pr]
					tl,tr:=wlmSiRawRepAliasInvarianceFalsificationR1Apply(leftSemantic,rightSemantic,opSemantic)
					if okL&&okR&&psl==tl&&psr==tr {
						sameHits++
					}
				}
			}
		}
	}

	crossHits:=0
	for valueFamily:=0;valueFamily<2;valueFamily++ {
		opFamily:=1-valueFamily
		for opSemantic:=0;opSemantic<5;opSemantic++ {
			for leftSemantic:=0;leftSemantic<4;leftSemantic++ {
				for rightSemantic:=0;rightSemantic<4;rightSemantic++ {
					raw:=wlmSiRawRepAliasInvarianceFalsificationR1EvalRecord(valueFamily,opFamily,leftSemantic,rightSemantic,opSemantic,2000+valueFamily*100+opSemantic*20+leftSemantic*4+rightSemantic)
					d:=wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw,reps)
					metrics["cross_family_evaluation_count"]++
					if len(d)!=3 {
						metrics["invalid_row_count"]++
						continue
					}
					left,ok0:=valueMap[d[0]]
					right,ok1:=valueMap[d[1]]
					op,ok2:=opMap[d[2]]
					if !ok0||!ok1||!ok2 {
						metrics["invalid_row_count"]++
						continue
					}
					if !transition.has(op,left,right) {
						metrics["cross_family_missing_transition_count"]++
					}
					pl,pr:=transition.predict(op,left,right)
					psl,okL:=semanticOfSurface[pl]
					psr,okR:=semanticOfSurface[pr]
					tl,tr:=wlmSiRawRepAliasInvarianceFalsificationR1Apply(leftSemantic,rightSemantic,opSemantic)
					if okL&&okR&&psl==tl&&psr==tr {
						crossHits++
					}
				}
			}
		}
	}

	if metrics["same_family_evaluation_count"]>0 {
		metrics["same_family_accuracy"]=float64(sameHits)/metrics["same_family_evaluation_count"]
	}
	if metrics["cross_family_evaluation_count"]>0 {
		metrics["cross_family_accuracy"]=float64(crossHits)/metrics["cross_family_evaluation_count"]
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) {
			metrics["invalid_row_count"]++
		}
	}
	return wlmSiRawRepAliasInvarianceFalsificationR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-SI-RAW-REP-ALIAS-INVARIANCE-FALSIFICATION-R1",
		Metrics:metrics,
	}
}
