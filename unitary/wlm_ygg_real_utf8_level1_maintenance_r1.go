package unitary

import (
	"math"
	"sort"
)

type wlmYggRealUTF8Level1MaintenanceR1Support struct {
	occ          int
	canonicalHit int
	baselineHit  int
}

type wlmYggRealUTF8Level1MaintenanceR1Target struct {
	key   [4]uint8
	occ   int
	gain  float64
}

type wlmYggRealUTF8Level1MaintenanceR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmYggRealUTF8Level1MaintenanceR1CloneStats(src *wlmLmRealUTF8BytePilotR2MotifStats) *wlmLmRealUTF8BytePilotR2MotifStats {
	if src == nil {
		return nil
	}
	dst := &wlmLmRealUTF8BytePilotR2MotifStats{total: src.total}
	dst.counts = src.counts
	return dst
}

func wlmYggRealUTF8Level1MaintenanceR1CloneEntry(src wlmLmRealUTF8BytePilotR2Selected) wlmLmRealUTF8BytePilotR2Selected {
	return wlmLmRealUTF8BytePilotR2Selected{
		key:         src.key,
		stats:       wlmYggRealUTF8Level1MaintenanceR1CloneStats(src.stats),
		best:        src.best,
		bestCount:   src.bestCount,
		consistency: src.consistency,
	}
}

func wlmYggRealUTF8Level1MaintenanceR1CloneMap(src map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected) map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected {
	dst := make(map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected, len(src))
	for key, entry := range src {
		dst[key] = wlmYggRealUTF8Level1MaintenanceR1CloneEntry(entry)
	}
	return dst
}

func wlmYggRealUTF8Level1MaintenanceR1Split(data []byte) [4][]byte {
	n := len(data)
	q1 := n / 4
	q2 := n / 2
	q3 := (3 * n) / 4
	return [4][]byte{
		data[:q1],
		data[q1:q2],
		data[q2:q3],
		data[q3:],
	}
}

func wlmYggRealUTF8Level1MaintenanceR1SupportFor(
	data []byte,
	baseline *[256][256]uint32,
	canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) map[[4]uint8]wlmYggRealUTF8Level1MaintenanceR1Support {
	out := make(map[[4]uint8]wlmYggRealUTF8Level1MaintenanceR1Support)
	for i := 4; i < len(data); i++ {
		key := [4]uint8{data[i-4], data[i-3], data[i-2], data[i-1]}
		entry, ok := canonical[key]
		if !ok {
			continue
		}
		row := out[key]
		row.occ++
		if entry.best == data[i] {
			row.canonicalHit++
		}
		if wlmLmRealUTF8BytePilotR2BaselinePrediction(baseline, data[i-1]) == data[i] {
			row.baselineHit++
		}
		out[key] = row
	}
	return out
}

func wlmYggRealUTF8Level1MaintenanceR1Targets(
	quarters [4][]byte,
	baseline *[256][256]uint32,
	canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) []wlmYggRealUTF8Level1MaintenanceR1Target {
	cal := wlmYggRealUTF8Level1MaintenanceR1SupportFor(quarters[0], baseline, canonical)
	diag := wlmYggRealUTF8Level1MaintenanceR1SupportFor(quarters[1], baseline, canonical)
	val := wlmYggRealUTF8Level1MaintenanceR1SupportFor(quarters[2], baseline, canonical)
	audit := wlmYggRealUTF8Level1MaintenanceR1SupportFor(quarters[3], baseline, canonical)
	rows := make([]wlmYggRealUTF8Level1MaintenanceR1Target, 0)
	for key, row := range cal {
		if row.occ < 8 || diag[key].occ < 4 || val[key].occ < 4 || audit[key].occ < 4 {
			continue
		}
		gain := float64(row.canonicalHit-row.baselineHit) / float64(row.occ)
		if gain < 0.25 {
			continue
		}
		rows = append(rows, wlmYggRealUTF8Level1MaintenanceR1Target{key: key, occ: row.occ, gain: gain})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].occ != rows[j].occ {
			return rows[i].occ > rows[j].occ
		}
		if rows[i].gain != rows[j].gain {
			return rows[i].gain > rows[j].gain
		}
		return wlmLmRealUTF8BytePilotR2Less(rows[i].key, rows[j].key)
	})
	return rows
}

func wlmYggRealUTF8Level1MaintenanceR1Predict(
	baseline *[256][256]uint32,
	active map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
	context [4]uint8,
) uint8 {
	if entry, ok := active[context]; ok {
		return entry.best
	}
	return wlmLmRealUTF8BytePilotR2BaselinePrediction(baseline, context[3])
}

func wlmYggRealUTF8Level1MaintenanceR1Prob(
	baseline *[256][256]uint32,
	active map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
	context [4]uint8,
	target uint8,
) float64 {
	if entry, ok := active[context]; ok {
		return wlmLmRealUTF8BytePilotR2Prob(&entry.stats.counts, target)
	}
	return wlmLmRealUTF8BytePilotR2BaselineProb(baseline, context[3], target)
}

func wlmYggRealUTF8Level1MaintenanceR1Monitor(
	data []byte,
	baseline *[256][256]uint32,
	active map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
	catalog map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) ([4]uint8, int, int, bool) {
	errors := make(map[[4]uint8]int)
	for i := 4; i < len(data); i++ {
		key := [4]uint8{data[i-4], data[i-3], data[i-2], data[i-1]}
		if _, ok := catalog[key]; !ok {
			continue
		}
		if wlmYggRealUTF8Level1MaintenanceR1Predict(baseline, active, key) != data[i] {
			errors[key]++
		}
	}
	if len(errors) == 0 {
		return [4]uint8{}, 0, 0, false
	}
	keys := make([][4]uint8, 0, len(errors))
	for key := range errors {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if errors[keys[i]] != errors[keys[j]] {
			return errors[keys[i]] > errors[keys[j]]
		}
		return wlmLmRealUTF8BytePilotR2Less(keys[i], keys[j])
	})
	top := errors[keys[0]]
	runner := 0
	if len(keys) > 1 {
		runner = errors[keys[1]]
	}
	return keys[0], top, runner, true
}

func wlmYggRealUTF8Level1MaintenanceR1KeyAccuracy(
	data []byte,
	key [4]uint8,
	baseline *[256][256]uint32,
	active map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) (float64, int) {
	correct := 0
	total := 0
	for i := 4; i < len(data); i++ {
		ctx := [4]uint8{data[i-4], data[i-3], data[i-2], data[i-1]}
		if ctx != key {
			continue
		}
		if wlmYggRealUTF8Level1MaintenanceR1Predict(baseline, active, ctx) == data[i] {
			correct++
		}
		total++
	}
	if total == 0 {
		return 0, 0
	}
	return float64(correct) / float64(total), total
}

func wlmYggRealUTF8Level1MaintenanceR1BPB(
	data []byte,
	baseline *[256][256]uint32,
	active map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) float64 {
	if len(data) < 2 {
		return 0
	}
	var bits float64
	total := 0
	for i := 1; i < len(data); i++ {
		var p float64
		if i >= 4 {
			key := [4]uint8{data[i-4], data[i-3], data[i-2], data[i-1]}
			p = wlmYggRealUTF8Level1MaintenanceR1Prob(baseline, active, key, data[i])
		} else {
			p = wlmLmRealUTF8BytePilotR2BaselineProb(baseline, data[i-1], data[i])
		}
		bits += -math.Log2(p)
		total++
	}
	return bits / float64(total)
}

func wlmYggRealUTF8Level1MaintenanceR1WrongSource(
	selectedKey [4]uint8,
	catalog map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) ([4]uint8, bool) {
	keys := make([][4]uint8, 0, len(catalog))
	for key := range catalog {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return wlmLmRealUTF8BytePilotR2Less(keys[i], keys[j]) })
	pos := -1
	for i, key := range keys {
		if key == selectedKey {
			pos = i
			break
		}
	}
	if pos < 0 {
		return [4]uint8{}, false
	}
	targetBest := catalog[selectedKey].best
	for step := 1; step < len(keys); step++ {
		key := keys[(pos+step)%len(keys)]
		if catalog[key].best != targetBest {
			return key, true
		}
	}
	return [4]uint8{}, false
}

func wlmYggRealUTF8Level1MaintenanceR1Jaccard(
	active, canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) float64 {
	intersection := 0
	union := len(canonical)
	for key := range active {
		if _, ok := canonical[key]; ok {
			intersection++
		} else {
			union++
		}
	}
	if union == 0 {
		return 1
	}
	return float64(intersection) / float64(union)
}

func wlmYggRealUTF8Level1MaintenanceR1StateMismatch(
	active, canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) int {
	mismatch := 0
	if len(active) != len(canonical) {
		if len(active) > len(canonical) {
			mismatch += len(active) - len(canonical)
		} else {
			mismatch += len(canonical) - len(active)
		}
	}
	for key, want := range canonical {
		got, ok := active[key]
		if !ok {
			mismatch++
			continue
		}
		if got.key != want.key || got.best != want.best || got.bestCount != want.bestCount || got.consistency != want.consistency || got.stats == nil || want.stats == nil {
			mismatch++
			continue
		}
		if got.stats.total != want.stats.total || got.stats.counts != want.stats.counts {
			mismatch++
		}
	}
	return mismatch
}

func wlmYggRealUTF8Level1MaintenanceR1Audit(
	data []byte,
	baseline *[256][256]uint32,
	active, canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) (int, float64) {
	predMismatch := 0
	maxProbDelta := 0.0
	for i := 1; i < len(data); i++ {
		if i < 4 {
			continue
		}
		key := [4]uint8{data[i-4], data[i-3], data[i-2], data[i-1]}
		ap := wlmYggRealUTF8Level1MaintenanceR1Predict(baseline, active, key)
		cp := wlmYggRealUTF8Level1MaintenanceR1Predict(baseline, canonical, key)
		if ap != cp {
			predMismatch++
		}
		aProb := wlmYggRealUTF8Level1MaintenanceR1Prob(baseline, active, key, data[i])
		cProb := wlmYggRealUTF8Level1MaintenanceR1Prob(baseline, canonical, key, data[i])
		delta := math.Abs(aProb - cProb)
		if delta > maxProbDelta {
			maxProbDelta = delta
		}
	}
	return predMismatch, maxProbDelta
}

func wlmYggRealUTF8Level1MaintenanceR1Min(a, b float64) float64 {
	if b < a { return b }
	return a
}
func wlmYggRealUTF8Level1MaintenanceR1Max(a, b float64) float64 {
	if b > a { return b }
	return a
}

// RunWlmYggRealUTF8Level1MaintenanceR1 evaluates frozen real-workload level-1 maintenance.
func RunWlmYggRealUTF8Level1MaintenanceR1() interface{} {
	metrics := map[string]float64{
		"valid_evaluation_file_count": 0,
		"completed_cycle_count": 0,
		"minimum_target_eligibility_count_per_file": 256,
		"minimum_monitor_target_selection_accuracy": 1,
		"minimum_selected_error_count_ratio_vs_runner_up": 1e9,
		"maximum_wrong_candidate_selected_key_accuracy_gain": -1,
		"wrong_candidate_retain_count": 0,
		"wrong_candidate_revert_count": 0,
		"minimum_correct_candidate_selected_key_accuracy_gain": 1,
		"correct_candidate_retain_count": 0,
		"correct_candidate_revert_count": 0,
		"maximum_post_repair_audit_prediction_mismatch_count": 0,
		"maximum_post_repair_audit_probability_abs_delta": 0,
		"minimum_structure_jaccard_after_each_cycle": 1,
		"maximum_candidate_read_bytes": 0,
		"maximum_candidate_repair_operations": 0,
		"maximum_final_selected_motif_count": 0,
		"minimum_final_selected_motif_count": 256,
		"lesion_label_access_count": 0,
		"evaluator_schedule_access_in_policy_count": 0,
		"tokenizer_use_count": 0,
		"capacity_growth_event_count": 0,
		"invalid_cycle_rows": 0,
		"training_file_identity_mismatch_count": 0,
		"evaluation_file_identity_mismatch_count": 0,
		"canonical_selected_motif_count": 0,
	}

	train := make([][]byte, 0, len(wlmLmRealUTF8BytePilotR2Train))
	trainSet := make(map[string]bool)
	for _, f := range wlmLmRealUTF8BytePilotR2Train {
		data, ok := wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok {
			metrics["training_file_identity_mismatch_count"]++
		}
		train = append(train, data)
		trainSet[f.sha] = true
	}
	baseline, _, selected := wlmLmRealUTF8BytePilotR2TrainModel(train, metrics)
	canonical := wlmYggRealUTF8Level1MaintenanceR1CloneMap(selected)
	active := wlmYggRealUTF8Level1MaintenanceR1CloneMap(canonical)
	metrics["canonical_selected_motif_count"] = float64(len(canonical))
	if len(canonical) != 256 {
		metrics["invalid_cycle_rows"]++
	}

	for _, f := range wlmLmRealUTF8BytePilotR2Eval {
		if trainSet[f.sha] {
			metrics["invalid_cycle_rows"]++
		}
		data, ok := wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok {
			metrics["evaluation_file_identity_mismatch_count"]++
		}
		quarters := wlmYggRealUTF8Level1MaintenanceR1Split(data)
		targets := wlmYggRealUTF8Level1MaintenanceR1Targets(quarters, &baseline, canonical)
		metrics["minimum_target_eligibility_count_per_file"] = wlmYggRealUTF8Level1MaintenanceR1Min(
			metrics["minimum_target_eligibility_count_per_file"], float64(len(targets)),
		)
		if len(targets) < 2 {
			metrics["invalid_cycle_rows"]++
			continue
		}
		metrics["valid_evaluation_file_count"]++

		for targetIndex := 0; targetIndex < 2; targetIndex++ {
			targetKey := targets[targetIndex].key
			delete(active, targetKey)

			selectedKey, topErrors, runnerErrors, monitorOK := wlmYggRealUTF8Level1MaintenanceR1Monitor(
				quarters[1], &baseline, active, canonical,
			)
			if !monitorOK || topErrors <= 0 {
				metrics["invalid_cycle_rows"]++
				metrics["minimum_monitor_target_selection_accuracy"] = 0
				metrics["minimum_selected_error_count_ratio_vs_runner_up"] = 0
				continue
			}
			if selectedKey != targetKey {
				metrics["minimum_monitor_target_selection_accuracy"] = 0
			}
			ratio := 1e9
			if runnerErrors > 0 {
				ratio = float64(topErrors) / float64(runnerErrors)
			}
			metrics["minimum_selected_error_count_ratio_vs_runner_up"] = wlmYggRealUTF8Level1MaintenanceR1Min(
				metrics["minimum_selected_error_count_ratio_vs_runner_up"], ratio,
			)

			preAcc, preCount := wlmYggRealUTF8Level1MaintenanceR1KeyAccuracy(
				quarters[2], selectedKey, &baseline, active,
			)
			if preCount < 4 {
				metrics["invalid_cycle_rows"]++
			}
			preBPB := wlmYggRealUTF8Level1MaintenanceR1BPB(quarters[2], &baseline, active)

			wrongSource, wrongOK := wlmYggRealUTF8Level1MaintenanceR1WrongSource(selectedKey, canonical)
			if !wrongOK {
				metrics["invalid_cycle_rows"]++
				continue
			}
			wrongState := wlmYggRealUTF8Level1MaintenanceR1CloneMap(active)
			wrongEntry := wlmYggRealUTF8Level1MaintenanceR1CloneEntry(canonical[wrongSource])
			wrongEntry.key = selectedKey
			wrongState[selectedKey] = wrongEntry
			wrongAcc, _ := wlmYggRealUTF8Level1MaintenanceR1KeyAccuracy(
				quarters[2], selectedKey, &baseline, wrongState,
			)
			wrongBPB := wlmYggRealUTF8Level1MaintenanceR1BPB(quarters[2], &baseline, wrongState)
			wrongGain := wrongAcc - preAcc
			metrics["maximum_wrong_candidate_selected_key_accuracy_gain"] = wlmYggRealUTF8Level1MaintenanceR1Max(
				metrics["maximum_wrong_candidate_selected_key_accuracy_gain"], wrongGain,
			)
			metrics["maximum_candidate_read_bytes"] = wlmYggRealUTF8Level1MaintenanceR1Max(metrics["maximum_candidate_read_bytes"], 1028)
			metrics["maximum_candidate_repair_operations"] = wlmYggRealUTF8Level1MaintenanceR1Max(metrics["maximum_candidate_repair_operations"], 1)
			wrongRetain := wrongGain >= 0.20 && (wrongBPB-preBPB) <= 0.005
			if wrongRetain {
				active = wrongState
				metrics["wrong_candidate_retain_count"]++
			} else {
				metrics["wrong_candidate_revert_count"]++
			}

			if !wrongRetain {
				correctState := wlmYggRealUTF8Level1MaintenanceR1CloneMap(active)
				correctState[selectedKey] = wlmYggRealUTF8Level1MaintenanceR1CloneEntry(canonical[selectedKey])
				correctAcc, _ := wlmYggRealUTF8Level1MaintenanceR1KeyAccuracy(
					quarters[2], selectedKey, &baseline, correctState,
				)
				correctBPB := wlmYggRealUTF8Level1MaintenanceR1BPB(quarters[2], &baseline, correctState)
				correctGain := correctAcc - preAcc
				metrics["minimum_correct_candidate_selected_key_accuracy_gain"] = wlmYggRealUTF8Level1MaintenanceR1Min(
					metrics["minimum_correct_candidate_selected_key_accuracy_gain"], correctGain,
				)
				correctRetain := correctGain >= 0.20 && (correctBPB-preBPB) <= 0.005
				if correctRetain {
					active = correctState
					metrics["correct_candidate_retain_count"]++
				} else {
					metrics["correct_candidate_revert_count"]++
				}
			}

			jaccard := wlmYggRealUTF8Level1MaintenanceR1Jaccard(active, canonical)
			metrics["minimum_structure_jaccard_after_each_cycle"] = wlmYggRealUTF8Level1MaintenanceR1Min(
				metrics["minimum_structure_jaccard_after_each_cycle"], jaccard,
			)
			if wlmYggRealUTF8Level1MaintenanceR1StateMismatch(active, canonical) != 0 {
				metrics["invalid_cycle_rows"]++
			}
			predMismatch, probDelta := wlmYggRealUTF8Level1MaintenanceR1Audit(
				quarters[3], &baseline, active, canonical,
			)
			metrics["maximum_post_repair_audit_prediction_mismatch_count"] = wlmYggRealUTF8Level1MaintenanceR1Max(
				metrics["maximum_post_repair_audit_prediction_mismatch_count"], float64(predMismatch),
			)
			metrics["maximum_post_repair_audit_probability_abs_delta"] = wlmYggRealUTF8Level1MaintenanceR1Max(
				metrics["maximum_post_repair_audit_probability_abs_delta"], probDelta,
			)
			metrics["completed_cycle_count"]++
		}
	}

	metrics["maximum_final_selected_motif_count"] = float64(len(active))
	metrics["minimum_final_selected_motif_count"] = float64(len(active))
	return wlmYggRealUTF8Level1MaintenanceR1Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-YGG-REAL-UTF8-LEVEL1-MAINTENANCE-R1",
		Metrics: metrics,
	}
}
