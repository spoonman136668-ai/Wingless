package unitary

import "math"

type wlmLmRawUTF8InstructionReasonGenerateR5GenContext [3]int

type wlmLmRawUTF8InstructionReasonGenerateR5Generator struct {
	counts [8]map[wlmLmRawUTF8InstructionReasonGenerateR5GenContext]*[256]uint32
}

type wlmLmRawUTF8InstructionReasonGenerateR5Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	OperationAnchors []string `json:"operation_anchors"`
	ValueAnchors []string `json:"value_anchors"`
}

func (g *wlmLmRawUTF8InstructionReasonGenerateR5Generator) observe(result int, output string, metrics map[string]float64) {
	if result < 0 || result >= 8 {
		metrics["invalid_row_count"]++
		return
	}
	if g.counts[result] == nil {
		g.counts[result] = make(map[wlmLmRawUTF8InstructionReasonGenerateR5GenContext]*[256]uint32)
	}
	ctx := wlmLmRawUTF8InstructionReasonGenerateR5GenContext{256,256,256}
	for _,b:=range []byte(output) {
		row:=g.counts[result][ctx]
		if row==nil {
			row=&[256]uint32{}
			g.counts[result][ctx]=row
		}
		if row[int(b)]==^uint32(0) {
			metrics["counter_overflow_count"]++
		} else {
			row[int(b)]++
		}
		ctx=wlmLmRawUTF8InstructionReasonGenerateR5GenContext{ctx[1],ctx[2],int(b)}
	}
}

func (g *wlmLmRawUTF8InstructionReasonGenerateR5Generator) generate(result int) string {
	if result < 0 || result >= 8 || g.counts[result]==nil {
		return ""
	}
	ctx:=wlmLmRawUTF8InstructionReasonGenerateR5GenContext{256,256,256}
	out:=make([]byte,0,32)
	for step:=0;step<32;step++ {
		row:=g.counts[result][ctx]
		if row==nil { break }
		best:=0
		bestCount:=row[0]
		for b:=1;b<256;b++ {
			if row[b]>bestCount {
				best=b
				bestCount=row[b]
			}
		}
		if bestCount==0 { break }
		out=append(out,byte(best))
		if byte(best)=='\n' { break }
		ctx=wlmLmRawUTF8InstructionReasonGenerateR5GenContext{ctx[1],ctx[2],best}
	}
	return string(out)
}

func (g *wlmLmRawUTF8InstructionReasonGenerateR5Generator) stateCount() int {
	total:=0
	for result:=0;result<8;result++ {
		total+=len(g.counts[result])
	}
	return total
}

// RunWlmLmRawUTF8InstructionReasonGenerateR5 tests the three-byte autoregressive answer generator.
func RunWlmLmRawUTF8InstructionReasonGenerateR5() interface{} {
	metrics:=map[string]float64{
		"training_example_count":0,
		"heldout_example_count":0,
		"training_context_count":0,
		"heldout_context_count":0,
		"heldout_exact_instruction_leak_count":0,
		"selected_operation_anchor_count":4,
		"selected_value_anchor_count":8,
		"value_anchor_duplicate_count":0,
		"value_operation_anchor_collision_count":0,
		"minimum_heldout_operation_accuracy":0,
		"minimum_heldout_left_accuracy":0,
		"minimum_heldout_right_accuracy":0,
		"maximum_heldout_operand_decode_failure_count":0,
		"minimum_heldout_full_representation_accuracy":0,
		"reasoning_counter_count":96,
		"minimum_heldout_reasoning_result_accuracy":0,
		"minimum_output_generator_class_coverage":0,
		"maximum_output_generator_state_count":0,
		"minimum_generator_canonical_class_exact_accuracy":0,
		"minimum_heldout_exact_raw_output_accuracy":0,
		"maximum_whole_instruction_lookup_exact_output_accuracy":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
	}

	train,held:=wlmLmRawUTF8InstructionReasonGenerateR1Corpus()
	metrics["training_example_count"]=float64(len(train))
	metrics["heldout_example_count"]=float64(len(held))
	trainContexts:=make(map[[3]int]bool)
	heldContexts:=make(map[[3]int]bool)
	trainStrings:=make(map[string]bool)
	for _,ex:=range train {
		trainContexts[[3]int{ex.op,ex.left,ex.right}]=true
		trainStrings[ex.text]=true
	}
	for _,ex:=range held {
		heldContexts[[3]int{ex.op,ex.left,ex.right}]=true
		if trainStrings[ex.text] { metrics["heldout_exact_instruction_leak_count"]++ }
	}
	metrics["training_context_count"]=float64(len(trainContexts))
	metrics["heldout_context_count"]=float64(len(heldContexts))

	opAnchors,_:=wlmLmRawUTF8InstructionReasonGenerateR3SelectAnchors(train)
	valueAnchors:=wlmLmRawUTF8InstructionReasonGenerateR4SelectValueAnchors(train)
	duplicates,collisions:=wlmLmRawUTF8InstructionReasonGenerateR4CollisionMetrics(opAnchors,valueAnchors)
	metrics["value_anchor_duplicate_count"]=float64(duplicates)
	metrics["value_operation_anchor_collision_count"]=float64(collisions)

	opModel:=wlmLmRawUTF8InstructionReasonGenerateR3TrainOpModel(train,opAnchors,metrics)
	heldOp,heldLeft,heldRight,heldFull,heldDecode:=wlmLmRawUTF8InstructionReasonGenerateR3Accuracy(held,&opModel,valueAnchors)
	metrics["minimum_heldout_operation_accuracy"]=heldOp
	metrics["minimum_heldout_left_accuracy"]=heldLeft
	metrics["minimum_heldout_right_accuracy"]=heldRight
	metrics["maximum_heldout_operand_decode_failure_count"]=float64(heldDecode)
	metrics["minimum_heldout_full_representation_accuracy"]=heldFull

	var reasoner wlmLmRawUTF8InstructionReasonGenerateR1Reasoner
	var generator wlmLmRawUTF8InstructionReasonGenerateR5Generator
	lookup:=make(map[string]string,len(train))
	for _,ex:=range train {
		op:=opModel.predict(ex.text)
		left,right,_:=wlmLmRawUTF8InstructionReasonGenerateR3FindOperands(ex.text,valueAnchors)
		reasoner.observe(op,left,right,ex.target,metrics)
		generator.observe(ex.target,ex.output,metrics)
		lookup[ex.text]=ex.output
	}

	coverage:=0
	canonicalHits:=0
	for result:=0;result<8;result++ {
		if generator.counts[result]!=nil&&len(generator.counts[result])>0 { coverage++ }
		if generator.generate(result)==wlmLmRawUTF8InstructionReasonGenerateR1Output(result) { canonicalHits++ }
	}
	metrics["minimum_output_generator_class_coverage"]=float64(coverage)
	metrics["maximum_output_generator_state_count"]=float64(generator.stateCount())
	metrics["minimum_generator_canonical_class_exact_accuracy"]=float64(canonicalHits)/8.0

	fallback:=wlmLmRawUTF8InstructionReasonGenerateR1Fallback(lookup)
	reasonHits,outputHits,lookupHits:=0,0,0
	for _,ex:=range held {
		op:=opModel.predict(ex.text)
		left,right,_:=wlmLmRawUTF8InstructionReasonGenerateR3FindOperands(ex.text,valueAnchors)
		result:=reasoner.predict(op,left,right)
		if result==ex.target { reasonHits++ }
		if generator.generate(result)==ex.output { outputHits++ }
		lookupOutput,ok:=lookup[ex.text]
		if !ok { lookupOutput=fallback }
		if lookupOutput==ex.output { lookupHits++ }
	}
	if len(held)>0 {
		den:=float64(len(held))
		metrics["minimum_heldout_reasoning_result_accuracy"]=float64(reasonHits)/den
		metrics["minimum_heldout_exact_raw_output_accuracy"]=float64(outputHits)/den
		metrics["maximum_whole_instruction_lookup_exact_output_accuracy"]=float64(lookupHits)/den
	}
	for _,value:=range metrics {
		if math.IsNaN(value)||math.IsInf(value,0) { metrics["invalid_row_count"]++ }
	}

	opOut:=make([]string,4)
	valueOut:=make([]string,8)
	copy(opOut,opAnchors[:])
	copy(valueOut,valueAnchors[:])
	return wlmLmRawUTF8InstructionReasonGenerateR5Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-UTF8-INSTRUCTION-REASON-GENERATE-R5",
		Metrics:metrics,
		OperationAnchors:opOut,
		ValueAnchors:valueOut,
	}
}
