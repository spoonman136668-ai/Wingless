package unitary

import "fmt"

const UPLM0QUpdateOrderSchema = "wingless.up-lm0q-update-order.v1"

type UPLM0QByteMetric struct {
	Schedule     string  `json:"schedule"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM0QRoutingMetric struct {
	Schedule string          `json:"schedule"`
	Metric   UPLM0JMetric `json:"metric"`
}

type UPLM0QUpdateOrderResult struct {
	Schema                 string                 `json:"schema"`
	Experiment             string                 `json:"experiment"`
	SourceUPLM0PSeal       string                 `json:"source_up_lm0p_seal"`
	StateDimension         int                    `json:"state_dimension"`
	ExactRecallCap         int                    `json:"exact_recall_cap"`
	BaseEpochs             int                    `json:"base_epochs"`
	AdaptationEpochs       int                    `json:"adaptation_epochs"`
	LearningRate           float64                `json:"learning_rate"`
	ParaphrasePassesPerEpoch int                  `json:"paraphrase_passes_per_epoch"`
	BasePassesPerEpoch     int                    `json:"base_passes_per_epoch"`
	RouterRetrainingUsed   bool                   `json:"router_retraining_used"`
	AttentionUsed          bool                   `json:"attention_used"`
	FutureOracleUsed       bool                   `json:"future_oracle_used"`
	ByteMetrics            []UPLM0QByteMetric     `json:"byte_metrics"`
	RoutingMetrics         []UPLM0QRoutingMetric  `json:"routing_metrics"`
}

func uplm0qTrainSchedule(model *uplm0aModel,baseTrain,paraTrain []uplm0fExample,schedule string) error {
	if len(baseTrain)!=len(paraTrain) { return fmt.Errorf("corpus count mismatch: base=%d paraphrase=%d",len(baseTrain),len(paraTrain)) }
	for epoch:=0;epoch<4;epoch++ {
		switch schedule {
		case "para_then_base":
			for _,ex:=range paraTrain { model.trainSentence(ex.text,0.08) }
			for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
		case "base_then_para":
			for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
			for _,ex:=range paraTrain { model.trainSentence(ex.text,0.08) }
		case "example_interleaved":
			for i:=range baseTrain {
				model.trainSentence(paraTrain[i].text,0.08)
				model.trainSentence(baseTrain[i].text,0.08)
			}
		default:
			return fmt.Errorf("unknown schedule %q",schedule)
		}
	}
	return nil
}

func RunUPLM0Q()(UPLM0QUpdateOrderResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) } }
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0QUpdateOrderResult{
		Schema:UPLM0QUpdateOrderSchema,Experiment:"UP-LM0Q-update-order",
		SourceUPLM0PSeal:"7321157517e23061a4bc15d498d535d717d8103a",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		ParaphrasePassesPerEpoch:1,BasePassesPerEpoch:1,
		RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	for _,schedule:=range []string{"para_then_base","base_then_para","example_interleaved"} {
		model:=uplm0oCloneModel(baseModel)
		if err:=uplm0qTrainSchedule(model,baseTrain,paraTrain,schedule);err!=nil{return result,err}
		b:=uplm0oByteEval(model,baseHeld,4,"base_heldout")
		p:=uplm0oByteEval(model,paraHeld,4,"paraphrase_block_heldout")
		result.ByteMetrics=append(result.ByteMetrics,
			UPLM0QByteMetric{Schedule:schedule,Split:b.Split,Top1Accuracy:b.Top1Accuracy,Perplexity:b.Perplexity,Tokens:b.Tokens},
			UPLM0QByteMetric{Schedule:schedule,Split:p.Split,Top1Accuracy:p.Top1Accuracy,Perplexity:p.Perplexity,Tokens:p.Tokens},
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0QRoutingMetric{Schedule:schedule,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0QRoutingMetric{Schedule:schedule,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
