package unitary

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
)

type wlmLmExternalRawRepPredR1Source struct {
	domain string
	data []byte
	sha256 string
	bytes int
}

type wlmLmExternalRawRepPredR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalRawRepPredR1SHA256(data []byte) string {
	sum:=sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func RunWlmLmExternalRawRepPredR1(code, structured, prose []byte) interface{} {
	sources:=[]wlmLmExternalRawRepPredR1Source{
		{domain:"code",data:code,sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
		{domain:"structured",data:structured,sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
		{domain:"technical_prose",data:prose,sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
	}
	metrics:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":float64(len(sources)),
		"total_source_bytes":0,
		"valid_evaluation_source_count":0,
		"selected_motif_count":0,
		"minimum_eval_representation_coverage_fraction":1,
		"minimum_covered_top1_accuracy_gain":1,
		"minimum_effective_event_reduction_fraction":1,
		"maximum_representation_minus_baseline_bits_per_byte":-100,
		"mean_representation_minus_baseline_bits_per_byte":0,
		"sources_with_positive_covered_top1_gain":0,
		"sources_with_nonworse_bits_per_byte":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"invalid_byte_rows":0,
		"counter_overflow_rows":0,
	}

	train:=make([][]byte,0,len(sources))
	evals:=make([][]byte,0,len(sources))
	for _,src:=range sources {
		metrics["total_source_bytes"]+=float64(len(src.data))
		if len(src.data)!=src.bytes || wlmLmExternalRawRepPredR1SHA256(src.data)!=src.sha256 {
			metrics["source_identity_mismatch_count"]++
		}
		split:=len(src.data)*3/5
		if split<5 || len(src.data)-split<5 {
			metrics["invalid_byte_rows"]++
		}
		train=append(train,src.data[:split])
		evals=append(evals,src.data[split:])
	}

	baseline,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(train,metrics)
	metrics["selected_motif_count"]=float64(len(selected))
	bpbSum:=0.0

	for i,src:=range sources {
		coverage,gain,reduction,bpbDelta:=wlmLmRawRepPredFreshHoldoutR1Evaluate(evals[i],&baseline,selected)
		metrics["eval_"+src.domain+"_representation_coverage_fraction"]=coverage
		metrics["eval_"+src.domain+"_covered_top1_accuracy_gain"]=gain
		metrics["eval_"+src.domain+"_effective_event_reduction_fraction"]=reduction
		metrics["eval_"+src.domain+"_representation_minus_baseline_bits_per_byte"]=bpbDelta
		metrics["minimum_eval_representation_coverage_fraction"]=wlmLmRawRepPredFreshHoldoutR1Min(metrics["minimum_eval_representation_coverage_fraction"],coverage)
		metrics["minimum_covered_top1_accuracy_gain"]=wlmLmRawRepPredFreshHoldoutR1Min(metrics["minimum_covered_top1_accuracy_gain"],gain)
		metrics["minimum_effective_event_reduction_fraction"]=wlmLmRawRepPredFreshHoldoutR1Min(metrics["minimum_effective_event_reduction_fraction"],reduction)
		metrics["maximum_representation_minus_baseline_bits_per_byte"]=wlmLmRawRepPredFreshHoldoutR1Max(metrics["maximum_representation_minus_baseline_bits_per_byte"],bpbDelta)
		bpbSum+=bpbDelta
		if gain>0 { metrics["sources_with_positive_covered_top1_gain"]++ }
		if bpbDelta<=0 { metrics["sources_with_nonworse_bits_per_byte"]++ }
		metrics["valid_evaluation_source_count"]++
	}
	metrics["mean_representation_minus_baseline_bits_per_byte"]=bpbSum/float64(len(sources))

	for _,value:=range metrics {
		if math.IsNaN(value)||math.IsInf(value,0) {
			metrics["invalid_byte_rows"]++
		}
	}

	return wlmLmExternalRawRepPredR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-RAW-REP-PRED-R1",
		Metrics:metrics,
	}
}
