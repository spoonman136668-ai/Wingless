package unitary

import "math"

const UP98BBalancedThreewaySchema = "wingless.up98b-balanced-threeway.v1"

type UP98BClassMetric struct {
	Arm       string  `json:"arm"`
	Split     string  `json:"split"`
	Label     string  `json:"label"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Examples  int     `json:"examples"`
}

type UP98BVerbMetric struct {
	Arm      string  `json:"arm"`
	Verb     string  `json:"verb"`
	Accuracy float64 `json:"accuracy"`
	Examples int     `json:"examples"`
}

type UP98BConfusion struct {
	Arm    string `json:"arm"`
	Target string `json:"target"`
	Pred   string `json:"pred"`
	Count  int    `json:"count"`
}

type UP98BRoutingPoint struct {
	Arm                    string  `json:"arm"`
	TotalWrites            int     `json:"total_writes"`
	TargetKeys             int     `json:"target_keys"`
	QueryAccuracy          float64 `json:"query_accuracy"`
	ExactTargetSetAccuracy float64 `json:"exact_target_set_accuracy"`
	RecallEntriesUsed      int     `json:"recall_entries_used"`
}

type UP98BBalancedThreewayResult struct {
	Schema                         string              `json:"schema"`
	Experiment                     string              `json:"experiment"`
	SourceUP97BSeal                string              `json:"source_up97b_seal"`
	StateDimension                 int                 `json:"state_dimension"`
	Epochs                         int                 `json:"epochs"`
	LearningRate                   float64             `json:"learning_rate"`
	ClassResamplingUsed            bool                `json:"class_resampling_used"`
	ExplicitClassAtLearnedInference bool               `json:"explicit_class_at_learned_inference"`
	ClassMetrics                   []UP98BClassMetric  `json:"class_metrics"`
	VerbMetrics                    []UP98BVerbMetric   `json:"verb_metrics"`
	Confusion                      []UP98BConfusion    `json:"confusion"`
	RoutingPoints                  []UP98BRoutingPoint `json:"routing_points"`
}

type up98bOVRClassifier struct {
	w [3][64]float64
	b [3]float64
}

func up98bSigmoid(x float64) float64 {
	if x >= 0 {
		z := math.Exp(-x)
		return 1 / (1 + z)
	}
	z := math.Exp(x)
	return z / (1 + z)
}

func (c *up98bOVRClassifier) probs(h [64]float64) [3]float64 {
	var p [3]float64
	for k := 0; k < 3; k++ {
		s := c.b[k]
		for i := 0; i < 64; i++ {
			s += c.w[k][i] * h[i]
		}
		p[k] = up98bSigmoid(s)
	}
	return p
}

func up98bArgmax(p [3]float64) int {
	best := 0
	for k := 1; k < 3; k++ {
		if p[k] > p[best] {
			best = k
		}
	}
	return best
}

func up98bTrainOVR() *up98bOVRClassifier {
	c := &up98bOVRClassifier{}
	for epoch := 0; epoch < 20; epoch++ {
		for ni := 0; ni < 6; ni++ {
			for vali := 0; vali < 6; vali++ {
				for vi := 0; vi < 9; vi++ {
					if !up97bTrainSplit(ni, vali, vi) {
						continue
					}
					h := up95bEncode(up97bSurface(ni, vi, vali))
					target := up97bVerbClass(vi)
					for k := 0; k < 3; k++ {
						y := 0.0
						if k == target {
							y = 1
						}
						s := c.b[k]
						for i := 0; i < 64; i++ {
							s += c.w[k][i] * h[i]
						}
						g := up98bSigmoid(s) - y
						for i := 0; i < 64; i++ {
							c.w[k][i] -= 0.08 * g * h[i]
						}
						c.b[k] -= 0.08 * g
					}
				}
			}
		}
	}
	return c
}

func up98bLabel(class int) string {
	switch class {
	case up97bStore:
		return "STORE"
	case up97bObserve:
		return "OBSERVE"
	default:
		return "REPORT"
	}
}

func up98bEvalClass(arm string, pred func(string) int, split string, class int, label string) UP98BClassMetric {
	total, hits, tp, fp, fn := 0, 0, 0, 0, 0
	for ni := 0; ni < 6; ni++ {
		for vali := 0; vali < 6; vali++ {
			for vi := 0; vi < 9; vi++ {
				isTrain := up97bTrainSplit(ni, vali, vi)
				if split == "train" && !isTrain {
					continue
				}
				if split == "heldout_recombination" && isTrain {
					continue
				}
				target := up97bVerbClass(vi)
				got := pred(up97bSurface(ni, vi, vali))
				total++
				if got == target {
					hits++
				}
				if class >= 0 {
					if got == class && target == class {
						tp++
					}
					if got == class && target != class {
						fp++
					}
					if got != class && target == class {
						fn++
					}
				}
			}
		}
	}
	precision, recall := 1.0, 1.0
	if class >= 0 {
		if tp+fp > 0 {
			precision = float64(tp) / float64(tp+fp)
		}
		if tp+fn > 0 {
			recall = float64(tp) / float64(tp+fn)
		}
	}
	return UP98BClassMetric{Arm: arm, Split: split, Label: label, Accuracy: float64(hits) / float64(total), Precision: precision, Recall: recall, Examples: total}
}

func up98bEvalVerb(arm string, pred func(string) int, vi int) UP98BVerbMetric {
	total, hits := 0, 0
	for ni := 0; ni < 6; ni++ {
		for vali := 0; vali < 6; vali++ {
			if up97bTrainSplit(ni, vali, vi) {
				continue
			}
			target := up97bVerbClass(vi)
			got := pred(up97bSurface(ni, vi, vali))
			total++
			if got == target {
				hits++
			}
		}
	}
	return UP98BVerbMetric{Arm: arm, Verb: up97bVerbs[vi], Accuracy: float64(hits) / float64(total), Examples: total}
}

func up98bConfusion(arm string, pred func(string) int) []UP98BConfusion {
	var counts [3][3]int
	for ni := 0; ni < 6; ni++ {
		for vali := 0; vali < 6; vali++ {
			for vi := 0; vi < 9; vi++ {
				if up97bTrainSplit(ni, vali, vi) {
					continue
				}
				target := up97bVerbClass(vi)
				got := pred(up97bSurface(ni, vi, vali))
				counts[target][got]++
			}
		}
	}
	out := make([]UP98BConfusion, 0, 9)
	for target := 0; target < 3; target++ {
		for got := 0; got < 3; got++ {
			out = append(out, UP98BConfusion{Arm: arm, Target: up98bLabel(target), Pred: up98bLabel(got), Count: counts[target][got]})
		}
	}
	return out
}

func up98bRoute(arm string, pred func(string) int, totalWrites, targetKeys int, seedBases []int) UP98BRoutingPoint {
	qHits, qTotal, exactHits, episodes := 0, 0, 0, 0
	maxRecall := 0
	for _, base := range seedBases {
		for ep := 0; ep < 64; ep++ {
			seed := sq0Seed(base, 1601+totalWrites*7+targetKeys, ep)
			rng := newSQ0RNG(seed)
			m := newSQ0Machine("transport_gated_correction", seed)
			truth := make(map[int]int, targetKeys)
			event := func(class, key, value int) {
				ni := key % 6
				if ni < 0 {
					ni = -ni
				}
				vali := value % 6
				if vali < 0 {
					vali = -vali
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
				if arm != "explicit_event_class" {
					predClass = pred(up97bSurface(ni, vi, vali))
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
				if arm != "explicit_event_class" {
					predClass = pred(up97bSurface(ni, vi, 0))
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
	return UP98BRoutingPoint{Arm: arm, TotalWrites: totalWrites, TargetKeys: targetKeys, QueryAccuracy: float64(qHits) / float64(qTotal), ExactTargetSetAccuracy: float64(exactHits) / float64(episodes), RecallEntriesUsed: maxRecall}
}

func RunUP98B() (UP98BBalancedThreewayResult, error) {
	softmax := up97bTrainClassifier()
	ovr := up98bTrainOVR()
	softmaxPred := func(s string) int { return up97bArgmax(softmax.probs(up95bEncode(s))) }
	ovrPred := func(s string) int { return up98bArgmax(ovr.probs(up95bEncode(s))) }

	result := UP98BBalancedThreewayResult{
		Schema: UP98BBalancedThreewaySchema,
		Experiment: "UP-98B-balanced-threeway",
		SourceUP97BSeal: "aca2c8ab1c465c3978f5bcb5078c663e789f4d02",
		StateDimension: 64,
		Epochs: 20,
		LearningRate: 0.08,
		ClassResamplingUsed: false,
		ExplicitClassAtLearnedInference: false,
	}

	for _, item := range []struct {
		arm  string
		pred func(string) int
	}{
		{"softmax_threeway", softmaxPred},
		{"ovr_logistic_threeway", ovrPred},
	} {
		result.ClassMetrics = append(result.ClassMetrics,
			up98bEvalClass(item.arm, item.pred, "train", -1, "all"),
			up98bEvalClass(item.arm, item.pred, "heldout_recombination", -1, "all"),
			up98bEvalClass(item.arm, item.pred, "heldout_recombination", up97bStore, "STORE"),
			up98bEvalClass(item.arm, item.pred, "heldout_recombination", up97bObserve, "OBSERVE"),
			up98bEvalClass(item.arm, item.pred, "heldout_recombination", up97bReport, "REPORT"),
		)
		for vi := 0; vi < 9; vi++ {
			result.VerbMetrics = append(result.VerbMetrics, up98bEvalVerb(item.arm, item.pred, vi))
		}
		result.Confusion = append(result.Confusion, up98bConfusion(item.arm, item.pred)...)
	}

	seedBases := []int{151000000, 152000000}
	for _, item := range []struct {
		arm  string
		pred func(string) int
	}{
		{"explicit_event_class", softmaxPred},
		{"softmax_threeway", softmaxPred},
		{"ovr_logistic_threeway", ovrPred},
	} {
		for _, writes := range []int{32, 64, 128, 256} {
			for _, targets := range []int{4, 8, 16} {
				result.RoutingPoints = append(result.RoutingPoints, up98bRoute(item.arm, item.pred, writes, targets, seedBases))
			}
		}
	}
	return result, nil
}
