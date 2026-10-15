package unitary

import (
	"fmt"
	"math"
	"sort"
)

type wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Result struct {
	Schema         string             `json:"schema"`
	Experiment     string             `json:"experiment"`
	Classification string             `json:"classification"`
	Metrics        map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Bucket(phase, level, offset int) (int, bool) {
	if phase < 0 || phase >= 16 || level < 1 || level > 8 || offset < 0 || offset >= 8 {
		return 0, false
	}
	rotated := (phase - offset + 16) % 16
	return (rotated * level / 8) % 2, true
}

func wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Induce(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, level, offset, derived int) (wlmLmRawRepContextualAliasInductionR2Classes, int) {
	grouped := map[string][]int{}
	invalid := 0
	if derived < 1 || derived > 4 {
		invalid++
	}
	for id, contexts := range counts {
		projected := map[string]bool{}
		for context := range contexts {
			bucket, ok := wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Bucket(int(context[0])-160, level, offset)
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

func RunWlmSiRawRepPhaseAnchorFragmentationSpectrumR20() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "actual108_record_count": 0, "actual108_records_per_surface_min": 0, "actual108_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0,
		"fragmentation_level_count": 8, "rotations_per_level": 8, "derived_coordinate_count": 4, "partition_arm_count": 256, "balanced_partition_count": 0,
		"partition_bucket0_phase_count_min": 16, "partition_bucket0_phase_count_max": 0, "transition_count_construction_mismatch_count": 0, "rotation_construction_mismatch_count": 0,
		"partition_full5_membership_match_count_total": 0, "fully_matching_level_count": 0, "zero_matching_level_count": 0, "partially_matching_level_count": 0, "restoration_after_zero_level_count": 0,
		"partition_alias_class_count_min": 27, "partition_alias_class_count_max": 0, "partition_correct_alias_triplet_count_total": 0, "partition_split_entity_triplet_count_total": 0, "partition_cross_entity_collision_pair_count_total": 0,
		"evaluation_entity_pair_count_total": 9216, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	for level := 1; level <= 8; level++ {
		m[fmt.Sprintf("level%d_full5_membership_match_count", level)] = 0
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
	seenZero := false
	for level := 1; level <= 8; level++ {
		levelMatches := 0
		for offset := 0; offset < 8; offset++ {
			bucket0, transitions := 0, 0
			for phase := 0; phase < 16; phase++ {
				bucket, ok := wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Bucket(phase, level, offset)
				next, ok2 := wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Bucket((phase+1)%16, level, offset)
				if !ok || !ok2 {
					invalid++
					continue
				}
				if bucket == 0 {
					bucket0++
				}
				if bucket != next {
					transitions++
				}
				reconstructed, ok3 := wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Bucket((phase-offset+16)%16, level, 0)
				if !ok3 || reconstructed != bucket {
					m["rotation_construction_mismatch_count"]++
				}
			}
			if bucket0 == 8 {
				m["balanced_partition_count"]++
			}
			if bucket0 < int(m["partition_bucket0_phase_count_min"]) {
				m["partition_bucket0_phase_count_min"] = float64(bucket0)
			}
			if bucket0 > int(m["partition_bucket0_phase_count_max"]) {
				m["partition_bucket0_phase_count_max"] = float64(bucket0)
			}
			if transitions != 2*level {
				m["transition_count_construction_mismatch_count"]++
			}
			for derived := 1; derived <= 4; derived++ {
				classes, bad := wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Induce(counts, level, offset, derived)
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
					levelMatches++
					m["partition_full5_membership_match_count_total"]++
				}
			}
		}
		m[fmt.Sprintf("level%d_full5_membership_match_count", level)] = float64(levelMatches)
		if levelMatches == 32 {
			m["fully_matching_level_count"]++
		} else if levelMatches == 0 {
			m["zero_matching_level_count"]++
			seenZero = true
		} else {
			m["partially_matching_level_count"]++
		}
		if seenZero && levelMatches > 0 {
			m["restoration_after_zero_level_count"]++
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	classification := "mixed"
	if m["partially_matching_level_count"] > 0 {
		classification = "rotation_or_coordinate_sensitive_fragmentation"
	} else if m["fully_matching_level_count"] == 1 && m["level1_full5_membership_match_count"] == 32 && m["zero_matching_level_count"] == 7 {
		classification = "immediate_fragmentation_failure"
	} else if m["restoration_after_zero_level_count"] > 0 {
		classification = "nonmonotone_fragmentation"
	} else if m["fully_matching_level_count"] >= 2 && m["zero_matching_level_count"] >= 2 {
		classification = "monotone_fragmentation_boundary"
	}
	return wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-FRAGMENTATION-SPECTRUM-R20", classification, m}
}
