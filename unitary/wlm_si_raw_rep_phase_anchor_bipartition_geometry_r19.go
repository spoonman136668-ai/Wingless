package unitary

import (
	"fmt"
	"math"
	"sort"
)

type wlmSiRawRepPhaseAnchorBipartitionGeometryR19Result struct {
	Schema         string             `json:"schema"`
	Experiment     string             `json:"experiment"`
	Classification string             `json:"classification"`
	Metrics        map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorBipartitionGeometryR19Bucket(phase, offset int, parity bool) (int, bool) {
	if phase < 0 || phase >= 16 || offset < 0 || offset >= 16 {
		return 0, false
	}
	if parity {
		return phase % 2, true
	}
	if (phase-offset+16)%16 < 8 {
		return 0, true
	}
	return 1, true
}

func wlmSiRawRepPhaseAnchorBipartitionGeometryR19Induce(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, offset, derived int, parity bool) (wlmLmRawRepContextualAliasInductionR2Classes, int) {
	grouped := map[string][]int{}
	invalid := 0
	if derived < 1 || derived > 4 {
		invalid++
	}
	for id, contexts := range counts {
		projected := map[string]bool{}
		for context := range contexts {
			bucket, ok := wlmSiRawRepPhaseAnchorBipartitionGeometryR19Bucket(int(context[0])-160, offset, parity)
			if !ok {
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

func RunWlmSiRawRepPhaseAnchorBipartitionGeometryR19() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "actual108_record_count": 0, "actual108_records_per_surface_min": 0, "actual108_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0,
		"r18_cardinality2_full5_membership_match_count": 0, "contiguous_partition_count": 8, "parity_partition_count": 1, "derived_coordinate_count": 4, "partition_arm_count": 36,
		"balanced_partition_count": 0, "partition_bucket0_phase_count_min": 16, "partition_bucket0_phase_count_max": 0, "contiguous_all_derived_match_count": 0, "parity_full5_membership_match_count": 0,
		"partition_full5_membership_match_count_total": 0, "partition_alias_class_count_min": 27, "partition_alias_class_count_max": 0, "partition_correct_alias_triplet_count_total": 0, "partition_split_entity_triplet_count_total": 0, "partition_cross_entity_collision_pair_count_total": 0,
		"partition_construction_mismatch_count": 0, "evaluation_entity_pair_count_total": 1296, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
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
	for partition := 0; partition < 9; partition++ {
		parity := partition == 8
		offset := partition % 8
		bucket0 := 0
		for phase := 0; phase < 16; phase++ {
			bucket, ok := wlmSiRawRepPhaseAnchorBipartitionGeometryR19Bucket(phase, offset, parity)
			if !ok {
				invalid++
			} else if bucket == 0 {
				bucket0++
			}
			if !parity {
				complement, ok2 := wlmSiRawRepPhaseAnchorBipartitionGeometryR19Bucket(phase, offset+8, false)
				if !ok2 || complement == bucket {
					m["partition_construction_mismatch_count"]++
				}
			}
		}
		if bucket0 == 8 {
			m["balanced_partition_count"]++
		}
		if float64(bucket0) < m["partition_bucket0_phase_count_min"] {
			m["partition_bucket0_phase_count_min"] = float64(bucket0)
		}
		if float64(bucket0) > m["partition_bucket0_phase_count_max"] {
			m["partition_bucket0_phase_count_max"] = float64(bucket0)
		}
		matches := 0
		for derived := 1; derived <= 4; derived++ {
			classes, bad := wlmSiRawRepPhaseAnchorBipartitionGeometryR19Induce(counts, offset, derived, parity)
			invalid += bad
			ac, correct, split, collision := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(classes, index)
			if float64(ac) < m["partition_alias_class_count_min"] {
				m["partition_alias_class_count_min"] = float64(ac)
			}
			if float64(ac) > m["partition_alias_class_count_max"] {
				m["partition_alias_class_count_max"] = float64(ac)
			}
			m["partition_correct_alias_triplet_count_total"] += float64(correct)
			m["partition_split_entity_triplet_count_total"] += float64(split)
			m["partition_cross_entity_collision_pair_count_total"] += float64(collision)
			if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, full) == 0 {
				matches++
				m["partition_full5_membership_match_count_total"]++
			}
			if partition == 0 {
				r18, bad2 := wlmSiRawRepPhaseAnchorResolutionR18Induce(counts, 2, derived)
				invalid += bad2
				if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, r18) == 0 {
					m["r18_cardinality2_full5_membership_match_count"]++
				}
			}
		}
		if parity {
			m["parity_full5_membership_match_count"] = float64(matches)
		} else if matches == 4 {
			m["contiguous_all_derived_match_count"]++
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	classification := "mixed"
	if m["contiguous_all_derived_match_count"] == 8 && m["parity_full5_membership_match_count"] == 4 {
		classification = "binary_partition_geometry_invariant"
	} else if m["contiguous_all_derived_match_count"] == 8 {
		classification = "contiguous_rotation_invariant_parity_specific"
	} else if m["parity_full5_membership_match_count"] > 0 || m["partition_full5_membership_match_count_total"] > 4 {
		classification = "derived_specific_partition_geometry"
	} else if m["r18_cardinality2_full5_membership_match_count"] == 4 {
		classification = "boundary_specific_binary_sufficiency"
	}
	return wlmSiRawRepPhaseAnchorBipartitionGeometryR19Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-BIPARTITION-GEOMETRY-R19", classification, m}
}
