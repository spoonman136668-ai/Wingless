package unitary

import "math"

type wlmSiRawRepContextSamplingFalsificationR2Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}
type wlmSiRawRepContextSamplingFalsificationR2Key struct{ op, left, right int }
type wlmSiRawRepContextSamplingFalsificationR2Pair struct{ left, right int }
type wlmSiRawRepContextSamplingFalsificationR2Model struct {
	counts      map[wlmSiRawRepContextSamplingFalsificationR2Key]map[wlmSiRawRepContextSamplingFalsificationR2Pair]uint32
	fallback    wlmSiRawRepContextSamplingFalsificationR2Pair
	hasFallback bool
}

func wlmSiRawRepContextSamplingFalsificationR2Bridge(independent bool) [][]uint8 {
	records := make([][]uint8, 0, 1728)
	for entity := 0; entity < 9; entity++ {
		for family := 0; family < 3; family++ {
			state := uint32(0x243f6a88) ^ uint32((entity+1)*1009)
			if independent {
				state ^= uint32(family+1) * uint32(0x9e3779b9)
			}
			for sample := 0; sample < 64; sample++ {
				state = state*1664525 + 1013904223
				contextIndex := int((state >> 16) % 16)
				noise := uint32(0x85ebca6b) ^ uint32((entity+1)*1009+(family+1)*10007+(sample+1)*65537)
				record := append([]uint8{}, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&noise, 3)...)
				record = append(record, uint8(192+(entity*7+contextIndex*3)%16), uint8(208+(entity*5+contextIndex*5)%16))
				motif := wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity)
				record = append(record, motif[:]...)
				record = append(record, uint8(224+(entity*11+contextIndex*7)%16), uint8(240+(entity*13+contextIndex*11)%16))
				record = append(record, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&noise, 3)...)
				records = append(records, record)
			}
		}
	}
	return records
}
func wlmSiRawRepContextSamplingFalsificationR2Counts(records [][]uint8, reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) []map[wlmLmRawRepContextualAliasInductionR2Context]uint32 {
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	out := make([]map[wlmLmRawRepContextualAliasInductionR2Context]uint32, len(reps))
	for i := range out {
		out[i] = make(map[wlmLmRawRepContextualAliasInductionR2Context]uint32)
	}
	for _, record := range records {
		for i := 2; i+6 <= len(record); i++ {
			key := wlmSiRawRepAliasInvarianceFalsificationR1Key{record[i], record[i+1], record[i+2], record[i+3]}
			id, ok := index[key]
			if ok {
				context := wlmLmRawRepContextualAliasInductionR2Context{record[i-2], record[i-1], record[i+4], record[i+5]}
				out[id][context]++
			}
		}
	}
	return out
}
func wlmSiRawRepContextSamplingFalsificationR2MaxL1(counts []map[wlmLmRawRepContextualAliasInductionR2Context]uint32, reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) int {
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	maximum := 0
	for entity := 0; entity < 9; entity++ {
		for a := 0; a < 3; a++ {
			for b := a + 1; b < 3; b++ {
				left := counts[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(a, entity)]]
				right := counts[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(b, entity)]]
				distance := 0
				seen := make(map[wlmLmRawRepContextualAliasInductionR2Context]bool)
				for k, v := range left {
					d := int(v) - int(right[k])
					if d < 0 {
						d = -d
					}
					distance += d
					seen[k] = true
				}
				for k, v := range right {
					if !seen[k] {
						distance += int(v)
					}
				}
				if distance > maximum {
					maximum = distance
				}
			}
		}
	}
	return maximum
}
func (m *wlmSiRawRepContextSamplingFalsificationR2Model) observe(key wlmSiRawRepContextSamplingFalsificationR2Key, pair wlmSiRawRepContextSamplingFalsificationR2Pair, metrics map[string]float64) {
	if m.counts == nil {
		m.counts = make(map[wlmSiRawRepContextSamplingFalsificationR2Key]map[wlmSiRawRepContextSamplingFalsificationR2Pair]uint32)
	}
	row := m.counts[key]
	if row == nil {
		row = make(map[wlmSiRawRepContextSamplingFalsificationR2Pair]uint32)
		m.counts[key] = row
	}
	if row[pair] == ^uint32(0) {
		metrics["counter_overflow_count"]++
		return
	}
	row[pair]++
	if !m.hasFallback || pair.left < m.fallback.left || (pair.left == m.fallback.left && pair.right < m.fallback.right) {
		m.fallback = pair
		m.hasFallback = true
	}
}
func (m *wlmSiRawRepContextSamplingFalsificationR2Model) predict(key wlmSiRawRepContextSamplingFalsificationR2Key) (wlmSiRawRepContextSamplingFalsificationR2Pair, bool) {
	row := m.counts[key]
	if len(row) == 0 {
		return m.fallback, false
	}
	first := true
	best := m.fallback
	var n uint32
	for pair, count := range row {
		if first || count > n || (count == n && (pair.left < best.left || (pair.left == best.left && pair.right < best.right))) {
			first = false
			best = pair
			n = count
		}
	}
	return best, true
}
func wlmSiRawRepContextSamplingFalsificationR2Build(records [][]uint8, reps []wlmSiRawRepAliasInvarianceFalsificationR1Key, classes wlmLmRawRepContextualAliasInductionR2Classes, metrics map[string]float64) wlmSiRawRepContextSamplingFalsificationR2Model {
	var model wlmSiRawRepContextSamplingFalsificationR2Model
	for _, record := range records {
		ids := wlmSiRawRepAliasInvarianceFalsificationR1Decode(record, reps)
		if len(ids) != 5 {
			metrics["invalid_row_count"]++
			continue
		}
		c := make([]int, 5)
		valid := true
		for i, id := range ids {
			if id < 0 || id >= len(classes.surfaceToClass) {
				valid = false
				break
			}
			c[i] = classes.surfaceToClass[id]
		}
		if !valid {
			metrics["invalid_row_count"]++
			continue
		}
		model.observe(wlmSiRawRepContextSamplingFalsificationR2Key{c[2], c[0], c[1]}, wlmSiRawRepContextSamplingFalsificationR2Pair{c[3], c[4]}, metrics)
	}
	return model
}
func wlmSiRawRepContextSamplingFalsificationR2Evaluate(reps []wlmSiRawRepAliasInvarianceFalsificationR1Key, classes wlmLmRawRepContextualAliasInductionR2Classes, model wlmSiRawRepContextSamplingFalsificationR2Model, metrics map[string]float64) (int, int) {
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	hits, missing := 0, 0
	for op := 0; op < 5; op++ {
		for left := 0; left < 4; left++ {
			for right := 0; right < 4; right++ {
				raw := wlmLmRawRepZeroshotAliasFamilyR1EvalRecord(left, right, op, 5000+op*20+left*4+right)
				ids := wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw, reps)
				if len(ids) != 3 {
					metrics["invalid_row_count"]++
					continue
				}
				pair, found := model.predict(wlmSiRawRepContextSamplingFalsificationR2Key{classes.surfaceToClass[ids[2]], classes.surfaceToClass[ids[0]], classes.surfaceToClass[ids[1]]})
				if !found {
					missing++
				}
				tl, tr := wlmSiRawRepAliasInvarianceFalsificationR1Apply(left, right, op)
				el := classes.surfaceToClass[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(2, tl)]]
				er := classes.surfaceToClass[index[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(2, tr)]]
				if pair.left == el && pair.right == er {
					hits++
				}
			}
		}
	}
	return hits, missing
}
func wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(classes wlmLmRawRepContextualAliasInductionR2Classes, index map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int) (int, int, int) {
	three := 0
	for _, members := range classes.members {
		if len(members) == 3 {
			three++
		}
	}
	return len(classes.members), three, wlmLmRawRepZeroshotAliasFamilyR1CorrectTriplets(index, classes)
}

func RunWlmSiRawRepContextSamplingFalsificationR2() interface{} {
	m := map[string]float64{"matched_bridge_record_count": 0, "independent_bridge_record_count": 0, "selected_surface_representation_count": 0, "selected_true_surface_motif_match_count": 0, "matched_max_context_count_l1_distance": 0, "independent_max_context_count_l1_distance": 0, "matched_alias_class_count": 0, "matched_three_member_class_count": 0, "matched_correct_alias_triplet_count": 0, "matched_unseen_family_accuracy": 0, "independent_alias_class_count": 0, "independent_three_member_class_count": 0, "independent_correct_alias_triplet_count": 0, "independent_missing_transition_count": 0, "independent_unseen_family_accuracy": 0, "exact_surface_unseen_family_accuracy": 0, "task_training_record_count": 0, "family2_task_training_record_count": 0, "unseen_family_evaluation_count": 0, "tokenizer_use_count": 0, "external_model_call_count": 0, "capacity_growth_event_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0}
	matched, independent := wlmSiRawRepContextSamplingFalsificationR2Bridge(false), wlmSiRawRepContextSamplingFalsificationR2Bridge(true)
	m["matched_bridge_record_count"] = float64(len(matched))
	m["independent_bridge_record_count"] = float64(len(independent))
	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(matched)
	m["selected_surface_representation_count"] = float64(len(reps))
	m["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	m["matched_max_context_count_l1_distance"] = float64(wlmSiRawRepContextSamplingFalsificationR2MaxL1(wlmSiRawRepContextSamplingFalsificationR2Counts(matched, reps), reps))
	m["independent_max_context_count_l1_distance"] = float64(wlmSiRawRepContextSamplingFalsificationR2MaxL1(wlmSiRawRepContextSamplingFalsificationR2Counts(independent, reps), reps))
	matchedClasses := wlmLmRawRepContextualAliasInductionR2Induce(reps, matched)
	independentClasses := wlmLmRawRepContextualAliasInductionR2Induce(reps, independent)
	if matchedClasses.decodeFailures != 0 || independentClasses.decodeFailures != 0 {
		m["invalid_row_count"] += float64(matchedClasses.decodeFailures + independentClasses.decodeFailures)
	}
	index := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)
	ma, mt, mc := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(matchedClasses, index)
	ia, it, ic := wlmSiRawRepContextSamplingFalsificationR2ClassMetrics(independentClasses, index)
	m["matched_alias_class_count"], m["matched_three_member_class_count"], m["matched_correct_alias_triplet_count"] = float64(ma), float64(mt), float64(mc)
	m["independent_alias_class_count"], m["independent_three_member_class_count"], m["independent_correct_alias_triplet_count"] = float64(ia), float64(it), float64(ic)
	tasks := wlmLmRawRepContextualAliasInductionR2TaskRecords()
	m["task_training_record_count"] = float64(len(tasks))
	matchedModel := wlmSiRawRepContextSamplingFalsificationR2Build(tasks, reps, matchedClasses, m)
	independentModel := wlmSiRawRepContextSamplingFalsificationR2Build(tasks, reps, independentClasses, m)
	matchedHits, _ := wlmSiRawRepContextSamplingFalsificationR2Evaluate(reps, matchedClasses, matchedModel, m)
	independentHits, missing := wlmSiRawRepContextSamplingFalsificationR2Evaluate(reps, independentClasses, independentModel, m)
	m["unseen_family_evaluation_count"] = 80
	m["matched_unseen_family_accuracy"] = float64(matchedHits) / 80
	m["independent_unseen_family_accuracy"] = float64(independentHits) / 80
	m["independent_missing_transition_count"] = float64(missing)
	surfaceModel, decodeFails := wlmLmRawRepZeroshotAliasFamilyR1BuildSurfaceModel(tasks, reps, m)
	if decodeFails != 0 {
		m["invalid_row_count"] += float64(decodeFails)
	}
	exactHits := 0
	for op := 0; op < 5; op++ {
		for left := 0; left < 4; left++ {
			for right := 0; right < 4; right++ {
				ok, _ := wlmLmRawRepZeroshotAliasFamilyR1EvaluateSurface(left, right, op, 6000+op*20+left*4+right, reps, surfaceModel, m)
				if ok {
					exactHits++
				}
			}
		}
	}
	m["exact_surface_unseen_family_accuracy"] = float64(exactHits) / 80
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmSiRawRepContextSamplingFalsificationR2Result{"wingless.research-scientific-result.v1", "WLM-SI-RAW-REP-CONTEXT-SAMPLING-FALSIFICATION-R2", m}
}
