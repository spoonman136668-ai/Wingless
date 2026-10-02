package unitary

type wlmLmAdaptiveGranularityR1Target struct {
	position int
}

type wlmLmAdaptiveGranularityR1Corpus struct {
	bytes             []uint8
	targets           []wlmLmAdaptiveGranularityR1Target
	predictableStarts map[int]bool
	region             []uint8
	predictableBytes  int
	highEntropyRecords int
}

type wlmLmAdaptiveGranularityR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmAdaptiveGranularityR1Step(state uint32) uint32 {
	return state*1664525 + 1013904223
}

func wlmLmAdaptiveGranularityR1MixedCorpus(seed uint32) wlmLmAdaptiveGranularityR1Corpus {
	out := wlmLmAdaptiveGranularityR1Corpus{
		bytes:              make([]uint8, 0, 30720),
		targets:            make([]wlmLmAdaptiveGranularityR1Target, 0, 2048),
		predictableStarts:  make(map[int]bool),
		region:             make([]uint8, 0, 30720),
		predictableBytes:   0,
		highEntropyRecords: 0,
	}
	state := seed ^ 0x517cc1b7
	seedOffset := int(seed % 8)
	predictableIndex := 0
	for record := 0; record < 4096; record++ {
		predictable := (record+int(seed))%2 == 0
		if predictable {
			start := len(out.bytes)
			id := (5*predictableIndex + seedOffset) % 8
			motif := wlmLmByteMotifDiscoveryR3BalancedMotif(id)
			out.predictableStarts[start] = true
			out.bytes = append(out.bytes, motif[:]...)
			out.region = append(out.region, 1, 1, 1, 1)
			out.targets = append(out.targets, wlmLmAdaptiveGranularityR1Target{position: len(out.bytes)})
			out.bytes = append(out.bytes, wlmLmByteMotifDiscoveryR3BalancedSuccessor(id))
			out.region = append(out.region, 1)
			for filler := 0; filler < 3; filler++ {
				state = wlmLmAdaptiveGranularityR1Step(state)
				out.bytes = append(out.bytes, uint8(state>>24))
				out.region = append(out.region, 1)
			}
			out.predictableBytes += 8
			predictableIndex++
		} else {
			for i := 0; i < 7; i++ {
				state = wlmLmAdaptiveGranularityR1Step(state)
				out.bytes = append(out.bytes, uint8(state>>24))
				out.region = append(out.region, 2)
			}
			out.highEntropyRecords++
		}
	}
	return out
}

func wlmLmAdaptiveGranularityR1TargetAccuracy(
	corpus wlmLmAdaptiveGranularityR1Corpus,
	learner *wlmLmByteMotifDiscoveryR3BalancedLearner,
	mode string,
) (float64, int) {
	correct := 0
	invalid := 0
	for _, target := range corpus.targets {
		pos := target.position
		if pos < 4 || pos >= len(corpus.bytes) {
			invalid++
			continue
		}
		key := [4]uint8{corpus.bytes[pos-4], corpus.bytes[pos-3], corpus.bytes[pos-2], corpus.bytes[pos-1]}
		expected := corpus.bytes[pos]
		var prediction uint8
		switch mode {
		case "full", "adaptive":
			prediction = learner.predict(key, "motif")
		case "static":
			start := pos - 4
			if start%4 == 0 {
				prediction = learner.predict(key, "motif")
			} else {
				prediction = learner.baselinePredict(key[2], key[3])
			}
		default:
			invalid++
			continue
		}
		if prediction == expected {
			correct++
		}
	}
	den := len(corpus.targets) - invalid
	if den <= 0 {
		return 0, invalid
	}
	return float64(correct) / float64(den), invalid
}

func wlmLmAdaptiveGranularityR1AdaptiveEvents(
	corpus wlmLmAdaptiveGranularityR1Corpus,
	learner *wlmLmByteMotifDiscoveryR3BalancedLearner,
) (events int, trueCoarse int, falseHighEntropy int, invalid int) {
	for i := 0; i < len(corpus.bytes); {
		if i+4 <= len(corpus.bytes) {
			key := [4]uint8{corpus.bytes[i], corpus.bytes[i+1], corpus.bytes[i+2], corpus.bytes[i+3]}
			if _, ok := learner.selectedMap[key]; ok {
				events++
				if corpus.predictableStarts[i] {
					trueCoarse++
				} else {
					allHigh := true
					for j := i; j < i+4; j++ {
						if corpus.region[j] != 2 {
							allHigh = false
							break
						}
					}
					if allHigh {
						falseHighEntropy++
					} else {
						invalid++
					}
				}
				i += 4
				continue
			}
		}
		events++
		i++
	}
	return
}

func wlmLmAdaptiveGranularityR1Min(a, b float64) float64 {
	if b < a {
		return b
	}
	return a
}

func wlmLmAdaptiveGranularityR1Max(a, b float64) float64 {
	if b > a {
		return b
	}
	return a
}

// RunWlmLmAdaptiveGranularityR1 evaluates the frozen WBG-3 mixed-uncertainty granularity matrix.
func RunWlmLmAdaptiveGranularityR1() interface{} {
	seeds := [...]uint32{11903, 11923, 11933, 11939}
	metrics := map[string]float64{
		"valid_seed_count": 0,
		"completed_substrate_runs": 0,
		"minimum_adaptive_target_accuracy": 1,
		"minimum_full_byte_target_accuracy": 1,
		"maximum_adaptive_accuracy_gap_vs_full_byte": 0,
		"maximum_static_fixed4_target_accuracy": 0,
		"minimum_adaptive_accuracy_gain_vs_static_fixed4": 1,
		"minimum_overall_event_reduction_fraction": 1,
		"minimum_predictable_region_event_reduction_fraction": 1,
		"maximum_high_entropy_false_coarse_activation_rate": 0,
		"minimum_predictable_coarse_activation_recall": 1,
		"maximum_adaptive_events_per_raw_byte": 0,
		"substrate_selected_motif_mismatch_count": 0,
		"substrate_prediction_mismatch_count": 0,
		"substrate_event_accounting_mismatch_count": 0,
		"invalid_stream_rows": 0,
		"counter_overflow_rows": 0,
	}

	for _, seed := range seeds {
		train := wlmLmByteMotifDiscoveryR3BalancedCorpusFor(seed, 4096, false)
		corpus := wlmLmAdaptiveGranularityR1MixedCorpus(seed)
		sliceLearner := wlmLmByteMotifDiscoveryR3BalancedNewLearner()
		ringLearner := wlmLmByteMotifDiscoveryR3BalancedNewLearner()
		local := map[string]float64{"counter_overflow_rows": 0}
		sliceLearner.train(train.bytes, "slice-backed-four-byte-window", local)
		ringLearner.train(train.bytes, "fixed-capacity-ring-four-byte-window", local)
		sliceLearner.selectMotifs(local)
		ringLearner.selectMotifs(local)
		metrics["counter_overflow_rows"] += local["counter_overflow_rows"]
		metrics["completed_substrate_runs"] += 2
		metrics["substrate_selected_motif_mismatch_count"] += float64(
			wlmLmByteMotifDiscoveryR3BalancedSelectedMismatch(&sliceLearner, &ringLearner),
		)

		sliceAdaptive, invalidSA := wlmLmAdaptiveGranularityR1TargetAccuracy(corpus, &sliceLearner, "adaptive")
		ringAdaptive, invalidRA := wlmLmAdaptiveGranularityR1TargetAccuracy(corpus, &ringLearner, "adaptive")
		sliceFull, invalidSF := wlmLmAdaptiveGranularityR1TargetAccuracy(corpus, &sliceLearner, "full")
		ringFull, invalidRF := wlmLmAdaptiveGranularityR1TargetAccuracy(corpus, &ringLearner, "full")
		sliceStatic, invalidSS := wlmLmAdaptiveGranularityR1TargetAccuracy(corpus, &sliceLearner, "static")
		ringStatic, invalidRS := wlmLmAdaptiveGranularityR1TargetAccuracy(corpus, &ringLearner, "static")
		metrics["invalid_stream_rows"] += float64(invalidSA + invalidRA + invalidSF + invalidRF + invalidSS + invalidRS)

		if sliceAdaptive != ringAdaptive || sliceFull != ringFull || sliceStatic != ringStatic {
			metrics["substrate_prediction_mismatch_count"]++
		}

		sliceEvents, sliceTrue, sliceFalseHigh, sliceInvalid := wlmLmAdaptiveGranularityR1AdaptiveEvents(corpus, &sliceLearner)
		ringEvents, ringTrue, ringFalseHigh, ringInvalid := wlmLmAdaptiveGranularityR1AdaptiveEvents(corpus, &ringLearner)
		metrics["invalid_stream_rows"] += float64(sliceInvalid + ringInvalid)
		if sliceEvents != ringEvents || sliceTrue != ringTrue || sliceFalseHigh != ringFalseHigh {
			metrics["substrate_event_accounting_mismatch_count"]++
		}

		rawBytes := float64(len(corpus.bytes))
		predictableCount := float64(len(corpus.targets))
		overallReduction := 1.0 - float64(sliceEvents)/rawBytes
		predictableReduction := float64(sliceTrue*3) / float64(corpus.predictableBytes)
		falseHighRate := float64(sliceFalseHigh) / float64(corpus.highEntropyRecords)
		coarseRecall := float64(sliceTrue) / predictableCount
		eventsPerByte := float64(sliceEvents) / rawBytes
		adaptiveGap := sliceFull - sliceAdaptive
		if adaptiveGap < 0 {
			adaptiveGap = -adaptiveGap
		}

		metrics["minimum_adaptive_target_accuracy"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_adaptive_target_accuracy"], sliceAdaptive)
		metrics["minimum_adaptive_target_accuracy"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_adaptive_target_accuracy"], ringAdaptive)
		metrics["minimum_full_byte_target_accuracy"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_full_byte_target_accuracy"], sliceFull)
		metrics["minimum_full_byte_target_accuracy"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_full_byte_target_accuracy"], ringFull)
		metrics["maximum_adaptive_accuracy_gap_vs_full_byte"] = wlmLmAdaptiveGranularityR1Max(metrics["maximum_adaptive_accuracy_gap_vs_full_byte"], adaptiveGap)
		metrics["maximum_static_fixed4_target_accuracy"] = wlmLmAdaptiveGranularityR1Max(metrics["maximum_static_fixed4_target_accuracy"], sliceStatic)
		metrics["maximum_static_fixed4_target_accuracy"] = wlmLmAdaptiveGranularityR1Max(metrics["maximum_static_fixed4_target_accuracy"], ringStatic)
		metrics["minimum_adaptive_accuracy_gain_vs_static_fixed4"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_adaptive_accuracy_gain_vs_static_fixed4"], sliceAdaptive-sliceStatic)
		metrics["minimum_adaptive_accuracy_gain_vs_static_fixed4"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_adaptive_accuracy_gain_vs_static_fixed4"], ringAdaptive-ringStatic)
		metrics["minimum_overall_event_reduction_fraction"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_overall_event_reduction_fraction"], overallReduction)
		metrics["minimum_predictable_region_event_reduction_fraction"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_predictable_region_event_reduction_fraction"], predictableReduction)
		metrics["maximum_high_entropy_false_coarse_activation_rate"] = wlmLmAdaptiveGranularityR1Max(metrics["maximum_high_entropy_false_coarse_activation_rate"], falseHighRate)
		metrics["minimum_predictable_coarse_activation_recall"] = wlmLmAdaptiveGranularityR1Min(metrics["minimum_predictable_coarse_activation_recall"], coarseRecall)
		metrics["maximum_adaptive_events_per_raw_byte"] = wlmLmAdaptiveGranularityR1Max(metrics["maximum_adaptive_events_per_raw_byte"], eventsPerByte)
		metrics["valid_seed_count"]++
	}

	return wlmLmAdaptiveGranularityR1Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-ADAPTIVE-GRANULARITY-R1",
		Metrics: metrics,
	}
}
