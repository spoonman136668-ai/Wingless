package unitary

import "math"

type wlmSiRawRepContextSupportIdentifiabilityR4Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepContextSupportIdentifiabilityR4Bridge() [][]uint8 {
	records := make([][]uint8, 0, 432)
	for entity := 0; entity < 9; entity++ {
		for family := 0; family < 3; family++ {
			for contextIndex := 0; contextIndex < 16; contextIndex++ {
				state := uint32(0x85ebca6b) ^ uint32((entity+1)*1009+(family+1)*10007+(contextIndex+1)*65537)
				record := append([]uint8{}, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
				record = append(record, uint8(192+(entity*7+contextIndex*3)%16), uint8(208+(entity*5+contextIndex*5)%16))
				motif := wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity)
				record = append(record, motif[:]...)
				record = append(record, uint8(224+(entity*11+contextIndex*7)%16), uint8(240+(entity*13+contextIndex*11)%16))
				record = append(record, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
				records = append(records, record)
			}
		}
	}
	return records
}

func wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(classes wlmLmRawRepContextualAliasInductionR2Classes, index map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int) int {
	pairs := 0
	for leftEntity := 0; leftEntity < 9; leftEntity++ {
		for rightEntity := leftEntity + 1; rightEntity < 9; rightEntity++ {
			collides := false
			for leftFamily := 0; leftFamily < 3 && !collides; leftFamily++ {
				for rightFamily := 0; rightFamily < 3; rightFamily++ {
					left := index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(leftFamily, leftEntity)]
					right := index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(rightFamily, rightEntity)]
					if classes.surfaceToClass[left] == classes.surfaceToClass[right] {
						collides = true
						break
					}
				}
			}
			if collides {
				pairs++
			}
		}
	}
	return pairs
}

func RunWlmSiRawRepContextSupportIdentifiabilityR4() interface{} {
	m := map[string]float64{
		"exhaustive_bridge_record_count": 0, "exhaustive_records_per_surface_min": 0, "exhaustive_records_per_surface_max": 0,
		"exhaustive_distinct_contexts_per_surface_min": 0, "exhaustive_distinct_contexts_per_surface_max": 0,
		"matched_control_bridge_record_count": 0, "selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"matched_control_alias_class_count": 0, "matched_control_three_member_class_count": 0, "matched_control_correct_alias_triplet_count": 0, "matched_control_cross_entity_collision_pair_count": 0,
		"exhaustive_alias_class_count": 0, "exhaustive_three_member_class_count": 0, "exhaustive_correct_alias_triplet_count": 0, "exhaustive_cross_entity_collision_pair_count": 0,
		"evaluation_entity_pair_count": 36, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0,
		"capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	matched := wlmSiRawRepContextSamplingFalsificationR2Bridge(false)
	exhaustive := wlmSiRawRepContextSupportIdentifiabilityR4Bridge()
	m["matched_control_bridge_record_count"] = float64(len(matched))
	m["exhaustive_bridge_record_count"] = float64(len(exhaustive))
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(matched)
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	matchedClasses := wlmSiRawRepContextSupportFalsificationR3Induce(matched, reps)
	exhaustiveClasses := wlmSiRawRepContextSupportFalsificationR3Induce(exhaustive, reps)
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	matchedAliases, matchedThree, matchedCorrect := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(matchedClasses, index)
	exhaustiveAliases, exhaustiveThree, exhaustiveCorrect := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(exhaustiveClasses, index)
	m["matched_control_alias_class_count"], m["matched_control_three_member_class_count"], m["matched_control_correct_alias_triplet_count"] = float64(matchedAliases), float64(matchedThree), float64(matchedCorrect)
	m["exhaustive_alias_class_count"], m["exhaustive_three_member_class_count"], m["exhaustive_correct_alias_triplet_count"] = float64(exhaustiveAliases), float64(exhaustiveThree), float64(exhaustiveCorrect)
	m["matched_control_cross_entity_collision_pair_count"] = float64(wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(matchedClasses, index))
	m["exhaustive_cross_entity_collision_pair_count"] = float64(wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(exhaustiveClasses, index))
	counts := wlmSiRawRepContextSamplingFalsificationR2Counts(exhaustive, reps)
	minRecords, maxRecords, minContexts, maxContexts := int(^uint(0)>>1), 0, int(^uint(0)>>1), 0
	for _, contexts := range counts {
		records := 0
		for _, count := range contexts {
			records += int(count)
		}
		if records < minRecords {
			minRecords = records
		}
		if records > maxRecords {
			maxRecords = records
		}
		if len(contexts) < minContexts {
			minContexts = len(contexts)
		}
		if len(contexts) > maxContexts {
			maxContexts = len(contexts)
		}
	}
	m["exhaustive_records_per_surface_min"], m["exhaustive_records_per_surface_max"] = float64(minRecords), float64(maxRecords)
	m["exhaustive_distinct_contexts_per_surface_min"], m["exhaustive_distinct_contexts_per_surface_max"] = float64(minContexts), float64(maxContexts)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepContextSupportIdentifiabilityR4Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-CONTEXT-SUPPORT-IDENTIFIABILITY-R4", m}
}
