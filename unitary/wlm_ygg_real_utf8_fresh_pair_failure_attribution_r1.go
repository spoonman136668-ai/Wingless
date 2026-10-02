package unitary

import (
	"math"
	"sort"
)

type wlmYggRealUTF8FreshPairFailureAttributionR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

type wlmYggRealUTF8FreshPairFailureAttributionR1Score struct {
	key   [4]uint8
	score float64
}

func wlmYggRealUTF8FreshPairFailureAttributionR1Rank(deltas map[[4]uint8]float64, target [4]uint8) (int, [4]uint8, float64, float64, bool) {
	rows := make([]wlmYggRealUTF8FreshPairFailureAttributionR1Score, 0, len(deltas))
	for key, score := range deltas {
		rows = append(rows, wlmYggRealUTF8FreshPairFailureAttributionR1Score{key, score})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].score != rows[j].score {
			return rows[i].score > rows[j].score
		}
		return wlmLmRealUTF8BytePilotR2Less(rows[i].key, rows[j].key)
	})
	if len(rows) < 2 {
		return 0, [4]uint8{}, 0, 0, false
	}
	rank := 0
	for i, row := range rows {
		if row.key == target {
			rank = i + 1
			break
		}
	}
	return rank, rows[0].key, rows[0].score, rows[0].score - rows[1].score, rank > 0
}

// RunWlmYggRealUTF8FreshPairFailureAttributionR1 independently attributes the frozen R4 negative.
func RunWlmYggRealUTF8FreshPairFailureAttributionR1() interface{} {
	m := map[string]float64{
		"valid_evaluation_file_count": 0, "completed_independent_row_count": 0, "minimum_target_eligibility_count_per_file": 256,
		"minimum_per_class_completed_row_count": 2, "monitor_top1_target_count": 0, "minimum_true_target_rank": 256, "maximum_true_target_rank": 0,
		"minimum_selected_paired_delta_score": math.Inf(1), "minimum_selected_minus_runner_up_paired_delta": math.Inf(1),
		"exact_target_repair_retain_count": 0, "exact_target_repair_revert_count": 0, "exact_target_repair_retain_count_plus_revert_count": 0,
		"minimum_exact_target_selected_key_accuracy_gain": 1, "maximum_exact_target_global_bpb_change": math.Inf(-1),
		"maximum_exact_target_post_repair_audit_prediction_mismatch_count": 0, "maximum_exact_target_post_repair_audit_probability_abs_delta": 0,
		"maximum_candidate_read_bytes": 0, "maximum_candidate_repair_operations": 0, "canonical_selected_motif_count": 0,
		"training_file_identity_mismatch_count": 0, "evaluation_file_identity_mismatch_count": 0, "fresh_eval_blob_overlap_count": 0,
		"invalid_attribution_rows": 0, "lesion_label_access_in_monitor_count": 0, "class_label_access_in_monitor_count": 0, "tokenizer_use_count": 0, "capacity_growth_event_count": 0,
	}
	classCounts := [4]int{}
	train := make([][]byte, 0, len(wlmLmRealUTF8BytePilotR2Train))
	known := make(map[string]bool)
	for _, f := range wlmLmRealUTF8BytePilotR2Train {
		data, ok := wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok {
			m["training_file_identity_mismatch_count"]++
		}
		train = append(train, data)
		known[f.sha] = true
	}
	for _, f := range wlmLmRealUTF8BytePilotR2Eval {
		known[f.sha] = true
	}
	baseline, _, selected := wlmLmRealUTF8BytePilotR2TrainModel(train, m)
	canonical := wlmYggRealUTF8Level1MaintenanceR3CloneMap(selected)
	m["canonical_selected_motif_count"] = float64(len(canonical))
	if len(canonical) != 256 {
		m["invalid_attribution_rows"]++
	}
	for _, f := range wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairEval {
		if known[f.sha] {
			m["fresh_eval_blob_overlap_count"]++
		}
		data, ok := wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok {
			m["evaluation_file_identity_mismatch_count"]++
		}
		q := wlmYggRealUTF8Level1MaintenanceR3Split(data)
		schedule, eligible := wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairSchedule(q, &baseline, canonical)
		m["minimum_target_eligibility_count_per_file"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMin(m["minimum_target_eligibility_count_per_file"], float64(eligible))
		if len(schedule) != 4 {
			m["invalid_attribution_rows"]++
			continue
		}
		background := wlmYggRealUTF8CalibratedMaintenanceR1BackgroundState(q[0], &baseline, canonical)
		healthy := wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(q[1], &baseline, canonical, canonical)
		healthyScores := wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Scores(background, healthy)
		if len(background) != 256 || len(healthyScores) == 0 {
			m["invalid_attribution_rows"]++
			continue
		}
		m["valid_evaluation_file_count"]++
		for ci, className := range wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairClasses {
			target := schedule[ci].key
			degraded, ok := wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairDegrade(canonical, target, className)
			if !ok {
				m["invalid_attribution_rows"]++
				continue
			}
			state := wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(q[1], &baseline, degraded, canonical)
			scores := wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Scores(background, state)
			if len(scores) != len(healthyScores) {
				m["invalid_attribution_rows"]++
				continue
			}
			deltas := make(map[[4]uint8]float64, len(healthyScores))
			valid := true
			for key, base := range healthyScores {
				value, present := scores[key]
				if !present {
					valid = false
					continue
				}
				deltas[key] = value - base
			}
			if !valid {
				m["invalid_attribution_rows"]++
				continue
			}
			rank, top, score, margin, ok := wlmYggRealUTF8FreshPairFailureAttributionR1Rank(deltas, target)
			if !ok {
				m["invalid_attribution_rows"]++
				continue
			}
			if top == target {
				m["monitor_top1_target_count"]++
			}
			m["minimum_true_target_rank"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMin(m["minimum_true_target_rank"], float64(rank))
			m["maximum_true_target_rank"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMax(m["maximum_true_target_rank"], float64(rank))
			m["minimum_selected_paired_delta_score"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMin(m["minimum_selected_paired_delta_score"], score)
			m["minimum_selected_minus_runner_up_paired_delta"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMin(m["minimum_selected_minus_runner_up_paired_delta"], margin)
			preAcc, count := wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(q[2], target, &baseline, degraded)
			if count < 3 {
				m["invalid_attribution_rows"]++
				continue
			}
			preBPB := wlmYggRealUTF8Level1MaintenanceR3BPB(q[2], &baseline, degraded)
			repaired := wlmYggRealUTF8Level1MaintenanceR3CloneMap(degraded)
			repaired[target] = wlmYggRealUTF8Level1MaintenanceR3CloneEntry(canonical[target])
			postAcc, _ := wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(q[2], target, &baseline, repaired)
			postBPB := wlmYggRealUTF8Level1MaintenanceR3BPB(q[2], &baseline, repaired)
			gain := postAcc - preAcc
			bpb := postBPB - preBPB
			m["minimum_exact_target_selected_key_accuracy_gain"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMin(m["minimum_exact_target_selected_key_accuracy_gain"], gain)
			m["maximum_exact_target_global_bpb_change"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMax(m["maximum_exact_target_global_bpb_change"], bpb)
			if gain >= 0.20 && bpb <= 0.005 {
				m["exact_target_repair_retain_count"]++
			} else {
				m["exact_target_repair_revert_count"]++
			}
			mismatch, prob := wlmYggRealUTF8Level1MaintenanceR3Audit(q[3], &baseline, repaired, canonical)
			m["maximum_exact_target_post_repair_audit_prediction_mismatch_count"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMax(m["maximum_exact_target_post_repair_audit_prediction_mismatch_count"], float64(mismatch))
			m["maximum_exact_target_post_repair_audit_probability_abs_delta"] = wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairMax(m["maximum_exact_target_post_repair_audit_probability_abs_delta"], prob)
			m["maximum_candidate_read_bytes"] = 1028
			m["maximum_candidate_repair_operations"] = 1
			classCounts[ci]++
			m["completed_independent_row_count"]++
		}
	}
	minClass := classCounts[0]
	for _, count := range classCounts[1:] {
		if count < minClass {
			minClass = count
		}
	}
	m["minimum_per_class_completed_row_count"] = float64(minClass)
	m["exact_target_repair_retain_count_plus_revert_count"] = m["exact_target_repair_retain_count"] + m["exact_target_repair_revert_count"]
	control := RunWlmYggRealUTF8HeterogeneousMaintenanceR4FreshPair().(wlmYggRealUTF8HeterogeneousMaintenanceR4FreshPairResult)
	controls := map[string]string{"sequential_control_valid_evaluation_file_count": "valid_evaluation_file_count", "sequential_control_completed_cycle_count": "completed_cycle_count", "sequential_control_minimum_monitor_target_selection_accuracy": "minimum_monitor_target_selection_accuracy", "sequential_control_wrong_candidate_retain_count": "wrong_candidate_retain_count", "sequential_control_wrong_candidate_revert_count": "wrong_candidate_revert_count", "sequential_control_correct_candidate_retain_count": "correct_candidate_retain_count", "sequential_control_correct_candidate_revert_count": "correct_candidate_revert_count", "sequential_control_maximum_post_repair_audit_prediction_mismatch_count": "maximum_post_repair_audit_prediction_mismatch_count", "sequential_control_maximum_post_repair_audit_probability_abs_delta": "maximum_post_repair_audit_probability_abs_delta", "sequential_control_minimum_structure_jaccard_after_each_cycle": "minimum_structure_jaccard_after_each_cycle"}
	for dst, src := range controls {
		m[dst] = control.Metrics[src]
	}
	return wlmYggRealUTF8FreshPairFailureAttributionR1Result{"wingless.research-scientific-result.v1", "WLM-YGG-REAL-UTF8-FRESH-PAIR-FAILURE-ATTRIBUTION-R1", m}
}
