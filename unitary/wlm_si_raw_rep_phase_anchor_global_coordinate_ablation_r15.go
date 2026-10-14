package unitary

import (
	"fmt"
	"math"
	"sort"
)

type wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Induce(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, omitted int) (wlmLmRawRepContextualAliasInductionR2Classes, int) {
	grouped := make(map[string][]int)
	invalid := 0
	for id, contexts := range counts {
		projected := make(map[[4]uint8]bool)
		for context := range contexts {
			var key [4]uint8
			at := 0
			for coordinate := 0; coordinate < 5; coordinate++ {
				if coordinate == omitted {
					continue
				}
				if at >= len(key) {
					invalid++
					continue
				}
				key[at] = context[coordinate]
				at++
			}
			if at != 4 {
				invalid++
			}
			projected[key] = true
		}
		keys := make([][4]uint8, 0, len(projected))
		for key := range projected {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			for coordinate := 0; coordinate < 4; coordinate++ {
				if keys[i][coordinate] != keys[j][coordinate] {
					return keys[i][coordinate] < keys[j][coordinate]
				}
			}
			return false
		})
		signature := ""
		for _, key := range keys {
			signature += fmt.Sprintf("%02x%02x%02x%02x;", key[0], key[1], key[2], key[3])
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

func wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(classes wlmLmRawRepContextualAliasInductionR2Classes, index map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int) (int, int, int, int) {
	aliasClasses, _, correct := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(classes, index)
	return aliasClasses, correct, wlmSiRawRepPhaseAnchorSamplingR6SplitTriplets(classes, index), wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(classes, index)
}

func RunWlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"actual108_record_count": 0, "actual108_records_per_surface_min": 0, "actual108_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0,
		"projection_arm_count": 5, "projection_construction_mismatch_count": 0,
		"leave_phase_alias_class_count": 0, "leave_phase_correct_alias_triplet_count": 0, "leave_phase_split_entity_triplet_count": 0, "leave_phase_cross_entity_collision_pair_count": 0,
		"derived_omission_arm_count": 4, "derived_omission_alias_class_count_total": 0, "derived_omission_correct_alias_triplet_count_total": 0,
		"derived_omission_split_entity_triplet_count_total": 0, "derived_omission_cross_entity_collision_pair_count_total": 0, "derived_omission_full5_membership_match_count": 0,
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
	for omitted := 0; omitted < 5; omitted++ {
		classes, bad := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Induce(counts, omitted)
		invalid += bad
		ac, correct, split, collision := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(classes, index)
		if omitted == 0 {
			m["leave_phase_alias_class_count"], m["leave_phase_correct_alias_triplet_count"] = float64(ac), float64(correct)
			m["leave_phase_split_entity_triplet_count"], m["leave_phase_cross_entity_collision_pair_count"] = float64(split), float64(collision)
		} else {
			m["derived_omission_alias_class_count_total"] += float64(ac)
			m["derived_omission_correct_alias_triplet_count_total"] += float64(correct)
			m["derived_omission_split_entity_triplet_count_total"] += float64(split)
			m["derived_omission_cross_entity_collision_pair_count_total"] += float64(collision)
			if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, full) == 0 {
				m["derived_omission_full5_membership_match_count"]++
			}
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-GLOBAL-COORDINATE-ABLATION-R15", m}
}
