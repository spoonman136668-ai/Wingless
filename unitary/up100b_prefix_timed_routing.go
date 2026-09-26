package unitary

const UP100BPrefixRoutingSchema = "wingless.up100b-prefix-timed-routing.v1"

type UP100BRoutingPoint struct {
	Arm                    string  `json:"arm"`
	TotalWrites            int     `json:"total_writes"`
	TargetKeys             int     `json:"target_keys"`
	QueryAccuracy          float64 `json:"query_accuracy"`
	ExactTargetSetAccuracy float64 `json:"exact_target_set_accuracy"`
	RecallEntriesUsed      int     `json:"recall_entries_used"`
}

type UP100BPrefixRoutingResult struct {
	Schema                         string               `json:"schema"`
	Experiment                     string               `json:"experiment"`
	SourceUP99BSeal                string               `json:"source_up99b_seal"`
	StateDimension                 int                  `json:"state_dimension"`
	Epochs                         int                  `json:"epochs"`
	LearningRate                   float64              `json:"learning_rate"`
	ExactRecallCap                 int                  `json:"exact_recall_cap"`
	ExplicitClassAtLearnedInference bool                `json:"explicit_class_at_learned_inference"`
	ClassMetrics                   []UP99BClassMetric   `json:"class_metrics"`
	VerbMetrics                    []UP99BVerbMetric    `json:"verb_metrics"`
	RoutingPoints                  []UP100BRoutingPoint `json:"routing_points"`
}

func up100bPredict(c *up97bClassifier, ni, vi int) int {
	return up97bArgmax(c.probs(up99bPrefixEncoder(ni, vi, 0)))
}

func up100bRoute(c *up97bClassifier, arm string, totalWrites, targetKeys int, seedBases []int) UP100BRoutingPoint {
	qHits, qTotal, exactHits, episodes := 0, 0, 0, 0
	maxRecall := 0
	for _, base := range seedBases {
		for ep := 0; ep < 64; ep++ {
			seed := sq0Seed(base, 1801+totalWrites*7+targetKeys, ep)
			rng := newSQ0RNG(seed)
			m := newSQ0Machine("transport_gated_correction", seed)
			truth := make(map[int]int, targetKeys)

			event := func(class, key, value int) {
				ni := key % 6
				if ni < 0 {
					ni = -ni
				}
				vi := 0
				switch class {
				case up97bStore:
					vi = (key + value + ep) % 3
				case up97bObserve:
					vi = 3 + (key+value+ep)%3
				case up97bReport:
					vi = 6 + (key+value+ep)%3
				}
				predClass := class
				if arm == "learned_prefix_threeway" {
					predClass = up100bPredict(c, ni, vi)
				}
				switch predClass {
				case up97bStore:
					m.write(key, value, 32)
					m.recallWrite(key, value)
				case up97bObserve:
					m.write(key, value, 32)
				case up97bReport:
				}
			}

			for k := 0; k < targetKeys; k++ {
				v := rng.intn(32)
				truth[k] = v
				event(up97bStore, k, v)
			}
			distractors := totalWrites - 2*targetKeys
			pre := distractors / 2
			post := distractors - pre
			for d := 0; d < pre; d++ {
				event(up97bObserve, 1000+ep*1000+d, rng.intn(32))
			}
			for k := 0; k < targetKeys; k++ {
				old := truth[k]
				v := (old + 1 + rng.intn(31)) % 32
				truth[k] = v
				event(up97bStore, k, v)
			}
			for d := 0; d < post; d++ {
				event(up97bObserve, 200000+ep*1000+d, rng.intn(32))
			}

			exact := true
			for k := 0; k < targetKeys; k++ {
				ni := k % 6
				vi := 6 + (k+ep)%3
				predClass := up97bReport
				if arm == "learned_prefix_threeway" {
					predClass = up100bPredict(c, ni, vi)
				}
				got, ok := 0, false
				if predClass == up97bReport {
					got, ok = m.recallQuery(k)
				}
				qTotal++
				if ok && got == truth[k] {
					qHits++
				} else {
					exact = false
				}
			}
			episodes++
			if exact {
				exactHits++
			}
			if m.recallUsed > maxRecall {
				maxRecall = m.recallUsed
			}
		}
	}
	return UP100BRoutingPoint{
		Arm: arm,
		TotalWrites: totalWrites,
		TargetKeys: targetKeys,
		QueryAccuracy: float64(qHits) / float64(qTotal),
		ExactTargetSetAccuracy: float64(exactHits) / float64(episodes),
		RecallEntriesUsed: maxRecall,
	}
}

func RunUP100B() (UP100BPrefixRoutingResult, error) {
	c := up99bTrain(up99bPrefixEncoder)
	result := UP100BPrefixRoutingResult{
		Schema: UP100BPrefixRoutingSchema,
		Experiment: "UP-100B-prefix-timed-routing",
		SourceUP99BSeal: "b54496cb0cd508acfb0a6fac65de96a57ed49a1f",
		StateDimension: 64,
		Epochs: 20,
		LearningRate: 0.08,
		ExactRecallCap: 16,
		ExplicitClassAtLearnedInference: false,
		ClassMetrics: []UP99BClassMetric{
			up99bEval("prefix_through_verb", "train", "all", -1, c, up99bPrefixEncoder),
			up99bEval("prefix_through_verb", "heldout_recombination", "all", -1, c, up99bPrefixEncoder),
			up99bEval("prefix_through_verb", "heldout_recombination", "STORE", up97bStore, c, up99bPrefixEncoder),
			up99bEval("prefix_through_verb", "heldout_recombination", "OBSERVE", up97bObserve, c, up99bPrefixEncoder),
			up99bEval("prefix_through_verb", "heldout_recombination", "REPORT", up97bReport, c, up99bPrefixEncoder),
		},
	}
	for vi := 0; vi < 9; vi++ {
		result.VerbMetrics = append(result.VerbMetrics, up99bVerbEval("prefix_through_verb", vi, c, up99bPrefixEncoder))
	}
	seedBases := []int{155000000, 156000000}
	for _, arm := range []string{"explicit_event_class", "learned_prefix_threeway"} {
		for _, writes := range []int{32, 64, 128, 256} {
			for _, targets := range []int{4, 8, 16} {
				result.RoutingPoints = append(result.RoutingPoints, up100bRoute(c, arm, writes, targets, seedBases))
			}
		}
	}
	return result, nil
}
