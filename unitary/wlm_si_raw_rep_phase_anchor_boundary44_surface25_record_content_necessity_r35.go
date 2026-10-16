package unitary

import (
	"bytes"
	"math"
)

type wlmSiRawRepPhaseAnchorBoundary44Surface25RecordContentNecessityR35Arm struct {
	SourceSampleIndex         int  `json:"source_sample_index"`
	Authentic                 bool `json:"authentic"`
	Full5AliasClassCount      int  `json:"full5_alias_class_count"`
	Full5CrossEntityCollision int  `json:"full5_cross_entity_collision_pair_count"`
	All4MatchCount            int  `json:"pair_preserving_all4_match_count"`
	PartialMatchCount         int  `json:"pair_preserving_partial_match_count"`
	ZeroMatchCount            int  `json:"pair_preserving_zero_match_count"`
}

type wlmSiRawRepPhaseAnchorBoundary44Surface25RecordContentNecessityR35Result struct {
	Schema         string                                                                  `json:"schema"`
	Experiment     string                                                                  `json:"experiment"`
	Classification string                                                                  `json:"classification"`
	Metrics        map[string]float64                                                      `json:"metrics"`
	Arms           []wlmSiRawRepPhaseAnchorBoundary44Surface25RecordContentNecessityR35Arm `json:"arms"`
}

func RunWlmSiRawRepPhaseAnchorBoundary44Surface25RecordContentNecessityR35() interface{} {
	m := map[string]float64{
		"experimental_arm_count": 44, "completed_arm_count": 0, "base_record_count": 0,
		"independent44_record_count": 0, "prefix_mismatch_count": 0, "authentic_arm_count": 0,
		"duplicate_control_arm_count": 0, "duplicate_identity_mismatch_count": 0,
		"arm_record_cardinality_mismatch_count": 0, "pair_preserving_partition_count_per_arm": 127,
		"partition_panel_arm_count": 0, "authentic_sufficient_arm_count": 0,
		"duplicate_sufficient_arm_count": 0, "duplicate_unchanged_arm_count": 0,
		"duplicate_other_arm_count": 0, "partial_partition_count_total": 0,
		"evaluation_entity_pair_count_total": 0, "task_training_record_count": 0,
		"tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0,
		"invalid_row_count": 0, "counter_overflow_count": 0,
	}
	base := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(43)
	full44 := wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(44)
	m["base_record_count"] = float64(len(base))
	m["independent44_record_count"] = float64(len(full44))
	for surface := 0; surface < 27; surface++ {
		for sample := 0; sample < 43; sample++ {
			if !bytes.Equal(base[surface*43+sample], full44[surface*44+sample]) {
				m["prefix_mismatch_count"]++
			}
		}
	}
	masks := make([]uint16, 0, 127)
	for subset := 1; subset < 128; subset++ {
		var mask uint16
		for q := 1; q < 8; q++ {
			if subset&(1<<(q-1)) != 0 {
				mask |= 1 << q
				mask |= 1 << (q + 8)
			}
		}
		masks = append(masks, mask)
	}
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	authentic := full44[25*44+43]
	for sample := 0; sample < 43; sample++ {
		if bytes.Equal(authentic, base[25*43+sample]) {
			m["duplicate_identity_mismatch_count"]++
		}
	}
	arms := make([]wlmSiRawRepPhaseAnchorBoundary44Surface25RecordContentNecessityR35Arm, 0, 44)
	invalid := 0
	for armIndex := 0; armIndex < 44; armIndex++ {
		sourceSample := 43
		appended := authentic
		isAuthentic := armIndex == 0
		if !isAuthentic {
			sourceSample = armIndex - 1
			appended = base[25*43+sourceSample]
		}
		records := make([][]uint8, 0, 1162)
		for surface := 0; surface < 27; surface++ {
			start := surface * 43
			records = append(records, base[start:start+43]...)
			if surface == 25 {
				records = append(records, appended)
			}
		}
		if len(records) != 1162 {
			m["arm_record_cardinality_mismatch_count"]++
		}
		if isAuthentic {
			m["authentic_arm_count"]++
			if !bytes.Equal(appended, authentic) {
				m["duplicate_identity_mismatch_count"]++
			}
		} else {
			m["duplicate_control_arm_count"]++
			if !bytes.Equal(appended, base[25*43+sourceSample]) {
				m["duplicate_identity_mismatch_count"]++
			}
		}
		counts, bad := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(records, reps)
		invalid += bad
		full := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(counts)
		classCount, _, _, collisions := wlmSiRawRepPhaseAnchorGlobalCoordinateAblationR15Audit(full, index)
		arm := wlmSiRawRepPhaseAnchorBoundary44Surface25RecordContentNecessityR35Arm{SourceSampleIndex: sourceSample, Authentic: isAuthentic, Full5AliasClassCount: classCount, Full5CrossEntityCollision: collisions}
		for _, mask := range masks {
			matches := 0
			for derived := 1; derived <= 4; derived++ {
				classes, rowBad := wlmSiRawRepPhaseAnchorExhaustiveTransferR22Induce(counts, mask, derived)
				invalid += rowBad
				if wlmSiRawRepPhaseAnchorSingleRecordSufficiencyR13MembershipMismatch(classes, full) == 0 {
					matches++
				}
				m["partition_panel_arm_count"]++
				m["evaluation_entity_pair_count_total"] += 36
			}
			if matches == 0 {
				arm.ZeroMatchCount++
			} else if matches == 4 {
				arm.All4MatchCount++
			} else {
				arm.PartialMatchCount++
				m["partial_partition_count_total"]++
			}
		}
		sufficient := arm.ZeroMatchCount == 127 && arm.PartialMatchCount == 0
		unchanged := arm.All4MatchCount == 120 && arm.PartialMatchCount == 0 && arm.ZeroMatchCount == 7
		if isAuthentic {
			if sufficient {
				m["authentic_sufficient_arm_count"]++
			}
		} else if sufficient {
			m["duplicate_sufficient_arm_count"]++
		} else if unchanged {
			m["duplicate_unchanged_arm_count"]++
		} else {
			m["duplicate_other_arm_count"]++
		}
		arms = append(arms, arm)
		m["completed_arm_count"]++
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	classification := "mixed"
	if m["authentic_sufficient_arm_count"] != 1 {
		classification = "authentic_control_failure"
	} else if m["duplicate_unchanged_arm_count"] == 43 {
		classification = "record_content_necessary"
	} else if m["duplicate_sufficient_arm_count"] == 43 {
		classification = "record_count_sufficient"
	} else if m["duplicate_sufficient_arm_count"] > 0 && m["duplicate_sufficient_arm_count"] < 43 && m["duplicate_sufficient_arm_count"]+m["duplicate_unchanged_arm_count"] == 43 && m["partial_partition_count_total"] == 0 {
		classification = "duplicate_identity_dependent"
	}
	return wlmSiRawRepPhaseAnchorBoundary44Surface25RecordContentNecessityR35Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-BOUNDARY44-SURFACE25-RECORD-CONTENT-NECESSITY-R35", classification, m, arms}
}
