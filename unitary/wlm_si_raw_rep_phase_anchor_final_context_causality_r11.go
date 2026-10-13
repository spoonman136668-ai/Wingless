package unitary

import (
	"bytes"
	"math"
)

type wlmSiRawRepPhaseAnchorFinalContextCausalityR11Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorFinalContextCausalityR11Sham(records107 [][]uint8) [][]uint8 {
	out := make([][]uint8, 0, 27*108)
	for surface := 0; surface < 27; surface++ {
		start := surface * 107
		for sample := 0; sample < 107; sample++ {
			out = append(out, append([]uint8{}, records107[start+sample]...))
		}
		out = append(out, append([]uint8{}, records107[start+106]...))
	}
	return out
}

func wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(before, after [][]uint8, reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) (int, int, int) {
	beforeCounts, beforeInvalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(before, reps)
	afterCounts, afterInvalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(after, reps)
	novel, surfaces := 0, 0
	for i := range beforeCounts {
		delta := 0
		for context := range afterCounts[i] {
			if beforeCounts[i][context] == 0 {
				delta++
			}
		}
		novel += delta
		if delta > 0 {
			surfaces++
		}
	}
	return novel, surfaces, beforeInvalid + afterInvalid
}

func RunWlmSiRawRepPhaseAnchorFinalContextCausalityR11() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"budget107_actual108_prefix_mismatch_count": 0, "sham_prefix_mismatch_count": 0, "sham_duplicate_mismatch_count": 0,
		"actual108_novel_support_addition_count": 0, "actual108_restoring_surface_count": 0, "sham108_novel_support_addition_count": 0,
		"evaluation_entity_pair_count_per_arm": 36, "task_training_record_count": 0, "tokenizer_use_count": 0,
		"external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	budget107 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(107)
	actual108 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(108)
	sham108 := wlmSiRawRepPhaseAnchorFinalContextCausalityR11Sham(budget107)
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	invalid := wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("budget107", budget107, reps, m)
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("actual108", actual108, reps, m)
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("sham108", sham108, reps, m)
	m["budget107_actual108_prefix_mismatch_count"] = float64(wlmSiRawRepPhaseAnchorSampleBudgetR7PrefixMismatch(budget107, actual108, 107, 108))
	m["sham_prefix_mismatch_count"] = float64(wlmSiRawRepPhaseAnchorSampleBudgetR7PrefixMismatch(budget107, sham108, 107, 108))
	for surface := 0; surface < 27; surface++ {
		if !bytes.Equal(sham108[surface*108+107], budget107[surface*107+106]) {
			m["sham_duplicate_mismatch_count"]++
		}
	}
	actualNovel, restoringSurfaces, actualInvalid := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, actual108, reps)
	shamNovel, _, shamInvalid := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, sham108, reps)
	m["actual108_novel_support_addition_count"] = float64(actualNovel)
	m["actual108_restoring_surface_count"] = float64(restoringSurfaces)
	m["sham108_novel_support_addition_count"] = float64(shamNovel)
	invalid += actualInvalid + shamInvalid
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepPhaseAnchorFinalContextCausalityR11Result{
		Schema: "wingless.research-scientific-result.v1", Experiment: "WLM-SI-RAW-REP-PHASE-ANCHOR-FINAL-CONTEXT-CAUSALITY-R11", Metrics: m,
	}
}
