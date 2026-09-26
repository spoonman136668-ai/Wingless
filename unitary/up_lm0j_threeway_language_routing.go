package unitary

import (
	"math"
	"strings"
)

const UPLM0JThreeWayRoutingSchema = "wingless.up-lm0j-threeway-language-routing.v1"

type UPLM0JClassifierMetric struct {
	Split     string  `json:"split"`
	Label     string  `json:"label"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Examples  int     `json:"examples"`
}

type UPLM0JMetric struct {
	Arm                        string  `json:"arm"`
	Split                      string  `json:"split"`
	Top1Accuracy               float64 `json:"top1_accuracy"`
	Perplexity                 float64 `json:"perplexity"`
	DependentFirstByteAccuracy float64 `json:"dependent_first_byte_accuracy"`
	QuerySetExactAccuracy      float64 `json:"query_set_exact_accuracy"`
	Update0ExactAccuracy       float64 `json:"update0_exact_accuracy"`
	Update1ExactAccuracy       float64 `json:"update1_exact_accuracy"`
	Update2ExactAccuracy       float64 `json:"update2_exact_accuracy"`
	Update4ExactAccuracy       float64 `json:"update4_exact_accuracy"`
	AdmissionPrecision         float64 `json:"admission_precision"`
	AdmissionRecall            float64 `json:"admission_recall"`
	EventRoutingAccuracy       float64 `json:"event_routing_accuracy"`
	ReportRoutingAccuracy      float64 `json:"report_routing_accuracy"`
	MaxRecallEntries           int     `json:"max_recall_entries"`
	ExactRecallBytes           int     `json:"exact_recall_bytes"`
	RecurrentStateBytes        int     `json:"recurrent_state_bytes"`
}

type UPLM0JThreeWayRoutingResult struct {
	Schema                         string                    `json:"schema"`
	Experiment                     string                    `json:"experiment"`
	SourceUPLM0ISeal               string                    `json:"source_up_lm0i_seal"`
	StateDimension                 int                       `json:"state_dimension"`
	ExactRecallCap                 int                       `json:"exact_recall_cap"`
	ClassifierEpochs               int                       `json:"classifier_epochs"`
	ClassifierLearningRate         float64                   `json:"classifier_learning_rate"`
	ExplicitTypeAtLearnedInference bool                      `json:"explicit_type_at_learned_inference"`
	ValueBytesAtClassifierInference bool                     `json:"value_bytes_at_classifier_inference"`
	AttentionUsed                  bool                      `json:"attention_used"`
	FutureOracleUsed               bool                      `json:"future_oracle_used"`
	ClassifierMetrics              []UPLM0JClassifierMetric `json:"classifier_metrics"`
	Metrics                        []UPLM0JMetric           `json:"metrics"`
}

func uplm0jAnchor(name string, class int) string {
	switch class {
	case up97bStore:
		return name + " stores"
	case up97bObserve:
		return name + " observes"
	default:
		return name + " reports"
	}
}

func uplm0jTrueClass(verb string) int {
	switch verb {
	case "stores":
		return up97bStore
	case "observes":
		return up97bObserve
	case "reports":
		return up97bReport
	default:
		return -1
	}
}

func uplm0jTrainClassifier() *up97bClassifier {
	c := &up97bClassifier{}
	names := uplm0gNames()
	for epoch := 0; epoch < 20; epoch++ {
		for ni := 0; ni < 4; ni++ {
			for class := 0; class < 3; class++ {
				h := uplm0fEncode(uplm0jAnchor(names[ni], class))
				p := c.probs(h)
				for k := 0; k < 3; k++ {
					g := p[k]
					if k == class {
						g -= 1
					}
					for i := 0; i < 64; i++ {
						c.w[k][i] -= 0.08 * g * h[i]
					}
					c.b[k] -= 0.08 * g
				}
			}
		}
	}
	return c
}

func uplm0jPredict(c *up97bClassifier, anchor string) int {
	return up97bArgmax(c.probs(uplm0fEncode(anchor)))
}

func uplm0jClassLabel(class int) string {
	switch class {
	case up97bStore:
		return "STORE"
	case up97bObserve:
		return "OBSERVE"
	default:
		return "REPORT"
	}
}

func uplm0jClassifierEval(c *up97bClassifier, split string, start, end, class int, label string) UPLM0JClassifierMetric {
	names := uplm0gNames()
	total, hits, tp, fp, fn := 0, 0, 0, 0, 0
	for ni := start; ni < end; ni++ {
		for target := 0; target < 3; target++ {
			pred := uplm0jPredict(c, uplm0jAnchor(names[ni], target))
			total++
			if pred == target {
				hits++
			}
			if class >= 0 {
				if pred == class && target == class {
					tp++
				}
				if pred == class && target != class {
					fp++
				}
				if pred != class && target == class {
					fn++
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
	return UPLM0JClassifierMetric{
		Split: split,
		Label: label,
		Accuracy: float64(hits) / float64(total),
		Precision: precision,
		Recall: recall,
		Examples: total,
	}
}

func uplm0jEvaluate(model *uplm0aModel, classifier *up97bClassifier, examples []uplm0fExample, stream4, learned bool, split string) UPLM0JMetric {
	hits, total := 0, 0
	depHits, depTotal := 0, 0
	exactParagraphs, paragraphs := 0, 0
	nll := 0.0
	maxEntries := 0
	tp, fp, fn := 0, 0, 0
	eventHits, eventTotal := 0, 0
	reportHits, reportTotal := 0, 0
	exactBy := map[int]int{0: 0, 1: 0, 2: 0, 4: 0}
	totalBy := map[int]int{0: 0, 1: 0, 2: 0, 4: 0}

	for start := 0; start < len(examples); {
		end := start + 1
		if stream4 {
			end = start + 4
			if end > len(examples) {
				end = len(examples)
			}
		}
		var h [64]float64
		recall := newUPLM0CRecall()

		for e := start; e < end; e++ {
			s := examples[e].text
			queryCorrect := [4]bool{}
			querySeen := [4]bool{}
			clause := ""
			routeKnown := false
			routeClass := -1
			trueClass := -1
			queryName := ""
			reportOverridePending := false

			for t := 0; t < len(s)-1; t++ {
				b := s[t]
				h = uplm0aStep(h, b)
				clause += string(b)
				trim := strings.TrimSpace(clause)

				if b == ' ' && !routeKnown {
					fields := strings.Fields(trim)
					if len(fields) == 2 {
						trueClass = uplm0jTrueClass(fields[1])
						if trueClass >= 0 {
							routeClass = trueClass
							if learned {
								routeClass = uplm0jPredict(classifier, fields[0]+" "+fields[1])
							}
							routeKnown = true
							eventTotal++
							if routeClass == trueClass {
								eventHits++
							}
							if trueClass == up97bReport {
								reportTotal++
								if routeClass == up97bReport {
									reportHits++
								}
							}
							if routeClass == up97bReport {
								queryName = fields[0]
								reportOverridePending = true
							}
						}
					}
				}

				targetByte := s[t+1]
				target := model.index[int(targetByte)]
				p := model.probs(h)
				pred := uplm0aArgmax(p)
				prob := p[target]

				if reportOverridePending && queryName != "" {
					if value, ok := recall.values[queryName]; ok && len(value) > 0 {
						memByte := value[0]
						pred = model.index[int(memByte)]
						if memByte == targetByte {
							prob = 1
						} else {
							prob = 1e-12
						}
					}
					reportOverridePending = false
				}

				if prob < 1e-12 {
					prob = 1e-12
				}
				total++
				nll -= math.Log(prob)
				if pred == target {
					hits++
				}
				if qi, ok := uplm0dIsTarget(t+1, examples[e].targetPos); ok {
					depTotal++
					querySeen[qi] = true
					if pred == target {
						depHits++
						queryCorrect[qi] = true
					}
				}

				if b == '.' {
					clean := strings.TrimSuffix(strings.TrimSpace(clause), ".")
					fields := strings.Fields(clean)
					if len(fields) >= 3 && trueClass >= 0 {
						name := fields[0]
						value := fields[2]
						predStore := routeKnown && routeClass == up97bStore
						trueStore := trueClass == up97bStore
						if predStore {
							recall.write(name, value)
							if trueStore {
								tp++
							} else {
								fp++
							}
						} else if trueStore {
							fn++
						}
					}
					clause = ""
					routeKnown = false
					routeClass = -1
					trueClass = -1
					queryName = ""
					reportOverridePending = false
				}

				if len(recall.order) > maxEntries {
					maxEntries = len(recall.order)
				}
			}

			all := true
			for i := 0; i < 4; i++ {
				if !querySeen[i] || !queryCorrect[i] {
					all = false
				}
			}
			paragraphs++
			totalBy[examples[e].updateCount]++
			if all {
				exactParagraphs++
				exactBy[examples[e].updateCount]++
			}
			if stream4 {
				h = uplm0aStep(h, s[len(s)-1])
			}
		}
		start = end
	}

	precision, recallScore := 1.0, 1.0
	if tp+fp > 0 {
		precision = float64(tp) / float64(tp+fp)
	}
	if tp+fn > 0 {
		recallScore = float64(tp) / float64(tp+fn)
	}
	acc := func(k int) float64 {
		if totalBy[k] == 0 {
			return 0
		}
		return float64(exactBy[k]) / float64(totalBy[k])
	}
	arm := "explicit_event_routing"
	if learned {
		arm = "learned_prefix_threeway"
	}
	eventAccuracy := 1.0
	if eventTotal > 0 {
		eventAccuracy = float64(eventHits) / float64(eventTotal)
	}
	reportAccuracy := 1.0
	if reportTotal > 0 {
		reportAccuracy = float64(reportHits) / float64(reportTotal)
	}

	return UPLM0JMetric{
		Arm: arm,
		Split: split,
		Top1Accuracy: float64(hits) / float64(total),
		Perplexity: math.Exp(nll / float64(total)),
		DependentFirstByteAccuracy: float64(depHits) / float64(depTotal),
		QuerySetExactAccuracy: float64(exactParagraphs) / float64(paragraphs),
		Update0ExactAccuracy: acc(0),
		Update1ExactAccuracy: acc(1),
		Update2ExactAccuracy: acc(2),
		Update4ExactAccuracy: acc(4),
		AdmissionPrecision: precision,
		AdmissionRecall: recallScore,
		EventRoutingAccuracy: eventAccuracy,
		ReportRoutingAccuracy: reportAccuracy,
		MaxRecallEntries: maxEntries,
		ExactRecallBytes: maxEntries * 16,
		RecurrentStateBytes: 512,
	}
}

func RunUPLM0J() (UPLM0JThreeWayRoutingResult, error) {
	train, held, alphabet := uplm0fCorpus()
	model := newUPLM0AModel(alphabet)
	for epoch := 0; epoch < 20; epoch++ {
		for _, ex := range train {
			model.trainSentence(ex.text, 0.08)
		}
	}
	classifier := uplm0jTrainClassifier()

	result := UPLM0JThreeWayRoutingResult{
		Schema: UPLM0JThreeWayRoutingSchema,
		Experiment: "UP-LM0J-threeway-language-routing",
		SourceUPLM0ISeal: "792ea898d6c33ca42b6b4031db0d9d4a82044cd9",
		StateDimension: 64,
		ExactRecallCap: 16,
		ClassifierEpochs: 20,
		ClassifierLearningRate: 0.08,
		ExplicitTypeAtLearnedInference: false,
		ValueBytesAtClassifierInference: false,
		AttentionUsed: false,
		FutureOracleUsed: false,
	}

	for _, split := range []struct {
		name       string
		start, end int
	}{
		{"train_subjects", 0, 4},
		{"heldout_subjects", 4, 6},
	} {
		result.ClassifierMetrics = append(result.ClassifierMetrics,
			uplm0jClassifierEval(classifier, split.name, split.start, split.end, -1, "all"),
			uplm0jClassifierEval(classifier, split.name, split.start, split.end, up97bStore, "STORE"),
			uplm0jClassifierEval(classifier, split.name, split.start, split.end, up97bObserve, "OBSERVE"),
			uplm0jClassifierEval(classifier, split.name, split.start, split.end, up97bReport, "REPORT"),
		)
	}

	result.Metrics = []UPLM0JMetric{
		uplm0jEvaluate(model, classifier, held, false, false, "heldout_recombination"),
		uplm0jEvaluate(model, classifier, held, true, false, "heldout_stream4"),
		uplm0jEvaluate(model, classifier, held, false, true, "heldout_recombination"),
		uplm0jEvaluate(model, classifier, held, true, true, "heldout_stream4"),
	}
	return result, nil
}
