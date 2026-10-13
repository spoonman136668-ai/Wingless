package unitary

import (
	"bytes"
	"math"
)

type wlmSiRawRepPhaseAnchorSampleBudgetR7Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(samples int) [][]uint8 {
	records := make([][]uint8, 0, 27*samples)
	for entity := 0; entity < 9; entity++ {
		for family := 0; family < 3; family++ {
			state := uint32(0x243f6a88) ^ uint32((entity+1)*1009) ^ uint32(family+1)*uint32(0x9e3779b9)
			for sample := 0; sample < samples; sample++ {
				state = state*1664525 + 1013904223
				contextIndex := int((state >> 16) % 16)
				noise := uint32(0x85ebca6b) ^ uint32((entity+1)*1009+(family+1)*10007+(sample+1)*65537)
				record := append([]uint8{}, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&noise, 3)...)
				record = append(record, uint8(160+contextIndex))
				record = append(record, uint8(192+(entity*7+contextIndex*3)%16), uint8(208+(entity*5+contextIndex*5)%16))
				motif := wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity)
				record = append(record, motif[:]...)
				record = append(record, uint8(224+(entity*11+contextIndex*7)%16), uint8(240+(entity*13+contextIndex*11)%16))
				record = append(record, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&noise, 3)...)
				records = append(records, record)
			}
		}
	}
	return records
}

func wlmSiRawRepPhaseAnchorSampleBudgetR7PrefixMismatch(shorter, longer [][]uint8, shorterSamples, longerSamples int) int {
	mismatches := 0
	for entityFamily := 0; entityFamily < 27; entityFamily++ {
		for sample := 0; sample < shorterSamples; sample++ {
			shortIndex := entityFamily*shorterSamples + sample
			longIndex := entityFamily*longerSamples + sample
			if !bytes.Equal(shorter[shortIndex], longer[longIndex]) {
				mismatches++
			}
		}
	}
	return mismatches
}

func wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics(prefix string, records [][]uint8, reps []wlmSiRawRepAliasInvarianceFalsificationR1Key, m map[string]float64) int {
	counts, invalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(records, reps)
	classes := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(counts)
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	aliasClasses, threeMember, correct := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(classes, index)
	recordMin, recordMax := wlmSiRawRepPhaseAnchorSamplingR6RecordRange(counts)
	contextMin, contextMax := wlmSiRawRepPhaseAnchorSamplingR6ContextRange(counts)
	m[prefix+"_bridge_record_count"] = float64(len(records))
	m[prefix+"_records_per_surface_min"] = float64(recordMin)
	m[prefix+"_records_per_surface_max"] = float64(recordMax)
	m[prefix+"_distinct_contexts_per_surface_min"] = float64(contextMin)
	m[prefix+"_distinct_contexts_per_surface_max"] = float64(contextMax)
	m[prefix+"_max_same_entity_support_symmetric_difference"] = float64(wlmSiRawRepPhaseAnchorSamplingR6MaxSymmetricDifference(counts, index))
	m[prefix+"_alias_class_count"] = float64(aliasClasses)
	m[prefix+"_three_member_class_count"] = float64(threeMember)
	m[prefix+"_correct_alias_triplet_count"] = float64(correct)
	m[prefix+"_split_entity_triplet_count"] = float64(wlmSiRawRepPhaseAnchorSamplingR6SplitTriplets(classes, index))
	m[prefix+"_cross_entity_collision_pair_count"] = float64(wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(classes, index))
	return invalid
}

func RunWlmSiRawRepPhaseAnchorSampleBudgetR7() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"nested_prefix_mismatch_count": 0, "evaluation_entity_pair_count_per_arm": 36,
		"task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0,
		"capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	exhaustive := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Bridge()
	budget64 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(64)
	budget128 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(128)
	budget256 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(256)
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	invalid := wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("exhaustive", exhaustive, reps, m)
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("budget64", budget64, reps, m)
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("budget128", budget128, reps, m)
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("budget256", budget256, reps, m)
	m["nested_prefix_mismatch_count"] = float64(
		wlmSiRawRepPhaseAnchorSampleBudgetR7PrefixMismatch(budget64, budget128, 64, 128) +
			wlmSiRawRepPhaseAnchorSampleBudgetR7PrefixMismatch(budget128, budget256, 128, 256),
	)
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepPhaseAnchorSampleBudgetR7Result{
		Schema: "wingless.research-scientific-result.v1", Experiment: "WLM-SI-RAW-REP-PHASE-ANCHOR-SAMPLE-BUDGET-R7", Metrics: m,
	}
}
