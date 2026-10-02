package unitary

type wlmLmByteMotifAttributionR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmByteMotifAttributionR1MotifID(key [4]uint8) int {
	for id := 0; id < 8; id++ {
		if key == wlmLmByteMotifDiscoveryR2Motif(id) {
			return id
		}
	}
	return -1
}

func wlmLmByteMotifAttributionR1SelectedCounts(
	learner *wlmLmByteMotifDiscoveryR2Learner,
	present map[int]bool,
) (selectedTrue int, selectedFalse int, presentRecovered int) {
	for _, selected := range learner.selected {
		id := wlmLmByteMotifAttributionR1MotifID(selected.key)
		if id >= 0 {
			selectedTrue++
			if present[id] {
				presentRecovered++
			}
		} else {
			selectedFalse++
		}
	}
	return
}

func wlmLmByteMotifAttributionR1Min(current, candidate float64) float64 {
	if candidate < current {
		return candidate
	}
	return current
}

func wlmLmByteMotifAttributionR1Max(current, candidate float64) float64 {
	if candidate > current {
		return candidate
	}
	return current
}

// RunWlmLmByteMotifAttributionR1 attributes the frozen WBG-2 negative without changing its corpus or selector.
func RunWlmLmByteMotifAttributionR1() interface{} {
	seeds := [...]uint32{11801, 11821, 11827, 11863}
	metrics := map[string]float64{
		"valid_seed_count": 0,
		"minimum_unique_hidden_motif_ids_per_seed": 8,
		"maximum_unique_hidden_motif_ids_per_seed": 0,
		"union_hidden_motif_id_count_across_seeds": 0,
		"minimum_present_motif_occurrence_count": 4096,
		"maximum_present_motif_occurrence_count": 0,
		"minimum_absent_motif_count_per_seed": 8,
		"maximum_absent_motif_count_per_seed": 0,
		"minimum_present_true_motif_selector_recall": 1,
		"maximum_present_true_motif_selector_recall": 0,
		"minimum_selected_true_motif_count_per_seed": 8,
		"maximum_selected_true_motif_count_per_seed": 0,
		"minimum_selected_false_motif_count_per_seed": 8,
		"maximum_selected_false_motif_count_per_seed": 0,
		"substrate_selected_motif_mismatch_count": 0,
		"counter_overflow_rows": 0,
	}
	union := make(map[int]bool)

	for _, seed := range seeds {
		corpus := wlmLmByteMotifDiscoveryR2CorpusFor(seed, 4096, false)
		counts := [8]int{}
		for _, pos := range corpus.targetPositions {
			if pos < 4 || pos >= len(corpus.bytes) {
				continue
			}
			key := [4]uint8{corpus.bytes[pos-4], corpus.bytes[pos-3], corpus.bytes[pos-2], corpus.bytes[pos-1]}
			id := wlmLmByteMotifAttributionR1MotifID(key)
			if id >= 0 {
				counts[id]++
				union[id] = true
			}
		}
		present := make(map[int]bool)
		presentCount := 0
		absentCount := 0
		for id, count := range counts {
			if count > 0 {
				present[id] = true
				presentCount++
				metrics["minimum_present_motif_occurrence_count"] = wlmLmByteMotifAttributionR1Min(metrics["minimum_present_motif_occurrence_count"], float64(count))
				metrics["maximum_present_motif_occurrence_count"] = wlmLmByteMotifAttributionR1Max(metrics["maximum_present_motif_occurrence_count"], float64(count))
			} else {
				absentCount++
			}
		}
		metrics["minimum_unique_hidden_motif_ids_per_seed"] = wlmLmByteMotifAttributionR1Min(metrics["minimum_unique_hidden_motif_ids_per_seed"], float64(presentCount))
		metrics["maximum_unique_hidden_motif_ids_per_seed"] = wlmLmByteMotifAttributionR1Max(metrics["maximum_unique_hidden_motif_ids_per_seed"], float64(presentCount))
		metrics["minimum_absent_motif_count_per_seed"] = wlmLmByteMotifAttributionR1Min(metrics["minimum_absent_motif_count_per_seed"], float64(absentCount))
		metrics["maximum_absent_motif_count_per_seed"] = wlmLmByteMotifAttributionR1Max(metrics["maximum_absent_motif_count_per_seed"], float64(absentCount))

		local := map[string]float64{"counter_overflow_rows": 0}
		sliceLearner := wlmLmByteMotifDiscoveryR2NewLearner()
		ringLearner := wlmLmByteMotifDiscoveryR2NewLearner()
		sliceLearner.train(corpus.bytes, "slice-backed-four-byte-window", local)
		ringLearner.train(corpus.bytes, "fixed-capacity-ring-four-byte-window", local)
		sliceLearner.selectMotifs(local)
		ringLearner.selectMotifs(local)
		metrics["counter_overflow_rows"] += local["counter_overflow_rows"]
		metrics["substrate_selected_motif_mismatch_count"] += float64(wlmLmByteMotifDiscoveryR2SelectedMismatch(&sliceLearner, &ringLearner))

		for _, learner := range []*wlmLmByteMotifDiscoveryR2Learner{&sliceLearner, &ringLearner} {
			selectedTrue, selectedFalse, presentRecovered := wlmLmByteMotifAttributionR1SelectedCounts(learner, present)
			recall := 0.0
			if presentCount > 0 {
				recall = float64(presentRecovered) / float64(presentCount)
			}
			metrics["minimum_present_true_motif_selector_recall"] = wlmLmByteMotifAttributionR1Min(metrics["minimum_present_true_motif_selector_recall"], recall)
			metrics["maximum_present_true_motif_selector_recall"] = wlmLmByteMotifAttributionR1Max(metrics["maximum_present_true_motif_selector_recall"], recall)
			metrics["minimum_selected_true_motif_count_per_seed"] = wlmLmByteMotifAttributionR1Min(metrics["minimum_selected_true_motif_count_per_seed"], float64(selectedTrue))
			metrics["maximum_selected_true_motif_count_per_seed"] = wlmLmByteMotifAttributionR1Max(metrics["maximum_selected_true_motif_count_per_seed"], float64(selectedTrue))
			metrics["minimum_selected_false_motif_count_per_seed"] = wlmLmByteMotifAttributionR1Min(metrics["minimum_selected_false_motif_count_per_seed"], float64(selectedFalse))
			metrics["maximum_selected_false_motif_count_per_seed"] = wlmLmByteMotifAttributionR1Max(metrics["maximum_selected_false_motif_count_per_seed"], float64(selectedFalse))
		}
		metrics["valid_seed_count"]++
	}

	metrics["union_hidden_motif_id_count_across_seeds"] = float64(len(union))
	return wlmLmByteMotifAttributionR1Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-BYTE-MOTIF-ATTRIBUTION-R1",
		Metrics: metrics,
	}
}
