package unitary

import "math"

type wlmLmFutureDataMarginalCurveGuidedPolicyTransferR16Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalFutureDataMarginalCurveGuidedPolicyTransferR16(code, structured, prose []byte) interface{} {
	type src struct {
		d string
		b []byte
		h string
		n int
		e int
		g int
	}
	ss := []src{
		{"code", code, "66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45", 41453, 582, 582},
		{"structured", structured, "95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932", 14365, 581, 436},
		{"technical_prose", prose, "48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b", 1454, 581, 726},
	}
	m := map[string]float64{
		"source_identity_mismatch_count": 0,
		"transfer_manifest_identity_mismatch_count": 0,
		"source_count": 3,
		"total_source_bytes": 0,
		"arm_count": 3,
		"non_target_training_byte_budget_per_arm": 15819,
		"motif_capacity": 256,
		"target_window_bytes_per_arm": 1454,
		"target_evaluation_bytes_per_arm": 582,
		"target_adaptation_reservoir_bytes_per_arm": 872,
		"total_adaptation_budget_equal": 0,
		"total_adaptation_budget_curve_guided": 0,
		"allocation_target_label_use_count": 0,
		"allocation_post_result_choice_count": 0,
		"packet0_aggregate_exact_hits": 0,
		"equal_aggregate_exact_hits": 0,
		"curve_guided_aggregate_exact_hits": 0,
		"curve_guided_aggregate_exact_hit_delta_vs_equal": 0,
		"curve_guided_technical_prose_exact_hit_delta_vs_equal": 0,
		"curve_guided_domain_below_packet0_count": 0,
		"capacity_growth_event_count": 0,
		"tokenizer_use_count": 0,
		"external_model_call_count": 0,
		"counter_overflow_count": 0,
		"invalid_row_count": 0,
	}
	for _, s := range ss {
		m["total_source_bytes"] += float64(len(s.b))
		m["total_adaptation_budget_equal"] += float64(s.e)
		m["total_adaptation_budget_curve_guided"] += float64(s.g)
		if len(s.b) != s.n || wlmLmExternalRawRepPredR1SHA256(s.b) != s.h {
			m["source_identity_mismatch_count"]++
		}
	}
	for ti, t := range ss {
		idx := []int{}
		for i := range ss {
			if i != ti {
				idx = append(idx, i)
			}
		}
		a, b := wlmLmTransferCapacityBudgetR2(ss[idx[0]].b, ss[idx[1]].b)
		if len(a)+len(b) != 15819 || len(t.b) < 1454 {
			m["invalid_row_count"]++
			continue
		}
		eval, res := t.b[:582], t.b[582:1454]
		br, bs := wlmLmR10Train256([][]byte{a, b}, m)
		er, es := wlmLmR10Train256([][]byte{a, b, res[:t.e]}, m)
		gr, gs := wlmLmR10Train256([][]byte{a, b, res[:t.g]}, m)
		if len(bs) != 256 || len(es) != 256 || len(gs) != 256 {
			m["invalid_row_count"]++
		}
		p0 := wlmLmR10Hits(eval, &br, bs)
		eh := wlmLmR10Hits(eval, &er, es)
		gh := wlmLmR10Hits(eval, &gr, gs)
		p := "arm_" + t.d + "_"
		m[p+"packet0_exact_hits"] = float64(p0)
		m[p+"equal_exact_hits"] = float64(eh)
		m[p+"curve_guided_exact_hits"] = float64(gh)
		m[p+"curve_guided_delta_vs_equal"] = float64(gh - eh)
		m[p+"curve_guided_delta_vs_packet0"] = float64(gh - p0)
		m["packet0_aggregate_exact_hits"] += float64(p0)
		m["equal_aggregate_exact_hits"] += float64(eh)
		m["curve_guided_aggregate_exact_hits"] += float64(gh)
		if gh < p0 {
			m["curve_guided_domain_below_packet0_count"]++
		}
		if t.d == "technical_prose" {
			m["curve_guided_technical_prose_exact_hit_delta_vs_equal"] = float64(gh - eh)
		}
	}
	m["curve_guided_aggregate_exact_hit_delta_vs_equal"] = m["curve_guided_aggregate_exact_hits"] - m["equal_aggregate_exact_hits"]
	if m["total_adaptation_budget_equal"] != 1744 || m["total_adaptation_budget_curve_guided"] != 1744 {
		m["invalid_row_count"]++
	}
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmLmFutureDataMarginalCurveGuidedPolicyTransferR16Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-MARGINAL-CURVE-GUIDED-POLICY-TRANSFER-R16",
		m,
	}
}
