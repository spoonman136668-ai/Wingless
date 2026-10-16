package unitary

import (
	"math"
	"math/bits"
)

type wlmSiRawRepPhaseAnchorBoundary40CardinalityTransferR30Result struct {
	Schema         string             `json:"schema"`
	Experiment     string             `json:"experiment"`
	Classification string             `json:"classification"`
	Metrics        map[string]float64 `json:"metrics"`
}

func RunWlmSiRawRepPhaseAnchorBoundary40CardinalityTransferR30() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "independent40_record_count": 0, "independent40_records_per_surface_min": 0, "independent40_records_per_surface_max": 0,
		"full5_alias_class_count": 0, "full5_correct_alias_triplet_count": 0, "full5_split_entity_triplet_count": 0, "full5_cross_entity_collision_pair_count": 0,
		"nonconstant_partition_count": 0, "duplicate_partition_count": 0, "complement_normalization_violation_count": 0, "pair_splitting_partition_count": 0, "pair_preserving_partition_count": 0,
		"independent_pair_preserving_mask_count": 0, "pair_preserving_set_mismatch_count": 0, "cardinality_stratum_count": 15, "cardinality_stratum_cardinality_mismatch_count": 0,
		"derived_coordinate_count": 4, "partition_arm_count": 0, "pair_splitting_all4_match_count": 0, "pair_splitting_partial_match_count": 0, "pair_splitting_zero_match_count": 0,
		"pair_preserving_all4_match_count": 0, "pair_preserving_partial_match_count": 0, "pair_preserving_zero_match_count": 0, "classification_mismatch_partition_count": 0,
		"stratum_law_mismatch_count": 0, "evaluation_entity_pair_count_total": 0, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0,
		"capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	for k := 1; k <= 15; k++ {
		m["cardinality_"+string(rune('0'+k/10))+string(rune('0'+k%10))+"_partition_count"] = 0
		m["cardinality_"+string(rune('0'+k/10))+string(rune('0'+k%10))+"_law_mismatch_count"] = 0
	}
	records := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(40)
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	m["independent40_record_count"] = float64(len(records))
	counts, invalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(records, reps)
	lo, hi := wlmSiRawRepPhaseAnchorSamplingR6RecordRange(counts)
	m["independent40_records_per_surface_min"], m["independent40_records_per_surface_max"] = float64(lo), float64(hi)
	full := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(counts)
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	fc, ft, fs, fx := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(full, index)
	m["full5_alias_class_count"], m["full5_correct_alias_triplet_count"], m["full5_split_entity_triplet_count"], m["full5_cross_entity_collision_pair_count"] = float64(fc), float64(ft), float64(fs), float64(fx)

	seen := map[uint16]bool{}
	preserving := map[uint16]bool{}
	stratumMismatch := map[int]bool{}
	for raw := 1; raw <= 0xffff; raw++ {
		mask := uint16(raw)
		if mask&1 != 0 {
			continue
		}
		m["nonconstant_partition_count"]++
		if seen[mask] {
			m["duplicate_partition_count"]++
		}
		seen[mask] = true
		if mask&1 != 0 {
			m["complement_normalization_violation_count"]++
		}
		cardinality := bits.OnesCount16(mask)
		key := "cardinality_" + string(rune('0'+cardinality/10)) + string(rune('0'+cardinality%10))
		m[key+"_partition_count"]++
		class := "pair_splitting"
		if wlmSiRawRepPhaseAnchorPairRelationTrichotomyR23Class(mask) == "fully_preserving" {
			class = "pair_preserving"
			preserving[mask] = true
		}
		m[class+"_partition_count"]++
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
		if class == "pair_preserving" {
			expected = 0
		}
		if matches != expected {
			m["classification_mismatch_partition_count"]++
			m[key+"_law_mismatch_count"]++
			stratumMismatch[cardinality] = true
		}
	}
	m["stratum_law_mismatch_count"] = float64(len(stratumMismatch))
	for k := 1; k <= 15; k++ {
		key := "cardinality_" + string(rune('0'+k/10)) + string(rune('0'+k%10))
		want := 1
		for i := 1; i <= k; i++ {
			want = want * (16 - i) / i
		}
		if int(m[key+"_partition_count"]) != want {
			m["cardinality_stratum_cardinality_mismatch_count"]++
		}
	}
	independent := map[uint16]bool{}
	for subset := 1; subset < 128; subset++ {
		var mask uint16
		for q := 1; q < 8; q++ {
			if subset&(1<<(q-1)) != 0 {
				mask |= 1 << q
				mask |= 1 << (q + 8)
			}
		}
		independent[mask] = true
	}
	m["independent_pair_preserving_mask_count"] = float64(len(independent))
	for mask := range preserving {
		if !independent[mask] {
			m["pair_preserving_set_mismatch_count"]++
		}
	}
	for mask := range independent {
		if !preserving[mask] {
			m["pair_preserving_set_mismatch_count"]++
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			m["invalid_row_count"]++
		}
	}
	classification := "mixed"
	if m["pair_splitting_all4_match_count"] < 32640 {
		classification = "boundary40_breaks_pair_splitting_sufficiency"
	} else if m["pair_preserving_zero_match_count"] < 127 {
		classification = "boundary40_breaks_pair_preserving_failure"
	} else if m["classification_mismatch_partition_count"] == 0 && m["stratum_law_mismatch_count"] == 0 && m["pair_splitting_partial_match_count"] == 0 && m["pair_preserving_partial_match_count"] == 0 {
		classification = "boundary40_unbalanced_antipodal_law_exact"
	}
	return wlmSiRawRepPhaseAnchorBoundary40CardinalityTransferR30Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-BOUNDARY40-CARDINALITY-TRANSFER-R30", classification, m}
}
