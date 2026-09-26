package unitary

const UPLM1AInitializationHistorySchema = "wingless.up-lm1a-initialization-history.v1"

type UPLM1AMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM1AInitializationHistoryResult struct {
	Schema                     string         `json:"schema"`
	Experiment                 string         `json:"experiment"`
	SourceUPLM0ZSeal           string         `json:"source_up_lm0z_seal"`
	StateDimension             int            `json:"state_dimension"`
	LearningRate               float64        `json:"learning_rate"`
	BasePretrainEpochs         int            `json:"base_pretrain_epochs"`
	ReadoutParameterCount      int            `json:"readout_parameter_count"`
	RecurrentParametersTrained bool           `json:"recurrent_parameters_trained"`
	RouterUsed                 bool           `json:"router_used"`
	ExactRecallUsed            bool           `json:"exact_recall_used"`
	AttentionUsed              bool           `json:"attention_used"`
	Metrics                    []UPLM1AMetric `json:"metrics"`
}

func uplm1aJointTrain(model *uplm0aModel,baseTrain,paraTrain []uplm0fExample,epochs int){
	for epoch:=0;epoch<epochs;epoch++{
		for i:=range baseTrain{
			model.trainSentence(baseTrain[i].text,0.08)
			model.trainSentence(paraTrain[i].text,0.08)
		}
	}
}

func uplm1aEvaluate(result *UPLM1AInitializationHistoryResult,arm string,model *uplm0aModel,baseTrain,baseHeld,paraTrain,paraHeld []uplm0fExample){
	for _,x:=range []struct{split string;data []uplm0fExample}{
		{"base_train",baseTrain},{"base_heldout",baseHeld},
		{"paraphrase_train",paraTrain},{"paraphrase_heldout",paraHeld},
	}{
		m:=uplm0oByteEval(model,x.data,0,x.split)
		result.Metrics=append(result.Metrics,UPLM1AMetric{Arm:arm,Split:x.split,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity,Tokens:m.Tokens})
	}
}

func RunUPLM1A()(UPLM1AInitializationHistoryResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	result:=UPLM1AInitializationHistoryResult{
		Schema:UPLM1AInitializationHistorySchema,
		Experiment:"UP-LM1A-initialization-history",
		SourceUPLM0ZSeal:"761ed2d9c275d6707602d5b658bf2e1d7e46d253",
		StateDimension:64,LearningRate:0.08,BasePretrainEpochs:20,
		ReadoutParameterCount:len(alphabet)*(64+1),
		RecurrentParametersTrained:false,RouterUsed:false,ExactRecallUsed:false,AttentionUsed:false,
	}

	for _,arm:=range []struct{name string;warm bool;jointEpochs int}{
		{"fresh_4",false,4},
		{"warm_4",true,4},
		{"fresh_20",false,20},
		{"warm_20",true,20},
	}{
		m:=newUPLM0AModel(alphabet)
		if arm.warm{
			for epoch:=0;epoch<20;epoch++{for _,ex:=range baseTrain{m.trainSentence(ex.text,0.08)}}
		}
		uplm1aJointTrain(m,baseTrain,paraTrain,arm.jointEpochs)
		uplm1aEvaluate(&result,arm.name,m,baseTrain,baseHeld,paraTrain,paraHeld)
	}
	return result,nil
}
