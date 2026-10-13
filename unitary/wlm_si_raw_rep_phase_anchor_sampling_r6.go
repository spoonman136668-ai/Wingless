package unitary

import "math"

type wlmSiRawRepPhaseAnchorSamplingR6Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorSamplingR6Bridge(independent bool) [][]uint8 {
	records := make([][]uint8, 0, 1728)
	for entity := 0; entity < 9; entity++ {
		for family := 0; family < 3; family++ {
			state := uint32(0x243f6a88) ^ uint32((entity+1)*1009)
			if independent {
				state ^= uint32(family+1) * uint32(0x9e3779b9)
			}
			for sample := 0; sample < 64; sample++ {
				state = state*1664525 + 1013904223
				contextIndex := int((state >> 16) % 16)
				noise := uint32(0x85ebca6b) ^ uint32((entity+1)*1009+(family+1)*10007+(sample+1)*65537)
				record := append([]uint8{}, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&noise, 3)...)
				record = append(record, uint8(160+contextIndex))
				record = append(record, uint8(192+(entity*7+contextIndex*3)%16), uint8(208+(entity*5+contextIndex*5)%16))
				motif := wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity)
				record = append(record, motif[:]...)
				record = append(record, uint8(224+(entity*11+contextIndex*7)%16), uint8(240+(entity*13+contextIndex*11)%16))
				record = append(record, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&noise, 3)...)
				records = append(records, record)
			}
		}
	}
	return records
}

func wlmSiRawRepPhaseAnchorSamplingR6RecordRange(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32) (int, int) {
	minimum, maximum := int(^uint(0)>>1), 0
	for _, support := range counts {
		total := 0
		for _, count := range support {
			total += int(count)
		}
		if total < minimum {
			minimum = total
		}
		if total > maximum {
			maximum = total
		}
	}
	return minimum, maximum
}

func wlmSiRawRepPhaseAnchorSamplingR6ContextRange(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32) (int, int) {
	minimum, maximum := int(^uint(0)>>1), 0
	for _, support := range counts {
		if len(support) < minimum {
			minimum = len(support)
		}
		if len(support) > maximum {
			maximum = len(support)
		}
	}
	return minimum, maximum
}

func wlmSiRawRepPhaseAnchorSamplingR6MaxSymmetricDifference(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, index map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int) int {
	maximum := 0
	for entity := 0; entity < 9; entity++ {
		for leftFamily := 0; leftFamily < 3; leftFamily++ {
			for rightFamily := leftFamily + 1; rightFamily < 3; rightFamily++ {
				left := counts[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(leftFamily, entity)]]
				right := counts[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(rightFamily, entity)]]
				difference := 0
				for context := range left {
					if _, ok := right[context]; !ok {
						difference++
					}
				}
				for context := range right {
					if _, ok := left[context]; !ok {
						difference++
					}
				}
				if difference > maximum {
					maximum = difference
				}
			}
		}
	}
	return maximum
}

func wlmSiRawRepPhaseAnchorSamplingR6SplitTriplets(classes wlmLmRawRepContextualAliasInductionR2Classes, index map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int) int {
	split := 0
	for entity := 0; entity < 9; entity++ {
		class0 := classes.surfaceToClass[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(0, entity)]]
		class1 := classes.surfaceToClass[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(1, entity)]]
		class2 := classes.surfaceToClass[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(2, entity)]]
		if class0 != class1 || class0 != class2 {
			split++
		}
	}
	return split
}

func RunWlmSiRawRepPhaseAnchorSamplingR6() interface{} {
	m := map[string]float64{
		"exhaustive_bridge_record_count": 0, "matched_bridge_record_count": 0, "independent_bridge_record_count": 0,
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"exhaustive_distinct_contexts_per_surface_min": 0, "exhaustive_distinct_contexts_per_surface_max": 0,
		"exhaustive_alias_class_count": 0, "exhaustive_three_member_class_count": 0, "exhaustive_correct_alias_triplet_count": 0,
		"exhaustive_split_entity_triplet_count": 0, "exhaustive_cross_entity_collision_pair_count": 0,
		"matched_records_per_surface_min": 0, "matched_records_per_surface_max": 0, "matched_max_same_entity_support_symmetric_difference": 0,
		"matched_alias_class_count": 0, "matched_three_member_class_count": 0, "matched_correct_alias_triplet_count": 0,
		"matched_split_entity_triplet_count": 0, "matched_cross_entity_collision_pair_count": 0,
		"independent_records_per_surface_min": 0, "independent_records_per_surface_max": 0, "independent_max_same_entity_support_symmetric_difference": 0,
		"independent_alias_class_count": 0, "independent_three_member_class_count": 0, "independent_correct_alias_triplet_count": 0,
		"independent_split_entity_triplet_count": 0, "independent_cross_entity_collision_pair_count": 0,
		"evaluation_entity_pair_count_per_arm": 36, "task_training_record_count": 0, "tokenizer_use_count": 0,
		"external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	exhaustive := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Bridge()
	matched := wlmSiRawRepPhaseAnchorSamplingR6Bridge(false)
	independent := wlmSiRawRepPhaseAnchorSamplingR6Bridge(true)
	m["exhaustive_bridge_record_count"] = float64(len(exhaustive))
	m["matched_bridge_record_count"] = float64(len(matched))
	m["independent_bridge_record_count"] = float64(len(independent))
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	exhaustiveCounts, exhaustiveInvalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(exhaustive, reps)
	matchedCounts, matchedInvalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(matched, reps)
	independentCounts, independentInvalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(independent, reps)
	m["invalid_row_count"] = float64(exhaustiveInvalid + matchedInvalid + independentInvalid)
	exhaustiveClasses := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(exhaustiveCounts)
	matchedClasses := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(matchedCounts)
	independentClasses := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(independentCounts)
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	ea, et, ec := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(exhaustiveClasses, index)
	ma, mt, mc := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(matchedClasses, index)
	ia, it, ic := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(independentClasses, index)
	m["exhaustive_alias_class_count"], m["exhaustive_three_member_class_count"], m["exhaustive_correct_alias_triplet_count"] = float64(ea), float64(et), float64(ec)
	m["matched_alias_class_count"], m["matched_three_member_class_count"], m["matched_correct_alias_triplet_count"] = float64(ma), float64(mt), float64(mc)
	m["independent_alias_class_count"], m["independent_three_member_class_count"], m["independent_correct_alias_triplet_count"] = float64(ia), float64(it), float64(ic)
	m["exhaustive_split_entity_triplet_count"] = float64(wlmSiRawRepPhaseAnchorSamplingR6SplitTriplets(exhaustiveClasses, index))
	m["matched_split_entity_triplet_count"] = float64(wlmSiRawRepPhaseAnchorSamplingR6SplitTriplets(matchedClasses, index))
	m["independent_split_entity_triplet_count"] = float64(wlmSiRawRepPhaseAnchorSamplingR6SplitTriplets(independentClasses, index))
	m["exhaustive_cross_entity_collision_pair_count"] = float64(wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(exhaustiveClasses, index))
	m["matched_cross_entity_collision_pair_count"] = float64(wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(matchedClasses, index))
	m["independent_cross_entity_collision_pair_count"] = float64(wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(independentClasses, index))
	exhaustiveMin, exhaustiveMax := wlmSiRawRepPhaseAnchorSamplingR6ContextRange(exhaustiveCounts)
	m["exhaustive_distinct_contexts_per_surface_min"], m["exhaustive_distinct_contexts_per_surface_max"] = float64(exhaustiveMin), float64(exhaustiveMax)
	matchedMin, matchedMax := wlmSiRawRepPhaseAnchorSamplingR6RecordRange(matchedCounts)
	independentMin, independentMax := wlmSiRawRepPhaseAnchorSamplingR6RecordRange(independentCounts)
	m["matched_records_per_surface_min"], m["matched_records_per_surface_max"] = float64(matchedMin), float64(matchedMax)
	m["independent_records_per_surface_min"], m["independent_records_per_surface_max"] = float64(independentMin), float64(independentMax)
	m["matched_max_same_entity_support_symmetric_difference"] = float64(wlmSiRawRepPhaseAnchorSamplingR6MaxSymmetricDifference(matchedCounts, index))
	m["independent_max_same_entity_support_symmetric_difference"] = float64(wlmSiRawRepPhaseAnchorSamplingR6MaxSymmetricDifference(independentCounts, index))
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepPhaseAnchorSamplingR6Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-PHASE-ANCHOR-SAMPLING-R6", m}
}
