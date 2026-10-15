package unitary

import (
	"fmt"
	"math"
	"sort"
)

type wlmSiRawRepPhaseAnchorResolutionR18Result struct {
	Schema         string             `json:"schema"`
	Experiment     string             `json:"experiment"`
	Classification string             `json:"classification"`
	Metrics        map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorResolutionR18Induce(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, cardinality, derived int) (wlmLmRawRepContextualAliasInductionR2Classes, int) {
	grouped := map[string][]int{}
	invalid := 0
	if cardinality < 1 || cardinality > 16 || 16%cardinality != 0 || derived < 1 || derived > 4 {
		invalid++
	}
	for id, contexts := range counts {
		projected := map[string]bool{}
		for context := range contexts {
			phase := int(context[0]) - 160
			if phase < 0 || phase >= 16 {
				invalid++
				continue
			}
			bucket := phase * cardinality / 16
			if bucket < 0 || bucket >= cardinality {
				invalid++
				continue
			}
			projected[fmt.Sprintf("%02x%02x", bucket, context[derived])] = true
		}
		keys := make([]string, 0, len(projected))
		for key := range projected {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		signature := ""
		for _, key := range keys {
			signature += key + ";"
		}
		grouped[signature] = append(grouped[signature], id)
	}
	signatures := make([]string, 0, len(grouped))
	for signature := range grouped {
		signatures = append(signatures, signature)
	}
	sort.Strings(signatures)
	classes := wlmLmRawRepContextualAliasInductionR2Classes{surfaceToClass: make([]int, len(counts)), members: make([][]int, 0, len(signatures))}
	for classID, signature := range signatures {
		members := append([]int(nil), grouped[signature]...)
		sort.Ints(members)
		classes.members = append(classes.members, members)
		for _, id := range members {
			classes.surfaceToClass[id] = classID
		}
	}
	return classes, invalid
}

func RunWlmSiRawRepPhaseAnchorResolutionR18() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "actual108_record_count": 0, "actual108_records_per_surface_min": 0, "actual108_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0,
		"r17_exact_pair_control_count": 4, "r17_exact_pair_full5_membership_match_count": 0, "phase_cardinality_count": 5, "derived_coordinate_count": 4, "resolution_arm_count": 20,
		"cardinality1_full5_membership_match_count": 0, "cardinality16_full5_membership_match_count": 0, "all_derived_match_cardinality_count": 0, "first_all_derived_match_cardinality": 0, "post_onset_failure_count": 0, "below16_partial_membership_match_count": 0,
		"resolution_alias_class_count_min": 27, "resolution_alias_class_count_max": 0, "resolution_correct_alias_triplet_count_total": 0, "resolution_split_entity_triplet_count_total": 0, "resolution_cross_entity_collision_pair_count_total": 0, "projection_construction_mismatch_count": 0,
		"evaluation_entity_pair_count_total": 900, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
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
	m["full5_alias_class_count"], m["full5_correct_alias_triplet_count"], m["full5_split_entity_triplet_count"], m["full5_cross_entity_collision_pair_count"] = float64(fc), float64(ft), float64(fs), float64(fx)
	for derived := 1; derived <= 4; derived++ {
		exact, bad := wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Induce(counts, []int{0, derived})
		invalid += bad
		if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(exact, full) == 0 {
			m["r17_exact_pair_full5_membership_match_count"]++
		}
	}
	onset := false
	for _, cardinality := range []int{1, 2, 4, 8, 16} {
		matches := 0
		for derived := 1; derived <= 4; derived++ {
			classes, bad := wlmSiRawRepPhaseAnchorResolutionR18Induce(counts, cardinality, derived)
			invalid += bad
			ac, correct, split, collision := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(classes, index)
			if float64(ac) < m["resolution_alias_class_count_min"] {
				m["resolution_alias_class_count_min"] = float64(ac)
			}
			if float64(ac) > m["resolution_alias_class_count_max"] {
				m["resolution_alias_class_count_max"] = float64(ac)
			}
			m["resolution_correct_alias_triplet_count_total"] += float64(correct)
			m["resolution_split_entity_triplet_count_total"] += float64(split)
			m["resolution_cross_entity_collision_pair_count_total"] += float64(collision)
			if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, full) == 0 {
				matches++
				if cardinality < 16 {
					m["below16_partial_membership_match_count"]++
				}
			}
			if cardinality == 1 {
				control, b := wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Induce(counts, []int{derived})
				invalid += b
				if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, control) != 0 {
					m["projection_construction_mismatch_count"]++
				}
			}
			if cardinality == 16 {
				control, b := wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Induce(counts, []int{0, derived})
				invalid += b
				if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, control) != 0 {
					m["projection_construction_mismatch_count"]++
				}
			}
		}
		if cardinality == 1 {
			m["cardinality1_full5_membership_match_count"] = float64(matches)
		}
		if cardinality == 16 {
			m["cardinality16_full5_membership_match_count"] = float64(matches)
		}
		if matches == 4 {
			m["all_derived_match_cardinality_count"]++
			if !onset {
				m["first_all_derived_match_cardinality"] = float64(cardinality)
				onset = true
			}
		} else if onset {
			m["post_onset_failure_count"]++
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	classification := "exact_phase_resolution_required"
	if m["post_onset_failure_count"] > 0 {
		classification = "nonmonotone"
	} else if m["first_all_derived_match_cardinality"] > 0 && m["first_all_derived_match_cardinality"] < 16 {
		classification = "coarsened_phase_resolution_supported"
	} else if m["first_all_derived_match_cardinality"] == 16 && m["below16_partial_membership_match_count"] > 0 {
		classification = "derived_specific_resolution"
	}
	return wlmSiRawRepPhaseAnchorResolutionR18Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-RESOLUTION-R18", classification, m}
}
