package unitary

import "math"

type wlmSiRawRepPhaseAnchorCoordinateNecessityR14Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func RunWlmSiRawRepPhaseAnchorCoordinateNecessityR14() interface{} {
	m := map[string]float64{"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "target_surface_count": 0, "target_context_coordinate_count": 5, "differing_target_context_coordinate_count": 0, "single_coordinate_arm_count": 5, "construction_mismatch_count": 0, "targeted_rescue108_novel_support_addition_count": 0, "coordinate_arm_novel_support_addition_count_total": 0, "coordinate_arm_destabilization_count": 0, "coordinate_arm_min_alias_class_count": math.MaxInt, "coordinate_arm_max_alias_class_count": 0, "coordinate_arm_correct_alias_triplet_count_total": 0, "coordinate_arm_split_entity_triplet_count_total": 0, "coordinate_arm_cross_entity_collision_pair_count_total": 0, "coordinate_arm_targeted_membership_match_count": 0, "coordinate_arm_all_sham_membership_match_count": 0, "evaluation_entity_pair_count_total": 216, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0}
	budget107 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(107)
	actual108 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(108)
	sham := wlmSiRawRepPhaseAnchorFinalContextCausalityR11Sham(budget107)
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	invalid := 0
	target := -1
	for surface := 0; surface < 27; surface++ {
		candidate := make([][]uint8, len(sham))
		copy(candidate, sham)
		candidate[surface*108+107] = actual108[surface*108+107]
		novel, _, bad := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, candidate, reps)
		invalid += bad
		if novel > 0 {
			m["target_surface_count"]++
			target = surface
		}
	}
	targeted := make([][]uint8, len(sham))
	for i := range sham {
		targeted[i] = append([]uint8{}, sham[i]...)
	}
	if target >= 0 {
		targeted[target*108+107] = append([]uint8{}, actual108[target*108+107]...)
	} else {
		m["construction_mismatch_count"]++
	}
	invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("targeted_rescue108", targeted, reps, m)
	novel, _, bad := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, targeted, reps)
	invalid += bad
	m["targeted_rescue108_novel_support_addition_count"] = float64(novel)
	targetedCounts, bad1 := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(targeted, reps)
	shamCounts, bad2 := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(sham, reps)
	invalid += bad1 + bad2
	targetedClasses := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(targetedCounts)
	shamClasses := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(shamCounts)
	offsets := []int{3, 4, 5, 10, 11}
	if target >= 0 {
		actualRecord := actual108[target*108+107]
		shamRecord := sham[target*108+107]
		for _, offset := range offsets {
			if actualRecord[offset] != shamRecord[offset] {
				m["differing_target_context_coordinate_count"]++
			}
			arm := make([][]uint8, len(targeted))
			for i := range targeted {
				arm[i] = append([]uint8{}, targeted[i]...)
			}
			arm[target*108+107][offset] = shamRecord[offset]
			diffs := 0
			for i := range actualRecord {
				if arm[target*108+107][i] != actualRecord[i] {
					diffs++
					if i != offset || arm[target*108+107][i] != shamRecord[i] {
						m["construction_mismatch_count"]++
					}
				}
			}
			if diffs != 1 {
				m["construction_mismatch_count"]++
			}
			am := map[string]float64{}
			invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics("arm", arm, reps, am)
			n, _, b := wlmSiRawRepPhaseAnchorFinalContextCausalityR11NovelSupport(budget107, arm, reps)
			invalid += b
			m["coordinate_arm_novel_support_addition_count_total"] += float64(n)
			classes := am["arm_alias_class_count"]
			if classes < m["coordinate_arm_min_alias_class_count"] {
				m["coordinate_arm_min_alias_class_count"] = classes
			}
			if classes > m["coordinate_arm_max_alias_class_count"] {
				m["coordinate_arm_max_alias_class_count"] = classes
			}
			m["coordinate_arm_correct_alias_triplet_count_total"] += am["arm_correct_alias_triplet_count"]
			m["coordinate_arm_split_entity_triplet_count_total"] += am["arm_split_entity_triplet_count"]
			m["coordinate_arm_cross_entity_collision_pair_count_total"] += am["arm_cross_entity_collision_pair_count"]
			if classes != 9 || am["arm_correct_alias_triplet_count"] != 9 || am["arm_split_entity_triplet_count"] != 0 {
				m["coordinate_arm_destabilization_count"]++
			}
			ac, ab := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(arm, reps)
			invalid += ab
			armClasses := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(ac)
			if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(armClasses, targetedClasses) == 0 {
				m["coordinate_arm_targeted_membership_match_count"]++
			}
			if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(armClasses, shamClasses) == 0 {
				m["coordinate_arm_all_sham_membership_match_count"]++
			}
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepPhaseAnchorCoordinateNecessityR14Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-COORDINATE-NECESSITY-R14", m}
}
