package unitary

const UPLM0UAnchorInterleavedSchema = "wingless.up-lm0u-anchor-interleaved.v1"

type UPLM0UByteMetric struct {
	Arm          string  `json:"arm"`
	Lambda       float64 `json:"lambda"`
	BaseReplay   bool    `json:"base_replay"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM0URoutingMetric struct {
	Arm    string       `json:"arm"`
	Lambda float64      `json:"lambda"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM0UAnchorInterleavedResult struct {
	Schema               string                  `json:"schema"`
	Experiment           string                  `json:"experiment"`
	SourceUPLM0TSeal     string                  `json:"source_up_lm0t_seal"`
	StateDimension       int                     `json:"state_dimension"`
	ExactRecallCap       int                     `json:"exact_recall_cap"`
	BaseEpochs           int                     `json:"base_epochs"`
	AdaptationEpochs     int                     `json:"adaptation_epochs"`
	LearningRate         float64                 `json:"learning_rate"`
	RouterRetrainingUsed bool                    `json:"router_retraining_used"`
	AttentionUsed        bool                    `json:"attention_used"`
	FutureOracleUsed     bool                    `json:"future_oracle_used"`
	ByteMetrics          []UPLM0UByteMetric      `json:"byte_metrics"`
	RoutingMetrics       []UPLM0URoutingMetric   `json:"routing_metrics"`
}

func uplm0uTrain(model,base *uplm0aModel,baseTrain,paraTrain []uplm0fExample,replay bool,lambda float64) {
	for epoch:=0;epoch<4;epoch++ {
		if !replay {
			for _,ex:=range paraTrain { uplm0tTrainAnchored(model,base,ex.text,0.08,lambda) }
			continue
		}
		for i:=range paraTrain {
			uplm0tTrainAnchored(model,base,paraTrain[i].text,0.08,lambda)
			uplm0tTrainAnchored(model,base,baseTrain[i].text,0.08,lambda)
		}
	}
}

func RunUPLM0U()(UPLM0UAnchorInterleavedResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) } }
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0UAnchorInterleavedResult{
		Schema:UPLM0UAnchorInterleavedSchema,Experiment:"UP-LM0U-anchor-interleaved",
		SourceUPLM0TSeal:"549017ab748a2fffffb49807088c1ad624d3c8ef",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	arms:=[]struct{name string; replay bool; lambda float64}{
		{"no_replay_anchor0",false,0},
		{"no_replay_anchor0001",false,0.0001},
		{"interleaved_anchor0",true,0},
		{"interleaved_anchor0001",true,0.0001},
	}
	for _,arm:=range arms {
		model:=uplm0oCloneModel(baseModel)
		uplm0uTrain(model,baseModel,baseTrain,paraTrain,arm.replay,arm.lambda)
		b:=uplm0oByteEval(model,baseHeld,4,"base_heldout")
		p:=uplm0oByteEval(model,paraHeld,4,"paraphrase_block_heldout")
		result.ByteMetrics=append(result.ByteMetrics,
			UPLM0UByteMetric{Arm:arm.name,Lambda:arm.lambda,BaseReplay:arm.replay,Split:b.Split,Top1Accuracy:b.Top1Accuracy,Perplexity:b.Perplexity,Tokens:b.Tokens},
			UPLM0UByteMetric{Arm:arm.name,Lambda:arm.lambda,BaseReplay:arm.replay,Split:p.Split,Top1Accuracy:p.Top1Accuracy,Perplexity:p.Perplexity,Tokens:p.Tokens},
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0URoutingMetric{Arm:arm.name,Lambda:arm.lambda,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0URoutingMetric{Arm:arm.name,Lambda:arm.lambda,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
