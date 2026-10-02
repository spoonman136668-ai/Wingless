package unitary

import (
	"math"
	"sort"
)

type wlmLmRawRepZeroshotAliasFamilyR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

type wlmLmRawRepZeroshotAliasFamilyR1SurfaceTransitionKey struct {
	op, left, right int
}

type wlmLmRawRepZeroshotAliasFamilyR1SurfacePair struct {
	left, right int
}

type wlmLmRawRepZeroshotAliasFamilyR1SurfaceModel struct {
	counts map[wlmLmRawRepZeroshotAliasFamilyR1SurfaceTransitionKey]map[wlmLmRawRepZeroshotAliasFamilyR1SurfacePair]uint32
}

func wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(id int) wlmSiRawRepAliasInvarianceFalsificationR1Key {
	switch id {
	case 0:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{0, 1, 2, 3}
	case 1:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{4, 5, 6, 7}
	case 2:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{8, 9, 10, 11}
	default:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{12, 13, 14, 15}
	}
}

func wlmLmRawRepZeroshotAliasFamilyR1Family2OpMotif(id int) wlmSiRawRepAliasInvarianceFalsificationR1Key {
	switch id {
	case 0:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{44, 45, 46, 47}
	case 1:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{48, 49, 50, 51}
	case 2:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{52, 53, 54, 55}
	case 3:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{84, 85, 86, 87}
	default:
		return wlmSiRawRepAliasInvarianceFalsificationR1Key{88, 89, 90, 91}
	}
}

func wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity int) wlmSiRawRepAliasInvarianceFalsificationR1Key {
	if family < 2 {
		return wlmLmRawRepContextualAliasInductionR2SurfaceKey(family, entity)
	}
	if entity < 4 {
		return wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(entity)
	}
	return wlmLmRawRepZeroshotAliasFamilyR1Family2OpMotif(entity - 4)
}

func wlmLmRawRepZeroshotAliasFamilyR1BridgeEntity(family, entity int, shuffled bool) int {
	if !shuffled || family < 2 {
		return entity
	}
	if entity < 4 {
		return (entity + 3) % 4
	}
	return 4 + ((entity - 4 + 1) % 5)
}

func wlmLmRawRepZeroshotAliasFamilyR1BridgeRecord(family, entity, contextIndex, repeat int, shuffled bool) []uint8 {
	assigned := wlmLmRawRepZeroshotAliasFamilyR1BridgeEntity(family, entity, shuffled)
	state := uint32(0x7f4a7c15 ^ uint32((family+1)*10007+(entity+1)*1009+(contextIndex+1)*97+(repeat+1)*65537))
	out := make([]uint8, 0, 18)
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
	out = append(out, uint8(192+assigned), uint8(208+contextIndex))
	key := wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity)
	out = append(out, key[:]...)
	out = append(out, uint8(224+assigned), uint8(240+contextIndex))
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 3)...)
	return out
}

func wlmLmRawRepZeroshotAliasFamilyR1BridgeRecords(shuffled bool) [][]uint8 {
	records := make([][]uint8, 0, 864)
	for entity := 0; entity < 9; entity++ {
		for family := 0; family < 3; family++ {
			for contextIndex := 0; contextIndex < 4; contextIndex++ {
				for repeat := 0; repeat < 8; repeat++ {
					records = append(records, wlmLmRawRepZeroshotAliasFamilyR1BridgeRecord(family, entity, contextIndex, repeat, shuffled))
				}
			}
		}
	}
	return records
}

func wlmLmRawRepZeroshotAliasFamilyR1Less(a, b wlmSiRawRepAliasInvarianceFalsificationR1Key) bool {
	for i := 0; i < 4; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

func wlmLmRawRepZeroshotAliasFamilyR1SelectSurface(records [][]uint8) []wlmSiRawRepAliasInvarianceFalsificationR1Key {
	counts := make(map[wlmSiRawRepAliasInvarianceFalsificationR1Key]uint32)
	for _, data := range records {
		for i := 0; i+4 <= len(data); i++ {
			key := wlmSiRawRepAliasInvarianceFalsificationR1Key{data[i], data[i+1], data[i+2], data[i+3]}
			counts[key]++
		}
	}
	type row struct {
		key   wlmSiRawRepAliasInvarianceFalsificationR1Key
		count uint32
	}
	rows := make([]row, 0, len(counts))
	for key, count := range counts {
		rows = append(rows, row{key: key, count: count})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		return wlmLmRawRepZeroshotAliasFamilyR1Less(rows[i].key, rows[j].key)
	})
	if len(rows) > 27 {
		rows = rows[:27]
	}
	out := make([]wlmSiRawRepAliasInvarianceFalsificationR1Key, len(rows))
	for i, row := range rows {
		out[i] = row.key
	}
	sort.Slice(out, func(i, j int) bool {
		return wlmLmRawRepZeroshotAliasFamilyR1Less(out[i], out[j])
	})
	return out
}

func wlmLmRawRepZeroshotAliasFamilyR1TrueMotifCount(reps []wlmSiRawRepAliasInvarianceFalsificationR1Key) int {
	count := 0
	for _, rep := range reps {
		matched := false
		for family := 0; family < 3; family++ {
			for entity := 0; entity < 9; entity++ {
				if rep == wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity) {
					matched = true
				}
			}
		}
		if matched {
			count++
		}
	}
	return count
}

func wlmLmRawRepZeroshotAliasFamilyR1CorrectTriplets(
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
) int {
	count := 0
	for entity := 0; entity < 9; entity++ {
		ids := [3]int{-1, -1, -1}
		ok := true
		for family := 0; family < 3; family++ {
			id, exists := repIndex[wlmLmRawRepZeroshotAliasFamilyR1SurfaceKey(family, entity)]
			if !exists {
				ok = false
				break
			}
			ids[family] = id
		}
		if ok &&
			classes.surfaceToClass[ids[0]] == classes.surfaceToClass[ids[1]] &&
			classes.surfaceToClass[ids[1]] == classes.surfaceToClass[ids[2]] {
			count++
		}
	}
	return count
}

func wlmLmRawRepZeroshotAliasFamilyR1EvalRecord(left, right, op, tag int) []uint8 {
	state := uint32(0x94d049bb ^ uint32((left+1)*131+(right+1)*977+(op+1)*8191+(tag+1)*65537))
	out := make([]uint8, 0, 60)
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 6+(left+tag)%5)...)
	leftKey := wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(left)
	out = append(out, leftKey[:]...)
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 8+(right+tag)%5)...)
	rightKey := wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(right)
	out = append(out, rightKey[:]...)
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 10+(op+tag)%5)...)
	opKey := wlmLmRawRepZeroshotAliasFamilyR1Family2OpMotif(op)
	out = append(out, opKey[:]...)
	out = append(out, wlmSiRawRepAliasInvarianceFalsificationR1Noise(&state, 7+(tag+left+right)%5)...)
	return out
}

func wlmLmRawRepZeroshotAliasFamilyR1ExpectedDenseValue(
	semantic int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	model wlmLmRawRepContextualAliasInductionR2ClassModel,
) (int, bool) {
	surfaceID, ok := repIndex[wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(semantic)]
	if !ok || surfaceID < 0 || surfaceID >= len(classes.surfaceToClass) {
		return 0, false
	}
	classID := classes.surfaceToClass[surfaceID]
	dense, ok := model.valueMap[classID]
	return dense, ok
}

func wlmLmRawRepZeroshotAliasFamilyR1EvaluateClass(
	leftSemantic, rightSemantic, opSemantic, tag int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	repIndex map[wlmSiRawRepAliasInvarianceFalsificationR1Key]int,
	classes wlmLmRawRepContextualAliasInductionR2Classes,
	model wlmLmRawRepContextualAliasInductionR2ClassModel,
	metrics map[string]float64,
) bool {
	raw := wlmLmRawRepZeroshotAliasFamilyR1EvalRecord(leftSemantic, rightSemantic, opSemantic, tag)
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
	targetLeft, ok3 := wlmLmRawRepZeroshotAliasFamilyR1ExpectedDenseValue(targetLeftSemantic, reps, repIndex, classes, model)
	targetRight, ok4 := wlmLmRawRepZeroshotAliasFamilyR1ExpectedDenseValue(targetRightSemantic, reps, repIndex, classes, model)
	if !ok3 || !ok4 {
		metrics["invalid_row_count"]++
		return false
	}
	return predLeft == targetLeft && predRight == targetRight
}

func (m *wlmLmRawRepZeroshotAliasFamilyR1SurfaceModel) observe(op, left, right, outLeft, outRight int, metrics map[string]float64) {
	if m.counts == nil {
		m.counts = make(map[wlmLmRawRepZeroshotAliasFamilyR1SurfaceTransitionKey]map[wlmLmRawRepZeroshotAliasFamilyR1SurfacePair]uint32)
	}
	key := wlmLmRawRepZeroshotAliasFamilyR1SurfaceTransitionKey{op: op, left: left, right: right}
	row := m.counts[key]
	if row == nil {
		row = make(map[wlmLmRawRepZeroshotAliasFamilyR1SurfacePair]uint32)
		m.counts[key] = row
	}
	pair := wlmLmRawRepZeroshotAliasFamilyR1SurfacePair{left: outLeft, right: outRight}
	if row[pair] == ^uint32(0) {
		metrics["counter_overflow_count"]++
		return
	}
	row[pair]++
}

func (m *wlmLmRawRepZeroshotAliasFamilyR1SurfaceModel) predict(op, left, right int) (int, int, bool) {
	key := wlmLmRawRepZeroshotAliasFamilyR1SurfaceTransitionKey{op: op, left: left, right: right}
	row := m.counts[key]
	if len(row) == 0 {
		return 0, 0, false
	}
	first := true
	best := wlmLmRawRepZeroshotAliasFamilyR1SurfacePair{}
	var bestCount uint32
	for pair, count := range row {
		if first || count > bestCount || (count == bestCount && (pair.left < best.left || (pair.left == best.left && pair.right < best.right))) {
			first = false
			best = pair
			bestCount = count
		}
	}
	return best.left, best.right, true
}

func wlmLmRawRepZeroshotAliasFamilyR1BuildSurfaceModel(
	taskRecords [][]uint8,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	metrics map[string]float64,
) (wlmLmRawRepZeroshotAliasFamilyR1SurfaceModel, int) {
	var model wlmLmRawRepZeroshotAliasFamilyR1SurfaceModel
	decodeFails := 0
	for _, raw := range taskRecords {
		row := wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw, reps)
		if len(row) != 5 {
			decodeFails++
			continue
		}
		model.observe(row[2], row[0], row[1], row[3], row[4], metrics)
	}
	return model, decodeFails
}

func wlmLmRawRepZeroshotAliasFamilyR1SemanticOfValueKey(key wlmSiRawRepAliasInvarianceFalsificationR1Key) int {
	for family := 0; family < 2; family++ {
		for semantic := 0; semantic < 4; semantic++ {
			if key == wlmSiRawRepAliasInvarianceFalsificationR1ValueMotif(family, semantic) {
				return semantic
			}
		}
	}
	for semantic := 0; semantic < 4; semantic++ {
		if key == wlmLmRawRepZeroshotAliasFamilyR1Family2ValueMotif(semantic) {
			return semantic
		}
	}
	return -1
}

func wlmLmRawRepZeroshotAliasFamilyR1EvaluateSurface(
	leftSemantic, rightSemantic, opSemantic, tag int,
	reps []wlmSiRawRepAliasInvarianceFalsificationR1Key,
	model wlmLmRawRepZeroshotAliasFamilyR1SurfaceModel,
	metrics map[string]float64,
) (bool, bool) {
	raw := wlmLmRawRepZeroshotAliasFamilyR1EvalRecord(leftSemantic, rightSemantic, opSemantic, tag)
	row := wlmSiRawRepAliasInvarianceFalsificationR1Decode(raw, reps)
	if len(row) != 3 {
		metrics["invalid_row_count"]++
		return false, false
	}
	predLeft, predRight, found := model.predict(row[2], row[0], row[1])
	if !found {
		// All transition counts are tied at zero for this exact unseen key.
		// Frozen control tie-break: lowest surface-representation pair.
		predLeft, predRight = 0, 0
	}
	if predLeft < 0 || predLeft >= len(reps) || predRight < 0 || predRight >= len(reps) {
		metrics["invalid_row_count"]++
		return false, !found
	}
	semanticLeft := wlmLmRawRepZeroshotAliasFamilyR1SemanticOfValueKey(reps[predLeft])
	semanticRight := wlmLmRawRepZeroshotAliasFamilyR1SemanticOfValueKey(reps[predRight])
	targetLeft, targetRight := wlmSiRawRepAliasInvarianceFalsificationR1Apply(leftSemantic, rightSemantic, opSemantic)
	return semanticLeft == targetLeft && semanticRight == targetRight, !found
}

// RunWlmLmRawRepZeroshotAliasFamilyR1 tests context-only transfer to an unseen surface family.
func RunWlmLmRawRepZeroshotAliasFamilyR1() interface{} {
	metrics := map[string]float64{
		"bridge_record_count":                              0,
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

	bridge := wlmLmRawRepZeroshotAliasFamilyR1BridgeRecords(false)
	metrics["bridge_record_count"] = float64(len(bridge))
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

	shuffledBridge := wlmLmRawRepZeroshotAliasFamilyR1BridgeRecords(true)
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

	return wlmLmRawRepZeroshotAliasFamilyR1Result{
		Schema:     "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-ZEROSHOT-ALIAS-FAMILY-R1",
		Metrics:    metrics,
	}
}
