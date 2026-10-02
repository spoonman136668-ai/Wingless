package unitary

import (
	"math"
	"sort"
	"strings"
)

type wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats struct {
	totalPresence int
	totalOccurrences int
	opPresence [4]int
	opOccurrences [4]int
	valuePresence [8]int
	valueOccurrences [8]int
}

type wlmLmRawUTF8InstructionReasonGenerateR3OpModel struct {
	anchors [4]string
	hist [4][4][9]uint32
	classN [4]uint32
}

type wlmLmRawUTF8InstructionReasonGenerateR3Match struct {
	start int
	class int
	width int
}

type wlmLmRawUTF8InstructionReasonGenerateR3Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	OperationAnchors []string `json:"operation_anchors"`
	ValueAnchors []string `json:"value_anchors"`
}

func wlmLmRawUTF8InstructionReasonGenerateR3BuildStats(
	train []wlmLmRawUTF8InstructionReasonGenerateR1Example,
)(map[string]*wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats,[4]int,[8]int) {
	stats:=make(map[string]*wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats)
	var opN [4]int
	var valueN [8]int
	for _,ex:=range train {
		opN[ex.op]++
		seenValue:=[8]bool{}
		seenValue[ex.left]=true
		seenValue[ex.right]=true
		for value:=0;value<8;value++ {
			if seenValue[value] { valueN[value]++ }
		}

		local:=make(map[string]int)
		data:=[]byte(ex.text)
		for width:=2;width<=7;width++ {
			for start:=0;start+width<=len(data);start++ {
				local[string(data[start:start+width])]++
			}
		}
		for candidate,count:=range local {
			row:=stats[candidate]
			if row==nil {
				row=&wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats{}
				stats[candidate]=row
			}
			row.totalPresence++
			row.totalOccurrences+=count
			row.opPresence[ex.op]++
			row.opOccurrences[ex.op]+=count
			for value:=0;value<8;value++ {
				if seenValue[value] {
					row.valuePresence[value]++
					row.valueOccurrences[value]+=count
				}
			}
		}
	}
	return stats,opN,valueN
}

func wlmLmRawUTF8InstructionReasonGenerateR3Better(
	candidate string,
	row *wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats,
	class int,
	totalExamples int,
	classN int,
	operation bool,
	best string,
	bestRow *wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats,
) bool {
	if bestRow==nil {
		return true
	}
	var posPresence,posOcc,bestPosPresence,bestPosOcc int
	if operation {
		posPresence=row.opPresence[class]
		posOcc=row.opOccurrences[class]
		bestPosPresence=bestRow.opPresence[class]
		bestPosOcc=bestRow.opOccurrences[class]
	} else {
		posPresence=row.valuePresence[class]
		posOcc=row.valueOccurrences[class]
		bestPosPresence=bestRow.valuePresence[class]
		bestPosOcc=bestRow.valueOccurrences[class]
	}

	coverage:=float64(posPresence)/float64(classN)
	bestCoverage:=float64(bestPosPresence)/float64(classN)
	if coverage!=bestCoverage {
		return coverage>bestCoverage
	}

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

	precision:=float64(posPresence)/float64(row.totalPresence)
	bestPrecision:=float64(bestPosPresence)/float64(bestRow.totalPresence)
	if precision!=bestPrecision {
		return precision>bestPrecision
	}
	if len(candidate)!=len(best) {
		return len(candidate)>len(best)
	}
	return candidate<best
}

func wlmLmRawUTF8InstructionReasonGenerateR3SelectAnchors(
	train []wlmLmRawUTF8InstructionReasonGenerateR1Example,
)([4]string,[8]string) {
	stats,opN,valueN:=wlmLmRawUTF8InstructionReasonGenerateR3BuildStats(train)
	var opAnchors [4]string
	var valueAnchors [8]string
	for class:=0;class<4;class++ {
		var bestRow *wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats
		for candidate,row:=range stats {
			if row.opPresence[class]==0 { continue }
			if wlmLmRawUTF8InstructionReasonGenerateR3Better(candidate,row,class,len(train),opN[class],true,opAnchors[class],bestRow) {
				opAnchors[class]=candidate
				bestRow=row
			}
		}
	}
	for class:=0;class<8;class++ {
		var bestRow *wlmLmRawUTF8InstructionReasonGenerateR3CandidateStats
		for candidate,row:=range stats {
			if row.valuePresence[class]==0 { continue }
			if wlmLmRawUTF8InstructionReasonGenerateR3Better(candidate,row,class,len(train),valueN[class],false,valueAnchors[class],bestRow) {
				valueAnchors[class]=candidate
				bestRow=row
			}
		}
	}
	return opAnchors,valueAnchors
}

func wlmLmRawUTF8InstructionReasonGenerateR3Count(text,anchor string) int {
	if anchor=="" { return 0 }
	count:=0
	start:=0
	for {
		index:=strings.Index(text[start:],anchor)
		if index<0 { break }
		count++
		start+=index+len(anchor)
		if start>=len(text) { break }
	}
	return count
}

func wlmLmRawUTF8InstructionReasonGenerateR3TrainOpModel(
	train []wlmLmRawUTF8InstructionReasonGenerateR1Example,
	anchors [4]string,
	metrics map[string]float64,
) wlmLmRawUTF8InstructionReasonGenerateR3OpModel {
	model:=wlmLmRawUTF8InstructionReasonGenerateR3OpModel{anchors:anchors}
	for _,ex:=range train {
		if model.classN[ex.op]==^uint32(0) {
			metrics["counter_overflow_count"]++
			continue
		}
		model.classN[ex.op]++
		for anchorID,anchor:=range anchors {
			count:=wlmLmRawUTF8InstructionReasonGenerateR3Count(ex.text,anchor)
			if count>8 { count=8 }
			if model.hist[ex.op][anchorID][count]==^uint32(0) {
				metrics["counter_overflow_count"]++
				continue
			}
			model.hist[ex.op][anchorID][count]++
		}
	}
	return model
}

func (m *wlmLmRawUTF8InstructionReasonGenerateR3OpModel) predict(text string) int {
	best:=0
	bestScore:=math.Inf(-1)
	var total uint32
	for _,n:=range m.classN { total+=n }
	for class:=0;class<4;class++ {
		score:=math.Log((float64(m.classN[class])+1)/(float64(total)+4))
		for anchorID,anchor:=range m.anchors {
			count:=wlmLmRawUTF8InstructionReasonGenerateR3Count(text,anchor)
			if count>8 { count=8 }
			score+=math.Log((float64(m.hist[class][anchorID][count])+1)/(float64(m.classN[class])+9))
		}
		if score>bestScore {
			best=class
			bestScore=score
		}
	}
	return best
}

func wlmLmRawUTF8InstructionReasonGenerateR3FindOperands(text string,anchors [8]string)(left,right int,ok bool) {
	matches:=make([]wlmLmRawUTF8InstructionReasonGenerateR3Match,0,8)
	data:=[]byte(text)
	for class,anchor:=range anchors {
		if anchor=="" { continue }
		a:=[]byte(anchor)
		for start:=0;start+len(a)<=len(data);start++ {
			equal:=true
			for i:=range a {
				if data[start+i]!=a[i] { equal=false;break }
			}
			if equal {
				matches=append(matches,wlmLmRawUTF8InstructionReasonGenerateR3Match{start:start,class:class,width:len(a)})
			}
		}
	}
	sort.Slice(matches,func(i,j int)bool{
		if matches[i].start!=matches[j].start { return matches[i].start<matches[j].start }
		if matches[i].width!=matches[j].width { return matches[i].width>matches[j].width }
		return matches[i].class<matches[j].class
	})
	unique:=make([]wlmLmRawUTF8InstructionReasonGenerateR3Match,0,2)
	lastStart:=-1
	for _,match:=range matches {
		if match.start==lastStart { continue }
		unique=append(unique,match)
		lastStart=match.start
		if len(unique)==2 { break }
	}
	if len(unique)<2 {
		if len(unique)==1 {
			return unique[0].class,0,false
		}
		return 0,0,false
	}
	return unique[0].class,unique[1].class,true
}

func wlmLmRawUTF8InstructionReasonGenerateR3Accuracy(
	examples []wlmLmRawUTF8InstructionReasonGenerateR1Example,
	opModel *wlmLmRawUTF8InstructionReasonGenerateR3OpModel,
	valueAnchors [8]string,
)(opAcc,leftAcc,rightAcc,fullAcc float64,decodeFailures int) {
	opHits,leftHits,rightHits,fullHits:=0,0,0,0
	for _,ex:=range examples {
		op:=opModel.predict(ex.text)
		left,right,ok:=wlmLmRawUTF8InstructionReasonGenerateR3FindOperands(ex.text,valueAnchors)
		if !ok { decodeFailures++ }
		if op==ex.op { opHits++ }
		if left==ex.left { leftHits++ }
		if right==ex.right { rightHits++ }
		if op==ex.op&&left==ex.left&&right==ex.right { fullHits++ }
	}
	den:=float64(len(examples))
	if den==0 { return 0,0,0,0,decodeFailures }
	return float64(opHits)/den,float64(leftHits)/den,float64(rightHits)/den,float64(fullHits)/den,decodeFailures
}

// RunWlmLmRawUTF8InstructionReasonGenerateR3 evaluates learned raw-byte anchor representations.
func RunWlmLmRawUTF8InstructionReasonGenerateR3() interface{} {
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

	opAnchors,valueAnchors:=wlmLmRawUTF8InstructionReasonGenerateR3SelectAnchors(train)
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
	return wlmLmRawUTF8InstructionReasonGenerateR3Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-UTF8-INSTRUCTION-REASON-GENERATE-R3",
		Metrics:metrics,
		OperationAnchors:opOut,
		ValueAnchors:valueOut,
	}
}
