package unitary

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"os"
	"path/filepath"
)

type wlmLmExternalRawByteGeneralizationR1Source struct {
	domain string
	file string
	sha256 string
	bytes int
}

var wlmLmExternalRawByteGeneralizationR1Sources = []wlmLmExternalRawByteGeneralizationR1Source{
	{domain:"code",file:"code.bin",sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
	{domain:"structured",file:"structured.bin",sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
	{domain:"technical-prose",file:"technical-prose.bin",sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
}

type wlmLmExternalRawByteGeneralizationR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalRawByteGeneralizationR1Load(root string, spec wlmLmExternalRawByteGeneralizationR1Source, metrics map[string]float64) []byte {
	data,err:=os.ReadFile(filepath.Join(root,spec.file))
	if err!=nil {
		metrics["source_identity_mismatch_count"]++
		return nil
	}
	sum:=sha256.Sum256(data)
	if hex.EncodeToString(sum[:])!=spec.sha256 || len(data)!=spec.bytes {
		metrics["source_identity_mismatch_count"]++
	}
	return data
}

func wlmLmExternalRawByteGeneralizationR1Split(data []byte)([]byte,[]byte) {
	n:=len(data)*60/100
	if n<1 || n>=len(data) {
		return nil,nil
	}
	return data[:n],data[n:]
}

func wlmLmExternalRawByteGeneralizationR1Min(a,b float64)float64 {
	if b<a{return b}
	return a
}

func wlmLmExternalRawByteGeneralizationR1Max(a,b float64)float64 {
	if b>a{return b}
	return a
}

// RunWlmLmExternalRawByteGeneralizationR1 evaluates the frozen external raw-byte experiment.
func RunWlmLmExternalRawByteGeneralizationR1(root string) interface{} {
	metrics:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":0,
		"total_source_bytes":0,
		"pooled_selected_motif_count":0,
		"minimum_pooled_motif_coverage_fraction":1,
		"minimum_pooled_motif_covered_accuracy_gain":1,
		"minimum_pooled_effective_event_reduction_fraction":1,
		"maximum_pooled_adaptive_minus_baseline_bits_per_byte":-100,
		"cross_domain_eval_count":0,
		"minimum_cross_domain_motif_coverage_fraction":1,
		"mean_cross_domain_motif_covered_accuracy_gain":0,
		"maximum_cross_domain_adaptive_minus_baseline_bits_per_byte":-100,
		"maximum_learned_motif_capacity":256,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"invalid_byte_rows":0,
		"counter_overflow_rows":0,
	}
	metrics["source_count"]=float64(len(wlmLmExternalRawByteGeneralizationR1Sources))
	data:=make([][]byte,len(wlmLmExternalRawByteGeneralizationR1Sources))
	train:=make([][]byte,len(data))
	held:=make([][]byte,len(data))
	for i,spec:=range wlmLmExternalRawByteGeneralizationR1Sources {
		data[i]=wlmLmExternalRawByteGeneralizationR1Load(root,spec,metrics)
		metrics["total_source_bytes"]+=float64(len(data[i]))
		train[i],held[i]=wlmLmExternalRawByteGeneralizationR1Split(data[i])
		if len(train[i])==0 || len(held[i])==0 {
			metrics["invalid_byte_rows"]++
		}
	}

	baseline,_,selected:=wlmLmRealUTF8BytePilotR2TrainModel(train,metrics)
	metrics["pooled_selected_motif_count"]=float64(len(selected))
	for _,eval:=range held {
		coverage,gain,reduction,delta:=wlmLmRealUTF8BytePilotR2EvalFile(eval,&baseline,selected)
		metrics["minimum_pooled_motif_coverage_fraction"]=wlmLmExternalRawByteGeneralizationR1Min(metrics["minimum_pooled_motif_coverage_fraction"],coverage)
		metrics["minimum_pooled_motif_covered_accuracy_gain"]=wlmLmExternalRawByteGeneralizationR1Min(metrics["minimum_pooled_motif_covered_accuracy_gain"],gain)
		metrics["minimum_pooled_effective_event_reduction_fraction"]=wlmLmExternalRawByteGeneralizationR1Min(metrics["minimum_pooled_effective_event_reduction_fraction"],reduction)
		metrics["maximum_pooled_adaptive_minus_baseline_bits_per_byte"]=wlmLmExternalRawByteGeneralizationR1Max(metrics["maximum_pooled_adaptive_minus_baseline_bits_per_byte"],delta)
	}

	crossGainSum:=0.0
	for heldIndex:=range held {
		crossTrain:=make([][]byte,0,len(train)-1)
		for trainIndex:=range train {
			if trainIndex!=heldIndex {
				crossTrain=append(crossTrain,train[trainIndex])
			}
		}
		crossBaseline,_,crossSelected:=wlmLmRealUTF8BytePilotR2TrainModel(crossTrain,metrics)
		coverage,gain,_,delta:=wlmLmRealUTF8BytePilotR2EvalFile(held[heldIndex],&crossBaseline,crossSelected)
		metrics["minimum_cross_domain_motif_coverage_fraction"]=wlmLmExternalRawByteGeneralizationR1Min(metrics["minimum_cross_domain_motif_coverage_fraction"],coverage)
		metrics["maximum_cross_domain_adaptive_minus_baseline_bits_per_byte"]=wlmLmExternalRawByteGeneralizationR1Max(metrics["maximum_cross_domain_adaptive_minus_baseline_bits_per_byte"],delta)
		crossGainSum+=gain
		metrics["cross_domain_eval_count"]++
	}
	if metrics["cross_domain_eval_count"]>0 {
		metrics["mean_cross_domain_motif_covered_accuracy_gain"]=crossGainSum/metrics["cross_domain_eval_count"]
	}

	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) {
			metrics["invalid_byte_rows"]++
		}
	}

	return wlmLmExternalRawByteGeneralizationR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-RAW-BYTE-GENERALIZATION-R1",
		Metrics:metrics,
	}
}
