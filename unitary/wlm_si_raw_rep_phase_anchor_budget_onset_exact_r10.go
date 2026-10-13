package unitary

import "math"

type wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Restored(prefix string, m map[string]float64) bool {
	return m[prefix+"_alias_class_count"] == 9 &&
		m[prefix+"_three_member_class_count"] == 9 &&
		m[prefix+"_correct_alias_triplet_count"] == 9 &&
		m[prefix+"_split_entity_triplet_count"] == 0 &&
		m[prefix+"_cross_entity_collision_pair_count"] == 0
}

func RunWlmSiRawRepPhaseAnchorBudgetOnsetExactR10() interface{} {
	m := map[string]float64{
		"selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0,
		"nested_prefix_mismatch_count": 0, "evaluation_entity_pair_count_per_arm": 36,
		"first_restored_budget": 0, "restoration_reversal_count": 0,
		"class_count_increase_count": 0, "split_triplet_increase_count": 0,
		"task_training_record_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0,
		"capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}
	budgets := []int{104, 105, 106, 107, 108}
	records := make([][][]uint8, len(budgets))
	for i, budget := range budgets {
		records[i] = wlmSiRawRepPhaseAnchorSampleBudgetR7Bridge(budget)
	}
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(wlmSiRawRepContextSamplingFalsificationR2Bridge(false))
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	invalid := 0
	for i, budget := range budgets {
		prefix := "budget" + wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Itoa(budget)
		invalid += wlmSiRawRepPhaseAnchorSampleBudgetR7ArmMetrics(prefix, records[i], reps, m)
		if i > 0 {
			m["nested_prefix_mismatch_count"] += float64(wlmSiRawRepPhaseAnchorSampleBudgetR7PrefixMismatch(records[i-1], records[i], budgets[i-1], budget))
			previous := "budget" + wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Itoa(budgets[i-1])
			if m[prefix+"_alias_class_count"] > m[previous+"_alias_class_count"] {
				m["class_count_increase_count"]++
			}
			if m[prefix+"_split_entity_triplet_count"] > m[previous+"_split_entity_triplet_count"] {
				m["split_triplet_increase_count"]++
			}
			if wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Restored(previous, m) && !wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Restored(prefix, m) {
				m["restoration_reversal_count"]++
			}
		}
	}
	for _, budget := range []int{105, 106, 107, 108} {
		prefix := "budget" + wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Itoa(budget)
		if wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Restored(prefix, m) {
			m["first_restored_budget"] = float64(budget)
			break
		}
	}
	m["invalid_row_count"] = float64(invalid)
	for _, value := range m {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Result{
		Schema: "wingless.research-scientific-result.v1", Experiment: "WLM-SI-RAW-REP-PHASE-ANCHOR-BUDGET-ONSET-EXACT-R10", Metrics: m,
	}
}

func wlmSiRawRepPhaseAnchorBudgetOnsetExactR10Itoa(value int) string {
	switch value {
	case 104:
		return "104"
	case 105:
		return "105"
	case 106:
		return "106"
	case 107:
		return "107"
	case 108:
		return "108"
	default:
		panic("unregistered budget")
	}
}
