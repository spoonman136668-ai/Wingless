package unitary

import (
	"crypto/sha1"
	"fmt"
	"math"
	"os/exec"
	"sort"
)

type wlmLmRealUTF8BytePilotR2MotifStats struct {
	total uint32
	counts [256]uint32
}

type wlmLmRealUTF8BytePilotR2Selected struct {
	key [4]uint8
	stats *wlmLmRealUTF8BytePilotR2MotifStats
	best uint8
	bestCount uint32
	consistency float64
}

type wlmLmRealUTF8BytePilotR2Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

type wlmLmRealUTF8BytePilotR2CorpusFile struct {
	path string
	sha string
}

var wlmLmRealUTF8BytePilotR2Train = []wlmLmRealUTF8BytePilotR2CorpusFile{
	{"README.md","50f1ab5ae468551ee48dd1143e534a0ebe01d6a8"},
	{"docs/architecture.md","dec5b6284ec69c1c94a4cf5b005022981df0e900"},
	{"unitary/wlm_lm_byte_motif_discovery_r3_balanced.go","f9ce41417b5000e3c2419b1b5bf0ce2fc4c41899"},
	{"unitary/wlm_lm_cross_domain_granularity_transfer_r1.go","07b3a6fe2dec98b09c91f770f8d5648defc49ae7"},
}
var wlmLmRealUTF8BytePilotR2Eval = []wlmLmRealUTF8BytePilotR2CorpusFile{
	{"docs/experiments/wlm-lm-adaptive-granularity-proposal-r1.json","32c1a462863a105ac3a42e7229f2efcb15050f30"},
	{"unitary/wlm_lm_multiscale_native_language_precursor_r2.go","ae9543d3fcd61602c14eae4cb22148610a701e0a"},
}

func wlmLmRealUTF8BytePilotR2GitBlobHash(data []byte) string {
	header:=[]byte(fmt.Sprintf("blob %d\x00",len(data)))
	h:=sha1.New();_,_=h.Write(header);_,_=h.Write(data)
	return fmt.Sprintf("%x",h.Sum(nil))
}

func wlmLmRealUTF8BytePilotR2LoadBlob(f wlmLmRealUTF8BytePilotR2CorpusFile) ([]byte,bool) {
	cmd:=exec.Command("git","cat-file","blob",f.sha)
	data,err:=cmd.Output()
	if err!=nil{return nil,false}
	return data,wlmLmRealUTF8BytePilotR2GitBlobHash(data)==f.sha
}

func wlmLmRealUTF8BytePilotR2Best(counts *[256]uint32)(uint8,uint32,uint64){
	best:=uint8(0);bestCount:=counts[0];var total uint64=uint64(counts[0])
	for i:=1;i<256;i++ {
		c:=counts[i];total+=uint64(c)
		if c>bestCount {best=uint8(i);bestCount=c}
	}
	return best,bestCount,total
}

func wlmLmRealUTF8BytePilotR2Less(a,b [4]uint8) bool {
	for i:=0;i<4;i++ {if a[i]!=b[i]{return a[i]<b[i]}}
	return false
}

func wlmLmRealUTF8BytePilotR2TrainModel(files [][]byte,metrics map[string]float64)([256][256]uint32,map[[4]uint8]*wlmLmRealUTF8BytePilotR2MotifStats,map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected){
	var baseline [256][256]uint32
	candidates:=make(map[[4]uint8]*wlmLmRealUTF8BytePilotR2MotifStats)
	for _,data:=range files {
		for i:=1;i<len(data);i++ {
			prev,next:=data[i-1],data[i]
			if baseline[prev][next]==^uint32(0){metrics["counter_overflow_rows"]++;continue}
			baseline[prev][next]++
		}
		for i:=4;i<len(data);i++ {
			key:=[4]uint8{data[i-4],data[i-3],data[i-2],data[i-1]}
			s:=candidates[key]
			if s==nil{s=&wlmLmRealUTF8BytePilotR2MotifStats{};candidates[key]=s}
			if s.total==^uint32(0)||s.counts[data[i]]==^uint32(0){metrics["counter_overflow_rows"]++;continue}
			s.total++;s.counts[data[i]]++
		}
	}
	rows:=make([]wlmLmRealUTF8BytePilotR2Selected,0)
	for key,s:=range candidates {
		if s.total<4{continue}
		best,bestCount,_:=wlmLmRealUTF8BytePilotR2Best(&s.counts)
		rows=append(rows,wlmLmRealUTF8BytePilotR2Selected{key:key,stats:s,best:best,bestCount:bestCount,consistency:float64(bestCount)/float64(s.total)})
	}
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].bestCount!=rows[j].bestCount{return rows[i].bestCount>rows[j].bestCount}
		if rows[i].consistency!=rows[j].consistency{return rows[i].consistency>rows[j].consistency}
		return wlmLmRealUTF8BytePilotR2Less(rows[i].key,rows[j].key)
	})
	if len(rows)>256{rows=rows[:256]}
	selected:=make(map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,len(rows))
	for _,r:=range rows{selected[r.key]=r}
	return baseline,candidates,selected
}

func wlmLmRealUTF8BytePilotR2BaselinePrediction(baseline *[256][256]uint32,prev uint8)uint8{
	best:=uint8(0);bestCount:=baseline[prev][0]
	for i:=1;i<256;i++ {if baseline[prev][i]>bestCount{best=uint8(i);bestCount=baseline[prev][i]}}
	return best
}

func wlmLmRealUTF8BytePilotR2Prob(counts *[256]uint32,target uint8)float64{
	_,_,total:=wlmLmRealUTF8BytePilotR2Best(counts)
	return (float64(counts[target])+0.5)/(float64(total)+128.0)
}

func wlmLmRealUTF8BytePilotR2BaselineProb(baseline *[256][256]uint32,prev,target uint8)float64{
	var total uint64
	for i:=0;i<256;i++{total+=uint64(baseline[prev][i])}
	return (float64(baseline[prev][target])+0.5)/(float64(total)+128.0)
}

func wlmLmRealUTF8BytePilotR2EventReduction(data []byte,selected map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected)float64{
	if len(data)==0{return 0}
	events:=0
	for i:=0;i<len(data); {
		if i+4<=len(data) {
			key:=[4]uint8{data[i],data[i+1],data[i+2],data[i+3]}
			if _,ok:=selected[key];ok{events++;i+=4;continue}
		}
		events++;i++
	}
	return 1-float64(events)/float64(len(data))
}

func wlmLmRealUTF8BytePilotR2Min(a,b float64)float64{if b<a{return b};return a}
func wlmLmRealUTF8BytePilotR2Max(a,b float64)float64{if b>a{return b};return a}

func wlmLmRealUTF8BytePilotR2EvalFile(data []byte,baseline *[256][256]uint32,selected map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected)(coverage,gain,reduction,bpbDelta float64){
	if len(data)<2{return 0,0,0,0}
	covered:=0;coveredBaseCorrect:=0;coveredAdaptiveCorrect:=0
	var baseBits,adaptBits float64
	total:=0
	for i:=1;i<len(data);i++ {
		target:=data[i];basePred:=wlmLmRealUTF8BytePilotR2BaselinePrediction(baseline,data[i-1])
		baseP:=wlmLmRealUTF8BytePilotR2BaselineProb(baseline,data[i-1],target)
		adaptP:=baseP
		if i>=4 {
			key:=[4]uint8{data[i-4],data[i-3],data[i-2],data[i-1]}
			if m,ok:=selected[key];ok {
				covered++
				if basePred==target{coveredBaseCorrect++}
				if m.best==target{coveredAdaptiveCorrect++}
				adaptP=wlmLmRealUTF8BytePilotR2Prob(&m.stats.counts,target)
			}
		}
		baseBits+=-math.Log2(baseP);adaptBits+=-math.Log2(adaptP);total++
	}
	eligible:=len(data)-4
	if eligible<1{eligible=1}
	coverage=float64(covered)/float64(eligible)
	if covered>0{gain=float64(coveredAdaptiveCorrect-coveredBaseCorrect)/float64(covered)}
	reduction=wlmLmRealUTF8BytePilotR2EventReduction(data,selected)
	bpbDelta=adaptBits/float64(total)-baseBits/float64(total)
	return
}

// RunWlmLmRealUTF8BytePilotR2 evaluates the frozen real UTF-8 raw-byte pilot.
func RunWlmLmRealUTF8BytePilotR2() interface{} {
	metrics:=map[string]float64{
		"training_file_identity_mismatch_count":0,
		"evaluation_file_identity_mismatch_count":0,
		"selected_motif_count":0,
		"minimum_eval_motif_coverage_fraction":1,
		"minimum_motif_covered_accuracy_gain":1,
		"minimum_effective_event_reduction_fraction":1,
		"maximum_adaptive_minus_baseline_bits_per_byte":-100,
		"maximum_learned_motif_capacity":256,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"invalid_byte_rows":0,
		"counter_overflow_rows":0,
	}
	train:=make([][]byte,0,len(wlmLmRealUTF8BytePilotR2Train))
	for _,f:=range wlmLmRealUTF8BytePilotR2Train {
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["training_file_identity_mismatch_count"]++}
		train=append(train,data)
	}
	baseline,_,selected:=wlmLmRealUTF8BytePilotR2TrainModel(train,metrics)
	metrics["selected_motif_count"]=float64(len(selected))
	for _,f:=range wlmLmRealUTF8BytePilotR2Eval {
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["evaluation_file_identity_mismatch_count"]++}
		coverage,gain,reduction,delta:=wlmLmRealUTF8BytePilotR2EvalFile(data,&baseline,selected)
		metrics["minimum_eval_motif_coverage_fraction"]=wlmLmRealUTF8BytePilotR2Min(metrics["minimum_eval_motif_coverage_fraction"],coverage)
		metrics["minimum_motif_covered_accuracy_gain"]=wlmLmRealUTF8BytePilotR2Min(metrics["minimum_motif_covered_accuracy_gain"],gain)
		metrics["minimum_effective_event_reduction_fraction"]=wlmLmRealUTF8BytePilotR2Min(metrics["minimum_effective_event_reduction_fraction"],reduction)
		metrics["maximum_adaptive_minus_baseline_bits_per_byte"]=wlmLmRealUTF8BytePilotR2Max(metrics["maximum_adaptive_minus_baseline_bits_per_byte"],delta)
	}
	return wlmLmRealUTF8BytePilotR2Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-REAL-UTF8-BYTE-PILOT-R2",Metrics:metrics}
}
