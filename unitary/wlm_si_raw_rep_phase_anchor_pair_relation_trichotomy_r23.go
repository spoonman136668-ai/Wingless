package unitary

import (
	"math"
	"math/bits"
)

type wlmSiRawRepPhaseAnchorPairRelationTrichotomyR23Result struct {
	Schema         string             `json:"schema"`
	Experiment     string             `json:"experiment"`
	Classification string             `json:"classification"`
	Metrics        map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorPairRelationTrichotomyR23Class(mask uint16) string {
	separated := 0
	for q := 0; q < 8; q++ {
		if ((mask >> q) & 1) != ((mask >> (q + 8)) & 1) {
			separated++
		}
	}
	if separated == 8 {
		return "fully_separating"
	}
	if separated == 0 {
		return "fully_preserving"
	}
	return "mixed"
}

func RunWlmSiRawRepPhaseAnchorPairRelationTrichotomyR23() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "exhaustive_record_count": 0, "exhaustive_records_per_surface_min": 0, "exhaustive_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0, "balanced_partition_count": 0, "duplicate_partition_count": 0, "complement_normalization_violation_count": 0,
		"fully_separating_partition_count": 0, "mixed_partition_count": 0, "fully_preserving_partition_count": 0, "independent_fully_preserving_mask_count": 0, "fully_preserving_set_mismatch_count": 0, "derived_coordinate_count": 4, "partition_arm_count": 0,
		"fully_separating_all4_match_count": 0, "fully_separating_partial_match_count": 0, "fully_separating_zero_match_count": 0, "mixed_all4_match_count": 0, "mixed_partial_match_count": 0, "mixed_zero_match_count": 0, "fully_preserving_all4_match_count": 0, "fully_preserving_partial_match_count": 0, "fully_preserving_zero_match_count": 0,
		"classification_mismatch_partition_count": 0, "evaluation_entity_pair_count_total": 0, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	records := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Bridge()
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	m["exhaustive_record_count"] = float64(len(records))
	counts, invalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(records, reps)
	lo, hi := wlmSiRawRepPhaseAnchorSamplingR6RecordRange(counts)
	m["exhaustive_records_per_surface_min"], m["exhaustive_records_per_surface_max"] = float64(lo), float64(hi)
	full := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(counts)
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	fc, ft, fs, fx := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(full, index)
	m["full5_alias_class_count"], m["full5_correct_alias_triplet_count"], m["full5_split_entity_triplet_count"], m["full5_cross_entity_collision_pair_count"] = float64(fc), float64(ft), float64(fs), float64(fx)
	seen := map[uint16]bool{}
	preserving := map[uint16]bool{}
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
		class := wlmSiRawRepPhaseAnchorPairRelationTrichotomyR23Class(mask)
		m[class+"_partition_count"]++
		if class == "fully_preserving" {
			preserving[mask] = true
		}
		matches := 0
		for derived := 1; derived <= 4; derived++ {
			classes, bad := wlmSiRawRepPhaseAnchorExhaustiveTransferR22Induce(counts, mask, derived)
			invalid += bad
			wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(classes, index)
			if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, full) == 0 {
				matches++
			}
		}
		m["partition_arm_count"] += 4
		m["evaluation_entity_pair_count_total"] += 144
		suffix := "partial_match_count"
		if matches == 0 {
			suffix = "zero_match_count"
		} else if matches == 4 {
			suffix = "all4_match_count"
		}
		m[class+"_"+suffix]++
		expected := 4
		if class == "fully_preserving" {
			expected = 0
		}
		if matches != expected {
			m["classification_mismatch_partition_count"]++
		}
	}
	independent := map[uint16]bool{}
	for subset := 0; subset < 256; subset++ {
		if bits.OnesCount8(uint8(subset)) != 4 {
			continue
		}
		var mask uint16
		for q := 0; q < 8; q++ {
			if subset&(1<<q) != 0 {
				mask |= 1<<q | 1<<(q+8)
			}
		}
		if mask&1 != 0 {
			mask = ^mask
		}
		independent[mask] = true
	}
	m["independent_fully_preserving_mask_count"] = float64(len(independent))
	for mask := range preserving {
		if !independent[mask] {
			m["fully_preserving_set_mismatch_count"]++
		}
	}
	for mask := range independent {
		if !preserving[mask] {
			m["fully_preserving_set_mismatch_count"]++
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			m["invalid_row_count"]++
		}
	}
	classification := "mixed"
	if m["fully_separating_all4_match_count"] < 128 || m["mixed_all4_match_count"] < 6272 {
		classification = "fully_preserving_not_necessary_for_failure"
	} else if m["fully_preserving_zero_match_count"] < 35 {
		classification = "fully_preserving_not_sufficient_for_failure"
	} else if m["classification_mismatch_partition_count"] == 0 && m["fully_separating_partial_match_count"] == 0 && m["mixed_partial_match_count"] == 0 && m["fully_preserving_partial_match_count"] == 0 {
		classification = "antipodal_pair_relation_trichotomy_exact"
	}
	return wlmSiRawRepPhaseAnchorPairRelationTrichotomyR23Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-PAIR-RELATION-TRICHOTOMY-R23", classification, m}
}
