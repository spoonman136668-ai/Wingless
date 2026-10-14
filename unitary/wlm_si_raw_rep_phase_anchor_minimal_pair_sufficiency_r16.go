package unitary

import (
	"fmt"
	"math"
	"sort"
)

type wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Result struct {
	Schema         string             `json:"schema"`
	Experiment     string             `json:"experiment"`
	Classification string             `json:"classification"`
	Metrics        map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Induce(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, coordinates []int) (wlmLmRawRepContextualAliasInductionR2Classes, int) {
	grouped := make(map[string][]int)
	invalid := 0
	for id, contexts := range counts {
		projected := make(map[string]bool)
		for context := range contexts {
			key := ""
			for _, coordinate := range coordinates {
				if coordinate < 0 || coordinate >= len(context) {
					invalid++
					continue
				}
				key += fmt.Sprintf("%02x", context[coordinate])
			}
			if len(key) != 2*len(coordinates) {
				invalid++
			}
			projected[key] = true
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

func RunWlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"actual108_record_count": 0, "actual108_records_per_surface_min": 0, "actual108_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0,
		"phase_only_distinct_support_value_count": 0, "phase_only_alias_class_count": 0, "phase_only_correct_alias_triplet_count": 0,
		"phase_only_split_entity_triplet_count": 0, "phase_only_cross_entity_collision_pair_count": 0,
		"phase_plus_derived_arm_count": 4, "phase_plus_derived_alias_class_count_total": 0, "phase_plus_derived_correct_alias_triplet_count_total": 0,
		"phase_plus_derived_split_entity_triplet_count_total": 0, "phase_plus_derived_cross_entity_collision_pair_count_total": 0,
		"phase_plus_derived_full5_membership_match_count": 0, "projection_construction_mismatch_count": 0,
		"evaluation_entity_pair_count_total": 216, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0,
		"capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
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
	phaseValues := make(map[uint8]bool)
	for _, contexts := range counts {
		for context := range contexts {
			phaseValues[context[0]] = true
		}
	}
	m["phase_only_distinct_support_value_count"] = float64(len(phaseValues))
	phaseOnly, bad := wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Induce(counts, []int{0})
	invalid += bad
	pc, pt, ps, px := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(phaseOnly, index)
	m["phase_only_alias_class_count"], m["phase_only_correct_alias_triplet_count"] = float64(pc), float64(pt)
	m["phase_only_split_entity_triplet_count"], m["phase_only_cross_entity_collision_pair_count"] = float64(ps), float64(px)
	for derived := 1; derived < 5; derived++ {
		classes, armBad := wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Induce(counts, []int{0, derived})
		invalid += armBad
		ac, correct, split, collision := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(classes, index)
		m["phase_plus_derived_alias_class_count_total"] += float64(ac)
		m["phase_plus_derived_correct_alias_triplet_count_total"] += float64(correct)
		m["phase_plus_derived_split_entity_triplet_count_total"] += float64(split)
		m["phase_plus_derived_cross_entity_collision_pair_count_total"] += float64(collision)
		if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, full) == 0 {
			m["phase_plus_derived_full5_membership_match_count"]++
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	phaseFull := wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(phaseOnly, full) == 0
	pairs := int(m["phase_plus_derived_full5_membership_match_count"])
	classification := "mixed"
	if phaseFull {
		classification = "phase_alone_sufficient"
	} else if pc == 1 && px == 36 && pairs == 4 {
		classification = "phase_plus_any_derived_minimal_sufficiency_supported"
	} else if pairs > 0 && pairs < 4 {
		classification = "derived_coordinate_specific_sufficiency"
	} else if pairs == 0 {
		classification = "no_tested_pair_sufficient"
	}
	return wlmSiRawRepPhaseAnchorMinimalPairSufficiencyR16Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-MINIMAL-PAIR-SUFFICIENCY-R16", classification, m}
}
