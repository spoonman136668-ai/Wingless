package unitary

import (
	"bytes"
	"math"
)

type wlmSiRawRepPhaseAnchorFinalRecordLocalityR12Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorFinalRecordLocalityR12Arm(actual108, budget107 [][]uint8, surface int) ([][]uint8, int) {
	arm := make([][]uint8, len(actual108))
	for i := range actual108 {
		arm[i] = append([]uint8{}, actual108[i]...)
	}
	arm[surface*108+107] = append([]uint8{}, budget107[surface*107+106]...)
	mismatch := 0
	if !bytes.Equal(arm[surface*108+107], budget107[surface*107+106]) {
		mismatch++
	}
	for i := range arm {
		if i != surface*108+107 && !bytes.Equal(arm[i], actual108[i]) {
			mismatch++
		}
	}
	return arm, mismatch
}

func RunWlmSiRawRepPhaseAnchorFinalRecordLocalityR12() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"budget107_actual108_prefix_mismatch_count": 0, "actual108_novel_support_addition_count": 0,
		"single_replacement_arm_count": 27, "replacement_mismatch_count": 0,
		"novel_removal_intervention_count": 0, "destabilizing_intervention_count": 0,
		"novel_removal_without_destabilization_count": 0, "nonnovel_removal_with_destabilization_count": 0,
		"destabilized_arm_min_alias_class_count": 0, "destabilized_arm_min_split_entity_triplet_count": 0,
		"stable_arm_class_or_membership_mismatch_count": 0, "intervention_cross_entity_collision_pair_count": 0,
		"evaluation_entity_pair_count_total": 972, "task_training_record_count": 0, "tokenizer_use_count": 0,
		"external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	budget107 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(107)
	actual108 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(108)
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	invalid := wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("budget107", budget107, reps, m)
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("actual108", actual108, reps, m)
	m["budget107_actual108_prefix_mismatch_count"] = float64(wlmSiRawRepPhaseAnchorSampleBudgetR7PrefixMismatch(budget107, actual108, 107, 108))
	actualNovel, _, noveltyInvalid := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, actual108, reps)
	m["actual108_novel_support_addition_count"] = float64(actualNovel)
	invalid += noveltyInvalid
	minClasses, minSplits := math.MaxInt, math.MaxInt
	for surface := 0; surface < 27; surface++ {
		arm, mismatch := wlmSiRawRepPhaseAnchorFinalRecordLocalityR12Arm(actual108, budget107, surface)
		m["replacement_mismatch_count"] += float64(mismatch)
		armMetrics := map[string]float64{}
		invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("arm", arm, reps, armMetrics)
		removedNovel, _, armNoveltyInvalid := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(arm, actual108, reps)
		invalid += armNoveltyInvalid
		novelRemoval := removedNovel > 0
		if novelRemoval {
			m["novel_removal_intervention_count"]++
		}
		stable := armMetrics["arm_alias_class_count"] == 9 && armMetrics["arm_three_member_class_count"] == 9 && armMetrics["arm_correct_alias_triplet_count"] == 9 && armMetrics["arm_split_entity_triplet_count"] == 0 && armMetrics["arm_cross_entity_collision_pair_count"] == 0
		if !stable {
			m["destabilizing_intervention_count"]++
			if int(armMetrics["arm_alias_class_count"]) < minClasses {
				minClasses = int(armMetrics["arm_alias_class_count"])
			}
			if int(armMetrics["arm_split_entity_triplet_count"]) < minSplits {
				minSplits = int(armMetrics["arm_split_entity_triplet_count"])
			}
		}
		if novelRemoval && stable {
			m["novel_removal_without_destabilization_count"]++
		}
		if !novelRemoval && !stable {
			m["nonnovel_removal_with_destabilization_count"]++
			m["stable_arm_class_or_membership_mismatch_count"]++
		}
		m["intervention_cross_entity_collision_pair_count"] += armMetrics["arm_cross_entity_collision_pair_count"]
	}
	if minClasses != math.MaxInt {
		m["destabilized_arm_min_alias_class_count"] = float64(minClasses)
	}
	if minSplits != math.MaxInt {
		m["destabilized_arm_min_split_entity_triplet_count"] = float64(minSplits)
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepPhaseAnchorFinalRecordLocalityR12Result{Schema: "wingless.research-scientific-result.v1", Experiment: "WLM-SI-RAW-REP-PHASE-ANCHOR-FINAL-RECORD-LOCALITY-R12", Metrics: m}
}
