package unitary

import "math"

type wlmLmRawUTF8InstructionReasonGenerateR4Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	OperationAnchors []string `json:"operation_anchors"`
	ValueAnchors []string `json:"value_anchors"`
}

func wlmLmRawUTF8InstructionReasonGenerateR4BetterValue(
	candidate string,
	row *wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats,
	class int,
	totalExamples int,
	classN int,
	best string,
	bestRow *wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats,
) bool {
	if bestRow==nil {
		return true
	}
	posPresence:=row.valuePresence[class]
	bestPosPresence:=bestRow.valuePresence[class]
	coverage:=float64(posPresence)/float64(classN)
	bestCoverage:=float64(bestPosPresence)/float64(classN)
	if coverage!=bestCoverage {
		return coverage>bestCoverage
	}

	precision:=float64(posPresence)/float64(row.totalPresence)
	bestPrecision:=float64(bestPosPresence)/float64(bestRow.totalPresence)
	if precision!=bestPrecision {
		return precision>bestPrecision
	}

	posOcc:=row.valueOccurrences[class]
	bestPosOcc:=bestRow.valueOccurrences[class]
	negN:=totalExamples-classN
	meanPos:=float64(posOcc)/float64(classN)
	meanNeg:=0.0
	bestMeanPos:=float64(bestPosOcc)/float64(classN)
	bestMeanNeg:=0.0
	if negN>0 {
		meanNeg=float64(row.totalOccurrences-posOcc)/float64(negN)
		bestMeanNeg=float64(bestRow.totalOccurrences-bestPosOcc)/float64(negN)
	}
	margin:=meanPos-meanNeg
	bestMargin:=bestMeanPos-bestMeanNeg
	if margin!=bestMargin {
		return margin>bestMargin
	}
	if len(candidate)!=len(best) {
		return len(candidate)>len(best)
	}
	return candidate<best
}

func wlmLmRawUTF8InstructionReasonGenerateR4SelectValueAnchors(
	train []wlmLmRawUTF8InstructionReasonGenerateR1Example,
)[8]string {
	stats,_,valueN:=wlmLmRawUTF8InstructionReasonGenerateR3BuildStats(train)
	var anchors [8]string
	for class:=0;class<8;class++ {
		var bestRow *wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats
		for candidate,row:=range stats {
			if row.valuePresence[class]==0 { continue }
			if wlmLmRawUTF8InstructionReasonGenerateR4BetterValue(candidate,row,class,len(train),valueN[class],anchors[class],bestRow) {
				anchors[class]=candidate
				bestRow=row
			}
		}
	}
	return anchors
}

func wlmLmRawUTF8InstructionReasonGenerateR4CollisionMetrics(opAnchors [4]string,valueAnchors [8]string)(duplicates,collisions int) {
	for i:=0;i<8;i++ {
		for j:=i+1;j<8;j++ {
			if valueAnchors[i]!=""&&valueAnchors[i]==valueAnchors[j] { duplicates++ }
		}
		for _,op:=range opAnchors {
			if valueAnchors[i]!=""&&valueAnchors[i]==op { collisions++ }
		}
	}
	return
}

// RunWlmLmRawUTF8InstructionReasonGenerateR4 tests precision-first value-anchor selection.
func RunWlmLmRawUTF8InstructionReasonGenerateR4() interface{} {
	metrics:=map[string]float64{
		"training_example_count":0,
		"heldout_example_count":0,
		"training_context_count":0,
		"heldout_context_count":0,
		"heldout_exact_instruction_leak_count":0,
		"selected_operation_anchor_count":4,
		"selected_value_anchor_count":8,
		"maximum_anchor_width_bytes":0,
		"maximum_total_anchor_bytes":0,
		"value_anchor_duplicate_count":0,
		"value_operation_anchor_collision_count":0,
		"minimum_training_operation_accuracy":0,
		"minimum_training_left_accuracy":0,
		"minimum_training_right_accuracy":0,
		"maximum_training_operand_decode_failure_count":0,
		"minimum_heldout_operation_accuracy":0,
		"minimum_heldout_left_accuracy":0,
		"minimum_heldout_right_accuracy":0,
		"maximum_heldout_operand_decode_failure_count":0,
		"minimum_heldout_full_representation_accuracy":0,
		"reasoning_counter_count":96,
		"minimum_heldout_reasoning_result_accuracy":0,
		"minimum_output_generator_class_coverage":0,
		"maximum_output_generator_state_count":0,
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
	maxWidth,totalBytes:=0,0
	for _,anchor:=range opAnchors {
		if len(anchor)>maxWidth { maxWidth=len(anchor) }
		totalBytes+=len(anchor)
		if anchor=="" { metrics["invalid_row_count"]++ }
	}
	for _,anchor:=range valueAnchors {
		if len(anchor)>maxWidth { maxWidth=len(anchor) }
		totalBytes+=len(anchor)
		if anchor=="" { metrics["invalid_row_count"]++ }
	}
	metrics["maximum_anchor_width_bytes"]=float64(maxWidth)
	metrics["maximum_total_anchor_bytes"]=float64(totalBytes)
	duplicates,collisions:=wlmLmRawUTF8InstructionReasonGenerateR4CollisionMetrics(opAnchors,valueAnchors)
	metrics["value_anchor_duplicate_count"]=float64(duplicates)
	metrics["value_operation_anchor_collision_count"]=float64(collisions)

	opModel:=wlmLmRawUTF8InstructionReasonGenerateR3TrainOpModel(train,opAnchors,metrics)
	trainOp,trainLeft,trainRight,_,trainDecode:=wlmLmRawUTF8InstructionReasonGenerateR3Accuracy(train,&opModel,valueAnchors)
	heldOp,heldLeft,heldRight,heldFull,heldDecode:=wlmLmRawUTF8InstructionReasonGenerateR3Accuracy(held,&opModel,valueAnchors)
	metrics["minimum_training_operation_accuracy"]=trainOp
	metrics["minimum_training_left_accuracy"]=trainLeft
	metrics["minimum_training_right_accuracy"]=trainRight
	metrics["maximum_training_operand_decode_failure_count"]=float64(trainDecode)
	metrics["minimum_heldout_operation_accuracy"]=heldOp
	metrics["minimum_heldout_left_accuracy"]=heldLeft
	metrics["minimum_heldout_right_accuracy"]=heldRight
	metrics["maximum_heldout_operand_decode_failure_count"]=float64(heldDecode)
	metrics["minimum_heldout_full_representation_accuracy"]=heldFull

	var reasoner wlmLmRawUTF8InstructionReasonGenerateR1Reasoner
	var generator wlmLmRawUTF8InstructionReasonGenerateR1Generator
	lookup:=make(map[string]string,len(train))
	for _,ex:=range train {
		op:=opModel.predict(ex.text)
		left,right,_:=wlmLmRawUTF8InstructionReasonGenerateR3FindOperands(ex.text,valueAnchors)
		reasoner.observe(op,left,right,ex.target,metrics)
		generator.observe(ex.target,ex.output,metrics)
		lookup[ex.text]=ex.output
	}

	coverage:=0
	for result:=0;result<8;result++ {
		if generator.counts[result]!=nil&&len(generator.counts[result])>0 { coverage++ }
	}
	metrics["minimum_output_generator_class_coverage"]=float64(coverage)
	metrics["maximum_output_generator_state_count"]=float64(generator.stateCount())

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
	return wlmLmRawUTF8InstructionReasonGenerateR4Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-UTF8-INSTRUCTION-REASON-GENERATE-R4",
		Metrics:metrics,
		OperationAnchors:opOut,
		ValueAnchors:valueOut,
	}
}
