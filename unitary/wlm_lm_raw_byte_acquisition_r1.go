package unitary

type wlmLmRawByteAcquisitionR1Context interface {
	Push(uint8)
	Snapshot() [2]uint8
	Ready() bool
}

type wlmLmRawByteAcquisitionR1SliceContext struct {
	values []uint8
}

func (c *wlmLmRawByteAcquisitionR1SliceContext) Push(v uint8) {
	c.values = append(c.values, v)
	if len(c.values) > 2 {
		c.values = c.values[len(c.values)-2:]
	}
}

func (c *wlmLmRawByteAcquisitionR1SliceContext) Snapshot() [2]uint8 {
	var out [2]uint8
	if len(c.values) == 2 {
		out[0], out[1] = c.values[0], c.values[1]
	}
	return out
}

func (c *wlmLmRawByteAcquisitionR1SliceContext) Ready() bool {
	return len(c.values) == 2
}

type wlmLmRawByteAcquisitionR1RingContext struct {
	values [2]uint8
	head   int
	length int
}

func (c *wlmLmRawByteAcquisitionR1RingContext) Push(v uint8) {
	if c.length < 2 {
		c.values[(c.head+c.length)%2] = v
		c.length++
		return
	}
	c.values[c.head] = v
	c.head = (c.head + 1) % 2
}

func (c *wlmLmRawByteAcquisitionR1RingContext) Snapshot() [2]uint8 {
	var out [2]uint8
	if c.length == 2 {
		out[0] = c.values[c.head]
		out[1] = c.values[(c.head+1)%2]
	}
	return out
}

func (c *wlmLmRawByteAcquisitionR1RingContext) Ready() bool {
	return c.length == 2
}

func wlmLmRawByteAcquisitionR1NewContext(substrate string) wlmLmRawByteAcquisitionR1Context {
	if substrate == "slice-backed-byte-context" {
		return &wlmLmRawByteAcquisitionR1SliceContext{}
	}
	return &wlmLmRawByteAcquisitionR1RingContext{}
}

type wlmLmRawByteAcquisitionR1Factorized struct {
	counts [8][2][2][2]uint32
}

func (l *wlmLmRawByteAcquisitionR1Factorized) Observe(a, b, next uint8, metrics map[string]float64) {
	for bit := 0; bit < 8; bit++ {
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

func (l *wlmLmRawByteAcquisitionR1Factorized) Predict(a, b uint8) uint8 {
	var out uint8
	for bit := 0; bit < 8; bit++ {
		x := (a >> bit) & 1
		y := (b >> bit) & 1
		if l.counts[bit][x][y][1] > l.counts[bit][x][y][0] {
			out |= 1 << bit
		}
	}
	return out
}

type wlmLmRawByteAcquisitionR1Memorizer struct {
	values [65536]uint8
	seen   [65536]bool
}

func wlmLmRawByteAcquisitionR1Index(a, b uint8) int {
	return int(a)<<8 | int(b)
}

func (m *wlmLmRawByteAcquisitionR1Memorizer) Observe(a, b, next uint8, metrics map[string]float64) {
	index := wlmLmRawByteAcquisitionR1Index(a, b)
	if m.seen[index] && m.values[index] != next {
		metrics["memorization_conflict_count"]++
		return
	}
	m.values[index] = next
	m.seen[index] = true
}

func (m *wlmLmRawByteAcquisitionR1Memorizer) Predict(a, b uint8) uint8 {
	index := wlmLmRawByteAcquisitionR1Index(a, b)
	if !m.seen[index] {
		return 0
	}
	return m.values[index]
}

type wlmLmRawByteAcquisitionR1Pair struct {
	a uint8
	b uint8
}

type wlmLmRawByteAcquisitionR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmRawByteAcquisitionR1Target(a, b uint8) uint8 {
	return a ^ b
}

func wlmLmRawByteAcquisitionR1Heldout(a, b uint8) bool {
	return (int(a)+int(b))%4 == 1
}

func wlmLmRawByteAcquisitionR1Permute(seed uint32, values []wlmLmRawByteAcquisitionR1Pair) {
	state := seed ^ 0x9e3779b9
	for i := len(values) - 1; i > 0; i-- {
		state = state*1664525 + 1013904223
		j := int(state % uint32(i+1))
		values[i], values[j] = values[j], values[i]
	}
}

func wlmLmRawByteAcquisitionR1Training(seed uint32) []wlmLmRawByteAcquisitionR1Pair {
	out := make([]wlmLmRawByteAcquisitionR1Pair, 0, 49152)
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			aa, bb := uint8(a), uint8(b)
			if wlmLmRawByteAcquisitionR1Heldout(aa, bb) {
				continue
			}
			out = append(out, wlmLmRawByteAcquisitionR1Pair{a: aa, b: bb})
		}
	}
	wlmLmRawByteAcquisitionR1Permute(seed, out)
	return out
}

func wlmLmRawByteAcquisitionR1HeldoutContexts(seed uint32) []wlmLmRawByteAcquisitionR1Pair {
	out := make([]wlmLmRawByteAcquisitionR1Pair, 0, 16384)
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			aa, bb := uint8(a), uint8(b)
			if wlmLmRawByteAcquisitionR1Heldout(aa, bb) {
				out = append(out, wlmLmRawByteAcquisitionR1Pair{a: aa, b: bb})
			}
		}
	}
	wlmLmRawByteAcquisitionR1Permute(seed^0x85ebca6b, out)
	return out
}

func wlmLmRawByteAcquisitionR1CoverageMissing() int {
	missing := 0
	for bit := 0; bit < 8; bit++ {
		seen := [2][2]bool{}
		for a := 0; a < 256; a++ {
			for b := 0; b < 256; b++ {
				aa, bb := uint8(a), uint8(b)
				if wlmLmRawByteAcquisitionR1Heldout(aa, bb) {
					continue
				}
				seen[(aa>>bit)&1][(bb>>bit)&1] = true
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

func wlmLmRawByteAcquisitionR1Train(seed uint32, substrate string, metrics map[string]float64) (wlmLmRawByteAcquisitionR1Factorized, wlmLmRawByteAcquisitionR1Memorizer) {
	var factorized wlmLmRawByteAcquisitionR1Factorized
	var memorizer wlmLmRawByteAcquisitionR1Memorizer
	for _, obs := range wlmLmRawByteAcquisitionR1Training(seed) {
		if wlmLmRawByteAcquisitionR1Heldout(obs.a, obs.b) {
			metrics["heldout_context_leak_count"]++
		}
		ctx := wlmLmRawByteAcquisitionR1NewContext(substrate)
		ctx.Push(obs.a)
		ctx.Push(obs.b)
		if !ctx.Ready() {
			metrics["invalid_context_rows"]++
			continue
		}
		snapshot := ctx.Snapshot()
		target := wlmLmRawByteAcquisitionR1Target(snapshot[0], snapshot[1])
		factorized.Observe(snapshot[0], snapshot[1], target, metrics)
		memorizer.Observe(snapshot[0], snapshot[1], target, metrics)
		metrics["training_transition_count"]++
	}
	return factorized, memorizer
}

func wlmLmRawByteAcquisitionR1FactorizedMismatch(a, b *wlmLmRawByteAcquisitionR1Factorized) int {
	mismatch := 0
	for bit := 0; bit < 8; bit++ {
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

func wlmLmRawByteAcquisitionR1Evaluate(
	seed uint32,
	sliceFactor, ringFactor *wlmLmRawByteAcquisitionR1Factorized,
	sliceMemo, ringMemo *wlmLmRawByteAcquisitionR1Memorizer,
	metrics map[string]float64,
) (float64, float64, float64, float64, float64, float64) {
	heldout := wlmLmRawByteAcquisitionR1HeldoutContexts(seed)
	sliceCorrect, ringCorrect := 0, 0
	sliceMemoCorrect, ringMemoCorrect := 0, 0
	for _, obs := range heldout {
		target := wlmLmRawByteAcquisitionR1Target(obs.a, obs.b)
		sp := sliceFactor.Predict(obs.a, obs.b)
		rp := ringFactor.Predict(obs.a, obs.b)
		if sp != rp {
			metrics["substrate_prediction_mismatch_count"]++
		}
		if sp == target {
			sliceCorrect++
		}
		if rp == target {
			ringCorrect++
		}
		if sliceMemo.Predict(obs.a, obs.b) == target {
			sliceMemoCorrect++
		}
		if ringMemo.Predict(obs.a, obs.b) == target {
			ringMemoCorrect++
		}
		metrics["heldout_one_step_evaluation_count"] += 2
	}

	starts := heldout
	if len(starts) > 1024 {
		starts = starts[:1024]
	}
	sliceRollCorrect, ringRollCorrect, rollTotal := 0, 0, 0
	for _, start := range starts {
		sliceCtx := wlmLmRawByteAcquisitionR1NewContext("slice-backed-byte-context")
		ringCtx := wlmLmRawByteAcquisitionR1NewContext("fixed-capacity-ring-byte-context")
		sliceCtx.Push(start.a)
		sliceCtx.Push(start.b)
		ringCtx.Push(start.a)
		ringCtx.Push(start.b)
		for step := 0; step < 16; step++ {
			if !sliceCtx.Ready() || !ringCtx.Ready() {
				metrics["invalid_context_rows"]++
				continue
			}
			ss := sliceCtx.Snapshot()
			rs := ringCtx.Snapshot()
			if ss != rs {
				metrics["substrate_context_observable_mismatch_count"]++
			}
			target := wlmLmRawByteAcquisitionR1Target(ss[0], ss[1])
			sp := sliceFactor.Predict(ss[0], ss[1])
			rp := ringFactor.Predict(rs[0], rs[1])
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
	return float64(sliceCorrect)/oneDen,
		float64(ringCorrect)/oneDen,
		float64(sliceMemoCorrect)/oneDen,
		float64(ringMemoCorrect)/oneDen,
		float64(sliceRollCorrect)/rollDen,
		float64(ringRollCorrect)/rollDen
}

// RunWlmLmRawByteAcquisitionR1 evaluates the frozen WBG-1 raw-byte acquisition matrix.
func RunWlmLmRawByteAcquisitionR1() interface{} {
	seeds := [...]uint32{11701, 11719, 11743, 11777}
	metrics := map[string]float64{
		"valid_seed_count": 0,
		"completed_substrate_runs": 0,
		"training_transition_count": 0,
		"heldout_one_step_evaluation_count": 0,
		"autoregressive_evaluation_transition_count": 0,
		"minimum_factorized_heldout_one_step_accuracy": 1,
		"maximum_memorization_heldout_one_step_accuracy": 0,
		"minimum_factorized_autoregressive_accuracy": 1,
		"substrate_prediction_mismatch_count": 0,
		"substrate_factorized_state_mismatch_count": 0,
		"substrate_context_observable_mismatch_count": 0,
		"heldout_context_leak_count": 0,
		"factorized_bit_coverage_missing_count": float64(wlmLmRawByteAcquisitionR1CoverageMissing()),
		"memorization_conflict_count": 0,
		"invalid_context_rows": 0,
		"counter_overflow_rows": 0,
	}
	for _, seed := range seeds {
		sliceFactor, sliceMemo := wlmLmRawByteAcquisitionR1Train(seed, "slice-backed-byte-context", metrics)
		ringFactor, ringMemo := wlmLmRawByteAcquisitionR1Train(seed, "fixed-capacity-ring-byte-context", metrics)
		metrics["completed_substrate_runs"] += 2
		metrics["substrate_factorized_state_mismatch_count"] += float64(wlmLmRawByteAcquisitionR1FactorizedMismatch(&sliceFactor, &ringFactor))
		sOne, rOne, sMemo, rMemo, sRoll, rRoll := wlmLmRawByteAcquisitionR1Evaluate(
			seed, &sliceFactor, &ringFactor, &sliceMemo, &ringMemo, metrics,
		)
		if sOne < metrics["minimum_factorized_heldout_one_step_accuracy"] {
			metrics["minimum_factorized_heldout_one_step_accuracy"] = sOne
		}
		if rOne < metrics["minimum_factorized_heldout_one_step_accuracy"] {
			metrics["minimum_factorized_heldout_one_step_accuracy"] = rOne
		}
		if sMemo > metrics["maximum_memorization_heldout_one_step_accuracy"] {
			metrics["maximum_memorization_heldout_one_step_accuracy"] = sMemo
		}
		if rMemo > metrics["maximum_memorization_heldout_one_step_accuracy"] {
			metrics["maximum_memorization_heldout_one_step_accuracy"] = rMemo
		}
		if sRoll < metrics["minimum_factorized_autoregressive_accuracy"] {
			metrics["minimum_factorized_autoregressive_accuracy"] = sRoll
		}
		if rRoll < metrics["minimum_factorized_autoregressive_accuracy"] {
			metrics["minimum_factorized_autoregressive_accuracy"] = rRoll
		}
		metrics["valid_seed_count"]++
	}
	return wlmLmRawByteAcquisitionR1Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-BYTE-ACQUISITION-R1",
		Metrics: metrics,
	}
}
