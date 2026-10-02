package unitary

import "math"

type wlmYggRealUTF8PairedDeltaMonitorR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

// RunWlmYggRealUTF8PairedDeltaMonitorR1 compares the frozen legacy score with a per-key paired delta.
func RunWlmYggRealUTF8PairedDeltaMonitorR1() interface{} {
	metrics := map[string]float64{
		"valid_evaluation_file_count": 0, "completed_row_count": 0,
		"minimum_target_eligibility_count_per_file": 256, "legacy_misrank_count": 0,
		"legacy_nonzero_counts_flattened_misrank_count": 0, "legacy_other_class_misrank_count": 0,
		"paired_delta_top1_target_count": 0, "paired_delta_nonpositive_target_count": 0,
		"paired_delta_positive_off_target_count": 0, "minimum_paired_delta_target_margin": 0,
		"invalid_row_count": 0, "training_file_identity_mismatch_count": 0,
		"evaluation_file_identity_mismatch_count": 0, "fresh_eval_blob_overlap_count": 0,
		"tokenizer_use_count": 0, "capacity_growth_event_count": 0,
	}
	minimumMargin := math.Inf(1)
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
		metrics["invalid_row_count"]++
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
			metrics["invalid_row_count"]++
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
				metrics["invalid_row_count"]++
				continue
			}
			degradedDiagnostic := wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(quarters[1], &baseline, degraded, canonical)
			degradedScores := wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Scores(background, degradedDiagnostic)
			legacyKey, legacyOK := wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Top(degradedScores)
			if !legacyOK || len(canonicalScores) != len(degradedScores) {
				metrics["invalid_row_count"]++
				continue
			}
			deltas := make(map[[4]uint8]float64, len(canonicalScores))
			rowValid := true
			for key, canonicalScore := range canonicalScores {
				degradedScore, present := degradedScores[key]
				if !present {
					rowValid = false
					continue
				}
				delta := degradedScore - canonicalScore
				deltas[key] = delta
				if key != targetKey && delta > 0 {
					metrics["paired_delta_positive_off_target_count"]++
				}
			}
			if !rowValid {
				metrics["invalid_row_count"]++
				continue
			}
			deltaKey, deltaOK := wlmYggRealUTF8HeterogeneousMonitorScoreDecompositionR1Top(deltas)
			targetDelta, targetOK := deltas[targetKey]
			if !deltaOK || !targetOK {
				metrics["invalid_row_count"]++
				continue
			}
			if legacyKey != targetKey {
				metrics["legacy_misrank_count"]++
				if className == "nonzero_counts_flattened" {
					metrics["legacy_nonzero_counts_flattened_misrank_count"]++
				} else {
					metrics["legacy_other_class_misrank_count"]++
				}
			}
			if targetDelta <= 0 {
				metrics["paired_delta_nonpositive_target_count"]++
			}
			if deltaKey == targetKey {
				metrics["paired_delta_top1_target_count"]++
			}
			secondScore := math.Inf(-1)
			for key, score := range deltas {
				if key != targetKey && score > secondScore {
					secondScore = score
				}
			}
			margin := targetDelta - secondScore
			if margin < minimumMargin {
				minimumMargin = margin
			}
			metrics["completed_row_count"]++
		}
	}
	if metrics["completed_row_count"] == 8 && !math.IsInf(minimumMargin, 0) {
		metrics["minimum_paired_delta_target_margin"] = minimumMargin
	} else {
		metrics["invalid_row_count"]++
	}
	return wlmYggRealUTF8PairedDeltaMonitorR1Result{Schema: "wingless.research-scientific-result.v1", Experiment: "WLM-YGG-REAL-UTF8-PAIRED-DELTA-MONITOR-R1", Metrics: metrics}
}
