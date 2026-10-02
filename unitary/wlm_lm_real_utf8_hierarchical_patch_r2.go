package unitary

import (
	"crypto/sha1"
	"fmt"
	"math"
	"os/exec"
	"sort"
)

type wlmLmRealUTF8HierarchicalPatchR2File struct{ path, sha string }
type wlmLmRealUTF8HierarchicalPatchR2Stats struct{ total uint32; counts [256]uint32 }
type wlmLmRealUTF8HierarchicalPatchR2L1 struct{ key [4]uint8; stats *wlmLmRealUTF8HierarchicalPatchR2Stats; best uint8; bestCount uint32; consistency float64 }
type wlmLmRealUTF8HierarchicalPatchR2L2 struct{ key [8]uint8; stats *wlmLmRealUTF8HierarchicalPatchR2Stats; best uint8; bestCount uint32; consistency float64 }
type wlmLmRealUTF8HierarchicalPatchR2Result struct{ Schema string `json:"schema"`; Experiment string `json:"experiment"`; Metrics map[string]float64 `json:"metrics"` }

var wlmLmRealUTF8HierarchicalPatchR2Train=[]wlmLmRealUTF8HierarchicalPatchR2File{
	{"README.md","50f1ab5ae468551ee48dd1143e534a0ebe01d6a8"},
	{"docs/architecture.md","dec5b6284ec69c1c94a4cf5b005022981df0e900"},
	{"docs/internal-stage-2.md","d71357155dc1b141369dfd61d971b0eb1d00c9c9"},
	{"docs/inference-backends.md","4e5b3c2c268e1ded129020fd22e8fbe07725ec30"},
	{"unitary/wlm_lm_byte_motif_discovery_r3_balanced.go","f9ce41417b5000e3c2419b1b5bf0ce2fc4c41899"},
	{"unitary/wlm_lm_cross_domain_granularity_transfer_r1.go","07b3a6fe2dec98b09c91f770f8d5648defc49ae7"},
	{"unitary/adaptive_continuous_fusion.go","60ae9ce88419cf2541205ad41d81b8c0d3cbe114"},
	{"unitary/anonymous_gram.go","8a85e5e09a0ffa4a0e8e849df1a63e78ecfd2a1e"},
	{"unitary/memory.go","24b132e0bf81e373d7b82b70719b9ee9419d0054"},
	{"unitary/training.go","f39dcd4cda2f2b78b97df77fc7098199afa55acb"},
}
var wlmLmRealUTF8HierarchicalPatchR2EvalFiles=[]wlmLmRealUTF8HierarchicalPatchR2File{
	{"docs/internal-stage-3.md","08dc1153ed401c23194dd5bdd6bb06ff5881c3ab"},
	{"docs/integration-seams.md","3b368fdb1c97dc5a24e41d437718c82f1c4b6ebc"},
	{"unitary/blind_channel.go","7323002862869b8073189b944fc645e4956a8f1b"},
	{"unitary/composite_state.go","29653a18fdf9eae88e95f10dd969835e3307b743"},
}

func wlmLmRealUTF8HierarchicalPatchR2Hash(data []byte) string {
	h:=sha1.New();_,_=h.Write([]byte(fmt.Sprintf("blob %d\x00",len(data))));_,_=h.Write(data);return fmt.Sprintf("%x",h.Sum(nil))
}
func wlmLmRealUTF8HierarchicalPatchR2Load(f wlmLmRealUTF8HierarchicalPatchR2File)([]byte,bool){
	data,err:=exec.Command("git","cat-file","blob",f.sha).Output();if err!=nil{return nil,false};return data,wlmLmRealUTF8HierarchicalPatchR2Hash(data)==f.sha
}
func wlmLmRealUTF8HierarchicalPatchR2Best(c *[256]uint32)(uint8,uint32,uint64){
	best:=uint8(0);bc:=c[0];tot:=uint64(c[0]);for i:=1;i<256;i++{tot+=uint64(c[i]);if c[i]>bc{best=uint8(i);bc=c[i]}};return best,bc,tot
}
func wlmLmRealUTF8HierarchicalPatchR2Less4(a,b [4]uint8)bool{for i:=0;i<4;i++{if a[i]!=b[i]{return a[i]<b[i]}};return false}
func wlmLmRealUTF8HierarchicalPatchR2Less8(a,b [8]uint8)bool{for i:=0;i<8;i++{if a[i]!=b[i]{return a[i]<b[i]}};return false}
func wlmLmRealUTF8HierarchicalPatchR2Prob(c *[256]uint32,target uint8)float64{_,_,tot:=wlmLmRealUTF8HierarchicalPatchR2Best(c);return(float64(c[target])+0.5)/(float64(tot)+128)}
func wlmLmRealUTF8HierarchicalPatchR2BaseProb(base *[256][256]uint32,prev,target uint8)float64{var tot uint64;for i:=0;i<256;i++{tot+=uint64(base[prev][i])};return(float64(base[prev][target])+0.5)/(float64(tot)+128)}
func wlmLmRealUTF8HierarchicalPatchR2BasePred(base *[256][256]uint32,prev uint8)uint8{best:=uint8(0);bc:=base[prev][0];for i:=1;i<256;i++{if base[prev][i]>bc{best=uint8(i);bc=base[prev][i]}};return best}

func wlmLmRealUTF8HierarchicalPatchR2TrainModel(files [][]byte,metrics map[string]float64)([256][256]uint32,map[[4]uint8]wlmLmRealUTF8HierarchicalPatchR2L1,map[[8]uint8]wlmLmRealUTF8HierarchicalPatchR2L2){
	var base [256][256]uint32
	c1:=make(map[[4]uint8]*wlmLmRealUTF8HierarchicalPatchR2Stats)
	for _,data:=range files{
		for i:=1;i<len(data);i++{if base[data[i-1]][data[i]]==^uint32(0){metrics["counter_overflow_rows"]++}else{base[data[i-1]][data[i]]++}}
		for i:=4;i<len(data);i++{
			k:=[4]uint8{data[i-4],data[i-3],data[i-2],data[i-1]};s:=c1[k];if s==nil{s=&wlmLmRealUTF8HierarchicalPatchR2Stats{};c1[k]=s}
			if s.total==^uint32(0)||s.counts[data[i]]==^uint32(0){metrics["counter_overflow_rows"]++;continue};s.total++;s.counts[data[i]]++
		}
	}
	rows1:=make([]wlmLmRealUTF8HierarchicalPatchR2L1,0)
	for k,s:=range c1{if s.total<4{continue};b,bc,_:=wlmLmRealUTF8HierarchicalPatchR2Best(&s.counts);rows1=append(rows1,wlmLmRealUTF8HierarchicalPatchR2L1{k,s,b,bc,float64(bc)/float64(s.total)})}
	sort.Slice(rows1,func(i,j int)bool{if rows1[i].bestCount!=rows1[j].bestCount{return rows1[i].bestCount>rows1[j].bestCount};if rows1[i].consistency!=rows1[j].consistency{return rows1[i].consistency>rows1[j].consistency};return wlmLmRealUTF8HierarchicalPatchR2Less4(rows1[i].key,rows1[j].key)})
	if len(rows1)>512{rows1=rows1[:512]}
	l1:=make(map[[4]uint8]wlmLmRealUTF8HierarchicalPatchR2L1,len(rows1));for _,r:=range rows1{l1[r.key]=r}

	c2:=make(map[[8]uint8]*wlmLmRealUTF8HierarchicalPatchR2Stats)
	for _,data:=range files{
		for i:=8;i<len(data);i++{
			a:=[4]uint8{data[i-8],data[i-7],data[i-6],data[i-5]};b:=[4]uint8{data[i-4],data[i-3],data[i-2],data[i-1]}
			if _,ok:=l1[a];!ok{continue};if _,ok:=l1[b];!ok{continue}
			k:=[8]uint8{data[i-8],data[i-7],data[i-6],data[i-5],data[i-4],data[i-3],data[i-2],data[i-1]}
			s:=c2[k];if s==nil{s=&wlmLmRealUTF8HierarchicalPatchR2Stats{};c2[k]=s}
			if s.total==^uint32(0)||s.counts[data[i]]==^uint32(0){metrics["counter_overflow_rows"]++;continue};s.total++;s.counts[data[i]]++
		}
	}
	rows2:=make([]wlmLmRealUTF8HierarchicalPatchR2L2,0)
	for k,s:=range c2{if s.total<2{continue};b,bc,_:=wlmLmRealUTF8HierarchicalPatchR2Best(&s.counts);rows2=append(rows2,wlmLmRealUTF8HierarchicalPatchR2L2{k,s,b,bc,float64(bc)/float64(s.total)})}
	sort.Slice(rows2,func(i,j int)bool{if rows2[i].bestCount!=rows2[j].bestCount{return rows2[i].bestCount>rows2[j].bestCount};if rows2[i].consistency!=rows2[j].consistency{return rows2[i].consistency>rows2[j].consistency};return wlmLmRealUTF8HierarchicalPatchR2Less8(rows2[i].key,rows2[j].key)})
	if len(rows2)>128{rows2=rows2[:128]}
	l2:=make(map[[8]uint8]wlmLmRealUTF8HierarchicalPatchR2L2,len(rows2));for _,r:=range rows2{l2[r.key]=r}
	return base,l1,l2
}

func wlmLmRealUTF8HierarchicalPatchR2EventReduction(data []byte,l1 map[[4]uint8]wlmLmRealUTF8HierarchicalPatchR2L1,l2 map[[8]uint8]wlmLmRealUTF8HierarchicalPatchR2L2)float64{
	if len(data)==0{return 0};events:=0
	for i:=0;i<len(data);{
		if i+8<=len(data){k:=[8]uint8{data[i],data[i+1],data[i+2],data[i+3],data[i+4],data[i+5],data[i+6],data[i+7]};if _,ok:=l2[k];ok{events++;i+=8;continue}}
		if i+4<=len(data){k:=[4]uint8{data[i],data[i+1],data[i+2],data[i+3]};if _,ok:=l1[k];ok{events++;i+=4;continue}}
		events++;i++
	}
	return 1-float64(events)/float64(len(data))
}
func wlmLmRealUTF8HierarchicalPatchR2Min(a,b float64)float64{if b<a{return b};return a}
func wlmLmRealUTF8HierarchicalPatchR2Max(a,b float64)float64{if b>a{return b};return a}

func wlmLmRealUTF8HierarchicalPatchR2EvaluateFile(data []byte,base *[256][256]uint32,l1 map[[4]uint8]wlmLmRealUTF8HierarchicalPatchR2L1,l2 map[[8]uint8]wlmLmRealUTF8HierarchicalPatchR2L2)(cov1,cov2,gain,reduction,hMinusL1,hMinusBase float64){
	if len(data)<2{return}
	c1,c2:=0,0;l1Corr,l2Corr:=0,0;total:=0
	var baseBits,l1Bits,hBits float64
	for i:=1;i<len(data);i++{
		target:=data[i];bp:=wlmLmRealUTF8HierarchicalPatchR2BaseProb(base,data[i-1],target);p1:=bp;ph:=bp
		pred1:=wlmLmRealUTF8HierarchicalPatchR2BasePred(base,data[i-1]);predh:=pred1
		if i>=4{k1:=[4]uint8{data[i-4],data[i-3],data[i-2],data[i-1]};if m,ok:=l1[k1];ok{c1++;p1=wlmLmRealUTF8HierarchicalPatchR2Prob(&m.stats.counts,target);ph=p1;pred1=m.best;predh=m.best}}
		if i>=8{k2:=[8]uint8{data[i-8],data[i-7],data[i-6],data[i-5],data[i-4],data[i-3],data[i-2],data[i-1]};if m,ok:=l2[k2];ok{c2++;ph=wlmLmRealUTF8HierarchicalPatchR2Prob(&m.stats.counts,target);predh=m.best;if pred1==target{l1Corr++};if predh==target{l2Corr++}}}
		baseBits+=-math.Log2(bp);l1Bits+=-math.Log2(p1);hBits+=-math.Log2(ph);total++
	}
	den1:=len(data)-4;if den1<1{den1=1};den2:=len(data)-8;if den2<1{den2=1}
	cov1=float64(c1)/float64(den1);cov2=float64(c2)/float64(den2);if c2>0{gain=float64(l2Corr-l1Corr)/float64(c2)}
	reduction=wlmLmRealUTF8HierarchicalPatchR2EventReduction(data,l1,l2)
	hMinusL1=hBits/float64(total)-l1Bits/float64(total);hMinusBase=hBits/float64(total)-baseBits/float64(total);return
}

// RunWlmLmRealUTF8HierarchicalPatchR2 evaluates the frozen fresh-holdout two-level raw-byte hierarchy.
func RunWlmLmRealUTF8HierarchicalPatchR2() interface{} {
	metrics:=map[string]float64{
		"training_file_identity_mismatch_count":0,"evaluation_file_identity_mismatch_count":0,"train_eval_blob_overlap_count":0,
		"selected_level1_count":0,"selected_level2_count":0,"minimum_eval_level1_coverage_fraction":1,"minimum_eval_level2_coverage_fraction":1,
		"minimum_level2_covered_accuracy_gain":1,"minimum_effective_event_reduction_fraction":1,"maximum_hierarchical_minus_level1_bits_per_byte":-100,
		"maximum_hierarchical_minus_baseline_bits_per_byte":-100,"maximum_total_learned_structure_count":0,"capacity_growth_event_count":0,
		"tokenizer_use_count":0,"invalid_byte_rows":0,"counter_overflow_rows":0,
	}
	trainSet:=map[string]bool{};for _,f:=range wlmLmRealUTF8HierarchicalPatchR2Train{trainSet[f.sha]=true}
	for _,f:=range wlmLmRealUTF8HierarchicalPatchR2EvalFiles{if trainSet[f.sha]{metrics["train_eval_blob_overlap_count"]++}}
	train:=make([][]byte,0,len(wlmLmRealUTF8HierarchicalPatchR2Train))
	for _,f:=range wlmLmRealUTF8HierarchicalPatchR2Train{d,ok:=wlmLmRealUTF8HierarchicalPatchR2Load(f);if !ok{metrics["training_file_identity_mismatch_count"]++};train=append(train,d)}
	base,l1,l2:=wlmLmRealUTF8HierarchicalPatchR2TrainModel(train,metrics)
	metrics["selected_level1_count"]=float64(len(l1));metrics["selected_level2_count"]=float64(len(l2));metrics["maximum_total_learned_structure_count"]=float64(len(l1)+len(l2))
	for _,f:=range wlmLmRealUTF8HierarchicalPatchR2EvalFiles{
		d,ok:=wlmLmRealUTF8HierarchicalPatchR2Load(f);if !ok{metrics["evaluation_file_identity_mismatch_count"]++}
		c1,c2,g,r,dl1,db:=wlmLmRealUTF8HierarchicalPatchR2EvaluateFile(d,&base,l1,l2)
		metrics["minimum_eval_level1_coverage_fraction"]=wlmLmRealUTF8HierarchicalPatchR2Min(metrics["minimum_eval_level1_coverage_fraction"],c1)
		metrics["minimum_eval_level2_coverage_fraction"]=wlmLmRealUTF8HierarchicalPatchR2Min(metrics["minimum_eval_level2_coverage_fraction"],c2)
		metrics["minimum_level2_covered_accuracy_gain"]=wlmLmRealUTF8HierarchicalPatchR2Min(metrics["minimum_level2_covered_accuracy_gain"],g)
		metrics["minimum_effective_event_reduction_fraction"]=wlmLmRealUTF8HierarchicalPatchR2Min(metrics["minimum_effective_event_reduction_fraction"],r)
		metrics["maximum_hierarchical_minus_level1_bits_per_byte"]=wlmLmRealUTF8HierarchicalPatchR2Max(metrics["maximum_hierarchical_minus_level1_bits_per_byte"],dl1)
		metrics["maximum_hierarchical_minus_baseline_bits_per_byte"]=wlmLmRealUTF8HierarchicalPatchR2Max(metrics["maximum_hierarchical_minus_baseline_bits_per_byte"],db)
	}
	return wlmLmRealUTF8HierarchicalPatchR2Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-REAL-UTF8-HIERARCHICAL-PATCH-R2",Metrics:metrics}
}
