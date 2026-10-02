package unitary

import (
	"math"
	"strings"
)

type wlmSiRawUTF8SurfaceInvarianceFalsificationR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmSiRawUTF8SurfaceInvarianceFalsificationR1Punctuation(text string) string {
	if strings.HasSuffix(text, "?\n") {
		return strings.TrimSuffix(text, "?\n") + ".\n"
	}
	if strings.HasSuffix(text, ".\n") {
		return strings.TrimSuffix(text, ".\n") + "!\n"
	}
	return text
}

func wlmSiRawUTF8SurfaceInvarianceFalsificationR1DoubleSpace(text string) string {
	return strings.ReplaceAll(text, " ", "  ")
}

func wlmSiRawUTF8SurfaceInvarianceFalsificationR1TabSpace(text string) string {
	return strings.ReplaceAll(text, " ", "\t")
}

func wlmSiRawUTF8SurfaceInvarianceFalsificationR1Upper(text string) string {
	data:=[]byte(text)
	for i,b:=range data {
		if b>='a'&&b<='z' {
			data[i]=b-'a'+'A'
		}
	}
	return string(data)
}

func wlmSiRawUTF8SurfaceInvarianceFalsificationR1Evaluate(
	examples []wlmLmRawUTF8InstructionReasonGenerateR1Example,
	transform func(string) string,
	opModel *wlmLmRawUTF8InstructionReasonGenerateR3OpModel,
	valueAnchors [8]string,
	reasoner *wlmLmRawUTF8InstructionReasonGenerateR1Reasoner,
	generator *wlmLmRawUTF8InstructionReasonGenerateR5Generator,
)(float64,float64) {
	reasonHits:=0
	outputHits:=0
	for _,ex:=range examples {
		text:=transform(ex.text)
		op:=opModel.predict(text)
		left,right,_:=wlmLmRawUTF8InstructionReasonGenerateR3FindOperands(text,valueAnchors)
		result:=reasoner.predict(op,left,right)
		if result==ex.target { reasonHits++ }
		if generator.generate(result)==ex.output { outputHits++ }
	}
	if len(examples)==0 { return 0,0 }
	den:=float64(len(examples))
	return float64(reasonHits)/den,float64(outputHits)/den
}

// RunWlmSiRawUTF8SurfaceInvarianceFalsificationR1 tests R5 against frozen raw-byte surface transforms.
func RunWlmSiRawUTF8SurfaceInvarianceFalsificationR1() interface{} {
	metrics:=map[string]float64{
		"evaluation_variant_count":5,
		"evaluation_examples_per_variant":0,
		"minimum_original_reasoning_accuracy":0,
		"minimum_original_exact_raw_output_accuracy":0,
		"minimum_punctuation_reasoning_accuracy":0,
		"minimum_punctuation_exact_raw_output_accuracy":0,
		"minimum_double_space_reasoning_accuracy":0,
		"minimum_double_space_exact_raw_output_accuracy":0,
		"maximum_tab_reasoning_accuracy":0,
		"maximum_tab_exact_raw_output_accuracy":0,
		"maximum_uppercase_reasoning_accuracy":0,
		"maximum_uppercase_exact_raw_output_accuracy":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
	}

	train,held:=wlmLmRawUTF8InstructionReasonGenerateR1Corpus()
	metrics["evaluation_examples_per_variant"]=float64(len(held))
	opAnchors,_:=wlmLmRawUTF8InstructionReasonGenerateR3SelectAnchors(train)
	valueAnchors:=wlmLmRawUTF8InstructionReasonGenerateR4SelectValueAnchors(train)
	opModel:=wlmLmRawUTF8InstructionReasonGenerateR3TrainOpModel(train,opAnchors,metrics)

	var reasoner wlmLmRawUTF8InstructionReasonGenerateR1Reasoner
	var generator wlmLmRawUTF8InstructionReasonGenerateR5Generator
	for _,ex:=range train {
		op:=opModel.predict(ex.text)
		left,right,_:=wlmLmRawUTF8InstructionReasonGenerateR3FindOperands(ex.text,valueAnchors)
		reasoner.observe(op,left,right,ex.target,metrics)
		generator.observe(ex.target,ex.output,metrics)
	}

	identity:=func(s string) string { return s }
	origR,origO:=wlmSiRawUTF8SurfaceInvarianceFalsificationR1Evaluate(held,identity,&opModel,valueAnchors,&reasoner,&generator)
	puncR,puncO:=wlmSiRawUTF8SurfaceInvarianceFalsificationR1Evaluate(held,wlmSiRawUTF8SurfaceInvarianceFalsificationR1Punctuation,&opModel,valueAnchors,&reasoner,&generator)
	doubleR,doubleO:=wlmSiRawUTF8SurfaceInvarianceFalsificationR1Evaluate(held,wlmSiRawUTF8SurfaceInvarianceFalsificationR1DoubleSpace,&opModel,valueAnchors,&reasoner,&generator)
	tabR,tabO:=wlmSiRawUTF8SurfaceInvarianceFalsificationR1Evaluate(held,wlmSiRawUTF8SurfaceInvarianceFalsificationR1TabSpace,&opModel,valueAnchors,&reasoner,&generator)
	upperR,upperO:=wlmSiRawUTF8SurfaceInvarianceFalsificationR1Evaluate(held,wlmSiRawUTF8SurfaceInvarianceFalsificationR1Upper,&opModel,valueAnchors,&reasoner,&generator)

	metrics["minimum_original_reasoning_accuracy"]=origR
	metrics["minimum_original_exact_raw_output_accuracy"]=origO
	metrics["minimum_punctuation_reasoning_accuracy"]=puncR
	metrics["minimum_punctuation_exact_raw_output_accuracy"]=puncO
	metrics["minimum_double_space_reasoning_accuracy"]=doubleR
	metrics["minimum_double_space_exact_raw_output_accuracy"]=doubleO
	metrics["maximum_tab_reasoning_accuracy"]=tabR
	metrics["maximum_tab_exact_raw_output_accuracy"]=tabO
	metrics["maximum_uppercase_reasoning_accuracy"]=upperR
	metrics["maximum_uppercase_exact_raw_output_accuracy"]=upperO

	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) {
			metrics["invalid_row_count"]++
		}
	}

	return wlmSiRawUTF8SurfaceInvarianceFalsificationR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-SI-RAW-UTF8-SURFACE-INVARIANCE-FALSIFICATION-R1",
		Metrics:metrics,
	}
}
