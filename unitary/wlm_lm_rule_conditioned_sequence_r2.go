package unitary

type wlmLmRuleConditionedSequenceR2Factorized struct {
	counts [4][3][2][2][2]uint32
}

type wlmLmRuleConditionedSequenceR2Lookup struct {
	counts [4][8][8][8]uint32
}

type wlmLmRuleConditionedSequenceR2Observation struct {
	rule int
	a    int
	b    int
}

type wlmLmRuleConditionedSequenceR2Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmRuleConditionedSequenceR2Target(rule, a, b int) int {
	switch rule {
	case 0:
		return a ^ b
	case 1:
		return (^(a ^ b)) & 7
	case 2:
		return a & b
	case 3:
		return a | b
	default:
		return 0
	}
}

func wlmLmRuleConditionedSequenceR2Heldout(rule, a, b int) bool {
	return (a+b+rule)%4 == 1
}

func wlmLmRuleConditionedSequenceR2Permute(seed uint32, values []wlmLmRuleConditionedSequenceR2Observation) {
	state := seed ^ 0x6d2b79f5
	for i := len(values) - 1; i > 0; i-- {
		state = state*1664525 + 1013904223
		j := int(state % uint32(i+1))
		values[i], values[j] = values[j], values[i]
	}
}

func wlmLmRuleConditionedSequenceR2Training(seed uint32) []wlmLmRuleConditionedSequenceR2Observation {
	out := make([]wlmLmRuleConditionedSequenceR2Observation, 0, 3072)
	for repeat := 0; repeat < 16; repeat++ {
		for rule := 0; rule < 4; rule++ {
			for a := 0; a < 8; a++ {
				for b := 0; b < 8; b++ {
					if wlmLmRuleConditionedSequenceR2Heldout(rule, a, b) {
						continue
					}
					out = append(out, wlmLmRuleConditionedSequenceR2Observation{rule: rule, a: a, b: b})
				}
			}
		}
	}
	wlmLmRuleConditionedSequenceR2Permute(seed, out)
	return out
}

func wlmLmRuleConditionedSequenceR2HeldoutContexts(seed uint32) []wlmLmRuleConditionedSequenceR2Observation {
	out := make([]wlmLmRuleConditionedSequenceR2Observation, 0, 64)
	for rule := 0; rule < 4; rule++ {
		for a := 0; a < 8; a++ {
			for b := 0; b < 8; b++ {
				if wlmLmRuleConditionedSequenceR2Heldout(rule, a, b) {
					out = append(out, wlmLmRuleConditionedSequenceR2Observation{rule: rule, a: a, b: b})
				}
			}
		}
	}
	wlmLmRuleConditionedSequenceR2Permute(seed^0x85ebca6b, out)
	return out
}

func (l *wlmLmRuleConditionedSequenceR2Factorized) Observe(rule, a, b, next int, metrics map[string]float64) {
	if rule < 0 || rule >= 4 || a < 0 || a >= 8 || b < 0 || b >= 8 || next < 0 || next >= 8 {
		metrics["invalid_context_rows"]++
		return
	}
	for bit := 0; bit < 3; bit++ {
		x := (a >> bit) & 1
		y := (b >> bit) & 1
		t := (next >> bit) & 1
		if l.counts[rule][bit][x][y][t] == ^uint32(0) {
			metrics["counter_overflow_rows"]++
			continue
		}
		l.counts[rule][bit][x][y][t]++
	}
}

func (l *wlmLmRuleConditionedSequenceR2Factorized) Predict(rule, a, b int) int {
	out := 0
	for bit := 0; bit < 3; bit++ {
		x := (a >> bit) & 1
		y := (b >> bit) & 1
		if l.counts[rule][bit][x][y][1] > l.counts[rule][bit][x][y][0] {
			out |= 1 << bit
		}
	}
	return out
}

func (l *wlmLmRuleConditionedSequenceR2Lookup) Observe(rule, a, b, next int, metrics map[string]float64) {
	if l.counts[rule][a][b][next] == ^uint32(0) {
		metrics["counter_overflow_rows"]++
		return
	}
	l.counts[rule][a][b][next]++
}

func (l *wlmLmRuleConditionedSequenceR2Lookup) Predict(rule, a, b int) int {
	best := 0
	bestCount := l.counts[rule][a][b][0]
	for token := 1; token < 8; token++ {
		if l.counts[rule][a][b][token] > bestCount {
			best = token
			bestCount = l.counts[rule][a][b][token]
		}
	}
	return best
}

func wlmLmRuleConditionedSequenceR2CoverageMissing() int {
	missing := 0
	for rule := 0; rule < 4; rule++ {
		for bit := 0; bit < 3; bit++ {
			seen := [2][2]bool{}
			for a := 0; a < 8; a++ {
				for b := 0; b < 8; b++ {
					if wlmLmRuleConditionedSequenceR2Heldout(rule, a, b) {
						continue
					}
					seen[(a>>bit)&1][(b>>bit)&1] = true
				}
			}
			for x := 0; x < 2; x++ {
				for y := 0; y < 2; y++ {
					if !seen[x][y] {
						missing++
					}
				}
			}
		}
	}
	return missing
}

func wlmLmRuleConditionedSequenceR2Train(seed uint32, substrate string, metrics map[string]float64) (wlmLmRuleConditionedSequenceR2Factorized, wlmLmRuleConditionedSequenceR2Lookup) {
	var factorized wlmLmRuleConditionedSequenceR2Factorized
	var lookup wlmLmRuleConditionedSequenceR2Lookup
	for _, obs := range wlmLmRuleConditionedSequenceR2Training(seed) {
		if wlmLmRuleConditionedSequenceR2Heldout(obs.rule, obs.a, obs.b) {
			metrics["heldout_rule_context_leak_count"]++
		}
		ctx := wlmLmNativeSequenceR2NewContext(substrate)
		ctx.Push(obs.a)
		ctx.Push(obs.b)
		snapshot := ctx.Snapshot()
		if !wlmLmNativeSequenceR2ContextPairValid(snapshot) {
			metrics["invalid_context_rows"]++
			continue
		}
		target := wlmLmRuleConditionedSequenceR2Target(obs.rule, snapshot[0], snapshot[1])
		factorized.Observe(obs.rule, snapshot[0], snapshot[1], target, metrics)
		lookup.Observe(obs.rule, snapshot[0], snapshot[1], target, metrics)
		metrics["training_transition_count"]++
	}
	return factorized, lookup
}

func wlmLmRuleConditionedSequenceR2FactorizedStateMismatch(a, b *wlmLmRuleConditionedSequenceR2Factorized) int {
	mismatch := 0
	for rule := 0; rule < 4; rule++ {
		for bit := 0; bit < 3; bit++ {
			for x := 0; x < 2; x++ {
				for y := 0; y < 2; y++ {
					for t := 0; t < 2; t++ {
						if a.counts[rule][bit][x][y][t] != b.counts[rule][bit][x][y][t] {
							mismatch++
						}
					}
				}
			}
		}
	}
	return mismatch
}

func wlmLmRuleConditionedSequenceR2Evaluate(
	seed uint32,
	sliceFactor, ringFactor *wlmLmRuleConditionedSequenceR2Factorized,
	sliceLookup, ringLookup *wlmLmRuleConditionedSequenceR2Lookup,
	metrics map[string]float64,
) (float64, float64, float64, float64, float64, float64) {
	heldout := wlmLmRuleConditionedSequenceR2HeldoutContexts(seed)
	sliceOneCorrect, ringOneCorrect := 0, 0
	sliceLookupCorrect, ringLookupCorrect := 0, 0

	for _, obs := range heldout {
		target := wlmLmRuleConditionedSequenceR2Target(obs.rule, obs.a, obs.b)
		sp := sliceFactor.Predict(obs.rule, obs.a, obs.b)
		rp := ringFactor.Predict(obs.rule, obs.a, obs.b)
		if sp != rp {
			metrics["substrate_prediction_mismatch_count"]++
		}
		if sp == target {
			sliceOneCorrect++
		}
		if rp == target {
			ringOneCorrect++
		}
		if sliceLookup.Predict(obs.rule, obs.a, obs.b) == target {
			sliceLookupCorrect++
		}
		if ringLookup.Predict(obs.rule, obs.a, obs.b) == target {
			ringLookupCorrect++
		}
		metrics["heldout_one_step_evaluation_count"] += 2
	}

	sliceRollCorrect, ringRollCorrect := 0, 0
	rollTotal := 0
	for _, start := range heldout {
		sliceCtx := wlmLmNativeSequenceR2NewContext("slice-backed-context")
		ringCtx := wlmLmNativeSequenceR2NewContext("fixed-capacity-ring-context")
		sliceCtx.Push(start.a)
		sliceCtx.Push(start.b)
		ringCtx.Push(start.a)
		ringCtx.Push(start.b)

		for step := 0; step < 16; step++ {
			sliceState := sliceCtx.Snapshot()
			ringState := ringCtx.Snapshot()
			if !wlmLmNativeSequenceR2ContextPairValid(sliceState) || !wlmLmNativeSequenceR2ContextPairValid(ringState) {
				metrics["invalid_context_rows"]++
				continue
			}
			if sliceState[0] != ringState[0] || sliceState[1] != ringState[1] {
				metrics["substrate_context_observable_mismatch_count"]++
			}
			target := wlmLmRuleConditionedSequenceR2Target(start.rule, sliceState[0], sliceState[1])
			sp := sliceFactor.Predict(start.rule, sliceState[0], sliceState[1])
			rp := ringFactor.Predict(start.rule, ringState[0], ringState[1])
			if sp != rp {
				metrics["substrate_prediction_mismatch_count"]++
			}
			if sp == target {
				sliceRollCorrect++
			}
			if rp == target {
				ringRollCorrect++
			}
			sliceCtx.Push(sp)
			ringCtx.Push(rp)
			rollTotal++
		}
	}
	metrics["autoregressive_evaluation_transition_count"] += float64(rollTotal * 2)
	oneDen := float64(len(heldout))
	rollDen := float64(rollTotal)
	return float64(sliceOneCorrect) / oneDen,
		float64(ringOneCorrect) / oneDen,
		float64(sliceLookupCorrect) / oneDen,
		float64(ringLookupCorrect) / oneDen,
		float64(sliceRollCorrect) / rollDen,
		float64(ringRollCorrect) / rollDen
}

// RunWlmLmRuleConditionedSequenceR2 evaluates the frozen multi-rule held-out sequence matrix.
func RunWlmLmRuleConditionedSequenceR2() interface{} {
	seeds := [...]uint32{11503, 11527, 11549, 11579, 11617, 11633, 11657, 11681}
	metrics := map[string]float64{
		"valid_seed_count": 0,
		"completed_substrate_runs": 0,
		"training_transition_count": 0,
		"heldout_one_step_evaluation_count": 0,
		"autoregressive_evaluation_transition_count": 0,
		"minimum_factorized_heldout_one_step_accuracy": 1,
		"maximum_pooled_lookup_heldout_one_step_accuracy": 0,
		"minimum_factorized_autoregressive_accuracy": 1,
		"substrate_prediction_mismatch_count": 0,
		"substrate_factorized_state_mismatch_count": 0,
		"substrate_context_observable_mismatch_count": 0,
		"heldout_rule_context_leak_count": 0,
		"factorized_bit_coverage_missing_count": float64(wlmLmRuleConditionedSequenceR2CoverageMissing()),
		"invalid_context_rows": 0,
		"counter_overflow_rows": 0,
	}

	for _, seed := range seeds {
		sliceFactor, sliceLookup := wlmLmRuleConditionedSequenceR2Train(seed, "slice-backed-context", metrics)
		ringFactor, ringLookup := wlmLmRuleConditionedSequenceR2Train(seed, "fixed-capacity-ring-context", metrics)
		metrics["completed_substrate_runs"] += 2
		metrics["substrate_factorized_state_mismatch_count"] += float64(
			wlmLmRuleConditionedSequenceR2FactorizedStateMismatch(&sliceFactor, &ringFactor),
		)

		sOne, rOne, sLookup, rLookup, sRoll, rRoll := wlmLmRuleConditionedSequenceR2Evaluate(
			seed, &sliceFactor, &ringFactor, &sliceLookup, &ringLookup, metrics,
		)
		if sOne < metrics["minimum_factorized_heldout_one_step_accuracy"] {
			metrics["minimum_factorized_heldout_one_step_accuracy"] = sOne
		}
		if rOne < metrics["minimum_factorized_heldout_one_step_accuracy"] {
			metrics["minimum_factorized_heldout_one_step_accuracy"] = rOne
		}
		if sLookup > metrics["maximum_pooled_lookup_heldout_one_step_accuracy"] {
			metrics["maximum_pooled_lookup_heldout_one_step_accuracy"] = sLookup
		}
		if rLookup > metrics["maximum_pooled_lookup_heldout_one_step_accuracy"] {
			metrics["maximum_pooled_lookup_heldout_one_step_accuracy"] = rLookup
		}
		if sRoll < metrics["minimum_factorized_autoregressive_accuracy"] {
			metrics["minimum_factorized_autoregressive_accuracy"] = sRoll
		}
		if rRoll < metrics["minimum_factorized_autoregressive_accuracy"] {
			metrics["minimum_factorized_autoregressive_accuracy"] = rRoll
		}
		metrics["valid_seed_count"]++
	}

	return wlmLmRuleConditionedSequenceR2Result{
		Schema:     "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RULE-CONDITIONED-SEQUENCE-R2",
		Metrics:    metrics,
	}
}
