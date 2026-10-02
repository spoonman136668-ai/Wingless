package unitary

type wlmLmCompositionalSequenceR1Factorized struct {
	counts [3][2][2][2]uint32
}

type wlmLmCompositionalSequenceR1Lookup struct {
	counts [8][8][8]uint32
}

type wlmLmCompositionalSequenceR1Observation struct {
	a int
	b int
}

type wlmLmCompositionalSequenceR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmCompositionalSequenceR1Target(a, b int) int {
	return a ^ b
}

func wlmLmCompositionalSequenceR1Heldout(a, b int) bool {
	return (a+b)%4 == 1
}

func wlmLmCompositionalSequenceR1Permute(seed uint32, values []wlmLmCompositionalSequenceR1Observation) {
	state := seed ^ 0x7f4a7c15
	for i := len(values) - 1; i > 0; i-- {
		state = state*1664525 + 1013904223
		j := int(state % uint32(i+1))
		values[i], values[j] = values[j], values[i]
	}
}

func wlmLmCompositionalSequenceR1Training(seed uint32) []wlmLmCompositionalSequenceR1Observation {
	out := make([]wlmLmCompositionalSequenceR1Observation, 0, 1536)
	for repeat := 0; repeat < 32; repeat++ {
		for a := 0; a < 8; a++ {
			for b := 0; b < 8; b++ {
				if wlmLmCompositionalSequenceR1Heldout(a, b) {
					continue
				}
				out = append(out, wlmLmCompositionalSequenceR1Observation{a: a, b: b})
			}
		}
	}
	wlmLmCompositionalSequenceR1Permute(seed, out)
	return out
}

func wlmLmCompositionalSequenceR1HeldoutContexts(seed uint32) []wlmLmCompositionalSequenceR1Observation {
	out := make([]wlmLmCompositionalSequenceR1Observation, 0, 16)
	for a := 0; a < 8; a++ {
		for b := 0; b < 8; b++ {
			if wlmLmCompositionalSequenceR1Heldout(a, b) {
				out = append(out, wlmLmCompositionalSequenceR1Observation{a: a, b: b})
			}
		}
	}
	wlmLmCompositionalSequenceR1Permute(seed^0x5bd1e995, out)
	return out
}

func (l *wlmLmCompositionalSequenceR1Factorized) Observe(a, b, next int, metrics map[string]float64) {
	if a < 0 || a >= 8 || b < 0 || b >= 8 || next < 0 || next >= 8 {
		metrics["invalid_context_rows"]++
		return
	}
	for bit := 0; bit < 3; bit++ {
		x := (a >> bit) & 1
		y := (b >> bit) & 1
		t := (next >> bit) & 1
		if l.counts[bit][x][y][t] == ^uint32(0) {
			metrics["counter_overflow_rows"]++
			continue
		}
		l.counts[bit][x][y][t]++
	}
}

func (l *wlmLmCompositionalSequenceR1Factorized) Predict(a, b int) int {
	out := 0
	for bit := 0; bit < 3; bit++ {
		x := (a >> bit) & 1
		y := (b >> bit) & 1
		if l.counts[bit][x][y][1] > l.counts[bit][x][y][0] {
			out |= 1 << bit
		}
	}
	return out
}

func (l *wlmLmCompositionalSequenceR1Lookup) Observe(a, b, next int, metrics map[string]float64) {
	if l.counts[a][b][next] == ^uint32(0) {
		metrics["counter_overflow_rows"]++
		return
	}
	l.counts[a][b][next]++
}

func (l *wlmLmCompositionalSequenceR1Lookup) Predict(a, b int) int {
	best := 0
	bestCount := l.counts[a][b][0]
	for token := 1; token < 8; token++ {
		if l.counts[a][b][token] > bestCount {
			best = token
			bestCount = l.counts[a][b][token]
		}
	}
	return best
}

func wlmLmCompositionalSequenceR1CoverageMissing() int {
	missing := 0
	for bit := 0; bit < 3; bit++ {
		seen := [2][2]bool{}
		for a := 0; a < 8; a++ {
			for b := 0; b < 8; b++ {
				if wlmLmCompositionalSequenceR1Heldout(a, b) {
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
	return missing
}

func wlmLmCompositionalSequenceR1Train(seed uint32, substrate string, metrics map[string]float64) (wlmLmCompositionalSequenceR1Factorized, wlmLmCompositionalSequenceR1Lookup) {
	var factorized wlmLmCompositionalSequenceR1Factorized
	var lookup wlmLmCompositionalSequenceR1Lookup
	for _, obs := range wlmLmCompositionalSequenceR1Training(seed) {
		if wlmLmCompositionalSequenceR1Heldout(obs.a, obs.b) {
			metrics["heldout_context_leak_count"]++
		}
		ctx := wlmLmNativeSequenceR2NewContext(substrate)
		ctx.Push(obs.a)
		ctx.Push(obs.b)
		snapshot := ctx.Snapshot()
		if !wlmLmNativeSequenceR2ContextPairValid(snapshot) {
			metrics["invalid_context_rows"]++
			continue
		}
		target := wlmLmCompositionalSequenceR1Target(snapshot[0], snapshot[1])
		factorized.Observe(snapshot[0], snapshot[1], target, metrics)
		lookup.Observe(snapshot[0], snapshot[1], target, metrics)
		metrics["training_transition_count"]++
	}
	return factorized, lookup
}

func wlmLmCompositionalSequenceR1FactorizedStateMismatch(a, b *wlmLmCompositionalSequenceR1Factorized) int {
	mismatch := 0
	for bit := 0; bit < 3; bit++ {
		for x := 0; x < 2; x++ {
			for y := 0; y < 2; y++ {
				for t := 0; t < 2; t++ {
					if a.counts[bit][x][y][t] != b.counts[bit][x][y][t] {
						mismatch++
					}
				}
			}
		}
	}
	return mismatch
}

func wlmLmCompositionalSequenceR1Evaluate(seed uint32, sliceFactor, ringFactor *wlmLmCompositionalSequenceR1Factorized, sliceLookup, ringLookup *wlmLmCompositionalSequenceR1Lookup, metrics map[string]float64) (float64, float64, float64, float64, float64, float64) {
	heldout := wlmLmCompositionalSequenceR1HeldoutContexts(seed)
	sliceOneCorrect, ringOneCorrect := 0, 0
	sliceLookupCorrect, ringLookupCorrect := 0, 0

	for _, obs := range heldout {
		target := wlmLmCompositionalSequenceR1Target(obs.a, obs.b)
		sp := sliceFactor.Predict(obs.a, obs.b)
		rp := ringFactor.Predict(obs.a, obs.b)
		if sp != rp {
			metrics["substrate_prediction_mismatch_count"]++
		}
		if sp == target {
			sliceOneCorrect++
		}
		if rp == target {
			ringOneCorrect++
		}
		if sliceLookup.Predict(obs.a, obs.b) == target {
			sliceLookupCorrect++
		}
		if ringLookup.Predict(obs.a, obs.b) == target {
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
			target := wlmLmCompositionalSequenceR1Target(sliceState[0], sliceState[1])
			sp := sliceFactor.Predict(sliceState[0], sliceState[1])
			rp := ringFactor.Predict(ringState[0], ringState[1])
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
	return float64(sliceOneCorrect) / oneDen, float64(ringOneCorrect) / oneDen,
		float64(sliceLookupCorrect) / oneDen, float64(ringLookupCorrect) / oneDen,
		float64(sliceRollCorrect) / rollDen, float64(ringRollCorrect) / rollDen
}

// RunWlmLmCompositionalSequenceR1 evaluates the frozen held-out compositional sequence matrix.
func RunWlmLmCompositionalSequenceR1() interface{} {
	seeds := [...]uint32{11261, 11287, 11311, 11329, 11351, 11369, 11383, 11411}
	metrics := map[string]float64{
		"valid_seed_count": 0,
		"completed_substrate_runs": 0,
		"training_transition_count": 0,
		"heldout_one_step_evaluation_count": 0,
		"autoregressive_evaluation_transition_count": 0,
		"minimum_factorized_heldout_one_step_accuracy": 1,
		"maximum_lookup_heldout_one_step_accuracy": 0,
		"minimum_factorized_autoregressive_accuracy": 1,
		"substrate_prediction_mismatch_count": 0,
		"substrate_factorized_state_mismatch_count": 0,
		"substrate_context_observable_mismatch_count": 0,
		"heldout_context_leak_count": 0,
		"factorized_bit_coverage_missing_count": float64(wlmLmCompositionalSequenceR1CoverageMissing()),
		"invalid_context_rows": 0,
		"counter_overflow_rows": 0,
	}

	for _, seed := range seeds {
		sliceFactor, sliceLookup := wlmLmCompositionalSequenceR1Train(seed, "slice-backed-context", metrics)
		ringFactor, ringLookup := wlmLmCompositionalSequenceR1Train(seed, "fixed-capacity-ring-context", metrics)
		metrics["completed_substrate_runs"] += 2
		metrics["substrate_factorized_state_mismatch_count"] += float64(wlmLmCompositionalSequenceR1FactorizedStateMismatch(&sliceFactor, &ringFactor))

		sOne, rOne, sLookup, rLookup, sRoll, rRoll := wlmLmCompositionalSequenceR1Evaluate(seed, &sliceFactor, &ringFactor, &sliceLookup, &ringLookup, metrics)
		if sOne < metrics["minimum_factorized_heldout_one_step_accuracy"] {
			metrics["minimum_factorized_heldout_one_step_accuracy"] = sOne
		}
		if rOne < metrics["minimum_factorized_heldout_one_step_accuracy"] {
			metrics["minimum_factorized_heldout_one_step_accuracy"] = rOne
		}
		if sLookup > metrics["maximum_lookup_heldout_one_step_accuracy"] {
			metrics["maximum_lookup_heldout_one_step_accuracy"] = sLookup
		}
		if rLookup > metrics["maximum_lookup_heldout_one_step_accuracy"] {
			metrics["maximum_lookup_heldout_one_step_accuracy"] = rLookup
		}
		if sRoll < metrics["minimum_factorized_autoregressive_accuracy"] {
			metrics["minimum_factorized_autoregressive_accuracy"] = sRoll
		}
		if rRoll < metrics["minimum_factorized_autoregressive_accuracy"] {
			metrics["minimum_factorized_autoregressive_accuracy"] = rRoll
		}
		metrics["valid_seed_count"]++
	}

	return wlmLmCompositionalSequenceR1Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-COMPOSITIONAL-SEQUENCE-R1",
		Metrics: metrics,
	}
}
