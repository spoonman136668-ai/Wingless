package unitary

import (
	"fmt"
	"math"
	"math/bits"
	"sort"
)

type wlmSiRawRepPhaseAnchorAntipodalLawR21Result struct {
	Schema         string             `json:"schema"`
	Experiment     string             `json:"experiment"`
	Classification string             `json:"classification"`
	Metrics        map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorAntipodalLawR21Induce(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, mask uint16, derived int) (wlmLmRawRepContextualAliasInductionR2Classes, int) {
	grouped := map[string][]int{}
	invalid := 0
	if derived < 1 || derived > 4 {
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
			bucket := (mask >> phase) & 1
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

func wlmSiRawRepPhaseAnchorAntipodalLawR21Separates(mask uint16) bool {
	for q := 0; q < 8; q++ {
		if ((mask >> q) & 1) == ((mask >> (q + 8)) & 1) {
			return false
		}
	}
	return true
}

func wlmSiRawRepPhaseAnchorAntipodalLawR21R20Mask(level, offset int) (uint16, bool) {
	var mask uint16
	for phase := 0; phase < 16; phase++ {
		bucket, ok := wlmSiRawRepPhaseAnchorFragmentationSpectrumR20Bucket(phase, level, offset)
		if !ok {
			return 0, false
		}
		if bucket == 1 {
			mask |= 1 << phase
		}
	}
	if mask&1 != 0 {
		mask = ^mask
	}
	return mask, bits.OnesCount16(mask) == 8 && mask&1 == 0
}

func RunWlmSiRawRepPhaseAnchorAntipodalLawR21() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "actual108_record_count": 0, "actual108_records_per_surface_min": 0, "actual108_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0,
		"balanced_partition_count": 0, "complement_normalization_violation_count": 0, "duplicate_partition_count": 0, "antipodal_separating_partition_count": 0, "nonseparating_partition_count": 0,
		"derived_coordinate_count": 4, "partition_arm_count": 0, "antipodal_all4_match_count": 0, "antipodal_partial_match_count": 0, "antipodal_zero_match_count": 0,
		"nonseparating_all4_match_count": 0, "nonseparating_partial_match_count": 0, "nonseparating_zero_match_count": 0, "law_false_negative_partition_count": 0, "law_false_positive_partition_count": 0,
		"r20_control_partition_count": 64, "r20_control_membership_mismatch_count": 0, "partition_alias_class_count_min": 27, "partition_alias_class_count_max": 0,
		"partition_correct_alias_triplet_count_total": 0, "partition_split_entity_triplet_count_total": 0, "partition_cross_entity_collision_pair_count_total": 0,
		"evaluation_entity_pair_count_total": 0, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
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
	matchCounts := make(map[uint16]int, 6435)
	seen := make(map[uint16]bool, 6435)
	for raw := 0; raw <= 0xffff; raw++ {
		mask := uint16(raw)
		if bits.OnesCount16(mask) != 8 || mask&1 != 0 {
			continue
		}
		m["balanced_partition_count"]++
		if seen[mask] {
			m["duplicate_partition_count"]++
		}
		seen[mask] = true
		if mask&1 != 0 {
			m["complement_normalization_violation_count"]++
		}
		separates := wlmSiRawRepPhaseAnchorAntipodalLawR21Separates(mask)
		if separates {
			m["antipodal_separating_partition_count"]++
		} else {
			m["nonseparating_partition_count"]++
		}
		matches := 0
		for derived := 1; derived <= 4; derived++ {
			classes, bad := wlmSiRawRepPhaseAnchorAntipodalLawR21Induce(counts, mask, derived)
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
			}
		}
		matchCounts[mask] = matches
		m["partition_arm_count"] += 4
		m["evaluation_entity_pair_count_total"] += 144
		prefix := "nonseparating"
		if separates {
			prefix = "antipodal"
		}
		switch matches {
		case 0:
			m[prefix+"_zero_match_count"]++
		case 4:
			m[prefix+"_all4_match_count"]++
		default:
			m[prefix+"_partial_match_count"]++
		}
		if separates && matches != 4 {
			m["law_false_negative_partition_count"]++
		}
		if !separates && matches != 0 {
			m["law_false_positive_partition_count"]++
		}
	}
	for level := 1; level <= 8; level++ {
		for offset := 0; offset < 8; offset++ {
			mask, ok := wlmSiRawRepPhaseAnchorAntipodalLawR21R20Mask(level, offset)
			if !ok {
				invalid++
				continue
			}
			expected := 0
			if level%2 == 1 {
				expected = 4
			}
			matches, exists := matchCounts[mask]
			if !exists || matches != expected {
				m["r20_control_membership_mismatch_count"]++
			}
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	classification := "mixed"
	if m["law_false_negative_partition_count"] > 0 {
		classification = "antipodal_separation_not_sufficient"
	} else if m["antipodal_all4_match_count"] == 128 && m["antipodal_partial_match_count"] == 0 && m["antipodal_zero_match_count"] == 0 && m["law_false_positive_partition_count"] > 0 {
		classification = "antipodal_separation_sufficient_not_necessary"
	} else if m["antipodal_all4_match_count"] == 128 && m["nonseparating_zero_match_count"] == 6307 && m["law_false_positive_partition_count"] == 0 {
		classification = "antipodal_separation_iff_sufficient"
	}
	return wlmSiRawRepPhaseAnchorAntipodalLawR21Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-ANTIPODAL-LAW-R21", classification, m}
}
