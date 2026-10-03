package unitary

import "math"

type wlmLmFutureDataMarginalCurveGuidedPolicyTransferAttributionR17Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmExternalFutureDataMarginalCurveGuidedPolicyTransferAttributionR17(code, structured, prose []byte) interface{} {
	type src struct {
		d string
		b []byte
		h string
		n int
		budgets []int
	}
	ss := []src{
		{"code", code, "66bb25b24a0316b4965c64798494de93a1d7332672b15b5f430ab6a2fb4b9d45", 41453, []int{0, 436, 582, 727, 872}},
		{"structured", structured, "95ddbd0eaef29aad5ecfc74f9da21b795481f58b2c59380324a445fcd4d08932", 14365, []int{0, 291, 436, 581, 726}},
		{"technical_prose", prose, "48c3d95b8b03864a4af41d892710675956cde85afd0d5d6c331594de9f17881b", 1454, []int{0, 291, 436, 581, 726}},
	}
	m := map[string]float64{
		"source_identity_mismatch_count": 0,
		"transfer_manifest_identity_mismatch_count": 0,
		"source_count": 3,
		"total_source_bytes": 0,
		"arm_count": 3,
		"non_target_training_byte_budget_per_arm": 15819,
		"motif_capacity": 256,
		"target_evaluation_bytes_per_arm": 582,
		"target_adaptation_reservoir_bytes_per_arm": 872,
		"budget_point_count_total": 15,
		"capacity_growth_event_count": 0,
		"tokenizer_use_count": 0,
		"external_model_call_count": 0,
		"counter_overflow_count": 0,
		"invalid_row_count": 0,
	}
	for _, s := range ss {
		m["total_source_bytes"] += float64(len(s.b))
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
		hits := make([]int, len(t.budgets))
		for bi, bud := range t.budgets {
			train := [][]byte{a, b}
			if bud > 0 {
				train = append(train, res[:bud])
			}
			rel, sel := wlmLmR10Train256(train, m)
			if len(sel) != 256 {
				m["invalid_row_count"]++
			}
			h := wlmLmR10Hits(eval, &rel, sel)
			hits[bi] = h
			m["arm_"+t.d+"_budget_"+itoaR17(bud)+"_exact_hits"] = float64(h)
		}
		for i := 1; i < len(t.budgets); i++ {
			lo, hi := t.budgets[i-1], t.budgets[i]
			gain := hits[i] - hits[i-1]
			eff := float64(gain) / float64(hi-lo)
			k := "arm_" + t.d + "_interval_" + itoaR17(lo) + "_" + itoaR17(hi)
			m[k+"_gain"] = float64(gain)
			m[k+"_efficiency"] = eff
		}
	}
	m["structured_gain_436_581"] = m["arm_structured_interval_436_581_gain"]
	m["technical_prose_gain_581_726"] = m["arm_technical_prose_interval_581_726_gain"]
	m["predicted_curve_net_delta"] = m["technical_prose_gain_581_726"] - m["structured_gain_436_581"]
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			m["invalid_row_count"]++
		}
	}
	return wlmLmFutureDataMarginalCurveGuidedPolicyTransferAttributionR17Result{
		"wingless.research-scientific-result.v1",
		"WLM-LM-EXTERNAL-FUTURE-DATA-MARGINAL-CURVE-GUIDED-POLICY-TRANSFER-ATTRIBUTION-R17",
		m,
	}
}

func itoaR17(n int) string {
	if n == 0 {
		return "0"
	}
	b := [20]byte{}
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
