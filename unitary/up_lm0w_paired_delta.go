package unitary

const UPLM0WPairedDeltaSchema = "wingless.up-lm0w-paired-delta.v1"

type UPLM0WByteMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM0WRoutingMetric struct {
	Arm    string       `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM0WPairedDeltaResult struct {
	Schema               string                 `json:"schema"`
	Experiment           string                 `json:"experiment"`
	SourceUPLM0VSeal     string                 `json:"source_up_lm0v_seal"`
	StateDimension       int                    `json:"state_dimension"`
	ExactRecallCap       int                    `json:"exact_recall_cap"`
	BaseEpochs           int                    `json:"base_epochs"`
	AdaptationEpochs     int                    `json:"adaptation_epochs"`
	LearningRate         float64                `json:"learning_rate"`
	RouterRetrainingUsed bool                   `json:"router_retraining_used"`
	AttentionUsed        bool                   `json:"attention_used"`
	FutureOracleUsed     bool                   `json:"future_oracle_used"`
	AnchoringUsed        bool                   `json:"anchoring_used"`
	ByteMetrics          []UPLM0WByteMetric     `json:"byte_metrics"`
	RoutingMetrics       []UPLM0WRoutingMetric  `json:"routing_metrics"`
}

func uplm0wApplyPairedDelta(model *uplm0aModel,paraText,baseText string) {
	p:=uplm0oCloneModel(model)
	b:=uplm0oCloneModel(model)
	p.trainSentence(paraText,0.08)
	b.trainSentence(baseText,0.08)
	for c:=range model.w {
		oldBias:=model.b[c]
		model.b[c]=oldBias+(p.b[c]-oldBias)+(b.b[c]-oldBias)
		for i:=0;i<64;i++ {
			old:=model.w[c][i]
			model.w[c][i]=old+(p.w[c][i]-old)+(b.w[c][i]-old)
		}
	}
}

func uplm0wTrain(model *uplm0aModel,baseTrain,paraTrain []uplm0fExample,arm string) {
	for epoch:=0;epoch<4;epoch++ {
		for i:=range paraTrain {
			switch arm {
			case "paraphrase_then_base":
				model.trainSentence(paraTrain[i].text,0.08)
				model.trainSentence(baseTrain[i].text,0.08)
			case "base_then_paraphrase":
				model.trainSentence(baseTrain[i].text,0.08)
				model.trainSentence(paraTrain[i].text,0.08)
			case "paired_delta_sum":
				uplm0wApplyPairedDelta(model,paraTrain[i].text,baseTrain[i].text)
			}
		}
	}
}

func RunUPLM0W()(UPLM0WPairedDeltaResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) } }
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0WPairedDeltaResult{
		Schema:UPLM0WPairedDeltaSchema,
		Experiment:"UP-LM0W-paired-delta",
		SourceUPLM0VSeal:"7d8d743c05ef86806713f26e1f41146618ce46bc",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,AnchoringUsed:false,
	}
	for _,arm:=range []string{"paraphrase_then_base","base_then_paraphrase","paired_delta_sum"} {
		model:=uplm0oCloneModel(baseModel)
		uplm0wTrain(model,baseTrain,paraTrain,arm)
		b:=uplm0oByteEval(model,baseHeld,4,"base_heldout")
		p:=uplm0oByteEval(model,paraHeld,4,"paraphrase_block_heldout")
		result.ByteMetrics=append(result.ByteMetrics,
			UPLM0WByteMetric{Arm:arm,Split:b.Split,Top1Accuracy:b.Top1Accuracy,Perplexity:b.Perplexity,Tokens:b.Tokens},
			UPLM0WByteMetric{Arm:arm,Split:p.Split,Top1Accuracy:p.Top1Accuracy,Perplexity:p.Perplexity,Tokens:p.Tokens},
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0WRoutingMetric{Arm:arm,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0WRoutingMetric{Arm:arm,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
