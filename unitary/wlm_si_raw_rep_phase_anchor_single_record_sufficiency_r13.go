package unitary

import (
	"bytes"
	"math"
)

type wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(a, b wlmLmRawRepContextualAliasInductionR2Classes) int {
	mismatch := 0
	for i := 0; i < len(a.surfaceToClass); i++ {
		for j := i + 1; j < len(a.surfaceToClass); j++ {
			if (a.surfaceToClass[i] == a.surfaceToClass[j]) != (b.surfaceToClass[i] == b.surfaceToClass[j]) {
				mismatch++
			}
		}
	}
	return mismatch
}

func RunWlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"budget107_actual108_prefix_mismatch_count": 0, "actual108_novel_support_addition_count": 0,
		"target_surface_count": 0, "targeted_actual_final_record_count": 0, "targeted_duplicate_final_record_count": 0,
		"construction_mismatch_count": 0, "all_sham108_novel_support_addition_count": 0,
		"targeted_rescue108_novel_support_addition_count": 0, "targeted_actual108_class_membership_mismatch_count": 0,
		"evaluation_entity_pair_count_total": 144, "task_training_record_count": 0, "tokenizer_use_count": 0,
		"external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	budget107 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(107)
	actual108 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(108)
	allSham108 := wlmSiRawRepPhaseAnchorFinalContextCausalityR11Sham(budget107)
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	invalid := wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("budget107", budget107, reps, m)
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("actual108", actual108, reps, m)
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("all_sham108", allSham108, reps, m)
	m["budget107_actual108_prefix_mismatch_count"] = float64(wlmSiRawRepPhaseAnchorSampleBudgetR7PrefixMismatch(budget107, actual108, 107, 108))
	actualNovel, _, actualInvalid := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, actual108, reps)
	shamNovel, _, shamInvalid := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, allSham108, reps)
	m["actual108_novel_support_addition_count"] = float64(actualNovel)
	m["all_sham108_novel_support_addition_count"] = float64(shamNovel)
	invalid += actualInvalid + shamInvalid
	targeted := make([][]uint8, len(allSham108))
	for i := range allSham108 {
		targeted[i] = append([]uint8{}, allSham108[i]...)
	}
	target := -1
	for surface := 0; surface < 27; surface++ {
		candidate := make([][]uint8, len(allSham108))
		copy(candidate, allSham108)
		candidate[surface*108+107] = actual108[surface*108+107]
		novel, _, candidateInvalid := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, candidate, reps)
		invalid += candidateInvalid
		if novel > 0 {
			m["target_surface_count"]++
			target = surface
		}
	}
	if target >= 0 {
		targeted[target*108+107] = append([]uint8{}, actual108[target*108+107]...)
	}
	for surface := 0; surface < 27; surface++ {
		final := targeted[surface*108+107]
		if surface == target {
			m["targeted_actual_final_record_count"]++
			if !bytes.Equal(final, actual108[surface*108+107]) {
				m["construction_mismatch_count"]++
			}
		} else {
			m["targeted_duplicate_final_record_count"]++
			if !bytes.Equal(final, budget107[surface*107+106]) {
				m["construction_mismatch_count"]++
			}
		}
	}
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("targeted_rescue108", targeted, reps, m)
	targetedNovel, _, targetedInvalid := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, targeted, reps)
	m["targeted_rescue108_novel_support_addition_count"] = float64(targetedNovel)
	invalid += targetedInvalid
	actualCounts, actualCountsInvalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(actual108, reps)
	targetedCounts, targetedCountsInvalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(targeted, reps)
	invalid += actualCountsInvalid + targetedCountsInvalid
	m["targeted_actual108_class_membership_mismatch_count"] = float64(wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(actualCounts), wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(targetedCounts)))
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-SINGLE-RECORD-SUFFICIENCY-R13", m}
}
