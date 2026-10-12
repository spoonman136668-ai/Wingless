package unitary

import (
	"fmt"
	"math"
	"sort"
)

type wlmSiRawRepContextSupportFalsificationR3Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepContextSupportFalsificationR3Induce(records [][]uint8, reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) wlmLmRawRepContextualAliasInductionR2Classes {
	counts := wlmSiRawRepContextSamplingFalsificationR2Counts(records, reps)
	grouped := make(map[string][]int)
	for id, contexts := range counts {
		keys := make([]wlmLmRawRepContextualAliasInductionR2Context, 0, len(contexts))
		for key := range contexts {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			for k := 0; k < 4; k++ {
				if keys[i][k] != keys[j][k] {
					return keys[i][k] < keys[j][k]
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
	classes := wlmLmRawRepContextualAliasInductionR2Classes{surfaceToClass: make([]int, len(reps)), members: make([][]int, 0, len(signatures))}
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
func wlmSiRawRepContextSupportFalsificationR3CrossEntityCollisions(classes wlmLmRawRepContextualAliasInductionR2Classes, index map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int) int {
	collisions := 0
	for le := 0; le < 9; le++ {
		for re := le + 1; re < 9; re++ {
			for lf := 0; lf < 3; lf++ {
				for rf := 0; rf < 3; rf++ {
					l := index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(lf, le)]
					r := index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(rf, re)]
					if classes.surfaceToClass[l] == classes.surfaceToClass[r] {
						collisions++
					}
				}
			}
		}
	}
	return collisions
}
func RunWlmSiRawRepContextSupportFalsificationR3() interface{} {
	m := map[string]float64{"matched_bridge_record_count": 0, "independent_bridge_record_count": 0, "selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "independent_max_context_count_l1_distance": 0, "exact_count_control_alias_class_count": 0, "exact_count_control_correct_alias_triplet_count": 0, "exact_count_control_missing_transition_count": 0, "exact_count_control_unseen_family_accuracy": 0, "matched_support_alias_class_count": 0, "matched_support_three_member_class_count": 0, "matched_support_correct_alias_triplet_count": 0, "matched_support_unseen_family_accuracy": 0, "independent_support_alias_class_count": 0, "independent_support_three_member_class_count": 0, "independent_support_correct_alias_triplet_count": 0, "independent_support_missing_transition_count": 0, "independent_support_unseen_family_accuracy": 0, "independent_support_cross_entity_collision_count": 0, "task_training_record_count": 0, "family2_task_training_record_count": 0, "unseen_family_evaluation_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0}
	matched, independent := wlmSiRawRepContextSamplingFalsificationR2Bridge(false), wlmSiRawRepContextSamplingFalsificationR2Bridge(true)
	m["matched_bridge_record_count"] = float64(len(matched))
	m["independent_bridge_record_count"] = float64(len(independent))
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(matched)
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	m["independent_max_context_count_l1_distance"] = float64(wlmSiRawRepContextSamplingFalsificationR2MaxL1(wlmSiRawRepContextSamplingFalsificationR2Counts(independent, reps), reps))
	exact := wlmLmRawRepContextualAliasInductionR2Induce(reps, independent)
	matchedSupport := wlmSiRawRepContextSupportFalsificationR3Induce(matched, reps)
	independentSupport := wlmSiRawRepContextSupportFalsificationR3Induce(independent, reps)
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	ea, _, ec := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(exact, index)
	ma, mt, mc := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(matchedSupport, index)
	ia, it, ic := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(independentSupport, index)
	m["exact_count_control_alias_class_count"], m["exact_count_control_correct_alias_triplet_count"] = float64(ea), float64(ec)
	m["matched_support_alias_class_count"], m["matched_support_three_member_class_count"], m["matched_support_correct_alias_triplet_count"] = float64(ma), float64(mt), float64(mc)
	m["independent_support_alias_class_count"], m["independent_support_three_member_class_count"], m["independent_support_correct_alias_triplet_count"] = float64(ia), float64(it), float64(ic)
	m["independent_support_cross_entity_collision_count"] = float64(wlmSiRawRepContextSupportFalsificationR3CrossEntityCollisions(independentSupport, index))
	tasks := wlmLmRawRepContextualAliasInductionR2TaskRecords()
	m["task_training_record_count"] = float64(len(tasks))
	exactModel := wlmSiRawRepContextSamplingFalsificationR2Build(tasks, reps, exact, m)
	matchedModel := wlmSiRawRepContextSamplingFalsificationR2Build(tasks, reps, matchedSupport, m)
	independentModel := wlmSiRawRepContextSamplingFalsificationR2Build(tasks, reps, independentSupport, m)
	exactHits, exactMissing := wlmSiRawRepContextSamplingFalsificationR2Evaluate(reps, exact, exactModel, m)
	matchedHits, _ := wlmSiRawRepContextSamplingFalsificationR2Evaluate(reps, matchedSupport, matchedModel, m)
	independentHits, independentMissing := wlmSiRawRepContextSamplingFalsificationR2Evaluate(reps, independentSupport, independentModel, m)
	m["unseen_family_evaluation_count"] = 80
	m["exact_count_control_missing_transition_count"] = float64(exactMissing)
	m["exact_count_control_unseen_family_accuracy"] = float64(exactHits) / 80
	m["matched_support_unseen_family_accuracy"] = float64(matchedHits) / 80
	m["independent_support_missing_transition_count"] = float64(independentMissing)
	m["independent_support_unseen_family_accuracy"] = float64(independentHits) / 80
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepContextSupportFalsificationR3Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-CONTEXT-SUPPORT-FALSIFICATION-R3", m}
}
