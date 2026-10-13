package unitary

import (
	"fmt"
	"math"
	"sort"
)

type wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context [5]uint8

type wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Bridge() [][]uint8 {
	records := make([][]uint8, 0, 432)
	for entity := 0; entity < 9; entity++ {
		for family := 0; family < 3; family++ {
			for contextIndex := 0; contextIndex < 16; contextIndex++ {
				state := uint32(0x85ebca6b) ^ uint32((entity+1)*1009+(family+1)*10007+(contextIndex+1)*65537)
				record := append([]uint8{}, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
				record = append(record, uint8(160+contextIndex))
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

func wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(records [][]uint8, reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) ([]map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, int) {
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	counts := make([]map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32, len(reps))
	for i := range counts {
		counts[i] = make(map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32)
	}
	invalid := 0
	for _, record := range records {
		matches := 0
		for i := 3; i+6 <= len(record); i++ {
			key := wlmSiRawRepAliasInvarianceFalsificationR1Key{record[i], record[i+1], record[i+2], record[i+3]}
			id, ok := index[key]
			if !ok {
				continue
			}
			context := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context{record[i-3], record[i-2], record[i-1], record[i+4], record[i+5]}
			counts[id][context]++
			matches++
			i += 3
		}
		if matches != 1 {
			invalid++
		}
	}
	return counts, invalid
}

func wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(counts []map[wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context]uint32) wlmLmRawRepContextualAliasInductionR2Classes {
	grouped := make(map[string][]int)
	for id, contexts := range counts {
		keys := make([]wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Context, 0, len(contexts))
		for key := range contexts {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			for k := 0; k < 5; k++ {
				if keys[i][k] != keys[j][k] {
					return keys[i][k] < keys[j][k]
				}
			}
			return false
		})
		signature := ""
		for _, key := range keys {
			signature += fmt.Sprintf("%02x%02x%02x%02x%02x;", key[0], key[1], key[2], key[3], key[4])
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
	return classes
}

func RunWlmSiRawRepContextPhaseAnchorIdentifiabilityR5() interface{} {
	m := map[string]float64{
		"anchored_bridge_record_count": 0, "anchored_records_per_surface_min": 0, "anchored_records_per_surface_max": 0,
		"anchored_distinct_contexts_per_surface_min": 0, "anchored_distinct_contexts_per_surface_max": 0,
		"anchored_distinct_phase_anchors_per_surface_min": 0, "anchored_distinct_phase_anchors_per_surface_max": 0,
		"unanchored_control_bridge_record_count": 0, "selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"unanchored_control_alias_class_count": 0, "unanchored_control_three_member_class_count": 0, "unanchored_control_correct_alias_triplet_count": 0, "unanchored_control_cross_entity_collision_pair_count": 0,
		"anchored_alias_class_count": 0, "anchored_three_member_class_count": 0, "anchored_correct_alias_triplet_count": 0, "anchored_cross_entity_collision_pair_count": 0,
		"evaluation_entity_pair_count": 36, "task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0,
		"capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	unanchored := wlmSiRawRepContextSupportIdentifiabilityR4Bridge()
	anchored := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Bridge()
	m["unanchored_control_bridge_record_count"] = float64(len(unanchored))
	m["anchored_bridge_record_count"] = float64(len(anchored))
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	unanchoredClasses := wlmSiRawRepContextSupportFalsificationR3Induce(unanchored, reps)
	anchoredCounts, invalid := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Counts(anchored, reps)
	m["invalid_row_count"] = float64(invalid)
	anchoredClasses := wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Induce(anchoredCounts)
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	ua, ut, uc := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(unanchoredClasses, index)
	aa, at, ac := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(anchoredClasses, index)
	m["unanchored_control_alias_class_count"], m["unanchored_control_three_member_class_count"], m["unanchored_control_correct_alias_triplet_count"] = float64(ua), float64(ut), float64(uc)
	m["anchored_alias_class_count"], m["anchored_three_member_class_count"], m["anchored_correct_alias_triplet_count"] = float64(aa), float64(at), float64(ac)
	m["unanchored_control_cross_entity_collision_pair_count"] = float64(wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(unanchoredClasses, index))
	m["anchored_cross_entity_collision_pair_count"] = float64(wlmSiRawRepContextSupportIdentifiabilityR4CrossEntityPairs(anchoredClasses, index))
	minRecords, maxRecords := int(^uint(0)>>1), 0
	minContexts, maxContexts := int(^uint(0)>>1), 0
	minAnchors, maxAnchors := int(^uint(0)>>1), 0
	for _, contexts := range anchoredCounts {
		records := 0
		anchors := make(map[uint8]bool)
		for context, count := range contexts {
			records += int(count)
			anchors[context[0]] = true
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
		if len(anchors) < minAnchors {
			minAnchors = len(anchors)
		}
		if len(anchors) > maxAnchors {
			maxAnchors = len(anchors)
		}
	}
	m["anchored_records_per_surface_min"], m["anchored_records_per_surface_max"] = float64(minRecords), float64(maxRecords)
	m["anchored_distinct_contexts_per_surface_min"], m["anchored_distinct_contexts_per_surface_max"] = float64(minContexts), float64(maxContexts)
	m["anchored_distinct_phase_anchors_per_surface_min"], m["anchored_distinct_phase_anchors_per_surface_max"] = float64(minAnchors), float64(maxAnchors)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepContextPhaseAnchorIdentifiabilityR5Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-CONTEXT-PHASE-ANCHOR-IDENTIFIABILITY-R5", m}
}
