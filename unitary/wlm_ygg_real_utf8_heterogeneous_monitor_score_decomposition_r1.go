package unitary

import (
	"math"
	"sort"
)

type wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Scores(background map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Background, diagnostic map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Diagnostic) map[[4]uint8]float64 {
	scores := make(map[[4]uint8]float64)
	for key, bg := range background {
		d := diagnostic[key]
		if bg.occ < 3 || d.occ < 1 {
			continue
		}
		scores[key] = float64(d.err) - float64(bg.err)*float64(d.occ)/float64(bg.occ)
	}
	return scores
}

func wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Top(scores map[[4]uint8]float64) ([4]uint8, bool) {
	keys := make([][4]uint8, 0, len(scores))
	for key := range scores {
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return [4]uint8{}, false
	}
	sort.Slice(keys, func(i, j int) bool {
		if scores[keys[i]] != scores[keys[j]] {
			return scores[keys[i]] > scores[keys[j]]
		}
		return wlmLmRealUTF8BytePilotR2Less(keys[i], keys[j])
	})
	return keys[0], true
}

// RunWlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1 decomposes the frozen R2 monitor miss.
func RunWlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1() interface{} {
	metrics := map[string]float64{
		"valid_evaluation_file_count": 0, "completed_paired_row_count": 0,
		"minimum_target_eligibility_count_per_file": 256, "degraded_monitor_misrank_count": 0,
		"nonzero_counts_flattened_misrank_count": 0, "other_class_misrank_count": 0,
		"maximum_off_target_score_abs_delta": 0, "off_target_score_change_count": 0,
		"missed_target_score_delta": 0, "missed_selected_canonical_minus_target_degraded_score": 0,
		"invalid_decomposition_rows": 0, "training_file_identity_mismatch_count": 0,
		"evaluation_file_identity_mismatch_count": 0, "fresh_eval_blob_overlap_count": 0,
		"tokenizer_use_count": 0, "capacity_growth_event_count": 0,
	}
	train := make([][]byte, 0, len(wlmLmRealUTF8BytePilotR2Train))
	known := make(map[string]bool)
	for _, f := range wlmLmRealUTF8BytePilotR2Train {
		data, ok := wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok {
			metrics["training_file_identity_mismatch_count"]++
		}
		train = append(train, data)
		known[f.sha] = true
	}
	for _, f := range wlmLmRealUTF8BytePilotR2Eval {
		known[f.sha] = true
	}
	baseline, _, selected := wlmLmRealUTF8BytePilotR2TrainModel(train, metrics)
	canonical := wlmYggRealUTF8Level1MaintenanceR3CloneMap(selected)
	if len(canonical) != 256 {
		metrics["invalid_decomposition_rows"]++
	}
	for _, f := range wlmYggRealUTF8HeterogeneousMaintenanceR1Eval {
		if known[f.sha] {
			metrics["fresh_eval_blob_overlap_count"]++
		}
		data, ok := wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok {
			metrics["evaluation_file_identity_mismatch_count"]++
		}
		quarters := wlmYggRealUTF8Level1MaintenanceR3Split(data)
		schedule, eligibilityCount := wlmYggRealUTF8HeterogeneousMaintenanceR1Schedule(quarters, &baseline, canonical)
		metrics["minimum_target_eligibility_count_per_file"] = wlmYggRealUTF8HeterogeneousFailureAttributionR2Min(metrics["minimum_target_eligibility_count_per_file"], float64(eligibilityCount))
		if len(schedule) != 4 {
			metrics["invalid_decomposition_rows"]++
			continue
		}
		background := wlmYggRealUTF8CalibratedMaintenanceR1BackgroundState(quarters[0], &baseline, canonical)
		canonicalDiagnostic := wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(quarters[1], &baseline, canonical, canonical)
		canonicalScores := wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Scores(background, canonicalDiagnostic)
		metrics["valid_evaluation_file_count"]++
		for classIndex, className := range wlmYggRealUTF8HeterogeneousMaintenanceR1Classes {
			targetKey := schedule[classIndex].key
			degraded, degradedOK := wlmYggRealUTF8HeterogeneousMaintenanceR1Degrade(canonical, targetKey, className)
			if !degradedOK {
				metrics["invalid_decomposition_rows"]++
				continue
			}
			degradedDiagnostic := wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(quarters[1], &baseline, degraded, canonical)
			degradedScores := wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Scores(background, degradedDiagnostic)
			selectedKey, selectedOK := wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Top(degradedScores)
			canonicalTargetScore, canonicalTargetOK := canonicalScores[targetKey]
			degradedTargetScore, degradedTargetOK := degradedScores[targetKey]
			if !selectedOK || !canonicalTargetOK || !degradedTargetOK || len(canonicalScores) != len(degradedScores) {
				metrics["invalid_decomposition_rows"]++
				continue
			}
			rowValid := true
			for key, canonicalScore := range canonicalScores {
				degradedScore, present := degradedScores[key]
				if !present {
					rowValid = false
					continue
				}
				if key == targetKey {
					continue
				}
				delta := math.Abs(degradedScore - canonicalScore)
				metrics["maximum_off_target_score_abs_delta"] = wlmYggRealUTF8HeterogeneousFailureAttributionR2Max(metrics["maximum_off_target_score_abs_delta"], delta)
				if delta != 0 {
					metrics["off_target_score_change_count"]++
				}
			}
			if !rowValid {
				metrics["invalid_decomposition_rows"]++
				continue
			}
			if selectedKey != targetKey {
				metrics["degraded_monitor_misrank_count"]++
				if className == "nonzero_counts_flattened" {
					metrics["nonzero_counts_flattened_misrank_count"]++
				} else {
					metrics["other_class_misrank_count"]++
				}
				metrics["missed_target_score_delta"] = degradedTargetScore - canonicalTargetScore
				selectedCanonicalScore, present := canonicalScores[selectedKey]
				if !present {
					metrics["invalid_decomposition_rows"]++
					continue
				}
				metrics["missed_selected_canonical_minus_target_degraded_score"] = selectedCanonicalScore - degradedTargetScore
			}
			metrics["completed_paired_row_count"]++
		}
	}
	return wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Result{Schema: "wingless.research-scientific-result.v1", Experiment: "WLM-YGG-REAL-UTF8-HETEROGENEOUS-MONITOR-SCORE-DECOMPOSITION-R1", Metrics: metrics}
}
