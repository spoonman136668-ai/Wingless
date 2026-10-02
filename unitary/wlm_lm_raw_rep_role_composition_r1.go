package unitary

import (
	"math"
	"sort"
)

type wlmLmRawRepRoleCompositionR1Key [4]uint8

type wlmLmRawRepRoleCompositionR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmRawRepRoleCompositionR1Sample struct {
	raw []uint8
}

type wlmLmRawRepRoleCompositionR1Transition struct {
	counts [5][4][4][16]uint32
}

type wlmLmRawRepRoleCompositionR1ProgramKey struct {
	left int
	right int
	ops string
}

type wlmLmRawRepRoleCompositionR1ProgramLookup struct {
	counts map[wlmLmRawRepRoleCompositionR1ProgramKey]*[16]uint32
}

func wlmLmRawRepRoleCompositionR1ValueMotif(id int) wlmLmRawRepRoleCompositionR1Key {
	base:=uint8(16+16*id)
	return wlmLmRawRepRoleCompositionR1Key{base,base+1,base+2,base+3}
}

func wlmLmRawRepRoleCompositionR1OpMotif(id int) wlmLmRawRepRoleCompositionR1Key {
	base:=uint8(96+8*id)
	return wlmLmRawRepRoleCompositionR1Key{base,base+1,base+2,base+3}
}

func wlmLmRawRepRoleCompositionR1Noise(state *uint32,n int) []uint8 {
	out:=make([]uint8,n)
	for i:=0;i<n;i++ {
		*state=*state*1664525+1013904223
		out[i]=uint8(192+((*state>>16)%64))
	}
	return out
}

func wlmLmRawRepRoleCompositionR1Apply(left,right,op int)(int,int) {
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

func wlmLmRawRepRoleCompositionR1Target(left,right int,ops []int)(int,int) {
	for _,op:=range ops {
		left,right=wlmLmRawRepRoleCompositionR1Apply(left,right,op)
	}
	return left,right
}

func wlmLmRawRepRoleCompositionR1AppendMotif(out []uint8,k wlmLmRawRepRoleCompositionR1Key) []uint8 {
	return append(out,k[:]...)
}

func wlmLmRawRepRoleCompositionR1TrainingRecord(left,right,op,repeat int) []uint8 {
	outLeft,outRight:=wlmLmRawRepRoleCompositionR1Apply(left,right,op)
	state:=uint32(0x85ebca6b ^ uint32((left+1)*131+(right+1)*977+(op+1)*8191+(repeat+1)*65537))
	out:=make([]uint8,0,90)
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,5+(repeat+left)%5)...)
	out=wlmLmRawRepRoleCompositionR1AppendMotif(out,wlmLmRawRepRoleCompositionR1ValueMotif(left))
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,7+(repeat+right)%5)...)
	out=wlmLmRawRepRoleCompositionR1AppendMotif(out,wlmLmRawRepRoleCompositionR1ValueMotif(right))
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,9+(repeat+op)%5)...)
	out=wlmLmRawRepRoleCompositionR1AppendMotif(out,wlmLmRawRepRoleCompositionR1OpMotif(op))
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,11+(repeat+left+right)%5)...)
	out=wlmLmRawRepRoleCompositionR1AppendMotif(out,wlmLmRawRepRoleCompositionR1ValueMotif(outLeft))
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,13+(repeat+op+right)%5)...)
	out=wlmLmRawRepRoleCompositionR1AppendMotif(out,wlmLmRawRepRoleCompositionR1ValueMotif(outRight))
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,5+(repeat+op+left)%5)...)
	return out
}

func wlmLmRawRepRoleCompositionR1ProgramRecord(left,right int,ops []int,tag int) []uint8 {
	state:=uint32(0xc2b2ae35 ^ uint32((left+1)*313+(right+1)*1999+(tag+1)*65599+len(ops)*8191))
	out:=make([]uint8,0,120)
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,6+(left+tag)%5)...)
	out=wlmLmRawRepRoleCompositionR1AppendMotif(out,wlmLmRawRepRoleCompositionR1ValueMotif(left))
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,8+(right+tag)%5)...)
	out=wlmLmRawRepRoleCompositionR1AppendMotif(out,wlmLmRawRepRoleCompositionR1ValueMotif(right))
	for i,op:=range ops {
		out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,7+(i+op+tag)%7)...)
		out=wlmLmRawRepRoleCompositionR1AppendMotif(out,wlmLmRawRepRoleCompositionR1OpMotif(op))
	}
	out=append(out,wlmLmRawRepRoleCompositionR1Noise(&state,9+(tag+len(ops))%7)...)
	return out
}

func wlmLmRawRepRoleCompositionR1Less(a,b wlmLmRawRepRoleCompositionR1Key) bool {
	for i:=0;i<4;i++ {
		if a[i]!=b[i] {
			return a[i]<b[i]
		}
	}
	return false
}

func wlmLmRawRepRoleCompositionR1LearnRepresentation(records [][]uint8) []wlmLmRawRepRoleCompositionR1Key {
	counts:=make(map[wlmLmRawRepRoleCompositionR1Key]uint32)
	for _,data:=range records {
		for i:=0;i+4<=len(data);i++ {
			k:=wlmLmRawRepRoleCompositionR1Key{data[i],data[i+1],data[i+2],data[i+3]}
			counts[k]++
		}
	}
	type row struct{key wlmLmRawRepRoleCompositionR1Key;count uint32}
	rows:=make([]row,0,len(counts))
	for k,c:=range counts {
		rows=append(rows,row{k,c})
	}
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count {
			return rows[i].count>rows[j].count
		}
		return wlmLmRawRepRoleCompositionR1Less(rows[i].key,rows[j].key)
	})
	if len(rows)>9 {
		rows=rows[:9]
	}
	out:=make([]wlmLmRawRepRoleCompositionR1Key,len(rows))
	for i,r:=range rows {
		out[i]=r.key
	}
	sort.Slice(out,func(i,j int)bool{return wlmLmRawRepRoleCompositionR1Less(out[i],out[j])})
	return out
}

func wlmLmRawRepRoleCompositionR1Decode(data []uint8,reps []wlmLmRawRepRoleCompositionR1Key) []int {
	index:=make(map[wlmLmRawRepRoleCompositionR1Key]int,len(reps))
	for i,k:=range reps {
		index[k]=i
	}
	out:=make([]int,0,10)
	for i:=0;i+4<=len(data); {
		k:=wlmLmRawRepRoleCompositionR1Key{data[i],data[i+1],data[i+2],data[i+3]}
		if id,ok:=index[k];ok {
			out=append(out,id)
			i+=4
			continue
		}
		i++
	}
	return out
}

func wlmLmRawRepRoleCompositionR1PairIndex(left,right int) int {
	return left*4+right
}

func wlmLmRawRepRoleCompositionR1PairDecode(index int)(int,int) {
	return index/4,index%4
}

func (t *wlmLmRawRepRoleCompositionR1Transition) observe(op,left,right,outLeft,outRight int,metrics map[string]float64) {
	idx:=wlmLmRawRepRoleCompositionR1PairIndex(outLeft,outRight)
	if t.counts[op][left][right][idx]==^uint32(0) {
		metrics["counter_overflow_count"]++
		return
	}
	t.counts[op][left][right][idx]++
}

func (t *wlmLmRawRepRoleCompositionR1Transition) predict(op,left,right int)(int,int) {
	best:=0
	bestCount:=t.counts[op][left][right][0]
	for i:=1;i<16;i++ {
		if t.counts[op][left][right][i]>bestCount {
			best=i
			bestCount=t.counts[op][left][right][i]
		}
	}
	return wlmLmRawRepRoleCompositionR1PairDecode(best)
}

func wlmLmRawRepRoleCompositionR1OpsKey(ops []int) string {
	b:=make([]byte,len(ops))
	for i,op:=range ops {
		b[i]=byte(op)
	}
	return string(b)
}

func (l *wlmLmRawRepRoleCompositionR1ProgramLookup) observe(left,right int,ops []int,outLeft,outRight int,metrics map[string]float64) {
	if l.counts==nil {
		l.counts=make(map[wlmLmRawRepRoleCompositionR1ProgramKey]*[16]uint32)
	}
	k:=wlmLmRawRepRoleCompositionR1ProgramKey{left:left,right:right,ops:wlmLmRawRepRoleCompositionR1OpsKey(ops)}
	row:=l.counts[k]
	if row==nil {
		row=&[16]uint32{}
		l.counts[k]=row
	}
	idx:=wlmLmRawRepRoleCompositionR1PairIndex(outLeft,outRight)
	if row[idx]==^uint32(0) {
		metrics["counter_overflow_count"]++
		return
	}
	row[idx]++
}

func (l *wlmLmRawRepRoleCompositionR1ProgramLookup) predict(left,right int,ops []int)(int,int) {
	k:=wlmLmRawRepRoleCompositionR1ProgramKey{left:left,right:right,ops:wlmLmRawRepRoleCompositionR1OpsKey(ops)}
	row:=l.counts[k]
	if row==nil {
		return 0,0
	}
	best:=0
	bestCount:=row[0]
	for i:=1;i<16;i++ {
		if row[i]>bestCount {
			best=i
			bestCount=row[i]
		}
	}
	return wlmLmRawRepRoleCompositionR1PairDecode(best)
}

func wlmLmRawRepRoleCompositionR1TrueMotifCount(reps []wlmLmRawRepRoleCompositionR1Key) int {
	count:=0
	for _,r:=range reps {
		matched:=false
		for i:=0;i<4;i++ {
			if r==wlmLmRawRepRoleCompositionR1ValueMotif(i) {
				matched=true
			}
		}
		for i:=0;i<5;i++ {
			if r==wlmLmRawRepRoleCompositionR1OpMotif(i) {
				matched=true
			}
		}
		if matched {
			count++
		}
	}
	return count
}

func wlmLmRawRepRoleCompositionR1ClassifyReps(reps []wlmLmRawRepRoleCompositionR1Key)([]int,[]int) {
	values:=make([]int,0,4)
	ops:=make([]int,0,5)
	for id,r:=range reps {
		for i:=0;i<4;i++ {
			if r==wlmLmRawRepRoleCompositionR1ValueMotif(i) {
				values=append(values,id)
			}
		}
		for i:=0;i<5;i++ {
			if r==wlmLmRawRepRoleCompositionR1OpMotif(i) {
				ops=append(ops,id)
			}
		}
	}
	sort.Ints(values)
	sort.Ints(ops)
	return values,ops
}

func wlmLmRawRepRoleCompositionR1IndexMap(ids []int) map[int]int {
	out:=make(map[int]int,len(ids))
	for i,id:=range ids {
		out[id]=i
	}
	return out
}

func wlmLmRawRepRoleCompositionR1LongProgram(left,right,seed,length int) []int {
	ops:=make([]int,length)
	for i:=range ops {
		if (i+seed+left+right)%3==0 {
			ops[i]=4
		} else {
			ops[i]=(i*i+3*i+seed+2*left+right)%4
		}
	}
	return ops
}

func wlmLmRawRepRoleCompositionR1EvaluateOne(
	left,right int,
	semanticOps []int,
	raw []uint8,
	reps []wlmLmRawRepRoleCompositionR1Key,
	valueMap,opMap map[int]int,
	valueIDs []int,
	transition *wlmLmRawRepRoleCompositionR1Transition,
	lookup *wlmLmRawRepRoleCompositionR1ProgramLookup,
	metrics map[string]float64,
)(bool,bool,bool) {
	decoded:=wlmLmRawRepRoleCompositionR1Decode(raw,reps)
	if len(decoded)!=2+len(semanticOps) {
		metrics["heldout_decode_failure_count"]++
		return false,false,false
	}
	dl,okL:=valueMap[decoded[0]]
	dr,okR:=valueMap[decoded[1]]
	if !okL||!okR {
		metrics["heldout_decode_failure_count"]++
		return false,false,false
	}
	decodedOps:=make([]int,len(semanticOps))
	for i:=range semanticOps {
		op,ok:=opMap[decoded[2+i]]
		if !ok {
			metrics["heldout_decode_failure_count"]++
			return false,false,false
		}
		decodedOps[i]=op
	}
	pl,pr:=dl,dr
	for _,op:=range decodedOps {
		pl,pr=transition.predict(op,pl,pr)
	}
	ll,lr:=lookup.predict(dl,dr,decodedOps)
	tl,tr:=wlmLmRawRepRoleCompositionR1Target(left,right,semanticOps)
	learnedCorrect:=pl==tl&&pr==tr
	lookupCorrect:=ll==tl&&lr==tr
	rawCorrect:=false
	if pl>=0&&pl<len(valueIDs)&&pr>=0&&pr<len(valueIDs) {
		leftKey:=reps[valueIDs[pl]]
		rightKey:=reps[valueIDs[pr]]
		rawCorrect=leftKey==wlmLmRawRepRoleCompositionR1ValueMotif(tl)&&rightKey==wlmLmRawRepRoleCompositionR1ValueMotif(tr)
	}
	return learnedCorrect,lookupCorrect,rawCorrect
}

// RunWlmLmRawRepRoleCompositionR1 tests raw representation feeding learned multi-step role composition.
func RunWlmLmRawRepRoleCompositionR1() interface{} {
	metrics:=map[string]float64{
		"training_record_count":0,
		"selected_representation_count":0,
		"selected_true_motif_match_count":0,
		"value_representation_count":0,
		"operation_representation_count":0,
		"training_decode_failure_count":0,
		"primitive_transition_missing_count":0,
		"heldout_program_count":0,
		"heldout_decode_failure_count":0,
		"conjugated_right_accuracy":0,
		"long_program_accuracy":0,
		"overall_multistep_accuracy":0,
		"whole_program_lookup_accuracy":0,
		"raw_output_pair_exact_accuracy":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	samples:=make([]wlmLmRawRepRoleCompositionR1Sample,0,3200)
	records:=make([][]uint8,0,3200)
	for repeat:=0;repeat<40;repeat++ {
		for op:=0;op<5;op++ {
			for left:=0;left<4;left++ {
				for right:=0;right<4;right++ {
					raw:=wlmLmRawRepRoleCompositionR1TrainingRecord(left,right,op,repeat)
					samples=append(samples,wlmLmRawRepRoleCompositionR1Sample{raw:raw})
					records=append(records,raw)
					metrics["training_record_count"]++
				}
			}
		}
	}

	reps:=wlmLmRawRepRoleCompositionR1LearnRepresentation(records)
	metrics["selected_representation_count"]=float64(len(reps))
	metrics["selected_true_motif_match_count"]=float64(wlmLmRawRepRoleCompositionR1TrueMotifCount(reps))
	valueIDs,opIDs:=wlmLmRawRepRoleCompositionR1ClassifyReps(reps)
	metrics["value_representation_count"]=float64(len(valueIDs))
	metrics["operation_representation_count"]=float64(len(opIDs))
	valueMap:=wlmLmRawRepRoleCompositionR1IndexMap(valueIDs)
	opMap:=wlmLmRawRepRoleCompositionR1IndexMap(opIDs)

	var transition wlmLmRawRepRoleCompositionR1Transition
	var lookup wlmLmRawRepRoleCompositionR1ProgramLookup
	for _,sample:=range samples {
		d:=wlmLmRawRepRoleCompositionR1Decode(sample.raw,reps)
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
		lookup.observe(left,right,[]int{op},outLeft,outRight,metrics)
	}

	for op:=0;op<5;op++ {
		for left:=0;left<4;left++ {
			for right:=0;right<4;right++ {
				total:=uint64(0)
				for i:=0;i<16;i++ {
					total+=uint64(transition.counts[op][left][right][i])
				}
				if total==0 {
					metrics["primitive_transition_missing_count"]++
				}
			}
		}
	}

	conjHits:=0
	conjTotal:=0
	longHits:=0
	longTotal:=0
	lookupHits:=0
	rawHits:=0
	overallHits:=0
	overallTotal:=0

	for left:=0;left<4;left++ {
		for right:=0;right<4;right++ {
			for start:=0;start<4;start++ {
				ops:=[]int{4,start,4}
				raw:=wlmLmRawRepRoleCompositionR1ProgramRecord(left,right,ops,100+left*100+right*10+start)
				ok,lk,rk:=wlmLmRawRepRoleCompositionR1EvaluateOne(left,right,ops,raw,reps,valueMap,opMap,valueIDs,&transition,&lookup,metrics)
				if ok {conjHits++;overallHits++}
				if lk {lookupHits++}
				if rk {rawHits++}
				conjTotal++;overallTotal++;metrics["heldout_program_count"]++
			}
		}
	}

	for _,length:=range []int{5,7} {
		for left:=0;left<4;left++ {
			for right:=0;right<4;right++ {
				for seed:=0;seed<4;seed++ {
					ops:=wlmLmRawRepRoleCompositionR1LongProgram(left,right,seed,length)
					raw:=wlmLmRawRepRoleCompositionR1ProgramRecord(left,right,ops,1000+length*100+left*20+right*4+seed)
					ok,lk,rk:=wlmLmRawRepRoleCompositionR1EvaluateOne(left,right,ops,raw,reps,valueMap,opMap,valueIDs,&transition,&lookup,metrics)
					if ok {longHits++;overallHits++}
					if lk {lookupHits++}
					if rk {rawHits++}
					longTotal++;overallTotal++;metrics["heldout_program_count"]++
				}
			}
		}
	}

	if conjTotal>0 {
		metrics["conjugated_right_accuracy"]=float64(conjHits)/float64(conjTotal)
	}
	if longTotal>0 {
		metrics["long_program_accuracy"]=float64(longHits)/float64(longTotal)
	}
	if overallTotal>0 {
		metrics["overall_multistep_accuracy"]=float64(overallHits)/float64(overallTotal)
		metrics["whole_program_lookup_accuracy"]=float64(lookupHits)/float64(overallTotal)
		metrics["raw_output_pair_exact_accuracy"]=float64(rawHits)/float64(overallTotal)
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) {
			metrics["invalid_row_count"]++
		}
	}
	return wlmLmRawRepRoleCompositionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-REP-ROLE-COMPOSITION-R1",
		Metrics:metrics,
	}
}
