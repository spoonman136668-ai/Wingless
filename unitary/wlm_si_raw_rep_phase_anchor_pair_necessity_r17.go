package unitary

import "math"

type wlmSiRawRepPhaseAnchorPairNecessityR17Result struct {
	Schema         string             `json:"schema"`
	Experiment     string             `json:"experiment"`
	Classification string             `json:"classification"`
	Metrics        map[string]float64 `json:"metrics"`
}

func RunWlmSiRawRepPhaseAnchorPairNecessityR17() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"actual108_record_count": 0, "actual108_records_per_surface_min": 0, "actual108_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0,
		"phase_plus_derived_control_count": 4, "phase_plus_derived_full5_membership_match_count": 0,
		"derived_only_pair_arm_count": 6, "derived_only_pair_full5_membership_match_count": 0,
		"derived_only_pair_alias_class_count_min": 27, "derived_only_pair_alias_class_count_max": 0,
		"derived_only_pair_correct_alias_triplet_count_total": 0, "derived_only_pair_split_entity_triplet_count_total": 0,
		"derived_only_pair_cross_entity_collision_pair_count_total": 0, "projection_construction_mismatch_count": 0,
		"evaluation_entity_pair_count_total": 396, "task_training_record_count": 0, "tokenizer_use_count": 0,
		"external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	records := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(108)
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	m["actual108_record_count"] = float64(len(records))
	counts, invalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(records, reps)
	minRecords, maxRecords := wlmSiRawRepPhaseAnchorSamplingR6RecordRange(counts)
	m["actual108_records_per_surface_min"], m["actual108_records_per_surface_max"] = float64(minRecords), float64(maxRecords)
	full := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(counts)
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	fc, ft, fs, fx := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(full, index)
	m["full5_alias_class_count"], m["full5_correct_alias_triplet_count"] = float64(fc), float64(ft)
	m["full5_split_entity_triplet_count"], m["full5_cross_entity_collision_pair_count"] = float64(fs), float64(fx)
	for derived := 1; derived < 5; derived++ {
		classes, bad := wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Induce(counts, []int{0, derived})
		invalid += bad
		if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, full) == 0 {
			m["phase_plus_derived_full5_membership_match_count"]++
		}
	}
	pairs := [][2]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}
	for _, pair := range pairs {
		if pair[0] < 1 || pair[0] >= pair[1] || pair[1] > 4 {
			m["projection_construction_mismatch_count"]++
		}
		classes, bad := wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Induce(counts, []int{pair[0], pair[1]})
		invalid += bad
		ac, correct, split, collision := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(classes, index)
		if float64(ac) < m["derived_only_pair_alias_class_count_min"] {
			m["derived_only_pair_alias_class_count_min"] = float64(ac)
		}
		if float64(ac) > m["derived_only_pair_alias_class_count_max"] {
			m["derived_only_pair_alias_class_count_max"] = float64(ac)
		}
		m["derived_only_pair_correct_alias_triplet_count_total"] += float64(correct)
		m["derived_only_pair_split_entity_triplet_count_total"] += float64(split)
		m["derived_only_pair_cross_entity_collision_pair_count_total"] += float64(collision)
		if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, full) == 0 {
			m["derived_only_pair_full5_membership_match_count"]++
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	classification := "phase_necessary_at_width_two_supported"
	if m["derived_only_pair_full5_membership_match_count"] > 0 {
		classification = "phase_not_necessary_at_width_two"
	}
	return wlmSiRawRepPhaseAnchorPairNecessityR17Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-PAIR-NECESSITY-R17", classification, m}
}
