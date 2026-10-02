package unitary

import (
	"math"
	"sort"
)

const wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim = 320

type wlmLmRawUTF8InstructionReasonGenerateR1Example struct {
	text   string
	op     int
	left   int
	right  int
	target int
	output string
}

type wlmLmRawUTF8InstructionReasonGenerateR1Head struct {
	classes int
	w       [][]float64
	b       []float64
}

type wlmLmRawUTF8InstructionReasonGenerateR1Reasoner struct {
	counts [4][3][2][2][2]uint32
}

type wlmLmRawUTF8InstructionReasonGenerateR1GenContext [2]int

type wlmLmRawUTF8InstructionReasonGenerateR1Generator struct {
	counts [8]map[wlmLmRawUTF8InstructionReasonGenerateR1GenContext]*[256]uint32
}

type wlmLmRawUTF8InstructionReasonGenerateR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmRawUTF8InstructionReasonGenerateR1ValueWord(id int) string {
	values := [...]string{"zero", "one", "two", "three", "four", "five", "six", "seven"}
	if id < 0 || id >= len(values) {
		return ""
	}
	return values[id]
}

func wlmLmRawUTF8InstructionReasonGenerateR1OpWord(id int) string {
	ops := [...]string{"xor", "xnor", "and", "or"}
	if id < 0 || id >= len(ops) {
		return ""
	}
	return ops[id]
}

func wlmLmRawUTF8InstructionReasonGenerateR1Output(result int) string {
	return "answer: " + wlmLmRawUTF8InstructionReasonGenerateR1ValueWord(result) + ".\n"
}

func wlmLmRawUTF8InstructionReasonGenerateR1Format(templateID, op, left, right int) string {
	o := wlmLmRawUTF8InstructionReasonGenerateR1OpWord(op)
	l := wlmLmRawUTF8InstructionReasonGenerateR1ValueWord(left)
	r := wlmLmRawUTF8InstructionReasonGenerateR1ValueWord(right)
	switch templateID {
	case 0:
		return "compute " + o + " of " + l + " and " + r + ".\n"
	case 1:
		return "please apply " + o + " to " + l + " with " + r + ".\n"
	case 2:
		return "using " + l + " and " + r + ", perform " + o + ".\n"
	case 3:
		return o + " on " + l + " plus " + r + "; give the answer.\n"
	case 4:
		return "what is " + l + " when you " + o + " it with " + r + "?\n"
	default:
		return "return the result of " + o + " for " + l + " and " + r + ".\n"
	}
}

func wlmLmRawUTF8InstructionReasonGenerateR1Heldout(op, left, right int) bool {
	return (op+left+right)%4 == 1
}

func wlmLmRawUTF8InstructionReasonGenerateR1Corpus() (train, held []wlmLmRawUTF8InstructionReasonGenerateR1Example) {
	for op := 0; op < 4; op++ {
		for left := 0; left < 8; left++ {
			for right := 0; right < 8; right++ {
				target := wlmLmRuleConditionedSequenceR2Target(op, left, right)
				if wlmLmRawUTF8InstructionReasonGenerateR1Heldout(op, left, right) {
					for templateID := 4; templateID < 6; templateID++ {
						held = append(held, wlmLmRawUTF8InstructionReasonGenerateR1Example{
							text:   wlmLmRawUTF8InstructionReasonGenerateR1Format(templateID, op, left, right),
							op:     op,
							left:   left,
							right:  right,
							target: target,
							output: wlmLmRawUTF8InstructionReasonGenerateR1Output(target),
						})
					}
					continue
				}
				for templateID := 0; templateID < 4; templateID++ {
					train = append(train, wlmLmRawUTF8InstructionReasonGenerateR1Example{
						text:   wlmLmRawUTF8InstructionReasonGenerateR1Format(templateID, op, left, right),
						op:     op,
						left:   left,
						right:  right,
						target: target,
						output: wlmLmRawUTF8InstructionReasonGenerateR1Output(target),
					})
				}
			}
		}
	}
	sort.Slice(train, func(i, j int) bool { return train[i].text < train[j].text })
	sort.Slice(held, func(i, j int) bool { return held[i].text < held[j].text })
	return
}

func wlmLmRawUTF8InstructionReasonGenerateR1Features(text string, metrics map[string]float64) [wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim]float64 {
	var out [wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim]float64
	data := []byte(text)
	if len(data) == 0 {
		metrics["invalid_row_count"]++
		return out
	}
	var h [64]float64
	var sums [4][64]float64
	var counts [4]int
	for i, b := range data {
		h = uplm0aStep(h, b)
		bucket := i * 4 / len(data)
		if bucket > 3 {
			bucket = 3
		}
		counts[bucket]++
		for d := 0; d < 64; d++ {
			sums[bucket][d] += h[d]
		}
	}
	for d := 0; d < 64; d++ {
		out[d] = h[d]
	}
	for bucket := 0; bucket < 4; bucket++ {
		if counts[bucket] == 0 {
			metrics["invalid_row_count"]++
			continue
		}
		den := float64(counts[bucket])
		for d := 0; d < 64; d++ {
			out[64+bucket*64+d] = sums[bucket][d] / den
		}
	}
	for _, value := range out {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			metrics["invalid_row_count"]++
			break
		}
	}
	return out
}

func wlmLmRawUTF8InstructionReasonGenerateR1NewHead(classes int) *wlmLmRawUTF8InstructionReasonGenerateR1Head {
	head := &wlmLmRawUTF8InstructionReasonGenerateR1Head{
		classes: classes,
		w:       make([][]float64, classes),
		b:       make([]float64, classes),
	}
	for class := 0; class < classes; class++ {
		head.w[class] = make([]float64, wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim)
	}
	return head
}

func (h *wlmLmRawUTF8InstructionReasonGenerateR1Head) probs(features *[wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim]float64) []float64 {
	logits := make([]float64, h.classes)
	maxLogit := math.Inf(-1)
	for class := 0; class < h.classes; class++ {
		score := h.b[class]
		for d := 0; d < wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim; d++ {
			score += h.w[class][d] * features[d]
		}
		logits[class] = score
		if score > maxLogit {
			maxLogit = score
		}
	}
	sum := 0.0
	for class := 0; class < h.classes; class++ {
		logits[class] = math.Exp(logits[class] - maxLogit)
		sum += logits[class]
	}
	if sum == 0 || math.IsNaN(sum) || math.IsInf(sum, 0) {
		return logits
	}
	for class := 0; class < h.classes; class++ {
		logits[class] /= sum
	}
	return logits
}

func (h *wlmLmRawUTF8InstructionReasonGenerateR1Head) train(features *[wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim]float64, target int, lr float64, metrics map[string]float64) {
	if target < 0 || target >= h.classes {
		metrics["invalid_row_count"]++
		return
	}
	p := h.probs(features)
	for class := 0; class < h.classes; class++ {
		if math.IsNaN(p[class]) || math.IsInf(p[class], 0) {
			metrics["invalid_row_count"]++
			return
		}
		gradient := p[class]
		if class == target {
			gradient -= 1
		}
		for d := 0; d < wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim; d++ {
			h.w[class][d] -= lr * gradient * features[d]
		}
		h.b[class] -= lr * gradient
	}
}

func (h *wlmLmRawUTF8InstructionReasonGenerateR1Head) predict(features *[wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim]float64) int {
	p := h.probs(features)
	best := 0
	for class := 1; class < len(p); class++ {
		if p[class] > p[best] {
			best = class
		}
	}
	return best
}

func (r *wlmLmRawUTF8InstructionReasonGenerateR1Reasoner) observe(op, left, right, target int, metrics map[string]float64) {
	if op < 0 || op >= 4 || left < 0 || left >= 8 || right < 0 || right >= 8 || target < 0 || target >= 8 {
		metrics["invalid_row_count"]++
		return
	}
	for bit := 0; bit < 3; bit++ {
		x := (left >> bit) & 1
		y := (right >> bit) & 1
		z := (target >> bit) & 1
		if r.counts[op][bit][x][y][z] == ^uint32(0) {
			metrics["counter_overflow_count"]++
			continue
		}
		r.counts[op][bit][x][y][z]++
	}
}

func (r *wlmLmRawUTF8InstructionReasonGenerateR1Reasoner) predict(op, left, right int) int {
	if op < 0 || op >= 4 || left < 0 || left >= 8 || right < 0 || right >= 8 {
		return 0
	}
	out := 0
	for bit := 0; bit < 3; bit++ {
		x := (left >> bit) & 1
		y := (right >> bit) & 1
		if r.counts[op][bit][x][y][1] > r.counts[op][bit][x][y][0] {
			out |= 1 << bit
		}
	}
	return out
}

func (g *wlmLmRawUTF8InstructionReasonGenerateR1Generator) observe(result int, output string, metrics map[string]float64) {
	if result < 0 || result >= 8 {
		metrics["invalid_row_count"]++
		return
	}
	if g.counts[result] == nil {
		g.counts[result] = make(map[wlmLmRawUTF8InstructionReasonGenerateR1GenContext]*[256]uint32)
	}
	ctx := wlmLmRawUTF8InstructionReasonGenerateR1GenContext{256, 256}
	for _, b := range []byte(output) {
		row := g.counts[result][ctx]
		if row == nil {
			row = &[256]uint32{}
			g.counts[result][ctx] = row
		}
		if row[int(b)] == ^uint32(0) {
			metrics["counter_overflow_count"]++
		} else {
			row[int(b)]++
		}
		ctx = wlmLmRawUTF8InstructionReasonGenerateR1GenContext{ctx[1], int(b)}
	}
}

func (g *wlmLmRawUTF8InstructionReasonGenerateR1Generator) generate(result int) string {
	if result < 0 || result >= 8 || g.counts[result] == nil {
		return ""
	}
	ctx := wlmLmRawUTF8InstructionReasonGenerateR1GenContext{256, 256}
	out := make([]byte, 0, 32)
	for step := 0; step < 32; step++ {
		row := g.counts[result][ctx]
		if row == nil {
			break
		}
		best := 0
		bestCount := row[0]
		for b := 1; b < 256; b++ {
			if row[b] > bestCount {
				best = b
				bestCount = row[b]
			}
		}
		if bestCount == 0 {
			break
		}
		out = append(out, byte(best))
		if byte(best) == '\n' {
			break
		}
		ctx = wlmLmRawUTF8InstructionReasonGenerateR1GenContext{ctx[1], best}
	}
	return string(out)
}

func (g *wlmLmRawUTF8InstructionReasonGenerateR1Generator) stateCount() int {
	total := 0
	for result := 0; result < 8; result++ {
		total += len(g.counts[result])
	}
	return total
}

func wlmLmRawUTF8InstructionReasonGenerateR1Fallback(lookup map[string]string) string {
	counts := make(map[string]int)
	for _, output := range lookup {
		counts[output]++
	}
	keys := make([]string, 0, len(counts))
	for output := range counts {
		keys = append(keys, output)
	}
	sort.Strings(keys)
	best := ""
	bestCount := -1
	for _, output := range keys {
		if counts[output] > bestCount {
			best = output
			bestCount = counts[output]
		}
	}
	return best
}

func wlmLmRawUTF8InstructionReasonGenerateR1Accuracy(
	examples []wlmLmRawUTF8InstructionReasonGenerateR1Example,
	features [][wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim]float64,
	opHead, leftHead, rightHead *wlmLmRawUTF8InstructionReasonGenerateR1Head,
) (float64, float64, float64, float64) {
	opHits, leftHits, rightHits, fullHits := 0, 0, 0, 0
	for i, ex := range examples {
		op := opHead.predict(&features[i])
		left := leftHead.predict(&features[i])
		right := rightHead.predict(&features[i])
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

// RunWlmLmRawUTF8InstructionReasonGenerateR1 evaluates raw UTF-8 instruction representation, reasoning, and generation.
func RunWlmLmRawUTF8InstructionReasonGenerateR1() interface{} {
	metrics := map[string]float64{
		"training_example_count": 0,
		"heldout_example_count": 0,
		"training_context_count": 0,
		"heldout_context_count": 0,
		"heldout_exact_instruction_leak_count": 0,
		"raw_feature_dimension": wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim,
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

	trainFeatures := make([][wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim]float64, len(train))
	heldFeatures := make([][wlmLmRawUTF8InstructionReasonGenerateR1FeatureDim]float64, len(held))
	for i, ex := range train {
		trainFeatures[i] = wlmLmRawUTF8InstructionReasonGenerateR1Features(ex.text, metrics)
	}
	for i, ex := range held {
		heldFeatures[i] = wlmLmRawUTF8InstructionReasonGenerateR1Features(ex.text, metrics)
	}

	opHead := wlmLmRawUTF8InstructionReasonGenerateR1NewHead(4)
	leftHead := wlmLmRawUTF8InstructionReasonGenerateR1NewHead(8)
	rightHead := wlmLmRawUTF8InstructionReasonGenerateR1NewHead(8)
	const epochs = 40
	const lr = 0.05
	for epoch := 0; epoch < epochs; epoch++ {
		for i, ex := range train {
			opHead.train(&trainFeatures[i], ex.op, lr, metrics)
			leftHead.train(&trainFeatures[i], ex.left, lr, metrics)
			rightHead.train(&trainFeatures[i], ex.right, lr, metrics)
		}
	}

	trainOp, trainLeft, trainRight, _ := wlmLmRawUTF8InstructionReasonGenerateR1Accuracy(train, trainFeatures, opHead, leftHead, rightHead)
	heldOp, heldLeft, heldRight, heldFull := wlmLmRawUTF8InstructionReasonGenerateR1Accuracy(held, heldFeatures, opHead, leftHead, rightHead)
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
		predOp := opHead.predict(&trainFeatures[i])
		predLeft := leftHead.predict(&trainFeatures[i])
		predRight := rightHead.predict(&trainFeatures[i])
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
		predOp := opHead.predict(&heldFeatures[i])
		predLeft := leftHead.predict(&heldFeatures[i])
		predRight := rightHead.predict(&heldFeatures[i])
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

	return wlmLmRawUTF8InstructionReasonGenerateR1Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-UTF8-INSTRUCTION-REASON-GENERATE-R1",
		Metrics: metrics,
	}
}
