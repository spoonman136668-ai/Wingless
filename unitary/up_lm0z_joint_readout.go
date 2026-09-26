package unitary

import "fmt"

const UPLM0ZJointReadoutSchema = "wingless.up-lm0z-joint-readout.v1"

type UPLM0ZMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM0ZJointReadoutResult struct {
	Schema                     string          `json:"schema"`
	Experiment                 string          `json:"experiment"`
	SourceUPLM0YSeal           string          `json:"source_up_lm0y_seal"`
	StateDimension             int             `json:"state_dimension"`
	LearningRate               float64         `json:"learning_rate"`
	BaseOnlyEpochs             int             `json:"base_only_epochs"`
	ParaphraseOnlyEpochs       int             `json:"paraphrase_only_epochs"`
	JointInterleavedEpochs     int             `json:"joint_interleaved_epochs"`
	JointFullBatchEpochs       int             `json:"joint_fullbatch_epochs"`
	ReadoutParameterCount      int             `json:"readout_parameter_count"`
	RecurrentParametersTrained bool            `json:"recurrent_parameters_trained"`
	RouterUsed                 bool            `json:"router_used"`
	ExactRecallUsed            bool            `json:"exact_recall_used"`
	AttentionUsed              bool            `json:"attention_used"`
	Metrics                    []UPLM0ZMetric `json:"metrics"`
}

func uplm0zMetric(arm,split string,m UPLM0OByteMetric) UPLM0ZMetric {
	return UPLM0ZMetric{Arm:arm,Split:split,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity,Tokens:m.Tokens}
}

func uplm0zTokenCount(examples []uplm0fExample) int {
	n:=0
	for _,ex:=range examples {
		if len(ex.text)>0 { n+=len(ex.text)-1 }
	}
	return n
}

func uplm0zAddGradient(dst *uplm0xGradient,src uplm0xGradient) {
	for c:=range dst.w {
		for i:=0;i<64;i++ { dst.w[c][i]+=src.w[c][i] }
		dst.b[c]+=src.b[c]
	}
}

func uplm0zFullBatchEpoch(model *uplm0aModel,baseTrain,paraTrain []uplm0fExample,lr float64) {
	g:=uplm0xNewGradient(len(model.alphabet))
	for _,ex:=range baseTrain { uplm0zAddGradient(&g,uplm0xSentenceGradient(model,ex.text)) }
	for _,ex:=range paraTrain { uplm0zAddGradient(&g,uplm0xSentenceGradient(model,ex.text)) }
	tokens:=uplm0zTokenCount(baseTrain)+uplm0zTokenCount(paraTrain)
	scale:=lr/float64(tokens)
	for c:=range model.w {
		for i:=0;i<64;i++ { model.w[c][i]-=scale*g.w[c][i] }
		model.b[c]-=scale*g.b[c]
	}
}

func uplm0zEvaluateArm(result *UPLM0ZJointReadoutResult,arm string,model *uplm0aModel,baseTrain,baseHeld,paraTrain,paraHeld []uplm0fExample) {
	for _,x:=range []struct{split string;data []uplm0fExample}{
		{"base_train",baseTrain},{"base_heldout",baseHeld},
		{"paraphrase_train",paraTrain},{"paraphrase_heldout",paraHeld},
	}{
		m:=uplm0oByteEval(model,x.data,0,x.split)
		result.Metrics=append(result.Metrics,uplm0zMetric(arm,x.split,m))
	}
}

func RunUPLM0Z()(UPLM0ZJointReadoutResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	result:=UPLM0ZJointReadoutResult{
		Schema:UPLM0ZJointReadoutSchema,
		Experiment:"UP-LM0Z-joint-readout",
		SourceUPLM0YSeal:"fc8d640b19fdd89fcdac2b79e423df2fbdcce6a6",
		StateDimension:64,LearningRate:0.08,
		BaseOnlyEpochs:20,ParaphraseOnlyEpochs:20,JointInterleavedEpochs:20,JointFullBatchEpochs:100,
		ReadoutParameterCount:len(alphabet)*(64+1),
		RecurrentParametersTrained:false,RouterUsed:false,ExactRecallUsed:false,AttentionUsed:false,
	}

	baseOnly:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { baseOnly.trainSentence(ex.text,0.08) } }
	uplm0zEvaluateArm(&result,"base_only",baseOnly,baseTrain,baseHeld,paraTrain,paraHeld)

	paraOnly:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range paraTrain { paraOnly.trainSentence(ex.text,0.08) } }
	uplm0zEvaluateArm(&result,"paraphrase_only",paraOnly,baseTrain,baseHeld,paraTrain,paraHeld)

	if len(baseTrain)!=len(paraTrain) {
		return result,fmt.Errorf("matched corpus count mismatch: base=%d paraphrase=%d",len(baseTrain),len(paraTrain))
	}
	jointInterleaved:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for i:=range baseTrain {
			jointInterleaved.trainSentence(baseTrain[i].text,0.08)
			jointInterleaved.trainSentence(paraTrain[i].text,0.08)
		}
	}
	uplm0zEvaluateArm(&result,"joint_interleaved",jointInterleaved,baseTrain,baseHeld,paraTrain,paraHeld)

	jointFullBatch:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<100;epoch++ { uplm0zFullBatchEpoch(jointFullBatch,baseTrain,paraTrain,0.08) }
	uplm0zEvaluateArm(&result,"joint_fullbatch",jointFullBatch,baseTrain,baseHeld,paraTrain,paraHeld)

	return result,nil
}
