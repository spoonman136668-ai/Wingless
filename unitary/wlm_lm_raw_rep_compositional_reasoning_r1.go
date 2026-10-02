package unitary

import (
	"math"
	"sort"
)

type wlmLmRawRepCompositionalReasoningR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmRawRepCompositionalReasoningR1Key [4]uint8

type wlmLmRawRepCompositionalReasoningR1Factor struct {
	counts [3][2][2][2]uint32
}

type wlmLmRawRepCompositionalReasoningR1Lookup struct {
	counts [8][8][8]uint32
}

func wlmLmRawRepCompositionalReasoningR1Motif(id int) wlmLmRawRepCompositionalReasoningR1Key {
	return wlmLmRawRepCompositionalReasoningR1Key{
		uint8(16 + 8*id),
		uint8(48 + 3*id),
		uint8(80 + 5*id),
		uint8(112 + 7*id),
	}
}

func wlmLmRawRepCompositionalReasoningR1Noise(state *uint32, n int) []uint8 {
	out:=make([]uint8,n)
	for i:=0;i<n;i++ {
		*state=*state*1664525+1013904223
		out[i]=uint8(192+(*state%64))
	}
	return out
}

func wlmLmRawRepCompositionalReasoningR1Record(a,b,target,repeat int,includeTarget bool) []uint8 {
	state:=uint32(0x9e3779b9 ^ uint32((a+1)*131+(b+1)*977+(target+1)*8191+(repeat+1)*65537))
	out:=make([]uint8,0,50)
	out=append(out,wlmLmRawRepCompositionalReasoningR1Noise(&state,7)...)
	ma:=wlmLmRawRepCompositionalReasoningR1Motif(a);out=append(out,ma[:]...)
	out=append(out,wlmLmRawRepCompositionalReasoningR1Noise(&state,11)...)
	mb:=wlmLmRawRepCompositionalReasoningR1Motif(b);out=append(out,mb[:]...)
	out=append(out,wlmLmRawRepCompositionalReasoningR1Noise(&state,11)...)
	if includeTarget {
		mt:=wlmLmRawRepCompositionalReasoningR1Motif(target);out=append(out,mt[:]...)
	}
	return out
}

func wlmLmRawRepCompositionalReasoningR1Heldout(a,b int) bool {
	return (a+b)%4==1
}

func wlmLmRawRepCompositionalReasoningR1Less(a,b wlmLmRawRepCompositionalReasoningR1Key) bool {
	for i:=0;i<4;i++ {
		if a[i]!=b[i] { return a[i]<b[i] }
	}
	return false
}

func wlmLmRawRepCompositionalReasoningR1LearnRepresentation(records [][]uint8) []wlmLmRawRepCompositionalReasoningR1Key {
	counts:=make(map[wlmLmRawRepCompositionalReasoningR1Key]uint32)
	for _,data:=range records {
		for i:=0;i+4<=len(data);i++ {
			key:=wlmLmRawRepCompositionalReasoningR1Key{data[i],data[i+1],data[i+2],data[i+3]}
			counts[key]++
		}
	}
	type row struct{ key wlmLmRawRepCompositionalReasoningR1Key; count uint32 }
	rows:=make([]row,0,len(counts))
	for k,c:=range counts { rows=append(rows,row{k,c}) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count { return rows[i].count>rows[j].count }
		return wlmLmRawRepCompositionalReasoningR1Less(rows[i].key,rows[j].key)
	})
	if len(rows)>8 { rows=rows[:8] }
	out:=make([]wlmLmRawRepCompositionalReasoningR1Key,len(rows))
	for i,r:=range rows { out[i]=r.key }
	sort.Slice(out,func(i,j int)bool{return wlmLmRawRepCompositionalReasoningR1Less(out[i],out[j])})
	return out
}

func wlmLmRawRepCompositionalReasoningR1Decode(data []uint8, reps []wlmLmRawRepCompositionalReasoningR1Key) []int {
	ids:=make(map[wlmLmRawRepCompositionalReasoningR1Key]int,len(reps))
	for i,k:=range reps { ids[k]=i }
	out:=make([]int,0,3)
	for i:=0;i+4<=len(data); {
		k:=wlmLmRawRepCompositionalReasoningR1Key{data[i],data[i+1],data[i+2],data[i+3]}
		if id,ok:=ids[k];ok {
			out=append(out,id)
			i+=4
			continue
		}
		i++
	}
	return out
}

func (f *wlmLmRawRepCompositionalReasoningR1Factor) observe(a,b,t int,metrics map[string]float64) {
	for bit:=0;bit<3;bit++ {
		x:=(a>>bit)&1;y:=(b>>bit)&1;z:=(t>>bit)&1
		if f.counts[bit][x][y][z]==^uint32(0) { metrics["counter_overflow_count"]++;continue }
		f.counts[bit][x][y][z]++
	}
}

func (f *wlmLmRawRepCompositionalReasoningR1Factor) predict(a,b int) int {
	out:=0
	for bit:=0;bit<3;bit++ {
		x:=(a>>bit)&1;y:=(b>>bit)&1
		if f.counts[bit][x][y][1]>f.counts[bit][x][y][0] { out|=1<<bit }
	}
	return out
}

func (l *wlmLmRawRepCompositionalReasoningR1Lookup) observe(a,b,t int,metrics map[string]float64) {
	if l.counts[a][b][t]==^uint32(0) { metrics["counter_overflow_count"]++;return }
	l.counts[a][b][t]++
}

func (l *wlmLmRawRepCompositionalReasoningR1Lookup) predict(a,b int) int {
	best:=0;bc:=l.counts[a][b][0]
	for t:=1;t<8;t++ {
		if l.counts[a][b][t]>bc { best=t;bc=l.counts[a][b][t] }
	}
	return best
}

func wlmLmRawRepCompositionalReasoningR1CoverageMissing(training [][2]int) int {
	missing:=0
	for bit:=0;bit<3;bit++ {
		var seen [2][2]bool
		for _,p:=range training { seen[(p[0]>>bit)&1][(p[1]>>bit)&1]=true }
		for x:=0;x<2;x++ { for y:=0;y<2;y++ { if !seen[x][y] { missing++ } } }
	}
	return missing
}

func wlmLmRawRepCompositionalReasoningR1KeyEqual(a,b wlmLmRawRepCompositionalReasoningR1Key) bool {
	for i:=0;i<4;i++ { if a[i]!=b[i] { return false } }
	return true
}

// RunWlmLmRawRepCompositionalReasoningR1 tests raw representation feeding held-out compositional inference.
func RunWlmLmRawRepCompositionalReasoningR1() interface{} {
	metrics:=map[string]float64{
		"training_record_count":0,
		"heldout_pair_count":0,
		"selected_representation_count":0,
		"selected_true_motif_match_count":0,
		"training_representation_decode_failure_count":0,
		"heldout_representation_decode_failure_count":0,
		"heldout_context_leak_count":0,
		"factorized_bit_coverage_missing_count":0,
		"factorized_heldout_accuracy":0,
		"lookup_heldout_accuracy":0,
		"raw_output_exact_accuracy":0,
		"factorized_counter_count":24,
		"lookup_counter_count":512,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}
	records:=make([][]uint8,0,1536)
	for repeat:=0;repeat<32;repeat++ {
		for a:=0;a<8;a++ {
			for b:=0;b<8;b++ {
				if wlmLmRawRepCompositionalReasoningR1Heldout(a,b) { continue }
				t:=a^b
				records=append(records,wlmLmRawRepCompositionalReasoningR1Record(a,b,t,repeat,true))
				metrics["training_record_count"]++
			}
		}
	}
	reps:=wlmLmRawRepCompositionalReasoningR1LearnRepresentation(records)
	metrics["selected_representation_count"]=float64(len(reps))
	for _,r:=range reps {
		for id:=0;id<8;id++ {
			if wlmLmRawRepCompositionalReasoningR1KeyEqual(r,wlmLmRawRepCompositionalReasoningR1Motif(id)) {
				metrics["selected_true_motif_match_count"]++
				break
			}
		}
	}

	var factor wlmLmRawRepCompositionalReasoningR1Factor
	var lookup wlmLmRawRepCompositionalReasoningR1Lookup
	trainingPairs:=make([][2]int,0,48)
	seenTrainingPair:=make(map[[2]int]bool)
	for repeat:=0;repeat<32;repeat++ {
		for a:=0;a<8;a++ {
			for b:=0;b<8;b++ {
				if wlmLmRawRepCompositionalReasoningR1Heldout(a,b) { continue }
				t:=a^b
				decoded:=wlmLmRawRepCompositionalReasoningR1Decode(wlmLmRawRepCompositionalReasoningR1Record(a,b,t,repeat,true),reps)
				if len(decoded)!=3 {
					metrics["training_representation_decode_failure_count"]++
					continue
				}
				if repeat==0 { trainingPairs=append(trainingPairs,[2]int{decoded[0],decoded[1]}) }
				seenTrainingPair[[2]int{decoded[0],decoded[1]}]=true
				factor.observe(decoded[0],decoded[1],decoded[2],metrics)
				lookup.observe(decoded[0],decoded[1],decoded[2],metrics)
			}
		}
	}
	metrics["factorized_bit_coverage_missing_count"]=float64(wlmLmRawRepCompositionalReasoningR1CoverageMissing(trainingPairs))

	factorHits:=0;lookupHits:=0;rawHits:=0;validHeldout:=0
	for a:=0;a<8;a++ {
		for b:=0;b<8;b++ {
			if !wlmLmRawRepCompositionalReasoningR1Heldout(a,b) { continue }
			metrics["heldout_pair_count"]++
			t:=a^b
			decoded:=wlmLmRawRepCompositionalReasoningR1Decode(wlmLmRawRepCompositionalReasoningR1Record(a,b,t,1000+a*8+b,false),reps)
			if len(decoded)!=2 {
				metrics["heldout_representation_decode_failure_count"]++
				continue
			}
			if seenTrainingPair[[2]int{decoded[0],decoded[1]}] { metrics["heldout_context_leak_count"]++ }
			fp:=factor.predict(decoded[0],decoded[1])
			lp:=lookup.predict(decoded[0],decoded[1])
			if fp==t { factorHits++ }
			if lp==t { lookupHits++ }
			if fp>=0 && fp<len(reps) && wlmLmRawRepCompositionalReasoningR1KeyEqual(reps[fp],wlmLmRawRepCompositionalReasoningR1Motif(t)) { rawHits++ }
			validHeldout++
		}
	}
	if validHeldout>0 {
		metrics["factorized_heldout_accuracy"]=float64(factorHits)/float64(validHeldout)
		metrics["lookup_heldout_accuracy"]=float64(lookupHits)/float64(validHeldout)
		metrics["raw_output_exact_accuracy"]=float64(rawHits)/float64(validHeldout)
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) { metrics["invalid_row_count"]++ }
	}
	return wlmLmRawRepCompositionalReasoningR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-REP-COMPOSITIONAL-REASONING-R1",
		Metrics:metrics,
	}
}
