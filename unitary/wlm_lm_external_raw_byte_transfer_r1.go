package unitary

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
)

type wlmLmExternalRawByteTransferR1Source struct {
	domain string
	data []byte
	sha256 string
	bytes int
}

type wlmLmExternalRawByteTransferR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalRawByteTransferR1SHA256(data []byte) string {
	sum:=sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func wlmLmExternalRawByteTransferR1Min(a,b float64) float64 {
	if b<a { return b }
	return a
}

func wlmLmExternalRawByteTransferR1Max(a,b float64) float64 {
	if b>a { return b }
	return a
}

// RunWlmLmExternalRawByteTransferR1 evaluates the frozen external raw-byte
// transfer experiment. Inputs must already have been fetched and verified by
// the hosted orchestration layer; the learner itself performs no network I/O.
func RunWlmLmExternalRawByteTransferR1(code, structured, prose []byte) interface{} {
	sources:=[]wlmLmExternalRawByteTransferR1Source{
		{
			domain:"code",
			data:code,
			sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",
			bytes:41453,
		},
		{
			domain:"structured",
			data:structured,
			sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",
			bytes:14365,
		},
		{
			domain:"technical_prose",
			data:prose,
			sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",
			bytes:1454,
		},
	}

	metrics:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":float64(len(sources)),
		"total_source_bytes":0,
		"cross_domain_arm_count":0,
		"within_domain_control_arm_count":0,
		"minimum_selected_motif_count_cross_domain":math.MaxFloat64,
		"minimum_selected_motif_count_within_domain":math.MaxFloat64,
		"mean_cross_domain_motif_coverage_fraction":0,
		"mean_cross_domain_motif_covered_accuracy_gain":0,
		"positive_cross_domain_gain_arm_count":0,
		"mean_cross_domain_adaptive_minus_baseline_bits_per_byte":0,
		"minimum_within_domain_motif_coverage_fraction":1,
		"minimum_within_domain_motif_covered_accuracy_gain":1,
		"maximum_total_learned_motif_capacity":0,
		"capacity_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"invalid_byte_rows":0,
		"counter_overflow_rows":0,
	}

	for _,src:=range sources {
		metrics["total_source_bytes"]+=float64(len(src.data))
		if len(src.data)!=src.bytes || wlmLmExternalRawByteTransferR1SHA256(src.data)!=src.sha256 {
			metrics["source_identity_mismatch_count"]++
		}
	}

	crossCoverageSum:=0.0
	crossGainSum:=0.0
	crossBpbDeltaSum:=0.0

	for targetIndex,target:=range sources {
		train:=make([][]byte,0,2)
		for i,src:=range sources {
			if i==targetIndex { continue }
			train=append(train,src.data)
		}
		baseline,_,selected:=wlmLmRealUTF8BytePilotR2TrainModel(train,metrics)
		selectedCount:=float64(len(selected))
		metrics["minimum_selected_motif_count_cross_domain"]=wlmLmExternalRawByteTransferR1Min(metrics["minimum_selected_motif_count_cross_domain"],selectedCount)
		metrics["maximum_total_learned_motif_capacity"]=wlmLmExternalRawByteTransferR1Max(metrics["maximum_total_learned_motif_capacity"],selectedCount)
		coverage,gain,reduction,delta:=wlmLmRealUTF8BytePilotR2EvalFile(target.data,&baseline,selected)
		metrics["cross_domain_"+target.domain+"_motif_coverage_fraction"]=coverage
		metrics["cross_domain_"+target.domain+"_motif_covered_accuracy_gain"]=gain
		metrics["cross_domain_"+target.domain+"_event_reduction_fraction"]=reduction
		metrics["cross_domain_"+target.domain+"_adaptive_minus_baseline_bits_per_byte"]=delta
		crossCoverageSum+=coverage
		crossGainSum+=gain
		crossBpbDeltaSum+=delta
		if gain>0 { metrics["positive_cross_domain_gain_arm_count"]++ }
		metrics["cross_domain_arm_count"]++
	}

	metrics["mean_cross_domain_motif_coverage_fraction"]=crossCoverageSum/3.0
	metrics["mean_cross_domain_motif_covered_accuracy_gain"]=crossGainSum/3.0
	metrics["mean_cross_domain_adaptive_minus_baseline_bits_per_byte"]=crossBpbDeltaSum/3.0

	for _,src:=range sources {
		mid:=len(src.data)/2
		if mid<4 || len(src.data)-mid<4 {
			metrics["invalid_byte_rows"]++
			continue
		}
		train:=[][]byte{src.data[:mid]}
		eval:=src.data[mid:]
		baseline,_,selected:=wlmLmRealUTF8BytePilotR2TrainModel(train,metrics)
		selectedCount:=float64(len(selected))
		metrics["minimum_selected_motif_count_within_domain"]=wlmLmExternalRawByteTransferR1Min(metrics["minimum_selected_motif_count_within_domain"],selectedCount)
		metrics["maximum_total_learned_motif_capacity"]=wlmLmExternalRawByteTransferR1Max(metrics["maximum_total_learned_motif_capacity"],selectedCount)
		coverage,gain,reduction,delta:=wlmLmRealUTF8BytePilotR2EvalFile(eval,&baseline,selected)
		metrics["within_domain_"+src.domain+"_motif_coverage_fraction"]=coverage
		metrics["within_domain_"+src.domain+"_motif_covered_accuracy_gain"]=gain
		metrics["within_domain_"+src.domain+"_event_reduction_fraction"]=reduction
		metrics["within_domain_"+src.domain+"_adaptive_minus_baseline_bits_per_byte"]=delta
		metrics["minimum_within_domain_motif_coverage_fraction"]=wlmLmExternalRawByteTransferR1Min(metrics["minimum_within_domain_motif_coverage_fraction"],coverage)
		metrics["minimum_within_domain_motif_covered_accuracy_gain"]=wlmLmExternalRawByteTransferR1Min(metrics["minimum_within_domain_motif_covered_accuracy_gain"],gain)
		metrics["within_domain_control_arm_count"]++
	}

	for _,value:=range metrics {
		if math.IsNaN(value)||math.IsInf(value,0) {
			metrics["invalid_byte_rows"]++
		}
	}

	return wlmLmExternalRawByteTransferR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-RAW-BYTE-TRANSFER-R1",
		Metrics:metrics,
	}
}
