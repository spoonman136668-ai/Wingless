package unitary

type wlmLmNativeSequenceR2Context interface {
	Push(int)
	Snapshot() []int
}

type wlmLmNativeSequenceR2SliceContext struct {
	values []int
}

func (c *wlmLmNativeSequenceR2SliceContext) Push(v int) {
	c.values = append(c.values, v)
	if len(c.values) > 2 {
		c.values = c.values[len(c.values)-2:]
	}
}

func (c *wlmLmNativeSequenceR2SliceContext) Snapshot() []int {
	out := make([]int, len(c.values))
	copy(out, c.values)
	return out
}

type wlmLmNativeSequenceR2RingContext struct {
	values [2]int
	head   int
	length int
}

func (c *wlmLmNativeSequenceR2RingContext) Push(v int) {
	if c.length < 2 {
		c.values[(c.head+c.length)%2] = v
		c.length++
		return
	}
	c.values[c.head] = v
	c.head = (c.head + 1) % 2
}

func (c *wlmLmNativeSequenceR2RingContext) Snapshot() []int {
	out := make([]int, c.length)
	for i := range out {
		out[i] = c.values[(c.head+i)%2]
	}
	return out
}

type wlmLmNativeSequenceR2Learner struct {
	counts [8][8][8]uint32
}

func (l *wlmLmNativeSequenceR2Learner) Observe(a, b, next int, metrics map[string]float64) {
	if a < 0 || a >= 8 || b < 0 || b >= 8 || next < 0 || next >= 8 {
		metrics["invalid_context_rows"]++
		return
	}
	if l.counts[a][b][next] == ^uint32(0) {
		metrics["counter_overflow_rows"]++
		return
	}
	l.counts[a][b][next]++
}

func (l *wlmLmNativeSequenceR2Learner) Predict(a, b int) int {
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

type wlmLmNativeSequenceR2Observation struct {
	a int
	b int
}

type wlmLmNativeSequenceR2Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmNativeSequenceR2Next(a, b int) int {
	return (3*a + 5*b + 1) % 8
}

func wlmLmNativeSequenceR2Permute(seed uint32, values []wlmLmNativeSequenceR2Observation) {
	state := seed ^ 0x9e3779b9
	for i := len(values) - 1; i > 0; i-- {
		state = state*1664525 + 1013904223
		j := int(state % uint32(i+1))
		values[i], values[j] = values[j], values[i]
	}
}

func wlmLmNativeSequenceR2Training(seed uint32) []wlmLmNativeSequenceR2Observation {
	out := make([]wlmLmNativeSequenceR2Observation, 0, 2048)
	for repeat := 0; repeat < 32; repeat++ {
		for a := 0; a < 8; a++ {
			for b := 0; b < 8; b++ {
				out = append(out, wlmLmNativeSequenceR2Observation{a: a, b: b})
			}
		}
	}
	wlmLmNativeSequenceR2Permute(seed, out)
	return out
}

func wlmLmNativeSequenceR2Starts(seed uint32) []wlmLmNativeSequenceR2Observation {
	out := make([]wlmLmNativeSequenceR2Observation, 0, 64)
	for a := 0; a < 8; a++ {
		for b := 0; b < 8; b++ {
			out = append(out, wlmLmNativeSequenceR2Observation{a: a, b: b})
		}
	}
	wlmLmNativeSequenceR2Permute(seed^0xa5a5a5a5, out)
	return out
}

func wlmLmNativeSequenceR2NewContext(substrate string) wlmLmNativeSequenceR2Context {
	if substrate == "slice-backed-context" {
		return &wlmLmNativeSequenceR2SliceContext{}
	}
	return &wlmLmNativeSequenceR2RingContext{}
}

func wlmLmNativeSequenceR2ContextPairValid(snapshot []int) bool {
	return len(snapshot) == 2 && snapshot[0] >= 0 && snapshot[0] < 8 && snapshot[1] >= 0 && snapshot[1] < 8
}

func wlmLmNativeSequenceR2Train(seed uint32, substrate string, metrics map[string]float64) wlmLmNativeSequenceR2Learner {
	var learner wlmLmNativeSequenceR2Learner
	for _, obs := range wlmLmNativeSequenceR2Training(seed) {
		ctx := wlmLmNativeSequenceR2NewContext(substrate)
		ctx.Push(obs.a)
		ctx.Push(obs.b)
		snapshot := ctx.Snapshot()
		if !wlmLmNativeSequenceR2ContextPairValid(snapshot) {
			metrics["invalid_context_rows"]++
			continue
		}
		learner.Observe(snapshot[0], snapshot[1], wlmLmNativeSequenceR2Next(snapshot[0], snapshot[1]), metrics)
		metrics["training_transition_count"]++
	}
	return learner
}

func wlmLmNativeSequenceR2EqualState(a, b *wlmLmNativeSequenceR2Learner) int {
	mismatches := 0
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			for z := 0; z < 8; z++ {
				if a.counts[x][y][z] != b.counts[x][y][z] {
					mismatches++
				}
			}
		}
	}
	return mismatches
}

func wlmLmNativeSequenceR2Evaluate(seed uint32, sliceLearner, ringLearner *wlmLmNativeSequenceR2Learner, metrics map[string]float64) (float64, float64, float64, float64) {
	starts := wlmLmNativeSequenceR2Starts(seed)
	sliceCorrect, ringCorrect := 0, 0
	sliceBaseCorrect, ringBaseCorrect := 0, 0
	total := 0

	for _, start := range starts {
		sliceCtx := wlmLmNativeSequenceR2NewContext("slice-backed-context")
		ringCtx := wlmLmNativeSequenceR2NewContext("fixed-capacity-ring-context")
		sliceBaseCtx := wlmLmNativeSequenceR2NewContext("slice-backed-context")
		ringBaseCtx := wlmLmNativeSequenceR2NewContext("fixed-capacity-ring-context")
		for _, ctx := range []wlmLmNativeSequenceR2Context{sliceCtx, ringCtx, sliceBaseCtx, ringBaseCtx} {
			ctx.Push(start.a)
			ctx.Push(start.b)
		}

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

			target := wlmLmNativeSequenceR2Next(sliceState[0], sliceState[1])
			slicePrediction := sliceLearner.Predict(sliceState[0], sliceState[1])
			ringPrediction := ringLearner.Predict(ringState[0], ringState[1])
			if slicePrediction != ringPrediction {
				metrics["substrate_prediction_mismatch_count"]++
			}
			if slicePrediction == target {
				sliceCorrect++
			}
			if ringPrediction == target {
				ringCorrect++
			}

			sliceBaseState := sliceBaseCtx.Snapshot()
			ringBaseState := ringBaseCtx.Snapshot()
			sliceBaseTarget := wlmLmNativeSequenceR2Next(sliceBaseState[0], sliceBaseState[1])
			ringBaseTarget := wlmLmNativeSequenceR2Next(ringBaseState[0], ringBaseState[1])
			sliceBasePrediction := 0
			ringBasePrediction := 0
			if sliceBasePrediction == sliceBaseTarget {
				sliceBaseCorrect++
			}
			if ringBasePrediction == ringBaseTarget {
				ringBaseCorrect++
			}

			sliceCtx.Push(slicePrediction)
			ringCtx.Push(ringPrediction)
			sliceBaseCtx.Push(sliceBasePrediction)
			ringBaseCtx.Push(ringBasePrediction)
			total++
		}
	}

	metrics["evaluation_transition_count"] += float64(total * 2)
	den := float64(total)
	return float64(sliceCorrect) / den, float64(ringCorrect) / den, float64(sliceBaseCorrect) / den, float64(ringBaseCorrect) / den
}

func wlmLmNativeSequenceR2Min(current, candidate float64) float64 {
	if candidate < current {
		return candidate
	}
	return current
}

func wlmLmNativeSequenceR2Max(current, candidate float64) float64 {
	if candidate > current {
		return candidate
	}
	return current
}

// RunWlmLmNativeSequenceR2 evaluates the frozen native sequence-learning matrix.
func RunWlmLmNativeSequenceR2() interface{} {
	seeds := [...]uint32{11003, 11027, 11047, 11071, 11083, 11113, 11131, 11149}
	metrics := map[string]float64{
		"valid_seed_count": 0,
		"completed_substrate_runs": 0,
		"training_transition_count": 0,
		"evaluation_transition_count": 0,
		"minimum_trained_next_token_accuracy": 1,
		"maximum_untrained_baseline_accuracy": 0,
		"minimum_accuracy_improvement_over_untrained": 1,
		"substrate_prediction_mismatch_count": 0,
		"substrate_learned_state_mismatch_count": 0,
		"substrate_context_observable_mismatch_count": 0,
		"invalid_context_rows": 0,
		"counter_overflow_rows": 0,
	}

	for _, seed := range seeds {
		sliceLearner := wlmLmNativeSequenceR2Train(seed, "slice-backed-context", metrics)
		ringLearner := wlmLmNativeSequenceR2Train(seed, "fixed-capacity-ring-context", metrics)
		metrics["completed_substrate_runs"] += 2
		metrics["substrate_learned_state_mismatch_count"] += float64(wlmLmNativeSequenceR2EqualState(&sliceLearner, &ringLearner))

		sliceAccuracy, ringAccuracy, sliceBaseline, ringBaseline := wlmLmNativeSequenceR2Evaluate(seed, &sliceLearner, &ringLearner, metrics)
		metrics["minimum_trained_next_token_accuracy"] = wlmLmNativeSequenceR2Min(metrics["minimum_trained_next_token_accuracy"], sliceAccuracy)
		metrics["minimum_trained_next_token_accuracy"] = wlmLmNativeSequenceR2Min(metrics["minimum_trained_next_token_accuracy"], ringAccuracy)
		metrics["maximum_untrained_baseline_accuracy"] = wlmLmNativeSequenceR2Max(metrics["maximum_untrained_baseline_accuracy"], sliceBaseline)
		metrics["maximum_untrained_baseline_accuracy"] = wlmLmNativeSequenceR2Max(metrics["maximum_untrained_baseline_accuracy"], ringBaseline)
		metrics["minimum_accuracy_improvement_over_untrained"] = wlmLmNativeSequenceR2Min(metrics["minimum_accuracy_improvement_over_untrained"], sliceAccuracy-sliceBaseline)
		metrics["minimum_accuracy_improvement_over_untrained"] = wlmLmNativeSequenceR2Min(metrics["minimum_accuracy_improvement_over_untrained"], ringAccuracy-ringBaseline)
		metrics["valid_seed_count"]++
	}

	return wlmLmNativeSequenceR2Result{
		Schema:     "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-NATIVE-SEQUENCE-R2",
		Metrics:    metrics,
	}
}
