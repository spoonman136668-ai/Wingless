package unitary

import "sort"

type wlmLmByteMotifDiscoveryR3BalancedWindow interface {
	Push(uint8)
	Ready() bool
	Snapshot() [4]uint8
}

type wlmLmByteMotifDiscoveryR3BalancedSliceWindow struct {
	values []uint8
}

func (w *wlmLmByteMotifDiscoveryR3BalancedSliceWindow) Push(v uint8) {
	w.values = append(w.values, v)
	if len(w.values) > 4 {
		w.values = w.values[len(w.values)-4:]
	}
}

func (w *wlmLmByteMotifDiscoveryR3BalancedSliceWindow) Ready() bool {
	return len(w.values) == 4
}

func (w *wlmLmByteMotifDiscoveryR3BalancedSliceWindow) Snapshot() [4]uint8 {
	var out [4]uint8
	if len(w.values) == 4 {
		copy(out[:], w.values)
	}
	return out
}

type wlmLmByteMotifDiscoveryR3BalancedRingWindow struct {
	values [4]uint8
	head   int
	length int
}

func (w *wlmLmByteMotifDiscoveryR3BalancedRingWindow) Push(v uint8) {
	if w.length < 4 {
		w.values[(w.head+w.length)%4] = v
		w.length++
		return
	}
	w.values[w.head] = v
	w.head = (w.head + 1) % 4
}

func (w *wlmLmByteMotifDiscoveryR3BalancedRingWindow) Ready() bool {
	return w.length == 4
}

func (w *wlmLmByteMotifDiscoveryR3BalancedRingWindow) Snapshot() [4]uint8 {
	var out [4]uint8
	if w.length == 4 {
		for i := 0; i < 4; i++ {
			out[i] = w.values[(w.head+i)%4]
		}
	}
	return out
}

func wlmLmByteMotifDiscoveryR3BalancedNewWindow(substrate string) wlmLmByteMotifDiscoveryR3BalancedWindow {
	if substrate == "slice-backed-four-byte-window" {
		return &wlmLmByteMotifDiscoveryR3BalancedSliceWindow{}
	}
	return &wlmLmByteMotifDiscoveryR3BalancedRingWindow{}
}

type wlmLmByteMotifDiscoveryR3BalancedNextStats struct {
	total  uint32
	counts map[uint8]uint32
}

func (s *wlmLmByteMotifDiscoveryR3BalancedNextStats) observe(next uint8, metrics map[string]float64) {
	if s.counts == nil {
		s.counts = make(map[uint8]uint32)
	}
	if s.total == ^uint32(0) || s.counts[next] == ^uint32(0) {
		metrics["counter_overflow_rows"]++
		return
	}
	s.total++
	s.counts[next]++
}

func (s *wlmLmByteMotifDiscoveryR3BalancedNextStats) best() (uint8, uint32) {
	var bestByte uint8
	var bestCount uint32
	first := true
	for b, count := range s.counts {
		if first || count > bestCount || (count == bestCount && b < bestByte) {
			bestByte = b
			bestCount = count
			first = false
		}
	}
	return bestByte, bestCount
}

type wlmLmByteMotifDiscoveryR3BalancedLearner struct {
	candidates map[[4]uint8]*wlmLmByteMotifDiscoveryR3BalancedNextStats
	baseline   map[uint16]*wlmLmByteMotifDiscoveryR3BalancedNextStats
	selected   []wlmLmByteMotifDiscoveryR3BalancedSelected
	selectedMap map[[4]uint8]uint8
	shuffledMap map[[4]uint8]uint8
}

type wlmLmByteMotifDiscoveryR3BalancedSelected struct {
	key         [4]uint8
	successor   uint8
	total       uint32
	bestCount   uint32
	consistency float64
}

type wlmLmByteMotifDiscoveryR3BalancedCorpus struct {
	bytes           []uint8
	targetPositions []int
}

type wlmLmByteMotifDiscoveryR3BalancedResult struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmByteMotifDiscoveryR3BalancedMotif(id int) [4]uint8 {
	return [4]uint8{uint8(128 + id), uint8(64 + ((7 * id) % 32)), 170, 85}
}

func wlmLmByteMotifDiscoveryR3BalancedSuccessor(id int) uint8 {
	return uint8(16 + 13*id)
}

func wlmLmByteMotifDiscoveryR3BalancedStep(state uint32) uint32 {
	return state*1664525 + 1013904223
}

func wlmLmByteMotifDiscoveryR3BalancedCorpusFor(seed uint32, records int, evaluation bool) wlmLmByteMotifDiscoveryR3BalancedCorpus {
	state := seed ^ 0x243f6a88
	if evaluation {
		state = seed ^ 0xb7e15162
	}
	out := wlmLmByteMotifDiscoveryR3BalancedCorpus{
		bytes:           make([]uint8, 0, records*8),
		targetPositions: make([]int, 0, records),
	}
	seedOffset := int(seed % 8)
	for record := 0; record < records; record++ {
		state = wlmLmByteMotifDiscoveryR3BalancedStep(state) // preserved R2 pre-motif LCG advance
		id := (5*record + seedOffset) % 8
		motif := wlmLmByteMotifDiscoveryR3BalancedMotif(id)
		out.bytes = append(out.bytes, motif[:]...)
		out.targetPositions = append(out.targetPositions, len(out.bytes))
		out.bytes = append(out.bytes, wlmLmByteMotifDiscoveryR3BalancedSuccessor(id))
		for filler := 0; filler < 3; filler++ {
			state = wlmLmByteMotifDiscoveryR3BalancedStep(state)
			out.bytes = append(out.bytes, uint8(192+(state%64)))
		}
	}
	return out
}

func wlmLmByteMotifDiscoveryR3BalancedNewLearner() wlmLmByteMotifDiscoveryR3BalancedLearner {
	return wlmLmByteMotifDiscoveryR3BalancedLearner{
		candidates: make(map[[4]uint8]*wlmLmByteMotifDiscoveryR3BalancedNextStats),
		baseline:   make(map[uint16]*wlmLmByteMotifDiscoveryR3BalancedNextStats),
		selectedMap: make(map[[4]uint8]uint8),
		shuffledMap: make(map[[4]uint8]uint8),
	}
}

func wlmLmByteMotifDiscoveryR3BalancedBaselineKey(a, b uint8) uint16 {
	return uint16(a)<<8 | uint16(b)
}

func (l *wlmLmByteMotifDiscoveryR3BalancedLearner) train(stream []uint8, substrate string, metrics map[string]float64) {
	window := wlmLmByteMotifDiscoveryR3BalancedNewWindow(substrate)
	for i, b := range stream {
		window.Push(b)
		if i >= 1 && i+1 < len(stream) {
			key := wlmLmByteMotifDiscoveryR3BalancedBaselineKey(stream[i-1], stream[i])
			stats := l.baseline[key]
			if stats == nil {
				stats = &wlmLmByteMotifDiscoveryR3BalancedNextStats{}
				l.baseline[key] = stats
			}
			stats.observe(stream[i+1], metrics)
		}
		if window.Ready() && i+1 < len(stream) {
			key := window.Snapshot()
			stats := l.candidates[key]
			if stats == nil {
				stats = &wlmLmByteMotifDiscoveryR3BalancedNextStats{}
				l.candidates[key] = stats
			}
			stats.observe(stream[i+1], metrics)
		}
	}
}

func wlmLmByteMotifDiscoveryR3BalancedKeyLess(a, b [4]uint8) bool {
	for i := 0; i < 4; i++ {
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return false
}

func (l *wlmLmByteMotifDiscoveryR3BalancedLearner) selectMotifs(metrics map[string]float64) {
	candidates := make([]wlmLmByteMotifDiscoveryR3BalancedSelected, 0)
	for key, stats := range l.candidates {
		if stats.total < 64 {
			continue
		}
		bestByte, bestCount := stats.best()
		if stats.total == 0 {
			metrics["invalid_candidate_rows"]++
			continue
		}
		consistency := float64(bestCount) / float64(stats.total)
		if consistency < 0.90 {
			continue
		}
		candidates = append(candidates, wlmLmByteMotifDiscoveryR3BalancedSelected{
			key: key, successor: bestByte, total: stats.total,
			bestCount: bestCount, consistency: consistency,
		})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].bestCount != candidates[j].bestCount {
			return candidates[i].bestCount > candidates[j].bestCount
		}
		if candidates[i].consistency != candidates[j].consistency {
			return candidates[i].consistency > candidates[j].consistency
		}
		return wlmLmByteMotifDiscoveryR3BalancedKeyLess(candidates[i].key, candidates[j].key)
	})
	if len(candidates) > 8 {
		candidates = candidates[:8]
	}
	l.selected = append([]wlmLmByteMotifDiscoveryR3BalancedSelected(nil), candidates...)
	for _, selected := range l.selected {
		l.selectedMap[selected.key] = selected.successor
		shuffled := [4]uint8{selected.key[2], selected.key[0], selected.key[3], selected.key[1]}
		if _, exists := l.shuffledMap[shuffled]; exists {
			metrics["duplicate_shuffled_control_key_count"]++
		}
		l.shuffledMap[shuffled] = selected.successor
	}
}

func (l *wlmLmByteMotifDiscoveryR3BalancedLearner) baselinePredict(a, b uint8) uint8 {
	stats := l.baseline[wlmLmByteMotifDiscoveryR3BalancedBaselineKey(a, b)]
	if stats == nil {
		return 0
	}
	best, _ := stats.best()
	return best
}

func (l *wlmLmByteMotifDiscoveryR3BalancedLearner) predict(window [4]uint8, mode string) uint8 {
	if mode == "motif" {
		if next, ok := l.selectedMap[window]; ok {
			return next
		}
	}
	if mode == "shuffled" {
		if next, ok := l.shuffledMap[window]; ok {
			return next
		}
	}
	return l.baselinePredict(window[2], window[3])
}

func wlmLmByteMotifDiscoveryR3BalancedSelectedMismatch(a, b *wlmLmByteMotifDiscoveryR3BalancedLearner) int {
	if len(a.selected) != len(b.selected) {
		if len(a.selected) > len(b.selected) {
			return len(a.selected) - len(b.selected)
		}
		return len(b.selected) - len(a.selected)
	}
	mismatch := 0
	for i := range a.selected {
		if a.selected[i].key != b.selected[i].key || a.selected[i].successor != b.selected[i].successor {
			mismatch++
		}
	}
	return mismatch
}

func wlmLmByteMotifDiscoveryR3BalancedTrueRecall(l *wlmLmByteMotifDiscoveryR3BalancedLearner) (float64, int) {
	trueSet := make(map[[4]uint8]bool)
	for id := 0; id < 8; id++ {
		trueSet[wlmLmByteMotifDiscoveryR3BalancedMotif(id)] = true
	}
	trueCount := 0
	falseCount := 0
	for _, selected := range l.selected {
		if trueSet[selected.key] {
			trueCount++
		} else {
			falseCount++
		}
	}
	return float64(trueCount) / 8.0, falseCount
}

func wlmLmByteMotifDiscoveryR3BalancedEventReduction(stream []uint8, l *wlmLmByteMotifDiscoveryR3BalancedLearner) float64 {
	events := 0
	for i := 0; i < len(stream); {
		if i+4 <= len(stream) {
			key := [4]uint8{stream[i], stream[i+1], stream[i+2], stream[i+3]}
			if _, ok := l.selectedMap[key]; ok {
				events++
				i += 4
				continue
			}
		}
		events++
		i++
	}
	return float64(len(stream)-events) / float64(len(stream))
}

func wlmLmByteMotifDiscoveryR3BalancedEvaluate(
	corpus wlmLmByteMotifDiscoveryR3BalancedCorpus,
	sliceLearner, ringLearner *wlmLmByteMotifDiscoveryR3BalancedLearner,
	metrics map[string]float64,
) (float64, float64, float64, float64, float64, float64, float64, float64) {
	sliceMotif, ringMotif := 0, 0
	sliceAblation, ringAblation := 0, 0
	sliceShuffled, ringShuffled := 0, 0
	total := 0
	for _, pos := range corpus.targetPositions {
		if pos < 4 || pos >= len(corpus.bytes) {
			metrics["invalid_candidate_rows"]++
			continue
		}
		key := [4]uint8{corpus.bytes[pos-4], corpus.bytes[pos-3], corpus.bytes[pos-2], corpus.bytes[pos-1]}
		target := corpus.bytes[pos]
		sm := sliceLearner.predict(key, "motif")
		rm := ringLearner.predict(key, "motif")
		if sm != rm {
			metrics["substrate_prediction_mismatch_count"]++
		}
		if sm == target { sliceMotif++ }
		if rm == target { ringMotif++ }
		if sliceLearner.predict(key, "ablation") == target { sliceAblation++ }
		if ringLearner.predict(key, "ablation") == target { ringAblation++ }
		if sliceLearner.predict(key, "shuffled") == target { sliceShuffled++ }
		if ringLearner.predict(key, "shuffled") == target { ringShuffled++ }
		total++
	}
	sliceReduction := wlmLmByteMotifDiscoveryR3BalancedEventReduction(corpus.bytes, sliceLearner)
	ringReduction := wlmLmByteMotifDiscoveryR3BalancedEventReduction(corpus.bytes, ringLearner)
	if sliceReduction != ringReduction {
		metrics["substrate_event_accounting_mismatch_count"]++
	}
	den := float64(total)
	return float64(sliceMotif)/den, float64(ringMotif)/den,
		float64(sliceAblation)/den, float64(ringAblation)/den,
		float64(sliceShuffled)/den, float64(ringShuffled)/den,
		sliceReduction, ringReduction
}

func wlmLmByteMotifDiscoveryR3BalancedMin(a, b float64) float64 {
	if b < a { return b }
	return a
}

func wlmLmByteMotifDiscoveryR3BalancedMax(a, b float64) float64 {
	if b > a { return b }
	return a
}

// RunWlmLmByteMotifDiscoveryR3Balanced evaluates the frozen WBG-2 motif-discovery matrix.
func RunWlmLmByteMotifDiscoveryR3Balanced() interface{} {
	seeds := [...]uint32{11801, 11821, 11827, 11863}
	metrics := map[string]float64{
		"valid_seed_count": 0,
		"completed_substrate_runs": 0,
		"selected_motif_count": 8,
		"minimum_selected_true_motif_recall": 1,
		"maximum_false_selected_motif_count": 0,
		"minimum_motif_enabled_successor_accuracy": 1,
		"maximum_byte_only_ablation_successor_accuracy": 0,
		"maximum_shuffled_motif_successor_accuracy": 0,
		"minimum_causal_accuracy_gain_over_ablation": 1,
		"minimum_effective_event_reduction_fraction": 1,
		"substrate_selected_motif_mismatch_count": 0,
		"substrate_prediction_mismatch_count": 0,
		"substrate_event_accounting_mismatch_count": 0,
		"duplicate_shuffled_control_key_count": 0,
		"invalid_candidate_rows": 0,
		"counter_overflow_rows": 0,
	}
	for _, seed := range seeds {
		train := wlmLmByteMotifDiscoveryR3BalancedCorpusFor(seed, 4096, false)
		eval := wlmLmByteMotifDiscoveryR3BalancedCorpusFor(seed, 2048, true)
		sliceLearner := wlmLmByteMotifDiscoveryR3BalancedNewLearner()
		ringLearner := wlmLmByteMotifDiscoveryR3BalancedNewLearner()
		sliceLearner.train(train.bytes, "slice-backed-four-byte-window", metrics)
		ringLearner.train(train.bytes, "fixed-capacity-ring-four-byte-window", metrics)
		sliceLearner.selectMotifs(metrics)
		ringLearner.selectMotifs(metrics)
		metrics["completed_substrate_runs"] += 2
		if float64(len(sliceLearner.selected)) < metrics["selected_motif_count"] {
			metrics["selected_motif_count"] = float64(len(sliceLearner.selected))
		}
		if float64(len(ringLearner.selected)) < metrics["selected_motif_count"] {
			metrics["selected_motif_count"] = float64(len(ringLearner.selected))
		}
		metrics["substrate_selected_motif_mismatch_count"] += float64(
			wlmLmByteMotifDiscoveryR3BalancedSelectedMismatch(&sliceLearner, &ringLearner),
		)
		sliceRecall, sliceFalse := wlmLmByteMotifDiscoveryR3BalancedTrueRecall(&sliceLearner)
		ringRecall, ringFalse := wlmLmByteMotifDiscoveryR3BalancedTrueRecall(&ringLearner)
		metrics["minimum_selected_true_motif_recall"] = wlmLmByteMotifDiscoveryR3BalancedMin(metrics["minimum_selected_true_motif_recall"], sliceRecall)
		metrics["minimum_selected_true_motif_recall"] = wlmLmByteMotifDiscoveryR3BalancedMin(metrics["minimum_selected_true_motif_recall"], ringRecall)
		metrics["maximum_false_selected_motif_count"] = wlmLmByteMotifDiscoveryR3BalancedMax(metrics["maximum_false_selected_motif_count"], float64(sliceFalse))
		metrics["maximum_false_selected_motif_count"] = wlmLmByteMotifDiscoveryR3BalancedMax(metrics["maximum_false_selected_motif_count"], float64(ringFalse))
		sm, rm, sa, ra, ss, rs, ser, rer := wlmLmByteMotifDiscoveryR3BalancedEvaluate(eval, &sliceLearner, &ringLearner, metrics)
		metrics["minimum_motif_enabled_successor_accuracy"] = wlmLmByteMotifDiscoveryR3BalancedMin(metrics["minimum_motif_enabled_successor_accuracy"], sm)
		metrics["minimum_motif_enabled_successor_accuracy"] = wlmLmByteMotifDiscoveryR3BalancedMin(metrics["minimum_motif_enabled_successor_accuracy"], rm)
		metrics["maximum_byte_only_ablation_successor_accuracy"] = wlmLmByteMotifDiscoveryR3BalancedMax(metrics["maximum_byte_only_ablation_successor_accuracy"], sa)
		metrics["maximum_byte_only_ablation_successor_accuracy"] = wlmLmByteMotifDiscoveryR3BalancedMax(metrics["maximum_byte_only_ablation_successor_accuracy"], ra)
		metrics["maximum_shuffled_motif_successor_accuracy"] = wlmLmByteMotifDiscoveryR3BalancedMax(metrics["maximum_shuffled_motif_successor_accuracy"], ss)
		metrics["maximum_shuffled_motif_successor_accuracy"] = wlmLmByteMotifDiscoveryR3BalancedMax(metrics["maximum_shuffled_motif_successor_accuracy"], rs)
		metrics["minimum_causal_accuracy_gain_over_ablation"] = wlmLmByteMotifDiscoveryR3BalancedMin(metrics["minimum_causal_accuracy_gain_over_ablation"], sm-sa)
		metrics["minimum_causal_accuracy_gain_over_ablation"] = wlmLmByteMotifDiscoveryR3BalancedMin(metrics["minimum_causal_accuracy_gain_over_ablation"], rm-ra)
		metrics["minimum_effective_event_reduction_fraction"] = wlmLmByteMotifDiscoveryR3BalancedMin(metrics["minimum_effective_event_reduction_fraction"], ser)
		metrics["minimum_effective_event_reduction_fraction"] = wlmLmByteMotifDiscoveryR3BalancedMin(metrics["minimum_effective_event_reduction_fraction"], rer)
		metrics["valid_seed_count"]++
	}
	return wlmLmByteMotifDiscoveryR3BalancedResult{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-BYTE-MOTIF-DISCOVERY-R3-BALANCED",
		Metrics: metrics,
	}
}
