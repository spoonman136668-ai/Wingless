package unitary

import (
	"fmt"
	"math"
	"sort"
)

type wlmLmRawRepContextualAliasInductionR2Context [4]uint8

type wlmLmRawRepContextualAliasInductionR2Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

type wlmLmRawRepContextualAliasInductionR2Classes struct {
	surfaceToClass []int
	members        [][]int
	decodeFailures int
}

type wlmLmRawRepContextualAliasInductionR2Transition struct {
	counts [5][4][4][16]uint32
}

func wlmLmRawRepContextualAliasInductionR2SurfaceKey(family, entity int) wlmSiRawRepAliasInvarianceFalsificationR1Key {
	if entity < 4 {
		return wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family, entity)
	}
	return wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(family, entity-4)
}

func wlmLmRawRepContextualAliasInductionR2BridgeEntity(family, entity int, shuffled bool) int {
	if !shuffled || family == 0 {
		return entity
	}
	if entity < 4 {
		return (entity + 2) % 4
	}
	return 4 + ((entity - 4 + 2) % 5)
}

func wlmLmRawRepContextualAliasInductionR2BridgeRecord(family, entity, contextIndex, repeat int, shuffled bool) []uint8 {
	assigned := wlmLmRawRepContextualAliasInductionR2BridgeEntity(family, entity, shuffled)
	state := uint32(0x9e3779b9 ^ uint32((family+1)*10007+(entity+1)*1009+(contextIndex+1)*97+(repeat+1)*65537))
	out := make([]uint8, 0, 18)
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
	out = append(out, uint8(192+assigned), uint8(208+contextIndex))
	key := wlmLmRawRepContextualAliasInductionR2SurfaceKey(family, entity)
	out = append(out, key[:]...)
	out = append(out, uint8(224+assigned), uint8(240+contextIndex))
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
	return out
}

func wlmLmRawRepContextualAliasInductionR2TaskRecords() [][]uint8 {
	records := make([][]uint8, 0, 3200)
	for repeat := 0; repeat < 20; repeat++ {
		for family := 0; family < 2; family++ {
			for op := 0; op < 5; op++ {
				for left := 0; left < 4; left++ {
					for right := 0; right < 4; right++ {
						records = append(records, wlmSiRawRepAliasInvarianceFalsificationR1TrainingRecord(family, left, right, op, repeat))
					}
				}
			}
		}
	}
	return records
}

func wlmLmRawRepContextualAliasInductionR2BridgeRecords(shuffled bool) [][]uint8 {
	records := make([][]uint8, 0, 576)
	for entity := 0; entity < 9; entity++ {
		for family := 0; family < 2; family++ {
			for contextIndex := 0; contextIndex < 4; contextIndex++ {
				for repeat := 0; repeat < 8; repeat++ {
					records = append(records, wlmLmRawRepContextualAliasInductionR2BridgeRecord(family, entity, contextIndex, repeat, shuffled))
				}
			}
		}
	}
	return records
}

func wlmLmRawRepContextualAliasInductionR2Signature(counts map[wlmLmRawRepContextualAliasInductionR2Context]uint32) string {
	keys := make([]wlmLmRawRepContextualAliasInductionR2Context, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		for k := 0; k < 4; k++ {
			if keys[i][k] != keys[j][k] {
				return keys[i][k] < keys[j][k]
			}
		}
		return false
	})
	out := ""
	for _, key := range keys {
		out += fmt.Sprintf("%02x%02x%02x%02x:%08x;", key[0], key[1], key[2], key[3], counts[key])
	}
	return out
}

func wlmLmRawRepContextualAliasInductionR2Induce(
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	bridge [][]uint8,
) wlmLmRawRepContextualAliasInductionR2Classes {
	repIndex := make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int, len(reps))
	for id, key := range reps {
		repIndex[key] = id
	}
	signatures := make([]map[wlmLmRawRepContextualAliasInductionR2Context]uint32, len(reps))
	for i := range signatures {
		signatures[i] = make(map[wlmLmRawRepContextualAliasInductionR2Context]uint32)
	}
	decodeFailures := 0
	for _, data := range bridge {
		matches := 0
		for i := 2; i+6 <= len(data); i++ {
			key := wlmSiRawRepAliasInvarianceFalsificationR1Key{data[i], data[i+1], data[i+2], data[i+3]}
			repID, ok := repIndex[key]
			if !ok {
				continue
			}
			context := wlmLmRawRepContextualAliasInductionR2Context{data[i-2], data[i-1], data[i+4], data[i+5]}
			signatures[repID][context]++
			matches++
			i += 3
		}
		if matches != 1 {
			decodeFailures++
		}
	}

	grouped := make(map[string][]int)
	for surfaceID, counts := range signatures {
		signature := wlmLmRawRepContextualAliasInductionR2Signature(counts)
		grouped[signature] = append(grouped[signature], surfaceID)
	}
	signatureKeys := make([]string, 0, len(grouped))
	for signature := range grouped {
		signatureKeys = append(signatureKeys, signature)
	}
	sort.Strings(signatureKeys)

	surfaceToClass := make([]int, len(reps))
	for i := range surfaceToClass {
		surfaceToClass[i] = -1
	}
	members := make([][]int, 0, len(signatureKeys))
	for classID, signature := range signatureKeys {
		group := append([]int(nil), grouped[signature]...)
		sort.Ints(group)
		members = append(members, group)
		for _, surfaceID := range group {
			surfaceToClass[surfaceID] = classID
		}
	}
	return wlmLmRawRepContextualAliasInductionR2Classes{
		surfaceToClass: surfaceToClass,
		members:        members,
		decodeFailures: decodeFailures,
	}
}

func wlmLmRawRepContextualAliasInductionR2IndexMap(ids []int) map[int]int {
	out := make(map[int]int, len(ids))
	for dense, id := range ids {
		out[id] = dense
	}
	return out
}

func (t *wlmLmRawRepContextualAliasInductionR2Transition) observe(op, left, right, outLeft, outRight int, metrics map[string]float64) {
	if op < 0 || op >= 5 || left < 0 || left >= 4 || right < 0 || right >= 4 || outLeft < 0 || outLeft >= 4 || outRight < 0 || outRight >= 4 {
		metrics["invalid_row_count"]++
		return
	}
	index := outLeft*4 + outRight
	if t.counts[op][left][right][index] == ^uint32(0) {
		metrics["counter_overflow_count"]++
		return
	}
	t.counts[op][left][right][index]++
}

func (t *wlmLmRawRepContextualAliasInductionR2Transition) predict(op, left, right int) (int, int) {
	if op < 0 || op >= 5 || left < 0 || left >= 4 || right < 0 || right >= 4 {
		return 0, 0
	}
	best := 0
	bestCount := t.counts[op][left][right][0]
	for index := 1; index < 16; index++ {
		if t.counts[op][left][right][index] > bestCount {
			best = index
			bestCount = t.counts[op][left][right][index]
		}
	}
	return best / 4, best % 4
}

func (t *wlmLmRawRepContextualAliasInductionR2Transition) missingCount() int {
	missing := 0
	for op := 0; op < 5; op++ {
		for left := 0; left < 4; left++ {
			for right := 0; right < 4; right++ {
				var total uint64
				for out := 0; out < 16; out++ {
					total += uint64(t.counts[op][left][right][out])
				}
				if total == 0 {
					missing++
				}
			}
		}
	}
	return missing
}

type wlmLmRawRepContextualAliasInductionR2ClassModel struct {
	transition  wlmLmRawRepContextualAliasInductionR2Transition
	valueIDs    []int
	opIDs       []int
	valueMap    map[int]int
	opMap       map[int]int
	decodeFails int
	valid       bool
}

func wlmLmRawRepContextualAliasInductionR2BuildClassModel(
	records [][]uint8,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	metrics map[string]float64,
) wlmLmRawRepContextualAliasInductionR2ClassModel {
	valueSet := make(map[int]bool)
	opSet := make(map[int]bool)
	decodedRows := make([][]int, 0, len(records))
	decodeFails := 0

	for _, raw := range records {
		surface := wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw, reps)
		if len(surface) != 5 {
			decodeFails++
			continue
		}
		row := make([]int, 5)
		valid := true
		for i, surfaceID := range surface {
			if surfaceID < 0 || surfaceID >= len(classes.surfaceToClass) {
				valid = false
				break
			}
			classID := classes.surfaceToClass[surfaceID]
			if classID < 0 {
				valid = false
				break
			}
			row[i] = classID
		}
		if !valid {
			decodeFails++
			continue
		}
		valueSet[row[0]] = true
		valueSet[row[1]] = true
		opSet[row[2]] = true
		valueSet[row[3]] = true
		valueSet[row[4]] = true
		decodedRows = append(decodedRows, row)
	}

	valueIDs := make([]int, 0, len(valueSet))
	for id := range valueSet {
		valueIDs = append(valueIDs, id)
	}
	opIDs := make([]int, 0, len(opSet))
	for id := range opSet {
		opIDs = append(opIDs, id)
	}
	sort.Ints(valueIDs)
	sort.Ints(opIDs)
	valueMap := wlmLmRawRepContextualAliasInductionR2IndexMap(valueIDs)
	opMap := wlmLmRawRepContextualAliasInductionR2IndexMap(opIDs)

	model := wlmLmRawRepContextualAliasInductionR2ClassModel{
		valueIDs: valueIDs,
		opIDs: opIDs,
		valueMap: valueMap,
		opMap: opMap,
		decodeFails: decodeFails,
		valid: len(valueIDs) == 4 && len(opIDs) == 5,
	}
	if !model.valid {
		metrics["invalid_row_count"]++
		return model
	}
	for _, row := range decodedRows {
		left, ok0 := valueMap[row[0]]
		right, ok1 := valueMap[row[1]]
		op, ok2 := opMap[row[2]]
		outLeft, ok3 := valueMap[row[3]]
		outRight, ok4 := valueMap[row[4]]
		if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 {
			metrics["invalid_row_count"]++
			continue
		}
		model.transition.observe(op, left, right, outLeft, outRight, metrics)
	}
	return model
}

func wlmLmRawRepContextualAliasInductionR2RepIndex(reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int {
	out := make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int, len(reps))
	for id, key := range reps {
		out[key] = id
	}
	return out
}

func wlmLmRawRepContextualAliasInductionR2ExpectedDenseValue(
	family, semantic int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	model wlmLmRawRepContextualAliasInductionR2ClassModel,
) (int, bool) {
	surfaceID, ok := repIndex[wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family, semantic)]
	if !ok || surfaceID < 0 || surfaceID >= len(classes.surfaceToClass) {
		return 0, false
	}
	classID := classes.surfaceToClass[surfaceID]
	dense, ok := model.valueMap[classID]
	return dense, ok
}

func wlmLmRawRepContextualAliasInductionR2EvaluateClass(
	valueFamily, opFamily, leftSemantic, rightSemantic, opSemantic, tag int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	model wlmLmRawRepContextualAliasInductionR2ClassModel,
	metrics map[string]float64,
) bool {
	raw := wlmSiRawRepAliasInvarianceFalsificationR1EvalRecord(valueFamily, opFamily, leftSemantic, rightSemantic, opSemantic, tag)
	surface := wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw, reps)
	if len(surface) != 3 {
		metrics["invalid_row_count"]++
		return false
	}
	classLeft := classes.surfaceToClass[surface[0]]
	classRight := classes.surfaceToClass[surface[1]]
	classOp := classes.surfaceToClass[surface[2]]
	left, ok0 := model.valueMap[classLeft]
	right, ok1 := model.valueMap[classRight]
	op, ok2 := model.opMap[classOp]
	if !ok0 || !ok1 || !ok2 {
		metrics["invalid_row_count"]++
		return false
	}
	predLeft, predRight := model.transition.predict(op, left, right)
	targetLeftSemantic, targetRightSemantic := wlmSiRawRepAliasInvarianceFalsificationR1Apply(leftSemantic, rightSemantic, opSemantic)
	targetLeft, ok3 := wlmLmRawRepContextualAliasInductionR2ExpectedDenseValue(valueFamily, targetLeftSemantic, reps, repIndex, classes, model)
	targetRight, ok4 := wlmLmRawRepContextualAliasInductionR2ExpectedDenseValue(valueFamily, targetRightSemantic, reps, repIndex, classes, model)
	if !ok3 || !ok4 {
		metrics["invalid_row_count"]++
		return false
	}
	return predLeft == targetLeft && predRight == targetRight
}

type wlmLmRawRepContextualAliasInductionR2SurfaceModel struct {
	transition wlmSiRawRepAliasInvarianceFalsificationR1Transition
	valueIDs   []int
	opIDs      []int
	valueMap   map[int]int
	opMap      map[int]int
	valid      bool
}

func wlmLmRawRepContextualAliasInductionR2BuildSurfaceModel(
	records [][]uint8,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	metrics map[string]float64,
) wlmLmRawRepContextualAliasInductionR2SurfaceModel {
	valueSet := make(map[int]bool)
	opSet := make(map[int]bool)
	decoded := make([][]int, 0, len(records))
	for _, raw := range records {
		row := wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw, reps)
		if len(row) != 5 {
			metrics["invalid_row_count"]++
			continue
		}
		valueSet[row[0]] = true
		valueSet[row[1]] = true
		opSet[row[2]] = true
		valueSet[row[3]] = true
		valueSet[row[4]] = true
		decoded = append(decoded, row)
	}
	valueIDs := make([]int, 0, len(valueSet))
	for id := range valueSet {
		valueIDs = append(valueIDs, id)
	}
	opIDs := make([]int, 0, len(opSet))
	for id := range opSet {
		opIDs = append(opIDs, id)
	}
	sort.Ints(valueIDs)
	sort.Ints(opIDs)
	valueMap := wlmSiRawRepAliasInvarianceFalsificationR1IndexMap(valueIDs)
	opMap := wlmSiRawRepAliasInvarianceFalsificationR1IndexMap(opIDs)
	model := wlmLmRawRepContextualAliasInductionR2SurfaceModel{
		valueIDs: valueIDs,
		opIDs: opIDs,
		valueMap: valueMap,
		opMap: opMap,
		valid: len(valueIDs) == 8 && len(opIDs) == 10,
	}
	if !model.valid {
		metrics["invalid_row_count"]++
		return model
	}
	for _, row := range decoded {
		left, ok0 := valueMap[row[0]]
		right, ok1 := valueMap[row[1]]
		op, ok2 := opMap[row[2]]
		outLeft, ok3 := valueMap[row[3]]
		outRight, ok4 := valueMap[row[4]]
		if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 {
			metrics["invalid_row_count"]++
			continue
		}
		model.transition.observe(op, left, right, outLeft, outRight, metrics)
	}
	return model
}

func wlmLmRawRepContextualAliasInductionR2EvaluateSurfaceCross(
	valueFamily, opFamily, leftSemantic, rightSemantic, opSemantic, tag int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	model wlmLmRawRepContextualAliasInductionR2SurfaceModel,
	metrics map[string]float64,
) (bool, bool) {
	raw := wlmSiRawRepAliasInvarianceFalsificationR1EvalRecord(valueFamily, opFamily, leftSemantic, rightSemantic, opSemantic, tag)
	row := wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw, reps)
	if len(row) != 3 {
		metrics["invalid_row_count"]++
		return false, false
	}
	left, ok0 := model.valueMap[row[0]]
	right, ok1 := model.valueMap[row[1]]
	op, ok2 := model.opMap[row[2]]
	if !ok0 || !ok1 || !ok2 {
		metrics["invalid_row_count"]++
		return false, false
	}
	missing := !model.transition.has(op, left, right)
	predLeft, predRight := model.transition.predict(op, left, right)
	targetLeftSemantic, targetRightSemantic := wlmSiRawRepAliasInvarianceFalsificationR1Apply(leftSemantic, rightSemantic, opSemantic)
	targetLeftSurface, ok3 := repIndex[wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(valueFamily, targetLeftSemantic)]
	targetRightSurface, ok4 := repIndex[wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(valueFamily, targetRightSemantic)]
	if !ok3 || !ok4 {
		metrics["invalid_row_count"]++
		return false, missing
	}
	targetLeft, ok5 := model.valueMap[targetLeftSurface]
	targetRight, ok6 := model.valueMap[targetRightSurface]
	if !ok5 || !ok6 {
		metrics["invalid_row_count"]++
		return false, missing
	}
	return predLeft == targetLeft && predRight == targetRight, missing
}

func wlmLmRawRepContextualAliasInductionR2CorrectAliasPairs(
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
) int {
	count := 0
	for semantic := 0; semantic < 4; semantic++ {
		a, ok0 := repIndex[wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(0, semantic)]
		b, ok1 := repIndex[wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(1, semantic)]
		if ok0 && ok1 && classes.surfaceToClass[a] == classes.surfaceToClass[b] {
			count++
		}
	}
	for semantic := 0; semantic < 5; semantic++ {
		a, ok0 := repIndex[wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(0, semantic)]
		b, ok1 := repIndex[wlmSiRawRepAliasInvarianceFalsificationR1OpMotif(1, semantic)]
		if ok0 && ok1 && classes.surfaceToClass[a] == classes.surfaceToClass[b] {
			count++
		}
	}
	return count
}

// RunWlmLmRawRepContextualAliasInductionR2 tests distributional raw-context alias induction.
func RunWlmLmRawRepContextualAliasInductionR2() interface{} {
	metrics := map[string]float64{
		"training_record_count":                                  0,
		"bridge_record_count":                                    0,
		"selected_surface_representation_count":                  0,
		"selected_true_surface_motif_match_count":                0,
		"bridge_decode_failure_count":                            0,
		"induced_alias_class_count":                              0,
		"induced_two_member_class_count":                         0,
		"correct_evaluator_alias_pair_count":                     0,
		"task_training_decode_failure_count":                     0,
		"class_transition_missing_count":                         0,
		"same_family_evaluation_count":                           0,
		"induced_same_family_accuracy":                           0,
		"cross_family_evaluation_count":                          0,
		"induced_cross_family_accuracy":                          0,
		"exact_surface_cross_family_accuracy":                    0,
		"exact_surface_cross_family_missing_transition_count":    0,
		"shuffled_bridge_alias_class_count":                      0,
		"shuffled_bridge_cross_family_accuracy":                  0,
		"tokenizer_use_count":                                    0,
		"external_model_call_count":                              0,
		"capacity_growth_event_count":                            0,
		"invalid_row_count":                                      0,
		"counter_overflow_count":                                 0,
	}

	taskRecords := wlmLmRawRepContextualAliasInductionR2TaskRecords()
	metrics["training_record_count"] = float64(len(taskRecords))
	reps := wlmSiRawRepAliasInvarianceFalsificationR1LearnRepresentation(taskRecords)
	metrics["selected_surface_representation_count"] = float64(len(reps))
	metrics["selected_true_surface_motif_match_count"] = float64(wlmSiRawRepAliasInvarianceFalsificationR1TrueMotifCount(reps))
	repIndex := wlmLmRawRepContextualAliasInductionR2RepIndex(reps)

	bridge := wlmLmRawRepContextualAliasInductionR2BridgeRecords(false)
	metrics["bridge_record_count"] = float64(len(bridge))
	classes := wlmLmRawRepContextualAliasInductionR2Induce(reps, bridge)
	metrics["bridge_decode_failure_count"] = float64(classes.decodeFailures)
	metrics["induced_alias_class_count"] = float64(len(classes.members))
	for _, members := range classes.members {
		if len(members) == 2 {
			metrics["induced_two_member_class_count"]++
		}
	}
	metrics["correct_evaluator_alias_pair_count"] = float64(wlmLmRawRepContextualAliasInductionR2CorrectAliasPairs(reps, repIndex, classes))

	classModel := wlmLmRawRepContextualAliasInductionR2BuildClassModel(taskRecords, reps, classes, metrics)
	metrics["task_training_decode_failure_count"] = float64(classModel.decodeFails)
	if classModel.valid {
		metrics["class_transition_missing_count"] = float64(classModel.transition.missingCount())
	}

	surfaceModel := wlmLmRawRepContextualAliasInductionR2BuildSurfaceModel(taskRecords, reps, metrics)

	shuffledBridge := wlmLmRawRepContextualAliasInductionR2BridgeRecords(true)
	shuffledClasses := wlmLmRawRepContextualAliasInductionR2Induce(reps, shuffledBridge)
	metrics["shuffled_bridge_alias_class_count"] = float64(len(shuffledClasses.members))
	if shuffledClasses.decodeFailures != 0 {
		metrics["invalid_row_count"] += float64(shuffledClasses.decodeFailures)
	}
	shuffledModel := wlmLmRawRepContextualAliasInductionR2BuildClassModel(taskRecords, reps, shuffledClasses, metrics)
	if shuffledModel.decodeFails != 0 {
		metrics["invalid_row_count"] += float64(shuffledModel.decodeFails)
	}

	sameHits := 0
	sameTotal := 0
	for family := 0; family < 2; family++ {
		for op := 0; op < 5; op++ {
			for left := 0; left < 4; left++ {
				for right := 0; right < 4; right++ {
					if wlmLmRawRepContextualAliasInductionR2EvaluateClass(
						family, family, left, right, op,
						1000+family*100+op*20+left*4+right,
						reps, repIndex, classes, classModel, metrics,
					) {
						sameHits++
					}
					sameTotal++
				}
			}
		}
	}
	metrics["same_family_evaluation_count"] = float64(sameTotal)
	if sameTotal > 0 {
		metrics["induced_same_family_accuracy"] = float64(sameHits) / float64(sameTotal)
	}

	crossHits := 0
	exactHits := 0
	exactMissing := 0
	shuffledHits := 0
	crossTotal := 0
	for valueFamily := 0; valueFamily < 2; valueFamily++ {
		opFamily := 1 - valueFamily
		for op := 0; op < 5; op++ {
			for left := 0; left < 4; left++ {
				for right := 0; right < 4; right++ {
					tag := 2000 + valueFamily*100 + op*20 + left*4 + right
					if wlmLmRawRepContextualAliasInductionR2EvaluateClass(
						valueFamily, opFamily, left, right, op, tag,
						reps, repIndex, classes, classModel, metrics,
					) {
						crossHits++
					}
					if ok, missing := wlmLmRawRepContextualAliasInductionR2EvaluateSurfaceCross(
						valueFamily, opFamily, left, right, op, tag,
						reps, repIndex, surfaceModel, metrics,
					); ok {
						exactHits++
						if missing {
							exactMissing++
						}
					} else if missing {
						exactMissing++
					}
					if wlmLmRawRepContextualAliasInductionR2EvaluateClass(
						valueFamily, opFamily, left, right, op, tag,
						reps, repIndex, shuffledClasses, shuffledModel, metrics,
					) {
						shuffledHits++
					}
					crossTotal++
				}
			}
		}
	}
	metrics["cross_family_evaluation_count"] = float64(crossTotal)
	if crossTotal > 0 {
		metrics["induced_cross_family_accuracy"] = float64(crossHits) / float64(crossTotal)
		metrics["exact_surface_cross_family_accuracy"] = float64(exactHits) / float64(crossTotal)
		metrics["shuffled_bridge_cross_family_accuracy"] = float64(shuffledHits) / float64(crossTotal)
	}
	metrics["exact_surface_cross_family_missing_transition_count"] = float64(exactMissing)

	for _, value := range metrics {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			metrics["invalid_row_count"]++
		}
	}
	return wlmLmRawRepContextualAliasInductionR2Result{
		Schema:     "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-CONTEXTUAL-ALIAS-INDUCTION-R2",
		Metrics:    metrics,
	}
}
