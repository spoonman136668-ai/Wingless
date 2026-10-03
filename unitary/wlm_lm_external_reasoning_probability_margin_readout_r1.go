package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalReasoningProbabilityMarginReadoutR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalReasoningProbabilityMarginReadoutR1(code, structured, prose []byte) interface{} {
	sources := []wlmLmExternalMotifRelationDistributionalR1Source{
		{domain: "code", data: code, sha256: "7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6", bytes: 41453},
		{domain: "structured", data: structured, sha256: "4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192", bytes: 14365},
		{domain: "technical_prose", data: prose, sha256: "8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3", bytes: 1454},
	}
	m := map[string]float64{
		"source_identity_mismatch_count": 0, "source_count": 3, "total_source_bytes": 0,
		"selected_motif_count": 0, "refined_reasoning_class_count": 8, "motif_class_assignment_count": 512,
		"minimum_refined_class_size": 512, "maximum_refined_class_size": 0,
		"depth3_query_count": 0, "training_margin_sample_count": 0, "training_margin_threshold": 0,
		"gate_activation_count": 0, "gate_activation_fraction": 0, "gate_exact_hit_count": 0, "gate_exact_precision": 0,
		"depth3_exact_repair_gain_query_count": 0, "hidden_gain_count": 0,
		"same_class_hidden_gate_activation_count": 0, "same_class_hidden_gate_activation_fraction": 0,
		"same_class_hidden_gain_count": 0, "cross_class_hidden_gain_count": 0,
		"same_class_truth_motif_share_increase_count": 0, "same_class_truth_motif_share_increase_fraction": 0,
		"hidden_gain_source_code_count": 0, "hidden_gain_source_structured_count": 0, "hidden_gain_source_technical_prose_count": 0,
		"same_class_source_code_count": 0, "same_class_source_structured_count": 0, "same_class_source_technical_prose_count": 0,
		"same_class_source_code_rate": 0, "same_class_source_structured_rate": 0, "same_class_source_technical_prose_rate": 0,
		"source_rate_denominator_nonzero_count": 0, "code_structured_same_class_rate_gap": 0,
		"minimum_baseline_truth_motif_share": 1, "minimum_repaired_truth_motif_share": 1,
		"maximum_truth_motif_share_increase": 0, "minimum_truth_motif_share_increase": 1,
		"attribution_accounting_error_count": 0, "maximum_probability_mass_error": 0, "minimum_repaired_retained_mass": 1,
		"heldout_selection_use_count": 0, "relation_capacity_growth_event_count": 0, "adaptive_readout_growth_event_count": 0,
		"tokenizer_use_count": 0, "external_model_call_count": 0, "invalid_row_count": 0, "counter_overflow_count": 0,
	}

	train := make([][]byte, 0, 3)
	evals := make([][]byte, 0, 3)
	for _, s := range sources {
		m["total_source_bytes"] += float64(len(s.data))
		if len(s.data) != s.bytes || wlmLmExternalRawRepPredR1SHA256(s.data) != s.sha256 {
			m["source_identity_mismatch_count"]++
		}
		q := len(s.data) * 3 / 5
		if q < 5 || len(s.data)-q < 5 {
			m["invalid_row_count"]++
		}
		train = append(train, s.data[:q])
		evals = append(evals, s.data[q:])
	}

	_, selected := wlmLmRawRepPredFreshHoldoutR1TrainModel(train, m)
	m["selected_motif_count"] = float64(len(selected))
	keys := wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
	index := make(map[[4]uint8]int, len(keys))
	for i, k := range keys {
		index[k] = i
	}

	streams := make([][]int, len(train))
	var counts [512]uint64
	var rel [512][512]uint32
	var totals [512]uint64
	var uni [512]uint32
	var unitotal uint64
	pairs := make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
	for si, data := range train {
		st := wlmLmExternalMotifRelationHoldoutR1Decode(data, index)
		streams[si] = st
		for _, x := range st {
			counts[x]++
		}
		for j := 0; j+1 < len(st); j++ {
			a, b := st[j], st[j+1]
			if rel[a][b] == ^uint32(0) || uni[b] == ^uint32(0) {
				m["counter_overflow_count"]++
				continue
			}
			rel[a][b]++
			totals[a]++
			uni[b]++
			unitotal++
		}
		for j := 1; j < len(st); j++ {
			p := wlmLmExternalContextStateHoldoutR1Pair{prev: st[j-1], curr: st[j]}
			if pairs[p] == ^uint32(0) {
				m["counter_overflow_count"]++
				continue
			}
			pairs[p]++
		}
	}
	if unitotal == 0 {
		m["invalid_row_count"]++
	}

	ids := make([]int, 512)
	for i := range ids {
		ids[i] = i
	}
	sort.Slice(ids, func(i, j int) bool {
		if counts[ids[i]] != counts[ids[j]] {
			return counts[ids[i]] > counts[ids[j]]
		}
		return ids[i] < ids[j]
	})
	var classes [512]int
	var sizes [8]int
	for rank, id := range ids {
		c := rank / 64
		classes[id] = c
		sizes[c]++
	}
	for _, n := range sizes {
		if float64(n) < m["minimum_refined_class_size"] {
			m["minimum_refined_class_size"] = float64(n)
		}
		if float64(n) > m["maximum_refined_class_size"] {
			m["maximum_refined_class_size"] = float64(n)
		}
	}

	rows := make([]wlmLmExternalContextStateHoldoutR1PairRow, 0, len(pairs))
	for p, c := range pairs {
		rows = append(rows, wlmLmExternalContextStateHoldoutR1PairRow{pair: p, count: c})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		if rows[i].pair.prev != rows[j].pair.prev {
			return rows[i].pair.prev < rows[j].pair.prev
		}
		return rows[i].pair.curr < rows[j].pair.curr
	})
	n := 511
	if len(rows) < n {
		n = len(rows)
	}
	exactID := make(map[wlmLmExternalContextStateHoldoutR1Pair]int, n)
	for i := 0; i < n; i++ {
		exactID[rows[i].pair] = i
	}
	var exact [512][512]uint32
	var exactTotals [512]uint64
	for _, st := range streams {
		for j := 1; j+1 < len(st); j++ {
			p := wlmLmExternalContextStateHoldoutR1Pair{prev: st[j-1], curr: st[j]}
			id := 511
			if v, ok := exactID[p]; ok {
				id = v
			}
			t := st[j+1]
			if exact[id][t] == ^uint32(0) {
				m["counter_overflow_count"]++
				continue
			}
			exact[id][t]++
			exactTotals[id]++
		}
	}

	var eligible [512]bool
	var marginsByRow [512]float64
	margins := make([]float64, 0, 512)
	for i := 0; i < 512; i++ {
		var total uint64
		var a, b uint32
		for _, v := range rel[i] {
			total += uint64(v)
			if v > a {
				b = a
				a = v
			} else if v > b {
				b = v
			}
		}
		if total >= 4 {
			eligible[i] = true
			marginsByRow[i] = (float64(a) - float64(b)) / float64(total)
			margins = append(margins, marginsByRow[i])
		}
	}
	sort.Float64s(margins)
	median := 0.0
	if len(margins) == 0 {
		m["invalid_row_count"]++
	} else if len(margins)%2 == 1 {
		median = margins[len(margins)/2]
	} else {
		median = (margins[len(margins)/2-1] + margins[len(margins)/2]) / 2
	}

	gateMargin := func(dist *[512]float64) float64 {
		top := wlmLmExternalReasoningReadoutRefinementR1Top1(dist)
		classID := classes[top]
		classMass := 0.0
		runner := 0.0
		for k, p := range dist {
			if classes[k] != classID {
				continue
			}
			classMass += p
			if k != top && p > runner {
				runner = p
			}
		}
		if classMass <= 0 {
			m["invalid_row_count"]++
			return 0
		}
		return (dist[top] - runner) / classMass
	}

	trainingMargins := make([]float64, 0)
	for _, st := range streams {
		for j := 0; j+4 < len(st); j++ {
			prev, curr := st[j], st[j+1]
			rd, re, rr := wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(prev, curr, 3, 2048, true, &rel, &totals, &exact, &exactTotals, exactID, pairs, &eligible, &marginsByRow, median, &uni, unitotal)
			if re > m["maximum_probability_mass_error"] {
				m["maximum_probability_mass_error"] = re
			}
			if rr < m["minimum_repaired_retained_mass"] {
				m["minimum_repaired_retained_mass"] = rr
			}
			trainingMargins = append(trainingMargins, gateMargin(&rd))
		}
	}
	sort.Float64s(trainingMargins)
	m["training_margin_sample_count"] = float64(len(trainingMargins))
	gateThreshold := 0.0
	if len(trainingMargins) == 0 {
		m["invalid_row_count"]++
	} else if len(trainingMargins)%2 == 1 {
		gateThreshold = trainingMargins[len(trainingMargins)/2]
	} else {
		gateThreshold = (trainingMargins[len(trainingMargins)/2-1] + trainingMargins[len(trainingMargins)/2]) / 2
	}
	m["training_margin_threshold"] = gateThreshold

	for si, data := range evals {
		st := wlmLmExternalMotifRelationHoldoutR1Decode(data, index)
		for j := 0; j+4 < len(st); j++ {
			prev, curr, target := st[j], st[j+1], st[j+4]
			bd, be, _ := wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(prev, curr, 3, 2048, false, &rel, &totals, &exact, &exactTotals, exactID, pairs, &eligible, &marginsByRow, median, &uni, unitotal)
			rd, re, rr := wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(prev, curr, 3, 2048, true, &rel, &totals, &exact, &exactTotals, exactID, pairs, &eligible, &marginsByRow, median, &uni, unitotal)
			if be > m["maximum_probability_mass_error"] {
				m["maximum_probability_mass_error"] = be
			}
			if re > m["maximum_probability_mass_error"] {
				m["maximum_probability_mass_error"] = re
			}
			if rr < m["minimum_repaired_retained_mass"] {
				m["minimum_repaired_retained_mass"] = rr
			}

			m["depth3_query_count"]++
			gateOpen := gateMargin(&rd) >= gateThreshold
			btop := wlmLmExternalReasoningReadoutRefinementR1Top1(&bd)
			rtop := wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
			if gateOpen {
				m["gate_activation_count"]++
				if rtop == target {
					m["gate_exact_hit_count"]++
				}
			}
			if btop == target || rtop != target {
				continue
			}
			m["depth3_exact_repair_gain_query_count"]++

			truth := classes[target]
			bc := wlmLmExternalReasoningReadoutRefinementR1Class(&bd, &classes)
			if classes[btop] != classes[rtop] && bc != truth {
				continue
			}
			m["hidden_gain_count"]++
			m["hidden_gain_source_"+sources[si].domain+"_count"]++

			if classes[btop] == classes[rtop] {
				m["same_class_hidden_gain_count"]++
				if gateOpen {
					m["same_class_hidden_gate_activation_count"]++
				}
				m["same_class_source_"+sources[si].domain+"_count"]++
				bm, rm := 0.0, 0.0
				for k := 0; k < 512; k++ {
					if classes[k] == truth {
						bm += bd[k]
						rm += rd[k]
					}
				}
				if bm <= 0 || rm <= 0 {
					m["invalid_row_count"]++
					continue
				}
				bshare := bd[target] / bm
				rshare := rd[target] / rm
				if bshare < m["minimum_baseline_truth_motif_share"] {
					m["minimum_baseline_truth_motif_share"] = bshare
				}
				if rshare < m["minimum_repaired_truth_motif_share"] {
					m["minimum_repaired_truth_motif_share"] = rshare
				}
				delta := rshare - bshare
				if delta > m["maximum_truth_motif_share_increase"] {
					m["maximum_truth_motif_share_increase"] = delta
				}
				if delta < m["minimum_truth_motif_share_increase"] {
					m["minimum_truth_motif_share_increase"] = delta
				}
				if delta > 1e-15 {
					m["same_class_truth_motif_share_increase_count"]++
				}
			} else {
				m["cross_class_hidden_gain_count"]++
			}
		}
	}

	if m["depth3_query_count"] > 0 {
		m["gate_activation_fraction"] = m["gate_activation_count"] / m["depth3_query_count"]
	}
	if m["gate_activation_count"] > 0 {
		m["gate_exact_precision"] = m["gate_exact_hit_count"] / m["gate_activation_count"]
	}
	h := m["hidden_gain_count"]
	same := m["same_class_hidden_gain_count"]
	if same > 0 {
		m["same_class_truth_motif_share_increase_fraction"] = m["same_class_truth_motif_share_increase_count"] / same
		m["same_class_hidden_gate_activation_fraction"] = m["same_class_hidden_gate_activation_count"] / same
	} else {
		m["invalid_row_count"]++
	}

	for _, domain := range []string{"code", "structured", "technical_prose"} {
		den := m["hidden_gain_source_"+domain+"_count"]
		if den > 0 {
			m["source_rate_denominator_nonzero_count"]++
			m["same_class_source_"+domain+"_rate"] = m["same_class_source_"+domain+"_count"] / den
		}
	}
	m["code_structured_same_class_rate_gap"] = math.Abs(m["same_class_source_code_rate"] - m["same_class_source_structured_rate"])

	if m["same_class_hidden_gain_count"]+m["cross_class_hidden_gain_count"] != h {
		m["attribution_accounting_error_count"]++
	}
	if m["hidden_gain_source_code_count"]+m["hidden_gain_source_structured_count"]+m["hidden_gain_source_technical_prose_count"] != h {
		m["attribution_accounting_error_count"]++
	}
	if m["same_class_source_code_count"]+m["same_class_source_structured_count"]+m["same_class_source_technical_prose_count"] != same {
		m["attribution_accounting_error_count"]++
	}
	if m["maximum_probability_mass_error"] > 1e-9 {
		m["invalid_row_count"]++
	}
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmLmExternalReasoningProbabilityMarginReadoutR1Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-EXTERNAL-REASONING-PROBABILITY-MARGIN-READOUT-R1",
		Metrics: m,
	}
}
