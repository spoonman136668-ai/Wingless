package unitary

import "math"

type wlmLmRawRepZeroshotAliasFamilyR5Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmRawRepZeroshotAliasFamilyR5BridgeRecord(family, entity, contextIndex, repeat int, shuffled bool) []uint8 {
	assigned := wlmLmRawRepZeroshotAliasFamilyR1BridgeEntity(family, entity, shuffled)
	state := uint32(0x7f4a7c15 ^ uint32((family+1)*10007+(entity+1)*1009+(contextIndex+1)*97+(repeat+1)*65537))
	out := make([]uint8, 0, 18)
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
	out = append(out,
		uint8(192+((assigned*7+contextIndex*3+repeat)%16)),
		uint8(208+((assigned*5+contextIndex+2*repeat)%16)),
	)
	key := wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity)
	out = append(out, key[:]...)
	out = append(out,
		uint8(224+((assigned*11+2*contextIndex+repeat)%16)),
		uint8(240+((assigned*13+contextIndex+3*repeat)%16)),
	)
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
	return out
}

func wlmLmRawRepZeroshotAliasFamilyR5BridgeRecords(shuffled bool) [][]uint8 {
	records := make([][]uint8, 0, 864)
	for entity := 0; entity < 9; entity++ {
		for family := 0; family < 3; family++ {
			for contextIndex := 0; contextIndex < 4; contextIndex++ {
				for repeat := 0; repeat < 8; repeat++ {
					records = append(records, wlmLmRawRepZeroshotAliasFamilyR5BridgeRecord(family, entity, contextIndex, repeat, shuffled))
				}
			}
		}
	}
	return records
}

func wlmLmRawRepZeroshotAliasFamilyR5ConstructionMetrics(records [][]uint8) (float64, float64) {
	counts := make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]uint32)
	for _, data := range records {
		for i := 0; i+4 <= len(data); i++ {
			key := wlmSiRawRepAliasInvarianceFalsificationR1Key{data[i], data[i+1], data[i+2], data[i+3]}
			counts[key]++
		}
	}
	trueSet := make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]bool, 27)
	for family := 0; family < 3; family++ {
		for entity := 0; entity < 9; entity++ {
			trueSet[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity)] = true
		}
	}
	minTrue := uint32(^uint32(0))
	maxNonTrue := uint32(0)
	for key, count := range counts {
		if trueSet[key] {
			if count < minTrue {
				minTrue = count
			}
		} else if count > maxNonTrue {
			maxNonTrue = count
		}
	}
	if minTrue == ^uint32(0) {
		minTrue = 0
	}
	return float64(minTrue), float64(maxNonTrue)
}

// RunWlmLmRawRepZeroshotAliasFamilyR5 tests repaired context-only transfer to an unseen surface family.
func RunWlmLmRawRepZeroshotAliasFamilyR5() interface{} {
	metrics := map[string]float64{
		"bridge_record_count":                              0,
		"minimum_true_surface_motif_occurrence_count":      0,
		"maximum_nontrue_bridge_window_occurrence_count":   0,
		"selected_surface_representation_count":            0,
		"selected_true_surface_motif_match_count":          0,
		"bridge_decode_failure_count":                      0,
		"induced_alias_class_count":                        0,
		"induced_three_member_class_count":                 0,
		"correct_evaluator_alias_triplet_count":            0,
		"task_training_record_count":                       0,
		"family2_task_training_record_count":               0,
		"task_training_decode_failure_count":               0,
		"class_transition_key_count":                       0,
		"unseen_family_evaluation_count":                   0,
		"unseen_family_accuracy":                           0,
		"exact_surface_unseen_family_accuracy":             0,
		"exact_surface_unseen_missing_transition_count":    0,
		"shuffled_family2_alias_class_count":               0,
		"shuffled_family2_unseen_accuracy":                 0,
		"tokenizer_use_count":                              0,
		"external_model_call_count":                        0,
		"capacity_growth_event_count":                      0,
		"invalid_row_count":                                0,
		"counter_overflow_count":                           0,
	}

	bridge := wlmLmRawRepZeroshotAliasFamilyR5BridgeRecords(false)
	metrics["bridge_record_count"] = float64(len(bridge))
	minTrue, maxNonTrue := wlmLmRawRepZeroshotAliasFamilyR5ConstructionMetrics(bridge)
	metrics["minimum_true_surface_motif_occurrence_count"] = minTrue
	metrics["maximum_nontrue_bridge_window_occurrence_count"] = maxNonTrue

	reps := wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(bridge)
	metrics["selected_surface_representation_count"] = float64(len(reps))
	metrics["selected_true_surface_motif_match_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps))
	repIndex := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)

	classes := wlmLmRawRepContextualAliasInductionR2Induce(reps, bridge)
	metrics["bridge_decode_failure_count"] = float64(classes.decodeFailures)
	metrics["induced_alias_class_count"] = float64(len(classes.members))
	for _, members := range classes.members {
		if len(members) == 3 {
			metrics["induced_three_member_class_count"]++
		}
	}
	metrics["correct_evaluator_alias_triplet_count"] = float64(wlmLmRawRepZeroshotAliasFamilyR1CorrectTriplets(repIndex, classes))

	taskRecords := wlmLmRawRepContextualAliasInductionR2TaskRecords()
	metrics["task_training_record_count"] = float64(len(taskRecords))
	classModel := wlmLmRawRepContextualAliasInductionR2BuildClassModel(taskRecords, reps, classes, metrics)
	metrics["task_training_decode_failure_count"] = float64(classModel.decodeFails)
	if classModel.valid {
		metrics["class_transition_key_count"] = float64(80 - classModel.transition.missingCount())
	}

	surfaceModel, surfaceDecodeFails := wlmLmRawRepZeroshotAliasFamilyR1BuildSurfaceModel(taskRecords, reps, metrics)
	if surfaceDecodeFails != 0 {
		metrics["invalid_row_count"] += float64(surfaceDecodeFails)
	}

	shuffledBridge := wlmLmRawRepZeroshotAliasFamilyR5BridgeRecords(true)
	shuffledClasses := wlmLmRawRepContextualAliasInductionR2Induce(reps, shuffledBridge)
	metrics["shuffled_family2_alias_class_count"] = float64(len(shuffledClasses.members))
	if shuffledClasses.decodeFailures != 0 {
		metrics["invalid_row_count"] += float64(shuffledClasses.decodeFailures)
	}
	shuffledModel := wlmLmRawRepContextualAliasInductionR2BuildClassModel(taskRecords, reps, shuffledClasses, metrics)
	if shuffledModel.decodeFails != 0 {
		metrics["invalid_row_count"] += float64(shuffledModel.decodeFails)
	}

	mainHits := 0
	exactHits := 0
	exactMissing := 0
	shuffledHits := 0
	total := 0
	for op := 0; op < 5; op++ {
		for left := 0; left < 4; left++ {
			for right := 0; right < 4; right++ {
				tag := 3000 + op*20 + left*4 + right
				if wlmLmRawRepZeroshotAliasFamilyR1EvaluateClass(
					left, right, op, tag, reps, repIndex, classes, classModel, metrics,
				) {
					mainHits++
				}
				if ok, missing := wlmLmRawRepZeroshotAliasFamilyR1EvaluateSurface(
					left, right, op, tag, reps, surfaceModel, metrics,
				); ok {
					exactHits++
					if missing {
						exactMissing++
					}
				} else if missing {
					exactMissing++
				}
				if wlmLmRawRepZeroshotAliasFamilyR1EvaluateClass(
					left, right, op, tag, reps, repIndex, shuffledClasses, shuffledModel, metrics,
				) {
					shuffledHits++
				}
				total++
			}
		}
	}
	metrics["unseen_family_evaluation_count"] = float64(total)
	if total > 0 {
		metrics["unseen_family_accuracy"] = float64(mainHits) / float64(total)
		metrics["exact_surface_unseen_family_accuracy"] = float64(exactHits) / float64(total)
		metrics["shuffled_family2_unseen_accuracy"] = float64(shuffledHits) / float64(total)
	}
	metrics["exact_surface_unseen_missing_transition_count"] = float64(exactMissing)

	for _, value := range metrics {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			metrics["invalid_row_count"]++
		}
	}

	return wlmLmRawRepZeroshotAliasFamilyR5Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-ZEROSHOT-ALIAS-FAMILY-R5",
		Metrics: metrics,
	}
}
