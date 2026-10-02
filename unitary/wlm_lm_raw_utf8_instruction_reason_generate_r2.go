package unitary

import (
	"math"
)

const wlmLmRawUTF8InstructionReasonGenerateR2FeatureDim = 1280

type wlmLmRawUTF8InstructionReasonGenerateR2Head struct {
	classes       int
	classExamples []uint32
	counts        [][]uint32
	totalExamples uint32
}

type wlmLmRawUTF8InstructionReasonGenerateR2Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmRawUTF8InstructionReasonGenerateR2Hash(window []byte) uint32 {
	h := uint32(2166136261)
	h ^= uint32(len(window))
	h *= 16777619
	for _, b := range window {
		h ^= uint32(b)
		h *= 16777619
	}
	h ^= h >> 16
	h *= 0x7feb352d
	h ^= h >> 15
	return h
}

func wlmLmRawUTF8InstructionReasonGenerateR2Features(text string, metrics map[string]float64) []int {
	data := []byte(text)
	if len(data) < 2 {
		metrics["invalid_row_count"]++
		return nil
	}
	var present [wlmLmRawUTF8InstructionReasonGenerateR2FeatureDim]bool
	active := make([]int, 0, 256)
	for width := 2; width <= 4; width++ {
		for start := 0; start+width <= len(data); start++ {
			h := wlmLmRawUTF8InstructionReasonGenerateR2Hash(data[start : start+width])
			global := int(h % 256)
			if !present[global] {
				present[global] = true
				active = append(active, global)
			}
			bucket := start * 8 / len(data)
			if bucket > 7 {
				bucket = 7
			}
			pos := 256 + bucket*128 + int((h>>8)%128)
			if !present[pos] {
				present[pos] = true
				active = append(active, pos)
			}
		}
	}
	if len(active) == 0 {
		metrics["invalid_row_count"]++
	}
	return active
}

func wlmLmRawUTF8InstructionReasonGenerateR2NewHead(classes int) *wlmLmRawUTF8InstructionReasonGenerateR2Head {
	head := &wlmLmRawUTF8InstructionReasonGenerateR2Head{
		classes:       classes,
		classExamples: make([]uint32, classes),
		counts:        make([][]uint32, classes),
	}
	for class := 0; class < classes; class++ {
		head.counts[class] = make([]uint32, wlmLmRawUTF8InstructionReasonGenerateR2FeatureDim)
	}
	return head
}

func (h *wlmLmRawUTF8InstructionReasonGenerateR2Head) observe(features []int, target int, metrics map[string]float64) {
	if target < 0 || target >= h.classes {
		metrics["invalid_row_count"]++
		return
	}
	if h.classExamples[target] == ^uint32(0) || h.totalExamples == ^uint32(0) {
		metrics["counter_overflow_count"]++
		return
	}
	h.classExamples[target]++
	h.totalExamples++
	for _, feature := range features {
		if feature < 0 || feature >= wlmLmRawUTF8InstructionReasonGenerateR2FeatureDim {
			metrics["invalid_row_count"]++
			continue
		}
		if h.counts[target][feature] == ^uint32(0) {
			metrics["counter_overflow_count"]++
			continue
		}
		h.counts[target][feature]++
	}
}

func (h *wlmLmRawUTF8InstructionReasonGenerateR2Head) predict(features []int) int {
	best := 0
	bestScore := math.Inf(-1)
	for class := 0; class < h.classes; class++ {
		classN := float64(h.classExamples[class])
		otherN := float64(h.totalExamples - h.classExamples[class])
		score := math.Log((classN + 1) / (float64(h.totalExamples) + float64(h.classes)))
		for _, feature := range features {
			classCount := float64(h.counts[class][feature])
			var allCount uint32
			for otherClass := 0; otherClass < h.classes; otherClass++ {
				allCount += h.counts[otherClass][feature]
			}
			otherCount := float64(allCount - h.counts[class][feature])
			pos := (classCount + 1) / (classN + 2)
			neg := (otherCount + 1) / (otherN + 2)
			score += math.Log(pos / neg)
		}
		if score > bestScore {
			best = class
			bestScore = score
		}
	}
	return best
}

func wlmLmRawUTF8InstructionReasonGenerateR2Accuracy(
	examples []wlmLmRawUTF8InstructionReasonGenerateR1Example,
	features [][]int,
	opHead, leftHead, rightHead *wlmLmRawUTF8InstructionReasonGenerateR2Head,
) (float64, float64, float64, float64) {
	opHits, leftHits, rightHits, fullHits := 0, 0, 0, 0
	for i, ex := range examples {
		op := opHead.predict(features[i])
		left := leftHead.predict(features[i])
		right := rightHead.predict(features[i])
		if op == ex.op {
			opHits++
		}
		if left == ex.left {
			leftHits++
		}
		if right == ex.right {
			rightHits++
		}
		if op == ex.op && left == ex.left && right == ex.right {
			fullHits++
		}
	}
	den := float64(len(examples))
	if den == 0 {
		return 0, 0, 0, 0
	}
	return float64(opHits) / den, float64(leftHits) / den, float64(rightHits) / den, float64(fullHits) / den
}

// RunWlmLmRawUTF8InstructionReasonGenerateR2 evaluates the fixed raw motif-evidence representation repair.
func RunWlmLmRawUTF8InstructionReasonGenerateR2() interface{} {
	metrics := map[string]float64{
		"training_example_count": 0,
		"heldout_example_count": 0,
		"training_context_count": 0,
		"heldout_context_count": 0,
		"heldout_exact_instruction_leak_count": 0,
		"raw_motif_feature_dimension": wlmLmRawUTF8InstructionReasonGenerateR2FeatureDim,
		"raw_motif_min_width": 2,
		"raw_motif_max_width": 4,
		"operation_head_class_count": 4,
		"value_head_class_count": 8,
		"minimum_training_operation_accuracy": 0,
		"minimum_training_left_accuracy": 0,
		"minimum_training_right_accuracy": 0,
		"minimum_heldout_operation_accuracy": 0,
		"minimum_heldout_left_accuracy": 0,
		"minimum_heldout_right_accuracy": 0,
		"minimum_heldout_full_representation_accuracy": 0,
		"reasoning_counter_count": 96,
		"minimum_heldout_reasoning_result_accuracy": 0,
		"minimum_output_generator_class_coverage": 0,
		"maximum_output_generator_state_count": 0,
		"minimum_heldout_exact_raw_output_accuracy": 0,
		"maximum_whole_instruction_lookup_exact_output_accuracy": 0,
		"tokenizer_use_count": 0,
		"external_model_call_count": 0,
		"capacity_growth_event_count": 0,
		"invalid_row_count": 0,
		"counter_overflow_count": 0,
	}

	train, held := wlmLmRawUTF8InstructionReasonGenerateR1Corpus()
	metrics["training_example_count"] = float64(len(train))
	metrics["heldout_example_count"] = float64(len(held))

	trainContexts := make(map[[3]int]bool)
	heldContexts := make(map[[3]int]bool)
	trainStrings := make(map[string]bool)
	for _, ex := range train {
		trainContexts[[3]int{ex.op, ex.left, ex.right}] = true
		trainStrings[ex.text] = true
	}
	for _, ex := range held {
		heldContexts[[3]int{ex.op, ex.left, ex.right}] = true
		if trainStrings[ex.text] {
			metrics["heldout_exact_instruction_leak_count"]++
		}
	}
	metrics["training_context_count"] = float64(len(trainContexts))
	metrics["heldout_context_count"] = float64(len(heldContexts))

	trainFeatures := make([][]int, len(train))
	heldFeatures := make([][]int, len(held))
	for i, ex := range train {
		trainFeatures[i] = wlmLmRawUTF8InstructionReasonGenerateR2Features(ex.text, metrics)
	}
	for i, ex := range held {
		heldFeatures[i] = wlmLmRawUTF8InstructionReasonGenerateR2Features(ex.text, metrics)
	}

	opHead := wlmLmRawUTF8InstructionReasonGenerateR2NewHead(4)
	leftHead := wlmLmRawUTF8InstructionReasonGenerateR2NewHead(8)
	rightHead := wlmLmRawUTF8InstructionReasonGenerateR2NewHead(8)
	for i, ex := range train {
		opHead.observe(trainFeatures[i], ex.op, metrics)
		leftHead.observe(trainFeatures[i], ex.left, metrics)
		rightHead.observe(trainFeatures[i], ex.right, metrics)
	}

	trainOp, trainLeft, trainRight, _ := wlmLmRawUTF8InstructionReasonGenerateR2Accuracy(train, trainFeatures, opHead, leftHead, rightHead)
	heldOp, heldLeft, heldRight, heldFull := wlmLmRawUTF8InstructionReasonGenerateR2Accuracy(held, heldFeatures, opHead, leftHead, rightHead)
	metrics["minimum_training_operation_accuracy"] = trainOp
	metrics["minimum_training_left_accuracy"] = trainLeft
	metrics["minimum_training_right_accuracy"] = trainRight
	metrics["minimum_heldout_operation_accuracy"] = heldOp
	metrics["minimum_heldout_left_accuracy"] = heldLeft
	metrics["minimum_heldout_right_accuracy"] = heldRight
	metrics["minimum_heldout_full_representation_accuracy"] = heldFull

	var reasoner wlmLmRawUTF8InstructionReasonGenerateR1Reasoner
	var generator wlmLmRawUTF8InstructionReasonGenerateR1Generator
	lookup := make(map[string]string, len(train))
	for i, ex := range train {
		predOp := opHead.predict(trainFeatures[i])
		predLeft := leftHead.predict(trainFeatures[i])
		predRight := rightHead.predict(trainFeatures[i])
		reasoner.observe(predOp, predLeft, predRight, ex.target, metrics)
		generator.observe(ex.target, ex.output, metrics)
		lookup[ex.text] = ex.output
	}

	coverage := 0
	for result := 0; result < 8; result++ {
		if generator.counts[result] != nil && len(generator.counts[result]) > 0 {
			coverage++
		}
	}
	metrics["minimum_output_generator_class_coverage"] = float64(coverage)
	metrics["maximum_output_generator_state_count"] = float64(generator.stateCount())

	fallback := wlmLmRawUTF8InstructionReasonGenerateR1Fallback(lookup)
	reasonHits, outputHits, lookupHits := 0, 0, 0
	for i, ex := range held {
		predOp := opHead.predict(heldFeatures[i])
		predLeft := leftHead.predict(heldFeatures[i])
		predRight := rightHead.predict(heldFeatures[i])
		predResult := reasoner.predict(predOp, predLeft, predRight)
		if predResult == ex.target {
			reasonHits++
		}
		if generator.generate(predResult) == ex.output {
			outputHits++
		}
		lookupOutput, ok := lookup[ex.text]
		if !ok {
			lookupOutput = fallback
		}
		if lookupOutput == ex.output {
			lookupHits++
		}
	}
	if len(held) > 0 {
		den := float64(len(held))
		metrics["minimum_heldout_reasoning_result_accuracy"] = float64(reasonHits) / den
		metrics["minimum_heldout_exact_raw_output_accuracy"] = float64(outputHits) / den
		metrics["maximum_whole_instruction_lookup_exact_output_accuracy"] = float64(lookupHits) / den
	}

	for _, value := range metrics {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			metrics["invalid_row_count"]++
		}
	}

	return wlmLmRawUTF8InstructionReasonGenerateR2Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-UTF8-INSTRUCTION-REASON-GENERATE-R2",
		Metrics: metrics,
	}
}
